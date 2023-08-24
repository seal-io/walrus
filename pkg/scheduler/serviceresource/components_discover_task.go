package serviceresource

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresource"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/operator"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/pkg/serviceresources"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

type ComponentsDiscoverTask struct {
	logger      log.Logger
	modelClient model.ClientSet
}

func NewComponentsDiscoverTask(logger log.Logger, mc model.ClientSet) (in *ComponentsDiscoverTask, err error) {
	in = &ComponentsDiscoverTask{
		logger:      logger,
		modelClient: mc,
	}

	return
}

func (in *ComponentsDiscoverTask) Process(ctx context.Context, args ...any) error {
	// NB(thxCode): connectors are usually less in number,
	// in case of reuse the connection built from a connector,
	// we can treat each connector as a task group,
	// group 100 resources of each connector into one task unit,
	// and then process sub resources syncing in task unit.
	cs, err := listCandidateConnectors(ctx, in.modelClient)
	if err != nil {
		return fmt.Errorf("cannot list all connectors: %w", err)
	}

	if len(cs) == 0 {
		return nil
	}

	wg := gopool.Group()

	for i := range cs {
		st := in.buildSyncTasks(ctx, cs[i])
		wg.Go(st)
	}

	return wg.Wait()
}

func (in *ComponentsDiscoverTask) buildSyncTasks(ctx context.Context, c *model.Connector) func() error {
	return func() error {
		op, err := operator.Get(ctx, optypes.CreateOptions{
			Connector: *c,
		})
		if err != nil {
			return err
		}

		if err = op.IsConnected(ctx); err != nil {
			// Warn out without breaking the whole syncing.
			in.logger.Warnf("unreachable connector %q", c.ID)
			// NB(thxCode): replace disconnected connector with unknown connector.
			op = operator.UnReachable()
		}

		cnt, err := c.QueryResources().
			Where(
				serviceresource.ModeNEQ(types.ServiceResourceModeDiscovered),
				serviceresource.Shape(types.ServiceResourceShapeInstance),
			).
			Count(ctx)
		if err != nil {
			return fmt.Errorf("cannot count not discovered resources of connector %q: %w", c.ID, err)
		}

		if cnt == 0 {
			return nil
		}

		const bks = 100

		bkc := cnt/bks + 1
		if bkc == 1 {
			at := in.buildSyncTask(ctx, op, c.ID, 0, bks)
			return at()
		}

		wg := gopool.Group()

		for bk := 0; bk < bkc; bk++ {
			at := in.buildSyncTask(ctx, op, c.ID, bk*bks, bks)
			wg.Go(at)
		}

		return wg.Wait()
	}
}

func (in *ComponentsDiscoverTask) buildSyncTask(
	ctx context.Context,
	op optypes.Operator,
	connectorID object.ID,
	offset,
	limit int,
) func() error {
	return func() error {
		rs, err := in.modelClient.ServiceResources().Query().
			Where(
				serviceresource.ConnectorID(connectorID),
				serviceresource.Shape(types.ServiceResourceShapeInstance),
				serviceresource.Mode(types.ServiceResourceModeManaged)).
			Order(model.Desc(serviceresource.FieldCreateTime)).
			Unique(false).
			Offset(offset).
			Limit(limit).
			Select(
				serviceresource.FieldShape,
				serviceresource.FieldMode,
				serviceresource.FieldID,
				serviceresource.FieldDeployerType,
				serviceresource.FieldType,
				serviceresource.FieldName,
				serviceresource.FieldProjectID,
				serviceresource.FieldEnvironmentID,
				serviceresource.FieldServiceID,
				serviceresource.FieldConnectorID).
			All(ctx)
		if err != nil {
			return fmt.Errorf("error listing service resources: %w", err)
		}

		_, err = serviceresources.Discover(ctx, op, in.modelClient, rs)

		return err
	}
}
