package resource

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis/resourcecomponent"
	"github.com/seal-io/walrus/pkg/apis/resourcerun"
	"github.com/seal-io/walrus/pkg/apis/runtime"
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
	return "Resource"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		runtime.Alias(
			resourcecomponent.Handle(h.modelClient),
			"Component"),
		runtime.Alias(
			resourcerun.Handle(h.modelClient, h.kubeConfig),
			"Run"),
	}
}
