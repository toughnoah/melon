package utils

import (
	"github.com/spf13/viper"
	"testing"
)

func Test_getNamingExpr(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test correct expr",
			args:    args{"../testdata"},
			want:    `^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)`,
			wantErr: false,
		},
		{
			name:    "test correct expr",
			args:    args{"../abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNamingExpr(tt.args.path)
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
				path: "../testdata",
			},
			wantErr: false,
		},
		{
			name: "test fail validating",
			args: args{
				name: "noah-test",
				path: "../testdata",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateNaming(tt.args.name, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("ValidateNaming() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
