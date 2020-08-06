// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/open-cluster-management/api/cluster/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ManagedClusterClaimLister helps list ManagedClusterClaims.
type ManagedClusterClaimLister interface {
	// List lists all ManagedClusterClaims in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.ManagedClusterClaim, err error)
	// ManagedClusterClaims returns an object that can list and get ManagedClusterClaims.
	ManagedClusterClaims(namespace string) ManagedClusterClaimNamespaceLister
	ManagedClusterClaimListerExpansion
}

// managedClusterClaimLister implements the ManagedClusterClaimLister interface.
type managedClusterClaimLister struct {
	indexer cache.Indexer
}

// NewManagedClusterClaimLister returns a new ManagedClusterClaimLister.
func NewManagedClusterClaimLister(indexer cache.Indexer) ManagedClusterClaimLister {
	return &managedClusterClaimLister{indexer: indexer}
}

// List lists all ManagedClusterClaims in the indexer.
func (s *managedClusterClaimLister) List(selector labels.Selector) (ret []*v1alpha1.ManagedClusterClaim, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ManagedClusterClaim))
	})
	return ret, err
}

// ManagedClusterClaims returns an object that can list and get ManagedClusterClaims.
func (s *managedClusterClaimLister) ManagedClusterClaims(namespace string) ManagedClusterClaimNamespaceLister {
	return managedClusterClaimNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ManagedClusterClaimNamespaceLister helps list and get ManagedClusterClaims.
type ManagedClusterClaimNamespaceLister interface {
	// List lists all ManagedClusterClaims in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.ManagedClusterClaim, err error)
	// Get retrieves the ManagedClusterClaim from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.ManagedClusterClaim, error)
	ManagedClusterClaimNamespaceListerExpansion
}

// managedClusterClaimNamespaceLister implements the ManagedClusterClaimNamespaceLister
// interface.
type managedClusterClaimNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ManagedClusterClaims in the indexer for a given namespace.
func (s managedClusterClaimNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ManagedClusterClaim, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ManagedClusterClaim))
	})
	return ret, err
}

// Get retrieves the ManagedClusterClaim from the indexer for a given namespace and name.
func (s managedClusterClaimNamespaceLister) Get(name string) (*v1alpha1.ManagedClusterClaim, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("managedclusterclaim"), name)
	}
	return obj.(*v1alpha1.ManagedClusterClaim), nil
}
