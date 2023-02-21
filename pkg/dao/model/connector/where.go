// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package connector

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ID filters vertices based on their ID field.
func ID(id types.ID) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id types.ID) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id types.ID) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...types.ID) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...types.ID) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id types.ID) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id types.ID) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id types.ID) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id types.ID) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldDescription, v))
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldStatus, v))
}

// StatusMessage applies equality check predicate on the "statusMessage" field. It's identical to StatusMessageEQ.
func StatusMessage(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldStatusMessage, v))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldUpdateTime, v))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldType, v))
}

// ConfigVersion applies equality check predicate on the "configVersion" field. It's identical to ConfigVersionEQ.
func ConfigVersion(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldConfigVersion, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.Connector {
	return predicate.Connector(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.Connector {
	return predicate.Connector(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldDescription, v))
}

// LabelsIsNil applies the IsNil predicate on the "labels" field.
func LabelsIsNil() predicate.Connector {
	return predicate.Connector(sql.FieldIsNull(FieldLabels))
}

// LabelsNotNil applies the NotNil predicate on the "labels" field.
func LabelsNotNil() predicate.Connector {
	return predicate.Connector(sql.FieldNotNull(FieldLabels))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldStatus, vs...))
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldStatus, v))
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldStatus, v))
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldStatus, v))
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldStatus, v))
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldStatus, v))
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldStatus, v))
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldStatus, v))
}

// StatusIsNil applies the IsNil predicate on the "status" field.
func StatusIsNil() predicate.Connector {
	return predicate.Connector(sql.FieldIsNull(FieldStatus))
}

// StatusNotNil applies the NotNil predicate on the "status" field.
func StatusNotNil() predicate.Connector {
	return predicate.Connector(sql.FieldNotNull(FieldStatus))
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldStatus, v))
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldStatus, v))
}

// StatusMessageEQ applies the EQ predicate on the "statusMessage" field.
func StatusMessageEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldStatusMessage, v))
}

// StatusMessageNEQ applies the NEQ predicate on the "statusMessage" field.
func StatusMessageNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldStatusMessage, v))
}

// StatusMessageIn applies the In predicate on the "statusMessage" field.
func StatusMessageIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldStatusMessage, vs...))
}

// StatusMessageNotIn applies the NotIn predicate on the "statusMessage" field.
func StatusMessageNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldStatusMessage, vs...))
}

// StatusMessageGT applies the GT predicate on the "statusMessage" field.
func StatusMessageGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldStatusMessage, v))
}

// StatusMessageGTE applies the GTE predicate on the "statusMessage" field.
func StatusMessageGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldStatusMessage, v))
}

// StatusMessageLT applies the LT predicate on the "statusMessage" field.
func StatusMessageLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldStatusMessage, v))
}

// StatusMessageLTE applies the LTE predicate on the "statusMessage" field.
func StatusMessageLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldStatusMessage, v))
}

// StatusMessageContains applies the Contains predicate on the "statusMessage" field.
func StatusMessageContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldStatusMessage, v))
}

// StatusMessageHasPrefix applies the HasPrefix predicate on the "statusMessage" field.
func StatusMessageHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldStatusMessage, v))
}

// StatusMessageHasSuffix applies the HasSuffix predicate on the "statusMessage" field.
func StatusMessageHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldStatusMessage, v))
}

// StatusMessageIsNil applies the IsNil predicate on the "statusMessage" field.
func StatusMessageIsNil() predicate.Connector {
	return predicate.Connector(sql.FieldIsNull(FieldStatusMessage))
}

// StatusMessageNotNil applies the NotNil predicate on the "statusMessage" field.
func StatusMessageNotNil() predicate.Connector {
	return predicate.Connector(sql.FieldNotNull(FieldStatusMessage))
}

// StatusMessageEqualFold applies the EqualFold predicate on the "statusMessage" field.
func StatusMessageEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldStatusMessage, v))
}

// StatusMessageContainsFold applies the ContainsFold predicate on the "statusMessage" field.
func StatusMessageContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldStatusMessage, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldUpdateTime, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldType, v))
}

// ConfigVersionEQ applies the EQ predicate on the "configVersion" field.
func ConfigVersionEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEQ(FieldConfigVersion, v))
}

// ConfigVersionNEQ applies the NEQ predicate on the "configVersion" field.
func ConfigVersionNEQ(v string) predicate.Connector {
	return predicate.Connector(sql.FieldNEQ(FieldConfigVersion, v))
}

// ConfigVersionIn applies the In predicate on the "configVersion" field.
func ConfigVersionIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldIn(FieldConfigVersion, vs...))
}

// ConfigVersionNotIn applies the NotIn predicate on the "configVersion" field.
func ConfigVersionNotIn(vs ...string) predicate.Connector {
	return predicate.Connector(sql.FieldNotIn(FieldConfigVersion, vs...))
}

// ConfigVersionGT applies the GT predicate on the "configVersion" field.
func ConfigVersionGT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGT(FieldConfigVersion, v))
}

// ConfigVersionGTE applies the GTE predicate on the "configVersion" field.
func ConfigVersionGTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldGTE(FieldConfigVersion, v))
}

// ConfigVersionLT applies the LT predicate on the "configVersion" field.
func ConfigVersionLT(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLT(FieldConfigVersion, v))
}

// ConfigVersionLTE applies the LTE predicate on the "configVersion" field.
func ConfigVersionLTE(v string) predicate.Connector {
	return predicate.Connector(sql.FieldLTE(FieldConfigVersion, v))
}

// ConfigVersionContains applies the Contains predicate on the "configVersion" field.
func ConfigVersionContains(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContains(FieldConfigVersion, v))
}

// ConfigVersionHasPrefix applies the HasPrefix predicate on the "configVersion" field.
func ConfigVersionHasPrefix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasPrefix(FieldConfigVersion, v))
}

// ConfigVersionHasSuffix applies the HasSuffix predicate on the "configVersion" field.
func ConfigVersionHasSuffix(v string) predicate.Connector {
	return predicate.Connector(sql.FieldHasSuffix(FieldConfigVersion, v))
}

// ConfigVersionEqualFold applies the EqualFold predicate on the "configVersion" field.
func ConfigVersionEqualFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldEqualFold(FieldConfigVersion, v))
}

// ConfigVersionContainsFold applies the ContainsFold predicate on the "configVersion" field.
func ConfigVersionContainsFold(v string) predicate.Connector {
	return predicate.Connector(sql.FieldContainsFold(FieldConfigVersion, v))
}

// ConfigDataIsNil applies the IsNil predicate on the "configData" field.
func ConfigDataIsNil() predicate.Connector {
	return predicate.Connector(sql.FieldIsNull(FieldConfigData))
}

// ConfigDataNotNil applies the NotNil predicate on the "configData" field.
func ConfigDataNotNil() predicate.Connector {
	return predicate.Connector(sql.FieldNotNull(FieldConfigData))
}

// HasEnvironments applies the HasEdge predicate on the "environments" edge.
func HasEnvironments() predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, EnvironmentsTable, EnvironmentsPrimaryKey...),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEnvironmentsWith applies the HasEdge predicate on the "environments" edge with a given conditions (other predicates).
func HasEnvironmentsWith(preds ...predicate.Environment) predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EnvironmentsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, EnvironmentsTable, EnvironmentsPrimaryKey...),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasResources applies the HasEdge predicate on the "resources" edge.
func HasResources() predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
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
func HasResourcesWith(preds ...predicate.ApplicationResource) predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
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

// HasEnvironmentConnectorRelationships applies the HasEdge predicate on the "environmentConnectorRelationships" edge.
func HasEnvironmentConnectorRelationships() predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, EnvironmentConnectorRelationshipsTable, EnvironmentConnectorRelationshipsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.EnvironmentConnectorRelationship
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEnvironmentConnectorRelationshipsWith applies the HasEdge predicate on the "environmentConnectorRelationships" edge with a given conditions (other predicates).
func HasEnvironmentConnectorRelationshipsWith(preds ...predicate.EnvironmentConnectorRelationship) predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(EnvironmentConnectorRelationshipsInverseTable, EnvironmentConnectorRelationshipsColumn),
			sqlgraph.Edge(sqlgraph.O2M, true, EnvironmentConnectorRelationshipsTable, EnvironmentConnectorRelationshipsColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.EnvironmentConnectorRelationship
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Connector) predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Connector) predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
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
func Not(p predicate.Connector) predicate.Connector {
	return predicate.Connector(func(s *sql.Selector) {
		p(s.Not())
	})
}
