#!/bin/bash

OC=oc
NAMESPACE=${1:-kubevirt}

if [ "${CI}" == "true" ] && [ "${TRAVIS}" == "true" ]; then
	sleep 5s
	${OC} get pods -n ${NAMESPACE}
	${OC} get pod -n ${NAMESPACE} -l kubevirt.io=virt-template-validator -o json
	sleep 5s
fi

for ix in $(seq 1 40); do
	VALIDATOR_POD_INFO=$(${OC} get pod -n ${NAMESPACE} -l kubevirt.io=virt-template-validator -o json )
	VALIDATOR_POD_STATUS=$( echo ${VALIDATOR_POD_INFO} | jq -r .items[0].status )
	VALIDATOR_STATUS=$( echo "${VALIDATOR_POD_INFO}" | jq -r .items[0].status.containerStatuses[0].ready)
	if [ "${VALIDATOR_STATUS}" = "true" ]; then
			exit 0
	fi
	echo "${VALIDATOR_POD_STATUS}"
	sleep 3
done
