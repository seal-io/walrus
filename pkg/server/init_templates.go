package server

import (
	"context"
	"database/sql"
	"errors"

	tmplbus "github.com/seal-io/seal/pkg/bus/template"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/settings"
)

func (r *Server) initTemplates(ctx context.Context, opts initOptions) error {
	ref, err := settings.ServeTemplateRefer.Value(ctx, opts.ModelClient)
	if err != nil {
		return err
	}

	builtin := []*model.Template{
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
			ID:          "rds",
			Description: "Provide a RDS instance of Kubernetes via Bitnami Charts.",
			Source:      "github.com/seal-io/modules//rds?ref=" + ref,
		},
		{
			ID:          "aws-rds",
			Description: "Provide a RDS instance of AWS Cloud.",
			Source:      "github.com/seal-io/modules//aws-rds?ref=" + ref,
		},
		{
			ID:          "alicloud-rds",
			Description: "Provide a RDS instance of Alibaba Cloud.",
			Source:      "github.com/seal-io/modules//alicloud-rds?ref=" + ref,
		},
		{
			ID:          "rds-seeder",
			Description: "Seed any RDS instances for Development or Testing.",
			Source:      "github.com/seal-io/modules//rds-seeder?ref=" + ref,
		},
	}

	creates, err := dao.TemplateCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}

	for i := range creates {
		// Do nothing if the template has been created.
		err = creates[i].
			OnConflictColumns(template.FieldID).
			DoNothing().
			Exec(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// No rows error is reasonable for nothing updating.
				continue
			}

			return err
		}

		err = tmplbus.Notify(ctx, builtin[i])
		if err != nil {
			return err
		}
	}

	return nil
}
