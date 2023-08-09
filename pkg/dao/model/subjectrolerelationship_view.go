// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"time"

	"github.com/seal-io/seal/pkg/dao/model/subjectrolerelationship"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// SubjectRoleRelationshipCreateInput holds the creation input of the SubjectRoleRelationship entity.
type SubjectRoleRelationshipCreateInput struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Subject *SubjectQueryInput `uri:"-" query:"-" json:"subject,omitempty"`
	Role    *RoleQueryInput    `uri:"-" query:"-" json:"role,omitempty"`
}

// Model returns the SubjectRoleRelationship entity for creating,
// after validating.
func (srrci *SubjectRoleRelationshipCreateInput) Model() *SubjectRoleRelationship {
	if srrci == nil {
		return nil
	}

	_srr := &SubjectRoleRelationship{}

	if srrci.Subject != nil {
		_srr.SubjectID = srrci.Subject.ID
	}
	if srrci.Role != nil {
		_srr.RoleID = srrci.Role.ID
	}
	return _srr
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (srrci *SubjectRoleRelationshipCreateInput) Load() error {
	if srrci == nil {
		return errors.New("nil receiver")
	}

	return srrci.LoadWith(srrci.inputConfig.Context, srrci.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (srrci *SubjectRoleRelationshipCreateInput) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if srrci == nil {
		return errors.New("nil receiver")
	}

	if srrci.Subject != nil {
		err = srrci.Subject.LoadWith(ctx, cs)
		if err != nil {
			return err
		}
	}
	if srrci.Role != nil {
		err = srrci.Role.LoadWith(ctx, cs)
		if err != nil {
			return err
		}
	}
	return nil
}

// SubjectRoleRelationshipCreateInputs holds the creation input item of the SubjectRoleRelationship entities.
type SubjectRoleRelationshipCreateInputsItem struct {
	Subject *SubjectQueryInput `uri:"-" query:"-" json:"subject,omitempty"`
	Role    *RoleQueryInput    `uri:"-" query:"-" json:"role,omitempty"`
}

// SubjectRoleRelationshipCreateInputs holds the creation input of the SubjectRoleRelationship entities.
type SubjectRoleRelationshipCreateInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Items []*SubjectRoleRelationshipCreateInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the SubjectRoleRelationship entities for creating,
// after validating.
func (srrci *SubjectRoleRelationshipCreateInputs) Model() []*SubjectRoleRelationship {
	if srrci == nil || len(srrci.Items) == 0 {
		return nil
	}

	_srrs := make([]*SubjectRoleRelationship, len(srrci.Items))

	for i := range srrci.Items {
		_srr := &SubjectRoleRelationship{}

		if srrci.Items[i].Subject != nil {
			_srr.SubjectID = srrci.Items[i].Subject.ID
		}
		if srrci.Items[i].Role != nil {
			_srr.RoleID = srrci.Items[i].Role.ID
		}

		_srrs[i] = _srr
	}

	return _srrs
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (srrci *SubjectRoleRelationshipCreateInputs) Load() error {
	if srrci == nil {
		return errors.New("nil receiver")
	}

	return srrci.LoadWith(srrci.inputConfig.Context, srrci.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (srrci *SubjectRoleRelationshipCreateInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if srrci == nil {
		return errors.New("nil receiver")
	}

	if len(srrci.Items) == 0 {
		return errors.New("empty items")
	}

	return nil
}

// SubjectRoleRelationshipDeleteInput holds the deletion input of the SubjectRoleRelationship entity.
type SubjectRoleRelationshipDeleteInput = SubjectRoleRelationshipQueryInput

// SubjectRoleRelationshipDeleteInputs holds the deletion input item of the SubjectRoleRelationship entities.
type SubjectRoleRelationshipDeleteInputsItem struct {
	ID object.ID `uri:"-" query:"-" json:"id"`
}

// SubjectRoleRelationshipDeleteInputs holds the deletion input of the SubjectRoleRelationship entities.
type SubjectRoleRelationshipDeleteInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Items []*SubjectRoleRelationshipDeleteInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the SubjectRoleRelationship entities for deleting,
// after validating.
func (srrdi *SubjectRoleRelationshipDeleteInputs) Model() []*SubjectRoleRelationship {
	if srrdi == nil || len(srrdi.Items) == 0 {
		return nil
	}

	_srrs := make([]*SubjectRoleRelationship, len(srrdi.Items))
	for i := range srrdi.Items {
		_srrs[i] = &SubjectRoleRelationship{
			ID: srrdi.Items[i].ID,
		}
	}
	return _srrs
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (srrdi *SubjectRoleRelationshipDeleteInputs) Load() error {
	if srrdi == nil {
		return errors.New("nil receiver")
	}

	return srrdi.LoadWith(srrdi.inputConfig.Context, srrdi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (srrdi *SubjectRoleRelationshipDeleteInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if srrdi == nil {
		return errors.New("nil receiver")
	}

	if len(srrdi.Items) == 0 {
		return errors.New("empty items")
	}

	q := cs.SubjectRoleRelationships().Query()

	ids := make([]object.ID, 0, len(srrdi.Items))

	for i := range srrdi.Items {
		if srrdi.Items[i] == nil {
			return errors.New("nil item")
		}

		if srrdi.Items[i].ID != "" {
			ids = append(ids, srrdi.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	idsLen := len(ids)

	idsCnt, err := q.Where(subjectrolerelationship.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != idsLen {
		return errors.New("found unrecognized item")
	}

	return nil
}

// SubjectRoleRelationshipQueryInput holds the query input of the SubjectRoleRelationship entity.
type SubjectRoleRelationshipQueryInput struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Refer *object.Refer `uri:"subjectrolerelationship,default=\"\"" query:"-" json:"-"`
	ID    object.ID     `uri:"id" query:"-" json:"id"` // TODO(thxCode): remove the uri:"id" after supporting hierarchical routes.
}

// Model returns the SubjectRoleRelationship entity for querying,
// after validating.
func (srrqi *SubjectRoleRelationshipQueryInput) Model() *SubjectRoleRelationship {
	if srrqi == nil {
		return nil
	}

	return &SubjectRoleRelationship{
		ID: srrqi.ID,
	}
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (srrqi *SubjectRoleRelationshipQueryInput) Load() error {
	if srrqi == nil {
		return errors.New("nil receiver")
	}

	return srrqi.LoadWith(srrqi.inputConfig.Context, srrqi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (srrqi *SubjectRoleRelationshipQueryInput) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if srrqi == nil {
		return errors.New("nil receiver")
	}

	if srrqi.Refer != nil && *srrqi.Refer == "" {
		return nil
	}

	q := cs.SubjectRoleRelationships().Query()

	if srrqi.Refer != nil {
		if srrqi.Refer.IsID() {
			q.Where(
				subjectrolerelationship.ID(srrqi.Refer.ID()))
		} else {
			return errors.New("invalid identify refer of subjectrolerelationship")
		}
	} else if srrqi.ID != "" {
		q.Where(
			subjectrolerelationship.ID(srrqi.ID))
	} else {
		return errors.New("invalid identify of subjectrolerelationship")
	}

	srrqi.ID, err = q.OnlyID(ctx)
	return err
}

// SubjectRoleRelationshipQueryInputs holds the query input of the SubjectRoleRelationship entities.
type SubjectRoleRelationshipQueryInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (srrqi *SubjectRoleRelationshipQueryInputs) Load() error {
	if srrqi == nil {
		return errors.New("nil receiver")
	}

	return srrqi.LoadWith(srrqi.inputConfig.Context, srrqi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (srrqi *SubjectRoleRelationshipQueryInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if srrqi == nil {
		return errors.New("nil receiver")
	}

	return err
}

// SubjectRoleRelationshipUpdateInput holds the modification input of the SubjectRoleRelationship entity.
type SubjectRoleRelationshipUpdateInput struct {
	SubjectRoleRelationshipQueryInput `uri:",inline" query:"-" json:",inline"`

	Subject *SubjectQueryInput `uri:"-" query:"-" json:"subject,omitempty"`
	Role    *RoleQueryInput    `uri:"-" query:"-" json:"role,omitempty"`
}

// Model returns the SubjectRoleRelationship entity for modifying,
// after validating.
func (srrui *SubjectRoleRelationshipUpdateInput) Model() *SubjectRoleRelationship {
	if srrui == nil {
		return nil
	}

	_srr := &SubjectRoleRelationship{
		ID: srrui.ID,
	}

	if srrui.Subject != nil {
		_srr.SubjectID = srrui.Subject.ID
	}
	if srrui.Role != nil {
		_srr.RoleID = srrui.Role.ID
	}
	return _srr
}

// SubjectRoleRelationshipUpdateInputs holds the modification input item of the SubjectRoleRelationship entities.
type SubjectRoleRelationshipUpdateInputsItem struct {
	ID object.ID `uri:"-" query:"-" json:"id"`

	Subject *SubjectQueryInput `uri:"-" query:"-" json:"subject,omitempty"`
	Role    *RoleQueryInput    `uri:"-" query:"-" json:"role,omitempty"`
}

// SubjectRoleRelationshipUpdateInputs holds the modification input of the SubjectRoleRelationship entities.
type SubjectRoleRelationshipUpdateInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Items []*SubjectRoleRelationshipUpdateInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the SubjectRoleRelationship entities for modifying,
// after validating.
func (srrui *SubjectRoleRelationshipUpdateInputs) Model() []*SubjectRoleRelationship {
	if srrui == nil || len(srrui.Items) == 0 {
		return nil
	}

	_srrs := make([]*SubjectRoleRelationship, len(srrui.Items))

	for i := range srrui.Items {
		_srr := &SubjectRoleRelationship{
			ID: srrui.Items[i].ID,
		}

		if srrui.Items[i].Subject != nil {
			_srr.SubjectID = srrui.Items[i].Subject.ID
		}
		if srrui.Items[i].Role != nil {
			_srr.RoleID = srrui.Items[i].Role.ID
		}

		_srrs[i] = _srr
	}

	return _srrs
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (srrui *SubjectRoleRelationshipUpdateInputs) Load() error {
	if srrui == nil {
		return errors.New("nil receiver")
	}

	return srrui.LoadWith(srrui.inputConfig.Context, srrui.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (srrui *SubjectRoleRelationshipUpdateInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if srrui == nil {
		return errors.New("nil receiver")
	}

	if len(srrui.Items) == 0 {
		return errors.New("empty items")
	}

	q := cs.SubjectRoleRelationships().Query()

	ids := make([]object.ID, 0, len(srrui.Items))

	for i := range srrui.Items {
		if srrui.Items[i] == nil {
			return errors.New("nil item")
		}

		if srrui.Items[i].ID != "" {
			ids = append(ids, srrui.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	idsLen := len(ids)

	idsCnt, err := q.Where(subjectrolerelationship.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != idsLen {
		return errors.New("found unrecognized item")
	}

	return nil
}

// SubjectRoleRelationshipOutput holds the output of the SubjectRoleRelationship entity.
type SubjectRoleRelationshipOutput struct {
	ID         object.ID  `json:"id,omitempty"`
	CreateTime *time.Time `json:"createTime,omitempty"`

	Project *ProjectOutput `json:"project,omitempty"`
	Subject *SubjectOutput `json:"subject,omitempty"`
	Role    *RoleOutput    `json:"role,omitempty"`
}

// View returns the output of SubjectRoleRelationship.
func (_srr *SubjectRoleRelationship) View() *SubjectRoleRelationshipOutput {
	return ExposeSubjectRoleRelationship(_srr)
}

// View returns the output of SubjectRoleRelationships.
func (_srrs SubjectRoleRelationships) View() []*SubjectRoleRelationshipOutput {
	return ExposeSubjectRoleRelationships(_srrs)
}

// ExposeSubjectRoleRelationship converts the SubjectRoleRelationship to SubjectRoleRelationshipOutput.
func ExposeSubjectRoleRelationship(_srr *SubjectRoleRelationship) *SubjectRoleRelationshipOutput {
	if _srr == nil {
		return nil
	}

	srro := &SubjectRoleRelationshipOutput{
		ID:         _srr.ID,
		CreateTime: _srr.CreateTime,
	}

	if _srr.Edges.Project != nil {
		srro.Project = ExposeProject(_srr.Edges.Project)
	} else if _srr.ProjectID != "" {
		srro.Project = &ProjectOutput{
			ID: _srr.ProjectID,
		}
	}
	if _srr.Edges.Subject != nil {
		srro.Subject = ExposeSubject(_srr.Edges.Subject)
	} else if _srr.SubjectID != "" {
		srro.Subject = &SubjectOutput{
			ID: _srr.SubjectID,
		}
	}
	if _srr.Edges.Role != nil {
		srro.Role = ExposeRole(_srr.Edges.Role)
	} else if _srr.RoleID != "" {
		srro.Role = &RoleOutput{
			ID: _srr.RoleID,
		}
	}
	return srro
}

// ExposeSubjectRoleRelationships converts the SubjectRoleRelationship slice to SubjectRoleRelationshipOutput pointer slice.
func ExposeSubjectRoleRelationships(_srrs []*SubjectRoleRelationship) []*SubjectRoleRelationshipOutput {
	if len(_srrs) == 0 {
		return nil
	}

	srros := make([]*SubjectRoleRelationshipOutput, len(_srrs))
	for i := range _srrs {
		srros[i] = ExposeSubjectRoleRelationship(_srrs[i])
	}
	return srros
}