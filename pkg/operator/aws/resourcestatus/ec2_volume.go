package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func getEc2Volume(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := ec2Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeVolumes(
		ctx,
		&ec2.DescribeVolumesInput{
			VolumeIds: []string{name},
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Volumes) == 0 {
		return nil, errNotFound
	}

	st := ec2VolumeStatusConverter.Convert(string(resp.Volumes[0].State), "")

	return st, nil
}
