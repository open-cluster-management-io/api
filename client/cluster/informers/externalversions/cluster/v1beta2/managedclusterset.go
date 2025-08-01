// Copyright Contributors to the Open Cluster Management project
// Code generated by informer-gen. DO NOT EDIT.

package v1beta2

import (
	context "context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	versioned "open-cluster-management.io/api/client/cluster/clientset/versioned"
	internalinterfaces "open-cluster-management.io/api/client/cluster/informers/externalversions/internalinterfaces"
	clusterv1beta2 "open-cluster-management.io/api/client/cluster/listers/cluster/v1beta2"
	apiclusterv1beta2 "open-cluster-management.io/api/cluster/v1beta2"
)

// ManagedClusterSetInformer provides access to a shared informer and lister for
// ManagedClusterSets.
type ManagedClusterSetInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() clusterv1beta2.ManagedClusterSetLister
}

type managedClusterSetInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewManagedClusterSetInformer constructs a new informer for ManagedClusterSet type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewManagedClusterSetInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredManagedClusterSetInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredManagedClusterSetInformer constructs a new informer for ManagedClusterSet type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredManagedClusterSetInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ClusterV1beta2().ManagedClusterSets().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ClusterV1beta2().ManagedClusterSets().Watch(context.TODO(), options)
			},
		},
		&apiclusterv1beta2.ManagedClusterSet{},
		resyncPeriod,
		indexers,
	)
}

func (f *managedClusterSetInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredManagedClusterSetInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *managedClusterSetInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiclusterv1beta2.ManagedClusterSet{}, f.defaultInformer)
}

func (f *managedClusterSetInformer) Lister() clusterv1beta2.ManagedClusterSetLister {
	return clusterv1beta2.NewManagedClusterSetLister(f.Informer().GetIndexer())
}
