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
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
)

func ApplyLabels(ctx context.Context, modelClient model.ClientSet, offset, limit int, appResourceIDs []types.ID) func() error {
	return func() (berr error) {
		filters := []predicate.ApplicationResource{
			applicationresource.TypeIn(intercept.TFLabeledTypes...),
		}
		if len(appResourceIDs) != 0 {
			filters = append(filters, applicationresource.IDIn(appResourceIDs...))
		}

		var entities, err = modelClient.ApplicationResources().Query().
			Order(model.Desc(applicationresource.FieldCreateTime)).
			Offset(offset).
			Limit(limit).
			Unique(false).
			Where(
				filters...,
			).
			Select(
				applicationresource.FieldID,
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
			}).
			All(ctx)
		if err != nil {
			return err
		}

		var connAppRes = make(map[*model.Connector][]*model.ApplicationResource)
		for i := 0; i < len(entities); i++ {
			conn := entities[i].Edges.Connector
			connAppRes[conn] = append(connAppRes[conn], entities[i])
		}

		for conn, ress := range connAppRes {
			var op, err = platform.GetOperator(ctx, operator.CreateOptions{
				Connector: *conn,
			})
			if multierr.AppendInto(&berr, err) {
				continue
			}

			for _, res := range ress {
				var (
					appName     string
					projectName string
					envName     string
				)
				instance := res.Edges.Instance
				if instance == nil {
					continue
				}

				// application
				app := instance.Edges.Application
				if app != nil {
					appName = app.Name

					// project
					if app.Edges.Project != nil {
						projectName = app.Edges.Project.Name
					}
				}

				// environment
				if instance.Edges.Environment != nil {
					envName = instance.Edges.Environment.Name
				}

				labels := map[string]string{
					types.LabelSealEnvironment: envName,
					types.LabelSealProject:     projectName,
					types.LabelSealApplication: appName,
				}
				err = op.Label(ctx, res, labels)
				if multierr.AppendInto(&berr, err) {
					continue
				}
			}
		}
		return
	}
}
