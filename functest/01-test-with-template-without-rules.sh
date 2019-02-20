#!/bin/bash
RET=1
$KUBECTL create -f manifests/01-vm-from-template-without-rules.yaml &> /dev/null
if $KUBECTL get vm vm-test-01 &> /dev/null; then
	RET=0
	$KUBECTL delete -f manifests/01-vm-from-template-without-rules.yaml &> /dev/null
fi
exit $RET	
