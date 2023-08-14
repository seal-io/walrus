package runtime

import (
	"github.com/gin-gonic/gin"
)

type (
	// RouteAdviceReceiver represents the type that can receive advice.
	RouteAdviceReceiver any

	// RouteAdviceProvider is a provider to provide advice to the request
	// of the reflected routes of a IHandler.
	RouteAdviceProvider interface {
		// CanSet validates the given RouteAdviceReceiver can be set or not in prepare phase,
		// returns true if the given RouteAdviceReceiver can be injected.
		// The given RouteAdviceReceiver is stateless,
		// please do not perform additional operations on it.
		CanSet(RouteAdviceReceiver) bool

		// Set injects the valid RouteAdviceReceiver by this provider before validating,
		// the provider should set the corresponding advice to the target.
		Set(RouteAdviceReceiver)
	}
)

// Built-in advice receivers.
type (
	// GinContextAdviceReceiver sets gin.Context
	// if the given request type implements this interface before validating.
	ginContextAdviceReceiver interface {
		// SetGinContext injects the session context before validating.
		SetGinContext(*gin.Context)
	}

	// UnidiStreamAdviceReceiver sets runtime.RequestUnidiStream
	// if the given request type implements this interface after validating.
	unidiStreamAdviceReceiver interface {
		// SetStream injects the runtime.RequestUnidiStream after validating.
		SetStream(RequestUnidiStream)
	}

	// BidiStreamAdviceReceiver sets inject runtime.RequestBidiStream
	// if the given request type implements this interface after validating.
	bidiStreamAdviceReceiver interface {
		// SetStream injects the runtime.RequestBidiStream after validating.
		SetStream(RequestBidiStream)
	}
)

// WithRouteAdviceProviders is a RouterOption to configure the advice providers
// for the reflected routes of a IHandler.
func WithRouteAdviceProviders(providers ...RouteAdviceProvider) RouterOption {
	return routerOption(func(r *Router) {
		for i := range providers {
			if providers[i] == nil {
				continue
			}

			r.adviceProviders = append(r.adviceProviders, providers[i])
		}
	})
}
