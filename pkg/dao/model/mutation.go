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
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/schema"
	"github.com/seal-io/seal/pkg/dao/types"

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
	TypeApplication                      = "Application"
	TypeApplicationModuleRelationship    = "ApplicationModuleRelationship"
	TypeApplicationResource              = "ApplicationResource"
	TypeApplicationRevision              = "ApplicationRevision"
	TypeConnector                        = "Connector"
	TypeEnvironment                      = "Environment"
	TypeEnvironmentConnectorRelationship = "EnvironmentConnectorRelationship"
	TypeModule                           = "Module"
	TypeProject                          = "Project"
	TypeRole                             = "Role"
	TypeSetting                          = "Setting"
	TypeSubject                          = "Subject"
	TypeToken                            = "Token"
)

// ApplicationMutation represents an operation that mutates the Application nodes in the graph.
type ApplicationMutation struct {
	config
	op                                    Op
	typ                                   string
	id                                    *types.ID
	name                                  *string
	description                           *string
	labels                                *map[string]string
	createTime                            *time.Time
	updateTime                            *time.Time
	clearedFields                         map[string]struct{}
	project                               *types.ID
	clearedproject                        bool
	environment                           *types.ID
	clearedenvironment                    bool
	resources                             map[types.ID]struct{}
	removedresources                      map[types.ID]struct{}
	clearedresources                      bool
	revisions                             map[types.ID]struct{}
	removedrevisions                      map[types.ID]struct{}
	clearedrevisions                      bool
	modules                               map[string]struct{}
	removedmodules                        map[string]struct{}
	clearedmodules                        bool
	applicationModuleRelationships        map[int]struct{}
	removedapplicationModuleRelationships map[int]struct{}
	clearedapplicationModuleRelationships bool
	done                                  bool
	oldValue                              func(context.Context) (*Application, error)
	predicates                            []predicate.Application
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
func withApplicationID(id types.ID) applicationOption {
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
func (m *ApplicationMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ApplicationMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ApplicationMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Application.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *ApplicationMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *ApplicationMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *ApplicationMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the "description" field.
func (m *ApplicationMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the value of the "description" field in the mutation.
func (m *ApplicationMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// OldDescription returns the old "description" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldDescription(ctx context.Context) (v string, err error) {
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

// ClearDescription clears the value of the "description" field.
func (m *ApplicationMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[application.FieldDescription] = struct{}{}
}

// DescriptionCleared returns if the "description" field was cleared in this mutation.
func (m *ApplicationMutation) DescriptionCleared() bool {
	_, ok := m.clearedFields[application.FieldDescription]
	return ok
}

// ResetDescription resets all changes to the "description" field.
func (m *ApplicationMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, application.FieldDescription)
}

// SetLabels sets the "labels" field.
func (m *ApplicationMutation) SetLabels(value map[string]string) {
	m.labels = &value
}

// Labels returns the value of the "labels" field in the mutation.
func (m *ApplicationMutation) Labels() (r map[string]string, exists bool) {
	v := m.labels
	if v == nil {
		return
	}
	return *v, true
}

// OldLabels returns the old "labels" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldLabels(ctx context.Context) (v map[string]string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLabels is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLabels requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLabels: %w", err)
	}
	return oldValue.Labels, nil
}

// ResetLabels resets all changes to the "labels" field.
func (m *ApplicationMutation) ResetLabels() {
	m.labels = nil
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
func (m *ApplicationMutation) SetProjectID(t types.ID) {
	m.project = &t
}

// ProjectID returns the value of the "projectID" field in the mutation.
func (m *ApplicationMutation) ProjectID() (r types.ID, exists bool) {
	v := m.project
	if v == nil {
		return
	}
	return *v, true
}

// OldProjectID returns the old "projectID" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldProjectID(ctx context.Context) (v types.ID, err error) {
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
	m.project = nil
}

// SetEnvironmentID sets the "environmentID" field.
func (m *ApplicationMutation) SetEnvironmentID(t types.ID) {
	m.environment = &t
}

// EnvironmentID returns the value of the "environmentID" field in the mutation.
func (m *ApplicationMutation) EnvironmentID() (r types.ID, exists bool) {
	v := m.environment
	if v == nil {
		return
	}
	return *v, true
}

// OldEnvironmentID returns the old "environmentID" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldEnvironmentID(ctx context.Context) (v types.ID, err error) {
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
	m.environment = nil
}

// ClearProject clears the "project" edge to the Project entity.
func (m *ApplicationMutation) ClearProject() {
	m.clearedproject = true
}

// ProjectCleared reports if the "project" edge to the Project entity was cleared.
func (m *ApplicationMutation) ProjectCleared() bool {
	return m.clearedproject
}

// ProjectIDs returns the "project" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ProjectID instead. It exists only for internal usage by the builders.
func (m *ApplicationMutation) ProjectIDs() (ids []types.ID) {
	if id := m.project; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetProject resets all changes to the "project" edge.
func (m *ApplicationMutation) ResetProject() {
	m.project = nil
	m.clearedproject = false
}

// ClearEnvironment clears the "environment" edge to the Environment entity.
func (m *ApplicationMutation) ClearEnvironment() {
	m.clearedenvironment = true
}

// EnvironmentCleared reports if the "environment" edge to the Environment entity was cleared.
func (m *ApplicationMutation) EnvironmentCleared() bool {
	return m.clearedenvironment
}

// EnvironmentIDs returns the "environment" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// EnvironmentID instead. It exists only for internal usage by the builders.
func (m *ApplicationMutation) EnvironmentIDs() (ids []types.ID) {
	if id := m.environment; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEnvironment resets all changes to the "environment" edge.
func (m *ApplicationMutation) ResetEnvironment() {
	m.environment = nil
	m.clearedenvironment = false
}

// AddResourceIDs adds the "resources" edge to the ApplicationResource entity by ids.
func (m *ApplicationMutation) AddResourceIDs(ids ...types.ID) {
	if m.resources == nil {
		m.resources = make(map[types.ID]struct{})
	}
	for i := range ids {
		m.resources[ids[i]] = struct{}{}
	}
}

// ClearResources clears the "resources" edge to the ApplicationResource entity.
func (m *ApplicationMutation) ClearResources() {
	m.clearedresources = true
}

// ResourcesCleared reports if the "resources" edge to the ApplicationResource entity was cleared.
func (m *ApplicationMutation) ResourcesCleared() bool {
	return m.clearedresources
}

// RemoveResourceIDs removes the "resources" edge to the ApplicationResource entity by IDs.
func (m *ApplicationMutation) RemoveResourceIDs(ids ...types.ID) {
	if m.removedresources == nil {
		m.removedresources = make(map[types.ID]struct{})
	}
	for i := range ids {
		delete(m.resources, ids[i])
		m.removedresources[ids[i]] = struct{}{}
	}
}

// RemovedResources returns the removed IDs of the "resources" edge to the ApplicationResource entity.
func (m *ApplicationMutation) RemovedResourcesIDs() (ids []types.ID) {
	for id := range m.removedresources {
		ids = append(ids, id)
	}
	return
}

// ResourcesIDs returns the "resources" edge IDs in the mutation.
func (m *ApplicationMutation) ResourcesIDs() (ids []types.ID) {
	for id := range m.resources {
		ids = append(ids, id)
	}
	return
}

// ResetResources resets all changes to the "resources" edge.
func (m *ApplicationMutation) ResetResources() {
	m.resources = nil
	m.clearedresources = false
	m.removedresources = nil
}

// AddRevisionIDs adds the "revisions" edge to the ApplicationRevision entity by ids.
func (m *ApplicationMutation) AddRevisionIDs(ids ...types.ID) {
	if m.revisions == nil {
		m.revisions = make(map[types.ID]struct{})
	}
	for i := range ids {
		m.revisions[ids[i]] = struct{}{}
	}
}

// ClearRevisions clears the "revisions" edge to the ApplicationRevision entity.
func (m *ApplicationMutation) ClearRevisions() {
	m.clearedrevisions = true
}

// RevisionsCleared reports if the "revisions" edge to the ApplicationRevision entity was cleared.
func (m *ApplicationMutation) RevisionsCleared() bool {
	return m.clearedrevisions
}

// RemoveRevisionIDs removes the "revisions" edge to the ApplicationRevision entity by IDs.
func (m *ApplicationMutation) RemoveRevisionIDs(ids ...types.ID) {
	if m.removedrevisions == nil {
		m.removedrevisions = make(map[types.ID]struct{})
	}
	for i := range ids {
		delete(m.revisions, ids[i])
		m.removedrevisions[ids[i]] = struct{}{}
	}
}

// RemovedRevisions returns the removed IDs of the "revisions" edge to the ApplicationRevision entity.
func (m *ApplicationMutation) RemovedRevisionsIDs() (ids []types.ID) {
	for id := range m.removedrevisions {
		ids = append(ids, id)
	}
	return
}

// RevisionsIDs returns the "revisions" edge IDs in the mutation.
func (m *ApplicationMutation) RevisionsIDs() (ids []types.ID) {
	for id := range m.revisions {
		ids = append(ids, id)
	}
	return
}

// ResetRevisions resets all changes to the "revisions" edge.
func (m *ApplicationMutation) ResetRevisions() {
	m.revisions = nil
	m.clearedrevisions = false
	m.removedrevisions = nil
}

// AddModuleIDs adds the "modules" edge to the Module entity by ids.
func (m *ApplicationMutation) AddModuleIDs(ids ...string) {
	if m.modules == nil {
		m.modules = make(map[string]struct{})
	}
	for i := range ids {
		m.modules[ids[i]] = struct{}{}
	}
}

// ClearModules clears the "modules" edge to the Module entity.
func (m *ApplicationMutation) ClearModules() {
	m.clearedmodules = true
}

// ModulesCleared reports if the "modules" edge to the Module entity was cleared.
func (m *ApplicationMutation) ModulesCleared() bool {
	return m.clearedmodules
}

// RemoveModuleIDs removes the "modules" edge to the Module entity by IDs.
func (m *ApplicationMutation) RemoveModuleIDs(ids ...string) {
	if m.removedmodules == nil {
		m.removedmodules = make(map[string]struct{})
	}
	for i := range ids {
		delete(m.modules, ids[i])
		m.removedmodules[ids[i]] = struct{}{}
	}
}

// RemovedModules returns the removed IDs of the "modules" edge to the Module entity.
func (m *ApplicationMutation) RemovedModulesIDs() (ids []string) {
	for id := range m.removedmodules {
		ids = append(ids, id)
	}
	return
}

// ModulesIDs returns the "modules" edge IDs in the mutation.
func (m *ApplicationMutation) ModulesIDs() (ids []string) {
	for id := range m.modules {
		ids = append(ids, id)
	}
	return
}

// ResetModules resets all changes to the "modules" edge.
func (m *ApplicationMutation) ResetModules() {
	m.modules = nil
	m.clearedmodules = false
	m.removedmodules = nil
}

// AddApplicationModuleRelationshipIDs adds the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity by ids.
func (m *ApplicationMutation) AddApplicationModuleRelationshipIDs(ids ...int) {
	if m.applicationModuleRelationships == nil {
		m.applicationModuleRelationships = make(map[int]struct{})
	}
	for i := range ids {
		m.applicationModuleRelationships[ids[i]] = struct{}{}
	}
}

// ClearApplicationModuleRelationships clears the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity.
func (m *ApplicationMutation) ClearApplicationModuleRelationships() {
	m.clearedapplicationModuleRelationships = true
}

// ApplicationModuleRelationshipsCleared reports if the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity was cleared.
func (m *ApplicationMutation) ApplicationModuleRelationshipsCleared() bool {
	return m.clearedapplicationModuleRelationships
}

// RemoveApplicationModuleRelationshipIDs removes the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity by IDs.
func (m *ApplicationMutation) RemoveApplicationModuleRelationshipIDs(ids ...int) {
	if m.removedapplicationModuleRelationships == nil {
		m.removedapplicationModuleRelationships = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.applicationModuleRelationships, ids[i])
		m.removedapplicationModuleRelationships[ids[i]] = struct{}{}
	}
}

// RemovedApplicationModuleRelationships returns the removed IDs of the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity.
func (m *ApplicationMutation) RemovedApplicationModuleRelationshipsIDs() (ids []int) {
	for id := range m.removedapplicationModuleRelationships {
		ids = append(ids, id)
	}
	return
}

// ApplicationModuleRelationshipsIDs returns the "applicationModuleRelationships" edge IDs in the mutation.
func (m *ApplicationMutation) ApplicationModuleRelationshipsIDs() (ids []int) {
	for id := range m.applicationModuleRelationships {
		ids = append(ids, id)
	}
	return
}

// ResetApplicationModuleRelationships resets all changes to the "applicationModuleRelationships" edge.
func (m *ApplicationMutation) ResetApplicationModuleRelationships() {
	m.applicationModuleRelationships = nil
	m.clearedapplicationModuleRelationships = false
	m.removedapplicationModuleRelationships = nil
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
	fields := make([]string, 0, 7)
	if m.name != nil {
		fields = append(fields, application.FieldName)
	}
	if m.description != nil {
		fields = append(fields, application.FieldDescription)
	}
	if m.labels != nil {
		fields = append(fields, application.FieldLabels)
	}
	if m.createTime != nil {
		fields = append(fields, application.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, application.FieldUpdateTime)
	}
	if m.project != nil {
		fields = append(fields, application.FieldProjectID)
	}
	if m.environment != nil {
		fields = append(fields, application.FieldEnvironmentID)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ApplicationMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case application.FieldName:
		return m.Name()
	case application.FieldDescription:
		return m.Description()
	case application.FieldLabels:
		return m.Labels()
	case application.FieldCreateTime:
		return m.CreateTime()
	case application.FieldUpdateTime:
		return m.UpdateTime()
	case application.FieldProjectID:
		return m.ProjectID()
	case application.FieldEnvironmentID:
		return m.EnvironmentID()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ApplicationMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case application.FieldName:
		return m.OldName(ctx)
	case application.FieldDescription:
		return m.OldDescription(ctx)
	case application.FieldLabels:
		return m.OldLabels(ctx)
	case application.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case application.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case application.FieldProjectID:
		return m.OldProjectID(ctx)
	case application.FieldEnvironmentID:
		return m.OldEnvironmentID(ctx)
	}
	return nil, fmt.Errorf("unknown Application field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationMutation) SetField(name string, value ent.Value) error {
	switch name {
	case application.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case application.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case application.FieldLabels:
		v, ok := value.(map[string]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLabels(v)
		return nil
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
		v, ok := value.(types.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetProjectID(v)
		return nil
	case application.FieldEnvironmentID:
		v, ok := value.(types.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnvironmentID(v)
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
	var fields []string
	if m.FieldCleared(application.FieldDescription) {
		fields = append(fields, application.FieldDescription)
	}
	return fields
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
	switch name {
	case application.FieldDescription:
		m.ClearDescription()
		return nil
	}
	return fmt.Errorf("unknown Application nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ApplicationMutation) ResetField(name string) error {
	switch name {
	case application.FieldName:
		m.ResetName()
		return nil
	case application.FieldDescription:
		m.ResetDescription()
		return nil
	case application.FieldLabels:
		m.ResetLabels()
		return nil
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
	}
	return fmt.Errorf("unknown Application field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationMutation) AddedEdges() []string {
	edges := make([]string, 0, 6)
	if m.project != nil {
		edges = append(edges, application.EdgeProject)
	}
	if m.environment != nil {
		edges = append(edges, application.EdgeEnvironment)
	}
	if m.resources != nil {
		edges = append(edges, application.EdgeResources)
	}
	if m.revisions != nil {
		edges = append(edges, application.EdgeRevisions)
	}
	if m.modules != nil {
		edges = append(edges, application.EdgeModules)
	}
	if m.applicationModuleRelationships != nil {
		edges = append(edges, application.EdgeApplicationModuleRelationships)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ApplicationMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case application.EdgeProject:
		if id := m.project; id != nil {
			return []ent.Value{*id}
		}
	case application.EdgeEnvironment:
		if id := m.environment; id != nil {
			return []ent.Value{*id}
		}
	case application.EdgeResources:
		ids := make([]ent.Value, 0, len(m.resources))
		for id := range m.resources {
			ids = append(ids, id)
		}
		return ids
	case application.EdgeRevisions:
		ids := make([]ent.Value, 0, len(m.revisions))
		for id := range m.revisions {
			ids = append(ids, id)
		}
		return ids
	case application.EdgeModules:
		ids := make([]ent.Value, 0, len(m.modules))
		for id := range m.modules {
			ids = append(ids, id)
		}
		return ids
	case application.EdgeApplicationModuleRelationships:
		ids := make([]ent.Value, 0, len(m.applicationModuleRelationships))
		for id := range m.applicationModuleRelationships {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationMutation) RemovedEdges() []string {
	edges := make([]string, 0, 6)
	if m.removedresources != nil {
		edges = append(edges, application.EdgeResources)
	}
	if m.removedrevisions != nil {
		edges = append(edges, application.EdgeRevisions)
	}
	if m.removedmodules != nil {
		edges = append(edges, application.EdgeModules)
	}
	if m.removedapplicationModuleRelationships != nil {
		edges = append(edges, application.EdgeApplicationModuleRelationships)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case application.EdgeResources:
		ids := make([]ent.Value, 0, len(m.removedresources))
		for id := range m.removedresources {
			ids = append(ids, id)
		}
		return ids
	case application.EdgeRevisions:
		ids := make([]ent.Value, 0, len(m.removedrevisions))
		for id := range m.removedrevisions {
			ids = append(ids, id)
		}
		return ids
	case application.EdgeModules:
		ids := make([]ent.Value, 0, len(m.removedmodules))
		for id := range m.removedmodules {
			ids = append(ids, id)
		}
		return ids
	case application.EdgeApplicationModuleRelationships:
		ids := make([]ent.Value, 0, len(m.removedapplicationModuleRelationships))
		for id := range m.removedapplicationModuleRelationships {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationMutation) ClearedEdges() []string {
	edges := make([]string, 0, 6)
	if m.clearedproject {
		edges = append(edges, application.EdgeProject)
	}
	if m.clearedenvironment {
		edges = append(edges, application.EdgeEnvironment)
	}
	if m.clearedresources {
		edges = append(edges, application.EdgeResources)
	}
	if m.clearedrevisions {
		edges = append(edges, application.EdgeRevisions)
	}
	if m.clearedmodules {
		edges = append(edges, application.EdgeModules)
	}
	if m.clearedapplicationModuleRelationships {
		edges = append(edges, application.EdgeApplicationModuleRelationships)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationMutation) EdgeCleared(name string) bool {
	switch name {
	case application.EdgeProject:
		return m.clearedproject
	case application.EdgeEnvironment:
		return m.clearedenvironment
	case application.EdgeResources:
		return m.clearedresources
	case application.EdgeRevisions:
		return m.clearedrevisions
	case application.EdgeModules:
		return m.clearedmodules
	case application.EdgeApplicationModuleRelationships:
		return m.clearedapplicationModuleRelationships
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationMutation) ClearEdge(name string) error {
	switch name {
	case application.EdgeProject:
		m.ClearProject()
		return nil
	case application.EdgeEnvironment:
		m.ClearEnvironment()
		return nil
	}
	return fmt.Errorf("unknown Application unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ApplicationMutation) ResetEdge(name string) error {
	switch name {
	case application.EdgeProject:
		m.ResetProject()
		return nil
	case application.EdgeEnvironment:
		m.ResetEnvironment()
		return nil
	case application.EdgeResources:
		m.ResetResources()
		return nil
	case application.EdgeRevisions:
		m.ResetRevisions()
		return nil
	case application.EdgeModules:
		m.ResetModules()
		return nil
	case application.EdgeApplicationModuleRelationships:
		m.ResetApplicationModuleRelationships()
		return nil
	}
	return fmt.Errorf("unknown Application edge %s", name)
}

// ApplicationModuleRelationshipMutation represents an operation that mutates the ApplicationModuleRelationship nodes in the graph.
type ApplicationModuleRelationshipMutation struct {
	config
	op                 Op
	typ                string
	id                 *int
	createTime         *time.Time
	updateTime         *time.Time
	name               *string
	variables          *map[string]interface{}
	clearedFields      map[string]struct{}
	application        *types.ID
	clearedapplication bool
	module             *string
	clearedmodule      bool
	done               bool
	oldValue           func(context.Context) (*ApplicationModuleRelationship, error)
	predicates         []predicate.ApplicationModuleRelationship
}

var _ ent.Mutation = (*ApplicationModuleRelationshipMutation)(nil)

// applicationmodulerelationshipOption allows management of the mutation configuration using functional options.
type applicationmodulerelationshipOption func(*ApplicationModuleRelationshipMutation)

// newApplicationModuleRelationshipMutation creates new mutation for the ApplicationModuleRelationship entity.
func newApplicationModuleRelationshipMutation(c config, op Op, opts ...applicationmodulerelationshipOption) *ApplicationModuleRelationshipMutation {
	m := &ApplicationModuleRelationshipMutation{
		config:        c,
		op:            op,
		typ:           TypeApplicationModuleRelationship,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withApplicationModuleRelationshipID sets the ID field of the mutation.
func withApplicationModuleRelationshipID(id int) applicationmodulerelationshipOption {
	return func(m *ApplicationModuleRelationshipMutation) {
		var (
			err   error
			once  sync.Once
			value *ApplicationModuleRelationship
		)
		m.oldValue = func(ctx context.Context) (*ApplicationModuleRelationship, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().ApplicationModuleRelationship.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withApplicationModuleRelationship sets the old ApplicationModuleRelationship of the mutation.
func withApplicationModuleRelationship(node *ApplicationModuleRelationship) applicationmodulerelationshipOption {
	return func(m *ApplicationModuleRelationshipMutation) {
		m.oldValue = func(context.Context) (*ApplicationModuleRelationship, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ApplicationModuleRelationshipMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ApplicationModuleRelationshipMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ApplicationModuleRelationshipMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ApplicationModuleRelationshipMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().ApplicationModuleRelationship.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *ApplicationModuleRelationshipMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *ApplicationModuleRelationshipMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the ApplicationModuleRelationship entity.
// If the ApplicationModuleRelationship object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationModuleRelationshipMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *ApplicationModuleRelationshipMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *ApplicationModuleRelationshipMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *ApplicationModuleRelationshipMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the ApplicationModuleRelationship entity.
// If the ApplicationModuleRelationship object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationModuleRelationshipMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *ApplicationModuleRelationshipMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetApplicationID sets the "application_id" field.
func (m *ApplicationModuleRelationshipMutation) SetApplicationID(t types.ID) {
	m.application = &t
}

// ApplicationID returns the value of the "application_id" field in the mutation.
func (m *ApplicationModuleRelationshipMutation) ApplicationID() (r types.ID, exists bool) {
	v := m.application
	if v == nil {
		return
	}
	return *v, true
}

// OldApplicationID returns the old "application_id" field's value of the ApplicationModuleRelationship entity.
// If the ApplicationModuleRelationship object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationModuleRelationshipMutation) OldApplicationID(ctx context.Context) (v types.ID, err error) {
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

// ResetApplicationID resets all changes to the "application_id" field.
func (m *ApplicationModuleRelationshipMutation) ResetApplicationID() {
	m.application = nil
}

// SetModuleID sets the "module_id" field.
func (m *ApplicationModuleRelationshipMutation) SetModuleID(s string) {
	m.module = &s
}

// ModuleID returns the value of the "module_id" field in the mutation.
func (m *ApplicationModuleRelationshipMutation) ModuleID() (r string, exists bool) {
	v := m.module
	if v == nil {
		return
	}
	return *v, true
}

// OldModuleID returns the old "module_id" field's value of the ApplicationModuleRelationship entity.
// If the ApplicationModuleRelationship object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationModuleRelationshipMutation) OldModuleID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldModuleID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldModuleID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldModuleID: %w", err)
	}
	return oldValue.ModuleID, nil
}

// ResetModuleID resets all changes to the "module_id" field.
func (m *ApplicationModuleRelationshipMutation) ResetModuleID() {
	m.module = nil
}

// SetName sets the "name" field.
func (m *ApplicationModuleRelationshipMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *ApplicationModuleRelationshipMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the ApplicationModuleRelationship entity.
// If the ApplicationModuleRelationship object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationModuleRelationshipMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *ApplicationModuleRelationshipMutation) ResetName() {
	m.name = nil
}

// SetVariables sets the "variables" field.
func (m *ApplicationModuleRelationshipMutation) SetVariables(value map[string]interface{}) {
	m.variables = &value
}

// Variables returns the value of the "variables" field in the mutation.
func (m *ApplicationModuleRelationshipMutation) Variables() (r map[string]interface{}, exists bool) {
	v := m.variables
	if v == nil {
		return
	}
	return *v, true
}

// OldVariables returns the old "variables" field's value of the ApplicationModuleRelationship entity.
// If the ApplicationModuleRelationship object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationModuleRelationshipMutation) OldVariables(ctx context.Context) (v map[string]interface{}, err error) {
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
func (m *ApplicationModuleRelationshipMutation) ClearVariables() {
	m.variables = nil
	m.clearedFields[applicationmodulerelationship.FieldVariables] = struct{}{}
}

// VariablesCleared returns if the "variables" field was cleared in this mutation.
func (m *ApplicationModuleRelationshipMutation) VariablesCleared() bool {
	_, ok := m.clearedFields[applicationmodulerelationship.FieldVariables]
	return ok
}

// ResetVariables resets all changes to the "variables" field.
func (m *ApplicationModuleRelationshipMutation) ResetVariables() {
	m.variables = nil
	delete(m.clearedFields, applicationmodulerelationship.FieldVariables)
}

// ClearApplication clears the "application" edge to the Application entity.
func (m *ApplicationModuleRelationshipMutation) ClearApplication() {
	m.clearedapplication = true
}

// ApplicationCleared reports if the "application" edge to the Application entity was cleared.
func (m *ApplicationModuleRelationshipMutation) ApplicationCleared() bool {
	return m.clearedapplication
}

// ApplicationIDs returns the "application" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ApplicationID instead. It exists only for internal usage by the builders.
func (m *ApplicationModuleRelationshipMutation) ApplicationIDs() (ids []types.ID) {
	if id := m.application; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetApplication resets all changes to the "application" edge.
func (m *ApplicationModuleRelationshipMutation) ResetApplication() {
	m.application = nil
	m.clearedapplication = false
}

// ClearModule clears the "module" edge to the Module entity.
func (m *ApplicationModuleRelationshipMutation) ClearModule() {
	m.clearedmodule = true
}

// ModuleCleared reports if the "module" edge to the Module entity was cleared.
func (m *ApplicationModuleRelationshipMutation) ModuleCleared() bool {
	return m.clearedmodule
}

// ModuleIDs returns the "module" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ModuleID instead. It exists only for internal usage by the builders.
func (m *ApplicationModuleRelationshipMutation) ModuleIDs() (ids []string) {
	if id := m.module; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetModule resets all changes to the "module" edge.
func (m *ApplicationModuleRelationshipMutation) ResetModule() {
	m.module = nil
	m.clearedmodule = false
}

// Where appends a list predicates to the ApplicationModuleRelationshipMutation builder.
func (m *ApplicationModuleRelationshipMutation) Where(ps ...predicate.ApplicationModuleRelationship) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ApplicationModuleRelationshipMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ApplicationModuleRelationshipMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.ApplicationModuleRelationship, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ApplicationModuleRelationshipMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ApplicationModuleRelationshipMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (ApplicationModuleRelationship).
func (m *ApplicationModuleRelationshipMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ApplicationModuleRelationshipMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.createTime != nil {
		fields = append(fields, applicationmodulerelationship.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, applicationmodulerelationship.FieldUpdateTime)
	}
	if m.application != nil {
		fields = append(fields, applicationmodulerelationship.FieldApplicationID)
	}
	if m.module != nil {
		fields = append(fields, applicationmodulerelationship.FieldModuleID)
	}
	if m.name != nil {
		fields = append(fields, applicationmodulerelationship.FieldName)
	}
	if m.variables != nil {
		fields = append(fields, applicationmodulerelationship.FieldVariables)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ApplicationModuleRelationshipMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case applicationmodulerelationship.FieldCreateTime:
		return m.CreateTime()
	case applicationmodulerelationship.FieldUpdateTime:
		return m.UpdateTime()
	case applicationmodulerelationship.FieldApplicationID:
		return m.ApplicationID()
	case applicationmodulerelationship.FieldModuleID:
		return m.ModuleID()
	case applicationmodulerelationship.FieldName:
		return m.Name()
	case applicationmodulerelationship.FieldVariables:
		return m.Variables()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ApplicationModuleRelationshipMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case applicationmodulerelationship.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case applicationmodulerelationship.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case applicationmodulerelationship.FieldApplicationID:
		return m.OldApplicationID(ctx)
	case applicationmodulerelationship.FieldModuleID:
		return m.OldModuleID(ctx)
	case applicationmodulerelationship.FieldName:
		return m.OldName(ctx)
	case applicationmodulerelationship.FieldVariables:
		return m.OldVariables(ctx)
	}
	return nil, fmt.Errorf("unknown ApplicationModuleRelationship field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationModuleRelationshipMutation) SetField(name string, value ent.Value) error {
	switch name {
	case applicationmodulerelationship.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case applicationmodulerelationship.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case applicationmodulerelationship.FieldApplicationID:
		v, ok := value.(types.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetApplicationID(v)
		return nil
	case applicationmodulerelationship.FieldModuleID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetModuleID(v)
		return nil
	case applicationmodulerelationship.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case applicationmodulerelationship.FieldVariables:
		v, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVariables(v)
		return nil
	}
	return fmt.Errorf("unknown ApplicationModuleRelationship field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ApplicationModuleRelationshipMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ApplicationModuleRelationshipMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationModuleRelationshipMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ApplicationModuleRelationship numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ApplicationModuleRelationshipMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(applicationmodulerelationship.FieldVariables) {
		fields = append(fields, applicationmodulerelationship.FieldVariables)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ApplicationModuleRelationshipMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ApplicationModuleRelationshipMutation) ClearField(name string) error {
	switch name {
	case applicationmodulerelationship.FieldVariables:
		m.ClearVariables()
		return nil
	}
	return fmt.Errorf("unknown ApplicationModuleRelationship nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ApplicationModuleRelationshipMutation) ResetField(name string) error {
	switch name {
	case applicationmodulerelationship.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case applicationmodulerelationship.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case applicationmodulerelationship.FieldApplicationID:
		m.ResetApplicationID()
		return nil
	case applicationmodulerelationship.FieldModuleID:
		m.ResetModuleID()
		return nil
	case applicationmodulerelationship.FieldName:
		m.ResetName()
		return nil
	case applicationmodulerelationship.FieldVariables:
		m.ResetVariables()
		return nil
	}
	return fmt.Errorf("unknown ApplicationModuleRelationship field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationModuleRelationshipMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.application != nil {
		edges = append(edges, applicationmodulerelationship.EdgeApplication)
	}
	if m.module != nil {
		edges = append(edges, applicationmodulerelationship.EdgeModule)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ApplicationModuleRelationshipMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case applicationmodulerelationship.EdgeApplication:
		if id := m.application; id != nil {
			return []ent.Value{*id}
		}
	case applicationmodulerelationship.EdgeModule:
		if id := m.module; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationModuleRelationshipMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationModuleRelationshipMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationModuleRelationshipMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedapplication {
		edges = append(edges, applicationmodulerelationship.EdgeApplication)
	}
	if m.clearedmodule {
		edges = append(edges, applicationmodulerelationship.EdgeModule)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationModuleRelationshipMutation) EdgeCleared(name string) bool {
	switch name {
	case applicationmodulerelationship.EdgeApplication:
		return m.clearedapplication
	case applicationmodulerelationship.EdgeModule:
		return m.clearedmodule
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationModuleRelationshipMutation) ClearEdge(name string) error {
	switch name {
	case applicationmodulerelationship.EdgeApplication:
		m.ClearApplication()
		return nil
	case applicationmodulerelationship.EdgeModule:
		m.ClearModule()
		return nil
	}
	return fmt.Errorf("unknown ApplicationModuleRelationship unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ApplicationModuleRelationshipMutation) ResetEdge(name string) error {
	switch name {
	case applicationmodulerelationship.EdgeApplication:
		m.ResetApplication()
		return nil
	case applicationmodulerelationship.EdgeModule:
		m.ResetModule()
		return nil
	}
	return fmt.Errorf("unknown ApplicationModuleRelationship edge %s", name)
}

// ApplicationResourceMutation represents an operation that mutates the ApplicationResource nodes in the graph.
type ApplicationResourceMutation struct {
	config
	op                 Op
	typ                string
	id                 *types.ID
	status             *string
	statusMessage      *string
	createTime         *time.Time
	updateTime         *time.Time
	module             *string
	mode               *string
	_type              *string
	name               *string
	clearedFields      map[string]struct{}
	application        *types.ID
	clearedapplication bool
	connector          *types.ID
	clearedconnector   bool
	done               bool
	oldValue           func(context.Context) (*ApplicationResource, error)
	predicates         []predicate.ApplicationResource
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
func withApplicationResourceID(id types.ID) applicationresourceOption {
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
func (m *ApplicationResourceMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ApplicationResourceMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ApplicationResourceMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
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
func (m *ApplicationResourceMutation) SetApplicationID(t types.ID) {
	m.application = &t
}

// ApplicationID returns the value of the "applicationID" field in the mutation.
func (m *ApplicationResourceMutation) ApplicationID() (r types.ID, exists bool) {
	v := m.application
	if v == nil {
		return
	}
	return *v, true
}

// OldApplicationID returns the old "applicationID" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldApplicationID(ctx context.Context) (v types.ID, err error) {
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
	m.application = nil
}

// SetConnectorID sets the "connectorID" field.
func (m *ApplicationResourceMutation) SetConnectorID(t types.ID) {
	m.connector = &t
}

// ConnectorID returns the value of the "connectorID" field in the mutation.
func (m *ApplicationResourceMutation) ConnectorID() (r types.ID, exists bool) {
	v := m.connector
	if v == nil {
		return
	}
	return *v, true
}

// OldConnectorID returns the old "connectorID" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldConnectorID(ctx context.Context) (v types.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldConnectorID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldConnectorID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldConnectorID: %w", err)
	}
	return oldValue.ConnectorID, nil
}

// ResetConnectorID resets all changes to the "connectorID" field.
func (m *ApplicationResourceMutation) ResetConnectorID() {
	m.connector = nil
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

// SetMode sets the "mode" field.
func (m *ApplicationResourceMutation) SetMode(s string) {
	m.mode = &s
}

// Mode returns the value of the "mode" field in the mutation.
func (m *ApplicationResourceMutation) Mode() (r string, exists bool) {
	v := m.mode
	if v == nil {
		return
	}
	return *v, true
}

// OldMode returns the old "mode" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldMode(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldMode is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldMode requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldMode: %w", err)
	}
	return oldValue.Mode, nil
}

// ResetMode resets all changes to the "mode" field.
func (m *ApplicationResourceMutation) ResetMode() {
	m.mode = nil
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

// SetName sets the "name" field.
func (m *ApplicationResourceMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *ApplicationResourceMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *ApplicationResourceMutation) ResetName() {
	m.name = nil
}

// ClearApplication clears the "application" edge to the Application entity.
func (m *ApplicationResourceMutation) ClearApplication() {
	m.clearedapplication = true
}

// ApplicationCleared reports if the "application" edge to the Application entity was cleared.
func (m *ApplicationResourceMutation) ApplicationCleared() bool {
	return m.clearedapplication
}

// ApplicationIDs returns the "application" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ApplicationID instead. It exists only for internal usage by the builders.
func (m *ApplicationResourceMutation) ApplicationIDs() (ids []types.ID) {
	if id := m.application; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetApplication resets all changes to the "application" edge.
func (m *ApplicationResourceMutation) ResetApplication() {
	m.application = nil
	m.clearedapplication = false
}

// ClearConnector clears the "connector" edge to the Connector entity.
func (m *ApplicationResourceMutation) ClearConnector() {
	m.clearedconnector = true
}

// ConnectorCleared reports if the "connector" edge to the Connector entity was cleared.
func (m *ApplicationResourceMutation) ConnectorCleared() bool {
	return m.clearedconnector
}

// ConnectorIDs returns the "connector" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ConnectorID instead. It exists only for internal usage by the builders.
func (m *ApplicationResourceMutation) ConnectorIDs() (ids []types.ID) {
	if id := m.connector; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetConnector resets all changes to the "connector" edge.
func (m *ApplicationResourceMutation) ResetConnector() {
	m.connector = nil
	m.clearedconnector = false
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
	fields := make([]string, 0, 10)
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
	if m.application != nil {
		fields = append(fields, applicationresource.FieldApplicationID)
	}
	if m.connector != nil {
		fields = append(fields, applicationresource.FieldConnectorID)
	}
	if m.module != nil {
		fields = append(fields, applicationresource.FieldModule)
	}
	if m.mode != nil {
		fields = append(fields, applicationresource.FieldMode)
	}
	if m._type != nil {
		fields = append(fields, applicationresource.FieldType)
	}
	if m.name != nil {
		fields = append(fields, applicationresource.FieldName)
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
	case applicationresource.FieldConnectorID:
		return m.ConnectorID()
	case applicationresource.FieldModule:
		return m.Module()
	case applicationresource.FieldMode:
		return m.Mode()
	case applicationresource.FieldType:
		return m.GetType()
	case applicationresource.FieldName:
		return m.Name()
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
	case applicationresource.FieldConnectorID:
		return m.OldConnectorID(ctx)
	case applicationresource.FieldModule:
		return m.OldModule(ctx)
	case applicationresource.FieldMode:
		return m.OldMode(ctx)
	case applicationresource.FieldType:
		return m.OldType(ctx)
	case applicationresource.FieldName:
		return m.OldName(ctx)
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
		v, ok := value.(types.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetApplicationID(v)
		return nil
	case applicationresource.FieldConnectorID:
		v, ok := value.(types.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetConnectorID(v)
		return nil
	case applicationresource.FieldModule:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetModule(v)
		return nil
	case applicationresource.FieldMode:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMode(v)
		return nil
	case applicationresource.FieldType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	case applicationresource.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
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
	case applicationresource.FieldConnectorID:
		m.ResetConnectorID()
		return nil
	case applicationresource.FieldModule:
		m.ResetModule()
		return nil
	case applicationresource.FieldMode:
		m.ResetMode()
		return nil
	case applicationresource.FieldType:
		m.ResetType()
		return nil
	case applicationresource.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationResourceMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.application != nil {
		edges = append(edges, applicationresource.EdgeApplication)
	}
	if m.connector != nil {
		edges = append(edges, applicationresource.EdgeConnector)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ApplicationResourceMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case applicationresource.EdgeApplication:
		if id := m.application; id != nil {
			return []ent.Value{*id}
		}
	case applicationresource.EdgeConnector:
		if id := m.connector; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationResourceMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationResourceMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationResourceMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedapplication {
		edges = append(edges, applicationresource.EdgeApplication)
	}
	if m.clearedconnector {
		edges = append(edges, applicationresource.EdgeConnector)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationResourceMutation) EdgeCleared(name string) bool {
	switch name {
	case applicationresource.EdgeApplication:
		return m.clearedapplication
	case applicationresource.EdgeConnector:
		return m.clearedconnector
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationResourceMutation) ClearEdge(name string) error {
	switch name {
	case applicationresource.EdgeApplication:
		m.ClearApplication()
		return nil
	case applicationresource.EdgeConnector:
		m.ClearConnector()
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ApplicationResourceMutation) ResetEdge(name string) error {
	switch name {
	case applicationresource.EdgeApplication:
		m.ResetApplication()
		return nil
	case applicationresource.EdgeConnector:
		m.ResetConnector()
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource edge %s", name)
}

// ApplicationRevisionMutation represents an operation that mutates the ApplicationRevision nodes in the graph.
type ApplicationRevisionMutation struct {
	config
	op                 Op
	typ                string
	id                 *types.ID
	status             *string
	statusMessage      *string
	createTime         *time.Time
	modules            *[]types.ApplicationModule
	appendmodules      []types.ApplicationModule
	inputVariables     *map[string]interface{}
	inputPlan          *string
	output             *string
	deployerType       *string
	duration           *int
	addduration        *int
	clearedFields      map[string]struct{}
	application        *types.ID
	clearedapplication bool
	environment        *types.ID
	clearedenvironment bool
	done               bool
	oldValue           func(context.Context) (*ApplicationRevision, error)
	predicates         []predicate.ApplicationRevision
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
func withApplicationRevisionID(id types.ID) applicationrevisionOption {
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
func (m *ApplicationRevisionMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ApplicationRevisionMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ApplicationRevisionMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
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

// SetApplicationID sets the "applicationID" field.
func (m *ApplicationRevisionMutation) SetApplicationID(t types.ID) {
	m.application = &t
}

// ApplicationID returns the value of the "applicationID" field in the mutation.
func (m *ApplicationRevisionMutation) ApplicationID() (r types.ID, exists bool) {
	v := m.application
	if v == nil {
		return
	}
	return *v, true
}

// OldApplicationID returns the old "applicationID" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldApplicationID(ctx context.Context) (v types.ID, err error) {
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
	m.application = nil
}

// SetEnvironmentID sets the "environmentID" field.
func (m *ApplicationRevisionMutation) SetEnvironmentID(t types.ID) {
	m.environment = &t
}

// EnvironmentID returns the value of the "environmentID" field in the mutation.
func (m *ApplicationRevisionMutation) EnvironmentID() (r types.ID, exists bool) {
	v := m.environment
	if v == nil {
		return
	}
	return *v, true
}

// OldEnvironmentID returns the old "environmentID" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldEnvironmentID(ctx context.Context) (v types.ID, err error) {
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
	m.environment = nil
}

// SetModules sets the "modules" field.
func (m *ApplicationRevisionMutation) SetModules(tm []types.ApplicationModule) {
	m.modules = &tm
	m.appendmodules = nil
}

// Modules returns the value of the "modules" field in the mutation.
func (m *ApplicationRevisionMutation) Modules() (r []types.ApplicationModule, exists bool) {
	v := m.modules
	if v == nil {
		return
	}
	return *v, true
}

// OldModules returns the old "modules" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldModules(ctx context.Context) (v []types.ApplicationModule, err error) {
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

// AppendModules adds tm to the "modules" field.
func (m *ApplicationRevisionMutation) AppendModules(tm []types.ApplicationModule) {
	m.appendmodules = append(m.appendmodules, tm...)
}

// AppendedModules returns the list of values that were appended to the "modules" field in this mutation.
func (m *ApplicationRevisionMutation) AppendedModules() ([]types.ApplicationModule, bool) {
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

// SetDeployerType sets the "deployerType" field.
func (m *ApplicationRevisionMutation) SetDeployerType(s string) {
	m.deployerType = &s
}

// DeployerType returns the value of the "deployerType" field in the mutation.
func (m *ApplicationRevisionMutation) DeployerType() (r string, exists bool) {
	v := m.deployerType
	if v == nil {
		return
	}
	return *v, true
}

// OldDeployerType returns the old "deployerType" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldDeployerType(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDeployerType is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDeployerType requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDeployerType: %w", err)
	}
	return oldValue.DeployerType, nil
}

// ResetDeployerType resets all changes to the "deployerType" field.
func (m *ApplicationRevisionMutation) ResetDeployerType() {
	m.deployerType = nil
}

// SetDuration sets the "duration" field.
func (m *ApplicationRevisionMutation) SetDuration(i int) {
	m.duration = &i
	m.addduration = nil
}

// Duration returns the value of the "duration" field in the mutation.
func (m *ApplicationRevisionMutation) Duration() (r int, exists bool) {
	v := m.duration
	if v == nil {
		return
	}
	return *v, true
}

// OldDuration returns the old "duration" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldDuration(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDuration is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDuration requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDuration: %w", err)
	}
	return oldValue.Duration, nil
}

// AddDuration adds i to the "duration" field.
func (m *ApplicationRevisionMutation) AddDuration(i int) {
	if m.addduration != nil {
		*m.addduration += i
	} else {
		m.addduration = &i
	}
}

// AddedDuration returns the value that was added to the "duration" field in this mutation.
func (m *ApplicationRevisionMutation) AddedDuration() (r int, exists bool) {
	v := m.addduration
	if v == nil {
		return
	}
	return *v, true
}

// ResetDuration resets all changes to the "duration" field.
func (m *ApplicationRevisionMutation) ResetDuration() {
	m.duration = nil
	m.addduration = nil
}

// ClearApplication clears the "application" edge to the Application entity.
func (m *ApplicationRevisionMutation) ClearApplication() {
	m.clearedapplication = true
}

// ApplicationCleared reports if the "application" edge to the Application entity was cleared.
func (m *ApplicationRevisionMutation) ApplicationCleared() bool {
	return m.clearedapplication
}

// ApplicationIDs returns the "application" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ApplicationID instead. It exists only for internal usage by the builders.
func (m *ApplicationRevisionMutation) ApplicationIDs() (ids []types.ID) {
	if id := m.application; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetApplication resets all changes to the "application" edge.
func (m *ApplicationRevisionMutation) ResetApplication() {
	m.application = nil
	m.clearedapplication = false
}

// ClearEnvironment clears the "environment" edge to the Environment entity.
func (m *ApplicationRevisionMutation) ClearEnvironment() {
	m.clearedenvironment = true
}

// EnvironmentCleared reports if the "environment" edge to the Environment entity was cleared.
func (m *ApplicationRevisionMutation) EnvironmentCleared() bool {
	return m.clearedenvironment
}

// EnvironmentIDs returns the "environment" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// EnvironmentID instead. It exists only for internal usage by the builders.
func (m *ApplicationRevisionMutation) EnvironmentIDs() (ids []types.ID) {
	if id := m.environment; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEnvironment resets all changes to the "environment" edge.
func (m *ApplicationRevisionMutation) ResetEnvironment() {
	m.environment = nil
	m.clearedenvironment = false
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
	fields := make([]string, 0, 11)
	if m.status != nil {
		fields = append(fields, applicationrevision.FieldStatus)
	}
	if m.statusMessage != nil {
		fields = append(fields, applicationrevision.FieldStatusMessage)
	}
	if m.createTime != nil {
		fields = append(fields, applicationrevision.FieldCreateTime)
	}
	if m.application != nil {
		fields = append(fields, applicationrevision.FieldApplicationID)
	}
	if m.environment != nil {
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
	if m.deployerType != nil {
		fields = append(fields, applicationrevision.FieldDeployerType)
	}
	if m.duration != nil {
		fields = append(fields, applicationrevision.FieldDuration)
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
	case applicationrevision.FieldDeployerType:
		return m.DeployerType()
	case applicationrevision.FieldDuration:
		return m.Duration()
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
	case applicationrevision.FieldDeployerType:
		return m.OldDeployerType(ctx)
	case applicationrevision.FieldDuration:
		return m.OldDuration(ctx)
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
	case applicationrevision.FieldApplicationID:
		v, ok := value.(types.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetApplicationID(v)
		return nil
	case applicationrevision.FieldEnvironmentID:
		v, ok := value.(types.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnvironmentID(v)
		return nil
	case applicationrevision.FieldModules:
		v, ok := value.([]types.ApplicationModule)
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
	case applicationrevision.FieldDeployerType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeployerType(v)
		return nil
	case applicationrevision.FieldDuration:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDuration(v)
		return nil
	}
	return fmt.Errorf("unknown ApplicationRevision field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ApplicationRevisionMutation) AddedFields() []string {
	var fields []string
	if m.addduration != nil {
		fields = append(fields, applicationrevision.FieldDuration)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ApplicationRevisionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case applicationrevision.FieldDuration:
		return m.AddedDuration()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationRevisionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case applicationrevision.FieldDuration:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddDuration(v)
		return nil
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
	case applicationrevision.FieldDeployerType:
		m.ResetDeployerType()
		return nil
	case applicationrevision.FieldDuration:
		m.ResetDuration()
		return nil
	}
	return fmt.Errorf("unknown ApplicationRevision field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationRevisionMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.application != nil {
		edges = append(edges, applicationrevision.EdgeApplication)
	}
	if m.environment != nil {
		edges = append(edges, applicationrevision.EdgeEnvironment)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ApplicationRevisionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case applicationrevision.EdgeApplication:
		if id := m.application; id != nil {
			return []ent.Value{*id}
		}
	case applicationrevision.EdgeEnvironment:
		if id := m.environment; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationRevisionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationRevisionMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationRevisionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedapplication {
		edges = append(edges, applicationrevision.EdgeApplication)
	}
	if m.clearedenvironment {
		edges = append(edges, applicationrevision.EdgeEnvironment)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationRevisionMutation) EdgeCleared(name string) bool {
	switch name {
	case applicationrevision.EdgeApplication:
		return m.clearedapplication
	case applicationrevision.EdgeEnvironment:
		return m.clearedenvironment
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationRevisionMutation) ClearEdge(name string) error {
	switch name {
	case applicationrevision.EdgeApplication:
		m.ClearApplication()
		return nil
	case applicationrevision.EdgeEnvironment:
		m.ClearEnvironment()
		return nil
	}
	return fmt.Errorf("unknown ApplicationRevision unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ApplicationRevisionMutation) ResetEdge(name string) error {
	switch name {
	case applicationrevision.EdgeApplication:
		m.ResetApplication()
		return nil
	case applicationrevision.EdgeEnvironment:
		m.ResetEnvironment()
		return nil
	}
	return fmt.Errorf("unknown ApplicationRevision edge %s", name)
}

// ConnectorMutation represents an operation that mutates the Connector nodes in the graph.
type ConnectorMutation struct {
	config
	op                                       Op
	typ                                      string
	id                                       *types.ID
	name                                     *string
	description                              *string
	labels                                   *map[string]string
	status                                   *string
	statusMessage                            *string
	createTime                               *time.Time
	updateTime                               *time.Time
	_type                                    *string
	configVersion                            *string
	configData                               *map[string]interface{}
	enableFinOps                             *bool
	finOpsStatus                             *string
	finOpsStatusMessage                      *string
	clearedFields                            map[string]struct{}
	environments                             map[types.ID]struct{}
	removedenvironments                      map[types.ID]struct{}
	clearedenvironments                      bool
	resources                                map[types.ID]struct{}
	removedresources                         map[types.ID]struct{}
	clearedresources                         bool
	environmentConnectorRelationships        map[int]struct{}
	removedenvironmentConnectorRelationships map[int]struct{}
	clearedenvironmentConnectorRelationships bool
	done                                     bool
	oldValue                                 func(context.Context) (*Connector, error)
	predicates                               []predicate.Connector
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
func withConnectorID(id types.ID) connectorOption {
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
func (m *ConnectorMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ConnectorMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ConnectorMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Connector.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *ConnectorMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *ConnectorMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *ConnectorMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the "description" field.
func (m *ConnectorMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the value of the "description" field in the mutation.
func (m *ConnectorMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// OldDescription returns the old "description" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldDescription(ctx context.Context) (v string, err error) {
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

// ClearDescription clears the value of the "description" field.
func (m *ConnectorMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[connector.FieldDescription] = struct{}{}
}

// DescriptionCleared returns if the "description" field was cleared in this mutation.
func (m *ConnectorMutation) DescriptionCleared() bool {
	_, ok := m.clearedFields[connector.FieldDescription]
	return ok
}

// ResetDescription resets all changes to the "description" field.
func (m *ConnectorMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, connector.FieldDescription)
}

// SetLabels sets the "labels" field.
func (m *ConnectorMutation) SetLabels(value map[string]string) {
	m.labels = &value
}

// Labels returns the value of the "labels" field in the mutation.
func (m *ConnectorMutation) Labels() (r map[string]string, exists bool) {
	v := m.labels
	if v == nil {
		return
	}
	return *v, true
}

// OldLabels returns the old "labels" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldLabels(ctx context.Context) (v map[string]string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLabels is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLabels requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLabels: %w", err)
	}
	return oldValue.Labels, nil
}

// ResetLabels resets all changes to the "labels" field.
func (m *ConnectorMutation) ResetLabels() {
	m.labels = nil
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

// SetType sets the "type" field.
func (m *ConnectorMutation) SetType(s string) {
	m._type = &s
}

// GetType returns the value of the "type" field in the mutation.
func (m *ConnectorMutation) GetType() (r string, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// OldType returns the old "type" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldType(ctx context.Context) (v string, err error) {
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
func (m *ConnectorMutation) ResetType() {
	m._type = nil
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

// ResetConfigData resets all changes to the "configData" field.
func (m *ConnectorMutation) ResetConfigData() {
	m.configData = nil
}

// SetEnableFinOps sets the "enableFinOps" field.
func (m *ConnectorMutation) SetEnableFinOps(b bool) {
	m.enableFinOps = &b
}

// EnableFinOps returns the value of the "enableFinOps" field in the mutation.
func (m *ConnectorMutation) EnableFinOps() (r bool, exists bool) {
	v := m.enableFinOps
	if v == nil {
		return
	}
	return *v, true
}

// OldEnableFinOps returns the old "enableFinOps" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldEnableFinOps(ctx context.Context) (v bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEnableFinOps is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEnableFinOps requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEnableFinOps: %w", err)
	}
	return oldValue.EnableFinOps, nil
}

// ResetEnableFinOps resets all changes to the "enableFinOps" field.
func (m *ConnectorMutation) ResetEnableFinOps() {
	m.enableFinOps = nil
}

// SetFinOpsStatus sets the "finOpsStatus" field.
func (m *ConnectorMutation) SetFinOpsStatus(s string) {
	m.finOpsStatus = &s
}

// FinOpsStatus returns the value of the "finOpsStatus" field in the mutation.
func (m *ConnectorMutation) FinOpsStatus() (r string, exists bool) {
	v := m.finOpsStatus
	if v == nil {
		return
	}
	return *v, true
}

// OldFinOpsStatus returns the old "finOpsStatus" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldFinOpsStatus(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFinOpsStatus is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFinOpsStatus requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFinOpsStatus: %w", err)
	}
	return oldValue.FinOpsStatus, nil
}

// ClearFinOpsStatus clears the value of the "finOpsStatus" field.
func (m *ConnectorMutation) ClearFinOpsStatus() {
	m.finOpsStatus = nil
	m.clearedFields[connector.FieldFinOpsStatus] = struct{}{}
}

// FinOpsStatusCleared returns if the "finOpsStatus" field was cleared in this mutation.
func (m *ConnectorMutation) FinOpsStatusCleared() bool {
	_, ok := m.clearedFields[connector.FieldFinOpsStatus]
	return ok
}

// ResetFinOpsStatus resets all changes to the "finOpsStatus" field.
func (m *ConnectorMutation) ResetFinOpsStatus() {
	m.finOpsStatus = nil
	delete(m.clearedFields, connector.FieldFinOpsStatus)
}

// SetFinOpsStatusMessage sets the "finOpsStatusMessage" field.
func (m *ConnectorMutation) SetFinOpsStatusMessage(s string) {
	m.finOpsStatusMessage = &s
}

// FinOpsStatusMessage returns the value of the "finOpsStatusMessage" field in the mutation.
func (m *ConnectorMutation) FinOpsStatusMessage() (r string, exists bool) {
	v := m.finOpsStatusMessage
	if v == nil {
		return
	}
	return *v, true
}

// OldFinOpsStatusMessage returns the old "finOpsStatusMessage" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldFinOpsStatusMessage(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFinOpsStatusMessage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFinOpsStatusMessage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFinOpsStatusMessage: %w", err)
	}
	return oldValue.FinOpsStatusMessage, nil
}

// ClearFinOpsStatusMessage clears the value of the "finOpsStatusMessage" field.
func (m *ConnectorMutation) ClearFinOpsStatusMessage() {
	m.finOpsStatusMessage = nil
	m.clearedFields[connector.FieldFinOpsStatusMessage] = struct{}{}
}

// FinOpsStatusMessageCleared returns if the "finOpsStatusMessage" field was cleared in this mutation.
func (m *ConnectorMutation) FinOpsStatusMessageCleared() bool {
	_, ok := m.clearedFields[connector.FieldFinOpsStatusMessage]
	return ok
}

// ResetFinOpsStatusMessage resets all changes to the "finOpsStatusMessage" field.
func (m *ConnectorMutation) ResetFinOpsStatusMessage() {
	m.finOpsStatusMessage = nil
	delete(m.clearedFields, connector.FieldFinOpsStatusMessage)
}

// AddEnvironmentIDs adds the "environments" edge to the Environment entity by ids.
func (m *ConnectorMutation) AddEnvironmentIDs(ids ...types.ID) {
	if m.environments == nil {
		m.environments = make(map[types.ID]struct{})
	}
	for i := range ids {
		m.environments[ids[i]] = struct{}{}
	}
}

// ClearEnvironments clears the "environments" edge to the Environment entity.
func (m *ConnectorMutation) ClearEnvironments() {
	m.clearedenvironments = true
}

// EnvironmentsCleared reports if the "environments" edge to the Environment entity was cleared.
func (m *ConnectorMutation) EnvironmentsCleared() bool {
	return m.clearedenvironments
}

// RemoveEnvironmentIDs removes the "environments" edge to the Environment entity by IDs.
func (m *ConnectorMutation) RemoveEnvironmentIDs(ids ...types.ID) {
	if m.removedenvironments == nil {
		m.removedenvironments = make(map[types.ID]struct{})
	}
	for i := range ids {
		delete(m.environments, ids[i])
		m.removedenvironments[ids[i]] = struct{}{}
	}
}

// RemovedEnvironments returns the removed IDs of the "environments" edge to the Environment entity.
func (m *ConnectorMutation) RemovedEnvironmentsIDs() (ids []types.ID) {
	for id := range m.removedenvironments {
		ids = append(ids, id)
	}
	return
}

// EnvironmentsIDs returns the "environments" edge IDs in the mutation.
func (m *ConnectorMutation) EnvironmentsIDs() (ids []types.ID) {
	for id := range m.environments {
		ids = append(ids, id)
	}
	return
}

// ResetEnvironments resets all changes to the "environments" edge.
func (m *ConnectorMutation) ResetEnvironments() {
	m.environments = nil
	m.clearedenvironments = false
	m.removedenvironments = nil
}

// AddResourceIDs adds the "resources" edge to the ApplicationResource entity by ids.
func (m *ConnectorMutation) AddResourceIDs(ids ...types.ID) {
	if m.resources == nil {
		m.resources = make(map[types.ID]struct{})
	}
	for i := range ids {
		m.resources[ids[i]] = struct{}{}
	}
}

// ClearResources clears the "resources" edge to the ApplicationResource entity.
func (m *ConnectorMutation) ClearResources() {
	m.clearedresources = true
}

// ResourcesCleared reports if the "resources" edge to the ApplicationResource entity was cleared.
func (m *ConnectorMutation) ResourcesCleared() bool {
	return m.clearedresources
}

// RemoveResourceIDs removes the "resources" edge to the ApplicationResource entity by IDs.
func (m *ConnectorMutation) RemoveResourceIDs(ids ...types.ID) {
	if m.removedresources == nil {
		m.removedresources = make(map[types.ID]struct{})
	}
	for i := range ids {
		delete(m.resources, ids[i])
		m.removedresources[ids[i]] = struct{}{}
	}
}

// RemovedResources returns the removed IDs of the "resources" edge to the ApplicationResource entity.
func (m *ConnectorMutation) RemovedResourcesIDs() (ids []types.ID) {
	for id := range m.removedresources {
		ids = append(ids, id)
	}
	return
}

// ResourcesIDs returns the "resources" edge IDs in the mutation.
func (m *ConnectorMutation) ResourcesIDs() (ids []types.ID) {
	for id := range m.resources {
		ids = append(ids, id)
	}
	return
}

// ResetResources resets all changes to the "resources" edge.
func (m *ConnectorMutation) ResetResources() {
	m.resources = nil
	m.clearedresources = false
	m.removedresources = nil
}

// AddEnvironmentConnectorRelationshipIDs adds the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity by ids.
func (m *ConnectorMutation) AddEnvironmentConnectorRelationshipIDs(ids ...int) {
	if m.environmentConnectorRelationships == nil {
		m.environmentConnectorRelationships = make(map[int]struct{})
	}
	for i := range ids {
		m.environmentConnectorRelationships[ids[i]] = struct{}{}
	}
}

// ClearEnvironmentConnectorRelationships clears the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity.
func (m *ConnectorMutation) ClearEnvironmentConnectorRelationships() {
	m.clearedenvironmentConnectorRelationships = true
}

// EnvironmentConnectorRelationshipsCleared reports if the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity was cleared.
func (m *ConnectorMutation) EnvironmentConnectorRelationshipsCleared() bool {
	return m.clearedenvironmentConnectorRelationships
}

// RemoveEnvironmentConnectorRelationshipIDs removes the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity by IDs.
func (m *ConnectorMutation) RemoveEnvironmentConnectorRelationshipIDs(ids ...int) {
	if m.removedenvironmentConnectorRelationships == nil {
		m.removedenvironmentConnectorRelationships = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.environmentConnectorRelationships, ids[i])
		m.removedenvironmentConnectorRelationships[ids[i]] = struct{}{}
	}
}

// RemovedEnvironmentConnectorRelationships returns the removed IDs of the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity.
func (m *ConnectorMutation) RemovedEnvironmentConnectorRelationshipsIDs() (ids []int) {
	for id := range m.removedenvironmentConnectorRelationships {
		ids = append(ids, id)
	}
	return
}

// EnvironmentConnectorRelationshipsIDs returns the "environmentConnectorRelationships" edge IDs in the mutation.
func (m *ConnectorMutation) EnvironmentConnectorRelationshipsIDs() (ids []int) {
	for id := range m.environmentConnectorRelationships {
		ids = append(ids, id)
	}
	return
}

// ResetEnvironmentConnectorRelationships resets all changes to the "environmentConnectorRelationships" edge.
func (m *ConnectorMutation) ResetEnvironmentConnectorRelationships() {
	m.environmentConnectorRelationships = nil
	m.clearedenvironmentConnectorRelationships = false
	m.removedenvironmentConnectorRelationships = nil
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
	fields := make([]string, 0, 13)
	if m.name != nil {
		fields = append(fields, connector.FieldName)
	}
	if m.description != nil {
		fields = append(fields, connector.FieldDescription)
	}
	if m.labels != nil {
		fields = append(fields, connector.FieldLabels)
	}
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
	if m._type != nil {
		fields = append(fields, connector.FieldType)
	}
	if m.configVersion != nil {
		fields = append(fields, connector.FieldConfigVersion)
	}
	if m.configData != nil {
		fields = append(fields, connector.FieldConfigData)
	}
	if m.enableFinOps != nil {
		fields = append(fields, connector.FieldEnableFinOps)
	}
	if m.finOpsStatus != nil {
		fields = append(fields, connector.FieldFinOpsStatus)
	}
	if m.finOpsStatusMessage != nil {
		fields = append(fields, connector.FieldFinOpsStatusMessage)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ConnectorMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case connector.FieldName:
		return m.Name()
	case connector.FieldDescription:
		return m.Description()
	case connector.FieldLabels:
		return m.Labels()
	case connector.FieldStatus:
		return m.Status()
	case connector.FieldStatusMessage:
		return m.StatusMessage()
	case connector.FieldCreateTime:
		return m.CreateTime()
	case connector.FieldUpdateTime:
		return m.UpdateTime()
	case connector.FieldType:
		return m.GetType()
	case connector.FieldConfigVersion:
		return m.ConfigVersion()
	case connector.FieldConfigData:
		return m.ConfigData()
	case connector.FieldEnableFinOps:
		return m.EnableFinOps()
	case connector.FieldFinOpsStatus:
		return m.FinOpsStatus()
	case connector.FieldFinOpsStatusMessage:
		return m.FinOpsStatusMessage()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ConnectorMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case connector.FieldName:
		return m.OldName(ctx)
	case connector.FieldDescription:
		return m.OldDescription(ctx)
	case connector.FieldLabels:
		return m.OldLabels(ctx)
	case connector.FieldStatus:
		return m.OldStatus(ctx)
	case connector.FieldStatusMessage:
		return m.OldStatusMessage(ctx)
	case connector.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case connector.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case connector.FieldType:
		return m.OldType(ctx)
	case connector.FieldConfigVersion:
		return m.OldConfigVersion(ctx)
	case connector.FieldConfigData:
		return m.OldConfigData(ctx)
	case connector.FieldEnableFinOps:
		return m.OldEnableFinOps(ctx)
	case connector.FieldFinOpsStatus:
		return m.OldFinOpsStatus(ctx)
	case connector.FieldFinOpsStatusMessage:
		return m.OldFinOpsStatusMessage(ctx)
	}
	return nil, fmt.Errorf("unknown Connector field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ConnectorMutation) SetField(name string, value ent.Value) error {
	switch name {
	case connector.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case connector.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case connector.FieldLabels:
		v, ok := value.(map[string]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLabels(v)
		return nil
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
	case connector.FieldType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
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
	case connector.FieldEnableFinOps:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnableFinOps(v)
		return nil
	case connector.FieldFinOpsStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFinOpsStatus(v)
		return nil
	case connector.FieldFinOpsStatusMessage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFinOpsStatusMessage(v)
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
	if m.FieldCleared(connector.FieldDescription) {
		fields = append(fields, connector.FieldDescription)
	}
	if m.FieldCleared(connector.FieldStatus) {
		fields = append(fields, connector.FieldStatus)
	}
	if m.FieldCleared(connector.FieldStatusMessage) {
		fields = append(fields, connector.FieldStatusMessage)
	}
	if m.FieldCleared(connector.FieldFinOpsStatus) {
		fields = append(fields, connector.FieldFinOpsStatus)
	}
	if m.FieldCleared(connector.FieldFinOpsStatusMessage) {
		fields = append(fields, connector.FieldFinOpsStatusMessage)
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
	case connector.FieldDescription:
		m.ClearDescription()
		return nil
	case connector.FieldStatus:
		m.ClearStatus()
		return nil
	case connector.FieldStatusMessage:
		m.ClearStatusMessage()
		return nil
	case connector.FieldFinOpsStatus:
		m.ClearFinOpsStatus()
		return nil
	case connector.FieldFinOpsStatusMessage:
		m.ClearFinOpsStatusMessage()
		return nil
	}
	return fmt.Errorf("unknown Connector nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ConnectorMutation) ResetField(name string) error {
	switch name {
	case connector.FieldName:
		m.ResetName()
		return nil
	case connector.FieldDescription:
		m.ResetDescription()
		return nil
	case connector.FieldLabels:
		m.ResetLabels()
		return nil
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
	case connector.FieldType:
		m.ResetType()
		return nil
	case connector.FieldConfigVersion:
		m.ResetConfigVersion()
		return nil
	case connector.FieldConfigData:
		m.ResetConfigData()
		return nil
	case connector.FieldEnableFinOps:
		m.ResetEnableFinOps()
		return nil
	case connector.FieldFinOpsStatus:
		m.ResetFinOpsStatus()
		return nil
	case connector.FieldFinOpsStatusMessage:
		m.ResetFinOpsStatusMessage()
		return nil
	}
	return fmt.Errorf("unknown Connector field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ConnectorMutation) AddedEdges() []string {
	edges := make([]string, 0, 3)
	if m.environments != nil {
		edges = append(edges, connector.EdgeEnvironments)
	}
	if m.resources != nil {
		edges = append(edges, connector.EdgeResources)
	}
	if m.environmentConnectorRelationships != nil {
		edges = append(edges, connector.EdgeEnvironmentConnectorRelationships)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ConnectorMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case connector.EdgeEnvironments:
		ids := make([]ent.Value, 0, len(m.environments))
		for id := range m.environments {
			ids = append(ids, id)
		}
		return ids
	case connector.EdgeResources:
		ids := make([]ent.Value, 0, len(m.resources))
		for id := range m.resources {
			ids = append(ids, id)
		}
		return ids
	case connector.EdgeEnvironmentConnectorRelationships:
		ids := make([]ent.Value, 0, len(m.environmentConnectorRelationships))
		for id := range m.environmentConnectorRelationships {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ConnectorMutation) RemovedEdges() []string {
	edges := make([]string, 0, 3)
	if m.removedenvironments != nil {
		edges = append(edges, connector.EdgeEnvironments)
	}
	if m.removedresources != nil {
		edges = append(edges, connector.EdgeResources)
	}
	if m.removedenvironmentConnectorRelationships != nil {
		edges = append(edges, connector.EdgeEnvironmentConnectorRelationships)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ConnectorMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case connector.EdgeEnvironments:
		ids := make([]ent.Value, 0, len(m.removedenvironments))
		for id := range m.removedenvironments {
			ids = append(ids, id)
		}
		return ids
	case connector.EdgeResources:
		ids := make([]ent.Value, 0, len(m.removedresources))
		for id := range m.removedresources {
			ids = append(ids, id)
		}
		return ids
	case connector.EdgeEnvironmentConnectorRelationships:
		ids := make([]ent.Value, 0, len(m.removedenvironmentConnectorRelationships))
		for id := range m.removedenvironmentConnectorRelationships {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ConnectorMutation) ClearedEdges() []string {
	edges := make([]string, 0, 3)
	if m.clearedenvironments {
		edges = append(edges, connector.EdgeEnvironments)
	}
	if m.clearedresources {
		edges = append(edges, connector.EdgeResources)
	}
	if m.clearedenvironmentConnectorRelationships {
		edges = append(edges, connector.EdgeEnvironmentConnectorRelationships)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ConnectorMutation) EdgeCleared(name string) bool {
	switch name {
	case connector.EdgeEnvironments:
		return m.clearedenvironments
	case connector.EdgeResources:
		return m.clearedresources
	case connector.EdgeEnvironmentConnectorRelationships:
		return m.clearedenvironmentConnectorRelationships
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ConnectorMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Connector unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ConnectorMutation) ResetEdge(name string) error {
	switch name {
	case connector.EdgeEnvironments:
		m.ResetEnvironments()
		return nil
	case connector.EdgeResources:
		m.ResetResources()
		return nil
	case connector.EdgeEnvironmentConnectorRelationships:
		m.ResetEnvironmentConnectorRelationships()
		return nil
	}
	return fmt.Errorf("unknown Connector edge %s", name)
}

// EnvironmentMutation represents an operation that mutates the Environment nodes in the graph.
type EnvironmentMutation struct {
	config
	op                                       Op
	typ                                      string
	id                                       *types.ID
	name                                     *string
	description                              *string
	labels                                   *map[string]string
	createTime                               *time.Time
	updateTime                               *time.Time
	variables                                *map[string]interface{}
	clearedFields                            map[string]struct{}
	connectors                               map[types.ID]struct{}
	removedconnectors                        map[types.ID]struct{}
	clearedconnectors                        bool
	applications                             map[types.ID]struct{}
	removedapplications                      map[types.ID]struct{}
	clearedapplications                      bool
	revisions                                map[types.ID]struct{}
	removedrevisions                         map[types.ID]struct{}
	clearedrevisions                         bool
	environmentConnectorRelationships        map[int]struct{}
	removedenvironmentConnectorRelationships map[int]struct{}
	clearedenvironmentConnectorRelationships bool
	done                                     bool
	oldValue                                 func(context.Context) (*Environment, error)
	predicates                               []predicate.Environment
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
func withEnvironmentID(id types.ID) environmentOption {
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
func (m *EnvironmentMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *EnvironmentMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *EnvironmentMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Environment.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *EnvironmentMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *EnvironmentMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Environment entity.
// If the Environment object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *EnvironmentMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the "description" field.
func (m *EnvironmentMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the value of the "description" field in the mutation.
func (m *EnvironmentMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// OldDescription returns the old "description" field's value of the Environment entity.
// If the Environment object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentMutation) OldDescription(ctx context.Context) (v string, err error) {
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

// ClearDescription clears the value of the "description" field.
func (m *EnvironmentMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[environment.FieldDescription] = struct{}{}
}

// DescriptionCleared returns if the "description" field was cleared in this mutation.
func (m *EnvironmentMutation) DescriptionCleared() bool {
	_, ok := m.clearedFields[environment.FieldDescription]
	return ok
}

// ResetDescription resets all changes to the "description" field.
func (m *EnvironmentMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, environment.FieldDescription)
}

// SetLabels sets the "labels" field.
func (m *EnvironmentMutation) SetLabels(value map[string]string) {
	m.labels = &value
}

// Labels returns the value of the "labels" field in the mutation.
func (m *EnvironmentMutation) Labels() (r map[string]string, exists bool) {
	v := m.labels
	if v == nil {
		return
	}
	return *v, true
}

// OldLabels returns the old "labels" field's value of the Environment entity.
// If the Environment object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentMutation) OldLabels(ctx context.Context) (v map[string]string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLabels is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLabels requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLabels: %w", err)
	}
	return oldValue.Labels, nil
}

// ResetLabels resets all changes to the "labels" field.
func (m *EnvironmentMutation) ResetLabels() {
	m.labels = nil
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

// AddConnectorIDs adds the "connectors" edge to the Connector entity by ids.
func (m *EnvironmentMutation) AddConnectorIDs(ids ...types.ID) {
	if m.connectors == nil {
		m.connectors = make(map[types.ID]struct{})
	}
	for i := range ids {
		m.connectors[ids[i]] = struct{}{}
	}
}

// ClearConnectors clears the "connectors" edge to the Connector entity.
func (m *EnvironmentMutation) ClearConnectors() {
	m.clearedconnectors = true
}

// ConnectorsCleared reports if the "connectors" edge to the Connector entity was cleared.
func (m *EnvironmentMutation) ConnectorsCleared() bool {
	return m.clearedconnectors
}

// RemoveConnectorIDs removes the "connectors" edge to the Connector entity by IDs.
func (m *EnvironmentMutation) RemoveConnectorIDs(ids ...types.ID) {
	if m.removedconnectors == nil {
		m.removedconnectors = make(map[types.ID]struct{})
	}
	for i := range ids {
		delete(m.connectors, ids[i])
		m.removedconnectors[ids[i]] = struct{}{}
	}
}

// RemovedConnectors returns the removed IDs of the "connectors" edge to the Connector entity.
func (m *EnvironmentMutation) RemovedConnectorsIDs() (ids []types.ID) {
	for id := range m.removedconnectors {
		ids = append(ids, id)
	}
	return
}

// ConnectorsIDs returns the "connectors" edge IDs in the mutation.
func (m *EnvironmentMutation) ConnectorsIDs() (ids []types.ID) {
	for id := range m.connectors {
		ids = append(ids, id)
	}
	return
}

// ResetConnectors resets all changes to the "connectors" edge.
func (m *EnvironmentMutation) ResetConnectors() {
	m.connectors = nil
	m.clearedconnectors = false
	m.removedconnectors = nil
}

// AddApplicationIDs adds the "applications" edge to the Application entity by ids.
func (m *EnvironmentMutation) AddApplicationIDs(ids ...types.ID) {
	if m.applications == nil {
		m.applications = make(map[types.ID]struct{})
	}
	for i := range ids {
		m.applications[ids[i]] = struct{}{}
	}
}

// ClearApplications clears the "applications" edge to the Application entity.
func (m *EnvironmentMutation) ClearApplications() {
	m.clearedapplications = true
}

// ApplicationsCleared reports if the "applications" edge to the Application entity was cleared.
func (m *EnvironmentMutation) ApplicationsCleared() bool {
	return m.clearedapplications
}

// RemoveApplicationIDs removes the "applications" edge to the Application entity by IDs.
func (m *EnvironmentMutation) RemoveApplicationIDs(ids ...types.ID) {
	if m.removedapplications == nil {
		m.removedapplications = make(map[types.ID]struct{})
	}
	for i := range ids {
		delete(m.applications, ids[i])
		m.removedapplications[ids[i]] = struct{}{}
	}
}

// RemovedApplications returns the removed IDs of the "applications" edge to the Application entity.
func (m *EnvironmentMutation) RemovedApplicationsIDs() (ids []types.ID) {
	for id := range m.removedapplications {
		ids = append(ids, id)
	}
	return
}

// ApplicationsIDs returns the "applications" edge IDs in the mutation.
func (m *EnvironmentMutation) ApplicationsIDs() (ids []types.ID) {
	for id := range m.applications {
		ids = append(ids, id)
	}
	return
}

// ResetApplications resets all changes to the "applications" edge.
func (m *EnvironmentMutation) ResetApplications() {
	m.applications = nil
	m.clearedapplications = false
	m.removedapplications = nil
}

// AddRevisionIDs adds the "revisions" edge to the ApplicationRevision entity by ids.
func (m *EnvironmentMutation) AddRevisionIDs(ids ...types.ID) {
	if m.revisions == nil {
		m.revisions = make(map[types.ID]struct{})
	}
	for i := range ids {
		m.revisions[ids[i]] = struct{}{}
	}
}

// ClearRevisions clears the "revisions" edge to the ApplicationRevision entity.
func (m *EnvironmentMutation) ClearRevisions() {
	m.clearedrevisions = true
}

// RevisionsCleared reports if the "revisions" edge to the ApplicationRevision entity was cleared.
func (m *EnvironmentMutation) RevisionsCleared() bool {
	return m.clearedrevisions
}

// RemoveRevisionIDs removes the "revisions" edge to the ApplicationRevision entity by IDs.
func (m *EnvironmentMutation) RemoveRevisionIDs(ids ...types.ID) {
	if m.removedrevisions == nil {
		m.removedrevisions = make(map[types.ID]struct{})
	}
	for i := range ids {
		delete(m.revisions, ids[i])
		m.removedrevisions[ids[i]] = struct{}{}
	}
}

// RemovedRevisions returns the removed IDs of the "revisions" edge to the ApplicationRevision entity.
func (m *EnvironmentMutation) RemovedRevisionsIDs() (ids []types.ID) {
	for id := range m.removedrevisions {
		ids = append(ids, id)
	}
	return
}

// RevisionsIDs returns the "revisions" edge IDs in the mutation.
func (m *EnvironmentMutation) RevisionsIDs() (ids []types.ID) {
	for id := range m.revisions {
		ids = append(ids, id)
	}
	return
}

// ResetRevisions resets all changes to the "revisions" edge.
func (m *EnvironmentMutation) ResetRevisions() {
	m.revisions = nil
	m.clearedrevisions = false
	m.removedrevisions = nil
}

// AddEnvironmentConnectorRelationshipIDs adds the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity by ids.
func (m *EnvironmentMutation) AddEnvironmentConnectorRelationshipIDs(ids ...int) {
	if m.environmentConnectorRelationships == nil {
		m.environmentConnectorRelationships = make(map[int]struct{})
	}
	for i := range ids {
		m.environmentConnectorRelationships[ids[i]] = struct{}{}
	}
}

// ClearEnvironmentConnectorRelationships clears the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity.
func (m *EnvironmentMutation) ClearEnvironmentConnectorRelationships() {
	m.clearedenvironmentConnectorRelationships = true
}

// EnvironmentConnectorRelationshipsCleared reports if the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity was cleared.
func (m *EnvironmentMutation) EnvironmentConnectorRelationshipsCleared() bool {
	return m.clearedenvironmentConnectorRelationships
}

// RemoveEnvironmentConnectorRelationshipIDs removes the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity by IDs.
func (m *EnvironmentMutation) RemoveEnvironmentConnectorRelationshipIDs(ids ...int) {
	if m.removedenvironmentConnectorRelationships == nil {
		m.removedenvironmentConnectorRelationships = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.environmentConnectorRelationships, ids[i])
		m.removedenvironmentConnectorRelationships[ids[i]] = struct{}{}
	}
}

// RemovedEnvironmentConnectorRelationships returns the removed IDs of the "environmentConnectorRelationships" edge to the EnvironmentConnectorRelationship entity.
func (m *EnvironmentMutation) RemovedEnvironmentConnectorRelationshipsIDs() (ids []int) {
	for id := range m.removedenvironmentConnectorRelationships {
		ids = append(ids, id)
	}
	return
}

// EnvironmentConnectorRelationshipsIDs returns the "environmentConnectorRelationships" edge IDs in the mutation.
func (m *EnvironmentMutation) EnvironmentConnectorRelationshipsIDs() (ids []int) {
	for id := range m.environmentConnectorRelationships {
		ids = append(ids, id)
	}
	return
}

// ResetEnvironmentConnectorRelationships resets all changes to the "environmentConnectorRelationships" edge.
func (m *EnvironmentMutation) ResetEnvironmentConnectorRelationships() {
	m.environmentConnectorRelationships = nil
	m.clearedenvironmentConnectorRelationships = false
	m.removedenvironmentConnectorRelationships = nil
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
	fields := make([]string, 0, 6)
	if m.name != nil {
		fields = append(fields, environment.FieldName)
	}
	if m.description != nil {
		fields = append(fields, environment.FieldDescription)
	}
	if m.labels != nil {
		fields = append(fields, environment.FieldLabels)
	}
	if m.createTime != nil {
		fields = append(fields, environment.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, environment.FieldUpdateTime)
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
	case environment.FieldName:
		return m.Name()
	case environment.FieldDescription:
		return m.Description()
	case environment.FieldLabels:
		return m.Labels()
	case environment.FieldCreateTime:
		return m.CreateTime()
	case environment.FieldUpdateTime:
		return m.UpdateTime()
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
	case environment.FieldName:
		return m.OldName(ctx)
	case environment.FieldDescription:
		return m.OldDescription(ctx)
	case environment.FieldLabels:
		return m.OldLabels(ctx)
	case environment.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case environment.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
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
	case environment.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case environment.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case environment.FieldLabels:
		v, ok := value.(map[string]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLabels(v)
		return nil
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
	if m.FieldCleared(environment.FieldDescription) {
		fields = append(fields, environment.FieldDescription)
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
	case environment.FieldDescription:
		m.ClearDescription()
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
	case environment.FieldName:
		m.ResetName()
		return nil
	case environment.FieldDescription:
		m.ResetDescription()
		return nil
	case environment.FieldLabels:
		m.ResetLabels()
		return nil
	case environment.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case environment.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case environment.FieldVariables:
		m.ResetVariables()
		return nil
	}
	return fmt.Errorf("unknown Environment field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *EnvironmentMutation) AddedEdges() []string {
	edges := make([]string, 0, 4)
	if m.connectors != nil {
		edges = append(edges, environment.EdgeConnectors)
	}
	if m.applications != nil {
		edges = append(edges, environment.EdgeApplications)
	}
	if m.revisions != nil {
		edges = append(edges, environment.EdgeRevisions)
	}
	if m.environmentConnectorRelationships != nil {
		edges = append(edges, environment.EdgeEnvironmentConnectorRelationships)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *EnvironmentMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case environment.EdgeConnectors:
		ids := make([]ent.Value, 0, len(m.connectors))
		for id := range m.connectors {
			ids = append(ids, id)
		}
		return ids
	case environment.EdgeApplications:
		ids := make([]ent.Value, 0, len(m.applications))
		for id := range m.applications {
			ids = append(ids, id)
		}
		return ids
	case environment.EdgeRevisions:
		ids := make([]ent.Value, 0, len(m.revisions))
		for id := range m.revisions {
			ids = append(ids, id)
		}
		return ids
	case environment.EdgeEnvironmentConnectorRelationships:
		ids := make([]ent.Value, 0, len(m.environmentConnectorRelationships))
		for id := range m.environmentConnectorRelationships {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *EnvironmentMutation) RemovedEdges() []string {
	edges := make([]string, 0, 4)
	if m.removedconnectors != nil {
		edges = append(edges, environment.EdgeConnectors)
	}
	if m.removedapplications != nil {
		edges = append(edges, environment.EdgeApplications)
	}
	if m.removedrevisions != nil {
		edges = append(edges, environment.EdgeRevisions)
	}
	if m.removedenvironmentConnectorRelationships != nil {
		edges = append(edges, environment.EdgeEnvironmentConnectorRelationships)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *EnvironmentMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case environment.EdgeConnectors:
		ids := make([]ent.Value, 0, len(m.removedconnectors))
		for id := range m.removedconnectors {
			ids = append(ids, id)
		}
		return ids
	case environment.EdgeApplications:
		ids := make([]ent.Value, 0, len(m.removedapplications))
		for id := range m.removedapplications {
			ids = append(ids, id)
		}
		return ids
	case environment.EdgeRevisions:
		ids := make([]ent.Value, 0, len(m.removedrevisions))
		for id := range m.removedrevisions {
			ids = append(ids, id)
		}
		return ids
	case environment.EdgeEnvironmentConnectorRelationships:
		ids := make([]ent.Value, 0, len(m.removedenvironmentConnectorRelationships))
		for id := range m.removedenvironmentConnectorRelationships {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *EnvironmentMutation) ClearedEdges() []string {
	edges := make([]string, 0, 4)
	if m.clearedconnectors {
		edges = append(edges, environment.EdgeConnectors)
	}
	if m.clearedapplications {
		edges = append(edges, environment.EdgeApplications)
	}
	if m.clearedrevisions {
		edges = append(edges, environment.EdgeRevisions)
	}
	if m.clearedenvironmentConnectorRelationships {
		edges = append(edges, environment.EdgeEnvironmentConnectorRelationships)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *EnvironmentMutation) EdgeCleared(name string) bool {
	switch name {
	case environment.EdgeConnectors:
		return m.clearedconnectors
	case environment.EdgeApplications:
		return m.clearedapplications
	case environment.EdgeRevisions:
		return m.clearedrevisions
	case environment.EdgeEnvironmentConnectorRelationships:
		return m.clearedenvironmentConnectorRelationships
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *EnvironmentMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Environment unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *EnvironmentMutation) ResetEdge(name string) error {
	switch name {
	case environment.EdgeConnectors:
		m.ResetConnectors()
		return nil
	case environment.EdgeApplications:
		m.ResetApplications()
		return nil
	case environment.EdgeRevisions:
		m.ResetRevisions()
		return nil
	case environment.EdgeEnvironmentConnectorRelationships:
		m.ResetEnvironmentConnectorRelationships()
		return nil
	}
	return fmt.Errorf("unknown Environment edge %s", name)
}

// EnvironmentConnectorRelationshipMutation represents an operation that mutates the EnvironmentConnectorRelationship nodes in the graph.
type EnvironmentConnectorRelationshipMutation struct {
	config
	op                 Op
	typ                string
	id                 *int
	createTime         *time.Time
	clearedFields      map[string]struct{}
	environment        *types.ID
	clearedenvironment bool
	connector          *types.ID
	clearedconnector   bool
	done               bool
	oldValue           func(context.Context) (*EnvironmentConnectorRelationship, error)
	predicates         []predicate.EnvironmentConnectorRelationship
}

var _ ent.Mutation = (*EnvironmentConnectorRelationshipMutation)(nil)

// environmentconnectorrelationshipOption allows management of the mutation configuration using functional options.
type environmentconnectorrelationshipOption func(*EnvironmentConnectorRelationshipMutation)

// newEnvironmentConnectorRelationshipMutation creates new mutation for the EnvironmentConnectorRelationship entity.
func newEnvironmentConnectorRelationshipMutation(c config, op Op, opts ...environmentconnectorrelationshipOption) *EnvironmentConnectorRelationshipMutation {
	m := &EnvironmentConnectorRelationshipMutation{
		config:        c,
		op:            op,
		typ:           TypeEnvironmentConnectorRelationship,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withEnvironmentConnectorRelationshipID sets the ID field of the mutation.
func withEnvironmentConnectorRelationshipID(id int) environmentconnectorrelationshipOption {
	return func(m *EnvironmentConnectorRelationshipMutation) {
		var (
			err   error
			once  sync.Once
			value *EnvironmentConnectorRelationship
		)
		m.oldValue = func(ctx context.Context) (*EnvironmentConnectorRelationship, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().EnvironmentConnectorRelationship.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withEnvironmentConnectorRelationship sets the old EnvironmentConnectorRelationship of the mutation.
func withEnvironmentConnectorRelationship(node *EnvironmentConnectorRelationship) environmentconnectorrelationshipOption {
	return func(m *EnvironmentConnectorRelationshipMutation) {
		m.oldValue = func(context.Context) (*EnvironmentConnectorRelationship, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EnvironmentConnectorRelationshipMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EnvironmentConnectorRelationshipMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *EnvironmentConnectorRelationshipMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *EnvironmentConnectorRelationshipMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().EnvironmentConnectorRelationship.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *EnvironmentConnectorRelationshipMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *EnvironmentConnectorRelationshipMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the EnvironmentConnectorRelationship entity.
// If the EnvironmentConnectorRelationship object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentConnectorRelationshipMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *EnvironmentConnectorRelationshipMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetEnvironmentID sets the "environment_id" field.
func (m *EnvironmentConnectorRelationshipMutation) SetEnvironmentID(t types.ID) {
	m.environment = &t
}

// EnvironmentID returns the value of the "environment_id" field in the mutation.
func (m *EnvironmentConnectorRelationshipMutation) EnvironmentID() (r types.ID, exists bool) {
	v := m.environment
	if v == nil {
		return
	}
	return *v, true
}

// OldEnvironmentID returns the old "environment_id" field's value of the EnvironmentConnectorRelationship entity.
// If the EnvironmentConnectorRelationship object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentConnectorRelationshipMutation) OldEnvironmentID(ctx context.Context) (v types.ID, err error) {
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

// ResetEnvironmentID resets all changes to the "environment_id" field.
func (m *EnvironmentConnectorRelationshipMutation) ResetEnvironmentID() {
	m.environment = nil
}

// SetConnectorID sets the "connector_id" field.
func (m *EnvironmentConnectorRelationshipMutation) SetConnectorID(t types.ID) {
	m.connector = &t
}

// ConnectorID returns the value of the "connector_id" field in the mutation.
func (m *EnvironmentConnectorRelationshipMutation) ConnectorID() (r types.ID, exists bool) {
	v := m.connector
	if v == nil {
		return
	}
	return *v, true
}

// OldConnectorID returns the old "connector_id" field's value of the EnvironmentConnectorRelationship entity.
// If the EnvironmentConnectorRelationship object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *EnvironmentConnectorRelationshipMutation) OldConnectorID(ctx context.Context) (v types.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldConnectorID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldConnectorID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldConnectorID: %w", err)
	}
	return oldValue.ConnectorID, nil
}

// ResetConnectorID resets all changes to the "connector_id" field.
func (m *EnvironmentConnectorRelationshipMutation) ResetConnectorID() {
	m.connector = nil
}

// ClearEnvironment clears the "environment" edge to the Environment entity.
func (m *EnvironmentConnectorRelationshipMutation) ClearEnvironment() {
	m.clearedenvironment = true
}

// EnvironmentCleared reports if the "environment" edge to the Environment entity was cleared.
func (m *EnvironmentConnectorRelationshipMutation) EnvironmentCleared() bool {
	return m.clearedenvironment
}

// EnvironmentIDs returns the "environment" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// EnvironmentID instead. It exists only for internal usage by the builders.
func (m *EnvironmentConnectorRelationshipMutation) EnvironmentIDs() (ids []types.ID) {
	if id := m.environment; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEnvironment resets all changes to the "environment" edge.
func (m *EnvironmentConnectorRelationshipMutation) ResetEnvironment() {
	m.environment = nil
	m.clearedenvironment = false
}

// ClearConnector clears the "connector" edge to the Connector entity.
func (m *EnvironmentConnectorRelationshipMutation) ClearConnector() {
	m.clearedconnector = true
}

// ConnectorCleared reports if the "connector" edge to the Connector entity was cleared.
func (m *EnvironmentConnectorRelationshipMutation) ConnectorCleared() bool {
	return m.clearedconnector
}

// ConnectorIDs returns the "connector" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ConnectorID instead. It exists only for internal usage by the builders.
func (m *EnvironmentConnectorRelationshipMutation) ConnectorIDs() (ids []types.ID) {
	if id := m.connector; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetConnector resets all changes to the "connector" edge.
func (m *EnvironmentConnectorRelationshipMutation) ResetConnector() {
	m.connector = nil
	m.clearedconnector = false
}

// Where appends a list predicates to the EnvironmentConnectorRelationshipMutation builder.
func (m *EnvironmentConnectorRelationshipMutation) Where(ps ...predicate.EnvironmentConnectorRelationship) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the EnvironmentConnectorRelationshipMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *EnvironmentConnectorRelationshipMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.EnvironmentConnectorRelationship, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *EnvironmentConnectorRelationshipMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *EnvironmentConnectorRelationshipMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (EnvironmentConnectorRelationship).
func (m *EnvironmentConnectorRelationshipMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *EnvironmentConnectorRelationshipMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.createTime != nil {
		fields = append(fields, environmentconnectorrelationship.FieldCreateTime)
	}
	if m.environment != nil {
		fields = append(fields, environmentconnectorrelationship.FieldEnvironmentID)
	}
	if m.connector != nil {
		fields = append(fields, environmentconnectorrelationship.FieldConnectorID)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *EnvironmentConnectorRelationshipMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case environmentconnectorrelationship.FieldCreateTime:
		return m.CreateTime()
	case environmentconnectorrelationship.FieldEnvironmentID:
		return m.EnvironmentID()
	case environmentconnectorrelationship.FieldConnectorID:
		return m.ConnectorID()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *EnvironmentConnectorRelationshipMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case environmentconnectorrelationship.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case environmentconnectorrelationship.FieldEnvironmentID:
		return m.OldEnvironmentID(ctx)
	case environmentconnectorrelationship.FieldConnectorID:
		return m.OldConnectorID(ctx)
	}
	return nil, fmt.Errorf("unknown EnvironmentConnectorRelationship field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *EnvironmentConnectorRelationshipMutation) SetField(name string, value ent.Value) error {
	switch name {
	case environmentconnectorrelationship.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case environmentconnectorrelationship.FieldEnvironmentID:
		v, ok := value.(types.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnvironmentID(v)
		return nil
	case environmentconnectorrelationship.FieldConnectorID:
		v, ok := value.(types.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetConnectorID(v)
		return nil
	}
	return fmt.Errorf("unknown EnvironmentConnectorRelationship field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *EnvironmentConnectorRelationshipMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *EnvironmentConnectorRelationshipMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *EnvironmentConnectorRelationshipMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown EnvironmentConnectorRelationship numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *EnvironmentConnectorRelationshipMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *EnvironmentConnectorRelationshipMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *EnvironmentConnectorRelationshipMutation) ClearField(name string) error {
	return fmt.Errorf("unknown EnvironmentConnectorRelationship nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *EnvironmentConnectorRelationshipMutation) ResetField(name string) error {
	switch name {
	case environmentconnectorrelationship.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case environmentconnectorrelationship.FieldEnvironmentID:
		m.ResetEnvironmentID()
		return nil
	case environmentconnectorrelationship.FieldConnectorID:
		m.ResetConnectorID()
		return nil
	}
	return fmt.Errorf("unknown EnvironmentConnectorRelationship field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *EnvironmentConnectorRelationshipMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.environment != nil {
		edges = append(edges, environmentconnectorrelationship.EdgeEnvironment)
	}
	if m.connector != nil {
		edges = append(edges, environmentconnectorrelationship.EdgeConnector)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *EnvironmentConnectorRelationshipMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case environmentconnectorrelationship.EdgeEnvironment:
		if id := m.environment; id != nil {
			return []ent.Value{*id}
		}
	case environmentconnectorrelationship.EdgeConnector:
		if id := m.connector; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *EnvironmentConnectorRelationshipMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *EnvironmentConnectorRelationshipMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *EnvironmentConnectorRelationshipMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedenvironment {
		edges = append(edges, environmentconnectorrelationship.EdgeEnvironment)
	}
	if m.clearedconnector {
		edges = append(edges, environmentconnectorrelationship.EdgeConnector)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *EnvironmentConnectorRelationshipMutation) EdgeCleared(name string) bool {
	switch name {
	case environmentconnectorrelationship.EdgeEnvironment:
		return m.clearedenvironment
	case environmentconnectorrelationship.EdgeConnector:
		return m.clearedconnector
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *EnvironmentConnectorRelationshipMutation) ClearEdge(name string) error {
	switch name {
	case environmentconnectorrelationship.EdgeEnvironment:
		m.ClearEnvironment()
		return nil
	case environmentconnectorrelationship.EdgeConnector:
		m.ClearConnector()
		return nil
	}
	return fmt.Errorf("unknown EnvironmentConnectorRelationship unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *EnvironmentConnectorRelationshipMutation) ResetEdge(name string) error {
	switch name {
	case environmentconnectorrelationship.EdgeEnvironment:
		m.ResetEnvironment()
		return nil
	case environmentconnectorrelationship.EdgeConnector:
		m.ResetConnector()
		return nil
	}
	return fmt.Errorf("unknown EnvironmentConnectorRelationship edge %s", name)
}

// ModuleMutation represents an operation that mutates the Module nodes in the graph.
type ModuleMutation struct {
	config
	op                                    Op
	typ                                   string
	id                                    *string
	status                                *string
	statusMessage                         *string
	createTime                            *time.Time
	updateTime                            *time.Time
	description                           *string
	labels                                *map[string]string
	source                                *string
	version                               *string
	schema                                **types.ModuleSchema
	clearedFields                         map[string]struct{}
	application                           map[types.ID]struct{}
	removedapplication                    map[types.ID]struct{}
	clearedapplication                    bool
	applicationModuleRelationships        map[int]struct{}
	removedapplicationModuleRelationships map[int]struct{}
	clearedapplicationModuleRelationships bool
	done                                  bool
	oldValue                              func(context.Context) (*Module, error)
	predicates                            []predicate.Module
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
func withModuleID(id string) moduleOption {
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
func (m *ModuleMutation) SetID(id string) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ModuleMutation) ID() (id string, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ModuleMutation) IDs(ctx context.Context) ([]string, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []string{id}, nil
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

// SetDescription sets the "description" field.
func (m *ModuleMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the value of the "description" field in the mutation.
func (m *ModuleMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// OldDescription returns the old "description" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldDescription(ctx context.Context) (v string, err error) {
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

// ClearDescription clears the value of the "description" field.
func (m *ModuleMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[module.FieldDescription] = struct{}{}
}

// DescriptionCleared returns if the "description" field was cleared in this mutation.
func (m *ModuleMutation) DescriptionCleared() bool {
	_, ok := m.clearedFields[module.FieldDescription]
	return ok
}

// ResetDescription resets all changes to the "description" field.
func (m *ModuleMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, module.FieldDescription)
}

// SetLabels sets the "labels" field.
func (m *ModuleMutation) SetLabels(value map[string]string) {
	m.labels = &value
}

// Labels returns the value of the "labels" field in the mutation.
func (m *ModuleMutation) Labels() (r map[string]string, exists bool) {
	v := m.labels
	if v == nil {
		return
	}
	return *v, true
}

// OldLabels returns the old "labels" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldLabels(ctx context.Context) (v map[string]string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLabels is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLabels requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLabels: %w", err)
	}
	return oldValue.Labels, nil
}

// ResetLabels resets all changes to the "labels" field.
func (m *ModuleMutation) ResetLabels() {
	m.labels = nil
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

// ClearVersion clears the value of the "version" field.
func (m *ModuleMutation) ClearVersion() {
	m.version = nil
	m.clearedFields[module.FieldVersion] = struct{}{}
}

// VersionCleared returns if the "version" field was cleared in this mutation.
func (m *ModuleMutation) VersionCleared() bool {
	_, ok := m.clearedFields[module.FieldVersion]
	return ok
}

// ResetVersion resets all changes to the "version" field.
func (m *ModuleMutation) ResetVersion() {
	m.version = nil
	delete(m.clearedFields, module.FieldVersion)
}

// SetSchema sets the "schema" field.
func (m *ModuleMutation) SetSchema(ts *types.ModuleSchema) {
	m.schema = &ts
}

// Schema returns the value of the "schema" field in the mutation.
func (m *ModuleMutation) Schema() (r *types.ModuleSchema, exists bool) {
	v := m.schema
	if v == nil {
		return
	}
	return *v, true
}

// OldSchema returns the old "schema" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldSchema(ctx context.Context) (v *types.ModuleSchema, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSchema is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSchema requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSchema: %w", err)
	}
	return oldValue.Schema, nil
}

// ResetSchema resets all changes to the "schema" field.
func (m *ModuleMutation) ResetSchema() {
	m.schema = nil
}

// AddApplicationIDs adds the "application" edge to the Application entity by ids.
func (m *ModuleMutation) AddApplicationIDs(ids ...types.ID) {
	if m.application == nil {
		m.application = make(map[types.ID]struct{})
	}
	for i := range ids {
		m.application[ids[i]] = struct{}{}
	}
}

// ClearApplication clears the "application" edge to the Application entity.
func (m *ModuleMutation) ClearApplication() {
	m.clearedapplication = true
}

// ApplicationCleared reports if the "application" edge to the Application entity was cleared.
func (m *ModuleMutation) ApplicationCleared() bool {
	return m.clearedapplication
}

// RemoveApplicationIDs removes the "application" edge to the Application entity by IDs.
func (m *ModuleMutation) RemoveApplicationIDs(ids ...types.ID) {
	if m.removedapplication == nil {
		m.removedapplication = make(map[types.ID]struct{})
	}
	for i := range ids {
		delete(m.application, ids[i])
		m.removedapplication[ids[i]] = struct{}{}
	}
}

// RemovedApplication returns the removed IDs of the "application" edge to the Application entity.
func (m *ModuleMutation) RemovedApplicationIDs() (ids []types.ID) {
	for id := range m.removedapplication {
		ids = append(ids, id)
	}
	return
}

// ApplicationIDs returns the "application" edge IDs in the mutation.
func (m *ModuleMutation) ApplicationIDs() (ids []types.ID) {
	for id := range m.application {
		ids = append(ids, id)
	}
	return
}

// ResetApplication resets all changes to the "application" edge.
func (m *ModuleMutation) ResetApplication() {
	m.application = nil
	m.clearedapplication = false
	m.removedapplication = nil
}

// AddApplicationModuleRelationshipIDs adds the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity by ids.
func (m *ModuleMutation) AddApplicationModuleRelationshipIDs(ids ...int) {
	if m.applicationModuleRelationships == nil {
		m.applicationModuleRelationships = make(map[int]struct{})
	}
	for i := range ids {
		m.applicationModuleRelationships[ids[i]] = struct{}{}
	}
}

// ClearApplicationModuleRelationships clears the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity.
func (m *ModuleMutation) ClearApplicationModuleRelationships() {
	m.clearedapplicationModuleRelationships = true
}

// ApplicationModuleRelationshipsCleared reports if the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity was cleared.
func (m *ModuleMutation) ApplicationModuleRelationshipsCleared() bool {
	return m.clearedapplicationModuleRelationships
}

// RemoveApplicationModuleRelationshipIDs removes the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity by IDs.
func (m *ModuleMutation) RemoveApplicationModuleRelationshipIDs(ids ...int) {
	if m.removedapplicationModuleRelationships == nil {
		m.removedapplicationModuleRelationships = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.applicationModuleRelationships, ids[i])
		m.removedapplicationModuleRelationships[ids[i]] = struct{}{}
	}
}

// RemovedApplicationModuleRelationships returns the removed IDs of the "applicationModuleRelationships" edge to the ApplicationModuleRelationship entity.
func (m *ModuleMutation) RemovedApplicationModuleRelationshipsIDs() (ids []int) {
	for id := range m.removedapplicationModuleRelationships {
		ids = append(ids, id)
	}
	return
}

// ApplicationModuleRelationshipsIDs returns the "applicationModuleRelationships" edge IDs in the mutation.
func (m *ModuleMutation) ApplicationModuleRelationshipsIDs() (ids []int) {
	for id := range m.applicationModuleRelationships {
		ids = append(ids, id)
	}
	return
}

// ResetApplicationModuleRelationships resets all changes to the "applicationModuleRelationships" edge.
func (m *ModuleMutation) ResetApplicationModuleRelationships() {
	m.applicationModuleRelationships = nil
	m.clearedapplicationModuleRelationships = false
	m.removedapplicationModuleRelationships = nil
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
	fields := make([]string, 0, 9)
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
	if m.description != nil {
		fields = append(fields, module.FieldDescription)
	}
	if m.labels != nil {
		fields = append(fields, module.FieldLabels)
	}
	if m.source != nil {
		fields = append(fields, module.FieldSource)
	}
	if m.version != nil {
		fields = append(fields, module.FieldVersion)
	}
	if m.schema != nil {
		fields = append(fields, module.FieldSchema)
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
	case module.FieldDescription:
		return m.Description()
	case module.FieldLabels:
		return m.Labels()
	case module.FieldSource:
		return m.Source()
	case module.FieldVersion:
		return m.Version()
	case module.FieldSchema:
		return m.Schema()
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
	case module.FieldDescription:
		return m.OldDescription(ctx)
	case module.FieldLabels:
		return m.OldLabels(ctx)
	case module.FieldSource:
		return m.OldSource(ctx)
	case module.FieldVersion:
		return m.OldVersion(ctx)
	case module.FieldSchema:
		return m.OldSchema(ctx)
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
	case module.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case module.FieldLabels:
		v, ok := value.(map[string]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLabels(v)
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
	case module.FieldSchema:
		v, ok := value.(*types.ModuleSchema)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSchema(v)
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
	if m.FieldCleared(module.FieldDescription) {
		fields = append(fields, module.FieldDescription)
	}
	if m.FieldCleared(module.FieldVersion) {
		fields = append(fields, module.FieldVersion)
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
	case module.FieldDescription:
		m.ClearDescription()
		return nil
	case module.FieldVersion:
		m.ClearVersion()
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
	case module.FieldDescription:
		m.ResetDescription()
		return nil
	case module.FieldLabels:
		m.ResetLabels()
		return nil
	case module.FieldSource:
		m.ResetSource()
		return nil
	case module.FieldVersion:
		m.ResetVersion()
		return nil
	case module.FieldSchema:
		m.ResetSchema()
		return nil
	}
	return fmt.Errorf("unknown Module field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ModuleMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.application != nil {
		edges = append(edges, module.EdgeApplication)
	}
	if m.applicationModuleRelationships != nil {
		edges = append(edges, module.EdgeApplicationModuleRelationships)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ModuleMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case module.EdgeApplication:
		ids := make([]ent.Value, 0, len(m.application))
		for id := range m.application {
			ids = append(ids, id)
		}
		return ids
	case module.EdgeApplicationModuleRelationships:
		ids := make([]ent.Value, 0, len(m.applicationModuleRelationships))
		for id := range m.applicationModuleRelationships {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ModuleMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	if m.removedapplication != nil {
		edges = append(edges, module.EdgeApplication)
	}
	if m.removedapplicationModuleRelationships != nil {
		edges = append(edges, module.EdgeApplicationModuleRelationships)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ModuleMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case module.EdgeApplication:
		ids := make([]ent.Value, 0, len(m.removedapplication))
		for id := range m.removedapplication {
			ids = append(ids, id)
		}
		return ids
	case module.EdgeApplicationModuleRelationships:
		ids := make([]ent.Value, 0, len(m.removedapplicationModuleRelationships))
		for id := range m.removedapplicationModuleRelationships {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ModuleMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedapplication {
		edges = append(edges, module.EdgeApplication)
	}
	if m.clearedapplicationModuleRelationships {
		edges = append(edges, module.EdgeApplicationModuleRelationships)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ModuleMutation) EdgeCleared(name string) bool {
	switch name {
	case module.EdgeApplication:
		return m.clearedapplication
	case module.EdgeApplicationModuleRelationships:
		return m.clearedapplicationModuleRelationships
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ModuleMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Module unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ModuleMutation) ResetEdge(name string) error {
	switch name {
	case module.EdgeApplication:
		m.ResetApplication()
		return nil
	case module.EdgeApplicationModuleRelationships:
		m.ResetApplicationModuleRelationships()
		return nil
	}
	return fmt.Errorf("unknown Module edge %s", name)
}

// ProjectMutation represents an operation that mutates the Project nodes in the graph.
type ProjectMutation struct {
	config
	op                  Op
	typ                 string
	id                  *types.ID
	name                *string
	description         *string
	labels              *map[string]string
	createTime          *time.Time
	updateTime          *time.Time
	clearedFields       map[string]struct{}
	applications        map[types.ID]struct{}
	removedapplications map[types.ID]struct{}
	clearedapplications bool
	done                bool
	oldValue            func(context.Context) (*Project, error)
	predicates          []predicate.Project
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
func withProjectID(id types.ID) projectOption {
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
func (m *ProjectMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ProjectMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ProjectMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Project.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *ProjectMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *ProjectMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Project entity.
// If the Project object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProjectMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *ProjectMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the "description" field.
func (m *ProjectMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the value of the "description" field in the mutation.
func (m *ProjectMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// OldDescription returns the old "description" field's value of the Project entity.
// If the Project object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProjectMutation) OldDescription(ctx context.Context) (v string, err error) {
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

// ClearDescription clears the value of the "description" field.
func (m *ProjectMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[project.FieldDescription] = struct{}{}
}

// DescriptionCleared returns if the "description" field was cleared in this mutation.
func (m *ProjectMutation) DescriptionCleared() bool {
	_, ok := m.clearedFields[project.FieldDescription]
	return ok
}

// ResetDescription resets all changes to the "description" field.
func (m *ProjectMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, project.FieldDescription)
}

// SetLabels sets the "labels" field.
func (m *ProjectMutation) SetLabels(value map[string]string) {
	m.labels = &value
}

// Labels returns the value of the "labels" field in the mutation.
func (m *ProjectMutation) Labels() (r map[string]string, exists bool) {
	v := m.labels
	if v == nil {
		return
	}
	return *v, true
}

// OldLabels returns the old "labels" field's value of the Project entity.
// If the Project object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ProjectMutation) OldLabels(ctx context.Context) (v map[string]string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLabels is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLabels requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLabels: %w", err)
	}
	return oldValue.Labels, nil
}

// ResetLabels resets all changes to the "labels" field.
func (m *ProjectMutation) ResetLabels() {
	m.labels = nil
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

// AddApplicationIDs adds the "applications" edge to the Application entity by ids.
func (m *ProjectMutation) AddApplicationIDs(ids ...types.ID) {
	if m.applications == nil {
		m.applications = make(map[types.ID]struct{})
	}
	for i := range ids {
		m.applications[ids[i]] = struct{}{}
	}
}

// ClearApplications clears the "applications" edge to the Application entity.
func (m *ProjectMutation) ClearApplications() {
	m.clearedapplications = true
}

// ApplicationsCleared reports if the "applications" edge to the Application entity was cleared.
func (m *ProjectMutation) ApplicationsCleared() bool {
	return m.clearedapplications
}

// RemoveApplicationIDs removes the "applications" edge to the Application entity by IDs.
func (m *ProjectMutation) RemoveApplicationIDs(ids ...types.ID) {
	if m.removedapplications == nil {
		m.removedapplications = make(map[types.ID]struct{})
	}
	for i := range ids {
		delete(m.applications, ids[i])
		m.removedapplications[ids[i]] = struct{}{}
	}
}

// RemovedApplications returns the removed IDs of the "applications" edge to the Application entity.
func (m *ProjectMutation) RemovedApplicationsIDs() (ids []types.ID) {
	for id := range m.removedapplications {
		ids = append(ids, id)
	}
	return
}

// ApplicationsIDs returns the "applications" edge IDs in the mutation.
func (m *ProjectMutation) ApplicationsIDs() (ids []types.ID) {
	for id := range m.applications {
		ids = append(ids, id)
	}
	return
}

// ResetApplications resets all changes to the "applications" edge.
func (m *ProjectMutation) ResetApplications() {
	m.applications = nil
	m.clearedapplications = false
	m.removedapplications = nil
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
	fields := make([]string, 0, 5)
	if m.name != nil {
		fields = append(fields, project.FieldName)
	}
	if m.description != nil {
		fields = append(fields, project.FieldDescription)
	}
	if m.labels != nil {
		fields = append(fields, project.FieldLabels)
	}
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
	case project.FieldName:
		return m.Name()
	case project.FieldDescription:
		return m.Description()
	case project.FieldLabels:
		return m.Labels()
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
	case project.FieldName:
		return m.OldName(ctx)
	case project.FieldDescription:
		return m.OldDescription(ctx)
	case project.FieldLabels:
		return m.OldLabels(ctx)
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
	case project.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case project.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case project.FieldLabels:
		v, ok := value.(map[string]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLabels(v)
		return nil
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
	var fields []string
	if m.FieldCleared(project.FieldDescription) {
		fields = append(fields, project.FieldDescription)
	}
	return fields
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
	switch name {
	case project.FieldDescription:
		m.ClearDescription()
		return nil
	}
	return fmt.Errorf("unknown Project nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ProjectMutation) ResetField(name string) error {
	switch name {
	case project.FieldName:
		m.ResetName()
		return nil
	case project.FieldDescription:
		m.ResetDescription()
		return nil
	case project.FieldLabels:
		m.ResetLabels()
		return nil
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
	edges := make([]string, 0, 1)
	if m.applications != nil {
		edges = append(edges, project.EdgeApplications)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ProjectMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case project.EdgeApplications:
		ids := make([]ent.Value, 0, len(m.applications))
		for id := range m.applications {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ProjectMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedapplications != nil {
		edges = append(edges, project.EdgeApplications)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ProjectMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case project.EdgeApplications:
		ids := make([]ent.Value, 0, len(m.removedapplications))
		for id := range m.removedapplications {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ProjectMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedapplications {
		edges = append(edges, project.EdgeApplications)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ProjectMutation) EdgeCleared(name string) bool {
	switch name {
	case project.EdgeApplications:
		return m.clearedapplications
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ProjectMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Project unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ProjectMutation) ResetEdge(name string) error {
	switch name {
	case project.EdgeApplications:
		m.ResetApplications()
		return nil
	}
	return fmt.Errorf("unknown Project edge %s", name)
}

// RoleMutation represents an operation that mutates the Role nodes in the graph.
type RoleMutation struct {
	config
	op             Op
	typ            string
	id             *types.ID
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
func withRoleID(id types.ID) roleOption {
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
func (m *RoleMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *RoleMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *RoleMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
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

// ClearDescription clears the value of the "description" field.
func (m *RoleMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[role.FieldDescription] = struct{}{}
}

// DescriptionCleared returns if the "description" field was cleared in this mutation.
func (m *RoleMutation) DescriptionCleared() bool {
	_, ok := m.clearedFields[role.FieldDescription]
	return ok
}

// ResetDescription resets all changes to the "description" field.
func (m *RoleMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, role.FieldDescription)
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
	var fields []string
	if m.FieldCleared(role.FieldDescription) {
		fields = append(fields, role.FieldDescription)
	}
	return fields
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
	switch name {
	case role.FieldDescription:
		m.ClearDescription()
		return nil
	}
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
	id            *types.ID
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
func withSettingID(id types.ID) settingOption {
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
func (m *SettingMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SettingMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SettingMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
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
	id            *types.ID
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
func withSubjectID(id types.ID) subjectOption {
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
func (m *SubjectMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SubjectMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SubjectMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
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

// ClearDescription clears the value of the "description" field.
func (m *SubjectMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[subject.FieldDescription] = struct{}{}
}

// DescriptionCleared returns if the "description" field was cleared in this mutation.
func (m *SubjectMutation) DescriptionCleared() bool {
	_, ok := m.clearedFields[subject.FieldDescription]
	return ok
}

// ResetDescription resets all changes to the "description" field.
func (m *SubjectMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, subject.FieldDescription)
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
	var fields []string
	if m.FieldCleared(subject.FieldDescription) {
		fields = append(fields, subject.FieldDescription)
	}
	return fields
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
	switch name {
	case subject.FieldDescription:
		m.ClearDescription()
		return nil
	}
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
	id                *types.ID
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
func withTokenID(id types.ID) tokenOption {
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
func (m *TokenMutation) SetID(id types.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *TokenMutation) ID() (id types.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *TokenMutation) IDs(ctx context.Context) ([]types.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []types.ID{id}, nil
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
