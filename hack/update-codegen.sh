# Copyright Contributors to the Open Cluster Management project
#!/bin/bash

# the path of this repo should be go/src/open-cluster-management.io/api which is set in kube-codegen.

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../../../k8s.io/code-generator)}

verify="${VERIFY:-}"

source "${CODEGEN_PKG}/kube_codegen.sh"

kube::codegen::gen_client \
  --output-pkg "open-cluster-management.io/api/client/cluster" \
  --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.txt" \
  --output-dir ${SCRIPT_ROOT}/client/cluster \
  --one-input-api cluster \
  --with-watch \
  .

kube::codegen::gen_client \
  --output-pkg "open-cluster-management.io/api/client/work" \
  --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.txt" \
  --output-dir ${SCRIPT_ROOT}/client/work \
  --one-input-api work \
  --with-watch \
  .

kube::codegen::gen_client \
  --output-pkg "open-cluster-management.io/api/client/operator" \
  --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.txt" \
  --output-dir ${SCRIPT_ROOT}/client/operator \
  --one-input-api operator \
  --with-watch \
  .

kube::codegen::gen_client \
  --output-pkg "open-cluster-management.io/api/client/addon" \
  --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.txt" \
  --output-dir ${SCRIPT_ROOT}/client/addon \
  --one-input-api addon \
  --with-watch \
  .
