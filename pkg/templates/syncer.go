package templates

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/bus/template"
	"github.com/seal-io/walrus/pkg/bus/templateversion"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

func SchemaSync(mc model.ClientSet) schemaSyncer {
	return schemaSyncer{mc: mc}
}

type schemaSyncer struct {
	mc model.ClientSet
}

// Do fetches and updates the schema of the given template,
// within 5 mins in the background.
func (s schemaSyncer) Do(_ context.Context, message template.BusMessage) error {
	logger := log.WithName("template")

	gopool.Go(func() {
		ctx := context.Background()

		m := message.Refer
		logger.Debugf("syncing schema for template %s", m.ID)

		// Sync schema.
		err := syncSchema(ctx, s.mc, m)
		if err == nil {
			return
		}

		logger.Warnf("recording syncing template %s schema failed: %v", m.ID, err)

		status.TemplateStatusInitialized.False(m, fmt.Sprintf("sync schema failed: %v", err))
		m.Status.SetSummary(status.WalkTemplate(&m.Status))

		// State template.
		err = s.mc.Templates().UpdateOne(m).
			SetStatus(m.Status).
			Exec(ctx)
		if err != nil {
			logger.Errorf("failed to update template %s: %v", m.ID, err)
		}
	})

	return nil
}

func syncSchema(ctx context.Context, mc model.ClientSet, t *model.Template) error {
	repo, err := vcs.ParseURLToRepo(t.Source)
	if err != nil {
		return err
	}

	if t.CatalogID.Valid() {
		c, err := mc.Catalogs().Get(ctx, t.CatalogID)
		if err != nil {
			return err
		}

		// Use the catalog type as the vcs repository type.
		repo.Driver = c.Type
	}

	if repo.Reference != "" {
		return SyncTemplateFromGitRef(ctx, mc, t, repo)
	}

	return SyncTemplateFromGitRepo(ctx, mc, t, repo)
}

func VersionSchemaSync(mc model.ClientSet) versionSchemaSyncer {
	return versionSchemaSyncer{mc: mc}
}

type versionSchemaSyncer struct {
	mc model.ClientSet
}

func (s versionSchemaSyncer) Do(_ context.Context, message templateversion.BusMessage) error {
	logger := log.WithName("template-version")

	gopool.Go(func() {
		ctx := context.Background()

		m := message.Refer
		logger.Debugf("syncing schema for template version %s", m.ID)

		// Sync schema.
		err := SyncTemplateVersion(ctx, s.mc, m)
		if err != nil {
			logger.Warnf("syncing template version %s schema failed: %v", m.ID, err)
			return
		}
	})

	return nil
}
