package servicerevision

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
)

type (
	GetRequest = model.ServiceRevisionQueryInput

	GetResponse = *model.ServiceRevisionOutput
)

type (
	CollectionGetRequest struct {
		model.ServiceRevisionQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.ServiceRevision, servicerevision.OrderOption,
		] `query:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.ServiceRevisionOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest struct {
	model.ServiceRevisionDeleteInputs `path:",inline" json:",inline"`
}

func (r *CollectionDeleteRequest) Validate() error {
	if err := r.ServiceRevisionDeleteInputs.Validate(); err != nil {
		return err
	}

	latestRevision, err := r.Client.ServiceRevisions().Query().
		Where(servicerevision.ServiceID(r.Service.ID)).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		Select(servicerevision.FieldID).
		First(r.Context)
	if err != nil {
		return runtime.Errorw(err, "failed to get latest revision")
	}

	for i := range r.Items {
		// Prevent deleting the latest revision.
		if r.Items[i].ID == latestRevision.ID {
			return errors.New("invalid ids: can not delete latest revision")
		}
	}

	return nil
}
