#!/bin/bash
# Copyright Contributors to the Open Cluster Management project

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
VERIFY=true "${SCRIPT_ROOT}/hack/update-copyright.sh"
