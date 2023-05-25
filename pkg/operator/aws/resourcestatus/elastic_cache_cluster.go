package resourcestatus

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func getElasticCacheCluster(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cli, err := elasticCacheClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cli.DescribeCacheClusters(ctx, &elasticache.DescribeCacheClustersInput{
		CacheClusterId: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("error describe aws resource %s %s: %w", resourceType, name, err)
	}

	if len(resp.CacheClusters) == 0 {
		return nil, errNotFound
	}

	cache := resp.CacheClusters[0]
	if cache.CacheClusterStatus == nil {
		return &status.Status{}, nil
	}

	st := elasticCacheStatusConverter.Convert(*cache.CacheClusterStatus, "")

	return st, nil
}
