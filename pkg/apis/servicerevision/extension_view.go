package servicerevision

import (
	"errors"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/deployer/terraform"
	"github.com/seal-io/walrus/utils/json"
)

type RevisionDiff struct {
	TemplateName    string          `json:"templateName"`
	TemplateVersion string          `json:"templateVersion"`
	Attributes      property.Values `json:"attributes"`
}

type (
	RouteGetTerraformStatesRequest struct {
		_ struct{} `route:"GET=/terraform-states"`

		model.ServiceRevisionQueryInput `path:",inline"`
	}

	RouteGetTerraformStatesResponse = json.RawMessage
)

type RouteUpdateTerraformStatesRequest struct {
	_ struct{} `route:"PUT=/terraform-states"`

	model.ServiceRevisionQueryInput `path:",inline"`

	json.RawMessage `path:"-" json:",inline"`
}

type RouteLogRequest struct {
	_ struct{} `route:"GET=/log"`

	model.ServiceRevisionQueryInput `path:",inline"`

	JobType string `query:"jobType,omitempty"`

	Stream *runtime.RequestUnidiStream
}

func (r *RouteLogRequest) Validate() error {
	if err := r.ServiceRevisionQueryInput.Validate(); err != nil {
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

		model.ServiceRevisionQueryInput `path:",inline"`
	}

	RouteGetDiffLatestResponse struct {
		Old RevisionDiff `json:"old"`
		New RevisionDiff `json:"new"`
	}
)

type (
	RouteGetDiffPreviousRequest struct {
		_ struct{} `route:"GET=/diff-previous"`

		model.ServiceRevisionQueryInput `path:",inline"`
	}

	RouteGetDiffPreviousResponse struct {
		Old RevisionDiff `json:"old"`
		New RevisionDiff `json:"new"`
	}
)
