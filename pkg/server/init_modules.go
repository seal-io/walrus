package server

import (
	"context"
	"database/sql"
	"errors"

	modbus "github.com/seal-io/seal/pkg/bus/module"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/module"
)

func (r *Server) initModules(ctx context.Context, opts initOptions) error {
	var builtin = []*model.Module{
		{
			ID:          "webservice",
			Description: "A long-running, scalable, containerized service that have a stable network endpoint to receive external network traffic.",
			Source:      "github.com/seal-io/modules//webservice",
			Icon:        "https://raw.githubusercontent.com/opencontainers/artwork/d8ccfe94471a0236b1d4a3f0f90862c4fe5486ce/icons/oci_icon_web.svg",
		},
		{
			ID:          "build-container-image",
			Description: "Build a container image from source code.",
			Source:      "github.com/seal-io/modules//build-container-image",
			Icon:        "https://raw.githubusercontent.com/opencontainers/artwork/d8ccfe94471a0236b1d4a3f0f90862c4fe5486ce/icons/oci_icon_containerimage.svg",
		},
		{
			ID:          "aws-rds",
			Description: "An AWS RDS instance.",
			Source:      "github.com/seal-io/modules//aws-rds",
			Icon:        "https://raw.githubusercontent.com/sashee/aws-svg-icons/ddf2928b65d8f18c20c6a792740ec934804e7a25/docs/Architecture-Service-Icons_07302021/Arch_Database/64/Arch_Amazon-RDS_64.svg",
		},
		{
			ID:          "mysql",
			Description: "A containerized mysql instance.",
			Source:      "github.com/seal-io/modules//mysql?ref=063c14c6774b359349a0910535b4d60fcd40b810",
			Icon:        "https://www.mysql.com/common/logos/logo-mysql-170x115.png",
		},
	}

	var creates, err = dao.ModuleCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}
	for i := range creates {
		// do nothing if the module has been created.
		err = creates[i].
			OnConflictColumns(module.FieldID).
			DoNothing().
			Exec(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// no rows error is reasonable for nothing updating.
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
