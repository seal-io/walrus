package environment

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/service"
	"github.com/seal-io/seal/pkg/apis/variable"
	"github.com/seal-io/seal/pkg/dao/model"
)

func Handle(mc model.ClientSet, kc *rest.Config, tc bool) Handler {
	return Handler{
		modelClient:  mc,
		kubeConfig:   kc,
		tlsCertified: tc,
	}
}

type Handler struct {
	modelClient  model.ClientSet
	kubeConfig   *rest.Config
	tlsCertified bool
}

func (Handler) Kind() string {
	return "Environment"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		variable.Handle(h.modelClient),
		service.Handle(h.modelClient, h.kubeConfig, h.tlsCertified),
	}
}
