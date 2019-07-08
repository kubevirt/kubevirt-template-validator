#!/bin/bash
{
RET=0
$KUBECTL apply -n default -f manifests/template-with-rules-incorrect-just-warning.yaml  || exit 2
sleep 1s
if $KUBECTL apply -f manifests/08-vm-from-template-with-incorrect-rules-unfulfilled.yaml ; then
	RET=1
	$KUBECTL delete vm vm-test-08
fi
$KUBECTL delete -n default -f manifests/template-with-rules-incorrect-just-warning.yaml
exit $RET
}
