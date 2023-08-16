package types

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
)

// CreateOptions holds the options for creating Operator.
type CreateOptions struct {
	// Connector indicates the model.Connector for creating.
	Connector model.Connector
}

// Creator is a factory func to create Operator.
type Creator func(context.Context, CreateOptions) (Operator, error)
