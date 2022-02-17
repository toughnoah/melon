package controller

import (
	"context"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_setUpLimitRange(t *testing.T) {
	type args struct {
		ctx           context.Context
		path          string
		c             client.Client
		namespaceName types.NamespacedName
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setUpLimitRange(tt.args.ctx, tt.args.path, tt.args.c, tt.args.namespaceName); (err != nil) != tt.wantErr {
				t.Errorf("setUpLimitRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setUpResourceQuota(t *testing.T) {
	type args struct {
		ctx           context.Context
		path          string
		c             client.Client
		namespaceName types.NamespacedName
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setUpResourceQuota(tt.args.ctx, tt.args.path, tt.args.c, tt.args.namespaceName); (err != nil) != tt.wantErr {
				t.Errorf("setUpResourceQuota() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAllNamespace(t *testing.T) {
	type args struct {
		c client.Client
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.NamespaceList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllNamespace(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}
