package resourcestatus

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"

	"github.com/seal-io/walrus/pkg/dao/types/status"
)

func getVirtualMachine(ctx context.Context, resourceType, name string) (*status.Status, error) {
	client, err := virtualMachineClient(ctx)
	if err != nil {
		return &status.Status{}, err
	}

	r, err := azure.ParseResourceID(name)
	if err != nil {
		return nil, err
	}

	resp, err := client.InstanceView(ctx, r.ResourceGroup, r.ResourceName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get azure resource %s %s: %w", resourceType, r.ResourceName, err)
	}

	if resp.Statuses == nil || len(resp.Statuses) == 0 {
		return nil, errors.New("not found")
	}

	for _, vmStatus := range resp.Statuses {
		if strings.HasPrefix(*vmStatus.Code, "PowerState/") {
			state := strings.TrimPrefix(*vmStatus.Code, "PowerState/")

			return virtualMachineStatusConverter.Convert(state, *vmStatus.Message), nil
		}
	}

	return nil, errors.New("not found")
}
