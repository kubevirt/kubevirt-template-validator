#!/bin/bash
set -ex

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )

# we can't guess.
[ -z "$1" ] && exit 1

oc patch kv kubevirt -n kubevirt --type=json -p "[{ \"op\": \"add\", \"path\": \"/spec/imageTag\", \"value\": \"$1\" }]"

