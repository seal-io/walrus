package datamessage

import (
	"context"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/topic"
)

// Topics.
const (
	// Service is the topic for service.
	Service topic.Topic = "Service"

	// ServiceRevision is the topic for service revision.
	ServiceRevision topic.Topic = "ServiceRevision"

	// ServiceResource is the topic for service resource.
	ServiceResource topic.Topic = "ServiceResource"

	// Connector is the topic for connector.
	Connector topic.Topic = "Connector"

	// Template is the topic for module.
	Template topic.Topic = "Template"
)

type Message[T any] struct {
	Type EventType
	Data []T
}

var allowed = sets.New(
	ServiceResource,
	Service,
	ServiceRevision,
	Connector,
	Template,
)

func IsAllowed(mutationType string) bool {
	return allowed.Has(topic.Topic(mutationType))
}

func Publish[T any](ctx context.Context, mutationType string, op model.Op, ids []T) error {
	if len(ids) == 0 || !IsAllowed(mutationType) {
		return nil
	}

	m := Message[T]{
		Type: EventTypeFor(op),
		Data: ids,
	}

	return topic.Publish(ctx, topic.Topic(mutationType), m)
}
