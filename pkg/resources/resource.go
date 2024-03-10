package resources

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/errorx"
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

func UpdateMetadata(ctx context.Context, mc model.ClientSet, entity *model.Resource) error {
	if err := SetSubjectID(ctx, entity); err != nil {
		return err
	}

	if err := SetDefaultLabels(ctx, mc, entity); err != nil {
		return err
	}

	return mc.Resources().UpdateOne(entity).
		SetAnnotations(entity.Annotations).
		SetLabels(entity.Labels).
		SetDescription(entity.Description).
		Exec(ctx)
}
