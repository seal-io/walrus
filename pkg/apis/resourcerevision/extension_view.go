package resourcerevision

import (
	"errors"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/types"
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

		model.ResourceRevisionQueryInput `path:",inline"`
	}

	RouteGetTerraformStatesResponse = json.RawMessage
)

type RouteUpdateTerraformStatesRequest struct {
	_ struct{} `route:"PUT=/terraform-states"`

	model.ResourceRevisionQueryInput `path:",inline"`

	json.RawMessage `path:"-" json:",inline"`
}

type RouteLogRequest struct {
	_ struct{} `route:"GET=/log"`

	model.ResourceRevisionQueryInput `path:",inline"`

	JobType string `query:"jobType,omitempty"`

	Stream *runtime.RequestUnidiStream
}

func (r *RouteLogRequest) Validate() error {
	if err := r.ResourceRevisionQueryInput.Validate(); err != nil {
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

		model.ResourceRevisionQueryInput `path:",inline"`
	}

	RouteGetDiffLatestResponse struct {
		Old RevisionDiff `json:"old"`
		New RevisionDiff `json:"new"`
	}
)

type (
	RouteGetDiffPreviousRequest struct {
		_ struct{} `route:"GET=/diff-previous"`

		model.ResourceRevisionQueryInput `path:",inline"`
	}

	RouteGetDiffPreviousResponse struct {
		Old RevisionDiff `json:"old"`
		New RevisionDiff `json:"new"`
	}
)

type RouteUpdateDriftRequest struct {
	_ struct{} `route:"PUT=/drift"`

	model.ResourceRevisionQueryInput `path:",inline"`

	json.RawMessage `path:"-" json:",inline"`
}

func (r *RouteUpdateDriftRequest) Validate() error {
	if err := r.ResourceRevisionQueryInput.Validate(); err != nil {
		return err
	}

	revision, err := r.Client.ResourceRevisions().Query().
		Select(resourcerevision.FieldID, resourcerevision.FieldType).
		Where(resourcerevision.ID(r.ID)).
		Only(r.Context)
	if err != nil {
		return err
	}

	switch revision.Type {
	case types.ResourceRevisionTypeDetect, types.ResourceRevisionTypeSync:
	default:
		return errors.New("invalid revision type")
	}

	return nil
}
