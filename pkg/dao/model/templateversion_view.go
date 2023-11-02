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
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// TemplateVersionCreateInput holds the creation input of the TemplateVersion entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type TemplateVersionCreateInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to create TemplateVersion entity CAN under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`

	// Source of the template.
	Source string `path:"-" query:"-" json:"source"`
	// Version of the template.
	Version string `path:"-" query:"-" json:"version"`
	// Name of the template.
	Name string `path:"-" query:"-" json:"name"`
	// Generated schema and data of the template.
	Schema types.Schema `path:"-" query:"-" json:"schema,omitempty"`
	// ui schema of the template.
	UiSchema types.UISchema `path:"-" query:"-" json:"uiSchema,omitempty"`
}

// Model returns the TemplateVersion entity for creating,
// after validating.
func (tvci *TemplateVersionCreateInput) Model() *TemplateVersion {
	if tvci == nil {
		return nil
	}

	_tv := &TemplateVersion{
		Source:   tvci.Source,
		Version:  tvci.Version,
		Name:     tvci.Name,
		Schema:   tvci.Schema,
		UiSchema: tvci.UiSchema,
	}

	if tvci.Project != nil {
		_tv.ProjectID = tvci.Project.ID
	}

	return _tv
}

// Validate checks the TemplateVersionCreateInput entity.
func (tvci *TemplateVersionCreateInput) Validate() error {
	if tvci == nil {
		return errors.New("nil receiver")
	}

	return tvci.ValidateWith(tvci.inputConfig.Context, tvci.inputConfig.Client, nil)
}

// ValidateWith checks the TemplateVersionCreateInput entity with the given context and client set.
func (tvci *TemplateVersionCreateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tvci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when creating under the Project route.
	if tvci.Project != nil {
		if err := tvci.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				tvci.Project = nil
			}
		}
	}

	return nil
}

// TemplateVersionCreateInputs holds the creation input item of the TemplateVersion entities.
type TemplateVersionCreateInputsItem struct {
	// Source of the template.
	Source string `path:"-" query:"-" json:"source"`
	// Version of the template.
	Version string `path:"-" query:"-" json:"version"`
	// Name of the template.
	Name string `path:"-" query:"-" json:"name"`
	// Generated schema and data of the template.
	Schema types.Schema `path:"-" query:"-" json:"schema,omitempty"`
	// ui schema of the template.
	UiSchema types.UISchema `path:"-" query:"-" json:"uiSchema,omitempty"`
}

// ValidateWith checks the TemplateVersionCreateInputsItem entity with the given context and client set.
func (tvci *TemplateVersionCreateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tvci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// TemplateVersionCreateInputs holds the creation input of the TemplateVersion entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type TemplateVersionCreateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to create TemplateVersion entity CAN under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*TemplateVersionCreateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the TemplateVersion entities for creating,
// after validating.
func (tvci *TemplateVersionCreateInputs) Model() []*TemplateVersion {
	if tvci == nil || len(tvci.Items) == 0 {
		return nil
	}

	_tvs := make([]*TemplateVersion, len(tvci.Items))

	for i := range tvci.Items {
		_tv := &TemplateVersion{
			Source:   tvci.Items[i].Source,
			Version:  tvci.Items[i].Version,
			Name:     tvci.Items[i].Name,
			Schema:   tvci.Items[i].Schema,
			UiSchema: tvci.Items[i].UiSchema,
		}

		if tvci.Project != nil {
			_tv.ProjectID = tvci.Project.ID
		}

		_tvs[i] = _tv
	}

	return _tvs
}

// Validate checks the TemplateVersionCreateInputs entity .
func (tvci *TemplateVersionCreateInputs) Validate() error {
	if tvci == nil {
		return errors.New("nil receiver")
	}

	return tvci.ValidateWith(tvci.inputConfig.Context, tvci.inputConfig.Client, nil)
}

// ValidateWith checks the TemplateVersionCreateInputs entity with the given context and client set.
func (tvci *TemplateVersionCreateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tvci == nil {
		return errors.New("nil receiver")
	}

	if len(tvci.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when creating under the Project route.
	if tvci.Project != nil {
		if err := tvci.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				tvci.Project = nil
			}
		}
	}

	for i := range tvci.Items {
		if tvci.Items[i] == nil {
			continue
		}

		if err := tvci.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// TemplateVersionDeleteInput holds the deletion input of the TemplateVersion entity,
// please tags with `path:",inline"` if embedding.
type TemplateVersionDeleteInput struct {
	TemplateVersionQueryInput `path:",inline"`
}

// TemplateVersionDeleteInputs holds the deletion input item of the TemplateVersion entities.
type TemplateVersionDeleteInputsItem struct {
	// ID of the TemplateVersion entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the TemplateVersion entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`
	// Version of the TemplateVersion entity, a part of the unique index.
	Version string `path:"-" query:"-" json:"version,omitempty"`
}

// TemplateVersionDeleteInputs holds the deletion input of the TemplateVersion entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type TemplateVersionDeleteInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to delete TemplateVersion entity CAN under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*TemplateVersionDeleteInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the TemplateVersion entities for deleting,
// after validating.
func (tvdi *TemplateVersionDeleteInputs) Model() []*TemplateVersion {
	if tvdi == nil || len(tvdi.Items) == 0 {
		return nil
	}

	_tvs := make([]*TemplateVersion, len(tvdi.Items))
	for i := range tvdi.Items {
		_tvs[i] = &TemplateVersion{
			ID: tvdi.Items[i].ID,
		}
	}
	return _tvs
}

// IDs returns the ID list of the TemplateVersion entities for deleting,
// after validating.
func (tvdi *TemplateVersionDeleteInputs) IDs() []object.ID {
	if tvdi == nil || len(tvdi.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(tvdi.Items))
	for i := range tvdi.Items {
		ids[i] = tvdi.Items[i].ID
	}
	return ids
}

// Validate checks the TemplateVersionDeleteInputs entity.
func (tvdi *TemplateVersionDeleteInputs) Validate() error {
	if tvdi == nil {
		return errors.New("nil receiver")
	}

	return tvdi.ValidateWith(tvdi.inputConfig.Context, tvdi.inputConfig.Client, nil)
}

// ValidateWith checks the TemplateVersionDeleteInputs entity with the given context and client set.
func (tvdi *TemplateVersionDeleteInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tvdi == nil {
		return errors.New("nil receiver")
	}

	if len(tvdi.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.TemplateVersions().Query()

	// Validate when deleting under the Project route.
	if tvdi.Project != nil {
		if err := tvdi.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				tvdi.Project = nil
				q.Where(
					templateversion.ProjectIDIsNil())
			}
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				templateversion.ProjectID(tvdi.Project.ID))
		}
	} else {
		q.Where(
			templateversion.ProjectIDIsNil())
	}

	ids := make([]object.ID, 0, len(tvdi.Items))
	ors := make([]predicate.TemplateVersion, 0, len(tvdi.Items))
	indexers := make(map[any][]int)

	for i := range tvdi.Items {
		if tvdi.Items[i] == nil {
			return errors.New("nil item")
		}

		if tvdi.Items[i].ID != "" {
			ids = append(ids, tvdi.Items[i].ID)
			ors = append(ors, templateversion.ID(tvdi.Items[i].ID))
			indexers[tvdi.Items[i].ID] = append(indexers[tvdi.Items[i].ID], i)
		} else if (tvdi.Items[i].Name != "") && (tvdi.Items[i].Version != "") {
			ors = append(ors, templateversion.And(
				templateversion.Name(tvdi.Items[i].Name),
				templateversion.Version(tvdi.Items[i].Version)))
			indexerKey := fmt.Sprint("/", tvdi.Items[i].Name, "/", tvdi.Items[i].Version)
			indexers[indexerKey] = append(indexers[indexerKey], i)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	p := templateversion.IDIn(ids...)
	if len(ids) != cap(ids) {
		p = templateversion.Or(ors...)
	}

	es, err := q.
		Where(p).
		Select(
			templateversion.FieldID,
			templateversion.FieldName,
			templateversion.FieldVersion,
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
			indexerKey := fmt.Sprint("/", es[i].Name, "/", es[i].Version)
			indexer = indexers[indexerKey]
		}
		for _, j := range indexer {
			tvdi.Items[j].ID = es[i].ID
			tvdi.Items[j].Name = es[i].Name
			tvdi.Items[j].Version = es[i].Version
		}
	}

	return nil
}

// TemplateVersionQueryInput holds the query input of the TemplateVersion entity,
// please tags with `path:",inline"` if embedding.
type TemplateVersionQueryInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to query TemplateVersion entity CAN under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"project,omitempty"`

	// Refer holds the route path reference of the TemplateVersion entity.
	Refer *object.Refer `path:"templateversion,default=" query:"-" json:"-"`
	// ID of the TemplateVersion entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the TemplateVersion entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`
	// Version of the TemplateVersion entity, a part of the unique index.
	Version string `path:"-" query:"-" json:"version,omitempty"`
}

// Model returns the TemplateVersion entity for querying,
// after validating.
func (tvqi *TemplateVersionQueryInput) Model() *TemplateVersion {
	if tvqi == nil {
		return nil
	}

	return &TemplateVersion{
		ID:      tvqi.ID,
		Name:    tvqi.Name,
		Version: tvqi.Version,
	}
}

// Validate checks the TemplateVersionQueryInput entity.
func (tvqi *TemplateVersionQueryInput) Validate() error {
	if tvqi == nil {
		return errors.New("nil receiver")
	}

	return tvqi.ValidateWith(tvqi.inputConfig.Context, tvqi.inputConfig.Client, nil)
}

// ValidateWith checks the TemplateVersionQueryInput entity with the given context and client set.
func (tvqi *TemplateVersionQueryInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tvqi == nil {
		return errors.New("nil receiver")
	}

	if tvqi.Refer != nil && *tvqi.Refer == "" {
		return fmt.Errorf("model: %s : %w", templateversion.Label, ErrBlankResourceRefer)
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.TemplateVersions().Query()

	// Validate when querying under the Project route.
	if tvqi.Project != nil {
		if err := tvqi.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				tvqi.Project = nil
				q.Where(
					templateversion.ProjectIDIsNil())
			}
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				templateversion.ProjectID(tvqi.Project.ID))
		}
	} else {
		q.Where(
			templateversion.ProjectIDIsNil())
	}

	if tvqi.Refer != nil {
		if tvqi.Refer.IsID() {
			q.Where(
				templateversion.ID(tvqi.Refer.ID()))
		} else if refers := tvqi.Refer.Split(2); len(refers) == 2 {
			q.Where(
				templateversion.Name(refers[0].String()),
				templateversion.Version(refers[1].String()))
		} else {
			return errors.New("invalid identify refer of templateversion")
		}
	} else if tvqi.ID != "" {
		q.Where(
			templateversion.ID(tvqi.ID))
	} else if (tvqi.Name != "") && (tvqi.Version != "") {
		q.Where(
			templateversion.Name(tvqi.Name),
			templateversion.Version(tvqi.Version))
	} else {
		return errors.New("invalid identify of templateversion")
	}

	q.Select(
		templateversion.FieldID,
		templateversion.FieldName,
		templateversion.FieldVersion,
	)

	var e *TemplateVersion
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
			e = cv.(*TemplateVersion)
		}
	}

	tvqi.ID = e.ID
	tvqi.Name = e.Name
	tvqi.Version = e.Version
	return nil
}

// TemplateVersionQueryInputs holds the query input of the TemplateVersion entities,
// please tags with `path:",inline" query:",inline"` if embedding.
type TemplateVersionQueryInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to query TemplateVersion entity CAN under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
}

// Validate checks the TemplateVersionQueryInputs entity.
func (tvqi *TemplateVersionQueryInputs) Validate() error {
	if tvqi == nil {
		return errors.New("nil receiver")
	}

	return tvqi.ValidateWith(tvqi.inputConfig.Context, tvqi.inputConfig.Client, nil)
}

// ValidateWith checks the TemplateVersionQueryInputs entity with the given context and client set.
func (tvqi *TemplateVersionQueryInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tvqi == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when querying under the Project route.
	if tvqi.Project != nil {
		if err := tvqi.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				tvqi.Project = nil
			}
		}
	}

	return nil
}

// TemplateVersionUpdateInput holds the modification input of the TemplateVersion entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type TemplateVersionUpdateInput struct {
	TemplateVersionQueryInput `path:",inline" query:"-" json:"-"`

	// Generated schema and data of the template.
	Schema types.Schema `path:"-" query:"-" json:"schema,omitempty"`
	// ui schema of the template.
	UiSchema types.UISchema `path:"-" query:"-" json:"uiSchema,omitempty"`
}

// Model returns the TemplateVersion entity for modifying,
// after validating.
func (tvui *TemplateVersionUpdateInput) Model() *TemplateVersion {
	if tvui == nil {
		return nil
	}

	_tv := &TemplateVersion{
		ID:       tvui.ID,
		Name:     tvui.Name,
		Version:  tvui.Version,
		Schema:   tvui.Schema,
		UiSchema: tvui.UiSchema,
	}

	return _tv
}

// Validate checks the TemplateVersionUpdateInput entity.
func (tvui *TemplateVersionUpdateInput) Validate() error {
	if tvui == nil {
		return errors.New("nil receiver")
	}

	return tvui.ValidateWith(tvui.inputConfig.Context, tvui.inputConfig.Client, nil)
}

// ValidateWith checks the TemplateVersionUpdateInput entity with the given context and client set.
func (tvui *TemplateVersionUpdateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := tvui.TemplateVersionQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	return nil
}

// TemplateVersionUpdateInputs holds the modification input item of the TemplateVersion entities.
type TemplateVersionUpdateInputsItem struct {
	// ID of the TemplateVersion entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the TemplateVersion entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`
	// Version of the TemplateVersion entity, a part of the unique index.
	Version string `path:"-" query:"-" json:"version,omitempty"`

	// Generated schema and data of the template.
	Schema types.Schema `path:"-" query:"-" json:"schema"`
	// ui schema of the template.
	UiSchema types.UISchema `path:"-" query:"-" json:"uiSchema"`
}

// ValidateWith checks the TemplateVersionUpdateInputsItem entity with the given context and client set.
func (tvui *TemplateVersionUpdateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tvui == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	return nil
}

// TemplateVersionUpdateInputs holds the modification input of the TemplateVersion entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type TemplateVersionUpdateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to update TemplateVersion entity CAN under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*TemplateVersionUpdateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the TemplateVersion entities for modifying,
// after validating.
func (tvui *TemplateVersionUpdateInputs) Model() []*TemplateVersion {
	if tvui == nil || len(tvui.Items) == 0 {
		return nil
	}

	_tvs := make([]*TemplateVersion, len(tvui.Items))

	for i := range tvui.Items {
		_tv := &TemplateVersion{
			ID:       tvui.Items[i].ID,
			Name:     tvui.Items[i].Name,
			Version:  tvui.Items[i].Version,
			Schema:   tvui.Items[i].Schema,
			UiSchema: tvui.Items[i].UiSchema,
		}

		_tvs[i] = _tv
	}

	return _tvs
}

// IDs returns the ID list of the TemplateVersion entities for modifying,
// after validating.
func (tvui *TemplateVersionUpdateInputs) IDs() []object.ID {
	if tvui == nil || len(tvui.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(tvui.Items))
	for i := range tvui.Items {
		ids[i] = tvui.Items[i].ID
	}
	return ids
}

// Validate checks the TemplateVersionUpdateInputs entity.
func (tvui *TemplateVersionUpdateInputs) Validate() error {
	if tvui == nil {
		return errors.New("nil receiver")
	}

	return tvui.ValidateWith(tvui.inputConfig.Context, tvui.inputConfig.Client, nil)
}

// ValidateWith checks the TemplateVersionUpdateInputs entity with the given context and client set.
func (tvui *TemplateVersionUpdateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if tvui == nil {
		return errors.New("nil receiver")
	}

	if len(tvui.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.TemplateVersions().Query()

	// Validate when updating under the Project route.
	if tvui.Project != nil {
		if err := tvui.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				tvui.Project = nil
				q.Where(
					templateversion.ProjectIDIsNil())
			}
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				templateversion.ProjectID(tvui.Project.ID))
		}
	} else {
		q.Where(
			templateversion.ProjectIDIsNil())
	}

	ids := make([]object.ID, 0, len(tvui.Items))
	ors := make([]predicate.TemplateVersion, 0, len(tvui.Items))
	indexers := make(map[any][]int)

	for i := range tvui.Items {
		if tvui.Items[i] == nil {
			return errors.New("nil item")
		}

		if tvui.Items[i].ID != "" {
			ids = append(ids, tvui.Items[i].ID)
			ors = append(ors, templateversion.ID(tvui.Items[i].ID))
			indexers[tvui.Items[i].ID] = append(indexers[tvui.Items[i].ID], i)
		} else if (tvui.Items[i].Name != "") && (tvui.Items[i].Version != "") {
			ors = append(ors, templateversion.And(
				templateversion.Name(tvui.Items[i].Name),
				templateversion.Version(tvui.Items[i].Version)))
			indexerKey := fmt.Sprint("/", tvui.Items[i].Name, "/", tvui.Items[i].Version)
			indexers[indexerKey] = append(indexers[indexerKey], i)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	p := templateversion.IDIn(ids...)
	if len(ids) != cap(ids) {
		p = templateversion.Or(ors...)
	}

	es, err := q.
		Where(p).
		Select(
			templateversion.FieldID,
			templateversion.FieldName,
			templateversion.FieldVersion,
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
			indexerKey := fmt.Sprint("/", es[i].Name, "/", es[i].Version)
			indexer = indexers[indexerKey]
		}
		for _, j := range indexer {
			tvui.Items[j].ID = es[i].ID
			tvui.Items[j].Name = es[i].Name
			tvui.Items[j].Version = es[i].Version
		}
	}

	for i := range tvui.Items {
		if err := tvui.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// TemplateVersionOutput holds the output of the TemplateVersion entity.
type TemplateVersionOutput struct {
	ID         object.ID      `json:"id,omitempty"`
	CreateTime *time.Time     `json:"createTime,omitempty"`
	UpdateTime *time.Time     `json:"updateTime,omitempty"`
	Name       string         `json:"name,omitempty"`
	Version    string         `json:"version,omitempty"`
	Source     string         `json:"source,omitempty"`
	Schema     types.Schema   `json:"schema,omitempty"`
	UiSchema   types.UISchema `json:"uiSchema,omitempty"`

	Template *TemplateOutput `json:"template,omitempty"`
	Project  *ProjectOutput  `json:"project,omitempty"`
}

// View returns the output of TemplateVersion entity.
func (_tv *TemplateVersion) View() *TemplateVersionOutput {
	return ExposeTemplateVersion(_tv)
}

// View returns the output of TemplateVersion entities.
func (_tvs TemplateVersions) View() []*TemplateVersionOutput {
	return ExposeTemplateVersions(_tvs)
}

// ExposeTemplateVersion converts the TemplateVersion to TemplateVersionOutput.
func ExposeTemplateVersion(_tv *TemplateVersion) *TemplateVersionOutput {
	if _tv == nil {
		return nil
	}

	tvo := &TemplateVersionOutput{
		ID:         _tv.ID,
		CreateTime: _tv.CreateTime,
		UpdateTime: _tv.UpdateTime,
		Name:       _tv.Name,
		Version:    _tv.Version,
		Source:     _tv.Source,
		Schema:     _tv.Schema,
		UiSchema:   _tv.UiSchema,
	}

	if _tv.Edges.Template != nil {
		tvo.Template = ExposeTemplate(_tv.Edges.Template)
	} else if _tv.TemplateID != "" {
		tvo.Template = &TemplateOutput{
			ID: _tv.TemplateID,
		}
	}
	if _tv.Edges.Project != nil {
		tvo.Project = ExposeProject(_tv.Edges.Project)
	} else if _tv.ProjectID != "" {
		tvo.Project = &ProjectOutput{
			ID: _tv.ProjectID,
		}
	}
	return tvo
}

// ExposeTemplateVersions converts the TemplateVersion slice to TemplateVersionOutput pointer slice.
func ExposeTemplateVersions(_tvs []*TemplateVersion) []*TemplateVersionOutput {
	if len(_tvs) == 0 {
		return nil
	}

	tvos := make([]*TemplateVersionOutput, len(_tvs))
	for i := range _tvs {
		tvos[i] = ExposeTemplateVersion(_tvs[i])
	}
	return tvos
}
