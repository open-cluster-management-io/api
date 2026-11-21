#!/bin/bash
# Copyright Contributors to the Open Cluster Management project

# shellcheck source=hack/lib/init.sh
source "$(dirname "${BASH_SOURCE[0]}")/lib/init.sh"

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
VERIFY=--verify-only "${SCRIPT_ROOT}/hack/update-deepcopy.sh"
