#!/bin/bash

set -e

if [ -z "$1" ]; then
	echo "usage: $0 <tag>"
	exit 1
fi
if [ -z "${GITHUB_TOKEN}" ] || [ -z "${GITHUB_USER}" ]; then
	echo "make sure to set GITHUB_TOKEN and GITHUB_USER env vars"
	exit 2
fi

TAG="$1"  #TODO: validate tag is vX.Y.Z

./hack/build/build.sh ${TAG}
if [ -d _out ]; then
	rm -rf _out;
fi
mkdir -p _out
cp cmd/kubevirt-template-validator/kubevirt-template-validator _out/kubevirt-template-validator-${TAG}-linux-amd64
git add cmd/kubevirt-template-validator/kubevirt-template-validator && git ci -s -m "binary: rebuild for tag ${TAG}"
git tag -a -m "kubevirt-template-validator ${TAG}" ${TAG}
if  which github-release 2> /dev/null; then
	github-release release -t ${TAG}
	github-release upload -t ${TAG} \
		-n kubevirt-template-validator-${TAG}-linux-amd64 \
		-f _out/kubevirt-template-validator-${TAG}-linux-amd64
fi
