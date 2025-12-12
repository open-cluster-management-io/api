// Copyright Contributors to the Open Cluster Management project
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
	addonv1beta1 "open-cluster-management.io/api/addon/v1beta1"
)

const installNamespaceMaxLength = 63

var _ = ginkgo.Describe("ManagedClusterAddOn API test", func() {
	var managedClusterAddOnName string

	ginkgo.BeforeEach(func() {
		managedClusterAddOnName = fmt.Sprintf("mca-%s", rand.String(5))
	})

	ginkgo.It("Should create a ManagedClusterAddOn", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{
				InstallNamespace: testNamespace,
				Configs: []addonv1alpha1.AddOnConfig{
					{
						ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
							Group:    "test.group",
							Resource: "tests",
						},
						ConfigReferent: addonv1alpha1.ConfigReferent{
							Namespace: testNamespace,
							Name:      "test",
						},
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		mca, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Get(
			context.TODO(),
			managedClusterAddOnName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(mca.Spec.InstallNamespace).To(gomega.BeEquivalentTo(testNamespace)) //nolint:staticcheck
	})

	ginkgo.It("Should create a ManagedClusterAddOn with empty spec", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		mca, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Get(
			context.TODO(),
			managedClusterAddOnName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(mca.Spec.InstallNamespace).To(gomega.BeEquivalentTo("open-cluster-management-agent-addon")) //nolint:staticcheck
	})

	ginkgo.It("Should update the ManagedClusterAddOn status without config", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		mca, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Get(
			context.TODO(),
			managedClusterAddOnName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		mca.Status.Registrations = []addonv1alpha1.RegistrationConfig{
			{
				SignerName: "open-cluster-management.io/addontest",
			},
		}

		_, err = hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).UpdateStatus(
			context.TODO(),
			mca,
			metav1.UpdateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Update failed with wrong signer name in the ManagedClusterAddOn", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		mca, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Get(
			context.TODO(),
			managedClusterAddOnName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		mca.Status.Registrations = []addonv1alpha1.RegistrationConfig{
			{
				SignerName: "addontest",
			},
		}

		_, err = hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).UpdateStatus(
			context.TODO(),
			mca,
			metav1.UpdateOptions{},
		)
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.It("Should update the ManagedClusterAddOn status with config", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{
				InstallNamespace: testNamespace,
				Configs: []addonv1alpha1.AddOnConfig{
					{
						ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
							Group:    "test.group",
							Resource: "tests",
						},
						ConfigReferent: addonv1alpha1.ConfigReferent{
							Namespace: testNamespace,
							Name:      "test",
						},
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		mca, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Get(
			context.TODO(),
			managedClusterAddOnName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(mca.Spec.InstallNamespace).To(gomega.BeEquivalentTo(testNamespace)) //nolint:staticcheck

		mca.Status.ConfigReferences = []addonv1alpha1.ConfigReference{
			{
				ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
					Group:    "test.group",
					Resource: "tests",
				},
				ConfigReferent: addonv1alpha1.ConfigReferent{
					Namespace: testNamespace,
					Name:      "test",
				},
				LastObservedGeneration: 1,
			},
		}

		_, err = hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).UpdateStatus(
			context.TODO(),
			mca,
			metav1.UpdateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

	})

	ginkgo.It("Should create a ManagedClusterAddOn with install namespace", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{
				InstallNamespace: rand.String(installNamespaceMaxLength),
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Should not create a ManagedClusterAddOn with a wrong install namespace", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{
				InstallNamespace: "#test",
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.It("Should not create a ManagedClusterAddOn with a long install namespace", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{
				InstallNamespace: rand.String(installNamespaceMaxLength + 1),
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.It("Should not create a ManagedClusterAddOn when its config type is empty", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{
				Configs: []addonv1alpha1.AddOnConfig{
					{
						ConfigReferent: addonv1alpha1.ConfigReferent{
							Namespace: testNamespace,
							Name:      "test",
						},
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.It("Should not create a ManagedClusterAddOn when its config name is empty", func() {
		managedClusterAddOn := &addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				Name: managedClusterAddOnName,
			},
			Spec: addonv1alpha1.ManagedClusterAddOnSpec{
				Configs: []addonv1alpha1.AddOnConfig{
					{
						ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
							Group:    "test.group",
							Resource: "tests",
						},
						ConfigReferent: addonv1alpha1.ConfigReferent{
							Namespace: testNamespace,
						},
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(
			context.TODO(),
			managedClusterAddOn,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.Describe("ManagedClusterAddOn v1beta1 API test", func() {
		ginkgo.It("Should create a ManagedClusterAddOn with v1beta1", func() {
			managedClusterAddOn := &addonv1beta1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: managedClusterAddOnName,
				},
				Spec: addonv1beta1.ManagedClusterAddOnSpec{
					Configs: []addonv1beta1.AddOnConfig{
						{
							ConfigGroupResource: addonv1beta1.ConfigGroupResource{
								Group:    "test.group",
								Resource: "tests",
							},
							ConfigReferent: addonv1beta1.ConfigReferent{
								Namespace: testNamespace,
								Name:      "test",
							},
						},
					},
				},
			}

			_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(
				context.TODO(),
				managedClusterAddOn,
				metav1.CreateOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Get(
				context.TODO(),
				managedClusterAddOnName,
				metav1.GetOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(mca).ToNot(gomega.BeNil())
		})

		ginkgo.It("Should create a ManagedClusterAddOn with empty spec v1beta1", func() {
			managedClusterAddOn := &addonv1beta1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: managedClusterAddOnName,
				},
				Spec: addonv1beta1.ManagedClusterAddOnSpec{},
			}

			_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(
				context.TODO(),
				managedClusterAddOn,
				metav1.CreateOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Get(
				context.TODO(),
				managedClusterAddOnName,
				metav1.GetOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(mca).ToNot(gomega.BeNil())
		})

		ginkgo.It("Should update the ManagedClusterAddOn status without config v1beta1", func() {
			managedClusterAddOn := &addonv1beta1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: managedClusterAddOnName,
				},
				Spec: addonv1beta1.ManagedClusterAddOnSpec{},
			}

			_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(
				context.TODO(),
				managedClusterAddOn,
				metav1.CreateOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Get(
				context.TODO(),
				managedClusterAddOnName,
				metav1.GetOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca.Status.Namespace = testNamespace
			mca.Status.Registrations = []addonv1beta1.RegistrationConfig{
				{
					Type: addonv1beta1.CSR,
					CSR: &addonv1beta1.CSRConfig{
						SignerName: "open-cluster-management.io/addontest",
					},
				},
			}

			_, err = hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).UpdateStatus(
				context.TODO(),
				mca,
				metav1.UpdateOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("Should update the ManagedClusterAddOn status with KubeClient registration v1beta1", func() {
			managedClusterAddOn := &addonv1beta1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: managedClusterAddOnName,
				},
				Spec: addonv1beta1.ManagedClusterAddOnSpec{},
			}

			_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(
				context.TODO(),
				managedClusterAddOn,
				metav1.CreateOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Get(
				context.TODO(),
				managedClusterAddOnName,
				metav1.GetOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca.Status.Namespace = testNamespace
			mca.Status.Registrations = []addonv1beta1.RegistrationConfig{
				{
					Type: addonv1beta1.KubeClient,
					KubeClient: &addonv1beta1.KubeClientConfig{
						Subject: addonv1beta1.KubeClientSubject{
							BaseSubject: addonv1beta1.BaseSubject{
								User:   "test-user",
								Groups: []string{"test-group"},
							},
						},
					},
				},
			}

			_, err = hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).UpdateStatus(
				context.TODO(),
				mca,
				metav1.UpdateOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("Update failed with wrong signer name in the ManagedClusterAddOn v1beta1", func() {
			managedClusterAddOn := &addonv1beta1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: managedClusterAddOnName,
				},
				Spec: addonv1beta1.ManagedClusterAddOnSpec{},
			}

			_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(
				context.TODO(),
				managedClusterAddOn,
				metav1.CreateOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Get(
				context.TODO(),
				managedClusterAddOnName,
				metav1.GetOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca.Status.Registrations = []addonv1beta1.RegistrationConfig{
				{
					Type: addonv1beta1.CSR,
					CSR: &addonv1beta1.CSRConfig{
						SignerName: "addontest",
					},
				},
			}

			_, err = hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).UpdateStatus(
				context.TODO(),
				mca,
				metav1.UpdateOptions{},
			)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("Should update the ManagedClusterAddOn status with config v1beta1", func() {
			managedClusterAddOn := &addonv1beta1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: managedClusterAddOnName,
				},
				Spec: addonv1beta1.ManagedClusterAddOnSpec{
					Configs: []addonv1beta1.AddOnConfig{
						{
							ConfigGroupResource: addonv1beta1.ConfigGroupResource{
								Group:    "test.group",
								Resource: "tests",
							},
							ConfigReferent: addonv1beta1.ConfigReferent{
								Namespace: testNamespace,
								Name:      "test",
							},
						},
					},
				},
			}

			_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(
				context.TODO(),
				managedClusterAddOn,
				metav1.CreateOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Get(
				context.TODO(),
				managedClusterAddOnName,
				metav1.GetOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			mca.Status.ConfigReferences = []addonv1beta1.ConfigReference{
				{
					ConfigGroupResource: addonv1beta1.ConfigGroupResource{
						Group:    "test.group",
						Resource: "tests",
					},
					LastObservedGeneration: 1,
					DesiredConfig: &addonv1beta1.ConfigSpecHash{
						ConfigReferent: addonv1beta1.ConfigReferent{
							Namespace: testNamespace,
							Name:      "test",
						},
						SpecHash: "test-spec-hash",
					},
				},
			}

			_, err = hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).UpdateStatus(
				context.TODO(),
				mca,
				metav1.UpdateOptions{},
			)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.It("Should not create a ManagedClusterAddOn when its config type is empty v1beta1", func() {
			managedClusterAddOn := &addonv1beta1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: managedClusterAddOnName,
				},
				Spec: addonv1beta1.ManagedClusterAddOnSpec{
					Configs: []addonv1beta1.AddOnConfig{
						{
							ConfigReferent: addonv1beta1.ConfigReferent{
								Namespace: testNamespace,
								Name:      "test",
							},
						},
					},
				},
			}

			_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(
				context.TODO(),
				managedClusterAddOn,
				metav1.CreateOptions{},
			)
			gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
		})

		ginkgo.It("Should not create a ManagedClusterAddOn when its config name is empty v1beta1", func() {
			managedClusterAddOn := &addonv1beta1.ManagedClusterAddOn{
				ObjectMeta: metav1.ObjectMeta{
					Name: managedClusterAddOnName,
				},
				Spec: addonv1beta1.ManagedClusterAddOnSpec{
					Configs: []addonv1beta1.AddOnConfig{
						{
							ConfigGroupResource: addonv1beta1.ConfigGroupResource{
								Group:    "test.group",
								Resource: "tests",
							},
							ConfigReferent: addonv1beta1.ConfigReferent{
								Namespace: testNamespace,
							},
						},
					},
				},
			}

			_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(
				context.TODO(),
				managedClusterAddOn,
				metav1.CreateOptions{},
			)
			gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
		})
	})
})
