#!/bin/bash
set -ex

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )

(
  kubectl wait --timeout=240s --for=condition=Ready -n kubevirt kv/kubevirt ;
) || {
  echo "Something went wrong"
  kubectl describe -n kubevirt kv/kubevirt
  kubectl describe pods -n kubevirt
  exit 1
}

# Give kvm some time to be announced
sleep 12

kubectl describe nodes
