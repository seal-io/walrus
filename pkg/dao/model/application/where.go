// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package application

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ID filters vertices based on their ID field.
func ID(id types.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id types.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id types.ID) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...types.ID) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...types.ID) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id types.ID) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id types.ID) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id types.ID) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id types.ID) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldDescription, v))
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
func ProjectID(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldProjectID, v))
}

// EnvironmentID applies equality check predicate on the "environmentID" field. It's identical to EnvironmentIDEQ.
func EnvironmentID(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldEnvironmentID, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Application {
	return predicate.Application(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Application {
	return predicate.Application(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Application {
	return predicate.Application(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Application {
	return predicate.Application(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Application {
	return predicate.Application(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Application {
	return predicate.Application(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Application {
	return predicate.Application(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Application {
	return predicate.Application(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Application {
	return predicate.Application(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Application {
	return predicate.Application(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Application {
	return predicate.Application(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Application {
	return predicate.Application(sql.FieldContainsFold(FieldDescription, v))
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
func ProjectIDEQ(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldProjectID, v))
}

// ProjectIDNEQ applies the NEQ predicate on the "projectID" field.
func ProjectIDNEQ(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldProjectID, v))
}

// ProjectIDIn applies the In predicate on the "projectID" field.
func ProjectIDIn(vs ...types.ID) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldProjectID, vs...))
}

// ProjectIDNotIn applies the NotIn predicate on the "projectID" field.
func ProjectIDNotIn(vs ...types.ID) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldProjectID, vs...))
}

// ProjectIDGT applies the GT predicate on the "projectID" field.
func ProjectIDGT(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldProjectID, v))
}

// ProjectIDGTE applies the GTE predicate on the "projectID" field.
func ProjectIDGTE(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldProjectID, v))
}

// ProjectIDLT applies the LT predicate on the "projectID" field.
func ProjectIDLT(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldProjectID, v))
}

// ProjectIDLTE applies the LTE predicate on the "projectID" field.
func ProjectIDLTE(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldProjectID, v))
}

// ProjectIDContains applies the Contains predicate on the "projectID" field.
func ProjectIDContains(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldContains(FieldProjectID, vc))
}

// ProjectIDHasPrefix applies the HasPrefix predicate on the "projectID" field.
func ProjectIDHasPrefix(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldHasPrefix(FieldProjectID, vc))
}

// ProjectIDHasSuffix applies the HasSuffix predicate on the "projectID" field.
func ProjectIDHasSuffix(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldHasSuffix(FieldProjectID, vc))
}

// ProjectIDEqualFold applies the EqualFold predicate on the "projectID" field.
func ProjectIDEqualFold(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldEqualFold(FieldProjectID, vc))
}

// ProjectIDContainsFold applies the ContainsFold predicate on the "projectID" field.
func ProjectIDContainsFold(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldContainsFold(FieldProjectID, vc))
}

// EnvironmentIDEQ applies the EQ predicate on the "environmentID" field.
func EnvironmentIDEQ(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldEQ(FieldEnvironmentID, v))
}

// EnvironmentIDNEQ applies the NEQ predicate on the "environmentID" field.
func EnvironmentIDNEQ(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldNEQ(FieldEnvironmentID, v))
}

// EnvironmentIDIn applies the In predicate on the "environmentID" field.
func EnvironmentIDIn(vs ...types.ID) predicate.Application {
	return predicate.Application(sql.FieldIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDNotIn applies the NotIn predicate on the "environmentID" field.
func EnvironmentIDNotIn(vs ...types.ID) predicate.Application {
	return predicate.Application(sql.FieldNotIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDGT applies the GT predicate on the "environmentID" field.
func EnvironmentIDGT(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldGT(FieldEnvironmentID, v))
}

// EnvironmentIDGTE applies the GTE predicate on the "environmentID" field.
func EnvironmentIDGTE(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldGTE(FieldEnvironmentID, v))
}

// EnvironmentIDLT applies the LT predicate on the "environmentID" field.
func EnvironmentIDLT(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldLT(FieldEnvironmentID, v))
}

// EnvironmentIDLTE applies the LTE predicate on the "environmentID" field.
func EnvironmentIDLTE(v types.ID) predicate.Application {
	return predicate.Application(sql.FieldLTE(FieldEnvironmentID, v))
}

// EnvironmentIDContains applies the Contains predicate on the "environmentID" field.
func EnvironmentIDContains(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldContains(FieldEnvironmentID, vc))
}

// EnvironmentIDHasPrefix applies the HasPrefix predicate on the "environmentID" field.
func EnvironmentIDHasPrefix(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldHasPrefix(FieldEnvironmentID, vc))
}

// EnvironmentIDHasSuffix applies the HasSuffix predicate on the "environmentID" field.
func EnvironmentIDHasSuffix(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldHasSuffix(FieldEnvironmentID, vc))
}

// EnvironmentIDEqualFold applies the EqualFold predicate on the "environmentID" field.
func EnvironmentIDEqualFold(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldEqualFold(FieldEnvironmentID, vc))
}

// EnvironmentIDContainsFold applies the ContainsFold predicate on the "environmentID" field.
func EnvironmentIDContainsFold(v types.ID) predicate.Application {
	vc := string(v)
	return predicate.Application(sql.FieldContainsFold(FieldEnvironmentID, vc))
}

// HasProject applies the HasEdge predicate on the "project" edge.
func HasProject() predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.Application
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectWith applies the HasEdge predicate on the "project" edge with a given conditions (other predicates).
func HasProjectWith(preds ...predicate.Project) predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProjectInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.Application
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEnvironment applies the HasEdge predicate on the "environment" edge.
func HasEnvironment() predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EnvironmentTable, EnvironmentColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.Application
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEnvironmentWith applies the HasEdge predicate on the "environment" edge with a given conditions (other predicates).
func HasEnvironmentWith(preds ...predicate.Environment) predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EnvironmentInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EnvironmentTable, EnvironmentColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.Application
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasResources applies the HasEdge predicate on the "resources" edge.
func HasResources() predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ResourcesTable, ResourcesColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationResource
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasResourcesWith applies the HasEdge predicate on the "resources" edge with a given conditions (other predicates).
func HasResourcesWith(preds ...predicate.ApplicationResource) predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ResourcesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ResourcesTable, ResourcesColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationResource
		step.Edge.Schema = schemaConfig.ApplicationResource
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRevisions applies the HasEdge predicate on the "revisions" edge.
func HasRevisions() predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, RevisionsTable, RevisionsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationRevision
		step.Edge.Schema = schemaConfig.ApplicationRevision
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRevisionsWith applies the HasEdge predicate on the "revisions" edge with a given conditions (other predicates).
func HasRevisionsWith(preds ...predicate.ApplicationRevision) predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(RevisionsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, RevisionsTable, RevisionsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationRevision
		step.Edge.Schema = schemaConfig.ApplicationRevision
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasModules applies the HasEdge predicate on the "modules" edge.
func HasModules() predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ModulesTable, ModulesPrimaryKey...),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Module
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasModulesWith applies the HasEdge predicate on the "modules" edge with a given conditions (other predicates).
func HasModulesWith(preds ...predicate.Module) predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ModulesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ModulesTable, ModulesPrimaryKey...),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Module
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasApplicationModuleRelationships applies the HasEdge predicate on the "applicationModuleRelationships" edge.
func HasApplicationModuleRelationships() predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ApplicationModuleRelationshipsTable, ApplicationModuleRelationshipsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationModuleRelationship
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasApplicationModuleRelationshipsWith applies the HasEdge predicate on the "applicationModuleRelationships" edge with a given conditions (other predicates).
func HasApplicationModuleRelationshipsWith(preds ...predicate.ApplicationModuleRelationship) predicate.Application {
	return predicate.Application(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ApplicationModuleRelationshipsInverseTable, ApplicationModuleRelationshipsColumn),
			sqlgraph.Edge(sqlgraph.O2M, true, ApplicationModuleRelationshipsTable, ApplicationModuleRelationshipsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationModuleRelationship
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
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
