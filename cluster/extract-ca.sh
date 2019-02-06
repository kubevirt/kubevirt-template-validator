#!/bin/sh

ROOT=$(cd $(dirname $0)/../../; pwd)

set -o errexit
set -o nounset
set -o pipefail

[ -z ${KUBECTL} ] && KUBECTL=kubectl

CSR="virtualmachine-template-validator.kubevirt"

export CA_BUNDLE=$( $KUBECTL get csr $CSR -o json | jq -r '.status.certificate' )

if command -v envsubst >/dev/null 2>&1; then
    envsubst
else
    sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g"
fi
