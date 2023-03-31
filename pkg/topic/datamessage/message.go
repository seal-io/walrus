package datamessage

import (
	"context"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
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
)

type Message struct {
	Type EventType
	Data []oid.ID
}

var allowed = sets.New(
	Application,
	ApplicationResource,
	ApplicationInstance,
	ApplicationRevision,
	Connector,
)

func Publish(ctx context.Context, mutationType string, op model.Op, ids []types.ID) error {
	if !allowed.Has(topic.Topic(mutationType)) {
		return nil
	}

	var m = Message{
		Type: EventTypeFor(op),
		Data: ids,
	}
	return topic.Publish(ctx, topic.Topic(mutationType), m)
}
