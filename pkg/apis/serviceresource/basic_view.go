package serviceresource

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresource"
)

type (
	CollectionGetRequest struct {
		model.ServiceResourceQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.ServiceResource, serviceresource.OrderOption,
		] `query:",inline"`

		WithoutKeys bool `query:"withoutKeys,omitempty"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.ServiceResourceOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}
