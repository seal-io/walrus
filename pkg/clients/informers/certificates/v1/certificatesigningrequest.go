// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	clientset "github.com/seal-io/walrus/pkg/clients/clientset"
	internalinterfaces "github.com/seal-io/walrus/pkg/clients/informers/internalinterfaces"
	v1 "github.com/seal-io/walrus/pkg/clients/listers/certificates/v1"
	certificatesv1 "k8s.io/api/certificates/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CertificateSigningRequestInformer provides access to a shared informer and lister for
// CertificateSigningRequests.
type CertificateSigningRequestInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.CertificateSigningRequestLister
}

type certificateSigningRequestInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewCertificateSigningRequestInformer constructs a new informer for CertificateSigningRequest type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCertificateSigningRequestInformer(client clientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCertificateSigningRequestInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredCertificateSigningRequestInformer constructs a new informer for CertificateSigningRequest type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCertificateSigningRequestInformer(client clientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CertificatesV1().CertificateSigningRequests().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CertificatesV1().CertificateSigningRequests().Watch(context.TODO(), options)
			},
		},
		&certificatesv1.CertificateSigningRequest{},
		resyncPeriod,
		indexers,
	)
}

func (f *certificateSigningRequestInformer) defaultInformer(client clientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCertificateSigningRequestInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *certificateSigningRequestInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&certificatesv1.CertificateSigningRequest{}, f.defaultInformer)
}

func (f *certificateSigningRequestInformer) Lister() v1.CertificateSigningRequestLister {
	return v1.NewCertificateSigningRequestLister(f.Informer().GetIndexer())
}