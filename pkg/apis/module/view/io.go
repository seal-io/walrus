package view

import (
	"fmt"

	"github.com/hashicorp/go-getter"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs

type CreateRequest struct {
	*model.Module `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	return validate(r.Module)
}

type CreateResponse = *model.Module

type UpdateRequest struct {
	UriID string `uri:"id"`

	*model.Module `json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	r.ID = r.UriID

	return validate(r.Module)
}

func validate(m *model.Module) error {
	if err := validation.IsQualifiedName(m.ID); err != nil {
		return err
	}
	if _, err := getter.Detect(m.Source, "", getter.Detectors); err != nil {
		return fmt.Errorf("invalid source: %w", err)
	}
	return nil
}

type GetRequest struct {
	ID string `uri:"id"`
}

func (r *GetRequest) Validate() error {
	return validation.IsQualifiedName(r.ID)
}

type GetResponse = *model.Module

type DeleteRequest = GetRequest

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	Group string `query:"_group,omitempty"`
}

type CollectionGetResponse = []GetResponse

// Extensional APIs

type RefreshRequest struct {
	_ struct{} `route:"POST=/refresh"`

	ID string `uri:"id"`
}

func (r *RefreshRequest) Validate() error {
	return validation.IsQualifiedName(r.ID)
}
