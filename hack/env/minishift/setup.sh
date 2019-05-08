#!/bin/bash
set -e

SELF=$( realpath $0 )
BASEPATH=$( dirname $SELF )

minishift start
# see https://github.com/minishift/minishift/pull/3044 for details
minishift addons install --defaults
minishift addons apply admissions-webhook
