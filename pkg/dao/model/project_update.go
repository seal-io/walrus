// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ProjectUpdate is the builder for updating Project entities.
type ProjectUpdate struct {
	config
	hooks     []Hook
	mutation  *ProjectMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ProjectUpdate builder.
func (pu *ProjectUpdate) Where(ps ...predicate.Project) *ProjectUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetName sets the "name" field.
func (pu *ProjectUpdate) SetName(s string) *ProjectUpdate {
	pu.mutation.SetName(s)
	return pu
}

// SetDescription sets the "description" field.
func (pu *ProjectUpdate) SetDescription(s string) *ProjectUpdate {
	pu.mutation.SetDescription(s)
	return pu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (pu *ProjectUpdate) SetNillableDescription(s *string) *ProjectUpdate {
	if s != nil {
		pu.SetDescription(*s)
	}
	return pu
}

// ClearDescription clears the value of the "description" field.
func (pu *ProjectUpdate) ClearDescription() *ProjectUpdate {
	pu.mutation.ClearDescription()
	return pu
}

// SetLabels sets the "labels" field.
func (pu *ProjectUpdate) SetLabels(m map[string]string) *ProjectUpdate {
	pu.mutation.SetLabels(m)
	return pu
}

// SetUpdateTime sets the "updateTime" field.
func (pu *ProjectUpdate) SetUpdateTime(t time.Time) *ProjectUpdate {
	pu.mutation.SetUpdateTime(t)
	return pu
}

// AddApplicationIDs adds the "applications" edge to the Application entity by IDs.
func (pu *ProjectUpdate) AddApplicationIDs(ids ...oid.ID) *ProjectUpdate {
	pu.mutation.AddApplicationIDs(ids...)
	return pu
}

// AddApplications adds the "applications" edges to the Application entity.
func (pu *ProjectUpdate) AddApplications(a ...*Application) *ProjectUpdate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return pu.AddApplicationIDs(ids...)
}

// AddSecretIDs adds the "secrets" edge to the Secret entity by IDs.
func (pu *ProjectUpdate) AddSecretIDs(ids ...oid.ID) *ProjectUpdate {
	pu.mutation.AddSecretIDs(ids...)
	return pu
}

// AddSecrets adds the "secrets" edges to the Secret entity.
func (pu *ProjectUpdate) AddSecrets(s ...*Secret) *ProjectUpdate {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pu.AddSecretIDs(ids...)
}

// Mutation returns the ProjectMutation object of the builder.
func (pu *ProjectUpdate) Mutation() *ProjectMutation {
	return pu.mutation
}

// ClearApplications clears all "applications" edges to the Application entity.
func (pu *ProjectUpdate) ClearApplications() *ProjectUpdate {
	pu.mutation.ClearApplications()
	return pu
}

// RemoveApplicationIDs removes the "applications" edge to Application entities by IDs.
func (pu *ProjectUpdate) RemoveApplicationIDs(ids ...oid.ID) *ProjectUpdate {
	pu.mutation.RemoveApplicationIDs(ids...)
	return pu
}

// RemoveApplications removes "applications" edges to Application entities.
func (pu *ProjectUpdate) RemoveApplications(a ...*Application) *ProjectUpdate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return pu.RemoveApplicationIDs(ids...)
}

// ClearSecrets clears all "secrets" edges to the Secret entity.
func (pu *ProjectUpdate) ClearSecrets() *ProjectUpdate {
	pu.mutation.ClearSecrets()
	return pu
}

// RemoveSecretIDs removes the "secrets" edge to Secret entities by IDs.
func (pu *ProjectUpdate) RemoveSecretIDs(ids ...oid.ID) *ProjectUpdate {
	pu.mutation.RemoveSecretIDs(ids...)
	return pu
}

// RemoveSecrets removes "secrets" edges to Secret entities.
func (pu *ProjectUpdate) RemoveSecrets(s ...*Secret) *ProjectUpdate {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return pu.RemoveSecretIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *ProjectUpdate) Save(ctx context.Context) (int, error) {
	if err := pu.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, ProjectMutation](ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *ProjectUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *ProjectUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *ProjectUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pu *ProjectUpdate) defaults() error {
	if _, ok := pu.mutation.UpdateTime(); !ok {
		if project.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized project.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := project.UpdateDefaultUpdateTime()
		pu.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (pu *ProjectUpdate) check() error {
	if v, ok := pu.mutation.Name(); ok {
		if err := project.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Project.name": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (pu *ProjectUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ProjectUpdate {
	pu.modifiers = append(pu.modifiers, modifiers...)
	return pu
}

func (pu *ProjectUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(project.Table, project.Columns, sqlgraph.NewFieldSpec(project.FieldID, field.TypeString))
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.Name(); ok {
		_spec.SetField(project.FieldName, field.TypeString, value)
	}
	if value, ok := pu.mutation.Description(); ok {
		_spec.SetField(project.FieldDescription, field.TypeString, value)
	}
	if pu.mutation.DescriptionCleared() {
		_spec.ClearField(project.FieldDescription, field.TypeString)
	}
	if value, ok := pu.mutation.Labels(); ok {
		_spec.SetField(project.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := pu.mutation.UpdateTime(); ok {
		_spec.SetField(project.FieldUpdateTime, field.TypeTime, value)
	}
	if pu.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.ApplicationsTable,
			Columns: []string{project.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = pu.schemaConfig.Application
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedApplicationsIDs(); len(nodes) > 0 && !pu.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.ApplicationsTable,
			Columns: []string{project.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = pu.schemaConfig.Application
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.ApplicationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.ApplicationsTable,
			Columns: []string{project.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = pu.schemaConfig.Application
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.mutation.SecretsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.SecretsTable,
			Columns: []string{project.SecretsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: secret.FieldID,
				},
			},
		}
		edge.Schema = pu.schemaConfig.Secret
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedSecretsIDs(); len(nodes) > 0 && !pu.mutation.SecretsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.SecretsTable,
			Columns: []string{project.SecretsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: secret.FieldID,
				},
			},
		}
		edge.Schema = pu.schemaConfig.Secret
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.SecretsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.SecretsTable,
			Columns: []string{project.SecretsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: secret.FieldID,
				},
			},
		}
		edge.Schema = pu.schemaConfig.Secret
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = pu.schemaConfig.Project
	ctx = internal.NewSchemaConfigContext(ctx, pu.schemaConfig)
	_spec.AddModifiers(pu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{project.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// ProjectUpdateOne is the builder for updating a single Project entity.
type ProjectUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ProjectMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetName sets the "name" field.
func (puo *ProjectUpdateOne) SetName(s string) *ProjectUpdateOne {
	puo.mutation.SetName(s)
	return puo
}

// SetDescription sets the "description" field.
func (puo *ProjectUpdateOne) SetDescription(s string) *ProjectUpdateOne {
	puo.mutation.SetDescription(s)
	return puo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (puo *ProjectUpdateOne) SetNillableDescription(s *string) *ProjectUpdateOne {
	if s != nil {
		puo.SetDescription(*s)
	}
	return puo
}

// ClearDescription clears the value of the "description" field.
func (puo *ProjectUpdateOne) ClearDescription() *ProjectUpdateOne {
	puo.mutation.ClearDescription()
	return puo
}

// SetLabels sets the "labels" field.
func (puo *ProjectUpdateOne) SetLabels(m map[string]string) *ProjectUpdateOne {
	puo.mutation.SetLabels(m)
	return puo
}

// SetUpdateTime sets the "updateTime" field.
func (puo *ProjectUpdateOne) SetUpdateTime(t time.Time) *ProjectUpdateOne {
	puo.mutation.SetUpdateTime(t)
	return puo
}

// AddApplicationIDs adds the "applications" edge to the Application entity by IDs.
func (puo *ProjectUpdateOne) AddApplicationIDs(ids ...oid.ID) *ProjectUpdateOne {
	puo.mutation.AddApplicationIDs(ids...)
	return puo
}

// AddApplications adds the "applications" edges to the Application entity.
func (puo *ProjectUpdateOne) AddApplications(a ...*Application) *ProjectUpdateOne {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return puo.AddApplicationIDs(ids...)
}

// AddSecretIDs adds the "secrets" edge to the Secret entity by IDs.
func (puo *ProjectUpdateOne) AddSecretIDs(ids ...oid.ID) *ProjectUpdateOne {
	puo.mutation.AddSecretIDs(ids...)
	return puo
}

// AddSecrets adds the "secrets" edges to the Secret entity.
func (puo *ProjectUpdateOne) AddSecrets(s ...*Secret) *ProjectUpdateOne {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return puo.AddSecretIDs(ids...)
}

// Mutation returns the ProjectMutation object of the builder.
func (puo *ProjectUpdateOne) Mutation() *ProjectMutation {
	return puo.mutation
}

// ClearApplications clears all "applications" edges to the Application entity.
func (puo *ProjectUpdateOne) ClearApplications() *ProjectUpdateOne {
	puo.mutation.ClearApplications()
	return puo
}

// RemoveApplicationIDs removes the "applications" edge to Application entities by IDs.
func (puo *ProjectUpdateOne) RemoveApplicationIDs(ids ...oid.ID) *ProjectUpdateOne {
	puo.mutation.RemoveApplicationIDs(ids...)
	return puo
}

// RemoveApplications removes "applications" edges to Application entities.
func (puo *ProjectUpdateOne) RemoveApplications(a ...*Application) *ProjectUpdateOne {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return puo.RemoveApplicationIDs(ids...)
}

// ClearSecrets clears all "secrets" edges to the Secret entity.
func (puo *ProjectUpdateOne) ClearSecrets() *ProjectUpdateOne {
	puo.mutation.ClearSecrets()
	return puo
}

// RemoveSecretIDs removes the "secrets" edge to Secret entities by IDs.
func (puo *ProjectUpdateOne) RemoveSecretIDs(ids ...oid.ID) *ProjectUpdateOne {
	puo.mutation.RemoveSecretIDs(ids...)
	return puo
}

// RemoveSecrets removes "secrets" edges to Secret entities.
func (puo *ProjectUpdateOne) RemoveSecrets(s ...*Secret) *ProjectUpdateOne {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return puo.RemoveSecretIDs(ids...)
}

// Where appends a list predicates to the ProjectUpdate builder.
func (puo *ProjectUpdateOne) Where(ps ...predicate.Project) *ProjectUpdateOne {
	puo.mutation.Where(ps...)
	return puo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *ProjectUpdateOne) Select(field string, fields ...string) *ProjectUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Project entity.
func (puo *ProjectUpdateOne) Save(ctx context.Context) (*Project, error) {
	if err := puo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Project, ProjectMutation](ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *ProjectUpdateOne) SaveX(ctx context.Context) *Project {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *ProjectUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *ProjectUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (puo *ProjectUpdateOne) defaults() error {
	if _, ok := puo.mutation.UpdateTime(); !ok {
		if project.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized project.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := project.UpdateDefaultUpdateTime()
		puo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (puo *ProjectUpdateOne) check() error {
	if v, ok := puo.mutation.Name(); ok {
		if err := project.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Project.name": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (puo *ProjectUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ProjectUpdateOne {
	puo.modifiers = append(puo.modifiers, modifiers...)
	return puo
}

func (puo *ProjectUpdateOne) sqlSave(ctx context.Context) (_node *Project, err error) {
	if err := puo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(project.Table, project.Columns, sqlgraph.NewFieldSpec(project.FieldID, field.TypeString))
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Project.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, project.FieldID)
		for _, f := range fields {
			if !project.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != project.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.Name(); ok {
		_spec.SetField(project.FieldName, field.TypeString, value)
	}
	if value, ok := puo.mutation.Description(); ok {
		_spec.SetField(project.FieldDescription, field.TypeString, value)
	}
	if puo.mutation.DescriptionCleared() {
		_spec.ClearField(project.FieldDescription, field.TypeString)
	}
	if value, ok := puo.mutation.Labels(); ok {
		_spec.SetField(project.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := puo.mutation.UpdateTime(); ok {
		_spec.SetField(project.FieldUpdateTime, field.TypeTime, value)
	}
	if puo.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.ApplicationsTable,
			Columns: []string{project.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = puo.schemaConfig.Application
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedApplicationsIDs(); len(nodes) > 0 && !puo.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.ApplicationsTable,
			Columns: []string{project.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = puo.schemaConfig.Application
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.ApplicationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.ApplicationsTable,
			Columns: []string{project.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = puo.schemaConfig.Application
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.mutation.SecretsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.SecretsTable,
			Columns: []string{project.SecretsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: secret.FieldID,
				},
			},
		}
		edge.Schema = puo.schemaConfig.Secret
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedSecretsIDs(); len(nodes) > 0 && !puo.mutation.SecretsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.SecretsTable,
			Columns: []string{project.SecretsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: secret.FieldID,
				},
			},
		}
		edge.Schema = puo.schemaConfig.Secret
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.SecretsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   project.SecretsTable,
			Columns: []string{project.SecretsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: secret.FieldID,
				},
			},
		}
		edge.Schema = puo.schemaConfig.Secret
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = puo.schemaConfig.Project
	ctx = internal.NewSchemaConfigContext(ctx, puo.schemaConfig)
	_spec.AddModifiers(puo.modifiers...)
	_node = &Project{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{project.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}
