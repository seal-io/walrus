package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstage"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstep"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/strs"
)

func WorkflowStagesEdgeSave(ctx context.Context, mc model.ClientSet, entity *model.Workflow) error {
	if len(entity.Edges.Stages) == 0 {
		return nil
	}

	// StageIDs that should be kept.
	// These stage will be updated.
	stageIDs := make([]object.ID, 0, len(entity.Edges.Stages))

	for i := range entity.Edges.Stages {
		stage := entity.Edges.Stages[i]
		if !stage.ID.Valid() {
			continue
		}

		stageIDs = append(stageIDs, stage.ID)
	}

	// Delete stale items.
	_, err := mc.WorkflowStages().Delete().
		Where(
			workflowstage.WorkflowID(entity.ID),
			workflowstage.IDNotIn(stageIDs...)).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Add new items or update existing items.
	newItems := entity.Edges.Stages

	for i := range newItems {
		if newItems[i] == nil {
			return errors.New("invalid input: nil workflow stage")
		}
		newItems[i].WorkflowID = entity.ID
		newItems[i].ProjectID = entity.ProjectID
		newItems[i].Order = i

		if newItems[i].ID.Valid() {
			newItems[i], err = mc.WorkflowStages().UpdateOne(newItems[i]).
				Set(newItems[i]).
				SaveE(ctx, WorkflowStageStepsEdgeSave)
		} else {
			newItems[i], err = mc.WorkflowStages().Create().
				Set(newItems[i]).
				SaveE(ctx, WorkflowStageStepsEdgeSave)
		}

		if err != nil {
			return err
		}
	}

	entity.Edges.Stages = newItems // Feedback.

	return nil
}

// WorkflowStageStepsEdgeSave saves the edge steps of model.WorkflowStage entity.
func WorkflowStageStepsEdgeSave(ctx context.Context, mc model.ClientSet, entity *model.WorkflowStage) error {
	if len(entity.Edges.Steps) == 0 {
		return nil
	}

	// StepIDs that should be kept.
	// These steps will be updated.
	stepIDs := make([]object.ID, 0, len(entity.Edges.Steps))

	for i := range entity.Edges.Steps {
		step := entity.Edges.Steps[i]
		if !step.ID.Valid() {
			continue
		}

		stepIDs = append(stepIDs, step.ID)
	}

	// Delete stale items.
	_, err := mc.WorkflowSteps().Delete().
		Where(
			workflowstep.WorkflowStageID(entity.ID),
			workflowstep.IDNotIn(stepIDs...),
		).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Add new items or update existing items.
	newItems := entity.Edges.Steps

	for i := range newItems {
		if newItems[i] == nil {
			return errors.New("invalid input: nil workflow step")
		}
		newItems[i].WorkflowStageID = entity.ID
		newItems[i].ProjectID = entity.ProjectID
		newItems[i].WorkflowID = entity.WorkflowID
		newItems[i].Order = i

		if newItems[i].ID.Valid() {
			newItems[i], err = mc.WorkflowSteps().UpdateOne(newItems[i]).
				Set(newItems[i]).
				Save(ctx)
		} else {
			newItems[i], err = mc.WorkflowSteps().Create().
				Set(newItems[i]).
				Save(ctx)
		}

		if err != nil {
			return err
		}
	}

	entity.Edges.Steps = newItems // Feedback.

	return nil
}

type (
	CreateWorkflowExecutionOptions struct {
		Workflow *model.Workflow
		// Execution parameters that replace workflow variables config value.
		Variables map[string]string
		// Description is the description of the workflow execution.
		Description string
	}

	CreateWorkflowStageExecutionOptions struct {
		CreateWorkflowExecutionOptions `json:",inline"`

		WorkflowExecution *model.WorkflowExecution
		Stage             *model.WorkflowStage
	}

	CreateWorkflowStepExecutionOptions struct {
		CreateWorkflowExecutionOptions `json:",inline"`

		StageExecution *model.WorkflowStageExecution
		Step           *model.WorkflowStep
	}
)

func CreateWorkflowExecution(
	ctx context.Context,
	mc model.ClientSet,
	opts CreateWorkflowExecutionOptions,
) (*model.WorkflowExecution, error) {
	s := session.MustGetSubject(ctx)

	var trigger types.WorkflowExecutionTrigger

	wf := opts.Workflow

	switch wf.Type {
	case types.WorkflowTypeDefault:
		userSubject, err := mc.Subjects().Query().
			Where(subject.ID(s.ID)).
			Only(ctx)
		if err != nil {
			return nil, err
		}
		trigger = types.WorkflowExecutionTrigger{
			Type: types.WorkflowExecutionTriggerTypeManual,
			User: userSubject.Name,
		}
	default:
		return nil, fmt.Errorf("invalid workflow type: %s", wf.Type)
	}

	workflowExecution := &model.WorkflowExecution{
		Name:        wf.Name,
		Description: opts.Description,
		Type:        wf.Type,
		ProjectID:   wf.ProjectID,
		WorkflowID:  wf.ID,
		SubjectID:   s.ID,
		Parallelism: wf.Parallelism,
		// When creating a workflow execution, the execute times is always 1.
		Times:   1,
		Timeout: wf.Timeout,
		Trigger: trigger,
		Version: wf.Version + 1,
	}

	status.WorkflowExecutionStatusPending.Unknown(workflowExecution, "")
	workflowExecution.Status.SetSummary(status.WalkWorkflowExecution(&workflowExecution.Status))

	stageMap := make(map[object.ID]*model.WorkflowStage)
	for i := range wf.Edges.Stages {
		stageMap[wf.Edges.Stages[i].ID] = wf.Edges.Stages[i]
	}

	entity, err := mc.WorkflowExecutions().Create().
		Set(workflowExecution).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	stageExecutions := make(model.WorkflowStageExecutions, len(wf.Edges.Stages))

	for i, stage := range wf.Edges.Stages {
		// Create workflow stage execution.
		stageExecution, err := CreateWorkflowStageExecution(ctx, mc, CreateWorkflowStageExecutionOptions{
			CreateWorkflowExecutionOptions: opts,
			WorkflowExecution:              entity,
			Stage:                          stage,
		})
		if err != nil {
			return nil, err
		}

		stageExecutions[i] = stageExecution
	}

	entity.Edges.Stages = stageExecutions

	return entity, nil
}

func CreateWorkflowStageExecution(
	ctx context.Context,
	mc model.ClientSet,
	opts CreateWorkflowStageExecutionOptions,
) (*model.WorkflowStageExecution, error) {
	stage := opts.Stage
	stageExec := &model.WorkflowStageExecution{
		Name:                stage.Name,
		ProjectID:           stage.ProjectID,
		WorkflowID:          stage.WorkflowID,
		WorkflowStageID:     stage.ID,
		WorkflowExecutionID: opts.WorkflowExecution.ID,
		Order:               stage.Order,
	}

	status.WorkflowStageExecutionStatusPending.Unknown(stageExec, "")
	stageExec.Status.SetSummary(status.WalkWorkflowStageExecution(&stageExec.Status))

	entity, err := mc.WorkflowStageExecutions().Create().
		Set(stageExec).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	stepMap := make(map[object.ID]*model.WorkflowStep)
	for i := range stage.Edges.Steps {
		stepMap[stage.Edges.Steps[i].ID] = stage.Edges.Steps[i]
	}

	stepExecutions := make(model.WorkflowStepExecutions, len(stage.Edges.Steps))

	for i, step := range stage.Edges.Steps {
		// Create workflow step execution.
		stepExecution, err := CreateWorkflowStepExecution(ctx, mc, CreateWorkflowStepExecutionOptions{
			CreateWorkflowExecutionOptions: opts.CreateWorkflowExecutionOptions,
			StageExecution:                 entity,
			Step:                           step,
		})
		if err != nil {
			return nil, err
		}

		stepExecutions[i] = stepExecution
	}

	entity.Edges.Steps = stepExecutions

	return entity, nil
}

func CreateWorkflowStepExecution(
	ctx context.Context,
	mc model.ClientSet,
	opts CreateWorkflowStepExecutionOptions,
) (*model.WorkflowStepExecution, error) {
	var (
		step = opts.Step
		wse  = opts.StageExecution
		vars = opts.Variables
	)

	attrs, err := parseWorkflowVariables(step.Attributes, vars, opts.Workflow.Variables)
	if err != nil {
		return nil, err
	}

	stepExec := &model.WorkflowStepExecution{
		Name:                     step.Name,
		Type:                     step.Type,
		Order:                    step.Order,
		ProjectID:                step.ProjectID,
		WorkflowID:               step.WorkflowID,
		WorkflowStepID:           step.ID,
		WorkflowExecutionID:      wse.WorkflowExecutionID,
		WorkflowStageExecutionID: wse.ID,
		Attributes:               attrs,
		Times:                    1,
		Timeout:                  step.Timeout,
		RetryStrategy:            step.RetryStrategy,
	}

	status.WorkflowStepExecutionStatusPending.Unknown(stepExec, "")
	stepExec.Status.SetSummary(status.WalkWorkflowStepExecution(&stepExec.Status))

	return mc.WorkflowStepExecutions().Create().
		Set(stepExec).
		Save(ctx)
}

// parseWorkflowVariables parses the params into the attributes.
// The params are the key-value pairs that are used to replace the keywords in the attributes.
// The keywords are the strings that are wrapped by dollar sign and curly braces, like "${workflow.var.key}".
// The keyword is "${workflow.var.key}" will be replaced by the value of the key in the params.
func parseWorkflowVariables(
	attrs map[string]any,
	params map[string]string,
	config types.WorkflowVariables,
) (map[string]any, error) {
	params, err := OverwriteWorkflowVariables(params, config)
	if err != nil {
		return nil, err
	}

	bs, err := json.Marshal(attrs)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\$\{\s*workflow\.var\.([-\w]+)\s*}`)
	replaced := re.ReplaceAllFunc(bs, func(s []byte) []byte {
		name := re.FindSubmatch(s)[1]
		name = bytes.TrimSpace(name)

		if val, ok := params[string(name)]; ok {
			return []byte(val)
		}

		return s
	})

	var replacedSpec map[string]any

	err = json.Unmarshal(replaced, &replacedSpec)
	if err != nil {
		return nil, err
	}

	return replacedSpec, nil
}

// OverwriteWorkflowVariables merges the variables into the config.
func OverwriteWorkflowVariables(vars map[string]string, config types.WorkflowVariables) (map[string]string, error) {
	mergedVars := make(map[string]string, len(config)+len(vars))
	overwriteKeys := sets.NewString()

	if vars == nil {
		vars = make(map[string]string)
	}

	for i := range config {
		c := config[i]
		// Only overwritable params can be replaced.
		if c.Overwrite {
			overwriteKeys.Insert(c.Name)
		}

		name, value := c.Name, c.Value
		if _, ok := vars[name]; ok && c.Overwrite {
			value = vars[name]
		}

		mergedVars[name] = value
	}

	for k := range vars {
		if !overwriteKeys.Has(k) {
			return nil, fmt.Errorf("invalid variables: %s", k)
		}
	}

	return mergedVars, nil
}

// GetRejectMessage returns the reject message of the workflow approval step execution.
func GetRejectMessage(ctx context.Context, mc model.ClientSet, entity *model.WorkflowStepExecution) (string, error) {
	message := ""

	approvalSpec, err := types.NewWorkflowStepApprovalSpec(entity.Attributes)
	if err != nil {
		return "", err
	}

	if approvalSpec.IsRejected() {
		rejectedUsers := approvalSpec.RejectedUsers
		rejectedUserNames := make([]string, len(rejectedUsers))

		subjects, err := mc.Subjects().Query().
			Where(subject.IDIn(rejectedUsers...)).
			All(ctx)
		if err != nil {
			return "", err
		}

		for i := range subjects {
			rejectedUserNames[i] = subjects[i].Name
		}

		message = "rejected by " + strs.Join[string](",", rejectedUserNames...)
	}

	return message, nil
}

// ResetWorkflowExecution resets the workflow execution and add execution times by 1.
// The workflow execution will be reset when rerun the workflow execution.
func ResetWorkflowExecution(ctx context.Context, mc model.ClientSet, workflowExecution *model.WorkflowExecution) error {
	status.WorkflowExecutionStatusPending.Reset(workflowExecution, "")
	workflowExecution.Status.SetSummary(status.WalkWorkflowExecution(&workflowExecution.Status))

	return mc.WorkflowExecutions().UpdateOne(workflowExecution).
		SetStatus(workflowExecution.Status).
		AddTimes(1).
		ClearExecuteTime().
		SetDuration(0).
		Exec(ctx)
}

// ResetWorkflowStageExecution resets the workflow stage execution and add execution times by 1.
// The workflow stage execution will be reset when rerun the workflow execution.
func ResetWorkflowStageExecution(
	ctx context.Context,
	mc model.ClientSet,
	stageExecution *model.WorkflowStageExecution,
) error {
	status.WorkflowStageExecutionStatusPending.Reset(stageExecution, "")
	stageExecution.Status.SetSummary(status.WalkWorkflowStageExecution(&stageExecution.Status))

	return mc.WorkflowStageExecutions().UpdateOne(stageExecution).
		SetStatus(stageExecution.Status).
		ClearExecuteTime().
		SetDuration(0).
		Exec(ctx)
}

// ResetWorkflowStepExecution resets the workflow step execution and add execution times by 1.
// The workflow step execution will be reset when rerun the workflow execution.
func ResetWorkflowStepExecution(
	ctx context.Context,
	mc model.ClientSet,
	stepExecution *model.WorkflowStepExecution,
) error {
	status.WorkflowStepExecutionStatusPending.Reset(stepExecution, "")
	stepExecution.Status.SetSummary(status.WalkWorkflowStepExecution(&stepExecution.Status))

	// Reset the approval spec.
	if stepExecution.Type == types.WorkflowStepTypeApproval {
		approvalSpec, err := types.NewWorkflowStepApprovalSpec(stepExecution.Attributes)
		if err != nil {
			return err
		}

		approvalSpec.Reset()

		stepExecution.Attributes = approvalSpec.ToAttributes()
	}

	return mc.WorkflowStepExecutions().UpdateOne(stepExecution).
		SetStatus(stepExecution.Status).
		AddTimes(1).
		ClearExecuteTime().
		SetAttributes(stepExecution.Attributes).
		SetDuration(0).
		Exec(ctx)
}
