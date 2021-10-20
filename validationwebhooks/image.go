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
	"fmt"
	"net/http"

	. "github.com/toughnoah/melon/internal/utils"

	"k8s.io/api/imagepolicy/v1alpha1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DeploymentValidator validates deployment
type ImageValidator struct {
	Client   client.Client
	decoder  *admission.Decoder
	ConfPath string
}

// Handle imageValidator admits a pod if a specific annotation exists.
func (v *ImageValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	image := &v1alpha1.ImageReview{}
	err := v.decoder.Decode(req, image)
	if err != nil {
		klog.Errorf(decodeError, err.Error())
		return admission.Errored(http.StatusBadRequest, err)
	}
	for _, container := range image.Spec.Containers {
		if err := ValidateNaming(container.Image, v.ConfPath, ImageKind); err != nil {
			return admission.Denied(fmt.Sprintf(namingCheckError, ImageKind, err.Error()))
		}
	}
	return admission.Allowed("")

}

// InjectDecoder injects the decoder.
func (v *ImageValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

//
