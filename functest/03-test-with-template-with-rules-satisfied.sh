#!/bin/bash
RET=1
$KUBECTL create -f manifests/template-with-rules.yaml &> /dev/null || exit 2
$KUBECTL create -f manifests/03-vm-from-template-with-rules-satisfied.yaml &> /dev/null
if $KUBECTL get vm vm-test-03 &> /dev/null; then
	RET=0
	$KUBECTL delete vm vm-test-03 &> /dev/null
fi
$KUBECTL delete -f manifests/template-with-rules.yaml &> /dev/null
exit $RET	
