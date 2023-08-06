package service

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/serviceresource"
	"github.com/seal-io/seal/pkg/apis/servicerevision"
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
	return "Service"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		runtime.Alias(
			serviceresource.Handle(h.modelClient),
			"Resource"),
		runtime.Alias(
			servicerevision.Handle(h.modelClient, h.kubeConfig, h.tlsCertified),
			"Revision"),
	}
}
