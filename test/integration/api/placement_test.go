// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
)

var _ = ginkgo.Describe("Placement API test", func() {
	var placementName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		placementName = fmt.Sprintf("placement-%s", suffix)
	})

	ginkgo.AfterEach(func() {
		err := hubClusterClient.ClusterV1beta1().Placements(testNamespace).Delete(context.TODO(), placementName, metav1.DeleteOptions{})
		if err != nil {
			// Ignore not found errors during cleanup
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		}
	})

	ginkgo.Context("Placement creation and validation", func() {
		ginkgo.It("should create placement with empty spec", func() {
			placement := &clusterv1beta1.Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      placementName,
					Namespace: testNamespace,
				},
				Spec: clusterv1beta1.PlacementSpec{},
			}

			createdPlacement, err := hubClusterClient.ClusterV1beta1().Placements(testNamespace).Create(context.TODO(), placement, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(createdPlacement.Name).Should(gomega.Equal(placementName))
			gomega.Expect(createdPlacement.Namespace).Should(gomega.Equal(testNamespace))
		})

		ginkgo.It("should create placement with cluster sets", func() {
			placement := &clusterv1beta1.Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      placementName,
					Namespace: testNamespace,
				},
				Spec: clusterv1beta1.PlacementSpec{
					ClusterSets: []string{"clusterset1", "clusterset2"},
				},
			}

			createdPlacement, err := hubClusterClient.ClusterV1beta1().Placements(testNamespace).Create(context.TODO(), placement, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(createdPlacement.Spec.ClusterSets).Should(gomega.Equal([]string{"clusterset1", "clusterset2"}))
		})

		ginkgo.It("should create placement with number of clusters", func() {
			numberOfClusters := int32(3)
			placement := &clusterv1beta1.Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      placementName,
					Namespace: testNamespace,
				},
				Spec: clusterv1beta1.PlacementSpec{
					NumberOfClusters: &numberOfClusters,
				},
			}

			createdPlacement, err := hubClusterClient.ClusterV1beta1().Placements(testNamespace).Create(context.TODO(), placement, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(*createdPlacement.Spec.NumberOfClusters).Should(gomega.Equal(int32(3)))
		})

		ginkgo.It("should create placement with label selector predicate", func() {
			placement := &clusterv1beta1.Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      placementName,
					Namespace: testNamespace,
				},
				Spec: clusterv1beta1.PlacementSpec{
					Predicates: []clusterv1beta1.ClusterPredicate{
						{
							RequiredClusterSelector: clusterv1beta1.ClusterSelector{
								LabelSelector: metav1.LabelSelector{
									MatchLabels: map[string]string{
										"environment": "production",
									},
								},
							},
						},
					},
				},
			}

			createdPlacement, err := hubClusterClient.ClusterV1beta1().Placements(testNamespace).Create(context.TODO(), placement, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(len(createdPlacement.Spec.Predicates)).Should(gomega.Equal(1))
			gomega.Expect(createdPlacement.Spec.Predicates[0].RequiredClusterSelector.LabelSelector.MatchLabels["environment"]).Should(gomega.Equal("production"))
		})

		ginkgo.It("should create placement with tolerations", func() {
			placement := &clusterv1beta1.Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      placementName,
					Namespace: testNamespace,
				},
				Spec: clusterv1beta1.PlacementSpec{
					Tolerations: []clusterv1beta1.Toleration{
						{
							Key:      "node.kubernetes.io/unreachable",
							Operator: clusterv1beta1.TolerationOpExists,
							Effect:   clusterv1.TaintEffectNoSelect,
						},
					},
				},
			}

			createdPlacement, err := hubClusterClient.ClusterV1beta1().Placements(testNamespace).Create(context.TODO(), placement, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(len(createdPlacement.Spec.Tolerations)).Should(gomega.Equal(1))
			gomega.Expect(createdPlacement.Spec.Tolerations[0].Key).Should(gomega.Equal("node.kubernetes.io/unreachable"))
		})
	})

	ginkgo.Context("Placement validation", func() {
		ginkgo.It("should accept valid cluster predicate", func() {
			placement := &clusterv1beta1.Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      placementName,
					Namespace: testNamespace,
				},
				Spec: clusterv1beta1.PlacementSpec{
					Predicates: []clusterv1beta1.ClusterPredicate{
						{
							RequiredClusterSelector: clusterv1beta1.ClusterSelector{
								LabelSelector: metav1.LabelSelector{
									MatchLabels: map[string]string{
										"app": "test",
									},
								},
							},
						},
					},
				},
			}

			createdPlacement, err := hubClusterClient.ClusterV1beta1().Placements(testNamespace).Create(context.TODO(), placement, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(createdPlacement.Spec.Predicates[0].RequiredClusterSelector.LabelSelector.MatchLabels["app"]).Should(gomega.Equal("test"))
		})

		ginkgo.It("should accept positive number of clusters", func() {
			numberOfClusters := int32(5)
			placement := &clusterv1beta1.Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      placementName,
					Namespace: testNamespace,
				},
				Spec: clusterv1beta1.PlacementSpec{
					NumberOfClusters: &numberOfClusters,
				},
			}

			createdPlacement, err := hubClusterClient.ClusterV1beta1().Placements(testNamespace).Create(context.TODO(), placement, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(*createdPlacement.Spec.NumberOfClusters).Should(gomega.Equal(int32(5)))
		})
	})

	ginkgo.Context("Placement updates", func() {
		var createdPlacement *clusterv1beta1.Placement

		ginkgo.BeforeEach(func() {
			placement := &clusterv1beta1.Placement{
				ObjectMeta: metav1.ObjectMeta{
					Name:      placementName,
					Namespace: testNamespace,
				},
				Spec: clusterv1beta1.PlacementSpec{},
			}

			var err error
			createdPlacement, err = hubClusterClient.ClusterV1beta1().Placements(testNamespace).Create(context.TODO(), placement, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("should update placement spec", func() {
			createdPlacement.Spec.ClusterSets = []string{"updated-clusterset"}
			numberOfClusters := int32(5)
			createdPlacement.Spec.NumberOfClusters = &numberOfClusters

			updatedPlacement, err := hubClusterClient.ClusterV1beta1().Placements(testNamespace).Update(context.TODO(), createdPlacement, metav1.UpdateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(updatedPlacement.Spec.ClusterSets).Should(gomega.Equal([]string{"updated-clusterset"}))
			gomega.Expect(*updatedPlacement.Spec.NumberOfClusters).Should(gomega.Equal(int32(5)))
		})
	})
})
