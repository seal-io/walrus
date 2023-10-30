// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package workflowstep

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/walrus/pkg/dao/model/internal"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// ID filters vertices based on their ID field.
func ID(id object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldDescription, v))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldUpdateTime, v))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldType, v))
}

// ProjectID applies equality check predicate on the "project_id" field. It's identical to ProjectIDEQ.
func ProjectID(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldProjectID, v))
}

// WorkflowID applies equality check predicate on the "workflow_id" field. It's identical to WorkflowIDEQ.
func WorkflowID(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldWorkflowID, v))
}

// WorkflowStageID applies equality check predicate on the "workflow_stage_id" field. It's identical to WorkflowStageIDEQ.
func WorkflowStageID(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldWorkflowStageID, v))
}

// Order applies equality check predicate on the "order" field. It's identical to OrderEQ.
func Order(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldOrder, v))
}

// Timeout applies equality check predicate on the "timeout" field. It's identical to TimeoutEQ.
func Timeout(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldTimeout, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldContainsFold(FieldDescription, v))
}

// LabelsIsNil applies the IsNil predicate on the "labels" field.
func LabelsIsNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIsNull(FieldLabels))
}

// LabelsNotNil applies the NotNil predicate on the "labels" field.
func LabelsNotNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotNull(FieldLabels))
}

// AnnotationsIsNil applies the IsNil predicate on the "annotations" field.
func AnnotationsIsNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIsNull(FieldAnnotations))
}

// AnnotationsNotNil applies the NotNil predicate on the "annotations" field.
func AnnotationsNotNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotNull(FieldAnnotations))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldUpdateTime, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldContainsFold(FieldType, v))
}

// ProjectIDEQ applies the EQ predicate on the "project_id" field.
func ProjectIDEQ(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldProjectID, v))
}

// ProjectIDNEQ applies the NEQ predicate on the "project_id" field.
func ProjectIDNEQ(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldProjectID, v))
}

// ProjectIDIn applies the In predicate on the "project_id" field.
func ProjectIDIn(vs ...object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldProjectID, vs...))
}

// ProjectIDNotIn applies the NotIn predicate on the "project_id" field.
func ProjectIDNotIn(vs ...object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldProjectID, vs...))
}

// ProjectIDGT applies the GT predicate on the "project_id" field.
func ProjectIDGT(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldProjectID, v))
}

// ProjectIDGTE applies the GTE predicate on the "project_id" field.
func ProjectIDGTE(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldProjectID, v))
}

// ProjectIDLT applies the LT predicate on the "project_id" field.
func ProjectIDLT(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldProjectID, v))
}

// ProjectIDLTE applies the LTE predicate on the "project_id" field.
func ProjectIDLTE(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldProjectID, v))
}

// ProjectIDContains applies the Contains predicate on the "project_id" field.
func ProjectIDContains(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldContains(FieldProjectID, vc))
}

// ProjectIDHasPrefix applies the HasPrefix predicate on the "project_id" field.
func ProjectIDHasPrefix(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldHasPrefix(FieldProjectID, vc))
}

// ProjectIDHasSuffix applies the HasSuffix predicate on the "project_id" field.
func ProjectIDHasSuffix(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldHasSuffix(FieldProjectID, vc))
}

// ProjectIDEqualFold applies the EqualFold predicate on the "project_id" field.
func ProjectIDEqualFold(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldEqualFold(FieldProjectID, vc))
}

// ProjectIDContainsFold applies the ContainsFold predicate on the "project_id" field.
func ProjectIDContainsFold(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldContainsFold(FieldProjectID, vc))
}

// WorkflowIDEQ applies the EQ predicate on the "workflow_id" field.
func WorkflowIDEQ(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldWorkflowID, v))
}

// WorkflowIDNEQ applies the NEQ predicate on the "workflow_id" field.
func WorkflowIDNEQ(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldWorkflowID, v))
}

// WorkflowIDIn applies the In predicate on the "workflow_id" field.
func WorkflowIDIn(vs ...object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldWorkflowID, vs...))
}

// WorkflowIDNotIn applies the NotIn predicate on the "workflow_id" field.
func WorkflowIDNotIn(vs ...object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldWorkflowID, vs...))
}

// WorkflowIDGT applies the GT predicate on the "workflow_id" field.
func WorkflowIDGT(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldWorkflowID, v))
}

// WorkflowIDGTE applies the GTE predicate on the "workflow_id" field.
func WorkflowIDGTE(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldWorkflowID, v))
}

// WorkflowIDLT applies the LT predicate on the "workflow_id" field.
func WorkflowIDLT(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldWorkflowID, v))
}

// WorkflowIDLTE applies the LTE predicate on the "workflow_id" field.
func WorkflowIDLTE(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldWorkflowID, v))
}

// WorkflowIDContains applies the Contains predicate on the "workflow_id" field.
func WorkflowIDContains(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldContains(FieldWorkflowID, vc))
}

// WorkflowIDHasPrefix applies the HasPrefix predicate on the "workflow_id" field.
func WorkflowIDHasPrefix(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldHasPrefix(FieldWorkflowID, vc))
}

// WorkflowIDHasSuffix applies the HasSuffix predicate on the "workflow_id" field.
func WorkflowIDHasSuffix(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldHasSuffix(FieldWorkflowID, vc))
}

// WorkflowIDEqualFold applies the EqualFold predicate on the "workflow_id" field.
func WorkflowIDEqualFold(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldEqualFold(FieldWorkflowID, vc))
}

// WorkflowIDContainsFold applies the ContainsFold predicate on the "workflow_id" field.
func WorkflowIDContainsFold(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldContainsFold(FieldWorkflowID, vc))
}

// WorkflowStageIDEQ applies the EQ predicate on the "workflow_stage_id" field.
func WorkflowStageIDEQ(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldWorkflowStageID, v))
}

// WorkflowStageIDNEQ applies the NEQ predicate on the "workflow_stage_id" field.
func WorkflowStageIDNEQ(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldWorkflowStageID, v))
}

// WorkflowStageIDIn applies the In predicate on the "workflow_stage_id" field.
func WorkflowStageIDIn(vs ...object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldWorkflowStageID, vs...))
}

// WorkflowStageIDNotIn applies the NotIn predicate on the "workflow_stage_id" field.
func WorkflowStageIDNotIn(vs ...object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldWorkflowStageID, vs...))
}

// WorkflowStageIDGT applies the GT predicate on the "workflow_stage_id" field.
func WorkflowStageIDGT(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldWorkflowStageID, v))
}

// WorkflowStageIDGTE applies the GTE predicate on the "workflow_stage_id" field.
func WorkflowStageIDGTE(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldWorkflowStageID, v))
}

// WorkflowStageIDLT applies the LT predicate on the "workflow_stage_id" field.
func WorkflowStageIDLT(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldWorkflowStageID, v))
}

// WorkflowStageIDLTE applies the LTE predicate on the "workflow_stage_id" field.
func WorkflowStageIDLTE(v object.ID) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldWorkflowStageID, v))
}

// WorkflowStageIDContains applies the Contains predicate on the "workflow_stage_id" field.
func WorkflowStageIDContains(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldContains(FieldWorkflowStageID, vc))
}

// WorkflowStageIDHasPrefix applies the HasPrefix predicate on the "workflow_stage_id" field.
func WorkflowStageIDHasPrefix(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldHasPrefix(FieldWorkflowStageID, vc))
}

// WorkflowStageIDHasSuffix applies the HasSuffix predicate on the "workflow_stage_id" field.
func WorkflowStageIDHasSuffix(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldHasSuffix(FieldWorkflowStageID, vc))
}

// WorkflowStageIDEqualFold applies the EqualFold predicate on the "workflow_stage_id" field.
func WorkflowStageIDEqualFold(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldEqualFold(FieldWorkflowStageID, vc))
}

// WorkflowStageIDContainsFold applies the ContainsFold predicate on the "workflow_stage_id" field.
func WorkflowStageIDContainsFold(v object.ID) predicate.WorkflowStep {
	vc := string(v)
	return predicate.WorkflowStep(sql.FieldContainsFold(FieldWorkflowStageID, vc))
}

// AttributesIsNil applies the IsNil predicate on the "attributes" field.
func AttributesIsNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIsNull(FieldAttributes))
}

// AttributesNotNil applies the NotNil predicate on the "attributes" field.
func AttributesNotNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotNull(FieldAttributes))
}

// InputsIsNil applies the IsNil predicate on the "inputs" field.
func InputsIsNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIsNull(FieldInputs))
}

// InputsNotNil applies the NotNil predicate on the "inputs" field.
func InputsNotNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotNull(FieldInputs))
}

// OutputsIsNil applies the IsNil predicate on the "outputs" field.
func OutputsIsNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIsNull(FieldOutputs))
}

// OutputsNotNil applies the NotNil predicate on the "outputs" field.
func OutputsNotNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotNull(FieldOutputs))
}

// OrderEQ applies the EQ predicate on the "order" field.
func OrderEQ(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldOrder, v))
}

// OrderNEQ applies the NEQ predicate on the "order" field.
func OrderNEQ(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldOrder, v))
}

// OrderIn applies the In predicate on the "order" field.
func OrderIn(vs ...int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldOrder, vs...))
}

// OrderNotIn applies the NotIn predicate on the "order" field.
func OrderNotIn(vs ...int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldOrder, vs...))
}

// OrderGT applies the GT predicate on the "order" field.
func OrderGT(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldOrder, v))
}

// OrderGTE applies the GTE predicate on the "order" field.
func OrderGTE(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldOrder, v))
}

// OrderLT applies the LT predicate on the "order" field.
func OrderLT(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldOrder, v))
}

// OrderLTE applies the LTE predicate on the "order" field.
func OrderLTE(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldOrder, v))
}

// RetryStrategyIsNil applies the IsNil predicate on the "retryStrategy" field.
func RetryStrategyIsNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIsNull(FieldRetryStrategy))
}

// RetryStrategyNotNil applies the NotNil predicate on the "retryStrategy" field.
func RetryStrategyNotNil() predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotNull(FieldRetryStrategy))
}

// TimeoutEQ applies the EQ predicate on the "timeout" field.
func TimeoutEQ(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldEQ(FieldTimeout, v))
}

// TimeoutNEQ applies the NEQ predicate on the "timeout" field.
func TimeoutNEQ(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNEQ(FieldTimeout, v))
}

// TimeoutIn applies the In predicate on the "timeout" field.
func TimeoutIn(vs ...int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldIn(FieldTimeout, vs...))
}

// TimeoutNotIn applies the NotIn predicate on the "timeout" field.
func TimeoutNotIn(vs ...int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldNotIn(FieldTimeout, vs...))
}

// TimeoutGT applies the GT predicate on the "timeout" field.
func TimeoutGT(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGT(FieldTimeout, v))
}

// TimeoutGTE applies the GTE predicate on the "timeout" field.
func TimeoutGTE(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldGTE(FieldTimeout, v))
}

// TimeoutLT applies the LT predicate on the "timeout" field.
func TimeoutLT(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLT(FieldTimeout, v))
}

// TimeoutLTE applies the LTE predicate on the "timeout" field.
func TimeoutLTE(v int) predicate.WorkflowStep {
	return predicate.WorkflowStep(sql.FieldLTE(FieldTimeout, v))
}

// HasProject applies the HasEdge predicate on the "project" edge.
func HasProject() predicate.WorkflowStep {
	return predicate.WorkflowStep(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.WorkflowStep
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectWith applies the HasEdge predicate on the "project" edge with a given conditions (other predicates).
func HasProjectWith(preds ...predicate.Project) predicate.WorkflowStep {
	return predicate.WorkflowStep(func(s *sql.Selector) {
		step := newProjectStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.WorkflowStep
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasStage applies the HasEdge predicate on the "stage" edge.
func HasStage() predicate.WorkflowStep {
	return predicate.WorkflowStep(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, StageTable, StageColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.WorkflowStage
		step.Edge.Schema = schemaConfig.WorkflowStep
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStageWith applies the HasEdge predicate on the "stage" edge with a given conditions (other predicates).
func HasStageWith(preds ...predicate.WorkflowStage) predicate.WorkflowStep {
	return predicate.WorkflowStep(func(s *sql.Selector) {
		step := newStageStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.WorkflowStage
		step.Edge.Schema = schemaConfig.WorkflowStep
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.WorkflowStep) predicate.WorkflowStep {
	return predicate.WorkflowStep(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.WorkflowStep) predicate.WorkflowStep {
	return predicate.WorkflowStep(func(s *sql.Selector) {
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
func Not(p predicate.WorkflowStep) predicate.WorkflowStep {
	return predicate.WorkflowStep(func(s *sql.Selector) {
		p(s.Not())
	})
}