package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/azure/resourcestatus"
	aztypes "github.com/seal-io/walrus/pkg/operator/azure/types"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/hash"
)

const OperatorType = "Azure"

func New(ctx context.Context, opts optypes.CreateOptions) (optypes.Operator, error) {
	name := opts.Connector.ID.String()

	cred, err := aztypes.GetCredential(opts.Connector.ConfigData)
	if err != nil {
		return nil, err
	}

	return Operator{
		name:       name,
		cred:       cred,
		identifier: hash.SumStrings("azure:", cred.SubscriptionID, cred.TenantID, cred.ClientID),
	}, nil
}

type Operator struct {
	name       string
	cred       *aztypes.Credential
	identifier string
}

func (o Operator) Type() optypes.Type {
	return OperatorType
}

func (o Operator) IsConnected(ctx context.Context) error {
	cred, err := azidentity.NewClientSecretCredential(o.cred.TenantID, o.cred.ClientID, o.cred.ClientSecret, nil)
	if err != nil {
		return err
	}

	clientFactory, err := armresources.NewClientFactory(o.cred.SubscriptionID, cred, nil)
	if err != nil {
		return err
	}

	client := clientFactory.NewResourceGroupsClient()

	pager := client.NewListPager(nil)

	_, err = pager.NextPage(ctx)
	if err != nil {
		return fmt.Errorf("error connect to azure: %w", err)
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
	resource *model.ResourceComponent,
) (*types.ResourceComponentOperationKeys, error) {
	return nil, nil
}

func (o Operator) GetStatus(ctx context.Context, component *model.ResourceComponent) (*status.Status, error) {
	st := &status.Status{}
	if !resourcestatus.IsSupported(component.Type) {
		return st, nil
	}

	newCtx := context.WithValue(ctx, optypes.CredentialKey, o.cred)

	nst, err := resourcestatus.Get(newCtx, component.Type, component.Name)
	if err != nil {
		return st, fmt.Errorf("error get resource %s:%s from %s: %w", component.Type, component.Name, o.name, err)
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
