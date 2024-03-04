package systemkuberes

import (
	"context"
	"fmt"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
	"github.com/seal-io/walrus/pkg/system"
)

// DefaultEnvironmentName is the Kubernetes Namespace name for the default environment.
const DefaultEnvironmentName = DefaultProjectName + "-local"

// InstallDefaultEnvironment creates the default environment, alias to Kubernetes Namespace default-local.
func InstallDefaultEnvironment(ctx context.Context, cli clientset.Interface) error {
	err := review.CanDoUpdate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    walrus.SchemeGroupVersion.Group,
				Version:  walrus.SchemeGroupVersion.Version,
				Resource: "environments",
			},
		},
		review.WithCreateIfNotExisted(),
	)
	if err != nil {
		return err
	}

	envCli := cli.WalrusV1().Environments(DefaultProjectName)
	env := &walrus.Environment{
		ObjectMeta: meta.ObjectMeta{
			Namespace: DefaultProjectName,
			Name:      DefaultEnvironmentName,
		},
		Spec: walrus.EnvironmentSpec{
			Type: func() walrus.EnvironmentType {
				if system.LoopbackKubeInside.Get() {
					return walrus.EnvironmentTypeProduction
				}
				return walrus.EnvironmentTypeDevelopment
			}(),
			DisplayName: "Default Environment",
			Description: "The default environment created by Walrus.",
		},
	}

	_, err = kubeclientset.Create(ctx, envCli, env)
	if err != nil {
		return fmt.Errorf("install default environment: %w", err)
	}

	return nil
}
