#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../../../k8s.io/code-generator)}

verify="${VERIFY:-}"

set -x
# Because go mod sux, we have to fake the vendor for generator in order to be able to build it...
mv ${CODEGEN_PKG}/generate-groups.sh ${CODEGEN_PKG}/generate-groups.sh.orig
sed 's/ go install/#go install/g' ${CODEGEN_PKG}/generate-groups.sh.orig > ${CODEGEN_PKG}/generate-groups.sh
function cleanup {
  mv ${CODEGEN_PKG}/generate-groups.sh.orig ${CODEGEN_PKG}/generate-groups.sh
}
trap cleanup EXIT

go install -mod=vendor ./${CODEGEN_PKG}/cmd/{defaulter-gen,client-gen,lister-gen,informer-gen,deepcopy-gen}

for group in cluster; do
  bash ${CODEGEN_PKG}/generate-groups.sh "client,lister,informer" \
    github.com/open-cluster-management/api/client/${group} \
    github.com/open-cluster-management/api \
    "${group}:v1,v1alpha1" \
    --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.txt \
    ${verify}
done

for group in work operator; do
  bash ${CODEGEN_PKG}/generate-groups.sh "client,lister,informer" \
    github.com/open-cluster-management/api/client/${group} \
    github.com/open-cluster-management/api \
    "${group}:v1" \
    --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.txt \
    ${verify}
done

for group in addon; do
  bash ${CODEGEN_PKG}/generate-groups.sh "client,lister,informer" \
    github.com/open-cluster-management/api/client/${group} \
    github.com/open-cluster-management/api \
    "${group}:v1alpha1" \
    --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.txt \
    ${verify}
done

