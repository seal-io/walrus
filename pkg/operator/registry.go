package operator

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/pkg/operatorany"
	"github.com/seal-io/seal/pkg/platformk8s"
)

var opCreators map[types.Type]types.Creator

func init() {
	// Register operator creators as below.
	opCreators = map[types.Type]types.Creator{
		platformk8s.OperatorType: platformk8s.NewOperator,
	}
}

// Get returns types.Operator with the given types.CreateOptions.
func Get(ctx context.Context, opts types.CreateOptions) (op types.Operator, err error) {
	f, exist := opCreators[opts.Connector.Type]
	if !exist {
		// Try to create an any operator.
		op, err = operatorany.NewOperator(ctx, opts)
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
