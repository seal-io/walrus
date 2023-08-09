package view

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-getter"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs.

type CreateRequest struct {
	model.TemplateCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	return validate(r.Model())
}

func (r *CreateRequest) Model() *model.Template {
	entity := r.TemplateCreateInput.Model()

	return entity
}

type CreateResponse = *model.TemplateOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	model.TemplateUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	return validate(r.Model())
}

func validate(m *model.Template) error {
	if err := validation.IsDNSLabel(m.Name); err != nil {
		return err
	}

	if _, err := getter.Detect(m.Source, "", getter.Detectors); err != nil {
		return fmt.Errorf("invalid source: %w", err)
	}

	return nil
}

type GetRequest struct {
	model.TemplateQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.TemplateOutput

type StreamResponse struct {
	Type       datamessage.EventType   `json:"type"`
	IDs        []object.ID             `json:"ids,omitempty"`
	Collection []*model.TemplateOutput `json:"collection,omitempty"`
}

// Batch APIs.

type CollectionDeleteRequest []*model.TemplateQueryInput

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

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Template, template.OrderOption] `query:",inline"`

	CatalogID object.ID `query:"catalogID,omitempty"`
}

func (r *CollectionGetRequest) Validate() error {
	if r.CatalogID != "" && !r.CatalogID.Valid() {
		return errors.New("invalid catalogId: blank")
	}

	return nil
}

type CollectionGetResponse = []*model.TemplateOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`
}

// Extensional APIs.

type RefreshRequest struct {
	_ struct{} `route:"POST=/refresh"`

	ID object.ID `uri:"id"`
}

func (r *RefreshRequest) Validate() error {
	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	return nil
}
