package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func getEc2Instance(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := ec2Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeInstances(
		ctx,
		&ec2.DescribeInstancesInput{
			InstanceIds: []string{name},
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return nil, errNotFound
	}

	st := ec2InstanceStatusConverter.Convert(string(resp.Reservations[0].Instances[0].State.Name), "")

	return st, nil
}
