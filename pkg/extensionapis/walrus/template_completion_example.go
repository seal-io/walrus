package walrus

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/seal-io/utils/stringx"
	"golang.org/x/text/language"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/registry/rest"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/extensionapi"
	"github.com/seal-io/walrus/pkg/i18n"
	"github.com/seal-io/walrus/pkg/systemkuberes"
)

// TemplateCompletionExampleHandler handles v1.TemplateCompletionExample objects.
//
// TemplateCompletionExampleHandler maintains an in-memory v1.TemplateCompletionExampleList object to serve the get/list operation.
//
// All the v1.TemplateCompletionExample objects are hard-coded in the source code.
//
// For list operation, it supports the `accept-language` label selector,
// to return the localized prompt of each v1.TemplateCompletionExample object.
type TemplateCompletionExampleHandler struct {
	extensionapi.ObjectInfo
	extensionapi.GetOperation
	extensionapi.ListOperation

	InMemoryListObject walrus.TemplateCompletionExampleList
}

func (h *TemplateCompletionExampleHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Declare GVR.
	gvr = walrus.SchemeGroupVersionResource("templatecompletionexamples")

	// Create table convertor to pretty the kubectl's output.
	var tc rest.TableConvertor
	{
		tc, err = extensionapi.NewJSONPathTableConvertor(
			extensionapi.JSONPathTableColumnDefinition{
				TableColumnDefinition: meta.TableColumnDefinition{
					Name: "Purpose",
					Type: "string",
				},
				JSONPath: ".status.purpose",
			})
		if err != nil {
			return
		}
	}

	// As storage.
	h.ObjectInfo = &walrus.TemplateCompletionExample{}
	h.GetOperation = extensionapi.WithGet(h)
	h.ListOperation = extensionapi.WithList(tc, h)

	// Hard-code the v1.TemplateCompletionExample objects.
	h.InMemoryListObject.Items = loadTemplateCompletionExamples(systemkuberes.SystemNamespaceName)

	return
}

var (
	_ rest.Storage = (*TemplateCompletionExampleHandler)(nil)
	_ rest.Lister  = (*TemplateCompletionExampleHandler)(nil)
	_ rest.Getter  = (*TemplateCompletionExampleHandler)(nil)
)

func (h *TemplateCompletionExampleHandler) New() runtime.Object {
	return &walrus.TemplateCompletionExample{}
}

func (h *TemplateCompletionExampleHandler) Destroy() {}

func (h *TemplateCompletionExampleHandler) NewList() runtime.Object {
	return &walrus.TemplateCompletionExampleList{}
}

func (h *TemplateCompletionExampleHandler) OnList(ctx context.Context, opts ctrlcli.ListOptions) (runtime.Object, error) {
	// Support watch with `kubectl get -A`.
	if opts.Namespace == "" {
		opts.Namespace = systemkuberes.SystemNamespaceName
	}

	// Only support list in system namespace.
	if opts.Namespace == systemkuberes.SystemNamespaceName {
		listObj := h.InMemoryListObject.DeepCopy()

		if ls := opts.Raw.LabelSelector; ls != "" && strings.HasPrefix(ls, "accept-language=") {
			accLang := strings.TrimPrefix(ls, "accept-language=")
			tags, _, _ := language.ParseAcceptLanguage(accLang)
			for i := range listObj.Items {
				listObj.Items[i].Status.Purpose = i18n.T(listObj.Items[i].Status.Purpose, tags...)
				listObj.Items[i].Status.Prompt = i18n.T(listObj.Items[i].Status.Prompt, tags...)
			}
		}

		return listObj, nil
	}

	return &walrus.TemplateCompletionExampleList{}, nil
}

func (h *TemplateCompletionExampleHandler) OnGet(ctx context.Context, key types.NamespacedName, opts ctrlcli.GetOptions) (runtime.Object, error) {
	// Only support get in system namespace.
	if key.Namespace == systemkuberes.SystemNamespaceName {
		for _, item := range h.InMemoryListObject.Items {
			if item.Name == key.Name {
				return item.DeepCopy(), nil
			}
		}
	}

	return nil, kerrors.NewNotFound(walrus.SchemeResource("templatecompletionexamples"), key.Name)
}

func loadTemplateCompletionExamples(namespace string) []walrus.TemplateCompletionExample {
	// TODO(thxCode): move this hard-cord map outside?
	entries := []walrus.TemplateCompletionExampleStatus{
		{
			Purpose: i18n.ExampleKubernetesName,
			Prompt:  i18n.ExampleKubernetesPrompt,
		},
		{
			Purpose: i18n.ExampleAlibabaCloudName,
			Prompt:  i18n.ExampleAlibabaCloudPrompt,
		},
		{
			Purpose: i18n.ExampleELKName,
			Prompt:  i18n.ExampleELKPrompt,
		},
	}

	exps := make([]walrus.TemplateCompletionExample, 0, len(entries))
	for i := range entries {
		exp := walrus.TemplateCompletionExample{
			ObjectMeta: meta.ObjectMeta{
				Namespace:         namespace,
				Name:              stringx.Dasherize(entries[i].Purpose),
				UID:               types.UID(uuid.NewMD5(uuid.Nil, []byte(entries[i].Purpose)).String()), // Create a deterministic UID.
				CreationTimestamp: meta.Now(),
				ResourceVersion:   stringx.FromInt(i + 1),
			},
			Status: entries[i],
		}

		exps = append(exps, exp)
	}

	return exps
}
