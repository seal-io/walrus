package runtime

import (
	"github.com/gin-gonic/gin"
)

type (
	// ResourceRouteAdviceReceiver represents the type that can receive advice.
	ResourceRouteAdviceReceiver any

	// ResourceRouteAdviceProvider is a provider to provide advice to the request
	// of the routes belongs to a IResourceHandler.
	ResourceRouteAdviceProvider interface {
		// CanSet validates the given ResourceRouteAdviceReceiver can be set or not in prepare phase,
		// returns true if the given ResourceRouteAdviceReceiver can be injected.
		// The given ResourceRouteAdviceReceiver is stateless,
		// please do not perform additional operations on it.
		CanSet(ResourceRouteAdviceReceiver) bool

		// Set injects the valid ResourceRouteAdviceReceiver by this provider before validating,
		// the provider should set the corresponding advice to the target.
		Set(ResourceRouteAdviceReceiver)
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

// WithResourceRouteAdviceProviders is a RouterOption to configure the advice providers
// for the routes of IResourceHandler.
func WithResourceRouteAdviceProviders(providers ...ResourceRouteAdviceProvider) RouterOption {
	return routerOption(func(r *Router) {
		for i := range providers {
			if providers[i] == nil {
				continue
			}

			r.resourceRouteAdviceProviders = append(r.resourceRouteAdviceProviders, providers[i])
		}
	})
}
