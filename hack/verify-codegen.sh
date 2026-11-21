#!/bin/bash
# Copyright Contributors to the Open Cluster Management project

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
VERIFY=--verify-only "${SCRIPT_ROOT}/hack/update-codegen.sh"
