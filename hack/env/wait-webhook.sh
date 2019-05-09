#!/bin/bash

OC=oc

while true; do
	VALIDATOR_POD=$(${OC} get pod -n kubevirt -l kubevirt.io=virt-template-validator -o json | jq -r .items[0].metadata.name)
	VALIDATOR_STATUS=$(${OC} get pod -n kubevirt ${VALIDATOR_POD} -o json | jq -r .status.containerStatuses[0].ready)
	if [ "${VALIDATOR_STATUS}" = "true" ]; then
			exit 0
	fi
	echo "${VALIDATOR_POD} ready=${VALIDATOR_STATUS}"
	sleep 3
done
