// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package workflowexecution

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/walrus/pkg/dao/model/internal"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// ID filters vertices based on their ID field.
func ID(id object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldDescription, v))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldUpdateTime, v))
}

// ProjectID applies equality check predicate on the "project_id" field. It's identical to ProjectIDEQ.
func ProjectID(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldProjectID, v))
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldVersion, v))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldType, v))
}

// WorkflowID applies equality check predicate on the "workflow_id" field. It's identical to WorkflowIDEQ.
func WorkflowID(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldWorkflowID, v))
}

// SubjectID applies equality check predicate on the "subject_id" field. It's identical to SubjectIDEQ.
func SubjectID(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldSubjectID, v))
}

// ExecuteTime applies equality check predicate on the "execute_time" field. It's identical to ExecuteTimeEQ.
func ExecuteTime(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldExecuteTime, v))
}

// Times applies equality check predicate on the "times" field. It's identical to TimesEQ.
func Times(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldTimes, v))
}

// Duration applies equality check predicate on the "duration" field. It's identical to DurationEQ.
func Duration(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldDuration, v))
}

// Parallelism applies equality check predicate on the "parallelism" field. It's identical to ParallelismEQ.
func Parallelism(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldParallelism, v))
}

// Timeout applies equality check predicate on the "timeout" field. It's identical to TimeoutEQ.
func Timeout(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldTimeout, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionIsNil applies the IsNil predicate on the "description" field.
func DescriptionIsNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIsNull(FieldDescription))
}

// DescriptionNotNil applies the NotNil predicate on the "description" field.
func DescriptionNotNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotNull(FieldDescription))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldContainsFold(FieldDescription, v))
}

// LabelsIsNil applies the IsNil predicate on the "labels" field.
func LabelsIsNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIsNull(FieldLabels))
}

// LabelsNotNil applies the NotNil predicate on the "labels" field.
func LabelsNotNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotNull(FieldLabels))
}

// AnnotationsIsNil applies the IsNil predicate on the "annotations" field.
func AnnotationsIsNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIsNull(FieldAnnotations))
}

// AnnotationsNotNil applies the NotNil predicate on the "annotations" field.
func AnnotationsNotNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotNull(FieldAnnotations))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldUpdateTime, v))
}

// StatusIsNil applies the IsNil predicate on the "status" field.
func StatusIsNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIsNull(FieldStatus))
}

// StatusNotNil applies the NotNil predicate on the "status" field.
func StatusNotNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotNull(FieldStatus))
}

// ProjectIDEQ applies the EQ predicate on the "project_id" field.
func ProjectIDEQ(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldProjectID, v))
}

// ProjectIDNEQ applies the NEQ predicate on the "project_id" field.
func ProjectIDNEQ(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldProjectID, v))
}

// ProjectIDIn applies the In predicate on the "project_id" field.
func ProjectIDIn(vs ...object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldProjectID, vs...))
}

// ProjectIDNotIn applies the NotIn predicate on the "project_id" field.
func ProjectIDNotIn(vs ...object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldProjectID, vs...))
}

// ProjectIDGT applies the GT predicate on the "project_id" field.
func ProjectIDGT(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldProjectID, v))
}

// ProjectIDGTE applies the GTE predicate on the "project_id" field.
func ProjectIDGTE(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldProjectID, v))
}

// ProjectIDLT applies the LT predicate on the "project_id" field.
func ProjectIDLT(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldProjectID, v))
}

// ProjectIDLTE applies the LTE predicate on the "project_id" field.
func ProjectIDLTE(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldProjectID, v))
}

// ProjectIDContains applies the Contains predicate on the "project_id" field.
func ProjectIDContains(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldContains(FieldProjectID, vc))
}

// ProjectIDHasPrefix applies the HasPrefix predicate on the "project_id" field.
func ProjectIDHasPrefix(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldHasPrefix(FieldProjectID, vc))
}

// ProjectIDHasSuffix applies the HasSuffix predicate on the "project_id" field.
func ProjectIDHasSuffix(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldHasSuffix(FieldProjectID, vc))
}

// ProjectIDEqualFold applies the EqualFold predicate on the "project_id" field.
func ProjectIDEqualFold(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldEqualFold(FieldProjectID, vc))
}

// ProjectIDContainsFold applies the ContainsFold predicate on the "project_id" field.
func ProjectIDContainsFold(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldContainsFold(FieldProjectID, vc))
}

// VersionEQ applies the EQ predicate on the "version" field.
func VersionEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldVersion, v))
}

// VersionNEQ applies the NEQ predicate on the "version" field.
func VersionNEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldVersion, v))
}

// VersionIn applies the In predicate on the "version" field.
func VersionIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldVersion, vs...))
}

// VersionNotIn applies the NotIn predicate on the "version" field.
func VersionNotIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldVersion, vs...))
}

// VersionGT applies the GT predicate on the "version" field.
func VersionGT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldVersion, v))
}

// VersionGTE applies the GTE predicate on the "version" field.
func VersionGTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldVersion, v))
}

// VersionLT applies the LT predicate on the "version" field.
func VersionLT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldVersion, v))
}

// VersionLTE applies the LTE predicate on the "version" field.
func VersionLTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldVersion, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldContainsFold(FieldType, v))
}

// WorkflowIDEQ applies the EQ predicate on the "workflow_id" field.
func WorkflowIDEQ(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldWorkflowID, v))
}

// WorkflowIDNEQ applies the NEQ predicate on the "workflow_id" field.
func WorkflowIDNEQ(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldWorkflowID, v))
}

// WorkflowIDIn applies the In predicate on the "workflow_id" field.
func WorkflowIDIn(vs ...object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldWorkflowID, vs...))
}

// WorkflowIDNotIn applies the NotIn predicate on the "workflow_id" field.
func WorkflowIDNotIn(vs ...object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldWorkflowID, vs...))
}

// WorkflowIDGT applies the GT predicate on the "workflow_id" field.
func WorkflowIDGT(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldWorkflowID, v))
}

// WorkflowIDGTE applies the GTE predicate on the "workflow_id" field.
func WorkflowIDGTE(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldWorkflowID, v))
}

// WorkflowIDLT applies the LT predicate on the "workflow_id" field.
func WorkflowIDLT(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldWorkflowID, v))
}

// WorkflowIDLTE applies the LTE predicate on the "workflow_id" field.
func WorkflowIDLTE(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldWorkflowID, v))
}

// WorkflowIDContains applies the Contains predicate on the "workflow_id" field.
func WorkflowIDContains(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldContains(FieldWorkflowID, vc))
}

// WorkflowIDHasPrefix applies the HasPrefix predicate on the "workflow_id" field.
func WorkflowIDHasPrefix(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldHasPrefix(FieldWorkflowID, vc))
}

// WorkflowIDHasSuffix applies the HasSuffix predicate on the "workflow_id" field.
func WorkflowIDHasSuffix(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldHasSuffix(FieldWorkflowID, vc))
}

// WorkflowIDEqualFold applies the EqualFold predicate on the "workflow_id" field.
func WorkflowIDEqualFold(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldEqualFold(FieldWorkflowID, vc))
}

// WorkflowIDContainsFold applies the ContainsFold predicate on the "workflow_id" field.
func WorkflowIDContainsFold(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldContainsFold(FieldWorkflowID, vc))
}

// SubjectIDEQ applies the EQ predicate on the "subject_id" field.
func SubjectIDEQ(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldSubjectID, v))
}

// SubjectIDNEQ applies the NEQ predicate on the "subject_id" field.
func SubjectIDNEQ(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldSubjectID, v))
}

// SubjectIDIn applies the In predicate on the "subject_id" field.
func SubjectIDIn(vs ...object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldSubjectID, vs...))
}

// SubjectIDNotIn applies the NotIn predicate on the "subject_id" field.
func SubjectIDNotIn(vs ...object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldSubjectID, vs...))
}

// SubjectIDGT applies the GT predicate on the "subject_id" field.
func SubjectIDGT(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldSubjectID, v))
}

// SubjectIDGTE applies the GTE predicate on the "subject_id" field.
func SubjectIDGTE(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldSubjectID, v))
}

// SubjectIDLT applies the LT predicate on the "subject_id" field.
func SubjectIDLT(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldSubjectID, v))
}

// SubjectIDLTE applies the LTE predicate on the "subject_id" field.
func SubjectIDLTE(v object.ID) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldSubjectID, v))
}

// SubjectIDContains applies the Contains predicate on the "subject_id" field.
func SubjectIDContains(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldContains(FieldSubjectID, vc))
}

// SubjectIDHasPrefix applies the HasPrefix predicate on the "subject_id" field.
func SubjectIDHasPrefix(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldHasPrefix(FieldSubjectID, vc))
}

// SubjectIDHasSuffix applies the HasSuffix predicate on the "subject_id" field.
func SubjectIDHasSuffix(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldHasSuffix(FieldSubjectID, vc))
}

// SubjectIDEqualFold applies the EqualFold predicate on the "subject_id" field.
func SubjectIDEqualFold(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldEqualFold(FieldSubjectID, vc))
}

// SubjectIDContainsFold applies the ContainsFold predicate on the "subject_id" field.
func SubjectIDContainsFold(v object.ID) predicate.WorkflowExecution {
	vc := string(v)
	return predicate.WorkflowExecution(sql.FieldContainsFold(FieldSubjectID, vc))
}

// ExecuteTimeEQ applies the EQ predicate on the "execute_time" field.
func ExecuteTimeEQ(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldExecuteTime, v))
}

// ExecuteTimeNEQ applies the NEQ predicate on the "execute_time" field.
func ExecuteTimeNEQ(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldExecuteTime, v))
}

// ExecuteTimeIn applies the In predicate on the "execute_time" field.
func ExecuteTimeIn(vs ...time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldExecuteTime, vs...))
}

// ExecuteTimeNotIn applies the NotIn predicate on the "execute_time" field.
func ExecuteTimeNotIn(vs ...time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldExecuteTime, vs...))
}

// ExecuteTimeGT applies the GT predicate on the "execute_time" field.
func ExecuteTimeGT(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldExecuteTime, v))
}

// ExecuteTimeGTE applies the GTE predicate on the "execute_time" field.
func ExecuteTimeGTE(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldExecuteTime, v))
}

// ExecuteTimeLT applies the LT predicate on the "execute_time" field.
func ExecuteTimeLT(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldExecuteTime, v))
}

// ExecuteTimeLTE applies the LTE predicate on the "execute_time" field.
func ExecuteTimeLTE(v time.Time) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldExecuteTime, v))
}

// ExecuteTimeIsNil applies the IsNil predicate on the "execute_time" field.
func ExecuteTimeIsNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIsNull(FieldExecuteTime))
}

// ExecuteTimeNotNil applies the NotNil predicate on the "execute_time" field.
func ExecuteTimeNotNil() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotNull(FieldExecuteTime))
}

// TimesEQ applies the EQ predicate on the "times" field.
func TimesEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldTimes, v))
}

// TimesNEQ applies the NEQ predicate on the "times" field.
func TimesNEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldTimes, v))
}

// TimesIn applies the In predicate on the "times" field.
func TimesIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldTimes, vs...))
}

// TimesNotIn applies the NotIn predicate on the "times" field.
func TimesNotIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldTimes, vs...))
}

// TimesGT applies the GT predicate on the "times" field.
func TimesGT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldTimes, v))
}

// TimesGTE applies the GTE predicate on the "times" field.
func TimesGTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldTimes, v))
}

// TimesLT applies the LT predicate on the "times" field.
func TimesLT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldTimes, v))
}

// TimesLTE applies the LTE predicate on the "times" field.
func TimesLTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldTimes, v))
}

// DurationEQ applies the EQ predicate on the "duration" field.
func DurationEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldDuration, v))
}

// DurationNEQ applies the NEQ predicate on the "duration" field.
func DurationNEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldDuration, v))
}

// DurationIn applies the In predicate on the "duration" field.
func DurationIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldDuration, vs...))
}

// DurationNotIn applies the NotIn predicate on the "duration" field.
func DurationNotIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldDuration, vs...))
}

// DurationGT applies the GT predicate on the "duration" field.
func DurationGT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldDuration, v))
}

// DurationGTE applies the GTE predicate on the "duration" field.
func DurationGTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldDuration, v))
}

// DurationLT applies the LT predicate on the "duration" field.
func DurationLT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldDuration, v))
}

// DurationLTE applies the LTE predicate on the "duration" field.
func DurationLTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldDuration, v))
}

// ParallelismEQ applies the EQ predicate on the "parallelism" field.
func ParallelismEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldParallelism, v))
}

// ParallelismNEQ applies the NEQ predicate on the "parallelism" field.
func ParallelismNEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldParallelism, v))
}

// ParallelismIn applies the In predicate on the "parallelism" field.
func ParallelismIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldParallelism, vs...))
}

// ParallelismNotIn applies the NotIn predicate on the "parallelism" field.
func ParallelismNotIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldParallelism, vs...))
}

// ParallelismGT applies the GT predicate on the "parallelism" field.
func ParallelismGT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldParallelism, v))
}

// ParallelismGTE applies the GTE predicate on the "parallelism" field.
func ParallelismGTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldParallelism, v))
}

// ParallelismLT applies the LT predicate on the "parallelism" field.
func ParallelismLT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldParallelism, v))
}

// ParallelismLTE applies the LTE predicate on the "parallelism" field.
func ParallelismLTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldParallelism, v))
}

// TimeoutEQ applies the EQ predicate on the "timeout" field.
func TimeoutEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldEQ(FieldTimeout, v))
}

// TimeoutNEQ applies the NEQ predicate on the "timeout" field.
func TimeoutNEQ(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNEQ(FieldTimeout, v))
}

// TimeoutIn applies the In predicate on the "timeout" field.
func TimeoutIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldIn(FieldTimeout, vs...))
}

// TimeoutNotIn applies the NotIn predicate on the "timeout" field.
func TimeoutNotIn(vs ...int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldNotIn(FieldTimeout, vs...))
}

// TimeoutGT applies the GT predicate on the "timeout" field.
func TimeoutGT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGT(FieldTimeout, v))
}

// TimeoutGTE applies the GTE predicate on the "timeout" field.
func TimeoutGTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldGTE(FieldTimeout, v))
}

// TimeoutLT applies the LT predicate on the "timeout" field.
func TimeoutLT(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLT(FieldTimeout, v))
}

// TimeoutLTE applies the LTE predicate on the "timeout" field.
func TimeoutLTE(v int) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.FieldLTE(FieldTimeout, v))
}

// HasProject applies the HasEdge predicate on the "project" edge.
func HasProject() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.WorkflowExecution
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectWith applies the HasEdge predicate on the "project" edge with a given conditions (other predicates).
func HasProjectWith(preds ...predicate.Project) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(func(s *sql.Selector) {
		step := newProjectStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.WorkflowExecution
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasStages applies the HasEdge predicate on the "stages" edge.
func HasStages() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, StagesTable, StagesColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.WorkflowStageExecution
		step.Edge.Schema = schemaConfig.WorkflowStageExecution
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStagesWith applies the HasEdge predicate on the "stages" edge with a given conditions (other predicates).
func HasStagesWith(preds ...predicate.WorkflowStageExecution) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(func(s *sql.Selector) {
		step := newStagesStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.WorkflowStageExecution
		step.Edge.Schema = schemaConfig.WorkflowStageExecution
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWorkflow applies the HasEdge predicate on the "workflow" edge.
func HasWorkflow() predicate.WorkflowExecution {
	return predicate.WorkflowExecution(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, WorkflowTable, WorkflowColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Workflow
		step.Edge.Schema = schemaConfig.WorkflowExecution
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWorkflowWith applies the HasEdge predicate on the "workflow" edge with a given conditions (other predicates).
func HasWorkflowWith(preds ...predicate.Workflow) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(func(s *sql.Selector) {
		step := newWorkflowStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Workflow
		step.Edge.Schema = schemaConfig.WorkflowExecution
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.WorkflowExecution) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.WorkflowExecution) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.WorkflowExecution) predicate.WorkflowExecution {
	return predicate.WorkflowExecution(sql.NotPredicates(p))
}
