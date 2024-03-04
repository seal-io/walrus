package systemkuberes

import (
	"context"
	"fmt"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
)

// DefaultProjectName is the Kubernetes Namespace name for the default project.
const DefaultProjectName = core.NamespaceDefault

// InstallDefaultProject creates the default project, alias of the Kubernetes Namespace default.
func InstallDefaultProject(ctx context.Context, cli clientset.Interface) error {
	err := review.CanDoUpdate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    walrus.SchemeGroupVersion.Group,
				Version:  walrus.SchemeGroupVersion.Version,
				Resource: "projects",
			},
		},
		review.WithCreateIfNotExisted(),
	)
	if err != nil {
		return err
	}

	projCli := cli.WalrusV1().Projects(SystemNamespaceName)
	proj := &walrus.Project{
		ObjectMeta: meta.ObjectMeta{
			Namespace: SystemNamespaceName,
			Name:      DefaultProjectName,
		},
		Spec: walrus.ProjectSpec{
			DisplayName: "Default Project",
			Description: "The default project created by Walrus.",
		},
	}

	_, err = kubeclientset.Create(ctx, projCli, proj)
	if err != nil {
		return fmt.Errorf("install default project: %w", err)
	}

	return nil
}
