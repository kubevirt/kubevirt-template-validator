#!/bin/bash
{
RET=0
$KUBECTL create -f manifests/template-with-rules.yaml || exit 2
if $KUBECTL create -f manifests/04-vm-from-template-with-rules-unfulfilled.yaml ; then
	RET=1
	$KUBECTL delete vm vm-test-04
fi
$KUBECTL delete -f manifests/template-with-rules.yaml
exit $RET
} &> /dev/null
