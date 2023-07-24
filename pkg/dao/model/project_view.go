// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"time"

	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// ProjectCreateInput holds the creation input of the Project entity.
type ProjectCreateInput struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Name        string            `uri:"-" query:"-" json:"name"`
	Description string            `uri:"-" query:"-" json:"description,omitempty"`
	Labels      map[string]string `uri:"-" query:"-" json:"labels,omitempty"`
}

// Model returns the Project entity for creating,
// after validating.
func (pci *ProjectCreateInput) Model() *Project {
	if pci == nil {
		return nil
	}

	p := &Project{
		Name:        pci.Name,
		Description: pci.Description,
		Labels:      pci.Labels,
	}

	return p
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pci *ProjectCreateInput) Load() error {
	if pci == nil {
		return errors.New("nil receiver")
	}

	return pci.LoadWith(pci.inputConfig.Context, pci.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pci *ProjectCreateInput) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pci == nil {
		return errors.New("nil receiver")
	}

	return nil
}

// ProjectCreateInputs holds the creation input item of the Project entities.
type ProjectCreateInputsItem struct {
	Name        string            `uri:"-" query:"-" json:"name"`
	Description string            `uri:"-" query:"-" json:"description,omitempty"`
	Labels      map[string]string `uri:"-" query:"-" json:"labels,omitempty"`
}

// ProjectCreateInputs holds the creation input of the Project entities.
type ProjectCreateInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Items []*ProjectCreateInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the Project entities for creating,
// after validating.
func (pci *ProjectCreateInputs) Model() []*Project {
	if pci == nil || len(pci.Items) == 0 {
		return nil
	}

	ps := make([]*Project, len(pci.Items))

	for i := range pci.Items {
		p := &Project{
			Name:        pci.Items[i].Name,
			Description: pci.Items[i].Description,
			Labels:      pci.Items[i].Labels,
		}

		ps[i] = p
	}

	return ps
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pci *ProjectCreateInputs) Load() error {
	if pci == nil {
		return errors.New("nil receiver")
	}

	return pci.LoadWith(pci.inputConfig.Context, pci.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pci *ProjectCreateInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pci == nil {
		return errors.New("nil receiver")
	}

	if len(pci.Items) == 0 {
		return errors.New("empty items")
	}

	return nil
}

// ProjectDeleteInput holds the deletion input of the Project entity.
type ProjectDeleteInput = ProjectQueryInput

// ProjectDeleteInputs holds the deletion input item of the Project entities.
type ProjectDeleteInputsItem struct {
	ID object.ID `uri:"-" query:"-" json:"id"`
}

// ProjectDeleteInputs holds the deletion input of the Project entities.
type ProjectDeleteInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Items []*ProjectDeleteInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the Project entities for deleting,
// after validating.
func (pdi *ProjectDeleteInputs) Model() []*Project {
	if pdi == nil || len(pdi.Items) == 0 {
		return nil
	}

	ps := make([]*Project, len(pdi.Items))
	for i := range pdi.Items {
		ps[i] = &Project{
			ID: pdi.Items[i].ID,
		}
	}
	return ps
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pdi *ProjectDeleteInputs) Load() error {
	if pdi == nil {
		return errors.New("nil receiver")
	}

	return pdi.LoadWith(pdi.inputConfig.Context, pdi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pdi *ProjectDeleteInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pdi == nil {
		return errors.New("nil receiver")
	}

	if len(pdi.Items) == 0 {
		return errors.New("empty items")
	}

	q := cs.Projects().Query()

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

	idsCnt, err := q.Where(project.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != idsLen {
		return errors.New("found unrecognized item")
	}

	return nil
}

// ProjectQueryInput holds the query input of the Project entity.
type ProjectQueryInput struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Refer *object.Refer `uri:"project,default=\"\"" query:"-" json:"-"`
	ID    object.ID     `uri:"id" query:"-" json:"id"` // TODO(thxCode): remove the uri:"id" after supporting hierarchical routes.
}

// Model returns the Project entity for querying,
// after validating.
func (pqi *ProjectQueryInput) Model() *Project {
	if pqi == nil {
		return nil
	}

	return &Project{
		ID: pqi.ID,
	}
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pqi *ProjectQueryInput) Load() error {
	if pqi == nil {
		return errors.New("nil receiver")
	}

	return pqi.LoadWith(pqi.inputConfig.Context, pqi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pqi *ProjectQueryInput) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pqi == nil {
		return errors.New("nil receiver")
	}

	if pqi.Refer != nil && *pqi.Refer == "" {
		return nil
	}

	q := cs.Projects().Query()

	if pqi.Refer != nil {
		if pqi.Refer.IsID() {
			q.Where(
				project.ID(pqi.Refer.ID()))
		} else {
			return errors.New("invalid identify refer of project")
		}
	} else if pqi.ID != "" {
		q.Where(
			project.ID(pqi.ID))
	} else {
		return errors.New("invalid identify of project")
	}

	pqi.ID, err = q.OnlyID(ctx)
	return err
}

// ProjectQueryInputs holds the query input of the Project entities.
type ProjectQueryInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pqi *ProjectQueryInputs) Load() error {
	if pqi == nil {
		return errors.New("nil receiver")
	}

	return pqi.LoadWith(pqi.inputConfig.Context, pqi.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pqi *ProjectQueryInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pqi == nil {
		return errors.New("nil receiver")
	}

	return err
}

// ProjectUpdateInput holds the modification input of the Project entity.
type ProjectUpdateInput struct {
	ProjectQueryInput `uri:",inline" query:"-" json:",inline"`

	Name        string            `uri:"-" query:"-" json:"name,omitempty"`
	Description string            `uri:"-" query:"-" json:"description,omitempty"`
	Labels      map[string]string `uri:"-" query:"-" json:"labels,omitempty"`
}

// Model returns the Project entity for modifying,
// after validating.
func (pui *ProjectUpdateInput) Model() *Project {
	if pui == nil {
		return nil
	}

	p := &Project{
		ID:          pui.ID,
		Name:        pui.Name,
		Description: pui.Description,
		Labels:      pui.Labels,
	}

	return p
}

// ProjectUpdateInputs holds the modification input item of the Project entities.
type ProjectUpdateInputsItem struct {
	ID object.ID `uri:"-" query:"-" json:"id"`

	Name        string            `uri:"-" query:"-" json:"name,omitempty"`
	Description string            `uri:"-" query:"-" json:"description,omitempty"`
	Labels      map[string]string `uri:"-" query:"-" json:"labels,omitempty"`
}

// ProjectUpdateInputs holds the modification input of the Project entities.
type ProjectUpdateInputs struct {
	inputConfig `uri:"-" query:"-" json:"-"`

	Items []*ProjectUpdateInputsItem `uri:"-" query:"-" json:"items"`
}

// Model returns the Project entities for modifying,
// after validating.
func (pui *ProjectUpdateInputs) Model() []*Project {
	if pui == nil || len(pui.Items) == 0 {
		return nil
	}

	ps := make([]*Project, len(pui.Items))

	for i := range pui.Items {
		p := &Project{
			ID:          pui.Items[i].ID,
			Name:        pui.Items[i].Name,
			Description: pui.Items[i].Description,
			Labels:      pui.Items[i].Labels,
		}

		ps[i] = p
	}

	return ps
}

// Load checks the input.
// TODO(thxCode): rename to Validate after supporting hierarchical routes.
func (pui *ProjectUpdateInputs) Load() error {
	if pui == nil {
		return errors.New("nil receiver")
	}

	return pui.LoadWith(pui.inputConfig.Context, pui.inputConfig.ClientSet)
}

// LoadWith checks the input with the given context and client set.
// TODO(thxCode): rename to ValidateWith after supporting hierarchical routes.
func (pui *ProjectUpdateInputs) LoadWith(ctx context.Context, cs ClientSet) (err error) {
	if pui == nil {
		return errors.New("nil receiver")
	}

	if len(pui.Items) == 0 {
		return errors.New("empty items")
	}

	q := cs.Projects().Query()

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

	idsCnt, err := q.Where(project.IDIn(ids...)).
		Count(ctx)
	if err != nil {
		return err
	}

	if idsCnt != idsLen {
		return errors.New("found unrecognized item")
	}

	return nil
}

// ProjectOutput holds the output of the Project entity.
type ProjectOutput struct {
	ID          object.ID         `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreateTime  *time.Time        `json:"createTime,omitempty"`
	UpdateTime  *time.Time        `json:"updateTime,omitempty"`
}

// View returns the output of Project.
func (p *Project) View() *ProjectOutput {
	return ExposeProject(p)
}

// View returns the output of Projects.
func (ps Projects) View() []*ProjectOutput {
	return ExposeProjects(ps)
}

// ExposeProject converts the Project to ProjectOutput.
func ExposeProject(p *Project) *ProjectOutput {
	if p == nil {
		return nil
	}

	po := &ProjectOutput{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Labels:      p.Labels,
		CreateTime:  p.CreateTime,
		UpdateTime:  p.UpdateTime,
	}

	return po
}

// ExposeProjects converts the Project slice to ProjectOutput pointer slice.
func ExposeProjects(ps []*Project) []*ProjectOutput {
	if len(ps) == 0 {
		return nil
	}

	pos := make([]*ProjectOutput, len(ps))
	for i := range ps {
		pos[i] = ExposeProject(ps[i])
	}
	return pos
}
