package topic

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/topic/platformtf"
)

type SetupOptions struct {
	ModelClient model.ClientSet
}

func Setup(ctx context.Context, opts SetupOptions) error {
	if err := platformtf.AddSubscriber(ctx, platformtf.Name); err != nil {
		return err
	}

	return nil
}
