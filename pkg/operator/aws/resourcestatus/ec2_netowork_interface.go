package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func getEc2NetworkInterface(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := ec2Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []string{name},
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.NetworkInterfaces) == 0 {
		return nil, errNotFound
	}

	st := ec2NetworkInterfaceStatusConverter.Convert(string(resp.NetworkInterfaces[0].Status), "")

	return st, nil
}
