package catalog

import (
	"github.com/seal-io/walrus/pkg/dao/model"
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (Handler) Kind() string {
	return "Catalog"
}
