#!/bin/bash
TAGS="$*"
for TAG in $TAGS; do
	buildah push fromani/kubevirt-template-validator:$TAG docker://quay.io/fromani/kubevirt-template-validator:$TAG
done
