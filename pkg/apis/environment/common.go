package environment

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/auths/session"
	envbus "github.com/seal-io/walrus/pkg/bus/environment"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/validation"
)

func createEnvironment(
	ctx *gin.Context,
	mc model.ClientSet,
	entity *model.Environment,
) (*model.EnvironmentOutput, error) {
	// Validate the creating environment has the same use with subject.
	sj := session.MustGetSubject(ctx)
	if !sj.IsApplicableEnvironmentType(entity.Type) {
		return nil, errorx.HttpErrorf(http.StatusForbidden,
			"cannot create an environment that type not in: %s", sj.ApplicableEnvironmentTypes)
	}

	err := mc.WithTx(ctx, func(tx *model.Tx) (err error) {
		entity, err = tx.Environments().Create().
			Set(entity).
			SaveE(ctx, dao.EnvironmentConnectorsEdgeSave, dao.EnvironmentVariablesEdgeSave)
		if err != nil {
			return err
		}

		// TODO(thxCode): move the following codes into DAO.

		resourceInputs := entity.Edges.Resources

		for _, svc := range resourceInputs {
			if svc == nil {
				return errors.New("invalid input: nil resource")
			}
			svc.ProjectID = entity.ProjectID
			svc.EnvironmentID = entity.ID
		}

		if err = pkgresource.SetSubjectID(ctx, resourceInputs...); err != nil {
			return err
		}

		resources, err := pkgresource.CreateScheduledResources(ctx, tx, resourceInputs)
		if err != nil {
			return err
		}

		entity.Edges.Resources = resources

		return envbus.NotifyIDs(ctx, tx, envbus.EventCreate, entity.ID)
	})
	if err != nil {
		return nil, errorx.Wrap(err, "failed to create environment")
	}

	return model.ExposeEnvironment(entity), nil
}

func validateEnvironmentCreateInput(r model.EnvironmentCreateInput) error {
	if err := r.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if !types.IsEnvironmentType(r.Type) {
		return fmt.Errorf("invalid type: %s", r.Type)
	}

	// Verify connections.
	connectorIDs := make([]object.ID, len(r.Connectors))
	for i, c := range r.Connectors {
		connectorIDs[i] = c.Connector.ID
	}

	if err := validateConnectors(r.Context, r.Client, connectorIDs); err != nil {
		return err
	}

	// Verify resources.
	for i := range r.Resources {
		if r.Resources[i] == nil {
			return errors.New("empty resource")
		}

		if err := validation.IsDNSLabel(r.Resources[i].Name); err != nil {
			return fmt.Errorf("invalid resource name: %w", err)
		}
	}

	// Get template versions.
	tvIDs := make([]object.ID, len(r.Resources))
	for i := range r.Resources {
		tvIDs[i] = r.Resources[i].Template.ID
	}

	tvs, err := r.Client.TemplateVersions().Query().
		Where(templateversion.IDIn(tvIDs...)).
		Select(
			templateversion.FieldID,
			templateversion.FieldName,
			templateversion.FieldVersion,
			templateversion.FieldSchema,
			templateversion.FieldUiSchema).
		All(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get template version: %w", err)
	}

	// Map template version by ID for resource validation.
	tvm := make(map[object.ID]*model.TemplateVersion, len(tvs))
	for i := range tvs {
		tvm[tvs[i].ID] = tvs[i]
	}

	// Verify resource's variables with variables schema that defined on the template version.
	for _, res := range r.Resources {
		err = res.Attributes.ValidateWith(tvm[res.Template.ID].Schema.VariableSchema())
		if err != nil {
			return fmt.Errorf("invalid variables: %w", err)
		}
	}

	// Verify variables.
	for i := range r.Variables {
		if r.Variables[i] == nil {
			return errors.New("empty variable")
		}

		if err := validation.IsDNSLabel(r.Variables[i].Name); err != nil {
			return fmt.Errorf("invalid variable name: %w", err)
		}
	}

	return nil
}

// validateConnectors checks if given connector IDs are valid within the same project or globally.
func validateConnectors(ctx context.Context, mc model.ClientSet, ids []object.ID) error {
	if len(ids) == 0 {
		return nil
	}

	var typeCount []struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
	}

	err := mc.Connectors().Query().
		Where(connector.IDIn(ids...)).
		GroupBy(connector.FieldType).
		Aggregate(model.Count()).
		Scan(ctx, &typeCount)
	if err != nil {
		return fmt.Errorf("failed to get connector type count: %w", err)
	}

	// Validate connector type is duplicated,
	// only one connector type is allowed in one environment.
	for _, c := range typeCount {
		if c.Count > 1 {
			return fmt.Errorf("invalid connectors: duplicated connector type %s", c.Type)
		}
	}

	return nil
}
