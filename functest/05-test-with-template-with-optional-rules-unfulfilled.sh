#!/bin/bash
{
RET=1
$KUBECTL create -n default -f manifests/template-with-rules-optional.yaml || exit 2
$KUBECTL create -f manifests/05-vm-from-template-with-optional-rules-unfulfilled.yaml
if $KUBECTL get vm vm-test-05; then
	RET=0
	$KUBECTL delete vm vm-test-05
fi
$KUBECTL delete -n default -f manifests/template-with-rules-optional.yaml
exit $RET
}
