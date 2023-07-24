package view

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/drone/go-scm/scm"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/pkg/vcs"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs.

type CreateRequest struct {
	model.ConnectorCreateInput `json:",inline"`

	ProjectID   object.ID `query:"projectID,omitempty"`
	ProjectName string    `query:"projectName,omitempty"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if r.Type == "" {
		return errors.New("invalid type: blank")
	}

	switch {
	case r.ProjectID != "":
		if !r.ProjectID.Valid() {
			return errors.New("invalid project id: blank")
		}

		r.Project = &model.ProjectQueryInput{
			ID: r.ProjectID,
		}
	case r.ProjectName != "":
		projectID, err := modelClient.Projects().Query().
			Where(project.Name(r.ProjectName)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get project")
		}

		r.ProjectID = projectID
		r.Project = &model.ProjectQueryInput{
			ID: projectID,
		}
	}

	entity := r.Model()
	if err := validateConnector(ctx, entity); err != nil {
		return err
	}

	return nil
}

type CreateResponse = *model.ConnectorOutput

type DeleteRequest struct {
	model.ConnectorQueryInput `uri:",inline"`

	ProjectID   object.ID `query:"projectID,omitempty"`
	ProjectName string    `query:"projectName,omitempty"`
}

func (r *DeleteRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	switch {
	case r.ProjectID != "":
		if !r.ProjectID.Valid() {
			return errors.New("invalid project id: blank")
		}
	case r.ProjectName != "":
		projectID, err := modelClient.Projects().Query().
			Where(project.Name(r.ProjectName)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get project")
		}

		r.ProjectID = projectID
	}

	// FIXME(thxCode): a workaround to protect general user deleting global connector,
	//   returns a not found error instead of forbidden.
	if r.ProjectID != "" {
		exist, err := modelClient.Connectors().Query().
			Where(
				connector.ID(r.ID),
				connector.ProjectID(r.ProjectID)).
			Exist(ctx)
		if err != nil {
			return err
		}

		if !exist {
			return runtime.Errorc(http.StatusNotFound)
		}
	}

	return nil
}

type UpdateRequest struct {
	model.ConnectorUpdateInput `uri:",inline" json:",inline"`

	ProjectID   object.ID `query:"projectID,omitempty"`
	ProjectName string    `query:"projectName,omitempty"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if r.Type == "" {
		return errors.New("invalid type: blank")
	}

	if r.ConfigData != nil {
		entity := r.Model()
		if err := validateConnector(ctx, entity); err != nil {
			return err
		}
	}

	switch {
	case r.ProjectID != "":
		if !r.ProjectID.Valid() {
			return errors.New("invalid project id: blank")
		}
	case r.ProjectName != "":
		projectID, err := modelClient.Projects().Query().
			Where(project.Name(r.ProjectName)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get project")
		}

		r.ProjectID = projectID
	}

	// FIXME(thxCode): a workaround to protect general user deleting global connector,
	//   returns a not found error instead of forbidden.
	if r.ProjectID != "" {
		exist, err := modelClient.Connectors().Query().
			Where(
				connector.ID(r.ID),
				connector.ProjectID(r.ProjectID)).
			Exist(ctx)
		if err != nil {
			return err
		}

		if !exist {
			return runtime.Errorc(http.StatusNotFound)
		}
	}

	return nil
}

type GetRequest struct {
	model.ConnectorQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.ConnectorOutput

type StreamResponse struct {
	Type       datamessage.EventType    `json:"type"`
	IDs        []object.ID              `json:"ids,omitempty"`
	Collection []*model.ConnectorOutput `json:"collection,omitempty"`
}
type StreamRequest struct {
	ID object.ID `uri:"id"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	modelClient := input.(model.ClientSet)

	exist, err := modelClient.Connectors().
		Query().
		Where(connector.IDEQ(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

// Batch APIs.

type CollectionDeleteRequest struct {
	Items []*model.ConnectorQueryInput `json:"items"`

	ProjectID   object.ID `query:"projectID,omitempty"`
	ProjectName string    `query:"projectName,omitempty"`
}

func (r CollectionDeleteRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if len(r.Items) == 0 {
		return errors.New("invalid input: empty")
	}

	for _, i := range r.Items {
		if !i.ID.Valid() {
			return errors.New("invalid id: blank")
		}
	}

	switch {
	case r.ProjectID != "":
		if !r.ProjectID.Valid() {
			return errors.New("invalid project id: blank")
		}
	case r.ProjectName != "":
		projectID, err := modelClient.Projects().Query().
			Where(project.Name(r.ProjectName)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get project")
		}

		r.ProjectID = projectID
	}

	// FIXME(thxCode): a workaround to protect general user deleting global connector,
	//   returns a not found error instead of forbidden.
	if r.ProjectID != "" {
		ids := make([]object.ID, len(r.Items))
		for i := range r.Items {
			ids[i] = r.Items[i].ID
		}

		cnt, err := modelClient.Connectors().Query().
			Where(
				connector.IDIn(ids...),
				connector.ProjectID(r.ProjectID)).
			Count(ctx)
		if err != nil {
			return err
		}

		if cnt != len(ids) {
			return runtime.Errorc(http.StatusNotFound)
		}
	}

	return nil
}

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Connector, connector.OrderOption] `query:",inline"`

	Category string `query:"category,omitempty"`
	Type     string `query:"type,omitempty"`

	ProjectID   object.ID `query:"projectID,omitempty"`
	ProjectName string    `query:"projectName,omitempty"`
	WithGlobal  bool      `query:"withGlobal,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	if r.Category != "" {
		switch r.Category {
		case types.ConnectorCategoryKubernetes, types.ConnectorCategoryCustom,
			types.ConnectorCategoryVersionControl, types.ConnectorCategoryCloudProvider:
		default:
			return fmt.Errorf("invalid category: %s", r.Category)
		}
	}

	// Query global scope connectors if the given `ProjectID` is empty,
	// otherwise, query project scope connectors.
	modelClient := input.(model.ClientSet)

	switch {
	case r.ProjectID != "":
		if !r.ProjectID.Valid() {
			return errors.New("invalid project id")
		}
	case r.ProjectName != "":
		projectID, err := modelClient.Projects().Query().
			Where(project.Name(r.ProjectName)).
			OnlyID(ctx)
		if err == nil {
			r.ProjectID = projectID
		}
	}

	return nil
}

type CollectionGetResponse = []*model.ConnectorOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`
}

// Extensional APIs.

type ApplyCostToolsRequest struct {
	_ struct{} `route:"POST=/apply-cost-tools"`

	ID object.ID `uri:"id"`
}

func (r *ApplyCostToolsRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	return validateConnectorType(ctx, modelClient, r.ID)
}

type SyncCostDataRequest struct {
	_ struct{} `route:"POST=/sync-cost-data"`

	ID object.ID `uri:"id"`
}

func (r *SyncCostDataRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	return validateConnectorType(ctx, modelClient, r.ID)
}

func validateConnectorType(ctx context.Context, modelClient model.ClientSet, id object.ID) error {
	conn, err := modelClient.Connectors().Query().
		Select(connector.FieldType).
		Where(connector.ID(id)).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get connector")
	}

	if conn.Type != types.ConnectorTypeK8s {
		return errors.New("invalid type: not support")
	}

	return nil
}

type GetRepositoriesRequest struct {
	_ struct{} `route:"GET=/repositories"`

	runtime.RequestCollection[predicate.Connector, connector.OrderOption] `query:",inline"`

	ID object.ID `uri:"id"`
}

type GetRepositoriesResponse = []*scm.Repository

type GetBranchesRequest struct {
	_ struct{} `route:"GET=/repository-branches"`

	runtime.RequestCollection[predicate.Connector, connector.OrderOption] `query:",inline"`

	ID         object.ID `uri:"id"`
	Repository string    `query:"repository"`
}

type GetBranchesResponse = []*scm.Reference

func validateConnector(ctx context.Context, entity *model.Connector) error {
	ops := optypes.CreateOptions{
		Connector: *entity,
	}

	switch entity.Category {
	case types.ConnectorCategoryKubernetes, types.ConnectorCategoryCloudProvider:
		op, err := operator.Get(ctx, ops)
		if err != nil {
			return fmt.Errorf("invalid connector config: %w", err)
		}

		if err = op.IsConnected(ctx); err != nil {
			return fmt.Errorf("unreachable connector: %w", err)
		}
	case types.ConnectorCategoryVersionControl:
		vcsClient, err := vcs.NewClient(entity)
		if err != nil {
			return fmt.Errorf("invalid connector config: %w", err)
		}

		_, _, err = vcsClient.Users.Find(ctx)
		if err != nil {
			return fmt.Errorf("invalid connector: %w", err)
		}
	case types.ConnectorCategoryCustom:

	default:
		return errors.New("invalid connector category")
	}

	if entity.Category != types.ConnectorCategoryKubernetes && entity.EnableFinOps {
		return errors.New("invalid connector: finOps not supported")
	}

	return nil
}
