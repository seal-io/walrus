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
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
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

// SetName sets the "name" field.
func (cu *ConnectorUpdate) SetName(s string) *ConnectorUpdate {
	cu.mutation.SetName(s)
	return cu
}

// SetDescription sets the "description" field.
func (cu *ConnectorUpdate) SetDescription(s string) *ConnectorUpdate {
	cu.mutation.SetDescription(s)
	return cu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cu *ConnectorUpdate) SetNillableDescription(s *string) *ConnectorUpdate {
	if s != nil {
		cu.SetDescription(*s)
	}
	return cu
}

// ClearDescription clears the value of the "description" field.
func (cu *ConnectorUpdate) ClearDescription() *ConnectorUpdate {
	cu.mutation.ClearDescription()
	return cu
}

// SetLabels sets the "labels" field.
func (cu *ConnectorUpdate) SetLabels(m map[string]string) *ConnectorUpdate {
	cu.mutation.SetLabels(m)
	return cu
}

// SetUpdateTime sets the "updateTime" field.
func (cu *ConnectorUpdate) SetUpdateTime(t time.Time) *ConnectorUpdate {
	cu.mutation.SetUpdateTime(t)
	return cu
}

// SetStatus sets the "status" field.
func (cu *ConnectorUpdate) SetStatus(s status.Status) *ConnectorUpdate {
	cu.mutation.SetStatus(s)
	return cu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (cu *ConnectorUpdate) SetNillableStatus(s *status.Status) *ConnectorUpdate {
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

// SetConfigVersion sets the "configVersion" field.
func (cu *ConnectorUpdate) SetConfigVersion(s string) *ConnectorUpdate {
	cu.mutation.SetConfigVersion(s)
	return cu
}

// SetConfigData sets the "configData" field.
func (cu *ConnectorUpdate) SetConfigData(c crypto.Properties) *ConnectorUpdate {
	cu.mutation.SetConfigData(c)
	return cu
}

// SetEnableFinOps sets the "enableFinOps" field.
func (cu *ConnectorUpdate) SetEnableFinOps(b bool) *ConnectorUpdate {
	cu.mutation.SetEnableFinOps(b)
	return cu
}

// SetFinOpsCustomPricing sets the "finOpsCustomPricing" field.
func (cu *ConnectorUpdate) SetFinOpsCustomPricing(tocp *types.FinOpsCustomPricing) *ConnectorUpdate {
	cu.mutation.SetFinOpsCustomPricing(tocp)
	return cu
}

// ClearFinOpsCustomPricing clears the value of the "finOpsCustomPricing" field.
func (cu *ConnectorUpdate) ClearFinOpsCustomPricing() *ConnectorUpdate {
	cu.mutation.ClearFinOpsCustomPricing()
	return cu
}

// SetCategory sets the "category" field.
func (cu *ConnectorUpdate) SetCategory(s string) *ConnectorUpdate {
	cu.mutation.SetCategory(s)
	return cu
}

// AddResourceIDs adds the "resources" edge to the ApplicationResource entity by IDs.
func (cu *ConnectorUpdate) AddResourceIDs(ids ...oid.ID) *ConnectorUpdate {
	cu.mutation.AddResourceIDs(ids...)
	return cu
}

// AddResources adds the "resources" edges to the ApplicationResource entity.
func (cu *ConnectorUpdate) AddResources(a ...*ApplicationResource) *ConnectorUpdate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return cu.AddResourceIDs(ids...)
}

// AddClusterCostIDs adds the "clusterCosts" edge to the ClusterCost entity by IDs.
func (cu *ConnectorUpdate) AddClusterCostIDs(ids ...int) *ConnectorUpdate {
	cu.mutation.AddClusterCostIDs(ids...)
	return cu
}

// AddClusterCosts adds the "clusterCosts" edges to the ClusterCost entity.
func (cu *ConnectorUpdate) AddClusterCosts(c ...*ClusterCost) *ConnectorUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cu.AddClusterCostIDs(ids...)
}

// AddAllocationCostIDs adds the "allocationCosts" edge to the AllocationCost entity by IDs.
func (cu *ConnectorUpdate) AddAllocationCostIDs(ids ...int) *ConnectorUpdate {
	cu.mutation.AddAllocationCostIDs(ids...)
	return cu
}

// AddAllocationCosts adds the "allocationCosts" edges to the AllocationCost entity.
func (cu *ConnectorUpdate) AddAllocationCosts(a ...*AllocationCost) *ConnectorUpdate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return cu.AddAllocationCostIDs(ids...)
}

// Mutation returns the ConnectorMutation object of the builder.
func (cu *ConnectorUpdate) Mutation() *ConnectorMutation {
	return cu.mutation
}

// ClearResources clears all "resources" edges to the ApplicationResource entity.
func (cu *ConnectorUpdate) ClearResources() *ConnectorUpdate {
	cu.mutation.ClearResources()
	return cu
}

// RemoveResourceIDs removes the "resources" edge to ApplicationResource entities by IDs.
func (cu *ConnectorUpdate) RemoveResourceIDs(ids ...oid.ID) *ConnectorUpdate {
	cu.mutation.RemoveResourceIDs(ids...)
	return cu
}

// RemoveResources removes "resources" edges to ApplicationResource entities.
func (cu *ConnectorUpdate) RemoveResources(a ...*ApplicationResource) *ConnectorUpdate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return cu.RemoveResourceIDs(ids...)
}

// ClearClusterCosts clears all "clusterCosts" edges to the ClusterCost entity.
func (cu *ConnectorUpdate) ClearClusterCosts() *ConnectorUpdate {
	cu.mutation.ClearClusterCosts()
	return cu
}

// RemoveClusterCostIDs removes the "clusterCosts" edge to ClusterCost entities by IDs.
func (cu *ConnectorUpdate) RemoveClusterCostIDs(ids ...int) *ConnectorUpdate {
	cu.mutation.RemoveClusterCostIDs(ids...)
	return cu
}

// RemoveClusterCosts removes "clusterCosts" edges to ClusterCost entities.
func (cu *ConnectorUpdate) RemoveClusterCosts(c ...*ClusterCost) *ConnectorUpdate {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cu.RemoveClusterCostIDs(ids...)
}

// ClearAllocationCosts clears all "allocationCosts" edges to the AllocationCost entity.
func (cu *ConnectorUpdate) ClearAllocationCosts() *ConnectorUpdate {
	cu.mutation.ClearAllocationCosts()
	return cu
}

// RemoveAllocationCostIDs removes the "allocationCosts" edge to AllocationCost entities by IDs.
func (cu *ConnectorUpdate) RemoveAllocationCostIDs(ids ...int) *ConnectorUpdate {
	cu.mutation.RemoveAllocationCostIDs(ids...)
	return cu
}

// RemoveAllocationCosts removes "allocationCosts" edges to AllocationCost entities.
func (cu *ConnectorUpdate) RemoveAllocationCosts(a ...*AllocationCost) *ConnectorUpdate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return cu.RemoveAllocationCostIDs(ids...)
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

// check runs all checks and user-defined validators on the builder.
func (cu *ConnectorUpdate) check() error {
	if v, ok := cu.mutation.Name(); ok {
		if err := connector.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Connector.name": %w`, err)}
		}
	}
	if v, ok := cu.mutation.ConfigVersion(); ok {
		if err := connector.ConfigVersionValidator(v); err != nil {
			return &ValidationError{Name: "configVersion", err: fmt.Errorf(`model: validator failed for field "Connector.configVersion": %w`, err)}
		}
	}
	if v, ok := cu.mutation.Category(); ok {
		if err := connector.CategoryValidator(v); err != nil {
			return &ValidationError{Name: "category", err: fmt.Errorf(`model: validator failed for field "Connector.category": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cu *ConnectorUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ConnectorUpdate {
	cu.modifiers = append(cu.modifiers, modifiers...)
	return cu
}

func (cu *ConnectorUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := cu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(connector.Table, connector.Columns, sqlgraph.NewFieldSpec(connector.FieldID, field.TypeString))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.Name(); ok {
		_spec.SetField(connector.FieldName, field.TypeString, value)
	}
	if value, ok := cu.mutation.Description(); ok {
		_spec.SetField(connector.FieldDescription, field.TypeString, value)
	}
	if cu.mutation.DescriptionCleared() {
		_spec.ClearField(connector.FieldDescription, field.TypeString)
	}
	if value, ok := cu.mutation.Labels(); ok {
		_spec.SetField(connector.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := cu.mutation.UpdateTime(); ok {
		_spec.SetField(connector.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cu.mutation.Status(); ok {
		_spec.SetField(connector.FieldStatus, field.TypeJSON, value)
	}
	if cu.mutation.StatusCleared() {
		_spec.ClearField(connector.FieldStatus, field.TypeJSON)
	}
	if value, ok := cu.mutation.ConfigVersion(); ok {
		_spec.SetField(connector.FieldConfigVersion, field.TypeString, value)
	}
	if value, ok := cu.mutation.ConfigData(); ok {
		_spec.SetField(connector.FieldConfigData, field.TypeOther, value)
	}
	if value, ok := cu.mutation.EnableFinOps(); ok {
		_spec.SetField(connector.FieldEnableFinOps, field.TypeBool, value)
	}
	if value, ok := cu.mutation.FinOpsCustomPricing(); ok {
		_spec.SetField(connector.FieldFinOpsCustomPricing, field.TypeJSON, value)
	}
	if cu.mutation.FinOpsCustomPricingCleared() {
		_spec.ClearField(connector.FieldFinOpsCustomPricing, field.TypeJSON)
	}
	if value, ok := cu.mutation.Category(); ok {
		_spec.SetField(connector.FieldCategory, field.TypeString, value)
	}
	if cu.mutation.ResourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ResourcesTable,
			Columns: []string{connector.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = cu.schemaConfig.ApplicationResource
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedResourcesIDs(); len(nodes) > 0 && !cu.mutation.ResourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ResourcesTable,
			Columns: []string{connector.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = cu.schemaConfig.ApplicationResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.ResourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ResourcesTable,
			Columns: []string{connector.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = cu.schemaConfig.ApplicationResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cu.mutation.ClusterCostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ClusterCostsTable,
			Columns: []string{connector.ClusterCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cu.schemaConfig.ClusterCost
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedClusterCostsIDs(); len(nodes) > 0 && !cu.mutation.ClusterCostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ClusterCostsTable,
			Columns: []string{connector.ClusterCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cu.schemaConfig.ClusterCost
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.ClusterCostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ClusterCostsTable,
			Columns: []string{connector.ClusterCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cu.schemaConfig.ClusterCost
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cu.mutation.AllocationCostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.AllocationCostsTable,
			Columns: []string{connector.AllocationCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(allocationcost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cu.schemaConfig.AllocationCost
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedAllocationCostsIDs(); len(nodes) > 0 && !cu.mutation.AllocationCostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.AllocationCostsTable,
			Columns: []string{connector.AllocationCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(allocationcost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cu.schemaConfig.AllocationCost
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.AllocationCostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.AllocationCostsTable,
			Columns: []string{connector.AllocationCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(allocationcost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cu.schemaConfig.AllocationCost
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
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

// SetName sets the "name" field.
func (cuo *ConnectorUpdateOne) SetName(s string) *ConnectorUpdateOne {
	cuo.mutation.SetName(s)
	return cuo
}

// SetDescription sets the "description" field.
func (cuo *ConnectorUpdateOne) SetDescription(s string) *ConnectorUpdateOne {
	cuo.mutation.SetDescription(s)
	return cuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cuo *ConnectorUpdateOne) SetNillableDescription(s *string) *ConnectorUpdateOne {
	if s != nil {
		cuo.SetDescription(*s)
	}
	return cuo
}

// ClearDescription clears the value of the "description" field.
func (cuo *ConnectorUpdateOne) ClearDescription() *ConnectorUpdateOne {
	cuo.mutation.ClearDescription()
	return cuo
}

// SetLabels sets the "labels" field.
func (cuo *ConnectorUpdateOne) SetLabels(m map[string]string) *ConnectorUpdateOne {
	cuo.mutation.SetLabels(m)
	return cuo
}

// SetUpdateTime sets the "updateTime" field.
func (cuo *ConnectorUpdateOne) SetUpdateTime(t time.Time) *ConnectorUpdateOne {
	cuo.mutation.SetUpdateTime(t)
	return cuo
}

// SetStatus sets the "status" field.
func (cuo *ConnectorUpdateOne) SetStatus(s status.Status) *ConnectorUpdateOne {
	cuo.mutation.SetStatus(s)
	return cuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (cuo *ConnectorUpdateOne) SetNillableStatus(s *status.Status) *ConnectorUpdateOne {
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

// SetConfigVersion sets the "configVersion" field.
func (cuo *ConnectorUpdateOne) SetConfigVersion(s string) *ConnectorUpdateOne {
	cuo.mutation.SetConfigVersion(s)
	return cuo
}

// SetConfigData sets the "configData" field.
func (cuo *ConnectorUpdateOne) SetConfigData(c crypto.Properties) *ConnectorUpdateOne {
	cuo.mutation.SetConfigData(c)
	return cuo
}

// SetEnableFinOps sets the "enableFinOps" field.
func (cuo *ConnectorUpdateOne) SetEnableFinOps(b bool) *ConnectorUpdateOne {
	cuo.mutation.SetEnableFinOps(b)
	return cuo
}

// SetFinOpsCustomPricing sets the "finOpsCustomPricing" field.
func (cuo *ConnectorUpdateOne) SetFinOpsCustomPricing(tocp *types.FinOpsCustomPricing) *ConnectorUpdateOne {
	cuo.mutation.SetFinOpsCustomPricing(tocp)
	return cuo
}

// ClearFinOpsCustomPricing clears the value of the "finOpsCustomPricing" field.
func (cuo *ConnectorUpdateOne) ClearFinOpsCustomPricing() *ConnectorUpdateOne {
	cuo.mutation.ClearFinOpsCustomPricing()
	return cuo
}

// SetCategory sets the "category" field.
func (cuo *ConnectorUpdateOne) SetCategory(s string) *ConnectorUpdateOne {
	cuo.mutation.SetCategory(s)
	return cuo
}

// AddResourceIDs adds the "resources" edge to the ApplicationResource entity by IDs.
func (cuo *ConnectorUpdateOne) AddResourceIDs(ids ...oid.ID) *ConnectorUpdateOne {
	cuo.mutation.AddResourceIDs(ids...)
	return cuo
}

// AddResources adds the "resources" edges to the ApplicationResource entity.
func (cuo *ConnectorUpdateOne) AddResources(a ...*ApplicationResource) *ConnectorUpdateOne {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return cuo.AddResourceIDs(ids...)
}

// AddClusterCostIDs adds the "clusterCosts" edge to the ClusterCost entity by IDs.
func (cuo *ConnectorUpdateOne) AddClusterCostIDs(ids ...int) *ConnectorUpdateOne {
	cuo.mutation.AddClusterCostIDs(ids...)
	return cuo
}

// AddClusterCosts adds the "clusterCosts" edges to the ClusterCost entity.
func (cuo *ConnectorUpdateOne) AddClusterCosts(c ...*ClusterCost) *ConnectorUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cuo.AddClusterCostIDs(ids...)
}

// AddAllocationCostIDs adds the "allocationCosts" edge to the AllocationCost entity by IDs.
func (cuo *ConnectorUpdateOne) AddAllocationCostIDs(ids ...int) *ConnectorUpdateOne {
	cuo.mutation.AddAllocationCostIDs(ids...)
	return cuo
}

// AddAllocationCosts adds the "allocationCosts" edges to the AllocationCost entity.
func (cuo *ConnectorUpdateOne) AddAllocationCosts(a ...*AllocationCost) *ConnectorUpdateOne {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return cuo.AddAllocationCostIDs(ids...)
}

// Mutation returns the ConnectorMutation object of the builder.
func (cuo *ConnectorUpdateOne) Mutation() *ConnectorMutation {
	return cuo.mutation
}

// ClearResources clears all "resources" edges to the ApplicationResource entity.
func (cuo *ConnectorUpdateOne) ClearResources() *ConnectorUpdateOne {
	cuo.mutation.ClearResources()
	return cuo
}

// RemoveResourceIDs removes the "resources" edge to ApplicationResource entities by IDs.
func (cuo *ConnectorUpdateOne) RemoveResourceIDs(ids ...oid.ID) *ConnectorUpdateOne {
	cuo.mutation.RemoveResourceIDs(ids...)
	return cuo
}

// RemoveResources removes "resources" edges to ApplicationResource entities.
func (cuo *ConnectorUpdateOne) RemoveResources(a ...*ApplicationResource) *ConnectorUpdateOne {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return cuo.RemoveResourceIDs(ids...)
}

// ClearClusterCosts clears all "clusterCosts" edges to the ClusterCost entity.
func (cuo *ConnectorUpdateOne) ClearClusterCosts() *ConnectorUpdateOne {
	cuo.mutation.ClearClusterCosts()
	return cuo
}

// RemoveClusterCostIDs removes the "clusterCosts" edge to ClusterCost entities by IDs.
func (cuo *ConnectorUpdateOne) RemoveClusterCostIDs(ids ...int) *ConnectorUpdateOne {
	cuo.mutation.RemoveClusterCostIDs(ids...)
	return cuo
}

// RemoveClusterCosts removes "clusterCosts" edges to ClusterCost entities.
func (cuo *ConnectorUpdateOne) RemoveClusterCosts(c ...*ClusterCost) *ConnectorUpdateOne {
	ids := make([]int, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return cuo.RemoveClusterCostIDs(ids...)
}

// ClearAllocationCosts clears all "allocationCosts" edges to the AllocationCost entity.
func (cuo *ConnectorUpdateOne) ClearAllocationCosts() *ConnectorUpdateOne {
	cuo.mutation.ClearAllocationCosts()
	return cuo
}

// RemoveAllocationCostIDs removes the "allocationCosts" edge to AllocationCost entities by IDs.
func (cuo *ConnectorUpdateOne) RemoveAllocationCostIDs(ids ...int) *ConnectorUpdateOne {
	cuo.mutation.RemoveAllocationCostIDs(ids...)
	return cuo
}

// RemoveAllocationCosts removes "allocationCosts" edges to AllocationCost entities.
func (cuo *ConnectorUpdateOne) RemoveAllocationCosts(a ...*AllocationCost) *ConnectorUpdateOne {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return cuo.RemoveAllocationCostIDs(ids...)
}

// Where appends a list predicates to the ConnectorUpdate builder.
func (cuo *ConnectorUpdateOne) Where(ps ...predicate.Connector) *ConnectorUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
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

// check runs all checks and user-defined validators on the builder.
func (cuo *ConnectorUpdateOne) check() error {
	if v, ok := cuo.mutation.Name(); ok {
		if err := connector.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Connector.name": %w`, err)}
		}
	}
	if v, ok := cuo.mutation.ConfigVersion(); ok {
		if err := connector.ConfigVersionValidator(v); err != nil {
			return &ValidationError{Name: "configVersion", err: fmt.Errorf(`model: validator failed for field "Connector.configVersion": %w`, err)}
		}
	}
	if v, ok := cuo.mutation.Category(); ok {
		if err := connector.CategoryValidator(v); err != nil {
			return &ValidationError{Name: "category", err: fmt.Errorf(`model: validator failed for field "Connector.category": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cuo *ConnectorUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ConnectorUpdateOne {
	cuo.modifiers = append(cuo.modifiers, modifiers...)
	return cuo
}

func (cuo *ConnectorUpdateOne) sqlSave(ctx context.Context) (_node *Connector, err error) {
	if err := cuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(connector.Table, connector.Columns, sqlgraph.NewFieldSpec(connector.FieldID, field.TypeString))
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
	if value, ok := cuo.mutation.Name(); ok {
		_spec.SetField(connector.FieldName, field.TypeString, value)
	}
	if value, ok := cuo.mutation.Description(); ok {
		_spec.SetField(connector.FieldDescription, field.TypeString, value)
	}
	if cuo.mutation.DescriptionCleared() {
		_spec.ClearField(connector.FieldDescription, field.TypeString)
	}
	if value, ok := cuo.mutation.Labels(); ok {
		_spec.SetField(connector.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := cuo.mutation.UpdateTime(); ok {
		_spec.SetField(connector.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cuo.mutation.Status(); ok {
		_spec.SetField(connector.FieldStatus, field.TypeJSON, value)
	}
	if cuo.mutation.StatusCleared() {
		_spec.ClearField(connector.FieldStatus, field.TypeJSON)
	}
	if value, ok := cuo.mutation.ConfigVersion(); ok {
		_spec.SetField(connector.FieldConfigVersion, field.TypeString, value)
	}
	if value, ok := cuo.mutation.ConfigData(); ok {
		_spec.SetField(connector.FieldConfigData, field.TypeOther, value)
	}
	if value, ok := cuo.mutation.EnableFinOps(); ok {
		_spec.SetField(connector.FieldEnableFinOps, field.TypeBool, value)
	}
	if value, ok := cuo.mutation.FinOpsCustomPricing(); ok {
		_spec.SetField(connector.FieldFinOpsCustomPricing, field.TypeJSON, value)
	}
	if cuo.mutation.FinOpsCustomPricingCleared() {
		_spec.ClearField(connector.FieldFinOpsCustomPricing, field.TypeJSON)
	}
	if value, ok := cuo.mutation.Category(); ok {
		_spec.SetField(connector.FieldCategory, field.TypeString, value)
	}
	if cuo.mutation.ResourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ResourcesTable,
			Columns: []string{connector.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = cuo.schemaConfig.ApplicationResource
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedResourcesIDs(); len(nodes) > 0 && !cuo.mutation.ResourcesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ResourcesTable,
			Columns: []string{connector.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = cuo.schemaConfig.ApplicationResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.ResourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ResourcesTable,
			Columns: []string{connector.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = cuo.schemaConfig.ApplicationResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cuo.mutation.ClusterCostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ClusterCostsTable,
			Columns: []string{connector.ClusterCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cuo.schemaConfig.ClusterCost
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedClusterCostsIDs(); len(nodes) > 0 && !cuo.mutation.ClusterCostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ClusterCostsTable,
			Columns: []string{connector.ClusterCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cuo.schemaConfig.ClusterCost
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.ClusterCostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.ClusterCostsTable,
			Columns: []string{connector.ClusterCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cuo.schemaConfig.ClusterCost
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cuo.mutation.AllocationCostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.AllocationCostsTable,
			Columns: []string{connector.AllocationCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(allocationcost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cuo.schemaConfig.AllocationCost
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedAllocationCostsIDs(); len(nodes) > 0 && !cuo.mutation.AllocationCostsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.AllocationCostsTable,
			Columns: []string{connector.AllocationCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(allocationcost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cuo.schemaConfig.AllocationCost
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.AllocationCostsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   connector.AllocationCostsTable,
			Columns: []string{connector.AllocationCostsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(allocationcost.FieldID, field.TypeInt),
			},
		}
		edge.Schema = cuo.schemaConfig.AllocationCost
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
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
