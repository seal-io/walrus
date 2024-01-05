package walrusfilehub

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

type walrusFileItem struct {
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Readme  string `json:"readme"`
	Content string `json:"content"`
}

var (
	yamlFileSuffix = []string{".yaml", ".yml"}
	icons          = []string{
		"icon.png",
		"icon.jpg",
		"icon.jpeg",
		"icon.svg",
	}
	walrusFileList   []walrusFileItem
	walrusFileHubDir = "/var/lib/walrus/walrus-file-hub"
)

const (
	fileNameReadme = "README.md"
)

func Index(ctx context.Context, mc model.ClientSet) runtime.Handle {
	if err := initWalrusFileHub(ctx, mc); err != nil {
		log.Errorf("failed to init walrus file hub: %w", err)
	}

	return func(c *gin.Context) {
		p, _ := c.Params.Get("filepath")

		switch p {
		case "/":
			c.JSON(http.StatusOK, walrusFileList)
		default:
			// Assets(Icons).
			fs := runtime.StaticHttpFileSystem{
				FileSystem: http.FS(os.DirFS(walrusFileHubDir)),
			}

			req := c.Request.Clone(c.Request.Context())
			req.URL.Path = p
			req.URL.RawPath = p
			http.FileServer(fs).ServeHTTP(c.Writer, req)
			c.Next()
		}
	}
}

func initWalrusFileHub(ctx context.Context, mc model.ClientSet) error {
	walrusFileHubURL, err := settings.WalrusFileHubURL.Value(ctx, mc)
	if err != nil {
		return fmt.Errorf("failed to get walrus file hub URL: %w", err)
	}

	u, err := url.Parse(walrusFileHubURL)
	if err != nil {
		return fmt.Errorf("failed to parse walrus file hub URL: %w", err)
	}

	walrusFileList, err = getWalrusFileList(u)
	if err != nil {
		return fmt.Errorf("failed to get walrus file list: %w", err)
	}

	return nil
}

func getWalrusFileList(u *url.URL) ([]walrusFileItem, error) {
	switch u.Scheme {
	case "file":
		return local(u.Path)
	default:
		return remote(u.String())
	}
}

// local returns the walrus file list from local file system.
func local(root string) ([]walrusFileItem, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	items := make([]walrusFileItem, 0, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := filepath.Join(root, entry.Name())

		subEntries, err := os.ReadDir(entryPath)
		if err != nil {
			return nil, err
		}

		item := walrusFileItem{
			Name: entry.Name(),
		}

		for _, subEntry := range subEntries {
			fName := subEntry.Name()
			fPath := filepath.Join(entryPath, subEntry.Name())

			switch {
			case slices.Contains(icons, fName):
				item.Icon = fmt.Sprintf("./%s/%s", entry.Name(), fName)
			case fName == fileNameReadme:
				readme, err := os.ReadFile(fPath)
				if err != nil {
					return nil, fmt.Errorf("error reading walrus file content: %w", err)
				}
				item.Readme = string(readme)
			case strs.HasSuffix(fName, yamlFileSuffix...):
				content, err := os.ReadFile(fPath)
				if err != nil {
					return nil, fmt.Errorf("error reading walrus file content: %w", err)
				}
				item.Content = string(content)
			}
		}

		items = append(items, item)
	}

	return items, nil
}

// remote clones the walrus file hub repo and returns the walrus file list.
// It expects URL for the walrus file hub git repo, mainly used for development purpose.
func remote(url string) ([]walrusFileItem, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	walrusFileHubDir = filepath.Join(pwd, ".cache", "walrus-file-hub")

	if _, err = os.Stat(walrusFileHubDir); os.IsNotExist(err) {
		if _, err = vcs.CloneGitRepo(context.Background(), url, walrusFileHubDir, false); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return local(walrusFileHubDir)
}
