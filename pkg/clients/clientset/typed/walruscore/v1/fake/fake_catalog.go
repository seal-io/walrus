// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus", DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1 "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	walruscorev1 "github.com/seal-io/walrus/pkg/clients/applyconfiguration/walruscore/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCatalogs implements CatalogInterface
type FakeCatalogs struct {
	Fake *FakeWalruscoreV1
	ns   string
}

var catalogsResource = v1.SchemeGroupVersion.WithResource("catalogs")

var catalogsKind = v1.SchemeGroupVersion.WithKind("Catalog")

// Get takes name of the catalog, and returns the corresponding catalog object, and an error if there is any.
func (c *FakeCatalogs) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Catalog, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(catalogsResource, c.ns, name), &v1.Catalog{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Catalog), err
}

// List takes label and field selectors, and returns the list of Catalogs that match those selectors.
func (c *FakeCatalogs) List(ctx context.Context, opts metav1.ListOptions) (result *v1.CatalogList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(catalogsResource, catalogsKind, c.ns, opts), &v1.CatalogList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.CatalogList{ListMeta: obj.(*v1.CatalogList).ListMeta}
	for _, item := range obj.(*v1.CatalogList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested catalogs.
func (c *FakeCatalogs) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(catalogsResource, c.ns, opts))

}

// Create takes the representation of a catalog and creates it.  Returns the server's representation of the catalog, and an error, if there is any.
func (c *FakeCatalogs) Create(ctx context.Context, catalog *v1.Catalog, opts metav1.CreateOptions) (result *v1.Catalog, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(catalogsResource, c.ns, catalog), &v1.Catalog{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Catalog), err
}

// Update takes the representation of a catalog and updates it. Returns the server's representation of the catalog, and an error, if there is any.
func (c *FakeCatalogs) Update(ctx context.Context, catalog *v1.Catalog, opts metav1.UpdateOptions) (result *v1.Catalog, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(catalogsResource, c.ns, catalog), &v1.Catalog{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Catalog), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCatalogs) UpdateStatus(ctx context.Context, catalog *v1.Catalog, opts metav1.UpdateOptions) (*v1.Catalog, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(catalogsResource, "status", c.ns, catalog), &v1.Catalog{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Catalog), err
}

// Delete takes name of the catalog and deletes it. Returns an error if one occurs.
func (c *FakeCatalogs) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(catalogsResource, c.ns, name, opts), &v1.Catalog{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCatalogs) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(catalogsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.CatalogList{})
	return err
}

// Patch applies the patch and returns the patched catalog.
func (c *FakeCatalogs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Catalog, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(catalogsResource, c.ns, name, pt, data, subresources...), &v1.Catalog{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Catalog), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied catalog.
func (c *FakeCatalogs) Apply(ctx context.Context, catalog *walruscorev1.CatalogApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Catalog, err error) {
	if catalog == nil {
		return nil, fmt.Errorf("catalog provided to Apply must not be nil")
	}
	data, err := json.Marshal(catalog)
	if err != nil {
		return nil, err
	}
	name := catalog.Name
	if name == nil {
		return nil, fmt.Errorf("catalog.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(catalogsResource, c.ns, *name, types.ApplyPatchType, data), &v1.Catalog{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Catalog), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeCatalogs) ApplyStatus(ctx context.Context, catalog *walruscorev1.CatalogApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Catalog, err error) {
	if catalog == nil {
		return nil, fmt.Errorf("catalog provided to Apply must not be nil")
	}
	data, err := json.Marshal(catalog)
	if err != nil {
		return nil, err
	}
	name := catalog.Name
	if name == nil {
		return nil, fmt.Errorf("catalog.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(catalogsResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1.Catalog{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Catalog), err
}