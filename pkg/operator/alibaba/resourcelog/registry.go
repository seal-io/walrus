package resourcelog

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/operator/alibaba/key"
	"github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/utils/log"
)

// resourceTypes indicate supported resource type and their functions.
var resourceTypes map[string]getLoggableResource

type getLoggableResource func(ctx context.Context) (types.LoggableResource, error)

func init() {
	resourceTypes = map[string]getLoggableResource{
		"alicloud_instance": getEcsInstance,
	}
}

// Supported indicate whether the resource is supported to get log.
func Supported(_ context.Context, k string) (bool, error) {
	resourceType, _, ok := key.Decode(k)
	if !ok {
		return false, errors.New("invalid key")
	}

	_, exist := resourceTypes[resourceType]
	if !exist {
		return false, errUnsupported
	}

	return true, nil
}

// Log get resource log by key.
func Log(ctx context.Context, k string, options types.LogOptions) error {
	supported, err := Supported(ctx, k)
	if err != nil {
		return err
	}

	if !supported {
		return errUnsupported
	}

	resourceType, name, ok := key.Decode(k)
	if !ok {
		return errors.New("invalid key")
	}

	res := resourceTypes[resourceType]

	fs, err := res(ctx)
	if err != nil {
		return err
	}

	err = fs.Log(ctx, name, options)
	if err != nil {
		log.Warnf("error get log resource %s/%s: %v", resourceType, name, err)
		return err
	}

	return nil
}
