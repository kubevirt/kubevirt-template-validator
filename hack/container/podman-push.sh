#!/bin/bash
TAGS="$*"
for TAG in $TAGS; do
	podman push kubevirt/kubevirt-template-validator:$TAG docker://quay.io/kubevirt/kubevirt-template-validator:$TAG
done
