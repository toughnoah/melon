/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"reflect"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	. "github.com/toughnoah/melon/internal/utils"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// reconcileReplicaSet reconciles ReplicaSets
type ReconcileResourceQuota struct {
	// client can be used to retrieve objects from the APIServer.
	client.Client
	Path string
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &ReconcileResourceQuota{}
var rLog = log.Log.WithName("resourceQuota-controller")

func (r *ReconcileResourceQuota) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	resourceQuota := &v1.ResourceQuota{}
	newRq, err := GetResourceQuotaSepcFromConf(r.Path)
	if err != nil {
		rLog.Error(err, "error geting ResourceQuota from conf file", "NamespacedName", request.NamespacedName)
		return reconcile.Result{}, err
	}
	if newRq == nil {
		resourceQuota.Name = ResourceControlDefaultName
		resourceQuota.Namespace = request.Namespace
		if err = r.Delete(ctx, resourceQuota); err != nil && !k8serrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}
	}
	newRq.Name = ResourceControlDefaultName
	newRq.Namespace = request.Namespace
	newRq.Annotations = map[string]string{
		ControllerDefaultAnnotationsKey: ControllerDefaultAnnotationsValue,
	}
	err = r.Get(ctx, request.NamespacedName, resourceQuota)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			rLog.Info("Starting to create ResourceQuota", "NamespacedName", request.NamespacedName)
			if err = r.Create(ctx, newRq); err != nil {
				rLog.Error(err, "failed when creating ResourceQuota", "NamespacedName", request.NamespacedName)
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, err
	}
	if !reflect.DeepEqual(resourceQuota.Spec, newRq.Spec) && resourceQuota.Annotations[ControllerDefaultAnnotationsKey] == ControllerDefaultAnnotationsValue {
		rLog.Info("Starting to update ResourceQuota", "NamespacedName", request.NamespacedName)
		if err = r.Update(ctx, newRq); err != nil {
			rLog.Error(err, "failed when updating ResourceQuota", "NamespacedName", request.NamespacedName)
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}
