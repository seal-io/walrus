package modules

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/bus"
	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

// BusMessage wraps the changed model.Module as a bus.Message.
type BusMessage struct {
	ModelClient model.ClientSet
	Refer       *model.Module
}

// Notify notifies the changed model.Setting.
func Notify(ctx context.Context, mc model.ClientSet, refer *model.Module) error {
	return bus.Publish(ctx, BusMessage{ModelClient: mc, Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from proxy model.Setting.
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
			if updateErr := message.ModelClient.Modules().UpdateOne(module).Exec(ctx); updateErr != nil {
				log.Errorf("failed to update module %s: %v", module.ID, updateErr)
			}
		}
	})
	return nil
}

func syncSchema(ctx context.Context, message BusMessage) error {
	module, err := message.ModelClient.Modules().Get(ctx, message.Refer.ID)
	if err != nil {
		return err
	}

	if module.Schema != nil {
		// Short circuit when the schema is presented. To refresh the schema, set it to nil first.
		return nil
	}

	log.Debugf("syncing schema for module %s", message.Refer.ID)

	tmpDir := files.TempDir("seal-module-*")
	defer os.RemoveAll(tmpDir)

	if err := getter.Get(tmpDir, message.Refer.Source); err != nil {
		return err
	}
	mod, _ := tfconfig.LoadModule(tmpDir)

	readme, err := getReadme(tmpDir)
	if err != nil {
		return err
	}

	moduleSchema := &types.ModuleSchema{
		Readme: readme,
	}

	for _, v := range mod.Variables {
		// TODO more custom schema fields by parsing tf comments
		moduleSchema.Variables = append(moduleSchema.Variables, types.ModuleVariable{
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
			Default:     v.Default,
			Required:    v.Required,
			Sensitive:   v.Sensitive,
		})
	}

	for _, v := range mod.Outputs {
		moduleSchema.Outputs = append(moduleSchema.Outputs, types.ModuleOutput{
			Name:        v.Name,
			Description: v.Description,
			Sensitive:   v.Sensitive,
		})
	}

	for name := range mod.RequiredProviders {
		moduleSchema.RequiredConnectorTypes = append(moduleSchema.RequiredConnectorTypes, name)
	}

	module.Schema = moduleSchema
	module.Status = status.Ready

	update, err := dao.ModuleUpdate(message.ModelClient, module)
	if err != nil {
		return err
	}

	return update.Exec(ctx)
}

func getReadme(dir string) (string, error) {
	path := filepath.Join(dir, "README.md")
	if files.Exists(path) {
		content, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}

	return "", nil
}
