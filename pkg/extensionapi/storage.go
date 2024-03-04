package extensionapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/seal-io/utils/pools/gopool"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	kmeta "k8s.io/apimachinery/pkg/api/meta"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/handlers/fieldmanager"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/apiserver/pkg/util/dryrun"
	"k8s.io/utils/ptr"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	// CreateHandler is an interface for a creation handler.
	CreateHandler interface {
		// New returns a new object.
		New() runtime.Object
		// OnCreate creates the object with the given options,
		// and returns the created object.
		OnCreate(ctx context.Context, object runtime.Object, opts ctrlcli.CreateOptions) (created runtime.Object, err error)
	}

	// CreateOperation implements the rest.Creater interface.
	CreateOperation struct {
		Handler CreateHandler
	}
)

var _ rest.Creater = (*CreateOperation)(nil)

// WithCreate returns a CreateOperation for the given CreateHandler.
func WithCreate(h CreateHandler) CreateOperation {
	return CreateOperation{Handler: h}
}

func (s CreateOperation) New() runtime.Object {
	return s.Handler.New()
}

func (s CreateOperation) Create(
	ctx context.Context,
	obj runtime.Object,
	createValidation rest.ValidateObjectFunc,
	options *meta.CreateOptions,
) (runtime.Object, error) {
	if options == nil {
		options = &meta.CreateOptions{}
	}

	om, err := kmeta.Accessor(obj)
	if err != nil {
		return nil, kerrors.NewInternalError(err)
	}
	rest.FillObjectMetaSystemFields(om)
	if gn := om.GetGenerateName(); len(gn) > 0 && len(om.GetName()) == 0 {
		om.SetName(names.SimpleNameGenerator.GenerateName(gn))
	}

	if createValidation != nil {
		err = createValidation(ctx, obj.DeepCopyObject())
		if err != nil {
			return nil, err
		}
	}

	if dryrun.IsDryRun(options.DryRun) {
		getter, ok := s.Handler.(rest.Getter)
		if !ok {
			// If the handler does not support get, we cannot check for existence.
			return obj, nil
		}

		_, err = getter.Get(ctx, om.GetName(), &meta.GetOptions{ResourceVersion: "0"})
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return nil, wrapError(ctx, om.GetName(), err)
			}
			return obj, nil
		}

		return nil, kerrors.NewAlreadyExists(qualifiedResourceFromContext(ctx), om.GetName())
	}

	obj, err = s.Handler.OnCreate(ctx, obj, convertCtrlCreateOptionsFromMeta(options))
	if err != nil {
		return nil, wrapError(ctx, om.GetName(), err)
	}
	return obj, nil
}

type (
	// ListHandler is an interface for a listing handler.
	ListHandler interface {
		// NewList returns an empty object that can be used with the List call.
		NewList() runtime.Object
		// OnList returns an object list collection for the given options.
		OnList(ctx context.Context, opts ctrlcli.ListOptions) (list runtime.Object, err error)
	}

	// ListOperation implements the rest.Lister interface.
	ListOperation struct {
		WatchBookmark  WatchBookmark
		Handler        ListHandler
		TableConvertor rest.TableConvertor
	}
)

var _ rest.Lister = (*ListOperation)(nil)

// WithList returns a ListOperation for the given rest.TableConvertor and ListHandler.
func WithList(tc rest.TableConvertor, h ListHandler) ListOperation {
	if tc == nil {
		tc = NewDefaultTableConvertor()
	}
	return ListOperation{
		WatchBookmark:  NewWatchBookmark(),
		Handler:        h,
		TableConvertor: tc,
	}
}

func (s ListOperation) ConvertToTable(ctx context.Context, object, tableOptions runtime.Object) (*meta.Table, error) {
	return s.TableConvertor.ConvertToTable(ctx, object, tableOptions)
}

func (s ListOperation) NewList() runtime.Object {
	return s.Handler.NewList()
}

func (s ListOperation) List(
	ctx context.Context,
	options *metainternal.ListOptions,
) (runtime.Object, error) {
	if options == nil {
		options = &metainternal.ListOptions{ResourceVersion: "0"}
	}

	listObj, err := s.Handler.OnList(ctx, convertCtrlListOptionsFromMeta(ctx, options))
	if err != nil {
		return nil, wrapError(ctx, "", err)
	}

	if listObj == nil {
		listObj = s.Handler.NewList()
	} else {
		// Get the highest revision from the list items.
		items, _ := kmeta.ExtractList(listObj)
		if len(items) != 0 {
			om, err := kmeta.Accessor(items[len(items)-1])
			if err == nil {
				s.WatchBookmark.SwapResourceVersion(om.GetResourceVersion())
			}
		}
	}

	return listObj, nil
}

type (
	// WatchHandler is an interface for a watching handler.
	WatchHandler interface {
		// OnWatch returns a watch.Interface that watches the list of objects.
		OnWatch(ctx context.Context, opts ctrlcli.ListOptions) (watch.Interface, error)
	}

	// ListWatchHandler is an interface for a listing handler.
	ListWatchHandler interface {
		ListHandler
		WatchHandler
	}

	// ListWatchOperation implements the rest.Lister and rest.Watcher interfaces.
	ListWatchOperation struct {
		ListOperation
		Handler WatchHandler
	}
)

var (
	_ rest.Lister  = (*ListWatchOperation)(nil)
	_ rest.Watcher = (*ListWatchOperation)(nil)
)

// WithListWatch returns a ListWatchOperation for the given rest.TableConvertor and ListWatchHandler.
func WithListWatch(tc rest.TableConvertor, h ListWatchHandler) ListWatchOperation {
	return ListWatchOperation{
		ListOperation: WithList(tc, h),
		Handler:       h,
	}
}

func (s ListWatchOperation) Watch(
	ctx context.Context,
	options *metainternal.ListOptions,
) (watch.Interface, error) {
	if options == nil {
		options = &metainternal.ListOptions{ResourceVersion: "0", Watch: true}
	}

	uw, err := s.Handler.OnWatch(ctx, convertCtrlListOptionsFromMeta(ctx, options))
	if err != nil {
		return nil, wrapError(ctx, "", err)
	}

	wm := s.WatchBookmark.DeepCopy()
	c := make(chan watch.Event)
	dw := watch.NewProxyWatcher(c)
	gopool.Go(func() {
		defer close(c)
		defer uw.Stop()

		for {
			select {
			case <-ctx.Done():
				// Cancel by context.
				return
			case <-dw.StopChan():
				// Cancel by downstream.
				return
			case e, ok := <-uw.ResultChan():
				if !ok {
					// Close by upstream.
					return
				}

				// Nothing to do.
				if e.Object == nil {
					c <- e
					continue
				}

				// Process bookmark.
				if e.Type == watch.Bookmark {
					c <- e
					continue
				}

				// Record the highest revision.
				om, err := kmeta.Accessor(e.Object)
				if err == nil {
					swapped := wm.SwapResourceVersion(om.GetResourceVersion())
					if !swapped {
						// If the revision is not updated, ignore the event.
						continue
					}
				}

				c <- e
			}
		}
	})

	return dw, nil
}

type (
	// GetHandler is an interface for a getting handler.
	GetHandler interface {
		rest.Scoper
		// OnGet returns the object with the given key and options.
		OnGet(ctx context.Context, key types.NamespacedName, opts ctrlcli.GetOptions) (object runtime.Object, err error)
	}

	// GetOperation implements the rest.Getter interface.
	GetOperation struct {
		Handler GetHandler
	}
)

var _ rest.Getter = (*GetOperation)(nil)

// WithGet returns a GetOperation for the given GetHandler.
func WithGet(h GetHandler) GetOperation {
	return GetOperation{Handler: h}
}

func (s GetOperation) Get(
	ctx context.Context,
	name string,
	options *meta.GetOptions,
) (runtime.Object, error) {
	if options == nil {
		options = &meta.GetOptions{ResourceVersion: "0"}
	}

	keyFunc := keyFuncForClusterScope
	if s.Handler.NamespaceScoped() {
		keyFunc = keyFuncForNamespacedScope
	}
	key, err := keyFunc(ctx, name)
	if err != nil {
		return nil, err
	}

	obj, err := s.Handler.OnGet(ctx, key, convertCtrlGetOptionsFromMeta(options))
	if err != nil {
		return nil, wrapError(ctx, name, err)
	}
	return obj, nil
}

type (
	// UpdateHandler is an interface for an updating handler.
	UpdateHandler interface {
		GetHandler
		// New returns a new object.
		New() runtime.Object
		// OnUpdate updates the object with the given options,
		// and returns the updated object.
		OnUpdate(ctx context.Context, object, oldObject runtime.Object, opts ctrlcli.UpdateOptions) (updated runtime.Object, err error)
	}

	// BeforeUpdateFunc is a function that is called before update.
	BeforeUpdateFunc = func(ctx context.Context, object, existing runtime.Object) (runtime.Object, error)

	// UpdateOperation implements the rest.Updater interface.
	UpdateOperation struct {
		Handler      UpdateHandler
		BeforeUpdate BeforeUpdateFunc
	}
)

var _ rest.Updater = (*UpdateOperation)(nil)

// WithUpdate returns a UpdateOperation for the given UpdateHandler.
func WithUpdate(h UpdateHandler) UpdateOperation {
	var bf BeforeUpdateFunc
	if _, ok := h.New().(ObjectWithStatusSubResource); ok {
		bf = beforeUpdateFuncForPreventStatusModify
	}
	return UpdateOperation{
		Handler:      h,
		BeforeUpdate: bf,
	}
}

func (s UpdateOperation) New() runtime.Object {
	return s.Handler.New()
}

func (s UpdateOperation) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *meta.UpdateOptions,
) (runtime.Object, bool, error) {
	if options == nil {
		options = &meta.UpdateOptions{}
	}

	keyFunc := keyFuncForClusterScope
	if s.Handler.NamespaceScoped() {
		keyFunc = keyFuncForNamespacedScope
	}
	key, err := keyFunc(ctx, name)
	if err != nil {
		return nil, false, err
	}

	var creating bool
	existing, err := s.Handler.OnGet(ctx, key, ctrlcli.GetOptions{})
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, false, wrapError(ctx, name, err)
		}
		if !forceAllowCreate {
			return nil, false, kerrors.NewNotFound(qualifiedResourceFromContext(ctx), name)
		}
		existing = s.Handler.New()
		creating = true
	}

	obj, err := objInfo.UpdatedObject(ctx, existing)
	if err != nil {
		return nil, false, err
	}

	if creating {
		createHandler, ok := s.Handler.(CreateHandler)
		if !ok {
			// If the handler does not support get, just return an error.
			return nil, false, kerrors.NewNotFound(qualifiedResourceFromContext(ctx), name)
		}

		creator := WithCreate(createHandler)
		obj, err = creator.Create(ctx, obj, createValidation, newCreateOptionsFromUpdateOptions(options))
		if err != nil {
			return nil, false, err
		}
		return obj, true, nil
	}

	om, err := kmeta.Accessor(obj)
	if err != nil {
		return nil, false, kerrors.NewInternalError(err)
	}

	if om.GetResourceVersion() == "" {
		gk := qualifiedKindFromContext(ctx)
		dt := field.Invalid(
			field.NewPath("metadata", "resourceVersion"),
			0,
			"must be specified for an update",
		)
		return nil, false, kerrors.NewInvalid(gk, name, field.ErrorList{dt})
	}

	if s.BeforeUpdate != nil {
		obj, err = s.BeforeUpdate(ctx, obj, existing)
		if err != nil {
			return nil, false, kerrors.NewInternalError(err)
		}
	}

	obj, err = fieldmanager.IgnoreManagedFieldsTimestampsTransformer(ctx, obj, existing)
	if err != nil {
		return nil, false, err
	}

	if updateValidation != nil {
		err = updateValidation(ctx, obj.DeepCopyObject(), existing.DeepCopyObject())
		if err != nil {
			return nil, false, err
		}
	}

	if dryrun.IsDryRun(options.DryRun) {
		return obj, false, nil
	}

	obj, err = s.Handler.OnUpdate(ctx, obj, existing, convertCtrlUpdateOptionsFromMeta(options))
	if err != nil {
		return nil, false, wrapError(ctx, name, err)
	}
	return obj, false, err
}

func (s UpdateOperation) GetStatusSubResourceUpdater() UpdateOperation {
	s.BeforeUpdate = beforeUpdateFuncForStatusModifyOnly
	return s
}

type (
	// CreateUpdateHandler is an interface for creating and updating handler.
	CreateUpdateHandler interface {
		GetHandler
		// New returns a new object.
		New() runtime.Object
		// OnCreate creates the object with the given options,
		// and returns the created object.
		OnCreate(ctx context.Context, object runtime.Object, opts ctrlcli.CreateOptions) (created runtime.Object, err error)
		// OnUpdate updates the object with the given options,
		// and returns the updated object.
		OnUpdate(ctx context.Context, object, oldObject runtime.Object, opts ctrlcli.UpdateOptions) (updated runtime.Object, err error)
	}

	// CreateUpdateOperation implements the rest.Creater and rest.Updater interface.
	CreateUpdateOperation struct {
		UpdateOperation
		CreateOperation CreateOperation
	}
)

var (
	_ rest.Creater = (*CreateUpdateOperation)(nil)
	_ rest.Updater = (*CreateUpdateOperation)(nil)
)

// WithCreateUpdate returns a CreateUpdateOperation for the given CreateUpdateHandler.
func WithCreateUpdate(h CreateUpdateHandler) CreateUpdateOperation {
	return CreateUpdateOperation{
		UpdateOperation: WithUpdate(h),
		CreateOperation: WithCreate(h),
	}
}

func (s CreateUpdateOperation) Create(
	ctx context.Context,
	obj runtime.Object,
	createValidation rest.ValidateObjectFunc,
	options *meta.CreateOptions,
) (runtime.Object, error) {
	return s.CreateOperation.Create(ctx, obj, createValidation, options)
}

func (s CreateUpdateOperation) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *meta.UpdateOptions,
) (runtime.Object, bool, error) {
	return s.UpdateOperation.Update(ctx, name, objInfo, createValidation, updateValidation, forceAllowCreate, options)
}

type (
	// DeleteHandler is an interface for a deletion handler.
	DeleteHandler = interface {
		GetHandler
		// OnDelete deletes the object with the given options.
		OnDelete(ctx context.Context, object runtime.Object, opts ctrlcli.DeleteOptions) error
	}

	// DeleteOperation implements the rest.Deleter interface.
	DeleteOperation struct {
		Handler DeleteHandler
	}
)

var _ rest.GracefulDeleter = (*DeleteOperation)(nil)

// WithDelete returns a DeleteOperation for the given DeleteHandler.
func WithDelete(h DeleteHandler) DeleteOperation {
	return DeleteOperation{
		Handler: h,
	}
}

func (s DeleteOperation) Delete(
	ctx context.Context,
	name string,
	deleteValidation rest.ValidateObjectFunc,
	options *meta.DeleteOptions,
) (runtime.Object, bool, error) {
	if options == nil {
		options = &meta.DeleteOptions{GracePeriodSeconds: ptr.To[int64](0)}
	}

	keyFunc := keyFuncForClusterScope
	if s.Handler.NamespaceScoped() {
		keyFunc = keyFuncForNamespacedScope
	}
	key, err := keyFunc(ctx, name)
	if err != nil {
		return nil, false, err
	}

	existing, err := s.Handler.OnGet(ctx, key, ctrlcli.GetOptions{})
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, false, wrapError(ctx, name, err)
		}
		return nil, false, kerrors.NewNotFound(qualifiedResourceFromContext(ctx), name)
	}

	om, err := kmeta.Accessor(existing)
	if err != nil {
		return nil, false, kerrors.NewInternalError(err)
	}

	if deleteValidation != nil {
		err = deleteValidation(ctx, existing.DeepCopyObject())
		if err != nil {
			return nil, false, err
		}
	}

	gk := qualifiedKindFromContext(ctx)
	st := &meta.Status{
		Status: meta.StatusSuccess,
		Details: &meta.StatusDetails{
			Group: gk.Group,
			Kind:  gk.Kind,
			Name:  name,
			UID:   om.GetUID(),
		},
	}
	async := ptr.Deref(options.PropagationPolicy, meta.DeletePropagationBackground) == meta.DeletePropagationBackground
	if dryrun.IsDryRun(options.DryRun) {
		// TODO(thxCode): support delete dry run?
		return st, async, nil
	}

	err = s.Handler.OnDelete(ctx, existing, convertCtrlDeleteOptionsFromMeta(options))
	if err != nil {
		return nil, false, wrapError(ctx, name, err)
	}
	return st, async, err
}

// CollectionDeleteOperation implements the rest.Lister, rest.GracefulDeleter and rest.CollectionDeleter interfaces.
type CollectionDeleteOperation struct {
	ListWatchOperation
	DeleteOperation
}

var (
	_ rest.Lister            = (*CollectionDeleteOperation)(nil)
	_ rest.GracefulDeleter   = (*CollectionDeleteOperation)(nil)
	_ rest.CollectionDeleter = (*CollectionDeleteOperation)(nil)
)

// WithCollectionDelete returns a CollectionDeleteOperation for the given ListWatchHandler and DeleteHandler.
func WithCollectionDelete(l ListWatchOperation, d DeleteOperation) CollectionDeleteOperation {
	return CollectionDeleteOperation{
		ListWatchOperation: l,
		DeleteOperation:    d,
	}
}

func (s CollectionDeleteOperation) DeleteCollection(
	ctx context.Context,
	deleteValidation rest.ValidateObjectFunc,
	options *meta.DeleteOptions,
	listOptions *metainternal.ListOptions,
) (runtime.Object, error) {
	if options == nil {
		options = meta.NewDeleteOptions(0)
	}
	if listOptions == nil {
		listOptions = &metainternal.ListOptions{ResourceVersion: "0"}
	} else {
		listOptions = listOptions.DeepCopy()
	}

	hasLimit := listOptions.Limit > 0
	if listOptions.Limit == 0 {
		listOptions.Limit = 1000
	}

	var (
		items   []runtime.Object
		itemsCh = make(chan runtime.Object, 256)
	)
	gp := gopool.GroupWithContextIn(ctx)
	gp.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			listObj, err := s.List(ctx, listOptions)
			if err != nil {
				return err
			}
			if listObj == nil {
				return nil
			}

			newItems, err := kmeta.ExtractList(listObj)
			if err != nil {
				return kerrors.NewInternalError(err)
			}
			items = append(items, newItems...)

			for i := 0; i < len(newItems); i++ {
				select {
				case itemsCh <- newItems[i]:
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			// Done as reached limit.
			if hasLimit {
				close(itemsCh)
				return nil
			}

			// Or done as no more items.
			om, err := kmeta.ListAccessor(listObj)
			if err != nil {
				return kerrors.NewInternalError(err)
			}
			if om.GetContinue() == "" {
				close(itemsCh)
				return nil
			}

			// Otherwise, continue.
			listOptions.Continue = om.GetContinue()
			listOptions.ResourceVersion = ""
			listOptions.ResourceVersionMatch = ""
		}
	})
	gp.Go(func(ctx context.Context) error {
		for {
			var (
				item runtime.Object
				ok   bool
			)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case item, ok = <-itemsCh:
				if !ok {
					return nil
				}
			}

			om, err := kmeta.Accessor(item)
			if err != nil {
				return kerrors.NewInternalError(err)
			}

			_, _, err = s.Delete(ctx, om.GetName(), deleteValidation, options.DeepCopy())
			if err != nil && !kerrors.IsNotFound(err) {
				return err
			}
		}
	})
	err := gp.Wait()
	if err != nil {
		return nil, kerrors.NewInternalError(err)
	}

	listObj := s.NewList()
	_ = kmeta.SetList(listObj, items)
	return listObj, nil
}

type (
	// CurdOperationsHandler is an interface for CURD handler.
	CurdOperationsHandler interface {
		CreateUpdateHandler
		ListWatchHandler
		GetHandler
		DeleteHandler
	}

	// CurdOperations implements the rest.Creater, rest.Lister, rest.Watcher, rest.Getter,
	// rest.Updater, rest.GracefulDeleter and rest.CollectionDeleter interfaces.
	CurdOperations struct {
		CreateUpdateOperation
		ListWatchOperation
		GetOperation
		DeleteOperation
		CollectionDeleteOperation
	}
)

var (
	_ rest.Creater           = (*CurdOperations)(nil)
	_ rest.Lister            = (*CurdOperations)(nil)
	_ rest.Watcher           = (*CurdOperations)(nil)
	_ rest.Getter            = (*CurdOperations)(nil)
	_ rest.Updater           = (*CurdOperations)(nil)
	_ rest.GracefulDeleter   = (*CurdOperations)(nil)
	_ rest.CollectionDeleter = (*CurdOperations)(nil)
)

// WithCurd returns a CurdOperations for the given rest.TableConvertor and CurdOperationsHandler,
// which implements the rest.Creater, rest.Lister, rest.Watcher, rest.Getter, rest.Updater, rest.GracefulDeleter
// and rest.CollectionDeleter interfaces.
func WithCurd(tc rest.TableConvertor, h CurdOperationsHandler) CurdOperations {
	cu := WithCreateUpdate(h)
	l := WithListWatch(tc, h)
	g := WithGet(h)
	d := WithDelete(h)
	return CurdOperations{
		CreateUpdateOperation:     cu,
		ListWatchOperation:        l,
		GetOperation:              g,
		DeleteOperation:           d,
		CollectionDeleteOperation: WithCollectionDelete(l, d),
	}
}

type (
	// CurdProxyHandler is an handler for proxy CURD handler.
	CurdProxyHandler[DO MetaObject, DOL MetaObjectList, UO MetaObject, UOL MetaObjectList] interface {
		rest.Scoper
		rest.Storage
		// NewList returns an empty object that can be used with the List call.
		NewList() runtime.Object
		// CastObjectTo casts the object from downstream to upstream.
		CastObjectTo(DO) UO
		// CastObjectFrom casts the object from upstream to downstream.
		CastObjectFrom(UO) DO
		// CastObjectListTo casts the object list from downstream to upstream.
		CastObjectListTo(DOL) UOL
		// CastObjectListFrom casts the object list from upstream to downstream.
		CastObjectListFrom(UOL) DOL
	}

	// CurdProxyOnCreateAdvice is an interface for intercepting OnCreate.
	CurdProxyOnCreateAdvice[DO MetaObject] interface {
		// BeforeOnCreate is called before create.
		BeforeOnCreate(ctx context.Context, obj DO, opts *ctrlcli.CreateOptions) error
	}

	// CurdProxyOnListWatchAdvice is an interface for interception OnList/OnWatch.
	CurdProxyOnListWatchAdvice interface {
		// BeforeOnListWatch is called before list/watch.
		BeforeOnListWatch(ctx context.Context, options *ctrlcli.ListOptions) error
	}

	// CurdProxyOnGetAdvice is an interface for intercepting OnGet.
	CurdProxyOnGetAdvice[DO MetaObject] interface {
		// BeforeOnGet is called before get.
		BeforeOnGet(ctx context.Context, key types.NamespacedName, opts *ctrlcli.GetOptions) error
	}

	// CurdProxyOnUpdateAdvice is an interface for intercepting OnUpdate.
	CurdProxyOnUpdateAdvice[DO MetaObject] interface {
		// BeforeOnUpdate is called before update.
		BeforeOnUpdate(ctx context.Context, obj, oldObj DO, opts *ctrlcli.UpdateOptions) error
	}

	// CurdProxyOnDeleteAdvice is an interface for intercepting OnDelete.
	CurdProxyOnDeleteAdvice[DO MetaObject] interface {
		// BeforeOnDelete is called before delete.
		BeforeOnDelete(ctx context.Context, obj DO, opts *ctrlcli.DeleteOptions) error
	}

	// _CurdProxyHandler is a private struct that implements CurdProxyHandler.
	_CurdProxyHandler[DO MetaObject, DOL MetaObjectList, UO MetaObject, UOL MetaObjectList] struct {
		CurdProxyHandler[DO, DOL, UO, UOL]
		CtrlCli ctrlcli.WithWatch
	}
)

// WithCurdProxy returns a CurdOperations with the given CurdProxyHandler.
func WithCurdProxy[DO MetaObject, DOL MetaObjectList, UO MetaObject, UOL MetaObjectList](
	tc rest.TableConvertor,
	h CurdProxyHandler[DO, DOL, UO, UOL],
	ctrlCli ctrlcli.WithWatch,
) CurdOperations {
	ph := _CurdProxyHandler[DO, DOL, UO, UOL]{
		CurdProxyHandler: h,
		CtrlCli:          ctrlCli,
	}
	return WithCurd(tc, ph)
}

// _CreateHandlerWithoutNew is an interface for a creation handler without New function.
type _CreateHandlerWithoutNew interface {
	// OnCreate creates the object with the given options,
	// and returns the created object.
	OnCreate(ctx context.Context, object runtime.Object, opts ctrlcli.CreateOptions) (created runtime.Object, err error)
}

func (h _CurdProxyHandler[DO, DOL, UO, UOL]) OnCreate(
	ctx context.Context,
	obj runtime.Object,
	opts ctrlcli.CreateOptions,
) (runtime.Object, error) {
	if gh, ok := h.CurdProxyHandler.(_CreateHandlerWithoutNew); ok {
		return gh.OnCreate(ctx, obj, opts)
	}

	if bf, ok := h.CurdProxyHandler.(CurdProxyOnCreateAdvice[DO]); ok {
		err := bf.BeforeOnCreate(ctx, obj.(DO), &opts)
		if err != nil {
			return nil, err
		}
	}

	uo := h.CastObjectTo(obj.(DO))
	err := h.CtrlCli.Create(ctx, uo, &opts)
	return h.CastObjectFrom(uo), err
}

func (h _CurdProxyHandler[DO, DOL, UO, UOL]) OnList(
	ctx context.Context,
	opts ctrlcli.ListOptions,
) (runtime.Object, error) {
	if gh, ok := h.CurdProxyHandler.(ListHandler); ok {
		return gh.OnList(ctx, opts)
	}

	if bf, ok := h.CurdProxyHandler.(CurdProxyOnListWatchAdvice); ok {
		err := bf.BeforeOnListWatch(ctx, &opts)
		if err != nil {
			return nil, err
		}
	}

	uol := h.CastObjectListTo(h.NewList().(DOL))
	err := h.CtrlCli.List(ctx, uol, &opts)
	return h.CastObjectListFrom(uol), err
}

func (h _CurdProxyHandler[DO, DOL, UO, UOL]) OnWatch(
	ctx context.Context,
	opts ctrlcli.ListOptions,
) (watch.Interface, error) {
	if gh, ok := h.CurdProxyHandler.(WatchHandler); ok {
		return gh.OnWatch(ctx, opts)
	}

	if bf, ok := h.CurdProxyHandler.(CurdProxyOnListWatchAdvice); ok {
		err := bf.BeforeOnListWatch(ctx, &opts)
		if err != nil {
			return nil, err
		}
	}

	uol := h.CastObjectListTo(h.NewList().(DOL))
	uw, err := h.CtrlCli.Watch(ctx, uol, &opts)
	if err != nil {
		return nil, err
	}

	c := make(chan watch.Event)
	dw := watch.NewProxyWatcher(c)
	gopool.Go(func() {
		defer close(c)
		defer uw.Stop()

		for {
			select {
			case <-ctx.Done():
				// Cancel by context.
				return
			case <-dw.StopChan():
				// Cancel by downstream.
				return
			case e, ok := <-uw.ResultChan():
				if !ok {
					// Close by upstream.
					return
				}

				if e.Object == nil {
					// Nothing to do.
					c <- e
					continue
				}

				// Type assert.
				uo, ok := e.Object.(UO)
				if !ok {
					continue
				}

				// Cast.
				e.Object = h.CastObjectFrom(uo)
				c <- e
			}
		}
	})

	return dw, nil
}

func (h _CurdProxyHandler[DO, DOL, UO, UOL]) OnGet(
	ctx context.Context,
	key types.NamespacedName,
	opts ctrlcli.GetOptions,
) (runtime.Object, error) {
	if gh, ok := h.CurdProxyHandler.(GetHandler); ok {
		return gh.OnGet(ctx, key, opts)
	}

	if bf, ok := h.CurdProxyHandler.(CurdProxyOnGetAdvice[DO]); ok {
		err := bf.BeforeOnGet(ctx, key, &opts)
		if err != nil {
			return nil, err
		}
	}
	uo := h.CastObjectTo(h.New().(DO))
	err := h.CtrlCli.Get(ctx, key, uo, &opts)
	return h.CastObjectFrom(uo), err
}

// _UpdateHandlerWithoutNew is an interface for an updating handler without New function.
type _UpdateHandlerWithoutNew interface {
	// OnUpdate updates the object with the given options,
	// and returns the updated object.
	OnUpdate(ctx context.Context, object, oldObject runtime.Object, opts ctrlcli.UpdateOptions) (updated runtime.Object, err error)
}

func (h _CurdProxyHandler[DO, DOL, UO, UOL]) OnUpdate(
	ctx context.Context,
	obj, oldObj runtime.Object,
	opts ctrlcli.UpdateOptions,
) (runtime.Object, error) {
	if gh, ok := h.CurdProxyHandler.(_UpdateHandlerWithoutNew); ok {
		return gh.OnUpdate(ctx, obj, oldObj, opts)
	}

	if bf, ok := h.CurdProxyHandler.(CurdProxyOnUpdateAdvice[DO]); ok {
		err := bf.BeforeOnUpdate(ctx, obj.(DO), oldObj.(DO), &opts)
		if err != nil {
			return nil, err
		}
	}
	uo := h.CastObjectTo(obj.(DO))
	err := h.CtrlCli.Update(ctx, uo, &opts)
	return h.CastObjectFrom(uo), err
}

func (h _CurdProxyHandler[DO, DOL, UO, UOL]) OnDelete(
	ctx context.Context,
	obj runtime.Object,
	opts ctrlcli.DeleteOptions,
) error {
	if gh, ok := h.CurdProxyHandler.(DeleteHandler); ok {
		return gh.OnDelete(ctx, obj, opts)
	}

	if bf, ok := h.CurdProxyHandler.(CurdProxyOnDeleteAdvice[DO]); ok {
		err := bf.BeforeOnDelete(ctx, obj.(DO), &opts)
		if err != nil {
			return err
		}
	}
	uo := h.CastObjectTo(obj.(DO))
	return h.CtrlCli.Delete(ctx, uo, &opts)
}

func qualifiedResourceFromContext(ctx context.Context) schema.GroupResource {
	ri, _ := genericapirequest.RequestInfoFrom(ctx)
	return schema.GroupResource{
		Group:    ri.APIGroup,
		Resource: ri.Resource,
	}
}

func qualifiedKindFromContext(ctx context.Context) schema.GroupKind {
	ri, _ := genericapirequest.RequestInfoFrom(ctx)
	return schema.GroupKind{
		Group: ri.APIGroup,
		Kind:  ri.Resource, // We use the resource name as the kind name.
	}
}

func wrapError(ctx context.Context, name string, err error) error {
	if err == nil {
		return nil
	}

	gk := qualifiedKindFromContext(ctx)

	var es meta.Status
	if st, ok := err.(kerrors.APIStatus); ok || errors.As(err, &st) {
		ss := st.Status()
		ssd := ss.Details
		if ssd == nil {
			ssd = &meta.StatusDetails{}
		}
		es = meta.Status{
			Status: ss.Status,
			Code:   ss.Code,
			Reason: ss.Reason,
			Details: &meta.StatusDetails{
				Name:              name,
				Group:             gk.Group,
				Kind:              gk.Kind,
				UID:               ssd.UID,
				Causes:            ssd.Causes,
				RetryAfterSeconds: ssd.RetryAfterSeconds,
			},
			Message: ss.Message,
		}
	} else {
		es = meta.Status{
			Status: meta.StatusFailure,
			Code:   http.StatusInternalServerError,
			Reason: meta.StatusReasonInternalError,
			Details: &meta.StatusDetails{
				Name:   name,
				Group:  gk.Group,
				Kind:   gk.Kind,
				Causes: []meta.StatusCause{{Message: err.Error()}},
			},
			Message: fmt.Sprintf("Internal error occurred: %v", err),
		}
	}

	return &kerrors.StatusError{ErrStatus: es}
}
