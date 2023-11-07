package workflow

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/apis/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/model"
)

func Handle(mc model.ClientSet, k8sConfig *rest.Config) Handler {
	return Handler{
		modelClient: mc,
		k8sConfig:   k8sConfig,
	}
}

type Handler struct {
	modelClient model.ClientSet
	k8sConfig   *rest.Config
}

func (Handler) Kind() string {
	return "Workflow"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		runtime.Alias(
			workflowexecution.Handle(h.modelClient, h.k8sConfig),
			"Execution",
		),
	}
}
