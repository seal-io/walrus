package extensionapi

import (
	"context"

	autoscaling "k8s.io/api/autoscaling/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	// SetupOptions is the options for setting up a handler.
	SetupOptions struct {
		// Manager is the controller-runtime manager.
		Manager ctrl.Manager
	}

	// Setup is an interface for setting up a handler.
	Setup interface {
		rest.Storage
		// SetupHandler sets up the handler.
		//
		// SetupHandler is called before the Cache is started,
		// you should not do anything that requires the Cache to be started.
		// Instead, you can configure the Cache, like IndexField or something else.
		SetupHandler(ctx context.Context, opts SetupOptions) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) // nolint: lll
	}
)

type (
	// ObjectInfo is an interface for informing extension API objects.
	ObjectInfo interface {
		rest.Scoper
		rest.KindProvider
		rest.SingularNameProvider
		rest.CategoriesProvider
		rest.ShortNamesProvider
	}

	// ObjectWithStatusSubResource is an interface for extension API objects with status subresources,
	// which indicates the resource has a status subresource.
	ObjectWithStatusSubResource interface {
		meta.ObjectMetaAccessor
		runtime.Object
		CopyStatusTo(runtime.Object)
	}

	// ObjectWithScaleSubResource is an interface for extension API objects with scale subresources,
	// which indicates the resource has a scale subresource.
	ObjectWithScaleSubResource interface {
		meta.ObjectMetaAccessor
		runtime.Object
		GetScale() *autoscaling.Scale
		SetScale(*autoscaling.Scale)
	}
)

type (
	// MetaObject is the interface for the object with metadata.
	MetaObject = ctrlcli.Object
	// MetaObjectList is the interface for the list of objects with metadata.
	MetaObjectList = ctrlcli.ObjectList
)
