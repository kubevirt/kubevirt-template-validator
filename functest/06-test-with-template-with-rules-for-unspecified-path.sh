#!/bin/bash
{
#RET=1
#$KUBECTL create -f manifests/template-with-rules.yaml  || exit 2
# this is what we should do, and we don't because of
# https://github.com/fromanirh/kubevirt-template-validator/issues/2
#$KUBECTL create -f manifests/06-vm-from-template-with-rules-and-unspecified-paths.yaml
#if $KUBECTL get vm vm-test-06 ; then
#	RET=0
#	$KUBECTL delete vm vm-test-06
#fi
#
#$KUBECTL delete -f manifests/template-with-rules.yaml
#exit $RET
exit 99
} &> /dev/null
