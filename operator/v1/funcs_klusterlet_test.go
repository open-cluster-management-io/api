// Copyright Contributors to the Open Cluster Management project
package v1

import "testing"

func TestKlusterletGetResourceRequirement(t *testing.T) {
	rr := &ResourceRequirement{Type: ResourceQosClassDefault}
	k := &Klusterlet{Spec: KlusterletSpec{ResourceRequirement: rr}}
	if got := k.GetResourceRequirement(); got != rr {
		t.Errorf("GetResourceRequirement() = %v, want %v", got, rr)
	}
}

func TestKlusterletSPIFFEConfigIDTemplates(t *testing.T) {
	custom := "spiffe://{trustDomain}/k8s/{cluster}/ns/{namespace}/ocm/{role}"
	customAddOn := "spiffe://{trustDomain}/k8s/{cluster}/addon/{addon}/{agent}"
	cases := []struct {
		name                 string
		config               *KlusterletSPIFFEConfig
		wantAgent, wantAddOn string
	}{
		{"nil config", nil, DefaultSPIFFEAgentIDTemplate, DefaultSPIFFEAddOnAgentIDTemplate},
		{"empty", &KlusterletSPIFFEConfig{}, DefaultSPIFFEAgentIDTemplate, DefaultSPIFFEAddOnAgentIDTemplate},
		{"set", &KlusterletSPIFFEConfig{AgentIDTemplate: custom, AddOnAgentIDTemplate: customAddOn}, custom, customAddOn},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.config.GetAgentIDTemplate(); got != c.wantAgent {
				t.Errorf("GetAgentIDTemplate() = %q, want %q", got, c.wantAgent)
			}
			if got := c.config.GetAddOnAgentIDTemplate(); got != c.wantAddOn {
				t.Errorf("GetAddOnAgentIDTemplate() = %q, want %q", got, c.wantAddOn)
			}
		})
	}
}
