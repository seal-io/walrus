// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model/serviceresource"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// ServiceResourceRelationship is the model entity for the ServiceResourceRelationship schema.
type ServiceResourceRelationship struct {
	config `json:"-"`
	// ID of the ent.
	ID object.ID `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime *time.Time `json:"create_time,omitempty"`
	// ID of the service resource.
	ServiceResourceID object.ID `json:"service_resource_id,omitempty"`
	// ID of the resource that resource depends on.
	DependencyID object.ID `json:"dependency_id,omitempty"`
	// Type of the relationship.
	Type string `json:"type,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ServiceResourceRelationshipQuery when eager-loading is set.
	Edges        ServiceResourceRelationshipEdges `json:"edges,omitempty"`
	selectValues sql.SelectValues
}

// ServiceResourceRelationshipEdges holds the relations/edges for other nodes in the graph.
type ServiceResourceRelationshipEdges struct {
	// ServiceResource to which it currently belongs.
	ServiceResource *ServiceResource `json:"serviceResource,omitempty"`
	// ServiceResource to which the dependency belongs.
	Dependency *ServiceResource `json:"dependency,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ServiceResourceOrErr returns the ServiceResource value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceResourceRelationshipEdges) ServiceResourceOrErr() (*ServiceResource, error) {
	if e.loadedTypes[0] {
		if e.ServiceResource == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: serviceresource.Label}
		}
		return e.ServiceResource, nil
	}
	return nil, &NotLoadedError{edge: "serviceResource"}
}

// DependencyOrErr returns the Dependency value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceResourceRelationshipEdges) DependencyOrErr() (*ServiceResource, error) {
	if e.loadedTypes[1] {
		if e.Dependency == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: serviceresource.Label}
		}
		return e.Dependency, nil
	}
	return nil, &NotLoadedError{edge: "dependency"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ServiceResourceRelationship) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case serviceresourcerelationship.FieldID, serviceresourcerelationship.FieldServiceResourceID, serviceresourcerelationship.FieldDependencyID:
			values[i] = new(object.ID)
		case serviceresourcerelationship.FieldType:
			values[i] = new(sql.NullString)
		case serviceresourcerelationship.FieldCreateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ServiceResourceRelationship fields.
func (srr *ServiceResourceRelationship) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case serviceresourcerelationship.FieldID:
			if value, ok := values[i].(*object.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				srr.ID = *value
			}
		case serviceresourcerelationship.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				srr.CreateTime = new(time.Time)
				*srr.CreateTime = value.Time
			}
		case serviceresourcerelationship.FieldServiceResourceID:
			if value, ok := values[i].(*object.ID); !ok {
				return fmt.Errorf("unexpected type %T for field service_resource_id", values[i])
			} else if value != nil {
				srr.ServiceResourceID = *value
			}
		case serviceresourcerelationship.FieldDependencyID:
			if value, ok := values[i].(*object.ID); !ok {
				return fmt.Errorf("unexpected type %T for field dependency_id", values[i])
			} else if value != nil {
				srr.DependencyID = *value
			}
		case serviceresourcerelationship.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				srr.Type = value.String
			}
		default:
			srr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ServiceResourceRelationship.
// This includes values selected through modifiers, order, etc.
func (srr *ServiceResourceRelationship) Value(name string) (ent.Value, error) {
	return srr.selectValues.Get(name)
}

// QueryServiceResource queries the "serviceResource" edge of the ServiceResourceRelationship entity.
func (srr *ServiceResourceRelationship) QueryServiceResource() *ServiceResourceQuery {
	return NewServiceResourceRelationshipClient(srr.config).QueryServiceResource(srr)
}

// QueryDependency queries the "dependency" edge of the ServiceResourceRelationship entity.
func (srr *ServiceResourceRelationship) QueryDependency() *ServiceResourceQuery {
	return NewServiceResourceRelationshipClient(srr.config).QueryDependency(srr)
}

// Update returns a builder for updating this ServiceResourceRelationship.
// Note that you need to call ServiceResourceRelationship.Unwrap() before calling this method if this ServiceResourceRelationship
// was returned from a transaction, and the transaction was committed or rolled back.
func (srr *ServiceResourceRelationship) Update() *ServiceResourceRelationshipUpdateOne {
	return NewServiceResourceRelationshipClient(srr.config).UpdateOne(srr)
}

// Unwrap unwraps the ServiceResourceRelationship entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (srr *ServiceResourceRelationship) Unwrap() *ServiceResourceRelationship {
	_tx, ok := srr.config.driver.(*txDriver)
	if !ok {
		panic("model: ServiceResourceRelationship is not a transactional entity")
	}
	srr.config.driver = _tx.drv
	return srr
}

// String implements the fmt.Stringer.
func (srr *ServiceResourceRelationship) String() string {
	var builder strings.Builder
	builder.WriteString("ServiceResourceRelationship(")
	builder.WriteString(fmt.Sprintf("id=%v, ", srr.ID))
	if v := srr.CreateTime; v != nil {
		builder.WriteString("create_time=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("service_resource_id=")
	builder.WriteString(fmt.Sprintf("%v", srr.ServiceResourceID))
	builder.WriteString(", ")
	builder.WriteString("dependency_id=")
	builder.WriteString(fmt.Sprintf("%v", srr.DependencyID))
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(srr.Type)
	builder.WriteByte(')')
	return builder.String()
}

// ServiceResourceRelationships is a parsable slice of ServiceResourceRelationship.
type ServiceResourceRelationships []*ServiceResourceRelationship
