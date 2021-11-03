package utils

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"regexp"
	"strings"
)

const DefaultMatchAllExpr = ".*?"

type Config struct {
	Deployment *Deployment `yaml:"deployment"`
}

type Deployment struct {
	Image  *string `yaml:"image"`
	Naming *string `yaml:"naming"`
	Limit  *string `yaml:"limit"`
}

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
		klog.Errorf("can not read config file for webhook server! %v", err)
		return "", err
	}
	val := viper.Get(kind)
	if val != nil {
		return val, nil
	}
	if !strings.Contains(kind, "naming") {
		klog.Errorf(emptyValueError, kind)
		return nil, errors.New(fmt.Sprintf(emptyValueError, kind))
	}
	val = viper.Get("global.naming")
	if val != nil {
		klog.Info("get empty rule for validating naming: %s, using global rule %s", kind, val)
		return val, nil
	}
	klog.Errorf("get empty rule for validating: %s and global.naming", kind)
	return nil, errors.New(fmt.Sprintf("get empty rule for validating: %s and global.naming", kind))
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
	return errors.New(fmt.Sprintf(matchExprError, name, expr))
}

func ValidateNaming(name, path, kind string) error {
	expr, err := GetValFromConfig(path, kind)
	if err != nil {
		return err
	}
	exprStr, ok := expr.(string)
	if !ok {
		return errors.New(fmt.Sprintf(badValueTypeError, kind, "string"))
	}
	if err = validateNaming(name, exprStr); err != nil {
		return err
	}
	return nil
}
