package resourcerevision

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/resourcedefinitions"
	"github.com/seal-io/walrus/pkg/terraform/parser"

	revisionbus "github.com/seal-io/walrus/pkg/bus/resourcerevision"
)

const (
	// DeployerType of the revision.
	DeployerType = types.DeployerTypeTF

	JobTypeApply   = "apply"
	JobTypeDestroy = "destroy"

	// _backendAPI the API path to terraform deploy backend.
	// Terraform will get and update deployment states from this API.
	_backendAPI = "/v1/projects/%s/environments/%s/resources/%s/revisions/%s/terraform-states"
)

type Manager struct {
	Plan IPlan
}

func NewManager() *Manager {
	return &Manager{
		Plan: NewPlan(types.DeployerTypeTF),
	}
}

type CreateOptions struct {
	// ResourceID indicates the ID of resource which is for create the revision.
	ResourceID object.ID

	// JobType indicates the type of the job.
	JobType string
}

// Create creates a new resource revision.
// Get the latest revision, and check it if it is running.
// If not running, then apply the latest revision.
// If running, then wait for the latest revision to be applied.
func (m Manager) Create(ctx context.Context, mc model.ClientSet, opts CreateOptions) (*model.ResourceRevision, error) {
	// Validate if there is a running revision.
	prevEntity, err := mc.ResourceRevisions().Query().
		Where(resourcerevision.And(
			resourcerevision.ResourceID(opts.ResourceID),
			resourcerevision.DeployerType(DeployerType))).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if prevEntity != nil && status.ResourceRevisionStatusReady.IsUnknown(prevEntity) {
		return nil, errors.New("deployment is running")
	}

	// Get the corresponding resource and template version.
	res, err := mc.Resources().Query().
		Where(resource.ID(opts.ResourceID)).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldName,
				templateversion.FieldVersion,
				templateversion.FieldTemplateID)
		}).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		WithEnvironment(func(env *model.EnvironmentQuery) {
			env.Select(environment.FieldLabels)
			env.Select(environment.FieldName)
			env.Select(environment.FieldType)
		}).
		WithResourceDefinition(func(rd *model.ResourceDefinitionQuery) {
			rd.Select(resourcedefinition.FieldType)
			rd.WithMatchingRules(func(mrq *model.ResourceDefinitionMatchingRuleQuery) {
				mrq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
					Select(
						resourcedefinitionmatchingrule.FieldName,
						resourcedefinitionmatchingrule.FieldSelector,
						resourcedefinitionmatchingrule.FieldAttributes,
					).
					WithTemplate(func(tvq *model.TemplateVersionQuery) {
						tvq.Select(
							templateversion.FieldID,
							templateversion.FieldVersion,
							templateversion.FieldName,
						)
					})
			})
		}).
		Select(
			resource.FieldID,
			resource.FieldProjectID,
			resource.FieldEnvironmentID,
			resource.FieldType,
			resource.FieldLabels,
			resource.FieldAttributes).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	var (
		templateID                    object.ID
		templateName, templateVersion string
		attributes                    property.Values
	)

	switch {
	case res.TemplateID != nil:
		templateID = res.Edges.Template.TemplateID
		templateName = res.Edges.Template.Name
		templateVersion = res.Edges.Template.Version
		attributes = res.Attributes
	case res.ResourceDefinitionID != nil:
		rd := res.Edges.ResourceDefinition
		matchRule := resourcedefinitions.Match(
			rd.Edges.MatchingRules,
			res.Edges.Project.Name,
			res.Edges.Environment.Name,
			res.Edges.Environment.Type,
			res.Edges.Environment.Labels,
			res.Labels,
		)

		if matchRule == nil {
			return nil, fmt.Errorf("resource definition %s does not match resource %s", rd.Name, res.Name)
		}
		templateName = matchRule.Edges.Template.Name
		templateVersion = matchRule.Edges.Template.Version

		templateID, err = mc.Templates().Query().
			Where(
				template.Name(templateName),
				// Now we only support resource definition globally.
				template.ProjectIDIsNil(),
			).
			OnlyID(ctx)
		if err != nil {
			return nil, err
		}

		// Merge attributes. Resource attributes take precedence over resource definition attributes.
		attributes = matchRule.Attributes
		if attributes == nil {
			attributes = make(property.Values)
		}

		for k, v := range res.Attributes {
			attributes[k] = v
		}
	default:
		return nil, errors.New("service has no template or resource definition")
	}

	entity := &model.ResourceRevision{
		ProjectID:       res.ProjectID,
		EnvironmentID:   res.EnvironmentID,
		ResourceID:      res.ID,
		TemplateID:      templateID,
		TemplateName:    templateName,
		TemplateVersion: templateVersion,
		Attributes:      attributes,
		DeployerType:    DeployerType,
	}

	status.ResourceRevisionStatusReady.Unknown(entity, "")
	entity.Status.SetSummary(status.WalkResourceRevision(&entity.Status))

	// Inherit the output of previous revision to create a new one.
	if prevEntity != nil {
		entity.Output = prevEntity.Output
	}

	switch {
	case opts.JobType == JobTypeApply && entity.Output != "":
		// Get required providers from the previous output after first deployment.
		requiredProviders, err := m.getRequiredProviders(ctx, mc, opts.ResourceID, entity.Output)
		if err != nil {
			return nil, err
		}
		entity.PreviousRequiredProviders = requiredProviders
	case opts.JobType == JobTypeDestroy && entity.Output != "":
		if status.ResourceRevisionStatusReady.IsFalse(prevEntity) {
			// Get required providers from the previous output after first deployment.
			requiredProviders, err := m.getRequiredProviders(ctx, mc, opts.ResourceID, entity.Output)
			if err != nil {
				return nil, err
			}
			entity.PreviousRequiredProviders = requiredProviders
		} else {
			// Copy required providers from the previous revision.
			entity.PreviousRequiredProviders = prevEntity.PreviousRequiredProviders
			// Reuse other fields from the previous revision.
			entity.TemplateID = prevEntity.TemplateID
			entity.TemplateName = prevEntity.TemplateName
			entity.TemplateVersion = prevEntity.TemplateVersion
			entity.Attributes = prevEntity.Attributes
			entity.InputPlan = prevEntity.InputPlan
		}
	}

	// Create revision.
	entity, err = mc.ResourceRevisions().Create().
		Set(entity).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (m Manager) Update(
	ctx context.Context,
	mc model.ClientSet,
	revision *model.ResourceRevision,
) error {
	revision, err := mc.ResourceRevisions().UpdateOne(revision).
		Set(revision).
		Save(ctx)
	if err != nil {
		return err
	}

	return revisionbus.Notify(ctx, mc, revision)
}

func (m Manager) LoadConfigs(
	ctx context.Context,
	mc model.ClientSet,
	opts *PlanOptions,
) (map[string][]byte, error) {
	return m.Plan.LoadConfigs(ctx, mc, opts)
}

func (m Manager) LoadConnectorConfigs(
	connectors model.Connectors,
) (map[string][]byte, error) {
	return m.Plan.LoadConnectorConfigs(connectors)
}

// getRequiredProviders get required providers of the resource.
func (m Manager) getRequiredProviders(
	ctx context.Context,
	mc model.ClientSet,
	instanceID object.ID,
	previousOutput string,
) ([]types.ProviderRequirement, error) {
	stateRequiredProviderSet := sets.NewString()

	previousRequiredProviders, err := m.getPreviousRequiredProviders(ctx, mc, instanceID)
	if err != nil {
		return nil, err
	}

	stateRequiredProviders, err := parser.ParseStateProviders(previousOutput)
	if err != nil {
		return nil, err
	}

	stateRequiredProviderSet.Insert(stateRequiredProviders...)

	requiredProviders := make([]types.ProviderRequirement, 0, len(previousRequiredProviders))

	for _, p := range previousRequiredProviders {
		if stateRequiredProviderSet.Has(p.Name) {
			requiredProviders = append(requiredProviders, p)
		}
	}

	return requiredProviders, nil
}

// getPreviousRequiredProviders get previous succeed revision required providers.
// NB(alex): the previous revision may be failed, the failed revision may not contain required providers of states.
func (m Manager) getPreviousRequiredProviders(
	ctx context.Context,
	mc model.ClientSet,
	resourceID object.ID,
) ([]types.ProviderRequirement, error) {
	prevRequiredProviders := make([]types.ProviderRequirement, 0)

	entity, err := mc.ResourceRevisions().Query().
		Where(resourcerevision.ResourceID(resourceID)).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if entity == nil {
		return prevRequiredProviders, nil
	}

	templateVersion, err := mc.TemplateVersions().Query().
		Where(
			templateversion.TemplateID(entity.TemplateID),
			templateversion.Version(entity.TemplateVersion),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if len(templateVersion.Schema.RequiredProviders) != 0 {
		prevRequiredProviders = append(prevRequiredProviders, templateVersion.Schema.RequiredProviders...)
	}

	return prevRequiredProviders, nil
}
