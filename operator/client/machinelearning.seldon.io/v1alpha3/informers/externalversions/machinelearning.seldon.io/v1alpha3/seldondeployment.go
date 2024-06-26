/*
Copyright (c) 2024 Seldon Technologies Ltd.

Use of this software is governed BY
(1) the license included in the LICENSE file or
(2) if the license included in the LICENSE file is the Business Source License 1.1,
the Change License after the Change Date as each is defined in accordance with the LICENSE file.
*/
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha3

import (
	"context"
	time "time"

	machinelearningseldoniov1alpha3 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1alpha3"
	versioned "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1alpha3/clientset/versioned"
	internalinterfaces "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1alpha3/informers/externalversions/internalinterfaces"
	v1alpha3 "github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1alpha3/listers/machinelearning.seldon.io/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SeldonDeploymentInformer provides access to a shared informer and lister for
// SeldonDeployments.
type SeldonDeploymentInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha3.SeldonDeploymentLister
}

type seldonDeploymentInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSeldonDeploymentInformer constructs a new informer for SeldonDeployment type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSeldonDeploymentInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSeldonDeploymentInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSeldonDeploymentInformer constructs a new informer for SeldonDeployment type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSeldonDeploymentInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MachinelearningV1alpha3().SeldonDeployments(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MachinelearningV1alpha3().SeldonDeployments(namespace).Watch(context.TODO(), options)
			},
		},
		&machinelearningseldoniov1alpha3.SeldonDeployment{},
		resyncPeriod,
		indexers,
	)
}

func (f *seldonDeploymentInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSeldonDeploymentInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *seldonDeploymentInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&machinelearningseldoniov1alpha3.SeldonDeployment{}, f.defaultInformer)
}

func (f *seldonDeploymentInformer) Lister() v1alpha3.SeldonDeploymentLister {
	return v1alpha3.NewSeldonDeploymentLister(f.Informer().GetIndexer())
}
