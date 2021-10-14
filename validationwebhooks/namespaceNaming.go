package validationwebhooks

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"net/http"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type NamespaceValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (v *NamespaceValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	namespace := &corev1.Namespace{}

	err := v.decoder.Decode(req, namespace)
	if err != nil {
		klog.Errorf("decode error: %s", err.Error())
		return admission.Errored(http.StatusBadRequest, err)
	}

	exp := `^(?:mgb|fdc)-(?:dba|sre|bigdata)-.+?-(?:test|prod)`
	isValid, err := v.validateNaming(namespace.Name, exp)
	if err != nil {
		klog.Errorf("regexp expr compile error: %s", err.Error())
		return admission.Errored(http.StatusInternalServerError, err)
	}
	if !isValid {
		klog.Errorf("namespace naming dose not match the exp: %s", exp)
		return admission.Denied(fmt.Sprintf("namespace naming dose not match the exp: %s", exp))
	}
	return admission.Allowed("")
}

func (v *NamespaceValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

func (v *NamespaceValidator) validateNaming(str string, exp string) (bool, error) {
	reg, err := regexp.Compile(exp)
	if err != nil {
		return false, err
	}

	match := reg.FindAllStringSubmatch(str, -1)
	if match != nil {
		return true, nil
	} else {
		return false, nil
	}
}
