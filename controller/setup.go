package controller

import (
	"context"

	. "github.com/toughnoah/melon/internal/utils"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var sLog = log.Log.WithName("setup")

func SetUp(c client.Client, path string, setUpWhitelist []string) error {
	ctx := context.Background()
	namespaces, err := GetAllNamespace(c)
	if err != nil {
		return err

	}
	for _, ns := range namespaces.Items {
		if In(ns.Name, setUpWhitelist) {
			continue
		}
		namespaceName := types.NamespacedName{
			Name:      ResourceControlDefaultName,
			Namespace: ns.Name,
		}
		if err := setUpLimitRange(ctx, path, c, namespaceName); err != nil {
			sLog.Error(err, "error setup LimitRange", "Namespace", namespaceName.Name)
		}

		if err := setUpResourceQuota(ctx, path, c, namespaceName); err != nil {
			sLog.Error(err, "error setup ResourceQuota", "Namespace", namespaceName.Name)
		}
	}
	return nil
}

func setUpLimitRange(ctx context.Context, path string, c client.Client, namespaceName types.NamespacedName) error {
	lr := &v1.LimitRange{}

	newLr, err := GetLimitRangeSepcFromConf(path)
	if err != nil {
		return err
	}
	if newLr == nil {
		return nil
	}
	err = c.Get(ctx, namespaceName, lr)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			newLr.Name = ResourceControlDefaultName
			newLr.Namespace = namespaceName.Namespace
			newLr.Annotations = map[string]string{
				ControllerDefaultAnnotationsKey: ControllerDefaultAnnotationsValue,
			}
			sLog.Info("Start to set up LimitRange", "Namespace", namespaceName.Namespace)
			return c.Create(ctx, newLr)
		}
		return err
	}
	lr.Spec = newLr.Spec
	sLog.Info("Start to update LimitRange", "Namespace", namespaceName.Namespace)
	return c.Update(ctx, lr)
}

func setUpResourceQuota(ctx context.Context, path string, c client.Client, namespaceName types.NamespacedName) error {
	rq := &v1.ResourceQuota{}

	newLr, err := GetResourceQuotaSepcFromConf(path)
	if err != nil {
		return err
	}
	if newLr == nil {
		return nil
	}

	err = c.Get(ctx, namespaceName, rq)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			newLr.Name = ResourceControlDefaultName
			newLr.Namespace = namespaceName.Namespace
			newLr.Generation = 1
			newLr.Annotations = map[string]string{
				ControllerDefaultAnnotationsKey: ControllerDefaultAnnotationsValue,
			}
			sLog.Info("Start to set up ResourceQuota", "Namespace", namespaceName.Namespace)
			return c.Create(ctx, newLr)
		}
		return err
	}
	rq.Spec = newLr.Spec
	rq.Generation += 1
	sLog.Info("Start to update ResourceQuota", "Namespace", namespaceName.Namespace)

	return c.Update(ctx, rq)
}

func GetAllNamespace(c client.Client) (*v1.NamespaceList, error) {
	namespaces := &v1.NamespaceList{}
	if err := c.List(context.Background(), namespaces); err != nil {
		return nil, err
	}
	return namespaces, nil

}
