#!/bin/bash
set -e

minishift start

oc login -u system:admin

oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-handler
oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-controller
oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-apiserver
oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-operator

oc create -f ./hack/tests/kubevirt.yaml

oc create -f ./cluster/okd/manifests/template-view-role.yaml
oc create -f ./cluster/okd/manifests/service.yaml
./wait-webhook.sh
./cluster/okd/extract-ca.sh ./cluster/okd/manifests/validating-webhook.yaml | oc apply -f -
