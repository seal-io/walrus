package platformtf

import (
	"context"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformtf"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
	"github.com/seal-io/seal/utils/topic"
)

const Name topic.Topic = "Terraform"

type Message struct {
	ModelClient         model.ClientSet
	ApplicationRevision *model.ApplicationRevision
}

func Notify(ctx context.Context, name topic.Topic, message Message) error {
	return topic.Publish(ctx, name, message)
}

func AddSubscriber(ctx context.Context, name topic.Topic) error {
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

			message, ok := e.Data.(Message)
			if !ok {
				logger.Warnf("message type error, data: %v", e.Data)
				continue
			}

			if updateErr := updateResource(ctx, message); updateErr != nil {
				logger.Errorf("update resource err: %v", updateErr)
			}
		}
	})

	return nil
}

func updateResource(ctx context.Context, message Message) error {
	var parser platformtf.Parser
	applicationResources, err := parser.ParseAppRevision(message.ApplicationRevision)
	if err != nil {
		return err
	}

	return message.ModelClient.WithTx(ctx, func(tx *model.Tx) error {
		var (
			existResourceIDs = make([]types.ID, 0)
			newResources     = make(model.ApplicationResources, 0)
		)

		// fetch the old resources of the application
		oldResources, err := message.ModelClient.ApplicationResources().
			Query().
			Where(applicationresource.InstanceID(message.ApplicationRevision.InstanceID)).
			All(ctx)
		if err != nil {
			return err
		}
		oldResourceSet := sets.NewString()
		for _, r := range oldResources {
			uniqueKey := getFingerprint(r)
			oldResourceSet.Insert(uniqueKey)
		}

		for _, ar := range applicationResources {
			// check if the resource is exists.
			key := getFingerprint(ar)
			exists := oldResourceSet.Has(key)
			if exists {
				existResourceIDs = append(existResourceIDs, ar.ID)
			} else {
				newResources = append(newResources, ar)
			}
		}

		// diff application resource of this revision and the latest revision.
		// if the resource is not in the latest revision, delete it.
		_, err = message.ModelClient.ApplicationResources().
			Delete().
			Where(
				applicationresource.InstanceID(message.ApplicationRevision.InstanceID),
				applicationresource.IDNotIn(existResourceIDs...),
			).
			Exec(ctx)
		if err != nil {
			return err
		}

		// create newResource.
		if len(newResources) > 0 {
			resourcesToCreate, err := dao.ApplicationResourceCreates(message.ModelClient, newResources...)
			if err != nil {
				return err
			}
			if _, err = message.ModelClient.ApplicationResources().CreateBulk(resourcesToCreate...).Save(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

// TODO(thxCode): generate by entc.
func getFingerprint(r *model.ApplicationResource) string {
	// align to schema definition.
	return strs.Join("-", string(r.ConnectorID), r.Module, r.Mode, r.Type, r.Name)
}
