package catalog

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/drone/go-scm/scm"

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

// GetOrgRepos returns full repositories list from the given org.
func GetOrgRepos(ctx context.Context, client *scm.Client, orgName string) ([]*scm.Repository, error) {
	opts := scm.ListOptions{Size: 100}

	var list []*scm.Repository

	for {
		repos, meta, err := client.Organizations.ListRepositories(ctx, orgName, opts)
		if err != nil {
			return nil, err
		}

		for _, src := range repos {
			if src != nil {
				list = append(list, src)
			}
		}

		opts.Page = meta.Page.Next
		opts.URL = meta.Page.NextURL

		if opts.Page == 0 && opts.URL == "" {
			break
		}
	}

	return list, nil
}

// GetRepos returns org and a list of repositories from the given catalog.
func GetRepos(ctx context.Context, c *model.Catalog) ([]*scm.Repository, error) {
	var (
		client *scm.Client
		err    error
	)

	switch c.Type {
	case types.GitDriverGithub, types.GitDriverGitlab:
		client, err = vcs.NewClientFromURL(c.Type, c.Source, "")
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
		counts = []*struct {
			Status status.Status `json:"status"`
			Count  int           `json:"count"`
		}{}
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
func SyncTemplates(ctx context.Context, mc model.ClientSet, c *model.Catalog) (err error) {
	logger := log.WithName("catalog")

	repos, err := GetRepos(ctx, c)
	if err != nil {
		return err
	}

	logger.Infof("found %d repositories in %s", len(repos), c.Source)

	wg := gopool.Group()

	batchSize := 10
	for i := 0; i < batchSize; i++ {
		s := i

		wg.Go(func() error {
			for j := s; j < len(repos); j += batchSize {
				repo := &vcs.Repository{
					Namespace:   repos[j].Namespace,
					Name:        repos[j].Name,
					Description: repos[j].Description,
					Link:        repos[j].Link,
				}

				logger.Debugf("sync template %s", repo.Name)

				if err := templates.SyncTemplateFromGitRepo(ctx, mc, c, repo); err != nil {
					logger.Warnf("failed to sync template %s: %v", repo.Name, err)
				}
			}

			return nil
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
