#!/bin/bash
set -e

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )
ENVPATH="${BASEPATH}/../env"

${ENVPATH}/minishift/setup.sh
sleep 40s

${ENVPATH}/try-login.sh
sleep 10s

${ENVPATH}/install-cluster.sh

export KUBECTL=oc
export OC=oc