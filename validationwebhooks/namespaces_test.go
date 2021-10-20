package validationwebhooks

import (
	"context"
	"reflect"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.InjectDecoder(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("NamespaceValidator.InjectDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
