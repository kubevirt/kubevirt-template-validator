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

if [ "${CI}" == "true" ] && [ "${TRAVIS}" == "true" ]; then
	REGISTRY=$(minishift openshift registry)

	sed "s|image:.*|image: ${REGISTRY}/kubevirt/kubevirt-template-validator:devel|" < ${MANIFESTPATH}/manifests/service.yaml | \
		sed "s|imagePullPolicy: Always|imagePullPolicy: IfNotPresent|g" | \
		oc create -f -
else
	oc create -n ${NAMESPACE} -f ${MANIFESTPATH}/manifests/service.yaml
fi

${BASEPATH}/wait-webhook.sh
sleep 10s

${MANIFESTPATH}/extract-ca.sh ${MANIFESTPATH}/manifests/validating-webhook.yaml | oc apply -n ${NAMESPACE} -f -

if [ "${CI}" == "true" ] && [ "${TRAVIS}" == "true" ]; then
	sleep 5s

	oc get pods -n ${NAMESPACE}
	oc get pod -n ${NAMESPACE} -l "kubevirt.io=virt-template-validator" -o yaml

	sleep 5s
fi
