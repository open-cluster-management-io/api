// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"
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

	It("refuses an unknown endpoint protocol", func() {
		cm := newClusterManager(operatorv1.RegistrationDriverHub{AuthType: operatorv1.CSRAuthType})
		cm.Spec.ServerConfiguration = &operatorv1.ServerConfiguration{
			EndpointsExposure: []operatorv1.EndpointExposure{{Protocol: "grpc-other"}},
		}
		_, err := create(cm)
		Expect(err).To(HaveOccurred())
	})
})
