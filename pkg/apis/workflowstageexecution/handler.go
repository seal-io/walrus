package workflowstageexecution

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/apis/workflowstepexecution"
	"github.com/seal-io/walrus/pkg/dao/model"
)

type Handler struct {
	modelClient model.ClientSet
	k8sConfig   *rest.Config
}

func Handle(mc model.ClientSet, kc *rest.Config) Handler {
	return Handler{
		modelClient: mc,
		k8sConfig:   kc,
	}
}

func (Handler) Kind() string {
	return "WorkflowStageExecution"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		runtime.Alias(
			workflowstepexecution.Handle(h.modelClient, h.k8sConfig),
			"StepExecution",
		),
	}
}
