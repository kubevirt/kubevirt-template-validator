#!/bin/bash
{
RET=1
$KUBECTL create -n default -f manifests/template-with-rules.yaml || exit 2
sleep 1s
$KUBECTL create -f manifests/11-vm-with-ui-style-annotations.yaml
if $KUBECTL get vm vm-test-11; then
	RET=0
	$KUBECTL delete vm vm-test-11
fi
$KUBECTL delete -n default -f manifests/template-with-rules.yaml
exit $RET
}
