// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/json"
)

// ResourceDefinitionMatchingRuleCreateInput holds the creation input of the ResourceDefinitionMatchingRule entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceDefinitionMatchingRuleCreateInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Resource selector.
	Selector types.Selector `path:"-" query:"-" json:"selector"`
	// Name of the matching rule.
	Name string `path:"-" query:"-" json:"name"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`

	// Template specifies full inserting the new TemplateVersion entity of the ResourceDefinitionMatchingRule entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template"`
}

// Model returns the ResourceDefinitionMatchingRule entity for creating,
// after validating.
func (rdmrci *ResourceDefinitionMatchingRuleCreateInput) Model() *ResourceDefinitionMatchingRule {
	if rdmrci == nil {
		return nil
	}

	_rdmr := &ResourceDefinitionMatchingRule{
		Selector:   rdmrci.Selector,
		Name:       rdmrci.Name,
		Attributes: rdmrci.Attributes,
	}

	if rdmrci.Template != nil {
		_rdmr.TemplateID = rdmrci.Template.ID
	}
	return _rdmr
}

// Validate checks the ResourceDefinitionMatchingRuleCreateInput entity.
func (rdmrci *ResourceDefinitionMatchingRuleCreateInput) Validate() error {
	if rdmrci == nil {
		return errors.New("nil receiver")
	}

	return rdmrci.ValidateWith(rdmrci.inputConfig.Context, rdmrci.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceDefinitionMatchingRuleCreateInput entity with the given context and client set.
func (rdmrci *ResourceDefinitionMatchingRuleCreateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rdmrci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	if rdmrci.Template != nil {
		if err := rdmrci.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rdmrci.Template = nil
			}
		}
	}

	return nil
}

// ResourceDefinitionMatchingRuleCreateInputs holds the creation input item of the ResourceDefinitionMatchingRule entities.
type ResourceDefinitionMatchingRuleCreateInputsItem struct {
	// Resource selector.
	Selector types.Selector `path:"-" query:"-" json:"selector"`
	// Name of the matching rule.
	Name string `path:"-" query:"-" json:"name"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`

	// Template specifies full inserting the new TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template"`
}

// ValidateWith checks the ResourceDefinitionMatchingRuleCreateInputsItem entity with the given context and client set.
func (rdmrci *ResourceDefinitionMatchingRuleCreateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rdmrci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	if rdmrci.Template != nil {
		if err := rdmrci.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rdmrci.Template = nil
			}
		}
	}

	return nil
}

// ResourceDefinitionMatchingRuleCreateInputs holds the creation input of the ResourceDefinitionMatchingRule entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceDefinitionMatchingRuleCreateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*ResourceDefinitionMatchingRuleCreateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the ResourceDefinitionMatchingRule entities for creating,
// after validating.
func (rdmrci *ResourceDefinitionMatchingRuleCreateInputs) Model() []*ResourceDefinitionMatchingRule {
	if rdmrci == nil || len(rdmrci.Items) == 0 {
		return nil
	}

	_rdmrs := make([]*ResourceDefinitionMatchingRule, len(rdmrci.Items))

	for i := range rdmrci.Items {
		_rdmr := &ResourceDefinitionMatchingRule{
			Selector:   rdmrci.Items[i].Selector,
			Name:       rdmrci.Items[i].Name,
			Attributes: rdmrci.Items[i].Attributes,
		}

		if rdmrci.Items[i].Template != nil {
			_rdmr.TemplateID = rdmrci.Items[i].Template.ID
		}

		_rdmrs[i] = _rdmr
	}

	return _rdmrs
}

// Validate checks the ResourceDefinitionMatchingRuleCreateInputs entity .
func (rdmrci *ResourceDefinitionMatchingRuleCreateInputs) Validate() error {
	if rdmrci == nil {
		return errors.New("nil receiver")
	}

	return rdmrci.ValidateWith(rdmrci.inputConfig.Context, rdmrci.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceDefinitionMatchingRuleCreateInputs entity with the given context and client set.
func (rdmrci *ResourceDefinitionMatchingRuleCreateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rdmrci == nil {
		return errors.New("nil receiver")
	}

	if len(rdmrci.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	for i := range rdmrci.Items {
		if rdmrci.Items[i] == nil {
			continue
		}

		if err := rdmrci.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// ResourceDefinitionMatchingRuleDeleteInput holds the deletion input of the ResourceDefinitionMatchingRule entity,
// please tags with `path:",inline"` if embedding.
type ResourceDefinitionMatchingRuleDeleteInput struct {
	ResourceDefinitionMatchingRuleQueryInput `path:",inline"`
}

// ResourceDefinitionMatchingRuleDeleteInputs holds the deletion input item of the ResourceDefinitionMatchingRule entities.
type ResourceDefinitionMatchingRuleDeleteInputsItem struct {
	// ID of the ResourceDefinitionMatchingRule entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the ResourceDefinitionMatchingRule entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`
}

// ResourceDefinitionMatchingRuleDeleteInputs holds the deletion input of the ResourceDefinitionMatchingRule entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceDefinitionMatchingRuleDeleteInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*ResourceDefinitionMatchingRuleDeleteInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the ResourceDefinitionMatchingRule entities for deleting,
// after validating.
func (rdmrdi *ResourceDefinitionMatchingRuleDeleteInputs) Model() []*ResourceDefinitionMatchingRule {
	if rdmrdi == nil || len(rdmrdi.Items) == 0 {
		return nil
	}

	_rdmrs := make([]*ResourceDefinitionMatchingRule, len(rdmrdi.Items))
	for i := range rdmrdi.Items {
		_rdmrs[i] = &ResourceDefinitionMatchingRule{
			ID: rdmrdi.Items[i].ID,
		}
	}
	return _rdmrs
}

// IDs returns the ID list of the ResourceDefinitionMatchingRule entities for deleting,
// after validating.
func (rdmrdi *ResourceDefinitionMatchingRuleDeleteInputs) IDs() []object.ID {
	if rdmrdi == nil || len(rdmrdi.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(rdmrdi.Items))
	for i := range rdmrdi.Items {
		ids[i] = rdmrdi.Items[i].ID
	}
	return ids
}

// Validate checks the ResourceDefinitionMatchingRuleDeleteInputs entity.
func (rdmrdi *ResourceDefinitionMatchingRuleDeleteInputs) Validate() error {
	if rdmrdi == nil {
		return errors.New("nil receiver")
	}

	return rdmrdi.ValidateWith(rdmrdi.inputConfig.Context, rdmrdi.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceDefinitionMatchingRuleDeleteInputs entity with the given context and client set.
func (rdmrdi *ResourceDefinitionMatchingRuleDeleteInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rdmrdi == nil {
		return errors.New("nil receiver")
	}

	if len(rdmrdi.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.ResourceDefinitionMatchingRules().Query()

	ids := make([]object.ID, 0, len(rdmrdi.Items))
	ors := make([]predicate.ResourceDefinitionMatchingRule, 0, len(rdmrdi.Items))
	indexers := make(map[any][]int)

	for i := range rdmrdi.Items {
		if rdmrdi.Items[i] == nil {
			return errors.New("nil item")
		}

		if rdmrdi.Items[i].ID != "" {
			ids = append(ids, rdmrdi.Items[i].ID)
			ors = append(ors, resourcedefinitionmatchingrule.ID(rdmrdi.Items[i].ID))
			indexers[rdmrdi.Items[i].ID] = append(indexers[rdmrdi.Items[i].ID], i)
		} else if rdmrdi.Items[i].Name != "" {
			ors = append(ors, resourcedefinitionmatchingrule.And(
				resourcedefinitionmatchingrule.Name(rdmrdi.Items[i].Name)))
			indexerKey := fmt.Sprint("/", rdmrdi.Items[i].Name)
			indexers[indexerKey] = append(indexers[indexerKey], i)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	p := resourcedefinitionmatchingrule.IDIn(ids...)
	if len(ids) != cap(ids) {
		p = resourcedefinitionmatchingrule.Or(ors...)
	}

	es, err := q.
		Where(p).
		Select(
			resourcedefinitionmatchingrule.FieldID,
			resourcedefinitionmatchingrule.FieldName,
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
			rdmrdi.Items[j].ID = es[i].ID
			rdmrdi.Items[j].Name = es[i].Name
		}
	}

	return nil
}

// ResourceDefinitionMatchingRulePatchInput holds the patch input of the ResourceDefinitionMatchingRule entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceDefinitionMatchingRulePatchInput struct {
	ResourceDefinitionMatchingRuleQueryInput `path:",inline" query:"-" json:"-"`

	// CreateTime holds the value of the "create_time" field.
	CreateTime *time.Time `path:"-" query:"-" json:"createTime,omitempty"`
	// Name of the matching rule.
	Name string `path:"-" query:"-" json:"name,omitempty"`
	// Resource selector.
	Selector types.Selector `path:"-" query:"-" json:"selector,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`
	// Order of the matching rule.
	Order int `path:"-" query:"-" json:"order,omitempty"`
	// Default value generated from resource definition's schema, ui schema and attributes
	SchemaDefaultValue []byte `path:"-" query:"-" json:"schemaDefaultValue,omitempty"`

	// Template indicates replacing the stale TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template"`

	patchedEntity *ResourceDefinitionMatchingRule `path:"-" query:"-" json:"-"`
}

// PatchModel returns the ResourceDefinitionMatchingRule partition entity for patching.
func (rdmrpi *ResourceDefinitionMatchingRulePatchInput) PatchModel() *ResourceDefinitionMatchingRule {
	if rdmrpi == nil {
		return nil
	}

	_rdmr := &ResourceDefinitionMatchingRule{
		CreateTime:         rdmrpi.CreateTime,
		Name:               rdmrpi.Name,
		Selector:           rdmrpi.Selector,
		Attributes:         rdmrpi.Attributes,
		Order:              rdmrpi.Order,
		SchemaDefaultValue: rdmrpi.SchemaDefaultValue,
	}

	if rdmrpi.Template != nil {
		_rdmr.TemplateID = rdmrpi.Template.ID
	}
	return _rdmr
}

// Model returns the ResourceDefinitionMatchingRule patched entity,
// after validating.
func (rdmrpi *ResourceDefinitionMatchingRulePatchInput) Model() *ResourceDefinitionMatchingRule {
	if rdmrpi == nil {
		return nil
	}

	return rdmrpi.patchedEntity
}

// Validate checks the ResourceDefinitionMatchingRulePatchInput entity.
func (rdmrpi *ResourceDefinitionMatchingRulePatchInput) Validate() error {
	if rdmrpi == nil {
		return errors.New("nil receiver")
	}

	return rdmrpi.ValidateWith(rdmrpi.inputConfig.Context, rdmrpi.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceDefinitionMatchingRulePatchInput entity with the given context and client set.
func (rdmrpi *ResourceDefinitionMatchingRulePatchInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := rdmrpi.ResourceDefinitionMatchingRuleQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	q := cs.ResourceDefinitionMatchingRules().Query()

	if rdmrpi.Template != nil {
		if err := rdmrpi.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rdmrpi.Template = nil
			}
		}
	}

	if rdmrpi.Refer != nil {
		if rdmrpi.Refer.IsID() {
			q.Where(
				resourcedefinitionmatchingrule.ID(rdmrpi.Refer.ID()))
		} else if refers := rdmrpi.Refer.Split(1); len(refers) == 1 {
			q.Where(
				resourcedefinitionmatchingrule.Name(refers[0].String()))
		} else {
			return errors.New("invalid identify refer of resourcedefinitionmatchingrule")
		}
	} else if rdmrpi.ID != "" {
		q.Where(
			resourcedefinitionmatchingrule.ID(rdmrpi.ID))
	} else if rdmrpi.Name != "" {
		q.Where(
			resourcedefinitionmatchingrule.Name(rdmrpi.Name))
	} else {
		return errors.New("invalid identify of resourcedefinitionmatchingrule")
	}

	q.Select(
		resourcedefinitionmatchingrule.WithoutFields(
			resourcedefinitionmatchingrule.FieldCreateTime,
			resourcedefinitionmatchingrule.FieldOrder,
			resourcedefinitionmatchingrule.FieldSchemaDefaultValue,
		)...,
	)

	var e *ResourceDefinitionMatchingRule
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
			e = cv.(*ResourceDefinitionMatchingRule)
		}
	}

	_pm := rdmrpi.PatchModel()

	_po, err := json.PatchObject(*e, *_pm)
	if err != nil {
		return err
	}

	_obj := _po.(*ResourceDefinitionMatchingRule)

	if !reflect.DeepEqual(e.CreateTime, _obj.CreateTime) {
		return errors.New("field createTime is immutable")
	}
	if e.Name != _obj.Name {
		return errors.New("field name is immutable")
	}

	rdmrpi.patchedEntity = _obj
	return nil
}

// ResourceDefinitionMatchingRuleQueryInput holds the query input of the ResourceDefinitionMatchingRule entity,
// please tags with `path:",inline"` if embedding.
type ResourceDefinitionMatchingRuleQueryInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Refer holds the route path reference of the ResourceDefinitionMatchingRule entity.
	Refer *object.Refer `path:"resourcedefinitionmatchingrule,default=" query:"-" json:"-"`
	// ID of the ResourceDefinitionMatchingRule entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the ResourceDefinitionMatchingRule entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`
}

// Model returns the ResourceDefinitionMatchingRule entity for querying,
// after validating.
func (rdmrqi *ResourceDefinitionMatchingRuleQueryInput) Model() *ResourceDefinitionMatchingRule {
	if rdmrqi == nil {
		return nil
	}

	return &ResourceDefinitionMatchingRule{
		ID:   rdmrqi.ID,
		Name: rdmrqi.Name,
	}
}

// Validate checks the ResourceDefinitionMatchingRuleQueryInput entity.
func (rdmrqi *ResourceDefinitionMatchingRuleQueryInput) Validate() error {
	if rdmrqi == nil {
		return errors.New("nil receiver")
	}

	return rdmrqi.ValidateWith(rdmrqi.inputConfig.Context, rdmrqi.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceDefinitionMatchingRuleQueryInput entity with the given context and client set.
func (rdmrqi *ResourceDefinitionMatchingRuleQueryInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rdmrqi == nil {
		return errors.New("nil receiver")
	}

	if rdmrqi.Refer != nil && *rdmrqi.Refer == "" {
		return fmt.Errorf("model: %s : %w", resourcedefinitionmatchingrule.Label, ErrBlankResourceRefer)
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.ResourceDefinitionMatchingRules().Query()

	if rdmrqi.Refer != nil {
		if rdmrqi.Refer.IsID() {
			q.Where(
				resourcedefinitionmatchingrule.ID(rdmrqi.Refer.ID()))
		} else if refers := rdmrqi.Refer.Split(1); len(refers) == 1 {
			q.Where(
				resourcedefinitionmatchingrule.Name(refers[0].String()))
		} else {
			return errors.New("invalid identify refer of resourcedefinitionmatchingrule")
		}
	} else if rdmrqi.ID != "" {
		q.Where(
			resourcedefinitionmatchingrule.ID(rdmrqi.ID))
	} else if rdmrqi.Name != "" {
		q.Where(
			resourcedefinitionmatchingrule.Name(rdmrqi.Name))
	} else {
		return errors.New("invalid identify of resourcedefinitionmatchingrule")
	}

	q.Select(
		resourcedefinitionmatchingrule.FieldID,
		resourcedefinitionmatchingrule.FieldName,
	)

	var e *ResourceDefinitionMatchingRule
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
			e = cv.(*ResourceDefinitionMatchingRule)
		}
	}

	rdmrqi.ID = e.ID
	rdmrqi.Name = e.Name
	return nil
}

// ResourceDefinitionMatchingRuleQueryInputs holds the query input of the ResourceDefinitionMatchingRule entities,
// please tags with `path:",inline" query:",inline"` if embedding.
type ResourceDefinitionMatchingRuleQueryInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`
}

// Validate checks the ResourceDefinitionMatchingRuleQueryInputs entity.
func (rdmrqi *ResourceDefinitionMatchingRuleQueryInputs) Validate() error {
	if rdmrqi == nil {
		return errors.New("nil receiver")
	}

	return rdmrqi.ValidateWith(rdmrqi.inputConfig.Context, rdmrqi.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceDefinitionMatchingRuleQueryInputs entity with the given context and client set.
func (rdmrqi *ResourceDefinitionMatchingRuleQueryInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rdmrqi == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// ResourceDefinitionMatchingRuleUpdateInput holds the modification input of the ResourceDefinitionMatchingRule entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceDefinitionMatchingRuleUpdateInput struct {
	ResourceDefinitionMatchingRuleQueryInput `path:",inline" query:"-" json:"-"`

	// Name of the matching rule.
	Name string `path:"-" query:"-" json:"name,omitempty"`
	// Resource selector.
	Selector types.Selector `path:"-" query:"-" json:"selector,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`

	// Template indicates replacing the stale TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template"`
}

// Model returns the ResourceDefinitionMatchingRule entity for modifying,
// after validating.
func (rdmrui *ResourceDefinitionMatchingRuleUpdateInput) Model() *ResourceDefinitionMatchingRule {
	if rdmrui == nil {
		return nil
	}

	_rdmr := &ResourceDefinitionMatchingRule{
		ID:         rdmrui.ID,
		Name:       rdmrui.Name,
		Selector:   rdmrui.Selector,
		Attributes: rdmrui.Attributes,
	}

	if rdmrui.Template != nil {
		_rdmr.TemplateID = rdmrui.Template.ID
	}
	return _rdmr
}

// Validate checks the ResourceDefinitionMatchingRuleUpdateInput entity.
func (rdmrui *ResourceDefinitionMatchingRuleUpdateInput) Validate() error {
	if rdmrui == nil {
		return errors.New("nil receiver")
	}

	return rdmrui.ValidateWith(rdmrui.inputConfig.Context, rdmrui.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceDefinitionMatchingRuleUpdateInput entity with the given context and client set.
func (rdmrui *ResourceDefinitionMatchingRuleUpdateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := rdmrui.ResourceDefinitionMatchingRuleQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	if rdmrui.Template != nil {
		if err := rdmrui.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rdmrui.Template = nil
			}
		}
	}

	return nil
}

// ResourceDefinitionMatchingRuleUpdateInputs holds the modification input item of the ResourceDefinitionMatchingRule entities.
type ResourceDefinitionMatchingRuleUpdateInputsItem struct {
	// ID of the ResourceDefinitionMatchingRule entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the ResourceDefinitionMatchingRule entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`

	// Resource selector.
	Selector types.Selector `path:"-" query:"-" json:"selector"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`

	// Template indicates replacing the stale TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template"`
}

// ValidateWith checks the ResourceDefinitionMatchingRuleUpdateInputsItem entity with the given context and client set.
func (rdmrui *ResourceDefinitionMatchingRuleUpdateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rdmrui == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	if rdmrui.Template != nil {
		if err := rdmrui.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				rdmrui.Template = nil
			}
		}
	}

	return nil
}

// ResourceDefinitionMatchingRuleUpdateInputs holds the modification input of the ResourceDefinitionMatchingRule entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type ResourceDefinitionMatchingRuleUpdateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*ResourceDefinitionMatchingRuleUpdateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the ResourceDefinitionMatchingRule entities for modifying,
// after validating.
func (rdmrui *ResourceDefinitionMatchingRuleUpdateInputs) Model() []*ResourceDefinitionMatchingRule {
	if rdmrui == nil || len(rdmrui.Items) == 0 {
		return nil
	}

	_rdmrs := make([]*ResourceDefinitionMatchingRule, len(rdmrui.Items))

	for i := range rdmrui.Items {
		_rdmr := &ResourceDefinitionMatchingRule{
			ID:         rdmrui.Items[i].ID,
			Name:       rdmrui.Items[i].Name,
			Selector:   rdmrui.Items[i].Selector,
			Attributes: rdmrui.Items[i].Attributes,
		}

		if rdmrui.Items[i].Template != nil {
			_rdmr.TemplateID = rdmrui.Items[i].Template.ID
		}

		_rdmrs[i] = _rdmr
	}

	return _rdmrs
}

// IDs returns the ID list of the ResourceDefinitionMatchingRule entities for modifying,
// after validating.
func (rdmrui *ResourceDefinitionMatchingRuleUpdateInputs) IDs() []object.ID {
	if rdmrui == nil || len(rdmrui.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(rdmrui.Items))
	for i := range rdmrui.Items {
		ids[i] = rdmrui.Items[i].ID
	}
	return ids
}

// Validate checks the ResourceDefinitionMatchingRuleUpdateInputs entity.
func (rdmrui *ResourceDefinitionMatchingRuleUpdateInputs) Validate() error {
	if rdmrui == nil {
		return errors.New("nil receiver")
	}

	return rdmrui.ValidateWith(rdmrui.inputConfig.Context, rdmrui.inputConfig.Client, nil)
}

// ValidateWith checks the ResourceDefinitionMatchingRuleUpdateInputs entity with the given context and client set.
func (rdmrui *ResourceDefinitionMatchingRuleUpdateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if rdmrui == nil {
		return errors.New("nil receiver")
	}

	if len(rdmrui.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.ResourceDefinitionMatchingRules().Query()

	ids := make([]object.ID, 0, len(rdmrui.Items))
	ors := make([]predicate.ResourceDefinitionMatchingRule, 0, len(rdmrui.Items))
	indexers := make(map[any][]int)

	for i := range rdmrui.Items {
		if rdmrui.Items[i] == nil {
			return errors.New("nil item")
		}

		if rdmrui.Items[i].ID != "" {
			ids = append(ids, rdmrui.Items[i].ID)
			ors = append(ors, resourcedefinitionmatchingrule.ID(rdmrui.Items[i].ID))
			indexers[rdmrui.Items[i].ID] = append(indexers[rdmrui.Items[i].ID], i)
		} else if rdmrui.Items[i].Name != "" {
			ors = append(ors, resourcedefinitionmatchingrule.And(
				resourcedefinitionmatchingrule.Name(rdmrui.Items[i].Name)))
			indexerKey := fmt.Sprint("/", rdmrui.Items[i].Name)
			indexers[indexerKey] = append(indexers[indexerKey], i)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	p := resourcedefinitionmatchingrule.IDIn(ids...)
	if len(ids) != cap(ids) {
		p = resourcedefinitionmatchingrule.Or(ors...)
	}

	es, err := q.
		Where(p).
		Select(
			resourcedefinitionmatchingrule.FieldID,
			resourcedefinitionmatchingrule.FieldName,
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
			rdmrui.Items[j].ID = es[i].ID
			rdmrui.Items[j].Name = es[i].Name
		}
	}

	for i := range rdmrui.Items {
		if err := rdmrui.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// ResourceDefinitionMatchingRuleOutput holds the output of the ResourceDefinitionMatchingRule entity.
type ResourceDefinitionMatchingRuleOutput struct {
	ID         object.ID       `json:"id,omitempty"`
	CreateTime *time.Time      `json:"createTime,omitempty"`
	Name       string          `json:"name,omitempty"`
	Selector   types.Selector  `json:"selector,omitempty"`
	Attributes property.Values `json:"attributes,omitempty"`

	Template *TemplateVersionOutput `json:"template,omitempty"`
}

// View returns the output of ResourceDefinitionMatchingRule entity.
func (_rdmr *ResourceDefinitionMatchingRule) View() *ResourceDefinitionMatchingRuleOutput {
	return ExposeResourceDefinitionMatchingRule(_rdmr)
}

// View returns the output of ResourceDefinitionMatchingRule entities.
func (_rdmrs ResourceDefinitionMatchingRules) View() []*ResourceDefinitionMatchingRuleOutput {
	return ExposeResourceDefinitionMatchingRules(_rdmrs)
}

// ExposeResourceDefinitionMatchingRule converts the ResourceDefinitionMatchingRule to ResourceDefinitionMatchingRuleOutput.
func ExposeResourceDefinitionMatchingRule(_rdmr *ResourceDefinitionMatchingRule) *ResourceDefinitionMatchingRuleOutput {
	if _rdmr == nil {
		return nil
	}

	rdmro := &ResourceDefinitionMatchingRuleOutput{
		ID:         _rdmr.ID,
		CreateTime: _rdmr.CreateTime,
		Name:       _rdmr.Name,
		Selector:   _rdmr.Selector,
		Attributes: _rdmr.Attributes,
	}

	if _rdmr.Edges.Template != nil {
		rdmro.Template = ExposeTemplateVersion(_rdmr.Edges.Template)
	} else if _rdmr.TemplateID != "" {
		rdmro.Template = &TemplateVersionOutput{
			ID: _rdmr.TemplateID,
		}
	}
	return rdmro
}

// ExposeResourceDefinitionMatchingRules converts the ResourceDefinitionMatchingRule slice to ResourceDefinitionMatchingRuleOutput pointer slice.
func ExposeResourceDefinitionMatchingRules(_rdmrs []*ResourceDefinitionMatchingRule) []*ResourceDefinitionMatchingRuleOutput {
	if len(_rdmrs) == 0 {
		return nil
	}

	rdmros := make([]*ResourceDefinitionMatchingRuleOutput, len(_rdmrs))
	for i := range _rdmrs {
		rdmros[i] = ExposeResourceDefinitionMatchingRule(_rdmrs[i])
	}
	return rdmros
}
