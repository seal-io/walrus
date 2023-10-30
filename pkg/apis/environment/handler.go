package environment

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis/resource"
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/apis/variable"
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
	return "Environment"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		variable.Handle(h.modelClient),
		resource.Handle(h.modelClient, h.kubeConfig),
	}
}
