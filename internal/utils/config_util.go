package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	k8syaml "sigs.k8s.io/yaml"
)

var cLog = log.Log.WithName("configreader")

func init() {
	viper.SetDefault("deployment.image", ".*?")
}

func GetValFromConfig(path, kind string) (interface{}, error) {
	if len(path) == 0 {
		path = defaultMelonConfig
	}
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		cLog.Error(err, "can not read config file for webhook server!", "path", path)
		return "", err
	}
	val := viper.Get(kind)
	if val != nil {
		return val, nil
	}
	if (strings.Contains(kind, "resources.resourceQuotaSpec") || strings.Contains(kind, "resources.limitRangeSpec")) && val == nil {
		cLog.Info("spec not found, skip it", "kind", kind)
		return val, nil
	}
	if !strings.Contains(kind, "naming") {
		err := errors.New(fmt.Sprintf(emptyValueError, kind))
		cLog.Error(err, "config error", "kind", kind)
		return nil, err
	}
	val = viper.Get("global.naming")
	if val != nil {
		cLog.Info("get empty check rule for validating naming, using global rule", "kind", kind, "global naming rule", val)
		return val, nil
	}
	err := errors.New(fmt.Sprintf("get empty rule for validating: %s and global.naming", kind))
	cLog.Error(err, "get empty check rule for validating and global.naming", "kine", kind)
	return nil, err
}

func GetResourceQuotaSepcFromConf(path string) (*v1.ResourceQuota, error) {
	val, err := GetValFromConfig(path, "deployment.resources.resourceQuotaSpec")
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	rq := &v1.ResourceQuotaSpec{}
	// here should use yaml.v2 to Marshal interface
	byteVal, err := yaml.Marshal(val)
	if err != nil {
		return nil, err
	}
	//here should use k8s yaml for Marshaling
	if err := k8syaml.Unmarshal([]byte(byteVal), rq); err != nil {
		return nil, err
	}
	return &v1.ResourceQuota{
		Spec: *rq,
	}, nil
}

func GetLimitRangeSepcFromConf(path string) (*v1.LimitRange, error) {
	val, err := GetValFromConfig(path, "deployment.resources.limitRangeSpec")
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}

	lr := &v1.LimitRangeSpec{}
	// here should use yaml.v2 to Marshal interface
	byteVal, err := yaml.Marshal(val)
	if err != nil {
		return nil, err
	}
	//here should use k8s yaml for Marshaling
	if err := k8syaml.Unmarshal([]byte(byteVal), lr); err != nil {
		return nil, err
	}
	for i := range lr.Limits {
		SetDefaults_LimitRangeItem(&lr.Limits[i])
	}
	return &v1.LimitRange{
		Spec: *lr,
	}, nil
}

func In(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}

// Refer from k8s /pkg/apis/core/v1/defautls.go
func SetDefaults_LimitRangeItem(obj *v1.LimitRangeItem) {
	// for container limits, we apply default values
	if obj.Type == v1.LimitTypeContainer {

		if obj.Default == nil {
			obj.Default = make(v1.ResourceList)
		}
		if obj.DefaultRequest == nil {
			obj.DefaultRequest = make(v1.ResourceList)
		}

		// If a default limit is unspecified, but the max is specified, default the limit to the max
		for key, value := range obj.Max {
			if _, exists := obj.Default[key]; !exists {
				obj.Default[key] = value.DeepCopy()
			}
		}
		// If a default limit is specified, but the default request is not, default request to limit
		for key, value := range obj.Default {
			if _, exists := obj.DefaultRequest[key]; !exists {
				obj.DefaultRequest[key] = value.DeepCopy()
			}
		}
		// If a default request is not specified, but the min is provided, default request to the min
		for key, value := range obj.Min {
			if _, exists := obj.DefaultRequest[key]; !exists {
				obj.DefaultRequest[key] = value.DeepCopy()
			}
		}
	}
}
