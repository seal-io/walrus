package bus

import (
	"context"

	authstoken "github.com/seal-io/walrus/pkg/auths/token"
	"github.com/seal-io/walrus/pkg/bus/builtin"
	"github.com/seal-io/walrus/pkg/bus/catalog"
	"github.com/seal-io/walrus/pkg/bus/environment"
	"github.com/seal-io/walrus/pkg/bus/resourcerevision"
	"github.com/seal-io/walrus/pkg/bus/setting"
	"github.com/seal-io/walrus/pkg/bus/template"
	"github.com/seal-io/walrus/pkg/bus/templateversion"
	"github.com/seal-io/walrus/pkg/bus/token"
	pkgcatalog "github.com/seal-io/walrus/pkg/catalog"
	"github.com/seal-io/walrus/pkg/cron"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/deployer/terraform"
	pkgenv "github.com/seal-io/walrus/pkg/environment"
	"github.com/seal-io/walrus/pkg/resourcedefinitions"
	"github.com/seal-io/walrus/pkg/templates"
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

	// ResourceRevision.
	err = resourcerevision.AddSubscriber("terraform-sync-resource-revision-status",
		terraform.SyncResourceRevisionStatus)
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

	// TemplateVersion.
	err = templateversion.AddSubscriber("sync-template-version-schema",
		templates.VersionSchemaSync(opts.ModelClient).Do)
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
	err = catalog.AddSubscriber("sync-catalog",
		pkgcatalog.CatalogSync(opts.ModelClient).Do)
	if err != nil {
		return
	}

	// Builtin.
	err = builtin.AddSubscriber("sync-builtin-resource-definitions",
		resourcedefinitions.SyncBuiltinResourceDefinitions)
	if err != nil {
		return
	}

	return
}
