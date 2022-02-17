package utils

import (
	"k8s.io/klog/v2"
)

func IsToValidateLimits(path string, kind string) bool {
	isValidateLimits, _ := GetValFromConfig(path, kind)

	if isValidateLimits == nil {
		return false
	}
	is, ok := isValidateLimits.(bool)
	if !ok {
		klog.Errorf(badValueTypeError, kind, "bool")
		return false
	}
	return is
}
