package templateversion

import (
	"github.com/seal-io/walrus/pkg/dao/model"
)

type RouteResetRequest struct {
	_ struct{} `route:"POST=/reset"`

	model.TemplateVersionQueryInput `path:",inline"`
}
