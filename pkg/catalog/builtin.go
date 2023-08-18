package catalog

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
)

// BuiltinCatalog returns the seal builtin catalog.
func BuiltinCatalog() *model.Catalog {
	return &model.Catalog{
		Name:        "builtin",
		Description: "Walrus Builtin Catalog",
		Type:        types.GitDriverGithub,
		Source:      "https://catalog.seal.io/walrus",
	}
}
