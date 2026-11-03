// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

	Context("Create with aws auth and invalid arn", func() {
		It("should reject the klusterlet creation", func() {
			klusterlet.Spec.RegistrationConfiguration = &operatorv1.RegistrationConfiguration{
				RegistrationDriver: operatorv1.RegistrationDriver{
					AuthType: "awsirsa",
					AwsIrsa: &operatorv1.AwsIrsa{
						ManagedClusterArn: "arn:aws:bks:us-west-2:123456789012:cluster/managed-cluster1",
						HubClusterArn:     "arn:aws:eks:us-west-2:123456789012:cluster/hub-cluster1",
					},
				},
			}
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).NotTo(BeNil())
		})
	})

	Context("Create with aws auth and valid arn", func() {
		It("should create successfully", func() {
			klusterlet.Spec.RegistrationConfiguration = &operatorv1.RegistrationConfiguration{
				RegistrationDriver: operatorv1.RegistrationDriver{
					AuthType: "awsirsa",
					AwsIrsa: &operatorv1.AwsIrsa{
						ManagedClusterArn: "arn:aws:eks:us-west-2:123456789012:cluster/managed-cluster1",
						HubClusterArn:     "arn:aws:eks:us-west-2:123456789012:cluster/hub-cluster1",
					},
				},
			}
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})

	Context("Create with aws auth and a valid arn in a non-commercial partition", func() {
		It("should create successfully for aws-us-gov", func() {
			klusterlet.Spec.RegistrationConfiguration = awsIrsaRegistrationConfig(
				"arn:aws-us-gov:eks:us-gov-west-1:123456789012:cluster/managed-cluster1",
				"arn:aws-us-gov:eks:us-gov-west-1:123456789012:cluster/hub-cluster1")
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})

		It("should create successfully for aws-iso-b", func() {
			klusterlet.Spec.RegistrationConfiguration = awsIrsaRegistrationConfig(
				"arn:aws-iso-b:eks:us-isob-east-1:123456789012:cluster/managed-cluster1",
				"arn:aws-iso-b:eks:us-isob-east-1:123456789012:cluster/hub-cluster1")
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})

		// The schema validates each arn on its own; it has no cross-field rule tying the two
		// partitions together. IAM trust relationships do not actually cross partitions, so this
		// pins the validation boundary, not a supported topology.
		It("should not constrain the hub and the managed cluster to the same partition", func() {
			klusterlet.Spec.RegistrationConfiguration = awsIrsaRegistrationConfig(
				"arn:aws-cn:eks:cn-north-1:123456789012:cluster/managed-cluster1",
				"arn:aws:eks:us-west-2:123456789012:cluster/hub-cluster1")
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})

	Context("Create with aws auth and an arn outside the aws partitions", func() {
		It("should reject a partition that is not an aws partition", func() {
			klusterlet.Spec.RegistrationConfiguration = awsIrsaRegistrationConfig(
				"arn:notaws:eks:us-west-2:123456789012:cluster/managed-cluster1",
				"arn:aws:eks:us-west-2:123456789012:cluster/hub-cluster1")
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
		})

		It("should reject a trailing hyphen in the partition", func() {
			klusterlet.Spec.RegistrationConfiguration = awsIrsaRegistrationConfig(
				"arn:aws-:eks:us-west-2:123456789012:cluster/managed-cluster1",
				"arn:aws:eks:us-west-2:123456789012:cluster/hub-cluster1")
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
		})
	})

	Context("Create with aws auth and a credential source", func() {
		BeforeEach(func() {
			klusterlet.Spec.RegistrationConfiguration = awsIrsaRegistrationConfig(
				"arn:aws:eks:us-west-2:123456789012:cluster/managed-cluster1",
				"arn:aws:eks:us-west-2:123456789012:cluster/hub-cluster1")
		})

		It("should default to IRSA", func() {
			created, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
			Expect(created.Spec.RegistrationConfiguration.RegistrationDriver.AwsIrsa.CredentialSource).To(Equal(operatorv1.AwsCredentialSourceIRSA))
		})

		It("should accept PodIdentity", func() {
			klusterlet.Spec.RegistrationConfiguration.RegistrationDriver.AwsIrsa.CredentialSource = operatorv1.AwsCredentialSourcePodIdentity
			created, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
			Expect(created.Spec.RegistrationConfiguration.RegistrationDriver.AwsIrsa.CredentialSource).To(Equal(operatorv1.AwsCredentialSourcePodIdentity))
		})

		It("should reject an unknown value", func() {
			klusterlet.Spec.RegistrationConfiguration.RegistrationDriver.AwsIrsa.CredentialSource = "Bogus"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
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

var _ = Describe("Klusterlet API test with AgentConfiguration", func() {
	var klusterlet *operatorv1.Klusterlet

	BeforeEach(func() {
		klusterlet = &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("klusterlet-agentconfig-%s", rand.String(5)),
			},
		}
	})

	AfterEach(func() {
		err := operatorClient.OperatorV1().Klusterlets().Delete(context.TODO(), klusterlet.Name, metav1.DeleteOptions{})
		if !apierrors.IsNotFound(err) {
			Expect(err).ToNot(HaveOccurred())
		}
	})

	withBindConfiguration := func(bind *operatorv1.AgentBindConfiguration) {
		klusterlet.Spec.AgentConfiguration = &operatorv1.AgentConfiguration{BindConfiguration: bind}
	}

	It("should default the health probe port and keep a zero metrics port", func() {
		withBindConfiguration(&operatorv1.AgentBindConfiguration{})
		created, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(created.Spec.AgentConfiguration.BindConfiguration.HealthProbePort).Should(Equal(int32(8443)))
		Expect(created.Spec.AgentConfiguration.BindConfiguration.MetricsPort).Should(Equal(int32(0)))
		Expect(created.Spec.AgentConfiguration.BindConfiguration.HostNetwork).Should(BeFalse())
	})

	It("should accept custom ports with host network", func() {
		withBindConfiguration(&operatorv1.AgentBindConfiguration{
			HealthProbePort: 9443,
			MetricsPort:     9090,
			HostNetwork:     true,
		})
		created, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(created.Spec.AgentConfiguration.BindConfiguration.HealthProbePort).Should(Equal(int32(9443)))
		Expect(created.Spec.AgentConfiguration.BindConfiguration.MetricsPort).Should(Equal(int32(9090)))
		Expect(created.Spec.AgentConfiguration.BindConfiguration.HostNetwork).Should(BeTrue())
	})

	It("should reject a privileged health probe port", func() {
		withBindConfiguration(&operatorv1.AgentBindConfiguration{HealthProbePort: 80})
		_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).To(HaveOccurred())
	})

	It("should accept a separate bind configuration per agent", func() {
		klusterlet.Spec.AgentConfiguration = &operatorv1.AgentConfiguration{
			RegistrationAgentConfiguration: &operatorv1.AgentComponentConfiguration{
				BindConfiguration: &operatorv1.AgentBindConfiguration{HealthProbePort: 9443, HostNetwork: true},
			},
			WorkAgentConfiguration: &operatorv1.AgentComponentConfiguration{
				BindConfiguration: &operatorv1.AgentBindConfiguration{HealthProbePort: 9444, HostNetwork: true},
			},
		}
		created, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(created.Spec.AgentConfiguration.BindConfiguration).Should(BeNil())

		registration := created.Spec.AgentConfiguration.RegistrationAgentConfiguration.BindConfiguration
		Expect(registration.HealthProbePort).Should(Equal(int32(9443)))
		Expect(registration.HostNetwork).Should(BeTrue())

		work := created.Spec.AgentConfiguration.WorkAgentConfiguration.BindConfiguration
		Expect(work.HealthProbePort).Should(Equal(int32(9444)))
		Expect(work.HostNetwork).Should(BeTrue())
	})

	It("should default the health probe port of a per-agent bind configuration", func() {
		klusterlet.Spec.AgentConfiguration = &operatorv1.AgentConfiguration{
			WorkAgentConfiguration: &operatorv1.AgentComponentConfiguration{
				BindConfiguration: &operatorv1.AgentBindConfiguration{},
			},
		}
		created, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(created.Spec.AgentConfiguration.WorkAgentConfiguration.BindConfiguration.HealthProbePort).Should(Equal(int32(8443)))
	})

	It("should reject a privileged health probe port in a per-agent bind configuration", func() {
		klusterlet.Spec.AgentConfiguration = &operatorv1.AgentConfiguration{
			RegistrationAgentConfiguration: &operatorv1.AgentComponentConfiguration{
				BindConfiguration: &operatorv1.AgentBindConfiguration{HealthProbePort: 80},
			},
		}
		_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Klusterlet v1 Enhanced API test", func() {
	var klusterletName string

	BeforeEach(func() {
		suffix := rand.String(5)
		klusterletName = fmt.Sprintf("klusterlet-enhanced-%s", suffix)
	})

	AfterEach(func() {
		err := operatorClient.OperatorV1().Klusterlets().Delete(context.TODO(), klusterletName, metav1.DeleteOptions{})
		if !apierrors.IsForbidden(err) {
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("Klusterlet comprehensive configuration validation", func() {
		It("should handle complete configuration with all optional fields", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{
					RegistrationImagePullSpec: "quay.io/test/registration:latest",
					WorkImagePullSpec:         "quay.io/test/work:latest",
					ClusterName:               "test-cluster",
					Namespace:                 "open-cluster-management-agent",
					ExternalServerURLs: []operatorv1.ServerURL{
						{
							URL: "https://hub.example.com:6443",
						},
					},
					NodePlacement: operatorv1.NodePlacement{
						NodeSelector: map[string]string{
							"node-role.kubernetes.io/worker": "",
						},
						Tolerations: []v1.Toleration{
							{
								Key:      "node-role.kubernetes.io/worker",
								Operator: v1.TolerationOpExists,
								Effect:   v1.TaintEffectNoSchedule,
							},
						},
					},
					DeployOption: operatorv1.KlusterletDeployOption{
						Mode: operatorv1.InstallModeDefault,
					},
					RegistrationConfiguration: &operatorv1.RegistrationConfiguration{
						FeatureGates: []operatorv1.FeatureGate{
							{
								Feature: "AddonManagement",
								Mode:    operatorv1.FeatureGateModeTypeEnable,
							},
						},
					},
					WorkConfiguration: &operatorv1.WorkAgentConfiguration{
						FeatureGates: []operatorv1.FeatureGate{
							{
								Feature: "ManifestWorkReplicaSet",
								Mode:    operatorv1.FeatureGateModeTypeEnable,
							},
						},
					},
				},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdKlusterlet.Spec.ClusterName).Should(Equal("test-cluster"))
			Expect(createdKlusterlet.Spec.Namespace).Should(Equal("open-cluster-management-agent"))
			Expect(len(createdKlusterlet.Spec.ExternalServerURLs)).Should(Equal(1))
			Expect(createdKlusterlet.Spec.NodePlacement.NodeSelector["node-role.kubernetes.io/worker"]).Should(Equal(""))
			Expect(len(createdKlusterlet.Spec.NodePlacement.Tolerations)).Should(Equal(1))
		})

		It("should validate hosted mode configuration", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{
					DeployOption: operatorv1.KlusterletDeployOption{
						Mode: operatorv1.InstallModeHosted,
					},
				},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdKlusterlet.Spec.DeployOption.Mode).Should(Equal(operatorv1.InstallModeHosted))
		})

		It("should validate priority class configuration", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{
					PriorityClassName: "system-cluster-critical",
				},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdKlusterlet.Spec.PriorityClassName).Should(Equal("system-cluster-critical"))
		})
	})

	Context("Klusterlet resource requirements", func() {
		It("should handle resource requirements configuration", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{
					ResourceRequirement: &operatorv1.ResourceRequirement{
						Type: operatorv1.ResourceQosClassResourceRequirement,
					},
				},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdKlusterlet.Spec.ResourceRequirement.Type).Should(Equal(operatorv1.ResourceQosClassResourceRequirement))
		})
	})

	Context("Klusterlet status updates", func() {
		It("should allow status updates", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())

			// Update status
			createdKlusterlet.Status = operatorv1.KlusterletStatus{
				ObservedGeneration: 1,
				Conditions: []metav1.Condition{
					{
						Type:               "Applied",
						Status:             metav1.ConditionTrue,
						Reason:             "KlusterletDeployed",
						LastTransitionTime: metav1.Now(),
					},
				},
				Generations: []operatorv1.GenerationStatus{
					{
						Group:          "apps",
						Version:        "v1",
						Resource:       "deployments",
						Namespace:      "open-cluster-management-agent",
						Name:           "klusterlet-registration-agent",
						LastGeneration: 1,
					},
				},
				RelatedResources: []operatorv1.RelatedResourceMeta{
					{
						Group:     "apps",
						Version:   "v1",
						Resource:  "deployments",
						Namespace: "open-cluster-management-agent",
						Name:      "klusterlet-registration-agent",
					},
				},
			}

			_, err = operatorClient.OperatorV1().Klusterlets().UpdateStatus(context.TODO(), createdKlusterlet, metav1.UpdateOptions{})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

// awsIrsaRegistrationConfig builds an awsirsa registration config for the given cluster arns,
// so the partition cases above differ only by the arns under test.
func awsIrsaRegistrationConfig(managedClusterArn, hubClusterArn string) *operatorv1.RegistrationConfiguration {
	return &operatorv1.RegistrationConfiguration{
		RegistrationDriver: operatorv1.RegistrationDriver{
			AuthType: "awsirsa",
			AwsIrsa: &operatorv1.AwsIrsa{
				ManagedClusterArn: managedClusterArn,
				HubClusterArn:     hubClusterArn,
			},
		},
	}
}
