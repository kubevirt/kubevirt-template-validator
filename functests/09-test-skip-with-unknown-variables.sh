#!/bin/bash
{
RET=1
$KUBECTL create -f manifests/09-vm-with-extra-fields.yaml
if $KUBECTL get vm vm-test-09; then
	RET=0
	$KUBECTL delete vm vm-test-09
fi
exit $RET
}
