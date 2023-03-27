package topic

import (
	"context"

	"github.com/seal-io/seal/pkg/applicationresources"
	"github.com/seal-io/seal/pkg/dao/model"
	resourcetopic "github.com/seal-io/seal/pkg/topic/applicationresource"
)

type SetupOptions struct {
	ModelClient model.ClientSet
}

func Setup(ctx context.Context, opts SetupOptions) error {
	if err := resourcetopic.AddSubscriber(ctx, resourcetopic.Name, applicationresources.Update); err != nil {
		return err
	}

	return nil
}
