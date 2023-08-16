package apis

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
)

type (
	modelClientAdviceReceiver interface {
		SetModelClient(client *model.Client)
	}

	modelClientAdviceProvider struct {
		modelClient *model.Client
	}
)

func provideModelClient(modelClient *model.Client) runtime.RouteAdviceProvider {
	return modelClientAdviceProvider{modelClient: modelClient}
}

func (m modelClientAdviceProvider) CanSet(receiver runtime.RouteAdviceReceiver) bool {
	_, ok := receiver.(modelClientAdviceReceiver)
	return ok
}

func (m modelClientAdviceProvider) Set(receiver runtime.RouteAdviceReceiver) {
	receiver.(modelClientAdviceReceiver).SetModelClient(m.modelClient)
}
