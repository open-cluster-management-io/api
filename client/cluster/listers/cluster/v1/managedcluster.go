// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

// ManagedClusterLister helps list ManagedClusters.
// All objects returned here must be treated as read-only.
type ManagedClusterLister interface {
	// List lists all ManagedClusters in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*clusterv1.ManagedCluster, err error)
	// Get retrieves the ManagedCluster from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*clusterv1.ManagedCluster, error)
	ManagedClusterListerExpansion
}

// managedClusterLister implements the ManagedClusterLister interface.
type managedClusterLister struct {
	listers.ResourceIndexer[*clusterv1.ManagedCluster]
}

// NewManagedClusterLister returns a new ManagedClusterLister.
func NewManagedClusterLister(indexer cache.Indexer) ManagedClusterLister {
	return &managedClusterLister{listers.New[*clusterv1.ManagedCluster](indexer, clusterv1.Resource("managedcluster"))}
}
