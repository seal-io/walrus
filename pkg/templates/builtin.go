package templates

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/settings"
)

const (
	TemplateWebService       = "webservice"
	TemplateDeploySourceCode = "deploy-source-code"
	TemplateRDS              = "rds"
	TemplateAWSRDS           = "aws-rds"
	TemplateAliCloudRDS      = "alicloud-rds"
	TemplateRDSSeeder        = "rds-seeder"
)

var BuiltinTemplateNames = []string{
	TemplateWebService,
	TemplateDeploySourceCode,
	TemplateRDS,
	TemplateAWSRDS,
	TemplateAliCloudRDS,
	TemplateRDSSeeder,
}

func BuiltinTemplates(ctx context.Context, modelClient model.ClientSet) ([]*model.Template, error) {
	ref, err := settings.ServeTemplateRefer.Value(ctx, modelClient)
	if err != nil {
		return nil, err
	}

	builtin := []*model.Template{
		{
			ID: TemplateWebService,
			Description: "A long-running, scalable, containerized service that" +
				" have a stable network endpoint to receive external network traffic.",
			Source: "github.com/seal-io/modules//webservice?ref=" + ref,
			Icon: "https://raw.githubusercontent.com/" +
				"opencontainers/artwork/d8ccfe94471a0236b1d4a3f0f90862c4fe5486ce/icons/oci_icon_web.svg",
		},
		{
			ID:          TemplateDeploySourceCode,
			Description: "Build and deploy a container image from source code.",
			Source:      "github.com/seal-io/modules//deploy-source-code?ref=" + ref,
			Icon: "https://raw.githubusercontent.com/" +
				"opencontainers/artwork/d8ccfe94471a0236b1d4a3f0f90862c4fe5486ce/icons/oci_icon_containerimage.svg",
		},
		{
			ID:          TemplateRDS,
			Description: "Provide a RDS instance of Kubernetes via Bitnami Charts.",
			Source:      "github.com/seal-io/modules//rds?ref=" + ref,
		},
		{
			ID:          TemplateAWSRDS,
			Description: "Provide a RDS instance of AWS Cloud.",
			Source:      "github.com/seal-io/modules//aws-rds?ref=" + ref,
		},
		{
			ID:          TemplateAliCloudRDS,
			Description: "Provide a RDS instance of Alibaba Cloud.",
			Source:      "github.com/seal-io/modules//alicloud-rds?ref=" + ref,
		},
		{
			ID:          TemplateRDSSeeder,
			Description: "Seed any RDS instances for Development or Testing.",
			Source:      "github.com/seal-io/modules//rds-seeder?ref=" + ref,
		},
	}

	return builtin, nil
}
