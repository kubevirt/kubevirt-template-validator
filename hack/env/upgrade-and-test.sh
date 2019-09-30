#!/bin/bash
set -ex

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )

STEP="${1}"

CURRENT_VERSION="unknown"
TARGET_VERSION="unknown"

PREVIOUS_VERSION=$( awk -F= '/^previous=/ { print $2 }' versionsrc )
BUILTIN_VERSION=$( awk -F= '/^builtin=/ { print $2 }' versionsrc )
LAST_VERSION=$( awk -F= '/^last=/ { print $2 }' versionsrc )

case $STEP in
	old-to-builtin)
		CURRENT_VERSION="$PREVIOUS_VERSION"
		TARGET_VERSION="$BUILTIN_VERSION"
		;;
	builtin-to-new)
		CURRENT_VERSION="$BUILTIN_VERSION"
		TARGET_VERSION="$LAST_VERSION"
		;;
	*)
		exit 1
		;;
esac


if [ "${TARGET_VERSION}" == "${CURRENT_VERSION}" ]; then
	echo "nothing to do: target == current"
	exit 0
fi

echo "${CURRENT_VERSION} -> ${TARGET_VERSION}"

bash -x ${BASEPATH}/upgrade-kubevirt.sh ${TARGET_VERSION}
bash -x ${BASEPATH}/wait-kubevirt.sh
# -e is important to let make use the vars defined in "matrix:" above
make -e $TARGET
