package google

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/google/resourcestatus"
	gtypes "github.com/seal-io/walrus/pkg/operator/google/types"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/hash"
)

const OperatorType = "Google"

func New(ctx context.Context, opts optypes.CreateOptions) (optypes.Operator, error) {
	name := opts.Connector.ID.String()

	cred, err := gtypes.GetCredential(opts.Connector.ConfigData)
	if err != nil {
		return nil, err
	}

	return Operator{
		name:       name,
		cred:       cred,
		identifier: hash.SumStrings("google:", cred.Project, cred.Region, cred.Zone),
	}, nil
}

type Operator struct {
	name       string
	cred       *gtypes.Credential
	identifier string
}

func (o Operator) Type() optypes.Type {
	return OperatorType
}

func (o Operator) IsConnected(ctx context.Context) error {
	service, err := compute.NewService(ctx, option.WithCredentialsJSON([]byte(o.cred.Credentials)))
	if err != nil {
		return err
	}

	_, err = service.Regions.List(o.cred.Project).Do()

	if err != nil {
		return fmt.Errorf("error connect to google cloud: %w", err)
	}

	return nil
}

func (o Operator) Burst() int {
	return 100
}

func (o Operator) ID() string {
	return o.identifier
}

func (o Operator) GetKeys(
	ctx context.Context,
	component *model.ResourceComponent,
) (*types.ResourceComponentOperationKeys, error) {
	return nil, nil
}

func (o Operator) GetStatus(ctx context.Context, resource *model.ResourceComponent) (*status.Status, error) {
	st := &status.Status{}
	if !resourcestatus.IsSupported(resource.Type) {
		return st, nil
	}

	newCtx := context.WithValue(ctx, optypes.CredentialKey, o.cred)

	nst, err := resourcestatus.Get(newCtx, resource.Type, resource.Name)
	if err != nil {
		return st, fmt.Errorf("error get resource %s:%s from %s: %w", resource.Type, resource.Name, o.name, err)
	}

	return nst, nil
}

func (o Operator) GetComponents(
	ctx context.Context,
	component *model.ResourceComponent,
) ([]*model.ResourceComponent, error) {
	return nil, nil
}

func (o Operator) Log(ctx context.Context, s string, options optypes.LogOptions) error {
	return nil
}

func (o Operator) Exec(ctx context.Context, s string, options optypes.ExecOptions) error {
	return nil
}

func (o Operator) Label(ctx context.Context, component *model.ResourceComponent, m map[string]string) error {
	return nil
}
