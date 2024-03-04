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

// ResourceRunHandler handles v1.ResourceRun objects.
//
// ResourceRunHandler proxies the v1.ResourceRun objects to the walrus core.
type ResourceRunHandler struct {
	extensionapi.ObjectInfo
	extensionapi.CurdOperations
}

func (h *ResourceRunHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Declare GVR.
	gvr = walrus.SchemeGroupVersionResource("resourceruns")

	// As storage.
	h.ObjectInfo = &walrus.ResourceRun{}
	h.CurdOperations = extensionapi.WithCurdProxy[
		*walrus.ResourceRun, *walrus.ResourceRunList, *walruscore.ResourceRun, *walruscore.ResourceRunList,
	](nil, h, opts.Manager.GetClient().(ctrlcli.WithWatch))

	return
}

var (
	_ rest.Storage           = (*ResourceRunHandler)(nil)
	_ rest.Creater           = (*ResourceRunHandler)(nil)
	_ rest.Lister            = (*ResourceRunHandler)(nil)
	_ rest.Watcher           = (*ResourceRunHandler)(nil)
	_ rest.Getter            = (*ResourceRunHandler)(nil)
	_ rest.Updater           = (*ResourceRunHandler)(nil)
	_ rest.Patcher           = (*ResourceRunHandler)(nil)
	_ rest.GracefulDeleter   = (*ResourceRunHandler)(nil)
	_ rest.CollectionDeleter = (*ResourceRunHandler)(nil)
)

func (h *ResourceRunHandler) Destroy() {
}

func (h *ResourceRunHandler) New() runtime.Object {
	return &walrus.ResourceRun{}
}

func (h *ResourceRunHandler) NewList() runtime.Object {
	return &walrus.ResourceRunList{}
}

func (h *ResourceRunHandler) CastObjectTo(do *walrus.ResourceRun) (uo *walruscore.ResourceRun) {
	return (*walruscore.ResourceRun)(do)
}

func (h *ResourceRunHandler) CastObjectFrom(uo *walruscore.ResourceRun) (do *walrus.ResourceRun) {
	return (*walrus.ResourceRun)(uo)
}

func (h *ResourceRunHandler) CastObjectListTo(dol *walrus.ResourceRunList) (uol *walruscore.ResourceRunList) {
	return (*walruscore.ResourceRunList)(dol)
}

func (h *ResourceRunHandler) CastObjectListFrom(uol *walruscore.ResourceRunList) (dol *walrus.ResourceRunList) {
	return (*walrus.ResourceRunList)(uol)
}
