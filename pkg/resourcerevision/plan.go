package resourcerevision

import (
	"context"
	"errors"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	pkgenv "github.com/seal-io/walrus/pkg/environment"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
)

// IPlan is the interface of the planner to construct input plan for the revision deployment.
type IPlan interface {
	// Plan plans the revision.
	LoadPlan(context.Context, model.ClientSet, *PlanOptions) ([]byte, error)
	// LoadConfigs loads the plan configs of the plan options.
	LoadConfigs(context.Context, model.ClientSet, *PlanOptions) (map[string][]byte, error)
	// LoadConnectorConfigs loads the connector configs of the plan options.
	// Some connectors may be required to deploy the service.
	LoadConnectorConfigs(model.Connectors) (map[string][]byte, error)
}

// PlanOptions are the options for planning a revision.
type PlanOptions struct {
	// SecretMountPath of the deploy job.
	SecretMountPath string

	ResourceRevision *model.ResourceRevision
	Connectors       model.Connectors
	SubjectID        object.ID
	// Walrus Context.
	Context pkgresource.Context
}

// NewPlan creates a new plan with the plan type.
func NewPlan(planType string) IPlan {
	switch planType {
	case types.DeployerTypeTF:
		return NewTerraformPlan()
	default:
		return nil
	}
}

// GetPlanOptions sets the plan options.
func GetPlanOptions(
	ctx context.Context,
	mc model.ClientSet,
	resourceRevision *model.ResourceRevision,
	secretMountPath string,
) (*PlanOptions, error) {
	opts := &PlanOptions{
		ResourceRevision: resourceRevision,
		SecretMountPath:  secretMountPath,
	}

	if !status.ResourceRevisionStatusReady.IsUnknown(resourceRevision) {
		return nil, errors.New("service revision is not running")
	}

	connectors, err := dao.GetConnectors(ctx, mc, resourceRevision.EnvironmentID)
	if err != nil {
		return nil, err
	}

	proj, err := mc.Projects().Get(ctx, resourceRevision.ProjectID)
	if err != nil {
		return nil, err
	}

	env, err := dao.GetEnvironmentByID(ctx, mc, resourceRevision.EnvironmentID)
	if err != nil {
		return nil, err
	}

	res, err := mc.Resources().Get(ctx, resourceRevision.ResourceID)
	if err != nil {
		return nil, err
	}

	var subjectID object.ID

	sj, _ := session.GetSubject(ctx)
	if sj.ID != "" {
		subjectID = sj.ID
	} else {
		subjectID, err = pkgresource.GetSubjectID(res)
		if err != nil {
			return nil, err
		}
	}

	if subjectID == "" {
		return nil, errors.New("subject id is empty")
	}

	opts.Connectors = connectors
	opts.SubjectID = subjectID

	// Walrus Context.
	opts.Context = *pkgresource.NewContext().
		SetProject(proj.ID, proj.Name).
		SetEnvironment(env.ID, env.Name, pkgenv.GetManagedNamespaceName(env)).
		SetResource(res.ID, res.Name)

	return opts, nil
}
