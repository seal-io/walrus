package resourcerevision

import (
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
)

type (
	GetRequest = model.ResourceRevisionQueryInput

	GetResponse = *model.ResourceRevisionOutput
)

type (
	CollectionGetRequest struct {
		model.ResourceRevisionQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.ResourceRevision, resourcerevision.OrderOption,
		] `query:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.ResourceRevisionOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest struct {
	model.ResourceRevisionDeleteInputs `path:",inline" json:",inline"`
}

func (r *CollectionDeleteRequest) Validate() error {
	if err := r.ResourceRevisionDeleteInputs.Validate(); err != nil {
		return err
	}

	latestRevision, err := r.Client.ResourceRevisions().Query().
		Where(resourcerevision.ResourceID(r.Resource.ID)).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		Select(resourcerevision.FieldID).
		First(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get latest revision: %w", err)
	}

	for i := range r.Items {
		// Prevent deleting the latest revision.
		if r.Items[i].ID == latestRevision.ID {
			return errors.New("invalid ids: can not delete latest revision")
		}
	}

	return nil
}
