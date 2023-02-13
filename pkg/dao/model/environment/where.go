// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package environment

import (
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// ID filters vertices based on their ID field.
func ID(id oid.ID) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id oid.ID) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id oid.ID) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...oid.ID) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...oid.ID) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id oid.ID) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id oid.ID) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id oid.ID) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id oid.ID) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldUpdateTime, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.Environment {
	return predicate.Environment(sql.FieldLTE(FieldUpdateTime, v))
}

// ConnectorIDsIsNil applies the IsNil predicate on the "connectorIDs" field.
func ConnectorIDsIsNil() predicate.Environment {
	return predicate.Environment(sql.FieldIsNull(FieldConnectorIDs))
}

// ConnectorIDsNotNil applies the NotNil predicate on the "connectorIDs" field.
func ConnectorIDsNotNil() predicate.Environment {
	return predicate.Environment(sql.FieldNotNull(FieldConnectorIDs))
}

// VariablesIsNil applies the IsNil predicate on the "variables" field.
func VariablesIsNil() predicate.Environment {
	return predicate.Environment(sql.FieldIsNull(FieldVariables))
}

// VariablesNotNil applies the NotNil predicate on the "variables" field.
func VariablesNotNil() predicate.Environment {
	return predicate.Environment(sql.FieldNotNull(FieldVariables))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Environment) predicate.Environment {
	return predicate.Environment(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Environment) predicate.Environment {
	return predicate.Environment(func(s *sql.Selector) {
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
func Not(p predicate.Environment) predicate.Environment {
	return predicate.Environment(func(s *sql.Selector) {
		p(s.Not())
	})
}
