#!/bin/bash
set -ex

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )

TARGET_VERSION="${1}"
CURRENT_VERSION="${2}"
if [ "$1" == "latest" ]; then
	TARGET_VERSION=$( curl --silent -k "https://api.github.com/repos/kubevirt/kubevirt/releases" |  ${BASEPATH}/latest-versions.py )
fi

if [ "${TARGET_VERSION}" == "${CURRENT_VERSION}" ]; then
	echo "nothing to do: target == current"
	exit 0
fi

bash -x ${BASEPATH}/upgrade-kubevirt.sh ${TARGET_VERSION}
bash -x ${BASEPATH}/wait-kubevirt.sh
# -e is important to let make use the vars defined in "matrix:" above
make -e $TARGET
