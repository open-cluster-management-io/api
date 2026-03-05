# Copyright Contributors to the Open Cluster Management project
TEST_TMP :=/tmp

ENSURE_ENVTEST_SCRIPT := https://raw.githubusercontent.com/open-cluster-management-io/sdk-go/main/ci/envtest/ensure-envtest.sh

.PHONY: envtest-setup
envtest-setup:
	$(eval export KUBEBUILDER_ASSETS=$(shell curl -fsSL $(ENSURE_ENVTEST_SCRIPT) | bash))
	@echo "KUBEBUILDER_ASSETS=$(KUBEBUILDER_ASSETS)"

clean-integration-test:
	$(RM) '$(KB_TOOLS_ARCHIVE_PATH)'
	rm -rf $(TEST_TMP)/kubebuilder
	$(RM) ./integration.test
.PHONY: clean-integration-test

clean: clean-integration-test

test-integration: test-api-integration
.PHONY: test-integration

test-api-integration: envtest-setup
	go test -c ./test/integration/api 
	./api.test -ginkgo.slowSpecThreshold=15 -ginkgo.v -ginkgo.failFast
.PHONY: test-api-integration
