package applicationresources

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ListCandidatesByPage gets the candidates for Label or State by pagination params.
func ListCandidatesByPage(ctx context.Context, modelClient model.ClientSet, offset, limit int) ([]*model.ApplicationResource, error) {
	return queryCandidates(modelClient).
		Offset(offset).
		Limit(limit).
		All(ctx)
}

// ListCandidatesByIDs gets the candidates for Label or State by id list.
func ListCandidatesByIDs(ctx context.Context, modelClient model.ClientSet, ids []types.ID) ([]*model.ApplicationResource, error) {
	return queryCandidates(modelClient).
		Where(applicationresource.IDIn(ids...)).
		All(ctx)
}

func queryCandidates(modelClient model.ClientSet) *model.ApplicationResourceQuery {
	return modelClient.ApplicationResources().Query().
		Order(model.Desc(applicationresource.FieldCreateTime)).
		Unique(false).
		Select(
			applicationresource.FieldID,
			applicationresource.FieldStatus,
			applicationresource.FieldInstanceID,
			applicationresource.FieldConnectorID,
			applicationresource.FieldType,
			applicationresource.FieldName,
			applicationresource.FieldDeployerType).
		WithInstance(func(aiq *model.ApplicationInstanceQuery) {
			aiq.Select(
				applicationinstance.FieldApplicationID,
				applicationinstance.FieldEnvironmentID,
			).WithEnvironment(
				func(eq *model.EnvironmentQuery) {
					eq.Select(
						environment.FieldID,
						environment.FieldName,
					)
				},
			).WithApplication(
				func(aq *model.ApplicationQuery) {
					aq.Select(
						application.FieldID,
						application.FieldName,
						application.FieldProjectID,
					).WithProject(func(pq *model.ProjectQuery) {
						pq.Select(
							project.FieldID,
							project.FieldName,
						)
					})
				},
			)
		})
}
