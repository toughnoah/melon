package validationwebhooks

import (
	"context"
	"fmt"
	. "github.com/toughnoah/melon/internal/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type ServiceValidator struct {
	Client   client.Client
	decoder  *admission.Decoder
	ConfPath string
}

func (v *ServiceValidator) Handle(_ context.Context, req admission.Request) admission.Response {
	svc := &corev1.Service{}

	err := v.decoder.Decode(req, svc)
	if err != nil {
		klog.Errorf(decodeError, err.Error())
		return admission.Errored(http.StatusBadRequest, err)
	}

	err = ValidateNaming(svc.Name, v.ConfPath, ServiceNamingKind)
	if err != nil {
		klog.Errorf(namingCheckError, ServiceNamingKind, err.Error())
		return admission.Denied(fmt.Sprintf(namingCheckError, ServiceNamingKind, err.Error()))
	}
	return admission.Allowed("")
}

func (v *ServiceValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
