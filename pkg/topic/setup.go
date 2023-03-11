package topic

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platformtf"
	tftopic "github.com/seal-io/seal/pkg/topic/platformtf"
)

type SetupOptions struct {
	ModelClient model.ClientSet
}

func Setup(ctx context.Context, opts SetupOptions) error {
	if err := tftopic.AddSubscriber(ctx, tftopic.Name, platformtf.UpdateResource); err != nil {
		return err
	}

	return nil
}
