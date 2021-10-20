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
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"testing"
)

var (
	ctx        = context.Background()
	decoder, _ = admission.NewDecoder(scheme.Scheme)
)

const (
	testDeploymentAllowed = `{
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "name": "noah-dev-deployment-test",
    "labels": {
      "app": "nginx"
    }
  },
  "spec": {
    "replicas": 1,
    "selector": {
      "matchLabels": {
        "app": "nginx"
      }
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "nginx"
        }
      },
      "spec": {
        "containers": [
          {
            "name": "web",
            "image": "nginx",
            "ports": [
              {
                "name": "web",
                "containerPort": 80
              }
            ],
            "volumeMounts": [
              {
                "name": "html",
                "mountPath": "/usr/share/nginx/html"
              }
            ]
          }
        ],
        "volumes": [
          {
            "name": "html",
            "persistentVolumeClaim": {
              "claimName": "efs-claim-expand-test"
            }
          }
        ]
      }
    }
  }
}`
	testDeploymentFailed = `{
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "name": "nginx-deployment-test",
    "labels": {
      "app": "nginx"
    }
  },
  "spec": {
    "replicas": 1,
    "selector": {
      "matchLabels": {
        "app": "nginx"
      }
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "nginx"
        }
      },
      "spec": {
        "containers": [
          {
            "name": "web",
            "image": "nginx",
            "ports": [
              {
                "name": "web",
                "containerPort": 80
              }
            ],
            "volumeMounts": [
              {
                "name": "html",
                "mountPath": "/usr/share/nginx/html"
              }
            ]
          }
        ],
        "volumes": [
          {
            "name": "html",
            "persistentVolumeClaim": {
              "claimName": "efs-claim-expand-test"
            }
          }
        ]
      }
    }
  }
}`
	deniedErrorMessage = `not match the expr ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)`
)

func TestDeploymentValidator_Handle(t *testing.T) {
	type args struct {
		ctx context.Context
		req admission.Request
	}
	tests := []struct {
		name string
		v    *DeploymentValidator
		args args
		want admission.Response
	}{
		{
			name: "test validate passe",
			v: &DeploymentValidator{
				Client:   fake.NewClientBuilder().Build(),
				ConfPath: "../pkg/testdata",
				decoder:  decoder,
			},
			args: args{
				ctx: ctx,
				req: admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						UID: "fake_request_allowed",
						RequestKind: &metav1.GroupVersionKind{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Object: runtime.RawExtension{
							Raw:    []byte(testDeploymentAllowed),
							Object: &appsv1.Deployment{},
						},
					},
				},
			},
			want: admission.Allowed(""),
		},
		{
			name: "test validate failed",
			v: &DeploymentValidator{
				Client:   fake.NewClientBuilder().Build(),
				ConfPath: "../pkg/testdata",
				decoder:  decoder,
			},
			args: args{
				ctx: ctx,
				req: admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						UID: "fake_request_allowed",
						RequestKind: &metav1.GroupVersionKind{
							Group:   "apps",
							Version: "v1",
							Kind:    "Deployment",
						},
						Object: runtime.RawExtension{
							Raw:    []byte(testDeploymentFailed),
							Object: &appsv1.Deployment{},
						},
					},
				},
			},
			want: admission.Denied(fmt.Sprintf(namingCheckError, deniedErrorMessage)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Handle(tt.args.ctx, tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeploymentValidator.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeploymentValidator_InjectDecoder(t *testing.T) {
	type args struct {
		d *admission.Decoder
	}
	tests := []struct {
		name    string
		v       *DeploymentValidator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.InjectDecoder(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("DeploymentValidator.InjectDecoder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
