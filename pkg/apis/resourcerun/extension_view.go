package resourcerun

import (
	"errors"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/deployer/terraform"
	"github.com/seal-io/walrus/utils/json"
)

type RunDiff struct {
	TemplateName       string          `json:"templateName"`
	TemplateVersion    string          `json:"templateVersion"`
	Attributes         property.Values `json:"attributes"`
	ComputedAttributes property.Values `json:"computedAttributes"`
}

type (
	RouteGetTerraformStatesRequest struct {
		_ struct{} `route:"GET=/terraform-states"`

		model.ResourceRunQueryInput `path:",inline"`
	}

	RouteGetTerraformStatesResponse = json.RawMessage
)

type RouteUpdateTerraformStatesRequest struct {
	_ struct{} `route:"PUT=/terraform-states"`

	model.ResourceRunQueryInput `path:",inline"`

	json.RawMessage `path:"-" json:",inline"`
}

type RouteLogRequest struct {
	_ struct{} `route:"GET=/log"`

	model.ResourceRunQueryInput `path:",inline"`

	JobType string `query:"jobType,omitempty"`

	Stream *runtime.RequestUnidiStream
}

func (r *RouteLogRequest) Validate() error {
	if err := r.ResourceRunQueryInput.Validate(); err != nil {
		return err
	}

	if r.JobType == "" {
		r.JobType = terraform.JobTypeApply
	}

	if r.JobType != terraform.JobTypeApply && r.JobType != terraform.JobTypeDestroy {
		return errors.New("invalid job type")
	}

	return nil
}

func (r *RouteLogRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type (
	RouteGetDiffLatestRequest struct {
		_ struct{} `route:"GET=/diff-latest"`

		model.ResourceRunQueryInput `path:",inline"`
	}

	RouteGetDiffLatestResponse struct {
		Old RunDiff `json:"old"`
		New RunDiff `json:"new"`
	}
)

type (
	RouteGetDiffPreviousRequest struct {
		_ struct{} `route:"GET=/diff-previous"`

		model.ResourceRunQueryInput `path:",inline"`
	}

	RouteGetDiffPreviousResponse struct {
		Old RunDiff `json:"old"`
		New RunDiff `json:"new"`
	}
)
