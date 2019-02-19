# kubevirt-template-validator

`kubevirt-template-validator` is a [kubevirt](http://kubevirt.io) addon to check the [annotations on templates](https://github.com/kubevirt/common-templates/blob/master/templates/VALIDATION.md) and reject them if unvalid.
It is implemented using a [validating webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/).

[![Go Report Card](https://goreportcard.com/badge/github.com/fromanirh/kubevirt-template-validator)](https://goreportcard.com/report/github.com/fromanirh/kubevirt-template-validator)

## License

Apache v2

## Dependencies

* [kubernetes APIs](https://github.com/kubernetes/kubernetes)


## Installation

### installation on kubernetes (K8S)

1. (ODK/OCP) first, make sure you have the `template:view` cluster role binding in your cluster. If not, add it:
```bash
KUBECTL=oc create -f ./cluster/manifests/template-view-role.yaml
```

2. first, create and deploy the certificates in a Kubernetes Secret, to be used in the following steps:
```bash
KUBECTL=kubectl ./cluster/webhook-create-signed-cert.sh
```

2.a. check that the secret exists:
```bash
kubectl  get secret -n kubevirt virtualmachine-template-validator-certs
NAME                                      TYPE      DATA      AGE
virtualmachine-template-validator-certs   Opaque    2         1h
```

3. deploy the service:
```bash
kubectl create -f ./cluster/manifests/service.yaml
```

4. In order to set up the webhook, we need a CA bundle. We can reuse the one from the certs we create from the step #1.
```bash
cat ./cluster/manifests/validating-webhook.yaml | ./cluster/extract-ca.sh | kubectl apply -f -
```

Done!

### Disable the webhook

To disable the webhook, just de-register it from the apiserver:
```bash
kubectl delete -f ./cluster/manifests/validating-webhook.yaml
```

## Caveats & Gotchas

content pending
