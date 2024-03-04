package walrus

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/seal-io/utils/pools/gopool"
	"github.com/seal-io/utils/stringx"
	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/rest"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/extensionapi"
	"github.com/seal-io/walrus/pkg/kubemeta"
	"github.com/seal-io/walrus/pkg/systemauthz"
	"github.com/seal-io/walrus/pkg/systemkuberes"
	"github.com/seal-io/walrus/pkg/systemmeta"
)

// EnvironmentHandler handles v1.Environment objects.
//
// EnvironmentHandler maps the v1.Environment object to a Kubernetes Namespace resource,
// which is name as the environment's name.
//
// Each v1.Environment object will be controlled by a v1.Project object,
// which records in the OwnerReferences of the Namespace resource.
type EnvironmentHandler struct {
	extensionapi.ObjectInfo
	extensionapi.CurdOperations

	Client ctrlcli.Client
}

func (h *EnvironmentHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Configure field indexer.
	fi := opts.Manager.GetFieldIndexer()
	err = fi.IndexField(ctx, &core.Namespace{}, "metadata.name",
		func(obj ctrlcli.Object) []string {
			if obj == nil {
				return nil
			}
			return []string{obj.GetName()}
		})
	if err != nil {
		return
	}

	// Declare GVR.
	gvr = walrus.SchemeGroupVersionResource("environments")

	// Create table convertor to pretty the kubectl's output.
	var tc rest.TableConvertor
	{
		tc, err = extensionapi.NewJSONPathTableConvertor(
			extensionapi.JSONPathTableColumnDefinition{
				TableColumnDefinition: meta.TableColumnDefinition{
					Name: "Type",
					Type: "string",
				},
				JSONPath: ".spec.type",
			},
			extensionapi.JSONPathTableColumnDefinition{
				TableColumnDefinition: meta.TableColumnDefinition{
					Name: "Phase",
					Type: "string",
				},
				JSONPath: ".status.phase",
			},
			extensionapi.JSONPathTableColumnDefinition{
				TableColumnDefinition: meta.TableColumnDefinition{
					Name: "Project",
					Type: "string",
				},
				JSONPath: ".status.project",
			})
		if err != nil {
			return
		}
	}

	// As storage.
	h.ObjectInfo = &walrus.Environment{}
	h.CurdOperations = extensionapi.WithCurd(tc, h)

	// Set client.
	h.Client = opts.Manager.GetClient()

	return
}

var (
	_ rest.Storage           = (*EnvironmentHandler)(nil)
	_ rest.Creater           = (*EnvironmentHandler)(nil)
	_ rest.Lister            = (*EnvironmentHandler)(nil)
	_ rest.Watcher           = (*EnvironmentHandler)(nil)
	_ rest.Getter            = (*EnvironmentHandler)(nil)
	_ rest.Updater           = (*EnvironmentHandler)(nil)
	_ rest.Patcher           = (*EnvironmentHandler)(nil)
	_ rest.GracefulDeleter   = (*EnvironmentHandler)(nil)
	_ rest.CollectionDeleter = (*EnvironmentHandler)(nil)
)

func (h *EnvironmentHandler) New() runtime.Object {
	return &walrus.Environment{}
}

func (h *EnvironmentHandler) Destroy() {}

func (h *EnvironmentHandler) OnCreate(ctx context.Context, obj runtime.Object, opts ctrlcli.CreateOptions) (runtime.Object, error) {
	// Validate.
	env := obj.(*walrus.Environment)
	{
		var errs field.ErrorList
		if !strings.HasPrefix(env.Name, env.Namespace+"-") {
			errs = append(errs, field.Invalid(
				field.NewPath("name"), env.Name, "name must start with the namespace"))
		}
		if stringx.StringWidth(env.Name) > 30 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("name"), stringx.StringWidth(env.Name), 30))
		}
		if err := env.Spec.Type.Validate(); err != nil {
			errs = append(errs, field.Invalid(
				field.NewPath("spec.type"), env.Spec.Type, err.Error()))
		}
		if stringx.StringWidth(env.Spec.DisplayName) > 30 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("spec.displayName"), stringx.StringWidth(env.Spec.DisplayName), 30))
		}
		if stringx.StringWidth(env.Spec.Description) > 50 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("spec.description"), stringx.StringWidth(env.Spec.Description), 50))
		}
		if len(errs) > 0 {
			return nil, kerrors.NewInvalid(walrus.SchemeKind("environments"), env.Name, errs)
		}
	}

	// Get project.
	proj := &walrus.Project{
		ObjectMeta: meta.ObjectMeta{
			Namespace: systemkuberes.SystemNamespaceName,
			Name:      env.Namespace,
		},
	}
	err := h.Client.Get(ctx, ctrlcli.ObjectKeyFromObject(proj), proj)
	if err != nil {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("projects"), proj.Name)
	}

	// Create.
	{
		ns := convertNamespaceFromEnvironment(env)
		kubemeta.ControlOn(ns, proj, walrus.SchemeGroupVersion.WithKind("Project"))
		err = h.Client.Create(ctx, ns, &opts)
		if err != nil {
			return nil, err
		}
		env = convertEnvironmentFromNamespace(ns)
	}

	// Create RBAC.
	err = systemauthz.CreateEnvironmentSpace(ctx, h.Client, env)
	if err != nil {
		return nil, kerrors.NewInternalError(fmt.Errorf("create environment space: %w", err))
	}

	return env, nil
}

func (h *EnvironmentHandler) NewList() runtime.Object {
	return &walrus.EnvironmentList{}
}

func (h *EnvironmentHandler) OnList(ctx context.Context, opts ctrlcli.ListOptions) (runtime.Object, error) {
	// List.
	nsList := new(core.NamespaceList)
	err := h.Client.List(ctx, nsList,
		convertNamespaceListOptsFromEnvironmentListOpts(opts))
	if err != nil {
		return nil, err
	}

	// Convert.
	eList := convertEnvironmentListFromNamespaceList(nsList, opts)
	return eList, nil
}

func (h *EnvironmentHandler) OnWatch(ctx context.Context, opts ctrlcli.ListOptions) (watch.Interface, error) {
	// Watch.
	uw, err := h.Client.(ctrlcli.WithWatch).Watch(ctx, new(core.NamespaceList),
		convertNamespaceListOptsFromEnvironmentListOpts(opts))
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
				// Stop by downstream.
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

				// Type assert.
				ns, ok := e.Object.(*core.Namespace)
				if !ok {
					c <- e
					continue
				}

				// Process bookmark.
				if e.Type == watch.Bookmark {
					e.Object = &walrus.Environment{ObjectMeta: ns.ObjectMeta}
					c <- e
					continue
				}

				// Convert.
				env := safeConvertEnvironmentFromNamespace(ns, opts.Namespace)
				if env == nil {
					continue
				}

				// Ignore if not be selected by `kubectl get --field-selector=metadata.namespace=...`.
				if fs := opts.FieldSelector; fs != nil &&
					!fs.Matches(fields.Set{"metadata.namespace": env.Namespace}) {
					continue
				}

				// Dispatch.
				e.Object = env
				c <- e
			}
		}
	})

	return dw, nil
}

func (h *EnvironmentHandler) OnGet(ctx context.Context, key types.NamespacedName, opts ctrlcli.GetOptions) (runtime.Object, error) {
	// Get.
	ns := &core.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: key.Name,
		},
	}
	err := h.Client.Get(ctx, ctrlcli.ObjectKeyFromObject(ns), ns, &opts)
	if err != nil {
		return nil, err
	}

	// Convert.
	env := safeConvertEnvironmentFromNamespace(ns, key.Namespace)
	if env == nil {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("environments"), key.Name)
	}

	return env, nil
}

func (h *EnvironmentHandler) OnUpdate(ctx context.Context, obj, oldObj runtime.Object, opts ctrlcli.UpdateOptions) (runtime.Object, error) {
	// Validate.
	env := obj.(*walrus.Environment)
	{
		oldEnv := oldObj.(*walrus.Environment)
		var errs field.ErrorList
		if env.Spec.Type != oldEnv.Spec.Type {
			errs = append(errs, field.Invalid(
				field.NewPath("spec.type"), env.Spec.Type, "type is immutable"))
		}
		if stringx.StringWidth(env.Spec.DisplayName) > 30 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("spec.displayName"), stringx.StringWidth(env.Spec.DisplayName), 30))
		}
		if stringx.StringWidth(env.Spec.Description) > 50 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("spec.description"), stringx.StringWidth(env.Spec.Description), 50))
		}
		if len(errs) > 0 {
			return nil, kerrors.NewInvalid(walrus.SchemeKind("environments"), env.Name, errs)
		}
	}

	// Update.
	{
		ns := convertNamespaceFromEnvironment(env)
		err := h.Client.Update(ctx, ns, &opts)
		if err != nil {
			return nil, err
		}
		env = convertEnvironmentFromNamespace(ns)
	}

	return env, nil
}

func (h *EnvironmentHandler) OnDelete(ctx context.Context, obj runtime.Object, opts ctrlcli.DeleteOptions) error {
	env := obj.(*walrus.Environment)

	// Validate.
	{
		// Prevent deleting if it has resources.
		resList := new(walrus.ResourceList)
		err := h.Client.List(ctx, resList, &ctrlcli.ListOptions{
			Namespace: env.Name,
		})
		if err != nil {
			return kerrors.NewInternalError(fmt.Errorf("list resources below the environment: %w", err))
		}
		if len(resList.Items) != 0 {
			return kerrors.NewForbidden(walrus.SchemeResource("environments"), env.Name,
				errors.New("environment has resources"))
		}
	}

	// Unlock if needed.
	ns := convertNamespaceFromEnvironment(env)
	unlocked := systemmeta.Unlock(ns)
	if !unlocked {
		err := h.Client.Update(ctx, ns)
		if err != nil {
			return fmt.Errorf("unset finalizer: %w", err)
		}
	}

	// Delete.
	err := h.Client.Delete(ctx, ns, &opts)
	if err != nil && kerrors.IsNotFound(err) && !unlocked {
		// NB(thxCode): If deleting resource has been locked,
		// we ignore the not found error after we unlock it.
		return nil
	}

	// Delete RBAC.
	err = systemauthz.DeleteEnvironmentSpace(ctx, h.Client, env)
	if err != nil {
		return kerrors.NewInternalError(fmt.Errorf("delete environment space: %w", err))
	}

	return nil
}

func convertNamespaceListOptsFromEnvironmentListOpts(in ctrlcli.ListOptions) (out *ctrlcli.ListOptions) {
	// Ignore namespace selector.
	in.Namespace = ""
	if in.FieldSelector != nil {
		reqs := slices.DeleteFunc(in.FieldSelector.Requirements(), func(req fields.Requirement) bool {
			return req.Field == "metadata.namespace" &&
				((req.Operator == selection.Equals && req.Value == systemkuberes.SystemNamespaceName) ||
					(req.Operator == selection.NotEquals && req.Value != systemkuberes.SystemNamespaceName))
		})
		if len(reqs) == 0 {
			in.FieldSelector = nil
		} else {
			in.FieldSelector = kubemeta.FieldSelectorFromRequirements(reqs)
		}
	}

	// Add necessary label selector.
	if lbs := systemmeta.GetResourcesLabelSelectorOfType("environments"); in.LabelSelector == nil {
		in.LabelSelector = lbs
	} else {
		reqs, _ := lbs.Requirements()
		in.LabelSelector = in.LabelSelector.DeepCopySelector().Add(reqs...)
	}

	return &in
}

func convertNamespaceFromEnvironment(env *walrus.Environment) *core.Namespace {
	ns := &core.Namespace{
		ObjectMeta: env.ObjectMeta,
	}
	systemmeta.NoteResource(ns, "environments", map[string]string{
		"type":        env.Spec.Type.String(),
		"displayName": env.Spec.DisplayName,
		"description": env.Spec.Description,
	})
	ns.Namespace = ""
	if ns.DeletionTimestamp == nil {
		systemmeta.Lock(ns)
	}
	return ns
}

func convertEnvironmentFromNamespace(ns *core.Namespace) *walrus.Environment {
	if ns == nil {
		return nil
	}

	resType, notes := systemmeta.UnnoteResource(ns)
	if resType != "environments" {
		return nil
	}

	ref := meta.GetControllerOf(ns)
	if ref == nil ||
		ref.APIVersion != walrus.SchemeGroupVersion.String() ||
		ref.Kind != "Project" {
		return nil
	}

	env := &walrus.Environment{
		ObjectMeta: ns.ObjectMeta,
		Spec: walrus.EnvironmentSpec{
			Type:        walrus.EnvironmentType(notes["type"]),
			DisplayName: notes["displayName"],
			Description: notes["description"],
		},
		Status: walrus.EnvironmentStatus{
			Project: ref.Name,
			Phase:   ns.Status.Phase,
		},
	}
	env.Namespace = ref.Name
	return env
}

func safeConvertEnvironmentFromNamespace(ns *core.Namespace, reqNamespace string) *walrus.Environment {
	env := convertEnvironmentFromNamespace(ns)
	if env != nil && reqNamespace != "" && reqNamespace != env.Namespace {
		// NB(thxCode): sanitize if the environment's namespace doesn't match requested namespace.
		env = nil
	}
	return env
}

func convertEnvironmentListFromNamespaceList(nsList *core.NamespaceList, opts ctrlcli.ListOptions) *walrus.EnvironmentList {
	if nsList == nil {
		return &walrus.EnvironmentList{}
	}

	// Sort by resource version.
	sort.SliceStable(nsList.Items, func(i, j int) bool {
		l, r := nsList.Items[i].ResourceVersion, nsList.Items[j].ResourceVersion
		return len(l) < len(r) ||
			(len(l) == len(r) && l < r)
	})

	eList := &walrus.EnvironmentList{
		Items: make([]walrus.Environment, 0, len(nsList.Items)),
	}

	for i := range nsList.Items {
		env := safeConvertEnvironmentFromNamespace(&nsList.Items[i], opts.Namespace)
		if env == nil {
			continue
		}
		// Ignore if not be selected by `kubectl get --field-selector=metadata.namespace=...`.
		if fs := opts.FieldSelector; fs != nil &&
			!fs.Matches(fields.Set{"metadata.namespace": env.Namespace}) {
			continue
		}
		eList.Items = append(eList.Items, *env)
	}

	return eList
}
