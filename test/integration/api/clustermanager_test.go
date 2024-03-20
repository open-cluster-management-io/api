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

var _ = Describe("ClusterManager API test", func() {
	var clusterManagerName string

	BeforeEach(func() {
		suffix := rand.String(5)
		clusterManagerName = fmt.Sprintf("cm-%s", suffix)
	})

	It("Create a cluster manager with empty spec", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{},
		}
		clusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		Expect(clusterManager.Spec.DeployOption.Mode).Should(Equal(operatorv1.InstallModeDefault))
	})

	It("Create a cluster manager with wrong install mode", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				DeployOption: operatorv1.ClusterManagerDeployOption{
					Mode: "WrongMode",
				},
			},
		}
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Create Cluster Manager Hosted mode", func() {
	var clusterManager *operatorv1.ClusterManager

	BeforeEach(func() {
		suffix := rand.String(5)
		clusterManagerName := fmt.Sprintf("cm-%s", suffix)
		clusterManager = &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				DeployOption: operatorv1.ClusterManagerDeployOption{
					Mode:   operatorv1.InstallModeHosted,
					Hosted: &operatorv1.HostedClusterManagerConfiguration{},
				},
			},
		}
	})

	Context("Set nothing in addresses", func() {
		It("should return err", func() {
			_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Set wrong format address", func() {
		It("should return err", func() {
			clusterManager.Spec.DeployOption.Hosted = &operatorv1.HostedClusterManagerConfiguration{
				RegistrationWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "test:test",
				},
				WorkWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "test:test",
				},
			}
			_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Set IPV4 format addresses", func() {
		It("should create successfully", func() {
			clusterManager.Spec.DeployOption.Hosted = &operatorv1.HostedClusterManagerConfiguration{
				RegistrationWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "192.168.2.3",
				},
				WorkWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "192.168.2.4",
				},
			}
			_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})

	Context("Set FQDN format addresses", func() {
		It("should create successfully", func() {
			clusterManager.Spec.DeployOption.Hosted = &operatorv1.HostedClusterManagerConfiguration{
				RegistrationWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "localhost",
				},
				WorkWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "foo.com",
				},
			}
			_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})

	Context("Set nothing in ports", func() {
		It("should has 443 as default value", func() {
			clusterManager.Spec.DeployOption.Hosted = &operatorv1.HostedClusterManagerConfiguration{
				RegistrationWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "localhost",
				},
				WorkWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "foo.com",
				},
			}
			clusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).To(BeNil())
			Expect(clusterManager.Spec.DeployOption.Hosted.RegistrationWebhookConfiguration.Port).Should(Equal(int32(443)))
			Expect(clusterManager.Spec.DeployOption.Hosted.WorkWebhookConfiguration.Port).Should(Equal(int32(443)))
		})
	})

	Context("Set port bigger than 65535", func() {
		It("should return err", func() {
			clusterManager.Spec.DeployOption.Hosted = &operatorv1.HostedClusterManagerConfiguration{
				RegistrationWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "localhost",
					Port:    65536,
				},
				WorkWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "foo.com",
				},
			}
			_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Set customized WebhookConfiguration", func() {
		It("should have euqually value after create", func() {
			clusterManager.Spec.DeployOption.Hosted = &operatorv1.HostedClusterManagerConfiguration{
				RegistrationWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "foo1.com",
					Port:    1443,
				},
				WorkWebhookConfiguration: operatorv1.WebhookConfiguration{
					Address: "foo2.com",
					Port:    2443,
				},
			}

			clusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).To(BeNil())
			Expect(clusterManager.Spec.DeployOption.Hosted.RegistrationWebhookConfiguration.Address).Should(Equal("foo1.com"))
			Expect(clusterManager.Spec.DeployOption.Hosted.RegistrationWebhookConfiguration.Port).Should(Equal(int32(1443)))
			Expect(clusterManager.Spec.DeployOption.Hosted.WorkWebhookConfiguration.Address).Should(Equal("foo2.com"))
			Expect(clusterManager.Spec.DeployOption.Hosted.WorkWebhookConfiguration.Port).Should(Equal(int32(2443)))
		})
	})
})

var _ = Describe("ClusterManager API test with RegistrationConfiguration", func() {
	var clusterManagerName string

	BeforeEach(func() {
		suffix := rand.String(5)
		clusterManagerName = fmt.Sprintf("cm-%s", suffix)
	})

	It("Create a cluster manager with empty spec", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{},
		}
		clusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		Expect(clusterManager.Spec.RegistrationConfiguration).To(BeNil())
	})

	It("Create a cluster manager with empty registration feature gate mode", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				RegistrationConfiguration: &operatorv1.RegistrationHubConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
						},
					},
				},
			},
		}
		clusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(clusterManager.Spec.RegistrationConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeDisable))
	})

	It("Create a cluster manager with wrong registration feature gate mode", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				RegistrationConfiguration: &operatorv1.RegistrationHubConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
							Mode:    "WrongMode",
						},
					},
				},
			},
		}
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).To(HaveOccurred())
	})

	It("Create a cluster manager with right registration feature gate mode", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				RegistrationConfiguration: &operatorv1.RegistrationHubConfiguration{
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
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).To(BeNil())
		Expect(clusterManager.Spec.RegistrationConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeDisable))
		Expect(clusterManager.Spec.RegistrationConfiguration.FeatureGates[1].Mode).Should(Equal(operatorv1.FeatureGateModeTypeEnable))
	})
})

var _ = Describe("ClusterManager API test with WorkConfiguration", func() {
	var clusterManagerName string

	BeforeEach(func() {
		suffix := rand.String(5)
		clusterManagerName = fmt.Sprintf("cm-%s", suffix)
	})

	It("Create a cluster manager with empty spec", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{},
		}
		clusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		Expect(clusterManager.Spec.WorkConfiguration.WorkDriver).Should(Equal(operatorv1.WorkDriverTypeKube))
	})

	It("Create a cluster manager with wrong driver type", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				WorkConfiguration: &operatorv1.WorkConfiguration{
					WorkDriver: "WrongDriver",
				},
			},
		}
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).To(HaveOccurred())
	})

	It("Create a cluster manager with empty work feature gate mode", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				WorkConfiguration: &operatorv1.WorkConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
						},
					},
				},
			},
		}
		clusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(clusterManager.Spec.WorkConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeDisable))
	})

	It("Create a cluster manager with wrong work feature gate mode", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				WorkConfiguration: &operatorv1.WorkConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
							Mode:    "WrongMode",
						},
					},
				},
			},
		}
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).To(HaveOccurred())
	})

	It("Create a cluster manager with right work feature gate mode", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				WorkConfiguration: &operatorv1.WorkConfiguration{
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
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).To(BeNil())
		Expect(clusterManager.Spec.WorkConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeDisable))
		Expect(clusterManager.Spec.WorkConfiguration.FeatureGates[1].Mode).Should(Equal(operatorv1.FeatureGateModeTypeEnable))
	})
})
