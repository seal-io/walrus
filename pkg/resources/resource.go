package resources

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	resstatus "github.com/seal-io/walrus/pkg/resources/status"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/strs"
)

func GetSubjectID(entity *model.Resource) (object.ID, error) {
	if entity == nil {
		return "", fmt.Errorf("resource is nil")
	}

	subjectIDStr := entity.Annotations[types.AnnotationSubjectID]

	return object.ID(subjectIDStr), nil
}

func SetSubjectID(ctx context.Context, resources ...*model.Resource) error {
	sj, err := session.GetSubject(ctx)
	if err != nil {
		return err
	}

	for i := range resources {
		if resources[i].Annotations == nil {
			resources[i].Annotations = make(map[string]string)
		}
		resources[i].Annotations[types.AnnotationSubjectID] = string(sj.ID)
	}

	return nil
}

// UpdateResourceSubjectID updates the subject ID of the resources.
func UpdateResourceSubjectID(ctx context.Context, mc model.ClientSet, resources ...*model.Resource) error {
	if len(resources) == 0 {
		return nil
	}

	if err := SetSubjectID(ctx, resources...); err != nil {
		return err
	}

	for i := range resources {
		res := resources[i]

		err := mc.Resources().UpdateOne(res).
			SetAnnotations(res.Annotations).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	summaryStatusDeploying   = "Deploying"
	summaryStatusProgressing = "Progressing"
)

// CheckDependencyStatus check resource dependencies status is ready to apply.
// Resource with dependencies cannot be applied directly,
// it must wait for dependencies to be ready.
func CheckDependencyStatus(
	ctx context.Context,
	mc model.ClientSet,
	dp deptypes.Deployer,
	entity *model.Resource,
) (bool, error) {
	// Check dependencies.
	dependencies, err := mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceID(entity.ID),
			resourcerelationship.DependencyIDNEQ(entity.ID),
		).
		QueryDependency().
		All(ctx)
	if err != nil {
		return false, err
	}

	if len(dependencies) == 0 {
		return true, nil
	}

	var (
		// Dependencies that are not ready.
		notReadyDependencies = make([]*model.Resource, 0, len(dependencies))
		resIDs               = make([]object.ID, 0, len(dependencies))
		dependencyMap        = make(map[object.ID]*model.Resource, len(dependencies))
	)
	for _, d := range dependencies {
		// If the resource status is ready or deployed, dependency is ready.
		if status.ResourceStatusReady.IsTrue(d) || status.ResourceStatusDeployed.IsTrue(d) {
			continue
		}

		resIDs = append(resIDs, d.ID)
		dependencyMap[d.ID] = d
	}

	latestRuns, err := dao.GetResourcesLatestRuns(ctx, mc, resIDs...)
	if err != nil {
		return false, err
	}

	for i := range latestRuns {
		run := latestRuns[i]
		// If the resource run is not applied, resource is not ready.
		if !status.ResourceRunStatusApplied.IsTrue(run) {
			notReadyDependencies = append(notReadyDependencies, dependencyMap[run.ResourceID])
		}
	}

	if status.ResourceStatusProgressing.IsUnknown(entity) {
		return false, nil
	}

	// Set status to progressing with a dependency message.
	dependencyNames := sets.NewString()
	for _, d := range notReadyDependencies {
		dependencyNames.Insert(d.Name)
	}
	msg := fmt.Sprintf("Waiting for dependent resources to be ready: %s", strs.Join(", ", dependencyNames.List()...))
	status.ResourceStatusProgressing.Reset(entity, msg)

	if err = resstatus.UpdateStatus(ctx, mc, entity); err != nil {
		return false, fmt.Errorf("failed to update resource status: %w", err)
	}

	return false, nil
}

// CheckDependantStatus check resource dependants status is ready to delete or stop.
// Resource with dependants cannot be deleted or stopped directly,
// it must wait for dependants to be deleted or stopped.
func CheckDependantStatus(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Resource,
	actionType string,
) (bool, error) {
	query := mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceIDNEQ(entity.ID),
			resourcerelationship.DependencyID(entity.ID),
		).
		QueryDependency()

	// When stop resource, stopped dependant resource should be exclude.
	if actionType == "stop" {
		query.Where(func(s *sql.Selector) {
			s.Where(sqljson.ValueNEQ(
				resource.FieldStatus,
				status.ResourceStatusStopped.String(),
				sqljson.Path("summaryStatus"),
			))
		})
	}

	// Check dependants.
	dependants, err := query.All(ctx)
	if err != nil {
		return false, err
	}

	if len(dependants) == 0 {
		return true, nil
	}

	dependantsNames := sets.NewString()
	for _, d := range dependants {
		dependantsNames.Insert(d.Name)
	}

	switch actionType {
	case ActionDelete:
		msg := fmt.Sprintf("Waiting for dependants to be deleted: %s", strs.Join(", ", dependantsNames.List()...))
		if !status.ResourceStatusProgressing.IsUnknown(entity) ||
			status.ResourceStatusDeleted.GetMessage(entity) != msg {
			// Mark status to deleting with a dependency message.
			status.ResourceStatusDeleted.Reset(entity, msg)
			status.ResourceStatusProgressing.Unknown(entity, "")

			if err = resstatus.UpdateStatus(ctx, mc, entity); err != nil {
				return false, fmt.Errorf("failed to update resource status: %w", err)
			}
		}
	case ActionStop:
		msg := fmt.Sprintf("Waiting for dependants to be stopped: %s", strs.Join(", ", dependantsNames.List()...))
		if !status.ResourceStatusProgressing.IsUnknown(entity) ||
			status.ResourceStatusStopped.GetMessage(entity) != msg {
			// Mark status to stopping with a dependency message.
			status.ResourceStatusStopped.Reset(entity, "")
			status.ResourceStatusProgressing.Unknown(entity, msg)

			if err = resstatus.UpdateStatus(ctx, mc, entity); err != nil {
				return false, fmt.Errorf("failed to update resource status: %w", err)
			}
		}
	default:
		return false, fmt.Errorf("unsupported action type: %s", actionType)
	}

	return false, nil
}

// IsStoppable tells whether the given resource is stoppable.
func IsStoppable(r *model.Resource) bool {
	if r == nil {
		return false
	}

	if r.Labels[types.LabelResourceStoppable] == "true" ||
		(r.TemplateID != nil && r.Labels[types.LabelResourceStoppable] != "false") {
		return true
	}

	return false
}

// CanBeStopped tells whether the given resource can be stopped.
func CanBeStopped(r *model.Resource) bool {
	return status.ResourceStatusDeployed.IsTrue(r)
}

func SetEnvResourceDefaultLabels(env *model.Environment, r *model.Resource) error {
	if r == nil || env == nil {
		return errorx.Errorf("resource or environment is nil")
	}

	if r.Labels == nil {
		r.Labels = make(map[string]string)
	}

	// Only set default labels if labels stoppable are not set.
	if _, ok := r.Labels[types.LabelResourceStoppable]; ok {
		return nil
	}

	switch env.Type {
	// Dev and staging environments resources are stoppable by default.
	case types.EnvironmentDevelopment, types.EnvironmentStaging:
		r.Labels[types.LabelResourceStoppable] = "true"
	case types.EnvironmentProduction:
		r.Labels[types.LabelResourceStoppable] = "false"
	default:
	}

	return nil
}

// SetDefaultLabels sets default labels for the provided resources.
func SetDefaultLabels(ctx context.Context, mc model.ClientSet, entities ...*model.Resource) error {
	if len(entities) == 0 {
		return nil
	}

	envIDs := make([]object.ID, 0, len(entities))

	for _, entity := range entities {
		envIDs = append(envIDs, entity.EnvironmentID)
	}

	envs, err := mc.Environments().Query().
		Where(environment.IDIn(envIDs...)).
		All(ctx)
	if err != nil {
		return err
	}

	envMap := make(map[object.ID]*model.Environment, len(envs))
	for _, env := range envs {
		envMap[env.ID] = env
	}

	for _, entity := range entities {
		env, ok := envMap[entity.EnvironmentID]
		if !ok {
			return fmt.Errorf("environment %q not found", entity.EnvironmentID)
		}

		if err := SetEnvResourceDefaultLabels(env, entity); err != nil {
			return err
		}
	}

	return nil
}
