#!/bin/bash

# This script is meant to be the entrypoint for OpenShift Bash scripts to import all of the support
# libraries at once in order to make Bash script preambles as minimal as possible. This script recur-
# sively `source`s *.sh files in this directory tree. As such, no files should be `source`ed outside
# of this script to ensure that we do not attempt to overwrite read-only variables.

set -o errexit
set -o nounset
set -o pipefail

API_GROUP_VERSIONS="\
cluster/v1 \
work/v1 \
nucleus/v1 \
"

API_PACKAGES="\
github.com/open-cluster-management/api/cluster/v1,\
github.com/open-cluster-management/api/work/v1,\
github.com/open-cluster-management/api/nucleus/v1\
"
