package resourcestatus

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

// resourceTypes indicate supported resource type and function to get status.
var resourceTypes map[string]getStatusFunc

// getStatusFunc is function use resource id to get resource status.
type getStatusFunc func(ctx context.Context, resourceType, name string) (*status.Status, error)

func init() {
	resourceTypes = map[string]getStatusFunc{
		"aws_cloudfront_distribution": getCloudFrontDistribution,
		"aws_ami":                     getEc2Image,
		"aws_instance":                getEc2Instance,
		"aws_network_interface":       getEc2NetworkInterface,
		"aws_ebs_snapshot":            getEc2Snapshot,
		"aws_subnet":                  getEc2Subnet,
		"aws_ebs_volume":              getEc2Volume,
		"aws_vpc":                     getEc2Vpc,
		"aws_eks_cluster":             getEksCluster,
		"aws_elasticache_cluster":     getElasticCacheCluster,
		"aws_lb":                      getElbLoadBalancer,
		"aws_rds_cluster":             getRdsDBCluster,
		"aws_db_instance":             getRdsDBInstance,
	}
}

// IsSupported indicate whether the resource type is supported.
func IsSupported(resourceType string) bool {
	_, ok := resourceTypes[resourceType]
	return ok
}

// Get resource status by resource type and name.
func Get(ctx context.Context, resourceType, name string) (*status.Status, error) {
	getFunc, exist := resourceTypes[resourceType]
	if !exist {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	st, err := getFunc(ctx, resourceType, name)
	if err != nil {
		return &status.Status{}, err
	}

	return st, nil
}
