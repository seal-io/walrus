package resourcedefinition

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
)

func Handle(mc model.ClientSet, config *rest.Config) Handler {
	return Handler{
		modelClient: mc,
		kubeConfig:  config,
	}
}

type Handler struct {
	modelClient model.ClientSet
	kubeConfig  *rest.Config
}

func (Handler) Kind() string {
	return "ResourceDefinition"
}
