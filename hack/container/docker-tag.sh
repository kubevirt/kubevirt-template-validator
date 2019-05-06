#!/bin/bash

TAG="${1:-devel}"

docker build -t quay.io/fromani/kubevirt-template-validator:$TAG .
