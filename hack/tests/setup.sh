#!/bin/bash
set -e

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )
ENVPATH="${BASEPATH}/../env"

${ENVPATH}/minishift/setup.sh
sleep 50s

${ENVPATH}/try-login.sh
sleep 10s

${ENVPATH}/install-cluster.sh
${ENVPATH}/install-validator.sh

export KUBECTL=oc
export OC=oc
