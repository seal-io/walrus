package walrus

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"github.com/seal-io/walrus/pkg/extensionapi"
)

// CatalogHandler handles v1.Catalog objects.
//
// CatalogHandler proxies the v1.Catalog objects to the walrus core.
type CatalogHandler struct {
	extensionapi.ObjectInfo
	extensionapi.CurdOperations
}

func (h *CatalogHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Declare GVR.
	gvr = walrus.SchemeGroupVersionResource("catalogs")

	// As storage.
	h.ObjectInfo = &walrus.Catalog{}
	h.CurdOperations = extensionapi.WithCurdProxy[
		*walrus.Catalog, *walrus.CatalogList, *walruscore.Catalog, *walruscore.CatalogList,
	](nil, h, opts.Manager.GetClient().(ctrlcli.WithWatch))

	return
}

var (
	_ rest.Storage           = (*CatalogHandler)(nil)
	_ rest.Creater           = (*CatalogHandler)(nil)
	_ rest.Lister            = (*CatalogHandler)(nil)
	_ rest.Watcher           = (*CatalogHandler)(nil)
	_ rest.Getter            = (*CatalogHandler)(nil)
	_ rest.Updater           = (*CatalogHandler)(nil)
	_ rest.GracefulDeleter   = (*CatalogHandler)(nil)
	_ rest.CollectionDeleter = (*CatalogHandler)(nil)
)

func (h *CatalogHandler) Destroy() {
}

func (h *CatalogHandler) New() runtime.Object {
	return &walrus.Catalog{}
}

func (h *CatalogHandler) NewList() runtime.Object {
	return &walrus.CatalogList{}
}

func (h *CatalogHandler) CastObjectTo(do *walrus.Catalog) (uo *walruscore.Catalog) {
	return (*walruscore.Catalog)(do)
}

func (h *CatalogHandler) CastObjectFrom(uo *walruscore.Catalog) (do *walrus.Catalog) {
	return (*walrus.Catalog)(uo)
}

func (h *CatalogHandler) CastObjectListTo(dol *walrus.CatalogList) (uol *walruscore.CatalogList) {
	return (*walruscore.CatalogList)(dol)
}

func (h *CatalogHandler) CastObjectListFrom(uol *walruscore.CatalogList) (dol *walrus.CatalogList) {
	return (*walrus.CatalogList)(uol)
}
