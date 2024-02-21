package project

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis/catalog"
	"github.com/seal-io/walrus/pkg/apis/connector"
	"github.com/seal-io/walrus/pkg/apis/environment"
	"github.com/seal-io/walrus/pkg/apis/projectsubject"
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/apis/template"
	"github.com/seal-io/walrus/pkg/apis/templateversion"
	"github.com/seal-io/walrus/pkg/apis/variable"
	"github.com/seal-io/walrus/pkg/apis/workflow"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/storage"
	pkgworkflow "github.com/seal-io/walrus/pkg/workflow"
)

func Handle(mc model.ClientSet, kc *rest.Config, wc pkgworkflow.Client, sm *storage.Manager) Handler {
	return Handler{
		modelClient:    mc,
		kubeConfig:     kc,
		workflowClient: wc,
		storageManager: sm,
	}
}

type Handler struct {
	modelClient    model.ClientSet
	kubeConfig     *rest.Config
	workflowClient pkgworkflow.Client
	storageManager *storage.Manager
}

func (Handler) Kind() string {
	return "Project"
}

func (h Handler) SubResourceHandlers() []runtime.IResourceHandler {
	return []runtime.IResourceHandler{
		connector.Handle(h.modelClient),
		environment.Handle(h.modelClient, h.kubeConfig, h.storageManager),
		variable.Handle(h.modelClient),
		workflow.Handle(h.modelClient, h.kubeConfig, h.workflowClient),
		catalog.Handle(h.modelClient),
		template.Handle(h.modelClient),
		templateversion.Handle(h.modelClient),
		runtime.Alias(
			projectsubject.Handle(h.modelClient),
			"Subject"),
	}
}
