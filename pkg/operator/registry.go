package operator

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/operator/alibaba"
	"github.com/seal-io/walrus/pkg/operator/aws"
	"github.com/seal-io/walrus/pkg/operator/k8s"
	"github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/pkg/operator/unknown"
	"github.com/seal-io/walrus/pkg/operator/unreachable"
)

var opCreators map[types.Type]types.Creator

func init() {
	// Register operator creators as below.
	opCreators = map[types.Type]types.Creator{
		k8s.OperatorType:     k8s.New,
		aws.OperatorType:     aws.New,
		alibaba.OperatorType: alibaba.New,
	}
}

// Get returns types.Operator with the given types.CreateOptions.
func Get(ctx context.Context, opts types.CreateOptions) (op types.Operator, err error) {
	f, exist := opCreators[opts.Connector.Type]
	if !exist {
		// Try to create an any operator.
		op, err = unknown.New(ctx, opts)
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

func UnReachable() types.Operator {
	return unreachable.Operator{}
}
