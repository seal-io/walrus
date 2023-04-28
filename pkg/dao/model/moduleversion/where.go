// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package moduleversion

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ID filters vertices based on their ID field.
func ID(id oid.ID) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id oid.ID) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id oid.ID) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...oid.ID) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...oid.ID) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id oid.ID) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id oid.ID) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id oid.ID) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id oid.ID) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldUpdateTime, v))
}

// ModuleID applies equality check predicate on the "moduleID" field. It's identical to ModuleIDEQ.
func ModuleID(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldModuleID, v))
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldVersion, v))
}

// Source applies equality check predicate on the "source" field. It's identical to SourceEQ.
func Source(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldSource, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLTE(FieldUpdateTime, v))
}

// ModuleIDEQ applies the EQ predicate on the "moduleID" field.
func ModuleIDEQ(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldModuleID, v))
}

// ModuleIDNEQ applies the NEQ predicate on the "moduleID" field.
func ModuleIDNEQ(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNEQ(FieldModuleID, v))
}

// ModuleIDIn applies the In predicate on the "moduleID" field.
func ModuleIDIn(vs ...string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldIn(FieldModuleID, vs...))
}

// ModuleIDNotIn applies the NotIn predicate on the "moduleID" field.
func ModuleIDNotIn(vs ...string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNotIn(FieldModuleID, vs...))
}

// ModuleIDGT applies the GT predicate on the "moduleID" field.
func ModuleIDGT(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGT(FieldModuleID, v))
}

// ModuleIDGTE applies the GTE predicate on the "moduleID" field.
func ModuleIDGTE(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGTE(FieldModuleID, v))
}

// ModuleIDLT applies the LT predicate on the "moduleID" field.
func ModuleIDLT(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLT(FieldModuleID, v))
}

// ModuleIDLTE applies the LTE predicate on the "moduleID" field.
func ModuleIDLTE(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLTE(FieldModuleID, v))
}

// ModuleIDContains applies the Contains predicate on the "moduleID" field.
func ModuleIDContains(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldContains(FieldModuleID, v))
}

// ModuleIDHasPrefix applies the HasPrefix predicate on the "moduleID" field.
func ModuleIDHasPrefix(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldHasPrefix(FieldModuleID, v))
}

// ModuleIDHasSuffix applies the HasSuffix predicate on the "moduleID" field.
func ModuleIDHasSuffix(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldHasSuffix(FieldModuleID, v))
}

// ModuleIDEqualFold applies the EqualFold predicate on the "moduleID" field.
func ModuleIDEqualFold(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEqualFold(FieldModuleID, v))
}

// ModuleIDContainsFold applies the ContainsFold predicate on the "moduleID" field.
func ModuleIDContainsFold(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldContainsFold(FieldModuleID, v))
}

// VersionEQ applies the EQ predicate on the "version" field.
func VersionEQ(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldVersion, v))
}

// VersionNEQ applies the NEQ predicate on the "version" field.
func VersionNEQ(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNEQ(FieldVersion, v))
}

// VersionIn applies the In predicate on the "version" field.
func VersionIn(vs ...string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldIn(FieldVersion, vs...))
}

// VersionNotIn applies the NotIn predicate on the "version" field.
func VersionNotIn(vs ...string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNotIn(FieldVersion, vs...))
}

// VersionGT applies the GT predicate on the "version" field.
func VersionGT(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGT(FieldVersion, v))
}

// VersionGTE applies the GTE predicate on the "version" field.
func VersionGTE(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGTE(FieldVersion, v))
}

// VersionLT applies the LT predicate on the "version" field.
func VersionLT(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLT(FieldVersion, v))
}

// VersionLTE applies the LTE predicate on the "version" field.
func VersionLTE(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLTE(FieldVersion, v))
}

// VersionContains applies the Contains predicate on the "version" field.
func VersionContains(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldContains(FieldVersion, v))
}

// VersionHasPrefix applies the HasPrefix predicate on the "version" field.
func VersionHasPrefix(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldHasPrefix(FieldVersion, v))
}

// VersionHasSuffix applies the HasSuffix predicate on the "version" field.
func VersionHasSuffix(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldHasSuffix(FieldVersion, v))
}

// VersionEqualFold applies the EqualFold predicate on the "version" field.
func VersionEqualFold(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEqualFold(FieldVersion, v))
}

// VersionContainsFold applies the ContainsFold predicate on the "version" field.
func VersionContainsFold(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldContainsFold(FieldVersion, v))
}

// SourceEQ applies the EQ predicate on the "source" field.
func SourceEQ(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEQ(FieldSource, v))
}

// SourceNEQ applies the NEQ predicate on the "source" field.
func SourceNEQ(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNEQ(FieldSource, v))
}

// SourceIn applies the In predicate on the "source" field.
func SourceIn(vs ...string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldIn(FieldSource, vs...))
}

// SourceNotIn applies the NotIn predicate on the "source" field.
func SourceNotIn(vs ...string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldNotIn(FieldSource, vs...))
}

// SourceGT applies the GT predicate on the "source" field.
func SourceGT(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGT(FieldSource, v))
}

// SourceGTE applies the GTE predicate on the "source" field.
func SourceGTE(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldGTE(FieldSource, v))
}

// SourceLT applies the LT predicate on the "source" field.
func SourceLT(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLT(FieldSource, v))
}

// SourceLTE applies the LTE predicate on the "source" field.
func SourceLTE(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldLTE(FieldSource, v))
}

// SourceContains applies the Contains predicate on the "source" field.
func SourceContains(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldContains(FieldSource, v))
}

// SourceHasPrefix applies the HasPrefix predicate on the "source" field.
func SourceHasPrefix(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldHasPrefix(FieldSource, v))
}

// SourceHasSuffix applies the HasSuffix predicate on the "source" field.
func SourceHasSuffix(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldHasSuffix(FieldSource, v))
}

// SourceEqualFold applies the EqualFold predicate on the "source" field.
func SourceEqualFold(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldEqualFold(FieldSource, v))
}

// SourceContainsFold applies the ContainsFold predicate on the "source" field.
func SourceContainsFold(v string) predicate.ModuleVersion {
	return predicate.ModuleVersion(sql.FieldContainsFold(FieldSource, v))
}

// HasModule applies the HasEdge predicate on the "module" edge.
func HasModule() predicate.ModuleVersion {
	return predicate.ModuleVersion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ModuleTable, ModuleColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Module
		step.Edge.Schema = schemaConfig.ModuleVersion
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasModuleWith applies the HasEdge predicate on the "module" edge with a given conditions (other predicates).
func HasModuleWith(preds ...predicate.Module) predicate.ModuleVersion {
	return predicate.ModuleVersion(func(s *sql.Selector) {
		step := newModuleStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Module
		step.Edge.Schema = schemaConfig.ModuleVersion
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ModuleVersion) predicate.ModuleVersion {
	return predicate.ModuleVersion(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ModuleVersion) predicate.ModuleVersion {
	return predicate.ModuleVersion(func(s *sql.Selector) {
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
func Not(p predicate.ModuleVersion) predicate.ModuleVersion {
	return predicate.ModuleVersion(func(s *sql.Selector) {
		p(s.Not())
	})
}
