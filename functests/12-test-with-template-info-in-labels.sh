#!/bin/bash
{
RET=1
$KUBECTL create -n default -f manifests/template-with-rules.yaml || exit 2
sleep 1s
$KUBECTL create -f manifests/12-vm-with-template-info-in-labels.yaml
if $KUBECTL get vm vm-test-12; then
	RET=0
	$KUBECTL delete vm vm-test-12
fi
$KUBECTL delete -n default -f manifests/template-with-rules.yaml
exit $RET
}
