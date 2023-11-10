package step

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/workflow/step/types"
)

//nolint:lll
const stepSource = `#!/bin/sh
set -e
set -o pipefail

serverURL="{{workflow.parameters.server}}"
projectID="{{workflow.parameters.projectID}}"
environmentID="{{inputs.parameters.environmentID}}"
token="{{inputs.parameters.token}}"
jobType="{{inputs.parameters.jobType}}"
resourceName="{{inputs.parameters.resourceName}}"
commonPath="$serverURL/v1/projects/$projectID/environments/$environmentID"

# if skip tls verify
tlsVerify="-k"
if [ "{{workflow.parameters.tlsVerify}}" == "true" ]; then
	tlsVerify=""
fi

# If jobType create resource.
if [ "$jobType" == "create" ]; then
	response=$(curl -s "$commonPath/resources" -X "POST" -H "content-type: application/json" -H "Authorization: Bearer $token" -d '{{inputs.parameters.attributes}}' $tlsVerify)

	resourceName=$(echo $response | jq -r '.name')
	if [ "$resourceName" == "null" ]; then
		echo "failed create resource, response: $response"
		exit 1
	fi
fi

# If jobType upgrade resource.
if [ "$jobType" == "upgrade" ]; then
	response=$(curl -s "$commonPath/resources/$resourceName/upgrade" -X "PUT" -H "content-type: application/json" -H "Authorization: Bearer $token" -d '{{inputs.parameters.attributes}}' $tlsVerify)
	resourceName=$(echo $response | jq -r '.name')
	if [ "$resourceName" == "null" ]; then
		echo "failed upgrade resource, response: $response"
		exit 1
	fi
fi

# Get latest revision id
revisionResponse=$(curl -s "$commonPath/resources/$resourceName/revisions?page=1&perPage=1&sort=-createTime" -X GET -H "Authorization: Bearer $token" $tlsVerify)
revisionID=$(echo $revisionResponse | jq -r '.items[0].id')

# Watch service logs until the service finished.
curl -o - -s "$commonPath/resources/$resourceName/revisions/$revisionID/log?jobType=$watchType&watch=true" -X GET -H "Authorization: Bearer $token" $tlsVerify --compressed
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
) (main *v1alpha1.Template, subTemplates []*v1alpha1.Template, err error) {
	deployerImage := settings.WorkflowStepServiceImage.ShouldValue(ctx, s.mc)

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

	// If resource exist in environment, job type is upgrade.
	// Otherwise, job type is create.
	svc, err := s.mc.Resources().Query().
		Select(
			resource.FieldID,
			resource.FieldName,
			resource.FieldEnvironmentID,
		).
		Where(
			resource.EnvironmentID(object.ID(environmentID)),
			resource.Name(resourceName),
		).
		Only(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, nil, fmt.Errorf("failed to get service: %w", err)
	}

	jobType := "create"
	if svc != nil {
		jobType = "upgrade"
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
	main = &v1alpha1.Template{
		Name: fmt.Sprintf("step-execution-%s", stepExecution.ID.String()),
		Metadata: v1alpha1.Metadata{
			Labels: map[string]string{
				"step-execution-id": stepExecution.ID.String(),
			},
		},
		Inputs: v1alpha1.Inputs{
			Parameters: []v1alpha1.Parameter{
				{
					Name:  "environmentID",
					Value: v1alpha1.AnyStringPtr(environmentID),
				},
				{
					Name:  "attributes",
					Value: v1alpha1.AnyStringPtr(string(attrs)),
				},
				{
					Name:  "jobType",
					Value: v1alpha1.AnyStringPtr(jobType),
				},
				{
					Name: "token",
				},
				{
					Name:  "resourceName",
					Value: v1alpha1.AnyStringPtr(resourceName),
				},
			},
		},
		Script: &v1alpha1.ScriptTemplate{
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
		main.RetryStrategy = &v1alpha1.RetryStrategy{
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
