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

package validationwebhooks

import (
	"context"
	"errors"
	"fmt"
	. "github.com/toughnoah/melon/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/klog/v2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DeploymentValidator validates deployment
type DeploymentValidator struct {
	Client   client.Client
	decoder  *admission.Decoder
	ConfPath string
}

// Handle podValidator admits a pod if a specific annotation exists.
func (v *DeploymentValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	deploy := &appsv1.Deployment{}
	err := v.decoder.Decode(req, deploy)
	if err != nil {
		klog.Errorf(decodeError, err.Error())
		return admission.Errored(http.StatusBadRequest, err)
	}

	if err = ValidateNaming(deploy.Name, v.ConfPath, DeploymentNamingKind); err != nil {
		klog.Errorf(namingCheckError, DeploymentNamingKind, err.Error())
		return admission.Denied(fmt.Sprintf(namingCheckError, DeploymentNamingKind, err.Error()))
	}
	if IsToValidateLimits(v.ConfPath, DeploymentLimitsKind) {
		err = validateResources(deploy)
		if err != nil {
			return admission.Denied(err.Error())
		}
	}
	if err = validateImageNaming(deploy, v.ConfPath); err != nil {
		klog.Errorf(namingCheckError, DeploymentImageKind, err.Error())
		return admission.Denied(fmt.Sprintf(namingCheckError, DeploymentImageKind, err.Error()))
	}
	return admission.Allowed("")

}

// InjectDecoder injects the decoder.
func (v *DeploymentValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

func validateResources(deploy *appsv1.Deployment) error {

	containerArray := deploy.Spec.Template.Spec.Containers
	if len(containerArray) == 0 {
		return errors.New(noContainerError)
	}
	for _, container := range containerArray {
		if len(container.Resources.Limits) == 0 {
			return errors.New(noResourcesLimitsError)
		}
	}
	return nil
}

//
func validateImageNaming(deploy *appsv1.Deployment, confPath string) error {
	containerArray := deploy.Spec.Template.Spec.Containers
	if len(containerArray) == 0 {
		return errors.New(noContainerError)
	}
	for _, container := range containerArray {
		err := ValidateNaming(container.Image, confPath, DeploymentImageKind)
		if err != nil {
			return err
		}
	}
	return nil
}
