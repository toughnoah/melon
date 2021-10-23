package utils

import (
	"github.com/spf13/viper"
	"testing"
)

func TestIsToValidateLimits(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "test to validate",
			want: true,
		},
		{
			name: "test not to validate",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want {
				viper.Set("is_validate_deploy_limits", true)
			}
			if got := IsToValidateLimits(); got != tt.want {
				t.Errorf("IsToValidateLimits() = %v, want %v", got, tt.want)
			}
			viper.Reset()
		})
	}
}

func TestIsToValidateLimitsWithSetFalse(t *testing.T) {
	viper.Set("is_validate_deploy_limits", "string")
	if got := IsToValidateLimits(); got != false {
		t.Errorf("IsToValidateLimits() = %v, want %v", got, false)
	}
	viper.Reset()
}
