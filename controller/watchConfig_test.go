package controller

import (
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestWatchConfig(t *testing.T) {
	type args struct {
		c client.Client
		r reconcile.Reconciler
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WatchConfig(tt.args.c, tt.args.r)
		})
	}
}
