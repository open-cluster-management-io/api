#!/bin/bash

# the path of this repo should be go/src/open-cluster-management.io/api which is set in kube-codegen.

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../../../k8s.io/code-generator)}

source "${CODEGEN_PKG}/kube_codegen.sh"

for group in operator work addon; do
kube::codegen::gen_helpers \
    --input-pkg-root open-cluster-management.io/api/${group} \
    --output-base "${SCRIPT_ROOT}/../.." \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.txt"
done

# skip cluster/v1alpha1 since failed to handle ClusterRolloutStatusFunc RolloutHandler in helper.go
# TODO: will added back until helper.go is removed.
for version in v1 v1beta1 v1beta2; do
kube::codegen::gen_helpers \
    --input-pkg-root open-cluster-management.io/api/cluster/${version} \
    --output-base "${SCRIPT_ROOT}/../.." \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.txt"
done
