package operator

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/operator/any"
	"github.com/seal-io/seal/pkg/operator/k8s"
	"github.com/seal-io/seal/pkg/operator/types"
)

var opCreators map[types.Type]types.Creator

func init() {
	// Register operator creators as below.
	opCreators = map[types.Type]types.Creator{
		k8s.OperatorType: k8s.NewOperator,
	}
}

// Get returns types.Operator with the given types.CreateOptions.
func Get(ctx context.Context, opts types.CreateOptions) (op types.Operator, err error) {
	f, exist := opCreators[opts.Connector.Type]
	if !exist {
		// Try to create an any operator.
		op, err = any.NewOperator(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("unknown operator: %s", opts.Connector.Type)
		}
	} else {
		op, err = f(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("error connecting %s operator: %w", opts.Connector.Type, err)
		}
	}

	return op, nil
}
