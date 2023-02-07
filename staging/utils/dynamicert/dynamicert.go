package dynamicert

import (
	"net"
	"net/http"

	"github.com/rancher/dynamiclistener"
	"github.com/rancher/dynamiclistener/factory"
)

type HostFilter func(candidates ...string) []string

type Manager struct {
	HostFilter HostFilter
	Cache      Cache
}

func (m *Manager) Handle(ls net.Listener, h http.Handler) (net.Listener, http.Handler, error) {
	var caCert, caKey, err = factory.LoadOrGenCA()
	if err != nil {
		return nil, nil, err
	}
	var dynConfig = dynamiclistener.Config{
		FilterCN: m.HostFilter,
	}
	dynLs, dynH, err := dynamiclistener.NewListener(ls, cacheToStorage(m.Cache), caCert, caKey, dynConfig)
	if err != nil {
		return nil, nil, err
	}
	var wh = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dynH.ServeHTTP(w, r)
		h.ServeHTTP(w, r)
	})
	return dynLs, wh, nil
}
