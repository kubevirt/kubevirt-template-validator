#!/bin/bash
{
RET=1
$KUBECTL create -f manifests/template-with-rules-incorrect.yaml  || exit 2
$KUBECTL create -f manifests/07-vm-from-template-with-incorrect-rules-satisfied.yaml
if $KUBECTL get vm vm-test-07 ; then
	RET=0
	$KUBECTL delete vm vm-test-07
fi
$KUBECTL delete -f manifests/template-with-rules-incorrect.yaml
exit $RET
} &> /dev/null
