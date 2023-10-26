// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package templateversion

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/walrus/pkg/dao/model/internal"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/types/object"
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

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldName, v))
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldVersion, v))
}

// Source applies equality check predicate on the "source" field. It's identical to SourceEQ.
func Source(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldSource, v))
}

// ProjectID applies equality check predicate on the "project_id" field. It's identical to ProjectIDEQ.
func ProjectID(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldProjectID, v))
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

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldContainsFold(FieldName, v))
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

// ProjectIDEQ applies the EQ predicate on the "project_id" field.
func ProjectIDEQ(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldEQ(FieldProjectID, v))
}

// ProjectIDNEQ applies the NEQ predicate on the "project_id" field.
func ProjectIDNEQ(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNEQ(FieldProjectID, v))
}

// ProjectIDIn applies the In predicate on the "project_id" field.
func ProjectIDIn(vs ...object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIn(FieldProjectID, vs...))
}

// ProjectIDNotIn applies the NotIn predicate on the "project_id" field.
func ProjectIDNotIn(vs ...object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotIn(FieldProjectID, vs...))
}

// ProjectIDGT applies the GT predicate on the "project_id" field.
func ProjectIDGT(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGT(FieldProjectID, v))
}

// ProjectIDGTE applies the GTE predicate on the "project_id" field.
func ProjectIDGTE(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldGTE(FieldProjectID, v))
}

// ProjectIDLT applies the LT predicate on the "project_id" field.
func ProjectIDLT(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLT(FieldProjectID, v))
}

// ProjectIDLTE applies the LTE predicate on the "project_id" field.
func ProjectIDLTE(v object.ID) predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldLTE(FieldProjectID, v))
}

// ProjectIDContains applies the Contains predicate on the "project_id" field.
func ProjectIDContains(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldContains(FieldProjectID, vc))
}

// ProjectIDHasPrefix applies the HasPrefix predicate on the "project_id" field.
func ProjectIDHasPrefix(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldHasPrefix(FieldProjectID, vc))
}

// ProjectIDHasSuffix applies the HasSuffix predicate on the "project_id" field.
func ProjectIDHasSuffix(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldHasSuffix(FieldProjectID, vc))
}

// ProjectIDIsNil applies the IsNil predicate on the "project_id" field.
func ProjectIDIsNil() predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldIsNull(FieldProjectID))
}

// ProjectIDNotNil applies the NotNil predicate on the "project_id" field.
func ProjectIDNotNil() predicate.TemplateVersion {
	return predicate.TemplateVersion(sql.FieldNotNull(FieldProjectID))
}

// ProjectIDEqualFold applies the EqualFold predicate on the "project_id" field.
func ProjectIDEqualFold(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldEqualFold(FieldProjectID, vc))
}

// ProjectIDContainsFold applies the ContainsFold predicate on the "project_id" field.
func ProjectIDContainsFold(v object.ID) predicate.TemplateVersion {
	vc := string(v)
	return predicate.TemplateVersion(sql.FieldContainsFold(FieldProjectID, vc))
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

// HasServices applies the HasEdge predicate on the "services" edge.
func HasServices() predicate.TemplateVersion {
	return predicate.TemplateVersion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ServicesTable, ServicesColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Service
		step.Edge.Schema = schemaConfig.Service
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServicesWith applies the HasEdge predicate on the "services" edge with a given conditions (other predicates).
func HasServicesWith(preds ...predicate.Service) predicate.TemplateVersion {
	return predicate.TemplateVersion(func(s *sql.Selector) {
		step := newServicesStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Service
		step.Edge.Schema = schemaConfig.Service
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProject applies the HasEdge predicate on the "project" edge.
func HasProject() predicate.TemplateVersion {
	return predicate.TemplateVersion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.TemplateVersion
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectWith applies the HasEdge predicate on the "project" edge with a given conditions (other predicates).
func HasProjectWith(preds ...predicate.Project) predicate.TemplateVersion {
	return predicate.TemplateVersion(func(s *sql.Selector) {
		step := newProjectStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
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
