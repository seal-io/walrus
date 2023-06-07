package alibaba

import (
	"context"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator/alibaba/key"
	"github.com/seal-io/seal/pkg/operator/alibaba/resourceexec"
	"github.com/seal-io/seal/pkg/operator/alibaba/resourcelog"
	"github.com/seal-io/seal/pkg/operator/alibaba/resourcestatus"
	optypes "github.com/seal-io/seal/pkg/operator/types"
)

const OperatorType = "Alibaba"

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

func (o Operator) IsConnected(_ context.Context) error {
	client, err := ecs.NewClientWithAccessKey(
		o.cred.Region,
		o.cred.AccessKey,
		o.cred.AccessSecret,
	)
	if err != nil {
		return fmt.Errorf("error create alibaba client %s: %w", o.name, err)
	}

	// Use DescribeRegion API to check reachable and user has access to region.
	// https://www.alibabacloud.com/help/en/elastic-compute-service/latest/regions-describeregions
	req := ecs.CreateDescribeRegionsRequest()
	req.Scheme = "HTTPS"

	_, err = client.DescribeRegions(req)
	if err != nil {
		return fmt.Errorf("error connect to %s: %w", o.name, err)
	}

	return nil
}

func (o Operator) GetStatus(_ context.Context, resource *model.ApplicationResource) (*status.Status, error) {
	st := &status.Status{}

	if !resourcestatus.IsSupported(resource.Type) {
		return st, nil
	}

	nst, err := resourcestatus.Get(*o.cred, resource.Type, resource.Name)
	if err != nil {
		return st, fmt.Errorf("error get resource %s:%s from %s: %w", resource.Type, resource.Name, o.name, err)
	}

	return nst, nil
}

func (o Operator) GetKeys(ctx context.Context, resource *model.ApplicationResource) (*optypes.Keys, error) {
	var (
		subCtx  = context.WithValue(ctx, optypes.CredentialKey, o.cred)
		keyName = key.Encode(resource.Type, resource.Name)
	)

	executable, err := resourceexec.Supported(subCtx, keyName)
	if err != nil {
		return nil, err
	}

	loggable, err := resourcelog.Supported(subCtx, keyName)
	if err != nil {
		return nil, err
	}

	k := optypes.Key{
		Name:       keyName,
		Executable: &executable,
		Loggable:   &loggable,
	}

	return &optypes.Keys{
		Keys:   []optypes.Key{k},
		Labels: []string{"Resource"},
	}, nil
}

func (o Operator) Exec(ctx context.Context, s string, options optypes.ExecOptions) (err error) {
	subCtx := context.WithValue(ctx, optypes.CredentialKey, o.cred)
	err = resourceexec.Exec(subCtx, s, options)

	return err
}

func (o Operator) Log(ctx context.Context, s string, options optypes.LogOptions) error {
	subCtx := context.WithValue(ctx, optypes.CredentialKey, o.cred)
	return resourcelog.Log(subCtx, s, options)
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

func (o Operator) Label(ctx context.Context, resource *model.ApplicationResource, m map[string]string) error {
	return nil
}
