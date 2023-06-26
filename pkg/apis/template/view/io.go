package view

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/go-getter"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs.

type CreateRequest struct {
	model.TemplateCreateInput `json:",inline"`

	ID string `json:"id"`
}

func (r *CreateRequest) Validate() error {
	return validate(r.Model())
}

func (r *CreateRequest) Model() *model.Template {
	entity := r.TemplateCreateInput.Model()
	entity.ID = r.ID

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
	if err := validation.IsQualifiedName(m.ID); err != nil {
		return err
	}

	if _, err := getter.Detect(m.Source, "", getter.Detectors); err != nil {
		return fmt.Errorf("invalid source: %w", err)
	}

	if m.Icon != "" {
		if _, err := url.ParseRequestURI(m.Icon); err != nil {
			return fmt.Errorf("invalid icon URL: %w", err)
		}
	}

	return nil
}

type GetRequest struct {
	model.TemplateQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	return validation.IsQualifiedName(r.ID)
}

type GetResponse = *model.TemplateOutput

type StreamResponse struct {
	Type       datamessage.EventType   `json:"type"`
	IDs        []string                `json:"ids,omitempty"`
	Collection []*model.TemplateOutput `json:"collection,omitempty"`
}

// Batch APIs.

type CollectionDeleteRequest []*model.TemplateQueryInput

func (r CollectionDeleteRequest) Validate() error {
	if len(r) == 0 {
		return errors.New("invalid input: empty")
	}

	for _, i := range r {
		if err := validation.IsQualifiedName(i.ID); err != nil {
			return fmt.Errorf("invalid id: %w", err)
		}
	}

	return nil
}

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Template, template.OrderOption] `query:",inline"`
}

type CollectionGetResponse = []*model.TemplateOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`
}

// Extensional APIs.

type RefreshRequest struct {
	_ struct{} `route:"POST=/refresh"`

	ID string `uri:"id"`
}

func (r *RefreshRequest) Validate() error {
	return validation.IsQualifiedName(r.ID)
}
