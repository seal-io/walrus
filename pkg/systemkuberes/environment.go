package systemkuberes

import (
	"context"
	"fmt"

	"golang.org/x/exp/maps"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
	"github.com/seal-io/walrus/pkg/kubemeta"
	"github.com/seal-io/walrus/pkg/system"
	"github.com/seal-io/walrus/pkg/systemmeta"
)

// DefaultEnvironmentName is the Kubernetes Namespace name for the default environment.
const DefaultEnvironmentName = DefaultProjectName + "-local"

// InstallDefaultEnvironment creates the default environment, alias to Kubernetes Namespace default-local.
func InstallDefaultEnvironment(ctx context.Context, cli clientset.Interface) error {
	err := review.CanDoUpdate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    core.SchemeGroupVersion.Group,
				Version:  core.SchemeGroupVersion.Version,
				Resource: "namespaces",
				Name:     DefaultEnvironmentName,
			},
		},
		review.WithCreate(),
	)
	if err != nil {
		return err
	}

	nsCli := cli.CoreV1().Namespaces()

	// Find owner.
	owner, err := nsCli.Get(ctx, DefaultProjectName, meta.GetOptions{})
	if err != nil {
		return fmt.Errorf("get owner namespace %q: %w", DefaultProjectName, err)
	}

	eNs := &core.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: DefaultEnvironmentName,
		},
	}
	systemmeta.NoteResource(eNs, "environments", map[string]string{
		"type": func() string {
			if system.LoopbackKubeInside.Get() {
				return walrus.EnvironmentTypeProduction.String()
			}
			return walrus.EnvironmentTypeDevelopment.String()
		}(),
		"displayName": "Local Environment",
		"description": "The default environment created by Walrus.",
	})
	systemmeta.Lock(eNs)
	kubemeta.ControlOn(eNs, owner, walrus.SchemeGroupVersion.WithKind("Project"))
	alignFn := func(aNs *core.Namespace) (_ *core.Namespace, skip bool, err error) {
		// Align owner.
		if !kubemeta.IsControlledBy(aNs, owner) {
			aNs.OwnerReferences = append(aNs.OwnerReferences,
				*kubemeta.NewControllerRef(owner, walrus.SchemeGroupVersion.WithKind("Project")))
			skip = false
		}
		// Align delegated info.
		aResType, aNotes := systemmeta.DescribeResource(aNs)
		eResType, eNotes := systemmeta.DescribeResource(eNs)
		if aResType != eResType || !sets.KeySet(aNotes).HasAll(maps.Keys(eNotes)...) {
			eNotes := maps.Clone(eNotes)
			maps.Copy(eNotes, aNotes)
			systemmeta.NoteResource(aNs, eResType, eNotes)
			skip = false
		}
		// Align locker.
		if !systemmeta.Lock(aNs) {
			skip = false
		}
		return aNs, skip, nil
	}

	_, err = kubeclientset.Update(ctx, nsCli, eNs,
		kubeclientset.UpdateAfterAlign(alignFn),
		kubeclientset.UpdateOrCreate[*core.Namespace]())
	if err != nil {
		return fmt.Errorf("install builtin environment: %w", err)
	}

	return nil
}
