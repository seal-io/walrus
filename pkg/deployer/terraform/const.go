package terraform

import "github.com/seal-io/walrus/pkg/dao/types"

const WalrusContextVariableName = "context"

const (
	JobTypeApply   = types.ResourceRevisionTypeApply
	JobTypeDestroy = types.ResourceRevisionTypeDestory
	JobTypeSync    = types.ResourceRevisionTypeSync
	JobTypeDetect  = types.ResourceRevisionTypeDetect
)

// DeployerType the type of deployer.
const DeployerType = types.DeployerTypeTF

const (
	// _backendAPI the API path to terraform deploy backend.
	// Terraform will get and update deployment states from this API.
	_backendAPI = "/v1/projects/%s/environments/%s/resources/%s/revisions/%s/terraform-states"

	// _driftAPI the API path to update revision drift.
	_driftAPI = "/v1/projects/%s/environments/%s/resources/%s/revisions/%s/drift"

	// _variablePrefix the prefix of the variable name.
	_variablePrefix = "_walrus_var_"

	// _resourcePrefix the prefix of the service output name.
	_resourcePrefix = "_walrus_res_"
)

const (
	// _podName the name of the pod.
	_podName = "deployer"

	// _driftFile the file name of the drift output.
	_driftFile = "plan.json"

	// _resourceRevisionIDLabel pod template label key for resource revision id.
	_resourceRevisionIDLabel = "walrus.seal.io/resource-revision-id"

	// _jobNameFormat the format of job name.
	_jobNameFormat = "tf-job-%s-%s"
	// _jobSecretPrefix the prefix of secret name.
	_jobSecretPrefix = "tf-secret-"
	// _secretMountPath the path to mount the secret.
	_secretMountPath = "/var/terraform/secrets"
	// _workdir the working directory of the job.
	_workdir = "/var/terraform/workspace"
)

const (
	// _applyCommands the commands to apply deployment of the resource.
	_applyCommands = "terraform init -no-color && terraform apply -auto-approve -no-color"
	// _destroyCommands the commands to destroy deployment of the resource.
	_destroyCommands = "terraform init -no-color && terraform destroy -auto-approve -no-color"
	// _detectCommands the commands to detect drift of the revision.
	_detectCommands = "terraform init -no-color && terraform plan -refresh-only -no-color -out=plan.out %s" +
		" && TF_LOG=ERROR terraform show -json plan.out > " + _driftFile
	// _syncCommands the commands to sync state of the resource.
	_syncCommands = "terraform init -no-color && terraform apply -refresh-only -auto-approve -no-color"
)
