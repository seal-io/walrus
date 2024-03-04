package walrus

import (
	"context"
	"errors"

	"github.com/seal-io/utils/pools/gopool"
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
	"github.com/seal-io/walrus/pkg/systemmeta"
	"github.com/seal-io/walrus/pkg/systemsetting"
)

// SettingHandler handles v1.Setting objects.
//
// SettingHandler maps all v1.Setting objects to a Kubernetes Secret resource,
// which is named as "walrus-system/walrus-settings".
//
// Each v1.Setting object records as a key-value pair in the Secret's Data field.
type SettingHandler struct {
	extensionapi.ObjectInfo
	extensionapi.ListWatchOperation
	extensionapi.GetOperation
	extensionapi.UpdateOperation

	Client ctrlcli.Client
}

func (h *SettingHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Configure field indexer.
	fi := opts.Manager.GetFieldIndexer()
	err = fi.IndexField(ctx, &core.Secret{}, "metadata.name",
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
	gvr = walrus.SchemeGroupVersionResource("settings")

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
			})
		if err != nil {
			return
		}
	}

	// As storage.
	h.ObjectInfo = &walrus.Setting{}
	h.ListWatchOperation = extensionapi.WithListWatch(tc, h)
	h.GetOperation = extensionapi.WithGet(h)
	h.UpdateOperation = extensionapi.WithUpdate(h)

	// Set client.
	h.Client = opts.Manager.GetClient()

	return
}

var (
	_ rest.Storage = (*SettingHandler)(nil)
	_ rest.Lister  = (*SettingHandler)(nil)
	_ rest.Watcher = (*SettingHandler)(nil)
	_ rest.Getter  = (*SettingHandler)(nil)
	_ rest.Updater = (*SettingHandler)(nil)
	_ rest.Patcher = (*SettingHandler)(nil)
)

func (h *SettingHandler) New() runtime.Object {
	return &walrus.Setting{}
}

func (h *SettingHandler) Destroy() {}

func (h *SettingHandler) NewList() runtime.Object {
	return &walrus.SettingList{}
}

func (h *SettingHandler) OnList(ctx context.Context, opts ctrlcli.ListOptions) (runtime.Object, error) {
	// Support watch with `kubectl get -A`.
	if opts.Namespace == "" {
		opts.Namespace = systemsetting.DelegatedSecretNamespace
	}

	// Get.
	sec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: opts.Namespace,
			Name:      systemsetting.DelegatedSecretName,
		},
	}
	err := h.Client.Get(ctx, ctrlcli.ObjectKeyFromObject(sec), sec)
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, err
		}
		// We return an empty list if the secret is not found.
		return &walrus.SettingList{}, nil
	}

	// Convert.
	sList := convertSettingListFromSecret(sec, opts)
	return sList, nil
}

func (h *SettingHandler) OnWatch(ctx context.Context, opts ctrlcli.ListOptions) (watch.Interface, error) {
	// Support watch with `kubectl get -A`.
	if opts.Namespace == "" {
		opts.Namespace = systemsetting.DelegatedSecretNamespace
	}

	// List and index.
	setIndexer := map[string]walrus.Setting{} // [sn] -> set
	{
		listObj, err := h.OnList(ctx, opts)
		if err != nil {
			return nil, err
		}
		sList := listObj.(*walrus.SettingList)
		for i := range sList.Items {
			setIndexKey := sList.Items[i].Name
			setIndexer[setIndexKey] = sList.Items[i]
		}
	}

	// Watch.
	uw, err := h.Client.(ctrlcli.WithWatch).Watch(ctx, new(core.SecretList),
		convertSecretListOptsFromSettingListOpts(opts))
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
				sec, ok := e.Object.(*core.Secret)
				if !ok {
					c <- e
					continue
				}

				// Process bookmark.
				if e.Type == watch.Bookmark {
					e.Object = &walrus.Setting{ObjectMeta: sec.ObjectMeta}
					c <- e
					continue
				}

				// Send.
				for name := range sec.Data {
					// Ignore if not be selected by `kubectl get --field-selector=metadata.name=...`.
					if fs := opts.FieldSelector; fs != nil &&
						!fs.Matches(fields.Set{"metadata.name": name}) {
						continue
					}

					// Convert.
					set := convertSettingFromSecret(sec, name)
					if set == nil {
						continue
					}

					// Ignore if the same as previous.
					setIndexKey := set.Name
					prevSet, ok := setIndexer[setIndexKey]
					switch {
					default:
						// ignore
						continue
					case !ok:
						// insert
						setIndexer[setIndexKey] = *set
					case !set.Equal(&prevSet):
						// update
						setIndexer[setIndexKey] = *set
					}

					// Dispatch.
					e2 := e.DeepCopy()
					e2.Object = set
					c <- *e2
				}
			}
		}
	})

	return dw, nil
}

func (h *SettingHandler) OnGet(ctx context.Context, key types.NamespacedName, opts ctrlcli.GetOptions) (runtime.Object, error) {
	// Get.
	sec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: key.Namespace,
			Name:      systemsetting.DelegatedSecretName,
		},
	}
	err := h.Client.Get(ctx, ctrlcli.ObjectKeyFromObject(sec), sec)
	if err != nil {
		return nil, err
	}

	// Convert.
	set := convertSettingFromSecret(sec, key.Name)
	if set == nil {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("settings"), key.Name)
	}
	return set, nil
}

func (h *SettingHandler) OnUpdate(ctx context.Context, obj, _ runtime.Object, opts ctrlcli.UpdateOptions) (runtime.Object, error) {
	// Validate.
	set := obj.(*walrus.Setting)
	if set.Namespace != systemsetting.DelegatedSecretNamespace {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("settings"), set.Name)
	}
	s, ok := systemsetting.Index(set.Name)
	if !ok || !s.Editable() {
		return nil, kerrors.NewForbidden(walrus.SchemeResource("settings"), set.Name,
			errors.New("setting is not editable"))
	}
	if set.Spec.Value == nil {
		return nil, kerrors.NewInvalid(walrus.SchemeKind("settings"), set.Name,
			field.ErrorList{field.Required(field.NewPath("spec.value"), "setting value is required")})
	}

	// Update.
	err := s.Configure(ctx, *set.Spec.Value)
	if err != nil {
		return nil, kerrors.NewConflict(walrus.SchemeResource("settings"), set.Name, err)
	}

	// Get.
	return h.OnGet(ctx, ctrlcli.ObjectKeyFromObject(set), ctrlcli.GetOptions{})
}

func convertSecretListOptsFromSettingListOpts(in ctrlcli.ListOptions) (out *ctrlcli.ListOptions) {
	// Lock field selector.
	in.FieldSelector = fields.SelectorFromSet(fields.Set{
		"metadata.namespace": in.Namespace,
		"metadata.name":      systemsetting.DelegatedSecretName,
	})

	// Add necessary label selector.
	if lbs := systemmeta.GetResourcesLabelSelectorOfType("settings"); in.LabelSelector == nil {
		in.LabelSelector = lbs
	} else {
		reqs, _ := lbs.Requirements()
		in.LabelSelector = in.LabelSelector.DeepCopySelector().Add(reqs...)
	}

	return &in
}

func convertSettingFromSecret(sec *core.Secret, reqName string) *walrus.Setting {
	resType := systemmeta.DescribeResourceType(sec)
	if resType != "settings" {
		return nil
	}

	// Filter out.
	s, ok := systemsetting.Index(reqName)
	if !ok || s.Private() || sec.Data == nil {
		return nil
	}

	uid := sec.UID
	if uidS := systemmeta.DescribeResourceNote(sec, reqName+"-uid"); len(uidS) != 0 {
		uid = types.UID(uidS)
	}
	var (
		value  = []byte("")
		value_ = sec.Data[reqName]
	)
	if len(value_) != 0 && s.Sensitive() {
		value = []byte("(sensitive)")
	} else if len(value_) != 0 {
		value = value_
	}

	return &walrus.Setting{
		ObjectMeta: meta.ObjectMeta{
			Namespace:         sec.Namespace,
			Name:              reqName,
			UID:               uid,
			ResourceVersion:   sec.ResourceVersion,
			CreationTimestamp: sec.CreationTimestamp,
		},
		Status: walrus.SettingStatus{
			Description: s.Description(),
			Hidden:      s.Hidden(),
			Editable:    s.Editable(),
			Sensitive:   s.Sensitive(),
			Value:       string(value),
			Value_:      string(value_),
		},
	}
}

func convertSettingListFromSecret(sec *core.Secret, opts ctrlcli.ListOptions) *walrus.SettingList {
	resType, notes := systemmeta.DescribeResource(sec)
	if resType != "settings" {
		return &walrus.SettingList{}
	}

	sList := &walrus.SettingList{
		Items: make([]walrus.Setting, 0, len(sec.Data)),
	}

	// Sort by name.
	for _, name := range sets.List(sets.KeySet(sec.Data)) {
		// Ignore if not be selected by `kubectl get --field-selector=metadata.name=...`.
		if fs := opts.FieldSelector; fs != nil &&
			!fs.Matches(fields.Set{"metadata.name": name}) {
			continue
		}

		// Filter out.
		s, ok := systemsetting.Index(name)
		if !ok || s.Private() {
			continue
		}

		uid := sec.UID
		if uidS := notes[name+"-uid"]; len(uidS) != 0 {
			uid = types.UID(uidS)
		}
		var (
			value  = []byte("")
			value_ = sec.Data[name]
		)
		if len(value_) != 0 && s.Sensitive() {
			value = []byte("(sensitive)")
		} else if len(value_) != 0 {
			value = value_
		}

		sList.Items = append(sList.Items, walrus.Setting{
			ObjectMeta: meta.ObjectMeta{
				Namespace:         sec.Namespace,
				Name:              name,
				UID:               uid,
				ResourceVersion:   sec.ResourceVersion,
				CreationTimestamp: sec.CreationTimestamp,
			},
			Status: walrus.SettingStatus{
				Description: s.Description(),
				Hidden:      s.Hidden(),
				Editable:    s.Editable(),
				Sensitive:   s.Sensitive(),
				Value:       string(value),
				Value_:      string(value_),
			},
		})
	}

	return sList
}
