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
cluster/v1alpha1 \
work/v1 \
operator/v1 \
addon/v1alpha1 \
"

API_PACKAGES="\
github.com/open-cluster-management/api/cluster/v1,\
github.com/open-cluster-management/api/cluster/v1alpha1,\
github.com/open-cluster-management/api/work/v1,\
github.com/open-cluster-management/api/operator/v1,\
github.com/open-cluster-management/api/addon/v1alpha1\
"
