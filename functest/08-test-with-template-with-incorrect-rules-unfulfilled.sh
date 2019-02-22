#!/bin/bash
RET=1
$KUBECTL create -f manifests/template-with-incorrect-rules.yaml &> /dev/null || exit 2
$KUBECTL create -f manifests/08-vm-from-template-with-incorrect-rules-unfulfilled.yaml &> /dev/null
if $KUBECTL get vm vm-test-08 &> /dev/null; then
	RET=0
	$KUBECTL delete vm vm-test-08 &> /dev/null
fi
$KUBECTL delete -f manifests/template-with-incorrect-rules.yaml &> /dev/null
exit $RET	
