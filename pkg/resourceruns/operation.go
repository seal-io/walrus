package resourceruns

import (
	"context"
	"errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao"
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
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	"github.com/seal-io/walrus/pkg/resourceruns/annotations"
	runjob "github.com/seal-io/walrus/pkg/resourceruns/job"
	runstatus "github.com/seal-io/walrus/pkg/resourceruns/status"
	"github.com/seal-io/walrus/pkg/storage"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

type CreateOptions struct {
	// ResourceID is the ID of the resource.
	ResourceID object.ID

	// DeployerType is the type of the deployer that run uses.
	// +required: true
	DeployerType string

	// RunType the type of the run, create, delete, etc.
	Type types.RunType

	// ChangeComment is the comment of the change.
	ChangeComment string

	// Preview is the run need preview.
	Preview bool
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

	if prevEntity != nil && runstatus.IsStatusRunning(prevEntity) {
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
		WithState().
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

	s, err := session.GetSubject(ctx)
	if err != nil {
		return nil, err
	}

	userSubject, err := mc.Subjects().Query().
		Where(subject.ID(s.ID)).
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
		Type:               opts.Type.String(),
		Preview:            opts.Preview,
	}

	status.ResourceRunStatusPending.Unknown(entity, "")
	entity.Status.SetSummary(status.WalkResourceRun(&entity.Status))

	output := res.Edges.State.Data

	if prevEntity != nil && output != "" {
		switch {
		case opts.Type == types.RunTypeCreate ||
			opts.Type == types.RunTypeUpdate ||
			opts.Type == types.RunTypeStart ||
			opts.Type == types.RunTypeRollback:
			// Get required providers from the previous output after first deployment.
			requiredProviders, err := getRequiredProviders(ctx, mc, opts.ResourceID, output)
			if err != nil {
				return nil, err
			}
			entity.PreviousRequiredProviders = requiredProviders

		case opts.Type == types.RunTypeDelete ||
			opts.Type == types.RunTypeStop:
			if status.ResourceRunStatusApplied.IsFalse(prevEntity) {
				// Get required providers from the previous output after first deployment.
				requiredProviders, err := getRequiredProviders(ctx, mc, opts.ResourceID, output)
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
	}

	// Set subject ID.
	err = annotations.SetSubjectID(ctx, entity)
	if err != nil {
		return nil, err
	}

	// Create run.
	entity, err = mc.ResourceRuns().Create().
		Set(entity).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Cancel the previous planed run if new run is created.
	if prevEntity != nil && prevEntity.Status.SummaryStatus == string(status.ResourceRunStatusPlanned) {
		status.ResourceRunStatusCanceled.Reset(prevEntity, "canceled by new run")
		entity.Status.SetSummary(status.WalkResourceRun(&prevEntity.Status))

		_, err = runstatus.UpdateStatus(ctx, mc, prevEntity)
		if err != nil {
			return nil, err
		}
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

	previousRequiredProviders, err := dao.GetPreviousRequiredProviders(ctx, mc, instanceID)
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

// Apply the resource run in planned status.
func Apply(ctx context.Context, mc model.ClientSet, dp deptypes.Deployer, run *model.ResourceRun) error {
	resourceLatestRuns, err := dao.GetResourcesLatestRuns(ctx, mc, run.ResourceID)
	if err != nil {
		return err
	}

	if len(resourceLatestRuns) == 0 {
		return errorx.Errorf("Latest resource run not found")
	}

	resourceLatestRun := resourceLatestRuns[0]

	if resourceLatestRun.ID != run.ID {
		return errorx.Errorf("Only the latest resource run can be applied")
	}

	if !runstatus.IsStatusPlanned(run) {
		return fmt.Errorf("can not apply run in status: %s", run.Status.SummaryStatus)
	}

	return runjob.PerformRunJob(ctx, mc, dp, run)
}

func CleanPlanFiles(ctx context.Context, mc model.ClientSet, sm *storage.Manager, ids ...object.ID) error {
	logger := log.WithName("resource-run")

	runs, err := mc.ResourceRuns().Query().
		Where(resourcerun.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return err
	}

	gopool.Go(func() {
		subCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		for i := range runs {
			if rerr := sm.DeleteRunPlan(subCtx, runs[i]); err != nil {
				logger.Error(rerr, "failed to delete run plan", "run", runs[i].ID)
			}
		}
	})

	return nil
}
