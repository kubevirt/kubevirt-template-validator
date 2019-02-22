#!/bin/bash
RET=1
$KUBECTL create -f manifests/template-with-rules.yaml &> /dev/null || exit 2
$KUBECTL create -f manifests/06-vm-from-template-with-rules-and-unspecified-paths.yaml &> /dev/null
if $KUBECTL get vm vm-test-06 &> /dev/null; then
	RET=0
	$KUBECTL delete vm vm-test-06 &> /dev/null
fi
$KUBECTL delete -f manifests/template-with-rules.yaml &> /dev/null
exit $RET	
