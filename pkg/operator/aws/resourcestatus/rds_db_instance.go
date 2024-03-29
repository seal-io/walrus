package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"

	"github.com/seal-io/walrus/pkg/dao/types/status"
)

func getRdsDBInstance(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := rdsClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeDBInstances(
		ctx,
		&rds.DescribeDBInstancesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("dbi-resource-id"),
					Values: []string{name},
				},
			},
		})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.DBInstances) == 0 {
		return nil, errNotFound
	}

	if resp.DBInstances[0].DBInstanceStatus == nil {
		return &status.Status{}, nil
	}

	st := rdsDBInstanceStatusConverter.Convert(*resp.DBInstances[0].DBInstanceStatus, "")

	return st, nil
}
