package controller

import (
	"context"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func WatchConfig(c client.Client, r ...reconcile.Reconciler) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		rLog.Info("start to reconcile resourcesSpec from conf")
		namespaces, err := GetAllNamespace(c)
		if err != nil {
			return
		}
		for _, rs := range r {

			for _, ns := range namespaces.Items {
				namespaceName := types.NamespacedName{
					Name:      ResourceControlDefaultName,
					Namespace: ns.Name,
				}
				go rs.Reconcile(context.Background(), reconcile.Request{
					NamespacedName: namespaceName,
				})
			}
		}
	})

}
