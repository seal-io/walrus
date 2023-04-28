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
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

// ID filters vertices based on their ID field.
func ID(id oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id oid.ID) predicate.ApplicationRevision {
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

// InstanceID applies equality check predicate on the "instanceID" field. It's identical to InstanceIDEQ.
func InstanceID(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldInstanceID, v))
}

// EnvironmentID applies equality check predicate on the "environmentID" field. It's identical to EnvironmentIDEQ.
func EnvironmentID(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldEnvironmentID, v))
}

// Secrets applies equality check predicate on the "secrets" field. It's identical to SecretsEQ.
func Secrets(v crypto.Map[string, string]) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldSecrets, v))
}

// Variables applies equality check predicate on the "variables" field. It's identical to VariablesEQ.
func Variables(v property.Schemas) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldVariables, v))
}

// InputVariables applies equality check predicate on the "inputVariables" field. It's identical to InputVariablesEQ.
func InputVariables(v property.Values) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldInputVariables, v))
}

// InputPlan applies equality check predicate on the "inputPlan" field. It's identical to InputPlanEQ.
func InputPlan(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldInputPlan, v))
}

// Output applies equality check predicate on the "output" field. It's identical to OutputEQ.
func Output(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldOutput, v))
}

// DeployerType applies equality check predicate on the "deployerType" field. It's identical to DeployerTypeEQ.
func DeployerType(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldDeployerType, v))
}

// Duration applies equality check predicate on the "duration" field. It's identical to DurationEQ.
func Duration(v int) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldDuration, v))
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

// InstanceIDEQ applies the EQ predicate on the "instanceID" field.
func InstanceIDEQ(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldInstanceID, v))
}

// InstanceIDNEQ applies the NEQ predicate on the "instanceID" field.
func InstanceIDNEQ(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldInstanceID, v))
}

// InstanceIDIn applies the In predicate on the "instanceID" field.
func InstanceIDIn(vs ...oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldInstanceID, vs...))
}

// InstanceIDNotIn applies the NotIn predicate on the "instanceID" field.
func InstanceIDNotIn(vs ...oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldInstanceID, vs...))
}

// InstanceIDGT applies the GT predicate on the "instanceID" field.
func InstanceIDGT(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldInstanceID, v))
}

// InstanceIDGTE applies the GTE predicate on the "instanceID" field.
func InstanceIDGTE(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldInstanceID, v))
}

// InstanceIDLT applies the LT predicate on the "instanceID" field.
func InstanceIDLT(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldInstanceID, v))
}

// InstanceIDLTE applies the LTE predicate on the "instanceID" field.
func InstanceIDLTE(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldInstanceID, v))
}

// InstanceIDContains applies the Contains predicate on the "instanceID" field.
func InstanceIDContains(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldContains(FieldInstanceID, vc))
}

// InstanceIDHasPrefix applies the HasPrefix predicate on the "instanceID" field.
func InstanceIDHasPrefix(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldHasPrefix(FieldInstanceID, vc))
}

// InstanceIDHasSuffix applies the HasSuffix predicate on the "instanceID" field.
func InstanceIDHasSuffix(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldHasSuffix(FieldInstanceID, vc))
}

// InstanceIDEqualFold applies the EqualFold predicate on the "instanceID" field.
func InstanceIDEqualFold(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldEqualFold(FieldInstanceID, vc))
}

// InstanceIDContainsFold applies the ContainsFold predicate on the "instanceID" field.
func InstanceIDContainsFold(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldContainsFold(FieldInstanceID, vc))
}

// EnvironmentIDEQ applies the EQ predicate on the "environmentID" field.
func EnvironmentIDEQ(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldEnvironmentID, v))
}

// EnvironmentIDNEQ applies the NEQ predicate on the "environmentID" field.
func EnvironmentIDNEQ(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldEnvironmentID, v))
}

// EnvironmentIDIn applies the In predicate on the "environmentID" field.
func EnvironmentIDIn(vs ...oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDNotIn applies the NotIn predicate on the "environmentID" field.
func EnvironmentIDNotIn(vs ...oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDGT applies the GT predicate on the "environmentID" field.
func EnvironmentIDGT(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldEnvironmentID, v))
}

// EnvironmentIDGTE applies the GTE predicate on the "environmentID" field.
func EnvironmentIDGTE(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldEnvironmentID, v))
}

// EnvironmentIDLT applies the LT predicate on the "environmentID" field.
func EnvironmentIDLT(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldEnvironmentID, v))
}

// EnvironmentIDLTE applies the LTE predicate on the "environmentID" field.
func EnvironmentIDLTE(v oid.ID) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldEnvironmentID, v))
}

// EnvironmentIDContains applies the Contains predicate on the "environmentID" field.
func EnvironmentIDContains(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldContains(FieldEnvironmentID, vc))
}

// EnvironmentIDHasPrefix applies the HasPrefix predicate on the "environmentID" field.
func EnvironmentIDHasPrefix(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldHasPrefix(FieldEnvironmentID, vc))
}

// EnvironmentIDHasSuffix applies the HasSuffix predicate on the "environmentID" field.
func EnvironmentIDHasSuffix(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldHasSuffix(FieldEnvironmentID, vc))
}

// EnvironmentIDEqualFold applies the EqualFold predicate on the "environmentID" field.
func EnvironmentIDEqualFold(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldEqualFold(FieldEnvironmentID, vc))
}

// EnvironmentIDContainsFold applies the ContainsFold predicate on the "environmentID" field.
func EnvironmentIDContainsFold(v oid.ID) predicate.ApplicationRevision {
	vc := string(v)
	return predicate.ApplicationRevision(sql.FieldContainsFold(FieldEnvironmentID, vc))
}

// SecretsEQ applies the EQ predicate on the "secrets" field.
func SecretsEQ(v crypto.Map[string, string]) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldSecrets, v))
}

// SecretsNEQ applies the NEQ predicate on the "secrets" field.
func SecretsNEQ(v crypto.Map[string, string]) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldSecrets, v))
}

// SecretsIn applies the In predicate on the "secrets" field.
func SecretsIn(vs ...crypto.Map[string, string]) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldSecrets, vs...))
}

// SecretsNotIn applies the NotIn predicate on the "secrets" field.
func SecretsNotIn(vs ...crypto.Map[string, string]) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldSecrets, vs...))
}

// SecretsGT applies the GT predicate on the "secrets" field.
func SecretsGT(v crypto.Map[string, string]) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldSecrets, v))
}

// SecretsGTE applies the GTE predicate on the "secrets" field.
func SecretsGTE(v crypto.Map[string, string]) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldSecrets, v))
}

// SecretsLT applies the LT predicate on the "secrets" field.
func SecretsLT(v crypto.Map[string, string]) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldSecrets, v))
}

// SecretsLTE applies the LTE predicate on the "secrets" field.
func SecretsLTE(v crypto.Map[string, string]) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldSecrets, v))
}

// VariablesEQ applies the EQ predicate on the "variables" field.
func VariablesEQ(v property.Schemas) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldVariables, v))
}

// VariablesNEQ applies the NEQ predicate on the "variables" field.
func VariablesNEQ(v property.Schemas) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldVariables, v))
}

// VariablesIn applies the In predicate on the "variables" field.
func VariablesIn(vs ...property.Schemas) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldVariables, vs...))
}

// VariablesNotIn applies the NotIn predicate on the "variables" field.
func VariablesNotIn(vs ...property.Schemas) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldVariables, vs...))
}

// VariablesGT applies the GT predicate on the "variables" field.
func VariablesGT(v property.Schemas) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldVariables, v))
}

// VariablesGTE applies the GTE predicate on the "variables" field.
func VariablesGTE(v property.Schemas) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldVariables, v))
}

// VariablesLT applies the LT predicate on the "variables" field.
func VariablesLT(v property.Schemas) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldVariables, v))
}

// VariablesLTE applies the LTE predicate on the "variables" field.
func VariablesLTE(v property.Schemas) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldVariables, v))
}

// VariablesIsNil applies the IsNil predicate on the "variables" field.
func VariablesIsNil() predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIsNull(FieldVariables))
}

// VariablesNotNil applies the NotNil predicate on the "variables" field.
func VariablesNotNil() predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotNull(FieldVariables))
}

// InputVariablesEQ applies the EQ predicate on the "inputVariables" field.
func InputVariablesEQ(v property.Values) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldInputVariables, v))
}

// InputVariablesNEQ applies the NEQ predicate on the "inputVariables" field.
func InputVariablesNEQ(v property.Values) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldInputVariables, v))
}

// InputVariablesIn applies the In predicate on the "inputVariables" field.
func InputVariablesIn(vs ...property.Values) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldInputVariables, vs...))
}

// InputVariablesNotIn applies the NotIn predicate on the "inputVariables" field.
func InputVariablesNotIn(vs ...property.Values) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldInputVariables, vs...))
}

// InputVariablesGT applies the GT predicate on the "inputVariables" field.
func InputVariablesGT(v property.Values) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldInputVariables, v))
}

// InputVariablesGTE applies the GTE predicate on the "inputVariables" field.
func InputVariablesGTE(v property.Values) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldInputVariables, v))
}

// InputVariablesLT applies the LT predicate on the "inputVariables" field.
func InputVariablesLT(v property.Values) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldInputVariables, v))
}

// InputVariablesLTE applies the LTE predicate on the "inputVariables" field.
func InputVariablesLTE(v property.Values) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldInputVariables, v))
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

// DeployerTypeEQ applies the EQ predicate on the "deployerType" field.
func DeployerTypeEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldDeployerType, v))
}

// DeployerTypeNEQ applies the NEQ predicate on the "deployerType" field.
func DeployerTypeNEQ(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldDeployerType, v))
}

// DeployerTypeIn applies the In predicate on the "deployerType" field.
func DeployerTypeIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldDeployerType, vs...))
}

// DeployerTypeNotIn applies the NotIn predicate on the "deployerType" field.
func DeployerTypeNotIn(vs ...string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldDeployerType, vs...))
}

// DeployerTypeGT applies the GT predicate on the "deployerType" field.
func DeployerTypeGT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldDeployerType, v))
}

// DeployerTypeGTE applies the GTE predicate on the "deployerType" field.
func DeployerTypeGTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldDeployerType, v))
}

// DeployerTypeLT applies the LT predicate on the "deployerType" field.
func DeployerTypeLT(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldDeployerType, v))
}

// DeployerTypeLTE applies the LTE predicate on the "deployerType" field.
func DeployerTypeLTE(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldDeployerType, v))
}

// DeployerTypeContains applies the Contains predicate on the "deployerType" field.
func DeployerTypeContains(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContains(FieldDeployerType, v))
}

// DeployerTypeHasPrefix applies the HasPrefix predicate on the "deployerType" field.
func DeployerTypeHasPrefix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasPrefix(FieldDeployerType, v))
}

// DeployerTypeHasSuffix applies the HasSuffix predicate on the "deployerType" field.
func DeployerTypeHasSuffix(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldHasSuffix(FieldDeployerType, v))
}

// DeployerTypeEqualFold applies the EqualFold predicate on the "deployerType" field.
func DeployerTypeEqualFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEqualFold(FieldDeployerType, v))
}

// DeployerTypeContainsFold applies the ContainsFold predicate on the "deployerType" field.
func DeployerTypeContainsFold(v string) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldContainsFold(FieldDeployerType, v))
}

// DurationEQ applies the EQ predicate on the "duration" field.
func DurationEQ(v int) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldEQ(FieldDuration, v))
}

// DurationNEQ applies the NEQ predicate on the "duration" field.
func DurationNEQ(v int) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNEQ(FieldDuration, v))
}

// DurationIn applies the In predicate on the "duration" field.
func DurationIn(vs ...int) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldIn(FieldDuration, vs...))
}

// DurationNotIn applies the NotIn predicate on the "duration" field.
func DurationNotIn(vs ...int) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldNotIn(FieldDuration, vs...))
}

// DurationGT applies the GT predicate on the "duration" field.
func DurationGT(v int) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGT(FieldDuration, v))
}

// DurationGTE applies the GTE predicate on the "duration" field.
func DurationGTE(v int) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldGTE(FieldDuration, v))
}

// DurationLT applies the LT predicate on the "duration" field.
func DurationLT(v int) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLT(FieldDuration, v))
}

// DurationLTE applies the LTE predicate on the "duration" field.
func DurationLTE(v int) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(sql.FieldLTE(FieldDuration, v))
}

// HasInstance applies the HasEdge predicate on the "instance" edge.
func HasInstance() predicate.ApplicationRevision {
	return predicate.ApplicationRevision(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, InstanceTable, InstanceColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationInstance
		step.Edge.Schema = schemaConfig.ApplicationRevision
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasInstanceWith applies the HasEdge predicate on the "instance" edge with a given conditions (other predicates).
func HasInstanceWith(preds ...predicate.ApplicationInstance) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(InstanceInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, InstanceTable, InstanceColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ApplicationInstance
		step.Edge.Schema = schemaConfig.ApplicationRevision
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasEnvironment applies the HasEdge predicate on the "environment" edge.
func HasEnvironment() predicate.ApplicationRevision {
	return predicate.ApplicationRevision(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EnvironmentTable, EnvironmentColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.ApplicationRevision
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEnvironmentWith applies the HasEdge predicate on the "environment" edge with a given conditions (other predicates).
func HasEnvironmentWith(preds ...predicate.Environment) predicate.ApplicationRevision {
	return predicate.ApplicationRevision(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EnvironmentInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, EnvironmentTable, EnvironmentColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
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
