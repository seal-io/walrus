// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package templateversion

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// ID filters vertices based on their ID field.
func ID(id object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldUpdateTime, v))
}

// TemplateID applies equality check predicate on the "template_id" field. It's identical to TemplateIDEQ.
func TemplateID(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldTemplateID, v))
}

// TemplateName applies equality check predicate on the "template_name" field. It's identical to TemplateNameEQ.
func TemplateName(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldTemplateName, v))
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldVersion, v))
}

// Source applies equality check predicate on the "source" field. It's identical to SourceEQ.
func Source(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldSource, v))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLTE(FieldUpdateTime, v))
}

// TemplateIDEQ applies the EQ predicate on the "template_id" field.
func TemplateIDEQ(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldTemplateID, v))
}

// TemplateIDNEQ applies the NEQ predicate on the "template_id" field.
func TemplateIDNEQ(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNEQ(FieldTemplateID, v))
}

// TemplateIDIn applies the In predicate on the "template_id" field.
func TemplateIDIn(vs ...object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIn(FieldTemplateID, vs...))
}

// TemplateIDNotIn applies the NotIn predicate on the "template_id" field.
func TemplateIDNotIn(vs ...object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotIn(FieldTemplateID, vs...))
}

// TemplateIDGT applies the GT predicate on the "template_id" field.
func TemplateIDGT(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGT(FieldTemplateID, v))
}

// TemplateIDGTE applies the GTE predicate on the "template_id" field.
func TemplateIDGTE(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGTE(FieldTemplateID, v))
}

// TemplateIDLT applies the LT predicate on the "template_id" field.
func TemplateIDLT(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLT(FieldTemplateID, v))
}

// TemplateIDLTE applies the LTE predicate on the "template_id" field.
func TemplateIDLTE(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLTE(FieldTemplateID, v))
}

// TemplateIDContains applies the Contains predicate on the "template_id" field.
func TemplateIDContains(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldContains(FieldTemplateID, vc))
}

// TemplateIDHasPrefix applies the HasPrefix predicate on the "template_id" field.
func TemplateIDHasPrefix(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldHasPrefix(FieldTemplateID, vc))
}

// TemplateIDHasSuffix applies the HasSuffix predicate on the "template_id" field.
func TemplateIDHasSuffix(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldHasSuffix(FieldTemplateID, vc))
}

// TemplateIDEqualFold applies the EqualFold predicate on the "template_id" field.
func TemplateIDEqualFold(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldEqualFold(FieldTemplateID, vc))
}

// TemplateIDContainsFold applies the ContainsFold predicate on the "template_id" field.
func TemplateIDContainsFold(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldContainsFold(FieldTemplateID, vc))
}

// TemplateNameEQ applies the EQ predicate on the "template_name" field.
func TemplateNameEQ(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldTemplateName, v))
}

// TemplateNameNEQ applies the NEQ predicate on the "template_name" field.
func TemplateNameNEQ(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNEQ(FieldTemplateName, v))
}

// TemplateNameIn applies the In predicate on the "template_name" field.
func TemplateNameIn(vs ...string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIn(FieldTemplateName, vs...))
}

// TemplateNameNotIn applies the NotIn predicate on the "template_name" field.
func TemplateNameNotIn(vs ...string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotIn(FieldTemplateName, vs...))
}

// TemplateNameGT applies the GT predicate on the "template_name" field.
func TemplateNameGT(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGT(FieldTemplateName, v))
}

// TemplateNameGTE applies the GTE predicate on the "template_name" field.
func TemplateNameGTE(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGTE(FieldTemplateName, v))
}

// TemplateNameLT applies the LT predicate on the "template_name" field.
func TemplateNameLT(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLT(FieldTemplateName, v))
}

// TemplateNameLTE applies the LTE predicate on the "template_name" field.
func TemplateNameLTE(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLTE(FieldTemplateName, v))
}

// TemplateNameContains applies the Contains predicate on the "template_name" field.
func TemplateNameContains(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldContains(FieldTemplateName, v))
}

// TemplateNameHasPrefix applies the HasPrefix predicate on the "template_name" field.
func TemplateNameHasPrefix(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldHasPrefix(FieldTemplateName, v))
}

// TemplateNameHasSuffix applies the HasSuffix predicate on the "template_name" field.
func TemplateNameHasSuffix(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldHasSuffix(FieldTemplateName, v))
}

// TemplateNameEqualFold applies the EqualFold predicate on the "template_name" field.
func TemplateNameEqualFold(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEqualFold(FieldTemplateName, v))
}

// TemplateNameContainsFold applies the ContainsFold predicate on the "template_name" field.
func TemplateNameContainsFold(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldContainsFold(FieldTemplateName, v))
}

// VersionEQ applies the EQ predicate on the "version" field.
func VersionEQ(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldVersion, v))
}

// VersionNEQ applies the NEQ predicate on the "version" field.
func VersionNEQ(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNEQ(FieldVersion, v))
}

// VersionIn applies the In predicate on the "version" field.
func VersionIn(vs ...string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIn(FieldVersion, vs...))
}

// VersionNotIn applies the NotIn predicate on the "version" field.
func VersionNotIn(vs ...string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotIn(FieldVersion, vs...))
}

// VersionGT applies the GT predicate on the "version" field.
func VersionGT(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGT(FieldVersion, v))
}

// VersionGTE applies the GTE predicate on the "version" field.
func VersionGTE(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGTE(FieldVersion, v))
}

// VersionLT applies the LT predicate on the "version" field.
func VersionLT(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLT(FieldVersion, v))
}

// VersionLTE applies the LTE predicate on the "version" field.
func VersionLTE(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLTE(FieldVersion, v))
}

// VersionContains applies the Contains predicate on the "version" field.
func VersionContains(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldContains(FieldVersion, v))
}

// VersionHasPrefix applies the HasPrefix predicate on the "version" field.
func VersionHasPrefix(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldHasPrefix(FieldVersion, v))
}

// VersionHasSuffix applies the HasSuffix predicate on the "version" field.
func VersionHasSuffix(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldHasSuffix(FieldVersion, v))
}

// VersionEqualFold applies the EqualFold predicate on the "version" field.
func VersionEqualFold(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEqualFold(FieldVersion, v))
}

// VersionContainsFold applies the ContainsFold predicate on the "version" field.
func VersionContainsFold(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldContainsFold(FieldVersion, v))
}

// SourceEQ applies the EQ predicate on the "source" field.
func SourceEQ(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldSource, v))
}

// SourceNEQ applies the NEQ predicate on the "source" field.
func SourceNEQ(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNEQ(FieldSource, v))
}

// SourceIn applies the In predicate on the "source" field.
func SourceIn(vs ...string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIn(FieldSource, vs...))
}

// SourceNotIn applies the NotIn predicate on the "source" field.
func SourceNotIn(vs ...string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotIn(FieldSource, vs...))
}

// SourceGT applies the GT predicate on the "source" field.
func SourceGT(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGT(FieldSource, v))
}

// SourceGTE applies the GTE predicate on the "source" field.
func SourceGTE(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGTE(FieldSource, v))
}

// SourceLT applies the LT predicate on the "source" field.
func SourceLT(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLT(FieldSource, v))
}

// SourceLTE applies the LTE predicate on the "source" field.
func SourceLTE(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLTE(FieldSource, v))
}

// SourceContains applies the Contains predicate on the "source" field.
func SourceContains(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldContains(FieldSource, v))
}

// SourceHasPrefix applies the HasPrefix predicate on the "source" field.
func SourceHasPrefix(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldHasPrefix(FieldSource, v))
}

// SourceHasSuffix applies the HasSuffix predicate on the "source" field.
func SourceHasSuffix(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldHasSuffix(FieldSource, v))
}

// SourceEqualFold applies the EqualFold predicate on the "source" field.
func SourceEqualFold(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEqualFold(FieldSource, v))
}

// SourceContainsFold applies the ContainsFold predicate on the "source" field.
func SourceContainsFold(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldContainsFold(FieldSource, v))
}

// HasTemplate applies the HasEdge predicate on the "template" edge.
func HasTemplate() predicate.TemplateVersion {
	return predicate.TemplateVersion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, TemplateTable, TemplateColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Template
		step.Edge.Schema = schemaConfig.TemplateVersion
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTemplateWith applies the HasEdge predicate on the "template" edge with a given conditions (other predicates).
func HasTemplateWith(preds ...predicate.Template) predicate.TemplateVersion {
	return predicate.TemplateVersion(func(s *sql.Selector) {
		step := newTemplateStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Template
		step.Edge.Schema = schemaConfig.TemplateVersion
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.TemplateVersion) predicate.TemplateVersion {
	return predicate.TemplateVersion(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.TemplateVersion) predicate.TemplateVersion {
	return predicate.TemplateVersion(func(s *sql.Selector) {
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
func Not(p predicate.TemplateVersion) predicate.TemplateVersion {
	return predicate.TemplateVersion(func(s *sql.Selector) {
		p(s.Not())
	})
}
