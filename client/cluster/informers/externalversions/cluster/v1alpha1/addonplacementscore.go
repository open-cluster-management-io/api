// Copyright Contributors to the Open Cluster Management project
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	versioned "open-cluster-management.io/api/client/cluster/clientset/versioned"
	internalinterfaces "open-cluster-management.io/api/client/cluster/informers/externalversions/internalinterfaces"
	clusterv1alpha1 "open-cluster-management.io/api/client/cluster/listers/cluster/v1alpha1"
	apiclusterv1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
)

// AddOnPlacementScoreInformer provides access to a shared informer and lister for
// AddOnPlacementScores.
type AddOnPlacementScoreInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() clusterv1alpha1.AddOnPlacementScoreLister
}

type addOnPlacementScoreInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewAddOnPlacementScoreInformer constructs a new informer for AddOnPlacementScore type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewAddOnPlacementScoreInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredAddOnPlacementScoreInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredAddOnPlacementScoreInformer constructs a new informer for AddOnPlacementScore type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredAddOnPlacementScoreInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ClusterV1alpha1().AddOnPlacementScores(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ClusterV1alpha1().AddOnPlacementScores(namespace).Watch(context.TODO(), options)
			},
		},
		&apiclusterv1alpha1.AddOnPlacementScore{},
		resyncPeriod,
		indexers,
	)
}

func (f *addOnPlacementScoreInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredAddOnPlacementScoreInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *addOnPlacementScoreInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiclusterv1alpha1.AddOnPlacementScore{}, f.defaultInformer)
}

func (f *addOnPlacementScoreInformer) Lister() clusterv1alpha1.AddOnPlacementScoreLister {
	return clusterv1alpha1.NewAddOnPlacementScoreLister(f.Informer().GetIndexer())
}
