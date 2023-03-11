package platformtf

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/topic"
)

const Name topic.Topic = "Terraform"

type TopicMessage struct {
	ModelClient         model.ClientSet
	ApplicationRevision *model.ApplicationRevision
}

func Notify(ctx context.Context, name topic.Topic, message TopicMessage) error {
	return topic.Publish(ctx, name, message)
}

func AddSubscriber(ctx context.Context, name topic.Topic, h func(context.Context, TopicMessage) error) error {
	logger := log.WithName("topic").WithName(string(name))
	var t, err = topic.Subscribe(name)
	if err != nil {
		return err
	}

	gopool.Go(func() {
		for {
			e, err := t.Receive(ctx)
			if err != nil {
				logger.Errorf("receive message err: %v", err)
				t.Unsubscribe()
				return
			}

			message, ok := e.Data.(TopicMessage)
			if !ok {
				logger.Warnf("message type error, data: %v", e.Data)
				continue
			}

			if updateErr := h(ctx, message); updateErr != nil {
				logger.Errorf("update resource err: %v", updateErr)
			}
		}
	})

	return nil
}
