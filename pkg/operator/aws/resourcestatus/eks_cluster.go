package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/seal-io/walrus/pkg/dao/types/status"
)

func getEksCluster(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := eksClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if resp.Cluster == nil {
		return &status.Status{}, nil
	}

	st := eksClusterStatusConverter.Convert(string(resp.Cluster.Status), "")

	return st, nil
}
