package resourcelog

import (
	"context"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/seal-io/walrus/pkg/operator/types"
)

func ecsClient(ctx context.Context) (*ecs.Client, error) {
	cred, err := types.CredentialFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	cli, err := ecs.NewClientWithAccessKey(cred.Region, cred.AccessKey, cred.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("error create alibaba ecs client %s: %w", cred.AccessKey, err)
	}

	cli.EnableAsync(10, 10)

	return cli, nil
}
