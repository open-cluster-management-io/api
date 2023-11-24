// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
	v1beta1 "open-cluster-management.io/api/client/cluster/clientset/versioned/typed/cluster/v1beta1"
)

type FakeClusterV1beta1 struct {
	*testing.Fake
}

func (c *FakeClusterV1beta1) ManagedClusterSets() v1beta1.ManagedClusterSetInterface {
	return &FakeManagedClusterSets{c}
}

func (c *FakeClusterV1beta1) ManagedClusterSetBindings(namespace string) v1beta1.ManagedClusterSetBindingInterface {
	return &FakeManagedClusterSetBindings{c, namespace}
}

func (c *FakeClusterV1beta1) Placements(namespace string) v1beta1.PlacementInterface {
	return &FakePlacements{c, namespace}
}

func (c *FakeClusterV1beta1) PlacementDecisions(namespace string) v1beta1.PlacementDecisionInterface {
	return &FakePlacementDecisions{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeClusterV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
