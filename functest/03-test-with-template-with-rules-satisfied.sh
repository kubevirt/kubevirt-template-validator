#!/bin/bash
{
RET=1
$KUBECTL create -f manifests/template-with-rules.yaml || exit 2
$KUBECTL create -f manifests/03-vm-from-template-with-rules-satisfied.yaml
if $KUBECTL get vm vm-test-03; then
	RET=0
	$KUBECTL delete vm vm-test-03
fi
$KUBECTL delete -f manifests/template-with-rules.yaml
exit $RET
} &> /dev/null
