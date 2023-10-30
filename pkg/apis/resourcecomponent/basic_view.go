package resourcecomponent

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
)

type (
	CollectionGetRequest struct {
		model.ResourceComponentQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.ResourceComponent, resourcecomponent.OrderOption,
		] `query:",inline"`

		WithoutKeys bool `query:"withoutKeys,omitempty"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.ResourceComponentOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}
