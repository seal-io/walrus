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
	"github.com/seal-io/walrus/pkg/dao/model/service"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
)

// ServiceCreateInput holds the creation input of the Service entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type ServiceCreateInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to create Service entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to create Service entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`

	// Name holds the value of the "name" field.
	Name string `path:"-" query:"-" json:"name"`
	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`

	// Template specifies full inserting the new TemplateVersion entity of the Service entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template"`
}

// Model returns the Service entity for creating,
// after validating.
func (sci *ServiceCreateInput) Model() *Service {
	if sci == nil {
		return nil
	}

	_s := &Service{
		Name:        sci.Name,
		Description: sci.Description,
		Labels:      sci.Labels,
		Attributes:  sci.Attributes,
	}

	if sci.Project != nil {
		_s.ProjectID = sci.Project.ID
	}
	if sci.Environment != nil {
		_s.EnvironmentID = sci.Environment.ID
	}

	if sci.Template != nil {
		_s.TemplateID = sci.Template.ID
	}
	return _s
}

// Validate checks the ServiceCreateInput entity.
func (sci *ServiceCreateInput) Validate() error {
	if sci == nil {
		return errors.New("nil receiver")
	}

	return sci.ValidateWith(sci.inputConfig.Context, sci.inputConfig.Client, nil)
}

// ValidateWith checks the ServiceCreateInput entity with the given context and client set.
func (sci *ServiceCreateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if sci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when creating under the Project route.
	if sci.Project != nil {
		if err := sci.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}
	// Validate when creating under the Environment route.
	if sci.Environment != nil {
		if err := sci.Environment.ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	if sci.Template != nil {
		if err := sci.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				sci.Template = nil
			}
		}
	}

	return nil
}

// ServiceCreateInputs holds the creation input item of the Service entities.
type ServiceCreateInputsItem struct {
	// Name holds the value of the "name" field.
	Name string `path:"-" query:"-" json:"name"`
	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`

	// Template specifies full inserting the new TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template"`
}

// ValidateWith checks the ServiceCreateInputsItem entity with the given context and client set.
func (sci *ServiceCreateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if sci == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	if sci.Template != nil {
		if err := sci.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				sci.Template = nil
			}
		}
	}

	return nil
}

// ServiceCreateInputs holds the creation input of the Service entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type ServiceCreateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to create Service entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to create Service entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*ServiceCreateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the Service entities for creating,
// after validating.
func (sci *ServiceCreateInputs) Model() []*Service {
	if sci == nil || len(sci.Items) == 0 {
		return nil
	}

	_ss := make([]*Service, len(sci.Items))

	for i := range sci.Items {
		_s := &Service{
			Name:        sci.Items[i].Name,
			Description: sci.Items[i].Description,
			Labels:      sci.Items[i].Labels,
			Attributes:  sci.Items[i].Attributes,
		}

		if sci.Project != nil {
			_s.ProjectID = sci.Project.ID
		}
		if sci.Environment != nil {
			_s.EnvironmentID = sci.Environment.ID
		}

		if sci.Items[i].Template != nil {
			_s.TemplateID = sci.Items[i].Template.ID
		}

		_ss[i] = _s
	}

	return _ss
}

// Validate checks the ServiceCreateInputs entity .
func (sci *ServiceCreateInputs) Validate() error {
	if sci == nil {
		return errors.New("nil receiver")
	}

	return sci.ValidateWith(sci.inputConfig.Context, sci.inputConfig.Client, nil)
}

// ValidateWith checks the ServiceCreateInputs entity with the given context and client set.
func (sci *ServiceCreateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if sci == nil {
		return errors.New("nil receiver")
	}

	if len(sci.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when creating under the Project route.
	if sci.Project != nil {
		if err := sci.Project.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				sci.Project = nil
			}
		}
	}
	// Validate when creating under the Environment route.
	if sci.Environment != nil {
		if err := sci.Environment.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				sci.Environment = nil
			}
		}
	}

	for i := range sci.Items {
		if sci.Items[i] == nil {
			continue
		}

		if err := sci.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// ServiceDeleteInput holds the deletion input of the Service entity,
// please tags with `path:",inline"` if embedding.
type ServiceDeleteInput struct {
	ServiceQueryInput `path:",inline"`
}

// ServiceDeleteInputs holds the deletion input item of the Service entities.
type ServiceDeleteInputsItem struct {
	// ID of the Service entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the Service entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`
}

// ServiceDeleteInputs holds the deletion input of the Service entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type ServiceDeleteInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to delete Service entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to delete Service entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*ServiceDeleteInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the Service entities for deleting,
// after validating.
func (sdi *ServiceDeleteInputs) Model() []*Service {
	if sdi == nil || len(sdi.Items) == 0 {
		return nil
	}

	_ss := make([]*Service, len(sdi.Items))
	for i := range sdi.Items {
		_ss[i] = &Service{
			ID: sdi.Items[i].ID,
		}
	}
	return _ss
}

// IDs returns the ID list of the Service entities for deleting,
// after validating.
func (sdi *ServiceDeleteInputs) IDs() []object.ID {
	if sdi == nil || len(sdi.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(sdi.Items))
	for i := range sdi.Items {
		ids[i] = sdi.Items[i].ID
	}
	return ids
}

// Validate checks the ServiceDeleteInputs entity.
func (sdi *ServiceDeleteInputs) Validate() error {
	if sdi == nil {
		return errors.New("nil receiver")
	}

	return sdi.ValidateWith(sdi.inputConfig.Context, sdi.inputConfig.Client, nil)
}

// ValidateWith checks the ServiceDeleteInputs entity with the given context and client set.
func (sdi *ServiceDeleteInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if sdi == nil {
		return errors.New("nil receiver")
	}

	if len(sdi.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.Services().Query()

	// Validate when deleting under the Project route.
	if sdi.Project != nil {
		if err := sdi.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				service.ProjectID(sdi.Project.ID))
		}
	}

	// Validate when deleting under the Environment route.
	if sdi.Environment != nil {
		if err := sdi.Environment.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			q.Where(
				service.EnvironmentID(sdi.Environment.ID))
		}
	}

	ids := make([]object.ID, 0, len(sdi.Items))
	ors := make([]predicate.Service, 0, len(sdi.Items))
	indexers := make(map[any][]int)

	for i := range sdi.Items {
		if sdi.Items[i] == nil {
			return errors.New("nil item")
		}

		if sdi.Items[i].ID != "" {
			ids = append(ids, sdi.Items[i].ID)
			ors = append(ors, service.ID(sdi.Items[i].ID))
			indexers[sdi.Items[i].ID] = append(indexers[sdi.Items[i].ID], i)
		} else if sdi.Items[i].Name != "" {
			ors = append(ors, service.And(
				service.Name(sdi.Items[i].Name)))
			indexerKey := fmt.Sprint("/", sdi.Items[i].Name)
			indexers[indexerKey] = append(indexers[indexerKey], i)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	p := service.IDIn(ids...)
	if len(ids) != cap(ids) {
		p = service.Or(ors...)
	}

	es, err := q.
		Where(p).
		Select(
			service.FieldID,
			service.FieldName,
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
			indexerKey := fmt.Sprint("/", sdi.Items[i].Name)
			indexer = indexers[indexerKey]
		}
		for _, j := range indexer {
			sdi.Items[j].ID = es[i].ID
			sdi.Items[j].Name = es[i].Name
		}
	}

	return nil
}

// ServiceQueryInput holds the query input of the Service entity,
// please tags with `path:",inline"` if embedding.
type ServiceQueryInput struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to query Service entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"project"`
	// Environment indicates to query Service entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"environment"`

	// Refer holds the route path reference of the Service entity.
	Refer *object.Refer `path:"service,default=" query:"-" json:"-"`
	// ID of the Service entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the Service entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`
}

// Model returns the Service entity for querying,
// after validating.
func (sqi *ServiceQueryInput) Model() *Service {
	if sqi == nil {
		return nil
	}

	return &Service{
		ID:   sqi.ID,
		Name: sqi.Name,
	}
}

// Validate checks the ServiceQueryInput entity.
func (sqi *ServiceQueryInput) Validate() error {
	if sqi == nil {
		return errors.New("nil receiver")
	}

	return sqi.ValidateWith(sqi.inputConfig.Context, sqi.inputConfig.Client, nil)
}

// ValidateWith checks the ServiceQueryInput entity with the given context and client set.
func (sqi *ServiceQueryInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if sqi == nil {
		return errors.New("nil receiver")
	}

	if sqi.Refer != nil && *sqi.Refer == "" {
		return fmt.Errorf("model: %s : %w", service.Label, ErrBlankResourceRefer)
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.Services().Query()

	// Validate when querying under the Project route.
	if sqi.Project != nil {
		if err := sqi.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				service.ProjectID(sqi.Project.ID))
		}
	}

	// Validate when querying under the Environment route.
	if sqi.Environment != nil {
		if err := sqi.Environment.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			q.Where(
				service.EnvironmentID(sqi.Environment.ID))
		}
	}

	if sqi.Refer != nil {
		if sqi.Refer.IsID() {
			q.Where(
				service.ID(sqi.Refer.ID()))
		} else if refers := sqi.Refer.Split(1); len(refers) == 1 {
			q.Where(
				service.Name(refers[0].String()))
		} else {
			return errors.New("invalid identify refer of service")
		}
	} else if sqi.ID != "" {
		q.Where(
			service.ID(sqi.ID))
	} else if sqi.Name != "" {
		q.Where(
			service.Name(sqi.Name))
	} else {
		return errors.New("invalid identify of service")
	}

	q.Select(
		service.FieldID,
		service.FieldName,
	)

	var e *Service
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
			e = cv.(*Service)
		}
	}

	sqi.ID = e.ID
	sqi.Name = e.Name
	return nil
}

// ServiceQueryInputs holds the query input of the Service entities,
// please tags with `path:",inline" query:",inline"` if embedding.
type ServiceQueryInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to query Service entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to query Service entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`
}

// Validate checks the ServiceQueryInputs entity.
func (sqi *ServiceQueryInputs) Validate() error {
	if sqi == nil {
		return errors.New("nil receiver")
	}

	return sqi.ValidateWith(sqi.inputConfig.Context, sqi.inputConfig.Client, nil)
}

// ValidateWith checks the ServiceQueryInputs entity with the given context and client set.
func (sqi *ServiceQueryInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if sqi == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	// Validate when querying under the Project route.
	if sqi.Project != nil {
		if err := sqi.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	// Validate when querying under the Environment route.
	if sqi.Environment != nil {
		if err := sqi.Environment.ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// ServiceUpdateInput holds the modification input of the Service entity,
// please tags with `path:",inline" json:",inline"` if embedding.
type ServiceUpdateInput struct {
	ServiceQueryInput `path:",inline" query:"-" json:"-"`

	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`

	// Template indicates replacing the stale TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template"`
}

// Model returns the Service entity for modifying,
// after validating.
func (sui *ServiceUpdateInput) Model() *Service {
	if sui == nil {
		return nil
	}

	_s := &Service{
		ID:          sui.ID,
		Name:        sui.Name,
		Description: sui.Description,
		Labels:      sui.Labels,
		Attributes:  sui.Attributes,
	}

	if sui.Template != nil {
		_s.TemplateID = sui.Template.ID
	}
	return _s
}

// Validate checks the ServiceUpdateInput entity.
func (sui *ServiceUpdateInput) Validate() error {
	if sui == nil {
		return errors.New("nil receiver")
	}

	return sui.ValidateWith(sui.inputConfig.Context, sui.inputConfig.Client, nil)
}

// ValidateWith checks the ServiceUpdateInput entity with the given context and client set.
func (sui *ServiceUpdateInput) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if cache == nil {
		cache = map[string]any{}
	}

	if err := sui.ServiceQueryInput.ValidateWith(ctx, cs, cache); err != nil {
		return err
	}

	if sui.Template != nil {
		if err := sui.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				sui.Template = nil
			}
		}
	}

	return nil
}

// ServiceUpdateInputs holds the modification input item of the Service entities.
type ServiceUpdateInputsItem struct {
	// ID of the Service entity, tries to retrieve the entity with the following unique index parts if no ID provided.
	ID object.ID `path:"-" query:"-" json:"id,omitempty"`
	// Name of the Service entity, a part of the unique index.
	Name string `path:"-" query:"-" json:"name,omitempty"`

	// Description holds the value of the "description" field.
	Description string `path:"-" query:"-" json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `path:"-" query:"-" json:"labels,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `path:"-" query:"-" json:"attributes,omitempty"`

	// Template indicates replacing the stale TemplateVersion entity.
	Template *TemplateVersionQueryInput `uri:"-" query:"-" json:"template"`
}

// ValidateWith checks the ServiceUpdateInputsItem entity with the given context and client set.
func (sui *ServiceUpdateInputsItem) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if sui == nil {
		return errors.New("nil receiver")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	if sui.Template != nil {
		if err := sui.Template.ValidateWith(ctx, cs, cache); err != nil {
			if !IsBlankResourceReferError(err) {
				return err
			} else {
				sui.Template = nil
			}
		}
	}

	return nil
}

// ServiceUpdateInputs holds the modification input of the Service entities,
// please tags with `path:",inline" json:",inline"` if embedding.
type ServiceUpdateInputs struct {
	inputConfig `path:"-" query:"-" json:"-"`

	// Project indicates to update Service entity MUST under the Project route.
	Project *ProjectQueryInput `path:",inline" query:"-" json:"-"`
	// Environment indicates to update Service entity MUST under the Environment route.
	Environment *EnvironmentQueryInput `path:",inline" query:"-" json:"-"`

	// Items holds the entities to create, which MUST not be empty.
	Items []*ServiceUpdateInputsItem `path:"-" query:"-" json:"items"`
}

// Model returns the Service entities for modifying,
// after validating.
func (sui *ServiceUpdateInputs) Model() []*Service {
	if sui == nil || len(sui.Items) == 0 {
		return nil
	}

	_ss := make([]*Service, len(sui.Items))

	for i := range sui.Items {
		_s := &Service{
			ID:          sui.Items[i].ID,
			Name:        sui.Items[i].Name,
			Description: sui.Items[i].Description,
			Labels:      sui.Items[i].Labels,
			Attributes:  sui.Items[i].Attributes,
		}

		if sui.Items[i].Template != nil {
			_s.TemplateID = sui.Items[i].Template.ID
		}

		_ss[i] = _s
	}

	return _ss
}

// IDs returns the ID list of the Service entities for modifying,
// after validating.
func (sui *ServiceUpdateInputs) IDs() []object.ID {
	if sui == nil || len(sui.Items) == 0 {
		return nil
	}

	ids := make([]object.ID, len(sui.Items))
	for i := range sui.Items {
		ids[i] = sui.Items[i].ID
	}
	return ids
}

// Validate checks the ServiceUpdateInputs entity.
func (sui *ServiceUpdateInputs) Validate() error {
	if sui == nil {
		return errors.New("nil receiver")
	}

	return sui.ValidateWith(sui.inputConfig.Context, sui.inputConfig.Client, nil)
}

// ValidateWith checks the ServiceUpdateInputs entity with the given context and client set.
func (sui *ServiceUpdateInputs) ValidateWith(ctx context.Context, cs ClientSet, cache map[string]any) error {
	if sui == nil {
		return errors.New("nil receiver")
	}

	if len(sui.Items) == 0 {
		return errors.New("empty items")
	}

	if cache == nil {
		cache = map[string]any{}
	}

	q := cs.Services().Query()

	// Validate when updating under the Project route.
	if sui.Project != nil {
		if err := sui.Project.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			ctx = valueContext(ctx, intercept.WithProjectInterceptor)
			q.Where(
				service.ProjectID(sui.Project.ID))
		}
	}

	// Validate when updating under the Environment route.
	if sui.Environment != nil {
		if err := sui.Environment.ValidateWith(ctx, cs, cache); err != nil {
			return err
		} else {
			q.Where(
				service.EnvironmentID(sui.Environment.ID))
		}
	}

	ids := make([]object.ID, 0, len(sui.Items))
	ors := make([]predicate.Service, 0, len(sui.Items))
	indexers := make(map[any][]int)

	for i := range sui.Items {
		if sui.Items[i] == nil {
			return errors.New("nil item")
		}

		if sui.Items[i].ID != "" {
			ids = append(ids, sui.Items[i].ID)
			ors = append(ors, service.ID(sui.Items[i].ID))
			indexers[sui.Items[i].ID] = append(indexers[sui.Items[i].ID], i)
		} else if sui.Items[i].Name != "" {
			ors = append(ors, service.And(
				service.Name(sui.Items[i].Name)))
			indexerKey := fmt.Sprint("/", sui.Items[i].Name)
			indexers[indexerKey] = append(indexers[indexerKey], i)
		} else {
			return errors.New("found item hasn't identify")
		}
	}

	p := service.IDIn(ids...)
	if len(ids) != cap(ids) {
		p = service.Or(ors...)
	}

	es, err := q.
		Where(p).
		Select(
			service.FieldID,
			service.FieldName,
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
			indexerKey := fmt.Sprint("/", sui.Items[i].Name)
			indexer = indexers[indexerKey]
		}
		for _, j := range indexer {
			sui.Items[j].ID = es[i].ID
			sui.Items[j].Name = es[i].Name
		}
	}

	for i := range sui.Items {
		if err := sui.Items[i].ValidateWith(ctx, cs, cache); err != nil {
			return err
		}
	}

	return nil
}

// ServiceOutput holds the output of the Service entity.
type ServiceOutput struct {
	ID          object.ID         `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	CreateTime  *time.Time        `json:"createTime,omitempty"`
	UpdateTime  *time.Time        `json:"updateTime,omitempty"`
	Status      status.Status     `json:"status,omitempty"`
	Attributes  property.Values   `json:"attributes,omitempty"`

	Project     *ProjectOutput         `json:"project,omitempty"`
	Environment *EnvironmentOutput     `json:"environment,omitempty"`
	Template    *TemplateVersionOutput `json:"template,omitempty"`
}

// View returns the output of Service entity.
func (_s *Service) View() *ServiceOutput {
	return ExposeService(_s)
}

// View returns the output of Service entities.
func (_ss Services) View() []*ServiceOutput {
	return ExposeServices(_ss)
}

// ExposeService converts the Service to ServiceOutput.
func ExposeService(_s *Service) *ServiceOutput {
	if _s == nil {
		return nil
	}

	so := &ServiceOutput{
		ID:          _s.ID,
		Name:        _s.Name,
		Description: _s.Description,
		Labels:      _s.Labels,
		CreateTime:  _s.CreateTime,
		UpdateTime:  _s.UpdateTime,
		Status:      _s.Status,
		Attributes:  _s.Attributes,
	}

	if _s.Edges.Project != nil {
		so.Project = ExposeProject(_s.Edges.Project)
	} else if _s.ProjectID != "" {
		so.Project = &ProjectOutput{
			ID: _s.ProjectID,
		}
	}
	if _s.Edges.Environment != nil {
		so.Environment = ExposeEnvironment(_s.Edges.Environment)
	} else if _s.EnvironmentID != "" {
		so.Environment = &EnvironmentOutput{
			ID: _s.EnvironmentID,
		}
	}
	if _s.Edges.Template != nil {
		so.Template = ExposeTemplateVersion(_s.Edges.Template)
	} else if _s.TemplateID != "" {
		so.Template = &TemplateVersionOutput{
			ID: _s.TemplateID,
		}
	}
	return so
}

// ExposeServices converts the Service slice to ServiceOutput pointer slice.
func ExposeServices(_ss []*Service) []*ServiceOutput {
	if len(_ss) == 0 {
		return nil
	}

	sos := make([]*ServiceOutput, len(_ss))
	for i := range _ss {
		sos[i] = ExposeService(_ss[i])
	}
	return sos
}
