package validationwebhooks

import (
	"context"
	"fmt"
	. "github.com/toughnoah/melon/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	namingCheckError = "naming check failed: %s"

	decodeError = "decode error: %s"
)

type NamespaceValidator struct {
	Client   client.Client
	decoder  *admission.Decoder
	ConfPath string
}

func (v *NamespaceValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	namespace := &corev1.Namespace{}

	err := v.decoder.Decode(req, namespace)
	if err != nil {
		klog.Errorf(decodeError, err.Error())
		return admission.Errored(http.StatusBadRequest, err)
	}

	//exp := `^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)`
	err = ValidateNaming(v.ConfPath)
	if err != nil {
		klog.Errorf(decodeError, err.Error())
		return admission.Denied(fmt.Sprintf(decodeError, err.Error()))
	}
	return admission.Allowed("")
}

func (v *NamespaceValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
