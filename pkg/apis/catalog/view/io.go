package view

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/catalog"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// Basic APIs.

type CreateRequest struct {
	model.CatalogCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	_, err := url.Parse(r.Source)
	if err != nil {
		return err
	}

	switch r.Type {
	case types.GitDriverGithub, types.GitDriverGitlab:
	default:
		return fmt.Errorf("unsupported catalog type %q", r.Type)
	}

	return nil
}

type CreateResponse = *model.CatalogOutput

type UpdateRequest struct {
	model.CatalogUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if !r.ID.Valid() {
		return fmt.Errorf("invalid id: blank")
	}

	_, err := url.Parse(r.Source)
	if err != nil {
		return err
	}

	return nil
}

type GetRequest struct {
	model.CatalogQueryInput `json:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid() {
		return fmt.Errorf("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.CatalogOutput

type DeleteRequest struct {
	model.CatalogQueryInput `json:",inline"`
}

func (r *DeleteRequest) Validate() error {
	if !r.ID.Valid() {
		return fmt.Errorf("invalid id: blank")
	}

	return nil
}

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Catalog, catalog.OrderOption] `query:",inline"`
}

type CollectionGetResponse = []*model.CatalogOutput

type CollectionDeleteRequest []*model.CatalogQueryInput

func (r CollectionDeleteRequest) Validate() error {
	if len(r) == 0 {
		return errors.New("invalid input: empty")
	}

	for _, i := range r {
		if !i.ID.Valid() {
			return errors.New("invalid id: blank")
		}
	}

	return nil
}

// Extensional APIs.

type RouteSyncCatalogRequest struct {
	_ struct{} `route:"POST=/refresh"`

	ID object.ID `uri:"id"`
}

func (r *RouteSyncCatalogRequest) Validate() error {
	if !r.ID.Valid() {
		return fmt.Errorf("invalid id: blank")
	}

	return nil
}
