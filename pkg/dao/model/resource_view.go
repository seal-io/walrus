// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/json"
)

// ResourceCreateInput holds the creation input of the Resource entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceCreateInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to create Resource entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to create Resource entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`

	// Name holds the value of the "name" field.
	Name string `path:"-" query:"-" json:"name"`
	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Type of the resource referring to a resource definition type.
	Type string `path:"-" query:"-" json:"type,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`
	// Change comment of the resource.
	ChangeComment string `path:"-" query:"-" json:"changeComment,omitempty"`

	// Template specifies full inserting the new TemplateVersion entity of the Resource entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template,omitempty"`
	// ResourceDefinition specifies full inserting the new ResourceDefinition entity of the Resource entity.
	ResourceDefinition *ResourceDefinitionQueryInput `uri:"-" query:"-" json:"-"`
}

// Model returns the Resource entity for creating,
// after validating.
func (rci *ResourceCreateInput) Model() *Resource {
	if rci == nil {
		return nil
	}

	_r := &Resource{
		Name:          rci.Name,
		Description:   rci.Description,
		Labels:        rci.Labels,
		Type:          rci.Type,
		Attributes:    rci.Attributes,
		ChangeComment: rci.ChangeComment,
	}

	if rci.Project != nil {
		_r.ProjectID = rci.Project.ID
	}
	if rci.Environment != nil {
		_r.EnvironmentID = rci.Environment.ID
	}

	if rci.Template != nil {
		_r.TemplateID = &rci.Template.ID
	}
	if rci.ResourceDefinition != nil {
		_r.ResourceDefinitionID = &rci.ResourceDefinition.ID
	}
	return _r
}

// Validate checks the ResourceCreateInput entity.
func (rci *ResourceCreateInput) Validate() error {
	if rci == nil {
		return errors.New("nil receiver")
	}

	return rci.ValidateWith(rci.inputConfig.Context, rci.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceCreateInput entity with the given context and client set.
func (rci *ResourceCreateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when creating under the Project route.
	if rci.Project != nil {
		if err := rci.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rci.Project = nil
			}
		}
	}
	// Validate when creating under the Environment route.
	if rci.Environment != nil {
		if err := rci.Environment.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rci.Environment = nil
			}
		}
	}

	if rci.Template != nil {
		if err := rci.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rci.Template = nil
			}
		}
	}

	if rci.ResourceDefinition != nil {
		if err := rci.ResourceDefinition.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rci.ResourceDefinition = nil
			}
		}
	}

	return nil
}

// ResourceCreateInputs holds the creation input item of the Resource entities.
type ResourceCreateInputsItem struct {
	// Name holds the value of the "name" field.
	Name string `path:"-" query:"-" json:"name"`
	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Type of the resource referring to a resource definition type.
	Type string `path:"-" query:"-" json:"type,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`
	// Change comment of the resource.
	ChangeComment string `path:"-" query:"-" json:"changeComment,omitempty"`

	// Template specifies full inserting the new TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template,omitempty"`
	// ResourceDefinition specifies full inserting the new ResourceDefinition entity.
	ResourceDefinition *ResourceDefinitionQueryInput `uri:"-" query:"-" json:"-"`
}

// ValidateWith checks the ResourceCreateInputsItem entity with the given context and client set.
func (rci *ResourceCreateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	if rci.Template != nil {
		if err := rci.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rci.Template = nil
			}
		}
	}

	if rci.ResourceDefinition != nil {
		if err := rci.ResourceDefinition.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rci.ResourceDefinition = nil
			}
		}
	}

	return nil
}

// ResourceCreateInputs holds the creation input of the Resource entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceCreateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to create Resource entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to create Resource entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*ResourceCreateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the Resource entities for creating,
// after validating.
func (rci *ResourceCreateInputs) Model() []*Resource {
	if rci == nil || len(rci.Items) == 0 {
		return nil
	}

	_rs := make([]*Resource, len(rci.Items))

	for i := range rci.Items {
		_r := &Resource{
			Name:          rci.Items[i].Name,
			Description:   rci.Items[i].Description,
			Labels:        rci.Items[i].Labels,
			Type:          rci.Items[i].Type,
			Attributes:    rci.Items[i].Attributes,
			ChangeComment: rci.Items[i].ChangeComment,
		}

		if rci.Project != nil {
			_r.ProjectID = rci.Project.ID
		}
		if rci.Environment != nil {
			_r.EnvironmentID = rci.Environment.ID
		}

		if rci.Items[i].Template != nil {
			_r.TemplateID = &rci.Items[i].Template.ID
		}
		if rci.Items[i].ResourceDefinition != nil {
			_r.ResourceDefinitionID = &rci.Items[i].ResourceDefinition.ID
		}

		_rs[i] = _r
	}

	return _rs
}

// Validate checks the ResourceCreateInputs entity .
func (rci *ResourceCreateInputs) Validate() error {
	if rci == nil {
		return errors.New("nil receiver")
	}

	return rci.ValidateWith(rci.inputConfig.Context, rci.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceCreateInputs entity with the given context and client set.
func (rci *ResourceCreateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rci == nil {
		return errors.New("nil receiver")
	}

	if len(rci.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when creating under the Project route.
	if rci.Project != nil {
		if err := rci.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rci.Project = nil
			}
		}
	}
	// Validate when creating under the Environment route.
	if rci.Environment != nil {
		if err := rci.Environment.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rci.Environment = nil
			}
		}
	}

	for i := range rci.Items {
		if rci.Items[i] == nil {
			continue
		}

		if err := rci.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// ResourceDeleteInput holds the deletion input of the Resource entity,
// please tags with `path:",inline"` if embedding.
type ResourceDeleteInput struct {
	ResourceQueryInput `path:",inline"`
}

// ResourceDeleteInputs holds the deletion input item of the Resource entities.
type ResourceDeleteInputsItem struct {
	// ID of the Resource entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the Resource entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`
}

// ResourceDeleteInputs holds the deletion input of the Resource entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceDeleteInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to delete Resource entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to delete Resource entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*ResourceDeleteInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the Resource entities for deleting,
// after validating.
func (rdi *ResourceDeleteInputs) Model() []*Resource {
	if rdi == nil || len(rdi.Items) == 0 {
		return nil
	}

	_rs := make([]*Resource, len(rdi.Items))
	for i := range rdi.Items {
		_rs[i] = &Resource{
			ID: rdi.Items[i].ID,
		}
	}
	return _rs
}

// IDs returns the ID list of the Resource entities for deleting,
// after validating.
func (rdi *ResourceDeleteInputs) IDs() []object.ID {
	if rdi == nil || len(rdi.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(rdi.Items))
	for i := range rdi.Items {
		ids[i] = rdi.Items[i].ID
	}
	return ids
}

// Validate checks the ResourceDeleteInputs entity.
func (rdi *ResourceDeleteInputs) Validate() error {
	if rdi == nil {
		return errors.New("nil receiver")
	}

	return rdi.ValidateWith(rdi.inputConfig.Context, rdi.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceDeleteInputs entity with the given context and client set.
func (rdi *ResourceDeleteInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rdi == nil {
		return errors.New("nil receiver")
	}

	if len(rdi.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.Resources().Query()

	// Validate when deleting under the Project route.
	if rdi.Project != nil {
		if err := rdi.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rdi.Project = nil
			}
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				resource.ProjectID(rdi.Project.ID))
		}
	}

	// Validate when deleting under the Environment route.
	if rdi.Environment != nil {
		if err := rdi.Environment.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rdi.Environment = nil
			}
		} else {
			q.Where(
				resource.EnvironmentID(rdi.Environment.ID))
		}
	}

	ids := make([]object.ID, 0, len(rdi.Items))
	ors := make([]predicate.Resource, 0, len(rdi.Items))
	indexers := make(map[any][]int)

	for i := range rdi.Items {
		if rdi.Items[i] == nil {
			return errors.New("nil item")
		}

		if rdi.Items[i].ID != "" {
			ids = append(ids, rdi.Items[i].ID)
			ors = append(ors, resource.ID(rdi.Items[i].ID))
			indexers[rdi.Items[i].ID] = append(indexers[rdi.Items[i].ID], i)
		} else if rdi.Items[i].Name != "" {
			ors = append(ors, resource.And(
				resource.Name(rdi.Items[i].Name)))
			indexerKey := fmt.Sprint("/", rdi.Items[i].Name)
			indexers[indexerKey] = append(indexers[indexerKey], i)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	p := resource.IDIn(ids...)
	if len(ids) != cap(ids) {
		p = resource.Or(ors...)
	}

	es, err := q.
		Where(p).
		Select(
			resource.FieldID,
			resource.FieldName,
		).
		All(ctx)
	if err != nil {
		return err
	}

	if len(es) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	for i := range es {
		indexer := indexers[es[i].ID]
		if indexer == nil {
			indexerKey := fmt.Sprint("/", es[i].Name)
			indexer = indexers[indexerKey]
		}
		for _, j := range indexer {
			rdi.Items[j].ID = es[i].ID
			rdi.Items[j].Name = es[i].Name
		}
	}

	return nil
}

// ResourcePatchInput holds the patch input of the Resource entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourcePatchInput struct {
	ResourceUpdateInput `path:",inline" query:"-" json:",inline"`

	patchedEntity *Resource `path:"-" query:"-" json:"-"`
}

// Model returns the Resource patched entity,
// after validating.
func (rpi *ResourcePatchInput) Model() *Resource {
	if rpi == nil {
		return nil
	}

	return rpi.patchedEntity
}

// Validate checks the ResourcePatchInput entity.
func (rpi *ResourcePatchInput) Validate() error {
	if rpi == nil {
		return errors.New("nil receiver")
	}

	return rpi.ValidateWith(rpi.inputConfig.Context, rpi.inputConfig.Client, nil)
}

// ValidateWith checks the ResourcePatchInput entity with the given context and client set.
func (rpi *ResourcePatchInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := rpi.ResourceUpdateInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	q := cs.Resources().Query()

	// Validate when querying under the Project route.
	if rpi.Project != nil {
		if err := rpi.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rpi.Project = nil
			}
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				resource.ProjectID(rpi.Project.ID))
		}
	}

	// Validate when querying under the Environment route.
	if rpi.Environment != nil {
		if err := rpi.Environment.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rpi.Environment = nil
			}
		} else {
			q.Where(
				resource.EnvironmentID(rpi.Environment.ID))
		}
	}

	if rpi.Refer != nil {
		if rpi.Refer.IsID() {
			q.Where(
				resource.ID(rpi.Refer.ID()))
		} else if refers := rpi.Refer.Split(1); len(refers) == 1 {
			q.Where(
				resource.Name(refers[0].String()))
		} else {
			return errors.New("invalid identify refer of resource")
		}
	} else if rpi.ID != "" {
		q.Where(
			resource.ID(rpi.ID))
	} else if rpi.Name != "" {
		q.Where(
			resource.Name(rpi.Name))
	} else {
		return errors.New("invalid identify of resource")
	}

	q.Select(
		resource.WithoutFields(
			resource.FieldAnnotations,
			resource.FieldCreateTime,
			resource.FieldUpdateTime,
			resource.FieldStatus,
		)...,
	)

	var e *Resource
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
			e = cv.(*Resource)
		}
	}

	_r := rpi.ResourceUpdateInput.Model()

	_obj, err := json.PatchObject(e, _r)
	if err != nil {
		return err
	}

	rpi.patchedEntity = _obj.(*Resource)
	return nil
}

// ResourceQueryInput holds the query input of the Resource entity,
// please tags with `path:",inline"` if embedding.
type ResourceQueryInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to query Resource entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"project"`
	// Environment indicates to query Resource entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"environment"`

	// Refer holds the route path reference of the Resource entity.
	Refer *object.Refer `path:"resource,default=" query:"-" json:"-"`
	// ID of the Resource entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the Resource entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`
}

// Model returns the Resource entity for querying,
// after validating.
func (rqi *ResourceQueryInput) Model() *Resource {
	if rqi == nil {
		return nil
	}

	return &Resource{
		ID:   rqi.ID,
		Name: rqi.Name,
	}
}

// Validate checks the ResourceQueryInput entity.
func (rqi *ResourceQueryInput) Validate() error {
	if rqi == nil {
		return errors.New("nil receiver")
	}

	return rqi.ValidateWith(rqi.inputConfig.Context, rqi.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceQueryInput entity with the given context and client set.
func (rqi *ResourceQueryInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rqi == nil {
		return errors.New("nil receiver")
	}

	if rqi.Refer != nil && *rqi.Refer == "" {
		return fmt.Errorf("model: %s : %w", resource.Label, ErrBlankResourceRefer)
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.Resources().Query()

	// Validate when querying under the Project route.
	if rqi.Project != nil {
		if err := rqi.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rqi.Project = nil
			}
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				resource.ProjectID(rqi.Project.ID))
		}
	}

	// Validate when querying under the Environment route.
	if rqi.Environment != nil {
		if err := rqi.Environment.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rqi.Environment = nil
			}
		} else {
			q.Where(
				resource.EnvironmentID(rqi.Environment.ID))
		}
	}

	if rqi.Refer != nil {
		if rqi.Refer.IsID() {
			q.Where(
				resource.ID(rqi.Refer.ID()))
		} else if refers := rqi.Refer.Split(1); len(refers) == 1 {
			q.Where(
				resource.Name(refers[0].String()))
		} else {
			return errors.New("invalid identify refer of resource")
		}
	} else if rqi.ID != "" {
		q.Where(
			resource.ID(rqi.ID))
	} else if rqi.Name != "" {
		q.Where(
			resource.Name(rqi.Name))
	} else {
		return errors.New("invalid identify of resource")
	}

	q.Select(
		resource.FieldID,
		resource.FieldName,
	)

	var e *Resource
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
			e = cv.(*Resource)
		}
	}

	rqi.ID = e.ID
	rqi.Name = e.Name
	return nil
}

// ResourceQueryInputs holds the query input of the Resource entities,
// please tags with `path:",inline" query:",inline"` if embedding.
type ResourceQueryInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to query Resource entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to query Resource entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`

	// Type of the resource referring to a resource definition type.
	Type string `path:"-" query:"type,omitempty" json:"-"`
}

// Validate checks the ResourceQueryInputs entity.
func (rqi *ResourceQueryInputs) Validate() error {
	if rqi == nil {
		return errors.New("nil receiver")
	}

	return rqi.ValidateWith(rqi.inputConfig.Context, rqi.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceQueryInputs entity with the given context and client set.
func (rqi *ResourceQueryInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rqi == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when querying under the Project route.
	if rqi.Project != nil {
		if err := rqi.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rqi.Project = nil
			}
		}
	}

	// Validate when querying under the Environment route.
	if rqi.Environment != nil {
		if err := rqi.Environment.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rqi.Environment = nil
			}
		}
	}

	return nil
}

// ResourceUpdateInput holds the modification input of the Resource entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceUpdateInput struct {
	ResourceQueryInput `path:",inline" query:"-" json:"-"`

	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`
	// Change comment of the resource.
	ChangeComment string `path:"-" query:"-" json:"changeComment,omitempty"`

	// Template indicates replacing the stale TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template,omitempty"`
	// ResourceDefinition indicates replacing the stale ResourceDefinition entity.
	ResourceDefinition *ResourceDefinitionQueryInput `uri:"-" query:"-" json:"-"`
}

// Model returns the Resource entity for modifying,
// after validating.
func (rui *ResourceUpdateInput) Model() *Resource {
	if rui == nil {
		return nil
	}

	_r := &Resource{
		ID:            rui.ID,
		Name:          rui.Name,
		Description:   rui.Description,
		Labels:        rui.Labels,
		Attributes:    rui.Attributes,
		ChangeComment: rui.ChangeComment,
	}

	if rui.Template != nil {
		_r.TemplateID = &rui.Template.ID
	}
	if rui.ResourceDefinition != nil {
		_r.ResourceDefinitionID = &rui.ResourceDefinition.ID
	}
	return _r
}

// Validate checks the ResourceUpdateInput entity.
func (rui *ResourceUpdateInput) Validate() error {
	if rui == nil {
		return errors.New("nil receiver")
	}

	return rui.ValidateWith(rui.inputConfig.Context, rui.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceUpdateInput entity with the given context and client set.
func (rui *ResourceUpdateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := rui.ResourceQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	if rui.Template != nil {
		if err := rui.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rui.Template = nil
			}
		}
	}

	if rui.ResourceDefinition != nil {
		if err := rui.ResourceDefinition.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rui.ResourceDefinition = nil
			}
		}
	}

	return nil
}

// ResourceUpdateInputs holds the modification input item of the Resource entities.
type ResourceUpdateInputsItem struct {
	// ID of the Resource entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the Resource entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`

	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`
	// Change comment of the resource.
	ChangeComment string `path:"-" query:"-" json:"changeComment,omitempty"`

	// Template indicates replacing the stale TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template,omitempty"`
	// ResourceDefinition indicates replacing the stale ResourceDefinition entity.
	ResourceDefinition *ResourceDefinitionQueryInput `uri:"-" query:"-" json:"-"`
}

// ValidateWith checks the ResourceUpdateInputsItem entity with the given context and client set.
func (rui *ResourceUpdateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rui == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	if rui.Template != nil {
		if err := rui.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rui.Template = nil
			}
		}
	}

	if rui.ResourceDefinition != nil {
		if err := rui.ResourceDefinition.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rui.ResourceDefinition = nil
			}
		}
	}

	return nil
}

// ResourceUpdateInputs holds the modification input of the Resource entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceUpdateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to update Resource entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to update Resource entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*ResourceUpdateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the Resource entities for modifying,
// after validating.
func (rui *ResourceUpdateInputs) Model() []*Resource {
	if rui == nil || len(rui.Items) == 0 {
		return nil
	}

	_rs := make([]*Resource, len(rui.Items))

	for i := range rui.Items {
		_r := &Resource{
			ID:            rui.Items[i].ID,
			Name:          rui.Items[i].Name,
			Description:   rui.Items[i].Description,
			Labels:        rui.Items[i].Labels,
			Attributes:    rui.Items[i].Attributes,
			ChangeComment: rui.Items[i].ChangeComment,
		}

		if rui.Items[i].Template != nil {
			_r.TemplateID = &rui.Items[i].Template.ID
		}
		if rui.Items[i].ResourceDefinition != nil {
			_r.ResourceDefinitionID = &rui.Items[i].ResourceDefinition.ID
		}

		_rs[i] = _r
	}

	return _rs
}

// IDs returns the ID list of the Resource entities for modifying,
// after validating.
func (rui *ResourceUpdateInputs) IDs() []object.ID {
	if rui == nil || len(rui.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(rui.Items))
	for i := range rui.Items {
		ids[i] = rui.Items[i].ID
	}
	return ids
}

// Validate checks the ResourceUpdateInputs entity.
func (rui *ResourceUpdateInputs) Validate() error {
	if rui == nil {
		return errors.New("nil receiver")
	}

	return rui.ValidateWith(rui.inputConfig.Context, rui.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceUpdateInputs entity with the given context and client set.
func (rui *ResourceUpdateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rui == nil {
		return errors.New("nil receiver")
	}

	if len(rui.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.Resources().Query()

	// Validate when updating under the Project route.
	if rui.Project != nil {
		if err := rui.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rui.Project = nil
			}
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				resource.ProjectID(rui.Project.ID))
		}
	}

	// Validate when updating under the Environment route.
	if rui.Environment != nil {
		if err := rui.Environment.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rui.Environment = nil
			}
		} else {
			q.Where(
				resource.EnvironmentID(rui.Environment.ID))
		}
	}

	ids := make([]object.ID, 0, len(rui.Items))
	ors := make([]predicate.Resource, 0, len(rui.Items))
	indexers := make(map[any][]int)

	for i := range rui.Items {
		if rui.Items[i] == nil {
			return errors.New("nil item")
		}

		if rui.Items[i].ID != "" {
			ids = append(ids, rui.Items[i].ID)
			ors = append(ors, resource.ID(rui.Items[i].ID))
			indexers[rui.Items[i].ID] = append(indexers[rui.Items[i].ID], i)
		} else if rui.Items[i].Name != "" {
			ors = append(ors, resource.And(
				resource.Name(rui.Items[i].Name)))
			indexerKey := fmt.Sprint("/", rui.Items[i].Name)
			indexers[indexerKey] = append(indexers[indexerKey], i)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	p := resource.IDIn(ids...)
	if len(ids) != cap(ids) {
		p = resource.Or(ors...)
	}

	es, err := q.
		Where(p).
		Select(
			resource.FieldID,
			resource.FieldName,
		).
		All(ctx)
	if err != nil {
		return err
	}

	if len(es) != cap(ids) {
		return errors.New("found unrecognized item")
	}

	for i := range es {
		indexer := indexers[es[i].ID]
		if indexer == nil {
			indexerKey := fmt.Sprint("/", es[i].Name)
			indexer = indexers[indexerKey]
		}
		for _, j := range indexer {
			rui.Items[j].ID = es[i].ID
			rui.Items[j].Name = es[i].Name
		}
	}

	for i := range rui.Items {
		if err := rui.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// ResourceOutput holds the output of the Resource entity.
type ResourceOutput struct {
	ID          object.ID         `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreateTime  *time.Time        `json:"createTime,omitempty"`
	UpdateTime  *time.Time        `json:"updateTime,omitempty"`
	Status      status.Status     `json:"status,omitempty"`
	Type        string            `json:"type,omitempty"`
	Attributes  property.Values   `json:"attributes,omitempty"`

	Project     *ProjectOutput         `json:"project,omitempty"`
	Environment *EnvironmentOutput     `json:"environment,omitempty"`
	Template    *TemplateVersionOutput `json:"template,omitempty"`
}

// View returns the output of Resource entity.
func (_r *Resource) View() *ResourceOutput {
	return ExposeResource(_r)
}

// View returns the output of Resource entities.
func (_rs Resources) View() []*ResourceOutput {
	return ExposeResources(_rs)
}

// ExposeResource converts the Resource to ResourceOutput.
func ExposeResource(_r *Resource) *ResourceOutput {
	if _r == nil {
		return nil
	}

	ro := &ResourceOutput{
		ID:          _r.ID,
		Name:        _r.Name,
		Description: _r.Description,
		Labels:      _r.Labels,
		CreateTime:  _r.CreateTime,
		UpdateTime:  _r.UpdateTime,
		Status:      _r.Status,
		Type:        _r.Type,
		Attributes:  _r.Attributes,
	}

	if _r.Edges.Project != nil {
		ro.Project = ExposeProject(_r.Edges.Project)
	} else if _r.ProjectID != "" {
		ro.Project = &ProjectOutput{
			ID: _r.ProjectID,
		}
	}
	if _r.Edges.Environment != nil {
		ro.Environment = ExposeEnvironment(_r.Edges.Environment)
	} else if _r.EnvironmentID != "" {
		ro.Environment = &EnvironmentOutput{
			ID: _r.EnvironmentID,
		}
	}
	if _r.Edges.Template != nil {
		ro.Template = ExposeTemplateVersion(_r.Edges.Template)
	} else if _r.TemplateID != nil {
		ro.Template = &TemplateVersionOutput{
			ID: *_r.TemplateID,
		}
	}
	return ro
}

// ExposeResources converts the Resource slice to ResourceOutput pointer slice.
func ExposeResources(_rs []*Resource) []*ResourceOutput {
	if len(_rs) == 0 {
		return nil
	}

	ros := make([]*ResourceOutput, len(_rs))
	for i := range _rs {
		ros[i] = ExposeResource(_rs[i])
	}
	return ros
}
