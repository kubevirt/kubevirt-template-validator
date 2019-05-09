#!/bin/bash
{
RET=0
$KUBECTL create -n default -f manifests/template-with-rules.yaml  || exit 2
# will fail because the default value if not specified is the zero value, so this VM will have zero cores (!)
if $KUBECTL create -f manifests/06-vm-from-template-with-rules-and-unspecified-paths.yaml; then
	RET=1
	$KUBECTL delete vm vm-test-06
fi
$KUBECTL delete -n default -f manifests/template-with-rules.yaml
exit $RET
}
