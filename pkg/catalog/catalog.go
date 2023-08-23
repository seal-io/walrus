package catalog

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/drone/go-scm/scm"
	"go.uber.org/multierr"

	"github.com/seal-io/walrus/pkg/bus/catalog"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/templates"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

// getRepos returns org and a list of repositories from the given catalog.
func getRepos(ctx context.Context, c *model.Catalog, ua string) ([]*vcs.Repository, error) {
	var (
		client *scm.Client
		err    error
	)

	switch c.Type {
	case types.GitDriverGithub, types.GitDriverGitlab:
		client, err = vcs.NewClientFromURL(c.Type, c.Source, options.WithUserAgent(ua))
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported catalog type %q", c.Type)
	}

	orgName, err := GetOrgFromGitURL(c.Source)
	if err != nil {
		return nil, err
	}

	repos, err := GetOrgRepos(ctx, client, orgName)
	if err != nil {
		return nil, err
	}

	return repos, nil
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
		if status.CatalogStatusInitialized.IsTrue(v.Status) || status.CatalogStatusReady.IsTrue(v.Status) {
			catalogSync.Succeeded += v.Count
		}

		if status.CatalogStatusInitialized.IsFalse(v.Status) || status.CatalogStatusReady.IsFalse(v.Status) {
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

	wg := gopool.Group()

	batchSize := 10
	for i := 0; i < batchSize; i++ {
		s := i

		wg.Go(func() error {
			// Merge the errors to return them all at once,
			// instead of returning the first error.
			var berr error

			for j := s; j < len(repos); j += batchSize {
				repo := &vcs.Repository{
					Namespace:   repos[j].Namespace,
					Name:        repos[j].Name,
					Description: repos[j].Description,
					Link:        repos[j].Link,
				}

				logger.Debugf("syncing template %s of catalog %s",
					repo.Name, c.ID)

				serr := templates.SyncTemplateFromGitRepo(ctx, mc, c, repo)
				if serr != nil {
					berr = multierr.Append(berr,
						fmt.Errorf("error syncing template %s: %w",
							repo.Name, serr))
				}
			}

			return berr
		})
	}

	return wg.Wait()
}

// GetOrgFromGitURL parses the organization from the given git repository URL.
func GetOrgFromGitURL(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}

	// https://github.com/<org>/<repo>
	parts := strings.Split(u.Path, "/")

	if len(parts) >= 2 {
		return parts[1], nil
	}

	return "", fmt.Errorf("invalid git url")
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

		err := SyncTemplates(subCtx, cs.mc, c)
		if err != nil {
			status.CatalogStatusInitialized.False(c, err.Error())
			logger.Errorf("failed to sync catalog %s templates: %v", c.Name, err)
		} else {
			status.CatalogStatusReady.Reset(c, "")
			status.CatalogStatusReady.True(c, "")
		}

		c.Status.SetSummary(status.WalkCatalog(&c.Status))
		update := cs.mc.Catalogs().UpdateOne(c).
			SetStatus(c.Status)

		syncResult, err := getSyncResult(subCtx, cs.mc, c)
		if err != nil {
			logger.Errorf("failed to update sync info: %v", err)
		}

		if syncResult != nil {
			update.SetSync(syncResult)
		}

		rerr := update.Exec(subCtx)
		if rerr != nil {
			logger.Errorf("failed to update catalog %s status: %v", c.Name, rerr)
		}
	})

	return nil
}
