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
			if !tt.want {
				viper.Set("deployment.limits", false)
			}
			if got := IsToValidateLimits("../../tests/testdata", "deployment.limits"); got != tt.want {
				t.Errorf("IsToValidateLimits() = %v, want %v", got, tt.want)
			}
			viper.Reset()
		})
	}
}

func TestIsToValidateLimitsWithSetFalse(t *testing.T) {
	viper.Set("deployment.limits", "string")
	if got := IsToValidateLimits("../../tests/testdata", "deployment.limits"); got != false {
		t.Errorf("IsToValidateLimits() = %v, want %v", got, false)
	}
	viper.Reset()
}
