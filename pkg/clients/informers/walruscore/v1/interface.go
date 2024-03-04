// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/seal-io/walrus/pkg/clients/informers/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Catalogs returns a CatalogInformer.
	Catalogs() CatalogInformer
	// Connectors returns a ConnectorInformer.
	Connectors() ConnectorInformer
	// Resources returns a ResourceInformer.
	Resources() ResourceInformer
	// ResourceDefinitions returns a ResourceDefinitionInformer.
	ResourceDefinitions() ResourceDefinitionInformer
	// ResourceRuns returns a ResourceRunInformer.
	ResourceRuns() ResourceRunInformer
	// Templates returns a TemplateInformer.
	Templates() TemplateInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Catalogs returns a CatalogInformer.
func (v *version) Catalogs() CatalogInformer {
	return &catalogInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Connectors returns a ConnectorInformer.
func (v *version) Connectors() ConnectorInformer {
	return &connectorInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Resources returns a ResourceInformer.
func (v *version) Resources() ResourceInformer {
	return &resourceInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// ResourceDefinitions returns a ResourceDefinitionInformer.
func (v *version) ResourceDefinitions() ResourceDefinitionInformer {
	return &resourceDefinitionInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// ResourceRuns returns a ResourceRunInformer.
func (v *version) ResourceRuns() ResourceRunInformer {
	return &resourceRunInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Templates returns a TemplateInformer.
func (v *version) Templates() TemplateInformer {
	return &templateInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
