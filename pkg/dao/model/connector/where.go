// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package connector

import (
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// ID filters vertices based on their ID field.
func ID(id oid.ID) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id oid.ID) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id oid.ID) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...oid.ID) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...oid.ID) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id oid.ID) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id oid.ID) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id oid.ID) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id oid.ID) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldID, id))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldStatus, v))
}

// StatusMessage applies equality check predicate on the "statusMessage" field. It's identical to StatusMessageEQ.
func StatusMessage(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldStatusMessage, v))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldUpdateTime, v))
}

// Driver applies equality check predicate on the "driver" field. It's identical to DriverEQ.
func Driver(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldDriver, v))
}

// ConfigVersion applies equality check predicate on the "configVersion" field. It's identical to ConfigVersionEQ.
func ConfigVersion(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldConfigVersion, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldStatus, v))
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldStatus, v))
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldStatus, v))
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldStatus, v))
}

// StatusIsNil applies the IsNil predicate on the "status" field.
func StatusIsNil() predicate.Connector {
	return predicate.Connector(sql.FieldIsNull(FieldStatus))
}

// StatusNotNil applies the NotNil predicate on the "status" field.
func StatusNotNil() predicate.Connector {
	return predicate.Connector(sql.FieldNotNull(FieldStatus))
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldStatus, v))
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldStatus, v))
}

// StatusMessageEQ applies the EQ predicate on the "statusMessage" field.
func StatusMessageEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldStatusMessage, v))
}

// StatusMessageNEQ applies the NEQ predicate on the "statusMessage" field.
func StatusMessageNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldStatusMessage, v))
}

// StatusMessageIn applies the In predicate on the "statusMessage" field.
func StatusMessageIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldStatusMessage, vs...))
}

// StatusMessageNotIn applies the NotIn predicate on the "statusMessage" field.
func StatusMessageNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldStatusMessage, vs...))
}

// StatusMessageGT applies the GT predicate on the "statusMessage" field.
func StatusMessageGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldStatusMessage, v))
}

// StatusMessageGTE applies the GTE predicate on the "statusMessage" field.
func StatusMessageGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldStatusMessage, v))
}

// StatusMessageLT applies the LT predicate on the "statusMessage" field.
func StatusMessageLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldStatusMessage, v))
}

// StatusMessageLTE applies the LTE predicate on the "statusMessage" field.
func StatusMessageLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldStatusMessage, v))
}

// StatusMessageContains applies the Contains predicate on the "statusMessage" field.
func StatusMessageContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldStatusMessage, v))
}

// StatusMessageHasPrefix applies the HasPrefix predicate on the "statusMessage" field.
func StatusMessageHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldStatusMessage, v))
}

// StatusMessageHasSuffix applies the HasSuffix predicate on the "statusMessage" field.
func StatusMessageHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldStatusMessage, v))
}

// StatusMessageIsNil applies the IsNil predicate on the "statusMessage" field.
func StatusMessageIsNil() predicate.Connector {
	return predicate.Connector(sql.FieldIsNull(FieldStatusMessage))
}

// StatusMessageNotNil applies the NotNil predicate on the "statusMessage" field.
func StatusMessageNotNil() predicate.Connector {
	return predicate.Connector(sql.FieldNotNull(FieldStatusMessage))
}

// StatusMessageEqualFold applies the EqualFold predicate on the "statusMessage" field.
func StatusMessageEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldStatusMessage, v))
}

// StatusMessageContainsFold applies the ContainsFold predicate on the "statusMessage" field.
func StatusMessageContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldStatusMessage, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldUpdateTime, v))
}

// DriverEQ applies the EQ predicate on the "driver" field.
func DriverEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldDriver, v))
}

// DriverNEQ applies the NEQ predicate on the "driver" field.
func DriverNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldDriver, v))
}

// DriverIn applies the In predicate on the "driver" field.
func DriverIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldDriver, vs...))
}

// DriverNotIn applies the NotIn predicate on the "driver" field.
func DriverNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldDriver, vs...))
}

// DriverGT applies the GT predicate on the "driver" field.
func DriverGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldDriver, v))
}

// DriverGTE applies the GTE predicate on the "driver" field.
func DriverGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldDriver, v))
}

// DriverLT applies the LT predicate on the "driver" field.
func DriverLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldDriver, v))
}

// DriverLTE applies the LTE predicate on the "driver" field.
func DriverLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldDriver, v))
}

// DriverContains applies the Contains predicate on the "driver" field.
func DriverContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldDriver, v))
}

// DriverHasPrefix applies the HasPrefix predicate on the "driver" field.
func DriverHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldDriver, v))
}

// DriverHasSuffix applies the HasSuffix predicate on the "driver" field.
func DriverHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldDriver, v))
}

// DriverEqualFold applies the EqualFold predicate on the "driver" field.
func DriverEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldDriver, v))
}

// DriverContainsFold applies the ContainsFold predicate on the "driver" field.
func DriverContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldDriver, v))
}

// ConfigVersionEQ applies the EQ predicate on the "configVersion" field.
func ConfigVersionEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldConfigVersion, v))
}

// ConfigVersionNEQ applies the NEQ predicate on the "configVersion" field.
func ConfigVersionNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldConfigVersion, v))
}

// ConfigVersionIn applies the In predicate on the "configVersion" field.
func ConfigVersionIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldConfigVersion, vs...))
}

// ConfigVersionNotIn applies the NotIn predicate on the "configVersion" field.
func ConfigVersionNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldConfigVersion, vs...))
}

// ConfigVersionGT applies the GT predicate on the "configVersion" field.
func ConfigVersionGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldConfigVersion, v))
}

// ConfigVersionGTE applies the GTE predicate on the "configVersion" field.
func ConfigVersionGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldConfigVersion, v))
}

// ConfigVersionLT applies the LT predicate on the "configVersion" field.
func ConfigVersionLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldConfigVersion, v))
}

// ConfigVersionLTE applies the LTE predicate on the "configVersion" field.
func ConfigVersionLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldConfigVersion, v))
}

// ConfigVersionContains applies the Contains predicate on the "configVersion" field.
func ConfigVersionContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldConfigVersion, v))
}

// ConfigVersionHasPrefix applies the HasPrefix predicate on the "configVersion" field.
func ConfigVersionHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldConfigVersion, v))
}

// ConfigVersionHasSuffix applies the HasSuffix predicate on the "configVersion" field.
func ConfigVersionHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldConfigVersion, v))
}

// ConfigVersionEqualFold applies the EqualFold predicate on the "configVersion" field.
func ConfigVersionEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldConfigVersion, v))
}

// ConfigVersionContainsFold applies the ContainsFold predicate on the "configVersion" field.
func ConfigVersionContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldConfigVersion, v))
}

// ConfigDataIsNil applies the IsNil predicate on the "configData" field.
func ConfigDataIsNil() predicate.Connector {
	return predicate.Connector(sql.FieldIsNull(FieldConfigData))
}

// ConfigDataNotNil applies the NotNil predicate on the "configData" field.
func ConfigDataNotNil() predicate.Connector {
	return predicate.Connector(sql.FieldNotNull(FieldConfigData))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Connector) predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Connector) predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
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
func Not(p predicate.Connector) predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
		p(s.Not())
	})
}
