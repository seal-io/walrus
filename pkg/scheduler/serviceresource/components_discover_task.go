package serviceresource

import (
	"context"
	"fmt"

	"go.uber.org/multierr"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresource"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
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

// Process implements the Task interface,
// it will discover the components of managed instance resources.
//
// Process fetches all connectors at first,
// and constructs the operator related to each connector.
// Then it will query the resources belong to each connector,
// and process the resources in batches concurrently according to the operator burst size.
func (in *ComponentsDiscoverTask) Process(ctx context.Context, args ...any) error {
	// Retrieve operators.
	opIndexer, opLimiter, err := retrieveOperators(ctx, in.modelClient, in.logger)
	if err != nil || len(opIndexer) == 0 {
		return err
	}

	// Merge the errors to return them all at once,
	// instead of returning the first error.
	var berr error

	wg := gopool.Group()

	for cid := range opIndexer {
		// Count the total size of resources belong to the connector,
		// skip if no resources or error raising.
		total, err := in.modelClient.ServiceResources().Query().
			Where(
				serviceresource.ConnectorID(cid),
				serviceresource.Shape(types.ServiceResourceShapeInstance),
				serviceresource.Mode(types.ServiceResourceModeManaged)).
			Count(ctx)
		if multierr.AppendInto(&berr, err) || total == 0 {
			continue
		}

		op := opIndexer[cid]

		// Divide the resources in multiple batches according to the operator burst size.
		bs, bc := getBatches(total, op.Burst(), 10)
		for b := 0; b < bc; b++ {
			// Process the resources in batches concurrently.
			p := in.buildProcess(ctx, cid, op, opLimiter, b*bs, bs)
			wg.Go(p)
		}
	}

	return multierr.Append(berr, wg.Wait())
}

func (in *ComponentsDiscoverTask) buildProcess(
	ctx context.Context,
	cid object.ID,
	op optypes.Operator,
	opLimiter operatorLimiter,
	offset, limit int,
) func() error {
	return func() error {
		// Controls the concurrency of operators with the same ID,
		// avoids server instability or throttling caused by creating too many client connections.
		opLimiter.Acquire(op.ID())
		defer opLimiter.Release(op.ID())

		rs, err := in.modelClient.ServiceResources().Query().
			Where(
				serviceresource.ConnectorID(cid),
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
			return fmt.Errorf("error listing service resources with offset %d, limit %d: %w",
				offset, limit, err)
		}

		_, err = serviceresources.Discover(ctx, op, in.modelClient, rs)

		return err
	}
}
