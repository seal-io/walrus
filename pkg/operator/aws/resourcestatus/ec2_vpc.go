package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func getEc2Vpc(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := ec2Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		VpcIds: []string{name},
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Vpcs) == 0 {
		return nil, errNotFound
	}

	st := vpcStatusConverter.Convert(string(resp.Vpcs[0].State), "")

	return st, nil
}
