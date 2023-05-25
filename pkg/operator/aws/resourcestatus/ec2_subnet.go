package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func getEc2Subnet(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := ec2Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		SubnetIds: []string{name},
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Subnets) == 0 {
		return nil, errNotFound
	}

	st := ec2SubnetStatusConverter.Convert(string(resp.Subnets[0].State), "")

	return st, nil
}
