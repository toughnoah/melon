package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

//
func Test_GetValFromConfig(t *testing.T) {
	type args struct {
		path string
		kind string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test correct expr",
			args:    args{"../../tests/testdata", "deployment.naming"},
			want:    `^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)`,
			wantErr: false,
		},
		{
			name:    "test global naming",
			args:    args{"../../tests/testdata2", "deployment.naming"},
			want:    `^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)`,
			wantErr: false,
		},
		{
			name:    "test correct expr with wrong path",
			args:    args{"../abc", "deployment.naming"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValFromConfig(tt.args.path, tt.args.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNamingExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getNamingExpr() = %v, want %v", got, tt.want)
			}
			// reset read in config
			viper.Reset()
		})
	}
}

func Test_validateNaming(t *testing.T) {
	type args struct {
		name string
		expr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test right name",
			args: args{
				name: "noah-dev-melon-test",
				expr: `^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)`,
			},
			wantErr: false,
		},
		{
			name: "test wrong name",
			args: args{
				name: "noah-test",
				expr: `^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateNaming(tt.args.name, tt.args.expr); (err != nil) != tt.wantErr {
				t.Errorf("validateNaming() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNaming(t *testing.T) {
	type args struct {
		name string
		path string
		kind string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test pass validating",
			args: args{
				name: "noah-dev-melon-test",
				path: "../../tests/testdata",
				kind: "deployment.naming",
			},
			wantErr: false,
		},
		{
			name: "test fail validating",
			args: args{
				name: "noah-test",
				path: "../../tests/testdata",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateNaming(tt.args.name, tt.args.path, tt.args.kind); (err != nil) != tt.wantErr {
				t.Errorf("ValidateNaming() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	config, err := GetValFromConfig("../../tests/testdata", "deployment.naming")
	if err != nil {
		return
	}
	fmt.Println(config)
}
