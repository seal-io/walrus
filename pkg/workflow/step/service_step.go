package step

import (
	"context"
	"encoding/json"
	"errors"
	"path"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/workflow/step/types"
)

// nolint:lll
const stepSource = `#!/bin/sh
set -e
set -o pipefail

serverURL="{{workflow.parameters.server}}"
projectID="{{workflow.parameters.projectID}}"
environmentID="{{inputs.parameters.environmentID}}"
token="{{inputs.parameters.token}}"
jobType="create"
resourceName="{{inputs.parameters.resourceName}}"
commonPath="$serverURL/v1/projects/$projectID/environments/$environmentID"

# if skip tls verify
tlsVerify="-k"
if [ "{{workflow.parameters.tlsVerify}}" == "true" ]; then
	tlsVerify=""
fi

common_get() {
    curl -sS --fail-with-body "$1" -X GET -H "Authorization: Bearer $token" $tlsVerify --compressed
}

wait_run_status() {
    local resourceName="$1"
    local runID="$2"
    local breakSummaryStatus="$3"
    local timeout=30
    local factor=1

    while true; do
        statusResponse=$(common_get "$commonPath/resources/$resourceName/runs/$runID")
        statusSummary=$(echo $statusResponse | jq -r '.status.summaryStatus')

        if [ "$statusSummary" == "$breakSummaryStatus" ]; then
            break
        fi

        if [ "$statusSummary" == "Planned" ]; then
            echo "run planned, waiting for running"
            break
        fi

        # If status is failed or canceled, exit.
        if [ "$statusSummary" == "Failed" ] || [ "$statusSummary" == "Canceled" ]; then
            echo "run failed or was canceled"
            exit 1
        fi

        if [ "$timeout" -le 0 ]; then
            echo "run timeout"
            exit 1
        fi

        sleep $((factor * 2))
        factor=$((factor * 2))
        timeout=$((timeout - factor))
    done
}


# check resource exist.
status=$(curl -sS -w "%{http_code}" -o /dev/null "$commonPath/resources/$resourceName" -X "GET" -H "Authorization: Bearer $token" $tlsVerify)
# if status >= 200 and status < 300
if [ "$status" -ge 200 ] && [ "$status" -lt 300 ]; then
    jobType="upgrade"
fi

# If jobType create resource.
if [ "$jobType" == "create" ]; then
    response=$(curl -sS "$commonPath/resources" -X "POST" -H "content-type: application/json" -H "Authorization: Bearer $token" -d '{{inputs.parameters.attributes}}' $tlsVerify)
    name=$(echo $response | jq -r '.name')

    if [ "$name" == "null" ]; then
        echo "failed create resource, response: $response"
        exit 1
    fi
fi

# If jobType upgrade resource.
if [ "$jobType" == "upgrade" ]; then
    status=$(curl -sS -w "%{http_code}" -o /dev/null "$commonPath/resources/$resourceName/upgrade" -X "PUT" -H "content-type: application/json" -H "Authorization: Bearer $token" -d '{{inputs.parameters.attributes}}' $tlsVerify)

    # if not status >= 200 and status < 300
    if [ "$status" -lt 200 ] || [ "$status" -ge 300 ]; then
        echo "failed upgrade resource, response status: $status"
        exit 1
    fi
fi

# Get latest run id
runResponse=$(common_get "$commonPath/resources/$resourceName/runs?page=1&perPage=1&sort=-createTime")
runID=$(echo $runResponse | jq -r '.items[0].id')

# Watch run plan logs until the plan finished.
common_get "$commonPath/resources/$resourceName/runs/$runID/log?jobType=plan&watch=true"

# Check run status, wait until run status is running.
wait_run_status "$resourceName" "$runID" "Running"

# Watch run logs until the run finished.
common_get "$commonPath/resources/$resourceName/runs/$runID/log?jobType=apply&watch=true"

# Check run status, wait until run status is ready.
wait_run_status "$resourceName" "$runID" "Succeeded"
`

// ServiceStepManager is service to generate service configs.
type ServiceStepManager struct {
	mc model.ClientSet
}

// NewServiceStepManager.
func NewServiceStepManager(opts types.CreateOptions) types.StepManager {
	return &ServiceStepManager{
		mc: opts.ModelClient,
	}
}

// GenerateTemplates generate service templates.
// If service exist in environment, job type is upgrade.
// Otherwise, job type is create.
func (s *ServiceStepManager) GenerateTemplates(
	ctx context.Context,
	stepExecution *model.WorkflowStepExecution,
) (main *wfv1.Template, subTemplates []*wfv1.Template, err error) {
	deployerImage := settings.WorkflowStepServiceImage.ShouldValue(ctx, s.mc)
	imageRegistry := settings.ImageRegistry.ShouldValue(ctx, s.mc)

	deployerImage = path.Join(imageRegistry, deployerImage)

	environment, ok := stepExecution.Attributes["environment"].(map[string]any)
	if !ok {
		return nil, nil, errors.New("environment is not found")
	}

	environmentID, ok := environment["id"].(string)
	if !ok {
		return nil, nil, errors.New("environment id is not found")
	}

	resourceName, ok := stepExecution.Attributes["name"].(string)
	if !ok {
		return nil, nil, errors.New("service name is not found")
	}

	// Inject workflow step execution id to request.
	stepAttrs := stepExecution.Attributes
	stepAttrs["workflowStepExecutionID"] = stepExecution.ID.String()

	attrs, err := json.Marshal(stepAttrs)
	if err != nil {
		return nil, nil, err
	}

	// An service type workflow template
	// Interact with walrus server to create or update service.
	// Watch service logs until the service finished.
	main = &wfv1.Template{
		Name: StepTemplateName(stepExecution),
		Metadata: wfv1.Metadata{
			Labels: map[string]string{
				"step-execution-id": stepExecution.ID.String(),
			},
		},
		Inputs: wfv1.Inputs{
			Parameters: []wfv1.Parameter{
				{
					Name:  "environmentID",
					Value: wfv1.AnyStringPtr(environmentID),
				},
				{
					Name:  "attributes",
					Value: wfv1.AnyStringPtr(string(attrs)),
				},
				{
					Name: "token",
				},
				{
					Name:  "resourceName",
					Value: wfv1.AnyStringPtr(resourceName),
				},
			},
		},
		Script: &wfv1.ScriptTemplate{
			Container: apiv1.Container{
				Image:           deployerImage,
				ImagePullPolicy: apiv1.PullIfNotPresent,
				Command:         []string{"sh"},
			},
			Source: stepSource,
		},
	}

	if stepExecution.RetryStrategy != nil {
		limit := intstr.FromInt(stepExecution.RetryStrategy.Limit)
		main.RetryStrategy = &wfv1.RetryStrategy{
			Limit:       &limit,
			RetryPolicy: stepExecution.RetryStrategy.RetryPolicy,
			Backoff:     stepExecution.RetryStrategy.Backoff,
		}
	}

	if stepExecution.Timeout > 0 {
		timeout := intstr.FromInt(stepExecution.Timeout)
		main.ActiveDeadlineSeconds = &timeout
	}

	return main, nil, nil
}
