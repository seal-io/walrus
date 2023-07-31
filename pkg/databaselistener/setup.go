package databaselistener

import (
	"context"

	"github.com/seal-io/seal/pkg/databaselistener/setting"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/database"
)

type SetupOptions struct {
	ModelClient       model.ClientSet
	DataSourceAddress string
}

func Setup(ctx context.Context, opts SetupOptions) (err error) {
	// Setup listen channels.
	database.ListenChannel(setting.ChannelName, setting.Handler)

	// Start the listener.
	return database.StartListener(ctx, opts.DataSourceAddress, opts.ModelClient)
}
