// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package application

import (
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// ID filters vertices based on their ID field.
func ID(id oid.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id oid.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id oid.ID) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...oid.ID) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...oid.ID) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id oid.ID) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id oid.ID) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id oid.ID) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id oid.ID) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldUpdateTime, v))
}

// ProjectID applies equality check predicate on the "projectID" field. It's identical to ProjectIDEQ.
func ProjectID(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldProjectID, v))
}

// EnvironmentID applies equality check predicate on the "environmentID" field. It's identical to EnvironmentIDEQ.
func EnvironmentID(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldEnvironmentID, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldUpdateTime, v))
}

// ProjectIDEQ applies the EQ predicate on the "projectID" field.
func ProjectIDEQ(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldProjectID, v))
}

// ProjectIDNEQ applies the NEQ predicate on the "projectID" field.
func ProjectIDNEQ(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldProjectID, v))
}

// ProjectIDIn applies the In predicate on the "projectID" field.
func ProjectIDIn(vs ...oid.ID) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldProjectID, vs...))
}

// ProjectIDNotIn applies the NotIn predicate on the "projectID" field.
func ProjectIDNotIn(vs ...oid.ID) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldProjectID, vs...))
}

// ProjectIDGT applies the GT predicate on the "projectID" field.
func ProjectIDGT(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldProjectID, v))
}

// ProjectIDGTE applies the GTE predicate on the "projectID" field.
func ProjectIDGTE(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldProjectID, v))
}

// ProjectIDLT applies the LT predicate on the "projectID" field.
func ProjectIDLT(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldProjectID, v))
}

// ProjectIDLTE applies the LTE predicate on the "projectID" field.
func ProjectIDLTE(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldProjectID, v))
}

// EnvironmentIDEQ applies the EQ predicate on the "environmentID" field.
func EnvironmentIDEQ(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldEnvironmentID, v))
}

// EnvironmentIDNEQ applies the NEQ predicate on the "environmentID" field.
func EnvironmentIDNEQ(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldEnvironmentID, v))
}

// EnvironmentIDIn applies the In predicate on the "environmentID" field.
func EnvironmentIDIn(vs ...oid.ID) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDNotIn applies the NotIn predicate on the "environmentID" field.
func EnvironmentIDNotIn(vs ...oid.ID) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDGT applies the GT predicate on the "environmentID" field.
func EnvironmentIDGT(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldEnvironmentID, v))
}

// EnvironmentIDGTE applies the GTE predicate on the "environmentID" field.
func EnvironmentIDGTE(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldEnvironmentID, v))
}

// EnvironmentIDLT applies the LT predicate on the "environmentID" field.
func EnvironmentIDLT(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldEnvironmentID, v))
}

// EnvironmentIDLTE applies the LTE predicate on the "environmentID" field.
func EnvironmentIDLTE(v oid.ID) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldEnvironmentID, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Application) predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Application) predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
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
func Not(p predicate.Application) predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		p(s.Not())
	})
}
