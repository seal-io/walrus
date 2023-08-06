package template

import (
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
)

type RouteRefreshRequest struct {
	_ struct{} `route:"POST=/refresh"`

	model.TemplateQueryInput `path:",inline"`
}

type (
	RouteGetVersionsRequest struct {
		_ struct{} `route:"GET=/versions"`

		model.TemplateQueryInput `path:",inline"`

		runtime.RequestCollection[
			predicate.TemplateVersion, templateversion.OrderOption,
		] `query:",inline"`
	}

	RouteGetVersionsResponse = []*model.TemplateVersionOutput
)
