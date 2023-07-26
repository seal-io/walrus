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

	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/json"
)

// ServiceResource is the model entity for the ServiceResource schema.
type ServiceResource struct {
	config `json:"-"`
	// ID of the ent.
	ID oid.ID `json:"id,omitempty" sql:"id"`
	// CreateTime holds the value of the "createTime" field.
	CreateTime *time.Time `json:"createTime,omitempty" sql:"createTime"`
	// UpdateTime holds the value of the "updateTime" field.
	UpdateTime *time.Time `json:"updateTime,omitempty" sql:"updateTime"`
	// ID of the project to belong.
	ProjectID oid.ID `json:"projectID,omitempty" sql:"projectID"`
	// ID of the service to which the resource belongs.
	ServiceID oid.ID `json:"serviceID,omitempty" sql:"serviceID"`
	// ID of the connector to which the resource deploys.
	ConnectorID oid.ID `json:"connectorID,omitempty" sql:"connectorID"`
	// ID of the parent resource.
	CompositionID oid.ID `json:"compositionID,omitempty" sql:"compositionID"`
	// ID of the parent class of the resource realization.
	ClassID oid.ID `json:"classID,omitempty" sql:"classID"`
	// Mode that manages the generated resource, it is the management way of the deployer to the resource, which provides by deployer.
	Mode string `json:"mode,omitempty" sql:"mode"`
	// Type of the generated resource, it is the type of the resource which the deployer observes, which provides by deployer.
	Type string `json:"type,omitempty" sql:"type"`
	// Name of the generated resource, it is the real identifier of the resource, which provides by deployer.
	Name string `json:"name,omitempty" sql:"name"`
	// Type of deployer.
	DeployerType string `json:"deployerType,omitempty" sql:"deployerType"`
	// Shape of the resource, it can be class or instance shape.
	Shape string `json:"shape,omitempty" sql:"shape"`
	// Status of the resource.
	Status types.ServiceResourceStatus `json:"status,omitempty" sql:"status"`
	// Drift detection result.
	DriftResult *types.ServiceResourceDriftResult `json:"driftResult,omitempty" sql:"driftResult"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ServiceResourceQuery when eager-loading is set.
	Edges        ServiceResourceEdges `json:"edges"`
	selectValues sql.SelectValues

	// Keys is the list of key used for operating the service resource,
	// it does not store in the database and only records for additional operations.
	Keys *types.ServiceResourceOperationKeys `json:"keys,omitempty"`
}

// ServiceResourceEdges holds the relations/edges for other nodes in the graph.
type ServiceResourceEdges struct {
	// Service to which the resource belongs.
	Service *Service `json:"service,omitempty" sql:"service"`
	// Connector to which the resource deploys.
	Connector *Connector `json:"connector,omitempty" sql:"connector"`
	// Composition holds the value of the composition edge.
	Composition *ServiceResource `json:"composition,omitempty" sql:"composition"`
	// Sub-resources that make up the resource.
	Components []*ServiceResource `json:"components,omitempty" sql:"components"`
	// Class holds the value of the class edge.
	Class *ServiceResource `json:"class,omitempty" sql:"class"`
	// Service resource instances to which the resource defines.
	Instances []*ServiceResource `json:"instances,omitempty" sql:"instances"`
	// Dependencies holds the value of the dependencies edge.
	Dependencies []*ServiceResourceRelationship `json:"dependencies,omitempty" sql:"dependencies"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [7]bool
}

// ServiceOrErr returns the Service value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceResourceEdges) ServiceOrErr() (*Service, error) {
	if e.loadedTypes[0] {
		if e.Service == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: service.Label}
		}
		return e.Service, nil
	}
	return nil, &NotLoadedError{edge: "service"}
}

// ConnectorOrErr returns the Connector value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceResourceEdges) ConnectorOrErr() (*Connector, error) {
	if e.loadedTypes[1] {
		if e.Connector == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: connector.Label}
		}
		return e.Connector, nil
	}
	return nil, &NotLoadedError{edge: "connector"}
}

// CompositionOrErr returns the Composition value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceResourceEdges) CompositionOrErr() (*ServiceResource, error) {
	if e.loadedTypes[2] {
		if e.Composition == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: serviceresource.Label}
		}
		return e.Composition, nil
	}
	return nil, &NotLoadedError{edge: "composition"}
}

// ComponentsOrErr returns the Components value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceResourceEdges) ComponentsOrErr() ([]*ServiceResource, error) {
	if e.loadedTypes[3] {
		return e.Components, nil
	}
	return nil, &NotLoadedError{edge: "components"}
}

// ClassOrErr returns the Class value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceResourceEdges) ClassOrErr() (*ServiceResource, error) {
	if e.loadedTypes[4] {
		if e.Class == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: serviceresource.Label}
		}
		return e.Class, nil
	}
	return nil, &NotLoadedError{edge: "class"}
}

// InstancesOrErr returns the Instances value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceResourceEdges) InstancesOrErr() ([]*ServiceResource, error) {
	if e.loadedTypes[5] {
		return e.Instances, nil
	}
	return nil, &NotLoadedError{edge: "instances"}
}

// DependenciesOrErr returns the Dependencies value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceResourceEdges) DependenciesOrErr() ([]*ServiceResourceRelationship, error) {
	if e.loadedTypes[6] {
		return e.Dependencies, nil
	}
	return nil, &NotLoadedError{edge: "dependencies"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ServiceResource) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case serviceresource.FieldStatus, serviceresource.FieldDriftResult:
			values[i] = new([]byte)
		case serviceresource.FieldID, serviceresource.FieldProjectID, serviceresource.FieldServiceID, serviceresource.FieldConnectorID, serviceresource.FieldCompositionID, serviceresource.FieldClassID:
			values[i] = new(oid.ID)
		case serviceresource.FieldMode, serviceresource.FieldType, serviceresource.FieldName, serviceresource.FieldDeployerType, serviceresource.FieldShape:
			values[i] = new(sql.NullString)
		case serviceresource.FieldCreateTime, serviceresource.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ServiceResource fields.
func (sr *ServiceResource) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case serviceresource.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				sr.ID = *value
			}
		case serviceresource.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				sr.CreateTime = new(time.Time)
				*sr.CreateTime = value.Time
			}
		case serviceresource.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				sr.UpdateTime = new(time.Time)
				*sr.UpdateTime = value.Time
			}
		case serviceresource.FieldProjectID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field projectID", values[i])
			} else if value != nil {
				sr.ProjectID = *value
			}
		case serviceresource.FieldServiceID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field serviceID", values[i])
			} else if value != nil {
				sr.ServiceID = *value
			}
		case serviceresource.FieldConnectorID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field connectorID", values[i])
			} else if value != nil {
				sr.ConnectorID = *value
			}
		case serviceresource.FieldCompositionID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field compositionID", values[i])
			} else if value != nil {
				sr.CompositionID = *value
			}
		case serviceresource.FieldClassID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field classID", values[i])
			} else if value != nil {
				sr.ClassID = *value
			}
		case serviceresource.FieldMode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field mode", values[i])
			} else if value.Valid {
				sr.Mode = value.String
			}
		case serviceresource.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				sr.Type = value.String
			}
		case serviceresource.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				sr.Name = value.String
			}
		case serviceresource.FieldDeployerType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field deployerType", values[i])
			} else if value.Valid {
				sr.DeployerType = value.String
			}
		case serviceresource.FieldShape:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field shape", values[i])
			} else if value.Valid {
				sr.Shape = value.String
			}
		case serviceresource.FieldStatus:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &sr.Status); err != nil {
					return fmt.Errorf("unmarshal field status: %w", err)
				}
			}
		case serviceresource.FieldDriftResult:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field driftResult", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &sr.DriftResult); err != nil {
					return fmt.Errorf("unmarshal field driftResult: %w", err)
				}
			}
		default:
			sr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ServiceResource.
// This includes values selected through modifiers, order, etc.
func (sr *ServiceResource) Value(name string) (ent.Value, error) {
	return sr.selectValues.Get(name)
}

// QueryService queries the "service" edge of the ServiceResource entity.
func (sr *ServiceResource) QueryService() *ServiceQuery {
	return NewServiceResourceClient(sr.config).QueryService(sr)
}

// QueryConnector queries the "connector" edge of the ServiceResource entity.
func (sr *ServiceResource) QueryConnector() *ConnectorQuery {
	return NewServiceResourceClient(sr.config).QueryConnector(sr)
}

// QueryComposition queries the "composition" edge of the ServiceResource entity.
func (sr *ServiceResource) QueryComposition() *ServiceResourceQuery {
	return NewServiceResourceClient(sr.config).QueryComposition(sr)
}

// QueryComponents queries the "components" edge of the ServiceResource entity.
func (sr *ServiceResource) QueryComponents() *ServiceResourceQuery {
	return NewServiceResourceClient(sr.config).QueryComponents(sr)
}

// QueryClass queries the "class" edge of the ServiceResource entity.
func (sr *ServiceResource) QueryClass() *ServiceResourceQuery {
	return NewServiceResourceClient(sr.config).QueryClass(sr)
}

// QueryInstances queries the "instances" edge of the ServiceResource entity.
func (sr *ServiceResource) QueryInstances() *ServiceResourceQuery {
	return NewServiceResourceClient(sr.config).QueryInstances(sr)
}

// QueryDependencies queries the "dependencies" edge of the ServiceResource entity.
func (sr *ServiceResource) QueryDependencies() *ServiceResourceRelationshipQuery {
	return NewServiceResourceClient(sr.config).QueryDependencies(sr)
}

// Update returns a builder for updating this ServiceResource.
// Note that you need to call ServiceResource.Unwrap() before calling this method if this ServiceResource
// was returned from a transaction, and the transaction was committed or rolled back.
func (sr *ServiceResource) Update() *ServiceResourceUpdateOne {
	return NewServiceResourceClient(sr.config).UpdateOne(sr)
}

// Unwrap unwraps the ServiceResource entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (sr *ServiceResource) Unwrap() *ServiceResource {
	_tx, ok := sr.config.driver.(*txDriver)
	if !ok {
		panic("model: ServiceResource is not a transactional entity")
	}
	sr.config.driver = _tx.drv
	return sr
}

// String implements the fmt.Stringer.
func (sr *ServiceResource) String() string {
	var builder strings.Builder
	builder.WriteString("ServiceResource(")
	builder.WriteString(fmt.Sprintf("id=%v, ", sr.ID))
	if v := sr.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := sr.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("projectID=")
	builder.WriteString(fmt.Sprintf("%v", sr.ProjectID))
	builder.WriteString(", ")
	builder.WriteString("serviceID=")
	builder.WriteString(fmt.Sprintf("%v", sr.ServiceID))
	builder.WriteString(", ")
	builder.WriteString("connectorID=")
	builder.WriteString(fmt.Sprintf("%v", sr.ConnectorID))
	builder.WriteString(", ")
	builder.WriteString("compositionID=")
	builder.WriteString(fmt.Sprintf("%v", sr.CompositionID))
	builder.WriteString(", ")
	builder.WriteString("classID=")
	builder.WriteString(fmt.Sprintf("%v", sr.ClassID))
	builder.WriteString(", ")
	builder.WriteString("mode=")
	builder.WriteString(sr.Mode)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(sr.Type)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(sr.Name)
	builder.WriteString(", ")
	builder.WriteString("deployerType=")
	builder.WriteString(sr.DeployerType)
	builder.WriteString(", ")
	builder.WriteString("shape=")
	builder.WriteString(sr.Shape)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", sr.Status))
	builder.WriteString(", ")
	builder.WriteString("driftResult=")
	builder.WriteString(fmt.Sprintf("%v", sr.DriftResult))
	builder.WriteByte(')')
	return builder.String()
}

// ServiceResources is a parsable slice of ServiceResource.
type ServiceResources []*ServiceResource
