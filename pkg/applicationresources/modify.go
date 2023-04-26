package applicationresources

import (
	"context"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/types"
	resourcetopic "github.com/seal-io/seal/pkg/topic/applicationresource"
	"github.com/seal-io/seal/utils/strs"
)

func Update(ctx context.Context, message resourcetopic.TopicMessage) error {
	var updated []*model.ApplicationResource
	err := message.ModelClient.WithTx(ctx, func(tx *model.Tx) error {
		var (
			existResourceIDs = make([]types.ID, 0)
			newResources     = make(model.ApplicationResources, 0)
		)

		// fetch the old resources of the application
		oldResources, err := message.ModelClient.ApplicationResources().
			Query().
			Where(applicationresource.InstanceID(message.InstanceID)).
			All(ctx)
		if err != nil {
			return err
		}

		oldResourceKeys := make(map[string]types.ID)
		for _, r := range oldResources {
			uniqueKey := getFingerprint(r)
			oldResourceKeys[uniqueKey] = r.ID
		}

		for _, ar := range message.ApplicationResources {
			// check if the resource is exists.
			key := getFingerprint(ar)
			oldResourceID, exists := oldResourceKeys[key]
			if exists {
				existResourceIDs = append(existResourceIDs, oldResourceID)
			} else {
				newResources = append(newResources, ar)
			}
		}

		// diff application resource of this revision and the latest revision.
		// if the resource is not in the latest revision, delete it.
		_, err = message.ModelClient.ApplicationResources().
			Delete().
			Where(
				applicationresource.InstanceID(message.InstanceID),
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
			if updated, err = message.ModelClient.ApplicationResources().CreateBulk(resourcesToCreate...).Save(ctx); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	var ids = make([]types.ID, len(updated))
	for i := range updated {
		ids[i] = updated[i].ID
	}
	entities, err := ListLabelCandidatesByIDs(ctx, message.ModelClient, ids)
	if err != nil {
		return err
	}
	return Label(ctx, entities)
}

// TODO(thxCode): generate by entc.
func getFingerprint(r *model.ApplicationResource) string {
	// align to schema definition.
	return strs.Join("-", string(r.ConnectorID), r.Module, r.Mode, r.Type, r.Name)
}
