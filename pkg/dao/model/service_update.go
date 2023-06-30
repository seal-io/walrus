// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerelationship"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

// ServiceUpdate is the builder for updating Service entities.
type ServiceUpdate struct {
	config
	hooks     []Hook
	mutation  *ServiceMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ServiceUpdate builder.
func (su *ServiceUpdate) Where(ps ...predicate.Service) *ServiceUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetName sets the "name" field.
func (su *ServiceUpdate) SetName(s string) *ServiceUpdate {
	su.mutation.SetName(s)
	return su
}

// SetDescription sets the "description" field.
func (su *ServiceUpdate) SetDescription(s string) *ServiceUpdate {
	su.mutation.SetDescription(s)
	return su
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableDescription(s *string) *ServiceUpdate {
	if s != nil {
		su.SetDescription(*s)
	}
	return su
}

// ClearDescription clears the value of the "description" field.
func (su *ServiceUpdate) ClearDescription() *ServiceUpdate {
	su.mutation.ClearDescription()
	return su
}

// SetLabels sets the "labels" field.
func (su *ServiceUpdate) SetLabels(m map[string]string) *ServiceUpdate {
	su.mutation.SetLabels(m)
	return su
}

// ClearLabels clears the value of the "labels" field.
func (su *ServiceUpdate) ClearLabels() *ServiceUpdate {
	su.mutation.ClearLabels()
	return su
}

// SetAnnotations sets the "annotations" field.
func (su *ServiceUpdate) SetAnnotations(m map[string]string) *ServiceUpdate {
	su.mutation.SetAnnotations(m)
	return su
}

// ClearAnnotations clears the value of the "annotations" field.
func (su *ServiceUpdate) ClearAnnotations() *ServiceUpdate {
	su.mutation.ClearAnnotations()
	return su
}

// SetUpdateTime sets the "updateTime" field.
func (su *ServiceUpdate) SetUpdateTime(t time.Time) *ServiceUpdate {
	su.mutation.SetUpdateTime(t)
	return su
}

// SetTemplate sets the "template" field.
func (su *ServiceUpdate) SetTemplate(tvr types.TemplateVersionRef) *ServiceUpdate {
	su.mutation.SetTemplate(tvr)
	return su
}

// SetAttributes sets the "attributes" field.
func (su *ServiceUpdate) SetAttributes(pr property.Values) *ServiceUpdate {
	su.mutation.SetAttributes(pr)
	return su
}

// ClearAttributes clears the value of the "attributes" field.
func (su *ServiceUpdate) ClearAttributes() *ServiceUpdate {
	su.mutation.ClearAttributes()
	return su
}

// SetStatus sets the "status" field.
func (su *ServiceUpdate) SetStatus(s status.Status) *ServiceUpdate {
	su.mutation.SetStatus(s)
	return su
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (su *ServiceUpdate) SetNillableStatus(s *status.Status) *ServiceUpdate {
	if s != nil {
		su.SetStatus(*s)
	}
	return su
}

// ClearStatus clears the value of the "status" field.
func (su *ServiceUpdate) ClearStatus() *ServiceUpdate {
	su.mutation.ClearStatus()
	return su
}

// AddRevisionIDs adds the "revisions" edge to the ServiceRevision entity by IDs.
func (su *ServiceUpdate) AddRevisionIDs(ids ...oid.ID) *ServiceUpdate {
	su.mutation.AddRevisionIDs(ids...)
	return su
}

// AddRevisions adds the "revisions" edges to the ServiceRevision entity.
func (su *ServiceUpdate) AddRevisions(s ...*ServiceRevision) *ServiceUpdate {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddRevisionIDs(ids...)
}

// AddResourceIDs adds the "resources" edge to the ServiceResource entity by IDs.
func (su *ServiceUpdate) AddResourceIDs(ids ...oid.ID) *ServiceUpdate {
	su.mutation.AddResourceIDs(ids...)
	return su
}

// AddResources adds the "resources" edges to the ServiceResource entity.
func (su *ServiceUpdate) AddResources(s ...*ServiceResource) *ServiceUpdate {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddResourceIDs(ids...)
}

// AddDependencyIDs adds the "dependencies" edge to the ServiceRelationship entity by IDs.
func (su *ServiceUpdate) AddDependencyIDs(ids ...oid.ID) *ServiceUpdate {
	su.mutation.AddDependencyIDs(ids...)
	return su
}

// AddDependencies adds the "dependencies" edges to the ServiceRelationship entity.
func (su *ServiceUpdate) AddDependencies(s ...*ServiceRelationship) *ServiceUpdate {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddDependencyIDs(ids...)
}

// Mutation returns the ServiceMutation object of the builder.
func (su *ServiceUpdate) Mutation() *ServiceMutation {
	return su.mutation
}

// ClearRevisions clears all "revisions" edges to the ServiceRevision entity.
func (su *ServiceUpdate) ClearRevisions() *ServiceUpdate {
	su.mutation.ClearRevisions()
	return su
}

// RemoveRevisionIDs removes the "revisions" edge to ServiceRevision entities by IDs.
func (su *ServiceUpdate) RemoveRevisionIDs(ids ...oid.ID) *ServiceUpdate {
	su.mutation.RemoveRevisionIDs(ids...)
	return su
}

// RemoveRevisions removes "revisions" edges to ServiceRevision entities.
func (su *ServiceUpdate) RemoveRevisions(s ...*ServiceRevision) *ServiceUpdate {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveRevisionIDs(ids...)
}

// ClearResources clears all "resources" edges to the ServiceResource entity.
func (su *ServiceUpdate) ClearResources() *ServiceUpdate {
	su.mutation.ClearResources()
	return su
}

// RemoveResourceIDs removes the "resources" edge to ServiceResource entities by IDs.
func (su *ServiceUpdate) RemoveResourceIDs(ids ...oid.ID) *ServiceUpdate {
	su.mutation.RemoveResourceIDs(ids...)
	return su
}

// RemoveResources removes "resources" edges to ServiceResource entities.
func (su *ServiceUpdate) RemoveResources(s ...*ServiceResource) *ServiceUpdate {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveResourceIDs(ids...)
}

// ClearDependencies clears all "dependencies" edges to the ServiceRelationship entity.
func (su *ServiceUpdate) ClearDependencies() *ServiceUpdate {
	su.mutation.ClearDependencies()
	return su
}

// RemoveDependencyIDs removes the "dependencies" edge to ServiceRelationship entities by IDs.
func (su *ServiceUpdate) RemoveDependencyIDs(ids ...oid.ID) *ServiceUpdate {
	su.mutation.RemoveDependencyIDs(ids...)
	return su
}

// RemoveDependencies removes "dependencies" edges to ServiceRelationship entities.
func (su *ServiceUpdate) RemoveDependencies(s ...*ServiceRelationship) *ServiceUpdate {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveDependencyIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *ServiceUpdate) Save(ctx context.Context) (int, error) {
	if err := su.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, ServiceMutation](ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *ServiceUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *ServiceUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *ServiceUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (su *ServiceUpdate) defaults() error {
	if _, ok := su.mutation.UpdateTime(); !ok {
		if service.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized service.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := service.UpdateDefaultUpdateTime()
		su.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (su *ServiceUpdate) check() error {
	if v, ok := su.mutation.Name(); ok {
		if err := service.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Service.name": %w`, err)}
		}
	}
	if _, ok := su.mutation.ProjectID(); su.mutation.ProjectCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "Service.project"`)
	}
	if _, ok := su.mutation.EnvironmentID(); su.mutation.EnvironmentCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "Service.environment"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (su *ServiceUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ServiceUpdate {
	su.modifiers = append(su.modifiers, modifiers...)
	return su
}

func (su *ServiceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(service.Table, service.Columns, sqlgraph.NewFieldSpec(service.FieldID, field.TypeString))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.Name(); ok {
		_spec.SetField(service.FieldName, field.TypeString, value)
	}
	if value, ok := su.mutation.Description(); ok {
		_spec.SetField(service.FieldDescription, field.TypeString, value)
	}
	if su.mutation.DescriptionCleared() {
		_spec.ClearField(service.FieldDescription, field.TypeString)
	}
	if value, ok := su.mutation.Labels(); ok {
		_spec.SetField(service.FieldLabels, field.TypeJSON, value)
	}
	if su.mutation.LabelsCleared() {
		_spec.ClearField(service.FieldLabels, field.TypeJSON)
	}
	if value, ok := su.mutation.Annotations(); ok {
		_spec.SetField(service.FieldAnnotations, field.TypeJSON, value)
	}
	if su.mutation.AnnotationsCleared() {
		_spec.ClearField(service.FieldAnnotations, field.TypeJSON)
	}
	if value, ok := su.mutation.UpdateTime(); ok {
		_spec.SetField(service.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := su.mutation.Template(); ok {
		_spec.SetField(service.FieldTemplate, field.TypeJSON, value)
	}
	if value, ok := su.mutation.Attributes(); ok {
		_spec.SetField(service.FieldAttributes, field.TypeOther, value)
	}
	if su.mutation.AttributesCleared() {
		_spec.ClearField(service.FieldAttributes, field.TypeOther)
	}
	if value, ok := su.mutation.Status(); ok {
		_spec.SetField(service.FieldStatus, field.TypeJSON, value)
	}
	if su.mutation.StatusCleared() {
		_spec.ClearField(service.FieldStatus, field.TypeJSON)
	}
	if su.mutation.RevisionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.RevisionsTable,
			Columns: []string{service.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerevision.FieldID, field.TypeString),
			},
		}
		edge.Schema = su.schemaConfig.ServiceRevision
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedRevisionsIDs(); len(nodes) > 0 && !su.mutation.RevisionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.RevisionsTable,
			Columns: []string{service.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerevision.FieldID, field.TypeString),
			},
		}
		edge.Schema = su.schemaConfig.ServiceRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RevisionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.RevisionsTable,
			Columns: []string{service.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerevision.FieldID, field.TypeString),
			},
		}
		edge.Schema = su.schemaConfig.ServiceRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.ResourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.ResourcesTable,
			Columns: []string{service.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = su.schemaConfig.ServiceResource
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedResourcesIDs(); len(nodes) > 0 && !su.mutation.ResourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.ResourcesTable,
			Columns: []string{service.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = su.schemaConfig.ServiceResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.ResourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.ResourcesTable,
			Columns: []string{service.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = su.schemaConfig.ServiceResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.DependenciesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   service.DependenciesTable,
			Columns: []string{service.DependenciesColumn},
			Bidi:    true,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = su.schemaConfig.ServiceRelationship
		createE := &ServiceRelationshipCreate{config: su.config, mutation: newServiceRelationshipMutation(su.config, OpCreate)}
		_ = createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedDependenciesIDs(); len(nodes) > 0 && !su.mutation.DependenciesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   service.DependenciesTable,
			Columns: []string{service.DependenciesColumn},
			Bidi:    true,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = su.schemaConfig.ServiceRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &ServiceRelationshipCreate{config: su.config, mutation: newServiceRelationshipMutation(su.config, OpCreate)}
		_ = createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.DependenciesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   service.DependenciesTable,
			Columns: []string{service.DependenciesColumn},
			Bidi:    true,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = su.schemaConfig.ServiceRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &ServiceRelationshipCreate{config: su.config, mutation: newServiceRelationshipMutation(su.config, OpCreate)}
		_ = createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = su.schemaConfig.Service
	ctx = internal.NewSchemaConfigContext(ctx, su.schemaConfig)
	_spec.AddModifiers(su.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{service.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// ServiceUpdateOne is the builder for updating a single Service entity.
type ServiceUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ServiceMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetName sets the "name" field.
func (suo *ServiceUpdateOne) SetName(s string) *ServiceUpdateOne {
	suo.mutation.SetName(s)
	return suo
}

// SetDescription sets the "description" field.
func (suo *ServiceUpdateOne) SetDescription(s string) *ServiceUpdateOne {
	suo.mutation.SetDescription(s)
	return suo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableDescription(s *string) *ServiceUpdateOne {
	if s != nil {
		suo.SetDescription(*s)
	}
	return suo
}

// ClearDescription clears the value of the "description" field.
func (suo *ServiceUpdateOne) ClearDescription() *ServiceUpdateOne {
	suo.mutation.ClearDescription()
	return suo
}

// SetLabels sets the "labels" field.
func (suo *ServiceUpdateOne) SetLabels(m map[string]string) *ServiceUpdateOne {
	suo.mutation.SetLabels(m)
	return suo
}

// ClearLabels clears the value of the "labels" field.
func (suo *ServiceUpdateOne) ClearLabels() *ServiceUpdateOne {
	suo.mutation.ClearLabels()
	return suo
}

// SetAnnotations sets the "annotations" field.
func (suo *ServiceUpdateOne) SetAnnotations(m map[string]string) *ServiceUpdateOne {
	suo.mutation.SetAnnotations(m)
	return suo
}

// ClearAnnotations clears the value of the "annotations" field.
func (suo *ServiceUpdateOne) ClearAnnotations() *ServiceUpdateOne {
	suo.mutation.ClearAnnotations()
	return suo
}

// SetUpdateTime sets the "updateTime" field.
func (suo *ServiceUpdateOne) SetUpdateTime(t time.Time) *ServiceUpdateOne {
	suo.mutation.SetUpdateTime(t)
	return suo
}

// SetTemplate sets the "template" field.
func (suo *ServiceUpdateOne) SetTemplate(tvr types.TemplateVersionRef) *ServiceUpdateOne {
	suo.mutation.SetTemplate(tvr)
	return suo
}

// SetAttributes sets the "attributes" field.
func (suo *ServiceUpdateOne) SetAttributes(pr property.Values) *ServiceUpdateOne {
	suo.mutation.SetAttributes(pr)
	return suo
}

// ClearAttributes clears the value of the "attributes" field.
func (suo *ServiceUpdateOne) ClearAttributes() *ServiceUpdateOne {
	suo.mutation.ClearAttributes()
	return suo
}

// SetStatus sets the "status" field.
func (suo *ServiceUpdateOne) SetStatus(s status.Status) *ServiceUpdateOne {
	suo.mutation.SetStatus(s)
	return suo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (suo *ServiceUpdateOne) SetNillableStatus(s *status.Status) *ServiceUpdateOne {
	if s != nil {
		suo.SetStatus(*s)
	}
	return suo
}

// ClearStatus clears the value of the "status" field.
func (suo *ServiceUpdateOne) ClearStatus() *ServiceUpdateOne {
	suo.mutation.ClearStatus()
	return suo
}

// AddRevisionIDs adds the "revisions" edge to the ServiceRevision entity by IDs.
func (suo *ServiceUpdateOne) AddRevisionIDs(ids ...oid.ID) *ServiceUpdateOne {
	suo.mutation.AddRevisionIDs(ids...)
	return suo
}

// AddRevisions adds the "revisions" edges to the ServiceRevision entity.
func (suo *ServiceUpdateOne) AddRevisions(s ...*ServiceRevision) *ServiceUpdateOne {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddRevisionIDs(ids...)
}

// AddResourceIDs adds the "resources" edge to the ServiceResource entity by IDs.
func (suo *ServiceUpdateOne) AddResourceIDs(ids ...oid.ID) *ServiceUpdateOne {
	suo.mutation.AddResourceIDs(ids...)
	return suo
}

// AddResources adds the "resources" edges to the ServiceResource entity.
func (suo *ServiceUpdateOne) AddResources(s ...*ServiceResource) *ServiceUpdateOne {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddResourceIDs(ids...)
}

// AddDependencyIDs adds the "dependencies" edge to the ServiceRelationship entity by IDs.
func (suo *ServiceUpdateOne) AddDependencyIDs(ids ...oid.ID) *ServiceUpdateOne {
	suo.mutation.AddDependencyIDs(ids...)
	return suo
}

// AddDependencies adds the "dependencies" edges to the ServiceRelationship entity.
func (suo *ServiceUpdateOne) AddDependencies(s ...*ServiceRelationship) *ServiceUpdateOne {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddDependencyIDs(ids...)
}

// Mutation returns the ServiceMutation object of the builder.
func (suo *ServiceUpdateOne) Mutation() *ServiceMutation {
	return suo.mutation
}

// ClearRevisions clears all "revisions" edges to the ServiceRevision entity.
func (suo *ServiceUpdateOne) ClearRevisions() *ServiceUpdateOne {
	suo.mutation.ClearRevisions()
	return suo
}

// RemoveRevisionIDs removes the "revisions" edge to ServiceRevision entities by IDs.
func (suo *ServiceUpdateOne) RemoveRevisionIDs(ids ...oid.ID) *ServiceUpdateOne {
	suo.mutation.RemoveRevisionIDs(ids...)
	return suo
}

// RemoveRevisions removes "revisions" edges to ServiceRevision entities.
func (suo *ServiceUpdateOne) RemoveRevisions(s ...*ServiceRevision) *ServiceUpdateOne {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveRevisionIDs(ids...)
}

// ClearResources clears all "resources" edges to the ServiceResource entity.
func (suo *ServiceUpdateOne) ClearResources() *ServiceUpdateOne {
	suo.mutation.ClearResources()
	return suo
}

// RemoveResourceIDs removes the "resources" edge to ServiceResource entities by IDs.
func (suo *ServiceUpdateOne) RemoveResourceIDs(ids ...oid.ID) *ServiceUpdateOne {
	suo.mutation.RemoveResourceIDs(ids...)
	return suo
}

// RemoveResources removes "resources" edges to ServiceResource entities.
func (suo *ServiceUpdateOne) RemoveResources(s ...*ServiceResource) *ServiceUpdateOne {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveResourceIDs(ids...)
}

// ClearDependencies clears all "dependencies" edges to the ServiceRelationship entity.
func (suo *ServiceUpdateOne) ClearDependencies() *ServiceUpdateOne {
	suo.mutation.ClearDependencies()
	return suo
}

// RemoveDependencyIDs removes the "dependencies" edge to ServiceRelationship entities by IDs.
func (suo *ServiceUpdateOne) RemoveDependencyIDs(ids ...oid.ID) *ServiceUpdateOne {
	suo.mutation.RemoveDependencyIDs(ids...)
	return suo
}

// RemoveDependencies removes "dependencies" edges to ServiceRelationship entities.
func (suo *ServiceUpdateOne) RemoveDependencies(s ...*ServiceRelationship) *ServiceUpdateOne {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveDependencyIDs(ids...)
}

// Where appends a list predicates to the ServiceUpdate builder.
func (suo *ServiceUpdateOne) Where(ps ...predicate.Service) *ServiceUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *ServiceUpdateOne) Select(field string, fields ...string) *ServiceUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Service entity.
func (suo *ServiceUpdateOne) Save(ctx context.Context) (*Service, error) {
	if err := suo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Service, ServiceMutation](ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *ServiceUpdateOne) SaveX(ctx context.Context) *Service {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *ServiceUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *ServiceUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (suo *ServiceUpdateOne) defaults() error {
	if _, ok := suo.mutation.UpdateTime(); !ok {
		if service.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized service.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := service.UpdateDefaultUpdateTime()
		suo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (suo *ServiceUpdateOne) check() error {
	if v, ok := suo.mutation.Name(); ok {
		if err := service.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Service.name": %w`, err)}
		}
	}
	if _, ok := suo.mutation.ProjectID(); suo.mutation.ProjectCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "Service.project"`)
	}
	if _, ok := suo.mutation.EnvironmentID(); suo.mutation.EnvironmentCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "Service.environment"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (suo *ServiceUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ServiceUpdateOne {
	suo.modifiers = append(suo.modifiers, modifiers...)
	return suo
}

func (suo *ServiceUpdateOne) sqlSave(ctx context.Context) (_node *Service, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(service.Table, service.Columns, sqlgraph.NewFieldSpec(service.FieldID, field.TypeString))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Service.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, service.FieldID)
		for _, f := range fields {
			if !service.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != service.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.Name(); ok {
		_spec.SetField(service.FieldName, field.TypeString, value)
	}
	if value, ok := suo.mutation.Description(); ok {
		_spec.SetField(service.FieldDescription, field.TypeString, value)
	}
	if suo.mutation.DescriptionCleared() {
		_spec.ClearField(service.FieldDescription, field.TypeString)
	}
	if value, ok := suo.mutation.Labels(); ok {
		_spec.SetField(service.FieldLabels, field.TypeJSON, value)
	}
	if suo.mutation.LabelsCleared() {
		_spec.ClearField(service.FieldLabels, field.TypeJSON)
	}
	if value, ok := suo.mutation.Annotations(); ok {
		_spec.SetField(service.FieldAnnotations, field.TypeJSON, value)
	}
	if suo.mutation.AnnotationsCleared() {
		_spec.ClearField(service.FieldAnnotations, field.TypeJSON)
	}
	if value, ok := suo.mutation.UpdateTime(); ok {
		_spec.SetField(service.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := suo.mutation.Template(); ok {
		_spec.SetField(service.FieldTemplate, field.TypeJSON, value)
	}
	if value, ok := suo.mutation.Attributes(); ok {
		_spec.SetField(service.FieldAttributes, field.TypeOther, value)
	}
	if suo.mutation.AttributesCleared() {
		_spec.ClearField(service.FieldAttributes, field.TypeOther)
	}
	if value, ok := suo.mutation.Status(); ok {
		_spec.SetField(service.FieldStatus, field.TypeJSON, value)
	}
	if suo.mutation.StatusCleared() {
		_spec.ClearField(service.FieldStatus, field.TypeJSON)
	}
	if suo.mutation.RevisionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.RevisionsTable,
			Columns: []string{service.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerevision.FieldID, field.TypeString),
			},
		}
		edge.Schema = suo.schemaConfig.ServiceRevision
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedRevisionsIDs(); len(nodes) > 0 && !suo.mutation.RevisionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.RevisionsTable,
			Columns: []string{service.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerevision.FieldID, field.TypeString),
			},
		}
		edge.Schema = suo.schemaConfig.ServiceRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RevisionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.RevisionsTable,
			Columns: []string{service.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerevision.FieldID, field.TypeString),
			},
		}
		edge.Schema = suo.schemaConfig.ServiceRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.ResourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.ResourcesTable,
			Columns: []string{service.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = suo.schemaConfig.ServiceResource
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedResourcesIDs(); len(nodes) > 0 && !suo.mutation.ResourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.ResourcesTable,
			Columns: []string{service.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = suo.schemaConfig.ServiceResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.ResourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.ResourcesTable,
			Columns: []string{service.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = suo.schemaConfig.ServiceResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.DependenciesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   service.DependenciesTable,
			Columns: []string{service.DependenciesColumn},
			Bidi:    true,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = suo.schemaConfig.ServiceRelationship
		createE := &ServiceRelationshipCreate{config: suo.config, mutation: newServiceRelationshipMutation(suo.config, OpCreate)}
		_ = createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedDependenciesIDs(); len(nodes) > 0 && !suo.mutation.DependenciesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   service.DependenciesTable,
			Columns: []string{service.DependenciesColumn},
			Bidi:    true,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = suo.schemaConfig.ServiceRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &ServiceRelationshipCreate{config: suo.config, mutation: newServiceRelationshipMutation(suo.config, OpCreate)}
		_ = createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.DependenciesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   service.DependenciesTable,
			Columns: []string{service.DependenciesColumn},
			Bidi:    true,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(servicerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = suo.schemaConfig.ServiceRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &ServiceRelationshipCreate{config: suo.config, mutation: newServiceRelationshipMutation(suo.config, OpCreate)}
		_ = createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = suo.schemaConfig.Service
	ctx = internal.NewSchemaConfigContext(ctx, suo.schemaConfig)
	_spec.AddModifiers(suo.modifiers...)
	_node = &Service{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{service.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
