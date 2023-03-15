// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package applicationinstance

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ID filters vertices based on their ID field.
func ID(id oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLTE(FieldID, id))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldStatus, v))
}

// StatusMessage applies equality check predicate on the "statusMessage" field. It's identical to StatusMessageEQ.
func StatusMessage(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldStatusMessage, v))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldUpdateTime, v))
}

// ApplicationID applies equality check predicate on the "applicationID" field. It's identical to ApplicationIDEQ.
func ApplicationID(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldApplicationID, v))
}

// EnvironmentID applies equality check predicate on the "environmentID" field. It's identical to EnvironmentIDEQ.
func EnvironmentID(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldEnvironmentID, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldName, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLTE(FieldStatus, v))
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldContains(FieldStatus, v))
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldHasPrefix(FieldStatus, v))
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldHasSuffix(FieldStatus, v))
}

// StatusIsNil applies the IsNil predicate on the "status" field.
func StatusIsNil() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIsNull(FieldStatus))
}

// StatusNotNil applies the NotNil predicate on the "status" field.
func StatusNotNil() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotNull(FieldStatus))
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEqualFold(FieldStatus, v))
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldContainsFold(FieldStatus, v))
}

// StatusMessageEQ applies the EQ predicate on the "statusMessage" field.
func StatusMessageEQ(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldStatusMessage, v))
}

// StatusMessageNEQ applies the NEQ predicate on the "statusMessage" field.
func StatusMessageNEQ(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNEQ(FieldStatusMessage, v))
}

// StatusMessageIn applies the In predicate on the "statusMessage" field.
func StatusMessageIn(vs ...string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIn(FieldStatusMessage, vs...))
}

// StatusMessageNotIn applies the NotIn predicate on the "statusMessage" field.
func StatusMessageNotIn(vs ...string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotIn(FieldStatusMessage, vs...))
}

// StatusMessageGT applies the GT predicate on the "statusMessage" field.
func StatusMessageGT(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGT(FieldStatusMessage, v))
}

// StatusMessageGTE applies the GTE predicate on the "statusMessage" field.
func StatusMessageGTE(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGTE(FieldStatusMessage, v))
}

// StatusMessageLT applies the LT predicate on the "statusMessage" field.
func StatusMessageLT(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLT(FieldStatusMessage, v))
}

// StatusMessageLTE applies the LTE predicate on the "statusMessage" field.
func StatusMessageLTE(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLTE(FieldStatusMessage, v))
}

// StatusMessageContains applies the Contains predicate on the "statusMessage" field.
func StatusMessageContains(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldContains(FieldStatusMessage, v))
}

// StatusMessageHasPrefix applies the HasPrefix predicate on the "statusMessage" field.
func StatusMessageHasPrefix(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldHasPrefix(FieldStatusMessage, v))
}

// StatusMessageHasSuffix applies the HasSuffix predicate on the "statusMessage" field.
func StatusMessageHasSuffix(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldHasSuffix(FieldStatusMessage, v))
}

// StatusMessageIsNil applies the IsNil predicate on the "statusMessage" field.
func StatusMessageIsNil() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIsNull(FieldStatusMessage))
}

// StatusMessageNotNil applies the NotNil predicate on the "statusMessage" field.
func StatusMessageNotNil() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotNull(FieldStatusMessage))
}

// StatusMessageEqualFold applies the EqualFold predicate on the "statusMessage" field.
func StatusMessageEqualFold(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEqualFold(FieldStatusMessage, v))
}

// StatusMessageContainsFold applies the ContainsFold predicate on the "statusMessage" field.
func StatusMessageContainsFold(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldContainsFold(FieldStatusMessage, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLTE(FieldUpdateTime, v))
}

// ApplicationIDEQ applies the EQ predicate on the "applicationID" field.
func ApplicationIDEQ(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldApplicationID, v))
}

// ApplicationIDNEQ applies the NEQ predicate on the "applicationID" field.
func ApplicationIDNEQ(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNEQ(FieldApplicationID, v))
}

// ApplicationIDIn applies the In predicate on the "applicationID" field.
func ApplicationIDIn(vs ...oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIn(FieldApplicationID, vs...))
}

// ApplicationIDNotIn applies the NotIn predicate on the "applicationID" field.
func ApplicationIDNotIn(vs ...oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotIn(FieldApplicationID, vs...))
}

// ApplicationIDGT applies the GT predicate on the "applicationID" field.
func ApplicationIDGT(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGT(FieldApplicationID, v))
}

// ApplicationIDGTE applies the GTE predicate on the "applicationID" field.
func ApplicationIDGTE(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGTE(FieldApplicationID, v))
}

// ApplicationIDLT applies the LT predicate on the "applicationID" field.
func ApplicationIDLT(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLT(FieldApplicationID, v))
}

// ApplicationIDLTE applies the LTE predicate on the "applicationID" field.
func ApplicationIDLTE(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLTE(FieldApplicationID, v))
}

// ApplicationIDContains applies the Contains predicate on the "applicationID" field.
func ApplicationIDContains(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldContains(FieldApplicationID, vc))
}

// ApplicationIDHasPrefix applies the HasPrefix predicate on the "applicationID" field.
func ApplicationIDHasPrefix(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldHasPrefix(FieldApplicationID, vc))
}

// ApplicationIDHasSuffix applies the HasSuffix predicate on the "applicationID" field.
func ApplicationIDHasSuffix(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldHasSuffix(FieldApplicationID, vc))
}

// ApplicationIDEqualFold applies the EqualFold predicate on the "applicationID" field.
func ApplicationIDEqualFold(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldEqualFold(FieldApplicationID, vc))
}

// ApplicationIDContainsFold applies the ContainsFold predicate on the "applicationID" field.
func ApplicationIDContainsFold(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldContainsFold(FieldApplicationID, vc))
}

// EnvironmentIDEQ applies the EQ predicate on the "environmentID" field.
func EnvironmentIDEQ(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldEnvironmentID, v))
}

// EnvironmentIDNEQ applies the NEQ predicate on the "environmentID" field.
func EnvironmentIDNEQ(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNEQ(FieldEnvironmentID, v))
}

// EnvironmentIDIn applies the In predicate on the "environmentID" field.
func EnvironmentIDIn(vs ...oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDNotIn applies the NotIn predicate on the "environmentID" field.
func EnvironmentIDNotIn(vs ...oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDGT applies the GT predicate on the "environmentID" field.
func EnvironmentIDGT(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGT(FieldEnvironmentID, v))
}

// EnvironmentIDGTE applies the GTE predicate on the "environmentID" field.
func EnvironmentIDGTE(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGTE(FieldEnvironmentID, v))
}

// EnvironmentIDLT applies the LT predicate on the "environmentID" field.
func EnvironmentIDLT(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLT(FieldEnvironmentID, v))
}

// EnvironmentIDLTE applies the LTE predicate on the "environmentID" field.
func EnvironmentIDLTE(v oid.ID) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLTE(FieldEnvironmentID, v))
}

// EnvironmentIDContains applies the Contains predicate on the "environmentID" field.
func EnvironmentIDContains(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldContains(FieldEnvironmentID, vc))
}

// EnvironmentIDHasPrefix applies the HasPrefix predicate on the "environmentID" field.
func EnvironmentIDHasPrefix(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldHasPrefix(FieldEnvironmentID, vc))
}

// EnvironmentIDHasSuffix applies the HasSuffix predicate on the "environmentID" field.
func EnvironmentIDHasSuffix(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldHasSuffix(FieldEnvironmentID, vc))
}

// EnvironmentIDEqualFold applies the EqualFold predicate on the "environmentID" field.
func EnvironmentIDEqualFold(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldEqualFold(FieldEnvironmentID, vc))
}

// EnvironmentIDContainsFold applies the ContainsFold predicate on the "environmentID" field.
func EnvironmentIDContainsFold(v oid.ID) predicate.ApplicationInstance {
	vc := string(v)
	return predicate.ApplicationInstance(sql.FieldContainsFold(FieldEnvironmentID, vc))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldContainsFold(FieldName, v))
}

// VariablesIsNil applies the IsNil predicate on the "variables" field.
func VariablesIsNil() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldIsNull(FieldVariables))
}

// VariablesNotNil applies the NotNil predicate on the "variables" field.
func VariablesNotNil() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(sql.FieldNotNull(FieldVariables))
}

// HasApplication applies the HasEdge predicate on the "application" edge.
func HasApplication() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ApplicationTable, ApplicationColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationInstance
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasApplicationWith applies the HasEdge predicate on the "application" edge with a given conditions (other predicates).
func HasApplicationWith(preds ...predicate.Application) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ApplicationInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ApplicationTable, ApplicationColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationInstance
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEnvironment applies the HasEdge predicate on the "environment" edge.
func HasEnvironment() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EnvironmentTable, EnvironmentColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.ApplicationInstance
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEnvironmentWith applies the HasEdge predicate on the "environment" edge with a given conditions (other predicates).
func HasEnvironmentWith(preds ...predicate.Environment) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EnvironmentInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EnvironmentTable, EnvironmentColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.ApplicationInstance
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRevisions applies the HasEdge predicate on the "revisions" edge.
func HasRevisions() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
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
func HasRevisionsWith(preds ...predicate.ApplicationRevision) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
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

// HasResources applies the HasEdge predicate on the "resources" edge.
func HasResources() predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
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
func HasResourcesWith(preds ...predicate.ApplicationResource) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
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

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ApplicationInstance) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ApplicationInstance) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
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
func Not(p predicate.ApplicationInstance) predicate.ApplicationInstance {
	return predicate.ApplicationInstance(func(s *sql.Selector) {
		p(s.Not())
	})
}
