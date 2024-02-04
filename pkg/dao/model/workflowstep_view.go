// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/seal-io/walrus/pkg/dao/model/workflowstep"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/json"
)

// WorkflowStepCreateInput holds the creation input of the WorkflowStep entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type WorkflowStepCreateInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to create WorkflowStep entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Stage indicates to create WorkflowStep entity MUST under the Stage route.
	Stage *WorkflowStageQueryInput `path:",inline" query:"-" json:"-"`

	// ID of the workflow that this workflow step belongs to.
	WorkflowID object.ID `path:"-" query:"-" json:"workflowID"`
	// Type of the workflow step.
	Type string `path:"-" query:"-" json:"type"`
	// Name holds the value of the "name" field.
	Name string `path:"-" query:"-" json:"name"`
	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes of the workflow step.
	Attributes map[string]interface{} `path:"-" query:"-" json:"attributes,omitempty"`
	// Inputs of the workflow step.
	Inputs map[string]interface{} `path:"-" query:"-" json:"inputs,omitempty"`
	// Outputs of the workflow step.
	Outputs map[string]interface{} `path:"-" query:"-" json:"outputs,omitempty"`
	// ID list of the workflow steps that this workflow step depends on.
	Dependencies []object.ID `path:"-" query:"-" json:"dependencies,omitempty"`
	// Retry policy of the workflow step.
	RetryStrategy *types.RetryStrategy `path:"-" query:"-" json:"retryStrategy,omitempty"`
	// Timeout seconds of the workflow step, 0 means no timeout.
	Timeout int `path:"-" query:"-" json:"timeout,omitempty"`
}

// Model returns the WorkflowStep entity for creating,
// after validating.
func (wsci *WorkflowStepCreateInput) Model() *WorkflowStep {
	if wsci == nil {
		return nil
	}

	_ws := &WorkflowStep{
		WorkflowID:    wsci.WorkflowID,
		Type:          wsci.Type,
		Name:          wsci.Name,
		Description:   wsci.Description,
		Labels:        wsci.Labels,
		Attributes:    wsci.Attributes,
		Inputs:        wsci.Inputs,
		Outputs:       wsci.Outputs,
		Dependencies:  wsci.Dependencies,
		RetryStrategy: wsci.RetryStrategy,
		Timeout:       wsci.Timeout,
	}

	if wsci.Project != nil {
		_ws.ProjectID = wsci.Project.ID
	}
	if wsci.Stage != nil {
		_ws.WorkflowStageID = wsci.Stage.ID
	}

	return _ws
}

// Validate checks the WorkflowStepCreateInput entity.
func (wsci *WorkflowStepCreateInput) Validate() error {
	if wsci == nil {
		return errors.New("nil receiver")
	}

	return wsci.ValidateWith(wsci.inputConfig.Context, wsci.inputConfig.Client, nil)
}

// ValidateWith checks the WorkflowStepCreateInput entity with the given context and client set.
func (wsci *WorkflowStepCreateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if wsci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when creating under the Project route.
	if wsci.Project != nil {
		if err := wsci.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}
	// Validate when creating under the Stage route.
	if wsci.Stage != nil {
		if err := wsci.Stage.ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// WorkflowStepCreateInputs holds the creation input item of the WorkflowStep entities.
type WorkflowStepCreateInputsItem struct {
	// ID of the workflow that this workflow step belongs to.
	WorkflowID object.ID `path:"-" query:"-" json:"workflowID"`
	// Type of the workflow step.
	Type string `path:"-" query:"-" json:"type"`
	// Name holds the value of the "name" field.
	Name string `path:"-" query:"-" json:"name"`
	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes of the workflow step.
	Attributes map[string]interface{} `path:"-" query:"-" json:"attributes,omitempty"`
	// Inputs of the workflow step.
	Inputs map[string]interface{} `path:"-" query:"-" json:"inputs,omitempty"`
	// Outputs of the workflow step.
	Outputs map[string]interface{} `path:"-" query:"-" json:"outputs,omitempty"`
	// ID list of the workflow steps that this workflow step depends on.
	Dependencies []object.ID `path:"-" query:"-" json:"dependencies,omitempty"`
	// Retry policy of the workflow step.
	RetryStrategy *types.RetryStrategy `path:"-" query:"-" json:"retryStrategy,omitempty"`
	// Timeout seconds of the workflow step, 0 means no timeout.
	Timeout int `path:"-" query:"-" json:"timeout,omitempty"`
}

// ValidateWith checks the WorkflowStepCreateInputsItem entity with the given context and client set.
func (wsci *WorkflowStepCreateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if wsci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// WorkflowStepCreateInputs holds the creation input of the WorkflowStep entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type WorkflowStepCreateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to create WorkflowStep entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Stage indicates to create WorkflowStep entity MUST under the Stage route.
	Stage *WorkflowStageQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*WorkflowStepCreateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the WorkflowStep entities for creating,
// after validating.
func (wsci *WorkflowStepCreateInputs) Model() []*WorkflowStep {
	if wsci == nil || len(wsci.Items) == 0 {
		return nil
	}

	_wss := make([]*WorkflowStep, len(wsci.Items))

	for i := range wsci.Items {
		_ws := &WorkflowStep{
			WorkflowID:    wsci.Items[i].WorkflowID,
			Type:          wsci.Items[i].Type,
			Name:          wsci.Items[i].Name,
			Description:   wsci.Items[i].Description,
			Labels:        wsci.Items[i].Labels,
			Attributes:    wsci.Items[i].Attributes,
			Inputs:        wsci.Items[i].Inputs,
			Outputs:       wsci.Items[i].Outputs,
			Dependencies:  wsci.Items[i].Dependencies,
			RetryStrategy: wsci.Items[i].RetryStrategy,
			Timeout:       wsci.Items[i].Timeout,
		}

		if wsci.Project != nil {
			_ws.ProjectID = wsci.Project.ID
		}
		if wsci.Stage != nil {
			_ws.WorkflowStageID = wsci.Stage.ID
		}

		_wss[i] = _ws
	}

	return _wss
}

// Validate checks the WorkflowStepCreateInputs entity .
func (wsci *WorkflowStepCreateInputs) Validate() error {
	if wsci == nil {
		return errors.New("nil receiver")
	}

	return wsci.ValidateWith(wsci.inputConfig.Context, wsci.inputConfig.Client, nil)
}

// ValidateWith checks the WorkflowStepCreateInputs entity with the given context and client set.
func (wsci *WorkflowStepCreateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if wsci == nil {
		return errors.New("nil receiver")
	}

	if len(wsci.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when creating under the Project route.
	if wsci.Project != nil {
		if err := wsci.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				wsci.Project = nil
			}
		}
	}
	// Validate when creating under the Stage route.
	if wsci.Stage != nil {
		if err := wsci.Stage.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				wsci.Stage = nil
			}
		}
	}

	for i := range wsci.Items {
		if wsci.Items[i] == nil {
			continue
		}

		if err := wsci.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// WorkflowStepDeleteInput holds the deletion input of the WorkflowStep entity,
// please tags with `path:",inline"` if embedding.
type WorkflowStepDeleteInput struct {
	WorkflowStepQueryInput `path:",inline"`
}

// WorkflowStepDeleteInputs holds the deletion input item of the WorkflowStep entities.
type WorkflowStepDeleteInputsItem struct {
	// ID of the WorkflowStep entity.
	ID object.ID `path:"-" query:"-" json:"id"`
}

// WorkflowStepDeleteInputs holds the deletion input of the WorkflowStep entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type WorkflowStepDeleteInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to delete WorkflowStep entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Stage indicates to delete WorkflowStep entity MUST under the Stage route.
	Stage *WorkflowStageQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*WorkflowStepDeleteInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the WorkflowStep entities for deleting,
// after validating.
func (wsdi *WorkflowStepDeleteInputs) Model() []*WorkflowStep {
	if wsdi == nil || len(wsdi.Items) == 0 {
		return nil
	}

	_wss := make([]*WorkflowStep, len(wsdi.Items))
	for i := range wsdi.Items {
		_wss[i] = &WorkflowStep{
			ID: wsdi.Items[i].ID,
		}
	}
	return _wss
}

// IDs returns the ID list of the WorkflowStep entities for deleting,
// after validating.
func (wsdi *WorkflowStepDeleteInputs) IDs() []object.ID {
	if wsdi == nil || len(wsdi.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(wsdi.Items))
	for i := range wsdi.Items {
		ids[i] = wsdi.Items[i].ID
	}
	return ids
}

// Validate checks the WorkflowStepDeleteInputs entity.
func (wsdi *WorkflowStepDeleteInputs) Validate() error {
	if wsdi == nil {
		return errors.New("nil receiver")
	}

	return wsdi.ValidateWith(wsdi.inputConfig.Context, wsdi.inputConfig.Client, nil)
}

// ValidateWith checks the WorkflowStepDeleteInputs entity with the given context and client set.
func (wsdi *WorkflowStepDeleteInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if wsdi == nil {
		return errors.New("nil receiver")
	}

	if len(wsdi.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.WorkflowSteps().Query()

	// Validate when deleting under the Project route.
	if wsdi.Project != nil {
		if err := wsdi.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				workflowstep.ProjectID(wsdi.Project.ID))
		}
	}

	// Validate when deleting under the Stage route.
	if wsdi.Stage != nil {
		if err := wsdi.Stage.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			q.Where(
				workflowstep.WorkflowStageID(wsdi.Stage.ID))
		}
	}

	ids := make([]object.ID, 0, len(wsdi.Items))

	for i := range wsdi.Items {
		if wsdi.Items[i] == nil {
			return errors.New("nil item")
		}

		if wsdi.Items[i].ID != "" {
			ids = append(ids, wsdi.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	if len(ids) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	idsCnt, err := q.Where(workflowstep.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != cap(ids) {
		return errors.New("found unrecognized item")
	}

	return nil
}

// WorkflowStepPatchInput holds the patch input of the WorkflowStep entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type WorkflowStepPatchInput struct {
	WorkflowStepQueryInput `path:",inline" query:"-" json:"-"`

	// Name holds the value of the "name" field.
	Name string `path:"-" query:"-" json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Annotations holds the value of the "annotations" field.
	Annotations map[string]string `path:"-" query:"-" json:"annotations,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime *time.Time `path:"-" query:"-" json:"createTime,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime *time.Time `path:"-" query:"-" json:"updateTime,omitempty"`
	// Type of the workflow step.
	Type string `path:"-" query:"-" json:"type,omitempty"`
	// ID of the workflow that this workflow step belongs to.
	WorkflowID object.ID `path:"-" query:"-" json:"workflowID,omitempty"`
	// Attributes of the workflow step.
	Attributes map[string]interface{} `path:"-" query:"-" json:"attributes,omitempty"`
	// Inputs of the workflow step.
	Inputs map[string]interface{} `path:"-" query:"-" json:"inputs,omitempty"`
	// Outputs of the workflow step.
	Outputs map[string]interface{} `path:"-" query:"-" json:"outputs,omitempty"`
	// Order of the workflow step.
	Order int `path:"-" query:"-" json:"order,omitempty"`
	// ID list of the workflow steps that this workflow step depends on.
	Dependencies []object.ID `path:"-" query:"-" json:"dependencies,omitempty"`
	// Retry policy of the workflow step.
	RetryStrategy *types.RetryStrategy `path:"-" query:"-" json:"retryStrategy,omitempty"`
	// Timeout seconds of the workflow step, 0 means no timeout.
	Timeout int `path:"-" query:"-" json:"timeout,omitempty"`

	patchedEntity *WorkflowStep `path:"-" query:"-" json:"-"`
}

// PatchModel returns the WorkflowStep partition entity for patching.
func (wspi *WorkflowStepPatchInput) PatchModel() *WorkflowStep {
	if wspi == nil {
		return nil
	}

	_ws := &WorkflowStep{
		Name:          wspi.Name,
		Description:   wspi.Description,
		Labels:        wspi.Labels,
		Annotations:   wspi.Annotations,
		CreateTime:    wspi.CreateTime,
		UpdateTime:    wspi.UpdateTime,
		Type:          wspi.Type,
		WorkflowID:    wspi.WorkflowID,
		Attributes:    wspi.Attributes,
		Inputs:        wspi.Inputs,
		Outputs:       wspi.Outputs,
		Order:         wspi.Order,
		Dependencies:  wspi.Dependencies,
		RetryStrategy: wspi.RetryStrategy,
		Timeout:       wspi.Timeout,
	}

	if wspi.Project != nil {
		_ws.ProjectID = wspi.Project.ID
	}
	if wspi.Stage != nil {
		_ws.WorkflowStageID = wspi.Stage.ID
	}

	return _ws
}

// Model returns the WorkflowStep patched entity,
// after validating.
func (wspi *WorkflowStepPatchInput) Model() *WorkflowStep {
	if wspi == nil {
		return nil
	}

	return wspi.patchedEntity
}

// Validate checks the WorkflowStepPatchInput entity.
func (wspi *WorkflowStepPatchInput) Validate() error {
	if wspi == nil {
		return errors.New("nil receiver")
	}

	return wspi.ValidateWith(wspi.inputConfig.Context, wspi.inputConfig.Client, nil)
}

// ValidateWith checks the WorkflowStepPatchInput entity with the given context and client set.
func (wspi *WorkflowStepPatchInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := wspi.WorkflowStepQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	q := cs.WorkflowSteps().Query()

	// Validate when querying under the Project route.
	if wspi.Project != nil {
		if err := wspi.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				workflowstep.ProjectID(wspi.Project.ID))
		}
	}

	// Validate when querying under the Stage route.
	if wspi.Stage != nil {
		if err := wspi.Stage.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			q.Where(
				workflowstep.WorkflowStageID(wspi.Stage.ID))
		}
	}

	if wspi.Refer != nil {
		if wspi.Refer.IsID() {
			q.Where(
				workflowstep.ID(wspi.Refer.ID()))
		} else {
			return errors.New("invalid identify refer of workflowstep")
		}
	} else if wspi.ID != "" {
		q.Where(
			workflowstep.ID(wspi.ID))
	} else {
		return errors.New("invalid identify of workflowstep")
	}

	q.Select(
		workflowstep.WithoutFields(
			workflowstep.FieldAnnotations,
			workflowstep.FieldCreateTime,
			workflowstep.FieldUpdateTime,
			workflowstep.FieldOrder,
		)...,
	)

	var e *WorkflowStep
	{
		// Get cache from previous validation.
		queryStmt, queryArgs := q.sqlQuery(setContextOp(ctx, q.ctx, "cache")).Query()
		ck := fmt.Sprintf("stmt=%v, args=%v", queryStmt, queryArgs)
		if cv, existed := cache[ck]; !existed {
			var err error
			e, err = q.Only(ctx)
			if err != nil {
				return err
			}

			// Set cache for other validation.
			cache[ck] = e
		} else {
			e = cv.(*WorkflowStep)
		}
	}

	_pm := wspi.PatchModel()

	_po, err := json.PatchObject(*e, *_pm)
	if err != nil {
		return err
	}

	_obj := _po.(*WorkflowStep)

	if e.Name != _obj.Name {
		return errors.New("field name is immutable")
	}
	if !reflect.DeepEqual(e.CreateTime, _obj.CreateTime) {
		return errors.New("field createTime is immutable")
	}
	if e.Type != _obj.Type {
		return errors.New("field type is immutable")
	}
	if e.WorkflowID != _obj.WorkflowID {
		return errors.New("field workflowID is immutable")
	}

	wspi.patchedEntity = _obj
	return nil
}

// WorkflowStepQueryInput holds the query input of the WorkflowStep entity,
// please tags with `path:",inline"` if embedding.
type WorkflowStepQueryInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to query WorkflowStep entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"project"`
	// Stage indicates to query WorkflowStep entity MUST under the Stage route.
	Stage *WorkflowStageQueryInput `path:",inline" query:"-" json:"stage"`

	// Refer holds the route path reference of the WorkflowStep entity.
	Refer *object.Refer `path:"workflowstep,default=" query:"-" json:"-"`
	// ID of the WorkflowStep entity.
	ID object.ID `path:"-" query:"-" json:"id"`
}

// Model returns the WorkflowStep entity for querying,
// after validating.
func (wsqi *WorkflowStepQueryInput) Model() *WorkflowStep {
	if wsqi == nil {
		return nil
	}

	return &WorkflowStep{
		ID: wsqi.ID,
	}
}

// Validate checks the WorkflowStepQueryInput entity.
func (wsqi *WorkflowStepQueryInput) Validate() error {
	if wsqi == nil {
		return errors.New("nil receiver")
	}

	return wsqi.ValidateWith(wsqi.inputConfig.Context, wsqi.inputConfig.Client, nil)
}

// ValidateWith checks the WorkflowStepQueryInput entity with the given context and client set.
func (wsqi *WorkflowStepQueryInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if wsqi == nil {
		return errors.New("nil receiver")
	}

	if wsqi.Refer != nil && *wsqi.Refer == "" {
		return fmt.Errorf("model: %s : %w", workflowstep.Label, ErrBlankResourceRefer)
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.WorkflowSteps().Query()

	// Validate when querying under the Project route.
	if wsqi.Project != nil {
		if err := wsqi.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				workflowstep.ProjectID(wsqi.Project.ID))
		}
	}

	// Validate when querying under the Stage route.
	if wsqi.Stage != nil {
		if err := wsqi.Stage.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			q.Where(
				workflowstep.WorkflowStageID(wsqi.Stage.ID))
		}
	}

	if wsqi.Refer != nil {
		if wsqi.Refer.IsID() {
			q.Where(
				workflowstep.ID(wsqi.Refer.ID()))
		} else {
			return errors.New("invalid identify refer of workflowstep")
		}
	} else if wsqi.ID != "" {
		q.Where(
			workflowstep.ID(wsqi.ID))
	} else {
		return errors.New("invalid identify of workflowstep")
	}

	q.Select(
		workflowstep.FieldID,
	)

	var e *WorkflowStep
	{
		// Get cache from previous validation.
		queryStmt, queryArgs := q.sqlQuery(setContextOp(ctx, q.ctx, "cache")).Query()
		ck := fmt.Sprintf("stmt=%v, args=%v", queryStmt, queryArgs)
		if cv, existed := cache[ck]; !existed {
			var err error
			e, err = q.Only(ctx)
			if err != nil {
				return err
			}

			// Set cache for other validation.
			cache[ck] = e
		} else {
			e = cv.(*WorkflowStep)
		}
	}

	wsqi.ID = e.ID
	return nil
}

// WorkflowStepQueryInputs holds the query input of the WorkflowStep entities,
// please tags with `path:",inline" query:",inline"` if embedding.
type WorkflowStepQueryInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to query WorkflowStep entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Stage indicates to query WorkflowStep entity MUST under the Stage route.
	Stage *WorkflowStageQueryInput `path:",inline" query:"-" json:"-"`
}

// Validate checks the WorkflowStepQueryInputs entity.
func (wsqi *WorkflowStepQueryInputs) Validate() error {
	if wsqi == nil {
		return errors.New("nil receiver")
	}

	return wsqi.ValidateWith(wsqi.inputConfig.Context, wsqi.inputConfig.Client, nil)
}

// ValidateWith checks the WorkflowStepQueryInputs entity with the given context and client set.
func (wsqi *WorkflowStepQueryInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if wsqi == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when querying under the Project route.
	if wsqi.Project != nil {
		if err := wsqi.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	// Validate when querying under the Stage route.
	if wsqi.Stage != nil {
		if err := wsqi.Stage.ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// WorkflowStepUpdateInput holds the modification input of the WorkflowStep entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type WorkflowStepUpdateInput struct {
	WorkflowStepQueryInput `path:",inline" query:"-" json:"-"`

	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes of the workflow step.
	Attributes map[string]interface{} `path:"-" query:"-" json:"attributes,omitempty"`
	// Inputs of the workflow step.
	Inputs map[string]interface{} `path:"-" query:"-" json:"inputs,omitempty"`
	// Outputs of the workflow step.
	Outputs map[string]interface{} `path:"-" query:"-" json:"outputs,omitempty"`
	// ID list of the workflow steps that this workflow step depends on.
	Dependencies []object.ID `path:"-" query:"-" json:"dependencies,omitempty"`
	// Retry policy of the workflow step.
	RetryStrategy *types.RetryStrategy `path:"-" query:"-" json:"retryStrategy,omitempty"`
	// Timeout seconds of the workflow step, 0 means no timeout.
	Timeout int `path:"-" query:"-" json:"timeout,omitempty"`
}

// Model returns the WorkflowStep entity for modifying,
// after validating.
func (wsui *WorkflowStepUpdateInput) Model() *WorkflowStep {
	if wsui == nil {
		return nil
	}

	_ws := &WorkflowStep{
		ID:            wsui.ID,
		Description:   wsui.Description,
		Labels:        wsui.Labels,
		Attributes:    wsui.Attributes,
		Inputs:        wsui.Inputs,
		Outputs:       wsui.Outputs,
		Dependencies:  wsui.Dependencies,
		RetryStrategy: wsui.RetryStrategy,
		Timeout:       wsui.Timeout,
	}

	return _ws
}

// Validate checks the WorkflowStepUpdateInput entity.
func (wsui *WorkflowStepUpdateInput) Validate() error {
	if wsui == nil {
		return errors.New("nil receiver")
	}

	return wsui.ValidateWith(wsui.inputConfig.Context, wsui.inputConfig.Client, nil)
}

// ValidateWith checks the WorkflowStepUpdateInput entity with the given context and client set.
func (wsui *WorkflowStepUpdateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := wsui.WorkflowStepQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	return nil
}

// WorkflowStepUpdateInputs holds the modification input item of the WorkflowStep entities.
type WorkflowStepUpdateInputsItem struct {
	// ID of the WorkflowStep entity.
	ID object.ID `path:"-" query:"-" json:"id"`

	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes of the workflow step.
	Attributes map[string]interface{} `path:"-" query:"-" json:"attributes,omitempty"`
	// Inputs of the workflow step.
	Inputs map[string]interface{} `path:"-" query:"-" json:"inputs,omitempty"`
	// Outputs of the workflow step.
	Outputs map[string]interface{} `path:"-" query:"-" json:"outputs,omitempty"`
	// ID list of the workflow steps that this workflow step depends on.
	Dependencies []object.ID `path:"-" query:"-" json:"dependencies"`
	// Retry policy of the workflow step.
	RetryStrategy *types.RetryStrategy `path:"-" query:"-" json:"retryStrategy,omitempty"`
	// Timeout seconds of the workflow step, 0 means no timeout.
	Timeout int `path:"-" query:"-" json:"timeout"`
}

// ValidateWith checks the WorkflowStepUpdateInputsItem entity with the given context and client set.
func (wsui *WorkflowStepUpdateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if wsui == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// WorkflowStepUpdateInputs holds the modification input of the WorkflowStep entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type WorkflowStepUpdateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to update WorkflowStep entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Stage indicates to update WorkflowStep entity MUST under the Stage route.
	Stage *WorkflowStageQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*WorkflowStepUpdateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the WorkflowStep entities for modifying,
// after validating.
func (wsui *WorkflowStepUpdateInputs) Model() []*WorkflowStep {
	if wsui == nil || len(wsui.Items) == 0 {
		return nil
	}

	_wss := make([]*WorkflowStep, len(wsui.Items))

	for i := range wsui.Items {
		_ws := &WorkflowStep{
			ID:            wsui.Items[i].ID,
			Description:   wsui.Items[i].Description,
			Labels:        wsui.Items[i].Labels,
			Attributes:    wsui.Items[i].Attributes,
			Inputs:        wsui.Items[i].Inputs,
			Outputs:       wsui.Items[i].Outputs,
			Dependencies:  wsui.Items[i].Dependencies,
			RetryStrategy: wsui.Items[i].RetryStrategy,
			Timeout:       wsui.Items[i].Timeout,
		}

		_wss[i] = _ws
	}

	return _wss
}

// IDs returns the ID list of the WorkflowStep entities for modifying,
// after validating.
func (wsui *WorkflowStepUpdateInputs) IDs() []object.ID {
	if wsui == nil || len(wsui.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(wsui.Items))
	for i := range wsui.Items {
		ids[i] = wsui.Items[i].ID
	}
	return ids
}

// Validate checks the WorkflowStepUpdateInputs entity.
func (wsui *WorkflowStepUpdateInputs) Validate() error {
	if wsui == nil {
		return errors.New("nil receiver")
	}

	return wsui.ValidateWith(wsui.inputConfig.Context, wsui.inputConfig.Client, nil)
}

// ValidateWith checks the WorkflowStepUpdateInputs entity with the given context and client set.
func (wsui *WorkflowStepUpdateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if wsui == nil {
		return errors.New("nil receiver")
	}

	if len(wsui.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.WorkflowSteps().Query()

	// Validate when updating under the Project route.
	if wsui.Project != nil {
		if err := wsui.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				workflowstep.ProjectID(wsui.Project.ID))
		}
	}

	// Validate when updating under the Stage route.
	if wsui.Stage != nil {
		if err := wsui.Stage.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			q.Where(
				workflowstep.WorkflowStageID(wsui.Stage.ID))
		}
	}

	ids := make([]object.ID, 0, len(wsui.Items))

	for i := range wsui.Items {
		if wsui.Items[i] == nil {
			return errors.New("nil item")
		}

		if wsui.Items[i].ID != "" {
			ids = append(ids, wsui.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	if len(ids) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	idsCnt, err := q.Where(workflowstep.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != cap(ids) {
		return errors.New("found unrecognized item")
	}

	for i := range wsui.Items {
		if err := wsui.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// WorkflowStepOutput holds the output of the WorkflowStep entity.
type WorkflowStepOutput struct {
	ID            object.ID              `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Labels        map[string]string      `json:"labels,omitempty"`
	CreateTime    *time.Time             `json:"createTime,omitempty"`
	UpdateTime    *time.Time             `json:"updateTime,omitempty"`
	Type          string                 `json:"type,omitempty"`
	WorkflowID    object.ID              `json:"workflowID,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
	Inputs        map[string]interface{} `json:"inputs,omitempty"`
	Outputs       map[string]interface{} `json:"outputs,omitempty"`
	Dependencies  []object.ID            `json:"dependencies,omitempty"`
	RetryStrategy *types.RetryStrategy   `json:"retryStrategy,omitempty"`
	Timeout       int                    `json:"timeout,omitempty"`

	Project *ProjectOutput       `json:"project,omitempty"`
	Stage   *WorkflowStageOutput `json:"stage,omitempty"`
}

// View returns the output of WorkflowStep entity.
func (_ws *WorkflowStep) View() *WorkflowStepOutput {
	return ExposeWorkflowStep(_ws)
}

// View returns the output of WorkflowStep entities.
func (_wss WorkflowSteps) View() []*WorkflowStepOutput {
	return ExposeWorkflowSteps(_wss)
}

// ExposeWorkflowStep converts the WorkflowStep to WorkflowStepOutput.
func ExposeWorkflowStep(_ws *WorkflowStep) *WorkflowStepOutput {
	if _ws == nil {
		return nil
	}

	wso := &WorkflowStepOutput{
		ID:            _ws.ID,
		Name:          _ws.Name,
		Description:   _ws.Description,
		Labels:        _ws.Labels,
		CreateTime:    _ws.CreateTime,
		UpdateTime:    _ws.UpdateTime,
		Type:          _ws.Type,
		WorkflowID:    _ws.WorkflowID,
		Attributes:    _ws.Attributes,
		Inputs:        _ws.Inputs,
		Outputs:       _ws.Outputs,
		Dependencies:  _ws.Dependencies,
		RetryStrategy: _ws.RetryStrategy,
		Timeout:       _ws.Timeout,
	}

	if _ws.Edges.Project != nil {
		wso.Project = ExposeProject(_ws.Edges.Project)
	} else if _ws.ProjectID != "" {
		wso.Project = &ProjectOutput{
			ID: _ws.ProjectID,
		}
	}
	if _ws.Edges.Stage != nil {
		wso.Stage = ExposeWorkflowStage(_ws.Edges.Stage)
	} else if _ws.WorkflowStageID != "" {
		wso.Stage = &WorkflowStageOutput{
			ID: _ws.WorkflowStageID,
		}
	}
	return wso
}

// ExposeWorkflowSteps converts the WorkflowStep slice to WorkflowStepOutput pointer slice.
func ExposeWorkflowSteps(_wss []*WorkflowStep) []*WorkflowStepOutput {
	if len(_wss) == 0 {
		return nil
	}

	wsos := make([]*WorkflowStepOutput, len(_wss))
	for i := range _wss {
		wsos[i] = ExposeWorkflowStep(_wss[i])
	}
	return wsos
}
