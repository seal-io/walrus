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

// ConnectorHandler handles v1.Connector objects.
//
// ConnectorHandler proxies the v1.Connector objects to the walrus core.
type ConnectorHandler struct {
	extensionapi.ObjectInfo
	extensionapi.CurdOperations
}

func (h *ConnectorHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Declare GVR.
	gvr = walrus.SchemeGroupVersionResource("connectors")

	// As storage.
	h.ObjectInfo = &walrus.Connector{}
	h.CurdOperations = extensionapi.WithCurdProxy[
		*walrus.Connector, *walrus.ConnectorList, *walruscore.Connector, *walruscore.ConnectorList,
	](nil, h, opts.Manager.GetClient().(ctrlcli.WithWatch))

	return
}

var (
	_ rest.Storage           = (*ConnectorHandler)(nil)
	_ rest.Creater           = (*ConnectorHandler)(nil)
	_ rest.Lister            = (*ConnectorHandler)(nil)
	_ rest.Watcher           = (*ConnectorHandler)(nil)
	_ rest.Getter            = (*ConnectorHandler)(nil)
	_ rest.Updater           = (*ConnectorHandler)(nil)
	_ rest.GracefulDeleter   = (*ConnectorHandler)(nil)
	_ rest.CollectionDeleter = (*ConnectorHandler)(nil)
)

func (h *ConnectorHandler) Destroy() {
}

func (h *ConnectorHandler) New() runtime.Object {
	return &walrus.Connector{}
}

func (h *ConnectorHandler) NewList() runtime.Object {
	return &walrus.ConnectorList{}
}

func (h *ConnectorHandler) CastObjectTo(do *walrus.Connector) (uo *walruscore.Connector) {
	return (*walruscore.Connector)(do)
}

func (h *ConnectorHandler) CastObjectFrom(uo *walruscore.Connector) (do *walrus.Connector) {
	return (*walrus.Connector)(uo)
}

func (h *ConnectorHandler) CastObjectListTo(dol *walrus.ConnectorList) (uol *walruscore.ConnectorList) {
	return (*walruscore.ConnectorList)(dol)
}

func (h *ConnectorHandler) CastObjectListFrom(uol *walruscore.ConnectorList) (dol *walrus.ConnectorList) {
	return (*walrus.ConnectorList)(uol)
}
