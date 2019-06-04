#!/bin/sh

ROOT=$(cd $(dirname $0)/../../; pwd)

set -o errexit
set -o nounset
#set -o pipefail

export OC=oc

VALIDATOR_POD=$( ${OC} get pod -n kubevirt -l kubevirt.io=virt-template-validator -o json | jq -r .items[0].metadata.name)
export CA_BUNDLE=$( ${OC} exec -ti -n kubevirt $VALIDATOR_POD -- /bin/cat /var/run/secrets/kubernetes.io/serviceaccount/service-ca.crt | base64 -w 0 )

if command -v envsubst >/dev/null 2>&1; then
    envsubst < $1
else
    sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g" < $1
fi
