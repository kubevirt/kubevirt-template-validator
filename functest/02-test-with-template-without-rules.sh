#!/bin/bash
RET=1
$KUBECTL create -f manifests/02-vm-from-template-without-rules.yaml &> /dev/null
if $KUBECTL get vm vm-test-02 &> /dev/null; then
	RET=0
	$KUBECTL delete vm vm-test-02 &> /dev/null
fi
exit $RET	
