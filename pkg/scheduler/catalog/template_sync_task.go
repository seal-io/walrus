package catalog

import (
	"context"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	pkgcatalog "github.com/seal-io/seal/pkg/catalog"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/catalog"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/log"
)

type TemplateSyncTask struct {
	mu sync.Mutex

	logger      log.Logger
	modelClient model.ClientSet
}

const summaryStatusReady = "Ready"

func NewCatalogTemplateSyncTask(mc model.ClientSet) (*TemplateSyncTask, error) {
	in := &TemplateSyncTask{
		modelClient: mc,
	}

	in.logger = log.WithName("task").WithName(in.Name())

	return in, nil
}

func (in *TemplateSyncTask) Name() string {
	return "catalog-template-sync-task"
}

func (in *TemplateSyncTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}

	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

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
