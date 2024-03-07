package resourcestatus

import (
	"context"
	"errors"
	"fmt"
	"strings"

	memorystore "cloud.google.com/go/redis/apiv1"
	"cloud.google.com/go/redis/apiv1/redispb"
	"google.golang.org/api/option"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	gtypes "github.com/seal-io/walrus/pkg/operator/google/types"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func getRedisInstance(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cred, ok := ctx.Value(types.CredentialKey).(*gtypes.Credential)
	if !ok {
		return nil, errors.New("not found credential from context")
	}

	client, err := memorystore.NewCloudRedisClient(ctx, option.WithCredentialsJSON([]byte(cred.Credentials)))
	if err != nil {
		return nil, err
	}

	req := &redispb.GetInstanceRequest{
		Name: name,
	}

	instance, err := client.GetInstance(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get google resource %s %s: %w", resourceType, name, err)
	}

	return redisInstanceStatusConverter.Convert(
		strings.ToLower(redispb.Instance_State_name[int32(instance.GetState())]),
		instance.GetStatusMessage(),
	), nil
}
