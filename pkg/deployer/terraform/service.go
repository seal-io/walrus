package terraform

import (
	"context"
	"time"

	revisionbus "github.com/seal-io/seal/pkg/bus/servicerevision"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/terraform/parser"
	"github.com/seal-io/seal/utils/strs"
)

// SyncServiceStatus sync service status with service revision bus message.
func SyncServiceStatus(ctx context.Context, bm revisionbus.BusMessage) (err error) {
	var (
		mc       = bm.TransactionalModelClient
		revision = bm.Refer
	)

	// Report to service.
	service, err := mc.Services().Query().
		Where(service.ID(revision.ServiceID)).
		Only(ctx)
	if err != nil {
		return err
	}

	var serviceStatusUpdate *model.ServiceUpdateOne

	switch revision.Status {
	case status.ServiceRevisionStatusSucceeded:
		if status.ServiceStatusDeleted.IsUnknown(service) {
			err = mc.Services().DeleteOne(service).
				Exec(ctx)
		} else {
			status.ServiceStatusDeployed.True(service, "")
			status.ServiceStatusReady.Unknown(service, "")

			err = updateServiceDriftResult(ctx, mc, service, revision)
		}
	case status.ServiceRevisionStatusFailed:
		if status.ServiceStatusDeleted.IsUnknown(service) {
			status.ServiceStatusDeleted.False(service, "")
		} else {
			status.ServiceStatusDeployed.False(service, "")
		}
		service.Status.SummaryStatusMessage = revision.StatusMessage

		serviceStatusUpdate, err = dao.ServiceStatusUpdate(mc, service)
		if err != nil {
			return err
		}
		err = serviceStatusUpdate.Exec(ctx)
	}

	return err
}

// updateServiceDriftResult update service and its drift result.
func updateServiceDriftResult(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Service,
	sr *model.ServiceRevision,
) error {
	var serviceDriftResult *types.ServiceDriftResult

	if sr.Type == types.ServiceRevisionTypeDetect {
		serviceDrift, err := parser.ParseDriftOutput(sr.StatusMessage)
		if err != nil {
			return err
		}

		if serviceDrift != nil && len(serviceDrift.ResourceDrifts) > 0 {
			rds := make(map[string]*types.ResourceDrift, len(serviceDrift.ResourceDrifts))
			for _, rd := range serviceDrift.ResourceDrifts {
				rds[strs.Join("/", rd.Type, rd.Name)] = rd
			}

			if err = updateResourceDriftResult(ctx, mc, entity.ID, rds); err != nil {
				return err
			}

			serviceDriftResult.Drifted = true
			serviceDriftResult.Time = time.Now()
			serviceDriftResult.Result = serviceDrift
		}
	}

	if !serviceDriftResult.Drifted {
		if err := resetResourceDriftResult(ctx, mc, entity.ID); err != nil {
			return err
		}
	}

	entity.DriftResult = serviceDriftResult

	serviceUpdate, err := dao.ServiceUpdate(mc, entity)
	if err != nil {
		return err
	}

	return serviceUpdate.Exec(ctx)
}

// updateResourceDriftResult update the drift detection result of the service's resources.
func updateResourceDriftResult(
	ctx context.Context,
	mc model.ClientSet,
	serviceID oid.ID,
	resourceDrifts map[string]*types.ResourceDrift,
) error {
	resources, err := mc.ServiceResources().Query().
		Where(serviceresource.ServiceID(serviceID)).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	for i := range resources {
		r := resources[i]

		key := strs.Join("/", r.Type, r.Name)
		if _, ok := resourceDrifts[key]; !ok {
			continue
		}

		r.DriftResult = &types.ServiceResourceDriftResult{
			Drifted: true,
			Time:    time.Now(),
			Result:  resourceDrifts[key],
		}
	}

	return updateResources(ctx, mc, resources)
}

// resetResourceDriftResult reset the drift detection result of the service's resources.
func resetResourceDriftResult(ctx context.Context, mc model.ClientSet, serviceID oid.ID) error {
	resources, err := mc.ServiceResources().Query().
		Where(serviceresource.ServiceID(serviceID)).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	for _, resource := range resources {
		resource.DriftResult = nil
	}

	return updateResources(ctx, mc, resources)
}

func updateResources(ctx context.Context, mc model.ClientSet, resources model.ServiceResources) error {
	updates, err := dao.ServiceResourceUpdates(mc, resources...)
	if err != nil {
		return err
	}

	for _, update := range updates {
		err = update.Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
