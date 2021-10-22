/*



Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0


distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"

	. "github.com/toughnoah/melon/internal/utils"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/spf13/pflag"
	ctrl "github.com/toughnoah/melon/controller"
	vw "github.com/toughnoah/melon/validationwebhooks"
	v1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	file := pflag.StringP("config", "f", "", "specify the config file for melon")
	controllerNamespaceWhitelist := pflag.StringSlice("namespace-whitelist", []string{"kube-system"}, "specify the namespace whitelist for controller")
	help := pflag.BoolP("help", "h", false, "Show help message")
	pflag.Parse()
	if *help {
		pflag.Usage()
		os.Exit(0)
	}

	entryLog := log.Log.WithName("entrypoint")
	confPath := *file

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile ReplicaSets
	entryLog.Info("Setting up LimitRange controller")

	resconfig := config.GetConfigOrDie()
	mapper, err := apiutil.NewDynamicRESTMapper(resconfig)
	client, err := client.New(resconfig, client.Options{
		Scheme: scheme.Scheme,
		Mapper: mapper,
	})
	if err != nil {
		entryLog.Error(err, "failed to setup test client")
		os.Exit(1)
	}

	LimitRangeReconciler := &ctrl.ReconcileLimitRange{Client: client}
	ResourceQuotaReconciler := &ctrl.ReconcileResourceQuota{Client: client}

	if err := ctrl.SetUp(client, confPath, *controllerNamespaceWhitelist); err != nil {
		entryLog.Error(err, "failed to setup LimitRange and ResourceQuota")
		os.Exit(1)

	}
	lc, err := controller.New("LimitRange-controller", mgr, controller.Options{
		Reconciler: LimitRangeReconciler,
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	predicateFuncs := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return !In(e.ObjectOld.GetNamespace(), *controllerNamespaceWhitelist)
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return !In(e.Object.GetNamespace(), *controllerNamespaceWhitelist)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return !In(e.Object.GetNamespace(), *controllerNamespaceWhitelist)
		},
	}
	if err := lc.Watch(&source.Kind{Type: &v1.LimitRange{}}, &handler.EnqueueRequestForObject{}, predicateFuncs); err != nil {
		entryLog.Error(err, "unable to watch LimitRange")
		os.Exit(1)
	}

	entryLog.Info("Setting up ResourceQuota controller")
	rc, err := controller.New("ResourceQuota-controller", mgr, controller.Options{
		Reconciler: ResourceQuotaReconciler,
	})
	if err != nil {
		entryLog.Error(err, "unable to set up ResourceQuota controller")
		os.Exit(1)
	}
	if err := rc.Watch(&source.Kind{Type: &v1.ResourceQuota{}}, &handler.EnqueueRequestForObject{}, predicateFuncs); err != nil {
		entryLog.Error(err, "unable to watch ResourceQuota")
		os.Exit(1)
	}

	// Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	hookServer.Register("/validate-apps-v1-deployment", &webhook.Admission{Handler: &vw.DeploymentValidator{Client: client, ConfPath: confPath}})
	hookServer.Register("/validate-v1-namespace", &webhook.Admission{Handler: &vw.NamespaceValidator{Client: client, ConfPath: confPath}})
	hookServer.Register("/validate-v1-service", &webhook.Admission{Handler: &vw.ServiceValidator{Client: client, ConfPath: confPath}})
	hookServer.Register("/validate-v1-configmap", &webhook.Admission{Handler: &vw.ConfigmapValidator{Client: client, ConfPath: confPath}})
	hookServer.Register("/validate-v1alpha1-image", &webhook.Admission{Handler: &vw.ImageValidator{Client: client, ConfPath: confPath}})
	entryLog.Info("starting watching config")
	go ctrl.WatchConfig(client, LimitRangeReconciler, ResourceQuotaReconciler)

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
