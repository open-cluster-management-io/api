// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
	operatorv1 "open-cluster-management.io/api/operator/v1"
)

// ClusterManagerLister helps list ClusterManagers.
// All objects returned here must be treated as read-only.
type ClusterManagerLister interface {
	// List lists all ClusterManagers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*operatorv1.ClusterManager, err error)
	// Get retrieves the ClusterManager from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*operatorv1.ClusterManager, error)
	ClusterManagerListerExpansion
}

// clusterManagerLister implements the ClusterManagerLister interface.
type clusterManagerLister struct {
	listers.ResourceIndexer[*operatorv1.ClusterManager]
}

// NewClusterManagerLister returns a new ClusterManagerLister.
func NewClusterManagerLister(indexer cache.Indexer) ClusterManagerLister {
	return &clusterManagerLister{listers.New[*operatorv1.ClusterManager](indexer, operatorv1.Resource("clustermanager"))}
}
