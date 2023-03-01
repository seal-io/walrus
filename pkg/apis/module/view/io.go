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
	*model.ModuleCreateInput `json:",inline"`

	ID string `json:"id"`
}

func (r *CreateRequest) Validate() error {
	return validate(r.ID, r.Source)
}

func (r *CreateRequest) Model() *model.Module {
	var entity = r.ModuleCreateInput.Model()
	entity.ID = r.ID
	return entity
}

type CreateResponse = *model.ModuleOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	*model.ModuleUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	return validate(r.ID, r.Source)
}

func validate(id string, source string) error {
	if err := validation.IsQualifiedName(id); err != nil {
		return err
	}
	if _, err := getter.Detect(source, "", getter.Detectors); err != nil {
		return fmt.Errorf("invalid source: %w", err)
	}
	return nil
}

type GetRequest struct {
	*model.ModuleQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	return validation.IsQualifiedName(r.ID)
}

type GetResponse = *model.ModuleOutput

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	Group string `query:"_group,omitempty"`
}

type CollectionGetResponse = []*model.ModuleOutput

// Extensional APIs

type RefreshRequest struct {
	_ struct{} `route:"POST=/refresh"`

	ID string `uri:"id"`
}

func (r *RefreshRequest) Validate() error {
	return validation.IsQualifiedName(r.ID)
}
