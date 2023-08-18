package catalog

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/drone/go-scm/scm"

	"github.com/seal-io/walrus/pkg/bus/catalog"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/templates"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/version"
)

// GetRepos returns org and a list of repositories from the given catalog.
func GetRepos(ctx context.Context, mc model.ClientSet, c *model.Catalog) ([]*vcs.Repository, error) {
	var (
		client *scm.Client
		err    error
	)

	// Some catalog source url may be redirected, so we need to get the real url.
	source, err := GetRedirectURL(c.Source,
		version.GetUserAgent()+"; uuid="+settings.InstallationUUID.ShouldValue(ctx, mc))
	if err != nil {
		return nil, err
	}

	orgName, err := vcs.GetOrgFromGitURL(source)
	if err != nil {
		return nil, err
	}

	switch c.Type {
	case types.GitDriverGithub, types.GitDriverGitlab:
		client, err = vcs.NewClientFromURL(c.Type, source, "")
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

	repos, err := GetRepos(ctx, mc, c)
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
				repo := repos[j]

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

func GetRedirectURL(rawURL, userAgent string) (string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// If the status code is not 301 or 302, it means that the url is not redirected.
	if resp.StatusCode != http.StatusMovedPermanently && resp.StatusCode != http.StatusFound {
		return rawURL, nil
	}

	return resp.Header.Get("Location"), nil
}
