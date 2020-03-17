#!/bin/bash
{
RET=0
echo '[test_id:2960] Negative test - Start a VM with memory restrictions violation'
$KUBECTL create -n default -f manifests/template-with-rules.yaml || exit 2
sleep 1s
if $KUBECTL create -f manifests/04-vm-from-template-with-rules-unfulfilled.yaml ; then
	RET=1
	$KUBECTL delete vm vm-test-04
fi
$KUBECTL delete -n default -f manifests/template-with-rules.yaml
exit $RET
}
