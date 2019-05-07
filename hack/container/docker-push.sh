#!/bin/bash

set -e

VERSION="${1:-devel}"

echo "$QUAY_BOT_PASS" | docker login -u="$QUAY_BOT_USER" --password-stdin quay.io
docker build -t quay.io/fromani/kubevirt-template-validator:$VERSION .
docker push quay.io/fromani/kubevirt-template-validator:$VERSION
