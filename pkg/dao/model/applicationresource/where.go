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
	"github.com/seal-io/seal/pkg/dao/types"
)

// ID filters vertices based on their ID field.
func ID(id types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldID, id))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldStatus, v))
}

// StatusMessage applies equality check predicate on the "statusMessage" field. It's identical to StatusMessageEQ.
func StatusMessage(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldStatusMessage, v))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldUpdateTime, v))
}

// ApplicationID applies equality check predicate on the "applicationID" field. It's identical to ApplicationIDEQ.
func ApplicationID(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldApplicationID, v))
}

// ConnectorID applies equality check predicate on the "connectorID" field. It's identical to ConnectorIDEQ.
func ConnectorID(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldConnectorID, v))
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

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldStatus, v))
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContains(FieldStatus, v))
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldStatus, v))
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldStatus, v))
}

// StatusIsNil applies the IsNil predicate on the "status" field.
func StatusIsNil() predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIsNull(FieldStatus))
}

// StatusNotNil applies the NotNil predicate on the "status" field.
func StatusNotNil() predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotNull(FieldStatus))
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldStatus, v))
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldStatus, v))
}

// StatusMessageEQ applies the EQ predicate on the "statusMessage" field.
func StatusMessageEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldStatusMessage, v))
}

// StatusMessageNEQ applies the NEQ predicate on the "statusMessage" field.
func StatusMessageNEQ(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldStatusMessage, v))
}

// StatusMessageIn applies the In predicate on the "statusMessage" field.
func StatusMessageIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldStatusMessage, vs...))
}

// StatusMessageNotIn applies the NotIn predicate on the "statusMessage" field.
func StatusMessageNotIn(vs ...string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldStatusMessage, vs...))
}

// StatusMessageGT applies the GT predicate on the "statusMessage" field.
func StatusMessageGT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldStatusMessage, v))
}

// StatusMessageGTE applies the GTE predicate on the "statusMessage" field.
func StatusMessageGTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldStatusMessage, v))
}

// StatusMessageLT applies the LT predicate on the "statusMessage" field.
func StatusMessageLT(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldStatusMessage, v))
}

// StatusMessageLTE applies the LTE predicate on the "statusMessage" field.
func StatusMessageLTE(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldStatusMessage, v))
}

// StatusMessageContains applies the Contains predicate on the "statusMessage" field.
func StatusMessageContains(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContains(FieldStatusMessage, v))
}

// StatusMessageHasPrefix applies the HasPrefix predicate on the "statusMessage" field.
func StatusMessageHasPrefix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldStatusMessage, v))
}

// StatusMessageHasSuffix applies the HasSuffix predicate on the "statusMessage" field.
func StatusMessageHasSuffix(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldStatusMessage, v))
}

// StatusMessageIsNil applies the IsNil predicate on the "statusMessage" field.
func StatusMessageIsNil() predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIsNull(FieldStatusMessage))
}

// StatusMessageNotNil applies the NotNil predicate on the "statusMessage" field.
func StatusMessageNotNil() predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotNull(FieldStatusMessage))
}

// StatusMessageEqualFold applies the EqualFold predicate on the "statusMessage" field.
func StatusMessageEqualFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldStatusMessage, v))
}

// StatusMessageContainsFold applies the ContainsFold predicate on the "statusMessage" field.
func StatusMessageContainsFold(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldStatusMessage, v))
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

// ApplicationIDEQ applies the EQ predicate on the "applicationID" field.
func ApplicationIDEQ(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldApplicationID, v))
}

// ApplicationIDNEQ applies the NEQ predicate on the "applicationID" field.
func ApplicationIDNEQ(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldApplicationID, v))
}

// ApplicationIDIn applies the In predicate on the "applicationID" field.
func ApplicationIDIn(vs ...types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldApplicationID, vs...))
}

// ApplicationIDNotIn applies the NotIn predicate on the "applicationID" field.
func ApplicationIDNotIn(vs ...types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldApplicationID, vs...))
}

// ApplicationIDGT applies the GT predicate on the "applicationID" field.
func ApplicationIDGT(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldApplicationID, v))
}

// ApplicationIDGTE applies the GTE predicate on the "applicationID" field.
func ApplicationIDGTE(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldApplicationID, v))
}

// ApplicationIDLT applies the LT predicate on the "applicationID" field.
func ApplicationIDLT(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldApplicationID, v))
}

// ApplicationIDLTE applies the LTE predicate on the "applicationID" field.
func ApplicationIDLTE(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldApplicationID, v))
}

// ApplicationIDContains applies the Contains predicate on the "applicationID" field.
func ApplicationIDContains(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContains(FieldApplicationID, vc))
}

// ApplicationIDHasPrefix applies the HasPrefix predicate on the "applicationID" field.
func ApplicationIDHasPrefix(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldApplicationID, vc))
}

// ApplicationIDHasSuffix applies the HasSuffix predicate on the "applicationID" field.
func ApplicationIDHasSuffix(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldApplicationID, vc))
}

// ApplicationIDEqualFold applies the EqualFold predicate on the "applicationID" field.
func ApplicationIDEqualFold(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldApplicationID, vc))
}

// ApplicationIDContainsFold applies the ContainsFold predicate on the "applicationID" field.
func ApplicationIDContainsFold(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldApplicationID, vc))
}

// ConnectorIDEQ applies the EQ predicate on the "connectorID" field.
func ConnectorIDEQ(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldConnectorID, v))
}

// ConnectorIDNEQ applies the NEQ predicate on the "connectorID" field.
func ConnectorIDNEQ(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldConnectorID, v))
}

// ConnectorIDIn applies the In predicate on the "connectorID" field.
func ConnectorIDIn(vs ...types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldConnectorID, vs...))
}

// ConnectorIDNotIn applies the NotIn predicate on the "connectorID" field.
func ConnectorIDNotIn(vs ...types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldConnectorID, vs...))
}

// ConnectorIDGT applies the GT predicate on the "connectorID" field.
func ConnectorIDGT(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldConnectorID, v))
}

// ConnectorIDGTE applies the GTE predicate on the "connectorID" field.
func ConnectorIDGTE(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldConnectorID, v))
}

// ConnectorIDLT applies the LT predicate on the "connectorID" field.
func ConnectorIDLT(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldConnectorID, v))
}

// ConnectorIDLTE applies the LTE predicate on the "connectorID" field.
func ConnectorIDLTE(v types.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldConnectorID, v))
}

// ConnectorIDContains applies the Contains predicate on the "connectorID" field.
func ConnectorIDContains(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContains(FieldConnectorID, vc))
}

// ConnectorIDHasPrefix applies the HasPrefix predicate on the "connectorID" field.
func ConnectorIDHasPrefix(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasPrefix(FieldConnectorID, vc))
}

// ConnectorIDHasSuffix applies the HasSuffix predicate on the "connectorID" field.
func ConnectorIDHasSuffix(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldHasSuffix(FieldConnectorID, vc))
}

// ConnectorIDEqualFold applies the EqualFold predicate on the "connectorID" field.
func ConnectorIDEqualFold(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldEqualFold(FieldConnectorID, vc))
}

// ConnectorIDContainsFold applies the ContainsFold predicate on the "connectorID" field.
func ConnectorIDContainsFold(v types.ID) predicate.ApplicationResource {
	vc := string(v)
	return predicate.ApplicationResource(sql.FieldContainsFold(FieldConnectorID, vc))
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

// HasApplication applies the HasEdge predicate on the "application" edge.
func HasApplication() predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ApplicationTable, ApplicationColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasApplicationWith applies the HasEdge predicate on the "application" edge with a given conditions (other predicates).
func HasApplicationWith(preds ...predicate.Application) predicate.ApplicationResource {
	return predicate.ApplicationResource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ApplicationInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ApplicationTable, ApplicationColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Application
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
