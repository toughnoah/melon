package utils

import (
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

func IsToValidateLimits() bool {
	isValidateLimits := viper.Get("is_validate_deploy_limits")

	if isValidateLimits == nil {
		return false
	}
	is, ok := isValidateLimits.(bool)
	if !ok {
		klog.Errorf(badValueTypeError, "bool")
		return false
	}
	return is
}
