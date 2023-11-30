package workflowexecution

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/apis/workflowstageexecution"
	"github.com/seal-io/walrus/pkg/dao/model"
	pkgworkflow "github.com/seal-io/walrus/pkg/workflow"
)

func Handle(mc model.ClientSet, kc *rest.Config, wc pkgworkflow.Client) Handler {
	return Handler{
		modelClient:    mc,
		k8sConfig:      kc,
		workflowClient: wc,
	}
}

type Handler struct {
	modelClient    model.ClientSet
	k8sConfig      *rest.Config
	workflowClient pkgworkflow.Client
}

func (Handler) Kind() string {
	return "WorkflowExecution"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		runtime.Alias(
			workflowstageexecution.Handle(h.modelClient, h.k8sConfig, h.workflowClient),
			"StageExecution",
		),
	}
}
