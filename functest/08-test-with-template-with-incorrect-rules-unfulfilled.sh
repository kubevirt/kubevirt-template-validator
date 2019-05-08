#!/bin/bash
{
RET=0
$KUBECTL create -n default -f manifests/template-with-rules-incorrect.yaml  || exit 2
if $KUBECTL create -f manifests/08-vm-from-template-with-incorrect-rules-unfulfilled.yaml ; then
	RET=1
	$KUBECTL delete vm vm-test-08
fi
$KUBECTL delete -n default -f manifests/template-with-rules-incorrect.yaml
exit $RET
}
