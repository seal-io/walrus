// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package environmentconnectorrelationship

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldEQ(FieldCreateTime, v))
}

// EnvironmentID applies equality check predicate on the "environment_id" field. It's identical to EnvironmentIDEQ.
func EnvironmentID(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldEQ(FieldEnvironmentID, v))
}

// ConnectorID applies equality check predicate on the "connector_id" field. It's identical to ConnectorIDEQ.
func ConnectorID(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldEQ(FieldConnectorID, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldLTE(FieldCreateTime, v))
}

// EnvironmentIDEQ applies the EQ predicate on the "environment_id" field.
func EnvironmentIDEQ(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldEQ(FieldEnvironmentID, v))
}

// EnvironmentIDNEQ applies the NEQ predicate on the "environment_id" field.
func EnvironmentIDNEQ(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldNEQ(FieldEnvironmentID, v))
}

// EnvironmentIDIn applies the In predicate on the "environment_id" field.
func EnvironmentIDIn(vs ...types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDNotIn applies the NotIn predicate on the "environment_id" field.
func EnvironmentIDNotIn(vs ...types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldNotIn(FieldEnvironmentID, vs...))
}

// EnvironmentIDGT applies the GT predicate on the "environment_id" field.
func EnvironmentIDGT(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldGT(FieldEnvironmentID, v))
}

// EnvironmentIDGTE applies the GTE predicate on the "environment_id" field.
func EnvironmentIDGTE(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldGTE(FieldEnvironmentID, v))
}

// EnvironmentIDLT applies the LT predicate on the "environment_id" field.
func EnvironmentIDLT(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldLT(FieldEnvironmentID, v))
}

// EnvironmentIDLTE applies the LTE predicate on the "environment_id" field.
func EnvironmentIDLTE(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldLTE(FieldEnvironmentID, v))
}

// EnvironmentIDContains applies the Contains predicate on the "environment_id" field.
func EnvironmentIDContains(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldContains(FieldEnvironmentID, vc))
}

// EnvironmentIDHasPrefix applies the HasPrefix predicate on the "environment_id" field.
func EnvironmentIDHasPrefix(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldHasPrefix(FieldEnvironmentID, vc))
}

// EnvironmentIDHasSuffix applies the HasSuffix predicate on the "environment_id" field.
func EnvironmentIDHasSuffix(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldHasSuffix(FieldEnvironmentID, vc))
}

// EnvironmentIDEqualFold applies the EqualFold predicate on the "environment_id" field.
func EnvironmentIDEqualFold(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldEqualFold(FieldEnvironmentID, vc))
}

// EnvironmentIDContainsFold applies the ContainsFold predicate on the "environment_id" field.
func EnvironmentIDContainsFold(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldContainsFold(FieldEnvironmentID, vc))
}

// ConnectorIDEQ applies the EQ predicate on the "connector_id" field.
func ConnectorIDEQ(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldEQ(FieldConnectorID, v))
}

// ConnectorIDNEQ applies the NEQ predicate on the "connector_id" field.
func ConnectorIDNEQ(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldNEQ(FieldConnectorID, v))
}

// ConnectorIDIn applies the In predicate on the "connector_id" field.
func ConnectorIDIn(vs ...types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldIn(FieldConnectorID, vs...))
}

// ConnectorIDNotIn applies the NotIn predicate on the "connector_id" field.
func ConnectorIDNotIn(vs ...types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldNotIn(FieldConnectorID, vs...))
}

// ConnectorIDGT applies the GT predicate on the "connector_id" field.
func ConnectorIDGT(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldGT(FieldConnectorID, v))
}

// ConnectorIDGTE applies the GTE predicate on the "connector_id" field.
func ConnectorIDGTE(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldGTE(FieldConnectorID, v))
}

// ConnectorIDLT applies the LT predicate on the "connector_id" field.
func ConnectorIDLT(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldLT(FieldConnectorID, v))
}

// ConnectorIDLTE applies the LTE predicate on the "connector_id" field.
func ConnectorIDLTE(v types.ID) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(sql.FieldLTE(FieldConnectorID, v))
}

// ConnectorIDContains applies the Contains predicate on the "connector_id" field.
func ConnectorIDContains(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldContains(FieldConnectorID, vc))
}

// ConnectorIDHasPrefix applies the HasPrefix predicate on the "connector_id" field.
func ConnectorIDHasPrefix(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldHasPrefix(FieldConnectorID, vc))
}

// ConnectorIDHasSuffix applies the HasSuffix predicate on the "connector_id" field.
func ConnectorIDHasSuffix(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldHasSuffix(FieldConnectorID, vc))
}

// ConnectorIDEqualFold applies the EqualFold predicate on the "connector_id" field.
func ConnectorIDEqualFold(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldEqualFold(FieldConnectorID, vc))
}

// ConnectorIDContainsFold applies the ContainsFold predicate on the "connector_id" field.
func ConnectorIDContainsFold(v types.ID) predicate.EnvironmentConnectorRelationship {
	vc := string(v)
	return predicate.EnvironmentConnectorRelationship(sql.FieldContainsFold(FieldConnectorID, vc))
}

// HasEnvironment applies the HasEdge predicate on the "environment" edge.
func HasEnvironment() predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, EnvironmentColumn),
			sqlgraph.Edge(sqlgraph.M2O, false, EnvironmentTable, EnvironmentColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasEnvironmentWith applies the HasEdge predicate on the "environment" edge with a given conditions (other predicates).
func HasEnvironmentWith(preds ...predicate.Environment) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, EnvironmentColumn),
			sqlgraph.To(EnvironmentInverseTable, EnvironmentFieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, EnvironmentTable, EnvironmentColumn),
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

// HasConnector applies the HasEdge predicate on the "connector" edge.
func HasConnector() predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, ConnectorColumn),
			sqlgraph.Edge(sqlgraph.M2O, false, ConnectorTable, ConnectorColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasConnectorWith applies the HasEdge predicate on the "connector" edge with a given conditions (other predicates).
func HasConnectorWith(preds ...predicate.Connector) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, ConnectorColumn),
			sqlgraph.To(ConnectorInverseTable, ConnectorFieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ConnectorTable, ConnectorColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.EnvironmentConnectorRelationship) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.EnvironmentConnectorRelationship) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(func(s *sql.Selector) {
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
func Not(p predicate.EnvironmentConnectorRelationship) predicate.EnvironmentConnectorRelationship {
	return predicate.EnvironmentConnectorRelationship(func(s *sql.Selector) {
		p(s.Not())
	})
}
