#!/bin/bash

TAG="${1:-devel}"

docker build -t quay.io/fromani/kubevirt-template-validator:$TAG . && \
docker push quay.io/fromani/kubevirt-template-validator:$TAG
