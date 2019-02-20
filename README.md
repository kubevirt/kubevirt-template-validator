# kubevirt-template-validator

`kubevirt-template-validator` is a [kubevirt](http://kubevirt.io) addon to check the [annotations on templates](https://github.com/kubevirt/common-templates/blob/master/templates/VALIDATION.md) and reject them if unvalid.
It is implemented using a [validating webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/).

[![Go Report Card](https://goreportcard.com/badge/github.com/fromanirh/kubevirt-template-validator)](https://goreportcard.com/report/github.com/fromanirh/kubevirt-template-validator)

## License

Apache v2

## Dependencies

* [kubernetes APIs](https://github.com/kubernetes/kubernetes)


## Installation

You need to pick the platform on which you want to install.
For kubernetes:
```bash
export PLATFORM=k8s
```
for OKD/OCP:
```bash
export PLATFORM=okd
```

now you can set which tool you need to use to interact with the cluster. Usually:
for kubernetes:
```bash
export KUBECTL=kubectl
```
for OKD/OCP:
```bash
export KUBECTL=oc
```

### installation on OKD/OCP

Make sure the validating webhooks are enabled. You either need to configure the platform when you install it
or to use OKD/OCP >= 4.0. See:
- https://github.com/openshift/origin/issues/20842
- https://github.com/openshift/openshift-ansible/issues/7983

Then, make sure you have the `template:view` cluster role binding in your cluster. If not, add it:
```bash
$KUBECTL create -f ./cluster/manifests/okd/template-view-role.yaml
```

### common installation instructions

1. first, create and deploy the certificates in a Kubernetes Secret, to be used in the following steps:
```bash
$KUBECTL ./cluster/$PLATFORM/webhook-create-signed-cert.sh
```

2.a. check that the secret exists:
```bash
$KUBECTL get secret -n kubevirt virtualmachine-template-validator-certs
NAME                                      TYPE      DATA      AGE
virtualmachine-template-validator-certs   Opaque    2         1h
```

3. deploy the service:
```bash
$KUBECTL create -f ./cluster/$PLATFORM/manifests/service.yaml
```

4. In order to set up the webhook, we need a CA bundle. We can reuse the one from the certs we create from the step #1.
```bash
cat ./cluster/$PLATFORM/manifests/validating-webhook.yaml | ./cluster/$PLATFORM/extract-ca.sh | $KUBECTL apply -f -
```

Done!

### Disable the webhook

To disable the webhook, just de-register it from the apiserver:
```bash
$KUBECTL delete -f ./cluster/$PLATFORM/manifests/validating-webhook.yaml
```

## Caveats & Gotchas

content pending
