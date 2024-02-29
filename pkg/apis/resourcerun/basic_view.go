package resourcerun

import (
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
)

type (
	GetRequest = model.ResourceRunQueryInput

	GetResponse = *model.ResourceRunOutput
)

type (
	CollectionGetRequest struct {
		model.ResourceRunQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.ResourceRun, resourcerun.OrderOption,
		] `query:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.ResourceRunOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest struct {
	model.ResourceRunDeleteInputs `path:",inline" json:",inline"`
}

func (r *CollectionDeleteRequest) Validate() error {
	if err := r.ResourceRunDeleteInputs.Validate(); err != nil {
		return err
	}

	latestRun, err := r.Client.ResourceRuns().Query().
		Where(resourcerun.ResourceID(r.Resource.ID)).
		Order(model.Desc(resourcerun.FieldCreateTime)).
		Select(resourcerun.FieldID).
		First(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get latest run: %w", err)
	}

	for i := range r.Items {
		// Prevent deleting the latest run.
		if r.Items[i].ID == latestRun.ID {
			return errors.New("invalid ids: can not delete latest run")
		}
	}

	return nil
}
