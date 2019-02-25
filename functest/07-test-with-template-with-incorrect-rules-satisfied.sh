#!/bin/bash
{
RET=0
$KUBECTL create -f manifests/template-with-rules-incorrect.yaml  || exit 2
if $KUBECTL create -f manifests/07-vm-from-template-with-incorrect-rules-satisfied.yaml ;  then
	RET=1
	$KUBECTL delete vm vm-test-07
fi
$KUBECTL delete -f manifests/template-with-rules-incorrect.yaml
exit $RET
} &> /dev/null
