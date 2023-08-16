package cost

import (
	"github.com/seal-io/walrus/pkg/costs/distributor"
	"github.com/seal-io/walrus/pkg/dao/model"
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
		distributor: distributor.New(mc),
	}
}

type Handler struct {
	modelClient model.ClientSet
	distributor *distributor.Distributor
}

func (Handler) Kind() string {
	return "Cost"
}
