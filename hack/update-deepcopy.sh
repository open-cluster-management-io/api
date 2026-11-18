#!/bin/bash
# Copyright Contributors to the Open Cluster Management project

# the path of this repo should be go/src/open-cluster-management.io/api which is set in kube-codegen.

# shellcheck source=hack/lib/init.sh
source "$(dirname "${BASH_SOURCE[0]}")/lib/init.sh"

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(
    cd "${SCRIPT_ROOT}"
    ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../../../k8s.io/code-generator
)}

source "${CODEGEN_PKG}/kube_codegen.sh"

# List of API groups
groups=("cluster" "operator" "work" "addon")

# Loop over each group
for group in "${groups[@]}"; do
    echo "Generating code for group: ${group}"

    # 1️⃣ Generate conversion functions
    kube::codegen::gen_helpers \
        --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.txt" \
        "${group}"

    # 2️⃣ Generate register (AddToScheme) functions
    kube::codegen::gen_register \
        --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.txt" \
        "${group}"
done
