package catalog

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	pkgcatalog "github.com/seal-io/walrus/pkg/catalog"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/catalog"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
)

type TemplateSyncTask struct {
	logger      log.Logger
	modelClient model.ClientSet
}

const summaryStatusReady = "Ready"

func NewCatalogTemplateSyncTask(logger log.Logger, mc model.ClientSet) (in *TemplateSyncTask, err error) {
	in = &TemplateSyncTask{
		logger:      logger,
		modelClient: mc,
	}

	return
}

func (in *TemplateSyncTask) Process(ctx context.Context, args ...any) error {
	return in.syncCatalogTemplates(ctx)
}

func (in *TemplateSyncTask) syncCatalogTemplates(ctx context.Context) error {
	if !settings.EnableSyncCatalog.ShouldValueBool(ctx, in.modelClient) {
		// Disable sync catalog.
		return nil
	}

	catalogs, err := in.modelClient.Catalogs().Query().
		Where(
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					catalog.FieldStatus,
					summaryStatusReady,
					sqljson.Path("summaryStatus"),
				))
			},
		).All(ctx)
	if err != nil {
		return err
	}

	for _, c := range catalogs {
		status.CatalogStatusInitialized.Reset(c, "Initializing catalog templates")
		c.Status.SetSummary(status.WalkCatalog(&c.Status))

		err := in.modelClient.Catalogs().UpdateOne(c).
			SetStatus(c.Status).
			Exec(ctx)
		if err != nil {
			return err
		}

		if err := pkgcatalog.SyncTemplates(ctx, in.modelClient, c); err != nil {
			in.logger.Errorf("failed to sync templates for catalog %s: %v", catalog.Name, err)
		}
	}

	return nil
}
