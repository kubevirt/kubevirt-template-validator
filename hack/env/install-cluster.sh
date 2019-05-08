#!/bin/bash
set -e

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )
MANIFESTPATH="${BASEPATH}/../../cluster/okd"

oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-handler
oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-controller
oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-apiserver
oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-operator

oc create -f ${BASEPATH}/kubevirt.yaml
oc project kubevirt
sleep 10s # more cool down

# required by functests
oc create -n default -f ${BASEPATH}/common-templates.yaml

sleep 10s # more cool down

oc create -f ${MANIFESTPATH}/manifests/template-view-role.yaml
oc create -f ${MANIFESTPATH}/manifests/service.yaml
${BASEPATH}/wait-webhook.sh
sleep 10s # more cool down
${MANIFESTPATH}/extract-ca.sh ${MANIFESTPATH}/manifests/validating-webhook.yaml | oc apply -f -
