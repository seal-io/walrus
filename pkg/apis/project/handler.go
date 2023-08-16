package project

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis/connector"
	"github.com/seal-io/walrus/pkg/apis/environment"
	"github.com/seal-io/walrus/pkg/apis/projectsubject"
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/apis/variable"
	"github.com/seal-io/walrus/pkg/dao/model"
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
	return "Project"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		connector.Handle(h.modelClient),
		environment.Handle(h.modelClient, h.kubeConfig, h.tlsCertified),
		variable.Handle(h.modelClient),
		runtime.Alias(
			projectsubject.Handle(h.modelClient),
			"Subject"),
	}
}
