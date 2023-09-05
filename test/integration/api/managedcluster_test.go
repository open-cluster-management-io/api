package api

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

		// Effect is not correct
		cluster.Spec.Taints = []clusterv1.Taint{
			{
				Key:    "test.io/test",
				Value:  "test",
				Effect: clusterv1.TaintEffectNoSelect,
			},
		}

		_, err = hubClusterClient.ClusterV1().ManagedClusters().Update(context.TODO(), cluster, metav1.UpdateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

	})
})
