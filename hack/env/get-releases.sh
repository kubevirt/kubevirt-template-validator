#!/bin/bash

set -ex

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )

BUILTIN_VERSION=$( awk ' /.*kubevirt.io\/client.go/ { print $2 }' < go.mod )

curl --silent -k "https://api.github.com/repos/kubevirt/kubevirt/releases" > releases.json
${BASEPATH}/find-versions.py "${BUILTIN_VERSION}" < releases.json > versionsrc
