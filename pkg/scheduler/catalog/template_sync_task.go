package catalog

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"go.uber.org/multierr"

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
	// Disable sync catalog.
	if !settings.EnableSyncCatalog.ShouldValueBool(ctx, in.modelClient) {
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

	if len(catalogs) == 0 {
		return nil
	}

	// Merge the errors to return them all at once,
	// instead of returning the first error.
	var berr error

	for i := range catalogs {
		c := catalogs[i]
		in.logger.Debugf("syncing templates of catalog %q", c.ID)

		status.CatalogStatusInitialized.Reset(c, "Initializing catalog templates")
		c.Status.SetSummary(status.WalkCatalog(&c.Status))

		uerr := in.modelClient.Catalogs().UpdateOne(c).
			SetStatus(c.Status).
			Exec(ctx)
		if multierr.AppendInto(&berr, uerr) {
			continue
		}

		serr := pkgcatalog.SyncTemplates(ctx, in.modelClient, c)

		uerr = pkgcatalog.UpdateStatusWithSyncErr(
			context.Background(), // Make sure status will be updated, in case of task timeout.
			in.modelClient,
			c,
			serr,
		)
		if uerr != nil {
			berr = multierr.Append(berr,
				fmt.Errorf("error syncing templates of catalog %s: %w",
					c.ID, uerr))
		}
	}

	return berr
}
