package template

import (
	"fmt"

	"github.com/hashicorp/go-getter"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/utils/validation"
)

type (
	CreateRequest struct {
		model.TemplateCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.TemplateOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.TemplateCreateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return err
	}

	if _, err := getter.Detect(r.Source, "", getter.Detectors); err != nil {
		return fmt.Errorf("invalid source: %w", err)
	}

	return nil
}

type (
	GetRequest = model.TemplateQueryInput

	GetResponse = *model.TemplateOutput
)

type UpdateRequest struct {
	model.TemplateUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.TemplateUpdateInput.Validate(); err != nil {
		return err
	}

	if _, err := getter.Detect(r.Source, "", getter.Detectors); err != nil {
		return fmt.Errorf("invalid source: %w", err)
	}

	return nil
}

type DeleteRequest = model.TemplateDeleteInput

type (
	CollectionGetRequest struct {
		model.TemplateQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Template, template.OrderOption,
		] `query:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.TemplateOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest = model.TemplateDeleteInputs
