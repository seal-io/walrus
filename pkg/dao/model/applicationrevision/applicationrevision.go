// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package applicationrevision

import (
	"time"

	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/types"
)

const (
	// Label holds the string label denoting the applicationrevision type in the database.
	Label = "application_revision"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldStatusMessage holds the string denoting the statusmessage field in the database.
	FieldStatusMessage = "status_message"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldApplicationID holds the string denoting the applicationid field in the database.
	FieldApplicationID = "application_id"
	// FieldEnvironmentID holds the string denoting the environmentid field in the database.
	FieldEnvironmentID = "environment_id"
	// FieldModules holds the string denoting the modules field in the database.
	FieldModules = "modules"
	// FieldInputVariables holds the string denoting the inputvariables field in the database.
	FieldInputVariables = "input_variables"
	// FieldInputPlan holds the string denoting the inputplan field in the database.
	FieldInputPlan = "input_plan"
	// FieldOutput holds the string denoting the output field in the database.
	FieldOutput = "output"
	// FieldDeployerType holds the string denoting the deployertype field in the database.
	FieldDeployerType = "deployer_type"
	// FieldDuration holds the string denoting the duration field in the database.
	FieldDuration = "duration"
	// EdgeApplication holds the string denoting the application edge name in mutations.
	EdgeApplication = "application"
	// EdgeEnvironment holds the string denoting the environment edge name in mutations.
	EdgeEnvironment = "environment"
	// Table holds the table name of the applicationrevision in the database.
	Table = "application_revisions"
	// ApplicationTable is the table that holds the application relation/edge.
	ApplicationTable = "application_revisions"
	// ApplicationInverseTable is the table name for the Application entity.
	// It exists in this package in order to avoid circular dependency with the "application" package.
	ApplicationInverseTable = "applications"
	// ApplicationColumn is the table column denoting the application relation/edge.
	ApplicationColumn = "application_id"
	// EnvironmentTable is the table that holds the environment relation/edge.
	EnvironmentTable = "application_revisions"
	// EnvironmentInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	EnvironmentInverseTable = "environments"
	// EnvironmentColumn is the table column denoting the environment relation/edge.
	EnvironmentColumn = "environment_id"
)

// Columns holds all SQL columns for applicationrevision fields.
var Columns = []string{
	FieldID,
	FieldStatus,
	FieldStatusMessage,
	FieldCreateTime,
	FieldUpdateTime,
	FieldApplicationID,
	FieldEnvironmentID,
	FieldModules,
	FieldInputVariables,
	FieldInputPlan,
	FieldOutput,
	FieldDeployerType,
	FieldDuration,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/seal-io/seal/pkg/dao/model/runtime"
var (
	Hooks [1]ent.Hook
	// DefaultCreateTime holds the default value on creation for the "createTime" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "updateTime" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "updateTime" field.
	UpdateDefaultUpdateTime func() time.Time
	// ApplicationIDValidator is a validator for the "applicationID" field. It is called by the builders before save.
	ApplicationIDValidator func(string) error
	// EnvironmentIDValidator is a validator for the "environmentID" field. It is called by the builders before save.
	EnvironmentIDValidator func(string) error
	// DefaultModules holds the default value on creation for the "modules" field.
	DefaultModules []types.ApplicationModule
	// DefaultInputVariables holds the default value on creation for the "inputVariables" field.
	DefaultInputVariables map[string]interface{}
	// DefaultDeployerType holds the default value on creation for the "deployerType" field.
	DefaultDeployerType string
	// DefaultDuration holds the default value on creation for the "duration" field.
	DefaultDuration int
)

// WithoutFields returns the fields ignored the given list.
func WithoutFields(ignores ...string) []string {
	if len(ignores) == 0 {
		return Columns
	}

	var s = make(map[string]bool, len(ignores))
	for i := range ignores {
		s[ignores[i]] = true
	}

	var r = make([]string, 0, len(Columns)-len(s))
	for i := range Columns {
		if s[Columns[i]] {
			continue
		}
		r = append(r, Columns[i])
	}
	return r
}
