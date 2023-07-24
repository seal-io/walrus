// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package serviceresourcerelationship

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// ID filters vertices based on their ID field.
func ID(id object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldCreateTime, v))
}

// ServiceResourceID applies equality check predicate on the "service_resource_id" field. It's identical to ServiceResourceIDEQ.
func ServiceResourceID(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldServiceResourceID, v))
}

// DependencyID applies equality check predicate on the "dependency_id" field. It's identical to DependencyIDEQ.
func DependencyID(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldDependencyID, v))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldType, v))
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLTE(FieldCreateTime, v))
}

// ServiceResourceIDEQ applies the EQ predicate on the "service_resource_id" field.
func ServiceResourceIDEQ(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldServiceResourceID, v))
}

// ServiceResourceIDNEQ applies the NEQ predicate on the "service_resource_id" field.
func ServiceResourceIDNEQ(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNEQ(FieldServiceResourceID, v))
}

// ServiceResourceIDIn applies the In predicate on the "service_resource_id" field.
func ServiceResourceIDIn(vs ...object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldIn(FieldServiceResourceID, vs...))
}

// ServiceResourceIDNotIn applies the NotIn predicate on the "service_resource_id" field.
func ServiceResourceIDNotIn(vs ...object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNotIn(FieldServiceResourceID, vs...))
}

// ServiceResourceIDGT applies the GT predicate on the "service_resource_id" field.
func ServiceResourceIDGT(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGT(FieldServiceResourceID, v))
}

// ServiceResourceIDGTE applies the GTE predicate on the "service_resource_id" field.
func ServiceResourceIDGTE(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGTE(FieldServiceResourceID, v))
}

// ServiceResourceIDLT applies the LT predicate on the "service_resource_id" field.
func ServiceResourceIDLT(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLT(FieldServiceResourceID, v))
}

// ServiceResourceIDLTE applies the LTE predicate on the "service_resource_id" field.
func ServiceResourceIDLTE(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLTE(FieldServiceResourceID, v))
}

// ServiceResourceIDContains applies the Contains predicate on the "service_resource_id" field.
func ServiceResourceIDContains(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldContains(FieldServiceResourceID, vc))
}

// ServiceResourceIDHasPrefix applies the HasPrefix predicate on the "service_resource_id" field.
func ServiceResourceIDHasPrefix(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldHasPrefix(FieldServiceResourceID, vc))
}

// ServiceResourceIDHasSuffix applies the HasSuffix predicate on the "service_resource_id" field.
func ServiceResourceIDHasSuffix(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldHasSuffix(FieldServiceResourceID, vc))
}

// ServiceResourceIDEqualFold applies the EqualFold predicate on the "service_resource_id" field.
func ServiceResourceIDEqualFold(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldEqualFold(FieldServiceResourceID, vc))
}

// ServiceResourceIDContainsFold applies the ContainsFold predicate on the "service_resource_id" field.
func ServiceResourceIDContainsFold(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldContainsFold(FieldServiceResourceID, vc))
}

// DependencyIDEQ applies the EQ predicate on the "dependency_id" field.
func DependencyIDEQ(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldDependencyID, v))
}

// DependencyIDNEQ applies the NEQ predicate on the "dependency_id" field.
func DependencyIDNEQ(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNEQ(FieldDependencyID, v))
}

// DependencyIDIn applies the In predicate on the "dependency_id" field.
func DependencyIDIn(vs ...object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldIn(FieldDependencyID, vs...))
}

// DependencyIDNotIn applies the NotIn predicate on the "dependency_id" field.
func DependencyIDNotIn(vs ...object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNotIn(FieldDependencyID, vs...))
}

// DependencyIDGT applies the GT predicate on the "dependency_id" field.
func DependencyIDGT(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGT(FieldDependencyID, v))
}

// DependencyIDGTE applies the GTE predicate on the "dependency_id" field.
func DependencyIDGTE(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGTE(FieldDependencyID, v))
}

// DependencyIDLT applies the LT predicate on the "dependency_id" field.
func DependencyIDLT(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLT(FieldDependencyID, v))
}

// DependencyIDLTE applies the LTE predicate on the "dependency_id" field.
func DependencyIDLTE(v object.ID) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLTE(FieldDependencyID, v))
}

// DependencyIDContains applies the Contains predicate on the "dependency_id" field.
func DependencyIDContains(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldContains(FieldDependencyID, vc))
}

// DependencyIDHasPrefix applies the HasPrefix predicate on the "dependency_id" field.
func DependencyIDHasPrefix(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldHasPrefix(FieldDependencyID, vc))
}

// DependencyIDHasSuffix applies the HasSuffix predicate on the "dependency_id" field.
func DependencyIDHasSuffix(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldHasSuffix(FieldDependencyID, vc))
}

// DependencyIDEqualFold applies the EqualFold predicate on the "dependency_id" field.
func DependencyIDEqualFold(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldEqualFold(FieldDependencyID, vc))
}

// DependencyIDContainsFold applies the ContainsFold predicate on the "dependency_id" field.
func DependencyIDContainsFold(v object.ID) predicate.ServiceResourceRelationship {
	vc := string(v)
	return predicate.ServiceResourceRelationship(sql.FieldContainsFold(FieldDependencyID, vc))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(sql.FieldContainsFold(FieldType, v))
}

// HasServiceResource applies the HasEdge predicate on the "serviceResource" edge.
func HasServiceResource() predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, ServiceResourceTable, ServiceResourceColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ServiceResource
		step.Edge.Schema = schemaConfig.ServiceResourceRelationship
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasServiceResourceWith applies the HasEdge predicate on the "serviceResource" edge with a given conditions (other predicates).
func HasServiceResourceWith(preds ...predicate.ServiceResource) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(func(s *sql.Selector) {
		step := newServiceResourceStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ServiceResource
		step.Edge.Schema = schemaConfig.ServiceResourceRelationship
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDependency applies the HasEdge predicate on the "dependency" edge.
func HasDependency() predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, DependencyTable, DependencyColumn),
		)
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ServiceResource
		step.Edge.Schema = schemaConfig.ServiceResourceRelationship
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDependencyWith applies the HasEdge predicate on the "dependency" edge with a given conditions (other predicates).
func HasDependencyWith(preds ...predicate.ServiceResource) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(func(s *sql.Selector) {
		step := newDependencyStep()
		schemaConfig := internal.SchemaConfigFromContext(s.Context())
		step.To.Schema = schemaConfig.ServiceResource
		step.Edge.Schema = schemaConfig.ServiceResourceRelationship
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ServiceResourceRelationship) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ServiceResourceRelationship) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(func(s *sql.Selector) {
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
func Not(p predicate.ServiceResourceRelationship) predicate.ServiceResourceRelationship {
	return predicate.ServiceResourceRelationship(func(s *sql.Selector) {
		p(s.Not())
	})
}
