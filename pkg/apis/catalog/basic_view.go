package catalog

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/catalog"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.CatalogCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.CatalogOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.CatalogCreateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsValidName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	_, err := url.Parse(r.Source)
	if err != nil {
		return err
	}

	switch r.Type {
	case types.GitDriverGithub, types.GitDriverGitlab, types.GitDriverGitee:
	default:
		return fmt.Errorf("unsupported catalog type %q", r.Type)
	}

	if _, err = regexp.Compile(r.FilterPattern); err != nil {
		return fmt.Errorf("invalid filter pattern: %w", err)
	}

	return nil
}

type (
	GetRequest struct {
		model.CatalogQueryInput `path:",inline"`
	}

	GetResponse = *model.CatalogOutput
)

type UpdateRequest struct {
	model.CatalogUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.CatalogUpdateInput.Validate(); err != nil {
		return err
	}

	if _, err := regexp.Compile(r.FilterPattern); err != nil {
		return fmt.Errorf("invalid filter pattern: %w", err)
	}

	return nil
}

type DeleteRequest = model.CatalogDeleteInput

type (
	CollectionGetRequest struct {
		model.CatalogQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Catalog, catalog.OrderOption,
		] `query:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.CatalogOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest = model.CatalogDeleteInputs
