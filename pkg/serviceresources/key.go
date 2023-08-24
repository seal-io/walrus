package serviceresources

import (
	"context"

	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/operator"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

// SetKeys sets the keys of the resources for operations like log and exec.
//
// The given model.ServiceResource item must specify the following fields:
// Shape, Mode, ID, DeployerType, Type, Name and ConnectorID.
func SetKeys(
	ctx context.Context,
	logger log.Logger,
	modelClient model.ClientSet,
	candidates []*model.ServiceResource,
	operators map[object.ID]optypes.Operator,
) {
	// Index the resource and its components via operator ID.
	connRess := make(map[object.ID][]*model.ServiceResource)
	{
		// Calculate the capacity.
		caps := make(map[object.ID]int)

		for _, c := range candidates {
			if c.ConnectorID == "" {
				continue
			}

			if c.Shape == types.ServiceResourceShapeClass {
				for _, cis := range c.Edges.Instances {
					caps[cis.ConnectorID] += 1 + len(cis.Edges.Components)
				}

				continue
			}

			caps[c.ConnectorID] += 1 + len(c.Edges.Components)
		}

		for id := range caps {
			connRess[id] = make([]*model.ServiceResource, 0, caps[id])
		}

		// Index the resources via connector ID.
		for i := range candidates {
			c := candidates[i]
			if c.ConnectorID == "" {
				continue
			}

			if c.Shape == types.ServiceResourceShapeClass {
				for j := range c.Edges.Instances {
					cis := c.Edges.Instances[j]
					connRess[cis.ConnectorID] = append(connRess[cis.ConnectorID], cis)
					connRess[cis.ConnectorID] = append(connRess[cis.ConnectorID],
						cis.Edges.Components...)
				}

				continue
			}

			connRess[c.ConnectorID] = append(connRess[c.ConnectorID], c)
			connRess[c.ConnectorID] = append(connRess[c.ConnectorID],
				c.Edges.Components...)
		}
	}

	// Construct the operator via connector.
	if operators == nil {
		operators = make(map[object.ID]optypes.Operator)
	}
	{
		cs, err := modelClient.Connectors().Query().
			Where(
				connector.CategoryNEQ(types.ConnectorCategoryCustom),
				connector.CategoryNEQ(types.ConnectorCategoryVersionControl),
				connector.IDIn(sets.KeySet(connRess).UnsortedList()...)).
			Select(
				connector.FieldID,
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData).
			All(ctx)
		if err != nil {
			logger.Errorf("cannot list connectors: %v", err)
			return
		}

		for i := range cs {
			if _, ok := operators[cs[i].ID]; ok {
				continue
			}

			var op optypes.Operator

			op, err = operator.Get(ctx, optypes.CreateOptions{
				Connector: *cs[i],
			})
			if err != nil {
				// Warn out without breaking the whole fetching.
				logger.Warnf("cannot get operator of connector %q: %v", cs[i].ID, err)
				continue
			}

			if err = op.IsConnected(ctx); err != nil {
				// Warn out without breaking the whole syncing.
				logger.Warnf("unreachable connector %q", cs[i].ID)
				// Replace disconnected connector with unknown connector.
				op = operator.UnReachable()
			}

			operators[cs[i].ID] = op
		}
	}

	// Index resources via operator.
	opRess := make(map[optypes.Operator][]*model.ServiceResource)
	{
		for id := range connRess {
			op, ok := operators[id]
			if !ok {
				continue
			}

			opRess[op] = connRess[id]
		}
	}

	if len(opRess) == 0 {
		return
	}

	// Get the keys of resources in parallel.
	wg := gopool.GroupWithContextIn(ctx)

	for op_ := range opRess {
		op := op_
		bks := op.Burst()
		ress := opRess[op]

		for i := 0; i < bks; i++ {
			s := i

			wg.Go(func(ctx context.Context) error {
				// Merge the errors to return them all at once,
				// instead of returning the first error.
				var berr error

				for j := s; j < len(ress); j += bks {
					if ress[j].Shape != types.ServiceResourceShapeInstance ||
						ress[j].Mode == types.ServiceResourceModeData {
						continue
					}

					var err error

					ress[j].Keys, err = op.GetKeys(ctx, ress[j])
					berr = multierr.Append(berr, err)
				}

				return berr
			})
		}
	}

	if err := wg.Wait(); err != nil {
		logger.Errorf("error getting keys of resources: %v", err)
	}
}
