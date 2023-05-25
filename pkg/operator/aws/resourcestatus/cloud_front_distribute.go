package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func getCloudFrontDistribution(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := cloudfrontClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.GetDistribution(ctx, &cloudfront.GetDistributionInput{
		Id: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	var sm string
	if resp.Distribution != nil && resp.Distribution.Status != nil {
		sm = *resp.Distribution.Status
	}

	st := cloudFrontStatusConverter.Convert(sm, "")

	return st, nil
}
