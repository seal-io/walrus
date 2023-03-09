package module

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/modules"
	"github.com/seal-io/seal/utils/bus"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

// BusMessage wraps the changed model.Module as a bus.Message.
type BusMessage struct {
	ModelClient model.ClientSet
	Refer       *model.Module
}

// Notify notifies the changed model.Module.
func Notify(ctx context.Context, mc model.ClientSet, refer *model.Module) error {
	return bus.Publish(ctx, BusMessage{ModelClient: mc, Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from model.Module.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}

// SyncSchema fetches a remote module and updates the module schema in the background.
func SyncSchema(ctx context.Context, message BusMessage) error {
	gopool.Go(func() {
		if err := syncSchema(ctx, message); err != nil {
			module := message.Refer
			module.Status = status.Error
			module.StatusMessage = fmt.Sprintf("sync schema failed: %v", err)
			update, updateErr := dao.ModuleUpdate(message.ModelClient, module)
			if updateErr != nil {
				log.Errorf("failed to prepare module update: %v", updateErr)
				return
			}
			if updateErr = update.Exec(ctx); updateErr != nil {
				log.Errorf("failed to update module %s: %v", module.ID, updateErr)
			}
		}
	})
	return nil
}

func syncSchema(ctx context.Context, message BusMessage) error {
	module := message.Refer

	if module.Schema != nil {
		// Short circuit when the schema is presented. To refresh the schema, set it to nil first.
		return nil
	}

	log.Debugf("syncing schema for module %s", message.Refer.ID)

	moduleSchema, err := modules.LoadTerraformModuleSchema(module.Source)
	if err != nil {
		return err
	}

	module.Schema = moduleSchema
	module.Status = status.Ready

	update, err := dao.ModuleUpdate(message.ModelClient, module)
	if err != nil {
		return err
	}

	return update.Exec(ctx)
}
