#!/bin/bash
RET=1
$KUBECTL create -f manifests/01-vm-without-template.yaml &> /dev/null
if $KUBECTL get vm vm-test-01 &> /dev/null; then
	RET=0
	$KUBECTL delete vm vm-test-01 &> /dev/null
fi
exit $RET	
