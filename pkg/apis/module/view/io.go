package view

import (
	"context"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs

type ModuleRequest struct {
	UriID string `uri:"id"`

	*model.Module `json:",inline"`
}

func (r *ModuleRequest) ValidateWith(ctx context.Context, input any) error {
	if err := validation.IsQualifiedName(r.UriID); err != nil {
		return err
	}
	r.ID = r.UriID
	return nil
}

type IDRequest struct {
	ID string `uri:"id"`
}

func (r *IDRequest) ValidateWith(ctx context.Context, input any) error {
	return validation.IsQualifiedName(r.ID)
}

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	Group string `query:"_group,omitempty"`
}

type CollectionGetResponse = []*model.Module

// Extensional APIs

type RefreshRequest struct {
	_ struct{} `route:"POST=/refresh"`

	ID string `uri:"id"`
}

func (r *RefreshRequest) ValidateWith(ctx context.Context, input any) error {
	return validation.IsQualifiedName(r.ID)
}
