package walrus

import (
	"context"
	"fmt"

	rbac "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	authnuser "k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/registry/rest"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/extensionapi"
	"github.com/seal-io/walrus/pkg/systemauthz"
	"github.com/seal-io/walrus/pkg/systemkuberes"
	"github.com/seal-io/walrus/pkg/systemmeta"
)

type ProjectSubjectsHandler struct {
	extensionapi.ObjectInfo
	extensionapi.GetOperation
	extensionapi.UpdateOperation

	Client ctrlcli.Client
}

func newProjectSubjectsHandler(opts extensionapi.SetupOptions) *ProjectSubjectsHandler {
	h := &ProjectSubjectsHandler{}

	// As storage.
	h.ObjectInfo = &walrus.ProjectSubjects{}
	h.GetOperation = extensionapi.WithGet(h)
	h.UpdateOperation = extensionapi.WithUpdate(h)

	// Set client.
	h.Client = opts.Manager.GetClient()

	return h
}

var (
	_ rest.Storage = (*ProjectSubjectsHandler)(nil)
	_ rest.Getter  = (*ProjectSubjectsHandler)(nil)
	_ rest.Updater = (*ProjectSubjectsHandler)(nil)
	_ rest.Patcher = (*ProjectSubjectsHandler)(nil)
)

func (h *ProjectSubjectsHandler) New() runtime.Object {
	return &walrus.ProjectSubjects{}
}

func (h *ProjectSubjectsHandler) Destroy() {}

func (h *ProjectSubjectsHandler) OnGet(ctx context.Context, key types.NamespacedName, opts ctrlcli.GetOptions) (runtime.Object, error) {
	// Validate.
	if key.Namespace != systemkuberes.SystemNamespaceName {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("projectsubjects"), key.Name)
	}

	// List.
	crbList := new(rbac.ClusterRoleBindingList)
	err := h.Client.List(ctx, crbList,
		ctrlcli.MatchingLabelsSelector{
			Selector: systemmeta.GetResourcesLabelSelectorOfType("rolebindings"),
		})
	if err != nil {
		return nil, kerrors.NewInternalError(err)
	}
	crbList = systemmeta.FilterResourceListByNotes(crbList, "project", key.Name)

	// Convert.
	psbjs := convertProjectSubjectsFromClusterRoleBindingList(crbList)
	if psbjs == nil {
		return nil, kerrors.NewNotFound(walrus.SchemeResource("projectsubjects"), key.Name)
	}

	// Get and refill.
	proj := new(walrus.Project)
	err = h.Client.Get(ctx, key, proj)
	if err != nil {
		return nil, kerrors.NewInternalError(err)
	}
	psbjs.ObjectMeta = proj.ObjectMeta

	return psbjs, nil
}

func (h *ProjectSubjectsHandler) OnUpdate(ctx context.Context, obj, objOld runtime.Object, opts ctrlcli.UpdateOptions) (runtime.Object, error) {
	psbjs, psbjsOld := obj.(*walrus.ProjectSubjects), objOld.(*walrus.ProjectSubjects)

	// Figure out delta.
	psbjsReverseIndex := make(map[walrus.ProjectSubject]int)
	{
		for i := range psbjs.Items {
			// Default.
			if psbjs.Items[i].Kind == "" {
				psbjs.Items[i].Kind = "User"
			}
			psbjsReverseIndex[psbjs.Items[i]] = i
		}
	}
	psbjsOldSet := make(map[walrus.ProjectSubject]sets.Empty)
	{
		for i := range psbjsOld.Items {
			psbjsOldSet[psbjsOld.Items[i]] = sets.Empty{}
		}
	}
	for psbj := range psbjsReverseIndex {
		// Delete the one exists in both of the new set and old set,
		// then the remaining items of the new set are need to create,
		// and the remaining items of the old set are need to delete.
		if _, existed := psbjsOldSet[psbj]; existed {
			delete(psbjsReverseIndex, psbj)
			delete(psbjsOldSet, psbj)
		}
	}

	// Validate.
	var errs field.ErrorList
	for psbj, psbjIdx := range psbjsReverseIndex {
		if psbj.Name == "" {
			errs = append(errs, field.Forbidden(
				field.NewPath(fmt.Sprintf("items[%d].name", psbjIdx)), "blank string"))
		}
		if err := psbj.Role.Validate(); err != nil {
			errs = append(errs, field.Invalid(
				field.NewPath(fmt.Sprintf("items[%d].role", psbjIdx)), psbj.Role, err.Error()))
		}
	}
	if len(errs) > 0 {
		return nil, kerrors.NewInvalid(walrus.SchemeKind("projectsubjects"), psbjs.Name, errs)
	}

	// Unbind.
	for psbj := range psbjsOldSet {
		uInfo := &authnuser.DefaultInfo{
			Name: psbj.Name,
		}
		err := systemauthz.UnbindProjectSubjectRoleFor(ctx, h.Client, psbj.Role, uInfo)
		if err != nil {
			return nil, kerrors.NewInternalError(fmt.Errorf("unbind project subject role: %w", err))
		}
	}

	// NB(thxCode): without retrieving again,
	// we can simply construct a Project object from the old ProjectSubjects.
	proj := &walrus.Project{
		ObjectMeta: psbjsOld.ObjectMeta,
	}

	// Bind.
	for psbj := range psbjsReverseIndex {
		uInfo := &authnuser.DefaultInfo{
			Name: psbj.Name,
		}
		err := systemauthz.BindProjectSubjectRoleFor(ctx, h.Client, proj, psbj.Role, uInfo)
		if err != nil {
			return nil, kerrors.NewInternalError(fmt.Errorf("bind project subject role: %w", err))
		}
	}

	// Get.
	return h.OnGet(ctx, ctrlcli.ObjectKeyFromObject(psbjs), ctrlcli.GetOptions{})
}

func convertProjectSubjectFromClusterRoleBinding(crb *rbac.ClusterRoleBinding) *walrus.ProjectSubject {
	if crb == nil || len(crb.Subjects) != 1 {
		return nil
	}

	psbjr := walrus.ProjectSubjectRole(crb.RoleRef.Name)
	if psbjr.Validate() != nil {
		return nil
	}

	psbj := &walrus.ProjectSubject{
		Name: crb.Subjects[0].Name,
		Kind: crb.Subjects[0].Kind,
		Role: psbjr,
	}
	return psbj
}

func convertProjectSubjectsFromClusterRoleBindingList(crbList *rbac.ClusterRoleBindingList) *walrus.ProjectSubjects {
	if crbList == nil {
		return nil
	}

	psbjs := &walrus.ProjectSubjects{
		Items: make([]walrus.ProjectSubject, 0, len(crbList.Items)),
	}

	for i := range crbList.Items {
		psbj := convertProjectSubjectFromClusterRoleBinding(&crbList.Items[i])
		if psbj == nil {
			continue
		}
		psbjs.Items = append(psbjs.Items, *psbj)
	}

	return psbjs
}
