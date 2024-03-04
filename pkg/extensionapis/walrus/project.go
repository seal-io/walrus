package walrus

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"

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

// ProjectHandler handles v1.Project objects.
//
// ProjectHandler maps the v1.Project object to a Kubernetes Namespace resource,
// which is named as the project's name.
type ProjectHandler struct {
	extensionapi.ObjectInfo
	extensionapi.CurdOperations

	Client ctrlcli.Client
}

func (h *ProjectHandler) SetupHandler(
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
	gvr = walrus.SchemeGroupVersionResource("projects")

	// Create table convertor to pretty the kubectl's output.
	var tc rest.TableConvertor
	{
		tc, err = extensionapi.NewJSONPathTableConvertor(
			extensionapi.JSONPathTableColumnDefinition{
				TableColumnDefinition: meta.TableColumnDefinition{
					Name: "Phase",
					Type: "string",
				},
				JSONPath: ".status.phase",
			})
		if err != nil {
			return
		}
	}

	// As storage.
	h.ObjectInfo = &walrus.Project{}
	h.CurdOperations = extensionapi.WithCurd(tc, h)

	// Set client.
	h.Client = opts.Manager.GetClient()

	// Create subresource handlers.
	srs = map[string]rest.Storage{}
	{
		// Handle /subjects.
		srs["subjects"] = newProjectSubjectsHandler(opts)
	}

	return
}

var (
	_ rest.Storage           = (*ProjectHandler)(nil)
	_ rest.Creater           = (*ProjectHandler)(nil)
	_ rest.Lister            = (*ProjectHandler)(nil)
	_ rest.Watcher           = (*ProjectHandler)(nil)
	_ rest.Getter            = (*ProjectHandler)(nil)
	_ rest.Updater           = (*ProjectHandler)(nil)
	_ rest.Patcher           = (*ProjectHandler)(nil)
	_ rest.GracefulDeleter   = (*ProjectHandler)(nil)
	_ rest.CollectionDeleter = (*ProjectHandler)(nil)
)

func (h *ProjectHandler) New() runtime.Object {
	return &walrus.Project{}
}

func (h *ProjectHandler) Destroy() {
}

func (h *ProjectHandler) OnCreate(ctx context.Context, obj runtime.Object, opts ctrlcli.CreateOptions) (runtime.Object, error) {
	// Validate.
	proj := obj.(*walrus.Project)
	{
		var errs field.ErrorList
		if proj.Namespace != systemkuberes.SystemNamespaceName {
			errs = append(errs, field.Invalid(field.NewPath("metadata.namespace"),
				proj.Namespace, "project namespace must be "+systemkuberes.SystemNamespaceName))
		}
		if slices.Contains([]string{"kube-system", "kube-public"}, proj.Name) {
			errs = append(errs, field.Invalid(field.NewPath("metadata.name"),
				proj.Name, "project name is reserved"))
		}
		if stringx.StringWidth(proj.Name) > 30 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("metadata.name"), stringx.StringWidth(proj.Name), 30))
		}
		if stringx.StringWidth(proj.Spec.DisplayName) > 30 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("spec.displayName"), stringx.StringWidth(proj.Spec.DisplayName), 30))
		}
		if stringx.StringWidth(proj.Spec.Description) > 50 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("spec.description"), stringx.StringWidth(proj.Spec.Description), 50))
		}
		if len(errs) > 0 {
			return nil, kerrors.NewInvalid(walrus.SchemeKind("projects"), proj.Name, errs)
		}
	}

	// Create.
	if proj.Name == systemkuberes.DefaultProjectName {
		// NB(thxCode): The default project is created by the system,
		// so we need another approach to adopt the default project.
		ns := convertNamespaceFromProject(proj)
		{
			// Refill UID and ResourceVersion.
			aNs := new(core.Namespace)
			err := h.Client.Get(ctx, ctrlcli.ObjectKeyFromObject(ns), aNs)
			if err != nil {
				return nil, kerrors.NewInternalError(fmt.Errorf("get default namespace: %w", err))
			}
			ns.UID = aNs.UID
			ns.ResourceVersion = aNs.ResourceVersion
		}
		err := h.Client.Update(ctx, ns)
		if err != nil {
			return nil, kerrors.NewInternalError(fmt.Errorf("create default project: %w", err))
		}
		proj = convertProjectFromNamespace(ns)
	} else {
		ns := convertNamespaceFromProject(proj)
		err := h.Client.Create(ctx, ns, &opts)
		if err != nil {
			return nil, err
		}
		proj = convertProjectFromNamespace(ns)
	}

	// Create RBAC.
	err := systemauthz.CreateProjectSpace(ctx, h.Client, proj)
	if err != nil {
		return nil, kerrors.NewInternalError(fmt.Errorf("create project space: %w", err))
	}

	return proj, nil
}

func (h *ProjectHandler) NewList() runtime.Object {
	return &walrus.ProjectList{}
}

func (h *ProjectHandler) OnList(ctx context.Context, opts ctrlcli.ListOptions) (runtime.Object, error) {
	// List.
	nsList := new(core.NamespaceList)
	err := h.Client.List(ctx, nsList,
		convertNamespaceListOptsFromProjectListOpts(opts))
	if err != nil {
		return nil, err
	}

	// TODO Validate RBAC

	// Convert.
	pList := convertProjectListFromNamespaceList(nsList, opts)
	return pList, nil
}

func (h *ProjectHandler) OnWatch(ctx context.Context, opts ctrlcli.ListOptions) (watch.Interface, error) {
	// Watch.
	uw, err := h.Client.(ctrlcli.WithWatch).Watch(ctx, new(core.NamespaceList),
		convertNamespaceListOptsFromProjectListOpts(opts))
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

				// TODO RBAC

				// Type assert.
				ns, ok := e.Object.(*core.Namespace)
				if !ok {
					c <- e
					continue
				}

				// Process bookmark.
				if e.Type == watch.Bookmark {
					e.Object = &walrus.Project{ObjectMeta: ns.ObjectMeta}
					c <- e
					continue
				}

				// Convert.
				proj := safeConvertProjectFromNamespace(ns, opts.Namespace)
				if proj == nil {
					continue
				}

				// Dispatch.
				e.Object = proj
				c <- e
			}
		}
	})

	return dw, nil
}

func (h *ProjectHandler) OnGet(ctx context.Context, key types.NamespacedName, opts ctrlcli.GetOptions) (runtime.Object, error) {
	// Validate.
	if key.Namespace != systemkuberes.SystemNamespaceName {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("projects"), key.Name)
	}

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

	// TODO Validate RBAC

	// Convert.
	proj := convertProjectFromNamespace(ns)
	if proj == nil {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("projects"), key.Name)
	}
	return proj, nil
}

func (h *ProjectHandler) OnUpdate(ctx context.Context, obj, _ runtime.Object, opts ctrlcli.UpdateOptions) (runtime.Object, error) {
	// Validate.
	proj := obj.(*walrus.Project)
	{
		var errs field.ErrorList
		if stringx.StringWidth(proj.Spec.DisplayName) > 30 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("spec.displayName"), stringx.StringWidth(proj.Spec.DisplayName), 30))
		}
		if stringx.StringWidth(proj.Spec.Description) > 50 {
			errs = append(errs, field.TooLongMaxLength(
				field.NewPath("spec.description"), stringx.StringWidth(proj.Spec.Description), 50))
		}
		if len(errs) > 0 {
			return nil, kerrors.NewInvalid(walrus.SchemeKind("projects"), proj.Name, errs)
		}
	}

	// TODO Validate RBAC

	// Update.
	{
		ns := convertNamespaceFromProject(proj)
		err := h.Client.Update(ctx, ns, &opts)
		if err != nil {
			return nil, err
		}
		proj = convertProjectFromNamespace(ns)
	}

	return proj, nil
}

func (h *ProjectHandler) OnDelete(ctx context.Context, obj runtime.Object, opts ctrlcli.DeleteOptions) error {
	proj := obj.(*walrus.Project)

	// Validate.
	{
		// Prevent deleting default project.
		if proj.Name == core.NamespaceDefault {
			return kerrors.NewBadRequest("cannot delete default project")
		}
		// Prevent deleting if it has environments.
		envList := new(walrus.EnvironmentList)
		err := h.Client.List(ctx, envList, &ctrlcli.ListOptions{
			Namespace: proj.Name,
		})
		if err != nil {
			return kerrors.NewInternalError(fmt.Errorf("list environments below the project: %w", err))
		}
		if len(envList.Items) != 0 {
			return kerrors.NewConflict(walrus.SchemeResource("projects"), proj.Name,
				errors.New("project has environments"))
		}
	}

	// Unlock if needed.
	ns := convertNamespaceFromProject(proj)
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
	return err
}

func convertNamespaceListOptsFromProjectListOpts(in ctrlcli.ListOptions) (out *ctrlcli.ListOptions) {
	if in.Namespace != systemkuberes.SystemNamespaceName {
		return &in
	}

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
	if lbs := systemmeta.GetResourcesLabelSelectorOfType("projects"); in.LabelSelector == nil {
		in.LabelSelector = lbs
	} else {
		reqs, _ := lbs.Requirements()
		in.LabelSelector = in.LabelSelector.DeepCopySelector().Add(reqs...)
	}

	return &in
}

func convertNamespaceFromProject(proj *walrus.Project) *core.Namespace {
	ns := &core.Namespace{
		ObjectMeta: proj.ObjectMeta,
	}
	systemmeta.NoteResource(ns, "projects", map[string]string{
		"displayName": proj.Spec.DisplayName,
		"description": proj.Spec.Description,
	})
	ns.Namespace = ""
	if ns.DeletionTimestamp == nil {
		systemmeta.Lock(ns)
	}
	return ns
}

func convertProjectFromNamespace(ns *core.Namespace) *walrus.Project {
	if ns == nil {
		return nil
	}

	resType, notes := systemmeta.UnnoteResource(ns)
	if resType != "projects" {
		return nil
	}

	proj := &walrus.Project{
		ObjectMeta: ns.ObjectMeta,
		Spec: walrus.ProjectSpec{
			DisplayName: notes["displayName"],
			Description: notes["description"],
		},
		Status: walrus.ProjectStatus{
			Phase: ns.Status.Phase,
		},
	}
	proj.Namespace = systemkuberes.SystemNamespaceName
	return proj
}

func safeConvertProjectFromNamespace(ns *core.Namespace, reqNamespace string) *walrus.Project {
	proj := convertProjectFromNamespace(ns)
	if proj != nil && reqNamespace != "" && reqNamespace != proj.Namespace {
		// NB(thxCode): sanitize if the project's namespace doesn't match requested namespace.
		proj = nil
	}
	return proj
}

func convertProjectListFromNamespaceList(nsList *core.NamespaceList, opts ctrlcli.ListOptions) *walrus.ProjectList {
	if nsList == nil {
		return &walrus.ProjectList{}
	}

	// Sort by resource version.
	sort.SliceStable(nsList.Items, func(i, j int) bool {
		l, r := nsList.Items[i].ResourceVersion, nsList.Items[j].ResourceVersion
		return len(l) < len(r) ||
			(len(l) == len(r) && l < r)
	})

	pList := &walrus.ProjectList{
		Items: make([]walrus.Project, 0, len(nsList.Items)),
	}

	for i := range nsList.Items {
		proj := safeConvertProjectFromNamespace(&nsList.Items[i], opts.Namespace)
		if proj == nil {
			continue
		}
		pList.Items = append(pList.Items, *proj)
	}

	return pList
}
