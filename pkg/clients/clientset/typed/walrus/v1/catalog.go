// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package v1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1 "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	walrusv1 "github.com/seal-io/walrus/pkg/clients/applyconfiguration/walrus/v1"
	scheme "github.com/seal-io/walrus/pkg/clients/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CatalogsGetter has a method to return a CatalogInterface.
// A group's client should implement this interface.
type CatalogsGetter interface {
	Catalogs(namespace string) CatalogInterface
}

// CatalogInterface has methods to work with Catalog resources.
type CatalogInterface interface {
	Create(ctx context.Context, catalog *v1.Catalog, opts metav1.CreateOptions) (*v1.Catalog, error)
	Update(ctx context.Context, catalog *v1.Catalog, opts metav1.UpdateOptions) (*v1.Catalog, error)
	UpdateStatus(ctx context.Context, catalog *v1.Catalog, opts metav1.UpdateOptions) (*v1.Catalog, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Catalog, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.CatalogList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Catalog, err error)
	Apply(ctx context.Context, catalog *walrusv1.CatalogApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Catalog, err error)
	ApplyStatus(ctx context.Context, catalog *walrusv1.CatalogApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Catalog, err error)
	CatalogExpansion
}

// catalogs implements CatalogInterface
type catalogs struct {
	client rest.Interface
	ns     string
}

// newCatalogs returns a Catalogs
func newCatalogs(c *WalrusV1Client, namespace string) *catalogs {
	return &catalogs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the catalog, and returns the corresponding catalog object, and an error if there is any.
func (c *catalogs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Catalog, err error) {
	result = &v1.Catalog{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("catalogs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Catalogs that match those selectors.
func (c *catalogs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.CatalogList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.CatalogList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("catalogs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested catalogs.
func (c *catalogs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("catalogs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a catalog and creates it.  Returns the server's representation of the catalog, and an error, if there is any.
func (c *catalogs) Create(ctx context.Context, catalog *v1.Catalog, opts metav1.CreateOptions) (result *v1.Catalog, err error) {
	result = &v1.Catalog{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("catalogs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(catalog).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a catalog and updates it. Returns the server's representation of the catalog, and an error, if there is any.
func (c *catalogs) Update(ctx context.Context, catalog *v1.Catalog, opts metav1.UpdateOptions) (result *v1.Catalog, err error) {
	result = &v1.Catalog{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("catalogs").
		Name(catalog.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(catalog).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *catalogs) UpdateStatus(ctx context.Context, catalog *v1.Catalog, opts metav1.UpdateOptions) (result *v1.Catalog, err error) {
	result = &v1.Catalog{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("catalogs").
		Name(catalog.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(catalog).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the catalog and deletes it. Returns an error if one occurs.
func (c *catalogs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("catalogs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *catalogs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("catalogs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched catalog.
func (c *catalogs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Catalog, err error) {
	result = &v1.Catalog{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("catalogs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied catalog.
func (c *catalogs) Apply(ctx context.Context, catalog *walrusv1.CatalogApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Catalog, err error) {
	if catalog == nil {
		return nil, fmt.Errorf("catalog provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(catalog)
	if err != nil {
		return nil, err
	}
	name := catalog.Name
	if name == nil {
		return nil, fmt.Errorf("catalog.Name must be provided to Apply")
	}
	result = &v1.Catalog{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("catalogs").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *catalogs) ApplyStatus(ctx context.Context, catalog *walrusv1.CatalogApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Catalog, err error) {
	if catalog == nil {
		return nil, fmt.Errorf("catalog provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(catalog)
	if err != nil {
		return nil, err
	}

	name := catalog.Name
	if name == nil {
		return nil, fmt.Errorf("catalog.Name must be provided to Apply")
	}

	result = &v1.Catalog{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("catalogs").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
