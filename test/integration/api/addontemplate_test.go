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
	workv1 "open-cluster-management.io/api/work/v1"
)

var _ = ginkgo.Describe("AddOnTemplate v1beta1 API test", func() {
	var addOnTemplateName string
	var addonName string

	ginkgo.BeforeEach(func() {
		suffix := rand.String(5)
		addOnTemplateName = fmt.Sprintf("addon-template-%s", suffix)
		addonName = fmt.Sprintf("addon-%s", suffix)
	})

	ginkgo.It("Should create an AddOnTemplate using v1beta1", func() {
		addOnTemplate := &addonv1beta1.AddOnTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name: addOnTemplateName,
			},
			Spec: addonv1beta1.AddOnTemplateSpec{
				AddonName: addonName,
				AgentSpec: workv1.ManifestWorkSpec{},
			},
		}

		_, err := hubAddonClient.AddonV1beta1().AddOnTemplates().Create(
			context.TODO(),
			addOnTemplate,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Should create v1beta1 and read as v1alpha1 (cross-version compatibility)", func() {
		addOnTemplateV1beta1 := &addonv1beta1.AddOnTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name: addOnTemplateName,
			},
			Spec: addonv1beta1.AddOnTemplateSpec{
				AddonName: addonName,
				AgentSpec: workv1.ManifestWorkSpec{},
			},
		}

		createdV1beta1, err := hubAddonClient.AddonV1beta1().AddOnTemplates().Create(
			context.TODO(),
			addOnTemplateV1beta1,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(createdV1beta1.Spec.AddonName).To(gomega.Equal(addonName))

		retrievedV1alpha1, err := hubAddonClient.AddonV1alpha1().AddOnTemplates().Get(
			context.TODO(),
			addOnTemplateName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(retrievedV1alpha1.Spec.AddonName).To(gomega.Equal(addonName))
	})

	ginkgo.It("Should create v1alpha1 and read as v1beta1 (cross-version compatibility)", func() {
		addOnTemplateV1alpha1 := &addonv1alpha1.AddOnTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name: addOnTemplateName,
			},
			Spec: addonv1alpha1.AddOnTemplateSpec{
				AddonName: addonName,
				AgentSpec: workv1.ManifestWorkSpec{},
			},
		}

		createdV1alpha1, err := hubAddonClient.AddonV1alpha1().AddOnTemplates().Create(
			context.TODO(),
			addOnTemplateV1alpha1,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(createdV1alpha1.Spec.AddonName).To(gomega.Equal(addonName))

		retrievedV1beta1, err := hubAddonClient.AddonV1beta1().AddOnTemplates().Get(
			context.TODO(),
			addOnTemplateName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(retrievedV1beta1.Spec.AddonName).To(gomega.Equal(addonName))
	})

	ginkgo.It("Should round-trip a CustomSigner subject's organizationUnit through v1beta1 and v1alpha1", func() {
		// Regression test: AddOnTemplate's CustomSigner.Subject must keep the exact
		// v1alpha1 JSON shape (organizationUnit) in v1beta1 too, since there is no
		// conversion webhook for this CRD (unlike ManagedClusterAddOn/
		// ClusterManagementAddOn) - conversion relies entirely on both versions
		// serializing identically.
		addOnTemplateV1beta1 := &addonv1beta1.AddOnTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name: addOnTemplateName,
			},
			Spec: addonv1beta1.AddOnTemplateSpec{
				AddonName: addonName,
				AgentSpec: workv1.ManifestWorkSpec{},
				Registration: []addonv1beta1.RegistrationSpec{
					{
						Type: addonv1beta1.RegistrationTypeCustomSigner,
						CustomSigner: &addonv1beta1.CustomSignerRegistrationConfig{
							SignerName: "example.com/my-signer",
							Subject: &addonv1beta1.AddOnTemplateSubject{
								User:              "my-user",
								OrganizationUnits: []string{"my-ou"},
							},
							SigningCA: addonv1beta1.SigningCARef{
								Name: "signing-ca",
							},
						},
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1beta1().AddOnTemplates().Create(
			context.TODO(),
			addOnTemplateV1beta1,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		retrievedV1alpha1, err := hubAddonClient.AddonV1alpha1().AddOnTemplates().Get(
			context.TODO(),
			addOnTemplateName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(retrievedV1alpha1.Spec.Registration).To(gomega.HaveLen(1))
		gomega.Expect(retrievedV1alpha1.Spec.Registration[0].CustomSigner.Subject.OrganizationUnits).
			To(gomega.Equal([]string{"my-ou"}))
	})

	ginkgo.It("Should not create an AddOnTemplate with an invalid custom signer name", func() {
		addOnTemplate := &addonv1beta1.AddOnTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name: addOnTemplateName,
			},
			Spec: addonv1beta1.AddOnTemplateSpec{
				AddonName: addonName,
				AgentSpec: workv1.ManifestWorkSpec{},
				Registration: []addonv1beta1.RegistrationSpec{
					{
						Type: addonv1beta1.RegistrationTypeCustomSigner,
						CustomSigner: &addonv1beta1.CustomSignerRegistrationConfig{
							// too short and missing the required "domain/path" pattern
							SignerName: "bad",
							SigningCA: addonv1beta1.SigningCARef{
								Name: "signing-ca",
							},
						},
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1beta1().AddOnTemplates().Create(
			context.TODO(),
			addOnTemplate,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})
})
