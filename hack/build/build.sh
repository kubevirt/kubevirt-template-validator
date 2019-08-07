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

export GO111MODULE=on
export GOPROXY=off
export GOFLAGS=-mod=vendor
cd cmd/kubevirt-template-validator && go build -v .
