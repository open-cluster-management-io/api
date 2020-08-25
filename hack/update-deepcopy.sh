#!/bin/bash

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../../../k8s.io/code-generator)}

verify="${VERIFY:-}"

GOFLAGS="" bash ${CODEGEN_PKG}/generate-groups.sh "deepcopy" \
  github.com/open-cluster-management/api/generated \
  github.com/open-cluster-management/api \
  "cluster:v1 cluster:v1alpha1 work:v1 operator:v1 addon:v1alpha1" \
  --go-header-file ${SCRIPT_ROOT}/hack/empty.txt \
  ${verify}
