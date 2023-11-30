package workflow

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/apis/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/model"
	pkgworkflow "github.com/seal-io/walrus/pkg/workflow"
)

func Handle(mc model.ClientSet, k8sConfig *rest.Config, wc pkgworkflow.Client) Handler {
	return Handler{
		modelClient:    mc,
		k8sConfig:      k8sConfig,
		workflowClient: wc,
	}
}

type Handler struct {
	modelClient    model.ClientSet
	k8sConfig      *rest.Config
	workflowClient pkgworkflow.Client
}

func (Handler) Kind() string {
	return "Workflow"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		runtime.Alias(
			workflowexecution.Handle(h.modelClient, h.k8sConfig, h.workflowClient),
			"Execution",
		),
	}
}
