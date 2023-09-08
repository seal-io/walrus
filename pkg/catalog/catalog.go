package catalog

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/drone/go-scm/scm"
	"go.uber.org/multierr"

	"github.com/seal-io/walrus/pkg/bus/catalog"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/templates"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/pkg/vcs/options"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/version"
)

// getRepos returns org and a list of repositories from the given catalog.
func getRepos(ctx context.Context, c *model.Catalog, ua string) ([]*vcs.Repository, error) {
	var (
		client *scm.Client
		err    error
	)

	orgName, err := vcs.GetOrgFromGitURL(c.Source)
	if err != nil {
		return nil, err
	}

	switch c.Type {
	case types.GitDriverGithub, types.GitDriverGitlab:
		client, err = vcs.NewClientFromURL(c.Type, c.Source, options.WithUserAgent(ua))
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported catalog type %q", c.Type)
	}

	repos, err := vcs.GetOrgRepos(ctx, client, orgName)
	if err != nil {
		return nil, err
	}

	list := make([]*vcs.Repository, len(repos))
	for i := range repos {
		list[i] = &vcs.Repository{
			Namespace:   repos[i].Namespace,
			Name:        repos[i].Name,
			Description: repos[i].Description,
			Link:        repos[i].Link,
		}
	}

	return list, nil
}

func getSyncResult(ctx context.Context, mc model.ClientSet, c *model.Catalog) (*types.CatalogSync, error) {
	var (
		catalogSync = &types.CatalogSync{
			Total:     0,
			Succeeded: 0,
			Failed:    0,
			Time:      time.Now(),
		}
		counts []*struct {
			Status status.Status `json:"status"`
			Count  int           `json:"count"`
		}
	)

	err := mc.Templates().Query().
		Select(template.FieldStatus).
		Where(template.CatalogID(c.ID)).
		GroupBy(template.FieldStatus).
		Aggregate(model.Count()).
		Scan(ctx, &counts)
	if err != nil {
		return nil, err
	}

	for _, v := range counts {
		if status.CatalogStatusInitialized.IsTrue(v) || status.CatalogStatusReady.IsTrue(v) {
			catalogSync.Succeeded += v.Count
		}

		if status.CatalogStatusInitialized.IsFalse(v) || status.CatalogStatusReady.IsFalse(v) {
			catalogSync.Failed += v.Count
		}
		catalogSync.Total += v.Count
	}

	return catalogSync, nil
}

// SyncTemplates fetch and update catalog templates.
func SyncTemplates(ctx context.Context, mc model.ClientSet, c *model.Catalog) error {
	logger := log.WithName("catalog")

	ua := version.GetUserAgent() + "; uuid=" + settings.InstallationUUID.ShouldValue(ctx, mc)

	repos, err := getRepos(ctx, c, ua)
	if err != nil {
		return err
	}

	logger.Infof("found %d repositories in %s", len(repos), c.Source)

	var (
		total     = len(repos)
		processed = int32(0)
		failed    = int32(0)

		wg = gopool.Group()
	)

	batchSize := 10
	for i := 0; i < batchSize; i++ {
		s := i

		wg.Go(func() error {
			// Merge the errors to return them all at once,
			// instead of returning the first error.
			var berr error

			for j := s; j < len(repos); j += batchSize {
				repo := repos[j]
				repo.Driver = c.Type

				t := &model.Template{
					Name:        repo.Name,
					Description: repo.Description,
					Source:      repo.Link,
					CatalogID:   c.ID,
				}

				logger.Debugf("syncing  \"%s:%s\" of catalog %q", c.Name, repo.Name, c.ID)

				serr := templates.SyncTemplateFromGitRepo(ctx, mc, t, repo)
				if serr != nil {
					logger.Debugf("failed sync \"%s:%s\" of catalog %q: %v", c.Name, repo.Name, c.ID, serr)
					berr = multierr.Append(berr,
						fmt.Errorf("error syncing \"%s:%s\" of catalog %q: %w",
							c.Name, repo.Name, c.ID, serr))

					atomic.AddInt32(&failed, 1)
				} else {
					atomic.AddInt32(&processed, 1)
				}

				logger.Debugf("synced catalog %s, total: %d, processed: %d, failed: %d",
					c.Name, total, processed, failed)
			}

			return berr
		})
	}

	return wg.Wait()
}

type catalogSyncer struct {
	mc model.ClientSet
}

func CatalogSync(mc model.ClientSet) catalogSyncer {
	return catalogSyncer{mc: mc}
}

// Do Sync the given catalog, it will create or update templates from the given catalog.
func (cs catalogSyncer) Do(_ context.Context, busMessage catalog.BusMessage) error {
	var (
		logger = log.WithName("catalog")

		c = busMessage.Refer
	)

	gopool.Go(func() {
		subCtx := context.Background()

		serr := SyncTemplates(subCtx, cs.mc, c)

		uerr := UpdateStatusWithSyncErr(
			subCtx,
			cs.mc,
			c,
			serr,
		)
		if uerr != nil {
			logger.Errorf("failed to update catalog %s status: %v", c.Name, uerr)
		}
	})

	return nil
}

// UpdateStatusWithSyncErr update catalog status with sync error.
func UpdateStatusWithSyncErr(ctx context.Context, mc model.ClientSet, c *model.Catalog, syncErr error) error {
	logger := log.WithName("catalog")

	if syncErr != nil {
		status.CatalogStatusInitialized.False(c, syncErr.Error())
		logger.Warnf("failed to sync catalog %s templates: %v", c.Name, syncErr)
	} else {
		status.CatalogStatusReady.Reset(c, "")
		status.CatalogStatusReady.True(c, "")
	}

	c.Status.SetSummary(status.WalkCatalog(&c.Status))
	update := mc.Catalogs().UpdateOne(c).
		SetStatus(c.Status)

	syncResult, err := getSyncResult(ctx, mc, c)
	if err != nil {
		return fmt.Errorf("failed to update sync info: %w", err)
	}

	if syncResult != nil {
		update.SetSync(syncResult)
	}

	return update.Exec(ctx)
}
