// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package applicationrevision

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ID filters vertices based on their ID field.
func ID(id types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldID, id))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldStatus, v))
}

// StatusMessage applies equality check predicate on the "statusMessage" field. It's identical to StatusMessageEQ.
func StatusMessage(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldStatusMessage, v))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldUpdateTime, v))
}

// ApplicationID applies equality check predicate on the "applicationID" field. It's identical to ApplicationIDEQ.
func ApplicationID(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldApplicationID, v))
}

// EnvironmentID applies equality check predicate on the "environmentID" field. It's identical to EnvironmentIDEQ.
func EnvironmentID(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldEnvironmentID, v))
}

// InputPlan applies equality check predicate on the "inputPlan" field. It's identical to InputPlanEQ.
func InputPlan(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldInputPlan, v))
}

// Output applies equality check predicate on the "output" field. It's identical to OutputEQ.
func Output(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldOutput, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldStatus, v))
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContains(FieldStatus, v))
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasPrefix(FieldStatus, v))
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasSuffix(FieldStatus, v))
}

// StatusIsNil applies the IsNil predicate on the "status" field.
func StatusIsNil() predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIsNull(FieldStatus))
}

// StatusNotNil applies the NotNil predicate on the "status" field.
func StatusNotNil() predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotNull(FieldStatus))
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEqualFold(FieldStatus, v))
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContainsFold(FieldStatus, v))
}

// StatusMessageEQ applies the EQ predicate on the "statusMessage" field.
func StatusMessageEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldStatusMessage, v))
}

// StatusMessageNEQ applies the NEQ predicate on the "statusMessage" field.
func StatusMessageNEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldStatusMessage, v))
}

// StatusMessageIn applies the In predicate on the "statusMessage" field.
func StatusMessageIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldStatusMessage, vs...))
}

// StatusMessageNotIn applies the NotIn predicate on the "statusMessage" field.
func StatusMessageNotIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldStatusMessage, vs...))
}

// StatusMessageGT applies the GT predicate on the "statusMessage" field.
func StatusMessageGT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldStatusMessage, v))
}

// StatusMessageGTE applies the GTE predicate on the "statusMessage" field.
func StatusMessageGTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldStatusMessage, v))
}

// StatusMessageLT applies the LT predicate on the "statusMessage" field.
func StatusMessageLT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldStatusMessage, v))
}

// StatusMessageLTE applies the LTE predicate on the "statusMessage" field.
func StatusMessageLTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldStatusMessage, v))
}

// StatusMessageContains applies the Contains predicate on the "statusMessage" field.
func StatusMessageContains(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContains(FieldStatusMessage, v))
}

// StatusMessageHasPrefix applies the HasPrefix predicate on the "statusMessage" field.
func StatusMessageHasPrefix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasPrefix(FieldStatusMessage, v))
}

// StatusMessageHasSuffix applies the HasSuffix predicate on the "statusMessage" field.
func StatusMessageHasSuffix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasSuffix(FieldStatusMessage, v))
}

// StatusMessageIsNil applies the IsNil predicate on the "statusMessage" field.
func StatusMessageIsNil() predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIsNull(FieldStatusMessage))
}

// StatusMessageNotNil applies the NotNil predicate on the "statusMessage" field.
func StatusMessageNotNil() predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotNull(FieldStatusMessage))
}

// StatusMessageEqualFold applies the EqualFold predicate on the "statusMessage" field.
func StatusMessageEqualFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEqualFold(FieldStatusMessage, v))
}

// StatusMessageContainsFold applies the ContainsFold predicate on the "statusMessage" field.
func StatusMessageContainsFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContainsFold(FieldStatusMessage, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldUpdateTime, v))
}

// ApplicationIDEQ applies the EQ predicate on the "applicationID" field.
func ApplicationIDEQ(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldApplicationID, v))
}

// ApplicationIDNEQ applies the NEQ predicate on the "applicationID" field.
func ApplicationIDNEQ(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldApplicationID, v))
}

// ApplicationIDIn applies the In predicate on the "applicationID" field.
func ApplicationIDIn(vs ...types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldApplicationID, vs...))
}

// ApplicationIDNotIn applies the NotIn predicate on the "applicationID" field.
func ApplicationIDNotIn(vs ...types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldApplicationID, vs...))
}

// ApplicationIDGT applies the GT predicate on the "applicationID" field.
func ApplicationIDGT(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldApplicationID, v))
}

// ApplicationIDGTE applies the GTE predicate on the "applicationID" field.
func ApplicationIDGTE(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldApplicationID, v))
}

// ApplicationIDLT applies the LT predicate on the "applicationID" field.
func ApplicationIDLT(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldApplicationID, v))
}

// ApplicationIDLTE applies the LTE predicate on the "applicationID" field.
func ApplicationIDLTE(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldApplicationID, v))
}

// ApplicationIDContains applies the Contains predicate on the "applicationID" field.
func ApplicationIDContains(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldContains(FieldApplicationID, vc))
}

// ApplicationIDHasPrefix applies the HasPrefix predicate on the "applicationID" field.
func ApplicationIDHasPrefix(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldHasPrefix(FieldApplicationID, vc))
}

// ApplicationIDHasSuffix applies the HasSuffix predicate on the "applicationID" field.
func ApplicationIDHasSuffix(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldHasSuffix(FieldApplicationID, vc))
}

// ApplicationIDEqualFold applies the EqualFold predicate on the "applicationID" field.
func ApplicationIDEqualFold(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldEqualFold(FieldApplicationID, vc))
}

// ApplicationIDContainsFold applies the ContainsFold predicate on the "applicationID" field.
func ApplicationIDContainsFold(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldContainsFold(FieldApplicationID, vc))
}

// EnvironmentIDEQ applies the EQ predicate on the "environmentID" field.
func EnvironmentIDEQ(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldEnvironmentID, v))
}

// EnvironmentIDNEQ applies the NEQ predicate on the "environmentID" field.
func EnvironmentIDNEQ(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldEnvironmentID, v))
}

// EnvironmentIDIn applies the In predicate on the "environmentID" field.
func EnvironmentIDIn(vs ...types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDNotIn applies the NotIn predicate on the "environmentID" field.
func EnvironmentIDNotIn(vs ...types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDGT applies the GT predicate on the "environmentID" field.
func EnvironmentIDGT(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldEnvironmentID, v))
}

// EnvironmentIDGTE applies the GTE predicate on the "environmentID" field.
func EnvironmentIDGTE(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldEnvironmentID, v))
}

// EnvironmentIDLT applies the LT predicate on the "environmentID" field.
func EnvironmentIDLT(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldEnvironmentID, v))
}

// EnvironmentIDLTE applies the LTE predicate on the "environmentID" field.
func EnvironmentIDLTE(v types.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldEnvironmentID, v))
}

// EnvironmentIDContains applies the Contains predicate on the "environmentID" field.
func EnvironmentIDContains(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldContains(FieldEnvironmentID, vc))
}

// EnvironmentIDHasPrefix applies the HasPrefix predicate on the "environmentID" field.
func EnvironmentIDHasPrefix(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldHasPrefix(FieldEnvironmentID, vc))
}

// EnvironmentIDHasSuffix applies the HasSuffix predicate on the "environmentID" field.
func EnvironmentIDHasSuffix(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldHasSuffix(FieldEnvironmentID, vc))
}

// EnvironmentIDEqualFold applies the EqualFold predicate on the "environmentID" field.
func EnvironmentIDEqualFold(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldEqualFold(FieldEnvironmentID, vc))
}

// EnvironmentIDContainsFold applies the ContainsFold predicate on the "environmentID" field.
func EnvironmentIDContainsFold(v types.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldContainsFold(FieldEnvironmentID, vc))
}

// InputPlanEQ applies the EQ predicate on the "inputPlan" field.
func InputPlanEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldInputPlan, v))
}

// InputPlanNEQ applies the NEQ predicate on the "inputPlan" field.
func InputPlanNEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldInputPlan, v))
}

// InputPlanIn applies the In predicate on the "inputPlan" field.
func InputPlanIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldInputPlan, vs...))
}

// InputPlanNotIn applies the NotIn predicate on the "inputPlan" field.
func InputPlanNotIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldInputPlan, vs...))
}

// InputPlanGT applies the GT predicate on the "inputPlan" field.
func InputPlanGT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldInputPlan, v))
}

// InputPlanGTE applies the GTE predicate on the "inputPlan" field.
func InputPlanGTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldInputPlan, v))
}

// InputPlanLT applies the LT predicate on the "inputPlan" field.
func InputPlanLT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldInputPlan, v))
}

// InputPlanLTE applies the LTE predicate on the "inputPlan" field.
func InputPlanLTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldInputPlan, v))
}

// InputPlanContains applies the Contains predicate on the "inputPlan" field.
func InputPlanContains(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContains(FieldInputPlan, v))
}

// InputPlanHasPrefix applies the HasPrefix predicate on the "inputPlan" field.
func InputPlanHasPrefix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasPrefix(FieldInputPlan, v))
}

// InputPlanHasSuffix applies the HasSuffix predicate on the "inputPlan" field.
func InputPlanHasSuffix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasSuffix(FieldInputPlan, v))
}

// InputPlanEqualFold applies the EqualFold predicate on the "inputPlan" field.
func InputPlanEqualFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEqualFold(FieldInputPlan, v))
}

// InputPlanContainsFold applies the ContainsFold predicate on the "inputPlan" field.
func InputPlanContainsFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContainsFold(FieldInputPlan, v))
}

// OutputEQ applies the EQ predicate on the "output" field.
func OutputEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldOutput, v))
}

// OutputNEQ applies the NEQ predicate on the "output" field.
func OutputNEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldOutput, v))
}

// OutputIn applies the In predicate on the "output" field.
func OutputIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldOutput, vs...))
}

// OutputNotIn applies the NotIn predicate on the "output" field.
func OutputNotIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldOutput, vs...))
}

// OutputGT applies the GT predicate on the "output" field.
func OutputGT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldOutput, v))
}

// OutputGTE applies the GTE predicate on the "output" field.
func OutputGTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldOutput, v))
}

// OutputLT applies the LT predicate on the "output" field.
func OutputLT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldOutput, v))
}

// OutputLTE applies the LTE predicate on the "output" field.
func OutputLTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldOutput, v))
}

// OutputContains applies the Contains predicate on the "output" field.
func OutputContains(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContains(FieldOutput, v))
}

// OutputHasPrefix applies the HasPrefix predicate on the "output" field.
func OutputHasPrefix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasPrefix(FieldOutput, v))
}

// OutputHasSuffix applies the HasSuffix predicate on the "output" field.
func OutputHasSuffix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasSuffix(FieldOutput, v))
}

// OutputEqualFold applies the EqualFold predicate on the "output" field.
func OutputEqualFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEqualFold(FieldOutput, v))
}

// OutputContainsFold applies the ContainsFold predicate on the "output" field.
func OutputContainsFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContainsFold(FieldOutput, v))
}

// HasApplication applies the HasEdge predicate on the "application" edge.
func HasApplication() predicate.ApplicationRevision {
	return predicate.ApplicationRevision(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ApplicationTable, ApplicationColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationRevision
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasApplicationWith applies the HasEdge predicate on the "application" edge with a given conditions (other predicates).
func HasApplicationWith(preds ...predicate.Application) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ApplicationInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ApplicationTable, ApplicationColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationRevision
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ApplicationRevision) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ApplicationRevision) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(func(s *sql.Selector) {
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
func Not(p predicate.ApplicationRevision) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(func(s *sql.Selector) {
		p(s.Not())
	})
}
