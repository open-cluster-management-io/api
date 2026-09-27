// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"encoding/json"
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
})

var _ = Describe("Create Klusterlet API with azure registration", func() {
	const (
		azureID  = "11111111-1111-1111-1111-111111111111"
		clientID = "22222222-2222-2222-2222-222222222222"
		tenantID = "33333333-3333-3333-3333-333333333333"
	)
	var klusterlet *operatorv1.Klusterlet
	BeforeEach(func() {
		klusterlet = &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("cm-%s", rand.String(5)),
			},
		}
	})

	create := func(azure *operatorv1.AzureAuth) error {
		klusterlet.Spec.RegistrationConfiguration = azureRegistrationConfig(azure)
		_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		return err
	}

	It("should reject authType azure without the azure configuration", func() {
		Expect(apierrors.IsInvalid(create(nil))).To(BeTrue())
	})

	It("should reject a missing credential", func() {
		Expect(apierrors.IsInvalid(create(&operatorv1.AzureAuth{
			ManagedClusterAzureID: azureID,
		}))).To(BeTrue())
	})

	It("should reject an unknown credential", func() {
		Expect(apierrors.IsInvalid(create(&operatorv1.AzureAuth{
			Credential:            "default-azure-credential",
			ManagedClusterAzureID: azureID,
		}))).To(BeTrue())
	})

	It("should reject a missing managedClusterAzureID", func() {
		Expect(apierrors.IsInvalid(create(&operatorv1.AzureAuth{
			Credential: operatorv1.AzureManagedIdentityCredential,
		}))).To(BeTrue())
	})

	It("should create managed-identity-credential without clientID, for a system-assigned identity", func() {
		Expect(create(&operatorv1.AzureAuth{
			Credential:            operatorv1.AzureManagedIdentityCredential,
			ManagedClusterAzureID: azureID,
		})).To(Succeed())
	})

	It("should create managed-identity-credential with clientID, for a user-assigned identity", func() {
		Expect(create(&operatorv1.AzureAuth{
			Credential:            operatorv1.AzureManagedIdentityCredential,
			ManagedClusterAzureID: azureID,
			ClientID:              clientID,
		})).To(Succeed())
	})

	It("should create workload-identity-credential with clientID", func() {
		Expect(create(&operatorv1.AzureAuth{
			Credential:            operatorv1.AzureWorkloadIdentityCredential,
			ManagedClusterAzureID: azureID,
			ClientID:              clientID,
			FederatedTokenFile:    "/var/run/secrets/azure/tokens/azure-identity-token",
		})).To(Succeed())
	})

	It("should reject workload-identity-credential without clientID", func() {
		Expect(apierrors.IsInvalid(create(&operatorv1.AzureAuth{
			Credential:            operatorv1.AzureWorkloadIdentityCredential,
			ManagedClusterAzureID: azureID,
		}))).To(BeTrue())
	})

	It("should create the environment-credential types with clientID and tenantID", func() {
		for _, credential := range []operatorv1.AzureCredentialType{
			operatorv1.AzureEnvironmentCredentialSecret, operatorv1.AzureEnvironmentCredentialCertificate,
		} {
			klusterlet.Name = fmt.Sprintf("cm-%s", rand.String(5))
			Expect(create(&operatorv1.AzureAuth{
				Credential:            credential,
				ManagedClusterAzureID: azureID,
				ClientID:              clientID,
				TenantID:              tenantID,
			})).To(Succeed(), string(credential))
		}
	})

	It("should reject the environment-credential types without clientID or tenantID", func() {
		for _, credential := range []operatorv1.AzureCredentialType{
			operatorv1.AzureEnvironmentCredentialSecret, operatorv1.AzureEnvironmentCredentialCertificate,
		} {
			Expect(apierrors.IsInvalid(create(&operatorv1.AzureAuth{
				Credential:            credential,
				ManagedClusterAzureID: azureID,
				ClientID:              clientID,
			}))).To(BeTrue(), string(credential)+" without tenantID")
			Expect(apierrors.IsInvalid(create(&operatorv1.AzureAuth{
				Credential:            credential,
				ManagedClusterAzureID: azureID,
				TenantID:              tenantID,
			}))).To(BeTrue(), string(credential)+" without clientID")
		}
	})

	It("should reject federatedTokenFile with a credential other than workload-identity-credential", func() {
		for _, credential := range []operatorv1.AzureCredentialType{
			operatorv1.AzureManagedIdentityCredential,
			operatorv1.AzureEnvironmentCredentialSecret,
			operatorv1.AzureEnvironmentCredentialCertificate,
		} {
			err := create(&operatorv1.AzureAuth{
				Credential:            credential,
				ManagedClusterAzureID: azureID,
				ClientID:              clientID,
				TenantID:              tenantID,
				FederatedTokenFile:    "/var/run/secrets/azure/tokens/azure-identity-token",
			})
			Expect(apierrors.IsInvalid(err)).To(BeTrue(), string(credential))
			Expect(err.Error()).To(ContainSubstring("federatedTokenFile applies only to workload-identity-credential"), string(credential))
		}
	})

	It("should persist every azure field", func() {
		azure := &operatorv1.AzureAuth{
			Credential:            operatorv1.AzureWorkloadIdentityCredential,
			ManagedClusterAzureID: azureID,
			ClientID:              clientID,
			TenantID:              tenantID,
			FederatedTokenFile:    "/var/run/secrets/azure/tokens/azure-identity-token",
			TokenAudience:         "api://hub-apiserver/.default",
		}
		Expect(create(azure)).To(Succeed())

		got, err := operatorClient.OperatorV1().Klusterlets().Get(context.TODO(), klusterlet.Name, metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(got.Spec.RegistrationConfiguration.RegistrationDriver.AuthType).To(Equal(operatorv1.AzureAuthType))
		Expect(got.Spec.RegistrationConfiguration.RegistrationDriver.Azure).To(Equal(azure))
	})

	Context("Update", func() {
		// update applies mutate to the stored Klusterlet's azure configuration.
		update := func(mutate func(azure *operatorv1.AzureAuth)) error {
			got, err := operatorClient.OperatorV1().Klusterlets().Get(context.TODO(), klusterlet.Name, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			mutate(got.Spec.RegistrationConfiguration.RegistrationDriver.Azure)
			_, err = operatorClient.OperatorV1().Klusterlets().Update(context.TODO(), got, metav1.UpdateOptions{})
			return err
		}

		It("should reject switching to workload-identity-credential without clientID", func() {
			Expect(create(&operatorv1.AzureAuth{
				Credential:            operatorv1.AzureManagedIdentityCredential,
				ManagedClusterAzureID: azureID,
			})).To(Succeed())

			err := update(func(azure *operatorv1.AzureAuth) {
				azure.Credential = operatorv1.AzureWorkloadIdentityCredential
			})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("clientID is required for workload-identity-credential"))
		})

		It("should reject removing tenantID from an environment credential", func() {
			Expect(create(&operatorv1.AzureAuth{
				Credential:            operatorv1.AzureEnvironmentCredentialSecret,
				ManagedClusterAzureID: azureID,
				ClientID:              clientID,
				TenantID:              tenantID,
			})).To(Succeed())

			err := update(func(azure *operatorv1.AzureAuth) {
				azure.TenantID = ""
			})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("clientID and tenantID are required"))
		})

		It("should accept switching to workload-identity-credential with clientID", func() {
			Expect(create(&operatorv1.AzureAuth{
				Credential:            operatorv1.AzureManagedIdentityCredential,
				ManagedClusterAzureID: azureID,
			})).To(Succeed())

			Expect(update(func(azure *operatorv1.AzureAuth) {
				azure.Credential = operatorv1.AzureWorkloadIdentityCredential
				azure.ClientID = clientID
			})).To(Succeed())
		})
	})

	// The typed client drops empty optional fields (omitempty) and always sends required
	// fields, so explicitly-empty and absent values - what a YAML or Helm user can send -
	// are only reachable with a raw request.
	Context("Raw requests", func() {
		createRaw := func(azure map[string]interface{}) error {
			body, err := json.Marshal(map[string]interface{}{
				"apiVersion": "operator.open-cluster-management.io/v1",
				"kind":       "Klusterlet",
				"metadata":   map[string]interface{}{"name": klusterlet.Name},
				"spec": map[string]interface{}{
					"registrationConfiguration": map[string]interface{}{
						"registrationDriver": map[string]interface{}{
							"authType": operatorv1.AzureAuthType,
							"azure":    azure,
						},
					},
				},
			})
			Expect(err).ToNot(HaveOccurred())
			return operatorClient.OperatorV1().RESTClient().Post().
				Resource("klusterlets").Body(body).Do(context.TODO()).Error()
		}

		It("should create a valid klusterlet, so the rejections below are not an artifact of the raw request", func() {
			Expect(createRaw(map[string]interface{}{
				"credential":            string(operatorv1.AzureWorkloadIdentityCredential),
				"managedClusterAzureID": azureID,
				"clientID":              clientID,
			})).To(Succeed())
		})

		It("should reject an explicitly empty clientID for workload-identity-credential", func() {
			err := createRaw(map[string]interface{}{
				"credential":            string(operatorv1.AzureWorkloadIdentityCredential),
				"managedClusterAzureID": azureID,
				"clientID":              "",
			})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("clientID is required for workload-identity-credential"))
		})

		It("should reject an explicitly empty clientID or tenantID for the environment credentials", func() {
			for _, credential := range []operatorv1.AzureCredentialType{
				operatorv1.AzureEnvironmentCredentialSecret, operatorv1.AzureEnvironmentCredentialCertificate,
			} {
				for field, azure := range map[string]map[string]interface{}{
					"clientID": {"credential": string(credential), "managedClusterAzureID": azureID, "clientID": "", "tenantID": tenantID},
					"tenantID": {"credential": string(credential), "managedClusterAzureID": azureID, "clientID": clientID, "tenantID": ""},
				} {
					err := createRaw(azure)
					Expect(apierrors.IsInvalid(err)).To(BeTrue(), string(credential)+" with empty "+field)
					Expect(err.Error()).To(ContainSubstring("clientID and tenantID are required"), string(credential)+" with empty "+field)
				}
			}
		})

		It("should reject an absent credential", func() {
			err := createRaw(map[string]interface{}{
				"managedClusterAzureID": azureID,
			})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("credential: Required value"))
		})

		It("should reject an absent managedClusterAzureID", func() {
			err := createRaw(map[string]interface{}{
				"credential": string(operatorv1.AzureManagedIdentityCredential),
			})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("managedClusterAzureID: Required value"))
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

func azureRegistrationConfig(azure *operatorv1.AzureAuth) *operatorv1.RegistrationConfiguration {
	return &operatorv1.RegistrationConfiguration{
		RegistrationDriver: operatorv1.RegistrationDriver{
			AuthType: operatorv1.AzureAuthType,
			Azure:    azure,
		},
	}
}
