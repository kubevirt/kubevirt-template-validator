#!/bin/bash

TAG="${1:-devel}"
XTAGS="${@:2}"  # see https://stackoverflow.com/questions/9057387/process-all-arguments-except-the-first-one-in-a-bash-script

buildah bud -t fromani/kubevirt-template-validator:$TAG .
for XTAG in $XTAGS; do
	buildah tag fromani/kubevirt-template-validator:$TAG fromani/kubevirt-template-validator:$XTAG
done
