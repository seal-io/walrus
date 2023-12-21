package resource

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
)

// IsService tells if the given resource is of service type.
func IsService(r *model.Resource) bool {
	if r == nil {
		return false
	}

	return r.TemplateID != nil
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

// IsInactive tells whether the given resource is inactive.
func IsInactive(r *model.Resource) bool {
	if r == nil {
		return false
	}

	return r.Status.SummaryStatus == status.ResourceStatusUnDeployed.String() ||
		r.Status.SummaryStatus == status.ResourceStatusStopped.String()
}

// IsStatusReady returns true if the resource is ready.
func IsStatusReady(entity *model.Resource) bool {
	switch entity.Status.SummaryStatus {
	case "Preparing", "NotReady", "Ready":
		return true
	}

	return false
}

// IsStatusFalse returns true if the resource is in error status.
func IsStatusFalse(entity *model.Resource) bool {
	switch entity.Status.SummaryStatus {
	case "DeployFailed", "DeleteFailed":
		return true
	case "Progressing":
		return entity.Status.Error
	}

	return false
}

// IsStatusDeleted returns true if the resource is deleted.
func IsStatusDeleted(entity *model.Resource) bool {
	switch entity.Status.SummaryStatus {
	case "Deleted", "Deleting":
		return true
	}

	return false
}

const (
	summaryStatusDeploying   = "Deploying"
	summaryStatusProgressing = "Progressing"
)

// CheckDependencyStatus check resource dependencies status is ready to apply.
func CheckDependencyStatus(ctx context.Context, mc model.ClientSet, entity *model.Resource) (bool, error) {
	// Check dependants.
	dependencies, err := mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceID(entity.ID),
			resourcerelationship.DependencyIDNEQ(entity.ID),
		).
		QueryDependency().
		Select(resource.FieldID).
		Where(
			resource.Or(
				func(s *sql.Selector) {
					s.Where(sqljson.ValueEQ(
						resource.FieldStatus,
						summaryStatusDeploying,
						sqljson.Path("summaryStatus"),
					))
				},
				resource.And(
					func(s *sql.Selector) {
						s.Where(sqljson.ValueEQ(
							resource.FieldStatus,
							summaryStatusProgressing,
							sqljson.Path("summaryStatus"),
						))
					},
					func(s *sql.Selector) {
						s.Where(sqljson.ValueEQ(
							resource.FieldStatus,
							true,
							sqljson.Path("transitioning"),
						))
					},
				),
			),
		).
		All(ctx)
	if err != nil {
		return false, err
	}

	if len(dependencies) > 0 {
		// If dependency resources is in deploying status.
		err = SetResourceStatusScheduled(ctx, mc, entity)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
