package utils

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"regexp"
)

const (
	Namespace = iota

	Deployment

	Service

	Pod

	Configmap

	Daemonset

	Secret

	Default
)

var resourceTypeKeyMap = map[int]string{
	Deployment: "deploy_expr",
	Namespace:  "ns_expr",
	Service:    "svc_expr",
	Pod:        "po_expr",
	Configmap:  "cm_expr",
	Daemonset:  "ds_expr",
	Secret:     "sc_expr",
	Default:    "default_expr",
}

func getNamingExpr(path string, kind int) (string, error) {
	if len(path) == 0 {
		path = defaultMelonConfig
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		klog.Errorf("can not read config file for webhook server! %v", err)
		return "", err
	}
	exprKey, ok := resourceTypeKeyMap[kind]

	if !ok {
		return "", errors.New(noSuchKindError)
	}
	expr, ok := viper.Get(exprKey).(string)
	if ok && len(expr) != 0 {
		return expr, nil
	}
	defaultKey := resourceTypeKeyMap[Default]
	expr, ok = viper.Get(defaultKey).(string)
	if ok && len(expr) != 0 {
		return expr, nil
	}
	return expr, errors.New(fmt.Sprintf(namingError, expr))
}

func validateNaming(name, expr string) error {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

	match := reg.FindAllStringSubmatch(name, -1)
	if match != nil {
		return nil
	}
	return errors.New(fmt.Sprintf(matchExprError, expr))
}

func ValidateNaming(name, path string, kind int) error {
	expr, err := getNamingExpr(path, kind)
	if err != nil {
		return err
	}
	if err = validateNaming(name, expr); err != nil {
		return err
	}
	return nil
}
