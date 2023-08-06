package apis

import (
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
)

type (
	modelClientAdviceReceiver interface {
		SetModelClient(client *model.Client)
	}

	modelClientAdviceProvider struct {
		modelClient *model.Client
	}
)

func provideModelClient(modelClient *model.Client) runtime.ResourceRouteAdviceProvider {
	return modelClientAdviceProvider{modelClient: modelClient}
}

func (m modelClientAdviceProvider) CanSet(receiver runtime.ResourceRouteAdviceReceiver) bool {
	_, ok := receiver.(modelClientAdviceReceiver)
	return ok
}

func (m modelClientAdviceProvider) Set(receiver runtime.ResourceRouteAdviceReceiver) {
	receiver.(modelClientAdviceReceiver).SetModelClient(m.modelClient)
}
