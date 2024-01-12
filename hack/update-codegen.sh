#!/bin/bash

# the path of this repo should be go/src/open-cluster-management.io/api which is set in kube-codegen.

source "$(dirname "${BASH_SOURCE}")/lib/init.sh"

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../../../k8s.io/code-generator)}

verify="${VERIFY:-}"

# HACK: For some reason this script is not executable.
${SED_CMD} -i 's,^exec \(.*/generate-internal-groups.sh\),bash \1,g' ${CODEGEN_PKG}/generate-groups.sh
# Because go mod sux, we have to fake the vendor for generator in order to be able to build it...
${SED_CMD} -i 's/GO111MODULE=on go install/#GO111MODULE=on go install/g' ${CODEGEN_PKG}/generate-internal-groups.sh

# For verification we need to ensure we don't remove files
# TODO: this should be properly resolved upstream so that we can get
# rid of the below if condition for verify scripts
if [ ! -z "$verify" ]; then
  ${SED_CMD} -i 's/xargs \-0 rm \-f/xargs -0 echo ""/g' ${CODEGEN_PKG}/generate-internal-groups.sh
fi

# ...but we have to put it back, or `verify` will puke.
trap "git checkout ${CODEGEN_PKG}" EXIT

go install -mod=vendor ./${CODEGEN_PKG}/cmd/{defaulter-gen,client-gen,lister-gen,informer-gen,deepcopy-gen}

for group in cluster; do
  bash ${CODEGEN_PKG}/generate-groups.sh "client,lister,informer" \
    open-cluster-management.io/api/client/${group} \
    open-cluster-management.io/api \
    "${group}:v1,v1alpha1,v1beta1,v1beta2" \
    --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.txt \
    ${verify}
done

for group in work; do
  bash ${CODEGEN_PKG}/generate-groups.sh "client,lister,informer" \
    open-cluster-management.io/api/client/${group} \
    open-cluster-management.io/api \
    "${group}:v1,v1alpha1" \
    --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.txt \
    ${verify}
done

for group in operator; do
  bash ${CODEGEN_PKG}/generate-groups.sh "client,lister,informer" \
    open-cluster-management.io/api/client/${group} \
    open-cluster-management.io/api \
    "${group}:v1" \
    --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.txt \
    ${verify}
done

for group in addon; do
  bash ${CODEGEN_PKG}/generate-groups.sh "client,lister,informer" \
    open-cluster-management.io/api/client/${group} \
    open-cluster-management.io/api \
    "${group}:v1alpha1" \
    --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.txt \
    ${verify}
done
