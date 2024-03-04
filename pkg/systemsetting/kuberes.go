package systemsetting

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
	"github.com/seal-io/walrus/pkg/systemkuberes"
	"github.com/seal-io/walrus/pkg/systemmeta"
)

const (
	// DelegatedSecretNamespace is the delegated Kubernetes Secret namespace for the settings.
	DelegatedSecretNamespace = systemkuberes.SystemNamespaceName

	// DelegatedSecretName is the delegated Kubernetes Secret name for the settings.
	DelegatedSecretName = "walrus-settings"
)

// Initialize initializes Kubernetes resources for settings.
//
// Initialize creates the delegated Kubernetes Secret for settings.
func Initialize(ctx context.Context, cli clientset.Interface) error {
	err := systemkuberes.InstallSystemNamespace(ctx, cli)
	if err != nil {
		return err
	}

	err = review.CanDoUpdate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    core.SchemeGroupVersion.Group,
				Version:  core.SchemeGroupVersion.Version,
				Resource: "secrets",
			},
		},
		review.WithCreate(),
	)
	if err != nil {
		return err
	}

	secCli := cli.CoreV1().Secrets(systemkuberes.SystemNamespaceName)

	eSec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name: DelegatedSecretName,
		},
		Data: map[string][]byte{},
	}
	eResType := "settings"
	eNotes := map[string]string{}
	for name, defVal := range Initials() {
		eSec.Data[name] = []byte(defVal)
		eNotes[name+"-uid"] = uuid.NewString()
	}
	systemmeta.NoteResource(eSec, eResType, eNotes)
	alignFn := func(aSec *core.Secret) (_ *core.Secret, skip bool, err error) {
		// Align delegated info.
		aResType, aNotes := systemmeta.DescribeResource(aSec)
		if aResType != eResType || !sets.KeySet(aNotes).HasAll(maps.Keys(eNotes)...) {
			eNotes := maps.Clone(eNotes)
			maps.Copy(eNotes, aNotes)
			systemmeta.NoteResource(aSec, eResType, eNotes)
			skip = false
		}
		// Align data.
		for k := range eSec.Data {
			if _, ok := aSec.Data[k]; !ok {
				aSec.Data[k] = eSec.Data[k]
				skip = false
			}
		}
		return aSec, skip, nil
	}

	_, err = kubeclientset.Update(ctx, secCli, eSec,
		kubeclientset.UpdateAfterAlign(alignFn),
		kubeclientset.UpdateOrCreate[*core.Secret]())
	if err != nil {
		return err
	}

	return nil
}
