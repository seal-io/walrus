package terraform

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/servervars"
)

const (
	// _planFileName the file to store the plan of the resource run.
	_planFileName = "plan.out"
	// _jsonPlanFileName the json file to show the plan of the resource run.
	_jsonPlanFileName = "plan.json"

	// _planCommands the commands to get the changes of the resource run.
	_planCommands = "terraform init -no-color && terraform plan %s -no-color -out=plan.out %s" +
		" && terraform show -json plan.out > " + _jsonPlanFileName
	// _applyCommands the commands to apply deployment of the resource run.
	_applyCommands = "terraform init -no-color && terraform apply %s -no-color"
	// _destroyCommands the commands to destroy deployment of the resource run.
	// As destroy planned in plan file, use apply command to execution the plan.
	_destroyCommands = "terraform init -no-color && terraform apply %s -no-color"

	// _planAPI.
	_planAPI = "/v1/projects/%s/environments/%s/resources/%s/runs/%s/plan"
)

func getPlanCommands(run *model.ResourceRun, opts JobCreateOptions) string {
	var (
		destroy string
		varfile = fmt.Sprintf(" -var-file=%s/terraform.tfvars", _secretMountPath)
	)

	if run.Type == types.RunTypeDelete.String() || run.Type == types.RunTypeStop.String() {
		destroy = "-destroy"
	}

	planCommands := fmt.Sprintf(_planCommands, destroy, varfile)

	planAPI := fmt.Sprintf("%s%s", opts.ServerURL,
		fmt.Sprintf(_planAPI,
			run.ProjectID,
			run.EnvironmentID,
			run.ResourceID,
			run.ID))

	planCommands += fmt.Sprintf(
		" && curl -sS -X POST -H \"Content-Type: multipart/form-data\" -H \"Authorization: Bearer $ACCESS_TOKEN\" %s -F jsonplan=@%s -F plan=@%s",
		planAPI,
		_jsonPlanFileName,
		_planFileName,
	)

	if !servervars.TlsCertified.Get() {
		planCommands += " -k"
	}

	return planCommands
}

func getApplyCommands(run *model.ResourceRun, opts JobCreateOptions) string {
	return fmt.Sprintf("%s && %s", getPlanFile(run, opts), fmt.Sprintf(_applyCommands, _planFileName))
}

func getDestroyCommands(run *model.ResourceRun, opts JobCreateOptions) string {
	return fmt.Sprintf("%s && %s", getPlanFile(run, opts), fmt.Sprintf(_destroyCommands, _planFileName))
}

func getPlanFile(run *model.ResourceRun, opts JobCreateOptions) string {
	getPlanAPI := fmt.Sprintf("%s%s", opts.ServerURL,
		fmt.Sprintf(_planAPI, run.ProjectID, run.EnvironmentID, run.ResourceID, run.ID))

	getPlan := fmt.Sprintf(
		"curl -sS -X GET -H \"Authorization: Bearer $ACCESS_TOKEN\" %s -o %s",
		getPlanAPI,
		_planFileName,
	)

	if !servervars.TlsCertified.Get() {
		getPlan += " -k"
	}

	return getPlan
}
