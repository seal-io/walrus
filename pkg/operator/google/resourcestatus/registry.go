package resourcestatus

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/types/status"
)

// resourceTypes indicate supported resource type and function to get status.
var resourceTypes map[string]getStatusFunc

// getStatusFunc is function use resource id to get resource status.
type getStatusFunc func(ctx context.Context, resourceType, name string) (*status.Status, error)

func init() {
	resourceTypes = map[string]getStatusFunc{
		"google_sql_database_instance": getSQLDatabaseInstance,
		"google_redis_instance":        getRedisInstance,
		"google_compute_instance":      getComputeInstance,
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
