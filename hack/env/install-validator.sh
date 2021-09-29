#!/bin/bash
set -e

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )
MANIFESTPATH="${BASEPATH}/../../cluster/okd"

NAMESPACE=${1:-kubevirt}

# intentionally put in the "default" namespace
oc create -n default -f ${BASEPATH}/common-templates.yaml
sleep 10s

oc create -f ${MANIFESTPATH}/manifests/template-view-role.yaml

sed "s|image:.*|image: ${VALIDATOR_IMAGE}|" < ${MANIFESTPATH}/manifests/service.yaml | \
	sed "s|imagePullPolicy: Always|imagePullPolicy: IfNotPresent|g" | \
	oc create -f -

${BASEPATH}/wait-webhook.sh
sleep 10s

${MANIFESTPATH}/extract-ca.sh ${MANIFESTPATH}/manifests/validating-webhook.yaml | oc apply -n ${NAMESPACE} -f -

oc get pods -n ${NAMESPACE}
oc describe pod -n ${NAMESPACE} -l "kubevirt.io=virt-template-validator"
