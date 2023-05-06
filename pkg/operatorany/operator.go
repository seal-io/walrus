package operatorany

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform/operator"
)

const OperatorType = "Any"

// NewOperator returns operator.Operator with the given options.
func NewOperator(ctx context.Context, opts operator.CreateOptions) (operator.Operator, error) {
	if opts.Connector.Category != types.ConnectorCategoryCustom {
		return nil, errors.New("not custom connector")
	}
	return Operator{}, nil
}

type Operator struct{}

func (Operator) Type() operator.Type {
	return OperatorType
}

func (Operator) IsConnected(ctx context.Context) error {
	return nil
}

func (Operator) GetKeys(ctx context.Context, resource *model.ApplicationResource) (*operator.Keys, error) {
	return nil, nil
}

func (Operator) GetStatus(ctx context.Context, resource *model.ApplicationResource) (*status.Status, error) {
	return &status.Status{
		Summary: status.Summary{
			SummaryStatus: "Ready",
		},
	}, nil
}

func (Operator) GetEndpoints(ctx context.Context, resource *model.ApplicationResource) ([]types.ApplicationResourceEndpoint, error) {
	return nil, nil
}

func (Operator) GetComponents(ctx context.Context, resource *model.ApplicationResource) ([]*model.ApplicationResource, error) {
	return nil, nil
}

func (Operator) Log(ctx context.Context, s string, options operator.LogOptions) error {
	return errors.New("cannot log")
}

func (Operator) Exec(ctx context.Context, s string, options operator.ExecOptions) error {
	return errors.New("cannot execute")
}

func (Operator) Label(ctx context.Context, resource *model.ApplicationResource, m map[string]string) error {
	return nil
}
