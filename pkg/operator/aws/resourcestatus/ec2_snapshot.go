package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func getEc2Snapshot(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := ec2Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeSnapshots(
		ctx,
		&ec2.DescribeSnapshotsInput{
			SnapshotIds: []string{name},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Snapshots) == 0 {
		return nil, errNotFound
	}

	var (
		msg      string
		snapshot = resp.Snapshots[0]
	)

	if snapshot.StateMessage != nil {
		msg = *snapshot.StateMessage
	}

	st := ec2SnapshotStatusConverter.Convert(string(snapshot.State), msg)

	return st, nil
}
