package resourcedefinition

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
)

type (
	RouteGetResourcesRequest struct {
		_ struct{} `route:"GET=/resources"`

		runtime.RequestCollection[
			predicate.Resource, resource.OrderOption,
		] `query:",inline"`

		model.ResourceDefinitionQueryInput `path:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	RouteGetResourcesResponse = []*model.ResourceOutput
)
