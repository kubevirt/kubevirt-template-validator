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
sleep 10s
