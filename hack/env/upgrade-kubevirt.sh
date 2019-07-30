#!/bin/bash
set -ex

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )

# we can't guess.
[ -z "$1" ] && exit 1
RELEASE="${1}"

# avoids https://github.com/kubevirt/kubevirt/issues/2533
oc apply -f https://github.com/kubevirt/kubevirt/releases/download/${RELEASE}/kubevirt-operator.yaml
# TODO: smarter wait
sleep 42s
oc patch kv kubevirt -n kubevirt --type=json -p "[{ \"op\": \"add\", \"path\": \"/spec/imageTag\", \"value\": \"${RELEASE}\" }]"

