// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	versioned "github.com/openshift-knative/serverless-operator/pkg/client/clientset/versioned"
	internalinterfaces "github.com/openshift-knative/serverless-operator/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/openshift-knative/serverless-operator/pkg/client/listers/config/v1"
	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ImageContentPolicyInformer provides access to a shared informer and lister for
// ImageContentPolicies.
type ImageContentPolicyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ImageContentPolicyLister
}

type imageContentPolicyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewImageContentPolicyInformer constructs a new informer for ImageContentPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewImageContentPolicyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredImageContentPolicyInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredImageContentPolicyInformer constructs a new informer for ImageContentPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredImageContentPolicyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().ImageContentPolicies().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().ImageContentPolicies().Watch(context.TODO(), options)
			},
		},
		&configv1.ImageContentPolicy{},
		resyncPeriod,
		indexers,
	)
}

func (f *imageContentPolicyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredImageContentPolicyInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *imageContentPolicyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configv1.ImageContentPolicy{}, f.defaultInformer)
}

func (f *imageContentPolicyInformer) Lister() v1.ImageContentPolicyLister {
	return v1.NewImageContentPolicyLister(f.Informer().GetIndexer())
}
