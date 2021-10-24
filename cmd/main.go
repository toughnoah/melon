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

package main

import (
	"github.com/spf13/pflag"
	"os"

	wb "github.com/toughnoah/melon/validationwebhooks"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	entryLog := log.Log.WithName("entrypoint")
	file := pflag.StringP("config", "f", "", "specify the config file for melon")
	help := pflag.BoolP("help", "h", false, "Show help message")
	pflag.Parse()
	if *help {
		pflag.Usage()
		os.Exit(0)
	}

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	entryLog.Info("registering webhooks to the webhook server")
	confPath := *file
	hookServer.Register("/validate-apps-v1-deployment", &webhook.Admission{Handler: &wb.DeploymentValidator{Client: mgr.GetClient(), ConfPath: confPath}})
	hookServer.Register("/validate-v1-namespace", &webhook.Admission{Handler: &wb.NamespaceValidator{Client: mgr.GetClient(), ConfPath: confPath}})
	hookServer.Register("/validate-v1-service", &webhook.Admission{Handler: &wb.ServiceValidator{Client: mgr.GetClient(), ConfPath: confPath}})
	hookServer.Register("/validate-v1-configmap", &webhook.Admission{Handler: &wb.ConfigmapValidator{Client: mgr.GetClient(), ConfPath: confPath}})

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
