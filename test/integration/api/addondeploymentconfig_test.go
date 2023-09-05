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
})
