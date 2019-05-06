#!/bin/bash

set -e

if [ -z "$1" ]; then
	echo "usage: $0 <tag>"
	exit 1
fi

TAG="$1"  #TODO: validate tag is vX.Y.Z

./hack/build/build.sh ${TAG}
git add cmd/kubevirt-template-validator/kubevirt-template-validator && git ci -s -m \"binary: rebuild for tag ${TAG}\"
git tag -a -m \"kubevirt-template-validator ${TAG}\" ${TAG}
