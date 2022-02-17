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

	. "github.com/toughnoah/melon/internal/utils"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// reconcileReplicaSet reconciles ReplicaSets
type ReconcileLimitRange struct {
	// client can be used to retrieve objects from the APIServer.
	client.Client
	Path string
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &ReconcileLimitRange{}
var lLog = log.Log.WithName("limitRange-controller")

func (r *ReconcileLimitRange) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

	limitRange := &v1.LimitRange{}
	newLr, err := GetLimitRangeSepcFromConf(r.Path)
	if err != nil {
		lLog.Error(err, "error geting LimitRange from conf file", "NamespacedName", request.NamespacedName)
		return reconcile.Result{}, err
	}
	if newLr == nil {
		limitRange.Name = ResourceControlDefaultName
		limitRange.Namespace = request.Namespace
		if err = r.Delete(ctx, limitRange); err != nil && !k8serrors.IsNotFound(err) {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	newLr.Name = ResourceControlDefaultName
	newLr.Namespace = request.Namespace
	newLr.Annotations = map[string]string{
		ControllerDefaultAnnotationsKey: ControllerDefaultAnnotationsValue,
	}
	err = r.Get(ctx, request.NamespacedName, limitRange)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			lLog.Info("Starting to create limitRange", "NamespacedName", request.NamespacedName)
			if err = r.Create(ctx, newLr); err != nil {
				lLog.Error(err, "failed when creating limitRange", "NamespacedName", request.NamespacedName)
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	if !reflect.DeepEqual(newLr.Spec, limitRange.Spec) && limitRange.Annotations[ControllerDefaultAnnotationsKey] == ControllerDefaultAnnotationsValue {
		lLog.Info("Starting to update limitRange", "NamespacedName", request.NamespacedName)
		if err = r.Update(ctx, newLr); err != nil {
			lLog.Error(err, "failed when updating limitRange", "NamespacedName", request.NamespacedName)
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}
