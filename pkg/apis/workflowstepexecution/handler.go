package workflowstepexecution

import (
	"k8s.io/client-go/rest"

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
	return "WorkflowStepExecution"
}
