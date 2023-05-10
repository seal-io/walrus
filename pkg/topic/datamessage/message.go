package datamessage

import (
	"context"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/topic"
)

// Application topics
const (
	// Application is the topic for application.
	Application topic.Topic = "Application"

	// ApplicationInstance is the topic for application instance.
	ApplicationInstance topic.Topic = "ApplicationInstance"

	// ApplicationRevision is the topic for application revision.
	ApplicationRevision topic.Topic = "ApplicationRevision"

	// ApplicationResource is the topic for application resource.
	ApplicationResource topic.Topic = "ApplicationResource"

	// Connector is the topic for connector.
	Connector topic.Topic = "Connector"

	// Module is the topic for module.
	Module topic.Topic = "Module"
)

type Message[T any] struct {
	Type EventType
	Data []T
}

var allowed = sets.New(
	Application,
	ApplicationResource,
	ApplicationInstance,
	ApplicationRevision,
	Connector,
	Module,
)

func Publish[T any](ctx context.Context, mutationType string, op model.Op, ids []T) error {
	if len(ids) == 0 {
		return nil
	}

	if !allowed.Has(topic.Topic(mutationType)) {
		return nil
	}

	var m = Message[T]{
		Type: EventTypeFor(op),
		Data: ids,
	}
	return topic.Publish(ctx, topic.Topic(mutationType), m)
}
