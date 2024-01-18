SHELL :=/bin/bash

all: build
.PHONY: all

# Include the library makefile
include $(addprefix ./vendor/github.com/openshift/build-machinery-go/make/, \
	golang.mk \
	targets/openshift/deps.mk \
	targets/openshift/crd-schema-gen.mk \
)

GO_PACKAGES :=$(addsuffix ...,$(addprefix ./,$(filter-out test/, $(filter-out vendor/,$(filter-out hack/,$(wildcard */))))))
GO_BUILD_PACKAGES :=$(GO_PACKAGES)
GO_BUILD_PACKAGES_EXPANDED :=$(GO_BUILD_PACKAGES)
# LDFLAGS are not needed for dummy builds (saving time on calling git commands)
GO_LD_FLAGS:=
# controller-gen setup
CONTROLLER_GEN_VERSION :=v0.11.3
CONTROLLER_GEN :=$(PERMANENT_TMP_GOPATH)/bin/controller-gen
ifneq "" "$(wildcard $(CONTROLLER_GEN))"
_controller_gen_installed_version = $(shell $(CONTROLLER_GEN) --version | awk '{print $$2}')
endif
controller_gen_dir :=$(abspath $(PERMANENT_TMP_GOPATH)/bin)

# $1 - target name
# $2 - apis
# $3 - manifests
# $4 - output
$(call add-crd-gen,clusterv1,./cluster/v1,./cluster/v1,./cluster/v1)
$(call add-crd-gen,clusterv1alpha1,./cluster/v1alpha1,./cluster/v1alpha1,./cluster/v1alpha1)
$(call add-crd-gen,clusterv1beta1,./cluster/v1beta1,./cluster/v1beta1,./cluster/v1beta1)
$(call add-crd-gen,clusterv1beta2,./cluster/v1beta1 ./cluster/v1beta2,./cluster/v1beta2,./cluster/v1beta2)
$(call add-crd-gen,workv1,./work/v1,./work/v1,./work/v1)
$(call add-crd-gen,workv1alpha1,./work/v1alpha1,./work/v1alpha1,./work/v1alpha1)
$(call add-crd-gen,operator,./operator/v1,./operator/v1,./operator/v1)
$(call add-crd-gen,addonv1alpha1,./addon/v1alpha1,./addon/v1alpha1,./addon/v1alpha1)

RUNTIME ?= podman
RUNTIME_IMAGE_NAME ?= openshift-api-generator

verify-gocilint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.0
	go vet ./...
	${GOPATH}/bin/golangci-lint run --timeout=3m ./...

verify-scripts:
	bash -x hack/verify-deepcopy.sh
	bash -x hack/verify-swagger-docs.sh
	bash -x hack/verify-crds.sh
	bash -x hack/verify-codegen.sh
.PHONY: verify-scripts
verify: check-env verify-scripts verify-codegen-crds verify-gocilint

update-scripts:
	hack/update-deepcopy.sh
	hack/update-swagger-docs.sh
	hack/update-codegen.sh
	hack/update-v1beta1-crds.sh
.PHONY: update-scripts
update: check-env update-scripts update-codegen-crds

build-runtime-image: Dockerfile.build
	$(RUNTIME) build -t $(RUNTIME_IMAGE_NAME) -f Dockerfile.build .

update-with-container: build-runtime-image
	$(RUNTIME) run -ti --rm -v $(PWD):/go/src/open-cluster-management.io/api:z -w /go/src/open-cluster-management.io/api $(RUNTIME_IMAGE_NAME) make update-scripts update-codegen-crds

include ./test/integration-test.mk

check-env:
ifeq ($(GOPATH),)
	$(warning "environment variable GOPATH is empty, auto set from go env GOPATH")
export GOPATH=$(shell go env GOPATH)
endif
.PHONY: check-env

# override ensure-controller-gen
ensure-controller-gen:
ifeq "" "$(wildcard $(CONTROLLER_GEN))"
	$(info Installing controller-gen into '$(CONTROLLER_GEN)')
	mkdir -p '$(controller_gen_dir)'
	GOBIN='$(controller_gen_dir)' go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)
	chmod +x '$(CONTROLLER_GEN)';
else
	$(info Using existing controller-gen from "$(CONTROLLER_GEN)")
	@[[ "$(_controller_gen_installed_version)" == $(CONTROLLER_GEN_VERSION) ]] || \
	echo "Warning: Installed controller-gen version $(_controller_gen_installed_version) does not match expected version $(CONTROLLER_GEN_VERSION)."
endif
