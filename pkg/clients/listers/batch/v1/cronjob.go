// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

import (
	v1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CronJobLister helps list CronJobs.
// All objects returned here must be treated as read-only.
type CronJobLister interface {
	// List lists all CronJobs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.CronJob, err error)
	// CronJobs returns an object that can list and get CronJobs.
	CronJobs(namespace string) CronJobNamespaceLister
	CronJobListerExpansion
}

// cronJobLister implements the CronJobLister interface.
type cronJobLister struct {
	indexer cache.Indexer
}

// NewCronJobLister returns a new CronJobLister.
func NewCronJobLister(indexer cache.Indexer) CronJobLister {
	return &cronJobLister{indexer: indexer}
}

// List lists all CronJobs in the indexer.
func (s *cronJobLister) List(selector labels.Selector) (ret []*v1.CronJob, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CronJob))
	})
	return ret, err
}

// CronJobs returns an object that can list and get CronJobs.
func (s *cronJobLister) CronJobs(namespace string) CronJobNamespaceLister {
	return cronJobNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CronJobNamespaceLister helps list and get CronJobs.
// All objects returned here must be treated as read-only.
type CronJobNamespaceLister interface {
	// List lists all CronJobs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.CronJob, err error)
	// Get retrieves the CronJob from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.CronJob, error)
	CronJobNamespaceListerExpansion
}

// cronJobNamespaceLister implements the CronJobNamespaceLister
// interface.
type cronJobNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CronJobs in the indexer for a given namespace.
func (s cronJobNamespaceLister) List(selector labels.Selector) (ret []*v1.CronJob, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CronJob))
	})
	return ret, err
}

// Get retrieves the CronJob from the indexer for a given namespace and name.
func (s cronJobNamespaceLister) Get(name string) (*v1.CronJob, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.SchemeResource("cronjob"), name)
	}
	return obj.(*v1.CronJob), nil
}
