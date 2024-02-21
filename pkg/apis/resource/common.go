package resource

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/deployer"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	"github.com/seal-io/walrus/utils/errorx"
)

func getDeployer(ctx context.Context, kubeConfig *rest.Config) (deptypes.Deployer, error) {
	dep, err := deployer.Get(ctx, deptypes.CreateOptions{
		Type:       types.DeployerTypeTF,
		KubeConfig: kubeConfig,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to get deployer")
	}

	return dep, nil
}
