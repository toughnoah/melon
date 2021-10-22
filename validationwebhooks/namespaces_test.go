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
	testNamespaceAllowed = `
apiVersion: v1
kind: Namespace
metadata:
  name: noah-dev-melon-test
`
	testNamespaceDenied = `
apiVersion: v1
kind: Namespace
metadata:
  name: noah-test
`
)

func TestNamespaceValidator_Handle(t *testing.T) {
	type args struct {
		ctx context.Context
		req admission.Request
	}
	tests := []struct {
		name string
		v    *NamespaceValidator
		args args
		want admission.Response
	}{

		{
			name: "test validate passe",
			v: &NamespaceValidator{
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
							Kind:    "Namespace",
						},
						Object: runtime.RawExtension{
							Raw:    []byte(testNamespaceAllowed),
							Object: &v1.Namespace{},
						},
					},
				},
			},
			want: admission.Allowed(""),
		},
		{
			name: "test validate denied",
			v: &NamespaceValidator{
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
							Kind:    "Namespace",
						},
						Object: runtime.RawExtension{
							Raw:    []byte(testNamespaceDenied),
							Object: &v1.Namespace{},
						},
					},
				},
			},
			want: admission.Denied(fmt.Sprintf(namingCheckError, deniedErrorMessage)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Handle(tt.args.ctx, tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NamespaceValidator.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespaceValidator_InjectDecoder(t *testing.T) {
	type args struct {
		d *admission.Decoder
	}
	tests := []struct {
		name    string
		v       *NamespaceValidator
		args    args
		wantErr bool
	}{
		{
			name: "test inject decoder",
			v: &NamespaceValidator{
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
				t.Errorf("NamespaceValidator.InjectDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
