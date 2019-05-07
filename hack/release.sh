#!/bin/bash

set -e

PROJECT=kubevirt-template-validator

if [ -z "$1" ]; then
	echo "usage: $0 <tag>"
	exit 1
fi
if [ -z "${GITHUB_TOKEN}" ] || [ -z "${GITHUB_USER}" ]; then
	echo "make sure to set GITHUB_TOKEN and GITHUB_USER env vars"
	exit 2
fi

TAG="$1"  #TODO: validate tag is vX.Y.Z
BRANCH=$(git rev-parse --abbrev-ref HEAD)

./hack/build/build.sh ${TAG}
if [ -d _out ]; then
	rm -rf _out;
fi
mkdir -p _out
cp cmd/${PROJECT}/${PROJECT} _out/${PROJECT}-${TAG}-linux-amd64
git add cmd/${PROJECT}/${PROJECT} && git ci -s -m "binaries: rebuild for tag ${TAG}"
git tag -a -m "${PROJECT} ${TAG}" ${TAG}
git push origin --tags ${BRANCH}
if  which github-release 2> /dev/null; then
	github-release release -t ${TAG} -r ${PROJECT}
	github-release upload -t ${TAG} -r ${PROJECT} \
		-n ${PROJECT}-${TAG}-linux-amd64 \
		-f _out/${PROJECT}-${TAG}-linux-amd64
fi
