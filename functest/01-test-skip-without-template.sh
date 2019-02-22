#!/bin/bash
{
RET=1
$KUBECTL create -f manifests/01-vm-without-template.yaml
if $KUBECTL get vm vm-test-01; then
	RET=0
	$KUBECTL delete vm vm-test-01
fi
exit $RET
} &> /dev/null
