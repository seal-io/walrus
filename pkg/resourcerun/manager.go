package resourcerun

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/auths/session"
	runbus "github.com/seal-io/walrus/pkg/bus/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/pkg/terraform/parser"
)

type CreateOptions struct {
	// ResourceID is the ID of the resource.
	ResourceID object.ID

	// DeployerType is the type of the deployer that run uses.
	// required: true
	DeployerType string

	// JobType is the type of the job, apply or destroy.
	JobType string

	// ChangeComment is the comment of the change.
	ChangeComment string
}

// Create creates a resource run.
func Create(ctx context.Context, mc model.ClientSet, opts CreateOptions) (*model.ResourceRun, error) {
	// Validate if there is a running run.
	prevEntity, err := mc.ResourceRuns().Query().
		Where(resourcerun.And(
			resourcerun.ResourceID(opts.ResourceID),
			resourcerun.DeployerType(opts.DeployerType))).
		Order(model.Desc(resourcerun.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if prevEntity != nil && status.ResourceRunStatusReady.IsUnknown(prevEntity) {
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
			pq.Select(project.FieldName, project.FieldLabels)
		}).
		WithEnvironment(func(env *model.EnvironmentQuery) {
			env.Select(environment.FieldLabels)
			env.Select(environment.FieldName)
			env.Select(environment.FieldType)
		}).
		WithResourceDefinitionMatchingRule(func(mrq *model.ResourceDefinitionMatchingRuleQuery) {
			mrq.Select(
				resourcedefinitionmatchingrule.FieldName,
				resourcedefinitionmatchingrule.FieldAttributes,
			).
				WithTemplate(func(tvq *model.TemplateVersionQuery) {
					tvq.Select(
						templateversion.FieldID,
						templateversion.FieldVersion,
						templateversion.FieldName,
					)
				})
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	var (
		templateID                    object.ID
		templateName, templateVersion string
		attributes                    = res.Attributes
		computedAttributes            = res.ComputedAttributes
	)

	switch {
	case res.TemplateID != nil:
		templateID = res.Edges.Template.TemplateID
		templateName = res.Edges.Template.Name
		templateVersion = res.Edges.Template.Version
	case res.ResourceDefinitionMatchingRuleID != nil:
		rule := res.Edges.ResourceDefinitionMatchingRule

		templateName = rule.Edges.Template.Name
		templateVersion = rule.Edges.Template.Version

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
	default:
		return nil, errors.New("missing template or resource definition")
	}

	var subjectID object.ID

	s, _ := session.GetSubject(ctx)
	if s.ID != "" {
		subjectID = s.ID
	} else {
		subjectID, err = pkgresource.GetSubjectID(res)
		if err != nil {
			return nil, err
		}
	}

	userSubject, err := mc.Subjects().Query().
		Where(subject.ID(subjectID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	entity := &model.ResourceRun{
		ProjectID:          res.ProjectID,
		EnvironmentID:      res.EnvironmentID,
		ResourceID:         res.ID,
		TemplateID:         templateID,
		TemplateName:       templateName,
		TemplateVersion:    templateVersion,
		Attributes:         attributes,
		ComputedAttributes: computedAttributes,
		DeployerType:       opts.DeployerType,
		CreatedBy:          userSubject.Name,
		ChangeComment:      opts.ChangeComment,
	}

	status.ResourceRunStatusReady.Unknown(entity, "")
	entity.Status.SetSummary(status.WalkResourceRun(&entity.Status))

	// Inherit the output of previous run to create a new one.
	if prevEntity != nil {
		entity.Output = prevEntity.Output
	}

	switch {
	case opts.JobType == types.RunJobTypeApply && entity.Output != "":
		// Get required providers from the previous output after first deployment.
		requiredProviders, err := getRequiredProviders(ctx, mc, opts.ResourceID, entity.Output)
		if err != nil {
			return nil, err
		}
		entity.PreviousRequiredProviders = requiredProviders
	case opts.JobType == types.RunJobTypeDestroy && entity.Output != "":
		if status.ResourceRunStatusReady.IsFalse(prevEntity) {
			// Get required providers from the previous output after first deployment.
			requiredProviders, err := getRequiredProviders(ctx, mc, opts.ResourceID, entity.Output)
			if err != nil {
				return nil, err
			}
			entity.PreviousRequiredProviders = requiredProviders
		} else {
			// Copy required providers from the previous run.
			entity.PreviousRequiredProviders = prevEntity.PreviousRequiredProviders
			// Reuse other fields from the previous run.
			entity.TemplateID = prevEntity.TemplateID
			entity.TemplateName = prevEntity.TemplateName
			entity.TemplateVersion = prevEntity.TemplateVersion
			entity.Attributes = prevEntity.Attributes
			entity.ComputedAttributes = prevEntity.ComputedAttributes
			entity.InputConfigs = prevEntity.InputConfigs
		}
	}

	// Create run.
	entity, err = mc.ResourceRuns().Create().
		Set(entity).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func getRequiredProviders(
	ctx context.Context,
	mc model.ClientSet,
	instanceID object.ID,
	previousOutput string,
) ([]types.ProviderRequirement, error) {
	stateRequiredProviderSet := sets.NewString()

	previousRequiredProviders, err := getPreviousRequiredProviders(ctx, mc, instanceID)
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

// getPreviousRequiredProviders get previous succeed run required providers.
// NB(alex): the previous run may be failed, the failed run may not contain required providers of states.
func getPreviousRequiredProviders(
	ctx context.Context,
	mc model.ClientSet,
	resourceID object.ID,
) ([]types.ProviderRequirement, error) {
	prevRequiredProviders := make([]types.ProviderRequirement, 0)

	entity, err := mc.ResourceRuns().Query().
		Where(resourcerun.ResourceID(resourceID)).
		Order(model.Desc(resourcerun.FieldCreateTime)).
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

// UpdateStatus updates the status of the resource run.
func UpdateStatus(ctx context.Context, mc model.ClientSet, run *model.ResourceRun) error {
	if run == nil {
		return nil
	}

	// Report to resource run.
	run.Status.SetSummary(status.WalkResourceRun(&run.Status))

	run, err := mc.ResourceRuns().UpdateOne(run).
		SetStatus(run.Status).
		Save(ctx)
	if err != nil {
		return err
	}

	if err = runbus.Notify(ctx, mc, run); err != nil {
		return err
	}

	return nil
}
