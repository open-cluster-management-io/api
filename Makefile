all: build
.PHONY: all

# Include the library makefile
include $(addprefix ./vendor/github.com/openshift/build-machinery-go/make/, \
	golang.mk \
	targets/openshift/deps.mk \
	targets/openshift/crd-schema-gen.mk \
)

GO_PACKAGES :=$(addsuffix ...,$(addprefix ./,$(filter-out vendor/,$(filter-out hack/,$(wildcard */)))))
GO_BUILD_PACKAGES :=$(GO_PACKAGES)
GO_BUILD_PACKAGES_EXPANDED :=$(GO_BUILD_PACKAGES)
# LDFLAGS are not needed for dummy builds (saving time on calling git commands)
GO_LD_FLAGS:=
CONTROLLER_GEN_VERSION :=v0.2.5

# $1 - target name
# $2 - apis
# $3 - manifests
# $4 - output
$(call add-crd-gen,clusterv1,./cluster/v1,./cluster/v1,./cluster/v1)
$(call add-crd-gen,clusterv1alpha1,./cluster/v1alpha1,./cluster/v1alpha1,./cluster/v1alpha1)
$(call add-crd-gen,work,./work/v1,./work/v1,./work/v1)
$(call add-crd-gen,operator,./operator/v1,./operator/v1,./operator/v1)
$(call add-crd-gen,addonv1alpha1,./addon/v1alpha1,./addon/v1alpha1,./addon/v1alpha1)

RUNTIME ?= podman
RUNTIME_IMAGE_NAME ?= openshift-api-generator

verify-scripts:
	bash -x hack/verify-deepcopy.sh
	bash -x hack/verify-swagger-docs.sh
	bash -x hack/verify-crds.sh
	bash -x hack/verify-codegen.sh
.PHONY: verify-scripts
verify: verify-scripts verify-codegen-crds

update-scripts:
	hack/update-deepcopy.sh
	hack/update-swagger-docs.sh
	hack/update-codegen.sh
.PHONY: update-scripts
update: update-scripts update-codegen-crds

generate-with-container: Dockerfile.build
	$(RUNTIME) build -t $(RUNTIME_IMAGE_NAME) -f Dockerfile.build .
	$(RUNTIME) run -ti --rm -v $(PWD):/go/src/github.com/open-cluster-management/api:z -w /go/src/github.com/open-cluster-management/api $(RUNTIME_IMAGE_NAME) make update-scripts
