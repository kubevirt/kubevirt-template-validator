#!/bin/bash
RET=1
$KUBECTL create -f manifests/template-with-optional-rules.yaml &> /dev/null || exit 2
$KUBECTL create -f manifests/05-vm-from-template-with-optional-rules-unfulfilled.yaml &> /dev/null
if $KUBECTL get vm vm-test-05 &> /dev/null; then
	RET=0
	$KUBECTL delete vm vm-test-05 &> /dev/null
fi
$KUBECTL delete -f manifests/template-with-optional-rules.yaml &> /dev/null
exit $RET	
