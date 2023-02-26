// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/types"
)

// EnvironmentConnectorRelationship is the model entity for the EnvironmentConnectorRelationship schema.
type EnvironmentConnectorRelationship struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// ID of the environment to which the relationship connects.
	EnvironmentID types.ID `json:"environmentID"`
	// ID of the connector to which the relationship connects.
	ConnectorID types.ID `json:"connectorID"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EnvironmentConnectorRelationshipQuery when eager-loading is set.
	Edges EnvironmentConnectorRelationshipEdges `json:"edges,omitempty"`
}

// EnvironmentConnectorRelationshipEdges holds the relations/edges for other nodes in the graph.
type EnvironmentConnectorRelationshipEdges struct {
	// Environments that connect to the relationship.
	Environment *Environment `json:"environment,omitempty"`
	// Connectors that connect to the relationship.
	Connector *Connector `json:"connector,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// EnvironmentOrErr returns the Environment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EnvironmentConnectorRelationshipEdges) EnvironmentOrErr() (*Environment, error) {
	if e.loadedTypes[0] {
		if e.Environment == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: environment.Label}
		}
		return e.Environment, nil
	}
	return nil, &NotLoadedError{edge: "environment"}
}

// ConnectorOrErr returns the Connector value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EnvironmentConnectorRelationshipEdges) ConnectorOrErr() (*Connector, error) {
	if e.loadedTypes[1] {
		if e.Connector == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: connector.Label}
		}
		return e.Connector, nil
	}
	return nil, &NotLoadedError{edge: "connector"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*EnvironmentConnectorRelationship) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case environmentconnectorrelationship.FieldID:
			values[i] = new(sql.NullInt64)
		case environmentconnectorrelationship.FieldCreateTime:
			values[i] = new(sql.NullTime)
		case environmentconnectorrelationship.FieldEnvironmentID, environmentconnectorrelationship.FieldConnectorID:
			values[i] = new(types.ID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type EnvironmentConnectorRelationship", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the EnvironmentConnectorRelationship fields.
func (ecr *EnvironmentConnectorRelationship) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case environmentconnectorrelationship.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ecr.ID = int(value.Int64)
		case environmentconnectorrelationship.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				ecr.CreateTime = new(time.Time)
				*ecr.CreateTime = value.Time
			}
		case environmentconnectorrelationship.FieldEnvironmentID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field environment_id", values[i])
			} else if value != nil {
				ecr.EnvironmentID = *value
			}
		case environmentconnectorrelationship.FieldConnectorID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field connector_id", values[i])
			} else if value != nil {
				ecr.ConnectorID = *value
			}
		}
	}
	return nil
}

// QueryEnvironment queries the "environment" edge of the EnvironmentConnectorRelationship entity.
func (ecr *EnvironmentConnectorRelationship) QueryEnvironment() *EnvironmentQuery {
	return NewEnvironmentConnectorRelationshipClient(ecr.config).QueryEnvironment(ecr)
}

// QueryConnector queries the "connector" edge of the EnvironmentConnectorRelationship entity.
func (ecr *EnvironmentConnectorRelationship) QueryConnector() *ConnectorQuery {
	return NewEnvironmentConnectorRelationshipClient(ecr.config).QueryConnector(ecr)
}

// Update returns a builder for updating this EnvironmentConnectorRelationship.
// Note that you need to call EnvironmentConnectorRelationship.Unwrap() before calling this method if this EnvironmentConnectorRelationship
// was returned from a transaction, and the transaction was committed or rolled back.
func (ecr *EnvironmentConnectorRelationship) Update() *EnvironmentConnectorRelationshipUpdateOne {
	return NewEnvironmentConnectorRelationshipClient(ecr.config).UpdateOne(ecr)
}

// Unwrap unwraps the EnvironmentConnectorRelationship entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ecr *EnvironmentConnectorRelationship) Unwrap() *EnvironmentConnectorRelationship {
	_tx, ok := ecr.config.driver.(*txDriver)
	if !ok {
		panic("model: EnvironmentConnectorRelationship is not a transactional entity")
	}
	ecr.config.driver = _tx.drv
	return ecr
}

// String implements the fmt.Stringer.
func (ecr *EnvironmentConnectorRelationship) String() string {
	var builder strings.Builder
	builder.WriteString("EnvironmentConnectorRelationship(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ecr.ID))
	if v := ecr.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("environment_id=")
	builder.WriteString(fmt.Sprintf("%v", ecr.EnvironmentID))
	builder.WriteString(", ")
	builder.WriteString("connector_id=")
	builder.WriteString(fmt.Sprintf("%v", ecr.ConnectorID))
	builder.WriteByte(')')
	return builder.String()
}

// EnvironmentConnectorRelationships is a parsable slice of EnvironmentConnectorRelationship.
type EnvironmentConnectorRelationships []*EnvironmentConnectorRelationship

func (ecr EnvironmentConnectorRelationships) config(cfg config) {
	for _i := range ecr {
		ecr[_i].config = cfg
	}
}
