// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package resource

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/walrus/pkg/dao/model/internal"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
)

// ID filters vertices based on their ID field.
func ID(id object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldDescription, v))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldUpdateTime, v))
}

// ProjectID applies equality check predicate on the "project_id" field. It's identical to ProjectIDEQ.
func ProjectID(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldProjectID, v))
}

// EnvironmentID applies equality check predicate on the "environment_id" field. It's identical to EnvironmentIDEQ.
func EnvironmentID(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldEnvironmentID, v))
}

// TemplateID applies equality check predicate on the "template_id" field. It's identical to TemplateIDEQ.
func TemplateID(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldTemplateID, v))
}

// Attributes applies equality check predicate on the "attributes" field. It's identical to AttributesEQ.
func Attributes(v property.Values) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldAttributes, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Resource {
	return predicate.Resource(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Resource {
	return predicate.Resource(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Resource {
	return predicate.Resource(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Resource {
	return predicate.Resource(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Resource {
	return predicate.Resource(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Resource {
	return predicate.Resource(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Resource {
	return predicate.Resource(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Resource {
	return predicate.Resource(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Resource {
	return predicate.Resource(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Resource {
	return predicate.Resource(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Resource {
	return predicate.Resource(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Resource {
	return predicate.Resource(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Resource {
	return predicate.Resource(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Resource {
	return predicate.Resource(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Resource {
	return predicate.Resource(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Resource {
	return predicate.Resource(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Resource {
	return predicate.Resource(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Resource {
	return predicate.Resource(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Resource {
	return predicate.Resource(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Resource {
	return predicate.Resource(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Resource {
	return predicate.Resource(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Resource {
	return predicate.Resource(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Resource {
	return predicate.Resource(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Resource {
	return predicate.Resource(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Resource {
	return predicate.Resource(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Resource {
	return predicate.Resource(sql.FieldContainsFold(FieldDescription, v))
}

// LabelsIsNil applies the IsNil predicate on the "labels" field.
func LabelsIsNil() predicate.Resource {
	return predicate.Resource(sql.FieldIsNull(FieldLabels))
}

// LabelsNotNil applies the NotNil predicate on the "labels" field.
func LabelsNotNil() predicate.Resource {
	return predicate.Resource(sql.FieldNotNull(FieldLabels))
}

// AnnotationsIsNil applies the IsNil predicate on the "annotations" field.
func AnnotationsIsNil() predicate.Resource {
	return predicate.Resource(sql.FieldIsNull(FieldAnnotations))
}

// AnnotationsNotNil applies the NotNil predicate on the "annotations" field.
func AnnotationsNotNil() predicate.Resource {
	return predicate.Resource(sql.FieldNotNull(FieldAnnotations))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.Resource {
	return predicate.Resource(sql.FieldLTE(FieldUpdateTime, v))
}

// StatusIsNil applies the IsNil predicate on the "status" field.
func StatusIsNil() predicate.Resource {
	return predicate.Resource(sql.FieldIsNull(FieldStatus))
}

// StatusNotNil applies the NotNil predicate on the "status" field.
func StatusNotNil() predicate.Resource {
	return predicate.Resource(sql.FieldNotNull(FieldStatus))
}

// ProjectIDEQ applies the EQ predicate on the "project_id" field.
func ProjectIDEQ(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldProjectID, v))
}

// ProjectIDNEQ applies the NEQ predicate on the "project_id" field.
func ProjectIDNEQ(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldNEQ(FieldProjectID, v))
}

// ProjectIDIn applies the In predicate on the "project_id" field.
func ProjectIDIn(vs ...object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldIn(FieldProjectID, vs...))
}

// ProjectIDNotIn applies the NotIn predicate on the "project_id" field.
func ProjectIDNotIn(vs ...object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldNotIn(FieldProjectID, vs...))
}

// ProjectIDGT applies the GT predicate on the "project_id" field.
func ProjectIDGT(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldGT(FieldProjectID, v))
}

// ProjectIDGTE applies the GTE predicate on the "project_id" field.
func ProjectIDGTE(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldGTE(FieldProjectID, v))
}

// ProjectIDLT applies the LT predicate on the "project_id" field.
func ProjectIDLT(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldLT(FieldProjectID, v))
}

// ProjectIDLTE applies the LTE predicate on the "project_id" field.
func ProjectIDLTE(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldLTE(FieldProjectID, v))
}

// ProjectIDContains applies the Contains predicate on the "project_id" field.
func ProjectIDContains(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldContains(FieldProjectID, vc))
}

// ProjectIDHasPrefix applies the HasPrefix predicate on the "project_id" field.
func ProjectIDHasPrefix(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldHasPrefix(FieldProjectID, vc))
}

// ProjectIDHasSuffix applies the HasSuffix predicate on the "project_id" field.
func ProjectIDHasSuffix(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldHasSuffix(FieldProjectID, vc))
}

// ProjectIDEqualFold applies the EqualFold predicate on the "project_id" field.
func ProjectIDEqualFold(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldEqualFold(FieldProjectID, vc))
}

// ProjectIDContainsFold applies the ContainsFold predicate on the "project_id" field.
func ProjectIDContainsFold(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldContainsFold(FieldProjectID, vc))
}

// EnvironmentIDEQ applies the EQ predicate on the "environment_id" field.
func EnvironmentIDEQ(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldEnvironmentID, v))
}

// EnvironmentIDNEQ applies the NEQ predicate on the "environment_id" field.
func EnvironmentIDNEQ(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldNEQ(FieldEnvironmentID, v))
}

// EnvironmentIDIn applies the In predicate on the "environment_id" field.
func EnvironmentIDIn(vs ...object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDNotIn applies the NotIn predicate on the "environment_id" field.
func EnvironmentIDNotIn(vs ...object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldNotIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDGT applies the GT predicate on the "environment_id" field.
func EnvironmentIDGT(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldGT(FieldEnvironmentID, v))
}

// EnvironmentIDGTE applies the GTE predicate on the "environment_id" field.
func EnvironmentIDGTE(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldGTE(FieldEnvironmentID, v))
}

// EnvironmentIDLT applies the LT predicate on the "environment_id" field.
func EnvironmentIDLT(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldLT(FieldEnvironmentID, v))
}

// EnvironmentIDLTE applies the LTE predicate on the "environment_id" field.
func EnvironmentIDLTE(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldLTE(FieldEnvironmentID, v))
}

// EnvironmentIDContains applies the Contains predicate on the "environment_id" field.
func EnvironmentIDContains(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldContains(FieldEnvironmentID, vc))
}

// EnvironmentIDHasPrefix applies the HasPrefix predicate on the "environment_id" field.
func EnvironmentIDHasPrefix(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldHasPrefix(FieldEnvironmentID, vc))
}

// EnvironmentIDHasSuffix applies the HasSuffix predicate on the "environment_id" field.
func EnvironmentIDHasSuffix(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldHasSuffix(FieldEnvironmentID, vc))
}

// EnvironmentIDEqualFold applies the EqualFold predicate on the "environment_id" field.
func EnvironmentIDEqualFold(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldEqualFold(FieldEnvironmentID, vc))
}

// EnvironmentIDContainsFold applies the ContainsFold predicate on the "environment_id" field.
func EnvironmentIDContainsFold(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldContainsFold(FieldEnvironmentID, vc))
}

// TemplateIDEQ applies the EQ predicate on the "template_id" field.
func TemplateIDEQ(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldTemplateID, v))
}

// TemplateIDNEQ applies the NEQ predicate on the "template_id" field.
func TemplateIDNEQ(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldNEQ(FieldTemplateID, v))
}

// TemplateIDIn applies the In predicate on the "template_id" field.
func TemplateIDIn(vs ...object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldIn(FieldTemplateID, vs...))
}

// TemplateIDNotIn applies the NotIn predicate on the "template_id" field.
func TemplateIDNotIn(vs ...object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldNotIn(FieldTemplateID, vs...))
}

// TemplateIDGT applies the GT predicate on the "template_id" field.
func TemplateIDGT(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldGT(FieldTemplateID, v))
}

// TemplateIDGTE applies the GTE predicate on the "template_id" field.
func TemplateIDGTE(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldGTE(FieldTemplateID, v))
}

// TemplateIDLT applies the LT predicate on the "template_id" field.
func TemplateIDLT(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldLT(FieldTemplateID, v))
}

// TemplateIDLTE applies the LTE predicate on the "template_id" field.
func TemplateIDLTE(v object.ID) predicate.Resource {
	return predicate.Resource(sql.FieldLTE(FieldTemplateID, v))
}

// TemplateIDContains applies the Contains predicate on the "template_id" field.
func TemplateIDContains(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldContains(FieldTemplateID, vc))
}

// TemplateIDHasPrefix applies the HasPrefix predicate on the "template_id" field.
func TemplateIDHasPrefix(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldHasPrefix(FieldTemplateID, vc))
}

// TemplateIDHasSuffix applies the HasSuffix predicate on the "template_id" field.
func TemplateIDHasSuffix(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldHasSuffix(FieldTemplateID, vc))
}

// TemplateIDEqualFold applies the EqualFold predicate on the "template_id" field.
func TemplateIDEqualFold(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldEqualFold(FieldTemplateID, vc))
}

// TemplateIDContainsFold applies the ContainsFold predicate on the "template_id" field.
func TemplateIDContainsFold(v object.ID) predicate.Resource {
	vc := string(v)
	return predicate.Resource(sql.FieldContainsFold(FieldTemplateID, vc))
}

// AttributesEQ applies the EQ predicate on the "attributes" field.
func AttributesEQ(v property.Values) predicate.Resource {
	return predicate.Resource(sql.FieldEQ(FieldAttributes, v))
}

// AttributesNEQ applies the NEQ predicate on the "attributes" field.
func AttributesNEQ(v property.Values) predicate.Resource {
	return predicate.Resource(sql.FieldNEQ(FieldAttributes, v))
}

// AttributesIn applies the In predicate on the "attributes" field.
func AttributesIn(vs ...property.Values) predicate.Resource {
	return predicate.Resource(sql.FieldIn(FieldAttributes, vs...))
}

// AttributesNotIn applies the NotIn predicate on the "attributes" field.
func AttributesNotIn(vs ...property.Values) predicate.Resource {
	return predicate.Resource(sql.FieldNotIn(FieldAttributes, vs...))
}

// AttributesGT applies the GT predicate on the "attributes" field.
func AttributesGT(v property.Values) predicate.Resource {
	return predicate.Resource(sql.FieldGT(FieldAttributes, v))
}

// AttributesGTE applies the GTE predicate on the "attributes" field.
func AttributesGTE(v property.Values) predicate.Resource {
	return predicate.Resource(sql.FieldGTE(FieldAttributes, v))
}

// AttributesLT applies the LT predicate on the "attributes" field.
func AttributesLT(v property.Values) predicate.Resource {
	return predicate.Resource(sql.FieldLT(FieldAttributes, v))
}

// AttributesLTE applies the LTE predicate on the "attributes" field.
func AttributesLTE(v property.Values) predicate.Resource {
	return predicate.Resource(sql.FieldLTE(FieldAttributes, v))
}

// AttributesIsNil applies the IsNil predicate on the "attributes" field.
func AttributesIsNil() predicate.Resource {
	return predicate.Resource(sql.FieldIsNull(FieldAttributes))
}

// AttributesNotNil applies the NotNil predicate on the "attributes" field.
func AttributesNotNil() predicate.Resource {
	return predicate.Resource(sql.FieldNotNull(FieldAttributes))
}

// HasProject applies the HasEdge predicate on the "project" edge.
func HasProject() predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.Resource
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectWith applies the HasEdge predicate on the "project" edge with a given conditions (other predicates).
func HasProjectWith(preds ...predicate.Project) predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := newProjectStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.Resource
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEnvironment applies the HasEdge predicate on the "environment" edge.
func HasEnvironment() predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EnvironmentTable, EnvironmentColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.Resource
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEnvironmentWith applies the HasEdge predicate on the "environment" edge with a given conditions (other predicates).
func HasEnvironmentWith(preds ...predicate.Environment) predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := newEnvironmentStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.Resource
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTemplate applies the HasEdge predicate on the "template" edge.
func HasTemplate() predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, TemplateTable, TemplateColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.TemplateVersion
		step.Edge.Schema = schemaConfig.Resource
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTemplateWith applies the HasEdge predicate on the "template" edge with a given conditions (other predicates).
func HasTemplateWith(preds ...predicate.TemplateVersion) predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := newTemplateStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.TemplateVersion
		step.Edge.Schema = schemaConfig.Resource
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRevisions applies the HasEdge predicate on the "revisions" edge.
func HasRevisions() predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, RevisionsTable, RevisionsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ResourceRevision
		step.Edge.Schema = schemaConfig.ResourceRevision
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRevisionsWith applies the HasEdge predicate on the "revisions" edge with a given conditions (other predicates).
func HasRevisionsWith(preds ...predicate.ResourceRevision) predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := newRevisionsStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ResourceRevision
		step.Edge.Schema = schemaConfig.ResourceRevision
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasComponents applies the HasEdge predicate on the "components" edge.
func HasComponents() predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ComponentsTable, ComponentsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ResourceComponent
		step.Edge.Schema = schemaConfig.ResourceComponent
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasComponentsWith applies the HasEdge predicate on the "components" edge with a given conditions (other predicates).
func HasComponentsWith(preds ...predicate.ResourceComponent) predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := newComponentsStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ResourceComponent
		step.Edge.Schema = schemaConfig.ResourceComponent
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDependencies applies the HasEdge predicate on the "dependencies" edge.
func HasDependencies() predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, DependenciesTable, DependenciesColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ResourceRelationship
		step.Edge.Schema = schemaConfig.ResourceRelationship
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDependenciesWith applies the HasEdge predicate on the "dependencies" edge with a given conditions (other predicates).
func HasDependenciesWith(preds ...predicate.ResourceRelationship) predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		step := newDependenciesStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ResourceRelationship
		step.Edge.Schema = schemaConfig.ResourceRelationship
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Resource) predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Resource) predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
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
func Not(p predicate.Resource) predicate.Resource {
	return predicate.Resource(func(s *sql.Selector) {
		p(s.Not())
	})
}
