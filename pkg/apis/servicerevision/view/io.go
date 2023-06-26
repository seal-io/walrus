package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/deployer/terraform"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/json"
)

// Basic APIs.

type GetRequest struct {
	model.ServiceRevisionQueryInput `uri:",inline"`

	ProjectID oid.ID `query:"projectID"`
}

func (r *GetRequest) Validate() error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.ServiceRevisionOutput

type StreamRequest struct {
	ID        oid.ID `uri:"id"`
	ProjectID oid.ID `query:"projectID"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	modelClient := input.(model.ClientSet)

	exist, err := modelClient.ServiceRevisions().Query().
		Where(servicerevision.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

type StreamResponse struct {
	Type       datamessage.EventType          `json:"type"`
	IDs        []oid.ID                       `json:"ids,omitempty"`
	Collection []*model.ServiceRevisionOutput `json:"collection,omitempty"`
}

// Batch APIs.

type CollectionDeleteRequest []*model.ServiceRevisionQueryInput

func (r CollectionDeleteRequest) ValidateWith(ctx context.Context, input any) error {
	if len(r) == 0 {
		return errors.New("invalid ids: blank")
	}

	var (
		ids         = make([]oid.ID, 0, len(r))
		modelClient = input.(model.ClientSet)
	)

	for _, i := range r {
		if !i.ID.Valid(0) {
			return errors.New("invalid ids: blank")
		}

		ids = append(ids, i.ID)
	}

	revisions, err := modelClient.ServiceRevisions().Query().
		Select(servicerevision.FieldID, servicerevision.FieldServiceID).
		Where(servicerevision.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get service revisions")
	}

	if len(revisions) != len(r) {
		return errors.New("invalid ids: some revisions are not found")
	}

	serviceID := revisions[0].ServiceID
	for _, revision := range revisions {
		if revision.ServiceID != serviceID {
			return errors.New("invalid ids: revision ids are not from the same service")
		}
	}

	latestRevision, err := modelClient.ServiceRevisions().Query().
		Select(servicerevision.FieldID).
		Where(servicerevision.ServiceID(serviceID)).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		First(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get latest revision")
	}

	for _, revision := range revisions {
		// Prevent deleting the latest revision.
		if revision.ID == latestRevision.ID {
			return errors.New("invalid ids: can not delete latest revision")
		}
	}

	return nil
}

type CollectionGetRequest struct {
	runtime.RequestPagination                           `query:",inline"`
	runtime.RequestExtracting                           `query:",inline"`
	runtime.RequestSorting[servicerevision.OrderOption] `query:",inline"`

	ProjectID       oid.ID `query:"projectID"`
	EnvironmentID   oid.ID `query:"environmentID,omitempty"`
	EnvironmentName string `query:"environmentName,omitempty"`
	ServiceID       oid.ID `query:"serviceID,omitempty"`
	ServiceName     string `query:"serviceName,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	switch {
	case r.ServiceID.Valid(0):
		_, err := modelClient.Services().Query().
			Where(service.ID(r.ServiceID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get service")
		}
	case r.ServiceName != "":
		switch {
		case r.EnvironmentID.Valid(0):
			id, err := modelClient.Services().Query().
				Where(
					service.ProjectID(r.ProjectID),
					service.EnvironmentID(r.EnvironmentID),
					service.Name(r.ServiceName),
				).
				OnlyID(ctx)
			if err != nil {
				return runtime.Errorw(err, "failed to get service by name")
			}

			r.ServiceID = id
		case r.EnvironmentName != "":
			id, err := modelClient.Services().Query().
				Where(
					service.ProjectID(r.ProjectID),
					service.HasEnvironmentWith(environment.Name(r.EnvironmentName)),
					service.Name(r.ServiceName),
				).
				OnlyID(ctx)
			if err != nil {
				return runtime.Errorw(err, "failed to get service by name")
			}

			r.ServiceID = id
		default:
			return errors.New("both environment id and environment name are blank, " +
				"one of them is required while query by service name")
		}
	}

	return nil
}

type CollectionGetResponse = []*model.ServiceRevisionOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`

	ProjectID oid.ID `query:"projectID"`
	ServiceID oid.ID `query:"serviceID,omitempty"`
}

func (r *CollectionStreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if r.ServiceID != "" {
		if !r.ServiceID.IsNaive() {
			return errors.New("invalid service id")
		}

		modelClient := input.(model.ClientSet)

		_, err := modelClient.Services().Query().
			Where(service.ID(r.ServiceID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get service")
		}
	}

	return nil
}

// Extensional APIs.

type GetTerraformStatesRequest = GetRequest

type GetTerraformStatesResponse = json.RawMessage

type UpdateTerraformStatesRequest struct {
	GetRequest `query:",inline" uri:",inline"`

	json.RawMessage `uri:"-" json:",inline"`
}

type StreamLogRequest struct {
	ID        oid.ID `uri:"id"`
	ProjectID oid.ID `query:"projectID"`
	JobType   string `query:"jobType,omitempty"`
}

func (r *StreamLogRequest) Validate() error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	if r.JobType == "" {
		r.JobType = terraform.JobTypeApply
	}

	if r.JobType != terraform.JobTypeApply && r.JobType != terraform.JobTypeDestroy {
		return errors.New("invalid job type")
	}

	return nil
}

type DiffLatestRequest struct {
	ID        oid.ID `uri:"id"`
	ProjectID oid.ID `query:"projectID"`
}

func (r *DiffLatestRequest) Validate() error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type RevisionDiffPreviousRequest struct {
	ID        oid.ID `uri:"id"`
	ProjectID oid.ID `query:"projectID"`
}

func (r *RevisionDiffPreviousRequest) Validate() error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type RevisionDiff struct {
	TemplateID      string          `json:"templateId"`
	TemplateVersion string          `json:"templateVersion"`
	Attributes      property.Values `json:"attributes"`
}

type RevisionDiffResponse struct {
	Old RevisionDiff `json:"old"`
	New RevisionDiff `json:"new"`
}
