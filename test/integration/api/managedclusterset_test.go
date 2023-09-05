package api

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	clusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
)

var _ = ginkgo.Describe("ManagedClusterSet API test", func() {
	var clusterSetName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		clusterSetName = fmt.Sprintf("clusterset-%s", suffix)
	})

	ginkgo.It("Create a clusterset with empty spec", func() {
		clusterset := &clusterv1beta2.ManagedClusterSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterSetName,
			},
		}
		clusterset, err := hubClusterClient.ClusterV1beta2().ManagedClusterSets().Create(context.TODO(), clusterset, metav1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		gomega.Expect(clusterset.Spec.ClusterSelector.SelectorType).Should(gomega.Equal(clusterv1beta2.ExclusiveClusterSetLabel))
	})

	ginkgo.It("Create a clusterset with wrong selector", func() {
		clusterset := &clusterv1beta2.ManagedClusterSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterSetName,
			},
			Spec: clusterv1beta2.ManagedClusterSetSpec{
				ClusterSelector: clusterv1beta2.ManagedClusterSelector{
					SelectorType: "WrongSelector",
				},
			},
		}
		_, err := hubClusterClient.ClusterV1beta2().ManagedClusterSets().Create(context.TODO(), clusterset, metav1.CreateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("Create a clusterset with null label selector", func() {
		clusterset := &clusterv1beta2.ManagedClusterSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterSetName,
			},
			Spec: clusterv1beta2.ManagedClusterSetSpec{
				ClusterSelector: clusterv1beta2.ManagedClusterSelector{
					SelectorType:  "LabelSelector",
					LabelSelector: &metav1.LabelSelector{},
				},
			},
		}
		_, err := hubClusterClient.ClusterV1beta2().ManagedClusterSets().Create(context.TODO(), clusterset, metav1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Create a clusterset with label selector(one label)", func() {
		clusterset := &clusterv1beta2.ManagedClusterSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterSetName,
			},
			Spec: clusterv1beta2.ManagedClusterSetSpec{
				ClusterSelector: clusterv1beta2.ManagedClusterSelector{
					SelectorType: "LabelSelector",
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"vendor": "OpenShift",
						},
					},
				},
			},
		}
		_, err := hubClusterClient.ClusterV1beta2().ManagedClusterSets().Create(context.TODO(), clusterset, metav1.CreateOptions{})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})
})
