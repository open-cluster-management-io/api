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
