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
	apiaddonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	versioned "open-cluster-management.io/api/client/addon/clientset/versioned"
	internalinterfaces "open-cluster-management.io/api/client/addon/informers/externalversions/internalinterfaces"
	addonv1alpha1 "open-cluster-management.io/api/client/addon/listers/addon/v1alpha1"
)

// AddOnDeploymentConfigInformer provides access to a shared informer and lister for
// AddOnDeploymentConfigs.
type AddOnDeploymentConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() addonv1alpha1.AddOnDeploymentConfigLister
}

type addOnDeploymentConfigInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewAddOnDeploymentConfigInformer constructs a new informer for AddOnDeploymentConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewAddOnDeploymentConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredAddOnDeploymentConfigInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredAddOnDeploymentConfigInformer constructs a new informer for AddOnDeploymentConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredAddOnDeploymentConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AddonV1alpha1().AddOnDeploymentConfigs(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AddonV1alpha1().AddOnDeploymentConfigs(namespace).Watch(context.TODO(), options)
			},
		},
		&apiaddonv1alpha1.AddOnDeploymentConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *addOnDeploymentConfigInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredAddOnDeploymentConfigInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *addOnDeploymentConfigInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiaddonv1alpha1.AddOnDeploymentConfig{}, f.defaultInformer)
}

func (f *addOnDeploymentConfigInformer) Lister() addonv1alpha1.AddOnDeploymentConfigLister {
	return addonv1alpha1.NewAddOnDeploymentConfigLister(f.Informer().GetIndexer())
}
