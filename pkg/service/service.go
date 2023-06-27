package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
	deptypes "github.com/seal-io/seal/pkg/deployer/types"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

const annotationSubjectIDName = "seal.io/subject-id"

// Options for deploy or destroy.
type Options struct {
	TlsCertified bool
	Tags         []string
}

func Create(
	ctx context.Context,
	mc model.ClientSet,
	dp deptypes.Deployer,
	entity *model.Service,
	opts Options,
) (*model.ServiceOutput, error) {
	err := mc.WithTx(ctx, func(tx *model.Tx) error {
		creates, err := dao.ServiceCreates(tx, entity)
		if err != nil {
			return err
		}

		entity, err = creates[0].Save(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Deploy service.
	err = Apply(ctx, mc, dp, entity, Options{
		TlsCertified: opts.TlsCertified,
		Tags:         opts.Tags,
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeService(entity), nil
}

func UpdateStatus(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Service,
) error {
	update, err := dao.ServiceStatusUpdate(mc, entity)
	if err != nil {
		return err
	}

	err = update.Exec(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	return nil
}

func Apply(
	ctx context.Context,
	mc model.ClientSet,
	dp deptypes.Deployer,
	entity *model.Service,
	opts Options,
) (err error) {
	logger := log.WithName("service")

	defer func() {
		if err == nil {
			return
		}
		// Update a failure status.
		status.ServiceStatusDeployed.False(entity, err.Error())

		uerr := UpdateStatus(ctx, mc, entity)
		if uerr != nil {
			logger.Errorf("error updating status of service %s: %v",
				entity.ID, uerr)
		}
	}()

	if !status.ServiceStatusDeployed.IsUnknown(entity) {
		return fmt.Errorf("service status is not deploying, service: %s", entity.ID)
	}

	if dp == nil {
		return errors.New("deployer is not set")
	}

	applyOpts := deptypes.ApplyOptions{
		SkipTLSVerify: !opts.TlsCertified,
		Tags:          opts.Tags,
	}

	return dp.Apply(ctx, entity, applyOpts)
}

func Destroy(
	ctx context.Context,
	mc model.ClientSet,
	dp deptypes.Deployer,
	entity *model.Service,
	opts Options,
) (err error) {
	logger := log.WithName("service")

	defer func() {
		if err == nil {
			return
		}
		// Update a failure status.
		status.ServiceStatusDeleted.False(entity, err.Error())

		uerr := UpdateStatus(ctx, mc, entity)
		if uerr != nil {
			logger.Errorf("error updating status of service %s: %v",
				entity.ID, uerr)
		}
	}()

	if !status.ServiceStatusDeleted.IsUnknown(entity) {
		return fmt.Errorf("service status is not deleting, service: %s", entity.ID)
	}

	if dp == nil {
		return errors.New("deployer is not set")
	}

	// Check dependants.
	dependants, err := dao.GetServiceDependantNames(ctx, mc, entity)
	if err != nil {
		return err
	}

	if len(dependants) > 0 {
		msg := fmt.Sprintf("Waiting for dependants to be deleted: %s", strs.Join(", ", dependants...))
		if !status.ServiceStatusProgressing.IsUnknown(entity) ||
			status.ServiceStatusProgressing.GetMessage(entity) != msg {
			status.ServiceStatusProgressing.Unknown(entity, msg)

			if err = UpdateStatus(ctx, mc, entity); err != nil {
				return fmt.Errorf("failed to update service status: %w", err)
			}
		}

		return nil
	} else {
		status.ServiceStatusProgressing.True(entity, "Resolved dependencies")

		if err = UpdateStatus(ctx, mc, entity); err != nil {
			return fmt.Errorf("failed to update service status: %w", err)
		}
	}

	destroyOpts := deptypes.DestroyOptions{
		SkipTLSVerify: !opts.TlsCertified,
	}

	return dp.Destroy(ctx, entity, destroyOpts)
}

func GetSubjectID(entity *model.Service) (oid.ID, error) {
	if entity == nil {
		return "", fmt.Errorf("service is nil")
	}

	subjectIDStr := entity.Annotations[annotationSubjectIDName]

	return oid.ID(subjectIDStr), nil
}

// SetServiceStatusScheduled sets the status of the service to scheduled.
func SetServiceStatusScheduled(ctx context.Context, mc model.ClientSet, entity *model.Service) error {
	if entity == nil {
		return fmt.Errorf("service is nil")
	}

	dependencyNames := dao.ServiceRelationshipGetDependencyNames(entity)

	msg := ""
	if len(dependencyNames) > 0 {
		msg = fmt.Sprintf("Waiting for dependent services to be ready: %s", strs.Join(", ", dependencyNames...))
	}

	status.ServiceStatusProgressing.Reset(
		entity,
		msg,
	)
	entity.Status.SetSummary(status.WalkService(&entity.Status))

	subject, err := session.GetSubject(ctx)
	if err != nil {
		return err
	}

	entity.Annotations[annotationSubjectIDName] = string(subject.ID)

	update, err := dao.ServiceUpdate(mc, entity)
	if err != nil {
		return err
	}

	return update.Exec(ctx)
}

// CreateScheduledServices creates scheduled services.
func CreateScheduledServices(ctx context.Context, mc model.ClientSet, entities model.Services) (model.Services, error) {
	results := make(model.Services, 0, len(entities))

	sortedServices, err := TopologicalSortServices(entities)
	if err != nil {
		return nil, err
	}

	for _, entity := range sortedServices {
		creates, err := dao.ServiceCreates(mc, entity)
		if err != nil {
			return nil, err
		}

		entity, err = creates[0].Save(ctx)
		if err != nil {
			return nil, err
		}

		err = SetServiceStatusScheduled(ctx, mc, entity)
		if err != nil {
			return nil, err
		}

		results = append(results, entity)
	}

	return results, nil
}

// IsStatusReady returns true if the service is ready.
func IsStatusReady(entity *model.Service) bool {
	switch entity.Status.SummaryStatus {
	case "Preparing", "Unready", "Ready":
		return true
	}

	return false
}

// IsStatusFalse returns true if the service is in error status.
func IsStatusFalse(entity *model.Service) bool {
	switch entity.Status.SummaryStatus {
	case "DeployFailed", "DeleteFailed":
		return true
	case "Progressing":
		return entity.Status.Error
	}

	return false
}

// IsStatusDeleted returns true if the service is deleted.
func IsStatusDeleted(entity *model.Service) bool {
	switch entity.Status.SummaryStatus {
	case "Deleted", "Deleting":
		return true
	}

	return false
}
