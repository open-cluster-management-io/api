// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	addonv1beta1 "open-cluster-management.io/api/addon/v1beta1"
	operatorv1 "open-cluster-management.io/api/operator/v1"
)

const (
	testHubEndpoint = "ocm-grpc.hub.internal:443"
	testHubSPIFFEID = "spiffe://ocm-hub.example.org/ocm/hub/grpc-server"
)

var _ = Describe("Klusterlet API test with the spiffe registration driver", func() {
	var klusterletName string

	BeforeEach(func() {
		klusterletName = fmt.Sprintf("klusterlet-%s", rand.String(5))
	})

	newKlusterlet := func(rc *operatorv1.RegistrationConfiguration) *operatorv1.Klusterlet {
		return &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{Name: klusterletName},
			Spec: operatorv1.KlusterletSpec{
				ClusterName:               "cluster1",
				RegistrationConfiguration: rc,
			},
		}
	}

	spiffeDriver := func() operatorv1.RegistrationDriver {
		return operatorv1.RegistrationDriver{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE: &operatorv1.KlusterletSPIFFEConfig{
				HubEndpoint: testHubEndpoint,
				HubSPIFFEID: testHubSPIFFEID,
			},
		}
	}

	spiffeAddOnDriver := func() *operatorv1.AddOnRegistrationDriver {
		return &operatorv1.AddOnRegistrationDriver{AuthType: operatorv1.SPIFFEAuthType}
	}

	create := func(k *operatorv1.Klusterlet) (*operatorv1.Klusterlet, error) {
		return operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), k, metav1.CreateOptions{})
	}

	It("accepts spiffe with the spiffe add-on driver and fills the defaults", func() {
		k, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                spiffeDriver(),
			AddOnKubeClientRegistrationDriver: spiffeAddOnDriver(),
		}))
		Expect(err).ToNot(HaveOccurred())

		s := k.Spec.RegistrationConfiguration.RegistrationDriver.SPIFFE
		Expect(s.HubEndpoint).Should(Equal(testHubEndpoint))
		Expect(s.HubSPIFFEID).Should(Equal(testHubSPIFFEID))
		Expect(s.WorkloadAPISocket).Should(Equal("unix:///spiffe-workload-api/spire-agent.sock"))
		Expect(k.Spec.RegistrationConfiguration.AddOnKubeClientRegistrationDriver.AuthType).Should(Equal("spiffe"))
	})

	It("accepts agents in a trust domain other than the hub's", func() {
		d := spiffeDriver()
		d.SPIFFE.TrustDomain = "east.example.org"
		k, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                d,
			AddOnKubeClientRegistrationDriver: spiffeAddOnDriver(),
		}))
		Expect(err).ToNot(HaveOccurred())
		Expect(k.Spec.RegistrationConfiguration.RegistrationDriver.SPIFFE.TrustDomain).Should(Equal("east.example.org"))
		Expect(k.Spec.RegistrationConfiguration.RegistrationDriver.SPIFFE.HubSPIFFEID).Should(Equal(testHubSPIFFEID))
	})

	withMode := func(mode operatorv1.InstallMode, authType string) *operatorv1.Klusterlet {
		rc := &operatorv1.RegistrationConfiguration{RegistrationDriver: operatorv1.RegistrationDriver{AuthType: authType}}
		if authType == operatorv1.SPIFFEAuthType {
			rc = &operatorv1.RegistrationConfiguration{RegistrationDriver: spiffeDriver(),
				AddOnKubeClientRegistrationDriver: spiffeAddOnDriver()}
		}
		k := newKlusterlet(rc)
		k.Spec.DeployOption.Mode = mode
		return k
	}

	It("accepts spiffe in the Default and Singleton modes, also with the mode left empty", func() {
		for i, mode := range []operatorv1.InstallMode{"", operatorv1.InstallModeDefault, operatorv1.InstallModeSingleton} {
			k := withMode(mode, operatorv1.SPIFFEAuthType)
			k.Name = fmt.Sprintf("%s-%d", klusterletName, i)
			_, err := create(k)
			Expect(err).ToNot(HaveOccurred(), "mode %q", mode)
		}
	})

	It("refuses spiffe in the Hosted and SingletonHosted modes", func() {
		for _, mode := range []operatorv1.InstallMode{operatorv1.InstallModeHosted, operatorv1.InstallModeSingletonHosted} {
			_, err := create(withMode(mode, operatorv1.SPIFFEAuthType))
			Expect(err).To(HaveOccurred(), "mode %q", mode)
			Expect(err.Error()).To(ContainSubstring("supports only the Default and Singleton install modes"))
		}
	})

	It("still accepts other drivers in the hosted modes", func() {
		_, err := create(withMode(operatorv1.InstallModeHosted, "csr"))
		Expect(err).ToNot(HaveOccurred())
	})

	It("refuses moving a spiffe Klusterlet to a hosted mode", func() {
		k, err := create(withMode(operatorv1.InstallModeDefault, operatorv1.SPIFFEAuthType))
		Expect(err).ToNot(HaveOccurred())
		k.Spec.DeployOption.Mode = operatorv1.InstallModeHosted
		_, err = operatorClient.OperatorV1().Klusterlets().Update(context.TODO(), k, metav1.UpdateOptions{})
		Expect(err).To(HaveOccurred())
	})

	It("refuses an invalid agent trust domain", func() {
		d := spiffeDriver()
		d.SPIFFE.TrustDomain = "East_Example"
		_, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                d,
			AddOnKubeClientRegistrationDriver: spiffeAddOnDriver(),
		}))
		Expect(err).To(HaveOccurred())
	})

	withDriver := func(mutate func(*operatorv1.KlusterletSPIFFEConfig)) *operatorv1.RegistrationConfiguration {
		d := spiffeDriver()
		mutate(d.SPIFFE)
		return &operatorv1.RegistrationConfiguration{RegistrationDriver: d, AddOnKubeClientRegistrationDriver: spiffeAddOnDriver()}
	}

	It("refuses a hub endpoint over 512 characters", func() {
		_, err := create(newKlusterlet(withDriver(func(s *operatorv1.KlusterletSPIFFEConfig) {
			s.HubEndpoint = strings.Repeat("a", 509) + ":443"
		})))
		Expect(err).To(HaveOccurred())
	})

	It("refuses a hub ID over 512 characters", func() {
		_, err := create(newKlusterlet(withDriver(func(s *operatorv1.KlusterletSPIFFEConfig) {
			s.HubSPIFFEID = "spiffe://ocm-hub.example.org/" + strings.Repeat("a", 512)
		})))
		Expect(err).To(HaveOccurred())
	})

	It("refuses an agent trust domain over 255 characters", func() {
		_, err := create(newKlusterlet(withDriver(func(s *operatorv1.KlusterletSPIFFEConfig) {
			s.TrustDomain = strings.Repeat("a", 252) + ".org"
		})))
		Expect(err).To(HaveOccurred())
	})

	It("leaves the agent trust domain empty when not set; the operator uses the hub ID's", func() {
		k, err := create(newKlusterlet(withDriver(func(s *operatorv1.KlusterletSPIFFEConfig) {})))
		Expect(err).ToNot(HaveOccurred())
		Expect(k.Spec.RegistrationConfiguration.RegistrationDriver.SPIFFE.TrustDomain).Should(BeEmpty())
	})

	It("refuses spiffe with an add-on driver that has no authType", func() {
		_, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                spiffeDriver(),
			AddOnKubeClientRegistrationDriver: &operatorv1.AddOnRegistrationDriver{},
		}))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("addOnKubeClientRegistrationDriver.authType must be spiffe"))
	})

	Context("spiffe status", func() {
		// updateStatus sets status.spiffe on the test's Klusterlet, creating it on first use.
		updateStatus := func(s *operatorv1.KlusterletSPIFFEStatus) error {
			k, err := operatorClient.OperatorV1().Klusterlets().Get(context.TODO(), klusterletName, metav1.GetOptions{})
			if err != nil {
				k, err = create(newKlusterlet(nil))
				Expect(err).ToNot(HaveOccurred())
			}
			k.Status.Conditions = []metav1.Condition{}
			k.Status.SPIFFE = s
			_, err = operatorClient.OperatorV1().Klusterlets().UpdateStatus(context.TODO(), k, metav1.UpdateOptions{})
			return err
		}
		policy := func(name string) operatorv1.SPIFFEAdmissionPolicyStatus {
			return operatorv1.SPIFFEAdmissionPolicyStatus{
				PolicyName: name, PolicyUID: "a", PolicyGeneration: 1, BindingUID: "b", BindingGeneration: 1,
			}
		}

		It("accepts 16 policies and refuses 17", func() {
			var ps []operatorv1.SPIFFEAdmissionPolicyStatus
			for i := 0; i < 16; i++ {
				ps = append(ps, policy(fmt.Sprintf("p%d", i)))
			}
			Expect(updateStatus(&operatorv1.KlusterletSPIFFEStatus{AdmissionPolicies: ps})).To(Succeed())
			ps = append(ps, policy("p16"))
			Expect(updateStatus(&operatorv1.KlusterletSPIFFEStatus{AdmissionPolicies: ps})).ToNot(Succeed())
		})

		It("refuses an empty or over-long policy name", func() {
			Expect(updateStatus(&operatorv1.KlusterletSPIFFEStatus{
				AdmissionPolicies: []operatorv1.SPIFFEAdmissionPolicyStatus{policy("")},
			})).ToNot(Succeed())
			Expect(updateStatus(&operatorv1.KlusterletSPIFFEStatus{
				AdmissionPolicies: []operatorv1.SPIFFEAdmissionPolicyStatus{policy(strings.Repeat("p", 254))},
			})).ToNot(Succeed())
		})

		It("refuses an empty or over-long add-on namespace name", func() {
			Expect(updateStatus(&operatorv1.KlusterletSPIFFEStatus{
				AddOnNamespace: &operatorv1.SPIFFENamespaceRef{Name: "", UID: "u"},
			})).ToNot(Succeed())
			Expect(updateStatus(&operatorv1.KlusterletSPIFFEStatus{
				AddOnNamespace: &operatorv1.SPIFFENamespaceRef{Name: strings.Repeat("n", 64), UID: "u"},
			})).ToNot(Succeed())
		})
	})

	It("refuses authType spiffe without the spiffe settings", func() {
		_, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                operatorv1.RegistrationDriver{AuthType: operatorv1.SPIFFEAuthType},
			AddOnKubeClientRegistrationDriver: spiffeAddOnDriver(),
		}))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("spiffe is required when authType is spiffe"))
	})

	It("refuses spiffe with the csr add-on driver", func() {
		_, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                spiffeDriver(),
			AddOnKubeClientRegistrationDriver: &operatorv1.AddOnRegistrationDriver{AuthType: "csr"},
		}))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(
			"addOnKubeClientRegistrationDriver.authType must be spiffe when registrationDriver.authType is spiffe"))
	})

	It("refuses spiffe without an add-on driver", func() {
		_, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver: spiffeDriver(),
		}))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("addOnKubeClientRegistrationDriver.authType must be spiffe"))
	})

	It("refuses a hub ID that is not a SPIFFE ID", func() {
		d := spiffeDriver()
		d.SPIFFE.HubSPIFFEID = "https://ocm-hub.example.org/ocm/hub/grpc-server"
		_, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                d,
			AddOnKubeClientRegistrationDriver: spiffeAddOnDriver(),
		}))
		Expect(err).To(HaveOccurred())
	})

	It("refuses an empty hub endpoint", func() {
		d := spiffeDriver()
		d.SPIFFE.HubEndpoint = ""
		_, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                d,
			AddOnKubeClientRegistrationDriver: spiffeAddOnDriver(),
		}))
		Expect(err).To(HaveOccurred())
	})

	It("accepts the spiffe add-on driver with another registration driver", func() {
		_, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                operatorv1.RegistrationDriver{AuthType: operatorv1.CSRAuthType},
			AddOnKubeClientRegistrationDriver: spiffeAddOnDriver(),
		}))
		Expect(err).ToNot(HaveOccurred())
	})

	It("still accepts a Klusterlet without a registration configuration", func() {
		_, err := create(newKlusterlet(nil))
		Expect(err).ToNot(HaveOccurred())
	})

	It("stores and returns the spiffe status", func() {
		k, err := create(newKlusterlet(&operatorv1.RegistrationConfiguration{
			RegistrationDriver:                spiffeDriver(),
			AddOnKubeClientRegistrationDriver: spiffeAddOnDriver(),
		}))
		Expect(err).ToNot(HaveOccurred())

		k.Status.Conditions = []metav1.Condition{{
			Type:               operatorv1.ConditionSPIFFEAdmissionReady,
			Status:             metav1.ConditionTrue,
			Reason:             "Recorded",
			Message:            "admission state recorded",
			LastTransitionTime: metav1.Now(),
		}}
		k.Status.SPIFFE = &operatorv1.KlusterletSPIFFEStatus{
			AdmissionPolicies: []operatorv1.SPIFFEAdmissionPolicyStatus{
				{
					PolicyName:        "ocm-spiffe-agents-klusterlet",
					PolicyUID:         "6f1c2d0e-0b5e-4c4f-9a57-3d1f0c8e2a11",
					PolicyGeneration:  3,
					BindingUID:        "9a2e7b44-5c1d-4e0a-8f3b-2c6d9e1f7a05",
					BindingGeneration: 1,
				},
				{
					PolicyName:        "ocm-spiffe-addons-klusterlet",
					PolicyUID:         "1b2c3d4e-0000-4c4f-9a57-3d1f0c8e2a11",
					PolicyGeneration:  1,
					BindingUID:        "5e6f7a8b-0000-4e0a-8f3b-2c6d9e1f7a05",
					BindingGeneration: 1,
				},
			},
			AddOnNamespace: &operatorv1.SPIFFENamespaceRef{
				Name: "open-cluster-management-agent-addon",
				UID:  "3b8d5f21-7e4a-4c19-b6d2-8a0f1e9c4d73",
			},
		}
		_, err = operatorClient.OperatorV1().Klusterlets().UpdateStatus(context.TODO(), k, metav1.UpdateOptions{})
		Expect(err).ToNot(HaveOccurred())

		got, err := operatorClient.OperatorV1().Klusterlets().Get(context.TODO(), klusterletName, metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(got.Status.SPIFFE).Should(Equal(k.Status.SPIFFE))
	})

	It("refuses two status entries for the same policy", func() {
		k, err := create(newKlusterlet(nil))
		Expect(err).ToNot(HaveOccurred())

		entry := operatorv1.SPIFFEAdmissionPolicyStatus{
			PolicyName: "ocm-spiffe-agents-klusterlet", PolicyUID: "a", PolicyGeneration: 1,
			BindingUID: "b", BindingGeneration: 1,
		}
		k.Status.Conditions = []metav1.Condition{}
		k.Status.SPIFFE = &operatorv1.KlusterletSPIFFEStatus{
			AdmissionPolicies: []operatorv1.SPIFFEAdmissionPolicyStatus{entry, entry},
		}
		_, err = operatorClient.OperatorV1().Klusterlets().UpdateStatus(context.TODO(), k, metav1.UpdateOptions{})
		Expect(err).To(HaveOccurred())
	})

	Context("agent ID templates", func() {
		It("fills both templates with the defaults", func() {
			k, err := create(newKlusterlet(withDriver(func(s *operatorv1.KlusterletSPIFFEConfig) {})))
			Expect(err).ToNot(HaveOccurred())
			s := k.Spec.RegistrationConfiguration.RegistrationDriver.SPIFFE
			Expect(s.AgentIDTemplate).Should(Equal(operatorv1.DefaultSPIFFEAgentIDTemplate))
			Expect(s.AddOnAgentIDTemplate).Should(Equal(operatorv1.DefaultSPIFFEAddOnAgentIDTemplate))
		})

		It("keeps a custom template with a namespace segment", func() {
			k, err := create(newKlusterlet(withDriver(func(s *operatorv1.KlusterletSPIFFEConfig) {
				s.AgentIDTemplate = "spiffe://{trustDomain}/k8s/{cluster}/ns/{namespace}/ocm/{role}"
			})))
			Expect(err).ToNot(HaveOccurred())
			Expect(k.Spec.RegistrationConfiguration.RegistrationDriver.SPIFFE.AgentIDTemplate).
				Should(Equal("spiffe://{trustDomain}/k8s/{cluster}/ns/{namespace}/ocm/{role}"))
		})

		for _, c := range []struct{ name, agent, addOn, message string }{
			{"agent without {role}", "spiffe://{trustDomain}/ocm/cluster/{cluster}/agent", "",
				"agentIDTemplate must contain {cluster} and {role}"},
			{"agent without {cluster}", "spiffe://{trustDomain}/ocm/agent/{role}", "",
				"agentIDTemplate must contain {cluster} and {role}"},
			{"add-on without {agent}", "", "spiffe://{trustDomain}/ocm/{cluster}/{addon}",
				"addOnAgentIDTemplate must contain {cluster}, {addon} and {agent}"},
			{"add-on without {addon}", "", "spiffe://{trustDomain}/ocm/{cluster}/x/{agent}",
				"addOnAgentIDTemplate must contain {cluster}, {addon} and {agent}"},
			{"fixed trust domain", "spiffe://east.example.org/ocm/{cluster}/{role}", "", "should match"},
		} {
			c := c
			It("refuses a bad template: "+c.name, func() {
				_, err := create(newKlusterlet(withDriver(func(s *operatorv1.KlusterletSPIFFEConfig) {
					s.AgentIDTemplate = c.agent
					s.AddOnAgentIDTemplate = c.addOn
				})))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(c.message))
			})
		}
	})

	Context("cluster name", func() {
		It("refuses spiffe without a cluster name", func() {
			k := newKlusterlet(withDriver(func(s *operatorv1.KlusterletSPIFFEConfig) {}))
			k.Spec.ClusterName = ""
			_, err := create(k)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("clusterName is required when registrationDriver.authType is spiffe"))
		})

		It("still allows other drivers without a cluster name", func() {
			k := newKlusterlet(&operatorv1.RegistrationConfiguration{
				RegistrationDriver: operatorv1.RegistrationDriver{AuthType: operatorv1.CSRAuthType},
			})
			k.Spec.ClusterName = ""
			_, err := create(k)
			Expect(err).ToNot(HaveOccurred())
		})

		It("still allows no registration configuration without a cluster name", func() {
			k := newKlusterlet(nil)
			k.Spec.ClusterName = ""
			_, err := create(k)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("agent IDs in the status", func() {
		updateAgentIDs := func(ids []operatorv1.SPIFFEAgentID) error {
			k, err := operatorClient.OperatorV1().Klusterlets().Get(context.TODO(), klusterletName, metav1.GetOptions{})
			if err != nil {
				k, err = create(newKlusterlet(nil))
				Expect(err).ToNot(HaveOccurred())
			}
			k.Status.Conditions = []metav1.Condition{}
			k.Status.SPIFFE = &operatorv1.KlusterletSPIFFEStatus{AgentIDs: ids}
			_, err = operatorClient.OperatorV1().Klusterlets().UpdateStatus(context.TODO(), k, metav1.UpdateOptions{})
			return err
		}
		id := func(role string) operatorv1.SPIFFEAgentID {
			return operatorv1.SPIFFEAgentID{
				Role:     role,
				SPIFFEID: "spiffe://ocm-hub.example.org/ocm/cluster/cluster1/agent/" + role,
			}
		}

		It("accepts the registration and work IDs", func() {
			Expect(updateAgentIDs([]operatorv1.SPIFFEAgentID{id("registration"), id("work")})).To(Succeed())
		})

		It("refuses an unknown role", func() {
			Expect(updateAgentIDs([]operatorv1.SPIFFEAgentID{id("admin")})).ToNot(Succeed())
		})

		It("refuses the same role twice", func() {
			Expect(updateAgentIDs([]operatorv1.SPIFFEAgentID{id("work"), id("work")})).ToNot(Succeed())
		})

		It("refuses an empty SPIFFE ID", func() {
			Expect(updateAgentIDs([]operatorv1.SPIFFEAgentID{{Role: "work"}})).ToNot(Succeed())
		})
	})
})

var _ = Describe("ManagedClusterAddOn API test with the spiffe driver", func() {
	var name string

	BeforeEach(func() {
		name = fmt.Sprintf("mca-%s", rand.String(5))
	})

	It("accepts driver spiffe in v1beta1", func() {
		_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(context.TODO(),
			&addonv1beta1.ManagedClusterAddOn{ObjectMeta: metav1.ObjectMeta{Name: name}}, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		mca, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Get(context.TODO(), name, metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		mca.Status.Namespace = testNamespace
		mca.Status.Registrations = []addonv1beta1.RegistrationConfig{{
			Type:       addonv1beta1.KubeClient,
			KubeClient: &addonv1beta1.KubeClientConfig{Driver: "spiffe"},
		}}
		// Only acceptance is checked: envtest has no conversion webhook, so v1beta1-only
		// registration details don't survive the round trip through the stored version.
		_, err = hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).UpdateStatus(
			context.TODO(), mca, metav1.UpdateOptions{})
		Expect(err).ToNot(HaveOccurred())
	})

	It("refuses an unknown driver in v1beta1", func() {
		_, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Create(context.TODO(),
			&addonv1beta1.ManagedClusterAddOn{ObjectMeta: metav1.ObjectMeta{Name: name}}, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		mca, err := hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).Get(context.TODO(), name, metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		mca.Status.Namespace = testNamespace
		mca.Status.Registrations = []addonv1beta1.RegistrationConfig{{
			Type:       addonv1beta1.KubeClient,
			KubeClient: &addonv1beta1.KubeClientConfig{Driver: "spire"},
		}}
		_, err = hubAddonClient.AddonV1beta1().ManagedClusterAddOns(testNamespace).UpdateStatus(
			context.TODO(), mca, metav1.UpdateOptions{})
		Expect(err).To(HaveOccurred())
	})

	It("accepts kubeClientDriver spiffe in v1alpha1", func() {
		_, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Create(context.TODO(),
			&addonv1alpha1.ManagedClusterAddOn{ObjectMeta: metav1.ObjectMeta{Name: name}}, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		mca, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).Get(context.TODO(), name, metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		mca.Status.KubeClientDriver = "spiffe"
		updated, err := hubAddonClient.AddonV1alpha1().ManagedClusterAddOns(testNamespace).UpdateStatus(
			context.TODO(), mca, metav1.UpdateOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(updated.Status.KubeClientDriver).Should(Equal("spiffe"))
	})
})
