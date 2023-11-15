package workflow

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	apiconfig "github.com/seal-io/walrus/pkg/apis/config"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/workflow/step"
	steptypes "github.com/seal-io/walrus/pkg/workflow/step/types"
	"github.com/seal-io/walrus/utils/pointer"
	"github.com/seal-io/walrus/utils/strs"
)

const (
	templateTypeStep     = "step"
	templateTypeStage    = "stage"
	templateTypeWorkflow = "workflow"

	// Stages represents the stage of the template task.
	// E.G. A step execution may have enter, main and exit stages
	// of its generated templates.
	templateStageEnter = "enter"
	templateStageMain  = "main"
	templateStageExit  = "exit"

	workflowEntrypointTemplateName = "entrypoint"
	workflowEnterTemplateName      = "workflowEnter"
	workflowExitTemplateName       = "workflowExit"

	executionUpdateURL = "{{workflow.parameters.server}}/v1/projects/{{workflow.parameters.projectID}}" +
		"/workflows/{{workflow.parameters.workflowID}}" +
		"/executions/{{inputs.parameters.id}}"
	stageExecutionUpdateURL = "{{workflow.parameters.server}}/v1/projects/{{workflow.parameters.projectID}}" +
		"/workflows/{{workflow.parameters.workflowID}}" +
		"/executions/{{inputs.parameters.workflowExecutionID}}" +
		"/stage-executions/{{inputs.parameters.id}}"
	stepExecutionUpdateURL = "{{workflow.parameters.server}}/v1/projects/{{workflow.parameters.projectID}}" +
		"/workflows/{{workflow.parameters.workflowID}}" +
		"/executions/{{inputs.parameters.workflowExecutionID}}" +
		"/stage-executions/{{inputs.parameters.workflowStageExecutionID}}" +
		"/step-executions/{{inputs.parameters.id}}"
)

const (
	statusRequestBody = `{
	"id": "{{inputs.parameters.id}}",
	"status": "{{inputs.parameters.status}}"
}`
)

var (
	limit  = intstr.FromInt(2)
	factor = intstr.FromInt(2)
	// The status retry strategy of updating status of workflow,
	// stage and step.
	statusUpdateRetryStrategy = &wfv1.RetryStrategy{
		Limit:       &limit,
		RetryPolicy: wfv1.RetryPolicyOnFailure,
		Backoff: &wfv1.Backoff{
			Duration:    "1",
			Factor:      &factor,
			MaxDuration: "1m",
		},
	}
)

// TemplateManager is the manager of workflow templates.
// Manager generate argo workflow definition with model.WorkflowExecution.
// It generates templates for workflow with workflow, stage and step executions.
type TemplateManager struct {
	mc model.ClientSet
}

func NewTemplateManager(mc model.ClientSet) *TemplateManager {
	return &TemplateManager{
		mc: mc,
	}
}

// ToArgoWorkflow returns an argo workflow for a workflow execution.
// The workflow execution MUST contains edges of stage and step executions.
func (t *TemplateManager) ToArgoWorkflow(
	ctx context.Context,
	workflowExecution *model.WorkflowExecution,
	token string,
) (*wfv1.Workflow, error) {
	// Prepare address for terraform backend.
	serverAddress, err := settings.ServeUrl.Value(ctx, t.mc)
	if err != nil {
		return nil, err
	}

	if serverAddress == "" {
		return nil, errors.New("server address is empty")
	}

	wfTemplates, err := t.GetWorkflowExecutionTemplates(ctx, workflowExecution.Edges.Stages)
	if err != nil {
		return nil, err
	}

	workflowTemplates := make([]wfv1.Template, 0, len(wfTemplates)+2)
	for _, tpl := range wfTemplates {
		workflowTemplates = append(workflowTemplates, *tpl)
	}

	wf := &wfv1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", workflowExecution.Name, workflowExecution.ID.String()),
			Labels: map[string]string{
				workflowExecutionIDLabel: workflowExecution.ID.String(),
			},
		},
		Spec: wfv1.WorkflowSpec{
			Entrypoint: workflowEntrypointTemplateName,
			Arguments: wfv1.Arguments{
				Parameters: []wfv1.Parameter{
					{
						Name:  "server",
						Value: wfv1.AnyStringPtr(serverAddress),
					},
					{
						Name:  "projectID",
						Value: wfv1.AnyStringPtr(workflowExecution.ProjectID.String()),
					},
					{
						Name:  "workflowID",
						Value: wfv1.AnyStringPtr(workflowExecution.WorkflowID.String()),
					},
					{
						Name:  "tlsVerify",
						Value: wfv1.AnyStringPtr(apiconfig.TlsCertified.Get()),
					},
					{
						Name:  "token",
						Value: wfv1.AnyStringPtr(token),
					},
				},
			},
			SecurityContext: &corev1.PodSecurityContext{
				RunAsNonRoot: pointer.Bool(true),
				RunAsUser:    pointer.Int64(1000),
			},
			TTLStrategy: &wfv1.TTLStrategy{
				SecondsAfterCompletion: pointer.Int32(600),
			},
			ServiceAccountName: types.WorkflowServiceAccountName,
			Templates:          workflowTemplates,
		},
	}

	if workflowExecution.Timeout > 0 {
		wf.Spec.ActiveDeadlineSeconds = pointer.Int64(int64(workflowExecution.Timeout))
	}

	if workflowExecution.Parallelism > 0 {
		wf.Spec.Parallelism = pointer.Int64(int64(workflowExecution.Parallelism))
	}

	return wf, nil
}

// GetWorkflowExecutionTemplates get workflow execution templates.
func (t *TemplateManager) GetWorkflowExecutionTemplates(
	ctx context.Context,
	stageExecutions model.WorkflowStageExecutions,
) ([]*wfv1.Template, error) {
	// Calculate the length of workflow templates.
	workflowTemplateLen := 1

	for i := range stageExecutions {
		stageExec := stageExecutions[i]
		workflowTemplateLen += len(stageExec.Edges.Steps)*2 + 1
	}

	workflowTemplates := make([]*wfv1.Template, 0, workflowTemplateLen)
	tasks := make([]wfv1.DAGTask, len(stageExecutions))
	entrypoint := &wfv1.Template{
		Name: workflowEntrypointTemplateName,
		DAG: &wfv1.DAGTemplate{
			Tasks: tasks,
		},
	}
	workflowTemplates = append(workflowTemplates, entrypoint)

	for i, stageExec := range stageExecutions {
		stageTemplate, stageSubTemplates, err := t.GetStageExecutionTemplates(ctx, stageExec)
		if err != nil {
			return nil, err
		}

		workflowTemplates = append(workflowTemplates, stageTemplate)
		workflowTemplates = append(workflowTemplates, stageSubTemplates...)

		var dependencies []string

		if i > 0 {
			// Add previous stage task as dependency.
			dependencies = append(dependencies, entrypoint.DAG.Tasks[i-1].Name)
		}

		entrypoint.DAG.Tasks[i] = wfv1.DAGTask{
			Name:         statusTemplateName(stageExec.ID, templateTypeStage, templateStageEnter),
			Template:     stageTemplate.Name,
			Dependencies: dependencies,
		}
	}

	return workflowTemplates, nil
}

// GetWorkflowExecutionStatusTemplate returns the status template of a workflow execution.
// The status template handler sync the status of the workflow execution to "Running", "Succeeded" or "Failed".
// It will be called with the lifecycle hook of the stage execution.
func (t *TemplateManager) GetWorkflowExecutionStatusTemplate(
	name string,
	wf *model.WorkflowExecution,
) *wfv1.Template {
	status := "{{workflow.status}}"
	if name == workflowEnterTemplateName {
		status = types.ExecutionStatusRunning
	}

	return &wfv1.Template{
		Name: name,
		Inputs: wfv1.Inputs{
			Parameters: []wfv1.Parameter{
				{
					Name:  "id",
					Value: wfv1.AnyStringPtr(wf.ID.String()),
				},
				{
					Name:  "status",
					Value: wfv1.AnyStringPtr(status),
				},
			},
		},
		HTTP: &wfv1.HTTP{
			URL: "{{workflow.parameters.server}}/v1/projects/{{workflow.parameters.projectID}}" +
				"/workflows/{{workflow.parameters.workflowID}}" +
				"/executions/{{inputs.parameters.id}}",
			Method: http.MethodPut,
			Headers: wfv1.HTTPHeaders{
				{
					Name:  "Content-Type",
					Value: "application/json",
				},
				{
					Name:  "Authorization",
					Value: "Bearer {{workflow.parameters.token}}",
				},
			},
			TimeoutSeconds:     pointer.Int64(5),
			InsecureSkipVerify: !apiconfig.TlsCertified.Get(),
			SuccessCondition:   "response.statusCode >= 200 && response.statusCode < 300",
			Body:               statusRequestBody,
		},
		RetryStrategy: statusUpdateRetryStrategy,
	}
}

// GetWorkflowExecutionEnterTemplate returns the enter template of a workflow execution.
func (t *TemplateManager) GetWorkflowExecutionEnterTemplate(wf *model.WorkflowExecution) *wfv1.Template {
	return t.GetWorkflowExecutionStatusTemplate(workflowEnterTemplateName, wf)
}

// getExitTemplate returns template for workflow exit handler.
func (t *TemplateManager) GetWorkflowExecutionExitTemplate(wf *model.WorkflowExecution) *wfv1.Template {
	return t.GetWorkflowExecutionStatusTemplate(workflowExitTemplateName, wf)
}

// GetStageExecutionTemplates extends one stage execution to template,
// return stage template and its step templates.
func (t *TemplateManager) GetStageExecutionTemplates(
	ctx context.Context,
	stageExecution *model.WorkflowStageExecution,
) (stageTemplate *wfv1.Template, subTemplates []*wfv1.Template, err error) {
	stageTemplate = &wfv1.Template{
		Name: statusTemplateName(stageExecution.ID, templateTypeStage, templateStageMain),
		DAG:  &wfv1.DAGTemplate{},
	}

	tasks := make([]wfv1.DAGTask, 0, len(stageExecution.Edges.Steps))

	// Get step templates with step executions.
	for _, stepExecution := range stageExecution.Edges.Steps {
		extendTemplate, stepTemplates, err := t.GetStepExecutionExtendTemplates(ctx, stepExecution)
		if err != nil {
			return nil, nil, err
		}

		subTemplates = append(subTemplates, stepTemplates...)
		subTemplates = append(subTemplates, extendTemplate)

		tasks = append(tasks, wfv1.DAGTask{
			Name:     statusTemplateName(stepExecution.ID, templateTypeStep, templateStageEnter),
			Template: extendTemplate.Name,
		})
	}

	stageTemplate.DAG.Tasks = tasks

	return
}

// GetStageExecutionStatusTemplate returns the status template of a stage execution.
// The status template handler sync the status of the stage execution to "Running", "Succeeded" or "Failed".
func (t *TemplateManager) GetStageExecutionStatusTemplate(
	name string,
	stageExecution *model.WorkflowStageExecution,
) *wfv1.Template {
	return &wfv1.Template{
		Name: name,
		Inputs: wfv1.Inputs{
			Parameters: []wfv1.Parameter{
				{
					Name:  "id",
					Value: wfv1.AnyStringPtr(stageExecution.ID.String()),
				},
				{
					Name:  "workflowExecutionID",
					Value: wfv1.AnyStringPtr(stageExecution.WorkflowExecutionID.String()),
				},
				{
					Name: "status",
				},
			},
		},
		HTTP: &wfv1.HTTP{
			URL:    stageExecutionUpdateURL,
			Method: http.MethodPut,
			Headers: wfv1.HTTPHeaders{
				{
					Name:  "Authorization",
					Value: "Bearer {{workflow.parameters.token}}",
				},
				{
					Name:  "Content-Type",
					Value: "application/json",
				},
			},
			TimeoutSeconds:     pointer.Int64(5),
			Body:               statusRequestBody,
			SuccessCondition:   "response.statusCode >= 200 && response.statusCode < 300",
			InsecureSkipVerify: !apiconfig.TlsCertified.Get(),
		},
		RetryStrategy: statusUpdateRetryStrategy,
	}
}

// GetStageExecutionEnterTemplate returns the enter template of a stage execution.
// The template handler sync the status of the stage execution to "Running".
func (t *TemplateManager) GetStageExecutionEnterTemplate(
	stageExecution *model.WorkflowStageExecution,
) *wfv1.Template {
	return t.GetStageExecutionStatusTemplate(
		statusTemplateName(
			stageExecution.ID,
			templateTypeStage,
			templateStageEnter,
		),
		stageExecution,
	)
}

// GetStageExecutionExitTemplate returns the exit template of a stage execution.
// The template handler sync the status of the stage execution to "Succeeded" or "Failed".
func (t *TemplateManager) GetStageExecutionExitTemplate(
	stageExecution *model.WorkflowStageExecution,
) *wfv1.Template {
	return t.GetStageExecutionStatusTemplate(
		statusTemplateName(
			stageExecution.ID,
			templateTypeStage,
			templateStageExit,
		),
		stageExecution,
	)
}

// GetStepExecutionExtendTemplates extends one step execution to three step executions, enter template, main template,
// exit step template, which are used to update the status of the step execution.
// The extend templates are used to manager lifecycle of the step execution.
func (t *TemplateManager) GetStepExecutionExtendTemplates(
	ctx context.Context,
	stepExecution *model.WorkflowStepExecution,
) (extendTemplate *wfv1.Template, stepTemplates []*wfv1.Template, err error) {
	secretName := fmt.Sprintf("workflow-execution-%s", stepExecution.WorkflowExecutionID.String())
	secretTemplate := getSecretTemplate(
		statusTemplateName(stepExecution.ID, templateTypeStep, "secret"),
		secretName,
		"token",
	)

	stepTemplates, err = t.GetStepExecutionTemplates(ctx, stepExecution)
	if err != nil {
		return nil, nil, err
	}

	stepTemplates = append(stepTemplates, secretTemplate)

	tokenRef := fmt.Sprintf("{{=fromBase64(steps['%s'].outputs.parameters.secretValue)}}", secretTemplate.Name)

	// Extend one step template to three step templates, enter template, main template,
	// and exit template.
	extendTemplate = &wfv1.Template{
		Name: fmt.Sprintf("%s-extend", step.StepTemplateName(stepExecution)),
		Steps: []wfv1.ParallelSteps{
			{
				Steps: []wfv1.WorkflowStep{
					{
						Name:     secretTemplate.Name,
						Template: secretTemplate.Name,
						Arguments: wfv1.Arguments{
							Parameters: []wfv1.Parameter{
								{
									Name:  "secretName",
									Value: wfv1.AnyStringPtr(secretName),
								},
								{
									Name:  "secretKey",
									Value: wfv1.AnyStringPtr("token"),
								},
							},
						},
					},
				},
			},
			{
				Steps: []wfv1.WorkflowStep{
					{
						Name:     statusTemplateName(stepExecution.ID, templateTypeStep, templateStageMain),
						Template: stepTemplates[0].Name,
						Arguments: wfv1.Arguments{
							Parameters: []wfv1.Parameter{
								{
									Name:  "token",
									Value: wfv1.AnyStringPtr(tokenRef),
								},
							},
						},
					},
				},
			},
		},
	}

	return extendTemplate, stepTemplates, nil
}

// GetStepExecutionTemplates extends one step execution to three step executions, enter template, main template,
// exit step template, which are used to update the status of the step execution.
func (t *TemplateManager) GetStepExecutionTemplates(
	ctx context.Context,
	stepExecution *model.WorkflowStepExecution,
) ([]*wfv1.Template, error) {
	stepService, err := step.GetStepManager(steptypes.CreateOptions{
		Type:        steptypes.Type(stepExecution.Type),
		ModelClient: t.mc,
	})
	if err != nil {
		return nil, err
	}

	// Generate service template.
	mainTemplate, subTemplates, err := stepService.GenerateTemplates(ctx, stepExecution)
	if err != nil {
		return nil, err
	}

	templates := []*wfv1.Template{
		mainTemplate,
	}

	templates = append(templates, subTemplates...)

	return templates, nil
}

// GetStepExecutionStatusTemplate returns the status template of a step execution.
// The status template handler sync the status of the step execution to
// "Running", "Succeeded" or "Failed".
func (t *TemplateManager) GetStepExecutionStatusTemplate(
	name string,
	stepExecution *model.WorkflowStepExecution,
) *wfv1.Template {
	return &wfv1.Template{
		Name: name,
		Inputs: wfv1.Inputs{
			Parameters: []wfv1.Parameter{
				{
					Name:  "id",
					Value: wfv1.AnyStringPtr(stepExecution.ID.String()),
				},
				{
					Name:  "workflowStageExecutionID",
					Value: wfv1.AnyStringPtr(stepExecution.WorkflowStageExecutionID.String()),
				},
				{
					Name:  "workflowExecutionID",
					Value: wfv1.AnyStringPtr(stepExecution.WorkflowExecutionID.String()),
				},
				{
					Name: "status",
				},
				{
					Name: "token",
				},
			},
		},
		HTTP: &wfv1.HTTP{
			URL:    stepExecutionUpdateURL,
			Method: http.MethodPut,
			Headers: wfv1.HTTPHeaders{
				{
					Name:  "Authorization",
					Value: "Bearer {{inputs.parameters.token}}",
				},
				{
					Name:  "Content-Type",
					Value: "application/json",
				},
			},
			TimeoutSeconds:     pointer.Int64(5),
			Body:               statusRequestBody,
			SuccessCondition:   "response.statusCode >= 200 && response.statusCode < 300",
			InsecureSkipVerify: !apiconfig.TlsCertified.Get(),
		},
		RetryStrategy: statusUpdateRetryStrategy,
	}
}

func statusTemplateName(id object.ID, templateType, stage string) string {
	return strs.Join("-", templateType, stage, id.String())
}

func parseTemplateName(s string) (templateType, stage string, id object.ID, ok bool) {
	splits := strings.Split(s, "-")
	if len(splits) != 3 {
		return
	}

	templateType, stage, id = splits[0], splits[1], object.ID(splits[2])

	if !id.Valid() {
		return "", "", "", false
	}

	ok = true

	return
}

// getSecretTemplate returns a template for getting a secret.
// It will get step execution token from the secret.
// This template will get k8s secret from k8s api server.
func getSecretTemplate(name, secretName, secretKey string) *wfv1.Template {
	return &wfv1.Template{
		Name: name,
		Inputs: wfv1.Inputs{
			Parameters: []wfv1.Parameter{
				{
					Name:  "secretName",
					Value: wfv1.AnyStringPtr(secretName),
				}, {
					Name:  "secretKey",
					Value: wfv1.AnyStringPtr(secretKey),
				}, {
					Name:  "secretNamespace",
					Value: wfv1.AnyStringPtr(types.WalrusSystemNamespace),
				},
			},
		},
		Outputs: wfv1.Outputs{
			Parameters: []wfv1.Parameter{
				{
					Name: "secretValue",
					ValueFrom: &wfv1.ValueFrom{
						JSONPath: "{.data.{{inputs.parameters.secretKey}}}",
					},
				},
			},
		},
		Timeout: "10s",
		RetryStrategy: &wfv1.RetryStrategy{
			Limit:       &limit,
			RetryPolicy: wfv1.RetryPolicyOnFailure,
			Backoff: &wfv1.Backoff{
				Duration:    "5",
				Factor:      &factor,
				MaxDuration: "1m",
			},
		},
		Resource: &wfv1.ResourceTemplate{
			Action: "get",
			Manifest: `apiVersion: v1
kind: Secret
metadata:
  name: "{{inputs.parameters.secretName}}"
  namespace: "{{inputs.parameters.secretNamespace}}"`,
		},
	}
}
