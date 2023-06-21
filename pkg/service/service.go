package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/status"
	deptypes "github.com/seal-io/seal/pkg/deployer/types"
	"github.com/seal-io/seal/utils/log"
)

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
		return fmt.Errorf("service status is not deploying, service: %s", entity.ID)
	}

	if dp == nil {
		return errors.New("deployer is not set")
	}

	destroyOpts := deptypes.DestroyOptions{
		SkipTLSVerify: !opts.TlsCertified,
	}

	return dp.Destroy(ctx, entity, destroyOpts)
}
