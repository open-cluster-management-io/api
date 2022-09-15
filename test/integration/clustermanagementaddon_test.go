package integration

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
)

var _ = ginkgo.Describe("ClusterManagementAddOn API test", func() {

	var clusterManagementAddOnName string

	ginkgo.BeforeEach(func() {
		clusterManagementAddOnName = fmt.Sprintf("cma-%s", rand.String(5))
	})

	ginkgo.It("Should create a ClusterManagementAddOn", func() {
		clusterManagementAddOn := &addonv1alpha1.ClusterManagementAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagementAddOnName,
			},
			Spec: addonv1alpha1.ClusterManagementAddOnSpec{
				AddOnMeta: addonv1alpha1.AddOnMeta{
					DisplayName: "test",
					Description: "for test",
				},
				SupportedConfigs: []addonv1alpha1.ConfigMeta{
					{
						ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
							Group:    "test.addon",
							Resource: "tests",
						},
						DefaultConfig: &addonv1alpha1.ConfigReferent{
							Namespace: testNamespace,
							Name:      "test",
						},
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ClusterManagementAddOns().Create(
			context.TODO(),
			clusterManagementAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Should create a ClusterManagementAddOn with empty spec", func() {
		clusterManagementAddOn := &addonv1alpha1.ClusterManagementAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagementAddOnName,
			},
			Spec: addonv1alpha1.ClusterManagementAddOnSpec{},
		}

		_, err := hubAddonClient.AddonV1alpha1().ClusterManagementAddOns().Create(
			context.TODO(),
			clusterManagementAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Should not create a ClusterManagementAddOn when its configuration resource is empty", func() {
		clusterManagementAddOn := &addonv1alpha1.ClusterManagementAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagementAddOnName,
			},
			Spec: addonv1alpha1.ClusterManagementAddOnSpec{
				SupportedConfigs: []addonv1alpha1.ConfigMeta{
					{
						ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "test.addon"},
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ClusterManagementAddOns().Create(
			context.TODO(),
			clusterManagementAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.It("Should not create a ClusterManagementAddOn when its configuration name is empty", func() {
		clusterManagementAddOn := &addonv1alpha1.ClusterManagementAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagementAddOnName,
			},
			Spec: addonv1alpha1.ClusterManagementAddOnSpec{
				SupportedConfigs: []addonv1alpha1.ConfigMeta{
					{
						ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
							Group:    "test.addon",
							Resource: "tests",
						},
						DefaultConfig: &addonv1alpha1.ConfigReferent{},
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ClusterManagementAddOns().Create(
			context.TODO(),
			clusterManagementAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})
})
