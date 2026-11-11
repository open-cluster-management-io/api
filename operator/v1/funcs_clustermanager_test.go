// Copyright Contributors to the Open Cluster Management project
package v1

import (
	"reflect"
	"testing"
)

func TestGetResourceRequirement(t *testing.T) {
	rr := &ResourceRequirement{Type: ResourceQosClassDefault}
	cm := &ClusterManager{Spec: ClusterManagerSpec{ResourceRequirement: rr}}
	if got := cm.GetResourceRequirement(); got != rr {
		t.Errorf("GetResourceRequirement() = %v, want %v", got, rr)
	}
}

func TestGetAgentTrustDomains(t *testing.T) {
	const hub = "ocm-hub.example.org"
	cases := []struct {
		name   string
		config *SPIFFERegistrationConfig
		want   []string
	}{
		{
			name:   "nil config",
			config: nil,
			want:   nil,
		},
		{
			name:   "no trust domain",
			config: &SPIFFERegistrationConfig{AgentTrustDomains: []string{"east.example.org"}},
			want:   nil,
		},
		{
			name:   "agent trust domains not set: only the hub's",
			config: &SPIFFERegistrationConfig{TrustDomain: hub},
			want:   []string{hub},
		},
		{
			name:   "agent trust domains empty: only the hub's",
			config: &SPIFFERegistrationConfig{TrustDomain: hub, AgentTrustDomains: []string{}},
			want:   []string{hub},
		},
		{
			name: "admin lists the hub's trust domain",
			config: &SPIFFERegistrationConfig{
				TrustDomain:       hub,
				AgentTrustDomains: []string{hub, "east.example.org", "west.example.org"},
			},
			want: []string{hub, "east.example.org", "west.example.org"},
		},
		{
			name: "admin leaves out the hub's trust domain: it is added first",
			config: &SPIFFERegistrationConfig{
				TrustDomain:       hub,
				AgentTrustDomains: []string{"east.example.org", "west.example.org"},
			},
			want: []string{hub, "east.example.org", "west.example.org"},
		},
		{
			name: "admin lists the hub's trust domain last: it moves first, once",
			config: &SPIFFERegistrationConfig{
				TrustDomain:       hub,
				AgentTrustDomains: []string{"east.example.org", hub},
			},
			want: []string{hub, "east.example.org"},
		},
		{
			name: "duplicates and empty entries are dropped",
			config: &SPIFFERegistrationConfig{
				TrustDomain:       hub,
				AgentTrustDomains: []string{"east.example.org", "", "east.example.org"},
			},
			want: []string{hub, "east.example.org"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.config.GetAgentTrustDomains()
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("GetAgentTrustDomains() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestGetAgentTrustDomainsDoesNotChangeTheSpec(t *testing.T) {
	config := &SPIFFERegistrationConfig{
		TrustDomain:       "ocm-hub.example.org",
		AgentTrustDomains: []string{"east.example.org"},
	}
	_ = config.GetAgentTrustDomains()
	if want := []string{"east.example.org"}; !reflect.DeepEqual(config.AgentTrustDomains, want) {
		t.Errorf("AgentTrustDomains changed to %v, want %v", config.AgentTrustDomains, want)
	}
}
