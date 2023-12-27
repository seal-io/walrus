package resourcestatus

import (
	"context"

	"github.com/docker/docker/client"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/strs"
)

func getContainerStatus(ctx context.Context, client *client.Client, _, containerID string) (*status.Status, error) {
	status := &status.Status{}

	containerJSON, err := client.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	status.SummaryStatus = strs.Capitalize(containerJSON.State.Status)
	status.Error = containerJSON.State.ExitCode != 0
	status.SummaryStatusMessage = containerJSON.State.Error

	return status, nil
}
