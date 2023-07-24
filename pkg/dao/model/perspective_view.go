// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"time"

	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// PerspectiveCreateInput holds the creation input of the Perspective entity.
type PerspectiveCreateInput struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	EndTime           string                 `uri:"-" query:"-" json:"endTime"`
	StartTime         string                 `uri:"-" query:"-" json:"startTime"`
	Name              string                 `uri:"-" query:"-" json:"name"`
	Description       string                 `uri:"-" query:"-" json:"description,omitempty"`
	Labels            map[string]string      `uri:"-" query:"-" json:"labels,omitempty"`
	Builtin           bool                   `uri:"-" query:"-" json:"builtin,omitempty"`
	AllocationQueries []types.QueryCondition `uri:"-" query:"-" json:"allocationQueries,omitempty"`
}

// Model returns the Perspective entity for creating,
// after validating.
func (pci *PerspectiveCreateInput) Model() *Perspective {
	if pci == nil {
		return nil
	}

	p := &Perspective{
		EndTime:           pci.EndTime,
		StartTime:         pci.StartTime,
		Name:              pci.Name,
		Description:       pci.Description,
		Labels:            pci.Labels,
		Builtin:           pci.Builtin,
		AllocationQueries: pci.AllocationQueries,
	}

	return p
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pci *PerspectiveCreateInput) Load() error {
	if pci == nil {
		return errors.New("nil receiver")
	}

	return pci.LoadWith(pci.inputConfig.Context, pci.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pci *PerspectiveCreateInput) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pci == nil {
		return errors.New("nil receiver")
	}

	return nil
}

// PerspectiveCreateInputs holds the creation input item of the Perspective entities.
type PerspectiveCreateInputsItem struct {
	EndTime           string                 `uri:"-" query:"-" json:"endTime"`
	StartTime         string                 `uri:"-" query:"-" json:"startTime"`
	Name              string                 `uri:"-" query:"-" json:"name"`
	Description       string                 `uri:"-" query:"-" json:"description,omitempty"`
	Labels            map[string]string      `uri:"-" query:"-" json:"labels,omitempty"`
	Builtin           bool                   `uri:"-" query:"-" json:"builtin,omitempty"`
	AllocationQueries []types.QueryCondition `uri:"-" query:"-" json:"allocationQueries,omitempty"`
}

// PerspectiveCreateInputs holds the creation input of the Perspective entities.
type PerspectiveCreateInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Items []*PerspectiveCreateInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the Perspective entities for creating,
// after validating.
func (pci *PerspectiveCreateInputs) Model() []*Perspective {
	if pci == nil || len(pci.Items) == 0 {
		return nil
	}

	ps := make([]*Perspective, len(pci.Items))

	for i := range pci.Items {
		p := &Perspective{
			EndTime:           pci.Items[i].EndTime,
			StartTime:         pci.Items[i].StartTime,
			Name:              pci.Items[i].Name,
			Description:       pci.Items[i].Description,
			Labels:            pci.Items[i].Labels,
			Builtin:           pci.Items[i].Builtin,
			AllocationQueries: pci.Items[i].AllocationQueries,
		}

		ps[i] = p
	}

	return ps
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pci *PerspectiveCreateInputs) Load() error {
	if pci == nil {
		return errors.New("nil receiver")
	}

	return pci.LoadWith(pci.inputConfig.Context, pci.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pci *PerspectiveCreateInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pci == nil {
		return errors.New("nil receiver")
	}

	if len(pci.Items) == 0 {
		return errors.New("empty items")
	}

	return nil
}

// PerspectiveDeleteInput holds the deletion input of the Perspective entity.
type PerspectiveDeleteInput = PerspectiveQueryInput

// PerspectiveDeleteInputs holds the deletion input item of the Perspective entities.
type PerspectiveDeleteInputsItem struct {
	ID object.ID `uri:"-" query:"-" json:"id"`
}

// PerspectiveDeleteInputs holds the deletion input of the Perspective entities.
type PerspectiveDeleteInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Items []*PerspectiveDeleteInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the Perspective entities for deleting,
// after validating.
func (pdi *PerspectiveDeleteInputs) Model() []*Perspective {
	if pdi == nil || len(pdi.Items) == 0 {
		return nil
	}

	ps := make([]*Perspective, len(pdi.Items))
	for i := range pdi.Items {
		ps[i] = &Perspective{
			ID: pdi.Items[i].ID,
		}
	}
	return ps
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pdi *PerspectiveDeleteInputs) Load() error {
	if pdi == nil {
		return errors.New("nil receiver")
	}

	return pdi.LoadWith(pdi.inputConfig.Context, pdi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pdi *PerspectiveDeleteInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pdi == nil {
		return errors.New("nil receiver")
	}

	if len(pdi.Items) == 0 {
		return errors.New("empty items")
	}

	q := cs.Perspectives().Query()

	ids := make([]object.ID, 0, len(pdi.Items))

	for i := range pdi.Items {
		if pdi.Items[i] == nil {
			return errors.New("nil item")
		}

		if pdi.Items[i].ID != "" {
			ids = append(ids, pdi.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	idsLen := len(ids)

	idsCnt, err := q.Where(perspective.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != idsLen {
		return errors.New("found unrecognized item")
	}

	return nil
}

// PerspectiveQueryInput holds the query input of the Perspective entity.
type PerspectiveQueryInput struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Refer *object.Refer `uri:"perspective,default=\"\"" query:"-" json:"-"`
	ID    object.ID     `uri:"id" query:"-" json:"id"` // TODO(thxCode): remove the uri:"id" after supporting hierarchical routes.
}

// Model returns the Perspective entity for querying,
// after validating.
func (pqi *PerspectiveQueryInput) Model() *Perspective {
	if pqi == nil {
		return nil
	}

	return &Perspective{
		ID: pqi.ID,
	}
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pqi *PerspectiveQueryInput) Load() error {
	if pqi == nil {
		return errors.New("nil receiver")
	}

	return pqi.LoadWith(pqi.inputConfig.Context, pqi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pqi *PerspectiveQueryInput) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pqi == nil {
		return errors.New("nil receiver")
	}

	if pqi.Refer != nil && *pqi.Refer == "" {
		return nil
	}

	q := cs.Perspectives().Query()

	if pqi.Refer != nil {
		if pqi.Refer.IsID() {
			q.Where(
				perspective.ID(pqi.Refer.ID()))
		} else {
			return errors.New("invalid identify refer of perspective")
		}
	} else if pqi.ID != "" {
		q.Where(
			perspective.ID(pqi.ID))
	} else {
		return errors.New("invalid identify of perspective")
	}

	pqi.ID, err = q.OnlyID(ctx)
	return err
}

// PerspectiveQueryInputs holds the query input of the Perspective entities.
type PerspectiveQueryInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pqi *PerspectiveQueryInputs) Load() error {
	if pqi == nil {
		return errors.New("nil receiver")
	}

	return pqi.LoadWith(pqi.inputConfig.Context, pqi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pqi *PerspectiveQueryInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pqi == nil {
		return errors.New("nil receiver")
	}

	return err
}

// PerspectiveUpdateInput holds the modification input of the Perspective entity.
type PerspectiveUpdateInput struct {
	PerspectiveQueryInput `uri:",inline" query:"-" json:",inline"`

	Name              string                 `uri:"-" query:"-" json:"name,omitempty"`
	Description       string                 `uri:"-" query:"-" json:"description,omitempty"`
	Labels            map[string]string      `uri:"-" query:"-" json:"labels,omitempty"`
	StartTime         string                 `uri:"-" query:"-" json:"startTime,omitempty"`
	EndTime           string                 `uri:"-" query:"-" json:"endTime,omitempty"`
	Builtin           bool                   `uri:"-" query:"-" json:"builtin,omitempty"`
	AllocationQueries []types.QueryCondition `uri:"-" query:"-" json:"allocationQueries,omitempty"`
}

// Model returns the Perspective entity for modifying,
// after validating.
func (pui *PerspectiveUpdateInput) Model() *Perspective {
	if pui == nil {
		return nil
	}

	p := &Perspective{
		ID:                pui.ID,
		Name:              pui.Name,
		Description:       pui.Description,
		Labels:            pui.Labels,
		StartTime:         pui.StartTime,
		EndTime:           pui.EndTime,
		Builtin:           pui.Builtin,
		AllocationQueries: pui.AllocationQueries,
	}

	return p
}

// PerspectiveUpdateInputs holds the modification input item of the Perspective entities.
type PerspectiveUpdateInputsItem struct {
	ID object.ID `uri:"-" query:"-" json:"id"`

	Name              string                 `uri:"-" query:"-" json:"name,omitempty"`
	Description       string                 `uri:"-" query:"-" json:"description,omitempty"`
	Labels            map[string]string      `uri:"-" query:"-" json:"labels,omitempty"`
	StartTime         string                 `uri:"-" query:"-" json:"startTime,omitempty"`
	EndTime           string                 `uri:"-" query:"-" json:"endTime,omitempty"`
	Builtin           bool                   `uri:"-" query:"-" json:"builtin,omitempty"`
	AllocationQueries []types.QueryCondition `uri:"-" query:"-" json:"allocationQueries,omitempty"`
}

// PerspectiveUpdateInputs holds the modification input of the Perspective entities.
type PerspectiveUpdateInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Items []*PerspectiveUpdateInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the Perspective entities for modifying,
// after validating.
func (pui *PerspectiveUpdateInputs) Model() []*Perspective {
	if pui == nil || len(pui.Items) == 0 {
		return nil
	}

	ps := make([]*Perspective, len(pui.Items))

	for i := range pui.Items {
		p := &Perspective{
			ID:                pui.Items[i].ID,
			Name:              pui.Items[i].Name,
			Description:       pui.Items[i].Description,
			Labels:            pui.Items[i].Labels,
			StartTime:         pui.Items[i].StartTime,
			EndTime:           pui.Items[i].EndTime,
			Builtin:           pui.Items[i].Builtin,
			AllocationQueries: pui.Items[i].AllocationQueries,
		}

		ps[i] = p
	}

	return ps
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pui *PerspectiveUpdateInputs) Load() error {
	if pui == nil {
		return errors.New("nil receiver")
	}

	return pui.LoadWith(pui.inputConfig.Context, pui.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pui *PerspectiveUpdateInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pui == nil {
		return errors.New("nil receiver")
	}

	if len(pui.Items) == 0 {
		return errors.New("empty items")
	}

	q := cs.Perspectives().Query()

	ids := make([]object.ID, 0, len(pui.Items))

	for i := range pui.Items {
		if pui.Items[i] == nil {
			return errors.New("nil item")
		}

		if pui.Items[i].ID != "" {
			ids = append(ids, pui.Items[i].ID)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	idsLen := len(ids)

	idsCnt, err := q.Where(perspective.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != idsLen {
		return errors.New("found unrecognized item")
	}

	return nil
}

// PerspectiveOutput holds the output of the Perspective entity.
type PerspectiveOutput struct {
	ID                object.ID              `json:"id,omitempty"`
	Name              string                 `json:"name,omitempty"`
	Description       string                 `json:"description,omitempty"`
	Labels            map[string]string      `json:"labels,omitempty"`
	CreateTime        *time.Time             `json:"createTime,omitempty"`
	UpdateTime        *time.Time             `json:"updateTime,omitempty"`
	StartTime         string                 `json:"startTime,omitempty"`
	EndTime           string                 `json:"endTime,omitempty"`
	Builtin           bool                   `json:"builtin,omitempty"`
	AllocationQueries []types.QueryCondition `json:"allocationQueries,omitempty"`
}

// View returns the output of Perspective.
func (p *Perspective) View() *PerspectiveOutput {
	return ExposePerspective(p)
}

// View returns the output of Perspectives.
func (ps Perspectives) View() []*PerspectiveOutput {
	return ExposePerspectives(ps)
}

// ExposePerspective converts the Perspective to PerspectiveOutput.
func ExposePerspective(p *Perspective) *PerspectiveOutput {
	if p == nil {
		return nil
	}

	po := &PerspectiveOutput{
		ID:                p.ID,
		Name:              p.Name,
		Description:       p.Description,
		Labels:            p.Labels,
		CreateTime:        p.CreateTime,
		UpdateTime:        p.UpdateTime,
		StartTime:         p.StartTime,
		EndTime:           p.EndTime,
		Builtin:           p.Builtin,
		AllocationQueries: p.AllocationQueries,
	}

	return po
}

// ExposePerspectives converts the Perspective slice to PerspectiveOutput pointer slice.
func ExposePerspectives(ps []*Perspective) []*PerspectiveOutput {
	if len(ps) == 0 {
		return nil
	}

	pos := make([]*PerspectiveOutput, len(ps))
	for i := range ps {
		pos[i] = ExposePerspective(ps[i])
	}
	return pos
}
