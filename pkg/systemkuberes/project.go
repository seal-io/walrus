package systemkuberes

import (
	"context"
	"fmt"

	"golang.org/x/exp/maps"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
	"github.com/seal-io/walrus/pkg/systemmeta"
)

// DefaultProjectName is the Kubernetes Namespace name for the default project.
const DefaultProjectName = core.NamespaceDefault

// InstallDefaultProject creates the default project, alias of the Kubernetes Namespace default.
func InstallDefaultProject(ctx context.Context, cli clientset.Interface) error {
	err := review.CanDoUpdate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    core.SchemeGroupVersion.Group,
				Version:  core.SchemeGroupVersion.Version,
				Resource: "namespaces",
				Name:     DefaultProjectName,
			},
		},
		review.WithCreate(),
	)
	if err != nil {
		return err
	}

	nsCli := cli.CoreV1().Namespaces()

	eNs := &core.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: DefaultProjectName,
		},
	}
	systemmeta.NoteResource(eNs, "projects", map[string]string{
		"displayName": "Default Project",
		"description": "The default project created by Walrus.",
	})
	systemmeta.Lock(eNs)
	alignFn := func(aNs *core.Namespace) (_ *core.Namespace, skip bool, err error) {
		// Align delegated info.
		aResType, aNotes := systemmeta.DescribeResource(aNs)
		eResType, eNotes := systemmeta.DescribeResource(eNs)
		if aResType != eResType || !sets.KeySet(aNotes).HasAll(maps.Keys(eNotes)...) {
			eNotes := maps.Clone(eNotes)
			maps.Copy(eNotes, aNotes)
			systemmeta.NoteResource(aNs, eResType, eNotes)
			skip = false
		}
		// Align lock.
		if !systemmeta.Lock(aNs) {
			skip = false
		}
		return aNs, skip, nil
	}

	_, err = kubeclientset.Update(ctx, nsCli, eNs,
		kubeclientset.UpdateAfterAlign(alignFn),
		kubeclientset.UpdateOrCreate[*core.Namespace]())
	if err != nil {
		return fmt.Errorf("install builtin project: %w", err)
	}

	return nil
}
