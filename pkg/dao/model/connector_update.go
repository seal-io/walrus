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

	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// ConnectorUpdate is the builder for updating Connector entities.
type ConnectorUpdate struct {
	config
	hooks     []Hook
	mutation  *ConnectorMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ConnectorUpdate builder.
func (cu *ConnectorUpdate) Where(ps ...predicate.Connector) *ConnectorUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetStatus sets the "status" field.
func (cu *ConnectorUpdate) SetStatus(s string) *ConnectorUpdate {
	cu.mutation.SetStatus(s)
	return cu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (cu *ConnectorUpdate) SetNillableStatus(s *string) *ConnectorUpdate {
	if s != nil {
		cu.SetStatus(*s)
	}
	return cu
}

// ClearStatus clears the value of the "status" field.
func (cu *ConnectorUpdate) ClearStatus() *ConnectorUpdate {
	cu.mutation.ClearStatus()
	return cu
}

// SetStatusMessage sets the "statusMessage" field.
func (cu *ConnectorUpdate) SetStatusMessage(s string) *ConnectorUpdate {
	cu.mutation.SetStatusMessage(s)
	return cu
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (cu *ConnectorUpdate) SetNillableStatusMessage(s *string) *ConnectorUpdate {
	if s != nil {
		cu.SetStatusMessage(*s)
	}
	return cu
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (cu *ConnectorUpdate) ClearStatusMessage() *ConnectorUpdate {
	cu.mutation.ClearStatusMessage()
	return cu
}

// SetUpdateTime sets the "updateTime" field.
func (cu *ConnectorUpdate) SetUpdateTime(t time.Time) *ConnectorUpdate {
	cu.mutation.SetUpdateTime(t)
	return cu
}

// SetDriver sets the "driver" field.
func (cu *ConnectorUpdate) SetDriver(s string) *ConnectorUpdate {
	cu.mutation.SetDriver(s)
	return cu
}

// SetConfigVersion sets the "configVersion" field.
func (cu *ConnectorUpdate) SetConfigVersion(s string) *ConnectorUpdate {
	cu.mutation.SetConfigVersion(s)
	return cu
}

// SetConfigData sets the "configData" field.
func (cu *ConnectorUpdate) SetConfigData(m map[string]interface{}) *ConnectorUpdate {
	cu.mutation.SetConfigData(m)
	return cu
}

// ClearConfigData clears the value of the "configData" field.
func (cu *ConnectorUpdate) ClearConfigData() *ConnectorUpdate {
	cu.mutation.ClearConfigData()
	return cu
}

// Mutation returns the ConnectorMutation object of the builder.
func (cu *ConnectorUpdate) Mutation() *ConnectorMutation {
	return cu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *ConnectorUpdate) Save(ctx context.Context) (int, error) {
	if err := cu.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, ConnectorMutation](ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *ConnectorUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *ConnectorUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *ConnectorUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cu *ConnectorUpdate) defaults() error {
	if _, ok := cu.mutation.UpdateTime(); !ok {
		if connector.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized connector.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := connector.UpdateDefaultUpdateTime()
		cu.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cu *ConnectorUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ConnectorUpdate {
	cu.modifiers = append(cu.modifiers, modifiers...)
	return cu
}

func (cu *ConnectorUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   connector.Table,
			Columns: connector.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: connector.FieldID,
			},
		},
	}
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.Status(); ok {
		_spec.SetField(connector.FieldStatus, field.TypeString, value)
	}
	if cu.mutation.StatusCleared() {
		_spec.ClearField(connector.FieldStatus, field.TypeString)
	}
	if value, ok := cu.mutation.StatusMessage(); ok {
		_spec.SetField(connector.FieldStatusMessage, field.TypeString, value)
	}
	if cu.mutation.StatusMessageCleared() {
		_spec.ClearField(connector.FieldStatusMessage, field.TypeString)
	}
	if value, ok := cu.mutation.UpdateTime(); ok {
		_spec.SetField(connector.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cu.mutation.Driver(); ok {
		_spec.SetField(connector.FieldDriver, field.TypeString, value)
	}
	if value, ok := cu.mutation.ConfigVersion(); ok {
		_spec.SetField(connector.FieldConfigVersion, field.TypeString, value)
	}
	if value, ok := cu.mutation.ConfigData(); ok {
		_spec.SetField(connector.FieldConfigData, field.TypeJSON, value)
	}
	if cu.mutation.ConfigDataCleared() {
		_spec.ClearField(connector.FieldConfigData, field.TypeJSON)
	}
	_spec.Node.Schema = cu.schemaConfig.Connector
	ctx = internal.NewSchemaConfigContext(ctx, cu.schemaConfig)
	_spec.AddModifiers(cu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{connector.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// ConnectorUpdateOne is the builder for updating a single Connector entity.
type ConnectorUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ConnectorMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetStatus sets the "status" field.
func (cuo *ConnectorUpdateOne) SetStatus(s string) *ConnectorUpdateOne {
	cuo.mutation.SetStatus(s)
	return cuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (cuo *ConnectorUpdateOne) SetNillableStatus(s *string) *ConnectorUpdateOne {
	if s != nil {
		cuo.SetStatus(*s)
	}
	return cuo
}

// ClearStatus clears the value of the "status" field.
func (cuo *ConnectorUpdateOne) ClearStatus() *ConnectorUpdateOne {
	cuo.mutation.ClearStatus()
	return cuo
}

// SetStatusMessage sets the "statusMessage" field.
func (cuo *ConnectorUpdateOne) SetStatusMessage(s string) *ConnectorUpdateOne {
	cuo.mutation.SetStatusMessage(s)
	return cuo
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (cuo *ConnectorUpdateOne) SetNillableStatusMessage(s *string) *ConnectorUpdateOne {
	if s != nil {
		cuo.SetStatusMessage(*s)
	}
	return cuo
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (cuo *ConnectorUpdateOne) ClearStatusMessage() *ConnectorUpdateOne {
	cuo.mutation.ClearStatusMessage()
	return cuo
}

// SetUpdateTime sets the "updateTime" field.
func (cuo *ConnectorUpdateOne) SetUpdateTime(t time.Time) *ConnectorUpdateOne {
	cuo.mutation.SetUpdateTime(t)
	return cuo
}

// SetDriver sets the "driver" field.
func (cuo *ConnectorUpdateOne) SetDriver(s string) *ConnectorUpdateOne {
	cuo.mutation.SetDriver(s)
	return cuo
}

// SetConfigVersion sets the "configVersion" field.
func (cuo *ConnectorUpdateOne) SetConfigVersion(s string) *ConnectorUpdateOne {
	cuo.mutation.SetConfigVersion(s)
	return cuo
}

// SetConfigData sets the "configData" field.
func (cuo *ConnectorUpdateOne) SetConfigData(m map[string]interface{}) *ConnectorUpdateOne {
	cuo.mutation.SetConfigData(m)
	return cuo
}

// ClearConfigData clears the value of the "configData" field.
func (cuo *ConnectorUpdateOne) ClearConfigData() *ConnectorUpdateOne {
	cuo.mutation.ClearConfigData()
	return cuo
}

// Mutation returns the ConnectorMutation object of the builder.
func (cuo *ConnectorUpdateOne) Mutation() *ConnectorMutation {
	return cuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *ConnectorUpdateOne) Select(field string, fields ...string) *ConnectorUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Connector entity.
func (cuo *ConnectorUpdateOne) Save(ctx context.Context) (*Connector, error) {
	if err := cuo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Connector, ConnectorMutation](ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *ConnectorUpdateOne) SaveX(ctx context.Context) *Connector {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *ConnectorUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *ConnectorUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cuo *ConnectorUpdateOne) defaults() error {
	if _, ok := cuo.mutation.UpdateTime(); !ok {
		if connector.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized connector.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := connector.UpdateDefaultUpdateTime()
		cuo.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cuo *ConnectorUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ConnectorUpdateOne {
	cuo.modifiers = append(cuo.modifiers, modifiers...)
	return cuo
}

func (cuo *ConnectorUpdateOne) sqlSave(ctx context.Context) (_node *Connector, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   connector.Table,
			Columns: connector.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: connector.FieldID,
			},
		},
	}
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Connector.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, connector.FieldID)
		for _, f := range fields {
			if !connector.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != connector.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.Status(); ok {
		_spec.SetField(connector.FieldStatus, field.TypeString, value)
	}
	if cuo.mutation.StatusCleared() {
		_spec.ClearField(connector.FieldStatus, field.TypeString)
	}
	if value, ok := cuo.mutation.StatusMessage(); ok {
		_spec.SetField(connector.FieldStatusMessage, field.TypeString, value)
	}
	if cuo.mutation.StatusMessageCleared() {
		_spec.ClearField(connector.FieldStatusMessage, field.TypeString)
	}
	if value, ok := cuo.mutation.UpdateTime(); ok {
		_spec.SetField(connector.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cuo.mutation.Driver(); ok {
		_spec.SetField(connector.FieldDriver, field.TypeString, value)
	}
	if value, ok := cuo.mutation.ConfigVersion(); ok {
		_spec.SetField(connector.FieldConfigVersion, field.TypeString, value)
	}
	if value, ok := cuo.mutation.ConfigData(); ok {
		_spec.SetField(connector.FieldConfigData, field.TypeJSON, value)
	}
	if cuo.mutation.ConfigDataCleared() {
		_spec.ClearField(connector.FieldConfigData, field.TypeJSON)
	}
	_spec.Node.Schema = cuo.schemaConfig.Connector
	ctx = internal.NewSchemaConfigContext(ctx, cuo.schemaConfig)
	_spec.AddModifiers(cuo.modifiers...)
	_node = &Connector{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{connector.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}
