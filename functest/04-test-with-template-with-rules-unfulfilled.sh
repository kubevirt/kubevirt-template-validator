#!/bin/bash
RET=0
$KUBECTL create -f manifests/template-with-rules.yaml &> /dev/null || exit 2
if $KUBECTL create -f manifests/04-vm-from-template-with-rules-unfulfilled.yaml &> /dev/null; then
	RET=1
	$KUBECTL delete vm vm-test-04 &> /dev/null
fi
$KUBECTL delete -f manifests/template-with-rules.yaml &> /dev/null
exit $RET	
