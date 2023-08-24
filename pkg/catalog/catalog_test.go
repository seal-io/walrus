package catalog

import (
	"context"
	"testing"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/version"
)

func TestGetRepos(t *testing.T) {
	cases := model.Catalogs{
		BuiltinCatalog(),
	}

	ctx := context.Background()
	for _, c := range cases {
		repos, err := getRepos(ctx, c, version.GetUserAgent())
		if err != nil {
			t.Fatal(err)
		}

		if len(repos) == 0 {
			t.Fatal("expected at least one repo")
		}
	}
}
