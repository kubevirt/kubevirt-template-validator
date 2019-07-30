#!/bin/sh

set -e

if [ -z "$1" ]; then
	echo "usage: $0 <tag>"
	exit 1
fi

TAG="$1"  #TODO: validate tag is vX.Y.Z
VERSIONDIR="internal/pkg/version"
VERSIONFILE="${VERSIONDIR}/version.go"

mkdir -p ${VERSIONDIR} && ./hack/build/genver.sh ${TAG} > ${VERSIONFILE}
cd cmd/kubevirt-template-validator && GO111MODULE=on go build -v .
