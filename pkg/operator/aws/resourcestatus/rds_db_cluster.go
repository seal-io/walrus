package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"

	"github.com/seal-io/walrus/pkg/dao/types/status"
)

func getRdsDBCluster(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := rdsClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeDBClusters(
		ctx,
		&rds.DescribeDBClustersInput{
			DBClusterIdentifier: &name,
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.DBClusters) == 0 {
		return nil, errNotFound
	}

	if resp.DBClusters[0].Status == nil {
		return &status.Status{}, nil
	}

	st := rdsDBClusterStatusConverter.Convert(*resp.DBClusters[0].Status, "")

	return st, nil
}
