package templateversion

import (
	tvbus "github.com/seal-io/walrus/pkg/bus/templateversion"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
)

func (h Handler) RouteReset(req RouteResetRequest) error {
	entity, err := h.modelClient.TemplateVersions().Query().
		Where(templateversion.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	return tvbus.Notify(req.Context, entity)
}
