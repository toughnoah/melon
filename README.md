![Go 1.16](https://img.shields.io/badge/Go-v1.16-blue)
[![CI Workflow](https://github.com/toughnoah/melon/actions/workflows/test-coverage.yaml/badge.svg?branch=master)](https://github.com/toughnoah/melon/actions/workflows/test-coverage.yaml)
[![codecov](https://codecov.io/gh/toughnoah/melon/branch/master/graph/badge.svg?token=Wa1IqU4OCF)](https://codecov.io/gh/toughnoah/melon)

# Melon
This is a project of validating admission webhook for Kubernetes to in memory of melon. The name melon comes from my another naughty fat house cat. Indeed, he had passed few days ago because of bad disease before this projects starts, and he had accompanied me for over four and a half years.

`Melon` is for validating some Kubernetes resources such as ***naming***`namespace`,`deployment`, `configmap`,`service`, ***and checking*** `contaners[].resources.limits`, `contaners[].image` by using regexp.

I am so miss my boy, and I will love him forever.

## Quick Start

### Sign cert
The webhook server is designed for the safest connection in cluster via `tls`, so we have to sign our own cert. 
```shell
openssl genrsa -out server.key 2048
```

```shell
openssl req -new -key server.key -out server.csr
```

after go 1.15, we have to use SAN, using the ca.key of apiserver for signing cert.
```shell
openssl x509 -req -extfile <(printf "subjectAltName=DNS:$(your webhook svc domain, such as melon.default.svc") -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
```

### Deploy
```shell
git clone https://github.com/toughnoah/melon.git
cd deploy
```
change the name and namespace of the configmap, deployment, service that you signed for in previous step.

Then
```shell
kubectl apply -k ./deploy
```

### Create webhook
Specify a ValidatingWebhookConfiguration
```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: melon-configmaps
webhooks:
  - admissionReviewVersions:
      - v1
      - v1beta1
    clientConfig:
      # change to your own caBundle
      caBundle: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeE1Ea3hPREEyTWpBeU9Gb1hEVE14TURreE5qQTJNakF5T0Zvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTzlaCmo5RDFobFBibG5PNG5sc09XbVVEQlB2ajF3bE9Pck16dWhpZjdGV0txeXErUlgwbngxcTU1RjFBNkxBYUxFTkQKSmZqcHJQS0tPbisva1Z1bmJnUmZIUWJucHl4UGcyZ2dObVhEc2VFOEo1ZzFnSGJVRkdtZVVvUGxiYzY1VzN6OQo5SWk0dlozNVp2blRDazZjZW5raDdnK0xLbTN3L2wyVUpNZDZ3UTV5ZkVzTVJYSGhJTGVtTXYyMjVnY3poRzZiCkxPQlJlcGM1RDlEU1QwSTZBeE1OUHR2cWkzcFcrdjNJZ0hQZnlMczJLM241NzlQZXRTbVhWWHg0eWVxdVZ5TkgKRExBYk1mbXl0Qk9SS3JJYzRvZjNhREwwdVF3TkxEZ3NIN3N2TDM1TnpnellTZ3h0SUdxK1RsRGltekp1Tk0wegpaa2IxbmptdnlZZ0xiTS9mNjY4Q0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZMUzBaRGNPbEl6bTRJZUlmZDh2YlZWY0NibVZNQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFCMVdSbW55akdQSnp3RWd5YnAzejZ4YmVXL2hHK2p1R1kwUC92WmFUQ2dKVnFCbTN5SgpoMUgrWnZkTGxQK0lkS2JUQUpmUUtnN05UcGRDWnlGWDdQOWxwaXNTeG9ENG1IRjFPbmRLTW5XOGdtejlzUTl1CmF1NTBHRi83N290VXlmL0pERmoyVWVIMDVvTllYOXVXSlpYcDdqajgyRGt5bjIxVnU4T1FrQlloNGdIc3RWM0MKUFRGZ09hTjZ0L05mZnBHanFWYXZBL0R1NU4xY2NXV0hLVmdFQ3F1ZXl5Ym5Xdjd5OThUNHkvVkF4RStNeHpQVgpyd0R4V3MxaDROL2ZUUVBnTzBnL3gvQlZvNGowblcwVzVQRTFDanVXcGRNSjk4Yk1pcHd3TWZNN0pFWVNqdHhmCkkvaWlRM1dyZnBqclJjU2dCeU5rQmVNeHJKcjdCNDBTcUowVwotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t
      service:
        name: melon
        namespace: default
        path: /validate-v1-configmap
        port: 9443
    name: admissionvalidationwebhook.melon.io
    namespaceSelector:
      matchExpressions:
        - key: kubernetes.io/metadata.name
          operator: NotIn
          values:
            - kube-system
            - isito-system
            - default
    rules:
      - apiGroups:
          - ""
        apiVersions:
          - v1
        operations:
          - CREATE
          - UPDATE
        resources:
          - configmaps
        scope: '*'
    sideEffects: None
    timeoutSeconds: 5
```

## More
Now melon only supports the keys shown below. I am working on it.

New feature is comming soon. Any advice is welcome. Please stay.
```yaml
#global expr for all Kind
global:
  naming: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)

# deployment naming expr can override default_expr
deployment:
  naming: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)
  image: ^(?:docker.io)/(?:toughnoah|test)/.+?:v1.0
  limits: true

# namespaces naming expr can override default_expr
namespace:
  naming: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)

# service naming expr can override default_expr
service:
  naming: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)

configmap:
  naming: ^(?:noah|blackbean|melon)-(?:dev|qa|sa)-.+?-(?:test|prod)
```