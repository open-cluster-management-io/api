// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

var _ = ginkgo.Describe("ManagedCluster API test", func() {
	var clusterName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		clusterName = fmt.Sprintf("managedcluster-%s", suffix)

		managedCluster := &clusterv1.ManagedCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterName,
			},
			Spec: clusterv1.ManagedClusterSpec{
				HubAcceptsClient: true,
			},
		}

		_, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.AfterEach(func() {
		err := hubClusterClient.ClusterV1().ManagedClusters().Delete(context.TODO(), clusterName, metav1.DeleteOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Lease duration should be set by default", func() {
		cluster, err := hubClusterClient.ClusterV1().ManagedClusters().Get(context.TODO(), clusterName, metav1.GetOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		gomega.Expect(cluster.Spec.LeaseDurationSeconds).Should(gomega.Equal(int32(60)))
	})

	ginkgo.It("Taint should be set correctly", func() {
		cluster, err := hubClusterClient.ClusterV1().ManagedClusters().Get(context.TODO(), clusterName, metav1.GetOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		// key, value, effect is required
		cluster.Spec.Taints = []clusterv1.Taint{
			{
				Key:    "",
				Value:  "",
				Effect: "",
			},
		}

		_, err = hubClusterClient.ClusterV1().ManagedClusters().Update(context.TODO(), cluster, metav1.UpdateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())

		// key format is not correct
		cluster.Spec.Taints = []clusterv1.Taint{
			{
				Key:    "test:test",
				Value:  "test",
				Effect: clusterv1.TaintEffectNoSelect,
			},
		}

		_, err = hubClusterClient.ClusterV1().ManagedClusters().Update(context.TODO(), cluster, metav1.UpdateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())

		// Effect is not correct
		cluster.Spec.Taints = []clusterv1.Taint{
			{
				Key:    "test.io/test",
				Value:  "test",
				Effect: "noop",
			},
		}

		_, err = hubClusterClient.ClusterV1().ManagedClusters().Update(context.TODO(), cluster, metav1.UpdateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())

		// key, value, effect is correct and update
		cluster.Spec.Taints = []clusterv1.Taint{
			{
				Key:    "test.io/test",
				Value:  "test",
				Effect: clusterv1.TaintEffectNoSelect,
			},
		}
		_, err = hubClusterClient.ClusterV1().ManagedClusters().Update(context.TODO(), cluster, metav1.UpdateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		// key, value, effect is correct and patch
		_, err = hubClusterClient.ClusterV1().ManagedClusters().Patch(
			context.TODO(),
			cluster.Name,
			types.MergePatchType,
			[]byte(`{"spec":{"taints":[{"key":"test.io/test","value":"testnew","effect":"NoSelect"}]}}`),
			metav1.PatchOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

	})
})

var _ = ginkgo.Describe("ManagedCluster v1 Enhanced API test", func() {
	var clusterName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		clusterName = fmt.Sprintf("managedcluster-enhanced-%s", suffix)
	})

	ginkgo.AfterEach(func() {
		err := hubClusterClient.ClusterV1().ManagedClusters().Delete(context.TODO(), clusterName, metav1.DeleteOptions{})
		if !apierrors.IsNotFound(err) {
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		}
	})

	ginkgo.Context("ManagedCluster ClientConfig validation", func() {
		ginkgo.It("should accept client config without validation", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient: true,
					ManagedClusterClientConfigs: []clusterv1.ClientConfig{
						{
							URL: "https://example.com:6443",
						},
					},
				},
			}

			createdCluster, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(createdCluster.Spec.ManagedClusterClientConfigs[0].URL).Should(gomega.Equal("https://example.com:6443"))
		})

		ginkgo.It("should accept valid HTTPS URL", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient: true,
					ManagedClusterClientConfigs: []clusterv1.ClientConfig{
						{
							URL:      "https://api.example.com:6443",
							CABundle: []byte("dummy-ca-bundle"),
						},
					},
				},
			}

			_, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("should handle multiple client configs", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient: true,
					ManagedClusterClientConfigs: []clusterv1.ClientConfig{
						{
							URL:      "https://api1.example.com:6443",
							CABundle: []byte("ca-bundle-1"),
						},
						{
							URL:      "https://api2.example.com:6443",
							CABundle: []byte("ca-bundle-2"),
						},
					},
				},
			}

			cluster, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(len(cluster.Spec.ManagedClusterClientConfigs)).Should(gomega.Equal(2))
		})
	})

	ginkgo.Context("ManagedCluster Taints advanced validation", func() {
		ginkgo.It("should reject taints with invalid key patterns", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient: true,
					Taints: []clusterv1.Taint{
						{
							Key:    "invalid/key/with/too/many/slashes",
							Value:  "test",
							Effect: clusterv1.TaintEffectNoSelect,
						},
					},
				},
			}

			_, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should accept valid taint with domain prefix", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient: true,
					Taints: []clusterv1.Taint{
						{
							Key:    "example.com/my-taint",
							Value:  "test-value",
							Effect: clusterv1.TaintEffectNoSelect,
						},
					},
				},
			}

			cluster, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(cluster.Spec.Taints[0].Key).Should(gomega.Equal("example.com/my-taint"))
		})

		ginkgo.It("should handle multiple taints with different effects", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient: true,
					Taints: []clusterv1.Taint{
						{
							Key:    "test.io/taint1",
							Value:  "value1",
							Effect: clusterv1.TaintEffectNoSelect,
						},
						{
							Key:    "test.io/taint2",
							Value:  "value2",
							Effect: clusterv1.TaintEffectNoSelectIfNew,
						},
					},
				},
			}

			cluster, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(len(cluster.Spec.Taints)).Should(gomega.Equal(2))
			gomega.Expect(cluster.Spec.Taints[0].Effect).Should(gomega.Equal(clusterv1.TaintEffectNoSelect))
			gomega.Expect(cluster.Spec.Taints[1].Effect).Should(gomega.Equal(clusterv1.TaintEffectNoSelectIfNew))
		})
	})

	ginkgo.Context("ManagedCluster lease duration validation", func() {
		ginkgo.It("should handle custom lease duration", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient:     true,
					LeaseDurationSeconds: 120,
				},
			}

			cluster, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(cluster.Spec.LeaseDurationSeconds).Should(gomega.Equal(int32(120)))
		})

		ginkgo.It("should handle zero lease duration with default", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient:     true,
					LeaseDurationSeconds: 0,
				},
			}

			cluster, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(cluster.Spec.LeaseDurationSeconds).Should(gomega.Equal(int32(60)))
		})
	})

	ginkgo.Context("ManagedCluster status and conditions", func() {
		ginkgo.It("should allow status updates", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient: true,
				},
			}

			cluster, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			// Update status
			now := metav1.NewTime(time.Now())
			cluster.Status = clusterv1.ManagedClusterStatus{
				Version: clusterv1.ManagedClusterVersion{
					Kubernetes: "v1.28.0",
				},
				Allocatable: clusterv1.ResourceList{
					"cpu":    *resource.NewQuantity(4, resource.DecimalSI),
					"memory": *resource.NewQuantity(8*1024*1024*1024, resource.BinarySI),
				},
				Capacity: clusterv1.ResourceList{
					"cpu":    *resource.NewQuantity(4, resource.DecimalSI),
					"memory": *resource.NewQuantity(8*1024*1024*1024, resource.BinarySI),
				},
				Conditions: []metav1.Condition{
					{
						Type:               clusterv1.ManagedClusterConditionAvailable,
						Status:             metav1.ConditionTrue,
						Reason:             "ManagedClusterAvailable",
						LastTransitionTime: now,
					},
				},
			}

			_, err = hubClusterClient.ClusterV1().ManagedClusters().UpdateStatus(context.TODO(), cluster, metav1.UpdateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("ManagedCluster patch operations", func() {
		ginkgo.It("should support strategic merge patch for taints", func() {
			managedCluster := &clusterv1.ManagedCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterName,
				},
				Spec: clusterv1.ManagedClusterSpec{
					HubAcceptsClient: true,
					Taints: []clusterv1.Taint{
						{
							Key:    "initial.io/taint",
							Value:  "initial",
							Effect: clusterv1.TaintEffectNoSelect,
						},
					},
				},
			}

			_, err := hubClusterClient.ClusterV1().ManagedClusters().Create(context.TODO(), managedCluster, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			// Patch to add another taint
			patchData := `{"spec":{"taints":[{"key":"initial.io/taint","value":"initial","effect":"NoSelect"},{"key":"new.io/taint","value":"new","effect":"NoSelect"}]}}`
			_, err = hubClusterClient.ClusterV1().ManagedClusters().Patch(
				context.TODO(),
				clusterName,
				types.MergePatchType,
				[]byte(patchData),
				metav1.PatchOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			// Verify patch was applied
			cluster, err := hubClusterClient.ClusterV1().ManagedClusters().Get(context.TODO(), clusterName, metav1.GetOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(len(cluster.Spec.Taints)).Should(gomega.Equal(2))
		})
	})
})
