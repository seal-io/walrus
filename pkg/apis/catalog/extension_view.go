package catalog

import (
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/template"
)

type RouteSyncCatalogRequest struct {
	_ struct{} `route:"POST=/refresh"`

	model.CatalogQueryInput `path:",inline" query:",inline"`
}

// NB(thxCode): the following cases is similar to pkg/apis/template/basic_view.go,
// but for multiple templates under a catalog at once query.
type (
	RouteGetTemplatesRequest struct {
		_ struct{} `route:"GET=/templates"`

		model.CatalogQueryInput `path:",inline"`

		runtime.RequestCollection[
			predicate.Template, template.OrderOption,
		] `query:",inline"`
	}

	RouteGetTemplatesResponse = []*model.TemplateOutput
)
