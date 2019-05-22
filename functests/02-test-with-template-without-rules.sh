#!/bin/bash
{
RET=1
$KUBECTL create -f manifests/02-vm-from-template-without-rules.yaml
if $KUBECTL get vm vm-test-02; then
	RET=0
	$KUBECTL delete vm vm-test-02
fi
exit $RET
}
