// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	operatorv1 "open-cluster-management.io/api/operator/v1"
)

var _ = Describe("ClusterManager API test with the spiffe registration driver", func() {
	var clusterManagerName string

	BeforeEach(func() {
		clusterManagerName = fmt.Sprintf("cm-%s", rand.String(5))
	})

	newClusterManager := func(driver operatorv1.RegistrationDriverHub) *operatorv1.ClusterManager {
		return &operatorv1.ClusterManager{
			ObjectMeta: metav1.ObjectMeta{Name: clusterManagerName},
			Spec: operatorv1.ClusterManagerSpec{
				RegistrationConfiguration: &operatorv1.RegistrationHubConfiguration{
					RegistrationDrivers: []operatorv1.RegistrationDriverHub{driver},
				},
			},
		}
	}

	create := func(cm *operatorv1.ClusterManager) (*operatorv1.ClusterManager, error) {
		return operatorClient.OperatorV1().ClusterManagers().Create(context.TODO(), cm, metav1.CreateOptions{})
	}

	It("accepts a minimal spiffe driver and fills the defaults", func() {
		cm, err := create(newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE:   &operatorv1.SPIFFERegistrationConfig{TrustDomain: "ocm-hub.example.org"},
		}))
		Expect(err).ToNot(HaveOccurred())

		s := cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE
		Expect(s.TrustDomain).Should(Equal("ocm-hub.example.org"))
		Expect(s.ServerSPIFFEID).Should(BeEmpty())
		Expect(s.WorkloadAPISocket).Should(Equal("unix:///spiffe-workload-api/spire-agent.sock"))
		Expect(s.MaxSVIDLifetime).ShouldNot(BeNil())
		Expect(s.MaxSVIDLifetime.Duration).Should(Equal(24 * time.Hour))
		Expect(s.MaxClockSkew).ShouldNot(BeNil())
		Expect(s.MaxClockSkew.Duration).Should(Equal(5 * time.Minute))
	})

	It("accepts all spiffe settings", func() {
		cm, err := create(newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE: &operatorv1.SPIFFERegistrationConfig{
				TrustDomain:       "ocm-hub.example.org",
				ServerSPIFFEID:    "spiffe://ocm-hub.example.org/ocm/hub/grpc-server",
				WorkloadAPISocket: "unix:///run/spire/agent.sock",
				MaxSVIDLifetime:   &metav1.Duration{Duration: time.Hour},
				MaxClockSkew:      &metav1.Duration{Duration: time.Minute},
			},
		}))
		Expect(err).ToNot(HaveOccurred())

		s := cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE
		Expect(s.ServerSPIFFEID).Should(Equal("spiffe://ocm-hub.example.org/ocm/hub/grpc-server"))
		Expect(s.WorkloadAPISocket).Should(Equal("unix:///run/spire/agent.sock"))
		Expect(s.MaxSVIDLifetime.Duration).Should(Equal(time.Hour))
		Expect(s.MaxClockSkew.Duration).Should(Equal(time.Minute))
	})

	It("accepts agents from several trust domains", func() {
		cm, err := create(newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE: &operatorv1.SPIFFERegistrationConfig{
				TrustDomain:       "ocm-hub.example.org",
				AgentTrustDomains: []string{"ocm-hub.example.org", "east.example.org", "west.example.org"},
			},
		}))
		Expect(err).ToNot(HaveOccurred())
		Expect(cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE.AgentTrustDomains).Should(
			Equal([]string{"ocm-hub.example.org", "east.example.org", "west.example.org"}))
	})

	It("leaves agentTrustDomains empty when not set; the operator uses trustDomain", func() {
		cm, err := create(newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE:   &operatorv1.SPIFFERegistrationConfig{TrustDomain: "ocm-hub.example.org"},
		}))
		Expect(err).ToNot(HaveOccurred())
		Expect(cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE.AgentTrustDomains).Should(BeEmpty())
	})

	It("refuses an invalid agent trust domain", func() {
		_, err := create(newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE: &operatorv1.SPIFFERegistrationConfig{
				TrustDomain:       "ocm-hub.example.org",
				AgentTrustDomains: []string{"east.example.org", "West_Example"},
			},
		}))
		Expect(err).To(HaveOccurred())
	})

	It("refuses a repeated agent trust domain", func() {
		_, err := create(newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE: &operatorv1.SPIFFERegistrationConfig{
				TrustDomain:       "ocm-hub.example.org",
				AgentTrustDomains: []string{"east.example.org", "east.example.org"},
			},
		}))
		Expect(err).To(HaveOccurred())
	})

	It("refuses authType spiffe without the spiffe settings", func() {
		_, err := create(newClusterManager(operatorv1.RegistrationDriverHub{AuthType: operatorv1.SPIFFEAuthType}))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("spiffe is required when authType is spiffe"))
	})

	// That serverSPIFFEID is in trustDomain is checked by the cluster-manager operator, not
	// by CEL: registrationDrivers has no maxItems, so a string rule there exceeds the CRD
	// validation cost budget.

	It("refuses an invalid trust domain", func() {
		_, err := create(newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE:   &operatorv1.SPIFFERegistrationConfig{TrustDomain: "OCM_Hub"},
		}))
		Expect(err).To(HaveOccurred())
	})

	It("refuses a server ID that is not a SPIFFE ID", func() {
		_, err := create(newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE: &operatorv1.SPIFFERegistrationConfig{
				TrustDomain:    "ocm-hub.example.org",
				ServerSPIFFEID: "https://ocm-hub.example.org/ocm/hub/grpc-server",
			},
		}))
		Expect(err).To(HaveOccurred())
	})

	It("refuses an invalid maxSVIDLifetime", func() {
		cm := newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE: &operatorv1.SPIFFERegistrationConfig{
				TrustDomain:     "ocm-hub.example.org",
				MaxSVIDLifetime: &metav1.Duration{Duration: 1500 * time.Millisecond},
			},
		})
		_, err := create(cm)
		Expect(err).To(HaveOccurred())
	})

	It("accepts spiffe next to grpc", func() {
		cm := newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE:   &operatorv1.SPIFFERegistrationConfig{TrustDomain: "ocm-hub.example.org"},
		})
		cm.Spec.RegistrationConfiguration.RegistrationDrivers = append(
			cm.Spec.RegistrationConfiguration.RegistrationDrivers,
			operatorv1.RegistrationDriverHub{AuthType: operatorv1.GRPCAuthType})
		_, err := create(cm)
		Expect(err).ToNot(HaveOccurred())
	})

	It("accepts a grpc-spiffe endpoint exposure", func() {
		cm := newClusterManager(operatorv1.RegistrationDriverHub{
			AuthType: operatorv1.SPIFFEAuthType,
			SPIFFE:   &operatorv1.SPIFFERegistrationConfig{TrustDomain: "ocm-hub.example.org"},
		})
		cm.Spec.ServerConfiguration = &operatorv1.ServerConfiguration{
			EndpointsExposure: []operatorv1.EndpointExposure{{
				Protocol: operatorv1.GRPCSPIFFEProtocol,
				GRPC:     &operatorv1.Endpoint{Type: operatorv1.EndpointTypeLoadBalancer},
			}},
		}
		cm, err := create(cm)
		Expect(err).ToNot(HaveOccurred())
		Expect(cm.Spec.ServerConfiguration.EndpointsExposure[0].Protocol).Should(Equal("grpc-spiffe"))
	})

	spiffeConfig := func(mutate func(*operatorv1.SPIFFERegistrationConfig)) operatorv1.RegistrationDriverHub {
		c := &operatorv1.SPIFFERegistrationConfig{TrustDomain: "ocm-hub.example.org"}
		mutate(c)
		return operatorv1.RegistrationDriverHub{AuthType: operatorv1.SPIFFEAuthType, SPIFFE: c}
	}

	// domainOfLength returns a valid trust domain name of exactly n characters.
	domainOfLength := func(n int) string { return strings.Repeat("a", n-4) + ".org" }

	// A duration Go can't hold, such as 2562048h, would be stored but no client could read it back.
	// Go clients can't send one, so it is posted as raw JSON.
	Context("maxSVIDLifetime and maxClockSkew bounds", func() {
		post := func(field, value string) error {
			body := fmt.Sprintf(`{"apiVersion":"operator.open-cluster-management.io/v1","kind":"ClusterManager",`+
				`"metadata":{"name":%q},"spec":{"registrationConfiguration":{"registrationDrivers":[`+
				`{"authType":"spiffe","spiffe":{"trustDomain":"ocm-hub.example.org",%q:%q}}]}}}`,
				fmt.Sprintf("cm-%s", rand.String(5)), field, value)
			return operatorClient.OperatorV1().RESTClient().Post().Resource("clustermanagers").
				Body([]byte(body)).Do(context.TODO()).Error()
		}
		for _, field := range []string{"maxSVIDLifetime", "maxClockSkew"} {
			for _, good := range []string{"24h", "1h30m", "99999h"} {
				It(field+" accepts "+good, func() {
					Expect(post(field, good)).To(Succeed())
				})
			}
			for _, bad := range []string{"2562048h", "999999h", "99999h99999h99999h99999h99999h"} {
				It(field+" refuses "+bad, func() {
					err := post(field, bad)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(field))
				})
			}
		}
	})

	// SPIFFE forbids the path segments . and ..; other segments with dots stay valid.
	Context("dot path segments", func() {
		for _, bad := range []struct{ name, value string }{
			{"server ID with ..", "spiffe://ocm-hub.example.org/ocm/../server"},
			{"server ID ending in .", "spiffe://ocm-hub.example.org/ocm/."},
		} {
			It("refuses a "+bad.name, func() {
				_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
					c.ServerSPIFFEID = bad.value
				})))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("serverSPIFFEID"))
			})
		}
		It("refuses . and .. in the ID templates and auto-approved identities", func() {
			for _, mutate := range []func(*operatorv1.SPIFFERegistrationConfig){
				func(c *operatorv1.SPIFFERegistrationConfig) {
					c.AgentIDTemplate = "spiffe://{trustDomain}/ocm/./{cluster}/{role}"
				},
				func(c *operatorv1.SPIFFERegistrationConfig) {
					c.AddOnAgentIDTemplate = "spiffe://{trustDomain}/ocm/../{cluster}/{addon}/{agent}"
				},
				func(c *operatorv1.SPIFFERegistrationConfig) {
					c.AutoApprovedIdentities = []string{"spiffe://ocm-hub.example.org/ocm/../*/agent/registration"}
				},
			} {
				clusterManagerName = fmt.Sprintf("cm-%s", rand.String(5))
				_, err := create(newClusterManager(spiffeConfig(mutate)))
				Expect(err).To(HaveOccurred())
			}
		})
		It("accepts segments that only contain dots among other characters, and ...", func() {
			cm, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.ServerSPIFFEID = "spiffe://ocm-hub.example.org/.well-known/.../grpc.server"
				c.AgentIDTemplate = "spiffe://{trustDomain}/ocm.v1/cluster/{cluster}/agent/{role}"
			})))
			Expect(err).ToNot(HaveOccurred())
			Expect(cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE.ServerSPIFFEID).
				Should(Equal("spiffe://ocm-hub.example.org/.well-known/.../grpc.server"))
		})
	})

	Context("trustDomain", func() {
		It("refuses a missing trust domain", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) { c.TrustDomain = "" })))
			Expect(err).To(HaveOccurred())
		})
		It("accepts 255 characters", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.TrustDomain = domainOfLength(255)
			})))
			Expect(err).ToNot(HaveOccurred())
		})
		It("refuses 256 characters", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.TrustDomain = domainOfLength(256)
			})))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("autoApprovedIdentities", func() {
		It("accepts SPIFFE IDs, with * for whole path segments", func() {
			ids := []string{"spiffe://ocm-hub.example.org/ocm/cluster/*/agent/registration",
				"spiffe://east.example.org/ocm/cluster/cluster1/agent/klusterlet"}
			cm, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.AutoApprovedIdentities = ids
			})))
			Expect(err).ToNot(HaveOccurred())
			Expect(cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE.AutoApprovedIdentities).Should(Equal(ids))
		})
		for _, bad := range []string{"ocm/cluster/*", "https://ocm-hub.example.org/x", "spiffe://ocm-hub.example.org",
			"spiffe://Hub/x", "spiffe://ocm-hub.example.org/ocm/clu*ster"} {
			It("refuses "+bad, func() {
				_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
					c.AutoApprovedIdentities = []string{bad}
				})))
				Expect(err).To(HaveOccurred())
			})
		}
		It("refuses a repeated entry", func() {
			id := "spiffe://ocm-hub.example.org/ocm/cluster/*/agent/registration"
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.AutoApprovedIdentities = []string{id, id}
			})))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("agentTrustDomains", func() {
		It("is stored as the admin wrote it when it lists the hub's trust domain", func() {
			cm, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.AgentTrustDomains = []string{"ocm-hub.example.org", "east.example.org"}
			})))
			Expect(err).ToNot(HaveOccurred())
			s := cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE
			Expect(s.AgentTrustDomains).Should(Equal([]string{"ocm-hub.example.org", "east.example.org"}))
			Expect(s.GetAgentTrustDomains()).Should(Equal([]string{"ocm-hub.example.org", "east.example.org"}))
		})
		It("adds the hub's trust domain when the admin leaves it out", func() {
			cm, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.AgentTrustDomains = []string{"east.example.org", "west.example.org"}
			})))
			Expect(err).ToNot(HaveOccurred())
			s := cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE
			// The stored spec is not rewritten; the accepted set always includes the hub's.
			Expect(s.AgentTrustDomains).Should(Equal([]string{"east.example.org", "west.example.org"}))
			Expect(s.GetAgentTrustDomains()).Should(Equal(
				[]string{"ocm-hub.example.org", "east.example.org", "west.example.org"}))
		})
		It("uses only the hub's trust domain when the admin sets none", func() {
			cm, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {})))
			Expect(err).ToNot(HaveOccurred())
			s := cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE
			Expect(s.GetAgentTrustDomains()).Should(Equal([]string{"ocm-hub.example.org"}))
		})
		It("accepts 64 entries", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				for i := 0; i < 64; i++ {
					c.AgentTrustDomains = append(c.AgentTrustDomains, fmt.Sprintf("td%d.example.org", i))
				}
			})))
			Expect(err).ToNot(HaveOccurred())
		})
		It("refuses 65 entries", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				for i := 0; i < 65; i++ {
					c.AgentTrustDomains = append(c.AgentTrustDomains, fmt.Sprintf("td%d.example.org", i))
				}
			})))
			Expect(err).To(HaveOccurred())
		})
		It("refuses an empty entry", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.AgentTrustDomains = []string{"east.example.org", ""}
			})))
			Expect(err).To(HaveOccurred())
		})
		It("refuses an entry of 256 characters", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.AgentTrustDomains = []string{domainOfLength(256)}
			})))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("serverSPIFFEID", func() {
		It("refuses an ID without a path", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.ServerSPIFFEID = "spiffe://ocm-hub.example.org"
			})))
			Expect(err).To(HaveOccurred())
		})
		It("refuses more than 512 characters", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.ServerSPIFFEID = "spiffe://ocm-hub.example.org/" + strings.Repeat("a", 512)
			})))
			Expect(err).To(HaveOccurred())
		})
	})

	It("refuses an invalid maxClockSkew", func() {
		_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
			c.MaxClockSkew = &metav1.Duration{Duration: 500 * time.Millisecond}
		})))
		Expect(err).To(HaveOccurred())
	})

	It("refuses an unknown endpoint protocol", func() {
		cm := newClusterManager(operatorv1.RegistrationDriverHub{AuthType: operatorv1.CSRAuthType})
		cm.Spec.ServerConfiguration = &operatorv1.ServerConfiguration{
			EndpointsExposure: []operatorv1.EndpointExposure{{Protocol: "grpc-other"}},
		}
		_, err := create(cm)
		Expect(err).To(HaveOccurred())
	})

	Context("agent ID templates", func() {
		It("fills both templates with the defaults", func() {
			cm, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {})))
			Expect(err).ToNot(HaveOccurred())
			s := cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE
			Expect(s.AgentIDTemplate).Should(Equal(operatorv1.DefaultSPIFFEAgentIDTemplate))
			Expect(s.AddOnAgentIDTemplate).Should(Equal(operatorv1.DefaultSPIFFEAddOnAgentIDTemplate))
		})

		It("keeps custom templates with a namespace segment", func() {
			cm, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.AgentIDTemplate = "spiffe://{trustDomain}/k8s/{cluster}/ns/{namespace}/ocm/{role}"
				c.AddOnAgentIDTemplate = "spiffe://{trustDomain}/k8s/{cluster}/ns/{namespace}/addon/{addon}/{agent}"
			})))
			Expect(err).ToNot(HaveOccurred())
			s := cm.Spec.RegistrationConfiguration.RegistrationDrivers[0].SPIFFE
			Expect(s.AgentIDTemplate).Should(Equal("spiffe://{trustDomain}/k8s/{cluster}/ns/{namespace}/ocm/{role}"))
			Expect(s.AddOnAgentIDTemplate).Should(Equal("spiffe://{trustDomain}/k8s/{cluster}/ns/{namespace}/addon/{addon}/{agent}"))
		})

		// The CRD checks only the shape. Unknown or missing placeholders, duplicates and overlap
		// are checked by the cluster-manager operator, which reports Degraded.
		It("leaves unknown placeholders to the operator", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.AgentIDTemplate = "spiffe://{trustDomain}/ocm/{clutser}/{role}"
			})))
			Expect(err).ToNot(HaveOccurred())
		})

		for _, c := range []struct{ name, agent, addOn string }{
			{"fixed trust domain", "spiffe://ocm-hub.example.org/ocm/{cluster}/{role}", ""},
			{"no path", "spiffe://{trustDomain}", ""},
			{"placeholder inside a segment", "spiffe://{trustDomain}/ocm/cluster-{cluster}/{role}", ""},
			{"empty segment", "spiffe://{trustDomain}/ocm//{cluster}/{role}", ""},
			{"trailing slash", "spiffe://{trustDomain}/ocm/{cluster}/{role}/", ""},
			{"query in add-on template", "", "spiffe://{trustDomain}/ocm/{cluster}/{addon}/{agent}?x"},
		} {
			c := c
			It("refuses a malformed template: "+c.name, func() {
				_, err := create(newClusterManager(spiffeConfig(func(s *operatorv1.SPIFFERegistrationConfig) {
					s.AgentIDTemplate = c.agent
					s.AddOnAgentIDTemplate = c.addOn
				})))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("should match"))
			})
		}

		It("refuses a template over 512 characters", func() {
			_, err := create(newClusterManager(spiffeConfig(func(c *operatorv1.SPIFFERegistrationConfig) {
				c.AgentIDTemplate = "spiffe://{trustDomain}/{cluster}/{role}/" + strings.Repeat("a", 512)
			})))
			Expect(err).To(HaveOccurred())
		})
	})
})
