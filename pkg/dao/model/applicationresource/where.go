// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package applicationresource

import (
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/oid"
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
func ApplicationID(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldApplicationID, v))
}

// Module applies equality check predicate on the "module" field. It's identical to ModuleEQ.
func Module(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldModule, v))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldType, v))
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
func ApplicationIDEQ(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldEQ(FieldApplicationID, v))
}

// ApplicationIDNEQ applies the NEQ predicate on the "applicationID" field.
func ApplicationIDNEQ(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNEQ(FieldApplicationID, v))
}

// ApplicationIDIn applies the In predicate on the "applicationID" field.
func ApplicationIDIn(vs ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldIn(FieldApplicationID, vs...))
}

// ApplicationIDNotIn applies the NotIn predicate on the "applicationID" field.
func ApplicationIDNotIn(vs ...oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldNotIn(FieldApplicationID, vs...))
}

// ApplicationIDGT applies the GT predicate on the "applicationID" field.
func ApplicationIDGT(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGT(FieldApplicationID, v))
}

// ApplicationIDGTE applies the GTE predicate on the "applicationID" field.
func ApplicationIDGTE(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldGTE(FieldApplicationID, v))
}

// ApplicationIDLT applies the LT predicate on the "applicationID" field.
func ApplicationIDLT(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLT(FieldApplicationID, v))
}

// ApplicationIDLTE applies the LTE predicate on the "applicationID" field.
func ApplicationIDLTE(v oid.ID) predicate.ApplicationResource {
	return predicate.ApplicationResource(sql.FieldLTE(FieldApplicationID, v))
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
