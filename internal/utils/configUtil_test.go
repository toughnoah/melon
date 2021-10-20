package utils

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

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

func TestGetResourceQuotaSepcFromConf(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.ResourceQuota
		wantErr bool
	}{
		{
			name: "test get ResourceQuota from config and marshal it",
			args: args{
				path: "../../tests/testdata",
			},
			want: &v1.ResourceQuota{
				Spec: v1.ResourceQuotaSpec{
					Hard: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("10"),
						v1.ResourceMemory: resource.MustParse("20Gi"),
						v1.ResourcePods:   resource.MustParse("10"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResourceQuotaSepcFromConf(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResourceQuotaSepcFromConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResourceQuotaSepcFromConf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLimitRangeSepcFromConf(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.LimitRange
		wantErr bool
	}{
		{
			name: "test get LimitRange from config and Marshal it",
			args: args{
				path: "../../tests/testdata",
			},
			want: &v1.LimitRange{
				Spec: v1.LimitRangeSpec{
					Limits: []v1.LimitRangeItem{
						{
							Type: v1.LimitTypeContainer,
							Default: v1.ResourceList{
								v1.ResourceCPU: resource.MustParse("800m"),
							},
							DefaultRequest: v1.ResourceList{
								v1.ResourceCPU: resource.MustParse("800m"),
							},
							Max: v1.ResourceList{
								v1.ResourceCPU: resource.MustParse("800m"),
							},
							Min: v1.ResourceList{
								v1.ResourceCPU: resource.MustParse("200m"),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLimitRangeSepcFromConf(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLimitRangeSepcFromConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLimitRangeSepcFromConf() = %v, want %v", got, tt.want)
			}
		})
	}
}
