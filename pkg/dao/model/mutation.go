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

	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"

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
	TypeAllocationCost                   = "AllocationCost"
	TypeApplication                      = "Application"
	TypeApplicationInstance              = "ApplicationInstance"
	TypeApplicationModuleRelationship    = "ApplicationModuleRelationship"
	TypeApplicationResource              = "ApplicationResource"
	TypeApplicationRevision              = "ApplicationRevision"
	TypeClusterCost                      = "ClusterCost"
	TypeConnector                        = "Connector"
	TypeEnvironment                      = "Environment"
	TypeEnvironmentConnectorRelationship = "EnvironmentConnectorRelationship"
	TypeModule                           = "Module"
	TypeModuleVersion                    = "ModuleVersion"
	TypePerspective                      = "Perspective"
	TypeProject                          = "Project"
	TypeRole                             = "Role"
	TypeSecret                           = "Secret"
	TypeSetting                          = "Setting"
	TypeSubject                          = "Subject"
	TypeToken                            = "Token"
)

// AllocationCostMutation represents an operation that mutates the AllocationCost nodes in the graph.
type AllocationCostMutation struct {
	config
	op                     Op
	typ                    string
	id                     *int
	startTime              *time.Time
	endTime                *time.Time
	minutes                *float64
	addminutes             *float64
	name                   *string
	fingerprint            *string
	clusterName            *string
	namespace              *string
	node                   *string
	controller             *string
	controllerKind         *string
	pod                    *string
	container              *string
	pvs                    *map[string]types.PVCost
	labels                 *map[string]string
	totalCost              *float64
	addtotalCost           *float64
	currency               *int
	addcurrency            *int
	cpuCost                *float64
	addcpuCost             *float64
	cpuCoreRequest         *float64
	addcpuCoreRequest      *float64
	gpuCost                *float64
	addgpuCost             *float64
	gpuCount               *float64
	addgpuCount            *float64
	ramCost                *float64
	addramCost             *float64
	ramByteRequest         *float64
	addramByteRequest      *float64
	pvCost                 *float64
	addpvCost              *float64
	pvBytes                *float64
	addpvBytes             *float64
	loadBalancerCost       *float64
	addloadBalancerCost    *float64
	cpuCoreUsageAverage    *float64
	addcpuCoreUsageAverage *float64
	cpuCoreUsageMax        *float64
	addcpuCoreUsageMax     *float64
	ramByteUsageAverage    *float64
	addramByteUsageAverage *float64
	ramByteUsageMax        *float64
	addramByteUsageMax     *float64
	clearedFields          map[string]struct{}
	connector              *oid.ID
	clearedconnector       bool
	done                   bool
	oldValue               func(context.Context) (*AllocationCost, error)
	predicates             []predicate.AllocationCost
}

var _ ent.Mutation = (*AllocationCostMutation)(nil)

// allocationCostOption allows management of the mutation configuration using functional options.
type allocationCostOption func(*AllocationCostMutation)

// newAllocationCostMutation creates new mutation for the AllocationCost entity.
func newAllocationCostMutation(c config, op Op, opts ...allocationCostOption) *AllocationCostMutation {
	m := &AllocationCostMutation{
		config:        c,
		op:            op,
		typ:           TypeAllocationCost,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withAllocationCostID sets the ID field of the mutation.
func withAllocationCostID(id int) allocationCostOption {
	return func(m *AllocationCostMutation) {
		var (
			err   error
			once  sync.Once
			value *AllocationCost
		)
		m.oldValue = func(ctx context.Context) (*AllocationCost, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().AllocationCost.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withAllocationCost sets the old AllocationCost of the mutation.
func withAllocationCost(node *AllocationCost) allocationCostOption {
	return func(m *AllocationCostMutation) {
		m.oldValue = func(context.Context) (*AllocationCost, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m AllocationCostMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m AllocationCostMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *AllocationCostMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *AllocationCostMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().AllocationCost.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetStartTime sets the "startTime" field.
func (m *AllocationCostMutation) SetStartTime(t time.Time) {
	m.startTime = &t
}

// StartTime returns the value of the "startTime" field in the mutation.
func (m *AllocationCostMutation) StartTime() (r time.Time, exists bool) {
	v := m.startTime
	if v == nil {
		return
	}
	return *v, true
}

// OldStartTime returns the old "startTime" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldStartTime(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStartTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStartTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStartTime: %w", err)
	}
	return oldValue.StartTime, nil
}

// ResetStartTime resets all changes to the "startTime" field.
func (m *AllocationCostMutation) ResetStartTime() {
	m.startTime = nil
}

// SetEndTime sets the "endTime" field.
func (m *AllocationCostMutation) SetEndTime(t time.Time) {
	m.endTime = &t
}

// EndTime returns the value of the "endTime" field in the mutation.
func (m *AllocationCostMutation) EndTime() (r time.Time, exists bool) {
	v := m.endTime
	if v == nil {
		return
	}
	return *v, true
}

// OldEndTime returns the old "endTime" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldEndTime(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEndTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEndTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEndTime: %w", err)
	}
	return oldValue.EndTime, nil
}

// ResetEndTime resets all changes to the "endTime" field.
func (m *AllocationCostMutation) ResetEndTime() {
	m.endTime = nil
}

// SetMinutes sets the "minutes" field.
func (m *AllocationCostMutation) SetMinutes(f float64) {
	m.minutes = &f
	m.addminutes = nil
}

// Minutes returns the value of the "minutes" field in the mutation.
func (m *AllocationCostMutation) Minutes() (r float64, exists bool) {
	v := m.minutes
	if v == nil {
		return
	}
	return *v, true
}

// OldMinutes returns the old "minutes" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldMinutes(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldMinutes is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldMinutes requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldMinutes: %w", err)
	}
	return oldValue.Minutes, nil
}

// AddMinutes adds f to the "minutes" field.
func (m *AllocationCostMutation) AddMinutes(f float64) {
	if m.addminutes != nil {
		*m.addminutes += f
	} else {
		m.addminutes = &f
	}
}

// AddedMinutes returns the value that was added to the "minutes" field in this mutation.
func (m *AllocationCostMutation) AddedMinutes() (r float64, exists bool) {
	v := m.addminutes
	if v == nil {
		return
	}
	return *v, true
}

// ResetMinutes resets all changes to the "minutes" field.
func (m *AllocationCostMutation) ResetMinutes() {
	m.minutes = nil
	m.addminutes = nil
}

// SetConnectorID sets the "connectorID" field.
func (m *AllocationCostMutation) SetConnectorID(o oid.ID) {
	m.connector = &o
}

// ConnectorID returns the value of the "connectorID" field in the mutation.
func (m *AllocationCostMutation) ConnectorID() (r oid.ID, exists bool) {
	v := m.connector
	if v == nil {
		return
	}
	return *v, true
}

// OldConnectorID returns the old "connectorID" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldConnectorID(ctx context.Context) (v oid.ID, err error) {
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
func (m *AllocationCostMutation) ResetConnectorID() {
	m.connector = nil
}

// SetName sets the "name" field.
func (m *AllocationCostMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *AllocationCostMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *AllocationCostMutation) ResetName() {
	m.name = nil
}

// SetFingerprint sets the "fingerprint" field.
func (m *AllocationCostMutation) SetFingerprint(s string) {
	m.fingerprint = &s
}

// Fingerprint returns the value of the "fingerprint" field in the mutation.
func (m *AllocationCostMutation) Fingerprint() (r string, exists bool) {
	v := m.fingerprint
	if v == nil {
		return
	}
	return *v, true
}

// OldFingerprint returns the old "fingerprint" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldFingerprint(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFingerprint is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFingerprint requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFingerprint: %w", err)
	}
	return oldValue.Fingerprint, nil
}

// ResetFingerprint resets all changes to the "fingerprint" field.
func (m *AllocationCostMutation) ResetFingerprint() {
	m.fingerprint = nil
}

// SetClusterName sets the "clusterName" field.
func (m *AllocationCostMutation) SetClusterName(s string) {
	m.clusterName = &s
}

// ClusterName returns the value of the "clusterName" field in the mutation.
func (m *AllocationCostMutation) ClusterName() (r string, exists bool) {
	v := m.clusterName
	if v == nil {
		return
	}
	return *v, true
}

// OldClusterName returns the old "clusterName" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldClusterName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldClusterName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldClusterName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldClusterName: %w", err)
	}
	return oldValue.ClusterName, nil
}

// ClearClusterName clears the value of the "clusterName" field.
func (m *AllocationCostMutation) ClearClusterName() {
	m.clusterName = nil
	m.clearedFields[allocationcost.FieldClusterName] = struct{}{}
}

// ClusterNameCleared returns if the "clusterName" field was cleared in this mutation.
func (m *AllocationCostMutation) ClusterNameCleared() bool {
	_, ok := m.clearedFields[allocationcost.FieldClusterName]
	return ok
}

// ResetClusterName resets all changes to the "clusterName" field.
func (m *AllocationCostMutation) ResetClusterName() {
	m.clusterName = nil
	delete(m.clearedFields, allocationcost.FieldClusterName)
}

// SetNamespace sets the "namespace" field.
func (m *AllocationCostMutation) SetNamespace(s string) {
	m.namespace = &s
}

// Namespace returns the value of the "namespace" field in the mutation.
func (m *AllocationCostMutation) Namespace() (r string, exists bool) {
	v := m.namespace
	if v == nil {
		return
	}
	return *v, true
}

// OldNamespace returns the old "namespace" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldNamespace(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldNamespace is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldNamespace requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldNamespace: %w", err)
	}
	return oldValue.Namespace, nil
}

// ClearNamespace clears the value of the "namespace" field.
func (m *AllocationCostMutation) ClearNamespace() {
	m.namespace = nil
	m.clearedFields[allocationcost.FieldNamespace] = struct{}{}
}

// NamespaceCleared returns if the "namespace" field was cleared in this mutation.
func (m *AllocationCostMutation) NamespaceCleared() bool {
	_, ok := m.clearedFields[allocationcost.FieldNamespace]
	return ok
}

// ResetNamespace resets all changes to the "namespace" field.
func (m *AllocationCostMutation) ResetNamespace() {
	m.namespace = nil
	delete(m.clearedFields, allocationcost.FieldNamespace)
}

// SetNode sets the "node" field.
func (m *AllocationCostMutation) SetNode(s string) {
	m.node = &s
}

// Node returns the value of the "node" field in the mutation.
func (m *AllocationCostMutation) Node() (r string, exists bool) {
	v := m.node
	if v == nil {
		return
	}
	return *v, true
}

// OldNode returns the old "node" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldNode(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldNode is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldNode requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldNode: %w", err)
	}
	return oldValue.Node, nil
}

// ClearNode clears the value of the "node" field.
func (m *AllocationCostMutation) ClearNode() {
	m.node = nil
	m.clearedFields[allocationcost.FieldNode] = struct{}{}
}

// NodeCleared returns if the "node" field was cleared in this mutation.
func (m *AllocationCostMutation) NodeCleared() bool {
	_, ok := m.clearedFields[allocationcost.FieldNode]
	return ok
}

// ResetNode resets all changes to the "node" field.
func (m *AllocationCostMutation) ResetNode() {
	m.node = nil
	delete(m.clearedFields, allocationcost.FieldNode)
}

// SetController sets the "controller" field.
func (m *AllocationCostMutation) SetController(s string) {
	m.controller = &s
}

// Controller returns the value of the "controller" field in the mutation.
func (m *AllocationCostMutation) Controller() (r string, exists bool) {
	v := m.controller
	if v == nil {
		return
	}
	return *v, true
}

// OldController returns the old "controller" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldController(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldController is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldController requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldController: %w", err)
	}
	return oldValue.Controller, nil
}

// ClearController clears the value of the "controller" field.
func (m *AllocationCostMutation) ClearController() {
	m.controller = nil
	m.clearedFields[allocationcost.FieldController] = struct{}{}
}

// ControllerCleared returns if the "controller" field was cleared in this mutation.
func (m *AllocationCostMutation) ControllerCleared() bool {
	_, ok := m.clearedFields[allocationcost.FieldController]
	return ok
}

// ResetController resets all changes to the "controller" field.
func (m *AllocationCostMutation) ResetController() {
	m.controller = nil
	delete(m.clearedFields, allocationcost.FieldController)
}

// SetControllerKind sets the "controllerKind" field.
func (m *AllocationCostMutation) SetControllerKind(s string) {
	m.controllerKind = &s
}

// ControllerKind returns the value of the "controllerKind" field in the mutation.
func (m *AllocationCostMutation) ControllerKind() (r string, exists bool) {
	v := m.controllerKind
	if v == nil {
		return
	}
	return *v, true
}

// OldControllerKind returns the old "controllerKind" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldControllerKind(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldControllerKind is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldControllerKind requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldControllerKind: %w", err)
	}
	return oldValue.ControllerKind, nil
}

// ClearControllerKind clears the value of the "controllerKind" field.
func (m *AllocationCostMutation) ClearControllerKind() {
	m.controllerKind = nil
	m.clearedFields[allocationcost.FieldControllerKind] = struct{}{}
}

// ControllerKindCleared returns if the "controllerKind" field was cleared in this mutation.
func (m *AllocationCostMutation) ControllerKindCleared() bool {
	_, ok := m.clearedFields[allocationcost.FieldControllerKind]
	return ok
}

// ResetControllerKind resets all changes to the "controllerKind" field.
func (m *AllocationCostMutation) ResetControllerKind() {
	m.controllerKind = nil
	delete(m.clearedFields, allocationcost.FieldControllerKind)
}

// SetPod sets the "pod" field.
func (m *AllocationCostMutation) SetPod(s string) {
	m.pod = &s
}

// Pod returns the value of the "pod" field in the mutation.
func (m *AllocationCostMutation) Pod() (r string, exists bool) {
	v := m.pod
	if v == nil {
		return
	}
	return *v, true
}

// OldPod returns the old "pod" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldPod(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPod is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPod requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPod: %w", err)
	}
	return oldValue.Pod, nil
}

// ClearPod clears the value of the "pod" field.
func (m *AllocationCostMutation) ClearPod() {
	m.pod = nil
	m.clearedFields[allocationcost.FieldPod] = struct{}{}
}

// PodCleared returns if the "pod" field was cleared in this mutation.
func (m *AllocationCostMutation) PodCleared() bool {
	_, ok := m.clearedFields[allocationcost.FieldPod]
	return ok
}

// ResetPod resets all changes to the "pod" field.
func (m *AllocationCostMutation) ResetPod() {
	m.pod = nil
	delete(m.clearedFields, allocationcost.FieldPod)
}

// SetContainer sets the "container" field.
func (m *AllocationCostMutation) SetContainer(s string) {
	m.container = &s
}

// Container returns the value of the "container" field in the mutation.
func (m *AllocationCostMutation) Container() (r string, exists bool) {
	v := m.container
	if v == nil {
		return
	}
	return *v, true
}

// OldContainer returns the old "container" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldContainer(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldContainer is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldContainer requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldContainer: %w", err)
	}
	return oldValue.Container, nil
}

// ClearContainer clears the value of the "container" field.
func (m *AllocationCostMutation) ClearContainer() {
	m.container = nil
	m.clearedFields[allocationcost.FieldContainer] = struct{}{}
}

// ContainerCleared returns if the "container" field was cleared in this mutation.
func (m *AllocationCostMutation) ContainerCleared() bool {
	_, ok := m.clearedFields[allocationcost.FieldContainer]
	return ok
}

// ResetContainer resets all changes to the "container" field.
func (m *AllocationCostMutation) ResetContainer() {
	m.container = nil
	delete(m.clearedFields, allocationcost.FieldContainer)
}

// SetPvs sets the "pvs" field.
func (m *AllocationCostMutation) SetPvs(mc map[string]types.PVCost) {
	m.pvs = &mc
}

// Pvs returns the value of the "pvs" field in the mutation.
func (m *AllocationCostMutation) Pvs() (r map[string]types.PVCost, exists bool) {
	v := m.pvs
	if v == nil {
		return
	}
	return *v, true
}

// OldPvs returns the old "pvs" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldPvs(ctx context.Context) (v map[string]types.PVCost, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPvs is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPvs requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPvs: %w", err)
	}
	return oldValue.Pvs, nil
}

// ResetPvs resets all changes to the "pvs" field.
func (m *AllocationCostMutation) ResetPvs() {
	m.pvs = nil
}

// SetLabels sets the "labels" field.
func (m *AllocationCostMutation) SetLabels(value map[string]string) {
	m.labels = &value
}

// Labels returns the value of the "labels" field in the mutation.
func (m *AllocationCostMutation) Labels() (r map[string]string, exists bool) {
	v := m.labels
	if v == nil {
		return
	}
	return *v, true
}

// OldLabels returns the old "labels" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldLabels(ctx context.Context) (v map[string]string, err error) {
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
func (m *AllocationCostMutation) ResetLabels() {
	m.labels = nil
}

// SetTotalCost sets the "totalCost" field.
func (m *AllocationCostMutation) SetTotalCost(f float64) {
	m.totalCost = &f
	m.addtotalCost = nil
}

// TotalCost returns the value of the "totalCost" field in the mutation.
func (m *AllocationCostMutation) TotalCost() (r float64, exists bool) {
	v := m.totalCost
	if v == nil {
		return
	}
	return *v, true
}

// OldTotalCost returns the old "totalCost" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldTotalCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTotalCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTotalCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTotalCost: %w", err)
	}
	return oldValue.TotalCost, nil
}

// AddTotalCost adds f to the "totalCost" field.
func (m *AllocationCostMutation) AddTotalCost(f float64) {
	if m.addtotalCost != nil {
		*m.addtotalCost += f
	} else {
		m.addtotalCost = &f
	}
}

// AddedTotalCost returns the value that was added to the "totalCost" field in this mutation.
func (m *AllocationCostMutation) AddedTotalCost() (r float64, exists bool) {
	v := m.addtotalCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetTotalCost resets all changes to the "totalCost" field.
func (m *AllocationCostMutation) ResetTotalCost() {
	m.totalCost = nil
	m.addtotalCost = nil
}

// SetCurrency sets the "currency" field.
func (m *AllocationCostMutation) SetCurrency(i int) {
	m.currency = &i
	m.addcurrency = nil
}

// Currency returns the value of the "currency" field in the mutation.
func (m *AllocationCostMutation) Currency() (r int, exists bool) {
	v := m.currency
	if v == nil {
		return
	}
	return *v, true
}

// OldCurrency returns the old "currency" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldCurrency(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCurrency is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCurrency requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCurrency: %w", err)
	}
	return oldValue.Currency, nil
}

// AddCurrency adds i to the "currency" field.
func (m *AllocationCostMutation) AddCurrency(i int) {
	if m.addcurrency != nil {
		*m.addcurrency += i
	} else {
		m.addcurrency = &i
	}
}

// AddedCurrency returns the value that was added to the "currency" field in this mutation.
func (m *AllocationCostMutation) AddedCurrency() (r int, exists bool) {
	v := m.addcurrency
	if v == nil {
		return
	}
	return *v, true
}

// ClearCurrency clears the value of the "currency" field.
func (m *AllocationCostMutation) ClearCurrency() {
	m.currency = nil
	m.addcurrency = nil
	m.clearedFields[allocationcost.FieldCurrency] = struct{}{}
}

// CurrencyCleared returns if the "currency" field was cleared in this mutation.
func (m *AllocationCostMutation) CurrencyCleared() bool {
	_, ok := m.clearedFields[allocationcost.FieldCurrency]
	return ok
}

// ResetCurrency resets all changes to the "currency" field.
func (m *AllocationCostMutation) ResetCurrency() {
	m.currency = nil
	m.addcurrency = nil
	delete(m.clearedFields, allocationcost.FieldCurrency)
}

// SetCpuCost sets the "cpuCost" field.
func (m *AllocationCostMutation) SetCpuCost(f float64) {
	m.cpuCost = &f
	m.addcpuCost = nil
}

// CpuCost returns the value of the "cpuCost" field in the mutation.
func (m *AllocationCostMutation) CpuCost() (r float64, exists bool) {
	v := m.cpuCost
	if v == nil {
		return
	}
	return *v, true
}

// OldCpuCost returns the old "cpuCost" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldCpuCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCpuCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCpuCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCpuCost: %w", err)
	}
	return oldValue.CpuCost, nil
}

// AddCpuCost adds f to the "cpuCost" field.
func (m *AllocationCostMutation) AddCpuCost(f float64) {
	if m.addcpuCost != nil {
		*m.addcpuCost += f
	} else {
		m.addcpuCost = &f
	}
}

// AddedCpuCost returns the value that was added to the "cpuCost" field in this mutation.
func (m *AllocationCostMutation) AddedCpuCost() (r float64, exists bool) {
	v := m.addcpuCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetCpuCost resets all changes to the "cpuCost" field.
func (m *AllocationCostMutation) ResetCpuCost() {
	m.cpuCost = nil
	m.addcpuCost = nil
}

// SetCpuCoreRequest sets the "cpuCoreRequest" field.
func (m *AllocationCostMutation) SetCpuCoreRequest(f float64) {
	m.cpuCoreRequest = &f
	m.addcpuCoreRequest = nil
}

// CpuCoreRequest returns the value of the "cpuCoreRequest" field in the mutation.
func (m *AllocationCostMutation) CpuCoreRequest() (r float64, exists bool) {
	v := m.cpuCoreRequest
	if v == nil {
		return
	}
	return *v, true
}

// OldCpuCoreRequest returns the old "cpuCoreRequest" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldCpuCoreRequest(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCpuCoreRequest is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCpuCoreRequest requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCpuCoreRequest: %w", err)
	}
	return oldValue.CpuCoreRequest, nil
}

// AddCpuCoreRequest adds f to the "cpuCoreRequest" field.
func (m *AllocationCostMutation) AddCpuCoreRequest(f float64) {
	if m.addcpuCoreRequest != nil {
		*m.addcpuCoreRequest += f
	} else {
		m.addcpuCoreRequest = &f
	}
}

// AddedCpuCoreRequest returns the value that was added to the "cpuCoreRequest" field in this mutation.
func (m *AllocationCostMutation) AddedCpuCoreRequest() (r float64, exists bool) {
	v := m.addcpuCoreRequest
	if v == nil {
		return
	}
	return *v, true
}

// ResetCpuCoreRequest resets all changes to the "cpuCoreRequest" field.
func (m *AllocationCostMutation) ResetCpuCoreRequest() {
	m.cpuCoreRequest = nil
	m.addcpuCoreRequest = nil
}

// SetGpuCost sets the "gpuCost" field.
func (m *AllocationCostMutation) SetGpuCost(f float64) {
	m.gpuCost = &f
	m.addgpuCost = nil
}

// GpuCost returns the value of the "gpuCost" field in the mutation.
func (m *AllocationCostMutation) GpuCost() (r float64, exists bool) {
	v := m.gpuCost
	if v == nil {
		return
	}
	return *v, true
}

// OldGpuCost returns the old "gpuCost" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldGpuCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldGpuCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldGpuCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldGpuCost: %w", err)
	}
	return oldValue.GpuCost, nil
}

// AddGpuCost adds f to the "gpuCost" field.
func (m *AllocationCostMutation) AddGpuCost(f float64) {
	if m.addgpuCost != nil {
		*m.addgpuCost += f
	} else {
		m.addgpuCost = &f
	}
}

// AddedGpuCost returns the value that was added to the "gpuCost" field in this mutation.
func (m *AllocationCostMutation) AddedGpuCost() (r float64, exists bool) {
	v := m.addgpuCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetGpuCost resets all changes to the "gpuCost" field.
func (m *AllocationCostMutation) ResetGpuCost() {
	m.gpuCost = nil
	m.addgpuCost = nil
}

// SetGpuCount sets the "gpuCount" field.
func (m *AllocationCostMutation) SetGpuCount(f float64) {
	m.gpuCount = &f
	m.addgpuCount = nil
}

// GpuCount returns the value of the "gpuCount" field in the mutation.
func (m *AllocationCostMutation) GpuCount() (r float64, exists bool) {
	v := m.gpuCount
	if v == nil {
		return
	}
	return *v, true
}

// OldGpuCount returns the old "gpuCount" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldGpuCount(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldGpuCount is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldGpuCount requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldGpuCount: %w", err)
	}
	return oldValue.GpuCount, nil
}

// AddGpuCount adds f to the "gpuCount" field.
func (m *AllocationCostMutation) AddGpuCount(f float64) {
	if m.addgpuCount != nil {
		*m.addgpuCount += f
	} else {
		m.addgpuCount = &f
	}
}

// AddedGpuCount returns the value that was added to the "gpuCount" field in this mutation.
func (m *AllocationCostMutation) AddedGpuCount() (r float64, exists bool) {
	v := m.addgpuCount
	if v == nil {
		return
	}
	return *v, true
}

// ResetGpuCount resets all changes to the "gpuCount" field.
func (m *AllocationCostMutation) ResetGpuCount() {
	m.gpuCount = nil
	m.addgpuCount = nil
}

// SetRamCost sets the "ramCost" field.
func (m *AllocationCostMutation) SetRamCost(f float64) {
	m.ramCost = &f
	m.addramCost = nil
}

// RamCost returns the value of the "ramCost" field in the mutation.
func (m *AllocationCostMutation) RamCost() (r float64, exists bool) {
	v := m.ramCost
	if v == nil {
		return
	}
	return *v, true
}

// OldRamCost returns the old "ramCost" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldRamCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRamCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRamCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRamCost: %w", err)
	}
	return oldValue.RamCost, nil
}

// AddRamCost adds f to the "ramCost" field.
func (m *AllocationCostMutation) AddRamCost(f float64) {
	if m.addramCost != nil {
		*m.addramCost += f
	} else {
		m.addramCost = &f
	}
}

// AddedRamCost returns the value that was added to the "ramCost" field in this mutation.
func (m *AllocationCostMutation) AddedRamCost() (r float64, exists bool) {
	v := m.addramCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetRamCost resets all changes to the "ramCost" field.
func (m *AllocationCostMutation) ResetRamCost() {
	m.ramCost = nil
	m.addramCost = nil
}

// SetRamByteRequest sets the "ramByteRequest" field.
func (m *AllocationCostMutation) SetRamByteRequest(f float64) {
	m.ramByteRequest = &f
	m.addramByteRequest = nil
}

// RamByteRequest returns the value of the "ramByteRequest" field in the mutation.
func (m *AllocationCostMutation) RamByteRequest() (r float64, exists bool) {
	v := m.ramByteRequest
	if v == nil {
		return
	}
	return *v, true
}

// OldRamByteRequest returns the old "ramByteRequest" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldRamByteRequest(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRamByteRequest is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRamByteRequest requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRamByteRequest: %w", err)
	}
	return oldValue.RamByteRequest, nil
}

// AddRamByteRequest adds f to the "ramByteRequest" field.
func (m *AllocationCostMutation) AddRamByteRequest(f float64) {
	if m.addramByteRequest != nil {
		*m.addramByteRequest += f
	} else {
		m.addramByteRequest = &f
	}
}

// AddedRamByteRequest returns the value that was added to the "ramByteRequest" field in this mutation.
func (m *AllocationCostMutation) AddedRamByteRequest() (r float64, exists bool) {
	v := m.addramByteRequest
	if v == nil {
		return
	}
	return *v, true
}

// ResetRamByteRequest resets all changes to the "ramByteRequest" field.
func (m *AllocationCostMutation) ResetRamByteRequest() {
	m.ramByteRequest = nil
	m.addramByteRequest = nil
}

// SetPvCost sets the "pvCost" field.
func (m *AllocationCostMutation) SetPvCost(f float64) {
	m.pvCost = &f
	m.addpvCost = nil
}

// PvCost returns the value of the "pvCost" field in the mutation.
func (m *AllocationCostMutation) PvCost() (r float64, exists bool) {
	v := m.pvCost
	if v == nil {
		return
	}
	return *v, true
}

// OldPvCost returns the old "pvCost" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldPvCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPvCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPvCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPvCost: %w", err)
	}
	return oldValue.PvCost, nil
}

// AddPvCost adds f to the "pvCost" field.
func (m *AllocationCostMutation) AddPvCost(f float64) {
	if m.addpvCost != nil {
		*m.addpvCost += f
	} else {
		m.addpvCost = &f
	}
}

// AddedPvCost returns the value that was added to the "pvCost" field in this mutation.
func (m *AllocationCostMutation) AddedPvCost() (r float64, exists bool) {
	v := m.addpvCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetPvCost resets all changes to the "pvCost" field.
func (m *AllocationCostMutation) ResetPvCost() {
	m.pvCost = nil
	m.addpvCost = nil
}

// SetPvBytes sets the "pvBytes" field.
func (m *AllocationCostMutation) SetPvBytes(f float64) {
	m.pvBytes = &f
	m.addpvBytes = nil
}

// PvBytes returns the value of the "pvBytes" field in the mutation.
func (m *AllocationCostMutation) PvBytes() (r float64, exists bool) {
	v := m.pvBytes
	if v == nil {
		return
	}
	return *v, true
}

// OldPvBytes returns the old "pvBytes" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldPvBytes(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPvBytes is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPvBytes requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPvBytes: %w", err)
	}
	return oldValue.PvBytes, nil
}

// AddPvBytes adds f to the "pvBytes" field.
func (m *AllocationCostMutation) AddPvBytes(f float64) {
	if m.addpvBytes != nil {
		*m.addpvBytes += f
	} else {
		m.addpvBytes = &f
	}
}

// AddedPvBytes returns the value that was added to the "pvBytes" field in this mutation.
func (m *AllocationCostMutation) AddedPvBytes() (r float64, exists bool) {
	v := m.addpvBytes
	if v == nil {
		return
	}
	return *v, true
}

// ResetPvBytes resets all changes to the "pvBytes" field.
func (m *AllocationCostMutation) ResetPvBytes() {
	m.pvBytes = nil
	m.addpvBytes = nil
}

// SetLoadBalancerCost sets the "loadBalancerCost" field.
func (m *AllocationCostMutation) SetLoadBalancerCost(f float64) {
	m.loadBalancerCost = &f
	m.addloadBalancerCost = nil
}

// LoadBalancerCost returns the value of the "loadBalancerCost" field in the mutation.
func (m *AllocationCostMutation) LoadBalancerCost() (r float64, exists bool) {
	v := m.loadBalancerCost
	if v == nil {
		return
	}
	return *v, true
}

// OldLoadBalancerCost returns the old "loadBalancerCost" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldLoadBalancerCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLoadBalancerCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLoadBalancerCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLoadBalancerCost: %w", err)
	}
	return oldValue.LoadBalancerCost, nil
}

// AddLoadBalancerCost adds f to the "loadBalancerCost" field.
func (m *AllocationCostMutation) AddLoadBalancerCost(f float64) {
	if m.addloadBalancerCost != nil {
		*m.addloadBalancerCost += f
	} else {
		m.addloadBalancerCost = &f
	}
}

// AddedLoadBalancerCost returns the value that was added to the "loadBalancerCost" field in this mutation.
func (m *AllocationCostMutation) AddedLoadBalancerCost() (r float64, exists bool) {
	v := m.addloadBalancerCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetLoadBalancerCost resets all changes to the "loadBalancerCost" field.
func (m *AllocationCostMutation) ResetLoadBalancerCost() {
	m.loadBalancerCost = nil
	m.addloadBalancerCost = nil
}

// SetCpuCoreUsageAverage sets the "cpuCoreUsageAverage" field.
func (m *AllocationCostMutation) SetCpuCoreUsageAverage(f float64) {
	m.cpuCoreUsageAverage = &f
	m.addcpuCoreUsageAverage = nil
}

// CpuCoreUsageAverage returns the value of the "cpuCoreUsageAverage" field in the mutation.
func (m *AllocationCostMutation) CpuCoreUsageAverage() (r float64, exists bool) {
	v := m.cpuCoreUsageAverage
	if v == nil {
		return
	}
	return *v, true
}

// OldCpuCoreUsageAverage returns the old "cpuCoreUsageAverage" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldCpuCoreUsageAverage(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCpuCoreUsageAverage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCpuCoreUsageAverage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCpuCoreUsageAverage: %w", err)
	}
	return oldValue.CpuCoreUsageAverage, nil
}

// AddCpuCoreUsageAverage adds f to the "cpuCoreUsageAverage" field.
func (m *AllocationCostMutation) AddCpuCoreUsageAverage(f float64) {
	if m.addcpuCoreUsageAverage != nil {
		*m.addcpuCoreUsageAverage += f
	} else {
		m.addcpuCoreUsageAverage = &f
	}
}

// AddedCpuCoreUsageAverage returns the value that was added to the "cpuCoreUsageAverage" field in this mutation.
func (m *AllocationCostMutation) AddedCpuCoreUsageAverage() (r float64, exists bool) {
	v := m.addcpuCoreUsageAverage
	if v == nil {
		return
	}
	return *v, true
}

// ResetCpuCoreUsageAverage resets all changes to the "cpuCoreUsageAverage" field.
func (m *AllocationCostMutation) ResetCpuCoreUsageAverage() {
	m.cpuCoreUsageAverage = nil
	m.addcpuCoreUsageAverage = nil
}

// SetCpuCoreUsageMax sets the "cpuCoreUsageMax" field.
func (m *AllocationCostMutation) SetCpuCoreUsageMax(f float64) {
	m.cpuCoreUsageMax = &f
	m.addcpuCoreUsageMax = nil
}

// CpuCoreUsageMax returns the value of the "cpuCoreUsageMax" field in the mutation.
func (m *AllocationCostMutation) CpuCoreUsageMax() (r float64, exists bool) {
	v := m.cpuCoreUsageMax
	if v == nil {
		return
	}
	return *v, true
}

// OldCpuCoreUsageMax returns the old "cpuCoreUsageMax" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldCpuCoreUsageMax(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCpuCoreUsageMax is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCpuCoreUsageMax requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCpuCoreUsageMax: %w", err)
	}
	return oldValue.CpuCoreUsageMax, nil
}

// AddCpuCoreUsageMax adds f to the "cpuCoreUsageMax" field.
func (m *AllocationCostMutation) AddCpuCoreUsageMax(f float64) {
	if m.addcpuCoreUsageMax != nil {
		*m.addcpuCoreUsageMax += f
	} else {
		m.addcpuCoreUsageMax = &f
	}
}

// AddedCpuCoreUsageMax returns the value that was added to the "cpuCoreUsageMax" field in this mutation.
func (m *AllocationCostMutation) AddedCpuCoreUsageMax() (r float64, exists bool) {
	v := m.addcpuCoreUsageMax
	if v == nil {
		return
	}
	return *v, true
}

// ResetCpuCoreUsageMax resets all changes to the "cpuCoreUsageMax" field.
func (m *AllocationCostMutation) ResetCpuCoreUsageMax() {
	m.cpuCoreUsageMax = nil
	m.addcpuCoreUsageMax = nil
}

// SetRamByteUsageAverage sets the "ramByteUsageAverage" field.
func (m *AllocationCostMutation) SetRamByteUsageAverage(f float64) {
	m.ramByteUsageAverage = &f
	m.addramByteUsageAverage = nil
}

// RamByteUsageAverage returns the value of the "ramByteUsageAverage" field in the mutation.
func (m *AllocationCostMutation) RamByteUsageAverage() (r float64, exists bool) {
	v := m.ramByteUsageAverage
	if v == nil {
		return
	}
	return *v, true
}

// OldRamByteUsageAverage returns the old "ramByteUsageAverage" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldRamByteUsageAverage(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRamByteUsageAverage is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRamByteUsageAverage requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRamByteUsageAverage: %w", err)
	}
	return oldValue.RamByteUsageAverage, nil
}

// AddRamByteUsageAverage adds f to the "ramByteUsageAverage" field.
func (m *AllocationCostMutation) AddRamByteUsageAverage(f float64) {
	if m.addramByteUsageAverage != nil {
		*m.addramByteUsageAverage += f
	} else {
		m.addramByteUsageAverage = &f
	}
}

// AddedRamByteUsageAverage returns the value that was added to the "ramByteUsageAverage" field in this mutation.
func (m *AllocationCostMutation) AddedRamByteUsageAverage() (r float64, exists bool) {
	v := m.addramByteUsageAverage
	if v == nil {
		return
	}
	return *v, true
}

// ResetRamByteUsageAverage resets all changes to the "ramByteUsageAverage" field.
func (m *AllocationCostMutation) ResetRamByteUsageAverage() {
	m.ramByteUsageAverage = nil
	m.addramByteUsageAverage = nil
}

// SetRamByteUsageMax sets the "ramByteUsageMax" field.
func (m *AllocationCostMutation) SetRamByteUsageMax(f float64) {
	m.ramByteUsageMax = &f
	m.addramByteUsageMax = nil
}

// RamByteUsageMax returns the value of the "ramByteUsageMax" field in the mutation.
func (m *AllocationCostMutation) RamByteUsageMax() (r float64, exists bool) {
	v := m.ramByteUsageMax
	if v == nil {
		return
	}
	return *v, true
}

// OldRamByteUsageMax returns the old "ramByteUsageMax" field's value of the AllocationCost entity.
// If the AllocationCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AllocationCostMutation) OldRamByteUsageMax(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRamByteUsageMax is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRamByteUsageMax requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRamByteUsageMax: %w", err)
	}
	return oldValue.RamByteUsageMax, nil
}

// AddRamByteUsageMax adds f to the "ramByteUsageMax" field.
func (m *AllocationCostMutation) AddRamByteUsageMax(f float64) {
	if m.addramByteUsageMax != nil {
		*m.addramByteUsageMax += f
	} else {
		m.addramByteUsageMax = &f
	}
}

// AddedRamByteUsageMax returns the value that was added to the "ramByteUsageMax" field in this mutation.
func (m *AllocationCostMutation) AddedRamByteUsageMax() (r float64, exists bool) {
	v := m.addramByteUsageMax
	if v == nil {
		return
	}
	return *v, true
}

// ResetRamByteUsageMax resets all changes to the "ramByteUsageMax" field.
func (m *AllocationCostMutation) ResetRamByteUsageMax() {
	m.ramByteUsageMax = nil
	m.addramByteUsageMax = nil
}

// ClearConnector clears the "connector" edge to the Connector entity.
func (m *AllocationCostMutation) ClearConnector() {
	m.clearedconnector = true
}

// ConnectorCleared reports if the "connector" edge to the Connector entity was cleared.
func (m *AllocationCostMutation) ConnectorCleared() bool {
	return m.clearedconnector
}

// ConnectorIDs returns the "connector" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ConnectorID instead. It exists only for internal usage by the builders.
func (m *AllocationCostMutation) ConnectorIDs() (ids []oid.ID) {
	if id := m.connector; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetConnector resets all changes to the "connector" edge.
func (m *AllocationCostMutation) ResetConnector() {
	m.connector = nil
	m.clearedconnector = false
}

// Where appends a list predicates to the AllocationCostMutation builder.
func (m *AllocationCostMutation) Where(ps ...predicate.AllocationCost) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the AllocationCostMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *AllocationCostMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.AllocationCost, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *AllocationCostMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *AllocationCostMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (AllocationCost).
func (m *AllocationCostMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *AllocationCostMutation) Fields() []string {
	fields := make([]string, 0, 30)
	if m.startTime != nil {
		fields = append(fields, allocationcost.FieldStartTime)
	}
	if m.endTime != nil {
		fields = append(fields, allocationcost.FieldEndTime)
	}
	if m.minutes != nil {
		fields = append(fields, allocationcost.FieldMinutes)
	}
	if m.connector != nil {
		fields = append(fields, allocationcost.FieldConnectorID)
	}
	if m.name != nil {
		fields = append(fields, allocationcost.FieldName)
	}
	if m.fingerprint != nil {
		fields = append(fields, allocationcost.FieldFingerprint)
	}
	if m.clusterName != nil {
		fields = append(fields, allocationcost.FieldClusterName)
	}
	if m.namespace != nil {
		fields = append(fields, allocationcost.FieldNamespace)
	}
	if m.node != nil {
		fields = append(fields, allocationcost.FieldNode)
	}
	if m.controller != nil {
		fields = append(fields, allocationcost.FieldController)
	}
	if m.controllerKind != nil {
		fields = append(fields, allocationcost.FieldControllerKind)
	}
	if m.pod != nil {
		fields = append(fields, allocationcost.FieldPod)
	}
	if m.container != nil {
		fields = append(fields, allocationcost.FieldContainer)
	}
	if m.pvs != nil {
		fields = append(fields, allocationcost.FieldPvs)
	}
	if m.labels != nil {
		fields = append(fields, allocationcost.FieldLabels)
	}
	if m.totalCost != nil {
		fields = append(fields, allocationcost.FieldTotalCost)
	}
	if m.currency != nil {
		fields = append(fields, allocationcost.FieldCurrency)
	}
	if m.cpuCost != nil {
		fields = append(fields, allocationcost.FieldCpuCost)
	}
	if m.cpuCoreRequest != nil {
		fields = append(fields, allocationcost.FieldCpuCoreRequest)
	}
	if m.gpuCost != nil {
		fields = append(fields, allocationcost.FieldGpuCost)
	}
	if m.gpuCount != nil {
		fields = append(fields, allocationcost.FieldGpuCount)
	}
	if m.ramCost != nil {
		fields = append(fields, allocationcost.FieldRamCost)
	}
	if m.ramByteRequest != nil {
		fields = append(fields, allocationcost.FieldRamByteRequest)
	}
	if m.pvCost != nil {
		fields = append(fields, allocationcost.FieldPvCost)
	}
	if m.pvBytes != nil {
		fields = append(fields, allocationcost.FieldPvBytes)
	}
	if m.loadBalancerCost != nil {
		fields = append(fields, allocationcost.FieldLoadBalancerCost)
	}
	if m.cpuCoreUsageAverage != nil {
		fields = append(fields, allocationcost.FieldCpuCoreUsageAverage)
	}
	if m.cpuCoreUsageMax != nil {
		fields = append(fields, allocationcost.FieldCpuCoreUsageMax)
	}
	if m.ramByteUsageAverage != nil {
		fields = append(fields, allocationcost.FieldRamByteUsageAverage)
	}
	if m.ramByteUsageMax != nil {
		fields = append(fields, allocationcost.FieldRamByteUsageMax)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *AllocationCostMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case allocationcost.FieldStartTime:
		return m.StartTime()
	case allocationcost.FieldEndTime:
		return m.EndTime()
	case allocationcost.FieldMinutes:
		return m.Minutes()
	case allocationcost.FieldConnectorID:
		return m.ConnectorID()
	case allocationcost.FieldName:
		return m.Name()
	case allocationcost.FieldFingerprint:
		return m.Fingerprint()
	case allocationcost.FieldClusterName:
		return m.ClusterName()
	case allocationcost.FieldNamespace:
		return m.Namespace()
	case allocationcost.FieldNode:
		return m.Node()
	case allocationcost.FieldController:
		return m.Controller()
	case allocationcost.FieldControllerKind:
		return m.ControllerKind()
	case allocationcost.FieldPod:
		return m.Pod()
	case allocationcost.FieldContainer:
		return m.Container()
	case allocationcost.FieldPvs:
		return m.Pvs()
	case allocationcost.FieldLabels:
		return m.Labels()
	case allocationcost.FieldTotalCost:
		return m.TotalCost()
	case allocationcost.FieldCurrency:
		return m.Currency()
	case allocationcost.FieldCpuCost:
		return m.CpuCost()
	case allocationcost.FieldCpuCoreRequest:
		return m.CpuCoreRequest()
	case allocationcost.FieldGpuCost:
		return m.GpuCost()
	case allocationcost.FieldGpuCount:
		return m.GpuCount()
	case allocationcost.FieldRamCost:
		return m.RamCost()
	case allocationcost.FieldRamByteRequest:
		return m.RamByteRequest()
	case allocationcost.FieldPvCost:
		return m.PvCost()
	case allocationcost.FieldPvBytes:
		return m.PvBytes()
	case allocationcost.FieldLoadBalancerCost:
		return m.LoadBalancerCost()
	case allocationcost.FieldCpuCoreUsageAverage:
		return m.CpuCoreUsageAverage()
	case allocationcost.FieldCpuCoreUsageMax:
		return m.CpuCoreUsageMax()
	case allocationcost.FieldRamByteUsageAverage:
		return m.RamByteUsageAverage()
	case allocationcost.FieldRamByteUsageMax:
		return m.RamByteUsageMax()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *AllocationCostMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case allocationcost.FieldStartTime:
		return m.OldStartTime(ctx)
	case allocationcost.FieldEndTime:
		return m.OldEndTime(ctx)
	case allocationcost.FieldMinutes:
		return m.OldMinutes(ctx)
	case allocationcost.FieldConnectorID:
		return m.OldConnectorID(ctx)
	case allocationcost.FieldName:
		return m.OldName(ctx)
	case allocationcost.FieldFingerprint:
		return m.OldFingerprint(ctx)
	case allocationcost.FieldClusterName:
		return m.OldClusterName(ctx)
	case allocationcost.FieldNamespace:
		return m.OldNamespace(ctx)
	case allocationcost.FieldNode:
		return m.OldNode(ctx)
	case allocationcost.FieldController:
		return m.OldController(ctx)
	case allocationcost.FieldControllerKind:
		return m.OldControllerKind(ctx)
	case allocationcost.FieldPod:
		return m.OldPod(ctx)
	case allocationcost.FieldContainer:
		return m.OldContainer(ctx)
	case allocationcost.FieldPvs:
		return m.OldPvs(ctx)
	case allocationcost.FieldLabels:
		return m.OldLabels(ctx)
	case allocationcost.FieldTotalCost:
		return m.OldTotalCost(ctx)
	case allocationcost.FieldCurrency:
		return m.OldCurrency(ctx)
	case allocationcost.FieldCpuCost:
		return m.OldCpuCost(ctx)
	case allocationcost.FieldCpuCoreRequest:
		return m.OldCpuCoreRequest(ctx)
	case allocationcost.FieldGpuCost:
		return m.OldGpuCost(ctx)
	case allocationcost.FieldGpuCount:
		return m.OldGpuCount(ctx)
	case allocationcost.FieldRamCost:
		return m.OldRamCost(ctx)
	case allocationcost.FieldRamByteRequest:
		return m.OldRamByteRequest(ctx)
	case allocationcost.FieldPvCost:
		return m.OldPvCost(ctx)
	case allocationcost.FieldPvBytes:
		return m.OldPvBytes(ctx)
	case allocationcost.FieldLoadBalancerCost:
		return m.OldLoadBalancerCost(ctx)
	case allocationcost.FieldCpuCoreUsageAverage:
		return m.OldCpuCoreUsageAverage(ctx)
	case allocationcost.FieldCpuCoreUsageMax:
		return m.OldCpuCoreUsageMax(ctx)
	case allocationcost.FieldRamByteUsageAverage:
		return m.OldRamByteUsageAverage(ctx)
	case allocationcost.FieldRamByteUsageMax:
		return m.OldRamByteUsageMax(ctx)
	}
	return nil, fmt.Errorf("unknown AllocationCost field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *AllocationCostMutation) SetField(name string, value ent.Value) error {
	switch name {
	case allocationcost.FieldStartTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStartTime(v)
		return nil
	case allocationcost.FieldEndTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEndTime(v)
		return nil
	case allocationcost.FieldMinutes:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMinutes(v)
		return nil
	case allocationcost.FieldConnectorID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetConnectorID(v)
		return nil
	case allocationcost.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case allocationcost.FieldFingerprint:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFingerprint(v)
		return nil
	case allocationcost.FieldClusterName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetClusterName(v)
		return nil
	case allocationcost.FieldNamespace:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetNamespace(v)
		return nil
	case allocationcost.FieldNode:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetNode(v)
		return nil
	case allocationcost.FieldController:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetController(v)
		return nil
	case allocationcost.FieldControllerKind:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetControllerKind(v)
		return nil
	case allocationcost.FieldPod:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPod(v)
		return nil
	case allocationcost.FieldContainer:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetContainer(v)
		return nil
	case allocationcost.FieldPvs:
		v, ok := value.(map[string]types.PVCost)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPvs(v)
		return nil
	case allocationcost.FieldLabels:
		v, ok := value.(map[string]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLabels(v)
		return nil
	case allocationcost.FieldTotalCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTotalCost(v)
		return nil
	case allocationcost.FieldCurrency:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCurrency(v)
		return nil
	case allocationcost.FieldCpuCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCpuCost(v)
		return nil
	case allocationcost.FieldCpuCoreRequest:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCpuCoreRequest(v)
		return nil
	case allocationcost.FieldGpuCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetGpuCost(v)
		return nil
	case allocationcost.FieldGpuCount:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetGpuCount(v)
		return nil
	case allocationcost.FieldRamCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRamCost(v)
		return nil
	case allocationcost.FieldRamByteRequest:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRamByteRequest(v)
		return nil
	case allocationcost.FieldPvCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPvCost(v)
		return nil
	case allocationcost.FieldPvBytes:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPvBytes(v)
		return nil
	case allocationcost.FieldLoadBalancerCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLoadBalancerCost(v)
		return nil
	case allocationcost.FieldCpuCoreUsageAverage:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCpuCoreUsageAverage(v)
		return nil
	case allocationcost.FieldCpuCoreUsageMax:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCpuCoreUsageMax(v)
		return nil
	case allocationcost.FieldRamByteUsageAverage:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRamByteUsageAverage(v)
		return nil
	case allocationcost.FieldRamByteUsageMax:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRamByteUsageMax(v)
		return nil
	}
	return fmt.Errorf("unknown AllocationCost field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *AllocationCostMutation) AddedFields() []string {
	var fields []string
	if m.addminutes != nil {
		fields = append(fields, allocationcost.FieldMinutes)
	}
	if m.addtotalCost != nil {
		fields = append(fields, allocationcost.FieldTotalCost)
	}
	if m.addcurrency != nil {
		fields = append(fields, allocationcost.FieldCurrency)
	}
	if m.addcpuCost != nil {
		fields = append(fields, allocationcost.FieldCpuCost)
	}
	if m.addcpuCoreRequest != nil {
		fields = append(fields, allocationcost.FieldCpuCoreRequest)
	}
	if m.addgpuCost != nil {
		fields = append(fields, allocationcost.FieldGpuCost)
	}
	if m.addgpuCount != nil {
		fields = append(fields, allocationcost.FieldGpuCount)
	}
	if m.addramCost != nil {
		fields = append(fields, allocationcost.FieldRamCost)
	}
	if m.addramByteRequest != nil {
		fields = append(fields, allocationcost.FieldRamByteRequest)
	}
	if m.addpvCost != nil {
		fields = append(fields, allocationcost.FieldPvCost)
	}
	if m.addpvBytes != nil {
		fields = append(fields, allocationcost.FieldPvBytes)
	}
	if m.addloadBalancerCost != nil {
		fields = append(fields, allocationcost.FieldLoadBalancerCost)
	}
	if m.addcpuCoreUsageAverage != nil {
		fields = append(fields, allocationcost.FieldCpuCoreUsageAverage)
	}
	if m.addcpuCoreUsageMax != nil {
		fields = append(fields, allocationcost.FieldCpuCoreUsageMax)
	}
	if m.addramByteUsageAverage != nil {
		fields = append(fields, allocationcost.FieldRamByteUsageAverage)
	}
	if m.addramByteUsageMax != nil {
		fields = append(fields, allocationcost.FieldRamByteUsageMax)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *AllocationCostMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case allocationcost.FieldMinutes:
		return m.AddedMinutes()
	case allocationcost.FieldTotalCost:
		return m.AddedTotalCost()
	case allocationcost.FieldCurrency:
		return m.AddedCurrency()
	case allocationcost.FieldCpuCost:
		return m.AddedCpuCost()
	case allocationcost.FieldCpuCoreRequest:
		return m.AddedCpuCoreRequest()
	case allocationcost.FieldGpuCost:
		return m.AddedGpuCost()
	case allocationcost.FieldGpuCount:
		return m.AddedGpuCount()
	case allocationcost.FieldRamCost:
		return m.AddedRamCost()
	case allocationcost.FieldRamByteRequest:
		return m.AddedRamByteRequest()
	case allocationcost.FieldPvCost:
		return m.AddedPvCost()
	case allocationcost.FieldPvBytes:
		return m.AddedPvBytes()
	case allocationcost.FieldLoadBalancerCost:
		return m.AddedLoadBalancerCost()
	case allocationcost.FieldCpuCoreUsageAverage:
		return m.AddedCpuCoreUsageAverage()
	case allocationcost.FieldCpuCoreUsageMax:
		return m.AddedCpuCoreUsageMax()
	case allocationcost.FieldRamByteUsageAverage:
		return m.AddedRamByteUsageAverage()
	case allocationcost.FieldRamByteUsageMax:
		return m.AddedRamByteUsageMax()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *AllocationCostMutation) AddField(name string, value ent.Value) error {
	switch name {
	case allocationcost.FieldMinutes:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddMinutes(v)
		return nil
	case allocationcost.FieldTotalCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddTotalCost(v)
		return nil
	case allocationcost.FieldCurrency:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddCurrency(v)
		return nil
	case allocationcost.FieldCpuCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddCpuCost(v)
		return nil
	case allocationcost.FieldCpuCoreRequest:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddCpuCoreRequest(v)
		return nil
	case allocationcost.FieldGpuCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddGpuCost(v)
		return nil
	case allocationcost.FieldGpuCount:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddGpuCount(v)
		return nil
	case allocationcost.FieldRamCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddRamCost(v)
		return nil
	case allocationcost.FieldRamByteRequest:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddRamByteRequest(v)
		return nil
	case allocationcost.FieldPvCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddPvCost(v)
		return nil
	case allocationcost.FieldPvBytes:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddPvBytes(v)
		return nil
	case allocationcost.FieldLoadBalancerCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLoadBalancerCost(v)
		return nil
	case allocationcost.FieldCpuCoreUsageAverage:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddCpuCoreUsageAverage(v)
		return nil
	case allocationcost.FieldCpuCoreUsageMax:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddCpuCoreUsageMax(v)
		return nil
	case allocationcost.FieldRamByteUsageAverage:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddRamByteUsageAverage(v)
		return nil
	case allocationcost.FieldRamByteUsageMax:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddRamByteUsageMax(v)
		return nil
	}
	return fmt.Errorf("unknown AllocationCost numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *AllocationCostMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(allocationcost.FieldClusterName) {
		fields = append(fields, allocationcost.FieldClusterName)
	}
	if m.FieldCleared(allocationcost.FieldNamespace) {
		fields = append(fields, allocationcost.FieldNamespace)
	}
	if m.FieldCleared(allocationcost.FieldNode) {
		fields = append(fields, allocationcost.FieldNode)
	}
	if m.FieldCleared(allocationcost.FieldController) {
		fields = append(fields, allocationcost.FieldController)
	}
	if m.FieldCleared(allocationcost.FieldControllerKind) {
		fields = append(fields, allocationcost.FieldControllerKind)
	}
	if m.FieldCleared(allocationcost.FieldPod) {
		fields = append(fields, allocationcost.FieldPod)
	}
	if m.FieldCleared(allocationcost.FieldContainer) {
		fields = append(fields, allocationcost.FieldContainer)
	}
	if m.FieldCleared(allocationcost.FieldCurrency) {
		fields = append(fields, allocationcost.FieldCurrency)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *AllocationCostMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *AllocationCostMutation) ClearField(name string) error {
	switch name {
	case allocationcost.FieldClusterName:
		m.ClearClusterName()
		return nil
	case allocationcost.FieldNamespace:
		m.ClearNamespace()
		return nil
	case allocationcost.FieldNode:
		m.ClearNode()
		return nil
	case allocationcost.FieldController:
		m.ClearController()
		return nil
	case allocationcost.FieldControllerKind:
		m.ClearControllerKind()
		return nil
	case allocationcost.FieldPod:
		m.ClearPod()
		return nil
	case allocationcost.FieldContainer:
		m.ClearContainer()
		return nil
	case allocationcost.FieldCurrency:
		m.ClearCurrency()
		return nil
	}
	return fmt.Errorf("unknown AllocationCost nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *AllocationCostMutation) ResetField(name string) error {
	switch name {
	case allocationcost.FieldStartTime:
		m.ResetStartTime()
		return nil
	case allocationcost.FieldEndTime:
		m.ResetEndTime()
		return nil
	case allocationcost.FieldMinutes:
		m.ResetMinutes()
		return nil
	case allocationcost.FieldConnectorID:
		m.ResetConnectorID()
		return nil
	case allocationcost.FieldName:
		m.ResetName()
		return nil
	case allocationcost.FieldFingerprint:
		m.ResetFingerprint()
		return nil
	case allocationcost.FieldClusterName:
		m.ResetClusterName()
		return nil
	case allocationcost.FieldNamespace:
		m.ResetNamespace()
		return nil
	case allocationcost.FieldNode:
		m.ResetNode()
		return nil
	case allocationcost.FieldController:
		m.ResetController()
		return nil
	case allocationcost.FieldControllerKind:
		m.ResetControllerKind()
		return nil
	case allocationcost.FieldPod:
		m.ResetPod()
		return nil
	case allocationcost.FieldContainer:
		m.ResetContainer()
		return nil
	case allocationcost.FieldPvs:
		m.ResetPvs()
		return nil
	case allocationcost.FieldLabels:
		m.ResetLabels()
		return nil
	case allocationcost.FieldTotalCost:
		m.ResetTotalCost()
		return nil
	case allocationcost.FieldCurrency:
		m.ResetCurrency()
		return nil
	case allocationcost.FieldCpuCost:
		m.ResetCpuCost()
		return nil
	case allocationcost.FieldCpuCoreRequest:
		m.ResetCpuCoreRequest()
		return nil
	case allocationcost.FieldGpuCost:
		m.ResetGpuCost()
		return nil
	case allocationcost.FieldGpuCount:
		m.ResetGpuCount()
		return nil
	case allocationcost.FieldRamCost:
		m.ResetRamCost()
		return nil
	case allocationcost.FieldRamByteRequest:
		m.ResetRamByteRequest()
		return nil
	case allocationcost.FieldPvCost:
		m.ResetPvCost()
		return nil
	case allocationcost.FieldPvBytes:
		m.ResetPvBytes()
		return nil
	case allocationcost.FieldLoadBalancerCost:
		m.ResetLoadBalancerCost()
		return nil
	case allocationcost.FieldCpuCoreUsageAverage:
		m.ResetCpuCoreUsageAverage()
		return nil
	case allocationcost.FieldCpuCoreUsageMax:
		m.ResetCpuCoreUsageMax()
		return nil
	case allocationcost.FieldRamByteUsageAverage:
		m.ResetRamByteUsageAverage()
		return nil
	case allocationcost.FieldRamByteUsageMax:
		m.ResetRamByteUsageMax()
		return nil
	}
	return fmt.Errorf("unknown AllocationCost field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *AllocationCostMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.connector != nil {
		edges = append(edges, allocationcost.EdgeConnector)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *AllocationCostMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case allocationcost.EdgeConnector:
		if id := m.connector; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *AllocationCostMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *AllocationCostMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *AllocationCostMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedconnector {
		edges = append(edges, allocationcost.EdgeConnector)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *AllocationCostMutation) EdgeCleared(name string) bool {
	switch name {
	case allocationcost.EdgeConnector:
		return m.clearedconnector
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *AllocationCostMutation) ClearEdge(name string) error {
	switch name {
	case allocationcost.EdgeConnector:
		m.ClearConnector()
		return nil
	}
	return fmt.Errorf("unknown AllocationCost unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *AllocationCostMutation) ResetEdge(name string) error {
	switch name {
	case allocationcost.EdgeConnector:
		m.ResetConnector()
		return nil
	}
	return fmt.Errorf("unknown AllocationCost edge %s", name)
}

// ApplicationMutation represents an operation that mutates the Application nodes in the graph.
type ApplicationMutation struct {
	config
	op               Op
	typ              string
	id               *oid.ID
	name             *string
	description      *string
	labels           *map[string]string
	createTime       *time.Time
	updateTime       *time.Time
	variables        *[]types.ApplicationVariable
	appendvariables  []types.ApplicationVariable
	clearedFields    map[string]struct{}
	project          *oid.ID
	clearedproject   bool
	instances        map[oid.ID]struct{}
	removedinstances map[oid.ID]struct{}
	clearedinstances bool
	done             bool
	oldValue         func(context.Context) (*Application, error)
	predicates       []predicate.Application
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
func (m *ApplicationMutation) SetProjectID(o oid.ID) {
	m.project = &o
}

// ProjectID returns the value of the "projectID" field in the mutation.
func (m *ApplicationMutation) ProjectID() (r oid.ID, exists bool) {
	v := m.project
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
	m.project = nil
}

// SetVariables sets the "variables" field.
func (m *ApplicationMutation) SetVariables(tv []types.ApplicationVariable) {
	m.variables = &tv
	m.appendvariables = nil
}

// Variables returns the value of the "variables" field in the mutation.
func (m *ApplicationMutation) Variables() (r []types.ApplicationVariable, exists bool) {
	v := m.variables
	if v == nil {
		return
	}
	return *v, true
}

// OldVariables returns the old "variables" field's value of the Application entity.
// If the Application object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationMutation) OldVariables(ctx context.Context) (v []types.ApplicationVariable, err error) {
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

// AppendVariables adds tv to the "variables" field.
func (m *ApplicationMutation) AppendVariables(tv []types.ApplicationVariable) {
	m.appendvariables = append(m.appendvariables, tv...)
}

// AppendedVariables returns the list of values that were appended to the "variables" field in this mutation.
func (m *ApplicationMutation) AppendedVariables() ([]types.ApplicationVariable, bool) {
	if len(m.appendvariables) == 0 {
		return nil, false
	}
	return m.appendvariables, true
}

// ClearVariables clears the value of the "variables" field.
func (m *ApplicationMutation) ClearVariables() {
	m.variables = nil
	m.appendvariables = nil
	m.clearedFields[application.FieldVariables] = struct{}{}
}

// VariablesCleared returns if the "variables" field was cleared in this mutation.
func (m *ApplicationMutation) VariablesCleared() bool {
	_, ok := m.clearedFields[application.FieldVariables]
	return ok
}

// ResetVariables resets all changes to the "variables" field.
func (m *ApplicationMutation) ResetVariables() {
	m.variables = nil
	m.appendvariables = nil
	delete(m.clearedFields, application.FieldVariables)
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
func (m *ApplicationMutation) ProjectIDs() (ids []oid.ID) {
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

// AddInstanceIDs adds the "instances" edge to the ApplicationInstance entity by ids.
func (m *ApplicationMutation) AddInstanceIDs(ids ...oid.ID) {
	if m.instances == nil {
		m.instances = make(map[oid.ID]struct{})
	}
	for i := range ids {
		m.instances[ids[i]] = struct{}{}
	}
}

// ClearInstances clears the "instances" edge to the ApplicationInstance entity.
func (m *ApplicationMutation) ClearInstances() {
	m.clearedinstances = true
}

// InstancesCleared reports if the "instances" edge to the ApplicationInstance entity was cleared.
func (m *ApplicationMutation) InstancesCleared() bool {
	return m.clearedinstances
}

// RemoveInstanceIDs removes the "instances" edge to the ApplicationInstance entity by IDs.
func (m *ApplicationMutation) RemoveInstanceIDs(ids ...oid.ID) {
	if m.removedinstances == nil {
		m.removedinstances = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.instances, ids[i])
		m.removedinstances[ids[i]] = struct{}{}
	}
}

// RemovedInstances returns the removed IDs of the "instances" edge to the ApplicationInstance entity.
func (m *ApplicationMutation) RemovedInstancesIDs() (ids []oid.ID) {
	for id := range m.removedinstances {
		ids = append(ids, id)
	}
	return
}

// InstancesIDs returns the "instances" edge IDs in the mutation.
func (m *ApplicationMutation) InstancesIDs() (ids []oid.ID) {
	for id := range m.instances {
		ids = append(ids, id)
	}
	return
}

// ResetInstances resets all changes to the "instances" edge.
func (m *ApplicationMutation) ResetInstances() {
	m.instances = nil
	m.clearedinstances = false
	m.removedinstances = nil
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
	if m.variables != nil {
		fields = append(fields, application.FieldVariables)
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
	case application.FieldVariables:
		return m.Variables()
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
	case application.FieldVariables:
		return m.OldVariables(ctx)
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
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetProjectID(v)
		return nil
	case application.FieldVariables:
		v, ok := value.([]types.ApplicationVariable)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVariables(v)
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
	if m.FieldCleared(application.FieldVariables) {
		fields = append(fields, application.FieldVariables)
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
	case application.FieldVariables:
		m.ClearVariables()
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
	case application.FieldVariables:
		m.ResetVariables()
		return nil
	}
	return fmt.Errorf("unknown Application field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.project != nil {
		edges = append(edges, application.EdgeProject)
	}
	if m.instances != nil {
		edges = append(edges, application.EdgeInstances)
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
	case application.EdgeInstances:
		ids := make([]ent.Value, 0, len(m.instances))
		for id := range m.instances {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	if m.removedinstances != nil {
		edges = append(edges, application.EdgeInstances)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case application.EdgeInstances:
		ids := make([]ent.Value, 0, len(m.removedinstances))
		for id := range m.removedinstances {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedproject {
		edges = append(edges, application.EdgeProject)
	}
	if m.clearedinstances {
		edges = append(edges, application.EdgeInstances)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationMutation) EdgeCleared(name string) bool {
	switch name {
	case application.EdgeProject:
		return m.clearedproject
	case application.EdgeInstances:
		return m.clearedinstances
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
	case application.EdgeInstances:
		m.ResetInstances()
		return nil
	}
	return fmt.Errorf("unknown Application edge %s", name)
}

// ApplicationInstanceMutation represents an operation that mutates the ApplicationInstance nodes in the graph.
type ApplicationInstanceMutation struct {
	config
	op                 Op
	typ                string
	id                 *oid.ID
	status             *string
	statusMessage      *string
	createTime         *time.Time
	updateTime         *time.Time
	name               *string
	variables          *map[string]interface{}
	clearedFields      map[string]struct{}
	application        *oid.ID
	clearedapplication bool
	environment        *oid.ID
	clearedenvironment bool
	revisions          map[oid.ID]struct{}
	removedrevisions   map[oid.ID]struct{}
	clearedrevisions   bool
	resources          map[oid.ID]struct{}
	removedresources   map[oid.ID]struct{}
	clearedresources   bool
	done               bool
	oldValue           func(context.Context) (*ApplicationInstance, error)
	predicates         []predicate.ApplicationInstance
}

var _ ent.Mutation = (*ApplicationInstanceMutation)(nil)

// applicationInstanceOption allows management of the mutation configuration using functional options.
type applicationInstanceOption func(*ApplicationInstanceMutation)

// newApplicationInstanceMutation creates new mutation for the ApplicationInstance entity.
func newApplicationInstanceMutation(c config, op Op, opts ...applicationInstanceOption) *ApplicationInstanceMutation {
	m := &ApplicationInstanceMutation{
		config:        c,
		op:            op,
		typ:           TypeApplicationInstance,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withApplicationInstanceID sets the ID field of the mutation.
func withApplicationInstanceID(id oid.ID) applicationInstanceOption {
	return func(m *ApplicationInstanceMutation) {
		var (
			err   error
			once  sync.Once
			value *ApplicationInstance
		)
		m.oldValue = func(ctx context.Context) (*ApplicationInstance, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().ApplicationInstance.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withApplicationInstance sets the old ApplicationInstance of the mutation.
func withApplicationInstance(node *ApplicationInstance) applicationInstanceOption {
	return func(m *ApplicationInstanceMutation) {
		m.oldValue = func(context.Context) (*ApplicationInstance, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ApplicationInstanceMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ApplicationInstanceMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of ApplicationInstance entities.
func (m *ApplicationInstanceMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ApplicationInstanceMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ApplicationInstanceMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().ApplicationInstance.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetStatus sets the "status" field.
func (m *ApplicationInstanceMutation) SetStatus(s string) {
	m.status = &s
}

// Status returns the value of the "status" field in the mutation.
func (m *ApplicationInstanceMutation) Status() (r string, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the ApplicationInstance entity.
// If the ApplicationInstance object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationInstanceMutation) OldStatus(ctx context.Context) (v string, err error) {
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
func (m *ApplicationInstanceMutation) ClearStatus() {
	m.status = nil
	m.clearedFields[applicationinstance.FieldStatus] = struct{}{}
}

// StatusCleared returns if the "status" field was cleared in this mutation.
func (m *ApplicationInstanceMutation) StatusCleared() bool {
	_, ok := m.clearedFields[applicationinstance.FieldStatus]
	return ok
}

// ResetStatus resets all changes to the "status" field.
func (m *ApplicationInstanceMutation) ResetStatus() {
	m.status = nil
	delete(m.clearedFields, applicationinstance.FieldStatus)
}

// SetStatusMessage sets the "statusMessage" field.
func (m *ApplicationInstanceMutation) SetStatusMessage(s string) {
	m.statusMessage = &s
}

// StatusMessage returns the value of the "statusMessage" field in the mutation.
func (m *ApplicationInstanceMutation) StatusMessage() (r string, exists bool) {
	v := m.statusMessage
	if v == nil {
		return
	}
	return *v, true
}

// OldStatusMessage returns the old "statusMessage" field's value of the ApplicationInstance entity.
// If the ApplicationInstance object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationInstanceMutation) OldStatusMessage(ctx context.Context) (v string, err error) {
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
func (m *ApplicationInstanceMutation) ClearStatusMessage() {
	m.statusMessage = nil
	m.clearedFields[applicationinstance.FieldStatusMessage] = struct{}{}
}

// StatusMessageCleared returns if the "statusMessage" field was cleared in this mutation.
func (m *ApplicationInstanceMutation) StatusMessageCleared() bool {
	_, ok := m.clearedFields[applicationinstance.FieldStatusMessage]
	return ok
}

// ResetStatusMessage resets all changes to the "statusMessage" field.
func (m *ApplicationInstanceMutation) ResetStatusMessage() {
	m.statusMessage = nil
	delete(m.clearedFields, applicationinstance.FieldStatusMessage)
}

// SetCreateTime sets the "createTime" field.
func (m *ApplicationInstanceMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *ApplicationInstanceMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the ApplicationInstance entity.
// If the ApplicationInstance object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationInstanceMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *ApplicationInstanceMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *ApplicationInstanceMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *ApplicationInstanceMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the ApplicationInstance entity.
// If the ApplicationInstance object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationInstanceMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *ApplicationInstanceMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetApplicationID sets the "applicationID" field.
func (m *ApplicationInstanceMutation) SetApplicationID(o oid.ID) {
	m.application = &o
}

// ApplicationID returns the value of the "applicationID" field in the mutation.
func (m *ApplicationInstanceMutation) ApplicationID() (r oid.ID, exists bool) {
	v := m.application
	if v == nil {
		return
	}
	return *v, true
}

// OldApplicationID returns the old "applicationID" field's value of the ApplicationInstance entity.
// If the ApplicationInstance object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationInstanceMutation) OldApplicationID(ctx context.Context) (v oid.ID, err error) {
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
func (m *ApplicationInstanceMutation) ResetApplicationID() {
	m.application = nil
}

// SetEnvironmentID sets the "environmentID" field.
func (m *ApplicationInstanceMutation) SetEnvironmentID(o oid.ID) {
	m.environment = &o
}

// EnvironmentID returns the value of the "environmentID" field in the mutation.
func (m *ApplicationInstanceMutation) EnvironmentID() (r oid.ID, exists bool) {
	v := m.environment
	if v == nil {
		return
	}
	return *v, true
}

// OldEnvironmentID returns the old "environmentID" field's value of the ApplicationInstance entity.
// If the ApplicationInstance object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationInstanceMutation) OldEnvironmentID(ctx context.Context) (v oid.ID, err error) {
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
func (m *ApplicationInstanceMutation) ResetEnvironmentID() {
	m.environment = nil
}

// SetName sets the "name" field.
func (m *ApplicationInstanceMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *ApplicationInstanceMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the ApplicationInstance entity.
// If the ApplicationInstance object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationInstanceMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *ApplicationInstanceMutation) ResetName() {
	m.name = nil
}

// SetVariables sets the "variables" field.
func (m *ApplicationInstanceMutation) SetVariables(value map[string]interface{}) {
	m.variables = &value
}

// Variables returns the value of the "variables" field in the mutation.
func (m *ApplicationInstanceMutation) Variables() (r map[string]interface{}, exists bool) {
	v := m.variables
	if v == nil {
		return
	}
	return *v, true
}

// OldVariables returns the old "variables" field's value of the ApplicationInstance entity.
// If the ApplicationInstance object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationInstanceMutation) OldVariables(ctx context.Context) (v map[string]interface{}, err error) {
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
func (m *ApplicationInstanceMutation) ClearVariables() {
	m.variables = nil
	m.clearedFields[applicationinstance.FieldVariables] = struct{}{}
}

// VariablesCleared returns if the "variables" field was cleared in this mutation.
func (m *ApplicationInstanceMutation) VariablesCleared() bool {
	_, ok := m.clearedFields[applicationinstance.FieldVariables]
	return ok
}

// ResetVariables resets all changes to the "variables" field.
func (m *ApplicationInstanceMutation) ResetVariables() {
	m.variables = nil
	delete(m.clearedFields, applicationinstance.FieldVariables)
}

// ClearApplication clears the "application" edge to the Application entity.
func (m *ApplicationInstanceMutation) ClearApplication() {
	m.clearedapplication = true
}

// ApplicationCleared reports if the "application" edge to the Application entity was cleared.
func (m *ApplicationInstanceMutation) ApplicationCleared() bool {
	return m.clearedapplication
}

// ApplicationIDs returns the "application" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ApplicationID instead. It exists only for internal usage by the builders.
func (m *ApplicationInstanceMutation) ApplicationIDs() (ids []oid.ID) {
	if id := m.application; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetApplication resets all changes to the "application" edge.
func (m *ApplicationInstanceMutation) ResetApplication() {
	m.application = nil
	m.clearedapplication = false
}

// ClearEnvironment clears the "environment" edge to the Environment entity.
func (m *ApplicationInstanceMutation) ClearEnvironment() {
	m.clearedenvironment = true
}

// EnvironmentCleared reports if the "environment" edge to the Environment entity was cleared.
func (m *ApplicationInstanceMutation) EnvironmentCleared() bool {
	return m.clearedenvironment
}

// EnvironmentIDs returns the "environment" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// EnvironmentID instead. It exists only for internal usage by the builders.
func (m *ApplicationInstanceMutation) EnvironmentIDs() (ids []oid.ID) {
	if id := m.environment; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEnvironment resets all changes to the "environment" edge.
func (m *ApplicationInstanceMutation) ResetEnvironment() {
	m.environment = nil
	m.clearedenvironment = false
}

// AddRevisionIDs adds the "revisions" edge to the ApplicationRevision entity by ids.
func (m *ApplicationInstanceMutation) AddRevisionIDs(ids ...oid.ID) {
	if m.revisions == nil {
		m.revisions = make(map[oid.ID]struct{})
	}
	for i := range ids {
		m.revisions[ids[i]] = struct{}{}
	}
}

// ClearRevisions clears the "revisions" edge to the ApplicationRevision entity.
func (m *ApplicationInstanceMutation) ClearRevisions() {
	m.clearedrevisions = true
}

// RevisionsCleared reports if the "revisions" edge to the ApplicationRevision entity was cleared.
func (m *ApplicationInstanceMutation) RevisionsCleared() bool {
	return m.clearedrevisions
}

// RemoveRevisionIDs removes the "revisions" edge to the ApplicationRevision entity by IDs.
func (m *ApplicationInstanceMutation) RemoveRevisionIDs(ids ...oid.ID) {
	if m.removedrevisions == nil {
		m.removedrevisions = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.revisions, ids[i])
		m.removedrevisions[ids[i]] = struct{}{}
	}
}

// RemovedRevisions returns the removed IDs of the "revisions" edge to the ApplicationRevision entity.
func (m *ApplicationInstanceMutation) RemovedRevisionsIDs() (ids []oid.ID) {
	for id := range m.removedrevisions {
		ids = append(ids, id)
	}
	return
}

// RevisionsIDs returns the "revisions" edge IDs in the mutation.
func (m *ApplicationInstanceMutation) RevisionsIDs() (ids []oid.ID) {
	for id := range m.revisions {
		ids = append(ids, id)
	}
	return
}

// ResetRevisions resets all changes to the "revisions" edge.
func (m *ApplicationInstanceMutation) ResetRevisions() {
	m.revisions = nil
	m.clearedrevisions = false
	m.removedrevisions = nil
}

// AddResourceIDs adds the "resources" edge to the ApplicationResource entity by ids.
func (m *ApplicationInstanceMutation) AddResourceIDs(ids ...oid.ID) {
	if m.resources == nil {
		m.resources = make(map[oid.ID]struct{})
	}
	for i := range ids {
		m.resources[ids[i]] = struct{}{}
	}
}

// ClearResources clears the "resources" edge to the ApplicationResource entity.
func (m *ApplicationInstanceMutation) ClearResources() {
	m.clearedresources = true
}

// ResourcesCleared reports if the "resources" edge to the ApplicationResource entity was cleared.
func (m *ApplicationInstanceMutation) ResourcesCleared() bool {
	return m.clearedresources
}

// RemoveResourceIDs removes the "resources" edge to the ApplicationResource entity by IDs.
func (m *ApplicationInstanceMutation) RemoveResourceIDs(ids ...oid.ID) {
	if m.removedresources == nil {
		m.removedresources = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.resources, ids[i])
		m.removedresources[ids[i]] = struct{}{}
	}
}

// RemovedResources returns the removed IDs of the "resources" edge to the ApplicationResource entity.
func (m *ApplicationInstanceMutation) RemovedResourcesIDs() (ids []oid.ID) {
	for id := range m.removedresources {
		ids = append(ids, id)
	}
	return
}

// ResourcesIDs returns the "resources" edge IDs in the mutation.
func (m *ApplicationInstanceMutation) ResourcesIDs() (ids []oid.ID) {
	for id := range m.resources {
		ids = append(ids, id)
	}
	return
}

// ResetResources resets all changes to the "resources" edge.
func (m *ApplicationInstanceMutation) ResetResources() {
	m.resources = nil
	m.clearedresources = false
	m.removedresources = nil
}

// Where appends a list predicates to the ApplicationInstanceMutation builder.
func (m *ApplicationInstanceMutation) Where(ps ...predicate.ApplicationInstance) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ApplicationInstanceMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ApplicationInstanceMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.ApplicationInstance, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ApplicationInstanceMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ApplicationInstanceMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (ApplicationInstance).
func (m *ApplicationInstanceMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ApplicationInstanceMutation) Fields() []string {
	fields := make([]string, 0, 8)
	if m.status != nil {
		fields = append(fields, applicationinstance.FieldStatus)
	}
	if m.statusMessage != nil {
		fields = append(fields, applicationinstance.FieldStatusMessage)
	}
	if m.createTime != nil {
		fields = append(fields, applicationinstance.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, applicationinstance.FieldUpdateTime)
	}
	if m.application != nil {
		fields = append(fields, applicationinstance.FieldApplicationID)
	}
	if m.environment != nil {
		fields = append(fields, applicationinstance.FieldEnvironmentID)
	}
	if m.name != nil {
		fields = append(fields, applicationinstance.FieldName)
	}
	if m.variables != nil {
		fields = append(fields, applicationinstance.FieldVariables)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ApplicationInstanceMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case applicationinstance.FieldStatus:
		return m.Status()
	case applicationinstance.FieldStatusMessage:
		return m.StatusMessage()
	case applicationinstance.FieldCreateTime:
		return m.CreateTime()
	case applicationinstance.FieldUpdateTime:
		return m.UpdateTime()
	case applicationinstance.FieldApplicationID:
		return m.ApplicationID()
	case applicationinstance.FieldEnvironmentID:
		return m.EnvironmentID()
	case applicationinstance.FieldName:
		return m.Name()
	case applicationinstance.FieldVariables:
		return m.Variables()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ApplicationInstanceMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case applicationinstance.FieldStatus:
		return m.OldStatus(ctx)
	case applicationinstance.FieldStatusMessage:
		return m.OldStatusMessage(ctx)
	case applicationinstance.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case applicationinstance.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case applicationinstance.FieldApplicationID:
		return m.OldApplicationID(ctx)
	case applicationinstance.FieldEnvironmentID:
		return m.OldEnvironmentID(ctx)
	case applicationinstance.FieldName:
		return m.OldName(ctx)
	case applicationinstance.FieldVariables:
		return m.OldVariables(ctx)
	}
	return nil, fmt.Errorf("unknown ApplicationInstance field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationInstanceMutation) SetField(name string, value ent.Value) error {
	switch name {
	case applicationinstance.FieldStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case applicationinstance.FieldStatusMessage:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatusMessage(v)
		return nil
	case applicationinstance.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case applicationinstance.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case applicationinstance.FieldApplicationID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetApplicationID(v)
		return nil
	case applicationinstance.FieldEnvironmentID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnvironmentID(v)
		return nil
	case applicationinstance.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case applicationinstance.FieldVariables:
		v, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVariables(v)
		return nil
	}
	return fmt.Errorf("unknown ApplicationInstance field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ApplicationInstanceMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ApplicationInstanceMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationInstanceMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ApplicationInstance numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ApplicationInstanceMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(applicationinstance.FieldStatus) {
		fields = append(fields, applicationinstance.FieldStatus)
	}
	if m.FieldCleared(applicationinstance.FieldStatusMessage) {
		fields = append(fields, applicationinstance.FieldStatusMessage)
	}
	if m.FieldCleared(applicationinstance.FieldVariables) {
		fields = append(fields, applicationinstance.FieldVariables)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ApplicationInstanceMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ApplicationInstanceMutation) ClearField(name string) error {
	switch name {
	case applicationinstance.FieldStatus:
		m.ClearStatus()
		return nil
	case applicationinstance.FieldStatusMessage:
		m.ClearStatusMessage()
		return nil
	case applicationinstance.FieldVariables:
		m.ClearVariables()
		return nil
	}
	return fmt.Errorf("unknown ApplicationInstance nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ApplicationInstanceMutation) ResetField(name string) error {
	switch name {
	case applicationinstance.FieldStatus:
		m.ResetStatus()
		return nil
	case applicationinstance.FieldStatusMessage:
		m.ResetStatusMessage()
		return nil
	case applicationinstance.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case applicationinstance.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case applicationinstance.FieldApplicationID:
		m.ResetApplicationID()
		return nil
	case applicationinstance.FieldEnvironmentID:
		m.ResetEnvironmentID()
		return nil
	case applicationinstance.FieldName:
		m.ResetName()
		return nil
	case applicationinstance.FieldVariables:
		m.ResetVariables()
		return nil
	}
	return fmt.Errorf("unknown ApplicationInstance field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationInstanceMutation) AddedEdges() []string {
	edges := make([]string, 0, 4)
	if m.application != nil {
		edges = append(edges, applicationinstance.EdgeApplication)
	}
	if m.environment != nil {
		edges = append(edges, applicationinstance.EdgeEnvironment)
	}
	if m.revisions != nil {
		edges = append(edges, applicationinstance.EdgeRevisions)
	}
	if m.resources != nil {
		edges = append(edges, applicationinstance.EdgeResources)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ApplicationInstanceMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case applicationinstance.EdgeApplication:
		if id := m.application; id != nil {
			return []ent.Value{*id}
		}
	case applicationinstance.EdgeEnvironment:
		if id := m.environment; id != nil {
			return []ent.Value{*id}
		}
	case applicationinstance.EdgeRevisions:
		ids := make([]ent.Value, 0, len(m.revisions))
		for id := range m.revisions {
			ids = append(ids, id)
		}
		return ids
	case applicationinstance.EdgeResources:
		ids := make([]ent.Value, 0, len(m.resources))
		for id := range m.resources {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationInstanceMutation) RemovedEdges() []string {
	edges := make([]string, 0, 4)
	if m.removedrevisions != nil {
		edges = append(edges, applicationinstance.EdgeRevisions)
	}
	if m.removedresources != nil {
		edges = append(edges, applicationinstance.EdgeResources)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationInstanceMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case applicationinstance.EdgeRevisions:
		ids := make([]ent.Value, 0, len(m.removedrevisions))
		for id := range m.removedrevisions {
			ids = append(ids, id)
		}
		return ids
	case applicationinstance.EdgeResources:
		ids := make([]ent.Value, 0, len(m.removedresources))
		for id := range m.removedresources {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationInstanceMutation) ClearedEdges() []string {
	edges := make([]string, 0, 4)
	if m.clearedapplication {
		edges = append(edges, applicationinstance.EdgeApplication)
	}
	if m.clearedenvironment {
		edges = append(edges, applicationinstance.EdgeEnvironment)
	}
	if m.clearedrevisions {
		edges = append(edges, applicationinstance.EdgeRevisions)
	}
	if m.clearedresources {
		edges = append(edges, applicationinstance.EdgeResources)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationInstanceMutation) EdgeCleared(name string) bool {
	switch name {
	case applicationinstance.EdgeApplication:
		return m.clearedapplication
	case applicationinstance.EdgeEnvironment:
		return m.clearedenvironment
	case applicationinstance.EdgeRevisions:
		return m.clearedrevisions
	case applicationinstance.EdgeResources:
		return m.clearedresources
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationInstanceMutation) ClearEdge(name string) error {
	switch name {
	case applicationinstance.EdgeApplication:
		m.ClearApplication()
		return nil
	case applicationinstance.EdgeEnvironment:
		m.ClearEnvironment()
		return nil
	}
	return fmt.Errorf("unknown ApplicationInstance unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ApplicationInstanceMutation) ResetEdge(name string) error {
	switch name {
	case applicationinstance.EdgeApplication:
		m.ResetApplication()
		return nil
	case applicationinstance.EdgeEnvironment:
		m.ResetEnvironment()
		return nil
	case applicationinstance.EdgeRevisions:
		m.ResetRevisions()
		return nil
	case applicationinstance.EdgeResources:
		m.ResetResources()
		return nil
	}
	return fmt.Errorf("unknown ApplicationInstance edge %s", name)
}

// ApplicationModuleRelationshipMutation represents an operation that mutates the ApplicationModuleRelationship nodes in the graph.
type ApplicationModuleRelationshipMutation struct {
	config
	op                 Op
	typ                string
	createTime         *time.Time
	updateTime         *time.Time
	version            *string
	name               *string
	attributes         *map[string]interface{}
	clearedFields      map[string]struct{}
	application        *oid.ID
	clearedapplication bool
	module             *string
	clearedmodule      bool
	done               bool
	oldValue           func(context.Context) (*ApplicationModuleRelationship, error)
	predicates         []predicate.ApplicationModuleRelationship
}

var _ ent.Mutation = (*ApplicationModuleRelationshipMutation)(nil)

// applicationModuleRelationshipOption allows management of the mutation configuration using functional options.
type applicationModuleRelationshipOption func(*ApplicationModuleRelationshipMutation)

// newApplicationModuleRelationshipMutation creates new mutation for the ApplicationModuleRelationship entity.
func newApplicationModuleRelationshipMutation(c config, op Op, opts ...applicationModuleRelationshipOption) *ApplicationModuleRelationshipMutation {
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

// ResetUpdateTime resets all changes to the "updateTime" field.
func (m *ApplicationModuleRelationshipMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetApplicationID sets the "application_id" field.
func (m *ApplicationModuleRelationshipMutation) SetApplicationID(o oid.ID) {
	m.application = &o
}

// ApplicationID returns the value of the "application_id" field in the mutation.
func (m *ApplicationModuleRelationshipMutation) ApplicationID() (r oid.ID, exists bool) {
	v := m.application
	if v == nil {
		return
	}
	return *v, true
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

// ResetModuleID resets all changes to the "module_id" field.
func (m *ApplicationModuleRelationshipMutation) ResetModuleID() {
	m.module = nil
}

// SetVersion sets the "version" field.
func (m *ApplicationModuleRelationshipMutation) SetVersion(s string) {
	m.version = &s
}

// Version returns the value of the "version" field in the mutation.
func (m *ApplicationModuleRelationshipMutation) Version() (r string, exists bool) {
	v := m.version
	if v == nil {
		return
	}
	return *v, true
}

// ResetVersion resets all changes to the "version" field.
func (m *ApplicationModuleRelationshipMutation) ResetVersion() {
	m.version = nil
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

// ResetName resets all changes to the "name" field.
func (m *ApplicationModuleRelationshipMutation) ResetName() {
	m.name = nil
}

// SetAttributes sets the "attributes" field.
func (m *ApplicationModuleRelationshipMutation) SetAttributes(value map[string]interface{}) {
	m.attributes = &value
}

// Attributes returns the value of the "attributes" field in the mutation.
func (m *ApplicationModuleRelationshipMutation) Attributes() (r map[string]interface{}, exists bool) {
	v := m.attributes
	if v == nil {
		return
	}
	return *v, true
}

// ClearAttributes clears the value of the "attributes" field.
func (m *ApplicationModuleRelationshipMutation) ClearAttributes() {
	m.attributes = nil
	m.clearedFields[applicationmodulerelationship.FieldAttributes] = struct{}{}
}

// AttributesCleared returns if the "attributes" field was cleared in this mutation.
func (m *ApplicationModuleRelationshipMutation) AttributesCleared() bool {
	_, ok := m.clearedFields[applicationmodulerelationship.FieldAttributes]
	return ok
}

// ResetAttributes resets all changes to the "attributes" field.
func (m *ApplicationModuleRelationshipMutation) ResetAttributes() {
	m.attributes = nil
	delete(m.clearedFields, applicationmodulerelationship.FieldAttributes)
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
func (m *ApplicationModuleRelationshipMutation) ApplicationIDs() (ids []oid.ID) {
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
	fields := make([]string, 0, 7)
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
	if m.version != nil {
		fields = append(fields, applicationmodulerelationship.FieldVersion)
	}
	if m.name != nil {
		fields = append(fields, applicationmodulerelationship.FieldName)
	}
	if m.attributes != nil {
		fields = append(fields, applicationmodulerelationship.FieldAttributes)
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
	case applicationmodulerelationship.FieldVersion:
		return m.Version()
	case applicationmodulerelationship.FieldName:
		return m.Name()
	case applicationmodulerelationship.FieldAttributes:
		return m.Attributes()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ApplicationModuleRelationshipMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	return nil, errors.New("edge schema ApplicationModuleRelationship does not support getting old values")
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
		v, ok := value.(oid.ID)
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
	case applicationmodulerelationship.FieldVersion:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVersion(v)
		return nil
	case applicationmodulerelationship.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case applicationmodulerelationship.FieldAttributes:
		v, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetAttributes(v)
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
	if m.FieldCleared(applicationmodulerelationship.FieldAttributes) {
		fields = append(fields, applicationmodulerelationship.FieldAttributes)
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
	case applicationmodulerelationship.FieldAttributes:
		m.ClearAttributes()
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
	case applicationmodulerelationship.FieldVersion:
		m.ResetVersion()
		return nil
	case applicationmodulerelationship.FieldName:
		m.ResetName()
		return nil
	case applicationmodulerelationship.FieldAttributes:
		m.ResetAttributes()
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
	id                 *oid.ID
	createTime         *time.Time
	updateTime         *time.Time
	module             *string
	mode               *string
	_type              *string
	name               *string
	deployerType       *string
	status             *types.ApplicationResourceStatus
	clearedFields      map[string]struct{}
	instance           *oid.ID
	clearedinstance    bool
	connector          *oid.ID
	clearedconnector   bool
	composition        *oid.ID
	clearedcomposition bool
	components         map[oid.ID]struct{}
	removedcomponents  map[oid.ID]struct{}
	clearedcomponents  bool
	done               bool
	oldValue           func(context.Context) (*ApplicationResource, error)
	predicates         []predicate.ApplicationResource
}

var _ ent.Mutation = (*ApplicationResourceMutation)(nil)

// applicationResourceOption allows management of the mutation configuration using functional options.
type applicationResourceOption func(*ApplicationResourceMutation)

// newApplicationResourceMutation creates new mutation for the ApplicationResource entity.
func newApplicationResourceMutation(c config, op Op, opts ...applicationResourceOption) *ApplicationResourceMutation {
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
func withApplicationResourceID(id oid.ID) applicationResourceOption {
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
func withApplicationResource(node *ApplicationResource) applicationResourceOption {
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

// SetInstanceID sets the "instanceID" field.
func (m *ApplicationResourceMutation) SetInstanceID(o oid.ID) {
	m.instance = &o
}

// InstanceID returns the value of the "instanceID" field in the mutation.
func (m *ApplicationResourceMutation) InstanceID() (r oid.ID, exists bool) {
	v := m.instance
	if v == nil {
		return
	}
	return *v, true
}

// OldInstanceID returns the old "instanceID" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldInstanceID(ctx context.Context) (v oid.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldInstanceID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldInstanceID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldInstanceID: %w", err)
	}
	return oldValue.InstanceID, nil
}

// ResetInstanceID resets all changes to the "instanceID" field.
func (m *ApplicationResourceMutation) ResetInstanceID() {
	m.instance = nil
}

// SetConnectorID sets the "connectorID" field.
func (m *ApplicationResourceMutation) SetConnectorID(o oid.ID) {
	m.connector = &o
}

// ConnectorID returns the value of the "connectorID" field in the mutation.
func (m *ApplicationResourceMutation) ConnectorID() (r oid.ID, exists bool) {
	v := m.connector
	if v == nil {
		return
	}
	return *v, true
}

// OldConnectorID returns the old "connectorID" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldConnectorID(ctx context.Context) (v oid.ID, err error) {
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

// SetCompositionID sets the "compositionID" field.
func (m *ApplicationResourceMutation) SetCompositionID(o oid.ID) {
	m.composition = &o
}

// CompositionID returns the value of the "compositionID" field in the mutation.
func (m *ApplicationResourceMutation) CompositionID() (r oid.ID, exists bool) {
	v := m.composition
	if v == nil {
		return
	}
	return *v, true
}

// OldCompositionID returns the old "compositionID" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldCompositionID(ctx context.Context) (v oid.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCompositionID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCompositionID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCompositionID: %w", err)
	}
	return oldValue.CompositionID, nil
}

// ClearCompositionID clears the value of the "compositionID" field.
func (m *ApplicationResourceMutation) ClearCompositionID() {
	m.composition = nil
	m.clearedFields[applicationresource.FieldCompositionID] = struct{}{}
}

// CompositionIDCleared returns if the "compositionID" field was cleared in this mutation.
func (m *ApplicationResourceMutation) CompositionIDCleared() bool {
	_, ok := m.clearedFields[applicationresource.FieldCompositionID]
	return ok
}

// ResetCompositionID resets all changes to the "compositionID" field.
func (m *ApplicationResourceMutation) ResetCompositionID() {
	m.composition = nil
	delete(m.clearedFields, applicationresource.FieldCompositionID)
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

// SetDeployerType sets the "deployerType" field.
func (m *ApplicationResourceMutation) SetDeployerType(s string) {
	m.deployerType = &s
}

// DeployerType returns the value of the "deployerType" field in the mutation.
func (m *ApplicationResourceMutation) DeployerType() (r string, exists bool) {
	v := m.deployerType
	if v == nil {
		return
	}
	return *v, true
}

// OldDeployerType returns the old "deployerType" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldDeployerType(ctx context.Context) (v string, err error) {
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
func (m *ApplicationResourceMutation) ResetDeployerType() {
	m.deployerType = nil
}

// SetStatus sets the "status" field.
func (m *ApplicationResourceMutation) SetStatus(trs types.ApplicationResourceStatus) {
	m.status = &trs
}

// Status returns the value of the "status" field in the mutation.
func (m *ApplicationResourceMutation) Status() (r types.ApplicationResourceStatus, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the ApplicationResource entity.
// If the ApplicationResource object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationResourceMutation) OldStatus(ctx context.Context) (v types.ApplicationResourceStatus, err error) {
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

// ClearInstance clears the "instance" edge to the ApplicationInstance entity.
func (m *ApplicationResourceMutation) ClearInstance() {
	m.clearedinstance = true
}

// InstanceCleared reports if the "instance" edge to the ApplicationInstance entity was cleared.
func (m *ApplicationResourceMutation) InstanceCleared() bool {
	return m.clearedinstance
}

// InstanceIDs returns the "instance" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// InstanceID instead. It exists only for internal usage by the builders.
func (m *ApplicationResourceMutation) InstanceIDs() (ids []oid.ID) {
	if id := m.instance; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetInstance resets all changes to the "instance" edge.
func (m *ApplicationResourceMutation) ResetInstance() {
	m.instance = nil
	m.clearedinstance = false
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
func (m *ApplicationResourceMutation) ConnectorIDs() (ids []oid.ID) {
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

// ClearComposition clears the "composition" edge to the ApplicationResource entity.
func (m *ApplicationResourceMutation) ClearComposition() {
	m.clearedcomposition = true
}

// CompositionCleared reports if the "composition" edge to the ApplicationResource entity was cleared.
func (m *ApplicationResourceMutation) CompositionCleared() bool {
	return m.CompositionIDCleared() || m.clearedcomposition
}

// CompositionIDs returns the "composition" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// CompositionID instead. It exists only for internal usage by the builders.
func (m *ApplicationResourceMutation) CompositionIDs() (ids []oid.ID) {
	if id := m.composition; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetComposition resets all changes to the "composition" edge.
func (m *ApplicationResourceMutation) ResetComposition() {
	m.composition = nil
	m.clearedcomposition = false
}

// AddComponentIDs adds the "components" edge to the ApplicationResource entity by ids.
func (m *ApplicationResourceMutation) AddComponentIDs(ids ...oid.ID) {
	if m.components == nil {
		m.components = make(map[oid.ID]struct{})
	}
	for i := range ids {
		m.components[ids[i]] = struct{}{}
	}
}

// ClearComponents clears the "components" edge to the ApplicationResource entity.
func (m *ApplicationResourceMutation) ClearComponents() {
	m.clearedcomponents = true
}

// ComponentsCleared reports if the "components" edge to the ApplicationResource entity was cleared.
func (m *ApplicationResourceMutation) ComponentsCleared() bool {
	return m.clearedcomponents
}

// RemoveComponentIDs removes the "components" edge to the ApplicationResource entity by IDs.
func (m *ApplicationResourceMutation) RemoveComponentIDs(ids ...oid.ID) {
	if m.removedcomponents == nil {
		m.removedcomponents = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.components, ids[i])
		m.removedcomponents[ids[i]] = struct{}{}
	}
}

// RemovedComponents returns the removed IDs of the "components" edge to the ApplicationResource entity.
func (m *ApplicationResourceMutation) RemovedComponentsIDs() (ids []oid.ID) {
	for id := range m.removedcomponents {
		ids = append(ids, id)
	}
	return
}

// ComponentsIDs returns the "components" edge IDs in the mutation.
func (m *ApplicationResourceMutation) ComponentsIDs() (ids []oid.ID) {
	for id := range m.components {
		ids = append(ids, id)
	}
	return
}

// ResetComponents resets all changes to the "components" edge.
func (m *ApplicationResourceMutation) ResetComponents() {
	m.components = nil
	m.clearedcomponents = false
	m.removedcomponents = nil
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
	fields := make([]string, 0, 11)
	if m.createTime != nil {
		fields = append(fields, applicationresource.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, applicationresource.FieldUpdateTime)
	}
	if m.instance != nil {
		fields = append(fields, applicationresource.FieldInstanceID)
	}
	if m.connector != nil {
		fields = append(fields, applicationresource.FieldConnectorID)
	}
	if m.composition != nil {
		fields = append(fields, applicationresource.FieldCompositionID)
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
	if m.deployerType != nil {
		fields = append(fields, applicationresource.FieldDeployerType)
	}
	if m.status != nil {
		fields = append(fields, applicationresource.FieldStatus)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ApplicationResourceMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case applicationresource.FieldCreateTime:
		return m.CreateTime()
	case applicationresource.FieldUpdateTime:
		return m.UpdateTime()
	case applicationresource.FieldInstanceID:
		return m.InstanceID()
	case applicationresource.FieldConnectorID:
		return m.ConnectorID()
	case applicationresource.FieldCompositionID:
		return m.CompositionID()
	case applicationresource.FieldModule:
		return m.Module()
	case applicationresource.FieldMode:
		return m.Mode()
	case applicationresource.FieldType:
		return m.GetType()
	case applicationresource.FieldName:
		return m.Name()
	case applicationresource.FieldDeployerType:
		return m.DeployerType()
	case applicationresource.FieldStatus:
		return m.Status()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ApplicationResourceMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case applicationresource.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case applicationresource.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case applicationresource.FieldInstanceID:
		return m.OldInstanceID(ctx)
	case applicationresource.FieldConnectorID:
		return m.OldConnectorID(ctx)
	case applicationresource.FieldCompositionID:
		return m.OldCompositionID(ctx)
	case applicationresource.FieldModule:
		return m.OldModule(ctx)
	case applicationresource.FieldMode:
		return m.OldMode(ctx)
	case applicationresource.FieldType:
		return m.OldType(ctx)
	case applicationresource.FieldName:
		return m.OldName(ctx)
	case applicationresource.FieldDeployerType:
		return m.OldDeployerType(ctx)
	case applicationresource.FieldStatus:
		return m.OldStatus(ctx)
	}
	return nil, fmt.Errorf("unknown ApplicationResource field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ApplicationResourceMutation) SetField(name string, value ent.Value) error {
	switch name {
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
	case applicationresource.FieldInstanceID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetInstanceID(v)
		return nil
	case applicationresource.FieldConnectorID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetConnectorID(v)
		return nil
	case applicationresource.FieldCompositionID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCompositionID(v)
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
	case applicationresource.FieldDeployerType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeployerType(v)
		return nil
	case applicationresource.FieldStatus:
		v, ok := value.(types.ApplicationResourceStatus)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
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
	if m.FieldCleared(applicationresource.FieldCompositionID) {
		fields = append(fields, applicationresource.FieldCompositionID)
	}
	if m.FieldCleared(applicationresource.FieldStatus) {
		fields = append(fields, applicationresource.FieldStatus)
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
	case applicationresource.FieldCompositionID:
		m.ClearCompositionID()
		return nil
	case applicationresource.FieldStatus:
		m.ClearStatus()
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ApplicationResourceMutation) ResetField(name string) error {
	switch name {
	case applicationresource.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case applicationresource.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case applicationresource.FieldInstanceID:
		m.ResetInstanceID()
		return nil
	case applicationresource.FieldConnectorID:
		m.ResetConnectorID()
		return nil
	case applicationresource.FieldCompositionID:
		m.ResetCompositionID()
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
	case applicationresource.FieldDeployerType:
		m.ResetDeployerType()
		return nil
	case applicationresource.FieldStatus:
		m.ResetStatus()
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationResourceMutation) AddedEdges() []string {
	edges := make([]string, 0, 4)
	if m.instance != nil {
		edges = append(edges, applicationresource.EdgeInstance)
	}
	if m.connector != nil {
		edges = append(edges, applicationresource.EdgeConnector)
	}
	if m.composition != nil {
		edges = append(edges, applicationresource.EdgeComposition)
	}
	if m.components != nil {
		edges = append(edges, applicationresource.EdgeComponents)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ApplicationResourceMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case applicationresource.EdgeInstance:
		if id := m.instance; id != nil {
			return []ent.Value{*id}
		}
	case applicationresource.EdgeConnector:
		if id := m.connector; id != nil {
			return []ent.Value{*id}
		}
	case applicationresource.EdgeComposition:
		if id := m.composition; id != nil {
			return []ent.Value{*id}
		}
	case applicationresource.EdgeComponents:
		ids := make([]ent.Value, 0, len(m.components))
		for id := range m.components {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ApplicationResourceMutation) RemovedEdges() []string {
	edges := make([]string, 0, 4)
	if m.removedcomponents != nil {
		edges = append(edges, applicationresource.EdgeComponents)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ApplicationResourceMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case applicationresource.EdgeComponents:
		ids := make([]ent.Value, 0, len(m.removedcomponents))
		for id := range m.removedcomponents {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ApplicationResourceMutation) ClearedEdges() []string {
	edges := make([]string, 0, 4)
	if m.clearedinstance {
		edges = append(edges, applicationresource.EdgeInstance)
	}
	if m.clearedconnector {
		edges = append(edges, applicationresource.EdgeConnector)
	}
	if m.clearedcomposition {
		edges = append(edges, applicationresource.EdgeComposition)
	}
	if m.clearedcomponents {
		edges = append(edges, applicationresource.EdgeComponents)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ApplicationResourceMutation) EdgeCleared(name string) bool {
	switch name {
	case applicationresource.EdgeInstance:
		return m.clearedinstance
	case applicationresource.EdgeConnector:
		return m.clearedconnector
	case applicationresource.EdgeComposition:
		return m.clearedcomposition
	case applicationresource.EdgeComponents:
		return m.clearedcomponents
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationResourceMutation) ClearEdge(name string) error {
	switch name {
	case applicationresource.EdgeInstance:
		m.ClearInstance()
		return nil
	case applicationresource.EdgeConnector:
		m.ClearConnector()
		return nil
	case applicationresource.EdgeComposition:
		m.ClearComposition()
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ApplicationResourceMutation) ResetEdge(name string) error {
	switch name {
	case applicationresource.EdgeInstance:
		m.ResetInstance()
		return nil
	case applicationresource.EdgeConnector:
		m.ResetConnector()
		return nil
	case applicationresource.EdgeComposition:
		m.ResetComposition()
		return nil
	case applicationresource.EdgeComponents:
		m.ResetComponents()
		return nil
	}
	return fmt.Errorf("unknown ApplicationResource edge %s", name)
}

// ApplicationRevisionMutation represents an operation that mutates the ApplicationRevision nodes in the graph.
type ApplicationRevisionMutation struct {
	config
	op                              Op
	typ                             string
	id                              *oid.ID
	status                          *string
	statusMessage                   *string
	createTime                      *time.Time
	modules                         *[]types.ApplicationModule
	appendmodules                   []types.ApplicationModule
	inputVariables                  *map[string]interface{}
	inputPlan                       *string
	output                          *string
	deployerType                    *string
	duration                        *int
	addduration                     *int
	previousRequiredProviders       *[]types.ProviderRequirement
	appendpreviousRequiredProviders []types.ProviderRequirement
	clearedFields                   map[string]struct{}
	instance                        *oid.ID
	clearedinstance                 bool
	environment                     *oid.ID
	clearedenvironment              bool
	done                            bool
	oldValue                        func(context.Context) (*ApplicationRevision, error)
	predicates                      []predicate.ApplicationRevision
}

var _ ent.Mutation = (*ApplicationRevisionMutation)(nil)

// applicationRevisionOption allows management of the mutation configuration using functional options.
type applicationRevisionOption func(*ApplicationRevisionMutation)

// newApplicationRevisionMutation creates new mutation for the ApplicationRevision entity.
func newApplicationRevisionMutation(c config, op Op, opts ...applicationRevisionOption) *ApplicationRevisionMutation {
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
func withApplicationRevisionID(id oid.ID) applicationRevisionOption {
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
func withApplicationRevision(node *ApplicationRevision) applicationRevisionOption {
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

// SetInstanceID sets the "instanceID" field.
func (m *ApplicationRevisionMutation) SetInstanceID(o oid.ID) {
	m.instance = &o
}

// InstanceID returns the value of the "instanceID" field in the mutation.
func (m *ApplicationRevisionMutation) InstanceID() (r oid.ID, exists bool) {
	v := m.instance
	if v == nil {
		return
	}
	return *v, true
}

// OldInstanceID returns the old "instanceID" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldInstanceID(ctx context.Context) (v oid.ID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldInstanceID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldInstanceID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldInstanceID: %w", err)
	}
	return oldValue.InstanceID, nil
}

// ResetInstanceID resets all changes to the "instanceID" field.
func (m *ApplicationRevisionMutation) ResetInstanceID() {
	m.instance = nil
}

// SetEnvironmentID sets the "environmentID" field.
func (m *ApplicationRevisionMutation) SetEnvironmentID(o oid.ID) {
	m.environment = &o
}

// EnvironmentID returns the value of the "environmentID" field in the mutation.
func (m *ApplicationRevisionMutation) EnvironmentID() (r oid.ID, exists bool) {
	v := m.environment
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

// SetPreviousRequiredProviders sets the "previousRequiredProviders" field.
func (m *ApplicationRevisionMutation) SetPreviousRequiredProviders(tr []types.ProviderRequirement) {
	m.previousRequiredProviders = &tr
	m.appendpreviousRequiredProviders = nil
}

// PreviousRequiredProviders returns the value of the "previousRequiredProviders" field in the mutation.
func (m *ApplicationRevisionMutation) PreviousRequiredProviders() (r []types.ProviderRequirement, exists bool) {
	v := m.previousRequiredProviders
	if v == nil {
		return
	}
	return *v, true
}

// OldPreviousRequiredProviders returns the old "previousRequiredProviders" field's value of the ApplicationRevision entity.
// If the ApplicationRevision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ApplicationRevisionMutation) OldPreviousRequiredProviders(ctx context.Context) (v []types.ProviderRequirement, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPreviousRequiredProviders is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPreviousRequiredProviders requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPreviousRequiredProviders: %w", err)
	}
	return oldValue.PreviousRequiredProviders, nil
}

// AppendPreviousRequiredProviders adds tr to the "previousRequiredProviders" field.
func (m *ApplicationRevisionMutation) AppendPreviousRequiredProviders(tr []types.ProviderRequirement) {
	m.appendpreviousRequiredProviders = append(m.appendpreviousRequiredProviders, tr...)
}

// AppendedPreviousRequiredProviders returns the list of values that were appended to the "previousRequiredProviders" field in this mutation.
func (m *ApplicationRevisionMutation) AppendedPreviousRequiredProviders() ([]types.ProviderRequirement, bool) {
	if len(m.appendpreviousRequiredProviders) == 0 {
		return nil, false
	}
	return m.appendpreviousRequiredProviders, true
}

// ResetPreviousRequiredProviders resets all changes to the "previousRequiredProviders" field.
func (m *ApplicationRevisionMutation) ResetPreviousRequiredProviders() {
	m.previousRequiredProviders = nil
	m.appendpreviousRequiredProviders = nil
}

// ClearInstance clears the "instance" edge to the ApplicationInstance entity.
func (m *ApplicationRevisionMutation) ClearInstance() {
	m.clearedinstance = true
}

// InstanceCleared reports if the "instance" edge to the ApplicationInstance entity was cleared.
func (m *ApplicationRevisionMutation) InstanceCleared() bool {
	return m.clearedinstance
}

// InstanceIDs returns the "instance" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// InstanceID instead. It exists only for internal usage by the builders.
func (m *ApplicationRevisionMutation) InstanceIDs() (ids []oid.ID) {
	if id := m.instance; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetInstance resets all changes to the "instance" edge.
func (m *ApplicationRevisionMutation) ResetInstance() {
	m.instance = nil
	m.clearedinstance = false
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
func (m *ApplicationRevisionMutation) EnvironmentIDs() (ids []oid.ID) {
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
	fields := make([]string, 0, 12)
	if m.status != nil {
		fields = append(fields, applicationrevision.FieldStatus)
	}
	if m.statusMessage != nil {
		fields = append(fields, applicationrevision.FieldStatusMessage)
	}
	if m.createTime != nil {
		fields = append(fields, applicationrevision.FieldCreateTime)
	}
	if m.instance != nil {
		fields = append(fields, applicationrevision.FieldInstanceID)
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
	if m.previousRequiredProviders != nil {
		fields = append(fields, applicationrevision.FieldPreviousRequiredProviders)
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
	case applicationrevision.FieldInstanceID:
		return m.InstanceID()
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
	case applicationrevision.FieldPreviousRequiredProviders:
		return m.PreviousRequiredProviders()
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
	case applicationrevision.FieldInstanceID:
		return m.OldInstanceID(ctx)
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
	case applicationrevision.FieldPreviousRequiredProviders:
		return m.OldPreviousRequiredProviders(ctx)
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
	case applicationrevision.FieldInstanceID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetInstanceID(v)
		return nil
	case applicationrevision.FieldEnvironmentID:
		v, ok := value.(oid.ID)
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
	case applicationrevision.FieldPreviousRequiredProviders:
		v, ok := value.([]types.ProviderRequirement)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPreviousRequiredProviders(v)
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
	case applicationrevision.FieldInstanceID:
		m.ResetInstanceID()
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
	case applicationrevision.FieldPreviousRequiredProviders:
		m.ResetPreviousRequiredProviders()
		return nil
	}
	return fmt.Errorf("unknown ApplicationRevision field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ApplicationRevisionMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.instance != nil {
		edges = append(edges, applicationrevision.EdgeInstance)
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
	case applicationrevision.EdgeInstance:
		if id := m.instance; id != nil {
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
	if m.clearedinstance {
		edges = append(edges, applicationrevision.EdgeInstance)
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
	case applicationrevision.EdgeInstance:
		return m.clearedinstance
	case applicationrevision.EdgeEnvironment:
		return m.clearedenvironment
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ApplicationRevisionMutation) ClearEdge(name string) error {
	switch name {
	case applicationrevision.EdgeInstance:
		m.ClearInstance()
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
	case applicationrevision.EdgeInstance:
		m.ResetInstance()
		return nil
	case applicationrevision.EdgeEnvironment:
		m.ResetEnvironment()
		return nil
	}
	return fmt.Errorf("unknown ApplicationRevision edge %s", name)
}

// ClusterCostMutation represents an operation that mutates the ClusterCost nodes in the graph.
type ClusterCostMutation struct {
	config
	op                Op
	typ               string
	id                *int
	startTime         *time.Time
	endTime           *time.Time
	minutes           *float64
	addminutes        *float64
	clusterName       *string
	totalCost         *float64
	addtotalCost      *float64
	currency          *int
	addcurrency       *int
	allocationCost    *float64
	addallocationCost *float64
	idleCost          *float64
	addidleCost       *float64
	managementCost    *float64
	addmanagementCost *float64
	clearedFields     map[string]struct{}
	connector         *oid.ID
	clearedconnector  bool
	done              bool
	oldValue          func(context.Context) (*ClusterCost, error)
	predicates        []predicate.ClusterCost
}

var _ ent.Mutation = (*ClusterCostMutation)(nil)

// clusterCostOption allows management of the mutation configuration using functional options.
type clusterCostOption func(*ClusterCostMutation)

// newClusterCostMutation creates new mutation for the ClusterCost entity.
func newClusterCostMutation(c config, op Op, opts ...clusterCostOption) *ClusterCostMutation {
	m := &ClusterCostMutation{
		config:        c,
		op:            op,
		typ:           TypeClusterCost,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withClusterCostID sets the ID field of the mutation.
func withClusterCostID(id int) clusterCostOption {
	return func(m *ClusterCostMutation) {
		var (
			err   error
			once  sync.Once
			value *ClusterCost
		)
		m.oldValue = func(ctx context.Context) (*ClusterCost, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().ClusterCost.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withClusterCost sets the old ClusterCost of the mutation.
func withClusterCost(node *ClusterCost) clusterCostOption {
	return func(m *ClusterCostMutation) {
		m.oldValue = func(context.Context) (*ClusterCost, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ClusterCostMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ClusterCostMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ClusterCostMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ClusterCostMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().ClusterCost.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetStartTime sets the "startTime" field.
func (m *ClusterCostMutation) SetStartTime(t time.Time) {
	m.startTime = &t
}

// StartTime returns the value of the "startTime" field in the mutation.
func (m *ClusterCostMutation) StartTime() (r time.Time, exists bool) {
	v := m.startTime
	if v == nil {
		return
	}
	return *v, true
}

// OldStartTime returns the old "startTime" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldStartTime(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStartTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStartTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStartTime: %w", err)
	}
	return oldValue.StartTime, nil
}

// ResetStartTime resets all changes to the "startTime" field.
func (m *ClusterCostMutation) ResetStartTime() {
	m.startTime = nil
}

// SetEndTime sets the "endTime" field.
func (m *ClusterCostMutation) SetEndTime(t time.Time) {
	m.endTime = &t
}

// EndTime returns the value of the "endTime" field in the mutation.
func (m *ClusterCostMutation) EndTime() (r time.Time, exists bool) {
	v := m.endTime
	if v == nil {
		return
	}
	return *v, true
}

// OldEndTime returns the old "endTime" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldEndTime(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEndTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEndTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEndTime: %w", err)
	}
	return oldValue.EndTime, nil
}

// ResetEndTime resets all changes to the "endTime" field.
func (m *ClusterCostMutation) ResetEndTime() {
	m.endTime = nil
}

// SetMinutes sets the "minutes" field.
func (m *ClusterCostMutation) SetMinutes(f float64) {
	m.minutes = &f
	m.addminutes = nil
}

// Minutes returns the value of the "minutes" field in the mutation.
func (m *ClusterCostMutation) Minutes() (r float64, exists bool) {
	v := m.minutes
	if v == nil {
		return
	}
	return *v, true
}

// OldMinutes returns the old "minutes" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldMinutes(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldMinutes is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldMinutes requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldMinutes: %w", err)
	}
	return oldValue.Minutes, nil
}

// AddMinutes adds f to the "minutes" field.
func (m *ClusterCostMutation) AddMinutes(f float64) {
	if m.addminutes != nil {
		*m.addminutes += f
	} else {
		m.addminutes = &f
	}
}

// AddedMinutes returns the value that was added to the "minutes" field in this mutation.
func (m *ClusterCostMutation) AddedMinutes() (r float64, exists bool) {
	v := m.addminutes
	if v == nil {
		return
	}
	return *v, true
}

// ResetMinutes resets all changes to the "minutes" field.
func (m *ClusterCostMutation) ResetMinutes() {
	m.minutes = nil
	m.addminutes = nil
}

// SetConnectorID sets the "connectorID" field.
func (m *ClusterCostMutation) SetConnectorID(o oid.ID) {
	m.connector = &o
}

// ConnectorID returns the value of the "connectorID" field in the mutation.
func (m *ClusterCostMutation) ConnectorID() (r oid.ID, exists bool) {
	v := m.connector
	if v == nil {
		return
	}
	return *v, true
}

// OldConnectorID returns the old "connectorID" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldConnectorID(ctx context.Context) (v oid.ID, err error) {
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
func (m *ClusterCostMutation) ResetConnectorID() {
	m.connector = nil
}

// SetClusterName sets the "clusterName" field.
func (m *ClusterCostMutation) SetClusterName(s string) {
	m.clusterName = &s
}

// ClusterName returns the value of the "clusterName" field in the mutation.
func (m *ClusterCostMutation) ClusterName() (r string, exists bool) {
	v := m.clusterName
	if v == nil {
		return
	}
	return *v, true
}

// OldClusterName returns the old "clusterName" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldClusterName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldClusterName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldClusterName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldClusterName: %w", err)
	}
	return oldValue.ClusterName, nil
}

// ResetClusterName resets all changes to the "clusterName" field.
func (m *ClusterCostMutation) ResetClusterName() {
	m.clusterName = nil
}

// SetTotalCost sets the "totalCost" field.
func (m *ClusterCostMutation) SetTotalCost(f float64) {
	m.totalCost = &f
	m.addtotalCost = nil
}

// TotalCost returns the value of the "totalCost" field in the mutation.
func (m *ClusterCostMutation) TotalCost() (r float64, exists bool) {
	v := m.totalCost
	if v == nil {
		return
	}
	return *v, true
}

// OldTotalCost returns the old "totalCost" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldTotalCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTotalCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTotalCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTotalCost: %w", err)
	}
	return oldValue.TotalCost, nil
}

// AddTotalCost adds f to the "totalCost" field.
func (m *ClusterCostMutation) AddTotalCost(f float64) {
	if m.addtotalCost != nil {
		*m.addtotalCost += f
	} else {
		m.addtotalCost = &f
	}
}

// AddedTotalCost returns the value that was added to the "totalCost" field in this mutation.
func (m *ClusterCostMutation) AddedTotalCost() (r float64, exists bool) {
	v := m.addtotalCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetTotalCost resets all changes to the "totalCost" field.
func (m *ClusterCostMutation) ResetTotalCost() {
	m.totalCost = nil
	m.addtotalCost = nil
}

// SetCurrency sets the "currency" field.
func (m *ClusterCostMutation) SetCurrency(i int) {
	m.currency = &i
	m.addcurrency = nil
}

// Currency returns the value of the "currency" field in the mutation.
func (m *ClusterCostMutation) Currency() (r int, exists bool) {
	v := m.currency
	if v == nil {
		return
	}
	return *v, true
}

// OldCurrency returns the old "currency" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldCurrency(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCurrency is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCurrency requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCurrency: %w", err)
	}
	return oldValue.Currency, nil
}

// AddCurrency adds i to the "currency" field.
func (m *ClusterCostMutation) AddCurrency(i int) {
	if m.addcurrency != nil {
		*m.addcurrency += i
	} else {
		m.addcurrency = &i
	}
}

// AddedCurrency returns the value that was added to the "currency" field in this mutation.
func (m *ClusterCostMutation) AddedCurrency() (r int, exists bool) {
	v := m.addcurrency
	if v == nil {
		return
	}
	return *v, true
}

// ClearCurrency clears the value of the "currency" field.
func (m *ClusterCostMutation) ClearCurrency() {
	m.currency = nil
	m.addcurrency = nil
	m.clearedFields[clustercost.FieldCurrency] = struct{}{}
}

// CurrencyCleared returns if the "currency" field was cleared in this mutation.
func (m *ClusterCostMutation) CurrencyCleared() bool {
	_, ok := m.clearedFields[clustercost.FieldCurrency]
	return ok
}

// ResetCurrency resets all changes to the "currency" field.
func (m *ClusterCostMutation) ResetCurrency() {
	m.currency = nil
	m.addcurrency = nil
	delete(m.clearedFields, clustercost.FieldCurrency)
}

// SetAllocationCost sets the "allocationCost" field.
func (m *ClusterCostMutation) SetAllocationCost(f float64) {
	m.allocationCost = &f
	m.addallocationCost = nil
}

// AllocationCost returns the value of the "allocationCost" field in the mutation.
func (m *ClusterCostMutation) AllocationCost() (r float64, exists bool) {
	v := m.allocationCost
	if v == nil {
		return
	}
	return *v, true
}

// OldAllocationCost returns the old "allocationCost" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldAllocationCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldAllocationCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldAllocationCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldAllocationCost: %w", err)
	}
	return oldValue.AllocationCost, nil
}

// AddAllocationCost adds f to the "allocationCost" field.
func (m *ClusterCostMutation) AddAllocationCost(f float64) {
	if m.addallocationCost != nil {
		*m.addallocationCost += f
	} else {
		m.addallocationCost = &f
	}
}

// AddedAllocationCost returns the value that was added to the "allocationCost" field in this mutation.
func (m *ClusterCostMutation) AddedAllocationCost() (r float64, exists bool) {
	v := m.addallocationCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetAllocationCost resets all changes to the "allocationCost" field.
func (m *ClusterCostMutation) ResetAllocationCost() {
	m.allocationCost = nil
	m.addallocationCost = nil
}

// SetIdleCost sets the "idleCost" field.
func (m *ClusterCostMutation) SetIdleCost(f float64) {
	m.idleCost = &f
	m.addidleCost = nil
}

// IdleCost returns the value of the "idleCost" field in the mutation.
func (m *ClusterCostMutation) IdleCost() (r float64, exists bool) {
	v := m.idleCost
	if v == nil {
		return
	}
	return *v, true
}

// OldIdleCost returns the old "idleCost" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldIdleCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldIdleCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldIdleCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldIdleCost: %w", err)
	}
	return oldValue.IdleCost, nil
}

// AddIdleCost adds f to the "idleCost" field.
func (m *ClusterCostMutation) AddIdleCost(f float64) {
	if m.addidleCost != nil {
		*m.addidleCost += f
	} else {
		m.addidleCost = &f
	}
}

// AddedIdleCost returns the value that was added to the "idleCost" field in this mutation.
func (m *ClusterCostMutation) AddedIdleCost() (r float64, exists bool) {
	v := m.addidleCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetIdleCost resets all changes to the "idleCost" field.
func (m *ClusterCostMutation) ResetIdleCost() {
	m.idleCost = nil
	m.addidleCost = nil
}

// SetManagementCost sets the "managementCost" field.
func (m *ClusterCostMutation) SetManagementCost(f float64) {
	m.managementCost = &f
	m.addmanagementCost = nil
}

// ManagementCost returns the value of the "managementCost" field in the mutation.
func (m *ClusterCostMutation) ManagementCost() (r float64, exists bool) {
	v := m.managementCost
	if v == nil {
		return
	}
	return *v, true
}

// OldManagementCost returns the old "managementCost" field's value of the ClusterCost entity.
// If the ClusterCost object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ClusterCostMutation) OldManagementCost(ctx context.Context) (v float64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldManagementCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldManagementCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldManagementCost: %w", err)
	}
	return oldValue.ManagementCost, nil
}

// AddManagementCost adds f to the "managementCost" field.
func (m *ClusterCostMutation) AddManagementCost(f float64) {
	if m.addmanagementCost != nil {
		*m.addmanagementCost += f
	} else {
		m.addmanagementCost = &f
	}
}

// AddedManagementCost returns the value that was added to the "managementCost" field in this mutation.
func (m *ClusterCostMutation) AddedManagementCost() (r float64, exists bool) {
	v := m.addmanagementCost
	if v == nil {
		return
	}
	return *v, true
}

// ResetManagementCost resets all changes to the "managementCost" field.
func (m *ClusterCostMutation) ResetManagementCost() {
	m.managementCost = nil
	m.addmanagementCost = nil
}

// ClearConnector clears the "connector" edge to the Connector entity.
func (m *ClusterCostMutation) ClearConnector() {
	m.clearedconnector = true
}

// ConnectorCleared reports if the "connector" edge to the Connector entity was cleared.
func (m *ClusterCostMutation) ConnectorCleared() bool {
	return m.clearedconnector
}

// ConnectorIDs returns the "connector" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ConnectorID instead. It exists only for internal usage by the builders.
func (m *ClusterCostMutation) ConnectorIDs() (ids []oid.ID) {
	if id := m.connector; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetConnector resets all changes to the "connector" edge.
func (m *ClusterCostMutation) ResetConnector() {
	m.connector = nil
	m.clearedconnector = false
}

// Where appends a list predicates to the ClusterCostMutation builder.
func (m *ClusterCostMutation) Where(ps ...predicate.ClusterCost) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ClusterCostMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ClusterCostMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.ClusterCost, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ClusterCostMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ClusterCostMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (ClusterCost).
func (m *ClusterCostMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ClusterCostMutation) Fields() []string {
	fields := make([]string, 0, 10)
	if m.startTime != nil {
		fields = append(fields, clustercost.FieldStartTime)
	}
	if m.endTime != nil {
		fields = append(fields, clustercost.FieldEndTime)
	}
	if m.minutes != nil {
		fields = append(fields, clustercost.FieldMinutes)
	}
	if m.connector != nil {
		fields = append(fields, clustercost.FieldConnectorID)
	}
	if m.clusterName != nil {
		fields = append(fields, clustercost.FieldClusterName)
	}
	if m.totalCost != nil {
		fields = append(fields, clustercost.FieldTotalCost)
	}
	if m.currency != nil {
		fields = append(fields, clustercost.FieldCurrency)
	}
	if m.allocationCost != nil {
		fields = append(fields, clustercost.FieldAllocationCost)
	}
	if m.idleCost != nil {
		fields = append(fields, clustercost.FieldIdleCost)
	}
	if m.managementCost != nil {
		fields = append(fields, clustercost.FieldManagementCost)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ClusterCostMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case clustercost.FieldStartTime:
		return m.StartTime()
	case clustercost.FieldEndTime:
		return m.EndTime()
	case clustercost.FieldMinutes:
		return m.Minutes()
	case clustercost.FieldConnectorID:
		return m.ConnectorID()
	case clustercost.FieldClusterName:
		return m.ClusterName()
	case clustercost.FieldTotalCost:
		return m.TotalCost()
	case clustercost.FieldCurrency:
		return m.Currency()
	case clustercost.FieldAllocationCost:
		return m.AllocationCost()
	case clustercost.FieldIdleCost:
		return m.IdleCost()
	case clustercost.FieldManagementCost:
		return m.ManagementCost()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ClusterCostMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case clustercost.FieldStartTime:
		return m.OldStartTime(ctx)
	case clustercost.FieldEndTime:
		return m.OldEndTime(ctx)
	case clustercost.FieldMinutes:
		return m.OldMinutes(ctx)
	case clustercost.FieldConnectorID:
		return m.OldConnectorID(ctx)
	case clustercost.FieldClusterName:
		return m.OldClusterName(ctx)
	case clustercost.FieldTotalCost:
		return m.OldTotalCost(ctx)
	case clustercost.FieldCurrency:
		return m.OldCurrency(ctx)
	case clustercost.FieldAllocationCost:
		return m.OldAllocationCost(ctx)
	case clustercost.FieldIdleCost:
		return m.OldIdleCost(ctx)
	case clustercost.FieldManagementCost:
		return m.OldManagementCost(ctx)
	}
	return nil, fmt.Errorf("unknown ClusterCost field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ClusterCostMutation) SetField(name string, value ent.Value) error {
	switch name {
	case clustercost.FieldStartTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStartTime(v)
		return nil
	case clustercost.FieldEndTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEndTime(v)
		return nil
	case clustercost.FieldMinutes:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMinutes(v)
		return nil
	case clustercost.FieldConnectorID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetConnectorID(v)
		return nil
	case clustercost.FieldClusterName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetClusterName(v)
		return nil
	case clustercost.FieldTotalCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTotalCost(v)
		return nil
	case clustercost.FieldCurrency:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCurrency(v)
		return nil
	case clustercost.FieldAllocationCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetAllocationCost(v)
		return nil
	case clustercost.FieldIdleCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIdleCost(v)
		return nil
	case clustercost.FieldManagementCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetManagementCost(v)
		return nil
	}
	return fmt.Errorf("unknown ClusterCost field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ClusterCostMutation) AddedFields() []string {
	var fields []string
	if m.addminutes != nil {
		fields = append(fields, clustercost.FieldMinutes)
	}
	if m.addtotalCost != nil {
		fields = append(fields, clustercost.FieldTotalCost)
	}
	if m.addcurrency != nil {
		fields = append(fields, clustercost.FieldCurrency)
	}
	if m.addallocationCost != nil {
		fields = append(fields, clustercost.FieldAllocationCost)
	}
	if m.addidleCost != nil {
		fields = append(fields, clustercost.FieldIdleCost)
	}
	if m.addmanagementCost != nil {
		fields = append(fields, clustercost.FieldManagementCost)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ClusterCostMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case clustercost.FieldMinutes:
		return m.AddedMinutes()
	case clustercost.FieldTotalCost:
		return m.AddedTotalCost()
	case clustercost.FieldCurrency:
		return m.AddedCurrency()
	case clustercost.FieldAllocationCost:
		return m.AddedAllocationCost()
	case clustercost.FieldIdleCost:
		return m.AddedIdleCost()
	case clustercost.FieldManagementCost:
		return m.AddedManagementCost()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ClusterCostMutation) AddField(name string, value ent.Value) error {
	switch name {
	case clustercost.FieldMinutes:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddMinutes(v)
		return nil
	case clustercost.FieldTotalCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddTotalCost(v)
		return nil
	case clustercost.FieldCurrency:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddCurrency(v)
		return nil
	case clustercost.FieldAllocationCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddAllocationCost(v)
		return nil
	case clustercost.FieldIdleCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIdleCost(v)
		return nil
	case clustercost.FieldManagementCost:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddManagementCost(v)
		return nil
	}
	return fmt.Errorf("unknown ClusterCost numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ClusterCostMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(clustercost.FieldCurrency) {
		fields = append(fields, clustercost.FieldCurrency)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ClusterCostMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ClusterCostMutation) ClearField(name string) error {
	switch name {
	case clustercost.FieldCurrency:
		m.ClearCurrency()
		return nil
	}
	return fmt.Errorf("unknown ClusterCost nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ClusterCostMutation) ResetField(name string) error {
	switch name {
	case clustercost.FieldStartTime:
		m.ResetStartTime()
		return nil
	case clustercost.FieldEndTime:
		m.ResetEndTime()
		return nil
	case clustercost.FieldMinutes:
		m.ResetMinutes()
		return nil
	case clustercost.FieldConnectorID:
		m.ResetConnectorID()
		return nil
	case clustercost.FieldClusterName:
		m.ResetClusterName()
		return nil
	case clustercost.FieldTotalCost:
		m.ResetTotalCost()
		return nil
	case clustercost.FieldCurrency:
		m.ResetCurrency()
		return nil
	case clustercost.FieldAllocationCost:
		m.ResetAllocationCost()
		return nil
	case clustercost.FieldIdleCost:
		m.ResetIdleCost()
		return nil
	case clustercost.FieldManagementCost:
		m.ResetManagementCost()
		return nil
	}
	return fmt.Errorf("unknown ClusterCost field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ClusterCostMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.connector != nil {
		edges = append(edges, clustercost.EdgeConnector)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ClusterCostMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case clustercost.EdgeConnector:
		if id := m.connector; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ClusterCostMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ClusterCostMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ClusterCostMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedconnector {
		edges = append(edges, clustercost.EdgeConnector)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ClusterCostMutation) EdgeCleared(name string) bool {
	switch name {
	case clustercost.EdgeConnector:
		return m.clearedconnector
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ClusterCostMutation) ClearEdge(name string) error {
	switch name {
	case clustercost.EdgeConnector:
		m.ClearConnector()
		return nil
	}
	return fmt.Errorf("unknown ClusterCost unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ClusterCostMutation) ResetEdge(name string) error {
	switch name {
	case clustercost.EdgeConnector:
		m.ResetConnector()
		return nil
	}
	return fmt.Errorf("unknown ClusterCost edge %s", name)
}

// ConnectorMutation represents an operation that mutates the Connector nodes in the graph.
type ConnectorMutation struct {
	config
	op                     Op
	typ                    string
	id                     *oid.ID
	name                   *string
	description            *string
	labels                 *map[string]string
	createTime             *time.Time
	updateTime             *time.Time
	status                 *status.Status
	_type                  *string
	configVersion          *string
	configData             *crypto.Map[string, interface{}]
	enableFinOps           *bool
	finOpsCustomPricing    *types.FinOpsCustomPricing
	category               *string
	clearedFields          map[string]struct{}
	resources              map[oid.ID]struct{}
	removedresources       map[oid.ID]struct{}
	clearedresources       bool
	clusterCosts           map[int]struct{}
	removedclusterCosts    map[int]struct{}
	clearedclusterCosts    bool
	allocationCosts        map[int]struct{}
	removedallocationCosts map[int]struct{}
	clearedallocationCosts bool
	done                   bool
	oldValue               func(context.Context) (*Connector, error)
	predicates             []predicate.Connector
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

// SetStatus sets the "status" field.
func (m *ConnectorMutation) SetStatus(s status.Status) {
	m.status = &s
}

// Status returns the value of the "status" field in the mutation.
func (m *ConnectorMutation) Status() (r status.Status, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldStatus(ctx context.Context) (v status.Status, err error) {
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
func (m *ConnectorMutation) SetConfigData(c crypto.Map[string, interface{}]) {
	m.configData = &c
}

// ConfigData returns the value of the "configData" field in the mutation.
func (m *ConnectorMutation) ConfigData() (r crypto.Map[string, interface{}], exists bool) {
	v := m.configData
	if v == nil {
		return
	}
	return *v, true
}

// OldConfigData returns the old "configData" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldConfigData(ctx context.Context) (v crypto.Map[string, interface{}], err error) {
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

// SetFinOpsCustomPricing sets the "finOpsCustomPricing" field.
func (m *ConnectorMutation) SetFinOpsCustomPricing(tocp types.FinOpsCustomPricing) {
	m.finOpsCustomPricing = &tocp
}

// FinOpsCustomPricing returns the value of the "finOpsCustomPricing" field in the mutation.
func (m *ConnectorMutation) FinOpsCustomPricing() (r types.FinOpsCustomPricing, exists bool) {
	v := m.finOpsCustomPricing
	if v == nil {
		return
	}
	return *v, true
}

// OldFinOpsCustomPricing returns the old "finOpsCustomPricing" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldFinOpsCustomPricing(ctx context.Context) (v types.FinOpsCustomPricing, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFinOpsCustomPricing is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFinOpsCustomPricing requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFinOpsCustomPricing: %w", err)
	}
	return oldValue.FinOpsCustomPricing, nil
}

// ClearFinOpsCustomPricing clears the value of the "finOpsCustomPricing" field.
func (m *ConnectorMutation) ClearFinOpsCustomPricing() {
	m.finOpsCustomPricing = nil
	m.clearedFields[connector.FieldFinOpsCustomPricing] = struct{}{}
}

// FinOpsCustomPricingCleared returns if the "finOpsCustomPricing" field was cleared in this mutation.
func (m *ConnectorMutation) FinOpsCustomPricingCleared() bool {
	_, ok := m.clearedFields[connector.FieldFinOpsCustomPricing]
	return ok
}

// ResetFinOpsCustomPricing resets all changes to the "finOpsCustomPricing" field.
func (m *ConnectorMutation) ResetFinOpsCustomPricing() {
	m.finOpsCustomPricing = nil
	delete(m.clearedFields, connector.FieldFinOpsCustomPricing)
}

// SetCategory sets the "category" field.
func (m *ConnectorMutation) SetCategory(s string) {
	m.category = &s
}

// Category returns the value of the "category" field in the mutation.
func (m *ConnectorMutation) Category() (r string, exists bool) {
	v := m.category
	if v == nil {
		return
	}
	return *v, true
}

// OldCategory returns the old "category" field's value of the Connector entity.
// If the Connector object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ConnectorMutation) OldCategory(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCategory is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCategory requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCategory: %w", err)
	}
	return oldValue.Category, nil
}

// ResetCategory resets all changes to the "category" field.
func (m *ConnectorMutation) ResetCategory() {
	m.category = nil
}

// AddResourceIDs adds the "resources" edge to the ApplicationResource entity by ids.
func (m *ConnectorMutation) AddResourceIDs(ids ...oid.ID) {
	if m.resources == nil {
		m.resources = make(map[oid.ID]struct{})
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
func (m *ConnectorMutation) RemoveResourceIDs(ids ...oid.ID) {
	if m.removedresources == nil {
		m.removedresources = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.resources, ids[i])
		m.removedresources[ids[i]] = struct{}{}
	}
}

// RemovedResources returns the removed IDs of the "resources" edge to the ApplicationResource entity.
func (m *ConnectorMutation) RemovedResourcesIDs() (ids []oid.ID) {
	for id := range m.removedresources {
		ids = append(ids, id)
	}
	return
}

// ResourcesIDs returns the "resources" edge IDs in the mutation.
func (m *ConnectorMutation) ResourcesIDs() (ids []oid.ID) {
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

// AddClusterCostIDs adds the "clusterCosts" edge to the ClusterCost entity by ids.
func (m *ConnectorMutation) AddClusterCostIDs(ids ...int) {
	if m.clusterCosts == nil {
		m.clusterCosts = make(map[int]struct{})
	}
	for i := range ids {
		m.clusterCosts[ids[i]] = struct{}{}
	}
}

// ClearClusterCosts clears the "clusterCosts" edge to the ClusterCost entity.
func (m *ConnectorMutation) ClearClusterCosts() {
	m.clearedclusterCosts = true
}

// ClusterCostsCleared reports if the "clusterCosts" edge to the ClusterCost entity was cleared.
func (m *ConnectorMutation) ClusterCostsCleared() bool {
	return m.clearedclusterCosts
}

// RemoveClusterCostIDs removes the "clusterCosts" edge to the ClusterCost entity by IDs.
func (m *ConnectorMutation) RemoveClusterCostIDs(ids ...int) {
	if m.removedclusterCosts == nil {
		m.removedclusterCosts = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.clusterCosts, ids[i])
		m.removedclusterCosts[ids[i]] = struct{}{}
	}
}

// RemovedClusterCosts returns the removed IDs of the "clusterCosts" edge to the ClusterCost entity.
func (m *ConnectorMutation) RemovedClusterCostsIDs() (ids []int) {
	for id := range m.removedclusterCosts {
		ids = append(ids, id)
	}
	return
}

// ClusterCostsIDs returns the "clusterCosts" edge IDs in the mutation.
func (m *ConnectorMutation) ClusterCostsIDs() (ids []int) {
	for id := range m.clusterCosts {
		ids = append(ids, id)
	}
	return
}

// ResetClusterCosts resets all changes to the "clusterCosts" edge.
func (m *ConnectorMutation) ResetClusterCosts() {
	m.clusterCosts = nil
	m.clearedclusterCosts = false
	m.removedclusterCosts = nil
}

// AddAllocationCostIDs adds the "allocationCosts" edge to the AllocationCost entity by ids.
func (m *ConnectorMutation) AddAllocationCostIDs(ids ...int) {
	if m.allocationCosts == nil {
		m.allocationCosts = make(map[int]struct{})
	}
	for i := range ids {
		m.allocationCosts[ids[i]] = struct{}{}
	}
}

// ClearAllocationCosts clears the "allocationCosts" edge to the AllocationCost entity.
func (m *ConnectorMutation) ClearAllocationCosts() {
	m.clearedallocationCosts = true
}

// AllocationCostsCleared reports if the "allocationCosts" edge to the AllocationCost entity was cleared.
func (m *ConnectorMutation) AllocationCostsCleared() bool {
	return m.clearedallocationCosts
}

// RemoveAllocationCostIDs removes the "allocationCosts" edge to the AllocationCost entity by IDs.
func (m *ConnectorMutation) RemoveAllocationCostIDs(ids ...int) {
	if m.removedallocationCosts == nil {
		m.removedallocationCosts = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.allocationCosts, ids[i])
		m.removedallocationCosts[ids[i]] = struct{}{}
	}
}

// RemovedAllocationCosts returns the removed IDs of the "allocationCosts" edge to the AllocationCost entity.
func (m *ConnectorMutation) RemovedAllocationCostsIDs() (ids []int) {
	for id := range m.removedallocationCosts {
		ids = append(ids, id)
	}
	return
}

// AllocationCostsIDs returns the "allocationCosts" edge IDs in the mutation.
func (m *ConnectorMutation) AllocationCostsIDs() (ids []int) {
	for id := range m.allocationCosts {
		ids = append(ids, id)
	}
	return
}

// ResetAllocationCosts resets all changes to the "allocationCosts" edge.
func (m *ConnectorMutation) ResetAllocationCosts() {
	m.allocationCosts = nil
	m.clearedallocationCosts = false
	m.removedallocationCosts = nil
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
	fields := make([]string, 0, 12)
	if m.name != nil {
		fields = append(fields, connector.FieldName)
	}
	if m.description != nil {
		fields = append(fields, connector.FieldDescription)
	}
	if m.labels != nil {
		fields = append(fields, connector.FieldLabels)
	}
	if m.createTime != nil {
		fields = append(fields, connector.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, connector.FieldUpdateTime)
	}
	if m.status != nil {
		fields = append(fields, connector.FieldStatus)
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
	if m.finOpsCustomPricing != nil {
		fields = append(fields, connector.FieldFinOpsCustomPricing)
	}
	if m.category != nil {
		fields = append(fields, connector.FieldCategory)
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
	case connector.FieldCreateTime:
		return m.CreateTime()
	case connector.FieldUpdateTime:
		return m.UpdateTime()
	case connector.FieldStatus:
		return m.Status()
	case connector.FieldType:
		return m.GetType()
	case connector.FieldConfigVersion:
		return m.ConfigVersion()
	case connector.FieldConfigData:
		return m.ConfigData()
	case connector.FieldEnableFinOps:
		return m.EnableFinOps()
	case connector.FieldFinOpsCustomPricing:
		return m.FinOpsCustomPricing()
	case connector.FieldCategory:
		return m.Category()
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
	case connector.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case connector.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case connector.FieldStatus:
		return m.OldStatus(ctx)
	case connector.FieldType:
		return m.OldType(ctx)
	case connector.FieldConfigVersion:
		return m.OldConfigVersion(ctx)
	case connector.FieldConfigData:
		return m.OldConfigData(ctx)
	case connector.FieldEnableFinOps:
		return m.OldEnableFinOps(ctx)
	case connector.FieldFinOpsCustomPricing:
		return m.OldFinOpsCustomPricing(ctx)
	case connector.FieldCategory:
		return m.OldCategory(ctx)
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
	case connector.FieldStatus:
		v, ok := value.(status.Status)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
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
		v, ok := value.(crypto.Map[string, interface{}])
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
	case connector.FieldFinOpsCustomPricing:
		v, ok := value.(types.FinOpsCustomPricing)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFinOpsCustomPricing(v)
		return nil
	case connector.FieldCategory:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCategory(v)
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
	if m.FieldCleared(connector.FieldFinOpsCustomPricing) {
		fields = append(fields, connector.FieldFinOpsCustomPricing)
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
	case connector.FieldFinOpsCustomPricing:
		m.ClearFinOpsCustomPricing()
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
	case connector.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case connector.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case connector.FieldStatus:
		m.ResetStatus()
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
	case connector.FieldFinOpsCustomPricing:
		m.ResetFinOpsCustomPricing()
		return nil
	case connector.FieldCategory:
		m.ResetCategory()
		return nil
	}
	return fmt.Errorf("unknown Connector field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ConnectorMutation) AddedEdges() []string {
	edges := make([]string, 0, 3)
	if m.resources != nil {
		edges = append(edges, connector.EdgeResources)
	}
	if m.clusterCosts != nil {
		edges = append(edges, connector.EdgeClusterCosts)
	}
	if m.allocationCosts != nil {
		edges = append(edges, connector.EdgeAllocationCosts)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ConnectorMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case connector.EdgeResources:
		ids := make([]ent.Value, 0, len(m.resources))
		for id := range m.resources {
			ids = append(ids, id)
		}
		return ids
	case connector.EdgeClusterCosts:
		ids := make([]ent.Value, 0, len(m.clusterCosts))
		for id := range m.clusterCosts {
			ids = append(ids, id)
		}
		return ids
	case connector.EdgeAllocationCosts:
		ids := make([]ent.Value, 0, len(m.allocationCosts))
		for id := range m.allocationCosts {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ConnectorMutation) RemovedEdges() []string {
	edges := make([]string, 0, 3)
	if m.removedresources != nil {
		edges = append(edges, connector.EdgeResources)
	}
	if m.removedclusterCosts != nil {
		edges = append(edges, connector.EdgeClusterCosts)
	}
	if m.removedallocationCosts != nil {
		edges = append(edges, connector.EdgeAllocationCosts)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ConnectorMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case connector.EdgeResources:
		ids := make([]ent.Value, 0, len(m.removedresources))
		for id := range m.removedresources {
			ids = append(ids, id)
		}
		return ids
	case connector.EdgeClusterCosts:
		ids := make([]ent.Value, 0, len(m.removedclusterCosts))
		for id := range m.removedclusterCosts {
			ids = append(ids, id)
		}
		return ids
	case connector.EdgeAllocationCosts:
		ids := make([]ent.Value, 0, len(m.removedallocationCosts))
		for id := range m.removedallocationCosts {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ConnectorMutation) ClearedEdges() []string {
	edges := make([]string, 0, 3)
	if m.clearedresources {
		edges = append(edges, connector.EdgeResources)
	}
	if m.clearedclusterCosts {
		edges = append(edges, connector.EdgeClusterCosts)
	}
	if m.clearedallocationCosts {
		edges = append(edges, connector.EdgeAllocationCosts)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ConnectorMutation) EdgeCleared(name string) bool {
	switch name {
	case connector.EdgeResources:
		return m.clearedresources
	case connector.EdgeClusterCosts:
		return m.clearedclusterCosts
	case connector.EdgeAllocationCosts:
		return m.clearedallocationCosts
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
	case connector.EdgeResources:
		m.ResetResources()
		return nil
	case connector.EdgeClusterCosts:
		m.ResetClusterCosts()
		return nil
	case connector.EdgeAllocationCosts:
		m.ResetAllocationCosts()
		return nil
	}
	return fmt.Errorf("unknown Connector edge %s", name)
}

// EnvironmentMutation represents an operation that mutates the Environment nodes in the graph.
type EnvironmentMutation struct {
	config
	op               Op
	typ              string
	id               *oid.ID
	name             *string
	description      *string
	labels           *map[string]string
	createTime       *time.Time
	updateTime       *time.Time
	clearedFields    map[string]struct{}
	instances        map[oid.ID]struct{}
	removedinstances map[oid.ID]struct{}
	clearedinstances bool
	revisions        map[oid.ID]struct{}
	removedrevisions map[oid.ID]struct{}
	clearedrevisions bool
	done             bool
	oldValue         func(context.Context) (*Environment, error)
	predicates       []predicate.Environment
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

// AddInstanceIDs adds the "instances" edge to the ApplicationInstance entity by ids.
func (m *EnvironmentMutation) AddInstanceIDs(ids ...oid.ID) {
	if m.instances == nil {
		m.instances = make(map[oid.ID]struct{})
	}
	for i := range ids {
		m.instances[ids[i]] = struct{}{}
	}
}

// ClearInstances clears the "instances" edge to the ApplicationInstance entity.
func (m *EnvironmentMutation) ClearInstances() {
	m.clearedinstances = true
}

// InstancesCleared reports if the "instances" edge to the ApplicationInstance entity was cleared.
func (m *EnvironmentMutation) InstancesCleared() bool {
	return m.clearedinstances
}

// RemoveInstanceIDs removes the "instances" edge to the ApplicationInstance entity by IDs.
func (m *EnvironmentMutation) RemoveInstanceIDs(ids ...oid.ID) {
	if m.removedinstances == nil {
		m.removedinstances = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.instances, ids[i])
		m.removedinstances[ids[i]] = struct{}{}
	}
}

// RemovedInstances returns the removed IDs of the "instances" edge to the ApplicationInstance entity.
func (m *EnvironmentMutation) RemovedInstancesIDs() (ids []oid.ID) {
	for id := range m.removedinstances {
		ids = append(ids, id)
	}
	return
}

// InstancesIDs returns the "instances" edge IDs in the mutation.
func (m *EnvironmentMutation) InstancesIDs() (ids []oid.ID) {
	for id := range m.instances {
		ids = append(ids, id)
	}
	return
}

// ResetInstances resets all changes to the "instances" edge.
func (m *EnvironmentMutation) ResetInstances() {
	m.instances = nil
	m.clearedinstances = false
	m.removedinstances = nil
}

// AddRevisionIDs adds the "revisions" edge to the ApplicationRevision entity by ids.
func (m *EnvironmentMutation) AddRevisionIDs(ids ...oid.ID) {
	if m.revisions == nil {
		m.revisions = make(map[oid.ID]struct{})
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
func (m *EnvironmentMutation) RemoveRevisionIDs(ids ...oid.ID) {
	if m.removedrevisions == nil {
		m.removedrevisions = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.revisions, ids[i])
		m.removedrevisions[ids[i]] = struct{}{}
	}
}

// RemovedRevisions returns the removed IDs of the "revisions" edge to the ApplicationRevision entity.
func (m *EnvironmentMutation) RemovedRevisionsIDs() (ids []oid.ID) {
	for id := range m.removedrevisions {
		ids = append(ids, id)
	}
	return
}

// RevisionsIDs returns the "revisions" edge IDs in the mutation.
func (m *EnvironmentMutation) RevisionsIDs() (ids []oid.ID) {
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
	fields := make([]string, 0, 5)
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
	}
	return fmt.Errorf("unknown Environment field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *EnvironmentMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.instances != nil {
		edges = append(edges, environment.EdgeInstances)
	}
	if m.revisions != nil {
		edges = append(edges, environment.EdgeRevisions)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *EnvironmentMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case environment.EdgeInstances:
		ids := make([]ent.Value, 0, len(m.instances))
		for id := range m.instances {
			ids = append(ids, id)
		}
		return ids
	case environment.EdgeRevisions:
		ids := make([]ent.Value, 0, len(m.revisions))
		for id := range m.revisions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *EnvironmentMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	if m.removedinstances != nil {
		edges = append(edges, environment.EdgeInstances)
	}
	if m.removedrevisions != nil {
		edges = append(edges, environment.EdgeRevisions)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *EnvironmentMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case environment.EdgeInstances:
		ids := make([]ent.Value, 0, len(m.removedinstances))
		for id := range m.removedinstances {
			ids = append(ids, id)
		}
		return ids
	case environment.EdgeRevisions:
		ids := make([]ent.Value, 0, len(m.removedrevisions))
		for id := range m.removedrevisions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *EnvironmentMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedinstances {
		edges = append(edges, environment.EdgeInstances)
	}
	if m.clearedrevisions {
		edges = append(edges, environment.EdgeRevisions)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *EnvironmentMutation) EdgeCleared(name string) bool {
	switch name {
	case environment.EdgeInstances:
		return m.clearedinstances
	case environment.EdgeRevisions:
		return m.clearedrevisions
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
	case environment.EdgeInstances:
		m.ResetInstances()
		return nil
	case environment.EdgeRevisions:
		m.ResetRevisions()
		return nil
	}
	return fmt.Errorf("unknown Environment edge %s", name)
}

// EnvironmentConnectorRelationshipMutation represents an operation that mutates the EnvironmentConnectorRelationship nodes in the graph.
type EnvironmentConnectorRelationshipMutation struct {
	config
	op                 Op
	typ                string
	createTime         *time.Time
	clearedFields      map[string]struct{}
	environment        *oid.ID
	clearedenvironment bool
	connector          *oid.ID
	clearedconnector   bool
	done               bool
	oldValue           func(context.Context) (*EnvironmentConnectorRelationship, error)
	predicates         []predicate.EnvironmentConnectorRelationship
}

var _ ent.Mutation = (*EnvironmentConnectorRelationshipMutation)(nil)

// environmentConnectorRelationshipOption allows management of the mutation configuration using functional options.
type environmentConnectorRelationshipOption func(*EnvironmentConnectorRelationshipMutation)

// newEnvironmentConnectorRelationshipMutation creates new mutation for the EnvironmentConnectorRelationship entity.
func newEnvironmentConnectorRelationshipMutation(c config, op Op, opts ...environmentConnectorRelationshipOption) *EnvironmentConnectorRelationshipMutation {
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

// ResetCreateTime resets all changes to the "createTime" field.
func (m *EnvironmentConnectorRelationshipMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetEnvironmentID sets the "environment_id" field.
func (m *EnvironmentConnectorRelationshipMutation) SetEnvironmentID(o oid.ID) {
	m.environment = &o
}

// EnvironmentID returns the value of the "environment_id" field in the mutation.
func (m *EnvironmentConnectorRelationshipMutation) EnvironmentID() (r oid.ID, exists bool) {
	v := m.environment
	if v == nil {
		return
	}
	return *v, true
}

// ResetEnvironmentID resets all changes to the "environment_id" field.
func (m *EnvironmentConnectorRelationshipMutation) ResetEnvironmentID() {
	m.environment = nil
}

// SetConnectorID sets the "connector_id" field.
func (m *EnvironmentConnectorRelationshipMutation) SetConnectorID(o oid.ID) {
	m.connector = &o
}

// ConnectorID returns the value of the "connector_id" field in the mutation.
func (m *EnvironmentConnectorRelationshipMutation) ConnectorID() (r oid.ID, exists bool) {
	v := m.connector
	if v == nil {
		return
	}
	return *v, true
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
func (m *EnvironmentConnectorRelationshipMutation) EnvironmentIDs() (ids []oid.ID) {
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
func (m *EnvironmentConnectorRelationshipMutation) ConnectorIDs() (ids []oid.ID) {
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
	return nil, errors.New("edge schema EnvironmentConnectorRelationship does not support getting old values")
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
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnvironmentID(v)
		return nil
	case environmentconnectorrelationship.FieldConnectorID:
		v, ok := value.(oid.ID)
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
	op              Op
	typ             string
	id              *string
	status          *string
	statusMessage   *string
	createTime      *time.Time
	updateTime      *time.Time
	description     *string
	icon            *string
	labels          *map[string]string
	source          *string
	clearedFields   map[string]struct{}
	versions        map[oid.ID]struct{}
	removedversions map[oid.ID]struct{}
	clearedversions bool
	done            bool
	oldValue        func(context.Context) (*Module, error)
	predicates      []predicate.Module
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

// SetIcon sets the "icon" field.
func (m *ModuleMutation) SetIcon(s string) {
	m.icon = &s
}

// Icon returns the value of the "icon" field in the mutation.
func (m *ModuleMutation) Icon() (r string, exists bool) {
	v := m.icon
	if v == nil {
		return
	}
	return *v, true
}

// OldIcon returns the old "icon" field's value of the Module entity.
// If the Module object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleMutation) OldIcon(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldIcon is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldIcon requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldIcon: %w", err)
	}
	return oldValue.Icon, nil
}

// ClearIcon clears the value of the "icon" field.
func (m *ModuleMutation) ClearIcon() {
	m.icon = nil
	m.clearedFields[module.FieldIcon] = struct{}{}
}

// IconCleared returns if the "icon" field was cleared in this mutation.
func (m *ModuleMutation) IconCleared() bool {
	_, ok := m.clearedFields[module.FieldIcon]
	return ok
}

// ResetIcon resets all changes to the "icon" field.
func (m *ModuleMutation) ResetIcon() {
	m.icon = nil
	delete(m.clearedFields, module.FieldIcon)
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

// AddVersionIDs adds the "versions" edge to the ModuleVersion entity by ids.
func (m *ModuleMutation) AddVersionIDs(ids ...oid.ID) {
	if m.versions == nil {
		m.versions = make(map[oid.ID]struct{})
	}
	for i := range ids {
		m.versions[ids[i]] = struct{}{}
	}
}

// ClearVersions clears the "versions" edge to the ModuleVersion entity.
func (m *ModuleMutation) ClearVersions() {
	m.clearedversions = true
}

// VersionsCleared reports if the "versions" edge to the ModuleVersion entity was cleared.
func (m *ModuleMutation) VersionsCleared() bool {
	return m.clearedversions
}

// RemoveVersionIDs removes the "versions" edge to the ModuleVersion entity by IDs.
func (m *ModuleMutation) RemoveVersionIDs(ids ...oid.ID) {
	if m.removedversions == nil {
		m.removedversions = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.versions, ids[i])
		m.removedversions[ids[i]] = struct{}{}
	}
}

// RemovedVersions returns the removed IDs of the "versions" edge to the ModuleVersion entity.
func (m *ModuleMutation) RemovedVersionsIDs() (ids []oid.ID) {
	for id := range m.removedversions {
		ids = append(ids, id)
	}
	return
}

// VersionsIDs returns the "versions" edge IDs in the mutation.
func (m *ModuleMutation) VersionsIDs() (ids []oid.ID) {
	for id := range m.versions {
		ids = append(ids, id)
	}
	return
}

// ResetVersions resets all changes to the "versions" edge.
func (m *ModuleMutation) ResetVersions() {
	m.versions = nil
	m.clearedversions = false
	m.removedversions = nil
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
	if m.description != nil {
		fields = append(fields, module.FieldDescription)
	}
	if m.icon != nil {
		fields = append(fields, module.FieldIcon)
	}
	if m.labels != nil {
		fields = append(fields, module.FieldLabels)
	}
	if m.source != nil {
		fields = append(fields, module.FieldSource)
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
	case module.FieldIcon:
		return m.Icon()
	case module.FieldLabels:
		return m.Labels()
	case module.FieldSource:
		return m.Source()
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
	case module.FieldIcon:
		return m.OldIcon(ctx)
	case module.FieldLabels:
		return m.OldLabels(ctx)
	case module.FieldSource:
		return m.OldSource(ctx)
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
	case module.FieldIcon:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIcon(v)
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
	if m.FieldCleared(module.FieldIcon) {
		fields = append(fields, module.FieldIcon)
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
	case module.FieldIcon:
		m.ClearIcon()
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
	case module.FieldIcon:
		m.ResetIcon()
		return nil
	case module.FieldLabels:
		m.ResetLabels()
		return nil
	case module.FieldSource:
		m.ResetSource()
		return nil
	}
	return fmt.Errorf("unknown Module field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ModuleMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.versions != nil {
		edges = append(edges, module.EdgeVersions)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ModuleMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case module.EdgeVersions:
		ids := make([]ent.Value, 0, len(m.versions))
		for id := range m.versions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ModuleMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedversions != nil {
		edges = append(edges, module.EdgeVersions)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ModuleMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case module.EdgeVersions:
		ids := make([]ent.Value, 0, len(m.removedversions))
		for id := range m.removedversions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ModuleMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedversions {
		edges = append(edges, module.EdgeVersions)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ModuleMutation) EdgeCleared(name string) bool {
	switch name {
	case module.EdgeVersions:
		return m.clearedversions
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
	case module.EdgeVersions:
		m.ResetVersions()
		return nil
	}
	return fmt.Errorf("unknown Module edge %s", name)
}

// ModuleVersionMutation represents an operation that mutates the ModuleVersion nodes in the graph.
type ModuleVersionMutation struct {
	config
	op            Op
	typ           string
	id            *oid.ID
	createTime    *time.Time
	updateTime    *time.Time
	version       *string
	source        *string
	schema        **types.ModuleSchema
	clearedFields map[string]struct{}
	module        *string
	clearedmodule bool
	done          bool
	oldValue      func(context.Context) (*ModuleVersion, error)
	predicates    []predicate.ModuleVersion
}

var _ ent.Mutation = (*ModuleVersionMutation)(nil)

// moduleVersionOption allows management of the mutation configuration using functional options.
type moduleVersionOption func(*ModuleVersionMutation)

// newModuleVersionMutation creates new mutation for the ModuleVersion entity.
func newModuleVersionMutation(c config, op Op, opts ...moduleVersionOption) *ModuleVersionMutation {
	m := &ModuleVersionMutation{
		config:        c,
		op:            op,
		typ:           TypeModuleVersion,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withModuleVersionID sets the ID field of the mutation.
func withModuleVersionID(id oid.ID) moduleVersionOption {
	return func(m *ModuleVersionMutation) {
		var (
			err   error
			once  sync.Once
			value *ModuleVersion
		)
		m.oldValue = func(ctx context.Context) (*ModuleVersion, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().ModuleVersion.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withModuleVersion sets the old ModuleVersion of the mutation.
func withModuleVersion(node *ModuleVersion) moduleVersionOption {
	return func(m *ModuleVersionMutation) {
		m.oldValue = func(context.Context) (*ModuleVersion, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ModuleVersionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ModuleVersionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of ModuleVersion entities.
func (m *ModuleVersionMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ModuleVersionMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *ModuleVersionMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().ModuleVersion.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *ModuleVersionMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *ModuleVersionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the ModuleVersion entity.
// If the ModuleVersion object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleVersionMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *ModuleVersionMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *ModuleVersionMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *ModuleVersionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the ModuleVersion entity.
// If the ModuleVersion object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleVersionMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *ModuleVersionMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetModuleID sets the "moduleID" field.
func (m *ModuleVersionMutation) SetModuleID(s string) {
	m.module = &s
}

// ModuleID returns the value of the "moduleID" field in the mutation.
func (m *ModuleVersionMutation) ModuleID() (r string, exists bool) {
	v := m.module
	if v == nil {
		return
	}
	return *v, true
}

// OldModuleID returns the old "moduleID" field's value of the ModuleVersion entity.
// If the ModuleVersion object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleVersionMutation) OldModuleID(ctx context.Context) (v string, err error) {
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

// ResetModuleID resets all changes to the "moduleID" field.
func (m *ModuleVersionMutation) ResetModuleID() {
	m.module = nil
}

// SetVersion sets the "version" field.
func (m *ModuleVersionMutation) SetVersion(s string) {
	m.version = &s
}

// Version returns the value of the "version" field in the mutation.
func (m *ModuleVersionMutation) Version() (r string, exists bool) {
	v := m.version
	if v == nil {
		return
	}
	return *v, true
}

// OldVersion returns the old "version" field's value of the ModuleVersion entity.
// If the ModuleVersion object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleVersionMutation) OldVersion(ctx context.Context) (v string, err error) {
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
func (m *ModuleVersionMutation) ResetVersion() {
	m.version = nil
}

// SetSource sets the "source" field.
func (m *ModuleVersionMutation) SetSource(s string) {
	m.source = &s
}

// Source returns the value of the "source" field in the mutation.
func (m *ModuleVersionMutation) Source() (r string, exists bool) {
	v := m.source
	if v == nil {
		return
	}
	return *v, true
}

// OldSource returns the old "source" field's value of the ModuleVersion entity.
// If the ModuleVersion object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleVersionMutation) OldSource(ctx context.Context) (v string, err error) {
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
func (m *ModuleVersionMutation) ResetSource() {
	m.source = nil
}

// SetSchema sets the "schema" field.
func (m *ModuleVersionMutation) SetSchema(ts *types.ModuleSchema) {
	m.schema = &ts
}

// Schema returns the value of the "schema" field in the mutation.
func (m *ModuleVersionMutation) Schema() (r *types.ModuleSchema, exists bool) {
	v := m.schema
	if v == nil {
		return
	}
	return *v, true
}

// OldSchema returns the old "schema" field's value of the ModuleVersion entity.
// If the ModuleVersion object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ModuleVersionMutation) OldSchema(ctx context.Context) (v *types.ModuleSchema, err error) {
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
func (m *ModuleVersionMutation) ResetSchema() {
	m.schema = nil
}

// ClearModule clears the "module" edge to the Module entity.
func (m *ModuleVersionMutation) ClearModule() {
	m.clearedmodule = true
}

// ModuleCleared reports if the "module" edge to the Module entity was cleared.
func (m *ModuleVersionMutation) ModuleCleared() bool {
	return m.clearedmodule
}

// ModuleIDs returns the "module" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ModuleID instead. It exists only for internal usage by the builders.
func (m *ModuleVersionMutation) ModuleIDs() (ids []string) {
	if id := m.module; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetModule resets all changes to the "module" edge.
func (m *ModuleVersionMutation) ResetModule() {
	m.module = nil
	m.clearedmodule = false
}

// Where appends a list predicates to the ModuleVersionMutation builder.
func (m *ModuleVersionMutation) Where(ps ...predicate.ModuleVersion) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the ModuleVersionMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *ModuleVersionMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.ModuleVersion, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *ModuleVersionMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *ModuleVersionMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (ModuleVersion).
func (m *ModuleVersionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ModuleVersionMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.createTime != nil {
		fields = append(fields, moduleversion.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, moduleversion.FieldUpdateTime)
	}
	if m.module != nil {
		fields = append(fields, moduleversion.FieldModuleID)
	}
	if m.version != nil {
		fields = append(fields, moduleversion.FieldVersion)
	}
	if m.source != nil {
		fields = append(fields, moduleversion.FieldSource)
	}
	if m.schema != nil {
		fields = append(fields, moduleversion.FieldSchema)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ModuleVersionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case moduleversion.FieldCreateTime:
		return m.CreateTime()
	case moduleversion.FieldUpdateTime:
		return m.UpdateTime()
	case moduleversion.FieldModuleID:
		return m.ModuleID()
	case moduleversion.FieldVersion:
		return m.Version()
	case moduleversion.FieldSource:
		return m.Source()
	case moduleversion.FieldSchema:
		return m.Schema()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ModuleVersionMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case moduleversion.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case moduleversion.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case moduleversion.FieldModuleID:
		return m.OldModuleID(ctx)
	case moduleversion.FieldVersion:
		return m.OldVersion(ctx)
	case moduleversion.FieldSource:
		return m.OldSource(ctx)
	case moduleversion.FieldSchema:
		return m.OldSchema(ctx)
	}
	return nil, fmt.Errorf("unknown ModuleVersion field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ModuleVersionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case moduleversion.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case moduleversion.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case moduleversion.FieldModuleID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetModuleID(v)
		return nil
	case moduleversion.FieldVersion:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVersion(v)
		return nil
	case moduleversion.FieldSource:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSource(v)
		return nil
	case moduleversion.FieldSchema:
		v, ok := value.(*types.ModuleSchema)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSchema(v)
		return nil
	}
	return fmt.Errorf("unknown ModuleVersion field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ModuleVersionMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ModuleVersionMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ModuleVersionMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ModuleVersion numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ModuleVersionMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ModuleVersionMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ModuleVersionMutation) ClearField(name string) error {
	return fmt.Errorf("unknown ModuleVersion nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ModuleVersionMutation) ResetField(name string) error {
	switch name {
	case moduleversion.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case moduleversion.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case moduleversion.FieldModuleID:
		m.ResetModuleID()
		return nil
	case moduleversion.FieldVersion:
		m.ResetVersion()
		return nil
	case moduleversion.FieldSource:
		m.ResetSource()
		return nil
	case moduleversion.FieldSchema:
		m.ResetSchema()
		return nil
	}
	return fmt.Errorf("unknown ModuleVersion field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ModuleVersionMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.module != nil {
		edges = append(edges, moduleversion.EdgeModule)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ModuleVersionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case moduleversion.EdgeModule:
		if id := m.module; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ModuleVersionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ModuleVersionMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ModuleVersionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedmodule {
		edges = append(edges, moduleversion.EdgeModule)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ModuleVersionMutation) EdgeCleared(name string) bool {
	switch name {
	case moduleversion.EdgeModule:
		return m.clearedmodule
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ModuleVersionMutation) ClearEdge(name string) error {
	switch name {
	case moduleversion.EdgeModule:
		m.ClearModule()
		return nil
	}
	return fmt.Errorf("unknown ModuleVersion unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ModuleVersionMutation) ResetEdge(name string) error {
	switch name {
	case moduleversion.EdgeModule:
		m.ResetModule()
		return nil
	}
	return fmt.Errorf("unknown ModuleVersion edge %s", name)
}

// PerspectiveMutation represents an operation that mutates the Perspective nodes in the graph.
type PerspectiveMutation struct {
	config
	op                      Op
	typ                     string
	id                      *oid.ID
	createTime              *time.Time
	updateTime              *time.Time
	name                    *string
	startTime               *string
	endTime                 *string
	builtin                 *bool
	allocationQueries       *[]types.QueryCondition
	appendallocationQueries []types.QueryCondition
	clearedFields           map[string]struct{}
	done                    bool
	oldValue                func(context.Context) (*Perspective, error)
	predicates              []predicate.Perspective
}

var _ ent.Mutation = (*PerspectiveMutation)(nil)

// perspectiveOption allows management of the mutation configuration using functional options.
type perspectiveOption func(*PerspectiveMutation)

// newPerspectiveMutation creates new mutation for the Perspective entity.
func newPerspectiveMutation(c config, op Op, opts ...perspectiveOption) *PerspectiveMutation {
	m := &PerspectiveMutation{
		config:        c,
		op:            op,
		typ:           TypePerspective,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withPerspectiveID sets the ID field of the mutation.
func withPerspectiveID(id oid.ID) perspectiveOption {
	return func(m *PerspectiveMutation) {
		var (
			err   error
			once  sync.Once
			value *Perspective
		)
		m.oldValue = func(ctx context.Context) (*Perspective, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Perspective.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withPerspective sets the old Perspective of the mutation.
func withPerspective(node *Perspective) perspectiveOption {
	return func(m *PerspectiveMutation) {
		m.oldValue = func(context.Context) (*Perspective, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m PerspectiveMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m PerspectiveMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Perspective entities.
func (m *PerspectiveMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *PerspectiveMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *PerspectiveMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Perspective.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *PerspectiveMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *PerspectiveMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Perspective entity.
// If the Perspective object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PerspectiveMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *PerspectiveMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *PerspectiveMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *PerspectiveMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Perspective entity.
// If the Perspective object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PerspectiveMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *PerspectiveMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetName sets the "name" field.
func (m *PerspectiveMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *PerspectiveMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Perspective entity.
// If the Perspective object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PerspectiveMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *PerspectiveMutation) ResetName() {
	m.name = nil
}

// SetStartTime sets the "startTime" field.
func (m *PerspectiveMutation) SetStartTime(s string) {
	m.startTime = &s
}

// StartTime returns the value of the "startTime" field in the mutation.
func (m *PerspectiveMutation) StartTime() (r string, exists bool) {
	v := m.startTime
	if v == nil {
		return
	}
	return *v, true
}

// OldStartTime returns the old "startTime" field's value of the Perspective entity.
// If the Perspective object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PerspectiveMutation) OldStartTime(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStartTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStartTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStartTime: %w", err)
	}
	return oldValue.StartTime, nil
}

// ResetStartTime resets all changes to the "startTime" field.
func (m *PerspectiveMutation) ResetStartTime() {
	m.startTime = nil
}

// SetEndTime sets the "endTime" field.
func (m *PerspectiveMutation) SetEndTime(s string) {
	m.endTime = &s
}

// EndTime returns the value of the "endTime" field in the mutation.
func (m *PerspectiveMutation) EndTime() (r string, exists bool) {
	v := m.endTime
	if v == nil {
		return
	}
	return *v, true
}

// OldEndTime returns the old "endTime" field's value of the Perspective entity.
// If the Perspective object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PerspectiveMutation) OldEndTime(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEndTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEndTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEndTime: %w", err)
	}
	return oldValue.EndTime, nil
}

// ResetEndTime resets all changes to the "endTime" field.
func (m *PerspectiveMutation) ResetEndTime() {
	m.endTime = nil
}

// SetBuiltin sets the "builtin" field.
func (m *PerspectiveMutation) SetBuiltin(b bool) {
	m.builtin = &b
}

// Builtin returns the value of the "builtin" field in the mutation.
func (m *PerspectiveMutation) Builtin() (r bool, exists bool) {
	v := m.builtin
	if v == nil {
		return
	}
	return *v, true
}

// OldBuiltin returns the old "builtin" field's value of the Perspective entity.
// If the Perspective object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PerspectiveMutation) OldBuiltin(ctx context.Context) (v bool, err error) {
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
func (m *PerspectiveMutation) ResetBuiltin() {
	m.builtin = nil
}

// SetAllocationQueries sets the "allocationQueries" field.
func (m *PerspectiveMutation) SetAllocationQueries(tc []types.QueryCondition) {
	m.allocationQueries = &tc
	m.appendallocationQueries = nil
}

// AllocationQueries returns the value of the "allocationQueries" field in the mutation.
func (m *PerspectiveMutation) AllocationQueries() (r []types.QueryCondition, exists bool) {
	v := m.allocationQueries
	if v == nil {
		return
	}
	return *v, true
}

// OldAllocationQueries returns the old "allocationQueries" field's value of the Perspective entity.
// If the Perspective object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *PerspectiveMutation) OldAllocationQueries(ctx context.Context) (v []types.QueryCondition, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldAllocationQueries is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldAllocationQueries requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldAllocationQueries: %w", err)
	}
	return oldValue.AllocationQueries, nil
}

// AppendAllocationQueries adds tc to the "allocationQueries" field.
func (m *PerspectiveMutation) AppendAllocationQueries(tc []types.QueryCondition) {
	m.appendallocationQueries = append(m.appendallocationQueries, tc...)
}

// AppendedAllocationQueries returns the list of values that were appended to the "allocationQueries" field in this mutation.
func (m *PerspectiveMutation) AppendedAllocationQueries() ([]types.QueryCondition, bool) {
	if len(m.appendallocationQueries) == 0 {
		return nil, false
	}
	return m.appendallocationQueries, true
}

// ResetAllocationQueries resets all changes to the "allocationQueries" field.
func (m *PerspectiveMutation) ResetAllocationQueries() {
	m.allocationQueries = nil
	m.appendallocationQueries = nil
}

// Where appends a list predicates to the PerspectiveMutation builder.
func (m *PerspectiveMutation) Where(ps ...predicate.Perspective) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the PerspectiveMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *PerspectiveMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Perspective, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *PerspectiveMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *PerspectiveMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Perspective).
func (m *PerspectiveMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *PerspectiveMutation) Fields() []string {
	fields := make([]string, 0, 7)
	if m.createTime != nil {
		fields = append(fields, perspective.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, perspective.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, perspective.FieldName)
	}
	if m.startTime != nil {
		fields = append(fields, perspective.FieldStartTime)
	}
	if m.endTime != nil {
		fields = append(fields, perspective.FieldEndTime)
	}
	if m.builtin != nil {
		fields = append(fields, perspective.FieldBuiltin)
	}
	if m.allocationQueries != nil {
		fields = append(fields, perspective.FieldAllocationQueries)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *PerspectiveMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case perspective.FieldCreateTime:
		return m.CreateTime()
	case perspective.FieldUpdateTime:
		return m.UpdateTime()
	case perspective.FieldName:
		return m.Name()
	case perspective.FieldStartTime:
		return m.StartTime()
	case perspective.FieldEndTime:
		return m.EndTime()
	case perspective.FieldBuiltin:
		return m.Builtin()
	case perspective.FieldAllocationQueries:
		return m.AllocationQueries()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *PerspectiveMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case perspective.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case perspective.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case perspective.FieldName:
		return m.OldName(ctx)
	case perspective.FieldStartTime:
		return m.OldStartTime(ctx)
	case perspective.FieldEndTime:
		return m.OldEndTime(ctx)
	case perspective.FieldBuiltin:
		return m.OldBuiltin(ctx)
	case perspective.FieldAllocationQueries:
		return m.OldAllocationQueries(ctx)
	}
	return nil, fmt.Errorf("unknown Perspective field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *PerspectiveMutation) SetField(name string, value ent.Value) error {
	switch name {
	case perspective.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case perspective.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case perspective.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case perspective.FieldStartTime:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStartTime(v)
		return nil
	case perspective.FieldEndTime:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEndTime(v)
		return nil
	case perspective.FieldBuiltin:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBuiltin(v)
		return nil
	case perspective.FieldAllocationQueries:
		v, ok := value.([]types.QueryCondition)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetAllocationQueries(v)
		return nil
	}
	return fmt.Errorf("unknown Perspective field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *PerspectiveMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *PerspectiveMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *PerspectiveMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Perspective numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *PerspectiveMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *PerspectiveMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *PerspectiveMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Perspective nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *PerspectiveMutation) ResetField(name string) error {
	switch name {
	case perspective.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case perspective.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case perspective.FieldName:
		m.ResetName()
		return nil
	case perspective.FieldStartTime:
		m.ResetStartTime()
		return nil
	case perspective.FieldEndTime:
		m.ResetEndTime()
		return nil
	case perspective.FieldBuiltin:
		m.ResetBuiltin()
		return nil
	case perspective.FieldAllocationQueries:
		m.ResetAllocationQueries()
		return nil
	}
	return fmt.Errorf("unknown Perspective field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *PerspectiveMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *PerspectiveMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *PerspectiveMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *PerspectiveMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *PerspectiveMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *PerspectiveMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *PerspectiveMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Perspective unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *PerspectiveMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Perspective edge %s", name)
}

// ProjectMutation represents an operation that mutates the Project nodes in the graph.
type ProjectMutation struct {
	config
	op                  Op
	typ                 string
	id                  *oid.ID
	name                *string
	description         *string
	labels              *map[string]string
	createTime          *time.Time
	updateTime          *time.Time
	clearedFields       map[string]struct{}
	applications        map[oid.ID]struct{}
	removedapplications map[oid.ID]struct{}
	clearedapplications bool
	secrets             map[oid.ID]struct{}
	removedsecrets      map[oid.ID]struct{}
	clearedsecrets      bool
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
func (m *ProjectMutation) AddApplicationIDs(ids ...oid.ID) {
	if m.applications == nil {
		m.applications = make(map[oid.ID]struct{})
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
func (m *ProjectMutation) RemoveApplicationIDs(ids ...oid.ID) {
	if m.removedapplications == nil {
		m.removedapplications = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.applications, ids[i])
		m.removedapplications[ids[i]] = struct{}{}
	}
}

// RemovedApplications returns the removed IDs of the "applications" edge to the Application entity.
func (m *ProjectMutation) RemovedApplicationsIDs() (ids []oid.ID) {
	for id := range m.removedapplications {
		ids = append(ids, id)
	}
	return
}

// ApplicationsIDs returns the "applications" edge IDs in the mutation.
func (m *ProjectMutation) ApplicationsIDs() (ids []oid.ID) {
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

// AddSecretIDs adds the "secrets" edge to the Secret entity by ids.
func (m *ProjectMutation) AddSecretIDs(ids ...oid.ID) {
	if m.secrets == nil {
		m.secrets = make(map[oid.ID]struct{})
	}
	for i := range ids {
		m.secrets[ids[i]] = struct{}{}
	}
}

// ClearSecrets clears the "secrets" edge to the Secret entity.
func (m *ProjectMutation) ClearSecrets() {
	m.clearedsecrets = true
}

// SecretsCleared reports if the "secrets" edge to the Secret entity was cleared.
func (m *ProjectMutation) SecretsCleared() bool {
	return m.clearedsecrets
}

// RemoveSecretIDs removes the "secrets" edge to the Secret entity by IDs.
func (m *ProjectMutation) RemoveSecretIDs(ids ...oid.ID) {
	if m.removedsecrets == nil {
		m.removedsecrets = make(map[oid.ID]struct{})
	}
	for i := range ids {
		delete(m.secrets, ids[i])
		m.removedsecrets[ids[i]] = struct{}{}
	}
}

// RemovedSecrets returns the removed IDs of the "secrets" edge to the Secret entity.
func (m *ProjectMutation) RemovedSecretsIDs() (ids []oid.ID) {
	for id := range m.removedsecrets {
		ids = append(ids, id)
	}
	return
}

// SecretsIDs returns the "secrets" edge IDs in the mutation.
func (m *ProjectMutation) SecretsIDs() (ids []oid.ID) {
	for id := range m.secrets {
		ids = append(ids, id)
	}
	return
}

// ResetSecrets resets all changes to the "secrets" edge.
func (m *ProjectMutation) ResetSecrets() {
	m.secrets = nil
	m.clearedsecrets = false
	m.removedsecrets = nil
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
	edges := make([]string, 0, 2)
	if m.applications != nil {
		edges = append(edges, project.EdgeApplications)
	}
	if m.secrets != nil {
		edges = append(edges, project.EdgeSecrets)
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
	case project.EdgeSecrets:
		ids := make([]ent.Value, 0, len(m.secrets))
		for id := range m.secrets {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ProjectMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	if m.removedapplications != nil {
		edges = append(edges, project.EdgeApplications)
	}
	if m.removedsecrets != nil {
		edges = append(edges, project.EdgeSecrets)
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
	case project.EdgeSecrets:
		ids := make([]ent.Value, 0, len(m.removedsecrets))
		for id := range m.removedsecrets {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ProjectMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedapplications {
		edges = append(edges, project.EdgeApplications)
	}
	if m.clearedsecrets {
		edges = append(edges, project.EdgeSecrets)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ProjectMutation) EdgeCleared(name string) bool {
	switch name {
	case project.EdgeApplications:
		return m.clearedapplications
	case project.EdgeSecrets:
		return m.clearedsecrets
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
	case project.EdgeSecrets:
		m.ResetSecrets()
		return nil
	}
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
	policies       *types.RolePolicies
	appendpolicies types.RolePolicies
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
func (m *RoleMutation) SetPolicies(tp types.RolePolicies) {
	m.policies = &tp
	m.appendpolicies = nil
}

// Policies returns the value of the "policies" field in the mutation.
func (m *RoleMutation) Policies() (r types.RolePolicies, exists bool) {
	v := m.policies
	if v == nil {
		return
	}
	return *v, true
}

// OldPolicies returns the old "policies" field's value of the Role entity.
// If the Role object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoleMutation) OldPolicies(ctx context.Context) (v types.RolePolicies, err error) {
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

// AppendPolicies adds tp to the "policies" field.
func (m *RoleMutation) AppendPolicies(tp types.RolePolicies) {
	m.appendpolicies = append(m.appendpolicies, tp...)
}

// AppendedPolicies returns the list of values that were appended to the "policies" field in this mutation.
func (m *RoleMutation) AppendedPolicies() (types.RolePolicies, bool) {
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
		v, ok := value.(types.RolePolicies)
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

// SecretMutation represents an operation that mutates the Secret nodes in the graph.
type SecretMutation struct {
	config
	op             Op
	typ            string
	id             *oid.ID
	createTime     *time.Time
	updateTime     *time.Time
	name           *string
	value          *crypto.String
	clearedFields  map[string]struct{}
	project        *oid.ID
	clearedproject bool
	done           bool
	oldValue       func(context.Context) (*Secret, error)
	predicates     []predicate.Secret
}

var _ ent.Mutation = (*SecretMutation)(nil)

// secretOption allows management of the mutation configuration using functional options.
type secretOption func(*SecretMutation)

// newSecretMutation creates new mutation for the Secret entity.
func newSecretMutation(c config, op Op, opts ...secretOption) *SecretMutation {
	m := &SecretMutation{
		config:        c,
		op:            op,
		typ:           TypeSecret,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withSecretID sets the ID field of the mutation.
func withSecretID(id oid.ID) secretOption {
	return func(m *SecretMutation) {
		var (
			err   error
			once  sync.Once
			value *Secret
		)
		m.oldValue = func(ctx context.Context) (*Secret, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Secret.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withSecret sets the old Secret of the mutation.
func withSecret(node *Secret) secretOption {
	return func(m *SecretMutation) {
		m.oldValue = func(context.Context) (*Secret, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SecretMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SecretMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("model: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Secret entities.
func (m *SecretMutation) SetID(id oid.ID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SecretMutation) ID() (id oid.ID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SecretMutation) IDs(ctx context.Context) ([]oid.ID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []oid.ID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Secret.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreateTime sets the "createTime" field.
func (m *SecretMutation) SetCreateTime(t time.Time) {
	m.createTime = &t
}

// CreateTime returns the value of the "createTime" field in the mutation.
func (m *SecretMutation) CreateTime() (r time.Time, exists bool) {
	v := m.createTime
	if v == nil {
		return
	}
	return *v, true
}

// OldCreateTime returns the old "createTime" field's value of the Secret entity.
// If the Secret object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecretMutation) OldCreateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *SecretMutation) ResetCreateTime() {
	m.createTime = nil
}

// SetUpdateTime sets the "updateTime" field.
func (m *SecretMutation) SetUpdateTime(t time.Time) {
	m.updateTime = &t
}

// UpdateTime returns the value of the "updateTime" field in the mutation.
func (m *SecretMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.updateTime
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdateTime returns the old "updateTime" field's value of the Secret entity.
// If the Secret object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecretMutation) OldUpdateTime(ctx context.Context) (v *time.Time, err error) {
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
func (m *SecretMutation) ResetUpdateTime() {
	m.updateTime = nil
}

// SetProjectID sets the "projectID" field.
func (m *SecretMutation) SetProjectID(o oid.ID) {
	m.project = &o
}

// ProjectID returns the value of the "projectID" field in the mutation.
func (m *SecretMutation) ProjectID() (r oid.ID, exists bool) {
	v := m.project
	if v == nil {
		return
	}
	return *v, true
}

// OldProjectID returns the old "projectID" field's value of the Secret entity.
// If the Secret object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecretMutation) OldProjectID(ctx context.Context) (v oid.ID, err error) {
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

// ClearProjectID clears the value of the "projectID" field.
func (m *SecretMutation) ClearProjectID() {
	m.project = nil
	m.clearedFields[secret.FieldProjectID] = struct{}{}
}

// ProjectIDCleared returns if the "projectID" field was cleared in this mutation.
func (m *SecretMutation) ProjectIDCleared() bool {
	_, ok := m.clearedFields[secret.FieldProjectID]
	return ok
}

// ResetProjectID resets all changes to the "projectID" field.
func (m *SecretMutation) ResetProjectID() {
	m.project = nil
	delete(m.clearedFields, secret.FieldProjectID)
}

// SetName sets the "name" field.
func (m *SecretMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *SecretMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Secret entity.
// If the Secret object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecretMutation) OldName(ctx context.Context) (v string, err error) {
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
func (m *SecretMutation) ResetName() {
	m.name = nil
}

// SetValue sets the "value" field.
func (m *SecretMutation) SetValue(c crypto.String) {
	m.value = &c
}

// Value returns the value of the "value" field in the mutation.
func (m *SecretMutation) Value() (r crypto.String, exists bool) {
	v := m.value
	if v == nil {
		return
	}
	return *v, true
}

// OldValue returns the old "value" field's value of the Secret entity.
// If the Secret object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecretMutation) OldValue(ctx context.Context) (v crypto.String, err error) {
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
func (m *SecretMutation) ResetValue() {
	m.value = nil
}

// ClearProject clears the "project" edge to the Project entity.
func (m *SecretMutation) ClearProject() {
	m.clearedproject = true
}

// ProjectCleared reports if the "project" edge to the Project entity was cleared.
func (m *SecretMutation) ProjectCleared() bool {
	return m.ProjectIDCleared() || m.clearedproject
}

// ProjectIDs returns the "project" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// ProjectID instead. It exists only for internal usage by the builders.
func (m *SecretMutation) ProjectIDs() (ids []oid.ID) {
	if id := m.project; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetProject resets all changes to the "project" edge.
func (m *SecretMutation) ResetProject() {
	m.project = nil
	m.clearedproject = false
}

// Where appends a list predicates to the SecretMutation builder.
func (m *SecretMutation) Where(ps ...predicate.Secret) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the SecretMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *SecretMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Secret, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *SecretMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *SecretMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Secret).
func (m *SecretMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *SecretMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.createTime != nil {
		fields = append(fields, secret.FieldCreateTime)
	}
	if m.updateTime != nil {
		fields = append(fields, secret.FieldUpdateTime)
	}
	if m.project != nil {
		fields = append(fields, secret.FieldProjectID)
	}
	if m.name != nil {
		fields = append(fields, secret.FieldName)
	}
	if m.value != nil {
		fields = append(fields, secret.FieldValue)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *SecretMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case secret.FieldCreateTime:
		return m.CreateTime()
	case secret.FieldUpdateTime:
		return m.UpdateTime()
	case secret.FieldProjectID:
		return m.ProjectID()
	case secret.FieldName:
		return m.Name()
	case secret.FieldValue:
		return m.Value()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *SecretMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case secret.FieldCreateTime:
		return m.OldCreateTime(ctx)
	case secret.FieldUpdateTime:
		return m.OldUpdateTime(ctx)
	case secret.FieldProjectID:
		return m.OldProjectID(ctx)
	case secret.FieldName:
		return m.OldName(ctx)
	case secret.FieldValue:
		return m.OldValue(ctx)
	}
	return nil, fmt.Errorf("unknown Secret field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SecretMutation) SetField(name string, value ent.Value) error {
	switch name {
	case secret.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case secret.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case secret.FieldProjectID:
		v, ok := value.(oid.ID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetProjectID(v)
		return nil
	case secret.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case secret.FieldValue:
		v, ok := value.(crypto.String)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetValue(v)
		return nil
	}
	return fmt.Errorf("unknown Secret field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *SecretMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *SecretMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SecretMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Secret numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *SecretMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(secret.FieldProjectID) {
		fields = append(fields, secret.FieldProjectID)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *SecretMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *SecretMutation) ClearField(name string) error {
	switch name {
	case secret.FieldProjectID:
		m.ClearProjectID()
		return nil
	}
	return fmt.Errorf("unknown Secret nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *SecretMutation) ResetField(name string) error {
	switch name {
	case secret.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case secret.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case secret.FieldProjectID:
		m.ResetProjectID()
		return nil
	case secret.FieldName:
		m.ResetName()
		return nil
	case secret.FieldValue:
		m.ResetValue()
		return nil
	}
	return fmt.Errorf("unknown Secret field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *SecretMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.project != nil {
		edges = append(edges, secret.EdgeProject)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *SecretMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case secret.EdgeProject:
		if id := m.project; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *SecretMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *SecretMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *SecretMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedproject {
		edges = append(edges, secret.EdgeProject)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *SecretMutation) EdgeCleared(name string) bool {
	switch name {
	case secret.EdgeProject:
		return m.clearedproject
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *SecretMutation) ClearEdge(name string) error {
	switch name {
	case secret.EdgeProject:
		m.ClearProject()
		return nil
	}
	return fmt.Errorf("unknown Secret unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *SecretMutation) ResetEdge(name string) error {
	switch name {
	case secret.EdgeProject:
		m.ResetProject()
		return nil
	}
	return fmt.Errorf("unknown Secret edge %s", name)
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
	roles         *types.SubjectRoles
	appendroles   types.SubjectRoles
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
func (m *SubjectMutation) SetRoles(tr types.SubjectRoles) {
	m.roles = &tr
	m.appendroles = nil
}

// Roles returns the value of the "roles" field in the mutation.
func (m *SubjectMutation) Roles() (r types.SubjectRoles, exists bool) {
	v := m.roles
	if v == nil {
		return
	}
	return *v, true
}

// OldRoles returns the old "roles" field's value of the Subject entity.
// If the Subject object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SubjectMutation) OldRoles(ctx context.Context) (v types.SubjectRoles, err error) {
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

// AppendRoles adds tr to the "roles" field.
func (m *SubjectMutation) AppendRoles(tr types.SubjectRoles) {
	m.appendroles = append(m.appendroles, tr...)
}

// AppendedRoles returns the list of values that were appended to the "roles" field in this mutation.
func (m *SubjectMutation) AppendedRoles() (types.SubjectRoles, bool) {
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
		v, ok := value.(types.SubjectRoles)
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
