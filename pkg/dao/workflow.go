package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstage"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstep"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
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
	ExecuteOptions struct {
		RestCfg *rest.Config
		Params  map[string]string
		// Description is the description of the workflow execution.
		Description string
	}
	CreateWorkflowExecutionOptions struct {
		ExecuteOptions `json:",inline"`

		Workflow *model.Workflow
	}

	CreateWorkflowStageExecutionOptions struct {
		ExecuteOptions `json:",inline"`

		WorkflowExecution *model.WorkflowExecution
		Stage             *model.WorkflowStage
	}

	CreateWorkflowStepExecutionOptions struct {
		ExecuteOptions `json:",inline"`

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
			ExecuteOptions:    opts.ExecuteOptions,
			WorkflowExecution: entity,
			Stage:             stage,
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
			ExecuteOptions: opts.ExecuteOptions,
			StageExecution: entity,
			Step:           step,
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
		step   = opts.Step
		wse    = opts.StageExecution
		params = opts.Params
	)

	attrs, err := parseParams(step.Attributes, params)
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

// parseParams parses the params into the attributes.
// The params are the key-value pairs that are used to replace the keywords in the attributes.
// The keywords are the strings that are wrapped by dollar sign and curly braces, like "${key}".
// The keyword is "${key}" will be replaced by the value of the key in the params.
func parseParams(attrs map[string]any, params map[string]string) (map[string]any, error) {
	if err := CheckParams(params); err != nil {
		return nil, err
	}

	bs, err := json.Marshal(attrs)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(bs))

	for k, v := range params {
		paramReg := regexp.MustCompile(fmt.Sprintf(`\$\{\s*%s\s*\}`, k))
		bs = paramReg.ReplaceAll(bs, []byte(v))
	}

	var replacedSpec map[string]any

	err = json.Unmarshal(bs, &replacedSpec)
	if err != nil {
		return nil, err
	}

	return replacedSpec, nil
}

// CheckParams checks if the params contain keywords.
// These keywords will reconginzed as the params like "{{input.parameters.xxx}}"
// of argo that may cause workflow execution failure.
func CheckParams(params map[string]string) error {
	keywordsReg := regexp.MustCompile(`{{.*}}`)

	for k, v := range params {
		if keywordsReg.MatchString(v) {
			return fmt.Errorf("params contain keywords: %s=%s", k, v)
		}
	}

	return nil
}
