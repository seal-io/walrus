package platformk8s

import (
	"context"
	"time"

	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/k8s"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/utils/log"
)

const OperatorType = types.ConnectorTypeK8s

// NewOperator returns operator.Operator with the given options.
func NewOperator(ctx context.Context, opts operator.CreateOptions) (operator.Operator, error) {
	// NB(thxCode): disable timeout as we don't know the maximum time-cost of once operation,
	// and rely on the session context timeout control of each operation.
	restConfig, err := GetConfig(opts.Connector, WithoutTimeout())
	if err != nil {
		return nil, err
	}
	op := Operator{
		Logger:     log.WithName("operator").WithName("k8s"),
		RestConfig: restConfig,
	}
	return op, nil
}

type Operator struct {
	Logger     log.Logger
	RestConfig *rest.Config
}

// Type implements operator.Operator.
func (Operator) Type() operator.Type {
	return OperatorType
}

// IsConnected implements operator.Operator.
func (op Operator) IsConnected(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return k8s.Wait(ctx, op.RestConfig)
}
