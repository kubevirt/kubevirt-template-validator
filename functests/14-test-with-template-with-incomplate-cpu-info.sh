#!/bin/bash
{
RET=1
$KUBECTL create -n default -f manifests/template-with-rules.yaml || exit 2
sleep 1s
$KUBECTL create -f manifests/14-vm-from-template-with-incomplete-cpu-info.yaml
if $KUBECTL get vm vm-test-14; then
	RET=0
	$KUBECTL delete vm vm-test-14
fi
$KUBECTL delete -n default -f manifests/template-with-rules.yaml
exit $RET
}
