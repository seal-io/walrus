package resource

import (
	"context"
	"time"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/strs"
)

// UpdateServiceDriftResult update service and its drift result.
func UpdateServiceDriftResult(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Resource,
	revision *model.ResourceRevision,
) error {
	resourceDriftResult := &types.ResourceDriftDetection{}

	if revision.Type == types.ResourceRevisionTypeDetect {
		resourceDrift, err := parser.ParseDriftOutput(revision.Drift)
		if err != nil {
			return err
		}

		if resourceDrift != nil && len(resourceDrift.ResourceComponentDrifts) > 0 {
			rds := make(map[string]*types.ResourceComponentDrift, len(resourceDrift.ResourceComponentDrifts))
			for _, rd := range resourceDrift.ResourceComponentDrifts {
				rds[strs.Join("/", rd.Type, rd.Name)] = rd
			}

			if err = updateResourceDriftResult(ctx, mc, entity.ID, rds); err != nil {
				return err
			}

			resourceDriftResult.Drifted = true
			resourceDriftResult.Time = time.Now()
			resourceDriftResult.Result = resourceDrift
		}
	}

	err := mc.Resources().UpdateOne(entity).
		SetDriftDetection(resourceDriftResult).
		Exec(ctx)
	if err != nil {
		return err
	}

	if resourceDriftResult.Drifted {
		return nil
	}

	return mc.ResourceComponents().Update().
		Where(resourcecomponent.ResourceID(entity.ID)).
		ClearDriftDetection().
		Exec(ctx)
}

// updateResourceDriftResult update the drift detection result of the service's resources.
func updateResourceDriftResult(
	ctx context.Context,
	mc model.ClientSet,
	resourceID object.ID,
	resourceDrifts map[string]*types.ResourceComponentDrift,
) error {
	resources, err := mc.ResourceComponents().Query().
		Where(resourcecomponent.ResourceID(resourceID)).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	for i := range resources {
		r := resources[i]

		key := strs.Join("/", r.Type, r.Name)
		if _, ok := resourceDrifts[key]; !ok {
			continue
		}

		err = mc.ResourceComponents().UpdateOne(r).
			SetDriftDetection(&types.ResourceComponentDriftDetection{
				Drifted: true,
				Time:    time.Now(),
				Result:  resourceDrifts[key],
			}).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
