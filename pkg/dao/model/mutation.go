// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeApplication         = "Application"
	TypeApplicationResource = "ApplicationResource"
	TypeApplicationRevision = "ApplicationRevision"
	TypeConnector           = "Connector"
	TypeEnvironment         = "Environment"
	TypeModule              = "Module"
	TypeProject             = "Project"
	TypeRole                = "Role"
	TypeSetting             = "Setting"
	TypeSubject             = "Subject"
	TypeToken               = "Token"
)

// ApplicationMutation represents an operation that mutates the Application nodes in the graph.
type ApplicationMutation struct {
	config
	op            Op
	typ           string
	id            *oid.ID
	createTime    *time.Time
	updateTime    *time.Time
	projectID     *oid.ID
	environmentID *oid.ID
	modules       *[]schema.ApplicationModule
	appendmodules []schema.ApplicationModule
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Application, error)
	predicates    []predicate.Application
}

var _ ent.Mutation = (*ApplicationMutation)(nil)

// applicationOption allows management of the mutation configuration using functional options.
type applicationOption func(*ApplicationMutation)

// newApplicationMutation creates new mutation for the Application entity.
func newApplicationMutation(c config, op Op, opts ...applicationOption) *ApplicationMutation {
	m := &ApplicationMutation{
		config:        c,
		op:            op,
		typ:           TypeApplication,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withApplicationID sets the ID field of the mutation.
func withApplicationID(id oid.ID) applicationOption {
	return func(m *ApplicationMutation) {
		var (
			err   error
			once  sync.Once
			value *Application
		)
		m.oldValue = func(ctx context.Context) (*Application, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Application.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withApplication sets the old Application of the mutation.
func withApplication(node *Application) applicationOption {
	return func(m *ApplicationMutation) {
		m.oldValue = func(context.Context) (*Application, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ApplicationMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ApplicationMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Application entities.
func (m *ApplicationMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ApplicationMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ApplicationMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Application.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *ApplicationMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *ApplicationMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *ApplicationMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *ApplicationMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *ApplicationMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *ApplicationMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetProjectID sets the "projectID" field.
func (m *ApplicationMutation) SetProjectID(o oid.ID) {
	m.projectID = &o
}

// ProjectID returns the value of the "projectID" field in the mutation.
func (m *ApplicationMutation) ProjectID() (r oid.ID, exists bool) {
	v := m.projectID
	if v == nil {
		return
	}
	return *v, true
}

// OldProjectID returns the old "projectID" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldProjectID(ctx context.Context) (v oid.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldProjectID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldProjectID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldProjectID: %w", err)
	}
	return oldValue.ProjectID, nil
}

// ResetProjectID resets all changes to the "projectID" field.
func (m *ApplicationMutation) ResetProjectID() {
	m.projectID = nil
}

// SetEnvironmentID sets the "environmentID" field.
func (m *ApplicationMutation) SetEnvironmentID(o oid.ID) {
	m.environmentID = &o
}

// EnvironmentID returns the value of the "environmentID" field in the mutation.
func (m *ApplicationMutation) EnvironmentID() (r oid.ID, exists bool) {
	v := m.environmentID
	if v == nil {
		return
	}
	return *v, true
}

// OldEnvironmentID returns the old "environmentID" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldEnvironmentID(ctx context.Context) (v oid.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEnvironmentID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEnvironmentID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEnvironmentID: %w", err)
	}
	return oldValue.EnvironmentID, nil
}

// ResetEnvironmentID resets all changes to the "environmentID" field.
func (m *ApplicationMutation) ResetEnvironmentID() {
	m.environmentID = nil
}

// SetModules sets the "modules" field.
func (m *ApplicationMutation) SetModules(sm []schema.ApplicationModule) {
	m.modules = &sm
	m.appendmodules = nil
}

// Modules returns the value of the "modules" field in the mutation.
func (m *ApplicationMutation) Modules() (r []schema.ApplicationModule, exists bool) {
	v := m.modules
	if v == nil {
		return
	}
	return *v, true
}

// OldModules returns the old "modules" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldModules(ctx context.Context) (v []schema.ApplicationModule, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldModules is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldModules requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldModules: %w", err)
	}
	return oldValue.Modules, nil
}

// AppendModules adds sm to the "modules" field.
func (m *ApplicationMutation) AppendModules(sm []schema.ApplicationModule) {
	m.appendmodules = append(m.appendmodules, sm...)
}

// AppendedModules returns the list of values that were appended to the "modules" field in this mutation.
func (m *ApplicationMutation) AppendedModules() ([]schema.ApplicationModule, bool) {
	if len(m.appendmodules) == 0 {
		return nil, false
	}
	return m.appendmodules, true
}

// ResetModules resets all changes to the "modules" field.
func (m *ApplicationMutation) ResetModules() {
	m.modules = nil
	m.appendmodules = nil
}

// Where appends a list predicates to the ApplicationMutation builder.
func (m *ApplicationMutation) Where(ps ...predicate.Application) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ApplicationMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ApplicationMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Application, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ApplicationMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ApplicationMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Application).
func (m *ApplicationMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ApplicationMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.createTime != nil {
		fields = append(fields, application.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, application.FieldUpdateTime)
	}
	if m.projectID != nil {
		fields = append(fields, application.FieldProjectID)
	}
	if m.environmentID != nil {
		fields = append(fields, application.FieldEnvironmentID)
	}
	if m.modules != nil {
		fields = append(fields, application.FieldModules)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ApplicationMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case application.FieldCreateTime:
		return m.CreateTime()
	case application.FieldUpdateTime:
		return m.UpdateTime()
	case application.FieldProjectID:
		return m.ProjectID()
	case application.FieldEnvironmentID:
		return m.EnvironmentID()
	case application.FieldModules:
		return m.Modules()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ApplicationMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case application.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case application.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case application.FieldProjectID:
		return m.OldProjectID(ctx)
	case application.FieldEnvironmentID:
		return m.OldEnvironmentID(ctx)
	case application.FieldModules:
		return m.OldModules(ctx)
	}
	return nil, fmt.Errorf("unknown Application field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationMutation) SetField(name string, value ent.Value) error {
	switch name {
	case application.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case application.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case application.FieldProjectID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetProjectID(v)
		return nil
	case application.FieldEnvironmentID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnvironmentID(v)
		return nil
	case application.FieldModules:
		v, ok := value.([]schema.ApplicationModule)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetModules(v)
		return nil
	}
	return fmt.Errorf("unknown Application field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ApplicationMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ApplicationMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Application numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ApplicationMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ApplicationMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ApplicationMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Application nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ApplicationMutation) ResetField(name string) error {
	switch name {
	case application.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case application.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case application.FieldProjectID:
		m.ResetProjectID()
		return nil
	case application.FieldEnvironmentID:
		m.ResetEnvironmentID()
		return nil
	case application.FieldModules:
		m.ResetModules()
		return nil
	}
	return fmt.Errorf("unknown Application field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ApplicationMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Application unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ApplicationMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Application edge %s", name)
}

// ApplicationResourceMutation represents an operation that mutates the ApplicationResource nodes in the graph.
type ApplicationResourceMutation struct {
	config
	op            Op
	typ           string
	id            *oid.ID
	status        *string
	statusMessage *string
	createTime    *time.Time
	updateTime    *time.Time
	applicationID *oid.ID
	module        *string
	_type         *string
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*ApplicationResource, error)
	predicates    []predicate.ApplicationResource
}

var _ ent.Mutation = (*ApplicationResourceMutation)(nil)

// applicationresourceOption allows management of the mutation configuration using functional options.
type applicationresourceOption func(*ApplicationResourceMutation)

// newApplicationResourceMutation creates new mutation for the ApplicationResource entity.
func newApplicationResourceMutation(c config, op Op, opts ...applicationresourceOption) *ApplicationResourceMutation {
	m := &ApplicationResourceMutation{
		config:        c,
		op:            op,
		typ:           TypeApplicationResource,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withApplicationResourceID sets the ID field of the mutation.
func withApplicationResourceID(id oid.ID) applicationresourceOption {
	return func(m *ApplicationResourceMutation) {
		var (
			err   error
			once  sync.Once
			value *ApplicationResource
		)
		m.oldValue = func(ctx context.Context) (*ApplicationResource, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().ApplicationResource.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withApplicationResource sets the old ApplicationResource of the mutation.
func withApplicationResource(node *ApplicationResource) applicationresourceOption {
	return func(m *ApplicationResourceMutation) {
		m.oldValue = func(context.Context) (*ApplicationResource, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ApplicationResourceMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ApplicationResourceMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of ApplicationResource entities.
func (m *ApplicationResourceMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ApplicationResourceMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ApplicationResourceMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().ApplicationResource.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetStatus sets the "status" field.
func (m *ApplicationResourceMutation) SetStatus(s string) {
	m.status = &s
}

// Status returns the value of the "status" field in the mutation.
func (m *ApplicationResourceMutation) Status() (r string, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldStatus(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatus is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatus requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatus: %w", err)
	}
	return oldValue.Status, nil
}

// ClearStatus clears the value of the "status" field.
func (m *ApplicationResourceMutation) ClearStatus() {
	m.status = nil
	m.clearedFields[applicationresource.FieldStatus] = struct{}{}
}

// StatusCleared returns if the "status" field was cleared in this mutation.
func (m *ApplicationResourceMutation) StatusCleared() bool {
	_, ok := m.clearedFields[applicationresource.FieldStatus]
	return ok
}

// ResetStatus resets all changes to the "status" field.
func (m *ApplicationResourceMutation) ResetStatus() {
	m.status = nil
	delete(m.clearedFields, applicationresource.FieldStatus)
}

// SetStatusMessage sets the "statusMessage" field.
func (m *ApplicationResourceMutation) SetStatusMessage(s string) {
	m.statusMessage = &s
}

// StatusMessage returns the value of the "statusMessage" field in the mutation.
func (m *ApplicationResourceMutation) StatusMessage() (r string, exists bool) {
	v := m.statusMessage
	if v == nil {
		return
	}
	return *v, true
}

// OldStatusMessage returns the old "statusMessage" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldStatusMessage(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatusMessage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatusMessage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatusMessage: %w", err)
	}
	return oldValue.StatusMessage, nil
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (m *ApplicationResourceMutation) ClearStatusMessage() {
	m.statusMessage = nil
	m.clearedFields[applicationresource.FieldStatusMessage] = struct{}{}
}

// StatusMessageCleared returns if the "statusMessage" field was cleared in this mutation.
func (m *ApplicationResourceMutation) StatusMessageCleared() bool {
	_, ok := m.clearedFields[applicationresource.FieldStatusMessage]
	return ok
}

// ResetStatusMessage resets all changes to the "statusMessage" field.
func (m *ApplicationResourceMutation) ResetStatusMessage() {
	m.statusMessage = nil
	delete(m.clearedFields, applicationresource.FieldStatusMessage)
}

// SetCreateTime sets the "createTime" field.
func (m *ApplicationResourceMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *ApplicationResourceMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *ApplicationResourceMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *ApplicationResourceMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *ApplicationResourceMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *ApplicationResourceMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetApplicationID sets the "applicationID" field.
func (m *ApplicationResourceMutation) SetApplicationID(o oid.ID) {
	m.applicationID = &o
}

// ApplicationID returns the value of the "applicationID" field in the mutation.
func (m *ApplicationResourceMutation) ApplicationID() (r oid.ID, exists bool) {
	v := m.applicationID
	if v == nil {
		return
	}
	return *v, true
}

// OldApplicationID returns the old "applicationID" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldApplicationID(ctx context.Context) (v oid.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldApplicationID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldApplicationID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldApplicationID: %w", err)
	}
	return oldValue.ApplicationID, nil
}

// ResetApplicationID resets all changes to the "applicationID" field.
func (m *ApplicationResourceMutation) ResetApplicationID() {
	m.applicationID = nil
}

// SetModule sets the "module" field.
func (m *ApplicationResourceMutation) SetModule(s string) {
	m.module = &s
}

// Module returns the value of the "module" field in the mutation.
func (m *ApplicationResourceMutation) Module() (r string, exists bool) {
	v := m.module
	if v == nil {
		return
	}
	return *v, true
}

// OldModule returns the old "module" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldModule(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldModule is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldModule requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldModule: %w", err)
	}
	return oldValue.Module, nil
}

// ResetModule resets all changes to the "module" field.
func (m *ApplicationResourceMutation) ResetModule() {
	m.module = nil
}

// SetType sets the "type" field.
func (m *ApplicationResourceMutation) SetType(s string) {
	m._type = &s
}

// GetType returns the value of the "type" field in the mutation.
func (m *ApplicationResourceMutation) GetType() (r string, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// OldType returns the old "type" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldType(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldType is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldType requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldType: %w", err)
	}
	return oldValue.Type, nil
}

// ResetType resets all changes to the "type" field.
func (m *ApplicationResourceMutation) ResetType() {
	m._type = nil
}

// Where appends a list predicates to the ApplicationResourceMutation builder.
func (m *ApplicationResourceMutation) Where(ps ...predicate.ApplicationResource) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ApplicationResourceMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ApplicationResourceMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.ApplicationResource, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ApplicationResourceMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ApplicationResourceMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (ApplicationResource).
func (m *ApplicationResourceMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ApplicationResourceMutation) Fields() []string {
	fields := make([]string, 0, 7)
	if m.status != nil {
		fields = append(fields, applicationresource.FieldStatus)
	}
	if m.statusMessage != nil {
		fields = append(fields, applicationresource.FieldStatusMessage)
	}
	if m.createTime != nil {
		fields = append(fields, applicationresource.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, applicationresource.FieldUpdateTime)
	}
	if m.applicationID != nil {
		fields = append(fields, applicationresource.FieldApplicationID)
	}
	if m.module != nil {
		fields = append(fields, applicationresource.FieldModule)
	}
	if m._type != nil {
		fields = append(fields, applicationresource.FieldType)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ApplicationResourceMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case applicationresource.FieldStatus:
		return m.Status()
	case applicationresource.FieldStatusMessage:
		return m.StatusMessage()
	case applicationresource.FieldCreateTime:
		return m.CreateTime()
	case applicationresource.FieldUpdateTime:
		return m.UpdateTime()
	case applicationresource.FieldApplicationID:
		return m.ApplicationID()
	case applicationresource.FieldModule:
		return m.Module()
	case applicationresource.FieldType:
		return m.GetType()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ApplicationResourceMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case applicationresource.FieldStatus:
		return m.OldStatus(ctx)
	case applicationresource.FieldStatusMessage:
		return m.OldStatusMessage(ctx)
	case applicationresource.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case applicationresource.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case applicationresource.FieldApplicationID:
		return m.OldApplicationID(ctx)
	case applicationresource.FieldModule:
		return m.OldModule(ctx)
	case applicationresource.FieldType:
		return m.OldType(ctx)
	}
	return nil, fmt.Errorf("unknown ApplicationResource field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationResourceMutation) SetField(name string, value ent.Value) error {
	switch name {
	case applicationresource.FieldStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case applicationresource.FieldStatusMessage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatusMessage(v)
		return nil
	case applicationresource.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case applicationresource.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case applicationresource.FieldApplicationID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetApplicationID(v)
		return nil
	case applicationresource.FieldModule:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetModule(v)
		return nil
	case applicationresource.FieldType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ApplicationResourceMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ApplicationResourceMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationResourceMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ApplicationResource numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ApplicationResourceMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(applicationresource.FieldStatus) {
		fields = append(fields, applicationresource.FieldStatus)
	}
	if m.FieldCleared(applicationresource.FieldStatusMessage) {
		fields = append(fields, applicationresource.FieldStatusMessage)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ApplicationResourceMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ApplicationResourceMutation) ClearField(name string) error {
	switch name {
	case applicationresource.FieldStatus:
		m.ClearStatus()
		return nil
	case applicationresource.FieldStatusMessage:
		m.ClearStatusMessage()
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ApplicationResourceMutation) ResetField(name string) error {
	switch name {
	case applicationresource.FieldStatus:
		m.ResetStatus()
		return nil
	case applicationresource.FieldStatusMessage:
		m.ResetStatusMessage()
		return nil
	case applicationresource.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case applicationresource.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case applicationresource.FieldApplicationID:
		m.ResetApplicationID()
		return nil
	case applicationresource.FieldModule:
		m.ResetModule()
		return nil
	case applicationresource.FieldType:
		m.ResetType()
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationResourceMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ApplicationResourceMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationResourceMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationResourceMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationResourceMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationResourceMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationResourceMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown ApplicationResource unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ApplicationResourceMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown ApplicationResource edge %s", name)
}

// ApplicationRevisionMutation represents an operation that mutates the ApplicationRevision nodes in the graph.
type ApplicationRevisionMutation struct {
	config
	op             Op
	typ            string
	id             *oid.ID
	status         *string
	statusMessage  *string
	createTime     *time.Time
	updateTime     *time.Time
	applicationID  *oid.ID
	environmentID  *oid.ID
	modules        *[]schema.ApplicationModule
	appendmodules  []schema.ApplicationModule
	inputVariables *map[string]interface{}
	inputPlan      *string
	output         *string
	clearedFields  map[string]struct{}
	done           bool
	oldValue       func(context.Context) (*ApplicationRevision, error)
	predicates     []predicate.ApplicationRevision
}

var _ ent.Mutation = (*ApplicationRevisionMutation)(nil)

// applicationrevisionOption allows management of the mutation configuration using functional options.
type applicationrevisionOption func(*ApplicationRevisionMutation)

// newApplicationRevisionMutation creates new mutation for the ApplicationRevision entity.
func newApplicationRevisionMutation(c config, op Op, opts ...applicationrevisionOption) *ApplicationRevisionMutation {
	m := &ApplicationRevisionMutation{
		config:        c,
		op:            op,
		typ:           TypeApplicationRevision,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withApplicationRevisionID sets the ID field of the mutation.
func withApplicationRevisionID(id oid.ID) applicationrevisionOption {
	return func(m *ApplicationRevisionMutation) {
		var (
			err   error
			once  sync.Once
			value *ApplicationRevision
		)
		m.oldValue = func(ctx context.Context) (*ApplicationRevision, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().ApplicationRevision.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withApplicationRevision sets the old ApplicationRevision of the mutation.
func withApplicationRevision(node *ApplicationRevision) applicationrevisionOption {
	return func(m *ApplicationRevisionMutation) {
		m.oldValue = func(context.Context) (*ApplicationRevision, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ApplicationRevisionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ApplicationRevisionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of ApplicationRevision entities.
func (m *ApplicationRevisionMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ApplicationRevisionMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ApplicationRevisionMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().ApplicationRevision.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetStatus sets the "status" field.
func (m *ApplicationRevisionMutation) SetStatus(s string) {
	m.status = &s
}

// Status returns the value of the "status" field in the mutation.
func (m *ApplicationRevisionMutation) Status() (r string, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldStatus(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatus is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatus requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatus: %w", err)
	}
	return oldValue.Status, nil
}

// ClearStatus clears the value of the "status" field.
func (m *ApplicationRevisionMutation) ClearStatus() {
	m.status = nil
	m.clearedFields[applicationrevision.FieldStatus] = struct{}{}
}

// StatusCleared returns if the "status" field was cleared in this mutation.
func (m *ApplicationRevisionMutation) StatusCleared() bool {
	_, ok := m.clearedFields[applicationrevision.FieldStatus]
	return ok
}

// ResetStatus resets all changes to the "status" field.
func (m *ApplicationRevisionMutation) ResetStatus() {
	m.status = nil
	delete(m.clearedFields, applicationrevision.FieldStatus)
}

// SetStatusMessage sets the "statusMessage" field.
func (m *ApplicationRevisionMutation) SetStatusMessage(s string) {
	m.statusMessage = &s
}

// StatusMessage returns the value of the "statusMessage" field in the mutation.
func (m *ApplicationRevisionMutation) StatusMessage() (r string, exists bool) {
	v := m.statusMessage
	if v == nil {
		return
	}
	return *v, true
}

// OldStatusMessage returns the old "statusMessage" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldStatusMessage(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatusMessage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatusMessage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatusMessage: %w", err)
	}
	return oldValue.StatusMessage, nil
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (m *ApplicationRevisionMutation) ClearStatusMessage() {
	m.statusMessage = nil
	m.clearedFields[applicationrevision.FieldStatusMessage] = struct{}{}
}

// StatusMessageCleared returns if the "statusMessage" field was cleared in this mutation.
func (m *ApplicationRevisionMutation) StatusMessageCleared() bool {
	_, ok := m.clearedFields[applicationrevision.FieldStatusMessage]
	return ok
}

// ResetStatusMessage resets all changes to the "statusMessage" field.
func (m *ApplicationRevisionMutation) ResetStatusMessage() {
	m.statusMessage = nil
	delete(m.clearedFields, applicationrevision.FieldStatusMessage)
}

// SetCreateTime sets the "createTime" field.
func (m *ApplicationRevisionMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *ApplicationRevisionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *ApplicationRevisionMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *ApplicationRevisionMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *ApplicationRevisionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *ApplicationRevisionMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetApplicationID sets the "applicationID" field.
func (m *ApplicationRevisionMutation) SetApplicationID(o oid.ID) {
	m.applicationID = &o
}

// ApplicationID returns the value of the "applicationID" field in the mutation.
func (m *ApplicationRevisionMutation) ApplicationID() (r oid.ID, exists bool) {
	v := m.applicationID
	if v == nil {
		return
	}
	return *v, true
}

// OldApplicationID returns the old "applicationID" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldApplicationID(ctx context.Context) (v oid.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldApplicationID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldApplicationID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldApplicationID: %w", err)
	}
	return oldValue.ApplicationID, nil
}

// ResetApplicationID resets all changes to the "applicationID" field.
func (m *ApplicationRevisionMutation) ResetApplicationID() {
	m.applicationID = nil
}

// SetEnvironmentID sets the "environmentID" field.
func (m *ApplicationRevisionMutation) SetEnvironmentID(o oid.ID) {
	m.environmentID = &o
}

// EnvironmentID returns the value of the "environmentID" field in the mutation.
func (m *ApplicationRevisionMutation) EnvironmentID() (r oid.ID, exists bool) {
	v := m.environmentID
	if v == nil {
		return
	}
	return *v, true
}

// OldEnvironmentID returns the old "environmentID" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldEnvironmentID(ctx context.Context) (v oid.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEnvironmentID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEnvironmentID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEnvironmentID: %w", err)
	}
	return oldValue.EnvironmentID, nil
}

// ResetEnvironmentID resets all changes to the "environmentID" field.
func (m *ApplicationRevisionMutation) ResetEnvironmentID() {
	m.environmentID = nil
}

// SetModules sets the "modules" field.
func (m *ApplicationRevisionMutation) SetModules(sm []schema.ApplicationModule) {
	m.modules = &sm
	m.appendmodules = nil
}

// Modules returns the value of the "modules" field in the mutation.
func (m *ApplicationRevisionMutation) Modules() (r []schema.ApplicationModule, exists bool) {
	v := m.modules
	if v == nil {
		return
	}
	return *v, true
}

// OldModules returns the old "modules" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldModules(ctx context.Context) (v []schema.ApplicationModule, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldModules is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldModules requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldModules: %w", err)
	}
	return oldValue.Modules, nil
}

// AppendModules adds sm to the "modules" field.
func (m *ApplicationRevisionMutation) AppendModules(sm []schema.ApplicationModule) {
	m.appendmodules = append(m.appendmodules, sm...)
}

// AppendedModules returns the list of values that were appended to the "modules" field in this mutation.
func (m *ApplicationRevisionMutation) AppendedModules() ([]schema.ApplicationModule, bool) {
	if len(m.appendmodules) == 0 {
		return nil, false
	}
	return m.appendmodules, true
}

// ResetModules resets all changes to the "modules" field.
func (m *ApplicationRevisionMutation) ResetModules() {
	m.modules = nil
	m.appendmodules = nil
}

// SetInputVariables sets the "inputVariables" field.
func (m *ApplicationRevisionMutation) SetInputVariables(value map[string]interface{}) {
	m.inputVariables = &value
}

// InputVariables returns the value of the "inputVariables" field in the mutation.
func (m *ApplicationRevisionMutation) InputVariables() (r map[string]interface{}, exists bool) {
	v := m.inputVariables
	if v == nil {
		return
	}
	return *v, true
}

// OldInputVariables returns the old "inputVariables" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldInputVariables(ctx context.Context) (v map[string]interface{}, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldInputVariables is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldInputVariables requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldInputVariables: %w", err)
	}
	return oldValue.InputVariables, nil
}

// ResetInputVariables resets all changes to the "inputVariables" field.
func (m *ApplicationRevisionMutation) ResetInputVariables() {
	m.inputVariables = nil
}

// SetInputPlan sets the "inputPlan" field.
func (m *ApplicationRevisionMutation) SetInputPlan(s string) {
	m.inputPlan = &s
}

// InputPlan returns the value of the "inputPlan" field in the mutation.
func (m *ApplicationRevisionMutation) InputPlan() (r string, exists bool) {
	v := m.inputPlan
	if v == nil {
		return
	}
	return *v, true
}

// OldInputPlan returns the old "inputPlan" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldInputPlan(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldInputPlan is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldInputPlan requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldInputPlan: %w", err)
	}
	return oldValue.InputPlan, nil
}

// ResetInputPlan resets all changes to the "inputPlan" field.
func (m *ApplicationRevisionMutation) ResetInputPlan() {
	m.inputPlan = nil
}

// SetOutput sets the "output" field.
func (m *ApplicationRevisionMutation) SetOutput(s string) {
	m.output = &s
}

// Output returns the value of the "output" field in the mutation.
func (m *ApplicationRevisionMutation) Output() (r string, exists bool) {
	v := m.output
	if v == nil {
		return
	}
	return *v, true
}

// OldOutput returns the old "output" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldOutput(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldOutput is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldOutput requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldOutput: %w", err)
	}
	return oldValue.Output, nil
}

// ResetOutput resets all changes to the "output" field.
func (m *ApplicationRevisionMutation) ResetOutput() {
	m.output = nil
}

// Where appends a list predicates to the ApplicationRevisionMutation builder.
func (m *ApplicationRevisionMutation) Where(ps ...predicate.ApplicationRevision) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ApplicationRevisionMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ApplicationRevisionMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.ApplicationRevision, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ApplicationRevisionMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ApplicationRevisionMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (ApplicationRevision).
func (m *ApplicationRevisionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ApplicationRevisionMutation) Fields() []string {
	fields := make([]string, 0, 10)
	if m.status != nil {
		fields = append(fields, applicationrevision.FieldStatus)
	}
	if m.statusMessage != nil {
		fields = append(fields, applicationrevision.FieldStatusMessage)
	}
	if m.createTime != nil {
		fields = append(fields, applicationrevision.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, applicationrevision.FieldUpdateTime)
	}
	if m.applicationID != nil {
		fields = append(fields, applicationrevision.FieldApplicationID)
	}
	if m.environmentID != nil {
		fields = append(fields, applicationrevision.FieldEnvironmentID)
	}
	if m.modules != nil {
		fields = append(fields, applicationrevision.FieldModules)
	}
	if m.inputVariables != nil {
		fields = append(fields, applicationrevision.FieldInputVariables)
	}
	if m.inputPlan != nil {
		fields = append(fields, applicationrevision.FieldInputPlan)
	}
	if m.output != nil {
		fields = append(fields, applicationrevision.FieldOutput)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ApplicationRevisionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case applicationrevision.FieldStatus:
		return m.Status()
	case applicationrevision.FieldStatusMessage:
		return m.StatusMessage()
	case applicationrevision.FieldCreateTime:
		return m.CreateTime()
	case applicationrevision.FieldUpdateTime:
		return m.UpdateTime()
	case applicationrevision.FieldApplicationID:
		return m.ApplicationID()
	case applicationrevision.FieldEnvironmentID:
		return m.EnvironmentID()
	case applicationrevision.FieldModules:
		return m.Modules()
	case applicationrevision.FieldInputVariables:
		return m.InputVariables()
	case applicationrevision.FieldInputPlan:
		return m.InputPlan()
	case applicationrevision.FieldOutput:
		return m.Output()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ApplicationRevisionMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case applicationrevision.FieldStatus:
		return m.OldStatus(ctx)
	case applicationrevision.FieldStatusMessage:
		return m.OldStatusMessage(ctx)
	case applicationrevision.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case applicationrevision.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case applicationrevision.FieldApplicationID:
		return m.OldApplicationID(ctx)
	case applicationrevision.FieldEnvironmentID:
		return m.OldEnvironmentID(ctx)
	case applicationrevision.FieldModules:
		return m.OldModules(ctx)
	case applicationrevision.FieldInputVariables:
		return m.OldInputVariables(ctx)
	case applicationrevision.FieldInputPlan:
		return m.OldInputPlan(ctx)
	case applicationrevision.FieldOutput:
		return m.OldOutput(ctx)
	}
	return nil, fmt.Errorf("unknown ApplicationRevision field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationRevisionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case applicationrevision.FieldStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case applicationrevision.FieldStatusMessage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatusMessage(v)
		return nil
	case applicationrevision.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case applicationrevision.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case applicationrevision.FieldApplicationID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetApplicationID(v)
		return nil
	case applicationrevision.FieldEnvironmentID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnvironmentID(v)
		return nil
	case applicationrevision.FieldModules:
		v, ok := value.([]schema.ApplicationModule)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetModules(v)
		return nil
	case applicationrevision.FieldInputVariables:
		v, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetInputVariables(v)
		return nil
	case applicationrevision.FieldInputPlan:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetInputPlan(v)
		return nil
	case applicationrevision.FieldOutput:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetOutput(v)
		return nil
	}
	return fmt.Errorf("unknown ApplicationRevision field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ApplicationRevisionMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ApplicationRevisionMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationRevisionMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ApplicationRevision numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ApplicationRevisionMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(applicationrevision.FieldStatus) {
		fields = append(fields, applicationrevision.FieldStatus)
	}
	if m.FieldCleared(applicationrevision.FieldStatusMessage) {
		fields = append(fields, applicationrevision.FieldStatusMessage)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ApplicationRevisionMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ApplicationRevisionMutation) ClearField(name string) error {
	switch name {
	case applicationrevision.FieldStatus:
		m.ClearStatus()
		return nil
	case applicationrevision.FieldStatusMessage:
		m.ClearStatusMessage()
		return nil
	}
	return fmt.Errorf("unknown ApplicationRevision nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ApplicationRevisionMutation) ResetField(name string) error {
	switch name {
	case applicationrevision.FieldStatus:
		m.ResetStatus()
		return nil
	case applicationrevision.FieldStatusMessage:
		m.ResetStatusMessage()
		return nil
	case applicationrevision.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case applicationrevision.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case applicationrevision.FieldApplicationID:
		m.ResetApplicationID()
		return nil
	case applicationrevision.FieldEnvironmentID:
		m.ResetEnvironmentID()
		return nil
	case applicationrevision.FieldModules:
		m.ResetModules()
		return nil
	case applicationrevision.FieldInputVariables:
		m.ResetInputVariables()
		return nil
	case applicationrevision.FieldInputPlan:
		m.ResetInputPlan()
		return nil
	case applicationrevision.FieldOutput:
		m.ResetOutput()
		return nil
	}
	return fmt.Errorf("unknown ApplicationRevision field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationRevisionMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ApplicationRevisionMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationRevisionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationRevisionMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationRevisionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationRevisionMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationRevisionMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown ApplicationRevision unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ApplicationRevisionMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown ApplicationRevision edge %s", name)
}

// ConnectorMutation represents an operation that mutates the Connector nodes in the graph.
type ConnectorMutation struct {
	config
	op            Op
	typ           string
	id            *oid.ID
	status        *string
	statusMessage *string
	createTime    *time.Time
	updateTime    *time.Time
	_driver       *string
	configVersion *string
	configData    *map[string]interface{}
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Connector, error)
	predicates    []predicate.Connector
}

var _ ent.Mutation = (*ConnectorMutation)(nil)

// connectorOption allows management of the mutation configuration using functional options.
type connectorOption func(*ConnectorMutation)

// newConnectorMutation creates new mutation for the Connector entity.
func newConnectorMutation(c config, op Op, opts ...connectorOption) *ConnectorMutation {
	m := &ConnectorMutation{
		config:        c,
		op:            op,
		typ:           TypeConnector,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withConnectorID sets the ID field of the mutation.
func withConnectorID(id oid.ID) connectorOption {
	return func(m *ConnectorMutation) {
		var (
			err   error
			once  sync.Once
			value *Connector
		)
		m.oldValue = func(ctx context.Context) (*Connector, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Connector.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withConnector sets the old Connector of the mutation.
func withConnector(node *Connector) connectorOption {
	return func(m *ConnectorMutation) {
		m.oldValue = func(context.Context) (*Connector, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ConnectorMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ConnectorMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Connector entities.
func (m *ConnectorMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ConnectorMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ConnectorMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Connector.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetStatus sets the "status" field.
func (m *ConnectorMutation) SetStatus(s string) {
	m.status = &s
}

// Status returns the value of the "status" field in the mutation.
func (m *ConnectorMutation) Status() (r string, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldStatus(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatus is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatus requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatus: %w", err)
	}
	return oldValue.Status, nil
}

// ClearStatus clears the value of the "status" field.
func (m *ConnectorMutation) ClearStatus() {
	m.status = nil
	m.clearedFields[connector.FieldStatus] = struct{}{}
}

// StatusCleared returns if the "status" field was cleared in this mutation.
func (m *ConnectorMutation) StatusCleared() bool {
	_, ok := m.clearedFields[connector.FieldStatus]
	return ok
}

// ResetStatus resets all changes to the "status" field.
func (m *ConnectorMutation) ResetStatus() {
	m.status = nil
	delete(m.clearedFields, connector.FieldStatus)
}

// SetStatusMessage sets the "statusMessage" field.
func (m *ConnectorMutation) SetStatusMessage(s string) {
	m.statusMessage = &s
}

// StatusMessage returns the value of the "statusMessage" field in the mutation.
func (m *ConnectorMutation) StatusMessage() (r string, exists bool) {
	v := m.statusMessage
	if v == nil {
		return
	}
	return *v, true
}

// OldStatusMessage returns the old "statusMessage" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldStatusMessage(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatusMessage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatusMessage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatusMessage: %w", err)
	}
	return oldValue.StatusMessage, nil
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (m *ConnectorMutation) ClearStatusMessage() {
	m.statusMessage = nil
	m.clearedFields[connector.FieldStatusMessage] = struct{}{}
}

// StatusMessageCleared returns if the "statusMessage" field was cleared in this mutation.
func (m *ConnectorMutation) StatusMessageCleared() bool {
	_, ok := m.clearedFields[connector.FieldStatusMessage]
	return ok
}

// ResetStatusMessage resets all changes to the "statusMessage" field.
func (m *ConnectorMutation) ResetStatusMessage() {
	m.statusMessage = nil
	delete(m.clearedFields, connector.FieldStatusMessage)
}

// SetCreateTime sets the "createTime" field.
func (m *ConnectorMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *ConnectorMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *ConnectorMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *ConnectorMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *ConnectorMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *ConnectorMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetDriver sets the "driver" field.
func (m *ConnectorMutation) SetDriver(s string) {
	m._driver = &s
}

// Driver returns the value of the "driver" field in the mutation.
func (m *ConnectorMutation) Driver() (r string, exists bool) {
	v := m._driver
	if v == nil {
		return
	}
	return *v, true
}

// OldDriver returns the old "driver" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldDriver(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDriver is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDriver requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDriver: %w", err)
	}
	return oldValue.Driver, nil
}

// ResetDriver resets all changes to the "driver" field.
func (m *ConnectorMutation) ResetDriver() {
	m._driver = nil
}

// SetConfigVersion sets the "configVersion" field.
func (m *ConnectorMutation) SetConfigVersion(s string) {
	m.configVersion = &s
}

// ConfigVersion returns the value of the "configVersion" field in the mutation.
func (m *ConnectorMutation) ConfigVersion() (r string, exists bool) {
	v := m.configVersion
	if v == nil {
		return
	}
	return *v, true
}

// OldConfigVersion returns the old "configVersion" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldConfigVersion(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldConfigVersion is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldConfigVersion requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldConfigVersion: %w", err)
	}
	return oldValue.ConfigVersion, nil
}

// ResetConfigVersion resets all changes to the "configVersion" field.
func (m *ConnectorMutation) ResetConfigVersion() {
	m.configVersion = nil
}

// SetConfigData sets the "configData" field.
func (m *ConnectorMutation) SetConfigData(value map[string]interface{}) {
	m.configData = &value
}

// ConfigData returns the value of the "configData" field in the mutation.
func (m *ConnectorMutation) ConfigData() (r map[string]interface{}, exists bool) {
	v := m.configData
	if v == nil {
		return
	}
	return *v, true
}

// OldConfigData returns the old "configData" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldConfigData(ctx context.Context) (v map[string]interface{}, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldConfigData is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldConfigData requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldConfigData: %w", err)
	}
	return oldValue.ConfigData, nil
}

// ClearConfigData clears the value of the "configData" field.
func (m *ConnectorMutation) ClearConfigData() {
	m.configData = nil
	m.clearedFields[connector.FieldConfigData] = struct{}{}
}

// ConfigDataCleared returns if the "configData" field was cleared in this mutation.
func (m *ConnectorMutation) ConfigDataCleared() bool {
	_, ok := m.clearedFields[connector.FieldConfigData]
	return ok
}

// ResetConfigData resets all changes to the "configData" field.
func (m *ConnectorMutation) ResetConfigData() {
	m.configData = nil
	delete(m.clearedFields, connector.FieldConfigData)
}

// Where appends a list predicates to the ConnectorMutation builder.
func (m *ConnectorMutation) Where(ps ...predicate.Connector) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ConnectorMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ConnectorMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Connector, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ConnectorMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ConnectorMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Connector).
func (m *ConnectorMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ConnectorMutation) Fields() []string {
	fields := make([]string, 0, 7)
	if m.status != nil {
		fields = append(fields, connector.FieldStatus)
	}
	if m.statusMessage != nil {
		fields = append(fields, connector.FieldStatusMessage)
	}
	if m.createTime != nil {
		fields = append(fields, connector.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, connector.FieldUpdateTime)
	}
	if m._driver != nil {
		fields = append(fields, connector.FieldDriver)
	}
	if m.configVersion != nil {
		fields = append(fields, connector.FieldConfigVersion)
	}
	if m.configData != nil {
		fields = append(fields, connector.FieldConfigData)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ConnectorMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case connector.FieldStatus:
		return m.Status()
	case connector.FieldStatusMessage:
		return m.StatusMessage()
	case connector.FieldCreateTime:
		return m.CreateTime()
	case connector.FieldUpdateTime:
		return m.UpdateTime()
	case connector.FieldDriver:
		return m.Driver()
	case connector.FieldConfigVersion:
		return m.ConfigVersion()
	case connector.FieldConfigData:
		return m.ConfigData()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ConnectorMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case connector.FieldStatus:
		return m.OldStatus(ctx)
	case connector.FieldStatusMessage:
		return m.OldStatusMessage(ctx)
	case connector.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case connector.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case connector.FieldDriver:
		return m.OldDriver(ctx)
	case connector.FieldConfigVersion:
		return m.OldConfigVersion(ctx)
	case connector.FieldConfigData:
		return m.OldConfigData(ctx)
	}
	return nil, fmt.Errorf("unknown Connector field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ConnectorMutation) SetField(name string, value ent.Value) error {
	switch name {
	case connector.FieldStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case connector.FieldStatusMessage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatusMessage(v)
		return nil
	case connector.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case connector.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case connector.FieldDriver:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDriver(v)
		return nil
	case connector.FieldConfigVersion:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetConfigVersion(v)
		return nil
	case connector.FieldConfigData:
		v, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetConfigData(v)
		return nil
	}
	return fmt.Errorf("unknown Connector field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ConnectorMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ConnectorMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ConnectorMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Connector numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ConnectorMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(connector.FieldStatus) {
		fields = append(fields, connector.FieldStatus)
	}
	if m.FieldCleared(connector.FieldStatusMessage) {
		fields = append(fields, connector.FieldStatusMessage)
	}
	if m.FieldCleared(connector.FieldConfigData) {
		fields = append(fields, connector.FieldConfigData)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ConnectorMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ConnectorMutation) ClearField(name string) error {
	switch name {
	case connector.FieldStatus:
		m.ClearStatus()
		return nil
	case connector.FieldStatusMessage:
		m.ClearStatusMessage()
		return nil
	case connector.FieldConfigData:
		m.ClearConfigData()
		return nil
	}
	return fmt.Errorf("unknown Connector nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ConnectorMutation) ResetField(name string) error {
	switch name {
	case connector.FieldStatus:
		m.ResetStatus()
		return nil
	case connector.FieldStatusMessage:
		m.ResetStatusMessage()
		return nil
	case connector.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case connector.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case connector.FieldDriver:
		m.ResetDriver()
		return nil
	case connector.FieldConfigVersion:
		m.ResetConfigVersion()
		return nil
	case connector.FieldConfigData:
		m.ResetConfigData()
		return nil
	}
	return fmt.Errorf("unknown Connector field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ConnectorMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ConnectorMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ConnectorMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ConnectorMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ConnectorMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ConnectorMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ConnectorMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Connector unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ConnectorMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Connector edge %s", name)
}

// EnvironmentMutation represents an operation that mutates the Environment nodes in the graph.
type EnvironmentMutation struct {
	config
	op                 Op
	typ                string
	id                 *oid.ID
	createTime         *time.Time
	updateTime         *time.Time
	connectorIDs       *[]oid.ID
	appendconnectorIDs []oid.ID
	variables          *map[string]interface{}
	clearedFields      map[string]struct{}
	done               bool
	oldValue           func(context.Context) (*Environment, error)
	predicates         []predicate.Environment
}

var _ ent.Mutation = (*EnvironmentMutation)(nil)

// environmentOption allows management of the mutation configuration using functional options.
type environmentOption func(*EnvironmentMutation)

// newEnvironmentMutation creates new mutation for the Environment entity.
func newEnvironmentMutation(c config, op Op, opts ...environmentOption) *EnvironmentMutation {
	m := &EnvironmentMutation{
		config:        c,
		op:            op,
		typ:           TypeEnvironment,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withEnvironmentID sets the ID field of the mutation.
func withEnvironmentID(id oid.ID) environmentOption {
	return func(m *EnvironmentMutation) {
		var (
			err   error
			once  sync.Once
			value *Environment
		)
		m.oldValue = func(ctx context.Context) (*Environment, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Environment.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withEnvironment sets the old Environment of the mutation.
func withEnvironment(node *Environment) environmentOption {
	return func(m *EnvironmentMutation) {
		m.oldValue = func(context.Context) (*Environment, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EnvironmentMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EnvironmentMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Environment entities.
func (m *EnvironmentMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *EnvironmentMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *EnvironmentMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Environment.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *EnvironmentMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *EnvironmentMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Environment entity.
// If the Environment object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *EnvironmentMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *EnvironmentMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *EnvironmentMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Environment entity.
// If the Environment object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *EnvironmentMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetConnectorIDs sets the "connectorIDs" field.
func (m *EnvironmentMutation) SetConnectorIDs(o []oid.ID) {
	m.connectorIDs = &o
	m.appendconnectorIDs = nil
}

// ConnectorIDs returns the value of the "connectorIDs" field in the mutation.
func (m *EnvironmentMutation) ConnectorIDs() (r []oid.ID, exists bool) {
	v := m.connectorIDs
	if v == nil {
		return
	}
	return *v, true
}

// OldConnectorIDs returns the old "connectorIDs" field's value of the Environment entity.
// If the Environment object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentMutation) OldConnectorIDs(ctx context.Context) (v []oid.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldConnectorIDs is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldConnectorIDs requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldConnectorIDs: %w", err)
	}
	return oldValue.ConnectorIDs, nil
}

// AppendConnectorIDs adds o to the "connectorIDs" field.
func (m *EnvironmentMutation) AppendConnectorIDs(o []oid.ID) {
	m.appendconnectorIDs = append(m.appendconnectorIDs, o...)
}

// AppendedConnectorIDs returns the list of values that were appended to the "connectorIDs" field in this mutation.
func (m *EnvironmentMutation) AppendedConnectorIDs() ([]oid.ID, bool) {
	if len(m.appendconnectorIDs) == 0 {
		return nil, false
	}
	return m.appendconnectorIDs, true
}

// ClearConnectorIDs clears the value of the "connectorIDs" field.
func (m *EnvironmentMutation) ClearConnectorIDs() {
	m.connectorIDs = nil
	m.appendconnectorIDs = nil
	m.clearedFields[environment.FieldConnectorIDs] = struct{}{}
}

// ConnectorIDsCleared returns if the "connectorIDs" field was cleared in this mutation.
func (m *EnvironmentMutation) ConnectorIDsCleared() bool {
	_, ok := m.clearedFields[environment.FieldConnectorIDs]
	return ok
}

// ResetConnectorIDs resets all changes to the "connectorIDs" field.
func (m *EnvironmentMutation) ResetConnectorIDs() {
	m.connectorIDs = nil
	m.appendconnectorIDs = nil
	delete(m.clearedFields, environment.FieldConnectorIDs)
}

// SetVariables sets the "variables" field.
func (m *EnvironmentMutation) SetVariables(value map[string]interface{}) {
	m.variables = &value
}

// Variables returns the value of the "variables" field in the mutation.
func (m *EnvironmentMutation) Variables() (r map[string]interface{}, exists bool) {
	v := m.variables
	if v == nil {
		return
	}
	return *v, true
}

// OldVariables returns the old "variables" field's value of the Environment entity.
// If the Environment object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentMutation) OldVariables(ctx context.Context) (v map[string]interface{}, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldVariables is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldVariables requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldVariables: %w", err)
	}
	return oldValue.Variables, nil
}

// ClearVariables clears the value of the "variables" field.
func (m *EnvironmentMutation) ClearVariables() {
	m.variables = nil
	m.clearedFields[environment.FieldVariables] = struct{}{}
}

// VariablesCleared returns if the "variables" field was cleared in this mutation.
func (m *EnvironmentMutation) VariablesCleared() bool {
	_, ok := m.clearedFields[environment.FieldVariables]
	return ok
}

// ResetVariables resets all changes to the "variables" field.
func (m *EnvironmentMutation) ResetVariables() {
	m.variables = nil
	delete(m.clearedFields, environment.FieldVariables)
}

// Where appends a list predicates to the EnvironmentMutation builder.
func (m *EnvironmentMutation) Where(ps ...predicate.Environment) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the EnvironmentMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *EnvironmentMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Environment, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *EnvironmentMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *EnvironmentMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Environment).
func (m *EnvironmentMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *EnvironmentMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.createTime != nil {
		fields = append(fields, environment.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, environment.FieldUpdateTime)
	}
	if m.connectorIDs != nil {
		fields = append(fields, environment.FieldConnectorIDs)
	}
	if m.variables != nil {
		fields = append(fields, environment.FieldVariables)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *EnvironmentMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case environment.FieldCreateTime:
		return m.CreateTime()
	case environment.FieldUpdateTime:
		return m.UpdateTime()
	case environment.FieldConnectorIDs:
		return m.ConnectorIDs()
	case environment.FieldVariables:
		return m.Variables()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *EnvironmentMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case environment.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case environment.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case environment.FieldConnectorIDs:
		return m.OldConnectorIDs(ctx)
	case environment.FieldVariables:
		return m.OldVariables(ctx)
	}
	return nil, fmt.Errorf("unknown Environment field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *EnvironmentMutation) SetField(name string, value ent.Value) error {
	switch name {
	case environment.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case environment.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case environment.FieldConnectorIDs:
		v, ok := value.([]oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetConnectorIDs(v)
		return nil
	case environment.FieldVariables:
		v, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVariables(v)
		return nil
	}
	return fmt.Errorf("unknown Environment field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *EnvironmentMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *EnvironmentMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *EnvironmentMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Environment numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *EnvironmentMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(environment.FieldConnectorIDs) {
		fields = append(fields, environment.FieldConnectorIDs)
	}
	if m.FieldCleared(environment.FieldVariables) {
		fields = append(fields, environment.FieldVariables)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *EnvironmentMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *EnvironmentMutation) ClearField(name string) error {
	switch name {
	case environment.FieldConnectorIDs:
		m.ClearConnectorIDs()
		return nil
	case environment.FieldVariables:
		m.ClearVariables()
		return nil
	}
	return fmt.Errorf("unknown Environment nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *EnvironmentMutation) ResetField(name string) error {
	switch name {
	case environment.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case environment.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case environment.FieldConnectorIDs:
		m.ResetConnectorIDs()
		return nil
	case environment.FieldVariables:
		m.ResetVariables()
		return nil
	}
	return fmt.Errorf("unknown Environment field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *EnvironmentMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *EnvironmentMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *EnvironmentMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *EnvironmentMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *EnvironmentMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *EnvironmentMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *EnvironmentMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Environment unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *EnvironmentMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Environment edge %s", name)
}

// ModuleMutation represents an operation that mutates the Module nodes in the graph.
type ModuleMutation struct {
	config
	op            Op
	typ           string
	id            *oid.ID
	status        *string
	statusMessage *string
	createTime    *time.Time
	updateTime    *time.Time
	source        *string
	version       *string
	inputSchema   *map[string]interface{}
	outputSchema  *map[string]interface{}
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Module, error)
	predicates    []predicate.Module
}

var _ ent.Mutation = (*ModuleMutation)(nil)

// moduleOption allows management of the mutation configuration using functional options.
type moduleOption func(*ModuleMutation)

// newModuleMutation creates new mutation for the Module entity.
func newModuleMutation(c config, op Op, opts ...moduleOption) *ModuleMutation {
	m := &ModuleMutation{
		config:        c,
		op:            op,
		typ:           TypeModule,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withModuleID sets the ID field of the mutation.
func withModuleID(id oid.ID) moduleOption {
	return func(m *ModuleMutation) {
		var (
			err   error
			once  sync.Once
			value *Module
		)
		m.oldValue = func(ctx context.Context) (*Module, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Module.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withModule sets the old Module of the mutation.
func withModule(node *Module) moduleOption {
	return func(m *ModuleMutation) {
		m.oldValue = func(context.Context) (*Module, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ModuleMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ModuleMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Module entities.
func (m *ModuleMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ModuleMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ModuleMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Module.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetStatus sets the "status" field.
func (m *ModuleMutation) SetStatus(s string) {
	m.status = &s
}

// Status returns the value of the "status" field in the mutation.
func (m *ModuleMutation) Status() (r string, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldStatus(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatus is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatus requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatus: %w", err)
	}
	return oldValue.Status, nil
}

// ClearStatus clears the value of the "status" field.
func (m *ModuleMutation) ClearStatus() {
	m.status = nil
	m.clearedFields[module.FieldStatus] = struct{}{}
}

// StatusCleared returns if the "status" field was cleared in this mutation.
func (m *ModuleMutation) StatusCleared() bool {
	_, ok := m.clearedFields[module.FieldStatus]
	return ok
}

// ResetStatus resets all changes to the "status" field.
func (m *ModuleMutation) ResetStatus() {
	m.status = nil
	delete(m.clearedFields, module.FieldStatus)
}

// SetStatusMessage sets the "statusMessage" field.
func (m *ModuleMutation) SetStatusMessage(s string) {
	m.statusMessage = &s
}

// StatusMessage returns the value of the "statusMessage" field in the mutation.
func (m *ModuleMutation) StatusMessage() (r string, exists bool) {
	v := m.statusMessage
	if v == nil {
		return
	}
	return *v, true
}

// OldStatusMessage returns the old "statusMessage" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldStatusMessage(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatusMessage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatusMessage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatusMessage: %w", err)
	}
	return oldValue.StatusMessage, nil
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (m *ModuleMutation) ClearStatusMessage() {
	m.statusMessage = nil
	m.clearedFields[module.FieldStatusMessage] = struct{}{}
}

// StatusMessageCleared returns if the "statusMessage" field was cleared in this mutation.
func (m *ModuleMutation) StatusMessageCleared() bool {
	_, ok := m.clearedFields[module.FieldStatusMessage]
	return ok
}

// ResetStatusMessage resets all changes to the "statusMessage" field.
func (m *ModuleMutation) ResetStatusMessage() {
	m.statusMessage = nil
	delete(m.clearedFields, module.FieldStatusMessage)
}

// SetCreateTime sets the "createTime" field.
func (m *ModuleMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *ModuleMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *ModuleMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *ModuleMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *ModuleMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *ModuleMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetSource sets the "source" field.
func (m *ModuleMutation) SetSource(s string) {
	m.source = &s
}

// Source returns the value of the "source" field in the mutation.
func (m *ModuleMutation) Source() (r string, exists bool) {
	v := m.source
	if v == nil {
		return
	}
	return *v, true
}

// OldSource returns the old "source" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldSource(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSource is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSource requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSource: %w", err)
	}
	return oldValue.Source, nil
}

// ResetSource resets all changes to the "source" field.
func (m *ModuleMutation) ResetSource() {
	m.source = nil
}

// SetVersion sets the "version" field.
func (m *ModuleMutation) SetVersion(s string) {
	m.version = &s
}

// Version returns the value of the "version" field in the mutation.
func (m *ModuleMutation) Version() (r string, exists bool) {
	v := m.version
	if v == nil {
		return
	}
	return *v, true
}

// OldVersion returns the old "version" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldVersion(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldVersion is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldVersion requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldVersion: %w", err)
	}
	return oldValue.Version, nil
}

// ResetVersion resets all changes to the "version" field.
func (m *ModuleMutation) ResetVersion() {
	m.version = nil
}

// SetInputSchema sets the "inputSchema" field.
func (m *ModuleMutation) SetInputSchema(value map[string]interface{}) {
	m.inputSchema = &value
}

// InputSchema returns the value of the "inputSchema" field in the mutation.
func (m *ModuleMutation) InputSchema() (r map[string]interface{}, exists bool) {
	v := m.inputSchema
	if v == nil {
		return
	}
	return *v, true
}

// OldInputSchema returns the old "inputSchema" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldInputSchema(ctx context.Context) (v map[string]interface{}, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldInputSchema is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldInputSchema requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldInputSchema: %w", err)
	}
	return oldValue.InputSchema, nil
}

// ClearInputSchema clears the value of the "inputSchema" field.
func (m *ModuleMutation) ClearInputSchema() {
	m.inputSchema = nil
	m.clearedFields[module.FieldInputSchema] = struct{}{}
}

// InputSchemaCleared returns if the "inputSchema" field was cleared in this mutation.
func (m *ModuleMutation) InputSchemaCleared() bool {
	_, ok := m.clearedFields[module.FieldInputSchema]
	return ok
}

// ResetInputSchema resets all changes to the "inputSchema" field.
func (m *ModuleMutation) ResetInputSchema() {
	m.inputSchema = nil
	delete(m.clearedFields, module.FieldInputSchema)
}

// SetOutputSchema sets the "outputSchema" field.
func (m *ModuleMutation) SetOutputSchema(value map[string]interface{}) {
	m.outputSchema = &value
}

// OutputSchema returns the value of the "outputSchema" field in the mutation.
func (m *ModuleMutation) OutputSchema() (r map[string]interface{}, exists bool) {
	v := m.outputSchema
	if v == nil {
		return
	}
	return *v, true
}

// OldOutputSchema returns the old "outputSchema" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldOutputSchema(ctx context.Context) (v map[string]interface{}, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldOutputSchema is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldOutputSchema requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldOutputSchema: %w", err)
	}
	return oldValue.OutputSchema, nil
}

// ClearOutputSchema clears the value of the "outputSchema" field.
func (m *ModuleMutation) ClearOutputSchema() {
	m.outputSchema = nil
	m.clearedFields[module.FieldOutputSchema] = struct{}{}
}

// OutputSchemaCleared returns if the "outputSchema" field was cleared in this mutation.
func (m *ModuleMutation) OutputSchemaCleared() bool {
	_, ok := m.clearedFields[module.FieldOutputSchema]
	return ok
}

// ResetOutputSchema resets all changes to the "outputSchema" field.
func (m *ModuleMutation) ResetOutputSchema() {
	m.outputSchema = nil
	delete(m.clearedFields, module.FieldOutputSchema)
}

// Where appends a list predicates to the ModuleMutation builder.
func (m *ModuleMutation) Where(ps ...predicate.Module) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ModuleMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ModuleMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Module, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ModuleMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ModuleMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Module).
func (m *ModuleMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ModuleMutation) Fields() []string {
	fields := make([]string, 0, 8)
	if m.status != nil {
		fields = append(fields, module.FieldStatus)
	}
	if m.statusMessage != nil {
		fields = append(fields, module.FieldStatusMessage)
	}
	if m.createTime != nil {
		fields = append(fields, module.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, module.FieldUpdateTime)
	}
	if m.source != nil {
		fields = append(fields, module.FieldSource)
	}
	if m.version != nil {
		fields = append(fields, module.FieldVersion)
	}
	if m.inputSchema != nil {
		fields = append(fields, module.FieldInputSchema)
	}
	if m.outputSchema != nil {
		fields = append(fields, module.FieldOutputSchema)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ModuleMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case module.FieldStatus:
		return m.Status()
	case module.FieldStatusMessage:
		return m.StatusMessage()
	case module.FieldCreateTime:
		return m.CreateTime()
	case module.FieldUpdateTime:
		return m.UpdateTime()
	case module.FieldSource:
		return m.Source()
	case module.FieldVersion:
		return m.Version()
	case module.FieldInputSchema:
		return m.InputSchema()
	case module.FieldOutputSchema:
		return m.OutputSchema()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ModuleMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case module.FieldStatus:
		return m.OldStatus(ctx)
	case module.FieldStatusMessage:
		return m.OldStatusMessage(ctx)
	case module.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case module.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case module.FieldSource:
		return m.OldSource(ctx)
	case module.FieldVersion:
		return m.OldVersion(ctx)
	case module.FieldInputSchema:
		return m.OldInputSchema(ctx)
	case module.FieldOutputSchema:
		return m.OldOutputSchema(ctx)
	}
	return nil, fmt.Errorf("unknown Module field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ModuleMutation) SetField(name string, value ent.Value) error {
	switch name {
	case module.FieldStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case module.FieldStatusMessage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatusMessage(v)
		return nil
	case module.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case module.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case module.FieldSource:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSource(v)
		return nil
	case module.FieldVersion:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVersion(v)
		return nil
	case module.FieldInputSchema:
		v, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetInputSchema(v)
		return nil
	case module.FieldOutputSchema:
		v, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetOutputSchema(v)
		return nil
	}
	return fmt.Errorf("unknown Module field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ModuleMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ModuleMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ModuleMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Module numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ModuleMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(module.FieldStatus) {
		fields = append(fields, module.FieldStatus)
	}
	if m.FieldCleared(module.FieldStatusMessage) {
		fields = append(fields, module.FieldStatusMessage)
	}
	if m.FieldCleared(module.FieldInputSchema) {
		fields = append(fields, module.FieldInputSchema)
	}
	if m.FieldCleared(module.FieldOutputSchema) {
		fields = append(fields, module.FieldOutputSchema)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ModuleMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ModuleMutation) ClearField(name string) error {
	switch name {
	case module.FieldStatus:
		m.ClearStatus()
		return nil
	case module.FieldStatusMessage:
		m.ClearStatusMessage()
		return nil
	case module.FieldInputSchema:
		m.ClearInputSchema()
		return nil
	case module.FieldOutputSchema:
		m.ClearOutputSchema()
		return nil
	}
	return fmt.Errorf("unknown Module nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ModuleMutation) ResetField(name string) error {
	switch name {
	case module.FieldStatus:
		m.ResetStatus()
		return nil
	case module.FieldStatusMessage:
		m.ResetStatusMessage()
		return nil
	case module.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case module.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case module.FieldSource:
		m.ResetSource()
		return nil
	case module.FieldVersion:
		m.ResetVersion()
		return nil
	case module.FieldInputSchema:
		m.ResetInputSchema()
		return nil
	case module.FieldOutputSchema:
		m.ResetOutputSchema()
		return nil
	}
	return fmt.Errorf("unknown Module field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ModuleMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ModuleMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ModuleMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ModuleMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ModuleMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ModuleMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ModuleMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Module unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ModuleMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Module edge %s", name)
}

// ProjectMutation represents an operation that mutates the Project nodes in the graph.
type ProjectMutation struct {
	config
	op            Op
	typ           string
	id            *oid.ID
	createTime    *time.Time
	updateTime    *time.Time
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Project, error)
	predicates    []predicate.Project
}

var _ ent.Mutation = (*ProjectMutation)(nil)

// projectOption allows management of the mutation configuration using functional options.
type projectOption func(*ProjectMutation)

// newProjectMutation creates new mutation for the Project entity.
func newProjectMutation(c config, op Op, opts ...projectOption) *ProjectMutation {
	m := &ProjectMutation{
		config:        c,
		op:            op,
		typ:           TypeProject,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withProjectID sets the ID field of the mutation.
func withProjectID(id oid.ID) projectOption {
	return func(m *ProjectMutation) {
		var (
			err   error
			once  sync.Once
			value *Project
		)
		m.oldValue = func(ctx context.Context) (*Project, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Project.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withProject sets the old Project of the mutation.
func withProject(node *Project) projectOption {
	return func(m *ProjectMutation) {
		m.oldValue = func(context.Context) (*Project, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ProjectMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ProjectMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Project entities.
func (m *ProjectMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ProjectMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ProjectMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Project.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *ProjectMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *ProjectMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Project entity.
// If the Project object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProjectMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *ProjectMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *ProjectMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *ProjectMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Project entity.
// If the Project object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProjectMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *ProjectMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// Where appends a list predicates to the ProjectMutation builder.
func (m *ProjectMutation) Where(ps ...predicate.Project) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ProjectMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ProjectMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Project, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ProjectMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ProjectMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Project).
func (m *ProjectMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ProjectMutation) Fields() []string {
	fields := make([]string, 0, 2)
	if m.createTime != nil {
		fields = append(fields, project.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, project.FieldUpdateTime)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ProjectMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case project.FieldCreateTime:
		return m.CreateTime()
	case project.FieldUpdateTime:
		return m.UpdateTime()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ProjectMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case project.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case project.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	}
	return nil, fmt.Errorf("unknown Project field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ProjectMutation) SetField(name string, value ent.Value) error {
	switch name {
	case project.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case project.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	}
	return fmt.Errorf("unknown Project field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ProjectMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ProjectMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ProjectMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Project numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ProjectMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ProjectMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ProjectMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Project nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ProjectMutation) ResetField(name string) error {
	switch name {
	case project.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case project.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	}
	return fmt.Errorf("unknown Project field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ProjectMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ProjectMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ProjectMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ProjectMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ProjectMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ProjectMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ProjectMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Project unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ProjectMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Project edge %s", name)
}

// RoleMutation represents an operation that mutates the Role nodes in the graph.
type RoleMutation struct {
	config
	op             Op
	typ            string
	id             *oid.ID
	createTime     *time.Time
	updateTime     *time.Time
	domain         *string
	name           *string
	description    *string
	policies       *schema.RolePolicies
	appendpolicies schema.RolePolicies
	builtin        *bool
	session        *bool
	clearedFields  map[string]struct{}
	done           bool
	oldValue       func(context.Context) (*Role, error)
	predicates     []predicate.Role
}

var _ ent.Mutation = (*RoleMutation)(nil)

// roleOption allows management of the mutation configuration using functional options.
type roleOption func(*RoleMutation)

// newRoleMutation creates new mutation for the Role entity.
func newRoleMutation(c config, op Op, opts ...roleOption) *RoleMutation {
	m := &RoleMutation{
		config:        c,
		op:            op,
		typ:           TypeRole,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withRoleID sets the ID field of the mutation.
func withRoleID(id oid.ID) roleOption {
	return func(m *RoleMutation) {
		var (
			err   error
			once  sync.Once
			value *Role
		)
		m.oldValue = func(ctx context.Context) (*Role, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Role.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withRole sets the old Role of the mutation.
func withRole(node *Role) roleOption {
	return func(m *RoleMutation) {
		m.oldValue = func(context.Context) (*Role, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m RoleMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m RoleMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Role entities.
func (m *RoleMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *RoleMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *RoleMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Role.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *RoleMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *RoleMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Role entity.
// If the Role object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoleMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *RoleMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *RoleMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *RoleMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Role entity.
// If the Role object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoleMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *RoleMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetDomain sets the "domain" field.
func (m *RoleMutation) SetDomain(s string) {
	m.domain = &s
}

// Domain returns the value of the "domain" field in the mutation.
func (m *RoleMutation) Domain() (r string, exists bool) {
	v := m.domain
	if v == nil {
		return
	}
	return *v, true
}

// OldDomain returns the old "domain" field's value of the Role entity.
// If the Role object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoleMutation) OldDomain(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDomain is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDomain requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDomain: %w", err)
	}
	return oldValue.Domain, nil
}

// ResetDomain resets all changes to the "domain" field.
func (m *RoleMutation) ResetDomain() {
	m.domain = nil
}

// SetName sets the "name" field.
func (m *RoleMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *RoleMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Role entity.
// If the Role object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoleMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *RoleMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the "description" field.
func (m *RoleMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the value of the "description" field in the mutation.
func (m *RoleMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// OldDescription returns the old "description" field's value of the Role entity.
// If the Role object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoleMutation) OldDescription(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDescription is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDescription requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDescription: %w", err)
	}
	return oldValue.Description, nil
}

// ResetDescription resets all changes to the "description" field.
func (m *RoleMutation) ResetDescription() {
	m.description = nil
}

// SetPolicies sets the "policies" field.
func (m *RoleMutation) SetPolicies(sp schema.RolePolicies) {
	m.policies = &sp
	m.appendpolicies = nil
}

// Policies returns the value of the "policies" field in the mutation.
func (m *RoleMutation) Policies() (r schema.RolePolicies, exists bool) {
	v := m.policies
	if v == nil {
		return
	}
	return *v, true
}

// OldPolicies returns the old "policies" field's value of the Role entity.
// If the Role object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoleMutation) OldPolicies(ctx context.Context) (v schema.RolePolicies, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPolicies is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPolicies requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPolicies: %w", err)
	}
	return oldValue.Policies, nil
}

// AppendPolicies adds sp to the "policies" field.
func (m *RoleMutation) AppendPolicies(sp schema.RolePolicies) {
	m.appendpolicies = append(m.appendpolicies, sp...)
}

// AppendedPolicies returns the list of values that were appended to the "policies" field in this mutation.
func (m *RoleMutation) AppendedPolicies() (schema.RolePolicies, bool) {
	if len(m.appendpolicies) == 0 {
		return nil, false
	}
	return m.appendpolicies, true
}

// ResetPolicies resets all changes to the "policies" field.
func (m *RoleMutation) ResetPolicies() {
	m.policies = nil
	m.appendpolicies = nil
}

// SetBuiltin sets the "builtin" field.
func (m *RoleMutation) SetBuiltin(b bool) {
	m.builtin = &b
}

// Builtin returns the value of the "builtin" field in the mutation.
func (m *RoleMutation) Builtin() (r bool, exists bool) {
	v := m.builtin
	if v == nil {
		return
	}
	return *v, true
}

// OldBuiltin returns the old "builtin" field's value of the Role entity.
// If the Role object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoleMutation) OldBuiltin(ctx context.Context) (v bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldBuiltin is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldBuiltin requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldBuiltin: %w", err)
	}
	return oldValue.Builtin, nil
}

// ResetBuiltin resets all changes to the "builtin" field.
func (m *RoleMutation) ResetBuiltin() {
	m.builtin = nil
}

// SetSession sets the "session" field.
func (m *RoleMutation) SetSession(b bool) {
	m.session = &b
}

// Session returns the value of the "session" field in the mutation.
func (m *RoleMutation) Session() (r bool, exists bool) {
	v := m.session
	if v == nil {
		return
	}
	return *v, true
}

// OldSession returns the old "session" field's value of the Role entity.
// If the Role object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoleMutation) OldSession(ctx context.Context) (v bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSession is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSession requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSession: %w", err)
	}
	return oldValue.Session, nil
}

// ResetSession resets all changes to the "session" field.
func (m *RoleMutation) ResetSession() {
	m.session = nil
}

// Where appends a list predicates to the RoleMutation builder.
func (m *RoleMutation) Where(ps ...predicate.Role) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the RoleMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *RoleMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Role, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *RoleMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *RoleMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Role).
func (m *RoleMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *RoleMutation) Fields() []string {
	fields := make([]string, 0, 8)
	if m.createTime != nil {
		fields = append(fields, role.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, role.FieldUpdateTime)
	}
	if m.domain != nil {
		fields = append(fields, role.FieldDomain)
	}
	if m.name != nil {
		fields = append(fields, role.FieldName)
	}
	if m.description != nil {
		fields = append(fields, role.FieldDescription)
	}
	if m.policies != nil {
		fields = append(fields, role.FieldPolicies)
	}
	if m.builtin != nil {
		fields = append(fields, role.FieldBuiltin)
	}
	if m.session != nil {
		fields = append(fields, role.FieldSession)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *RoleMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case role.FieldCreateTime:
		return m.CreateTime()
	case role.FieldUpdateTime:
		return m.UpdateTime()
	case role.FieldDomain:
		return m.Domain()
	case role.FieldName:
		return m.Name()
	case role.FieldDescription:
		return m.Description()
	case role.FieldPolicies:
		return m.Policies()
	case role.FieldBuiltin:
		return m.Builtin()
	case role.FieldSession:
		return m.Session()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *RoleMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case role.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case role.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case role.FieldDomain:
		return m.OldDomain(ctx)
	case role.FieldName:
		return m.OldName(ctx)
	case role.FieldDescription:
		return m.OldDescription(ctx)
	case role.FieldPolicies:
		return m.OldPolicies(ctx)
	case role.FieldBuiltin:
		return m.OldBuiltin(ctx)
	case role.FieldSession:
		return m.OldSession(ctx)
	}
	return nil, fmt.Errorf("unknown Role field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RoleMutation) SetField(name string, value ent.Value) error {
	switch name {
	case role.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case role.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case role.FieldDomain:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDomain(v)
		return nil
	case role.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case role.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case role.FieldPolicies:
		v, ok := value.(schema.RolePolicies)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPolicies(v)
		return nil
	case role.FieldBuiltin:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBuiltin(v)
		return nil
	case role.FieldSession:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSession(v)
		return nil
	}
	return fmt.Errorf("unknown Role field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *RoleMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *RoleMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RoleMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Role numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *RoleMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *RoleMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *RoleMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Role nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *RoleMutation) ResetField(name string) error {
	switch name {
	case role.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case role.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case role.FieldDomain:
		m.ResetDomain()
		return nil
	case role.FieldName:
		m.ResetName()
		return nil
	case role.FieldDescription:
		m.ResetDescription()
		return nil
	case role.FieldPolicies:
		m.ResetPolicies()
		return nil
	case role.FieldBuiltin:
		m.ResetBuiltin()
		return nil
	case role.FieldSession:
		m.ResetSession()
		return nil
	}
	return fmt.Errorf("unknown Role field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *RoleMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *RoleMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *RoleMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *RoleMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *RoleMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *RoleMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *RoleMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Role unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *RoleMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Role edge %s", name)
}

// SettingMutation represents an operation that mutates the Setting nodes in the graph.
type SettingMutation struct {
	config
	op            Op
	typ           string
	id            *oid.ID
	createTime    *time.Time
	updateTime    *time.Time
	name          *string
	value         *string
	hidden        *bool
	editable      *bool
	private       *bool
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Setting, error)
	predicates    []predicate.Setting
}

var _ ent.Mutation = (*SettingMutation)(nil)

// settingOption allows management of the mutation configuration using functional options.
type settingOption func(*SettingMutation)

// newSettingMutation creates new mutation for the Setting entity.
func newSettingMutation(c config, op Op, opts ...settingOption) *SettingMutation {
	m := &SettingMutation{
		config:        c,
		op:            op,
		typ:           TypeSetting,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withSettingID sets the ID field of the mutation.
func withSettingID(id oid.ID) settingOption {
	return func(m *SettingMutation) {
		var (
			err   error
			once  sync.Once
			value *Setting
		)
		m.oldValue = func(ctx context.Context) (*Setting, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Setting.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withSetting sets the old Setting of the mutation.
func withSetting(node *Setting) settingOption {
	return func(m *SettingMutation) {
		m.oldValue = func(context.Context) (*Setting, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SettingMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SettingMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Setting entities.
func (m *SettingMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SettingMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SettingMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Setting.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *SettingMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *SettingMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Setting entity.
// If the Setting object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SettingMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *SettingMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *SettingMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *SettingMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Setting entity.
// If the Setting object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SettingMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *SettingMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetName sets the "name" field.
func (m *SettingMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *SettingMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Setting entity.
// If the Setting object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SettingMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *SettingMutation) ResetName() {
	m.name = nil
}

// SetValue sets the "value" field.
func (m *SettingMutation) SetValue(s string) {
	m.value = &s
}

// Value returns the value of the "value" field in the mutation.
func (m *SettingMutation) Value() (r string, exists bool) {
	v := m.value
	if v == nil {
		return
	}
	return *v, true
}

// OldValue returns the old "value" field's value of the Setting entity.
// If the Setting object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SettingMutation) OldValue(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldValue is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldValue requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldValue: %w", err)
	}
	return oldValue.Value, nil
}

// ResetValue resets all changes to the "value" field.
func (m *SettingMutation) ResetValue() {
	m.value = nil
}

// SetHidden sets the "hidden" field.
func (m *SettingMutation) SetHidden(b bool) {
	m.hidden = &b
}

// Hidden returns the value of the "hidden" field in the mutation.
func (m *SettingMutation) Hidden() (r bool, exists bool) {
	v := m.hidden
	if v == nil {
		return
	}
	return *v, true
}

// OldHidden returns the old "hidden" field's value of the Setting entity.
// If the Setting object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SettingMutation) OldHidden(ctx context.Context) (v *bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldHidden is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldHidden requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldHidden: %w", err)
	}
	return oldValue.Hidden, nil
}

// ResetHidden resets all changes to the "hidden" field.
func (m *SettingMutation) ResetHidden() {
	m.hidden = nil
}

// SetEditable sets the "editable" field.
func (m *SettingMutation) SetEditable(b bool) {
	m.editable = &b
}

// Editable returns the value of the "editable" field in the mutation.
func (m *SettingMutation) Editable() (r bool, exists bool) {
	v := m.editable
	if v == nil {
		return
	}
	return *v, true
}

// OldEditable returns the old "editable" field's value of the Setting entity.
// If the Setting object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SettingMutation) OldEditable(ctx context.Context) (v *bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEditable is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEditable requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEditable: %w", err)
	}
	return oldValue.Editable, nil
}

// ResetEditable resets all changes to the "editable" field.
func (m *SettingMutation) ResetEditable() {
	m.editable = nil
}

// SetPrivate sets the "private" field.
func (m *SettingMutation) SetPrivate(b bool) {
	m.private = &b
}

// Private returns the value of the "private" field in the mutation.
func (m *SettingMutation) Private() (r bool, exists bool) {
	v := m.private
	if v == nil {
		return
	}
	return *v, true
}

// OldPrivate returns the old "private" field's value of the Setting entity.
// If the Setting object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SettingMutation) OldPrivate(ctx context.Context) (v *bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPrivate is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPrivate requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPrivate: %w", err)
	}
	return oldValue.Private, nil
}

// ResetPrivate resets all changes to the "private" field.
func (m *SettingMutation) ResetPrivate() {
	m.private = nil
}

// Where appends a list predicates to the SettingMutation builder.
func (m *SettingMutation) Where(ps ...predicate.Setting) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the SettingMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *SettingMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Setting, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *SettingMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *SettingMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Setting).
func (m *SettingMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *SettingMutation) Fields() []string {
	fields := make([]string, 0, 7)
	if m.createTime != nil {
		fields = append(fields, setting.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, setting.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, setting.FieldName)
	}
	if m.value != nil {
		fields = append(fields, setting.FieldValue)
	}
	if m.hidden != nil {
		fields = append(fields, setting.FieldHidden)
	}
	if m.editable != nil {
		fields = append(fields, setting.FieldEditable)
	}
	if m.private != nil {
		fields = append(fields, setting.FieldPrivate)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *SettingMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case setting.FieldCreateTime:
		return m.CreateTime()
	case setting.FieldUpdateTime:
		return m.UpdateTime()
	case setting.FieldName:
		return m.Name()
	case setting.FieldValue:
		return m.Value()
	case setting.FieldHidden:
		return m.Hidden()
	case setting.FieldEditable:
		return m.Editable()
	case setting.FieldPrivate:
		return m.Private()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *SettingMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case setting.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case setting.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case setting.FieldName:
		return m.OldName(ctx)
	case setting.FieldValue:
		return m.OldValue(ctx)
	case setting.FieldHidden:
		return m.OldHidden(ctx)
	case setting.FieldEditable:
		return m.OldEditable(ctx)
	case setting.FieldPrivate:
		return m.OldPrivate(ctx)
	}
	return nil, fmt.Errorf("unknown Setting field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SettingMutation) SetField(name string, value ent.Value) error {
	switch name {
	case setting.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case setting.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case setting.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case setting.FieldValue:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetValue(v)
		return nil
	case setting.FieldHidden:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetHidden(v)
		return nil
	case setting.FieldEditable:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEditable(v)
		return nil
	case setting.FieldPrivate:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPrivate(v)
		return nil
	}
	return fmt.Errorf("unknown Setting field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *SettingMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *SettingMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SettingMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Setting numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *SettingMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *SettingMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *SettingMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Setting nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *SettingMutation) ResetField(name string) error {
	switch name {
	case setting.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case setting.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case setting.FieldName:
		m.ResetName()
		return nil
	case setting.FieldValue:
		m.ResetValue()
		return nil
	case setting.FieldHidden:
		m.ResetHidden()
		return nil
	case setting.FieldEditable:
		m.ResetEditable()
		return nil
	case setting.FieldPrivate:
		m.ResetPrivate()
		return nil
	}
	return fmt.Errorf("unknown Setting field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *SettingMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *SettingMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *SettingMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *SettingMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *SettingMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *SettingMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *SettingMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Setting unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *SettingMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Setting edge %s", name)
}

// SubjectMutation represents an operation that mutates the Subject nodes in the graph.
type SubjectMutation struct {
	config
	op            Op
	typ           string
	id            *oid.ID
	createTime    *time.Time
	updateTime    *time.Time
	kind          *string
	group         *string
	name          *string
	description   *string
	mountTo       *bool
	loginTo       *bool
	roles         *schema.SubjectRoles
	appendroles   schema.SubjectRoles
	paths         *[]string
	appendpaths   []string
	builtin       *bool
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Subject, error)
	predicates    []predicate.Subject
}

var _ ent.Mutation = (*SubjectMutation)(nil)

// subjectOption allows management of the mutation configuration using functional options.
type subjectOption func(*SubjectMutation)

// newSubjectMutation creates new mutation for the Subject entity.
func newSubjectMutation(c config, op Op, opts ...subjectOption) *SubjectMutation {
	m := &SubjectMutation{
		config:        c,
		op:            op,
		typ:           TypeSubject,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withSubjectID sets the ID field of the mutation.
func withSubjectID(id oid.ID) subjectOption {
	return func(m *SubjectMutation) {
		var (
			err   error
			once  sync.Once
			value *Subject
		)
		m.oldValue = func(ctx context.Context) (*Subject, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Subject.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withSubject sets the old Subject of the mutation.
func withSubject(node *Subject) subjectOption {
	return func(m *SubjectMutation) {
		m.oldValue = func(context.Context) (*Subject, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SubjectMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SubjectMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Subject entities.
func (m *SubjectMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SubjectMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SubjectMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Subject.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *SubjectMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *SubjectMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *SubjectMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *SubjectMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *SubjectMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *SubjectMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetKind sets the "kind" field.
func (m *SubjectMutation) SetKind(s string) {
	m.kind = &s
}

// Kind returns the value of the "kind" field in the mutation.
func (m *SubjectMutation) Kind() (r string, exists bool) {
	v := m.kind
	if v == nil {
		return
	}
	return *v, true
}

// OldKind returns the old "kind" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldKind(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldKind is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldKind requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldKind: %w", err)
	}
	return oldValue.Kind, nil
}

// ResetKind resets all changes to the "kind" field.
func (m *SubjectMutation) ResetKind() {
	m.kind = nil
}

// SetGroup sets the "group" field.
func (m *SubjectMutation) SetGroup(s string) {
	m.group = &s
}

// Group returns the value of the "group" field in the mutation.
func (m *SubjectMutation) Group() (r string, exists bool) {
	v := m.group
	if v == nil {
		return
	}
	return *v, true
}

// OldGroup returns the old "group" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldGroup(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldGroup is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldGroup requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldGroup: %w", err)
	}
	return oldValue.Group, nil
}

// ResetGroup resets all changes to the "group" field.
func (m *SubjectMutation) ResetGroup() {
	m.group = nil
}

// SetName sets the "name" field.
func (m *SubjectMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *SubjectMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *SubjectMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the "description" field.
func (m *SubjectMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the value of the "description" field in the mutation.
func (m *SubjectMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// OldDescription returns the old "description" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldDescription(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDescription is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDescription requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDescription: %w", err)
	}
	return oldValue.Description, nil
}

// ResetDescription resets all changes to the "description" field.
func (m *SubjectMutation) ResetDescription() {
	m.description = nil
}

// SetMountTo sets the "mountTo" field.
func (m *SubjectMutation) SetMountTo(b bool) {
	m.mountTo = &b
}

// MountTo returns the value of the "mountTo" field in the mutation.
func (m *SubjectMutation) MountTo() (r bool, exists bool) {
	v := m.mountTo
	if v == nil {
		return
	}
	return *v, true
}

// OldMountTo returns the old "mountTo" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldMountTo(ctx context.Context) (v *bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldMountTo is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldMountTo requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldMountTo: %w", err)
	}
	return oldValue.MountTo, nil
}

// ResetMountTo resets all changes to the "mountTo" field.
func (m *SubjectMutation) ResetMountTo() {
	m.mountTo = nil
}

// SetLoginTo sets the "loginTo" field.
func (m *SubjectMutation) SetLoginTo(b bool) {
	m.loginTo = &b
}

// LoginTo returns the value of the "loginTo" field in the mutation.
func (m *SubjectMutation) LoginTo() (r bool, exists bool) {
	v := m.loginTo
	if v == nil {
		return
	}
	return *v, true
}

// OldLoginTo returns the old "loginTo" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldLoginTo(ctx context.Context) (v *bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLoginTo is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLoginTo requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLoginTo: %w", err)
	}
	return oldValue.LoginTo, nil
}

// ResetLoginTo resets all changes to the "loginTo" field.
func (m *SubjectMutation) ResetLoginTo() {
	m.loginTo = nil
}

// SetRoles sets the "roles" field.
func (m *SubjectMutation) SetRoles(sr schema.SubjectRoles) {
	m.roles = &sr
	m.appendroles = nil
}

// Roles returns the value of the "roles" field in the mutation.
func (m *SubjectMutation) Roles() (r schema.SubjectRoles, exists bool) {
	v := m.roles
	if v == nil {
		return
	}
	return *v, true
}

// OldRoles returns the old "roles" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldRoles(ctx context.Context) (v schema.SubjectRoles, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRoles is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRoles requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRoles: %w", err)
	}
	return oldValue.Roles, nil
}

// AppendRoles adds sr to the "roles" field.
func (m *SubjectMutation) AppendRoles(sr schema.SubjectRoles) {
	m.appendroles = append(m.appendroles, sr...)
}

// AppendedRoles returns the list of values that were appended to the "roles" field in this mutation.
func (m *SubjectMutation) AppendedRoles() (schema.SubjectRoles, bool) {
	if len(m.appendroles) == 0 {
		return nil, false
	}
	return m.appendroles, true
}

// ResetRoles resets all changes to the "roles" field.
func (m *SubjectMutation) ResetRoles() {
	m.roles = nil
	m.appendroles = nil
}

// SetPaths sets the "paths" field.
func (m *SubjectMutation) SetPaths(s []string) {
	m.paths = &s
	m.appendpaths = nil
}

// Paths returns the value of the "paths" field in the mutation.
func (m *SubjectMutation) Paths() (r []string, exists bool) {
	v := m.paths
	if v == nil {
		return
	}
	return *v, true
}

// OldPaths returns the old "paths" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldPaths(ctx context.Context) (v []string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPaths is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPaths requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPaths: %w", err)
	}
	return oldValue.Paths, nil
}

// AppendPaths adds s to the "paths" field.
func (m *SubjectMutation) AppendPaths(s []string) {
	m.appendpaths = append(m.appendpaths, s...)
}

// AppendedPaths returns the list of values that were appended to the "paths" field in this mutation.
func (m *SubjectMutation) AppendedPaths() ([]string, bool) {
	if len(m.appendpaths) == 0 {
		return nil, false
	}
	return m.appendpaths, true
}

// ResetPaths resets all changes to the "paths" field.
func (m *SubjectMutation) ResetPaths() {
	m.paths = nil
	m.appendpaths = nil
}

// SetBuiltin sets the "builtin" field.
func (m *SubjectMutation) SetBuiltin(b bool) {
	m.builtin = &b
}

// Builtin returns the value of the "builtin" field in the mutation.
func (m *SubjectMutation) Builtin() (r bool, exists bool) {
	v := m.builtin
	if v == nil {
		return
	}
	return *v, true
}

// OldBuiltin returns the old "builtin" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldBuiltin(ctx context.Context) (v bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldBuiltin is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldBuiltin requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldBuiltin: %w", err)
	}
	return oldValue.Builtin, nil
}

// ResetBuiltin resets all changes to the "builtin" field.
func (m *SubjectMutation) ResetBuiltin() {
	m.builtin = nil
}

// Where appends a list predicates to the SubjectMutation builder.
func (m *SubjectMutation) Where(ps ...predicate.Subject) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the SubjectMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *SubjectMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Subject, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *SubjectMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *SubjectMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Subject).
func (m *SubjectMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *SubjectMutation) Fields() []string {
	fields := make([]string, 0, 11)
	if m.createTime != nil {
		fields = append(fields, subject.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, subject.FieldUpdateTime)
	}
	if m.kind != nil {
		fields = append(fields, subject.FieldKind)
	}
	if m.group != nil {
		fields = append(fields, subject.FieldGroup)
	}
	if m.name != nil {
		fields = append(fields, subject.FieldName)
	}
	if m.description != nil {
		fields = append(fields, subject.FieldDescription)
	}
	if m.mountTo != nil {
		fields = append(fields, subject.FieldMountTo)
	}
	if m.loginTo != nil {
		fields = append(fields, subject.FieldLoginTo)
	}
	if m.roles != nil {
		fields = append(fields, subject.FieldRoles)
	}
	if m.paths != nil {
		fields = append(fields, subject.FieldPaths)
	}
	if m.builtin != nil {
		fields = append(fields, subject.FieldBuiltin)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *SubjectMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case subject.FieldCreateTime:
		return m.CreateTime()
	case subject.FieldUpdateTime:
		return m.UpdateTime()
	case subject.FieldKind:
		return m.Kind()
	case subject.FieldGroup:
		return m.Group()
	case subject.FieldName:
		return m.Name()
	case subject.FieldDescription:
		return m.Description()
	case subject.FieldMountTo:
		return m.MountTo()
	case subject.FieldLoginTo:
		return m.LoginTo()
	case subject.FieldRoles:
		return m.Roles()
	case subject.FieldPaths:
		return m.Paths()
	case subject.FieldBuiltin:
		return m.Builtin()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *SubjectMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case subject.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case subject.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case subject.FieldKind:
		return m.OldKind(ctx)
	case subject.FieldGroup:
		return m.OldGroup(ctx)
	case subject.FieldName:
		return m.OldName(ctx)
	case subject.FieldDescription:
		return m.OldDescription(ctx)
	case subject.FieldMountTo:
		return m.OldMountTo(ctx)
	case subject.FieldLoginTo:
		return m.OldLoginTo(ctx)
	case subject.FieldRoles:
		return m.OldRoles(ctx)
	case subject.FieldPaths:
		return m.OldPaths(ctx)
	case subject.FieldBuiltin:
		return m.OldBuiltin(ctx)
	}
	return nil, fmt.Errorf("unknown Subject field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SubjectMutation) SetField(name string, value ent.Value) error {
	switch name {
	case subject.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case subject.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case subject.FieldKind:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetKind(v)
		return nil
	case subject.FieldGroup:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetGroup(v)
		return nil
	case subject.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case subject.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case subject.FieldMountTo:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMountTo(v)
		return nil
	case subject.FieldLoginTo:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLoginTo(v)
		return nil
	case subject.FieldRoles:
		v, ok := value.(schema.SubjectRoles)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRoles(v)
		return nil
	case subject.FieldPaths:
		v, ok := value.([]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPaths(v)
		return nil
	case subject.FieldBuiltin:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBuiltin(v)
		return nil
	}
	return fmt.Errorf("unknown Subject field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *SubjectMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *SubjectMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SubjectMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Subject numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *SubjectMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *SubjectMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *SubjectMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Subject nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *SubjectMutation) ResetField(name string) error {
	switch name {
	case subject.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case subject.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case subject.FieldKind:
		m.ResetKind()
		return nil
	case subject.FieldGroup:
		m.ResetGroup()
		return nil
	case subject.FieldName:
		m.ResetName()
		return nil
	case subject.FieldDescription:
		m.ResetDescription()
		return nil
	case subject.FieldMountTo:
		m.ResetMountTo()
		return nil
	case subject.FieldLoginTo:
		m.ResetLoginTo()
		return nil
	case subject.FieldRoles:
		m.ResetRoles()
		return nil
	case subject.FieldPaths:
		m.ResetPaths()
		return nil
	case subject.FieldBuiltin:
		m.ResetBuiltin()
		return nil
	}
	return fmt.Errorf("unknown Subject field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *SubjectMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *SubjectMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *SubjectMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *SubjectMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *SubjectMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *SubjectMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *SubjectMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Subject unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *SubjectMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Subject edge %s", name)
}

// TokenMutation represents an operation that mutates the Token nodes in the graph.
type TokenMutation struct {
	config
	op                Op
	typ               string
	id                *oid.ID
	createTime        *time.Time
	updateTime        *time.Time
	casdoorTokenName  *string
	casdoorTokenOwner *string
	name              *string
	expiration        *int
	addexpiration     *int
	clearedFields     map[string]struct{}
	done              bool
	oldValue          func(context.Context) (*Token, error)
	predicates        []predicate.Token
}

var _ ent.Mutation = (*TokenMutation)(nil)

// tokenOption allows management of the mutation configuration using functional options.
type tokenOption func(*TokenMutation)

// newTokenMutation creates new mutation for the Token entity.
func newTokenMutation(c config, op Op, opts ...tokenOption) *TokenMutation {
	m := &TokenMutation{
		config:        c,
		op:            op,
		typ:           TypeToken,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withTokenID sets the ID field of the mutation.
func withTokenID(id oid.ID) tokenOption {
	return func(m *TokenMutation) {
		var (
			err   error
			once  sync.Once
			value *Token
		)
		m.oldValue = func(ctx context.Context) (*Token, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Token.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withToken sets the old Token of the mutation.
func withToken(node *Token) tokenOption {
	return func(m *TokenMutation) {
		m.oldValue = func(context.Context) (*Token, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m TokenMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m TokenMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Token entities.
func (m *TokenMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *TokenMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *TokenMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Token.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *TokenMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *TokenMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Token entity.
// If the Token object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TokenMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreateTime: %w", err)
	}
	return oldValue.CreateTime, nil
}

// ResetCreateTime resets all changes to the "createTime" field.
func (m *TokenMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *TokenMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *TokenMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Token entity.
// If the Token object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TokenMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdateTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdateTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdateTime: %w", err)
	}
	return oldValue.UpdateTime, nil
}

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *TokenMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetCasdoorTokenName sets the "casdoorTokenName" field.
func (m *TokenMutation) SetCasdoorTokenName(s string) {
	m.casdoorTokenName = &s
}

// CasdoorTokenName returns the value of the "casdoorTokenName" field in the mutation.
func (m *TokenMutation) CasdoorTokenName() (r string, exists bool) {
	v := m.casdoorTokenName
	if v == nil {
		return
	}
	return *v, true
}

// OldCasdoorTokenName returns the old "casdoorTokenName" field's value of the Token entity.
// If the Token object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TokenMutation) OldCasdoorTokenName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCasdoorTokenName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCasdoorTokenName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCasdoorTokenName: %w", err)
	}
	return oldValue.CasdoorTokenName, nil
}

// ResetCasdoorTokenName resets all changes to the "casdoorTokenName" field.
func (m *TokenMutation) ResetCasdoorTokenName() {
	m.casdoorTokenName = nil
}

// SetCasdoorTokenOwner sets the "casdoorTokenOwner" field.
func (m *TokenMutation) SetCasdoorTokenOwner(s string) {
	m.casdoorTokenOwner = &s
}

// CasdoorTokenOwner returns the value of the "casdoorTokenOwner" field in the mutation.
func (m *TokenMutation) CasdoorTokenOwner() (r string, exists bool) {
	v := m.casdoorTokenOwner
	if v == nil {
		return
	}
	return *v, true
}

// OldCasdoorTokenOwner returns the old "casdoorTokenOwner" field's value of the Token entity.
// If the Token object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TokenMutation) OldCasdoorTokenOwner(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCasdoorTokenOwner is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCasdoorTokenOwner requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCasdoorTokenOwner: %w", err)
	}
	return oldValue.CasdoorTokenOwner, nil
}

// ResetCasdoorTokenOwner resets all changes to the "casdoorTokenOwner" field.
func (m *TokenMutation) ResetCasdoorTokenOwner() {
	m.casdoorTokenOwner = nil
}

// SetName sets the "name" field.
func (m *TokenMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *TokenMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Token entity.
// If the Token object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TokenMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *TokenMutation) ResetName() {
	m.name = nil
}

// SetExpiration sets the "expiration" field.
func (m *TokenMutation) SetExpiration(i int) {
	m.expiration = &i
	m.addexpiration = nil
}

// Expiration returns the value of the "expiration" field in the mutation.
func (m *TokenMutation) Expiration() (r int, exists bool) {
	v := m.expiration
	if v == nil {
		return
	}
	return *v, true
}

// OldExpiration returns the old "expiration" field's value of the Token entity.
// If the Token object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TokenMutation) OldExpiration(ctx context.Context) (v *int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldExpiration is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldExpiration requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldExpiration: %w", err)
	}
	return oldValue.Expiration, nil
}

// AddExpiration adds i to the "expiration" field.
func (m *TokenMutation) AddExpiration(i int) {
	if m.addexpiration != nil {
		*m.addexpiration += i
	} else {
		m.addexpiration = &i
	}
}

// AddedExpiration returns the value that was added to the "expiration" field in this mutation.
func (m *TokenMutation) AddedExpiration() (r int, exists bool) {
	v := m.addexpiration
	if v == nil {
		return
	}
	return *v, true
}

// ClearExpiration clears the value of the "expiration" field.
func (m *TokenMutation) ClearExpiration() {
	m.expiration = nil
	m.addexpiration = nil
	m.clearedFields[token.FieldExpiration] = struct{}{}
}

// ExpirationCleared returns if the "expiration" field was cleared in this mutation.
func (m *TokenMutation) ExpirationCleared() bool {
	_, ok := m.clearedFields[token.FieldExpiration]
	return ok
}

// ResetExpiration resets all changes to the "expiration" field.
func (m *TokenMutation) ResetExpiration() {
	m.expiration = nil
	m.addexpiration = nil
	delete(m.clearedFields, token.FieldExpiration)
}

// Where appends a list predicates to the TokenMutation builder.
func (m *TokenMutation) Where(ps ...predicate.Token) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the TokenMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *TokenMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Token, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *TokenMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *TokenMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Token).
func (m *TokenMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *TokenMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.createTime != nil {
		fields = append(fields, token.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, token.FieldUpdateTime)
	}
	if m.casdoorTokenName != nil {
		fields = append(fields, token.FieldCasdoorTokenName)
	}
	if m.casdoorTokenOwner != nil {
		fields = append(fields, token.FieldCasdoorTokenOwner)
	}
	if m.name != nil {
		fields = append(fields, token.FieldName)
	}
	if m.expiration != nil {
		fields = append(fields, token.FieldExpiration)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *TokenMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case token.FieldCreateTime:
		return m.CreateTime()
	case token.FieldUpdateTime:
		return m.UpdateTime()
	case token.FieldCasdoorTokenName:
		return m.CasdoorTokenName()
	case token.FieldCasdoorTokenOwner:
		return m.CasdoorTokenOwner()
	case token.FieldName:
		return m.Name()
	case token.FieldExpiration:
		return m.Expiration()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *TokenMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case token.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case token.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case token.FieldCasdoorTokenName:
		return m.OldCasdoorTokenName(ctx)
	case token.FieldCasdoorTokenOwner:
		return m.OldCasdoorTokenOwner(ctx)
	case token.FieldName:
		return m.OldName(ctx)
	case token.FieldExpiration:
		return m.OldExpiration(ctx)
	}
	return nil, fmt.Errorf("unknown Token field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *TokenMutation) SetField(name string, value ent.Value) error {
	switch name {
	case token.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case token.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case token.FieldCasdoorTokenName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCasdoorTokenName(v)
		return nil
	case token.FieldCasdoorTokenOwner:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCasdoorTokenOwner(v)
		return nil
	case token.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case token.FieldExpiration:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetExpiration(v)
		return nil
	}
	return fmt.Errorf("unknown Token field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *TokenMutation) AddedFields() []string {
	var fields []string
	if m.addexpiration != nil {
		fields = append(fields, token.FieldExpiration)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *TokenMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case token.FieldExpiration:
		return m.AddedExpiration()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *TokenMutation) AddField(name string, value ent.Value) error {
	switch name {
	case token.FieldExpiration:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddExpiration(v)
		return nil
	}
	return fmt.Errorf("unknown Token numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *TokenMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(token.FieldExpiration) {
		fields = append(fields, token.FieldExpiration)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *TokenMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *TokenMutation) ClearField(name string) error {
	switch name {
	case token.FieldExpiration:
		m.ClearExpiration()
		return nil
	}
	return fmt.Errorf("unknown Token nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *TokenMutation) ResetField(name string) error {
	switch name {
	case token.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case token.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case token.FieldCasdoorTokenName:
		m.ResetCasdoorTokenName()
		return nil
	case token.FieldCasdoorTokenOwner:
		m.ResetCasdoorTokenOwner()
		return nil
	case token.FieldName:
		m.ResetName()
		return nil
	case token.FieldExpiration:
		m.ResetExpiration()
		return nil
	}
	return fmt.Errorf("unknown Token field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *TokenMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *TokenMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *TokenMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *TokenMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *TokenMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *TokenMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *TokenMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Token unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *TokenMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Token edge %s", name)
}
