#!/bin/bash
set -e

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )

eval $( minishift docker-env )

REGISTRY=$(minishift openshift registry)

docker --version

oc login -u developer
docker login -u $(oc whoami) -p $(oc whoami -t) ${REGISTRY}

docker images

make -C ${BASEPATH}/../.. binary VERSION=devel
docker build -f ${BASEPATH}/../../Dockerfile -t ${REGISTRY}/kubevirt/kubevirt-template-validator:devel ${BASEPATH}/../..
docker push ${REGISTRY}/kubevirt/kubevirt-template-validator:devel || :

sleep 5s

docker images
docker image inspect ${REGISTRY}/kubevirt/kubevirt-template-validator:devel

sleep 5s

oc login -u system:admin

kubectl get configmap -n kubevirt kubevirt-config -o yaml
