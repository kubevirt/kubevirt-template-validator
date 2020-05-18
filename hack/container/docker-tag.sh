#!/bin/bash

TAG="${1:-devel}"

docker build -t quay.io/kubevirt/kubevirt-template-validator:$TAG .
