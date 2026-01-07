// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	addonv1beta1 "open-cluster-management.io/api/addon/v1beta1"
)

const (
	variableNameMaxLength  = 255
	variableValueMaxLength = 1024
)

var _ = ginkgo.Describe("AddOnDeploymentConfig API test", func() {
	var addOnDeploymentConfigName string

	ginkgo.BeforeEach(func() {
		addOnDeploymentConfigName = fmt.Sprintf("adc-%s", rand.String(5))
	})

	ginkgo.It("Should create a AddOnDeploymentConfig", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				CustomizedVariables: []addonv1alpha1.CustomizedVariable{
					{
						Name:  fmt.Sprintf("v_%s", rand.String(variableNameMaxLength-2)), // avoid rand string start with number
						Value: rand.String(variableValueMaxLength),
					},
				},
				NodePlacement: &addonv1alpha1.NodePlacement{
					Tolerations: []corev1.Toleration{
						{
							Key:      "foo",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoExecute,
						},
					},
					NodeSelector: map[string]string{
						"kubernetes.io/os": "linux",
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Should create a AddOnDeploymentConfig without node placement", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				CustomizedVariables: []addonv1alpha1.CustomizedVariable{
					{
						Name:  fmt.Sprintf("v_%s", rand.String(variableNameMaxLength-2)), // avoid rand string start with number
						Value: rand.String(variableValueMaxLength),
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Should not create a AddOnDeploymentConfig with a wrong variable name", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				CustomizedVariables: []addonv1alpha1.CustomizedVariable{
					{
						Name:  "@test",
						Value: "test",
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.It("Should not create a AddOnDeploymentConfig with a long variable name", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				CustomizedVariables: []addonv1alpha1.CustomizedVariable{
					{
						Name:  fmt.Sprintf("v_%s", rand.String(variableNameMaxLength)),
						Value: "test",
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.It("Should not create a AddOnDeploymentConfig with a long variable value", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				CustomizedVariables: []addonv1alpha1.CustomizedVariable{
					{
						Name:  "test",
						Value: rand.String(variableValueMaxLength + 1),
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.It("Should create a AddOnDeploymentConfig with empty agentInstallNamespace", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				AgentInstallNamespace: "",
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Should create a AddOnDeploymentConfig with valid agentInstallNamespace", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				AgentInstallNamespace: "my-custom-namespace",
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Should not create a AddOnDeploymentConfig with invalid agentInstallNamespace (starts with hyphen)", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				AgentInstallNamespace: "-invalid",
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.It("Should not create a AddOnDeploymentConfig with invalid agentInstallNamespace (contains uppercase)", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				AgentInstallNamespace: "Invalid-Namespace",
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})

	ginkgo.It("Should not create a AddOnDeploymentConfig with agentInstallNamespace exceeding max length", func() {
		addOnDeploymentConfig := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				AgentInstallNamespace: rand.String(64), // max is 63
			},
		}

		_, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})
})

var _ = ginkgo.Describe("AddOnDeploymentConfig v1beta1 API test", func() {
	var addOnDeploymentConfigName string

	ginkgo.BeforeEach(func() {
		addOnDeploymentConfigName = fmt.Sprintf("adc-v1beta1-%s", rand.String(5))
	})

	ginkgo.It("Should create a AddOnDeploymentConfig using v1beta1", func() {
		addOnDeploymentConfig := &addonv1beta1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1beta1.AddOnDeploymentConfigSpec{
				CustomizedVariables: []addonv1beta1.CustomizedVariable{
					{
						Name:  fmt.Sprintf("v_%s", rand.String(variableNameMaxLength-2)),
						Value: rand.String(variableValueMaxLength),
					},
				},
				NodePlacement: &addonv1beta1.NodePlacement{
					Tolerations: []corev1.Toleration{
						{
							Key:      "foo",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoExecute,
						},
					},
					NodeSelector: map[string]string{
						"kubernetes.io/os": "linux",
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1beta1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.It("Should create v1beta1 and read as v1alpha1 (cross-version compatibility)", func() {
		// Create using v1beta1
		addOnDeploymentConfigV1beta1 := &addonv1beta1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1beta1.AddOnDeploymentConfigSpec{
				CustomizedVariables: []addonv1beta1.CustomizedVariable{
					{
						Name:  "test_var",
						Value: "test_value",
					},
				},
				AgentInstallNamespace: "test-namespace",
			},
		}

		createdV1beta1, err := hubAddonClient.AddonV1beta1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfigV1beta1,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(createdV1beta1.Spec.CustomizedVariables).To(gomega.HaveLen(1))
		gomega.Expect(createdV1beta1.Spec.CustomizedVariables[0].Name).To(gomega.Equal("test_var"))

		// Read using v1alpha1
		retrievedV1alpha1, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Get(
			context.TODO(),
			addOnDeploymentConfigName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(retrievedV1alpha1.Spec.CustomizedVariables).To(gomega.HaveLen(1))
		gomega.Expect(retrievedV1alpha1.Spec.CustomizedVariables[0].Name).To(gomega.Equal("test_var"))
		gomega.Expect(retrievedV1alpha1.Spec.AgentInstallNamespace).To(gomega.Equal("test-namespace"))
	})

	ginkgo.It("Should create v1alpha1 and read as v1beta1 (cross-version compatibility)", func() {
		// Create using v1alpha1
		addOnDeploymentConfigV1alpha1 := &addonv1alpha1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1alpha1.AddOnDeploymentConfigSpec{
				CustomizedVariables: []addonv1alpha1.CustomizedVariable{
					{
						Name:  "alpha_var",
						Value: "alpha_value",
					},
				},
				AgentInstallNamespace: "alpha-namespace",
			},
		}

		createdV1alpha1, err := hubAddonClient.AddonV1alpha1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfigV1alpha1,
			metav1.CreateOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(createdV1alpha1.Spec.CustomizedVariables).To(gomega.HaveLen(1))

		// Read using v1beta1
		retrievedV1beta1, err := hubAddonClient.AddonV1beta1().AddOnDeploymentConfigs(testNamespace).Get(
			context.TODO(),
			addOnDeploymentConfigName,
			metav1.GetOptions{},
		)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(retrievedV1beta1.Spec.CustomizedVariables).To(gomega.HaveLen(1))
		gomega.Expect(retrievedV1beta1.Spec.CustomizedVariables[0].Name).To(gomega.Equal("alpha_var"))
		gomega.Expect(retrievedV1beta1.Spec.AgentInstallNamespace).To(gomega.Equal("alpha-namespace"))
	})

	ginkgo.It("Should validate v1beta1 constraints (invalid variable name)", func() {
		addOnDeploymentConfig := &addonv1beta1.AddOnDeploymentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      addOnDeploymentConfigName,
				Namespace: testNamespace,
			},
			Spec: addonv1beta1.AddOnDeploymentConfigSpec{
				CustomizedVariables: []addonv1beta1.CustomizedVariable{
					{
						Name:  "@invalid",
						Value: "test",
					},
				},
			},
		}

		_, err := hubAddonClient.AddonV1beta1().AddOnDeploymentConfigs(testNamespace).Create(
			context.TODO(),
			addOnDeploymentConfig,
			metav1.CreateOptions{},
		)
		gomega.Expect(errors.IsInvalid(err)).To(gomega.BeTrue())
	})
})
