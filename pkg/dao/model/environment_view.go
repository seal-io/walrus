// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"time"

	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// EnvironmentCreateInput holds the creation input of the Environment entity.
type EnvironmentCreateInput struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Project *ProjectQueryInput `uri:",inline" query:"-" json:"project"`

	Name        string            `uri:"-" query:"-" json:"name"`
	Description string            `uri:"-" query:"-" json:"description,omitempty"`
	Labels      map[string]string `uri:"-" query:"-" json:"labels,omitempty"`

	Connectors []*EnvironmentConnectorRelationshipCreateInput `uri:"-" query:"-" json:"connectors,omitempty"`
}

// Model returns the Environment entity for creating,
// after validating.
func (eci *EnvironmentCreateInput) Model() *Environment {
	if eci == nil {
		return nil
	}

	_e := &Environment{
		Name:        eci.Name,
		Description: eci.Description,
		Labels:      eci.Labels,
	}

	if eci.Project != nil {
		_e.ProjectID = eci.Project.ID
	}

	for j := range eci.Connectors {
		if eci.Connectors[j] == nil {
			continue
		}
		_e.Edges.Connectors = append(_e.Edges.Connectors,
			eci.Connectors[j].Model())
	}
	return _e
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (eci *EnvironmentCreateInput) Load() error {
	if eci == nil {
		return errors.New("nil receiver")
	}

	return eci.LoadWith(eci.inputConfig.Context, eci.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (eci *EnvironmentCreateInput) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if eci == nil {
		return errors.New("nil receiver")
	}

	if eci.Project != nil {
		err = eci.Project.LoadWith(ctx, cs)
		if err != nil {
			return err
		}
	}

	for i := range eci.Connectors {
		if eci.Connectors[i] == nil {
			continue
		}
		err = eci.Connectors[i].LoadWith(ctx, cs)
		if err != nil {
			return err
		}
	}
	return nil
}

// EnvironmentCreateInputs holds the creation input item of the Environment entities.
type EnvironmentCreateInputsItem struct {
	Name        string            `uri:"-" query:"-" json:"name"`
	Description string            `uri:"-" query:"-" json:"description,omitempty"`
	Labels      map[string]string `uri:"-" query:"-" json:"labels,omitempty"`

	Connectors []*EnvironmentConnectorRelationshipCreateInput `uri:"-" query:"-" json:"connectors,omitempty"`
}

// EnvironmentCreateInputs holds the creation input of the Environment entities.
type EnvironmentCreateInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Project *ProjectQueryInput `uri:",inline" query:"-" json:"project"`

	Items []*EnvironmentCreateInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the Environment entities for creating,
// after validating.
func (eci *EnvironmentCreateInputs) Model() []*Environment {
	if eci == nil || len(eci.Items) == 0 {
		return nil
	}

	_es := make([]*Environment, len(eci.Items))

	for i := range eci.Items {
		_e := &Environment{
			Name:        eci.Items[i].Name,
			Description: eci.Items[i].Description,
			Labels:      eci.Items[i].Labels,
		}

		if eci.Project != nil {
			_e.ProjectID = eci.Project.ID
		}

		for j := range eci.Items[i].Connectors {
			if eci.Items[i].Connectors[j] == nil {
				continue
			}
			_e.Edges.Connectors = append(_e.Edges.Connectors,
				eci.Items[i].Connectors[j].Model())
		}

		_es[i] = _e
	}

	return _es
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (eci *EnvironmentCreateInputs) Load() error {
	if eci == nil {
		return errors.New("nil receiver")
	}

	return eci.LoadWith(eci.inputConfig.Context, eci.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (eci *EnvironmentCreateInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if eci == nil {
		return errors.New("nil receiver")
	}

	if len(eci.Items) == 0 {
		return errors.New("empty items")
	}

	if eci.Project != nil {
		err = eci.Project.LoadWith(ctx, cs)
		if err != nil {
			return err
		}
	}
	return nil
}

// EnvironmentDeleteInput holds the deletion input of the Environment entity.
type EnvironmentDeleteInput = EnvironmentQueryInput

// EnvironmentDeleteInputs holds the deletion input item of the Environment entities.
type EnvironmentDeleteInputsItem struct {
	ID object.ID `uri:"-" query:"-" json:"id"`
}

// EnvironmentDeleteInputs holds the deletion input of the Environment entities.
type EnvironmentDeleteInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Project *ProjectQueryInput `uri:",inline" query:"-" json:"project"`

	Items []*EnvironmentDeleteInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the Environment entities for deleting,
// after validating.
func (edi *EnvironmentDeleteInputs) Model() []*Environment {
	if edi == nil || len(edi.Items) == 0 {
		return nil
	}

	_es := make([]*Environment, len(edi.Items))
	for i := range edi.Items {
		_es[i] = &Environment{
			ID: edi.Items[i].ID,
		}
	}
	return _es
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (edi *EnvironmentDeleteInputs) Load() error {
	if edi == nil {
		return errors.New("nil receiver")
	}

	return edi.LoadWith(edi.inputConfig.Context, edi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (edi *EnvironmentDeleteInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if edi == nil {
		return errors.New("nil receiver")
	}

	if len(edi.Items) == 0 {
		return errors.New("empty items")
	}

	q := cs.Environments().Query()

	if edi.Project != nil {
		err = edi.Project.LoadWith(ctx, cs)
		if err != nil {
			return err
		}
		q.Where(
			environment.ProjectID(edi.Project.ID))
	}

	ids := make([]object.ID, 0, len(edi.Items))

	for i := range edi.Items {
		if edi.Items[i] == nil {
			return errors.New("nil item")
		}

		if edi.Items[i].ID != "" {
			ids = append(ids, edi.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	idsLen := len(ids)

	idsCnt, err := q.Where(environment.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != idsLen {
		return errors.New("found unrecognized item")
	}

	return nil
}

// EnvironmentQueryInput holds the query input of the Environment entity.
type EnvironmentQueryInput struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Project *ProjectQueryInput `uri:",inline" query:"-" json:"-"`

	Refer *object.Refer `uri:"environment,default=\"\"" query:"-" json:"-"`
	ID    object.ID     `uri:"id" query:"-" json:"id"` // TODO(thxCode): remove the uri:"id" after supporting hierarchical routes.
}

// Model returns the Environment entity for querying,
// after validating.
func (eqi *EnvironmentQueryInput) Model() *Environment {
	if eqi == nil {
		return nil
	}

	return &Environment{
		ID: eqi.ID,
	}
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (eqi *EnvironmentQueryInput) Load() error {
	if eqi == nil {
		return errors.New("nil receiver")
	}

	return eqi.LoadWith(eqi.inputConfig.Context, eqi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (eqi *EnvironmentQueryInput) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if eqi == nil {
		return errors.New("nil receiver")
	}

	if eqi.Refer != nil && *eqi.Refer == "" {
		return nil
	}

	q := cs.Environments().Query()

	if eqi.Project != nil {
		err = eqi.Project.LoadWith(ctx, cs)
		if err != nil {
			return err
		}
		q.Where(
			environment.ProjectID(eqi.Project.ID))
	}

	if eqi.Refer != nil {
		if eqi.Refer.IsID() {
			q.Where(
				environment.ID(eqi.Refer.ID()))
		} else {
			return errors.New("invalid identify refer of environment")
		}
	} else if eqi.ID != "" {
		q.Where(
			environment.ID(eqi.ID))
	} else {
		return errors.New("invalid identify of environment")
	}

	eqi.ID, err = q.OnlyID(ctx)
	return err
}

// EnvironmentQueryInputs holds the query input of the Environment entities.
type EnvironmentQueryInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Project *ProjectQueryInput `uri:",inline" query:"-" json:"project"`
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (eqi *EnvironmentQueryInputs) Load() error {
	if eqi == nil {
		return errors.New("nil receiver")
	}

	return eqi.LoadWith(eqi.inputConfig.Context, eqi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (eqi *EnvironmentQueryInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if eqi == nil {
		return errors.New("nil receiver")
	}

	if eqi.Project != nil {
		err = eqi.Project.LoadWith(ctx, cs)
		if err != nil {
			return err
		}
	}

	return err
}

// EnvironmentUpdateInput holds the modification input of the Environment entity.
type EnvironmentUpdateInput struct {
	EnvironmentQueryInput `uri:",inline" query:"-" json:",inline"`

	Name        string            `uri:"-" query:"-" json:"name,omitempty"`
	Description string            `uri:"-" query:"-" json:"description,omitempty"`
	Labels      map[string]string `uri:"-" query:"-" json:"labels,omitempty"`

	Connectors []*EnvironmentConnectorRelationshipUpdateInput `uri:"-" query:"-" json:"connectors,omitempty"`
}

// Model returns the Environment entity for modifying,
// after validating.
func (eui *EnvironmentUpdateInput) Model() *Environment {
	if eui == nil {
		return nil
	}

	_e := &Environment{
		ID:          eui.ID,
		Name:        eui.Name,
		Description: eui.Description,
		Labels:      eui.Labels,
	}

	for j := range eui.Connectors {
		if eui.Connectors[j] == nil {
			continue
		}
		_e.Edges.Connectors = append(_e.Edges.Connectors,
			eui.Connectors[j].Model())
	}
	return _e
}

// EnvironmentUpdateInputs holds the modification input item of the Environment entities.
type EnvironmentUpdateInputsItem struct {
	ID object.ID `uri:"-" query:"-" json:"id"`

	Name        string            `uri:"-" query:"-" json:"name,omitempty"`
	Description string            `uri:"-" query:"-" json:"description,omitempty"`
	Labels      map[string]string `uri:"-" query:"-" json:"labels,omitempty"`

	Connectors []*EnvironmentConnectorRelationshipUpdateInput `uri:"-" query:"-" json:"connectors,omitempty"`
}

// EnvironmentUpdateInputs holds the modification input of the Environment entities.
type EnvironmentUpdateInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Project *ProjectQueryInput `uri:",inline" query:"-" json:"project"`

	Items []*EnvironmentUpdateInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the Environment entities for modifying,
// after validating.
func (eui *EnvironmentUpdateInputs) Model() []*Environment {
	if eui == nil || len(eui.Items) == 0 {
		return nil
	}

	_es := make([]*Environment, len(eui.Items))

	for i := range eui.Items {
		_e := &Environment{
			ID:          eui.Items[i].ID,
			Name:        eui.Items[i].Name,
			Description: eui.Items[i].Description,
			Labels:      eui.Items[i].Labels,
		}

		for j := range eui.Items[i].Connectors {
			if eui.Items[i].Connectors[j] == nil {
				continue
			}
			_e.Edges.Connectors = append(_e.Edges.Connectors,
				eui.Items[i].Connectors[j].Model())
		}

		_es[i] = _e
	}

	return _es
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (eui *EnvironmentUpdateInputs) Load() error {
	if eui == nil {
		return errors.New("nil receiver")
	}

	return eui.LoadWith(eui.inputConfig.Context, eui.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (eui *EnvironmentUpdateInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if eui == nil {
		return errors.New("nil receiver")
	}

	if len(eui.Items) == 0 {
		return errors.New("empty items")
	}

	q := cs.Environments().Query()

	if eui.Project != nil {
		err = eui.Project.LoadWith(ctx, cs)
		if err != nil {
			return err
		}
		q.Where(
			environment.ProjectID(eui.Project.ID))
	}

	ids := make([]object.ID, 0, len(eui.Items))

	for i := range eui.Items {
		if eui.Items[i] == nil {
			return errors.New("nil item")
		}

		if eui.Items[i].ID != "" {
			ids = append(ids, eui.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	idsLen := len(ids)

	idsCnt, err := q.Where(environment.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != idsLen {
		return errors.New("found unrecognized item")
	}

	return nil
}

// EnvironmentOutput holds the output of the Environment entity.
type EnvironmentOutput struct {
	ID          object.ID         `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreateTime  *time.Time        `json:"createTime,omitempty"`
	UpdateTime  *time.Time        `json:"updateTime,omitempty"`

	Project    *ProjectOutput                            `json:"project,omitempty"`
	Connectors []*EnvironmentConnectorRelationshipOutput `json:"connectors,omitempty"`
}

// View returns the output of Environment.
func (_e *Environment) View() *EnvironmentOutput {
	return ExposeEnvironment(_e)
}

// View returns the output of Environments.
func (_es Environments) View() []*EnvironmentOutput {
	return ExposeEnvironments(_es)
}

// ExposeEnvironment converts the Environment to EnvironmentOutput.
func ExposeEnvironment(_e *Environment) *EnvironmentOutput {
	if _e == nil {
		return nil
	}

	eo := &EnvironmentOutput{
		ID:          _e.ID,
		Name:        _e.Name,
		Description: _e.Description,
		Labels:      _e.Labels,
		CreateTime:  _e.CreateTime,
		UpdateTime:  _e.UpdateTime,
	}

	if _e.Edges.Project != nil {
		eo.Project = ExposeProject(_e.Edges.Project)
	} else if _e.ProjectID != "" {
		eo.Project = &ProjectOutput{
			ID: _e.ProjectID,
		}
	}
	if _e.Edges.Connectors != nil {
		eo.Connectors = ExposeEnvironmentConnectorRelationships(_e.Edges.Connectors)
	}
	return eo
}

// ExposeEnvironments converts the Environment slice to EnvironmentOutput pointer slice.
func ExposeEnvironments(_es []*Environment) []*EnvironmentOutput {
	if len(_es) == 0 {
		return nil
	}

	eos := make([]*EnvironmentOutput, len(_es))
	for i := range _es {
		eos[i] = ExposeEnvironment(_es[i])
	}
	return eos
}
