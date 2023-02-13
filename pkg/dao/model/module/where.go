// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package module

import (
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// ID filters vertices based on their ID field.
func ID(id oid.ID) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id oid.ID) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id oid.ID) predicate.Module {
	return predicate.Module(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...oid.ID) predicate.Module {
	return predicate.Module(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...oid.ID) predicate.Module {
	return predicate.Module(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id oid.ID) predicate.Module {
	return predicate.Module(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id oid.ID) predicate.Module {
	return predicate.Module(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id oid.ID) predicate.Module {
	return predicate.Module(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id oid.ID) predicate.Module {
	return predicate.Module(sql.FieldLTE(FieldID, id))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldStatus, v))
}

// StatusMessage applies equality check predicate on the "statusMessage" field. It's identical to StatusMessageEQ.
func StatusMessage(v string) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldStatusMessage, v))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldUpdateTime, v))
}

// Source applies equality check predicate on the "source" field. It's identical to SourceEQ.
func Source(v string) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldSource, v))
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v string) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldVersion, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.Module {
	return predicate.Module(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.Module {
	return predicate.Module(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.Module {
	return predicate.Module(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.Module {
	return predicate.Module(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.Module {
	return predicate.Module(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.Module {
	return predicate.Module(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.Module {
	return predicate.Module(sql.FieldLTE(FieldStatus, v))
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.Module {
	return predicate.Module(sql.FieldContains(FieldStatus, v))
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.Module {
	return predicate.Module(sql.FieldHasPrefix(FieldStatus, v))
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.Module {
	return predicate.Module(sql.FieldHasSuffix(FieldStatus, v))
}

// StatusIsNil applies the IsNil predicate on the "status" field.
func StatusIsNil() predicate.Module {
	return predicate.Module(sql.FieldIsNull(FieldStatus))
}

// StatusNotNil applies the NotNil predicate on the "status" field.
func StatusNotNil() predicate.Module {
	return predicate.Module(sql.FieldNotNull(FieldStatus))
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.Module {
	return predicate.Module(sql.FieldEqualFold(FieldStatus, v))
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.Module {
	return predicate.Module(sql.FieldContainsFold(FieldStatus, v))
}

// StatusMessageEQ applies the EQ predicate on the "statusMessage" field.
func StatusMessageEQ(v string) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldStatusMessage, v))
}

// StatusMessageNEQ applies the NEQ predicate on the "statusMessage" field.
func StatusMessageNEQ(v string) predicate.Module {
	return predicate.Module(sql.FieldNEQ(FieldStatusMessage, v))
}

// StatusMessageIn applies the In predicate on the "statusMessage" field.
func StatusMessageIn(vs ...string) predicate.Module {
	return predicate.Module(sql.FieldIn(FieldStatusMessage, vs...))
}

// StatusMessageNotIn applies the NotIn predicate on the "statusMessage" field.
func StatusMessageNotIn(vs ...string) predicate.Module {
	return predicate.Module(sql.FieldNotIn(FieldStatusMessage, vs...))
}

// StatusMessageGT applies the GT predicate on the "statusMessage" field.
func StatusMessageGT(v string) predicate.Module {
	return predicate.Module(sql.FieldGT(FieldStatusMessage, v))
}

// StatusMessageGTE applies the GTE predicate on the "statusMessage" field.
func StatusMessageGTE(v string) predicate.Module {
	return predicate.Module(sql.FieldGTE(FieldStatusMessage, v))
}

// StatusMessageLT applies the LT predicate on the "statusMessage" field.
func StatusMessageLT(v string) predicate.Module {
	return predicate.Module(sql.FieldLT(FieldStatusMessage, v))
}

// StatusMessageLTE applies the LTE predicate on the "statusMessage" field.
func StatusMessageLTE(v string) predicate.Module {
	return predicate.Module(sql.FieldLTE(FieldStatusMessage, v))
}

// StatusMessageContains applies the Contains predicate on the "statusMessage" field.
func StatusMessageContains(v string) predicate.Module {
	return predicate.Module(sql.FieldContains(FieldStatusMessage, v))
}

// StatusMessageHasPrefix applies the HasPrefix predicate on the "statusMessage" field.
func StatusMessageHasPrefix(v string) predicate.Module {
	return predicate.Module(sql.FieldHasPrefix(FieldStatusMessage, v))
}

// StatusMessageHasSuffix applies the HasSuffix predicate on the "statusMessage" field.
func StatusMessageHasSuffix(v string) predicate.Module {
	return predicate.Module(sql.FieldHasSuffix(FieldStatusMessage, v))
}

// StatusMessageIsNil applies the IsNil predicate on the "statusMessage" field.
func StatusMessageIsNil() predicate.Module {
	return predicate.Module(sql.FieldIsNull(FieldStatusMessage))
}

// StatusMessageNotNil applies the NotNil predicate on the "statusMessage" field.
func StatusMessageNotNil() predicate.Module {
	return predicate.Module(sql.FieldNotNull(FieldStatusMessage))
}

// StatusMessageEqualFold applies the EqualFold predicate on the "statusMessage" field.
func StatusMessageEqualFold(v string) predicate.Module {
	return predicate.Module(sql.FieldEqualFold(FieldStatusMessage, v))
}

// StatusMessageContainsFold applies the ContainsFold predicate on the "statusMessage" field.
func StatusMessageContainsFold(v string) predicate.Module {
	return predicate.Module(sql.FieldContainsFold(FieldStatusMessage, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.Module {
	return predicate.Module(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Module {
	return predicate.Module(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.Module {
	return predicate.Module(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Module {
	return predicate.Module(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.Module {
	return predicate.Module(sql.FieldLTE(FieldUpdateTime, v))
}

// SourceEQ applies the EQ predicate on the "source" field.
func SourceEQ(v string) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldSource, v))
}

// SourceNEQ applies the NEQ predicate on the "source" field.
func SourceNEQ(v string) predicate.Module {
	return predicate.Module(sql.FieldNEQ(FieldSource, v))
}

// SourceIn applies the In predicate on the "source" field.
func SourceIn(vs ...string) predicate.Module {
	return predicate.Module(sql.FieldIn(FieldSource, vs...))
}

// SourceNotIn applies the NotIn predicate on the "source" field.
func SourceNotIn(vs ...string) predicate.Module {
	return predicate.Module(sql.FieldNotIn(FieldSource, vs...))
}

// SourceGT applies the GT predicate on the "source" field.
func SourceGT(v string) predicate.Module {
	return predicate.Module(sql.FieldGT(FieldSource, v))
}

// SourceGTE applies the GTE predicate on the "source" field.
func SourceGTE(v string) predicate.Module {
	return predicate.Module(sql.FieldGTE(FieldSource, v))
}

// SourceLT applies the LT predicate on the "source" field.
func SourceLT(v string) predicate.Module {
	return predicate.Module(sql.FieldLT(FieldSource, v))
}

// SourceLTE applies the LTE predicate on the "source" field.
func SourceLTE(v string) predicate.Module {
	return predicate.Module(sql.FieldLTE(FieldSource, v))
}

// SourceContains applies the Contains predicate on the "source" field.
func SourceContains(v string) predicate.Module {
	return predicate.Module(sql.FieldContains(FieldSource, v))
}

// SourceHasPrefix applies the HasPrefix predicate on the "source" field.
func SourceHasPrefix(v string) predicate.Module {
	return predicate.Module(sql.FieldHasPrefix(FieldSource, v))
}

// SourceHasSuffix applies the HasSuffix predicate on the "source" field.
func SourceHasSuffix(v string) predicate.Module {
	return predicate.Module(sql.FieldHasSuffix(FieldSource, v))
}

// SourceEqualFold applies the EqualFold predicate on the "source" field.
func SourceEqualFold(v string) predicate.Module {
	return predicate.Module(sql.FieldEqualFold(FieldSource, v))
}

// SourceContainsFold applies the ContainsFold predicate on the "source" field.
func SourceContainsFold(v string) predicate.Module {
	return predicate.Module(sql.FieldContainsFold(FieldSource, v))
}

// VersionEQ applies the EQ predicate on the "version" field.
func VersionEQ(v string) predicate.Module {
	return predicate.Module(sql.FieldEQ(FieldVersion, v))
}

// VersionNEQ applies the NEQ predicate on the "version" field.
func VersionNEQ(v string) predicate.Module {
	return predicate.Module(sql.FieldNEQ(FieldVersion, v))
}

// VersionIn applies the In predicate on the "version" field.
func VersionIn(vs ...string) predicate.Module {
	return predicate.Module(sql.FieldIn(FieldVersion, vs...))
}

// VersionNotIn applies the NotIn predicate on the "version" field.
func VersionNotIn(vs ...string) predicate.Module {
	return predicate.Module(sql.FieldNotIn(FieldVersion, vs...))
}

// VersionGT applies the GT predicate on the "version" field.
func VersionGT(v string) predicate.Module {
	return predicate.Module(sql.FieldGT(FieldVersion, v))
}

// VersionGTE applies the GTE predicate on the "version" field.
func VersionGTE(v string) predicate.Module {
	return predicate.Module(sql.FieldGTE(FieldVersion, v))
}

// VersionLT applies the LT predicate on the "version" field.
func VersionLT(v string) predicate.Module {
	return predicate.Module(sql.FieldLT(FieldVersion, v))
}

// VersionLTE applies the LTE predicate on the "version" field.
func VersionLTE(v string) predicate.Module {
	return predicate.Module(sql.FieldLTE(FieldVersion, v))
}

// VersionContains applies the Contains predicate on the "version" field.
func VersionContains(v string) predicate.Module {
	return predicate.Module(sql.FieldContains(FieldVersion, v))
}

// VersionHasPrefix applies the HasPrefix predicate on the "version" field.
func VersionHasPrefix(v string) predicate.Module {
	return predicate.Module(sql.FieldHasPrefix(FieldVersion, v))
}

// VersionHasSuffix applies the HasSuffix predicate on the "version" field.
func VersionHasSuffix(v string) predicate.Module {
	return predicate.Module(sql.FieldHasSuffix(FieldVersion, v))
}

// VersionEqualFold applies the EqualFold predicate on the "version" field.
func VersionEqualFold(v string) predicate.Module {
	return predicate.Module(sql.FieldEqualFold(FieldVersion, v))
}

// VersionContainsFold applies the ContainsFold predicate on the "version" field.
func VersionContainsFold(v string) predicate.Module {
	return predicate.Module(sql.FieldContainsFold(FieldVersion, v))
}

// InputSchemaIsNil applies the IsNil predicate on the "inputSchema" field.
func InputSchemaIsNil() predicate.Module {
	return predicate.Module(sql.FieldIsNull(FieldInputSchema))
}

// InputSchemaNotNil applies the NotNil predicate on the "inputSchema" field.
func InputSchemaNotNil() predicate.Module {
	return predicate.Module(sql.FieldNotNull(FieldInputSchema))
}

// OutputSchemaIsNil applies the IsNil predicate on the "outputSchema" field.
func OutputSchemaIsNil() predicate.Module {
	return predicate.Module(sql.FieldIsNull(FieldOutputSchema))
}

// OutputSchemaNotNil applies the NotNil predicate on the "outputSchema" field.
func OutputSchemaNotNil() predicate.Module {
	return predicate.Module(sql.FieldNotNull(FieldOutputSchema))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Module) predicate.Module {
	return predicate.Module(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Module) predicate.Module {
	return predicate.Module(func(s *sql.Selector) {
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
func Not(p predicate.Module) predicate.Module {
	return predicate.Module(func(s *sql.Selector) {
		p(s.Not())
	})
}
