package api

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	operatorv1 "open-cluster-management.io/api/operator/v1"
)

var _ = Describe("Create Klusterlet API", func() {
	var klusterlet *operatorv1.Klusterlet
	BeforeEach(func() {
		suffix := rand.String(5)
		klusterletName := fmt.Sprintf("cm-%s", suffix)
		klusterlet = &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{},
		}
	})

	Context("Create without nothing set", func() {
		It("should create successfully", func() {
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})

	Context("Create with invalid namespace", func() {
		It("should reject the klusterlet creation", func() {
			klusterlet.Spec.Namespace = "invalid-klusterlet-ns"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).NotTo(BeNil())
		})
	})
})

var _ = Describe("valid HubApiServerHostAlias", func() {
	var klusterlet *operatorv1.Klusterlet

	BeforeEach(func() {
		suffix := rand.String(5)
		klusterletName := fmt.Sprintf("cm-%s", suffix)
		klusterlet = &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{
				HubApiServerHostAlias: &operatorv1.HubApiServerHostAlias{},
			},
		}
	})

	Context("Empty IPV4 address", func() {
		It("should return err", func() {
			klusterlet.Spec.HubApiServerHostAlias.Hostname = "xxx.yyy.zzz"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Empty hostname", func() {
		It("should return err", func() {
			klusterlet.Spec.HubApiServerHostAlias.IP = "1.2.3.4"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Invalid IPV4 address and hostname", func() {
		It("should return err", func() {
			klusterlet.Spec.HubApiServerHostAlias.IP = "1.2.3.257"
			klusterlet.Spec.HubApiServerHostAlias.Hostname = "xxxyyyzzz"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Valid IPV4 address and hostname", func() {
		It("should create successfully", func() {
			klusterlet.Spec.HubApiServerHostAlias.IP = "1.2.3.4"
			klusterlet.Spec.HubApiServerHostAlias.Hostname = "xxx.yyy.zzz"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})
})

var _ = Describe("Klusterlet API test with WorkConfiguration", func() {
	var klusterletName string

	BeforeEach(func() {
		suffix := rand.String(5)
		klusterletName = fmt.Sprintf("cm-%s", suffix)
	})

	It("Create a klusterlet with empty spec", func() {
		klusterlet := &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{},
		}
		klusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		Expect(klusterlet.Spec.WorkConfiguration).To(BeNil())
	})

	It("Create a klusterlet with empty work feature gate mode", func() {
		klusterlet := &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{
				WorkConfiguration: &operatorv1.WorkAgentConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
						},
					},
				},
			},
		}
		klusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(klusterlet.Spec.WorkConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeDisable))
	})

	It("Create a klusterlet with wrong work feature gate mode", func() {
		klusterlet := &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{
				WorkConfiguration: &operatorv1.WorkAgentConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
							Mode:    "WrongMode",
						},
					},
				},
			},
		}
		_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).To(HaveOccurred())
	})

	It("Create a klusterlet with right work feature gate mode", func() {
		klusterlet := &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{
				WorkConfiguration: &operatorv1.WorkAgentConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
							Mode:    "Disable",
						},
						{
							Feature: "Bar",
							Mode:    "Enable",
						},
					},
				},
			},
		}
		_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).To(BeNil())
		Expect(klusterlet.Spec.WorkConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeDisable))
		Expect(klusterlet.Spec.WorkConfiguration.FeatureGates[1].Mode).Should(Equal(operatorv1.FeatureGateModeTypeEnable))
	})
})
