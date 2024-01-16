package resource

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer"
	deployertf "github.com/seal-io/walrus/pkg/deployer/terraform"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/utils/errorx"
)

func DeleteResources(req CollectionDeleteRequest, mc model.ClientSet, kubeConfig *rest.Config) error {
	ids := req.IDs()

	return mc.WithTx(req.Context, func(tx *model.Tx) error {
		if req.WithoutCleanup {
			// Do not clean deployed native resources.
			_, err := tx.Resources().Delete().
				Where(resource.IDIn(ids...)).
				Exec(req.Context)

			return err
		}

		resources, err := tx.Resources().Query().
			Where(resource.IDIn(ids...)).
			All(req.Context)
		if err != nil {
			return err
		}

		resources, err = pkgresource.ReverseTopologicalSortResources(resources)
		if err != nil {
			return err
		}

		deployerOpts := deptypes.CreateOptions{
			Type:       deployertf.DeployerType,
			KubeConfig: kubeConfig,
		}

		dp, err := deployer.Get(req.Context, deployerOpts)
		if err != nil {
			return err
		}

		destroyOpts := pkgresource.Options{
			Deployer: dp,
		}

		for _, s := range resources {
			if err := pkgresource.SetSubjectID(req.Context, s); err != nil {
				return err
			}

			s, err = tx.Resources().UpdateOne(s).
				SetAnnotations(s.Annotations).
				Save(req.Context)
			if err != nil {
				return err
			}

			err = pkgresource.Destroy(req.Context, tx, s, destroyOpts)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func UpgradeResources(req CollectionRouteUpgradeRequest, mc model.ClientSet, kubeConfig *rest.Config) error {
	// Make sure the resources are upgraded in topological order.
	entities, err := pkgresource.TopologicalSortResources(req.Model())
	if err != nil {
		return err
	}

	for _, entity := range entities {
		entity.ChangeComment = req.ChangeComment

		if err = upgrade(req.Context, kubeConfig, mc, entity, req.Draft); err != nil {
			return err
		}
	}

	return nil
}

func upgrade(
	ctx context.Context,
	kubeConfig *rest.Config,
	mc model.ClientSet,
	entity *model.Resource,
	draft bool,
) error {
	if draft {
		_, err := mc.Resources().
			UpdateOne(entity).
			Set(entity).
			Save(ctx)

		return err
	}

	// Update resource, mark status from deploying.
	status.ResourceStatusDeployed.Reset(entity, "Upgrading")
	entity.Status.SetSummary(status.WalkResource(&entity.Status))

	if err := pkgresource.SetSubjectID(ctx, entity); err != nil {
		return err
	}

	err := mc.WithTx(ctx, func(tx *model.Tx) (err error) {
		entity, err = tx.Resources().UpdateOne(entity).
			Set(entity).
			SaveE(ctx, dao.ResourceDependenciesEdgeSave)

		return err
	})
	if err != nil {
		return errorx.Wrap(err, "error updating resource")
	}

	return apply(ctx, mc, kubeConfig, entity)
}

func apply(ctx context.Context, mc model.ClientSet, kubeConfig *rest.Config, entity *model.Resource) error {
	dp, err := getDeployer(ctx, kubeConfig)
	if err != nil {
		return err
	}
	// Apply resource.
	applyOpts := pkgresource.Options{
		Deployer: dp,
	}

	ready, err := pkgresource.CheckDependencyStatus(ctx, mc, dp, entity)
	if err != nil {
		return errorx.Wrap(err, "error checking dependency status")
	}

	if ready {
		return pkgresource.Apply(
			ctx,
			mc,
			entity,
			applyOpts)
	}

	return nil
}

func getDeployer(ctx context.Context, kubeConfig *rest.Config) (deptypes.Deployer, error) {
	dep, err := deployer.Get(ctx, deptypes.CreateOptions{
		Type:       deployertf.DeployerType,
		KubeConfig: kubeConfig,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to get deployer")
	}

	return dep, nil
}
