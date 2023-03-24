// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// ClusterCostUpdate is the builder for updating ClusterCost entities.
type ClusterCostUpdate struct {
	config
	hooks     []Hook
	mutation  *ClusterCostMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ClusterCostUpdate builder.
func (ccu *ClusterCostUpdate) Where(ps ...predicate.ClusterCost) *ClusterCostUpdate {
	ccu.mutation.Where(ps...)
	return ccu
}

// SetTotalCost sets the "totalCost" field.
func (ccu *ClusterCostUpdate) SetTotalCost(f float64) *ClusterCostUpdate {
	ccu.mutation.ResetTotalCost()
	ccu.mutation.SetTotalCost(f)
	return ccu
}

// SetNillableTotalCost sets the "totalCost" field if the given value is not nil.
func (ccu *ClusterCostUpdate) SetNillableTotalCost(f *float64) *ClusterCostUpdate {
	if f != nil {
		ccu.SetTotalCost(*f)
	}
	return ccu
}

// AddTotalCost adds f to the "totalCost" field.
func (ccu *ClusterCostUpdate) AddTotalCost(f float64) *ClusterCostUpdate {
	ccu.mutation.AddTotalCost(f)
	return ccu
}

// SetCurrency sets the "currency" field.
func (ccu *ClusterCostUpdate) SetCurrency(i int) *ClusterCostUpdate {
	ccu.mutation.ResetCurrency()
	ccu.mutation.SetCurrency(i)
	return ccu
}

// SetNillableCurrency sets the "currency" field if the given value is not nil.
func (ccu *ClusterCostUpdate) SetNillableCurrency(i *int) *ClusterCostUpdate {
	if i != nil {
		ccu.SetCurrency(*i)
	}
	return ccu
}

// AddCurrency adds i to the "currency" field.
func (ccu *ClusterCostUpdate) AddCurrency(i int) *ClusterCostUpdate {
	ccu.mutation.AddCurrency(i)
	return ccu
}

// ClearCurrency clears the value of the "currency" field.
func (ccu *ClusterCostUpdate) ClearCurrency() *ClusterCostUpdate {
	ccu.mutation.ClearCurrency()
	return ccu
}

// SetAllocationCost sets the "allocationCost" field.
func (ccu *ClusterCostUpdate) SetAllocationCost(f float64) *ClusterCostUpdate {
	ccu.mutation.ResetAllocationCost()
	ccu.mutation.SetAllocationCost(f)
	return ccu
}

// SetNillableAllocationCost sets the "allocationCost" field if the given value is not nil.
func (ccu *ClusterCostUpdate) SetNillableAllocationCost(f *float64) *ClusterCostUpdate {
	if f != nil {
		ccu.SetAllocationCost(*f)
	}
	return ccu
}

// AddAllocationCost adds f to the "allocationCost" field.
func (ccu *ClusterCostUpdate) AddAllocationCost(f float64) *ClusterCostUpdate {
	ccu.mutation.AddAllocationCost(f)
	return ccu
}

// SetIdleCost sets the "idleCost" field.
func (ccu *ClusterCostUpdate) SetIdleCost(f float64) *ClusterCostUpdate {
	ccu.mutation.ResetIdleCost()
	ccu.mutation.SetIdleCost(f)
	return ccu
}

// SetNillableIdleCost sets the "idleCost" field if the given value is not nil.
func (ccu *ClusterCostUpdate) SetNillableIdleCost(f *float64) *ClusterCostUpdate {
	if f != nil {
		ccu.SetIdleCost(*f)
	}
	return ccu
}

// AddIdleCost adds f to the "idleCost" field.
func (ccu *ClusterCostUpdate) AddIdleCost(f float64) *ClusterCostUpdate {
	ccu.mutation.AddIdleCost(f)
	return ccu
}

// SetManagementCost sets the "managementCost" field.
func (ccu *ClusterCostUpdate) SetManagementCost(f float64) *ClusterCostUpdate {
	ccu.mutation.ResetManagementCost()
	ccu.mutation.SetManagementCost(f)
	return ccu
}

// SetNillableManagementCost sets the "managementCost" field if the given value is not nil.
func (ccu *ClusterCostUpdate) SetNillableManagementCost(f *float64) *ClusterCostUpdate {
	if f != nil {
		ccu.SetManagementCost(*f)
	}
	return ccu
}

// AddManagementCost adds f to the "managementCost" field.
func (ccu *ClusterCostUpdate) AddManagementCost(f float64) *ClusterCostUpdate {
	ccu.mutation.AddManagementCost(f)
	return ccu
}

// Mutation returns the ClusterCostMutation object of the builder.
func (ccu *ClusterCostUpdate) Mutation() *ClusterCostMutation {
	return ccu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ccu *ClusterCostUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, ClusterCostMutation](ctx, ccu.sqlSave, ccu.mutation, ccu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ccu *ClusterCostUpdate) SaveX(ctx context.Context) int {
	affected, err := ccu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ccu *ClusterCostUpdate) Exec(ctx context.Context) error {
	_, err := ccu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccu *ClusterCostUpdate) ExecX(ctx context.Context) {
	if err := ccu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ccu *ClusterCostUpdate) check() error {
	if v, ok := ccu.mutation.TotalCost(); ok {
		if err := clustercost.TotalCostValidator(v); err != nil {
			return &ValidationError{Name: "totalCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.totalCost": %w`, err)}
		}
	}
	if v, ok := ccu.mutation.AllocationCost(); ok {
		if err := clustercost.AllocationCostValidator(v); err != nil {
			return &ValidationError{Name: "allocationCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.allocationCost": %w`, err)}
		}
	}
	if v, ok := ccu.mutation.IdleCost(); ok {
		if err := clustercost.IdleCostValidator(v); err != nil {
			return &ValidationError{Name: "idleCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.idleCost": %w`, err)}
		}
	}
	if v, ok := ccu.mutation.ManagementCost(); ok {
		if err := clustercost.ManagementCostValidator(v); err != nil {
			return &ValidationError{Name: "managementCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.managementCost": %w`, err)}
		}
	}
	if _, ok := ccu.mutation.ConnectorID(); ccu.mutation.ConnectorCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ClusterCost.connector"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ccu *ClusterCostUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ClusterCostUpdate {
	ccu.modifiers = append(ccu.modifiers, modifiers...)
	return ccu
}

func (ccu *ClusterCostUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ccu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(clustercost.Table, clustercost.Columns, sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt))
	if ps := ccu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ccu.mutation.TotalCost(); ok {
		_spec.SetField(clustercost.FieldTotalCost, field.TypeFloat64, value)
	}
	if value, ok := ccu.mutation.AddedTotalCost(); ok {
		_spec.AddField(clustercost.FieldTotalCost, field.TypeFloat64, value)
	}
	if value, ok := ccu.mutation.Currency(); ok {
		_spec.SetField(clustercost.FieldCurrency, field.TypeInt, value)
	}
	if value, ok := ccu.mutation.AddedCurrency(); ok {
		_spec.AddField(clustercost.FieldCurrency, field.TypeInt, value)
	}
	if ccu.mutation.CurrencyCleared() {
		_spec.ClearField(clustercost.FieldCurrency, field.TypeInt)
	}
	if value, ok := ccu.mutation.AllocationCost(); ok {
		_spec.SetField(clustercost.FieldAllocationCost, field.TypeFloat64, value)
	}
	if value, ok := ccu.mutation.AddedAllocationCost(); ok {
		_spec.AddField(clustercost.FieldAllocationCost, field.TypeFloat64, value)
	}
	if value, ok := ccu.mutation.IdleCost(); ok {
		_spec.SetField(clustercost.FieldIdleCost, field.TypeFloat64, value)
	}
	if value, ok := ccu.mutation.AddedIdleCost(); ok {
		_spec.AddField(clustercost.FieldIdleCost, field.TypeFloat64, value)
	}
	if value, ok := ccu.mutation.ManagementCost(); ok {
		_spec.SetField(clustercost.FieldManagementCost, field.TypeFloat64, value)
	}
	if value, ok := ccu.mutation.AddedManagementCost(); ok {
		_spec.AddField(clustercost.FieldManagementCost, field.TypeFloat64, value)
	}
	_spec.Node.Schema = ccu.schemaConfig.ClusterCost
	ctx = internal.NewSchemaConfigContext(ctx, ccu.schemaConfig)
	_spec.AddModifiers(ccu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ccu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{clustercost.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ccu.mutation.done = true
	return n, nil
}

// ClusterCostUpdateOne is the builder for updating a single ClusterCost entity.
type ClusterCostUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ClusterCostMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetTotalCost sets the "totalCost" field.
func (ccuo *ClusterCostUpdateOne) SetTotalCost(f float64) *ClusterCostUpdateOne {
	ccuo.mutation.ResetTotalCost()
	ccuo.mutation.SetTotalCost(f)
	return ccuo
}

// SetNillableTotalCost sets the "totalCost" field if the given value is not nil.
func (ccuo *ClusterCostUpdateOne) SetNillableTotalCost(f *float64) *ClusterCostUpdateOne {
	if f != nil {
		ccuo.SetTotalCost(*f)
	}
	return ccuo
}

// AddTotalCost adds f to the "totalCost" field.
func (ccuo *ClusterCostUpdateOne) AddTotalCost(f float64) *ClusterCostUpdateOne {
	ccuo.mutation.AddTotalCost(f)
	return ccuo
}

// SetCurrency sets the "currency" field.
func (ccuo *ClusterCostUpdateOne) SetCurrency(i int) *ClusterCostUpdateOne {
	ccuo.mutation.ResetCurrency()
	ccuo.mutation.SetCurrency(i)
	return ccuo
}

// SetNillableCurrency sets the "currency" field if the given value is not nil.
func (ccuo *ClusterCostUpdateOne) SetNillableCurrency(i *int) *ClusterCostUpdateOne {
	if i != nil {
		ccuo.SetCurrency(*i)
	}
	return ccuo
}

// AddCurrency adds i to the "currency" field.
func (ccuo *ClusterCostUpdateOne) AddCurrency(i int) *ClusterCostUpdateOne {
	ccuo.mutation.AddCurrency(i)
	return ccuo
}

// ClearCurrency clears the value of the "currency" field.
func (ccuo *ClusterCostUpdateOne) ClearCurrency() *ClusterCostUpdateOne {
	ccuo.mutation.ClearCurrency()
	return ccuo
}

// SetAllocationCost sets the "allocationCost" field.
func (ccuo *ClusterCostUpdateOne) SetAllocationCost(f float64) *ClusterCostUpdateOne {
	ccuo.mutation.ResetAllocationCost()
	ccuo.mutation.SetAllocationCost(f)
	return ccuo
}

// SetNillableAllocationCost sets the "allocationCost" field if the given value is not nil.
func (ccuo *ClusterCostUpdateOne) SetNillableAllocationCost(f *float64) *ClusterCostUpdateOne {
	if f != nil {
		ccuo.SetAllocationCost(*f)
	}
	return ccuo
}

// AddAllocationCost adds f to the "allocationCost" field.
func (ccuo *ClusterCostUpdateOne) AddAllocationCost(f float64) *ClusterCostUpdateOne {
	ccuo.mutation.AddAllocationCost(f)
	return ccuo
}

// SetIdleCost sets the "idleCost" field.
func (ccuo *ClusterCostUpdateOne) SetIdleCost(f float64) *ClusterCostUpdateOne {
	ccuo.mutation.ResetIdleCost()
	ccuo.mutation.SetIdleCost(f)
	return ccuo
}

// SetNillableIdleCost sets the "idleCost" field if the given value is not nil.
func (ccuo *ClusterCostUpdateOne) SetNillableIdleCost(f *float64) *ClusterCostUpdateOne {
	if f != nil {
		ccuo.SetIdleCost(*f)
	}
	return ccuo
}

// AddIdleCost adds f to the "idleCost" field.
func (ccuo *ClusterCostUpdateOne) AddIdleCost(f float64) *ClusterCostUpdateOne {
	ccuo.mutation.AddIdleCost(f)
	return ccuo
}

// SetManagementCost sets the "managementCost" field.
func (ccuo *ClusterCostUpdateOne) SetManagementCost(f float64) *ClusterCostUpdateOne {
	ccuo.mutation.ResetManagementCost()
	ccuo.mutation.SetManagementCost(f)
	return ccuo
}

// SetNillableManagementCost sets the "managementCost" field if the given value is not nil.
func (ccuo *ClusterCostUpdateOne) SetNillableManagementCost(f *float64) *ClusterCostUpdateOne {
	if f != nil {
		ccuo.SetManagementCost(*f)
	}
	return ccuo
}

// AddManagementCost adds f to the "managementCost" field.
func (ccuo *ClusterCostUpdateOne) AddManagementCost(f float64) *ClusterCostUpdateOne {
	ccuo.mutation.AddManagementCost(f)
	return ccuo
}

// Mutation returns the ClusterCostMutation object of the builder.
func (ccuo *ClusterCostUpdateOne) Mutation() *ClusterCostMutation {
	return ccuo.mutation
}

// Where appends a list predicates to the ClusterCostUpdate builder.
func (ccuo *ClusterCostUpdateOne) Where(ps ...predicate.ClusterCost) *ClusterCostUpdateOne {
	ccuo.mutation.Where(ps...)
	return ccuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ccuo *ClusterCostUpdateOne) Select(field string, fields ...string) *ClusterCostUpdateOne {
	ccuo.fields = append([]string{field}, fields...)
	return ccuo
}

// Save executes the query and returns the updated ClusterCost entity.
func (ccuo *ClusterCostUpdateOne) Save(ctx context.Context) (*ClusterCost, error) {
	return withHooks[*ClusterCost, ClusterCostMutation](ctx, ccuo.sqlSave, ccuo.mutation, ccuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ccuo *ClusterCostUpdateOne) SaveX(ctx context.Context) *ClusterCost {
	node, err := ccuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ccuo *ClusterCostUpdateOne) Exec(ctx context.Context) error {
	_, err := ccuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccuo *ClusterCostUpdateOne) ExecX(ctx context.Context) {
	if err := ccuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ccuo *ClusterCostUpdateOne) check() error {
	if v, ok := ccuo.mutation.TotalCost(); ok {
		if err := clustercost.TotalCostValidator(v); err != nil {
			return &ValidationError{Name: "totalCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.totalCost": %w`, err)}
		}
	}
	if v, ok := ccuo.mutation.AllocationCost(); ok {
		if err := clustercost.AllocationCostValidator(v); err != nil {
			return &ValidationError{Name: "allocationCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.allocationCost": %w`, err)}
		}
	}
	if v, ok := ccuo.mutation.IdleCost(); ok {
		if err := clustercost.IdleCostValidator(v); err != nil {
			return &ValidationError{Name: "idleCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.idleCost": %w`, err)}
		}
	}
	if v, ok := ccuo.mutation.ManagementCost(); ok {
		if err := clustercost.ManagementCostValidator(v); err != nil {
			return &ValidationError{Name: "managementCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.managementCost": %w`, err)}
		}
	}
	if _, ok := ccuo.mutation.ConnectorID(); ccuo.mutation.ConnectorCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ClusterCost.connector"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ccuo *ClusterCostUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ClusterCostUpdateOne {
	ccuo.modifiers = append(ccuo.modifiers, modifiers...)
	return ccuo
}

func (ccuo *ClusterCostUpdateOne) sqlSave(ctx context.Context) (_node *ClusterCost, err error) {
	if err := ccuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(clustercost.Table, clustercost.Columns, sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt))
	id, ok := ccuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "ClusterCost.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ccuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, clustercost.FieldID)
		for _, f := range fields {
			if !clustercost.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != clustercost.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ccuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ccuo.mutation.TotalCost(); ok {
		_spec.SetField(clustercost.FieldTotalCost, field.TypeFloat64, value)
	}
	if value, ok := ccuo.mutation.AddedTotalCost(); ok {
		_spec.AddField(clustercost.FieldTotalCost, field.TypeFloat64, value)
	}
	if value, ok := ccuo.mutation.Currency(); ok {
		_spec.SetField(clustercost.FieldCurrency, field.TypeInt, value)
	}
	if value, ok := ccuo.mutation.AddedCurrency(); ok {
		_spec.AddField(clustercost.FieldCurrency, field.TypeInt, value)
	}
	if ccuo.mutation.CurrencyCleared() {
		_spec.ClearField(clustercost.FieldCurrency, field.TypeInt)
	}
	if value, ok := ccuo.mutation.AllocationCost(); ok {
		_spec.SetField(clustercost.FieldAllocationCost, field.TypeFloat64, value)
	}
	if value, ok := ccuo.mutation.AddedAllocationCost(); ok {
		_spec.AddField(clustercost.FieldAllocationCost, field.TypeFloat64, value)
	}
	if value, ok := ccuo.mutation.IdleCost(); ok {
		_spec.SetField(clustercost.FieldIdleCost, field.TypeFloat64, value)
	}
	if value, ok := ccuo.mutation.AddedIdleCost(); ok {
		_spec.AddField(clustercost.FieldIdleCost, field.TypeFloat64, value)
	}
	if value, ok := ccuo.mutation.ManagementCost(); ok {
		_spec.SetField(clustercost.FieldManagementCost, field.TypeFloat64, value)
	}
	if value, ok := ccuo.mutation.AddedManagementCost(); ok {
		_spec.AddField(clustercost.FieldManagementCost, field.TypeFloat64, value)
	}
	_spec.Node.Schema = ccuo.schemaConfig.ClusterCost
	ctx = internal.NewSchemaConfigContext(ctx, ccuo.schemaConfig)
	_spec.AddModifiers(ccuo.modifiers...)
	_node = &ClusterCost{config: ccuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ccuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{clustercost.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ccuo.mutation.done = true
	return _node, nil
}
