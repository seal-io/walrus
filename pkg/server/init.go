package server

import (
	"context"
	"database/sql"
	"fmt"

	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/cache"
	"github.com/seal-io/seal/pkg/dao/migration"
	"github.com/seal-io/seal/pkg/dao/model"
)

type initOptions struct {
	K8sConfig      *rest.Config
	K8sCacheReady  chan struct{}
	ModelClient    *model.Client
	SkipTLSVerify  bool
	DatabaseDriver *sql.DB
	CacheDriver    cache.Driver
}

func (r *Server) init(ctx context.Context, opts initOptions) error {
	// Create model schema.
	err := opts.ModelClient.Schema.Create(ctx, migration.Options...)
	if err != nil {
		return fmt.Errorf("error creating model schemas: %w", err)
	}

	// Initialize critical resources.
	type initor struct {
		name string
		init func(context.Context, initOptions) error
	}

	inits := []initor{
		{name: "settings", init: r.initSettings},
		{name: "configs", init: r.initConfigs},
		{name: "dispatches", init: r.initDispatches},
		{name: "metrics", init: r.initMetrics},
		{name: "healthCheckers", init: r.initHealthCheckers},
		{name: "backgroundJobs", init: r.initBackgroundJobs},
		{name: "subscribers", init: r.initSubscribers},
		{name: "rbac", init: r.initRbac},
	}
	inits = append(inits,
		initor{name: "modules", init: r.initTemplates},
		initor{name: "perspective", init: r.initPerspectives},
		initor{name: "projects", init: r.initProjects},
	)

	if r.EnableAuthn {
		inits = append(inits,
			initor{name: "casdoor", init: r.initCasdoor},
		)
	}

	for i := range inits {
		if err = inits[i].init(ctx, opts); err != nil {
			return fmt.Errorf("%s: %w", inits[i].name, err)
		}
	}

	return nil
}
