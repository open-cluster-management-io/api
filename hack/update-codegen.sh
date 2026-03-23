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

OPENAPI_TMPDIR=$(mktemp -d)
trap "rm -rf ${OPENAPI_TMPDIR}" EXIT

for pkg in \
  cluster/v1 \
  cluster/v1alpha1 \
  cluster/v1beta1 \
  cluster/v1beta2 \
  work/v1 \
  work/v1alpha1 \
  operator/v1 \
  addon/v1alpha1 \
  addon/v1beta1; do
  go run -mod=vendor k8s.io/kube-openapi/cmd/openapi-gen \
    --output-file zz_generated.openapi.go \
    --output-model-name-file zz_generated.model_name.go \
    --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.txt" \
    --output-dir "${OPENAPI_TMPDIR}" \
    --output-pkg "discard" \
    --report-filename /dev/null \
    "open-cluster-management.io/api/${pkg}"
done
