#!/bin/bash
{
RET=1
$KUBECTL apply -n default -f manifests/template-with-rules-incorrect-just-warning.yaml  || exit 2
sleep 1s
if $KUBECTL apply -f manifests/09-vm-from-template-with-incorrect-rules-just-warning.yaml ; then
		while read -r pod ; do
			#$KUBECTL logs -n kubevirt $pod
			if [ $($KUBECTL logs -n kubevirt $pod | grep "warning.*Memory size not within range:"| wc -l) -gt 0 ]; then
				RET=0
			fi
		done <<< "$($KUBECTL get pods --all-namespaces | grep virt-template-validator | awk -F ' ' '{print $2}')"
	$KUBECTL delete vm vm-test-09
fi

$KUBECTL delete -n default -f manifests/template-with-rules-incorrect-just-warning.yaml
exit $RET
}
