package resourcelog

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	opawstypes "github.com/seal-io/walrus/pkg/operator/aws/types"
)

func ec2Client(ctx context.Context) (*ec2.Client, error) {
	cfg, err := opawstypes.ConfigFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return ec2.NewFromConfig(*cfg), nil
}
