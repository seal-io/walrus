// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package secret

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ID filters vertices based on their ID field.
func ID(id oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldUpdateTime, v))
}

// ProjectID applies equality check predicate on the "projectID" field. It's identical to ProjectIDEQ.
func ProjectID(v oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldProjectID, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldName, v))
}

// Value applies equality check predicate on the "value" field. It's identical to ValueEQ.
func Value(v crypto.String) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldValue, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.Secret {
	return predicate.Secret(sql.FieldLTE(FieldUpdateTime, v))
}

// ProjectIDEQ applies the EQ predicate on the "projectID" field.
func ProjectIDEQ(v oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldProjectID, v))
}

// ProjectIDNEQ applies the NEQ predicate on the "projectID" field.
func ProjectIDNEQ(v oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldNEQ(FieldProjectID, v))
}

// ProjectIDIn applies the In predicate on the "projectID" field.
func ProjectIDIn(vs ...oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldIn(FieldProjectID, vs...))
}

// ProjectIDNotIn applies the NotIn predicate on the "projectID" field.
func ProjectIDNotIn(vs ...oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldNotIn(FieldProjectID, vs...))
}

// ProjectIDGT applies the GT predicate on the "projectID" field.
func ProjectIDGT(v oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldGT(FieldProjectID, v))
}

// ProjectIDGTE applies the GTE predicate on the "projectID" field.
func ProjectIDGTE(v oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldGTE(FieldProjectID, v))
}

// ProjectIDLT applies the LT predicate on the "projectID" field.
func ProjectIDLT(v oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldLT(FieldProjectID, v))
}

// ProjectIDLTE applies the LTE predicate on the "projectID" field.
func ProjectIDLTE(v oid.ID) predicate.Secret {
	return predicate.Secret(sql.FieldLTE(FieldProjectID, v))
}

// ProjectIDContains applies the Contains predicate on the "projectID" field.
func ProjectIDContains(v oid.ID) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldContains(FieldProjectID, vc))
}

// ProjectIDHasPrefix applies the HasPrefix predicate on the "projectID" field.
func ProjectIDHasPrefix(v oid.ID) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldHasPrefix(FieldProjectID, vc))
}

// ProjectIDHasSuffix applies the HasSuffix predicate on the "projectID" field.
func ProjectIDHasSuffix(v oid.ID) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldHasSuffix(FieldProjectID, vc))
}

// ProjectIDIsNil applies the IsNil predicate on the "projectID" field.
func ProjectIDIsNil() predicate.Secret {
	return predicate.Secret(sql.FieldIsNull(FieldProjectID))
}

// ProjectIDNotNil applies the NotNil predicate on the "projectID" field.
func ProjectIDNotNil() predicate.Secret {
	return predicate.Secret(sql.FieldNotNull(FieldProjectID))
}

// ProjectIDEqualFold applies the EqualFold predicate on the "projectID" field.
func ProjectIDEqualFold(v oid.ID) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldEqualFold(FieldProjectID, vc))
}

// ProjectIDContainsFold applies the ContainsFold predicate on the "projectID" field.
func ProjectIDContainsFold(v oid.ID) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldContainsFold(FieldProjectID, vc))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Secret {
	return predicate.Secret(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Secret {
	return predicate.Secret(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Secret {
	return predicate.Secret(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Secret {
	return predicate.Secret(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Secret {
	return predicate.Secret(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Secret {
	return predicate.Secret(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Secret {
	return predicate.Secret(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Secret {
	return predicate.Secret(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Secret {
	return predicate.Secret(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Secret {
	return predicate.Secret(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Secret {
	return predicate.Secret(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Secret {
	return predicate.Secret(sql.FieldContainsFold(FieldName, v))
}

// ValueEQ applies the EQ predicate on the "value" field.
func ValueEQ(v crypto.String) predicate.Secret {
	return predicate.Secret(sql.FieldEQ(FieldValue, v))
}

// ValueNEQ applies the NEQ predicate on the "value" field.
func ValueNEQ(v crypto.String) predicate.Secret {
	return predicate.Secret(sql.FieldNEQ(FieldValue, v))
}

// ValueIn applies the In predicate on the "value" field.
func ValueIn(vs ...crypto.String) predicate.Secret {
	return predicate.Secret(sql.FieldIn(FieldValue, vs...))
}

// ValueNotIn applies the NotIn predicate on the "value" field.
func ValueNotIn(vs ...crypto.String) predicate.Secret {
	return predicate.Secret(sql.FieldNotIn(FieldValue, vs...))
}

// ValueGT applies the GT predicate on the "value" field.
func ValueGT(v crypto.String) predicate.Secret {
	return predicate.Secret(sql.FieldGT(FieldValue, v))
}

// ValueGTE applies the GTE predicate on the "value" field.
func ValueGTE(v crypto.String) predicate.Secret {
	return predicate.Secret(sql.FieldGTE(FieldValue, v))
}

// ValueLT applies the LT predicate on the "value" field.
func ValueLT(v crypto.String) predicate.Secret {
	return predicate.Secret(sql.FieldLT(FieldValue, v))
}

// ValueLTE applies the LTE predicate on the "value" field.
func ValueLTE(v crypto.String) predicate.Secret {
	return predicate.Secret(sql.FieldLTE(FieldValue, v))
}

// ValueContains applies the Contains predicate on the "value" field.
func ValueContains(v crypto.String) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldContains(FieldValue, vc))
}

// ValueHasPrefix applies the HasPrefix predicate on the "value" field.
func ValueHasPrefix(v crypto.String) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldHasPrefix(FieldValue, vc))
}

// ValueHasSuffix applies the HasSuffix predicate on the "value" field.
func ValueHasSuffix(v crypto.String) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldHasSuffix(FieldValue, vc))
}

// ValueEqualFold applies the EqualFold predicate on the "value" field.
func ValueEqualFold(v crypto.String) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldEqualFold(FieldValue, vc))
}

// ValueContainsFold applies the ContainsFold predicate on the "value" field.
func ValueContainsFold(v crypto.String) predicate.Secret {
	vc := string(v)
	return predicate.Secret(sql.FieldContainsFold(FieldValue, vc))
}

// HasProject applies the HasEdge predicate on the "project" edge.
func HasProject() predicate.Secret {
	return predicate.Secret(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.Secret
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectWith applies the HasEdge predicate on the "project" edge with a given conditions (other predicates).
func HasProjectWith(preds ...predicate.Project) predicate.Secret {
	return predicate.Secret(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.Secret
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Secret) predicate.Secret {
	return predicate.Secret(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Secret) predicate.Secret {
	return predicate.Secret(func(s *sql.Selector) {
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
func Not(p predicate.Secret) predicate.Secret {
	return predicate.Secret(func(s *sql.Selector) {
		p(s.Not())
	})
}
