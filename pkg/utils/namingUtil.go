package utils

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"regexp"
)

const (
	defaultMelonConfig = "/etc/melon/config"

	namingError = "naming expr error. want string type, but get %s"

	matchExprError = "not match the expr %s"
)

func getNamingExpr(path string) (string, error) {
	if len(path) == 0 {
		path = defaultMelonConfig
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		klog.Fatalf("can not read config file for webhook server! %v", err)
		return "", err
	}

	expr, ok := viper.Get("naming_expr").(string)
	if !ok || len(expr) == 0 {
		return expr, errors.New(fmt.Sprintf(namingError, expr))
	}
	return expr, nil
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

func ValidateNaming(path string) error {
	expr, err := getNamingExpr(path)
	if err != nil {
		return err
	}
	if err = validateNaming(path, expr); err != nil {
		return err
	}
	return nil
}
