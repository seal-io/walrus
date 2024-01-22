package resourcedefinition

import (
	apiresource "github.com/seal-io/walrus/pkg/apis/resource"
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

		ProjectName      string `query:"projectName,omitempty"`
		MatchingRuleName string `query:"matchingRuleName,omitempty"`

		Stream *runtime.RequestUnidiStream
	}

	RouteGetResourcesResponse = []*model.ResourceOutput
)

func (r *RouteGetResourcesRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type (
	RouteDeleteResourcesRequest struct {
		_ struct{} `route:"DELETE=/resources"`

		apiresource.CollectionDeleteRequest `json:",inline" query:",inline"`
	}
)

type (
	RouteUpgradeResourcesRequest struct {
		_ struct{} `route:"POST=/resources/_/upgrade"`

		apiresource.CollectionRouteUpgradeRequest `json:",inline"`
	}
)
