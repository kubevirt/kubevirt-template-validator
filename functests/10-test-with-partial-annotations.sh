#!/bin/bash
{
RET=1
$KUBECTL create -n default -f manifests/template-with-rules-incorrect.yaml || exit 2
sleep 1s
$KUBECTL create -f manifests/10-vm-with-partial-annotations.yaml
if $KUBECTL get vm vm-test-10; then
	RET=0
	$KUBECTL delete vm vm-test-10
fi
$KUBECTL delete -n default -f manifests/template-with-rules-incorrect.yaml
exit $RET
}
