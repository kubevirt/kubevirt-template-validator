#!/bin/sh

ROOT=$(cd $(dirname $0)/../../; pwd)

set -o errexit
set -o nounset
set -o pipefail

[ -z ${KUBECTL} ] && KUBECTL=oc
[ -z ${NAMESPACE} ] && NAMESPACE=kubevirt

SECRET="virtualmachine-template-validator-certs"

export CA_BUNDLE=$( $KUBECTL get secret -n $NAMESPACE $SECRET -o json | jq -r '.data["cert.pem"]' )

if command -v envsubst >/dev/null 2>&1; then
    envsubst
else
    sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g"
fi
