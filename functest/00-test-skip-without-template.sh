#!/bin/bash
RET=1
$KUBECTL create -f manifests/00-vm-without-template.yaml &> /dev/null
if $KUBECTL get vm vm-test-00 &> /dev/null; then
	RET=0
	$KUBECTL delete -f manifests/00-vm-without-template.yaml &> /dev/null
fi
exit $RET	
