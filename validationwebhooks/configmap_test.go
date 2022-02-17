package validationwebhooks

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	testConfigNamingPass = `{
  "apiVersion": "v1",
  "data": {
    "config.yaml": "default_expr: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)\n## deployment naming expr can override default_expr\ndeploy_expr: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)\n## namespaces naming expr can override default_expr\nns_expr: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)\nimg_expr: ^(?:docker.io)/(?:toughnoah|test)/.+?:v1.0\nis_validate_deploy_limits: true\n"
  },
  "kind": "ConfigMap",
  "metadata": {
    "name": "noah-sa-melon-config-test"
  }
}`
	testConfigNamingFail = `{
  "apiVersion": "v1",
  "data": {
    "config.yaml": "default_expr: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)\n## deployment naming expr can override default_expr\ndeploy_expr: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)\n## namespaces naming expr can override default_expr\nns_expr: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)\nimg_expr: ^(?:docker.io)/(?:toughnoah|test)/.+?:v1.0\nis_validate_deploy_limits: true\n"
  },
  "kind": "ConfigMap",
  "metadata": {
    "name": "melon-config-test"
  }
}`
)

func TestConfigmapValidator_Handle(t *testing.T) {
	type args struct {
		ctx context.Context
		req admission.Request
	}
	tests := []struct {
		name string
		v    *ConfigmapValidator
		args args
		want admission.Response
	}{
		{
			name: "test validate configmaps pass",
			v: &ConfigmapValidator{
				Client:   fake.NewClientBuilder().Build(),
				ConfPath: "../tests/testdata",
				decoder:  decoder,
			},
			args: args{
				ctx: ctx,
				req: admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						UID: "fake_request_allowed",
						RequestKind: &metav1.GroupVersionKind{
							Group:   "",
							Version: "v1",
							Kind:    "ConfigMap",
						},
						Object: runtime.RawExtension{
							Raw:    []byte(testConfigNamingPass),
							Object: &v1.ConfigMap{},
						},
					},
				},
			},
			want: admission.Allowed(""),
		},
		{
			name: "test validate configmaps denied",
			v: &ConfigmapValidator{
				Client:   fake.NewClientBuilder().Build(),
				ConfPath: "../tests/testdata",
				decoder:  decoder,
			},
			args: args{
				ctx: ctx,
				req: admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						UID: "fake_request_allowed",
						RequestKind: &metav1.GroupVersionKind{
							Group:   "",
							Version: "v1",
							Kind:    "ConfigMap",
						},
						Object: runtime.RawExtension{
							Raw:    []byte(testConfigNamingFail),
							Object: &v1.ConfigMap{},
						},
					},
				},
			},
			want: admission.Denied(fmt.Sprintf(namingCheckError, "configmap.naming", "melon-config-test not match the expr ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Handle(tt.args.ctx, tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConfigmapValidator.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigmapValidator_InjectDecoder(t *testing.T) {
	type args struct {
		d *admission.Decoder
	}
	tests := []struct {
		name    string
		v       *ConfigmapValidator
		args    args
		wantErr bool
	}{
		{
			name: "test inject decoder",
			v: &ConfigmapValidator{
				Client:   fake.NewClientBuilder().Build(),
				ConfPath: "../tests/testdata",
				decoder:  decoder,
			},
			args: args{
				d: decoder,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.InjectDecoder(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("ConfigmapValidator.InjectDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
