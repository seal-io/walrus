package resourcestatus

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/go-autorest/autorest/azure"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	aztypes "github.com/seal-io/walrus/pkg/operator/azure/types"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func getVirtualNetwork(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cred, ok := ctx.Value(types.CredentialKey).(*aztypes.Credential)
	if !ok {
		return nil, errors.New("not found credential from context")
	}

	sc, err := azidentity.NewClientSecretCredential(cred.TenantID, cred.ClientID, cred.ClientSecret, nil)
	if err != nil {
		return nil, err
	}

	client, err := armnetwork.NewVirtualNetworksClient(cred.SubscriptionID, sc, nil)
	if err != nil {
		return nil, err
	}

	r, err := azure.ParseResourceID(name)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, r.ResourceGroup, r.ResourceName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get azure resource %s %s: %w", resourceType, r.ResourceName, err)
	}

	return virtualNetworkStatusConverter.Convert(string(*resp.Properties.ProvisioningState), ""), nil
}
