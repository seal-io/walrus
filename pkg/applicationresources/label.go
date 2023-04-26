package applicationresources

import (
	"context"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
)

// Label applies the labels to the given model.ApplicationResource according to its connector.
func Label(ctx context.Context, candidates []*model.ApplicationResource) (berr error) {
	for i := range candidates {
		// get operator.
		var op, err = platform.GetOperator(ctx, operator.CreateOptions{
			Connector: *candidates[i].Edges.Connector,
		})
		if multierr.AppendInto(&berr, err) {
			continue
		}

		// get label values.
		var (
			appName     string
			projectName string
			envName     string
		)
		if ins := candidates[i].Edges.Instance; ins == nil {
			continue
		} else {
			// application name
			if app := ins.Edges.Application; app != nil {
				appName = app.Name
				// project name
				if proj := app.Edges.Project; proj != nil {
					projectName = proj.Name
				}
			}
			// environment name
			if env := ins.Edges.Environment; env != nil {
				envName = env.Name
			}
		}

		var ls = map[string]string{
			types.LabelSealEnvironment: envName,
			types.LabelSealProject:     projectName,
			types.LabelSealApplication: appName,
		}
		err = op.Label(ctx, candidates[i], ls)
		if multierr.AppendInto(&berr, err) {
			continue
		}
	}
	return
}

// ListLabelCandidatesByPage gets the candidates for Label by pagination params.
func ListLabelCandidatesByPage(ctx context.Context, modelClient model.ClientSet, offset, limit int) ([]*model.ApplicationResource, error) {
	return queryLabelCandidates(modelClient).
		Offset(offset).
		Limit(limit).
		All(ctx)
}

// ListLabelCandidatesByIDs gets the candidates for Label by id list.
func ListLabelCandidatesByIDs(ctx context.Context, modelClient model.ClientSet, ids []types.ID) ([]*model.ApplicationResource, error) {
	return queryLabelCandidates(modelClient).
		Where(applicationresource.IDIn(ids...)).
		All(ctx)
}

func queryLabelCandidates(modelClient model.ClientSet) *model.ApplicationResourceQuery {
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
			).WithEnvironment(func(
				eq *model.EnvironmentQuery) {
				eq.Select(
					environment.FieldID,
					environment.FieldName,
				)
			}).WithApplication(func(aq *model.ApplicationQuery) {
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
			})
		}).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		})
}
