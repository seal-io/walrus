// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package applicationresource

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ID filters vertices based on their ID field.
func ID(id oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldUpdateTime, v))
}

// InstanceID applies equality check predicate on the "instanceID" field. It's identical to InstanceIDEQ.
func InstanceID(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldInstanceID, v))
}

// ConnectorID applies equality check predicate on the "connectorID" field. It's identical to ConnectorIDEQ.
func ConnectorID(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldConnectorID, v))
}

// CompositionID applies equality check predicate on the "compositionID" field. It's identical to CompositionIDEQ.
func CompositionID(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldCompositionID, v))
}

// Module applies equality check predicate on the "module" field. It's identical to ModuleEQ.
func Module(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldModule, v))
}

// Mode applies equality check predicate on the "mode" field. It's identical to ModeEQ.
func Mode(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldMode, v))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldType, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldName, v))
}

// DeployerType applies equality check predicate on the "deployerType" field. It's identical to DeployerTypeEQ.
func DeployerType(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldDeployerType, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldUpdateTime, v))
}

// InstanceIDEQ applies the EQ predicate on the "instanceID" field.
func InstanceIDEQ(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldInstanceID, v))
}

// InstanceIDNEQ applies the NEQ predicate on the "instanceID" field.
func InstanceIDNEQ(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldInstanceID, v))
}

// InstanceIDIn applies the In predicate on the "instanceID" field.
func InstanceIDIn(vs ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldInstanceID, vs...))
}

// InstanceIDNotIn applies the NotIn predicate on the "instanceID" field.
func InstanceIDNotIn(vs ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldInstanceID, vs...))
}

// InstanceIDGT applies the GT predicate on the "instanceID" field.
func InstanceIDGT(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldInstanceID, v))
}

// InstanceIDGTE applies the GTE predicate on the "instanceID" field.
func InstanceIDGTE(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldInstanceID, v))
}

// InstanceIDLT applies the LT predicate on the "instanceID" field.
func InstanceIDLT(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldInstanceID, v))
}

// InstanceIDLTE applies the LTE predicate on the "instanceID" field.
func InstanceIDLTE(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldInstanceID, v))
}

// InstanceIDContains applies the Contains predicate on the "instanceID" field.
func InstanceIDContains(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContains(FieldInstanceID, vc))
}

// InstanceIDHasPrefix applies the HasPrefix predicate on the "instanceID" field.
func InstanceIDHasPrefix(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldInstanceID, vc))
}

// InstanceIDHasSuffix applies the HasSuffix predicate on the "instanceID" field.
func InstanceIDHasSuffix(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldInstanceID, vc))
}

// InstanceIDEqualFold applies the EqualFold predicate on the "instanceID" field.
func InstanceIDEqualFold(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldInstanceID, vc))
}

// InstanceIDContainsFold applies the ContainsFold predicate on the "instanceID" field.
func InstanceIDContainsFold(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldInstanceID, vc))
}

// ConnectorIDEQ applies the EQ predicate on the "connectorID" field.
func ConnectorIDEQ(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldConnectorID, v))
}

// ConnectorIDNEQ applies the NEQ predicate on the "connectorID" field.
func ConnectorIDNEQ(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldConnectorID, v))
}

// ConnectorIDIn applies the In predicate on the "connectorID" field.
func ConnectorIDIn(vs ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldConnectorID, vs...))
}

// ConnectorIDNotIn applies the NotIn predicate on the "connectorID" field.
func ConnectorIDNotIn(vs ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldConnectorID, vs...))
}

// ConnectorIDGT applies the GT predicate on the "connectorID" field.
func ConnectorIDGT(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldConnectorID, v))
}

// ConnectorIDGTE applies the GTE predicate on the "connectorID" field.
func ConnectorIDGTE(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldConnectorID, v))
}

// ConnectorIDLT applies the LT predicate on the "connectorID" field.
func ConnectorIDLT(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldConnectorID, v))
}

// ConnectorIDLTE applies the LTE predicate on the "connectorID" field.
func ConnectorIDLTE(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldConnectorID, v))
}

// ConnectorIDContains applies the Contains predicate on the "connectorID" field.
func ConnectorIDContains(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContains(FieldConnectorID, vc))
}

// ConnectorIDHasPrefix applies the HasPrefix predicate on the "connectorID" field.
func ConnectorIDHasPrefix(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldConnectorID, vc))
}

// ConnectorIDHasSuffix applies the HasSuffix predicate on the "connectorID" field.
func ConnectorIDHasSuffix(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldConnectorID, vc))
}

// ConnectorIDEqualFold applies the EqualFold predicate on the "connectorID" field.
func ConnectorIDEqualFold(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldConnectorID, vc))
}

// ConnectorIDContainsFold applies the ContainsFold predicate on the "connectorID" field.
func ConnectorIDContainsFold(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldConnectorID, vc))
}

// CompositionIDEQ applies the EQ predicate on the "compositionID" field.
func CompositionIDEQ(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldCompositionID, v))
}

// CompositionIDNEQ applies the NEQ predicate on the "compositionID" field.
func CompositionIDNEQ(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldCompositionID, v))
}

// CompositionIDIn applies the In predicate on the "compositionID" field.
func CompositionIDIn(vs ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldCompositionID, vs...))
}

// CompositionIDNotIn applies the NotIn predicate on the "compositionID" field.
func CompositionIDNotIn(vs ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldCompositionID, vs...))
}

// CompositionIDGT applies the GT predicate on the "compositionID" field.
func CompositionIDGT(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldCompositionID, v))
}

// CompositionIDGTE applies the GTE predicate on the "compositionID" field.
func CompositionIDGTE(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldCompositionID, v))
}

// CompositionIDLT applies the LT predicate on the "compositionID" field.
func CompositionIDLT(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldCompositionID, v))
}

// CompositionIDLTE applies the LTE predicate on the "compositionID" field.
func CompositionIDLTE(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldCompositionID, v))
}

// CompositionIDContains applies the Contains predicate on the "compositionID" field.
func CompositionIDContains(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContains(FieldCompositionID, vc))
}

// CompositionIDHasPrefix applies the HasPrefix predicate on the "compositionID" field.
func CompositionIDHasPrefix(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldCompositionID, vc))
}

// CompositionIDHasSuffix applies the HasSuffix predicate on the "compositionID" field.
func CompositionIDHasSuffix(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldCompositionID, vc))
}

// CompositionIDIsNil applies the IsNil predicate on the "compositionID" field.
func CompositionIDIsNil() predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIsNull(FieldCompositionID))
}

// CompositionIDNotNil applies the NotNil predicate on the "compositionID" field.
func CompositionIDNotNil() predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotNull(FieldCompositionID))
}

// CompositionIDEqualFold applies the EqualFold predicate on the "compositionID" field.
func CompositionIDEqualFold(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldCompositionID, vc))
}

// CompositionIDContainsFold applies the ContainsFold predicate on the "compositionID" field.
func CompositionIDContainsFold(v oid.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldCompositionID, vc))
}

// ModuleEQ applies the EQ predicate on the "module" field.
func ModuleEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldModule, v))
}

// ModuleNEQ applies the NEQ predicate on the "module" field.
func ModuleNEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldModule, v))
}

// ModuleIn applies the In predicate on the "module" field.
func ModuleIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldModule, vs...))
}

// ModuleNotIn applies the NotIn predicate on the "module" field.
func ModuleNotIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldModule, vs...))
}

// ModuleGT applies the GT predicate on the "module" field.
func ModuleGT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldModule, v))
}

// ModuleGTE applies the GTE predicate on the "module" field.
func ModuleGTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldModule, v))
}

// ModuleLT applies the LT predicate on the "module" field.
func ModuleLT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldModule, v))
}

// ModuleLTE applies the LTE predicate on the "module" field.
func ModuleLTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldModule, v))
}

// ModuleContains applies the Contains predicate on the "module" field.
func ModuleContains(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContains(FieldModule, v))
}

// ModuleHasPrefix applies the HasPrefix predicate on the "module" field.
func ModuleHasPrefix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldModule, v))
}

// ModuleHasSuffix applies the HasSuffix predicate on the "module" field.
func ModuleHasSuffix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldModule, v))
}

// ModuleEqualFold applies the EqualFold predicate on the "module" field.
func ModuleEqualFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldModule, v))
}

// ModuleContainsFold applies the ContainsFold predicate on the "module" field.
func ModuleContainsFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldModule, v))
}

// ModeEQ applies the EQ predicate on the "mode" field.
func ModeEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldMode, v))
}

// ModeNEQ applies the NEQ predicate on the "mode" field.
func ModeNEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldMode, v))
}

// ModeIn applies the In predicate on the "mode" field.
func ModeIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldMode, vs...))
}

// ModeNotIn applies the NotIn predicate on the "mode" field.
func ModeNotIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldMode, vs...))
}

// ModeGT applies the GT predicate on the "mode" field.
func ModeGT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldMode, v))
}

// ModeGTE applies the GTE predicate on the "mode" field.
func ModeGTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldMode, v))
}

// ModeLT applies the LT predicate on the "mode" field.
func ModeLT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldMode, v))
}

// ModeLTE applies the LTE predicate on the "mode" field.
func ModeLTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldMode, v))
}

// ModeContains applies the Contains predicate on the "mode" field.
func ModeContains(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContains(FieldMode, v))
}

// ModeHasPrefix applies the HasPrefix predicate on the "mode" field.
func ModeHasPrefix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldMode, v))
}

// ModeHasSuffix applies the HasSuffix predicate on the "mode" field.
func ModeHasSuffix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldMode, v))
}

// ModeEqualFold applies the EqualFold predicate on the "mode" field.
func ModeEqualFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldMode, v))
}

// ModeContainsFold applies the ContainsFold predicate on the "mode" field.
func ModeContainsFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldMode, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldType, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldName, v))
}

// DeployerTypeEQ applies the EQ predicate on the "deployerType" field.
func DeployerTypeEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldDeployerType, v))
}

// DeployerTypeNEQ applies the NEQ predicate on the "deployerType" field.
func DeployerTypeNEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldDeployerType, v))
}

// DeployerTypeIn applies the In predicate on the "deployerType" field.
func DeployerTypeIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldDeployerType, vs...))
}

// DeployerTypeNotIn applies the NotIn predicate on the "deployerType" field.
func DeployerTypeNotIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldDeployerType, vs...))
}

// DeployerTypeGT applies the GT predicate on the "deployerType" field.
func DeployerTypeGT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldDeployerType, v))
}

// DeployerTypeGTE applies the GTE predicate on the "deployerType" field.
func DeployerTypeGTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldDeployerType, v))
}

// DeployerTypeLT applies the LT predicate on the "deployerType" field.
func DeployerTypeLT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldDeployerType, v))
}

// DeployerTypeLTE applies the LTE predicate on the "deployerType" field.
func DeployerTypeLTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldDeployerType, v))
}

// DeployerTypeContains applies the Contains predicate on the "deployerType" field.
func DeployerTypeContains(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContains(FieldDeployerType, v))
}

// DeployerTypeHasPrefix applies the HasPrefix predicate on the "deployerType" field.
func DeployerTypeHasPrefix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldDeployerType, v))
}

// DeployerTypeHasSuffix applies the HasSuffix predicate on the "deployerType" field.
func DeployerTypeHasSuffix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldDeployerType, v))
}

// DeployerTypeEqualFold applies the EqualFold predicate on the "deployerType" field.
func DeployerTypeEqualFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldDeployerType, v))
}

// DeployerTypeContainsFold applies the ContainsFold predicate on the "deployerType" field.
func DeployerTypeContainsFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldDeployerType, v))
}

// StatusIsNil applies the IsNil predicate on the "status" field.
func StatusIsNil() predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIsNull(FieldStatus))
}

// StatusNotNil applies the NotNil predicate on the "status" field.
func StatusNotNil() predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotNull(FieldStatus))
}

// HasInstance applies the HasEdge predicate on the "instance" edge.
func HasInstance() predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, InstanceTable, InstanceColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationInstance
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasInstanceWith applies the HasEdge predicate on the "instance" edge with a given conditions (other predicates).
func HasInstanceWith(preds ...predicate.ApplicationInstance) predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(InstanceInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, InstanceTable, InstanceColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationInstance
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasConnector applies the HasEdge predicate on the "connector" edge.
func HasConnector() predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ConnectorTable, ConnectorColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasConnectorWith applies the HasEdge predicate on the "connector" edge with a given conditions (other predicates).
func HasConnectorWith(preds ...predicate.Connector) predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ConnectorInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ConnectorTable, ConnectorColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasComposition applies the HasEdge predicate on the "composition" edge.
func HasComposition() predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, CompositionTable, CompositionColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationResource
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCompositionWith applies the HasEdge predicate on the "composition" edge with a given conditions (other predicates).
func HasCompositionWith(preds ...predicate.ApplicationResource) predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, CompositionTable, CompositionColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationResource
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasComponents applies the HasEdge predicate on the "components" edge.
func HasComponents() predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ComponentsTable, ComponentsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationResource
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasComponentsWith applies the HasEdge predicate on the "components" edge with a given conditions (other predicates).
func HasComponentsWith(preds ...predicate.ApplicationResource) predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ComponentsTable, ComponentsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationResource
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ApplicationResource) predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ApplicationResource) predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ApplicationResource) predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		p(s.Not())
	})
}
