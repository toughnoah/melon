package validationwebhooks

import (
	"context"
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	testServicesNamingAllowed = `{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "noah-dev-svc-test"
  },
  "spec": {
    "selector": {
      "app": "nginx"
    },
    "ports": [
      {
        "protocol": "TCP",
        "port": 80,
        "targetPort": 80
      }
    ]
  }
}`
	testServicesNamingFailed = `{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "nginx-svc"
  },
  "spec": {
    "selector": {
      "app": "nginx"
    },
    "ports": [
      {
        "protocol": "TCP",
        "port": 80,
        "targetPort": 80
      }
    ]
  }
}`
)

func TestServiceValidator_Handle(t *testing.T) {
	type args struct {
		ctx context.Context
		req admission.Request
	}
	tests := []struct {
		name string
		v    *ServiceValidator
		args args
		want admission.Response
	}{
		{
			name: "test validate services pass",
			v: &ServiceValidator{
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
							Kind:    "Services",
						},
						Object: runtime.RawExtension{
							Raw:    []byte(testServicesNamingAllowed),
							Object: &v1.Namespace{},
						},
					},
				},
			},
			want: admission.Allowed(""),
		},
		{
			name: "test validate services denied",
			v: &ServiceValidator{
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
							Kind:    "Service",
						},
						Object: runtime.RawExtension{
							Raw:    []byte(testServicesNamingFailed),
							Object: &v1.Service{},
						},
					},
				},
			},
			want: admission.Denied(fmt.Sprintf(namingCheckError, "Services", "nginx-svc not match the expr ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Handle(tt.args.ctx, tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceValidator.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceValidator_InjectDecoder(t *testing.T) {
	type args struct {
		d *admission.Decoder
	}
	tests := []struct {
		name    string
		v       *ServiceValidator
		args    args
		wantErr bool
	}{
		{
			name: "test inject decoder",
			v: &ServiceValidator{
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
				t.Errorf("ServiceValidator.InjectDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
