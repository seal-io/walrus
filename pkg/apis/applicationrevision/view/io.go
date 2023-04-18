package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformtf"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/json"
)

// Basic APIs

type GetRequest struct {
	*model.ApplicationRevisionQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetResponse = *model.ApplicationRevisionOutput

type StreamResponse struct {
	Type       datamessage.EventType              `json:"type"`
	IDs        []types.ID                         `json:"ids"`
	Collection []*model.ApplicationRevisionOutput `json:"collection"`
}

type StreamRequest struct {
	ID types.ID `uri:"id"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	var modelClient = input.(model.ClientSet)
	exist, err := modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

// Batch APIs

type CollectionDeleteRequest = []*model.ApplicationRevisionQueryInput

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	InstanceID types.ID `query:"instanceID,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if r.InstanceID != "" {
		if !r.InstanceID.IsNaive() {
			return errors.New("invalid instance id")
		}
		_, err := modelClient.ApplicationInstances().Query().
			Where(applicationinstance.ID(r.InstanceID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get application instance")
		}
	}

	return nil
}

type CollectionGetResponse = []*model.ApplicationRevisionOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`

	InstanceID types.ID `query:"instanceID,omitempty"`
}

func (r *CollectionStreamRequest) ValidateWith(ctx context.Context, input any) error {
	if r.InstanceID != "" {
		var modelClient = input.(model.ClientSet)
		if r.InstanceID != "" {
			if !r.InstanceID.IsNaive() {
				return errors.New("invalid instance id")
			}
			_, err := modelClient.ApplicationInstances().Query().
				Where(applicationinstance.ID(r.InstanceID)).
				OnlyID(ctx)
			if err != nil {
				return runtime.Errorw(err, "failed to get application instance")
			}
		}
	}

	return nil
}

// Extensional APIs

type GetTerraformStatesRequest = GetRequest

type GetTerraformStatesResponse = json.RawMessage

type UpdateTerraformStatesRequest struct {
	GetRequest      `uri:",inline"`
	json.RawMessage `json:",inline"`
}

type StreamLogRequest struct {
	ID      types.ID `uri:"id"`
	JobType string   `query:"jobType,omitempty"`
}

func (r *StreamLogRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	if r.JobType == "" {
		r.JobType = platformtf.JobTypeApply
	}
	if r.JobType != platformtf.JobTypeApply && r.JobType != platformtf.JobTypeDestroy {
		return errors.New("invalid job type")
	}

	return nil
}

type RollbackApplicationRequest struct {
	*model.ApplicationRevisionQueryInput `uri:",inline"`
}

func (r *RollbackApplicationRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}
