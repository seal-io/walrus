package walrus

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seal-io/utils/pools/gopool"
	"github.com/seal-io/utils/stringx"
	"golang.org/x/exp/maps"
	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/rest"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/extensionapi"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubemeta"
	"github.com/seal-io/walrus/pkg/systemkuberes"
	"github.com/seal-io/walrus/pkg/systemmeta"
)

// VariableHandler handles v1.Variable objects.
//
// VariableHandler maps all v1.Variable objects to a Kubernetes Secret resource,
// which is named as "${namespace}/walrus-variables".
//
// Each v1.Variable object records as a key-value pair in the Secret's Data field.
type VariableHandler struct {
	extensionapi.ObjectInfo
	extensionapi.CurdOperations

	Client ctrlcli.Client
}

func (h *VariableHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Configure field indexer.
	fi := opts.Manager.GetFieldIndexer()
	err = fi.IndexField(ctx, &core.Secret{}, "metadata.name", func(obj ctrlcli.Object) []string {
		if obj == nil {
			return nil
		}
		return []string{obj.GetName()}
	})
	if err != nil {
		return
	}

	// Declare GVR.
	gvr = walrus.SchemeGroupVersionResource("variables")

	// Create table convertor to pretty the kubectl's output.
	var tc rest.TableConvertor
	{
		tc, err = extensionapi.NewJSONPathTableConvertor(
			extensionapi.JSONPathTableColumnDefinition{
				TableColumnDefinition: meta.TableColumnDefinition{
					Name: "Value",
					Type: "string",
				},
				JSONPath: ".status.value",
			},
			extensionapi.JSONPathTableColumnDefinition{
				TableColumnDefinition: meta.TableColumnDefinition{
					Name: "Environment",
					Type: "string",
				},
				JSONPath: ".status.environment",
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
	h.ObjectInfo = &walrus.Variable{}
	h.CurdOperations = extensionapi.WithCurd(tc, h)

	// Set client.
	h.Client = opts.Manager.GetClient()

	return
}

var (
	_ rest.Storage           = (*VariableHandler)(nil)
	_ rest.Creater           = (*VariableHandler)(nil)
	_ rest.Lister            = (*VariableHandler)(nil)
	_ rest.Watcher           = (*VariableHandler)(nil)
	_ rest.Getter            = (*VariableHandler)(nil)
	_ rest.Updater           = (*VariableHandler)(nil)
	_ rest.GracefulDeleter   = (*VariableHandler)(nil)
	_ rest.CollectionDeleter = (*VariableHandler)(nil)
)

func (h *VariableHandler) New() runtime.Object {
	return &walrus.Variable{}
}

func (h *VariableHandler) Destroy() {}

func (h *VariableHandler) OnCreate(ctx context.Context, obj runtime.Object, opts ctrlcli.CreateOptions) (runtime.Object, error) {
	// Validate.
	vra := obj.(*walrus.Variable)
	if vra.Spec.Value == nil {
		return nil, kerrors.NewInvalid(walrus.SchemeKind("variables"),
			vra.Name, field.ErrorList{field.Required(field.NewPath("spec.value"), "variable value is required")})
	}

	// Check affiliation.
	var project, environment string
	if vra.Namespace != systemkuberes.SystemNamespaceName {
		owner := &core.Namespace{
			ObjectMeta: meta.ObjectMeta{
				Name: vra.Namespace,
			},
		}
		_ = h.Client.Get(ctx, ctrlcli.ObjectKeyFromObject(owner), owner)
		resType := systemmeta.DescribeResourceType(owner)
		switch resType {
		default:
			return nil, kerrors.NewInvalid(walrus.SchemeKind("variables"),
				vra.Name, field.ErrorList{field.Invalid(field.NewPath("metadata.namespace"), vra.Namespace, "namespace is not a project or environment")}) // nolint: lll
		case "projects":
			project = owner.Name
		case "environments":
			environment = owner.Name
			proj := kubemeta.GetOwnerRefOfNoCopy(owner, walrus.SchemeGroupVersionKind("Project"))
			if proj == nil {
				return nil, kerrors.NewInvalid(walrus.SchemeKind("variables"),
					vra.Name, field.ErrorList{field.Invalid(field.NewPath("metadata.namespace"), vra.Namespace, "environment is not belong to any project")}) // nolint: lll
			}
			project = proj.Name
		}
	}

	// Update or Create.
	eSec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: vra.Namespace,
			Name:      systemkuberes.VariablesDelegatedSecretName,
		},
		Data: map[string][]byte{
			vra.Name: []byte(*vra.Spec.Value),
		},
	}
	eResType := "variables"
	eNotes := map[string]string{
		"project":               project,
		"environment":           environment,
		vra.Name + "-uid":       uuid.NewString(),
		vra.Name + "-create-at": time.Now().Format(time.RFC3339),
		vra.Name + "-sensitive": strconv.FormatBool(vra.Spec.Sensitive),
	}
	systemmeta.NoteResource(eSec, eResType, eNotes)
	alignFn := func(aSec *core.Secret) (*core.Secret, bool, error) {
		// Validate.
		if aSec.Data == nil {
			aSec.Data = make(map[string][]byte)
		}
		if _, ok := aSec.Data[vra.Name]; ok {
			return nil, true, kerrors.NewAlreadyExists(walrus.SchemeResource("variables"), vra.Name)
		}
		// Align data.
		aSec.Data[vra.Name] = eSec.Data[vra.Name]
		// Align delegated info.
		_, aNotes := systemmeta.DescribeResource(aSec)
		eNotes := maps.Clone(eNotes)
		maps.Copy(eNotes, aNotes)
		systemmeta.NoteResource(aSec, eResType, eNotes)
		return aSec, false, nil
	}

	sec, err := kubeclientset.UpdateWithCtrlClient(ctx, h.Client, eSec,
		kubeclientset.UpdateAfterAlign(alignFn),
		kubeclientset.UpdateOrCreate[*core.Secret]())
	if err != nil {
		return nil, err
	}

	// Convert.
	vra = convertVariableFromSecret(sec, vra.Name)
	return vra, nil
}

func (h *VariableHandler) NewList() runtime.Object {
	return &walrus.VariableList{}
}

func (h *VariableHandler) OnList(ctx context.Context, opts ctrlcli.ListOptions) (runtime.Object, error) {
	// List.
	if opts.Namespace == "" {
		secList := new(core.SecretList)
		err := h.Client.List(ctx, secList,
			convertSecretListOptsFromVariableListOpts(opts))
		if err != nil {
			return nil, err
		}

		// Convert.
		return convertVariableListFromSecretList(secList, opts), nil
	}

	// Get.
	sec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      systemkuberes.VariablesDelegatedSecretName,
		},
	}
	err := h.Client.Get(ctx, ctrlcli.ObjectKeyFromObject(sec), sec)
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, err
		}
		// We return an empty list if the secret is not found.
		return &walrus.VariableList{}, nil
	}

	// Convert.
	vList := convertVariableListFromSecret(sec, opts)
	return vList, nil
}

func (h *VariableHandler) OnWatch(ctx context.Context, opts ctrlcli.ListOptions) (watch.Interface, error) {
	// Index.
	vraIndexer := map[string]walrus.Variable{} // [pn/en/vn] -> vra
	{
		listObj, err := h.OnList(ctx, opts)
		if err != nil {
			return nil, err
		}
		vList := listObj.(*walrus.VariableList)
		for i := range vList.Items {
			vraL1IndexKey := stringx.Join("/", vList.Items[i].Status.Project, vList.Items[i].Status.Environment)
			vraIndexKey := stringx.Join("/", vraL1IndexKey, vList.Items[i].Name)
			vraIndexer[vraIndexKey] = vList.Items[i]
		}
	}

	// Watch.
	uw, err := h.Client.(ctrlcli.WithWatch).Watch(ctx, new(core.SecretList),
		convertSecretListOptsFromVariableListOpts(opts))
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

				// Nothing to do
				if e.Object == nil {
					c <- e
					continue
				}

				// Type assert.
				sec, ok := e.Object.(*core.Secret)
				if !ok {
					c <- e
					continue
				}

				// Process bookmark.
				if e.Type == watch.Bookmark {
					e.Object = &walrus.Variable{ObjectMeta: sec.ObjectMeta}
					c <- e
					continue
				}

				notes := systemmeta.DescribeResourceNotes(sec, []string{"project", "environment"})
				vraL1IndexKey := stringx.Join("/", notes["project"], notes["environment"])
				vraIndexKeySet := sets.NewString()

				// Send.
				for name := range sec.Data {
					// Ignore if not be selected by `kubectl get --field-selector=metadata.namespace=...,metadata.name=...`.
					if fs := opts.FieldSelector; fs != nil &&
						!fs.Matches(fields.Set{"metadata.namespace": sec.Namespace, "metadata.name": name}) {
						continue
					}

					// Convert.
					vra := convertVariableFromSecret(sec, name)
					if vra == nil {
						continue
					}

					vraIndexKey := stringx.Join("/", vra.Status.Project, vra.Status.Environment, vra.Name)
					vraIndexKeySet.Insert(vraIndexKey)

					// Ignore if the same as previous.
					prevVra, ok := vraIndexer[vraIndexKey]
					switch {
					default:
						// ignore
						continue
					case !ok:
						// insert
						vraIndexer[vraIndexKey] = *vra
					case !vra.Equal(&prevVra):
						// update
						vraIndexer[vraIndexKey] = *vra
					}

					// Dispatch.
					e2 := e.DeepCopy()
					e2.Object = vra
					c <- *e2
				}

				// GC.
				for vraIndexKey := range vraIndexer {
					if !strings.HasPrefix(vraIndexKey, vraL1IndexKey+"/") {
						continue
					}
					switch {
					default:
						continue
					case e.Type == watch.Deleted:
					case !vraIndexKeySet.Has(vraIndexKey):
					}

					vra := vraIndexer[vraIndexKey]
					delete(vraIndexer, vraIndexKey)

					// Ignore if not be selected by `kubectl get --field-selector=metadata.namespace=...,metadata.name=...`.
					if fs := opts.FieldSelector; fs != nil &&
						!fs.Matches(fields.Set{"metadata.namespace": vra.Namespace, "metadata.name": vra.Name}) {
						continue
					}

					// Dispatch a delete event.
					e2 := e.DeepCopy()
					e2.Type = watch.Deleted
					e2.Object = &vra
					c <- *e2
				}
			}
		}
	})

	return dw, nil
}

func (h *VariableHandler) OnGet(ctx context.Context, name types.NamespacedName, opts ctrlcli.GetOptions) (runtime.Object, error) {
	// Get.
	sec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: name.Namespace,
			Name:      systemkuberes.VariablesDelegatedSecretName,
		},
	}
	err := h.Client.Get(ctx, ctrlcli.ObjectKeyFromObject(sec), sec)
	if err != nil {
		return nil, err
	}

	// Convert.
	vra := convertVariableFromSecret(sec, name.Name)
	if vra == nil {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("variables"), name.Name)
	}
	return vra, nil
}

func (h *VariableHandler) OnUpdate(ctx context.Context, obj, _ runtime.Object, opts ctrlcli.UpdateOptions) (runtime.Object, error) {
	// Validate.
	vra := obj.(*walrus.Variable)
	if vra.Spec.Value == nil {
		return nil, kerrors.NewInvalid(walrus.SchemeKind("variables"),
			vra.Name, field.ErrorList{field.Required(field.NewPath("spec.value"), "variable value is required")})
	}

	// Update.
	eSec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: vra.Namespace,
			Name:      systemkuberes.VariablesDelegatedSecretName,
		},
		Data: map[string][]byte{
			vra.Name: []byte(*vra.Spec.Value),
		},
	}
	eNotes := map[string]string{
		vra.Name + "-sensitive": strconv.FormatBool(vra.Spec.Sensitive),
	}
	systemmeta.NoteResource(eSec, "variables", eNotes)
	alignFn := func(aSec *core.Secret) (*core.Secret, bool, error) {
		// Validate.
		if aSec.Data == nil || aSec.Data[vra.Name] == nil {
			return nil, true, kerrors.NewNotFound(walrus.SchemeResource("variables"), vra.Name)
		}
		// Align data.
		aSec.Data[vra.Name] = eSec.Data[vra.Name]
		// Align delegated info.
		systemmeta.NoteResource(aSec, "variables", eNotes)
		return aSec, false, nil
	}

	sec, err := kubeclientset.UpdateWithCtrlClient(ctx, h.Client, eSec,
		kubeclientset.UpdateAfterAlign(alignFn))
	if err != nil {
		return nil, err
	}

	// Convert.
	vra = convertVariableFromSecret(sec, vra.Name)
	if vra == nil {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("variables"), vra.Name)
	}
	return vra, nil
}

func (h *VariableHandler) OnDelete(ctx context.Context, obj runtime.Object, opts ctrlcli.DeleteOptions) error {
	vra := obj.(*walrus.Variable)

	// Update.
	eSec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: vra.Namespace,
			Name:      systemkuberes.VariablesDelegatedSecretName,
		},
	}
	alignFn := func(aSec *core.Secret) (*core.Secret, bool, error) {
		// Validate.
		if aSec.Data == nil || aSec.Data[vra.Name] == nil {
			return nil, true, kerrors.NewNotFound(walrus.SchemeResource("variables"), vra.Name)
		}
		// Align data.
		delete(aSec.Data, vra.Name)
		// Align delegated info.
		systemmeta.PopResourceNotes(aSec, []string{
			vra.Name + "-uid",
			vra.Name + "-create-at",
			vra.Name + "-sensitive",
		})
		return aSec, false, nil
	}

	_, err := kubeclientset.UpdateWithCtrlClient(ctx, h.Client, eSec,
		kubeclientset.UpdateAfterAlign(alignFn))
	return err
}

func convertVariableFromSecret(sec *core.Secret, name string) *walrus.Variable {
	resType, notes := systemmeta.DescribeResource(sec)
	if resType != "variables" {
		return nil
	}

	// Filter out.
	if _, ok := sec.Data[name]; !ok {
		return nil
	}

	uid := sec.UID
	if uidS := notes[name+"-uid"]; len(uidS) != 0 {
		uid = types.UID(uidS)
	}
	createAt := sec.CreationTimestamp
	if createS := notes[name+"-create-at"]; len(createS) != 0 {
		if createAt_, err := time.Parse(time.RFC3339, createS); err == nil {
			createAt = meta.NewTime(createAt_)
		}
	}
	sensitive := notes[name+"-sensitive"] == "true"
	var (
		value  = []byte("")
		value_ = sec.Data[name]
	)
	if len(value_) != 0 && sensitive {
		value = []byte("(sensitive)")
	} else if len(value_) != 0 {
		value = value_
	}

	return &walrus.Variable{
		ObjectMeta: meta.ObjectMeta{
			Namespace:         sec.Namespace,
			Name:              name,
			UID:               uid,
			ResourceVersion:   sec.ResourceVersion,
			CreationTimestamp: createAt,
		},
		Spec: walrus.VariableSpec{
			Sensitive: sensitive,
		},
		Status: walrus.VariableStatus{
			Project:     notes["project"],
			Environment: notes["environment"],
			Value:       string(value),
			Value_:      string(value_),
		},
	}
}

func convertSecretListOptsFromVariableListOpts(in ctrlcli.ListOptions) (out *ctrlcli.ListOptions) {
	// Lock field selector.
	in.FieldSelector = fields.SelectorFromSet(fields.Set{
		"metadata.name": systemkuberes.VariablesDelegatedSecretName,
	})

	// Add necessary label selector.
	if lbs := systemmeta.LabelSelectorOf("variables"); in.LabelSelector == nil {
		in.LabelSelector = lbs
	} else {
		reqs, _ := lbs.Requirements()
		in.LabelSelector = in.LabelSelector.DeepCopySelector().Add(reqs...)
	}

	return &in
}

func convertVariableListFromSecret(sec *core.Secret, opts ctrlcli.ListOptions) *walrus.VariableList {
	resType, notes := systemmeta.DescribeResource(sec)
	if resType != "variables" {
		return &walrus.VariableList{}
	}

	vList := &walrus.VariableList{
		Items: make([]walrus.Variable, 0, len(sec.Data)),
	}

	for _, name := range sets.KeySet(sec.Data).UnsortedList() {
		// Ignore if not be selected by `kubectl get --field-selector=metadata.namespace=...,metadata.name=...`.
		if fs := opts.FieldSelector; fs != nil &&
			!fs.Matches(fields.Set{"metadata.namespace": sec.Namespace, "metadata.name": name}) {
			continue
		}

		uid := sec.UID
		if uidS := notes[name+"-uid"]; len(uidS) != 0 {
			uid = types.UID(uidS)
		}
		createAt := sec.CreationTimestamp
		if createS := notes[name+"-create-at"]; len(createS) != 0 {
			if createAt_, err := time.Parse(time.RFC3339, createS); err == nil {
				createAt = meta.NewTime(createAt_)
			}
		}
		sensitive := notes[name+"-sensitive"] == "true"
		var (
			value  = []byte("")
			value_ = sec.Data[name]
		)
		if len(value_) != 0 && sensitive {
			value = []byte("(sensitive)")
		} else if len(value_) != 0 {
			value = value_
		}

		vList.Items = append(vList.Items, walrus.Variable{
			ObjectMeta: meta.ObjectMeta{
				Namespace:         sec.Namespace,
				Name:              name,
				UID:               uid,
				ResourceVersion:   sec.ResourceVersion,
				CreationTimestamp: createAt,
			},
			Spec: walrus.VariableSpec{
				Sensitive: sensitive,
			},
			Status: walrus.VariableStatus{
				Project:     notes["project"],
				Environment: notes["environment"],
				Value:       string(value),
				Value_:      string(value_),
			},
		})
	}

	return vList
}

func convertVariableListFromSecretList(secList *core.SecretList, opts ctrlcli.ListOptions) *walrus.VariableList {
	// Sort by resource version.
	sort.SliceStable(secList.Items, func(i, j int) bool {
		l, r := secList.Items[i].ResourceVersion, secList.Items[j].ResourceVersion
		return len(l) < len(r) ||
			(len(l) == len(r) && l < r)
	})

	vList := &walrus.VariableList{
		Items: make([]walrus.Variable, 0, len(secList.Items)),
	}

	for i := range secList.Items {
		svList := convertVariableListFromSecret(&secList.Items[i], opts)
		if svList == nil {
			continue
		}
		vList.Items = append(vList.Items, svList.Items...)
	}

	return vList
}
