// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/utils/json"
)

// ServiceRevision is the model entity for the ServiceRevision schema.
type ServiceRevision struct {
	config `json:"-"`
	// ID of the ent.
	ID oid.ID `json:"id,omitempty" sql:"id"`
	// CreateTime holds the value of the "createTime" field.
	CreateTime *time.Time `json:"createTime,omitempty" sql:"createTime"`
	// ID of the project to belong.
	ProjectID oid.ID `json:"projectID,omitempty" sql:"projectID"`
	// Status holds the value of the "status" field.
	Status string `json:"status,omitempty" sql:"status"`
	// StatusMessage holds the value of the "statusMessage" field.
	StatusMessage string `json:"statusMessage,omitempty" sql:"statusMessage"`
	// Type of the revision.
	Type string `json:"type,omitempty" sql:"type"`
	// ID of the service to which the revision belongs.
	ServiceID oid.ID `json:"serviceID,omitempty" sql:"serviceID"`
	// ID of the environment to which the service deploys.
	EnvironmentID oid.ID `json:"environmentID,omitempty" sql:"environmentID"`
	// ID of the template.
	TemplateID string `json:"templateID,omitempty" sql:"templateID"`
	// Version of the template.
	TemplateVersion string `json:"templateVersion,omitempty" sql:"templateVersion"`
	// Attributes to configure the template.
	Attributes property.Values `json:"attributes,omitempty" sql:"attributes"`
	// Variables of the revision.
	Variables crypto.Map[string, string] `json:"variables,omitempty" sql:"variables"`
	// Input plan of the revision.
	InputPlan string `json:"-" sql:"inputPlan"`
	// Output of the revision.
	Output string `json:"-" sql:"output"`
	// Type of deployer.
	DeployerType string `json:"deployerType,omitempty" sql:"deployerType"`
	// Duration in seconds of the revision deploying.
	Duration int `json:"duration,omitempty" sql:"duration"`
	// Previous provider requirement of the revision.
	PreviousRequiredProviders []types.ProviderRequirement `json:"previousRequiredProviders,omitempty" sql:"previousRequiredProviders"`
	// Tags of the revision.
	Tags []string `json:"tags,omitempty" sql:"tags"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ServiceRevisionQuery when eager-loading is set.
	Edges        ServiceRevisionEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ServiceRevisionEdges holds the relations/edges for other nodes in the graph.
type ServiceRevisionEdges struct {
	// Project to which the revision belongs.
	Project *Project `json:"project,omitempty" sql:"project"`
	// Environment to which the revision deploys.
	Environment *Environment `json:"environment,omitempty" sql:"environment"`
	// Service to which the revision belongs.
	Service *Service `json:"service,omitempty" sql:"service"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceRevisionEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[0] {
		if e.Project == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// EnvironmentOrErr returns the Environment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceRevisionEdges) EnvironmentOrErr() (*Environment, error) {
	if e.loadedTypes[1] {
		if e.Environment == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: environment.Label}
		}
		return e.Environment, nil
	}
	return nil, &NotLoadedError{edge: "environment"}
}

// ServiceOrErr returns the Service value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceRevisionEdges) ServiceOrErr() (*Service, error) {
	if e.loadedTypes[2] {
		if e.Service == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: service.Label}
		}
		return e.Service, nil
	}
	return nil, &NotLoadedError{edge: "service"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ServiceRevision) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case servicerevision.FieldPreviousRequiredProviders, servicerevision.FieldTags:
			values[i] = new([]byte)
		case servicerevision.FieldVariables:
			values[i] = new(crypto.Map[string, string])
		case servicerevision.FieldID, servicerevision.FieldProjectID, servicerevision.FieldServiceID, servicerevision.FieldEnvironmentID:
			values[i] = new(oid.ID)
		case servicerevision.FieldAttributes:
			values[i] = new(property.Values)
		case servicerevision.FieldDuration:
			values[i] = new(sql.NullInt64)
		case servicerevision.FieldStatus, servicerevision.FieldStatusMessage, servicerevision.FieldType, servicerevision.FieldTemplateID, servicerevision.FieldTemplateVersion, servicerevision.FieldInputPlan, servicerevision.FieldOutput, servicerevision.FieldDeployerType:
			values[i] = new(sql.NullString)
		case servicerevision.FieldCreateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ServiceRevision fields.
func (sr *ServiceRevision) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case servicerevision.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				sr.ID = *value
			}
		case servicerevision.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				sr.CreateTime = new(time.Time)
				*sr.CreateTime = value.Time
			}
		case servicerevision.FieldProjectID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field projectID", values[i])
			} else if value != nil {
				sr.ProjectID = *value
			}
		case servicerevision.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				sr.Status = value.String
			}
		case servicerevision.FieldStatusMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field statusMessage", values[i])
			} else if value.Valid {
				sr.StatusMessage = value.String
			}
		case servicerevision.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				sr.Type = value.String
			}
		case servicerevision.FieldServiceID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field serviceID", values[i])
			} else if value != nil {
				sr.ServiceID = *value
			}
		case servicerevision.FieldEnvironmentID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field environmentID", values[i])
			} else if value != nil {
				sr.EnvironmentID = *value
			}
		case servicerevision.FieldTemplateID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field templateID", values[i])
			} else if value.Valid {
				sr.TemplateID = value.String
			}
		case servicerevision.FieldTemplateVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field templateVersion", values[i])
			} else if value.Valid {
				sr.TemplateVersion = value.String
			}
		case servicerevision.FieldAttributes:
			if value, ok := values[i].(*property.Values); !ok {
				return fmt.Errorf("unexpected type %T for field attributes", values[i])
			} else if value != nil {
				sr.Attributes = *value
			}
		case servicerevision.FieldVariables:
			if value, ok := values[i].(*crypto.Map[string, string]); !ok {
				return fmt.Errorf("unexpected type %T for field variables", values[i])
			} else if value != nil {
				sr.Variables = *value
			}
		case servicerevision.FieldInputPlan:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field inputPlan", values[i])
			} else if value.Valid {
				sr.InputPlan = value.String
			}
		case servicerevision.FieldOutput:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field output", values[i])
			} else if value.Valid {
				sr.Output = value.String
			}
		case servicerevision.FieldDeployerType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field deployerType", values[i])
			} else if value.Valid {
				sr.DeployerType = value.String
			}
		case servicerevision.FieldDuration:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field duration", values[i])
			} else if value.Valid {
				sr.Duration = int(value.Int64)
			}
		case servicerevision.FieldPreviousRequiredProviders:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field previousRequiredProviders", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &sr.PreviousRequiredProviders); err != nil {
					return fmt.Errorf("unmarshal field previousRequiredProviders: %w", err)
				}
			}
		case servicerevision.FieldTags:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field tags", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &sr.Tags); err != nil {
					return fmt.Errorf("unmarshal field tags: %w", err)
				}
			}
		default:
			sr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ServiceRevision.
// This includes values selected through modifiers, order, etc.
func (sr *ServiceRevision) Value(name string) (ent.Value, error) {
	return sr.selectValues.Get(name)
}

// QueryProject queries the "project" edge of the ServiceRevision entity.
func (sr *ServiceRevision) QueryProject() *ProjectQuery {
	return NewServiceRevisionClient(sr.config).QueryProject(sr)
}

// QueryEnvironment queries the "environment" edge of the ServiceRevision entity.
func (sr *ServiceRevision) QueryEnvironment() *EnvironmentQuery {
	return NewServiceRevisionClient(sr.config).QueryEnvironment(sr)
}

// QueryService queries the "service" edge of the ServiceRevision entity.
func (sr *ServiceRevision) QueryService() *ServiceQuery {
	return NewServiceRevisionClient(sr.config).QueryService(sr)
}

// Update returns a builder for updating this ServiceRevision.
// Note that you need to call ServiceRevision.Unwrap() before calling this method if this ServiceRevision
// was returned from a transaction, and the transaction was committed or rolled back.
func (sr *ServiceRevision) Update() *ServiceRevisionUpdateOne {
	return NewServiceRevisionClient(sr.config).UpdateOne(sr)
}

// Unwrap unwraps the ServiceRevision entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (sr *ServiceRevision) Unwrap() *ServiceRevision {
	_tx, ok := sr.config.driver.(*txDriver)
	if !ok {
		panic("model: ServiceRevision is not a transactional entity")
	}
	sr.config.driver = _tx.drv
	return sr
}

// String implements the fmt.Stringer.
func (sr *ServiceRevision) String() string {
	var builder strings.Builder
	builder.WriteString("ServiceRevision(")
	builder.WriteString(fmt.Sprintf("id=%v, ", sr.ID))
	if v := sr.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("projectID=")
	builder.WriteString(fmt.Sprintf("%v", sr.ProjectID))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(sr.Status)
	builder.WriteString(", ")
	builder.WriteString("statusMessage=")
	builder.WriteString(sr.StatusMessage)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(sr.Type)
	builder.WriteString(", ")
	builder.WriteString("serviceID=")
	builder.WriteString(fmt.Sprintf("%v", sr.ServiceID))
	builder.WriteString(", ")
	builder.WriteString("environmentID=")
	builder.WriteString(fmt.Sprintf("%v", sr.EnvironmentID))
	builder.WriteString(", ")
	builder.WriteString("templateID=")
	builder.WriteString(sr.TemplateID)
	builder.WriteString(", ")
	builder.WriteString("templateVersion=")
	builder.WriteString(sr.TemplateVersion)
	builder.WriteString(", ")
	builder.WriteString("attributes=")
	builder.WriteString(fmt.Sprintf("%v", sr.Attributes))
	builder.WriteString(", ")
	builder.WriteString("variables=")
	builder.WriteString(fmt.Sprintf("%v", sr.Variables))
	builder.WriteString(", ")
	builder.WriteString("inputPlan=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("output=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("deployerType=")
	builder.WriteString(sr.DeployerType)
	builder.WriteString(", ")
	builder.WriteString("duration=")
	builder.WriteString(fmt.Sprintf("%v", sr.Duration))
	builder.WriteString(", ")
	builder.WriteString("previousRequiredProviders=")
	builder.WriteString(fmt.Sprintf("%v", sr.PreviousRequiredProviders))
	builder.WriteString(", ")
	builder.WriteString("tags=")
	builder.WriteString(fmt.Sprintf("%v", sr.Tags))
	builder.WriteByte(')')
	return builder.String()
}

// ServiceRevisions is a parsable slice of ServiceRevision.
type ServiceRevisions []*ServiceRevision
