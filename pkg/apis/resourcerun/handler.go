package resourcerun

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
)

func Handle(mc model.ClientSet, kc *rest.Config) Handler {
	return Handler{
		modelClient: mc,
		kubeConfig:  kc,
	}
}

type Handler struct {
	modelClient model.ClientSet
	kubeConfig  *rest.Config
}

func (Handler) Kind() string {
	return "ResourceRun"
}
