package bus

import (
	"context"

	authstoken "github.com/seal-io/seal/pkg/auths/token"
	"github.com/seal-io/seal/pkg/bus/catalog"
	"github.com/seal-io/seal/pkg/bus/environment"
	"github.com/seal-io/seal/pkg/bus/servicerevision"
	"github.com/seal-io/seal/pkg/bus/setting"
	"github.com/seal-io/seal/pkg/bus/template"
	"github.com/seal-io/seal/pkg/bus/token"
	pkgcatalog "github.com/seal-io/seal/pkg/catalog"
	"github.com/seal-io/seal/pkg/cron"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/deployer/terraform"
	pkgenv "github.com/seal-io/seal/pkg/environment"
	"github.com/seal-io/seal/pkg/templates"
)

type SetupOptions struct {
	ModelClient model.ClientSet
}

func Setup(ctx context.Context, opts SetupOptions) (err error) {
	// Environment.
	err = environment.AddSubscriber("managed-kubernetes-namespace-sync",
		pkgenv.SyncManagedKubernetesNamespace)
	if err != nil {
		return
	}

	// ServiceRevision.
	err = servicerevision.AddSubscriber("terraform-sync-service-revision-status",
		terraform.SyncServiceRevisionStatus)
	if err != nil {
		return
	}

	// Setting.
	err = setting.AddSubscriber("cron-sync",
		cron.Sync)
	if err != nil {
		return
	}

	// Template.
	err = template.AddSubscriber("sync-template-schema",
		templates.SchemaSync(opts.ModelClient).Do)
	if err != nil {
		return
	}

	// Token.
	err = token.AddSubscriber("auths-token-delete-cached",
		authstoken.DelCached)
	if err != nil {
		return
	}

	// Catalog.
	err = catalog.AddSubscriber("sync-catalog", pkgcatalog.Sync)
	if err != nil {
		return
	}

	return
}
