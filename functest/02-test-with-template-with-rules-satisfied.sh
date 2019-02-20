#!/bin/bash
RET=1
$KUBECTL create -f manifests/template-with-rules.yaml &> /dev/null || exit 2
$KUBECTL create -f manifests/02-vm-from-template-with-rules-satisfied.yaml &> /dev/null
if $KUBECTL get vm vm-test-02 &> /dev/null; then
	RET=0
	$KUBECTL delete -f manifests/02-vm-from-template-with-rules-satisfied.yaml &> /dev/null
fi
$KUBECTL delete -f manifests/template-with-rules.yaml &> /dev/null
exit $RET	
