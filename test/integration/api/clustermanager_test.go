// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
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
				RegistrationWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
					Address: "test:test",
				},
				WorkWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
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
				RegistrationWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
					Address: "192.168.2.3",
				},
				WorkWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
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
				RegistrationWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
					Address: "localhost",
				},
				WorkWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
					Address: "foo.com",
				},
			}
			_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})

	Context("Set nothing in ports", func() {
		It("should have 443 as default value in hosted mode", func() {
			clusterManager.Spec.DeployOption.Hosted = &operatorv1.HostedClusterManagerConfiguration{
				RegistrationWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
					Address: "localhost",
				},
				WorkWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
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
				RegistrationWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
					Address: "localhost",
					Port:    65536,
				},
				WorkWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
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
				RegistrationWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
					Address: "foo1.com",
					Port:    1443,
				},
				WorkWebhookConfiguration: operatorv1.HostedWebhookConfiguration{
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

	It("Create a cluster manager with aws registration and invalid hubClusterArn", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				RegistrationConfiguration: &operatorv1.RegistrationHubConfiguration{
					RegistrationDrivers: []operatorv1.RegistrationDriverHub{
						{
							AuthType: "awsirsa",
							AwsIrsa: &operatorv1.AwsIrsaConfig{
								HubClusterArn: "arn:aws:bks:us-west-2:123456789012:cluster/hub-cluster1",
							},
						},
					},
				},
			},
		}
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).ToNot(BeNil())
	})

	It("Create a cluster manager with aws registration and valid hubClusterArn", func() {
		clusterManager := &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{
				Name: clusterManagerName,
			},
			Spec: operatorv1.ClusterManagerSpec{
				RegistrationConfiguration: &operatorv1.RegistrationHubConfiguration{
					RegistrationDrivers: []operatorv1.RegistrationDriverHub{
						{
							AuthType: "awsirsa",
							AwsIrsa: &operatorv1.AwsIrsaConfig{
								HubClusterArn: "arn:aws:eks:us-west-2:123456789012:cluster/hub-cluster1",
							},
						},
					},
				},
			},
		}
		_, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
		Expect(err).To(BeNil())
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

var _ = Describe("ClusterManager v1 Enhanced API test", func() {
	var clusterManagerName string

	BeforeEach(func() {
		suffix := rand.String(5)
		clusterManagerName = fmt.Sprintf("cm-enhanced-%s", suffix)
	})

	AfterEach(func() {
		err := operatorClient.OperatorV1().ClusterManagers().Delete(context.TODO(), clusterManagerName, metav1.DeleteOptions{})
		if !apierrors.IsForbidden(err) {
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("ClusterManager comprehensive configuration validation", func() {
		It("should handle complete configuration with all optional fields", func() {
			clusterManager := &operatorv1.ClusterManager{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterManagerName,
				},
				Spec: operatorv1.ClusterManagerSpec{
					RegistrationImagePullSpec: "quay.io/test/registration:latest",
					WorkImagePullSpec:         "quay.io/test/work:latest",
					PlacementImagePullSpec:    "quay.io/test/placement:latest",
					AddOnManagerImagePullSpec: "quay.io/test/addon-manager:latest",
					NodePlacement: operatorv1.NodePlacement{
						NodeSelector: map[string]string{
							"node-role.kubernetes.io/infra": "",
						},
						Tolerations: []v1.Toleration{
							{
								Key:      "node-role.kubernetes.io/infra",
								Operator: v1.TolerationOpExists,
								Effect:   v1.TaintEffectNoSchedule,
							},
						},
					},
					DeployOption: operatorv1.ClusterManagerDeployOption{
						Mode: operatorv1.InstallModeDefault,
					},
					RegistrationConfiguration: &operatorv1.RegistrationHubConfiguration{
						AutoApproveUsers: []string{"system:admin"},
						FeatureGates: []operatorv1.FeatureGate{
							{
								Feature: "DefaultClusterSet",
								Mode:    operatorv1.FeatureGateModeTypeEnable,
							},
						},
					},
					WorkConfiguration: &operatorv1.WorkConfiguration{
						WorkDriver: operatorv1.WorkDriverTypeKube,
						FeatureGates: []operatorv1.FeatureGate{
							{
								Feature: "ManifestWorkReplicaSet",
								Mode:    operatorv1.FeatureGateModeTypeEnable,
							},
						},
					},
				},
			}

			createdClusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdClusterManager.Spec.NodePlacement.NodeSelector["node-role.kubernetes.io/infra"]).Should(Equal(""))
			Expect(len(createdClusterManager.Spec.NodePlacement.Tolerations)).Should(Equal(1))
			Expect(createdClusterManager.Spec.RegistrationConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeEnable))
			Expect(createdClusterManager.Spec.WorkConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeEnable))
		})

		It("should validate addon manager configuration", func() {
			clusterManager := &operatorv1.ClusterManager{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterManagerName,
				},
				Spec: operatorv1.ClusterManagerSpec{
					AddOnManagerConfiguration: &operatorv1.AddOnManagerConfiguration{
						FeatureGates: []operatorv1.FeatureGate{
							{
								Feature: "AddonManagement",
								Mode:    operatorv1.FeatureGateModeTypeEnable,
							},
						},
					},
				},
			}

			createdClusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdClusterManager.Spec.AddOnManagerConfiguration.FeatureGates[0].Feature).Should(Equal("AddonManagement"))
		})

		It("should validate server configuration", func() {
			clusterManager := &operatorv1.ClusterManager{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterManagerName,
				},
				Spec: operatorv1.ClusterManagerSpec{
					ServerConfiguration: &operatorv1.ServerConfiguration{},
				},
			}

			createdClusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdClusterManager.Spec.ServerConfiguration).ShouldNot(BeNil())
		})
	})

	Context("ClusterManager resource requirements", func() {
		It("should handle resource requirements configuration", func() {
			clusterManager := &operatorv1.ClusterManager{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterManagerName,
				},
				Spec: operatorv1.ClusterManagerSpec{
					ResourceRequirement: &operatorv1.ResourceRequirement{
						Type: operatorv1.ResourceQosClassResourceRequirement,
					},
				},
			}

			createdClusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdClusterManager.Spec.ResourceRequirement.Type).Should(Equal(operatorv1.ResourceQosClassResourceRequirement))
		})
	})

	Context("ClusterManager status updates", func() {
		It("should allow status updates", func() {
			clusterManager := &operatorv1.ClusterManager{
				ObjectMeta: metav1.ObjectMeta{
					Name: clusterManagerName,
				},
				Spec: operatorv1.ClusterManagerSpec{},
			}

			createdClusterManager, err := operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), clusterManager, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())

			// Update status
			createdClusterManager.Status = operatorv1.ClusterManagerStatus{
				ObservedGeneration: 1,
				Conditions: []metav1.Condition{
					{
						Type:               "Applied",
						Status:             metav1.ConditionTrue,
						Reason:             "ClusterManagerDeployed",
						LastTransitionTime: metav1.Now(),
					},
				},
				Generations: []operatorv1.GenerationStatus{
					{
						Group:          "apps",
						Version:        "v1",
						Resource:       "deployments",
						Namespace:      "open-cluster-management-hub",
						Name:           "cluster-manager-registration-controller",
						LastGeneration: 1,
					},
				},
				RelatedResources: []operatorv1.RelatedResourceMeta{
					{
						Group:     "apps",
						Version:   "v1",
						Resource:  "deployments",
						Namespace: "open-cluster-management-hub",
						Name:      "cluster-manager-registration-controller",
					},
				},
			}

			_, err = operatorClient.OperatorV1().ClusterManagers().UpdateStatus(context.TODO(), createdClusterManager, metav1.UpdateOptions{})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
