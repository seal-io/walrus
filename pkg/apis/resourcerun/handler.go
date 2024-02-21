package resourcerun

import (
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/storage"
)

func Handle(mc model.ClientSet, kc *rest.Config, sm *storage.Manager) Handler {
	return Handler{
		modelClient:    mc,
		kubeConfig:     kc,
		storageManager: sm,
	}
}

type Handler struct {
	modelClient    model.ClientSet
	kubeConfig     *rest.Config
	storageManager *storage.Manager
}

func (Handler) Kind() string {
	return "ResourceRun"
}
