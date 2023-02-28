package connectors

import (
	"context"
	"fmt"
	"time"

	"github.com/seal-io/seal/pkg/costs/deployer"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/bus"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

const (
	timeout = 3 * time.Minute
)

// BusMessage wraps the changed model.Connector as a bus.Message.
type BusMessage struct {
	ModelClient model.ClientSet
	Refer       *model.Connector
	Replace     bool
}

// Notify notifies the changed connector.
func Notify(ctx context.Context, mc model.ClientSet, refer *model.Connector, replace bool) error {
	return bus.Publish(ctx, BusMessage{ModelClient: mc, Refer: refer, Replace: replace})
}

// AddSubscriber add the subscriber to handle the changed notification from connector.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}

// EnsureCostTools install the cost tools.
func EnsureCostTools(_ context.Context, message BusMessage) error {
	if !message.Refer.EnableFinOps {
		return nil
	}

	gopool.Go(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var (
			st  string
			msg string
		)
		err := ensureCostTools(ctx, message)
		if err != nil {
			st = status.Error
			msg = fmt.Sprintf("error ensure cost tools: %v", err)
			log.Errorf("error ensure cost tools for connector %s: %v", message.Refer.Name, err)
		} else {
			st = status.Ready
			msg = ""
		}

		updateErr := message.ModelClient.Connectors().
			UpdateOneID(message.Refer.ID).
			SetFinOpsStatus(st).
			SetFinOpsStatusMessage(msg).
			Exec(ctx)
		if updateErr != nil {
			log.Errorf("failed to update connector %s: %v", message.Refer.Name, updateErr)
		}

	})
	return nil
}

func ensureCostTools(ctx context.Context, message BusMessage) error {
	log.WithName("cost").Debugf("ensuring cost tools for connector %s", message.Refer.Name)

	conn, err := message.ModelClient.Connectors().Get(ctx, message.Refer.ID)
	if err != nil {
		return err
	}

	return deployer.DeployCostTools(ctx, conn, message.Replace)
}
