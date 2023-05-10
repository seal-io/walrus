package distributor

import (
	"context"
	"time"

	"github.com/seal-io/seal/pkg/apis/cost/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Distributor support use input condition to query resources costs,
// and split the shared costs base on the condition defined,
// while the condition with the step(day, month etc.),
// will return the resource cost for each time bucket, like namespace cost per day,
// without step will return the resource total cost within the queried time range.
type Distributor struct {
	ad accumulateDistributor
	sd stepDistributor
}

func New(client model.ClientSet) *Distributor {
	return &Distributor{
		ad: accumulateDistributor{
			client: client,
		},
		sd: stepDistributor{
			client: client,
		},
	}
}

func (d *Distributor) Distribute(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) ([]view.Resource, int, error) {
	cond = wrappedCondition(cond)
	switch {
	case cond.Step != "":
		return d.sd.distribute(ctx, startTime, endTime, cond)
	default:
		return d.ad.distribute(ctx, startTime, endTime, cond)
	}
}
