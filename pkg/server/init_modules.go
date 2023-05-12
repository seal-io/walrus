package server

import (
	"context"
	"database/sql"
	"errors"

	modbus "github.com/seal-io/seal/pkg/bus/module"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/settings"
)

func (r *Server) initModules(ctx context.Context, opts initOptions) error {
	ref, err := settings.ServeModuleRefer.Value(ctx, opts.ModelClient)
	if err != nil {
		return err
	}

	builtin := []*model.Module{
		{
			ID: "webservice",
			Description: "A long-running, scalable, containerized service that" +
				" have a stable network endpoint to receive external network traffic.",
			Source: "github.com/seal-io/modules//webservice?ref=" + ref,
			Icon: "https://raw.githubusercontent.com/" +
				"opencontainers/artwork/d8ccfe94471a0236b1d4a3f0f90862c4fe5486ce/icons/oci_icon_web.svg",
		},
		{
			ID:          "build-container-image",
			Description: "Build a container image from source code.",
			Source:      "github.com/seal-io/modules//build-container-image?ref=" + ref,
			Icon: "https://raw.githubusercontent.com/" +
				"opencontainers/artwork/d8ccfe94471a0236b1d4a3f0f90862c4fe5486ce/icons/oci_icon_containerimage.svg",
		},
		{
			ID:          "aws-rds",
			Description: "An AWS RDS instance.",
			Source:      "github.com/seal-io/modules//aws-rds?ref=" + ref,
			Icon: "https://raw.githubusercontent.com/sashee/aws-svg-icons/" +
				"ddf2928b65d8f18c20c6a792740ec934804e7a25/docs/" +
				"Architecture-Service-Icons_07302021/Arch_Database/64/Arch_Amazon-RDS_64.svg",
		},
		{
			ID:          "mysql",
			Description: "A containerized mysql instance.",
			Source:      "github.com/seal-io/modules//mysql?ref=" + ref,
			Icon:        "https://www.mysql.com/common/logos/logo-mysql-170x115.png",
		},
	}

	creates, err := dao.ModuleCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}

	for i := range creates {
		// Do nothing if the module has been created.
		err = creates[i].
			OnConflictColumns(module.FieldID).
			DoNothing().
			Exec(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// No rows error is reasonable for nothing updating.
				continue
			}

			return err
		}

		err = modbus.Notify(ctx, builtin[i])
		if err != nil {
			return err
		}
	}

	return nil
}
