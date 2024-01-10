package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/docker/resourcestatus"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/hash"
)

const OperatorType = "Docker"

// New returns operator.Operator with the given options.
func New(ctx context.Context, opts optypes.CreateOptions) (optypes.Operator, error) {
	name := opts.Connector.ID.String()

	host, _, err := property.GetString(opts.Connector.ConfigData["host"].Value)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return Operator{
		name:       name,
		identifier: hash.SumStrings("docker:", host),
		client:     cli,
	}, nil
}

type Operator struct {
	name       string
	identifier string
	client     *client.Client
}

func (o Operator) Type() optypes.Type {
	return OperatorType
}

func (o Operator) IsConnected(ctx context.Context) error {
	if _, err := o.client.ServerVersion(ctx); err != nil {
		return fmt.Errorf("error connect to docker daemon: %w", err)
	}

	return nil
}

func (o Operator) Burst() int {
	return 100
}

func (o Operator) ID() string {
	return o.identifier
}

func (o Operator) GetStatus(ctx context.Context, resource *model.ResourceComponent) (*status.Status, error) {
	st := &status.Status{}
	if !resourcestatus.IsSupported(resource.Type) {
		return st, nil
	}

	nst, err := resourcestatus.Get(ctx, o.client, resource.Type, resource.Name)
	if err != nil {
		return st, fmt.Errorf("error get resource %s:%s from %s: %w", resource.Type, resource.Name, o.name, err)
	}

	return nst, nil
}

func (o Operator) GetKeys(
	ctx context.Context,
	resource *model.ResourceComponent,
) (*types.ResourceComponentOperationKeys, error) {
	return nil, nil
}

func (o Operator) Exec(ctx context.Context, s string, options optypes.ExecOptions) error {
	return nil
}

func (o Operator) Log(ctx context.Context, s string, options optypes.LogOptions) error {
	return nil
}

func (o Operator) GetComponents(
	ctx context.Context,
	resource *model.ResourceComponent,
) ([]*model.ResourceComponent, error) {
	return nil, nil
}

func (o Operator) Label(ctx context.Context, resource *model.ResourceComponent, m map[string]string) error {
	return nil
}
