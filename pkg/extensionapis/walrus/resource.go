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

// ResourceHandler handles v1.Resource objects.
//
// ResourceHandler proxies the v1.Resource objects to the walrus core.
type ResourceHandler struct {
	extensionapi.ObjectInfo
	extensionapi.CurdOperations
}

func (h *ResourceHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Declare GVR.
	gvr = walrus.SchemeGroupVersionResource("resources")

	// As storage.
	h.ObjectInfo = &walrus.Resource{}
	h.CurdOperations = extensionapi.WithCurdProxy[
		*walrus.Resource, *walrus.ResourceList, *walruscore.Resource, *walruscore.ResourceList,
	](nil, h, opts.Manager.GetClient().(ctrlcli.WithWatch))

	return
}

var (
	_ rest.Storage           = (*ResourceHandler)(nil)
	_ rest.Creater           = (*ResourceHandler)(nil)
	_ rest.Lister            = (*ResourceHandler)(nil)
	_ rest.Watcher           = (*ResourceHandler)(nil)
	_ rest.Getter            = (*ResourceHandler)(nil)
	_ rest.Updater           = (*ResourceHandler)(nil)
	_ rest.GracefulDeleter   = (*ResourceHandler)(nil)
	_ rest.CollectionDeleter = (*ResourceHandler)(nil)
)

func (h *ResourceHandler) Destroy() {
}

func (h *ResourceHandler) New() runtime.Object {
	return &walrus.Resource{}
}

func (h *ResourceHandler) NewList() runtime.Object {
	return &walrus.ResourceList{}
}

func (h *ResourceHandler) CastObjectTo(do *walrus.Resource) (uo *walruscore.Resource) {
	return (*walruscore.Resource)(do)
}

func (h *ResourceHandler) CastObjectFrom(uo *walruscore.Resource) (do *walrus.Resource) {
	return (*walrus.Resource)(uo)
}

func (h *ResourceHandler) CastObjectListTo(dol *walrus.ResourceList) (uol *walruscore.ResourceList) {
	return (*walruscore.ResourceList)(dol)
}

func (h *ResourceHandler) CastObjectListFrom(uol *walruscore.ResourceList) (dol *walrus.ResourceList) {
	return (*walrus.ResourceList)(uol)
}
