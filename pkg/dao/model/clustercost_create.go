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

	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ClusterCostCreate is the builder for creating a ClusterCost entity.
type ClusterCostCreate struct {
	config
	mutation *ClusterCostMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetStartTime sets the "startTime" field.
func (ccc *ClusterCostCreate) SetStartTime(t time.Time) *ClusterCostCreate {
	ccc.mutation.SetStartTime(t)
	return ccc
}

// SetEndTime sets the "endTime" field.
func (ccc *ClusterCostCreate) SetEndTime(t time.Time) *ClusterCostCreate {
	ccc.mutation.SetEndTime(t)
	return ccc
}

// SetMinutes sets the "minutes" field.
func (ccc *ClusterCostCreate) SetMinutes(f float64) *ClusterCostCreate {
	ccc.mutation.SetMinutes(f)
	return ccc
}

// SetConnectorID sets the "connectorID" field.
func (ccc *ClusterCostCreate) SetConnectorID(o oid.ID) *ClusterCostCreate {
	ccc.mutation.SetConnectorID(o)
	return ccc
}

// SetClusterName sets the "clusterName" field.
func (ccc *ClusterCostCreate) SetClusterName(s string) *ClusterCostCreate {
	ccc.mutation.SetClusterName(s)
	return ccc
}

// SetTotalCost sets the "totalCost" field.
func (ccc *ClusterCostCreate) SetTotalCost(f float64) *ClusterCostCreate {
	ccc.mutation.SetTotalCost(f)
	return ccc
}

// SetNillableTotalCost sets the "totalCost" field if the given value is not nil.
func (ccc *ClusterCostCreate) SetNillableTotalCost(f *float64) *ClusterCostCreate {
	if f != nil {
		ccc.SetTotalCost(*f)
	}
	return ccc
}

// SetCurrency sets the "currency" field.
func (ccc *ClusterCostCreate) SetCurrency(i int) *ClusterCostCreate {
	ccc.mutation.SetCurrency(i)
	return ccc
}

// SetNillableCurrency sets the "currency" field if the given value is not nil.
func (ccc *ClusterCostCreate) SetNillableCurrency(i *int) *ClusterCostCreate {
	if i != nil {
		ccc.SetCurrency(*i)
	}
	return ccc
}

// SetCpuCost sets the "cpuCost" field.
func (ccc *ClusterCostCreate) SetCpuCost(f float64) *ClusterCostCreate {
	ccc.mutation.SetCpuCost(f)
	return ccc
}

// SetNillableCpuCost sets the "cpuCost" field if the given value is not nil.
func (ccc *ClusterCostCreate) SetNillableCpuCost(f *float64) *ClusterCostCreate {
	if f != nil {
		ccc.SetCpuCost(*f)
	}
	return ccc
}

// SetGpuCost sets the "gpuCost" field.
func (ccc *ClusterCostCreate) SetGpuCost(f float64) *ClusterCostCreate {
	ccc.mutation.SetGpuCost(f)
	return ccc
}

// SetNillableGpuCost sets the "gpuCost" field if the given value is not nil.
func (ccc *ClusterCostCreate) SetNillableGpuCost(f *float64) *ClusterCostCreate {
	if f != nil {
		ccc.SetGpuCost(*f)
	}
	return ccc
}

// SetRamCost sets the "ramCost" field.
func (ccc *ClusterCostCreate) SetRamCost(f float64) *ClusterCostCreate {
	ccc.mutation.SetRamCost(f)
	return ccc
}

// SetNillableRamCost sets the "ramCost" field if the given value is not nil.
func (ccc *ClusterCostCreate) SetNillableRamCost(f *float64) *ClusterCostCreate {
	if f != nil {
		ccc.SetRamCost(*f)
	}
	return ccc
}

// SetStorageCost sets the "storageCost" field.
func (ccc *ClusterCostCreate) SetStorageCost(f float64) *ClusterCostCreate {
	ccc.mutation.SetStorageCost(f)
	return ccc
}

// SetNillableStorageCost sets the "storageCost" field if the given value is not nil.
func (ccc *ClusterCostCreate) SetNillableStorageCost(f *float64) *ClusterCostCreate {
	if f != nil {
		ccc.SetStorageCost(*f)
	}
	return ccc
}

// SetAllocationCost sets the "allocationCost" field.
func (ccc *ClusterCostCreate) SetAllocationCost(f float64) *ClusterCostCreate {
	ccc.mutation.SetAllocationCost(f)
	return ccc
}

// SetNillableAllocationCost sets the "allocationCost" field if the given value is not nil.
func (ccc *ClusterCostCreate) SetNillableAllocationCost(f *float64) *ClusterCostCreate {
	if f != nil {
		ccc.SetAllocationCost(*f)
	}
	return ccc
}

// SetIdleCost sets the "idleCost" field.
func (ccc *ClusterCostCreate) SetIdleCost(f float64) *ClusterCostCreate {
	ccc.mutation.SetIdleCost(f)
	return ccc
}

// SetNillableIdleCost sets the "idleCost" field if the given value is not nil.
func (ccc *ClusterCostCreate) SetNillableIdleCost(f *float64) *ClusterCostCreate {
	if f != nil {
		ccc.SetIdleCost(*f)
	}
	return ccc
}

// SetManagementCost sets the "managementCost" field.
func (ccc *ClusterCostCreate) SetManagementCost(f float64) *ClusterCostCreate {
	ccc.mutation.SetManagementCost(f)
	return ccc
}

// SetNillableManagementCost sets the "managementCost" field if the given value is not nil.
func (ccc *ClusterCostCreate) SetNillableManagementCost(f *float64) *ClusterCostCreate {
	if f != nil {
		ccc.SetManagementCost(*f)
	}
	return ccc
}

// SetConnector sets the "connector" edge to the Connector entity.
func (ccc *ClusterCostCreate) SetConnector(c *Connector) *ClusterCostCreate {
	return ccc.SetConnectorID(c.ID)
}

// Mutation returns the ClusterCostMutation object of the builder.
func (ccc *ClusterCostCreate) Mutation() *ClusterCostMutation {
	return ccc.mutation
}

// Save creates the ClusterCost in the database.
func (ccc *ClusterCostCreate) Save(ctx context.Context) (*ClusterCost, error) {
	ccc.defaults()
	return withHooks[*ClusterCost, ClusterCostMutation](ctx, ccc.sqlSave, ccc.mutation, ccc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ccc *ClusterCostCreate) SaveX(ctx context.Context) *ClusterCost {
	v, err := ccc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccc *ClusterCostCreate) Exec(ctx context.Context) error {
	_, err := ccc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccc *ClusterCostCreate) ExecX(ctx context.Context) {
	if err := ccc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ccc *ClusterCostCreate) defaults() {
	if _, ok := ccc.mutation.TotalCost(); !ok {
		v := clustercost.DefaultTotalCost
		ccc.mutation.SetTotalCost(v)
	}
	if _, ok := ccc.mutation.CpuCost(); !ok {
		v := clustercost.DefaultCpuCost
		ccc.mutation.SetCpuCost(v)
	}
	if _, ok := ccc.mutation.GpuCost(); !ok {
		v := clustercost.DefaultGpuCost
		ccc.mutation.SetGpuCost(v)
	}
	if _, ok := ccc.mutation.RamCost(); !ok {
		v := clustercost.DefaultRamCost
		ccc.mutation.SetRamCost(v)
	}
	if _, ok := ccc.mutation.StorageCost(); !ok {
		v := clustercost.DefaultStorageCost
		ccc.mutation.SetStorageCost(v)
	}
	if _, ok := ccc.mutation.AllocationCost(); !ok {
		v := clustercost.DefaultAllocationCost
		ccc.mutation.SetAllocationCost(v)
	}
	if _, ok := ccc.mutation.IdleCost(); !ok {
		v := clustercost.DefaultIdleCost
		ccc.mutation.SetIdleCost(v)
	}
	if _, ok := ccc.mutation.ManagementCost(); !ok {
		v := clustercost.DefaultManagementCost
		ccc.mutation.SetManagementCost(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ccc *ClusterCostCreate) check() error {
	if _, ok := ccc.mutation.StartTime(); !ok {
		return &ValidationError{Name: "startTime", err: errors.New(`model: missing required field "ClusterCost.startTime"`)}
	}
	if _, ok := ccc.mutation.EndTime(); !ok {
		return &ValidationError{Name: "endTime", err: errors.New(`model: missing required field "ClusterCost.endTime"`)}
	}
	if _, ok := ccc.mutation.Minutes(); !ok {
		return &ValidationError{Name: "minutes", err: errors.New(`model: missing required field "ClusterCost.minutes"`)}
	}
	if _, ok := ccc.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connectorID", err: errors.New(`model: missing required field "ClusterCost.connectorID"`)}
	}
	if v, ok := ccc.mutation.ConnectorID(); ok {
		if err := clustercost.ConnectorIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "connectorID", err: fmt.Errorf(`model: validator failed for field "ClusterCost.connectorID": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.ClusterName(); !ok {
		return &ValidationError{Name: "clusterName", err: errors.New(`model: missing required field "ClusterCost.clusterName"`)}
	}
	if v, ok := ccc.mutation.ClusterName(); ok {
		if err := clustercost.ClusterNameValidator(v); err != nil {
			return &ValidationError{Name: "clusterName", err: fmt.Errorf(`model: validator failed for field "ClusterCost.clusterName": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.TotalCost(); !ok {
		return &ValidationError{Name: "totalCost", err: errors.New(`model: missing required field "ClusterCost.totalCost"`)}
	}
	if v, ok := ccc.mutation.TotalCost(); ok {
		if err := clustercost.TotalCostValidator(v); err != nil {
			return &ValidationError{Name: "totalCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.totalCost": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.CpuCost(); !ok {
		return &ValidationError{Name: "cpuCost", err: errors.New(`model: missing required field "ClusterCost.cpuCost"`)}
	}
	if v, ok := ccc.mutation.CpuCost(); ok {
		if err := clustercost.CpuCostValidator(v); err != nil {
			return &ValidationError{Name: "cpuCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.cpuCost": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.GpuCost(); !ok {
		return &ValidationError{Name: "gpuCost", err: errors.New(`model: missing required field "ClusterCost.gpuCost"`)}
	}
	if v, ok := ccc.mutation.GpuCost(); ok {
		if err := clustercost.GpuCostValidator(v); err != nil {
			return &ValidationError{Name: "gpuCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.gpuCost": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.RamCost(); !ok {
		return &ValidationError{Name: "ramCost", err: errors.New(`model: missing required field "ClusterCost.ramCost"`)}
	}
	if v, ok := ccc.mutation.RamCost(); ok {
		if err := clustercost.RamCostValidator(v); err != nil {
			return &ValidationError{Name: "ramCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.ramCost": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.StorageCost(); !ok {
		return &ValidationError{Name: "storageCost", err: errors.New(`model: missing required field "ClusterCost.storageCost"`)}
	}
	if v, ok := ccc.mutation.StorageCost(); ok {
		if err := clustercost.StorageCostValidator(v); err != nil {
			return &ValidationError{Name: "storageCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.storageCost": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.AllocationCost(); !ok {
		return &ValidationError{Name: "allocationCost", err: errors.New(`model: missing required field "ClusterCost.allocationCost"`)}
	}
	if v, ok := ccc.mutation.AllocationCost(); ok {
		if err := clustercost.AllocationCostValidator(v); err != nil {
			return &ValidationError{Name: "allocationCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.allocationCost": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.IdleCost(); !ok {
		return &ValidationError{Name: "idleCost", err: errors.New(`model: missing required field "ClusterCost.idleCost"`)}
	}
	if v, ok := ccc.mutation.IdleCost(); ok {
		if err := clustercost.IdleCostValidator(v); err != nil {
			return &ValidationError{Name: "idleCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.idleCost": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.ManagementCost(); !ok {
		return &ValidationError{Name: "managementCost", err: errors.New(`model: missing required field "ClusterCost.managementCost"`)}
	}
	if v, ok := ccc.mutation.ManagementCost(); ok {
		if err := clustercost.ManagementCostValidator(v); err != nil {
			return &ValidationError{Name: "managementCost", err: fmt.Errorf(`model: validator failed for field "ClusterCost.managementCost": %w`, err)}
		}
	}
	if _, ok := ccc.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connector", err: errors.New(`model: missing required edge "ClusterCost.connector"`)}
	}
	return nil
}

func (ccc *ClusterCostCreate) sqlSave(ctx context.Context) (*ClusterCost, error) {
	if err := ccc.check(); err != nil {
		return nil, err
	}
	_node, _spec := ccc.createSpec()
	if err := sqlgraph.CreateNode(ctx, ccc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ccc.mutation.id = &_node.ID
	ccc.mutation.done = true
	return _node, nil
}

func (ccc *ClusterCostCreate) createSpec() (*ClusterCost, *sqlgraph.CreateSpec) {
	var (
		_node = &ClusterCost{config: ccc.config}
		_spec = sqlgraph.NewCreateSpec(clustercost.Table, sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt))
	)
	_spec.Schema = ccc.schemaConfig.ClusterCost
	_spec.OnConflict = ccc.conflict
	if value, ok := ccc.mutation.StartTime(); ok {
		_spec.SetField(clustercost.FieldStartTime, field.TypeTime, value)
		_node.StartTime = value
	}
	if value, ok := ccc.mutation.EndTime(); ok {
		_spec.SetField(clustercost.FieldEndTime, field.TypeTime, value)
		_node.EndTime = value
	}
	if value, ok := ccc.mutation.Minutes(); ok {
		_spec.SetField(clustercost.FieldMinutes, field.TypeFloat64, value)
		_node.Minutes = value
	}
	if value, ok := ccc.mutation.ClusterName(); ok {
		_spec.SetField(clustercost.FieldClusterName, field.TypeString, value)
		_node.ClusterName = value
	}
	if value, ok := ccc.mutation.TotalCost(); ok {
		_spec.SetField(clustercost.FieldTotalCost, field.TypeFloat64, value)
		_node.TotalCost = value
	}
	if value, ok := ccc.mutation.Currency(); ok {
		_spec.SetField(clustercost.FieldCurrency, field.TypeInt, value)
		_node.Currency = value
	}
	if value, ok := ccc.mutation.CpuCost(); ok {
		_spec.SetField(clustercost.FieldCpuCost, field.TypeFloat64, value)
		_node.CpuCost = value
	}
	if value, ok := ccc.mutation.GpuCost(); ok {
		_spec.SetField(clustercost.FieldGpuCost, field.TypeFloat64, value)
		_node.GpuCost = value
	}
	if value, ok := ccc.mutation.RamCost(); ok {
		_spec.SetField(clustercost.FieldRamCost, field.TypeFloat64, value)
		_node.RamCost = value
	}
	if value, ok := ccc.mutation.StorageCost(); ok {
		_spec.SetField(clustercost.FieldStorageCost, field.TypeFloat64, value)
		_node.StorageCost = value
	}
	if value, ok := ccc.mutation.AllocationCost(); ok {
		_spec.SetField(clustercost.FieldAllocationCost, field.TypeFloat64, value)
		_node.AllocationCost = value
	}
	if value, ok := ccc.mutation.IdleCost(); ok {
		_spec.SetField(clustercost.FieldIdleCost, field.TypeFloat64, value)
		_node.IdleCost = value
	}
	if value, ok := ccc.mutation.ManagementCost(); ok {
		_spec.SetField(clustercost.FieldManagementCost, field.TypeFloat64, value)
		_node.ManagementCost = value
	}
	if nodes := ccc.mutation.ConnectorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   clustercost.ConnectorTable,
			Columns: []string{clustercost.ConnectorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = ccc.schemaConfig.ClusterCost
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ConnectorID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ClusterCost.Create().
//		SetStartTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ClusterCostUpsert) {
//			SetStartTime(v+v).
//		}).
//		Exec(ctx)
func (ccc *ClusterCostCreate) OnConflict(opts ...sql.ConflictOption) *ClusterCostUpsertOne {
	ccc.conflict = opts
	return &ClusterCostUpsertOne{
		create: ccc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ClusterCost.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ccc *ClusterCostCreate) OnConflictColumns(columns ...string) *ClusterCostUpsertOne {
	ccc.conflict = append(ccc.conflict, sql.ConflictColumns(columns...))
	return &ClusterCostUpsertOne{
		create: ccc,
	}
}

type (
	// ClusterCostUpsertOne is the builder for "upsert"-ing
	//  one ClusterCost node.
	ClusterCostUpsertOne struct {
		create *ClusterCostCreate
	}

	// ClusterCostUpsert is the "OnConflict" setter.
	ClusterCostUpsert struct {
		*sql.UpdateSet
	}
)

// SetTotalCost sets the "totalCost" field.
func (u *ClusterCostUpsert) SetTotalCost(v float64) *ClusterCostUpsert {
	u.Set(clustercost.FieldTotalCost, v)
	return u
}

// UpdateTotalCost sets the "totalCost" field to the value that was provided on create.
func (u *ClusterCostUpsert) UpdateTotalCost() *ClusterCostUpsert {
	u.SetExcluded(clustercost.FieldTotalCost)
	return u
}

// AddTotalCost adds v to the "totalCost" field.
func (u *ClusterCostUpsert) AddTotalCost(v float64) *ClusterCostUpsert {
	u.Add(clustercost.FieldTotalCost, v)
	return u
}

// SetCurrency sets the "currency" field.
func (u *ClusterCostUpsert) SetCurrency(v int) *ClusterCostUpsert {
	u.Set(clustercost.FieldCurrency, v)
	return u
}

// UpdateCurrency sets the "currency" field to the value that was provided on create.
func (u *ClusterCostUpsert) UpdateCurrency() *ClusterCostUpsert {
	u.SetExcluded(clustercost.FieldCurrency)
	return u
}

// AddCurrency adds v to the "currency" field.
func (u *ClusterCostUpsert) AddCurrency(v int) *ClusterCostUpsert {
	u.Add(clustercost.FieldCurrency, v)
	return u
}

// ClearCurrency clears the value of the "currency" field.
func (u *ClusterCostUpsert) ClearCurrency() *ClusterCostUpsert {
	u.SetNull(clustercost.FieldCurrency)
	return u
}

// SetCpuCost sets the "cpuCost" field.
func (u *ClusterCostUpsert) SetCpuCost(v float64) *ClusterCostUpsert {
	u.Set(clustercost.FieldCpuCost, v)
	return u
}

// UpdateCpuCost sets the "cpuCost" field to the value that was provided on create.
func (u *ClusterCostUpsert) UpdateCpuCost() *ClusterCostUpsert {
	u.SetExcluded(clustercost.FieldCpuCost)
	return u
}

// AddCpuCost adds v to the "cpuCost" field.
func (u *ClusterCostUpsert) AddCpuCost(v float64) *ClusterCostUpsert {
	u.Add(clustercost.FieldCpuCost, v)
	return u
}

// SetGpuCost sets the "gpuCost" field.
func (u *ClusterCostUpsert) SetGpuCost(v float64) *ClusterCostUpsert {
	u.Set(clustercost.FieldGpuCost, v)
	return u
}

// UpdateGpuCost sets the "gpuCost" field to the value that was provided on create.
func (u *ClusterCostUpsert) UpdateGpuCost() *ClusterCostUpsert {
	u.SetExcluded(clustercost.FieldGpuCost)
	return u
}

// AddGpuCost adds v to the "gpuCost" field.
func (u *ClusterCostUpsert) AddGpuCost(v float64) *ClusterCostUpsert {
	u.Add(clustercost.FieldGpuCost, v)
	return u
}

// SetRamCost sets the "ramCost" field.
func (u *ClusterCostUpsert) SetRamCost(v float64) *ClusterCostUpsert {
	u.Set(clustercost.FieldRamCost, v)
	return u
}

// UpdateRamCost sets the "ramCost" field to the value that was provided on create.
func (u *ClusterCostUpsert) UpdateRamCost() *ClusterCostUpsert {
	u.SetExcluded(clustercost.FieldRamCost)
	return u
}

// AddRamCost adds v to the "ramCost" field.
func (u *ClusterCostUpsert) AddRamCost(v float64) *ClusterCostUpsert {
	u.Add(clustercost.FieldRamCost, v)
	return u
}

// SetStorageCost sets the "storageCost" field.
func (u *ClusterCostUpsert) SetStorageCost(v float64) *ClusterCostUpsert {
	u.Set(clustercost.FieldStorageCost, v)
	return u
}

// UpdateStorageCost sets the "storageCost" field to the value that was provided on create.
func (u *ClusterCostUpsert) UpdateStorageCost() *ClusterCostUpsert {
	u.SetExcluded(clustercost.FieldStorageCost)
	return u
}

// AddStorageCost adds v to the "storageCost" field.
func (u *ClusterCostUpsert) AddStorageCost(v float64) *ClusterCostUpsert {
	u.Add(clustercost.FieldStorageCost, v)
	return u
}

// SetAllocationCost sets the "allocationCost" field.
func (u *ClusterCostUpsert) SetAllocationCost(v float64) *ClusterCostUpsert {
	u.Set(clustercost.FieldAllocationCost, v)
	return u
}

// UpdateAllocationCost sets the "allocationCost" field to the value that was provided on create.
func (u *ClusterCostUpsert) UpdateAllocationCost() *ClusterCostUpsert {
	u.SetExcluded(clustercost.FieldAllocationCost)
	return u
}

// AddAllocationCost adds v to the "allocationCost" field.
func (u *ClusterCostUpsert) AddAllocationCost(v float64) *ClusterCostUpsert {
	u.Add(clustercost.FieldAllocationCost, v)
	return u
}

// SetIdleCost sets the "idleCost" field.
func (u *ClusterCostUpsert) SetIdleCost(v float64) *ClusterCostUpsert {
	u.Set(clustercost.FieldIdleCost, v)
	return u
}

// UpdateIdleCost sets the "idleCost" field to the value that was provided on create.
func (u *ClusterCostUpsert) UpdateIdleCost() *ClusterCostUpsert {
	u.SetExcluded(clustercost.FieldIdleCost)
	return u
}

// AddIdleCost adds v to the "idleCost" field.
func (u *ClusterCostUpsert) AddIdleCost(v float64) *ClusterCostUpsert {
	u.Add(clustercost.FieldIdleCost, v)
	return u
}

// SetManagementCost sets the "managementCost" field.
func (u *ClusterCostUpsert) SetManagementCost(v float64) *ClusterCostUpsert {
	u.Set(clustercost.FieldManagementCost, v)
	return u
}

// UpdateManagementCost sets the "managementCost" field to the value that was provided on create.
func (u *ClusterCostUpsert) UpdateManagementCost() *ClusterCostUpsert {
	u.SetExcluded(clustercost.FieldManagementCost)
	return u
}

// AddManagementCost adds v to the "managementCost" field.
func (u *ClusterCostUpsert) AddManagementCost(v float64) *ClusterCostUpsert {
	u.Add(clustercost.FieldManagementCost, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.ClusterCost.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ClusterCostUpsertOne) UpdateNewValues() *ClusterCostUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.StartTime(); exists {
			s.SetIgnore(clustercost.FieldStartTime)
		}
		if _, exists := u.create.mutation.EndTime(); exists {
			s.SetIgnore(clustercost.FieldEndTime)
		}
		if _, exists := u.create.mutation.Minutes(); exists {
			s.SetIgnore(clustercost.FieldMinutes)
		}
		if _, exists := u.create.mutation.ConnectorID(); exists {
			s.SetIgnore(clustercost.FieldConnectorID)
		}
		if _, exists := u.create.mutation.ClusterName(); exists {
			s.SetIgnore(clustercost.FieldClusterName)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ClusterCost.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ClusterCostUpsertOne) Ignore() *ClusterCostUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ClusterCostUpsertOne) DoNothing() *ClusterCostUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ClusterCostCreate.OnConflict
// documentation for more info.
func (u *ClusterCostUpsertOne) Update(set func(*ClusterCostUpsert)) *ClusterCostUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ClusterCostUpsert{UpdateSet: update})
	}))
	return u
}

// SetTotalCost sets the "totalCost" field.
func (u *ClusterCostUpsertOne) SetTotalCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetTotalCost(v)
	})
}

// AddTotalCost adds v to the "totalCost" field.
func (u *ClusterCostUpsertOne) AddTotalCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddTotalCost(v)
	})
}

// UpdateTotalCost sets the "totalCost" field to the value that was provided on create.
func (u *ClusterCostUpsertOne) UpdateTotalCost() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateTotalCost()
	})
}

// SetCurrency sets the "currency" field.
func (u *ClusterCostUpsertOne) SetCurrency(v int) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetCurrency(v)
	})
}

// AddCurrency adds v to the "currency" field.
func (u *ClusterCostUpsertOne) AddCurrency(v int) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddCurrency(v)
	})
}

// UpdateCurrency sets the "currency" field to the value that was provided on create.
func (u *ClusterCostUpsertOne) UpdateCurrency() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateCurrency()
	})
}

// ClearCurrency clears the value of the "currency" field.
func (u *ClusterCostUpsertOne) ClearCurrency() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.ClearCurrency()
	})
}

// SetCpuCost sets the "cpuCost" field.
func (u *ClusterCostUpsertOne) SetCpuCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetCpuCost(v)
	})
}

// AddCpuCost adds v to the "cpuCost" field.
func (u *ClusterCostUpsertOne) AddCpuCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddCpuCost(v)
	})
}

// UpdateCpuCost sets the "cpuCost" field to the value that was provided on create.
func (u *ClusterCostUpsertOne) UpdateCpuCost() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateCpuCost()
	})
}

// SetGpuCost sets the "gpuCost" field.
func (u *ClusterCostUpsertOne) SetGpuCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetGpuCost(v)
	})
}

// AddGpuCost adds v to the "gpuCost" field.
func (u *ClusterCostUpsertOne) AddGpuCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddGpuCost(v)
	})
}

// UpdateGpuCost sets the "gpuCost" field to the value that was provided on create.
func (u *ClusterCostUpsertOne) UpdateGpuCost() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateGpuCost()
	})
}

// SetRamCost sets the "ramCost" field.
func (u *ClusterCostUpsertOne) SetRamCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetRamCost(v)
	})
}

// AddRamCost adds v to the "ramCost" field.
func (u *ClusterCostUpsertOne) AddRamCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddRamCost(v)
	})
}

// UpdateRamCost sets the "ramCost" field to the value that was provided on create.
func (u *ClusterCostUpsertOne) UpdateRamCost() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateRamCost()
	})
}

// SetStorageCost sets the "storageCost" field.
func (u *ClusterCostUpsertOne) SetStorageCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetStorageCost(v)
	})
}

// AddStorageCost adds v to the "storageCost" field.
func (u *ClusterCostUpsertOne) AddStorageCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddStorageCost(v)
	})
}

// UpdateStorageCost sets the "storageCost" field to the value that was provided on create.
func (u *ClusterCostUpsertOne) UpdateStorageCost() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateStorageCost()
	})
}

// SetAllocationCost sets the "allocationCost" field.
func (u *ClusterCostUpsertOne) SetAllocationCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetAllocationCost(v)
	})
}

// AddAllocationCost adds v to the "allocationCost" field.
func (u *ClusterCostUpsertOne) AddAllocationCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddAllocationCost(v)
	})
}

// UpdateAllocationCost sets the "allocationCost" field to the value that was provided on create.
func (u *ClusterCostUpsertOne) UpdateAllocationCost() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateAllocationCost()
	})
}

// SetIdleCost sets the "idleCost" field.
func (u *ClusterCostUpsertOne) SetIdleCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetIdleCost(v)
	})
}

// AddIdleCost adds v to the "idleCost" field.
func (u *ClusterCostUpsertOne) AddIdleCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddIdleCost(v)
	})
}

// UpdateIdleCost sets the "idleCost" field to the value that was provided on create.
func (u *ClusterCostUpsertOne) UpdateIdleCost() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateIdleCost()
	})
}

// SetManagementCost sets the "managementCost" field.
func (u *ClusterCostUpsertOne) SetManagementCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetManagementCost(v)
	})
}

// AddManagementCost adds v to the "managementCost" field.
func (u *ClusterCostUpsertOne) AddManagementCost(v float64) *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddManagementCost(v)
	})
}

// UpdateManagementCost sets the "managementCost" field to the value that was provided on create.
func (u *ClusterCostUpsertOne) UpdateManagementCost() *ClusterCostUpsertOne {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateManagementCost()
	})
}

// Exec executes the query.
func (u *ClusterCostUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ClusterCostCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ClusterCostUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ClusterCostUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ClusterCostUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ClusterCostCreateBulk is the builder for creating many ClusterCost entities in bulk.
type ClusterCostCreateBulk struct {
	config
	builders []*ClusterCostCreate
	conflict []sql.ConflictOption
}

// Save creates the ClusterCost entities in the database.
func (cccb *ClusterCostCreateBulk) Save(ctx context.Context) ([]*ClusterCost, error) {
	specs := make([]*sqlgraph.CreateSpec, len(cccb.builders))
	nodes := make([]*ClusterCost, len(cccb.builders))
	mutators := make([]Mutator, len(cccb.builders))
	for i := range cccb.builders {
		func(i int, root context.Context) {
			builder := cccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ClusterCostMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, cccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = cccb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, cccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, cccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (cccb *ClusterCostCreateBulk) SaveX(ctx context.Context) []*ClusterCost {
	v, err := cccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cccb *ClusterCostCreateBulk) Exec(ctx context.Context) error {
	_, err := cccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cccb *ClusterCostCreateBulk) ExecX(ctx context.Context) {
	if err := cccb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ClusterCost.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ClusterCostUpsert) {
//			SetStartTime(v+v).
//		}).
//		Exec(ctx)
func (cccb *ClusterCostCreateBulk) OnConflict(opts ...sql.ConflictOption) *ClusterCostUpsertBulk {
	cccb.conflict = opts
	return &ClusterCostUpsertBulk{
		create: cccb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ClusterCost.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (cccb *ClusterCostCreateBulk) OnConflictColumns(columns ...string) *ClusterCostUpsertBulk {
	cccb.conflict = append(cccb.conflict, sql.ConflictColumns(columns...))
	return &ClusterCostUpsertBulk{
		create: cccb,
	}
}

// ClusterCostUpsertBulk is the builder for "upsert"-ing
// a bulk of ClusterCost nodes.
type ClusterCostUpsertBulk struct {
	create *ClusterCostCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ClusterCost.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ClusterCostUpsertBulk) UpdateNewValues() *ClusterCostUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.StartTime(); exists {
				s.SetIgnore(clustercost.FieldStartTime)
			}
			if _, exists := b.mutation.EndTime(); exists {
				s.SetIgnore(clustercost.FieldEndTime)
			}
			if _, exists := b.mutation.Minutes(); exists {
				s.SetIgnore(clustercost.FieldMinutes)
			}
			if _, exists := b.mutation.ConnectorID(); exists {
				s.SetIgnore(clustercost.FieldConnectorID)
			}
			if _, exists := b.mutation.ClusterName(); exists {
				s.SetIgnore(clustercost.FieldClusterName)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ClusterCost.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ClusterCostUpsertBulk) Ignore() *ClusterCostUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ClusterCostUpsertBulk) DoNothing() *ClusterCostUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ClusterCostCreateBulk.OnConflict
// documentation for more info.
func (u *ClusterCostUpsertBulk) Update(set func(*ClusterCostUpsert)) *ClusterCostUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ClusterCostUpsert{UpdateSet: update})
	}))
	return u
}

// SetTotalCost sets the "totalCost" field.
func (u *ClusterCostUpsertBulk) SetTotalCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetTotalCost(v)
	})
}

// AddTotalCost adds v to the "totalCost" field.
func (u *ClusterCostUpsertBulk) AddTotalCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddTotalCost(v)
	})
}

// UpdateTotalCost sets the "totalCost" field to the value that was provided on create.
func (u *ClusterCostUpsertBulk) UpdateTotalCost() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateTotalCost()
	})
}

// SetCurrency sets the "currency" field.
func (u *ClusterCostUpsertBulk) SetCurrency(v int) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetCurrency(v)
	})
}

// AddCurrency adds v to the "currency" field.
func (u *ClusterCostUpsertBulk) AddCurrency(v int) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddCurrency(v)
	})
}

// UpdateCurrency sets the "currency" field to the value that was provided on create.
func (u *ClusterCostUpsertBulk) UpdateCurrency() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateCurrency()
	})
}

// ClearCurrency clears the value of the "currency" field.
func (u *ClusterCostUpsertBulk) ClearCurrency() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.ClearCurrency()
	})
}

// SetCpuCost sets the "cpuCost" field.
func (u *ClusterCostUpsertBulk) SetCpuCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetCpuCost(v)
	})
}

// AddCpuCost adds v to the "cpuCost" field.
func (u *ClusterCostUpsertBulk) AddCpuCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddCpuCost(v)
	})
}

// UpdateCpuCost sets the "cpuCost" field to the value that was provided on create.
func (u *ClusterCostUpsertBulk) UpdateCpuCost() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateCpuCost()
	})
}

// SetGpuCost sets the "gpuCost" field.
func (u *ClusterCostUpsertBulk) SetGpuCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetGpuCost(v)
	})
}

// AddGpuCost adds v to the "gpuCost" field.
func (u *ClusterCostUpsertBulk) AddGpuCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddGpuCost(v)
	})
}

// UpdateGpuCost sets the "gpuCost" field to the value that was provided on create.
func (u *ClusterCostUpsertBulk) UpdateGpuCost() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateGpuCost()
	})
}

// SetRamCost sets the "ramCost" field.
func (u *ClusterCostUpsertBulk) SetRamCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetRamCost(v)
	})
}

// AddRamCost adds v to the "ramCost" field.
func (u *ClusterCostUpsertBulk) AddRamCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddRamCost(v)
	})
}

// UpdateRamCost sets the "ramCost" field to the value that was provided on create.
func (u *ClusterCostUpsertBulk) UpdateRamCost() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateRamCost()
	})
}

// SetStorageCost sets the "storageCost" field.
func (u *ClusterCostUpsertBulk) SetStorageCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetStorageCost(v)
	})
}

// AddStorageCost adds v to the "storageCost" field.
func (u *ClusterCostUpsertBulk) AddStorageCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddStorageCost(v)
	})
}

// UpdateStorageCost sets the "storageCost" field to the value that was provided on create.
func (u *ClusterCostUpsertBulk) UpdateStorageCost() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateStorageCost()
	})
}

// SetAllocationCost sets the "allocationCost" field.
func (u *ClusterCostUpsertBulk) SetAllocationCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetAllocationCost(v)
	})
}

// AddAllocationCost adds v to the "allocationCost" field.
func (u *ClusterCostUpsertBulk) AddAllocationCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddAllocationCost(v)
	})
}

// UpdateAllocationCost sets the "allocationCost" field to the value that was provided on create.
func (u *ClusterCostUpsertBulk) UpdateAllocationCost() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateAllocationCost()
	})
}

// SetIdleCost sets the "idleCost" field.
func (u *ClusterCostUpsertBulk) SetIdleCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetIdleCost(v)
	})
}

// AddIdleCost adds v to the "idleCost" field.
func (u *ClusterCostUpsertBulk) AddIdleCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddIdleCost(v)
	})
}

// UpdateIdleCost sets the "idleCost" field to the value that was provided on create.
func (u *ClusterCostUpsertBulk) UpdateIdleCost() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateIdleCost()
	})
}

// SetManagementCost sets the "managementCost" field.
func (u *ClusterCostUpsertBulk) SetManagementCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.SetManagementCost(v)
	})
}

// AddManagementCost adds v to the "managementCost" field.
func (u *ClusterCostUpsertBulk) AddManagementCost(v float64) *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.AddManagementCost(v)
	})
}

// UpdateManagementCost sets the "managementCost" field to the value that was provided on create.
func (u *ClusterCostUpsertBulk) UpdateManagementCost() *ClusterCostUpsertBulk {
	return u.Update(func(s *ClusterCostUpsert) {
		s.UpdateManagementCost()
	})
}

// Exec executes the query.
func (u *ClusterCostUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ClusterCostCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ClusterCostCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ClusterCostUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
