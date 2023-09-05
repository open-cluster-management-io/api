package api

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
				InstallStrategy: addonv1alpha1.InstallStrategy{
					Type: addonv1alpha1.AddonInstallStrategyManual,
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

	ginkgo.It("Should not create a ClusterManagementAddOn with empty spec via client", func() {
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
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("Should not create a ClusterManagementAddOn when its configuration resource is empty", func() {
		clusterManagementAddOn := &addonv1alpha1.ClusterManagementAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagementAddOnName,
			},
			Spec: addonv1alpha1.ClusterManagementAddOnSpec{
				InstallStrategy: addonv1alpha1.InstallStrategy{
					Type: addonv1alpha1.AddonInstallStrategyManual,
				},
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
				InstallStrategy: addonv1alpha1.InstallStrategy{
					Type: addonv1alpha1.AddonInstallStrategyManual,
				},
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

	ginkgo.It("Should update the ClusterManagementAddOn status", func() {
		clusterManagementAddOn := &addonv1alpha1.ClusterManagementAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagementAddOnName,
			},
			Spec: addonv1alpha1.ClusterManagementAddOnSpec{
				InstallStrategy: addonv1alpha1.InstallStrategy{
					Type: addonv1alpha1.AddonInstallStrategyManual,
				},
			},
		}

		cma, err := hubAddonClient.AddonV1alpha1().ClusterManagementAddOns().Create(
			context.TODO(),
			clusterManagementAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		cma.Status.DefaultConfigReferences = []addonv1alpha1.DefaultConfigReference{
			{
				ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
					Group:    "test.group",
					Resource: "tests",
				},
				DesiredConfig: &addonv1alpha1.ConfigSpecHash{
					ConfigReferent: addonv1alpha1.ConfigReferent{
						Namespace: testNamespace,
						Name:      "test1",
					},
					SpecHash: "test-spec-hash",
				},
			},
		}
		cma.Status.InstallProgressions = []addonv1alpha1.InstallProgression{
			{
				PlacementRef: addonv1alpha1.PlacementRef{
					Name:      "test",
					Namespace: testNamespace,
				},
				ConfigReferences: []addonv1alpha1.InstallConfigReference{
					{
						ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
							Group:    "test.group",
							Resource: "tests",
						},
						DesiredConfig: &addonv1alpha1.ConfigSpecHash{
							ConfigReferent: addonv1alpha1.ConfigReferent{
								Namespace: testNamespace,
								Name:      "test2",
							},
							SpecHash: "test-spec-hash",
						},
					},
				},
			},
		}

		_, err = hubAddonClient.AddonV1alpha1().ClusterManagementAddOns().UpdateStatus(
			context.TODO(),
			cma,
			metav1.UpdateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

	})

	ginkgo.It("Should update the ClusterManagementAddOn status with empty status", func() {
		clusterManagementAddOn := &addonv1alpha1.ClusterManagementAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagementAddOnName,
			},
			Spec: addonv1alpha1.ClusterManagementAddOnSpec{
				InstallStrategy: addonv1alpha1.InstallStrategy{
					Type: addonv1alpha1.AddonInstallStrategyManual,
				},
			},
		}

		cma, err := hubAddonClient.AddonV1alpha1().ClusterManagementAddOns().Create(
			context.TODO(),
			clusterManagementAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		_, err = hubAddonClient.AddonV1alpha1().ClusterManagementAddOns().UpdateStatus(
			context.TODO(),
			cma,
			metav1.UpdateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

	})
})
