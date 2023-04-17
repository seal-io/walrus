// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ApplicationResourceCreate is the builder for creating a ApplicationResource entity.
type ApplicationResourceCreate struct {
	config
	mutation *ApplicationResourceMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "createTime" field.
func (arc *ApplicationResourceCreate) SetCreateTime(t time.Time) *ApplicationResourceCreate {
	arc.mutation.SetCreateTime(t)
	return arc
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (arc *ApplicationResourceCreate) SetNillableCreateTime(t *time.Time) *ApplicationResourceCreate {
	if t != nil {
		arc.SetCreateTime(*t)
	}
	return arc
}

// SetUpdateTime sets the "updateTime" field.
func (arc *ApplicationResourceCreate) SetUpdateTime(t time.Time) *ApplicationResourceCreate {
	arc.mutation.SetUpdateTime(t)
	return arc
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (arc *ApplicationResourceCreate) SetNillableUpdateTime(t *time.Time) *ApplicationResourceCreate {
	if t != nil {
		arc.SetUpdateTime(*t)
	}
	return arc
}

// SetInstanceID sets the "instanceID" field.
func (arc *ApplicationResourceCreate) SetInstanceID(o oid.ID) *ApplicationResourceCreate {
	arc.mutation.SetInstanceID(o)
	return arc
}

// SetConnectorID sets the "connectorID" field.
func (arc *ApplicationResourceCreate) SetConnectorID(o oid.ID) *ApplicationResourceCreate {
	arc.mutation.SetConnectorID(o)
	return arc
}

// SetCompositionID sets the "compositionID" field.
func (arc *ApplicationResourceCreate) SetCompositionID(o oid.ID) *ApplicationResourceCreate {
	arc.mutation.SetCompositionID(o)
	return arc
}

// SetNillableCompositionID sets the "compositionID" field if the given value is not nil.
func (arc *ApplicationResourceCreate) SetNillableCompositionID(o *oid.ID) *ApplicationResourceCreate {
	if o != nil {
		arc.SetCompositionID(*o)
	}
	return arc
}

// SetModule sets the "module" field.
func (arc *ApplicationResourceCreate) SetModule(s string) *ApplicationResourceCreate {
	arc.mutation.SetModule(s)
	return arc
}

// SetMode sets the "mode" field.
func (arc *ApplicationResourceCreate) SetMode(s string) *ApplicationResourceCreate {
	arc.mutation.SetMode(s)
	return arc
}

// SetType sets the "type" field.
func (arc *ApplicationResourceCreate) SetType(s string) *ApplicationResourceCreate {
	arc.mutation.SetType(s)
	return arc
}

// SetName sets the "name" field.
func (arc *ApplicationResourceCreate) SetName(s string) *ApplicationResourceCreate {
	arc.mutation.SetName(s)
	return arc
}

// SetDeployerType sets the "deployerType" field.
func (arc *ApplicationResourceCreate) SetDeployerType(s string) *ApplicationResourceCreate {
	arc.mutation.SetDeployerType(s)
	return arc
}

// SetStatus sets the "status" field.
func (arc *ApplicationResourceCreate) SetStatus(trs types.ApplicationResourceStatus) *ApplicationResourceCreate {
	arc.mutation.SetStatus(trs)
	return arc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (arc *ApplicationResourceCreate) SetNillableStatus(trs *types.ApplicationResourceStatus) *ApplicationResourceCreate {
	if trs != nil {
		arc.SetStatus(*trs)
	}
	return arc
}

// SetID sets the "id" field.
func (arc *ApplicationResourceCreate) SetID(o oid.ID) *ApplicationResourceCreate {
	arc.mutation.SetID(o)
	return arc
}

// SetInstance sets the "instance" edge to the ApplicationInstance entity.
func (arc *ApplicationResourceCreate) SetInstance(a *ApplicationInstance) *ApplicationResourceCreate {
	return arc.SetInstanceID(a.ID)
}

// SetConnector sets the "connector" edge to the Connector entity.
func (arc *ApplicationResourceCreate) SetConnector(c *Connector) *ApplicationResourceCreate {
	return arc.SetConnectorID(c.ID)
}

// SetComposition sets the "composition" edge to the ApplicationResource entity.
func (arc *ApplicationResourceCreate) SetComposition(a *ApplicationResource) *ApplicationResourceCreate {
	return arc.SetCompositionID(a.ID)
}

// AddComponentIDs adds the "components" edge to the ApplicationResource entity by IDs.
func (arc *ApplicationResourceCreate) AddComponentIDs(ids ...oid.ID) *ApplicationResourceCreate {
	arc.mutation.AddComponentIDs(ids...)
	return arc
}

// AddComponents adds the "components" edges to the ApplicationResource entity.
func (arc *ApplicationResourceCreate) AddComponents(a ...*ApplicationResource) *ApplicationResourceCreate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return arc.AddComponentIDs(ids...)
}

// Mutation returns the ApplicationResourceMutation object of the builder.
func (arc *ApplicationResourceCreate) Mutation() *ApplicationResourceMutation {
	return arc.mutation
}

// Save creates the ApplicationResource in the database.
func (arc *ApplicationResourceCreate) Save(ctx context.Context) (*ApplicationResource, error) {
	if err := arc.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*ApplicationResource, ApplicationResourceMutation](ctx, arc.sqlSave, arc.mutation, arc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (arc *ApplicationResourceCreate) SaveX(ctx context.Context) *ApplicationResource {
	v, err := arc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (arc *ApplicationResourceCreate) Exec(ctx context.Context) error {
	_, err := arc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (arc *ApplicationResourceCreate) ExecX(ctx context.Context) {
	if err := arc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (arc *ApplicationResourceCreate) defaults() error {
	if _, ok := arc.mutation.CreateTime(); !ok {
		if applicationresource.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized applicationresource.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := applicationresource.DefaultCreateTime()
		arc.mutation.SetCreateTime(v)
	}
	if _, ok := arc.mutation.UpdateTime(); !ok {
		if applicationresource.DefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized applicationresource.DefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := applicationresource.DefaultUpdateTime()
		arc.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (arc *ApplicationResourceCreate) check() error {
	if _, ok := arc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "ApplicationResource.createTime"`)}
	}
	if _, ok := arc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "ApplicationResource.updateTime"`)}
	}
	if _, ok := arc.mutation.InstanceID(); !ok {
		return &ValidationError{Name: "instanceID", err: errors.New(`model: missing required field "ApplicationResource.instanceID"`)}
	}
	if v, ok := arc.mutation.InstanceID(); ok {
		if err := applicationresource.InstanceIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "instanceID", err: fmt.Errorf(`model: validator failed for field "ApplicationResource.instanceID": %w`, err)}
		}
	}
	if _, ok := arc.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connectorID", err: errors.New(`model: missing required field "ApplicationResource.connectorID"`)}
	}
	if v, ok := arc.mutation.ConnectorID(); ok {
		if err := applicationresource.ConnectorIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "connectorID", err: fmt.Errorf(`model: validator failed for field "ApplicationResource.connectorID": %w`, err)}
		}
	}
	if _, ok := arc.mutation.Module(); !ok {
		return &ValidationError{Name: "module", err: errors.New(`model: missing required field "ApplicationResource.module"`)}
	}
	if v, ok := arc.mutation.Module(); ok {
		if err := applicationresource.ModuleValidator(v); err != nil {
			return &ValidationError{Name: "module", err: fmt.Errorf(`model: validator failed for field "ApplicationResource.module": %w`, err)}
		}
	}
	if _, ok := arc.mutation.Mode(); !ok {
		return &ValidationError{Name: "mode", err: errors.New(`model: missing required field "ApplicationResource.mode"`)}
	}
	if v, ok := arc.mutation.Mode(); ok {
		if err := applicationresource.ModeValidator(v); err != nil {
			return &ValidationError{Name: "mode", err: fmt.Errorf(`model: validator failed for field "ApplicationResource.mode": %w`, err)}
		}
	}
	if _, ok := arc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`model: missing required field "ApplicationResource.type"`)}
	}
	if v, ok := arc.mutation.GetType(); ok {
		if err := applicationresource.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`model: validator failed for field "ApplicationResource.type": %w`, err)}
		}
	}
	if _, ok := arc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "ApplicationResource.name"`)}
	}
	if v, ok := arc.mutation.Name(); ok {
		if err := applicationresource.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "ApplicationResource.name": %w`, err)}
		}
	}
	if _, ok := arc.mutation.DeployerType(); !ok {
		return &ValidationError{Name: "deployerType", err: errors.New(`model: missing required field "ApplicationResource.deployerType"`)}
	}
	if v, ok := arc.mutation.DeployerType(); ok {
		if err := applicationresource.DeployerTypeValidator(v); err != nil {
			return &ValidationError{Name: "deployerType", err: fmt.Errorf(`model: validator failed for field "ApplicationResource.deployerType": %w`, err)}
		}
	}
	if _, ok := arc.mutation.InstanceID(); !ok {
		return &ValidationError{Name: "instance", err: errors.New(`model: missing required edge "ApplicationResource.instance"`)}
	}
	if _, ok := arc.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connector", err: errors.New(`model: missing required edge "ApplicationResource.connector"`)}
	}
	return nil
}

func (arc *ApplicationResourceCreate) sqlSave(ctx context.Context) (*ApplicationResource, error) {
	if err := arc.check(); err != nil {
		return nil, err
	}
	_node, _spec := arc.createSpec()
	if err := sqlgraph.CreateNode(ctx, arc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*oid.ID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	arc.mutation.id = &_node.ID
	arc.mutation.done = true
	return _node, nil
}

func (arc *ApplicationResourceCreate) createSpec() (*ApplicationResource, *sqlgraph.CreateSpec) {
	var (
		_node = &ApplicationResource{config: arc.config}
		_spec = sqlgraph.NewCreateSpec(applicationresource.Table, sqlgraph.NewFieldSpec(applicationresource.FieldID, field.TypeString))
	)
	_spec.Schema = arc.schemaConfig.ApplicationResource
	_spec.OnConflict = arc.conflict
	if id, ok := arc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := arc.mutation.CreateTime(); ok {
		_spec.SetField(applicationresource.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := arc.mutation.UpdateTime(); ok {
		_spec.SetField(applicationresource.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := arc.mutation.Module(); ok {
		_spec.SetField(applicationresource.FieldModule, field.TypeString, value)
		_node.Module = value
	}
	if value, ok := arc.mutation.Mode(); ok {
		_spec.SetField(applicationresource.FieldMode, field.TypeString, value)
		_node.Mode = value
	}
	if value, ok := arc.mutation.GetType(); ok {
		_spec.SetField(applicationresource.FieldType, field.TypeString, value)
		_node.Type = value
	}
	if value, ok := arc.mutation.Name(); ok {
		_spec.SetField(applicationresource.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := arc.mutation.DeployerType(); ok {
		_spec.SetField(applicationresource.FieldDeployerType, field.TypeString, value)
		_node.DeployerType = value
	}
	if value, ok := arc.mutation.Status(); ok {
		_spec.SetField(applicationresource.FieldStatus, field.TypeJSON, value)
		_node.Status = value
	}
	if nodes := arc.mutation.InstanceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   applicationresource.InstanceTable,
			Columns: []string{applicationresource.InstanceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationinstance.FieldID,
				},
			},
		}
		edge.Schema = arc.schemaConfig.ApplicationResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.InstanceID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := arc.mutation.ConnectorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   applicationresource.ConnectorTable,
			Columns: []string{applicationresource.ConnectorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = arc.schemaConfig.ApplicationResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ConnectorID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := arc.mutation.CompositionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   applicationresource.CompositionTable,
			Columns: []string{applicationresource.CompositionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationresource.FieldID,
				},
			},
		}
		edge.Schema = arc.schemaConfig.ApplicationResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.CompositionID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := arc.mutation.ComponentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   applicationresource.ComponentsTable,
			Columns: []string{applicationresource.ComponentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationresource.FieldID,
				},
			},
		}
		edge.Schema = arc.schemaConfig.ApplicationResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ApplicationResource.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationResourceUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (arc *ApplicationResourceCreate) OnConflict(opts ...sql.ConflictOption) *ApplicationResourceUpsertOne {
	arc.conflict = opts
	return &ApplicationResourceUpsertOne{
		create: arc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ApplicationResource.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (arc *ApplicationResourceCreate) OnConflictColumns(columns ...string) *ApplicationResourceUpsertOne {
	arc.conflict = append(arc.conflict, sql.ConflictColumns(columns...))
	return &ApplicationResourceUpsertOne{
		create: arc,
	}
}

type (
	// ApplicationResourceUpsertOne is the builder for "upsert"-ing
	//  one ApplicationResource node.
	ApplicationResourceUpsertOne struct {
		create *ApplicationResourceCreate
	}

	// ApplicationResourceUpsert is the "OnConflict" setter.
	ApplicationResourceUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationResourceUpsert) SetUpdateTime(v time.Time) *ApplicationResourceUpsert {
	u.Set(applicationresource.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationResourceUpsert) UpdateUpdateTime() *ApplicationResourceUpsert {
	u.SetExcluded(applicationresource.FieldUpdateTime)
	return u
}

// SetStatus sets the "status" field.
func (u *ApplicationResourceUpsert) SetStatus(v types.ApplicationResourceStatus) *ApplicationResourceUpsert {
	u.Set(applicationresource.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ApplicationResourceUpsert) UpdateStatus() *ApplicationResourceUpsert {
	u.SetExcluded(applicationresource.FieldStatus)
	return u
}

// ClearStatus clears the value of the "status" field.
func (u *ApplicationResourceUpsert) ClearStatus() *ApplicationResourceUpsert {
	u.SetNull(applicationresource.FieldStatus)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.ApplicationResource.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(applicationresource.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ApplicationResourceUpsertOne) UpdateNewValues() *ApplicationResourceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(applicationresource.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(applicationresource.FieldCreateTime)
		}
		if _, exists := u.create.mutation.InstanceID(); exists {
			s.SetIgnore(applicationresource.FieldInstanceID)
		}
		if _, exists := u.create.mutation.ConnectorID(); exists {
			s.SetIgnore(applicationresource.FieldConnectorID)
		}
		if _, exists := u.create.mutation.CompositionID(); exists {
			s.SetIgnore(applicationresource.FieldCompositionID)
		}
		if _, exists := u.create.mutation.Module(); exists {
			s.SetIgnore(applicationresource.FieldModule)
		}
		if _, exists := u.create.mutation.Mode(); exists {
			s.SetIgnore(applicationresource.FieldMode)
		}
		if _, exists := u.create.mutation.GetType(); exists {
			s.SetIgnore(applicationresource.FieldType)
		}
		if _, exists := u.create.mutation.Name(); exists {
			s.SetIgnore(applicationresource.FieldName)
		}
		if _, exists := u.create.mutation.DeployerType(); exists {
			s.SetIgnore(applicationresource.FieldDeployerType)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ApplicationResource.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ApplicationResourceUpsertOne) Ignore() *ApplicationResourceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationResourceUpsertOne) DoNothing() *ApplicationResourceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationResourceCreate.OnConflict
// documentation for more info.
func (u *ApplicationResourceUpsertOne) Update(set func(*ApplicationResourceUpsert)) *ApplicationResourceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationResourceUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationResourceUpsertOne) SetUpdateTime(v time.Time) *ApplicationResourceUpsertOne {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationResourceUpsertOne) UpdateUpdateTime() *ApplicationResourceUpsertOne {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetStatus sets the "status" field.
func (u *ApplicationResourceUpsertOne) SetStatus(v types.ApplicationResourceStatus) *ApplicationResourceUpsertOne {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ApplicationResourceUpsertOne) UpdateStatus() *ApplicationResourceUpsertOne {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.UpdateStatus()
	})
}

// ClearStatus clears the value of the "status" field.
func (u *ApplicationResourceUpsertOne) ClearStatus() *ApplicationResourceUpsertOne {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.ClearStatus()
	})
}

// Exec executes the query.
func (u *ApplicationResourceUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationResourceCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationResourceUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ApplicationResourceUpsertOne) ID(ctx context.Context) (id oid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: ApplicationResourceUpsertOne.ID is not supported by MySQL driver. Use ApplicationResourceUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ApplicationResourceUpsertOne) IDX(ctx context.Context) oid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ApplicationResourceCreateBulk is the builder for creating many ApplicationResource entities in bulk.
type ApplicationResourceCreateBulk struct {
	config
	builders []*ApplicationResourceCreate
	conflict []sql.ConflictOption
}

// Save creates the ApplicationResource entities in the database.
func (arcb *ApplicationResourceCreateBulk) Save(ctx context.Context) ([]*ApplicationResource, error) {
	specs := make([]*sqlgraph.CreateSpec, len(arcb.builders))
	nodes := make([]*ApplicationResource, len(arcb.builders))
	mutators := make([]Mutator, len(arcb.builders))
	for i := range arcb.builders {
		func(i int, root context.Context) {
			builder := arcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ApplicationResourceMutation)
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
					_, err = mutators[i+1].Mutate(root, arcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = arcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, arcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
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
		if _, err := mutators[0].Mutate(ctx, arcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (arcb *ApplicationResourceCreateBulk) SaveX(ctx context.Context) []*ApplicationResource {
	v, err := arcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (arcb *ApplicationResourceCreateBulk) Exec(ctx context.Context) error {
	_, err := arcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (arcb *ApplicationResourceCreateBulk) ExecX(ctx context.Context) {
	if err := arcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ApplicationResource.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationResourceUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (arcb *ApplicationResourceCreateBulk) OnConflict(opts ...sql.ConflictOption) *ApplicationResourceUpsertBulk {
	arcb.conflict = opts
	return &ApplicationResourceUpsertBulk{
		create: arcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ApplicationResource.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (arcb *ApplicationResourceCreateBulk) OnConflictColumns(columns ...string) *ApplicationResourceUpsertBulk {
	arcb.conflict = append(arcb.conflict, sql.ConflictColumns(columns...))
	return &ApplicationResourceUpsertBulk{
		create: arcb,
	}
}

// ApplicationResourceUpsertBulk is the builder for "upsert"-ing
// a bulk of ApplicationResource nodes.
type ApplicationResourceUpsertBulk struct {
	create *ApplicationResourceCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ApplicationResource.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(applicationresource.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ApplicationResourceUpsertBulk) UpdateNewValues() *ApplicationResourceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(applicationresource.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(applicationresource.FieldCreateTime)
			}
			if _, exists := b.mutation.InstanceID(); exists {
				s.SetIgnore(applicationresource.FieldInstanceID)
			}
			if _, exists := b.mutation.ConnectorID(); exists {
				s.SetIgnore(applicationresource.FieldConnectorID)
			}
			if _, exists := b.mutation.CompositionID(); exists {
				s.SetIgnore(applicationresource.FieldCompositionID)
			}
			if _, exists := b.mutation.Module(); exists {
				s.SetIgnore(applicationresource.FieldModule)
			}
			if _, exists := b.mutation.Mode(); exists {
				s.SetIgnore(applicationresource.FieldMode)
			}
			if _, exists := b.mutation.GetType(); exists {
				s.SetIgnore(applicationresource.FieldType)
			}
			if _, exists := b.mutation.Name(); exists {
				s.SetIgnore(applicationresource.FieldName)
			}
			if _, exists := b.mutation.DeployerType(); exists {
				s.SetIgnore(applicationresource.FieldDeployerType)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ApplicationResource.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ApplicationResourceUpsertBulk) Ignore() *ApplicationResourceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationResourceUpsertBulk) DoNothing() *ApplicationResourceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationResourceCreateBulk.OnConflict
// documentation for more info.
func (u *ApplicationResourceUpsertBulk) Update(set func(*ApplicationResourceUpsert)) *ApplicationResourceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationResourceUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationResourceUpsertBulk) SetUpdateTime(v time.Time) *ApplicationResourceUpsertBulk {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationResourceUpsertBulk) UpdateUpdateTime() *ApplicationResourceUpsertBulk {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetStatus sets the "status" field.
func (u *ApplicationResourceUpsertBulk) SetStatus(v types.ApplicationResourceStatus) *ApplicationResourceUpsertBulk {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ApplicationResourceUpsertBulk) UpdateStatus() *ApplicationResourceUpsertBulk {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.UpdateStatus()
	})
}

// ClearStatus clears the value of the "status" field.
func (u *ApplicationResourceUpsertBulk) ClearStatus() *ApplicationResourceUpsertBulk {
	return u.Update(func(s *ApplicationResourceUpsert) {
		s.ClearStatus()
	})
}

// Exec executes the query.
func (u *ApplicationResourceUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ApplicationResourceCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationResourceCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationResourceUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
