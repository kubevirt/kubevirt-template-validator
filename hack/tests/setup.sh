#!/bin/bash
set -e

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )
ENVPATH="${BASEPATH}/../env"

${ENVPATH}/minishift/setup.sh
sleep 30s # to cool down
${ENVPATH}/try-login.sh

${ENVPATH}/install-cluster.sh

export KUBECTL=oc
export OC=oc
