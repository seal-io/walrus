package webhookserver

import (
	"context"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/healthz"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
)

func Dummy() ctrlwebhook.Server {
	return &dummy{}
}

type dummy struct{}

func (dummy) NeedLeaderElection() bool {
	return false
}

func (dummy) Register(path string, hook http.Handler) {
}

func (dummy) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (dummy) StartedChecker() healthz.Checker {
	return func(req *http.Request) error {
		return nil
	}
}

func (dummy) WebhookMux() *http.ServeMux {
	return http.NewServeMux()
}

func (dummy) IsDummy() bool {
	return true
}
