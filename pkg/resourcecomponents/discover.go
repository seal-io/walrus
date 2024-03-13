package resourcecomponents

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"go.uber.org/multierr"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/strs"
)

// Discover discovers the components of the given model.ResourceComponent with by the given operator.Operator,
// and returns all the discovered model.ResourceComponent items.
//
// The given model.ResourceComponent item must be instance shape and managed mode.
//
// The given model.ResourceComponent item must specify the following fields:
// Shape, Mode, ID, DeployerType, Type, Name, ProjectID, EnvironmentID, ResourceID and ConnectorID.
func Discover(
	ctx context.Context,
	op optypes.Operator,
	modelClient model.ClientSet,
	candidates []*model.ResourceComponent,
) ([]*model.ResourceComponent, error) {
	var ncrs []*model.ResourceComponent

	if op == nil {
		return ncrs, nil
	}

	// Merge the errors to return them all at once,
	// instead of returning the first error.
	var berr error

	for i := range candidates {
		// Give up the loop if the context is canceled.
		if multierr.AppendInto(&berr, ctx.Err()) {
			break
		}

		// Skip the resource if it is not instance shape or not managed mode.
		if candidates[i].Shape != types.ResourceComponentShapeInstance ||
			candidates[i].Mode != types.ResourceComponentModeManaged {
			continue
		}

		// Get observed components from remote.
		observedComps, err := op.GetComponents(ctx, candidates[i])
		if multierr.AppendInto(&berr, err) {
			continue
		}

		if observedComps == nil {
			continue
		}

		// Get record components from local.
		recordComps, err := modelClient.ResourceComponents().Query().
			Where(resourcecomponent.CompositionID(candidates[i].ID)).
			All(ctx)
		if multierr.AppendInto(&berr, err) {
			continue
		}

		// Calculate creating list and deleting list.
		observedCompsIndex := make(map[string]*model.ResourceComponent, len(observedComps))

		for j := range observedComps {
			c := observedComps[j]
			observedCompsIndex[strs.Join("/", c.Type, c.Name)] = c
		}

		deleteCompIDs := make([]object.ID, 0, len(recordComps))

		for _, c := range recordComps {
			k := strs.Join("/", c.Type, c.Name)
			if observedCompsIndex[k] != nil {
				delete(observedCompsIndex, k)
				continue
			}

			deleteCompIDs = append(deleteCompIDs, c.ID)
		}

		createComps := make([]*model.ResourceComponent, 0, len(observedCompsIndex))

		for k := range observedCompsIndex {
			observedCompsIndex[k].Shape = types.ResourceComponentShapeInstance
			createComps = append(createComps, observedCompsIndex[k])
		}

		// Create new components.
		if len(createComps) != 0 {
			createComps, err = modelClient.ResourceComponents().CreateBulk().
				Set(createComps...).
				Save(ctx)
			if multierr.AppendInto(&berr, err) {
				continue
			}

			ncrs = append(ncrs, createComps...)
		}

		// Delete stale components.
		if len(deleteCompIDs) != 0 {
			_, err = modelClient.ResourceComponents().Delete().
				Where(resourcecomponent.IDIn(deleteCompIDs...)).
				Exec(ctx)
			if multierr.AppendInto(&berr, err) {
				continue
			}
		}
	}

	return ncrs, berr
}

// FilterResourceComponentChange filters the given types.ResourceComponentChange items with current resource's model.ResourceComponents.
func FilterResourceComponentChange(
	ctx context.Context,
	mc model.ClientSet,
	resourceID object.ID,
	changes []*types.ResourceComponentChange,
) ([]*types.ResourceComponentChange, error) {
	// Get the resource components of the resource.
	rcs, err := mc.ResourceComponents().Query().
		Where(resourcecomponent.ResourceID(resourceID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Index the resource components by their type and name.
	rcsIndex := make(map[string]*model.ResourceComponent, len(rcs))

	for i := range rcs {
		if rcs[i].IndexKey == "" {
			continue
		}

		rcsIndex[rcs[i].IndexKey] = rcs[i]
	}

	// Filter the changes.
	var fchanges []*types.ResourceComponentChange

	for i := range changes {
		rcc := changes[i]
		if rcc.Type == parser.TerraformTypeData || rcc.Mode == tfjson.DataResourceMode {
			continue
		}

		// Create change is always accepted.
		if rcc.Change.Type == types.ResourceComponentChangeTypeCreate {
			fchanges = append(fchanges, rcc)
			continue
		}

		if rcc.Index != nil {
			switch {
			case parser.IsNumber(rcc.Index):
				rcsKey := strs.Join(".", rcc.ModuleAddress, rcc.Type, rcc.Name) + fmt.Sprintf("[%d]", rcc.Index)
				if rc, rcsOk := rcsIndex[rcsKey]; rcsOk {
					rcc.Name = rc.Name
					fchanges = append(fchanges, rcc)
					continue
				}
			case parser.IsString(rcc.Index):
				rcsKey := strs.Join(".", rcc.ModuleAddress, rcc.Type, rcc.Name) + fmt.Sprintf("[\"%s\"]", rcc.Index)
				if rc, rcsOk := rcsIndex[rcsKey]; rcsOk {
					rcc.Name = rc.Name
					fchanges = append(fchanges, rcc)
					continue
				}
			}
		}

		key := strs.Join(".", rcc.ModuleAddress, rcc.Type, rcc.Name)
		if rc, ok := rcsIndex[key]; ok {
			rcc.Name = rc.Name
			fchanges = append(fchanges, rcc)
		}
	}

	return fchanges, nil
}
