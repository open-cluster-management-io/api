// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	versioned "github.com/open-cluster-management/api/client/cluster/clientset/versioned"
	internalinterfaces "github.com/open-cluster-management/api/client/cluster/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/open-cluster-management/api/client/cluster/listers/cluster/v1alpha1"
	clusterv1alpha1 "github.com/open-cluster-management/api/cluster/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ManagedClusterSetBindingInformer provides access to a shared informer and lister for
// ManagedClusterSetBindings.
type ManagedClusterSetBindingInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ManagedClusterSetBindingLister
}

type managedClusterSetBindingInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewManagedClusterSetBindingInformer constructs a new informer for ManagedClusterSetBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewManagedClusterSetBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredManagedClusterSetBindingInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredManagedClusterSetBindingInformer constructs a new informer for ManagedClusterSetBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredManagedClusterSetBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ClusterV1alpha1().ManagedClusterSetBindings(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ClusterV1alpha1().ManagedClusterSetBindings(namespace).Watch(context.TODO(), options)
			},
		},
		&clusterv1alpha1.ManagedClusterSetBinding{},
		resyncPeriod,
		indexers,
	)
}

func (f *managedClusterSetBindingInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredManagedClusterSetBindingInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *managedClusterSetBindingInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&clusterv1alpha1.ManagedClusterSetBinding{}, f.defaultInformer)
}

func (f *managedClusterSetBindingInformer) Lister() v1alpha1.ManagedClusterSetBindingLister {
	return v1alpha1.NewManagedClusterSetBindingLister(f.Informer().GetIndexer())
}
