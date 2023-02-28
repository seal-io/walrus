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

	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
)

// AllocationCostCreate is the builder for creating a AllocationCost entity.
type AllocationCostCreate struct {
	config
	mutation *AllocationCostMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetStartTime sets the "startTime" field.
func (acc *AllocationCostCreate) SetStartTime(t time.Time) *AllocationCostCreate {
	acc.mutation.SetStartTime(t)
	return acc
}

// SetEndTime sets the "endTime" field.
func (acc *AllocationCostCreate) SetEndTime(t time.Time) *AllocationCostCreate {
	acc.mutation.SetEndTime(t)
	return acc
}

// SetMinutes sets the "minutes" field.
func (acc *AllocationCostCreate) SetMinutes(f float64) *AllocationCostCreate {
	acc.mutation.SetMinutes(f)
	return acc
}

// SetConnectorID sets the "connectorID" field.
func (acc *AllocationCostCreate) SetConnectorID(t types.ID) *AllocationCostCreate {
	acc.mutation.SetConnectorID(t)
	return acc
}

// SetName sets the "name" field.
func (acc *AllocationCostCreate) SetName(s string) *AllocationCostCreate {
	acc.mutation.SetName(s)
	return acc
}

// SetFingerprint sets the "fingerprint" field.
func (acc *AllocationCostCreate) SetFingerprint(s string) *AllocationCostCreate {
	acc.mutation.SetFingerprint(s)
	return acc
}

// SetClusterName sets the "clusterName" field.
func (acc *AllocationCostCreate) SetClusterName(s string) *AllocationCostCreate {
	acc.mutation.SetClusterName(s)
	return acc
}

// SetNillableClusterName sets the "clusterName" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableClusterName(s *string) *AllocationCostCreate {
	if s != nil {
		acc.SetClusterName(*s)
	}
	return acc
}

// SetNamespace sets the "namespace" field.
func (acc *AllocationCostCreate) SetNamespace(s string) *AllocationCostCreate {
	acc.mutation.SetNamespace(s)
	return acc
}

// SetNillableNamespace sets the "namespace" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableNamespace(s *string) *AllocationCostCreate {
	if s != nil {
		acc.SetNamespace(*s)
	}
	return acc
}

// SetNode sets the "node" field.
func (acc *AllocationCostCreate) SetNode(s string) *AllocationCostCreate {
	acc.mutation.SetNode(s)
	return acc
}

// SetNillableNode sets the "node" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableNode(s *string) *AllocationCostCreate {
	if s != nil {
		acc.SetNode(*s)
	}
	return acc
}

// SetController sets the "controller" field.
func (acc *AllocationCostCreate) SetController(s string) *AllocationCostCreate {
	acc.mutation.SetController(s)
	return acc
}

// SetNillableController sets the "controller" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableController(s *string) *AllocationCostCreate {
	if s != nil {
		acc.SetController(*s)
	}
	return acc
}

// SetControllerKind sets the "controllerKind" field.
func (acc *AllocationCostCreate) SetControllerKind(s string) *AllocationCostCreate {
	acc.mutation.SetControllerKind(s)
	return acc
}

// SetNillableControllerKind sets the "controllerKind" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableControllerKind(s *string) *AllocationCostCreate {
	if s != nil {
		acc.SetControllerKind(*s)
	}
	return acc
}

// SetPod sets the "pod" field.
func (acc *AllocationCostCreate) SetPod(s string) *AllocationCostCreate {
	acc.mutation.SetPod(s)
	return acc
}

// SetNillablePod sets the "pod" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillablePod(s *string) *AllocationCostCreate {
	if s != nil {
		acc.SetPod(*s)
	}
	return acc
}

// SetContainer sets the "container" field.
func (acc *AllocationCostCreate) SetContainer(s string) *AllocationCostCreate {
	acc.mutation.SetContainer(s)
	return acc
}

// SetNillableContainer sets the "container" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableContainer(s *string) *AllocationCostCreate {
	if s != nil {
		acc.SetContainer(*s)
	}
	return acc
}

// SetPvs sets the "pvs" field.
func (acc *AllocationCostCreate) SetPvs(mc map[string]types.PVCost) *AllocationCostCreate {
	acc.mutation.SetPvs(mc)
	return acc
}

// SetLabels sets the "labels" field.
func (acc *AllocationCostCreate) SetLabels(m map[string]string) *AllocationCostCreate {
	acc.mutation.SetLabels(m)
	return acc
}

// SetTotalCost sets the "totalCost" field.
func (acc *AllocationCostCreate) SetTotalCost(f float64) *AllocationCostCreate {
	acc.mutation.SetTotalCost(f)
	return acc
}

// SetNillableTotalCost sets the "totalCost" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableTotalCost(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetTotalCost(*f)
	}
	return acc
}

// SetCurrency sets the "currency" field.
func (acc *AllocationCostCreate) SetCurrency(i int) *AllocationCostCreate {
	acc.mutation.SetCurrency(i)
	return acc
}

// SetNillableCurrency sets the "currency" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableCurrency(i *int) *AllocationCostCreate {
	if i != nil {
		acc.SetCurrency(*i)
	}
	return acc
}

// SetCpuCost sets the "cpuCost" field.
func (acc *AllocationCostCreate) SetCpuCost(f float64) *AllocationCostCreate {
	acc.mutation.SetCpuCost(f)
	return acc
}

// SetNillableCpuCost sets the "cpuCost" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableCpuCost(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetCpuCost(*f)
	}
	return acc
}

// SetCpuCoreRequest sets the "cpuCoreRequest" field.
func (acc *AllocationCostCreate) SetCpuCoreRequest(f float64) *AllocationCostCreate {
	acc.mutation.SetCpuCoreRequest(f)
	return acc
}

// SetNillableCpuCoreRequest sets the "cpuCoreRequest" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableCpuCoreRequest(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetCpuCoreRequest(*f)
	}
	return acc
}

// SetGpuCost sets the "gpuCost" field.
func (acc *AllocationCostCreate) SetGpuCost(f float64) *AllocationCostCreate {
	acc.mutation.SetGpuCost(f)
	return acc
}

// SetNillableGpuCost sets the "gpuCost" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableGpuCost(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetGpuCost(*f)
	}
	return acc
}

// SetGpuCount sets the "gpuCount" field.
func (acc *AllocationCostCreate) SetGpuCount(f float64) *AllocationCostCreate {
	acc.mutation.SetGpuCount(f)
	return acc
}

// SetNillableGpuCount sets the "gpuCount" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableGpuCount(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetGpuCount(*f)
	}
	return acc
}

// SetRamCost sets the "ramCost" field.
func (acc *AllocationCostCreate) SetRamCost(f float64) *AllocationCostCreate {
	acc.mutation.SetRamCost(f)
	return acc
}

// SetNillableRamCost sets the "ramCost" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableRamCost(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetRamCost(*f)
	}
	return acc
}

// SetRamByteRequest sets the "ramByteRequest" field.
func (acc *AllocationCostCreate) SetRamByteRequest(f float64) *AllocationCostCreate {
	acc.mutation.SetRamByteRequest(f)
	return acc
}

// SetNillableRamByteRequest sets the "ramByteRequest" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableRamByteRequest(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetRamByteRequest(*f)
	}
	return acc
}

// SetPvCost sets the "pvCost" field.
func (acc *AllocationCostCreate) SetPvCost(f float64) *AllocationCostCreate {
	acc.mutation.SetPvCost(f)
	return acc
}

// SetNillablePvCost sets the "pvCost" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillablePvCost(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetPvCost(*f)
	}
	return acc
}

// SetPvBytes sets the "pvBytes" field.
func (acc *AllocationCostCreate) SetPvBytes(f float64) *AllocationCostCreate {
	acc.mutation.SetPvBytes(f)
	return acc
}

// SetNillablePvBytes sets the "pvBytes" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillablePvBytes(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetPvBytes(*f)
	}
	return acc
}

// SetCpuCoreUsageAverage sets the "cpuCoreUsageAverage" field.
func (acc *AllocationCostCreate) SetCpuCoreUsageAverage(f float64) *AllocationCostCreate {
	acc.mutation.SetCpuCoreUsageAverage(f)
	return acc
}

// SetNillableCpuCoreUsageAverage sets the "cpuCoreUsageAverage" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableCpuCoreUsageAverage(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetCpuCoreUsageAverage(*f)
	}
	return acc
}

// SetCpuCoreUsageMax sets the "cpuCoreUsageMax" field.
func (acc *AllocationCostCreate) SetCpuCoreUsageMax(f float64) *AllocationCostCreate {
	acc.mutation.SetCpuCoreUsageMax(f)
	return acc
}

// SetNillableCpuCoreUsageMax sets the "cpuCoreUsageMax" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableCpuCoreUsageMax(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetCpuCoreUsageMax(*f)
	}
	return acc
}

// SetRamByteUsageAverage sets the "ramByteUsageAverage" field.
func (acc *AllocationCostCreate) SetRamByteUsageAverage(f float64) *AllocationCostCreate {
	acc.mutation.SetRamByteUsageAverage(f)
	return acc
}

// SetNillableRamByteUsageAverage sets the "ramByteUsageAverage" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableRamByteUsageAverage(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetRamByteUsageAverage(*f)
	}
	return acc
}

// SetRamByteUsageMax sets the "ramByteUsageMax" field.
func (acc *AllocationCostCreate) SetRamByteUsageMax(f float64) *AllocationCostCreate {
	acc.mutation.SetRamByteUsageMax(f)
	return acc
}

// SetNillableRamByteUsageMax sets the "ramByteUsageMax" field if the given value is not nil.
func (acc *AllocationCostCreate) SetNillableRamByteUsageMax(f *float64) *AllocationCostCreate {
	if f != nil {
		acc.SetRamByteUsageMax(*f)
	}
	return acc
}

// SetConnector sets the "connector" edge to the Connector entity.
func (acc *AllocationCostCreate) SetConnector(c *Connector) *AllocationCostCreate {
	return acc.SetConnectorID(c.ID)
}

// Mutation returns the AllocationCostMutation object of the builder.
func (acc *AllocationCostCreate) Mutation() *AllocationCostMutation {
	return acc.mutation
}

// Save creates the AllocationCost in the database.
func (acc *AllocationCostCreate) Save(ctx context.Context) (*AllocationCost, error) {
	acc.defaults()
	return withHooks[*AllocationCost, AllocationCostMutation](ctx, acc.sqlSave, acc.mutation, acc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (acc *AllocationCostCreate) SaveX(ctx context.Context) *AllocationCost {
	v, err := acc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acc *AllocationCostCreate) Exec(ctx context.Context) error {
	_, err := acc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acc *AllocationCostCreate) ExecX(ctx context.Context) {
	if err := acc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (acc *AllocationCostCreate) defaults() {
	if _, ok := acc.mutation.Pvs(); !ok {
		v := allocationcost.DefaultPvs
		acc.mutation.SetPvs(v)
	}
	if _, ok := acc.mutation.Labels(); !ok {
		v := allocationcost.DefaultLabels
		acc.mutation.SetLabels(v)
	}
	if _, ok := acc.mutation.TotalCost(); !ok {
		v := allocationcost.DefaultTotalCost
		acc.mutation.SetTotalCost(v)
	}
	if _, ok := acc.mutation.CpuCost(); !ok {
		v := allocationcost.DefaultCpuCost
		acc.mutation.SetCpuCost(v)
	}
	if _, ok := acc.mutation.CpuCoreRequest(); !ok {
		v := allocationcost.DefaultCpuCoreRequest
		acc.mutation.SetCpuCoreRequest(v)
	}
	if _, ok := acc.mutation.GpuCost(); !ok {
		v := allocationcost.DefaultGpuCost
		acc.mutation.SetGpuCost(v)
	}
	if _, ok := acc.mutation.GpuCount(); !ok {
		v := allocationcost.DefaultGpuCount
		acc.mutation.SetGpuCount(v)
	}
	if _, ok := acc.mutation.RamCost(); !ok {
		v := allocationcost.DefaultRamCost
		acc.mutation.SetRamCost(v)
	}
	if _, ok := acc.mutation.RamByteRequest(); !ok {
		v := allocationcost.DefaultRamByteRequest
		acc.mutation.SetRamByteRequest(v)
	}
	if _, ok := acc.mutation.PvCost(); !ok {
		v := allocationcost.DefaultPvCost
		acc.mutation.SetPvCost(v)
	}
	if _, ok := acc.mutation.PvBytes(); !ok {
		v := allocationcost.DefaultPvBytes
		acc.mutation.SetPvBytes(v)
	}
	if _, ok := acc.mutation.CpuCoreUsageAverage(); !ok {
		v := allocationcost.DefaultCpuCoreUsageAverage
		acc.mutation.SetCpuCoreUsageAverage(v)
	}
	if _, ok := acc.mutation.CpuCoreUsageMax(); !ok {
		v := allocationcost.DefaultCpuCoreUsageMax
		acc.mutation.SetCpuCoreUsageMax(v)
	}
	if _, ok := acc.mutation.RamByteUsageAverage(); !ok {
		v := allocationcost.DefaultRamByteUsageAverage
		acc.mutation.SetRamByteUsageAverage(v)
	}
	if _, ok := acc.mutation.RamByteUsageMax(); !ok {
		v := allocationcost.DefaultRamByteUsageMax
		acc.mutation.SetRamByteUsageMax(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (acc *AllocationCostCreate) check() error {
	if _, ok := acc.mutation.StartTime(); !ok {
		return &ValidationError{Name: "startTime", err: errors.New(`model: missing required field "AllocationCost.startTime"`)}
	}
	if _, ok := acc.mutation.EndTime(); !ok {
		return &ValidationError{Name: "endTime", err: errors.New(`model: missing required field "AllocationCost.endTime"`)}
	}
	if _, ok := acc.mutation.Minutes(); !ok {
		return &ValidationError{Name: "minutes", err: errors.New(`model: missing required field "AllocationCost.minutes"`)}
	}
	if _, ok := acc.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connectorID", err: errors.New(`model: missing required field "AllocationCost.connectorID"`)}
	}
	if v, ok := acc.mutation.ConnectorID(); ok {
		if err := allocationcost.ConnectorIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "connectorID", err: fmt.Errorf(`model: validator failed for field "AllocationCost.connectorID": %w`, err)}
		}
	}
	if _, ok := acc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "AllocationCost.name"`)}
	}
	if _, ok := acc.mutation.Fingerprint(); !ok {
		return &ValidationError{Name: "fingerprint", err: errors.New(`model: missing required field "AllocationCost.fingerprint"`)}
	}
	if _, ok := acc.mutation.Pvs(); !ok {
		return &ValidationError{Name: "pvs", err: errors.New(`model: missing required field "AllocationCost.pvs"`)}
	}
	if _, ok := acc.mutation.Labels(); !ok {
		return &ValidationError{Name: "labels", err: errors.New(`model: missing required field "AllocationCost.labels"`)}
	}
	if _, ok := acc.mutation.TotalCost(); !ok {
		return &ValidationError{Name: "totalCost", err: errors.New(`model: missing required field "AllocationCost.totalCost"`)}
	}
	if v, ok := acc.mutation.TotalCost(); ok {
		if err := allocationcost.TotalCostValidator(v); err != nil {
			return &ValidationError{Name: "totalCost", err: fmt.Errorf(`model: validator failed for field "AllocationCost.totalCost": %w`, err)}
		}
	}
	if _, ok := acc.mutation.CpuCost(); !ok {
		return &ValidationError{Name: "cpuCost", err: errors.New(`model: missing required field "AllocationCost.cpuCost"`)}
	}
	if v, ok := acc.mutation.CpuCost(); ok {
		if err := allocationcost.CpuCostValidator(v); err != nil {
			return &ValidationError{Name: "cpuCost", err: fmt.Errorf(`model: validator failed for field "AllocationCost.cpuCost": %w`, err)}
		}
	}
	if _, ok := acc.mutation.CpuCoreRequest(); !ok {
		return &ValidationError{Name: "cpuCoreRequest", err: errors.New(`model: missing required field "AllocationCost.cpuCoreRequest"`)}
	}
	if v, ok := acc.mutation.CpuCoreRequest(); ok {
		if err := allocationcost.CpuCoreRequestValidator(v); err != nil {
			return &ValidationError{Name: "cpuCoreRequest", err: fmt.Errorf(`model: validator failed for field "AllocationCost.cpuCoreRequest": %w`, err)}
		}
	}
	if _, ok := acc.mutation.GpuCost(); !ok {
		return &ValidationError{Name: "gpuCost", err: errors.New(`model: missing required field "AllocationCost.gpuCost"`)}
	}
	if v, ok := acc.mutation.GpuCost(); ok {
		if err := allocationcost.GpuCostValidator(v); err != nil {
			return &ValidationError{Name: "gpuCost", err: fmt.Errorf(`model: validator failed for field "AllocationCost.gpuCost": %w`, err)}
		}
	}
	if _, ok := acc.mutation.GpuCount(); !ok {
		return &ValidationError{Name: "gpuCount", err: errors.New(`model: missing required field "AllocationCost.gpuCount"`)}
	}
	if v, ok := acc.mutation.GpuCount(); ok {
		if err := allocationcost.GpuCountValidator(v); err != nil {
			return &ValidationError{Name: "gpuCount", err: fmt.Errorf(`model: validator failed for field "AllocationCost.gpuCount": %w`, err)}
		}
	}
	if _, ok := acc.mutation.RamCost(); !ok {
		return &ValidationError{Name: "ramCost", err: errors.New(`model: missing required field "AllocationCost.ramCost"`)}
	}
	if v, ok := acc.mutation.RamCost(); ok {
		if err := allocationcost.RamCostValidator(v); err != nil {
			return &ValidationError{Name: "ramCost", err: fmt.Errorf(`model: validator failed for field "AllocationCost.ramCost": %w`, err)}
		}
	}
	if _, ok := acc.mutation.RamByteRequest(); !ok {
		return &ValidationError{Name: "ramByteRequest", err: errors.New(`model: missing required field "AllocationCost.ramByteRequest"`)}
	}
	if v, ok := acc.mutation.RamByteRequest(); ok {
		if err := allocationcost.RamByteRequestValidator(v); err != nil {
			return &ValidationError{Name: "ramByteRequest", err: fmt.Errorf(`model: validator failed for field "AllocationCost.ramByteRequest": %w`, err)}
		}
	}
	if _, ok := acc.mutation.PvCost(); !ok {
		return &ValidationError{Name: "pvCost", err: errors.New(`model: missing required field "AllocationCost.pvCost"`)}
	}
	if v, ok := acc.mutation.PvCost(); ok {
		if err := allocationcost.PvCostValidator(v); err != nil {
			return &ValidationError{Name: "pvCost", err: fmt.Errorf(`model: validator failed for field "AllocationCost.pvCost": %w`, err)}
		}
	}
	if _, ok := acc.mutation.PvBytes(); !ok {
		return &ValidationError{Name: "pvBytes", err: errors.New(`model: missing required field "AllocationCost.pvBytes"`)}
	}
	if v, ok := acc.mutation.PvBytes(); ok {
		if err := allocationcost.PvBytesValidator(v); err != nil {
			return &ValidationError{Name: "pvBytes", err: fmt.Errorf(`model: validator failed for field "AllocationCost.pvBytes": %w`, err)}
		}
	}
	if _, ok := acc.mutation.CpuCoreUsageAverage(); !ok {
		return &ValidationError{Name: "cpuCoreUsageAverage", err: errors.New(`model: missing required field "AllocationCost.cpuCoreUsageAverage"`)}
	}
	if v, ok := acc.mutation.CpuCoreUsageAverage(); ok {
		if err := allocationcost.CpuCoreUsageAverageValidator(v); err != nil {
			return &ValidationError{Name: "cpuCoreUsageAverage", err: fmt.Errorf(`model: validator failed for field "AllocationCost.cpuCoreUsageAverage": %w`, err)}
		}
	}
	if _, ok := acc.mutation.CpuCoreUsageMax(); !ok {
		return &ValidationError{Name: "cpuCoreUsageMax", err: errors.New(`model: missing required field "AllocationCost.cpuCoreUsageMax"`)}
	}
	if v, ok := acc.mutation.CpuCoreUsageMax(); ok {
		if err := allocationcost.CpuCoreUsageMaxValidator(v); err != nil {
			return &ValidationError{Name: "cpuCoreUsageMax", err: fmt.Errorf(`model: validator failed for field "AllocationCost.cpuCoreUsageMax": %w`, err)}
		}
	}
	if _, ok := acc.mutation.RamByteUsageAverage(); !ok {
		return &ValidationError{Name: "ramByteUsageAverage", err: errors.New(`model: missing required field "AllocationCost.ramByteUsageAverage"`)}
	}
	if v, ok := acc.mutation.RamByteUsageAverage(); ok {
		if err := allocationcost.RamByteUsageAverageValidator(v); err != nil {
			return &ValidationError{Name: "ramByteUsageAverage", err: fmt.Errorf(`model: validator failed for field "AllocationCost.ramByteUsageAverage": %w`, err)}
		}
	}
	if _, ok := acc.mutation.RamByteUsageMax(); !ok {
		return &ValidationError{Name: "ramByteUsageMax", err: errors.New(`model: missing required field "AllocationCost.ramByteUsageMax"`)}
	}
	if v, ok := acc.mutation.RamByteUsageMax(); ok {
		if err := allocationcost.RamByteUsageMaxValidator(v); err != nil {
			return &ValidationError{Name: "ramByteUsageMax", err: fmt.Errorf(`model: validator failed for field "AllocationCost.ramByteUsageMax": %w`, err)}
		}
	}
	if _, ok := acc.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connector", err: errors.New(`model: missing required edge "AllocationCost.connector"`)}
	}
	return nil
}

func (acc *AllocationCostCreate) sqlSave(ctx context.Context) (*AllocationCost, error) {
	if err := acc.check(); err != nil {
		return nil, err
	}
	_node, _spec := acc.createSpec()
	if err := sqlgraph.CreateNode(ctx, acc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	acc.mutation.id = &_node.ID
	acc.mutation.done = true
	return _node, nil
}

func (acc *AllocationCostCreate) createSpec() (*AllocationCost, *sqlgraph.CreateSpec) {
	var (
		_node = &AllocationCost{config: acc.config}
		_spec = sqlgraph.NewCreateSpec(allocationcost.Table, sqlgraph.NewFieldSpec(allocationcost.FieldID, field.TypeInt))
	)
	_spec.Schema = acc.schemaConfig.AllocationCost
	_spec.OnConflict = acc.conflict
	if value, ok := acc.mutation.StartTime(); ok {
		_spec.SetField(allocationcost.FieldStartTime, field.TypeTime, value)
		_node.StartTime = value
	}
	if value, ok := acc.mutation.EndTime(); ok {
		_spec.SetField(allocationcost.FieldEndTime, field.TypeTime, value)
		_node.EndTime = value
	}
	if value, ok := acc.mutation.Minutes(); ok {
		_spec.SetField(allocationcost.FieldMinutes, field.TypeFloat64, value)
		_node.Minutes = value
	}
	if value, ok := acc.mutation.Name(); ok {
		_spec.SetField(allocationcost.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := acc.mutation.Fingerprint(); ok {
		_spec.SetField(allocationcost.FieldFingerprint, field.TypeString, value)
		_node.Fingerprint = value
	}
	if value, ok := acc.mutation.ClusterName(); ok {
		_spec.SetField(allocationcost.FieldClusterName, field.TypeString, value)
		_node.ClusterName = value
	}
	if value, ok := acc.mutation.Namespace(); ok {
		_spec.SetField(allocationcost.FieldNamespace, field.TypeString, value)
		_node.Namespace = value
	}
	if value, ok := acc.mutation.Node(); ok {
		_spec.SetField(allocationcost.FieldNode, field.TypeString, value)
		_node.Node = value
	}
	if value, ok := acc.mutation.Controller(); ok {
		_spec.SetField(allocationcost.FieldController, field.TypeString, value)
		_node.Controller = value
	}
	if value, ok := acc.mutation.ControllerKind(); ok {
		_spec.SetField(allocationcost.FieldControllerKind, field.TypeString, value)
		_node.ControllerKind = value
	}
	if value, ok := acc.mutation.Pod(); ok {
		_spec.SetField(allocationcost.FieldPod, field.TypeString, value)
		_node.Pod = value
	}
	if value, ok := acc.mutation.Container(); ok {
		_spec.SetField(allocationcost.FieldContainer, field.TypeString, value)
		_node.Container = value
	}
	if value, ok := acc.mutation.Pvs(); ok {
		_spec.SetField(allocationcost.FieldPvs, field.TypeJSON, value)
		_node.Pvs = value
	}
	if value, ok := acc.mutation.Labels(); ok {
		_spec.SetField(allocationcost.FieldLabels, field.TypeJSON, value)
		_node.Labels = value
	}
	if value, ok := acc.mutation.TotalCost(); ok {
		_spec.SetField(allocationcost.FieldTotalCost, field.TypeFloat64, value)
		_node.TotalCost = value
	}
	if value, ok := acc.mutation.Currency(); ok {
		_spec.SetField(allocationcost.FieldCurrency, field.TypeInt, value)
		_node.Currency = value
	}
	if value, ok := acc.mutation.CpuCost(); ok {
		_spec.SetField(allocationcost.FieldCpuCost, field.TypeFloat64, value)
		_node.CpuCost = value
	}
	if value, ok := acc.mutation.CpuCoreRequest(); ok {
		_spec.SetField(allocationcost.FieldCpuCoreRequest, field.TypeFloat64, value)
		_node.CpuCoreRequest = value
	}
	if value, ok := acc.mutation.GpuCost(); ok {
		_spec.SetField(allocationcost.FieldGpuCost, field.TypeFloat64, value)
		_node.GpuCost = value
	}
	if value, ok := acc.mutation.GpuCount(); ok {
		_spec.SetField(allocationcost.FieldGpuCount, field.TypeFloat64, value)
		_node.GpuCount = value
	}
	if value, ok := acc.mutation.RamCost(); ok {
		_spec.SetField(allocationcost.FieldRamCost, field.TypeFloat64, value)
		_node.RamCost = value
	}
	if value, ok := acc.mutation.RamByteRequest(); ok {
		_spec.SetField(allocationcost.FieldRamByteRequest, field.TypeFloat64, value)
		_node.RamByteRequest = value
	}
	if value, ok := acc.mutation.PvCost(); ok {
		_spec.SetField(allocationcost.FieldPvCost, field.TypeFloat64, value)
		_node.PvCost = value
	}
	if value, ok := acc.mutation.PvBytes(); ok {
		_spec.SetField(allocationcost.FieldPvBytes, field.TypeFloat64, value)
		_node.PvBytes = value
	}
	if value, ok := acc.mutation.CpuCoreUsageAverage(); ok {
		_spec.SetField(allocationcost.FieldCpuCoreUsageAverage, field.TypeFloat64, value)
		_node.CpuCoreUsageAverage = value
	}
	if value, ok := acc.mutation.CpuCoreUsageMax(); ok {
		_spec.SetField(allocationcost.FieldCpuCoreUsageMax, field.TypeFloat64, value)
		_node.CpuCoreUsageMax = value
	}
	if value, ok := acc.mutation.RamByteUsageAverage(); ok {
		_spec.SetField(allocationcost.FieldRamByteUsageAverage, field.TypeFloat64, value)
		_node.RamByteUsageAverage = value
	}
	if value, ok := acc.mutation.RamByteUsageMax(); ok {
		_spec.SetField(allocationcost.FieldRamByteUsageMax, field.TypeFloat64, value)
		_node.RamByteUsageMax = value
	}
	if nodes := acc.mutation.ConnectorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   allocationcost.ConnectorTable,
			Columns: []string{allocationcost.ConnectorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = acc.schemaConfig.AllocationCost
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
//	client.AllocationCost.Create().
//		SetStartTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AllocationCostUpsert) {
//			SetStartTime(v+v).
//		}).
//		Exec(ctx)
func (acc *AllocationCostCreate) OnConflict(opts ...sql.ConflictOption) *AllocationCostUpsertOne {
	acc.conflict = opts
	return &AllocationCostUpsertOne{
		create: acc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AllocationCost.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (acc *AllocationCostCreate) OnConflictColumns(columns ...string) *AllocationCostUpsertOne {
	acc.conflict = append(acc.conflict, sql.ConflictColumns(columns...))
	return &AllocationCostUpsertOne{
		create: acc,
	}
}

type (
	// AllocationCostUpsertOne is the builder for "upsert"-ing
	//  one AllocationCost node.
	AllocationCostUpsertOne struct {
		create *AllocationCostCreate
	}

	// AllocationCostUpsert is the "OnConflict" setter.
	AllocationCostUpsert struct {
		*sql.UpdateSet
	}
)

// SetTotalCost sets the "totalCost" field.
func (u *AllocationCostUpsert) SetTotalCost(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldTotalCost, v)
	return u
}

// UpdateTotalCost sets the "totalCost" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdateTotalCost() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldTotalCost)
	return u
}

// AddTotalCost adds v to the "totalCost" field.
func (u *AllocationCostUpsert) AddTotalCost(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldTotalCost, v)
	return u
}

// SetCurrency sets the "currency" field.
func (u *AllocationCostUpsert) SetCurrency(v int) *AllocationCostUpsert {
	u.Set(allocationcost.FieldCurrency, v)
	return u
}

// UpdateCurrency sets the "currency" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdateCurrency() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldCurrency)
	return u
}

// AddCurrency adds v to the "currency" field.
func (u *AllocationCostUpsert) AddCurrency(v int) *AllocationCostUpsert {
	u.Add(allocationcost.FieldCurrency, v)
	return u
}

// ClearCurrency clears the value of the "currency" field.
func (u *AllocationCostUpsert) ClearCurrency() *AllocationCostUpsert {
	u.SetNull(allocationcost.FieldCurrency)
	return u
}

// SetCpuCost sets the "cpuCost" field.
func (u *AllocationCostUpsert) SetCpuCost(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldCpuCost, v)
	return u
}

// UpdateCpuCost sets the "cpuCost" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdateCpuCost() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldCpuCost)
	return u
}

// AddCpuCost adds v to the "cpuCost" field.
func (u *AllocationCostUpsert) AddCpuCost(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldCpuCost, v)
	return u
}

// SetGpuCost sets the "gpuCost" field.
func (u *AllocationCostUpsert) SetGpuCost(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldGpuCost, v)
	return u
}

// UpdateGpuCost sets the "gpuCost" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdateGpuCost() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldGpuCost)
	return u
}

// AddGpuCost adds v to the "gpuCost" field.
func (u *AllocationCostUpsert) AddGpuCost(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldGpuCost, v)
	return u
}

// SetRamCost sets the "ramCost" field.
func (u *AllocationCostUpsert) SetRamCost(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldRamCost, v)
	return u
}

// UpdateRamCost sets the "ramCost" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdateRamCost() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldRamCost)
	return u
}

// AddRamCost adds v to the "ramCost" field.
func (u *AllocationCostUpsert) AddRamCost(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldRamCost, v)
	return u
}

// SetPvCost sets the "pvCost" field.
func (u *AllocationCostUpsert) SetPvCost(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldPvCost, v)
	return u
}

// UpdatePvCost sets the "pvCost" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdatePvCost() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldPvCost)
	return u
}

// AddPvCost adds v to the "pvCost" field.
func (u *AllocationCostUpsert) AddPvCost(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldPvCost, v)
	return u
}

// SetPvBytes sets the "pvBytes" field.
func (u *AllocationCostUpsert) SetPvBytes(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldPvBytes, v)
	return u
}

// UpdatePvBytes sets the "pvBytes" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdatePvBytes() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldPvBytes)
	return u
}

// AddPvBytes adds v to the "pvBytes" field.
func (u *AllocationCostUpsert) AddPvBytes(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldPvBytes, v)
	return u
}

// SetCpuCoreUsageAverage sets the "cpuCoreUsageAverage" field.
func (u *AllocationCostUpsert) SetCpuCoreUsageAverage(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldCpuCoreUsageAverage, v)
	return u
}

// UpdateCpuCoreUsageAverage sets the "cpuCoreUsageAverage" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdateCpuCoreUsageAverage() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldCpuCoreUsageAverage)
	return u
}

// AddCpuCoreUsageAverage adds v to the "cpuCoreUsageAverage" field.
func (u *AllocationCostUpsert) AddCpuCoreUsageAverage(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldCpuCoreUsageAverage, v)
	return u
}

// SetCpuCoreUsageMax sets the "cpuCoreUsageMax" field.
func (u *AllocationCostUpsert) SetCpuCoreUsageMax(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldCpuCoreUsageMax, v)
	return u
}

// UpdateCpuCoreUsageMax sets the "cpuCoreUsageMax" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdateCpuCoreUsageMax() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldCpuCoreUsageMax)
	return u
}

// AddCpuCoreUsageMax adds v to the "cpuCoreUsageMax" field.
func (u *AllocationCostUpsert) AddCpuCoreUsageMax(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldCpuCoreUsageMax, v)
	return u
}

// SetRamByteUsageAverage sets the "ramByteUsageAverage" field.
func (u *AllocationCostUpsert) SetRamByteUsageAverage(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldRamByteUsageAverage, v)
	return u
}

// UpdateRamByteUsageAverage sets the "ramByteUsageAverage" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdateRamByteUsageAverage() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldRamByteUsageAverage)
	return u
}

// AddRamByteUsageAverage adds v to the "ramByteUsageAverage" field.
func (u *AllocationCostUpsert) AddRamByteUsageAverage(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldRamByteUsageAverage, v)
	return u
}

// SetRamByteUsageMax sets the "ramByteUsageMax" field.
func (u *AllocationCostUpsert) SetRamByteUsageMax(v float64) *AllocationCostUpsert {
	u.Set(allocationcost.FieldRamByteUsageMax, v)
	return u
}

// UpdateRamByteUsageMax sets the "ramByteUsageMax" field to the value that was provided on create.
func (u *AllocationCostUpsert) UpdateRamByteUsageMax() *AllocationCostUpsert {
	u.SetExcluded(allocationcost.FieldRamByteUsageMax)
	return u
}

// AddRamByteUsageMax adds v to the "ramByteUsageMax" field.
func (u *AllocationCostUpsert) AddRamByteUsageMax(v float64) *AllocationCostUpsert {
	u.Add(allocationcost.FieldRamByteUsageMax, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.AllocationCost.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AllocationCostUpsertOne) UpdateNewValues() *AllocationCostUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.StartTime(); exists {
			s.SetIgnore(allocationcost.FieldStartTime)
		}
		if _, exists := u.create.mutation.EndTime(); exists {
			s.SetIgnore(allocationcost.FieldEndTime)
		}
		if _, exists := u.create.mutation.Minutes(); exists {
			s.SetIgnore(allocationcost.FieldMinutes)
		}
		if _, exists := u.create.mutation.ConnectorID(); exists {
			s.SetIgnore(allocationcost.FieldConnectorID)
		}
		if _, exists := u.create.mutation.Name(); exists {
			s.SetIgnore(allocationcost.FieldName)
		}
		if _, exists := u.create.mutation.Fingerprint(); exists {
			s.SetIgnore(allocationcost.FieldFingerprint)
		}
		if _, exists := u.create.mutation.ClusterName(); exists {
			s.SetIgnore(allocationcost.FieldClusterName)
		}
		if _, exists := u.create.mutation.Namespace(); exists {
			s.SetIgnore(allocationcost.FieldNamespace)
		}
		if _, exists := u.create.mutation.Node(); exists {
			s.SetIgnore(allocationcost.FieldNode)
		}
		if _, exists := u.create.mutation.Controller(); exists {
			s.SetIgnore(allocationcost.FieldController)
		}
		if _, exists := u.create.mutation.ControllerKind(); exists {
			s.SetIgnore(allocationcost.FieldControllerKind)
		}
		if _, exists := u.create.mutation.Pod(); exists {
			s.SetIgnore(allocationcost.FieldPod)
		}
		if _, exists := u.create.mutation.Container(); exists {
			s.SetIgnore(allocationcost.FieldContainer)
		}
		if _, exists := u.create.mutation.Pvs(); exists {
			s.SetIgnore(allocationcost.FieldPvs)
		}
		if _, exists := u.create.mutation.Labels(); exists {
			s.SetIgnore(allocationcost.FieldLabels)
		}
		if _, exists := u.create.mutation.CpuCoreRequest(); exists {
			s.SetIgnore(allocationcost.FieldCpuCoreRequest)
		}
		if _, exists := u.create.mutation.GpuCount(); exists {
			s.SetIgnore(allocationcost.FieldGpuCount)
		}
		if _, exists := u.create.mutation.RamByteRequest(); exists {
			s.SetIgnore(allocationcost.FieldRamByteRequest)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.AllocationCost.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *AllocationCostUpsertOne) Ignore() *AllocationCostUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AllocationCostUpsertOne) DoNothing() *AllocationCostUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AllocationCostCreate.OnConflict
// documentation for more info.
func (u *AllocationCostUpsertOne) Update(set func(*AllocationCostUpsert)) *AllocationCostUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AllocationCostUpsert{UpdateSet: update})
	}))
	return u
}

// SetTotalCost sets the "totalCost" field.
func (u *AllocationCostUpsertOne) SetTotalCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetTotalCost(v)
	})
}

// AddTotalCost adds v to the "totalCost" field.
func (u *AllocationCostUpsertOne) AddTotalCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddTotalCost(v)
	})
}

// UpdateTotalCost sets the "totalCost" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdateTotalCost() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateTotalCost()
	})
}

// SetCurrency sets the "currency" field.
func (u *AllocationCostUpsertOne) SetCurrency(v int) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetCurrency(v)
	})
}

// AddCurrency adds v to the "currency" field.
func (u *AllocationCostUpsertOne) AddCurrency(v int) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddCurrency(v)
	})
}

// UpdateCurrency sets the "currency" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdateCurrency() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateCurrency()
	})
}

// ClearCurrency clears the value of the "currency" field.
func (u *AllocationCostUpsertOne) ClearCurrency() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.ClearCurrency()
	})
}

// SetCpuCost sets the "cpuCost" field.
func (u *AllocationCostUpsertOne) SetCpuCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetCpuCost(v)
	})
}

// AddCpuCost adds v to the "cpuCost" field.
func (u *AllocationCostUpsertOne) AddCpuCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddCpuCost(v)
	})
}

// UpdateCpuCost sets the "cpuCost" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdateCpuCost() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateCpuCost()
	})
}

// SetGpuCost sets the "gpuCost" field.
func (u *AllocationCostUpsertOne) SetGpuCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetGpuCost(v)
	})
}

// AddGpuCost adds v to the "gpuCost" field.
func (u *AllocationCostUpsertOne) AddGpuCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddGpuCost(v)
	})
}

// UpdateGpuCost sets the "gpuCost" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdateGpuCost() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateGpuCost()
	})
}

// SetRamCost sets the "ramCost" field.
func (u *AllocationCostUpsertOne) SetRamCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetRamCost(v)
	})
}

// AddRamCost adds v to the "ramCost" field.
func (u *AllocationCostUpsertOne) AddRamCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddRamCost(v)
	})
}

// UpdateRamCost sets the "ramCost" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdateRamCost() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateRamCost()
	})
}

// SetPvCost sets the "pvCost" field.
func (u *AllocationCostUpsertOne) SetPvCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetPvCost(v)
	})
}

// AddPvCost adds v to the "pvCost" field.
func (u *AllocationCostUpsertOne) AddPvCost(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddPvCost(v)
	})
}

// UpdatePvCost sets the "pvCost" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdatePvCost() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdatePvCost()
	})
}

// SetPvBytes sets the "pvBytes" field.
func (u *AllocationCostUpsertOne) SetPvBytes(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetPvBytes(v)
	})
}

// AddPvBytes adds v to the "pvBytes" field.
func (u *AllocationCostUpsertOne) AddPvBytes(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddPvBytes(v)
	})
}

// UpdatePvBytes sets the "pvBytes" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdatePvBytes() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdatePvBytes()
	})
}

// SetCpuCoreUsageAverage sets the "cpuCoreUsageAverage" field.
func (u *AllocationCostUpsertOne) SetCpuCoreUsageAverage(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetCpuCoreUsageAverage(v)
	})
}

// AddCpuCoreUsageAverage adds v to the "cpuCoreUsageAverage" field.
func (u *AllocationCostUpsertOne) AddCpuCoreUsageAverage(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddCpuCoreUsageAverage(v)
	})
}

// UpdateCpuCoreUsageAverage sets the "cpuCoreUsageAverage" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdateCpuCoreUsageAverage() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateCpuCoreUsageAverage()
	})
}

// SetCpuCoreUsageMax sets the "cpuCoreUsageMax" field.
func (u *AllocationCostUpsertOne) SetCpuCoreUsageMax(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetCpuCoreUsageMax(v)
	})
}

// AddCpuCoreUsageMax adds v to the "cpuCoreUsageMax" field.
func (u *AllocationCostUpsertOne) AddCpuCoreUsageMax(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddCpuCoreUsageMax(v)
	})
}

// UpdateCpuCoreUsageMax sets the "cpuCoreUsageMax" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdateCpuCoreUsageMax() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateCpuCoreUsageMax()
	})
}

// SetRamByteUsageAverage sets the "ramByteUsageAverage" field.
func (u *AllocationCostUpsertOne) SetRamByteUsageAverage(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetRamByteUsageAverage(v)
	})
}

// AddRamByteUsageAverage adds v to the "ramByteUsageAverage" field.
func (u *AllocationCostUpsertOne) AddRamByteUsageAverage(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddRamByteUsageAverage(v)
	})
}

// UpdateRamByteUsageAverage sets the "ramByteUsageAverage" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdateRamByteUsageAverage() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateRamByteUsageAverage()
	})
}

// SetRamByteUsageMax sets the "ramByteUsageMax" field.
func (u *AllocationCostUpsertOne) SetRamByteUsageMax(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetRamByteUsageMax(v)
	})
}

// AddRamByteUsageMax adds v to the "ramByteUsageMax" field.
func (u *AllocationCostUpsertOne) AddRamByteUsageMax(v float64) *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddRamByteUsageMax(v)
	})
}

// UpdateRamByteUsageMax sets the "ramByteUsageMax" field to the value that was provided on create.
func (u *AllocationCostUpsertOne) UpdateRamByteUsageMax() *AllocationCostUpsertOne {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateRamByteUsageMax()
	})
}

// Exec executes the query.
func (u *AllocationCostUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for AllocationCostCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AllocationCostUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *AllocationCostUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *AllocationCostUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// AllocationCostCreateBulk is the builder for creating many AllocationCost entities in bulk.
type AllocationCostCreateBulk struct {
	config
	builders []*AllocationCostCreate
	conflict []sql.ConflictOption
}

// Save creates the AllocationCost entities in the database.
func (accb *AllocationCostCreateBulk) Save(ctx context.Context) ([]*AllocationCost, error) {
	specs := make([]*sqlgraph.CreateSpec, len(accb.builders))
	nodes := make([]*AllocationCost, len(accb.builders))
	mutators := make([]Mutator, len(accb.builders))
	for i := range accb.builders {
		func(i int, root context.Context) {
			builder := accb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AllocationCostMutation)
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
					_, err = mutators[i+1].Mutate(root, accb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = accb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, accb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, accb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (accb *AllocationCostCreateBulk) SaveX(ctx context.Context) []*AllocationCost {
	v, err := accb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (accb *AllocationCostCreateBulk) Exec(ctx context.Context) error {
	_, err := accb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (accb *AllocationCostCreateBulk) ExecX(ctx context.Context) {
	if err := accb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.AllocationCost.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.AllocationCostUpsert) {
//			SetStartTime(v+v).
//		}).
//		Exec(ctx)
func (accb *AllocationCostCreateBulk) OnConflict(opts ...sql.ConflictOption) *AllocationCostUpsertBulk {
	accb.conflict = opts
	return &AllocationCostUpsertBulk{
		create: accb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.AllocationCost.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (accb *AllocationCostCreateBulk) OnConflictColumns(columns ...string) *AllocationCostUpsertBulk {
	accb.conflict = append(accb.conflict, sql.ConflictColumns(columns...))
	return &AllocationCostUpsertBulk{
		create: accb,
	}
}

// AllocationCostUpsertBulk is the builder for "upsert"-ing
// a bulk of AllocationCost nodes.
type AllocationCostUpsertBulk struct {
	create *AllocationCostCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.AllocationCost.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *AllocationCostUpsertBulk) UpdateNewValues() *AllocationCostUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.StartTime(); exists {
				s.SetIgnore(allocationcost.FieldStartTime)
			}
			if _, exists := b.mutation.EndTime(); exists {
				s.SetIgnore(allocationcost.FieldEndTime)
			}
			if _, exists := b.mutation.Minutes(); exists {
				s.SetIgnore(allocationcost.FieldMinutes)
			}
			if _, exists := b.mutation.ConnectorID(); exists {
				s.SetIgnore(allocationcost.FieldConnectorID)
			}
			if _, exists := b.mutation.Name(); exists {
				s.SetIgnore(allocationcost.FieldName)
			}
			if _, exists := b.mutation.Fingerprint(); exists {
				s.SetIgnore(allocationcost.FieldFingerprint)
			}
			if _, exists := b.mutation.ClusterName(); exists {
				s.SetIgnore(allocationcost.FieldClusterName)
			}
			if _, exists := b.mutation.Namespace(); exists {
				s.SetIgnore(allocationcost.FieldNamespace)
			}
			if _, exists := b.mutation.Node(); exists {
				s.SetIgnore(allocationcost.FieldNode)
			}
			if _, exists := b.mutation.Controller(); exists {
				s.SetIgnore(allocationcost.FieldController)
			}
			if _, exists := b.mutation.ControllerKind(); exists {
				s.SetIgnore(allocationcost.FieldControllerKind)
			}
			if _, exists := b.mutation.Pod(); exists {
				s.SetIgnore(allocationcost.FieldPod)
			}
			if _, exists := b.mutation.Container(); exists {
				s.SetIgnore(allocationcost.FieldContainer)
			}
			if _, exists := b.mutation.Pvs(); exists {
				s.SetIgnore(allocationcost.FieldPvs)
			}
			if _, exists := b.mutation.Labels(); exists {
				s.SetIgnore(allocationcost.FieldLabels)
			}
			if _, exists := b.mutation.CpuCoreRequest(); exists {
				s.SetIgnore(allocationcost.FieldCpuCoreRequest)
			}
			if _, exists := b.mutation.GpuCount(); exists {
				s.SetIgnore(allocationcost.FieldGpuCount)
			}
			if _, exists := b.mutation.RamByteRequest(); exists {
				s.SetIgnore(allocationcost.FieldRamByteRequest)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.AllocationCost.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *AllocationCostUpsertBulk) Ignore() *AllocationCostUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *AllocationCostUpsertBulk) DoNothing() *AllocationCostUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the AllocationCostCreateBulk.OnConflict
// documentation for more info.
func (u *AllocationCostUpsertBulk) Update(set func(*AllocationCostUpsert)) *AllocationCostUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&AllocationCostUpsert{UpdateSet: update})
	}))
	return u
}

// SetTotalCost sets the "totalCost" field.
func (u *AllocationCostUpsertBulk) SetTotalCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetTotalCost(v)
	})
}

// AddTotalCost adds v to the "totalCost" field.
func (u *AllocationCostUpsertBulk) AddTotalCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddTotalCost(v)
	})
}

// UpdateTotalCost sets the "totalCost" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdateTotalCost() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateTotalCost()
	})
}

// SetCurrency sets the "currency" field.
func (u *AllocationCostUpsertBulk) SetCurrency(v int) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetCurrency(v)
	})
}

// AddCurrency adds v to the "currency" field.
func (u *AllocationCostUpsertBulk) AddCurrency(v int) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddCurrency(v)
	})
}

// UpdateCurrency sets the "currency" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdateCurrency() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateCurrency()
	})
}

// ClearCurrency clears the value of the "currency" field.
func (u *AllocationCostUpsertBulk) ClearCurrency() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.ClearCurrency()
	})
}

// SetCpuCost sets the "cpuCost" field.
func (u *AllocationCostUpsertBulk) SetCpuCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetCpuCost(v)
	})
}

// AddCpuCost adds v to the "cpuCost" field.
func (u *AllocationCostUpsertBulk) AddCpuCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddCpuCost(v)
	})
}

// UpdateCpuCost sets the "cpuCost" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdateCpuCost() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateCpuCost()
	})
}

// SetGpuCost sets the "gpuCost" field.
func (u *AllocationCostUpsertBulk) SetGpuCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetGpuCost(v)
	})
}

// AddGpuCost adds v to the "gpuCost" field.
func (u *AllocationCostUpsertBulk) AddGpuCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddGpuCost(v)
	})
}

// UpdateGpuCost sets the "gpuCost" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdateGpuCost() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateGpuCost()
	})
}

// SetRamCost sets the "ramCost" field.
func (u *AllocationCostUpsertBulk) SetRamCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetRamCost(v)
	})
}

// AddRamCost adds v to the "ramCost" field.
func (u *AllocationCostUpsertBulk) AddRamCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddRamCost(v)
	})
}

// UpdateRamCost sets the "ramCost" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdateRamCost() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateRamCost()
	})
}

// SetPvCost sets the "pvCost" field.
func (u *AllocationCostUpsertBulk) SetPvCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetPvCost(v)
	})
}

// AddPvCost adds v to the "pvCost" field.
func (u *AllocationCostUpsertBulk) AddPvCost(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddPvCost(v)
	})
}

// UpdatePvCost sets the "pvCost" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdatePvCost() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdatePvCost()
	})
}

// SetPvBytes sets the "pvBytes" field.
func (u *AllocationCostUpsertBulk) SetPvBytes(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetPvBytes(v)
	})
}

// AddPvBytes adds v to the "pvBytes" field.
func (u *AllocationCostUpsertBulk) AddPvBytes(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddPvBytes(v)
	})
}

// UpdatePvBytes sets the "pvBytes" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdatePvBytes() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdatePvBytes()
	})
}

// SetCpuCoreUsageAverage sets the "cpuCoreUsageAverage" field.
func (u *AllocationCostUpsertBulk) SetCpuCoreUsageAverage(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetCpuCoreUsageAverage(v)
	})
}

// AddCpuCoreUsageAverage adds v to the "cpuCoreUsageAverage" field.
func (u *AllocationCostUpsertBulk) AddCpuCoreUsageAverage(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddCpuCoreUsageAverage(v)
	})
}

// UpdateCpuCoreUsageAverage sets the "cpuCoreUsageAverage" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdateCpuCoreUsageAverage() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateCpuCoreUsageAverage()
	})
}

// SetCpuCoreUsageMax sets the "cpuCoreUsageMax" field.
func (u *AllocationCostUpsertBulk) SetCpuCoreUsageMax(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetCpuCoreUsageMax(v)
	})
}

// AddCpuCoreUsageMax adds v to the "cpuCoreUsageMax" field.
func (u *AllocationCostUpsertBulk) AddCpuCoreUsageMax(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddCpuCoreUsageMax(v)
	})
}

// UpdateCpuCoreUsageMax sets the "cpuCoreUsageMax" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdateCpuCoreUsageMax() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateCpuCoreUsageMax()
	})
}

// SetRamByteUsageAverage sets the "ramByteUsageAverage" field.
func (u *AllocationCostUpsertBulk) SetRamByteUsageAverage(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetRamByteUsageAverage(v)
	})
}

// AddRamByteUsageAverage adds v to the "ramByteUsageAverage" field.
func (u *AllocationCostUpsertBulk) AddRamByteUsageAverage(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddRamByteUsageAverage(v)
	})
}

// UpdateRamByteUsageAverage sets the "ramByteUsageAverage" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdateRamByteUsageAverage() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateRamByteUsageAverage()
	})
}

// SetRamByteUsageMax sets the "ramByteUsageMax" field.
func (u *AllocationCostUpsertBulk) SetRamByteUsageMax(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.SetRamByteUsageMax(v)
	})
}

// AddRamByteUsageMax adds v to the "ramByteUsageMax" field.
func (u *AllocationCostUpsertBulk) AddRamByteUsageMax(v float64) *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.AddRamByteUsageMax(v)
	})
}

// UpdateRamByteUsageMax sets the "ramByteUsageMax" field to the value that was provided on create.
func (u *AllocationCostUpsertBulk) UpdateRamByteUsageMax() *AllocationCostUpsertBulk {
	return u.Update(func(s *AllocationCostUpsert) {
		s.UpdateRamByteUsageMax()
	})
}

// Exec executes the query.
func (u *AllocationCostUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the AllocationCostCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for AllocationCostCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *AllocationCostUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
