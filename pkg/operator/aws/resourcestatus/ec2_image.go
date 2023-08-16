package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/seal-io/walrus/pkg/dao/types/status"
)

func getEc2Image(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := ec2Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeImages(
		ctx,
		&ec2.DescribeImagesInput{
			ImageIds: []string{name},
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Images) == 0 {
		return nil, errNotFound
	}

	var (
		msg string
		img = resp.Images[0]
	)

	if img.StateReason != nil && img.StateReason.Message != nil {
		msg = *img.StateReason.Message
	}

	st := ec2ImageStatusConverter.Convert(string(img.State), msg)

	return st, nil
}
