// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

import (
	v1 "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ConnectorLister helps list Connectors.
// All objects returned here must be treated as read-only.
type ConnectorLister interface {
	// List lists all Connectors in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Connector, err error)
	// Connectors returns an object that can list and get Connectors.
	Connectors(namespace string) ConnectorNamespaceLister
	ConnectorListerExpansion
}

// connectorLister implements the ConnectorLister interface.
type connectorLister struct {
	indexer cache.Indexer
}

// NewConnectorLister returns a new ConnectorLister.
func NewConnectorLister(indexer cache.Indexer) ConnectorLister {
	return &connectorLister{indexer: indexer}
}

// List lists all Connectors in the indexer.
func (s *connectorLister) List(selector labels.Selector) (ret []*v1.Connector, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Connector))
	})
	return ret, err
}

// Connectors returns an object that can list and get Connectors.
func (s *connectorLister) Connectors(namespace string) ConnectorNamespaceLister {
	return connectorNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ConnectorNamespaceLister helps list and get Connectors.
// All objects returned here must be treated as read-only.
type ConnectorNamespaceLister interface {
	// List lists all Connectors in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Connector, err error)
	// Get retrieves the Connector from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Connector, error)
	ConnectorNamespaceListerExpansion
}

// connectorNamespaceLister implements the ConnectorNamespaceLister
// interface.
type connectorNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Connectors in the indexer for a given namespace.
func (s connectorNamespaceLister) List(selector labels.Selector) (ret []*v1.Connector, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Connector))
	})
	return ret, err
}

// Get retrieves the Connector from the indexer for a given namespace and name.
func (s connectorNamespaceLister) Get(name string) (*v1.Connector, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.SchemeResource("connector"), name)
	}
	return obj.(*v1.Connector), nil
}
