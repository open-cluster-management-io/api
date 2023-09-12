package cloudevents

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	apitypes "k8s.io/apimachinery/pkg/types"

	"open-cluster-management.io/api/cloudevents/generic/types"
	"open-cluster-management.io/api/cloudevents/work/payload"
	"open-cluster-management.io/api/test/integration/cloudevents/agent"
	"open-cluster-management.io/api/test/integration/cloudevents/source"
	workv1 "open-cluster-management.io/api/work/v1"
)

var _ = ginkgo.Describe("Cloudevents clients test", func() {
	ginkgo.Context("Resync resources", func() {
		ginkgo.It("resync resources between source and agent", func() {
			ginkgo.By("start an agent on cluster1")
			clusterName := "cluster1"

			clientHolder, err := agent.StartWorkAgent(context.TODO(), clusterName, mqttOptions)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			informer := clientHolder.ManifestWorkInformer()
			lister := informer.Lister().ManifestWorks(clusterName)
			agentWorkClient := clientHolder.ManifestWorks(clusterName)

			gomega.Eventually(func() error {
				list, err := lister.List(labels.Everything())
				if err != nil {
					return err
				}

				// ensure there is only one work was synced on the cluster1
				if len(list) != 1 {
					return fmt.Errorf("unexpected work list %v", list)
				}

				// ensure the work can be get by work client
				workName := source.ResourceID(clusterName, "resource1")
				work, err := agentWorkClient.Get(context.TODO(), workName, metav1.GetOptions{})
				if err != nil {
					return err
				}

				newWork := work.DeepCopy()
				newWork.Status = workv1.ManifestWorkStatus{Conditions: []metav1.Condition{{Type: "Created", Status: metav1.ConditionTrue}}}

				// only update the status on the agent local part
				store := informer.Informer().GetStore()
				if err := store.Update(newWork); err != nil {
					return err
				}

				return nil
			}, 10*time.Second, 1*time.Second).Should(gomega.Succeed())

			ginkgo.By("resync the status from source")
			err = sourceCloudEventsClient.Resync(context.TODO())
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			gomega.Eventually(func() error {
				resourceID := source.ResourceID(clusterName, "resource1")
				resource, err := source.GetStore().Get(resourceID)
				if err != nil {
					return err
				}

				// ensure the resource status is synced
				if !meta.IsStatusConditionTrue(resource.Status.Conditions, "Created") {
					return fmt.Errorf("unexpected status %v", resource.Status.Conditions)
				}

				return nil
			}, 10*time.Second, 1*time.Second).Should(gomega.Succeed())
		})
	})

	ginkgo.Context("Publish a resource", func() {
		ginkgo.It("send a resource to a cluster", func() {
			ginkgo.By("start an agent on cluster2")
			clusterName := "cluster2"

			clientHolder, err := agent.StartWorkAgent(context.TODO(), clusterName, mqttOptions)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			lister := clientHolder.ManifestWorkInformer().Lister().ManifestWorks(clusterName)
			agentWorkClient := clientHolder.ManifestWorks(clusterName)

			gomega.Eventually(func() error {
				list, err := lister.List(labels.Everything())
				if err != nil {
					return err
				}

				// ensure there is only one work was synced on the cluster2
				if len(list) != 1 {
					return fmt.Errorf("unexpected work list %v", list)
				}

				// ensure the work can be get by work client
				workName := source.ResourceID(clusterName, "resource1")
				_, err = agentWorkClient.Get(context.TODO(), workName, metav1.GetOptions{})
				if err != nil {
					return err
				}

				return nil
			}, 10*time.Second, 1*time.Second).Should(gomega.Succeed())

			newResourceName := "resource2"
			ginkgo.By("create a new resource on the source and send it to the cluster2", func() {
				newResource := source.NewResource(clusterName, newResourceName)
				source.GetStore().Add(newResource)

				err := sourceCloudEventsClient.Publish(context.TODO(), types.CloudEventsType{
					CloudEventsDataType: payload.ManifestEventDataType,
					SubResource:         types.SubResourceSpec,
					Action:              "test_create_request",
				}, newResource)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})

			ginkgo.By("receive the new resource on the cluster2", func() {
				gomega.Eventually(func() error {
					workName := source.ResourceID(clusterName, newResourceName)
					work, err := agentWorkClient.Get(context.TODO(), workName, metav1.GetOptions{})
					if err != nil {
						return err
					}

					// add finalizers firstly
					patchBytes, err := json.Marshal(map[string]interface{}{
						"metadata": map[string]interface{}{
							"uid":             work.GetUID(),
							"resourceVersion": work.GetResourceVersion(),
							"finalizers":      []string{"work-test-finalizer"},
						},
					})
					gomega.Expect(err).ToNot(gomega.HaveOccurred())

					_, err = agentWorkClient.Patch(context.TODO(), work.Name, apitypes.MergePatchType, patchBytes, metav1.PatchOptions{})
					gomega.Expect(err).ToNot(gomega.HaveOccurred())

					work, err = agentWorkClient.Get(context.TODO(), workName, metav1.GetOptions{})
					if err != nil {
						return err
					}

					if len(work.Finalizers) != 1 {
						return fmt.Errorf("expected finalizers on the work, but got %v", work.Finalizers)
					}

					// update the work status
					newWork := work.DeepCopy()
					newWork.Status = workv1.ManifestWorkStatus{Conditions: []metav1.Condition{{Type: "Created", Status: metav1.ConditionTrue}}}

					oldData, err := json.Marshal(work)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())

					newData, err := json.Marshal(newWork)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())

					patchBytes, err = jsonpatch.CreateMergePatch(oldData, newData)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())

					_, err = agentWorkClient.Patch(context.TODO(), work.Name, apitypes.MergePatchType, patchBytes, metav1.PatchOptions{}, "status")
					gomega.Expect(err).ToNot(gomega.HaveOccurred())

					return nil
				}, 10*time.Second, 1*time.Second).Should(gomega.Succeed())
			})

			ginkgo.By("update the resource on the source and send it to the cluster2", func() {
				var resource *source.Resource
				var err error

				// ensure the resource is created on the cluster
				resourceID := source.ResourceID(clusterName, newResourceName)
				gomega.Eventually(func() error {
					resource, err = source.GetStore().Get(resourceID)
					if err != nil {
						return err
					}

					if !meta.IsStatusConditionTrue(resource.Status.Conditions, "Created") {
						return fmt.Errorf("unexpected status %v", resource.Status.Conditions)
					}

					return nil
				}, 10*time.Second, 1*time.Second).Should(gomega.Succeed())

				resource.ResourceVersion = resource.ResourceVersion + 1
				resource.Spec.Object["data"] = "test"

				err = source.GetStore().Update(resource)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				err = sourceCloudEventsClient.Publish(context.TODO(), types.CloudEventsType{
					CloudEventsDataType: payload.ManifestEventDataType,
					SubResource:         types.SubResourceSpec,
					Action:              "test_update_request",
				}, resource)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})

			ginkgo.By("receive the updated resource on the cluster2", func() {
				gomega.Eventually(func() error {
					workName := source.ResourceID(clusterName, newResourceName)
					work, err := agentWorkClient.Get(context.TODO(), workName, metav1.GetOptions{})
					if err != nil {
						return err
					}

					if len(work.Spec.Workload.Manifests) != 1 {
						return fmt.Errorf("expected manifests in the work, but got %v", work)
					}

					workload := map[string]any{}
					if err := json.Unmarshal(work.Spec.Workload.Manifests[0].Raw, &workload); err != nil {
						return err
					}

					if workload["data"] != "test" {
						return fmt.Errorf("unexpected workload %v", workload)
					}

					return nil
				}, 10*time.Second, 1*time.Second).Should(gomega.Succeed())
			})

			ginkgo.By("mark the resource to deleting on the source and send it to cluster2", func() {
				resourceID := source.ResourceID(clusterName, newResourceName)
				resource, err := source.GetStore().Get(resourceID)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				resource.DeletionTimestamp = &metav1.Time{Time: time.Now()}

				err = source.GetStore().Update(resource)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				err = sourceCloudEventsClient.Publish(context.TODO(), types.CloudEventsType{
					CloudEventsDataType: payload.ManifestEventDataType,
					SubResource:         types.SubResourceSpec,
					Action:              "test_delete_request",
				}, resource)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
			})

			ginkgo.By("receive the deleting resource on the cluster2", func() {
				gomega.Eventually(func() error {
					workName := source.ResourceID(clusterName, newResourceName)
					work, err := agentWorkClient.Get(context.TODO(), workName, metav1.GetOptions{})
					if err != nil {
						return err
					}

					if work.DeletionTimestamp.IsZero() {
						return fmt.Errorf("expected work is deleting, but got %v", work)
					}

					// remove the finalizers
					patchBytes, err := json.Marshal(map[string]interface{}{
						"metadata": map[string]interface{}{
							"uid":             work.GetUID(),
							"resourceVersion": work.GetResourceVersion(),
							"finalizers":      []string{},
						},
					})
					gomega.Expect(err).ToNot(gomega.HaveOccurred())

					_, err = agentWorkClient.Patch(context.TODO(), work.Name, apitypes.MergePatchType, patchBytes, metav1.PatchOptions{})
					gomega.Expect(err).ToNot(gomega.HaveOccurred())

					return nil
				}, 10*time.Second, 1*time.Second).Should(gomega.Succeed())
			})

			ginkgo.By("delete the resource from the source", func() {
				gomega.Eventually(func() error {
					resourceID := source.ResourceID(clusterName, newResourceName)
					resource, err := source.GetStore().Get(resourceID)
					if err != nil {
						return err
					}

					if !meta.IsStatusConditionTrue(resource.Status.Conditions, "Deleted") {
						return fmt.Errorf("unexpected status %v", resource.Status.Conditions)
					}

					source.GetStore().Delete(resourceID)

					return nil
				}, 10*time.Second, 1*time.Second).Should(gomega.Succeed())
			})
		})
	})
})
