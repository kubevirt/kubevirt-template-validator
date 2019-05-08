#!/bin/bash
set -e

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )
MANIFESTPATH="${BASEPATH}/../../cluster/okd"

minishift start
# see https://github.com/minishift/minishift/pull/3044 for details
minishift addons install --defaults
minishift addons apply admissions-webhook

oc login -u system:admin

oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-handler
oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-controller
oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-apiserver
oc adm policy add-scc-to-user privileged -n kubevirt -z kubevirt-operator

oc create -f ${BASEPATH}/kubevirt.yaml

oc create -f ${MANIFESTPATH}/manifests/template-view-role.yaml
oc create -f ${MANIFESTPATH}/manifests/service.yaml
${BASEPATH}/wait-webhook.sh
${MANIFESTPATH}/extract-ca.sh ${MANIFESTPATH}/manifests/validating-webhook.yaml | oc apply -f -
