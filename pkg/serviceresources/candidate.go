package serviceresources

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ListCandidatesPageByConnector gets the candidates for Label or State by connector id in pagination.
func ListCandidatesPageByConnector(
	ctx context.Context,
	modelClient model.ClientSet,
	connectorID oid.ID,
	offset,
	limit int,
) ([]*model.ServiceResource, error) {
	return queryCandidates(modelClient).
		Where(serviceresource.ConnectorID(connectorID)).
		Offset(offset).
		Limit(limit).
		All(ctx)
}

// ListCandidatesByIDs gets the candidates for Label or State by id list.
func ListCandidatesByIDs(
	ctx context.Context,
	modelClient model.ClientSet,
	ids []oid.ID,
) ([]*model.ServiceResource, error) {
	return queryCandidates(modelClient).
		Where(serviceresource.IDIn(ids...)).
		All(ctx)
}

func queryCandidates(modelClient model.ClientSet) *model.ServiceResourceQuery {
	return modelClient.ServiceResources().Query().
		Order(model.Desc(serviceresource.FieldCreateTime)).
		Unique(false).
		Select(
			serviceresource.FieldID,
			serviceresource.FieldStatus,
			serviceresource.FieldServiceID,
			serviceresource.FieldConnectorID,
			serviceresource.FieldType,
			serviceresource.FieldName,
			serviceresource.FieldDeployerType).
		WithService(func(sq *model.ServiceQuery) {
			sq.Select(
				service.FieldName,
				service.FieldEnvironmentID,
				service.FieldProjectID,
			).WithEnvironment(
				func(eq *model.EnvironmentQuery) {
					eq.Select(
						environment.FieldID,
						environment.FieldName,
					)
				},
			).WithProject(
				func(pq *model.ProjectQuery) {
					pq.Select(
						project.FieldID,
						project.FieldName,
					)
				},
			)
		})
}
