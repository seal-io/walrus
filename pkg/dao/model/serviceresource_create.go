// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

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

	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ServiceResourceCreate is the builder for creating a ServiceResource entity.
type ServiceResourceCreate struct {
	config
	mutation *ServiceResourceMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetProjectID sets the "projectID" field.
func (src *ServiceResourceCreate) SetProjectID(o oid.ID) *ServiceResourceCreate {
	src.mutation.SetProjectID(o)
	return src
}

// SetCreateTime sets the "createTime" field.
func (src *ServiceResourceCreate) SetCreateTime(t time.Time) *ServiceResourceCreate {
	src.mutation.SetCreateTime(t)
	return src
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (src *ServiceResourceCreate) SetNillableCreateTime(t *time.Time) *ServiceResourceCreate {
	if t != nil {
		src.SetCreateTime(*t)
	}
	return src
}

// SetUpdateTime sets the "updateTime" field.
func (src *ServiceResourceCreate) SetUpdateTime(t time.Time) *ServiceResourceCreate {
	src.mutation.SetUpdateTime(t)
	return src
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (src *ServiceResourceCreate) SetNillableUpdateTime(t *time.Time) *ServiceResourceCreate {
	if t != nil {
		src.SetUpdateTime(*t)
	}
	return src
}

// SetServiceID sets the "serviceID" field.
func (src *ServiceResourceCreate) SetServiceID(o oid.ID) *ServiceResourceCreate {
	src.mutation.SetServiceID(o)
	return src
}

// SetConnectorID sets the "connectorID" field.
func (src *ServiceResourceCreate) SetConnectorID(o oid.ID) *ServiceResourceCreate {
	src.mutation.SetConnectorID(o)
	return src
}

// SetCompositionID sets the "compositionID" field.
func (src *ServiceResourceCreate) SetCompositionID(o oid.ID) *ServiceResourceCreate {
	src.mutation.SetCompositionID(o)
	return src
}

// SetNillableCompositionID sets the "compositionID" field if the given value is not nil.
func (src *ServiceResourceCreate) SetNillableCompositionID(o *oid.ID) *ServiceResourceCreate {
	if o != nil {
		src.SetCompositionID(*o)
	}
	return src
}

// SetMode sets the "mode" field.
func (src *ServiceResourceCreate) SetMode(s string) *ServiceResourceCreate {
	src.mutation.SetMode(s)
	return src
}

// SetType sets the "type" field.
func (src *ServiceResourceCreate) SetType(s string) *ServiceResourceCreate {
	src.mutation.SetType(s)
	return src
}

// SetName sets the "name" field.
func (src *ServiceResourceCreate) SetName(s string) *ServiceResourceCreate {
	src.mutation.SetName(s)
	return src
}

// SetDeployerType sets the "deployerType" field.
func (src *ServiceResourceCreate) SetDeployerType(s string) *ServiceResourceCreate {
	src.mutation.SetDeployerType(s)
	return src
}

// SetStatus sets the "status" field.
func (src *ServiceResourceCreate) SetStatus(trs types.ServiceResourceStatus) *ServiceResourceCreate {
	src.mutation.SetStatus(trs)
	return src
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (src *ServiceResourceCreate) SetNillableStatus(trs *types.ServiceResourceStatus) *ServiceResourceCreate {
	if trs != nil {
		src.SetStatus(*trs)
	}
	return src
}

// SetID sets the "id" field.
func (src *ServiceResourceCreate) SetID(o oid.ID) *ServiceResourceCreate {
	src.mutation.SetID(o)
	return src
}

// SetService sets the "service" edge to the Service entity.
func (src *ServiceResourceCreate) SetService(s *Service) *ServiceResourceCreate {
	return src.SetServiceID(s.ID)
}

// SetConnector sets the "connector" edge to the Connector entity.
func (src *ServiceResourceCreate) SetConnector(c *Connector) *ServiceResourceCreate {
	return src.SetConnectorID(c.ID)
}

// SetComposition sets the "composition" edge to the ServiceResource entity.
func (src *ServiceResourceCreate) SetComposition(s *ServiceResource) *ServiceResourceCreate {
	return src.SetCompositionID(s.ID)
}

// AddComponentIDs adds the "components" edge to the ServiceResource entity by IDs.
func (src *ServiceResourceCreate) AddComponentIDs(ids ...oid.ID) *ServiceResourceCreate {
	src.mutation.AddComponentIDs(ids...)
	return src
}

// AddComponents adds the "components" edges to the ServiceResource entity.
func (src *ServiceResourceCreate) AddComponents(s ...*ServiceResource) *ServiceResourceCreate {
	ids := make([]oid.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return src.AddComponentIDs(ids...)
}

// Mutation returns the ServiceResourceMutation object of the builder.
func (src *ServiceResourceCreate) Mutation() *ServiceResourceMutation {
	return src.mutation
}

// Save creates the ServiceResource in the database.
func (src *ServiceResourceCreate) Save(ctx context.Context) (*ServiceResource, error) {
	if err := src.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*ServiceResource, ServiceResourceMutation](ctx, src.sqlSave, src.mutation, src.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (src *ServiceResourceCreate) SaveX(ctx context.Context) *ServiceResource {
	v, err := src.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (src *ServiceResourceCreate) Exec(ctx context.Context) error {
	_, err := src.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (src *ServiceResourceCreate) ExecX(ctx context.Context) {
	if err := src.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (src *ServiceResourceCreate) defaults() error {
	if _, ok := src.mutation.CreateTime(); !ok {
		if serviceresource.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized serviceresource.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := serviceresource.DefaultCreateTime()
		src.mutation.SetCreateTime(v)
	}
	if _, ok := src.mutation.UpdateTime(); !ok {
		if serviceresource.DefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized serviceresource.DefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := serviceresource.DefaultUpdateTime()
		src.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (src *ServiceResourceCreate) check() error {
	if _, ok := src.mutation.ProjectID(); !ok {
		return &ValidationError{Name: "projectID", err: errors.New(`model: missing required field "ServiceResource.projectID"`)}
	}
	if v, ok := src.mutation.ProjectID(); ok {
		if err := serviceresource.ProjectIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "projectID", err: fmt.Errorf(`model: validator failed for field "ServiceResource.projectID": %w`, err)}
		}
	}
	if _, ok := src.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "ServiceResource.createTime"`)}
	}
	if _, ok := src.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "ServiceResource.updateTime"`)}
	}
	if _, ok := src.mutation.ServiceID(); !ok {
		return &ValidationError{Name: "serviceID", err: errors.New(`model: missing required field "ServiceResource.serviceID"`)}
	}
	if v, ok := src.mutation.ServiceID(); ok {
		if err := serviceresource.ServiceIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "serviceID", err: fmt.Errorf(`model: validator failed for field "ServiceResource.serviceID": %w`, err)}
		}
	}
	if _, ok := src.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connectorID", err: errors.New(`model: missing required field "ServiceResource.connectorID"`)}
	}
	if v, ok := src.mutation.ConnectorID(); ok {
		if err := serviceresource.ConnectorIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "connectorID", err: fmt.Errorf(`model: validator failed for field "ServiceResource.connectorID": %w`, err)}
		}
	}
	if _, ok := src.mutation.Mode(); !ok {
		return &ValidationError{Name: "mode", err: errors.New(`model: missing required field "ServiceResource.mode"`)}
	}
	if v, ok := src.mutation.Mode(); ok {
		if err := serviceresource.ModeValidator(v); err != nil {
			return &ValidationError{Name: "mode", err: fmt.Errorf(`model: validator failed for field "ServiceResource.mode": %w`, err)}
		}
	}
	if _, ok := src.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`model: missing required field "ServiceResource.type"`)}
	}
	if v, ok := src.mutation.GetType(); ok {
		if err := serviceresource.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`model: validator failed for field "ServiceResource.type": %w`, err)}
		}
	}
	if _, ok := src.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "ServiceResource.name"`)}
	}
	if v, ok := src.mutation.Name(); ok {
		if err := serviceresource.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "ServiceResource.name": %w`, err)}
		}
	}
	if _, ok := src.mutation.DeployerType(); !ok {
		return &ValidationError{Name: "deployerType", err: errors.New(`model: missing required field "ServiceResource.deployerType"`)}
	}
	if v, ok := src.mutation.DeployerType(); ok {
		if err := serviceresource.DeployerTypeValidator(v); err != nil {
			return &ValidationError{Name: "deployerType", err: fmt.Errorf(`model: validator failed for field "ServiceResource.deployerType": %w`, err)}
		}
	}
	if _, ok := src.mutation.ServiceID(); !ok {
		return &ValidationError{Name: "service", err: errors.New(`model: missing required edge "ServiceResource.service"`)}
	}
	if _, ok := src.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connector", err: errors.New(`model: missing required edge "ServiceResource.connector"`)}
	}
	return nil
}

func (src *ServiceResourceCreate) sqlSave(ctx context.Context) (*ServiceResource, error) {
	if err := src.check(); err != nil {
		return nil, err
	}
	_node, _spec := src.createSpec()
	if err := sqlgraph.CreateNode(ctx, src.driver, _spec); err != nil {
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
	src.mutation.id = &_node.ID
	src.mutation.done = true
	return _node, nil
}

func (src *ServiceResourceCreate) createSpec() (*ServiceResource, *sqlgraph.CreateSpec) {
	var (
		_node = &ServiceResource{config: src.config}
		_spec = sqlgraph.NewCreateSpec(serviceresource.Table, sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString))
	)
	_spec.Schema = src.schemaConfig.ServiceResource
	_spec.OnConflict = src.conflict
	if id, ok := src.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := src.mutation.ProjectID(); ok {
		_spec.SetField(serviceresource.FieldProjectID, field.TypeString, value)
		_node.ProjectID = value
	}
	if value, ok := src.mutation.CreateTime(); ok {
		_spec.SetField(serviceresource.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := src.mutation.UpdateTime(); ok {
		_spec.SetField(serviceresource.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := src.mutation.Mode(); ok {
		_spec.SetField(serviceresource.FieldMode, field.TypeString, value)
		_node.Mode = value
	}
	if value, ok := src.mutation.GetType(); ok {
		_spec.SetField(serviceresource.FieldType, field.TypeString, value)
		_node.Type = value
	}
	if value, ok := src.mutation.Name(); ok {
		_spec.SetField(serviceresource.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := src.mutation.DeployerType(); ok {
		_spec.SetField(serviceresource.FieldDeployerType, field.TypeString, value)
		_node.DeployerType = value
	}
	if value, ok := src.mutation.Status(); ok {
		_spec.SetField(serviceresource.FieldStatus, field.TypeJSON, value)
		_node.Status = value
	}
	if nodes := src.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceresource.ServiceTable,
			Columns: []string{serviceresource.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(service.FieldID, field.TypeString),
			},
		}
		edge.Schema = src.schemaConfig.ServiceResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ServiceID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := src.mutation.ConnectorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceresource.ConnectorTable,
			Columns: []string{serviceresource.ConnectorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(connector.FieldID, field.TypeString),
			},
		}
		edge.Schema = src.schemaConfig.ServiceResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ConnectorID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := src.mutation.CompositionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   serviceresource.CompositionTable,
			Columns: []string{serviceresource.CompositionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = src.schemaConfig.ServiceResource
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.CompositionID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := src.mutation.ComponentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   serviceresource.ComponentsTable,
			Columns: []string{serviceresource.ComponentsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString),
			},
		}
		edge.Schema = src.schemaConfig.ServiceResource
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
//	client.ServiceResource.Create().
//		SetProjectID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ServiceResourceUpsert) {
//			SetProjectID(v+v).
//		}).
//		Exec(ctx)
func (src *ServiceResourceCreate) OnConflict(opts ...sql.ConflictOption) *ServiceResourceUpsertOne {
	src.conflict = opts
	return &ServiceResourceUpsertOne{
		create: src,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ServiceResource.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (src *ServiceResourceCreate) OnConflictColumns(columns ...string) *ServiceResourceUpsertOne {
	src.conflict = append(src.conflict, sql.ConflictColumns(columns...))
	return &ServiceResourceUpsertOne{
		create: src,
	}
}

type (
	// ServiceResourceUpsertOne is the builder for "upsert"-ing
	//  one ServiceResource node.
	ServiceResourceUpsertOne struct {
		create *ServiceResourceCreate
	}

	// ServiceResourceUpsert is the "OnConflict" setter.
	ServiceResourceUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "updateTime" field.
func (u *ServiceResourceUpsert) SetUpdateTime(v time.Time) *ServiceResourceUpsert {
	u.Set(serviceresource.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ServiceResourceUpsert) UpdateUpdateTime() *ServiceResourceUpsert {
	u.SetExcluded(serviceresource.FieldUpdateTime)
	return u
}

// SetStatus sets the "status" field.
func (u *ServiceResourceUpsert) SetStatus(v types.ServiceResourceStatus) *ServiceResourceUpsert {
	u.Set(serviceresource.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ServiceResourceUpsert) UpdateStatus() *ServiceResourceUpsert {
	u.SetExcluded(serviceresource.FieldStatus)
	return u
}

// ClearStatus clears the value of the "status" field.
func (u *ServiceResourceUpsert) ClearStatus() *ServiceResourceUpsert {
	u.SetNull(serviceresource.FieldStatus)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.ServiceResource.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(serviceresource.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ServiceResourceUpsertOne) UpdateNewValues() *ServiceResourceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(serviceresource.FieldID)
		}
		if _, exists := u.create.mutation.ProjectID(); exists {
			s.SetIgnore(serviceresource.FieldProjectID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(serviceresource.FieldCreateTime)
		}
		if _, exists := u.create.mutation.ServiceID(); exists {
			s.SetIgnore(serviceresource.FieldServiceID)
		}
		if _, exists := u.create.mutation.ConnectorID(); exists {
			s.SetIgnore(serviceresource.FieldConnectorID)
		}
		if _, exists := u.create.mutation.CompositionID(); exists {
			s.SetIgnore(serviceresource.FieldCompositionID)
		}
		if _, exists := u.create.mutation.Mode(); exists {
			s.SetIgnore(serviceresource.FieldMode)
		}
		if _, exists := u.create.mutation.GetType(); exists {
			s.SetIgnore(serviceresource.FieldType)
		}
		if _, exists := u.create.mutation.Name(); exists {
			s.SetIgnore(serviceresource.FieldName)
		}
		if _, exists := u.create.mutation.DeployerType(); exists {
			s.SetIgnore(serviceresource.FieldDeployerType)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ServiceResource.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ServiceResourceUpsertOne) Ignore() *ServiceResourceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ServiceResourceUpsertOne) DoNothing() *ServiceResourceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ServiceResourceCreate.OnConflict
// documentation for more info.
func (u *ServiceResourceUpsertOne) Update(set func(*ServiceResourceUpsert)) *ServiceResourceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ServiceResourceUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ServiceResourceUpsertOne) SetUpdateTime(v time.Time) *ServiceResourceUpsertOne {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ServiceResourceUpsertOne) UpdateUpdateTime() *ServiceResourceUpsertOne {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetStatus sets the "status" field.
func (u *ServiceResourceUpsertOne) SetStatus(v types.ServiceResourceStatus) *ServiceResourceUpsertOne {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ServiceResourceUpsertOne) UpdateStatus() *ServiceResourceUpsertOne {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.UpdateStatus()
	})
}

// ClearStatus clears the value of the "status" field.
func (u *ServiceResourceUpsertOne) ClearStatus() *ServiceResourceUpsertOne {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.ClearStatus()
	})
}

// Exec executes the query.
func (u *ServiceResourceUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ServiceResourceCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ServiceResourceUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ServiceResourceUpsertOne) ID(ctx context.Context) (id oid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: ServiceResourceUpsertOne.ID is not supported by MySQL driver. Use ServiceResourceUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ServiceResourceUpsertOne) IDX(ctx context.Context) oid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ServiceResourceCreateBulk is the builder for creating many ServiceResource entities in bulk.
type ServiceResourceCreateBulk struct {
	config
	builders []*ServiceResourceCreate
	conflict []sql.ConflictOption
}

// Save creates the ServiceResource entities in the database.
func (srcb *ServiceResourceCreateBulk) Save(ctx context.Context) ([]*ServiceResource, error) {
	specs := make([]*sqlgraph.CreateSpec, len(srcb.builders))
	nodes := make([]*ServiceResource, len(srcb.builders))
	mutators := make([]Mutator, len(srcb.builders))
	for i := range srcb.builders {
		func(i int, root context.Context) {
			builder := srcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ServiceResourceMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, srcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = srcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, srcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, srcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (srcb *ServiceResourceCreateBulk) SaveX(ctx context.Context) []*ServiceResource {
	v, err := srcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (srcb *ServiceResourceCreateBulk) Exec(ctx context.Context) error {
	_, err := srcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (srcb *ServiceResourceCreateBulk) ExecX(ctx context.Context) {
	if err := srcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ServiceResource.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ServiceResourceUpsert) {
//			SetProjectID(v+v).
//		}).
//		Exec(ctx)
func (srcb *ServiceResourceCreateBulk) OnConflict(opts ...sql.ConflictOption) *ServiceResourceUpsertBulk {
	srcb.conflict = opts
	return &ServiceResourceUpsertBulk{
		create: srcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ServiceResource.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (srcb *ServiceResourceCreateBulk) OnConflictColumns(columns ...string) *ServiceResourceUpsertBulk {
	srcb.conflict = append(srcb.conflict, sql.ConflictColumns(columns...))
	return &ServiceResourceUpsertBulk{
		create: srcb,
	}
}

// ServiceResourceUpsertBulk is the builder for "upsert"-ing
// a bulk of ServiceResource nodes.
type ServiceResourceUpsertBulk struct {
	create *ServiceResourceCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ServiceResource.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(serviceresource.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ServiceResourceUpsertBulk) UpdateNewValues() *ServiceResourceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(serviceresource.FieldID)
			}
			if _, exists := b.mutation.ProjectID(); exists {
				s.SetIgnore(serviceresource.FieldProjectID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(serviceresource.FieldCreateTime)
			}
			if _, exists := b.mutation.ServiceID(); exists {
				s.SetIgnore(serviceresource.FieldServiceID)
			}
			if _, exists := b.mutation.ConnectorID(); exists {
				s.SetIgnore(serviceresource.FieldConnectorID)
			}
			if _, exists := b.mutation.CompositionID(); exists {
				s.SetIgnore(serviceresource.FieldCompositionID)
			}
			if _, exists := b.mutation.Mode(); exists {
				s.SetIgnore(serviceresource.FieldMode)
			}
			if _, exists := b.mutation.GetType(); exists {
				s.SetIgnore(serviceresource.FieldType)
			}
			if _, exists := b.mutation.Name(); exists {
				s.SetIgnore(serviceresource.FieldName)
			}
			if _, exists := b.mutation.DeployerType(); exists {
				s.SetIgnore(serviceresource.FieldDeployerType)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ServiceResource.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ServiceResourceUpsertBulk) Ignore() *ServiceResourceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ServiceResourceUpsertBulk) DoNothing() *ServiceResourceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ServiceResourceCreateBulk.OnConflict
// documentation for more info.
func (u *ServiceResourceUpsertBulk) Update(set func(*ServiceResourceUpsert)) *ServiceResourceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ServiceResourceUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ServiceResourceUpsertBulk) SetUpdateTime(v time.Time) *ServiceResourceUpsertBulk {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ServiceResourceUpsertBulk) UpdateUpdateTime() *ServiceResourceUpsertBulk {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetStatus sets the "status" field.
func (u *ServiceResourceUpsertBulk) SetStatus(v types.ServiceResourceStatus) *ServiceResourceUpsertBulk {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ServiceResourceUpsertBulk) UpdateStatus() *ServiceResourceUpsertBulk {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.UpdateStatus()
	})
}

// ClearStatus clears the value of the "status" field.
func (u *ServiceResourceUpsertBulk) ClearStatus() *ServiceResourceUpsertBulk {
	return u.Update(func(s *ServiceResourceUpsert) {
		s.ClearStatus()
	})
}

// Exec executes the query.
func (u *ServiceResourceUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ServiceResourceCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ServiceResourceCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ServiceResourceUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
