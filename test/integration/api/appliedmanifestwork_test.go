// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	workv1 "open-cluster-management.io/api/work/v1"
)

var _ = ginkgo.Describe("AppliedManifestWork v1 API test", func() {
	var appliedManifestWorkName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		appliedManifestWorkName = fmt.Sprintf("appliedmanifestwork-%s", suffix)
	})

	ginkgo.AfterEach(func() {
		err := hubWorkClient.WorkV1().AppliedManifestWorks().Delete(context.TODO(), appliedManifestWorkName, metav1.DeleteOptions{})
		if !errors.IsNotFound(err) {
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		}
	})

	ginkgo.Context("AppliedManifestWork creation and validation", func() {
		ginkgo.It("should create AppliedManifestWork with basic spec", func() {
			appliedWork := &workv1.AppliedManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: appliedManifestWorkName,
				},
				Spec: workv1.AppliedManifestWorkSpec{
					HubHash:          "test-hub-hash",
					AgentID:          "test-agent",
					ManifestWorkName: "test-manifestwork",
				},
			}

			_, err := hubWorkClient.WorkV1().AppliedManifestWorks().Create(context.TODO(), appliedWork, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("should handle AppliedManifestWork with applied resources", func() {
			appliedWork := &workv1.AppliedManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: appliedManifestWorkName,
				},
				Spec: workv1.AppliedManifestWorkSpec{
					HubHash:          "test-hub-hash",
					AgentID:          "test-agent",
					ManifestWorkName: "test-manifestwork",
				},
			}

			_, err := hubWorkClient.WorkV1().AppliedManifestWorks().Create(context.TODO(), appliedWork, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("AppliedManifestWork status validation", func() {
		ginkgo.It("should allow status updates with applied resource status", func() {
			appliedWork := &workv1.AppliedManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: appliedManifestWorkName,
				},
				Spec: workv1.AppliedManifestWorkSpec{
					HubHash:          "test-hub-hash",
					AgentID:          "test-agent",
					ManifestWorkName: "test-manifestwork",
				},
			}

			appliedManifestWork, err := hubWorkClient.WorkV1().AppliedManifestWorks().Create(context.TODO(), appliedWork, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			// Update status
			appliedManifestWork.Status = workv1.AppliedManifestWorkStatus{
				AppliedResources: []workv1.AppliedManifestResourceMeta{
					{
						ResourceIdentifier: workv1.ResourceIdentifier{
							Group:     "",
							Resource:  "configmaps",
							Name:      "test-configmap",
							Namespace: "default",
						},
						Version: "v1",
						UID:     "test-uid-123",
					},
				},
			}

			_, err = hubWorkClient.WorkV1().AppliedManifestWorks().UpdateStatus(context.TODO(), appliedManifestWork, metav1.UpdateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("should handle complex status with multiple resources", func() {
			appliedWork := &workv1.AppliedManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: appliedManifestWorkName,
				},
				Spec: workv1.AppliedManifestWorkSpec{
					HubHash:          "test-hub-hash",
					AgentID:          "test-agent",
					ManifestWorkName: "test-manifestwork",
				},
			}

			appliedManifestWork, err := hubWorkClient.WorkV1().AppliedManifestWorks().Create(context.TODO(), appliedWork, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			// Update with complex status
			appliedManifestWork.Status = workv1.AppliedManifestWorkStatus{
				AppliedResources: []workv1.AppliedManifestResourceMeta{
					{
						ResourceIdentifier: workv1.ResourceIdentifier{
							Group:     "",
							Resource:  "configmaps",
							Name:      "test-configmap",
							Namespace: "default",
						},
						Version: "v1",
						UID:     "configmap-uid-123",
					},
					{
						ResourceIdentifier: workv1.ResourceIdentifier{
							Group:     "apps",
							Resource:  "deployments",
							Name:      "test-deployment",
							Namespace: "default",
						},
						Version: "v1",
						UID:     "deployment-uid-456",
					},
				},
			}

			updatedWork, err := hubWorkClient.WorkV1().AppliedManifestWorks().UpdateStatus(context.TODO(), appliedManifestWork, metav1.UpdateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(len(updatedWork.Status.AppliedResources)).Should(gomega.Equal(2))
		})
	})

	ginkgo.Context("AppliedManifestWork validation edge cases", func() {
		ginkgo.It("should create with required fields", func() {
			appliedWork := &workv1.AppliedManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: appliedManifestWorkName,
				},
				Spec: workv1.AppliedManifestWorkSpec{
					HubHash:          "test-hub-hash",
					AgentID:          "test-agent",
					ManifestWorkName: "test-manifestwork",
				},
			}

			createdAppliedWork, err := hubWorkClient.WorkV1().AppliedManifestWorks().Create(context.TODO(), appliedWork, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(createdAppliedWork.Spec.HubHash).Should(gomega.Equal("test-hub-hash"))
			gomega.Expect(createdAppliedWork.Spec.AgentID).Should(gomega.Equal("test-agent"))
		})

		ginkgo.It("should handle empty applied resources list", func() {
			appliedWork := &workv1.AppliedManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Name: appliedManifestWorkName,
				},
				Spec: workv1.AppliedManifestWorkSpec{
					HubHash:          "test-hub-hash",
					AgentID:          "test-agent",
					ManifestWorkName: "test-manifestwork",
				},
			}

			appliedManifestWork, err := hubWorkClient.WorkV1().AppliedManifestWorks().Create(context.TODO(), appliedWork, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(appliedManifestWork.Spec.HubHash).Should(gomega.Equal("test-hub-hash"))
		})
	})
})
