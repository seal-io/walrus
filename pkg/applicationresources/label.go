package applicationresources

import (
	"context"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform/operator"
)

// Label applies the labels to the given model.ApplicationResource list with the given operator.Operator.
func Label(ctx context.Context, op operator.Operator, candidates []*model.ApplicationResource) (berr error) {
	if op == nil {
		return
	}

	for i := range candidates {
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
		var err = op.Label(ctx, candidates[i], ls)
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
			).WithEnvironment(func(eq *model.EnvironmentQuery) {
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
		})
}
