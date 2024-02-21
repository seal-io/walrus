package resourcestatus

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"

	aztypes "github.com/seal-io/walrus/pkg/operator/azure/types"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func virtualMachineClient(ctx context.Context) (*armcompute.VirtualMachinesClient, error) {
	cred, ok := ctx.Value(types.CredentialKey).(*aztypes.Credential)
	if !ok {
		return nil, errors.New("not found credential from context")
	}

	sc, err := azidentity.NewClientSecretCredential(cred.TenantID, cred.ClientID, cred.ClientSecret, nil)
	if err != nil {
		return nil, err
	}

	client, err := armcompute.NewVirtualMachinesClient(cred.SubscriptionID, sc, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
