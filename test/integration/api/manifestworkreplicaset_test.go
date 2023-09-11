package api

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	workv1alpha1 "open-cluster-management.io/api/work/v1alpha1"
)

var _ = ginkgo.Describe("ManifestWorkSet API test", func() {
	var manifestWorkSetName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		manifestWorkSetName = fmt.Sprintf("manifestworkset-%s", suffix)
	})

	ginkgo.Context("Placement refs", func() {
		ginkgo.It("Create a ManifestWorkSet with empty placement refs", func() {
			workSet := &workv1alpha1.ManifestWorkReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkSetName,
				},
				Spec: workv1alpha1.ManifestWorkReplicaSetSpec{
					PlacementRefs: []workv1alpha1.LocalPlacementReference{},
				},
			}
			_, err := hubWorkClient.WorkV1alpha1().ManifestWorkReplicaSets(testNamespace).
				Create(context.TODO(), workSet, metav1.CreateOptions{})
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("Create a ManifestWorkSet with a placement refs and empty name", func() {
			workSet := &workv1alpha1.ManifestWorkReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkSetName,
				},
				Spec: workv1alpha1.ManifestWorkReplicaSetSpec{
					PlacementRefs: []workv1alpha1.LocalPlacementReference{
						{
							Name: "",
						},
					},
				},
			}
			_, err := hubWorkClient.WorkV1alpha1().ManifestWorkReplicaSets(testNamespace).
				Create(context.TODO(), workSet, metav1.CreateOptions{})
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("Create a ManifestWorkSet successfully", func() {
			workSet := &workv1alpha1.ManifestWorkReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: manifestWorkSetName,
				},
				Spec: workv1alpha1.ManifestWorkReplicaSetSpec{
					PlacementRefs: []workv1alpha1.LocalPlacementReference{
						{
							Name: "testPlacement",
						},
					},
				},
			}
			_, err := hubWorkClient.WorkV1alpha1().ManifestWorkReplicaSets(testNamespace).
				Create(context.TODO(), workSet, metav1.CreateOptions{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})
	})
})
