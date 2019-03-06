#!/bin/bash

TAG="${1:-devel}"

docker build -t fromanirh/kubevirt-template-validator:$TAG .
