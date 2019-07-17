#!/bin/bash
{
RET=0
$KUBECTL apply -n default -f manifests/template-with-rules-incorrect-just-warning.yaml  || exit 2
sleep 1s
if $KUBECTL apply -f manifests/09-vm-from-template-with-incorrect-rules-just-warning.yaml ; then
	if [ $($KUBECTL logs -n kubevirt $($KUBECTL get pods --all-namespaces | grep virt-template-validator | awk -F ' ' '{print $2}') | grep "warning.*Memory size not within range:"| wc -l) -eq 0 ]; then
		RET=1
	fi

	$KUBECTL delete vm vm-test-09
fi
$KUBECTL delete -n default -f manifests/template-with-rules-incorrect-just-warning.yaml
exit $RET
}
