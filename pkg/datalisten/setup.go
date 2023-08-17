package datalisten

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/database"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
)

type StartOptions struct {
	ModelClient       model.ClientSet
	DataSourceAddress string
}

// Start leverages the database listener to establish the data changes.
func Start(ctx context.Context, opts StartOptions) error {
	// Create listener with datasource address.
	l, err := database.NewListener(opts.DataSourceAddress)
	if err != nil {
		return err
	}

	// Register listen handler.
	hs := []database.ListenHandler{
		modelchange.Handle(opts.ModelClient),
	}
	for i := range hs {
		if err = l.Register(hs[i]); err != nil {
			return err
		}
	}

	// Start listen.
	return l.Start(ctx)
}
