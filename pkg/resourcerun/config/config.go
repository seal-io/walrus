package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	pkgenv "github.com/seal-io/walrus/pkg/environment"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
)

// Loader is the interface to construct input configs and dependency connectors for the run.
type Loader interface {
	InputLoader
	ProviderLoader
}

// InputLoader is the interface to construct input configs for the run.
type InputLoader interface {
	// LoadMain loads the main config file of the config options.
	LoadMain(context.Context, model.ClientSet, *LoaderOptions) (types.ResourceRunConfigData, error)
	// LoadAll loads the configs files of the config options.
	LoadAll(context.Context, model.ClientSet, *LoaderOptions) (map[string]types.ResourceRunConfigData, error)
}

// ProviderLoader is the interface to construct dependency connectors files for the run.
type ProviderLoader interface {
	// LoadProviders loads the providers of the run required,
	// Some connectors may be required to deploy the service.
	LoadProviders(model.Connectors) (map[string]types.ResourceRunConfigData, error)
}

// LoaderOptions are the options for load a run config files.
type LoaderOptions struct {
	// SecretMountPath of the deployment job.
	SecretMountPath string

	ResourceRun *model.ResourceRun
	Connectors  model.Connectors
	SubjectID   object.ID
	// Walrus Context.
	Context types.Context
}

// NewInputLoader creates a new plan with the plan type.
func NewInputLoader(deployerType string) Loader {
	switch deployerType {
	case types.DeployerTypeTF:
		return NewTerraformLoader()
	default:
		return nil
	}
}

// GetConfigLoaderOptions sets the config loader options.
// It will fetch the resource run, environment, project, resource and subject.
func GetConfigLoaderOptions(
	ctx context.Context,
	mc model.ClientSet,
	run *model.ResourceRun,
	secretMountPath string,
) (*LoaderOptions, error) {
	opts := &LoaderOptions{
		ResourceRun:     run,
		SecretMountPath: secretMountPath,
	}

	if !status.ResourceRunStatusReady.IsUnknown(run) {
		return nil, errors.New("resource run is not running")
	}

	connectors, err := dao.GetConnectors(ctx, mc, run.EnvironmentID)
	if err != nil {
		return nil, err
	}

	proj, err := mc.Projects().Get(ctx, run.ProjectID)
	if err != nil {
		return nil, err
	}

	env, err := dao.GetEnvironmentByID(ctx, mc, run.EnvironmentID)
	if err != nil {
		return nil, err
	}

	res, err := mc.Resources().Get(ctx, run.ResourceID)
	if err != nil {
		return nil, err
	}

	sj, err := getSubject(ctx, mc, res)
	if err != nil {
		return nil, err
	}

	opts.Connectors = connectors
	opts.SubjectID = sj.ID

	// Walrus Context.
	opts.Context = *types.NewContext().
		SetProject(proj.ID, proj.Name).
		SetEnvironment(env.ID, env.Name, pkgenv.GetManagedNamespaceName(env)).
		SetResource(res.ID, res.Name)

	return opts, nil
}

// getSubject gets the subject of the given resource.
func getSubject(ctx context.Context, mc model.ClientSet, res *model.Resource) (*model.Subject, error) {
	var (
		subjectID object.ID
		err       error
	)

	s, _ := session.GetSubject(ctx)
	if s.ID != "" {
		subjectID = s.ID
	} else {
		subjectID, err = pkgresource.GetSubjectID(res)
		if err != nil {
			return nil, err
		}
	}

	if subjectID == "" {
		return nil, fmt.Errorf("subject id is empty")
	}

	return mc.Subjects().Get(ctx, subjectID)
}
