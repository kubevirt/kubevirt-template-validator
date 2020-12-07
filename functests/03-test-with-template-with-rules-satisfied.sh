#!/bin/bash
{
RET=1
echo "[test_id:5033]: Template with validations, VM without validations"
$KUBECTL create -n default -f manifests/template-with-rules.yaml || exit 2
sleep 1s
$KUBECTL create -f manifests/03-vm-from-template-with-rules-satisfied.yaml
if $KUBECTL get vm vm-test-03; then
	RET=0
	$KUBECTL delete vm vm-test-03
fi
$KUBECTL delete -n default -f manifests/template-with-rules.yaml
exit $RET
}
