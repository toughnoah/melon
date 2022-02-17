package validationwebhooks

import (
	"context"
	"fmt"

	"net/http"

	. "github.com/toughnoah/melon/internal/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type ConfigmapValidator struct {
	Client   client.Client
	decoder  *admission.Decoder
	ConfPath string
}

func (v *ConfigmapValidator) Handle(_ context.Context, req admission.Request) admission.Response {
	cm := &corev1.ConfigMap{}
	err := v.decoder.Decode(req, cm)
	if err != nil {
		klog.Errorf(decodeError, err.Error())
		return admission.Errored(http.StatusBadRequest, err)
	}

	err = ValidateNaming(cm.Name, v.ConfPath, ConfigmapNamingKind)
	if err != nil {
		klog.Errorf(namingCheckError, ConfigmapNamingKind, err.Error())
		return admission.Denied(fmt.Sprintf(namingCheckError, ConfigmapNamingKind, err.Error()))
	}
	return admission.Allowed("")
}

func (v *ConfigmapValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
