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

	"github.com/seal-io/seal/pkg/dao/model/predicate"
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
	TypeRole    = "Role"
	TypeSetting = "Setting"
	TypeSubject = "Subject"
	TypeToken   = "Token"
)

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
