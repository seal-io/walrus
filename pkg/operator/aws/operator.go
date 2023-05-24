package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator/aws/resourcestatus"
	optypes "github.com/seal-io/seal/pkg/operator/types"
)

const OperatorType = "AWS"

// New returns operator.Operator with the given options.
func New(ctx context.Context, opts optypes.CreateOptions) (optypes.Operator, error) {
	name := opts.Connector.ID.String()

	cred, err := optypes.GetCredential(opts.Connector.ConfigData)
	if err != nil {
		return nil, err
	}

	return Operator{
		name: name,
		cred: cred,
	}, nil
}

type Operator struct {
	name string
	cred *optypes.Credential
}

func (o Operator) Type() optypes.Type {
	return OperatorType
}

func (o Operator) IsConnected(ctx context.Context) error {
	cred := resourcestatus.Credential(*o.cred)

	cfg, err := cred.Config()
	if err != nil {
		return err
	}

	// Use DescribeRegions API to check reachable.
	cli := ec2.NewFromConfig(*cfg)

	_, err = cli.DescribeRegions(ctx, nil)
	if err != nil {
		return fmt.Errorf("error connect to aws: %w", err)
	}

	return nil
}

func (o Operator) GetStatus(ctx context.Context, resource *model.ApplicationResource) (*status.Status, error) {
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

func (o Operator) GetKeys(ctx context.Context, resource *model.ApplicationResource) (*optypes.Keys, error) {
	return nil, nil
}

func (o Operator) GetEndpoints(
	ctx context.Context,
	resource *model.ApplicationResource,
) ([]types.ApplicationResourceEndpoint, error) {
	return nil, nil
}

func (o Operator) GetComponents(
	ctx context.Context,
	resource *model.ApplicationResource,
) ([]*model.ApplicationResource, error) {
	return nil, nil
}

func (o Operator) Log(ctx context.Context, s string, options optypes.LogOptions) error {
	return errors.New("cannot log")
}

func (o Operator) Exec(ctx context.Context, s string, options optypes.ExecOptions) error {
	return errors.New("cannot execute")
}

func (o Operator) Label(ctx context.Context, resource *model.ApplicationResource, m map[string]string) error {
	return nil
}
