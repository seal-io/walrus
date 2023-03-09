package server

import (
	"context"
	"fmt"

	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao/model"
)

type initOptions struct {
	K8sConfig   *rest.Config
	ModelClient *model.Client
}

func (r *Server) init(ctx context.Context, opts initOptions) error {
	// create model schema.
	var err = opts.ModelClient.Schema.Create(ctx)
	if err != nil {
		return fmt.Errorf("error creating model schemas: %w", err)
	}

	// initialize critical resources.
	type initor struct {
		name string
		init func(context.Context, initOptions) error
	}

	var inits = []initor{
		{name: "settings", init: r.initSettings},
		{name: "cronjobs", init: r.initCronJobs},
		{name: "subscribers", init: r.initSubscribers},
	}
	inits = append(inits,
		initor{name: "modules", init: r.initModules},
		initor{name: "roles", init: r.initRoles},
		initor{name: "subjects", init: r.initSubjects},
		initor{name: "perspective", init: r.initPerspectives},
		initor{name: "projects", init: r.initDefaultProject},
		initor{name: "environments", init: r.initDefaultEnvironment},
		initor{name: "deployer-runtime", init: r.initDeployerRuntime},
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
