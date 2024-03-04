package walrus

import (
	"context"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"github.com/seal-io/walrus/pkg/extensionapi"
	"github.com/seal-io/walrus/pkg/systemkuberes"
)

// ResourceDefinitionHandler handles v1.ResourceDefinition objects.
//
// ResourceDefinitionHandler proxies the v1.ResourceDefinition objects to the walrus core.
type ResourceDefinitionHandler struct {
	extensionapi.ObjectInfo
	extensionapi.CurdOperations
}

func (h *ResourceDefinitionHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Declare GVR.
	gvr = walrus.SchemeGroupVersionResource("resourcedefinitions")

	// As storage.
	h.ObjectInfo = &walrus.ResourceDefinition{}
	h.CurdOperations = extensionapi.WithCurdProxy[
		*walrus.ResourceDefinition, *walrus.ResourceDefinitionList, *walruscore.ResourceDefinition, *walruscore.ResourceDefinitionList,
	](nil, h, opts.Manager.GetClient().(ctrlcli.WithWatch))

	return
}

var (
	_ rest.Storage           = (*ResourceDefinitionHandler)(nil)
	_ rest.Creater           = (*ResourceDefinitionHandler)(nil)
	_ rest.Lister            = (*ResourceDefinitionHandler)(nil)
	_ rest.Watcher           = (*ResourceDefinitionHandler)(nil)
	_ rest.Getter            = (*ResourceDefinitionHandler)(nil)
	_ rest.Updater           = (*ResourceDefinitionHandler)(nil)
	_ rest.GracefulDeleter   = (*ResourceDefinitionHandler)(nil)
	_ rest.CollectionDeleter = (*ResourceDefinitionHandler)(nil)
)

func (h *ResourceDefinitionHandler) BeforeOnCreate(ctx context.Context, obj *walrus.ResourceDefinition, opts *ctrlcli.CreateOptions) error {
	if obj.Namespace != systemkuberes.SystemNamespaceName {
		errs := field.ErrorList{
			field.Invalid(field.NewPath("metadata.namespace"), obj.Namespace,
				"resource definition namespace must be "+systemkuberes.SystemNamespaceName),
		}
		return kerrors.NewInvalid(walrus.SchemeKind("resourcedefinitions"), obj.Name, errs)
	}
	return nil
}

func (h *ResourceDefinitionHandler) BeforeOnGet(ctx context.Context, key types.NamespacedName, opts *ctrlcli.GetOptions) error {
	if key.Namespace != systemkuberes.SystemNamespaceName {
		return kerrors.NewNotFound(walrus.SchemeResource("resourcedefinitions"), key.Name)
	}
	return nil
}

func (h *ResourceDefinitionHandler) BeforeOnListWatch(ctx context.Context, opts *ctrlcli.ListOptions) error {
	if opts.Namespace == "" {
		opts.Namespace = systemkuberes.SystemNamespaceName
	}
	return nil
}

func (h *ResourceDefinitionHandler) Destroy() {
}

func (h *ResourceDefinitionHandler) New() runtime.Object {
	return &walrus.ResourceDefinition{}
}

func (h *ResourceDefinitionHandler) NewList() runtime.Object {
	return &walrus.ResourceDefinitionList{}
}

func (h *ResourceDefinitionHandler) CastObjectTo(do *walrus.ResourceDefinition) (uo *walruscore.ResourceDefinition) {
	return (*walruscore.ResourceDefinition)(do)
}

func (h *ResourceDefinitionHandler) CastObjectFrom(uo *walruscore.ResourceDefinition) (do *walrus.ResourceDefinition) {
	return (*walrus.ResourceDefinition)(uo)
}

func (h *ResourceDefinitionHandler) CastObjectListTo(dol *walrus.ResourceDefinitionList) (uol *walruscore.ResourceDefinitionList) {
	return (*walruscore.ResourceDefinitionList)(dol)
}

func (h *ResourceDefinitionHandler) CastObjectListFrom(uol *walruscore.ResourceDefinitionList) (dol *walrus.ResourceDefinitionList) {
	return (*walrus.ResourceDefinitionList)(uol)
}
