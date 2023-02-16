// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package token

import (
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ID filters vertices based on their ID field.
func ID(id types.ID) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id types.ID) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id types.ID) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...types.ID) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...types.ID) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id types.ID) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id types.ID) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id types.ID) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id types.ID) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldID, id))
}

// CreateTime applies equality check predicate on the "createTime" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldCreateTime, v))
}

// UpdateTime applies equality check predicate on the "updateTime" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldUpdateTime, v))
}

// CasdoorTokenName applies equality check predicate on the "casdoorTokenName" field. It's identical to CasdoorTokenNameEQ.
func CasdoorTokenName(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldCasdoorTokenName, v))
}

// CasdoorTokenOwner applies equality check predicate on the "casdoorTokenOwner" field. It's identical to CasdoorTokenOwnerEQ.
func CasdoorTokenOwner(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldCasdoorTokenOwner, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldName, v))
}

// Expiration applies equality check predicate on the "expiration" field. It's identical to ExpirationEQ.
func Expiration(v int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldExpiration, v))
}

// CreateTimeEQ applies the EQ predicate on the "createTime" field.
func CreateTimeEQ(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldCreateTime, v))
}

// CreateTimeNEQ applies the NEQ predicate on the "createTime" field.
func CreateTimeNEQ(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldCreateTime, v))
}

// CreateTimeIn applies the In predicate on the "createTime" field.
func CreateTimeIn(vs ...time.Time) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldCreateTime, vs...))
}

// CreateTimeNotIn applies the NotIn predicate on the "createTime" field.
func CreateTimeNotIn(vs ...time.Time) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldCreateTime, vs...))
}

// CreateTimeGT applies the GT predicate on the "createTime" field.
func CreateTimeGT(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldCreateTime, v))
}

// CreateTimeGTE applies the GTE predicate on the "createTime" field.
func CreateTimeGTE(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldCreateTime, v))
}

// CreateTimeLT applies the LT predicate on the "createTime" field.
func CreateTimeLT(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldCreateTime, v))
}

// CreateTimeLTE applies the LTE predicate on the "createTime" field.
func CreateTimeLTE(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldCreateTime, v))
}

// UpdateTimeEQ applies the EQ predicate on the "updateTime" field.
func UpdateTimeEQ(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldUpdateTime, v))
}

// UpdateTimeNEQ applies the NEQ predicate on the "updateTime" field.
func UpdateTimeNEQ(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldUpdateTime, v))
}

// UpdateTimeIn applies the In predicate on the "updateTime" field.
func UpdateTimeIn(vs ...time.Time) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldUpdateTime, vs...))
}

// UpdateTimeNotIn applies the NotIn predicate on the "updateTime" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldUpdateTime, vs...))
}

// UpdateTimeGT applies the GT predicate on the "updateTime" field.
func UpdateTimeGT(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldUpdateTime, v))
}

// UpdateTimeGTE applies the GTE predicate on the "updateTime" field.
func UpdateTimeGTE(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldUpdateTime, v))
}

// UpdateTimeLT applies the LT predicate on the "updateTime" field.
func UpdateTimeLT(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldUpdateTime, v))
}

// UpdateTimeLTE applies the LTE predicate on the "updateTime" field.
func UpdateTimeLTE(v time.Time) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldUpdateTime, v))
}

// CasdoorTokenNameEQ applies the EQ predicate on the "casdoorTokenName" field.
func CasdoorTokenNameEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameNEQ applies the NEQ predicate on the "casdoorTokenName" field.
func CasdoorTokenNameNEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameIn applies the In predicate on the "casdoorTokenName" field.
func CasdoorTokenNameIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldCasdoorTokenName, vs...))
}

// CasdoorTokenNameNotIn applies the NotIn predicate on the "casdoorTokenName" field.
func CasdoorTokenNameNotIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldCasdoorTokenName, vs...))
}

// CasdoorTokenNameGT applies the GT predicate on the "casdoorTokenName" field.
func CasdoorTokenNameGT(v string) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameGTE applies the GTE predicate on the "casdoorTokenName" field.
func CasdoorTokenNameGTE(v string) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameLT applies the LT predicate on the "casdoorTokenName" field.
func CasdoorTokenNameLT(v string) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameLTE applies the LTE predicate on the "casdoorTokenName" field.
func CasdoorTokenNameLTE(v string) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameContains applies the Contains predicate on the "casdoorTokenName" field.
func CasdoorTokenNameContains(v string) predicate.Token {
	return predicate.Token(sql.FieldContains(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameHasPrefix applies the HasPrefix predicate on the "casdoorTokenName" field.
func CasdoorTokenNameHasPrefix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasPrefix(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameHasSuffix applies the HasSuffix predicate on the "casdoorTokenName" field.
func CasdoorTokenNameHasSuffix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasSuffix(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameEqualFold applies the EqualFold predicate on the "casdoorTokenName" field.
func CasdoorTokenNameEqualFold(v string) predicate.Token {
	return predicate.Token(sql.FieldEqualFold(FieldCasdoorTokenName, v))
}

// CasdoorTokenNameContainsFold applies the ContainsFold predicate on the "casdoorTokenName" field.
func CasdoorTokenNameContainsFold(v string) predicate.Token {
	return predicate.Token(sql.FieldContainsFold(FieldCasdoorTokenName, v))
}

// CasdoorTokenOwnerEQ applies the EQ predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerNEQ applies the NEQ predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerNEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerIn applies the In predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldCasdoorTokenOwner, vs...))
}

// CasdoorTokenOwnerNotIn applies the NotIn predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerNotIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldCasdoorTokenOwner, vs...))
}

// CasdoorTokenOwnerGT applies the GT predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerGT(v string) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerGTE applies the GTE predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerGTE(v string) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerLT applies the LT predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerLT(v string) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerLTE applies the LTE predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerLTE(v string) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerContains applies the Contains predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerContains(v string) predicate.Token {
	return predicate.Token(sql.FieldContains(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerHasPrefix applies the HasPrefix predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerHasPrefix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasPrefix(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerHasSuffix applies the HasSuffix predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerHasSuffix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasSuffix(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerEqualFold applies the EqualFold predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerEqualFold(v string) predicate.Token {
	return predicate.Token(sql.FieldEqualFold(FieldCasdoorTokenOwner, v))
}

// CasdoorTokenOwnerContainsFold applies the ContainsFold predicate on the "casdoorTokenOwner" field.
func CasdoorTokenOwnerContainsFold(v string) predicate.Token {
	return predicate.Token(sql.FieldContainsFold(FieldCasdoorTokenOwner, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Token {
	return predicate.Token(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Token {
	return predicate.Token(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Token {
	return predicate.Token(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Token {
	return predicate.Token(sql.FieldContainsFold(FieldName, v))
}

// ExpirationEQ applies the EQ predicate on the "expiration" field.
func ExpirationEQ(v int) predicate.Token {
	return predicate.Token(sql.FieldEQ(FieldExpiration, v))
}

// ExpirationNEQ applies the NEQ predicate on the "expiration" field.
func ExpirationNEQ(v int) predicate.Token {
	return predicate.Token(sql.FieldNEQ(FieldExpiration, v))
}

// ExpirationIn applies the In predicate on the "expiration" field.
func ExpirationIn(vs ...int) predicate.Token {
	return predicate.Token(sql.FieldIn(FieldExpiration, vs...))
}

// ExpirationNotIn applies the NotIn predicate on the "expiration" field.
func ExpirationNotIn(vs ...int) predicate.Token {
	return predicate.Token(sql.FieldNotIn(FieldExpiration, vs...))
}

// ExpirationGT applies the GT predicate on the "expiration" field.
func ExpirationGT(v int) predicate.Token {
	return predicate.Token(sql.FieldGT(FieldExpiration, v))
}

// ExpirationGTE applies the GTE predicate on the "expiration" field.
func ExpirationGTE(v int) predicate.Token {
	return predicate.Token(sql.FieldGTE(FieldExpiration, v))
}

// ExpirationLT applies the LT predicate on the "expiration" field.
func ExpirationLT(v int) predicate.Token {
	return predicate.Token(sql.FieldLT(FieldExpiration, v))
}

// ExpirationLTE applies the LTE predicate on the "expiration" field.
func ExpirationLTE(v int) predicate.Token {
	return predicate.Token(sql.FieldLTE(FieldExpiration, v))
}

// ExpirationIsNil applies the IsNil predicate on the "expiration" field.
func ExpirationIsNil() predicate.Token {
	return predicate.Token(sql.FieldIsNull(FieldExpiration))
}

// ExpirationNotNil applies the NotNil predicate on the "expiration" field.
func ExpirationNotNil() predicate.Token {
	return predicate.Token(sql.FieldNotNull(FieldExpiration))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Token) predicate.Token {
	return predicate.Token(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Token) predicate.Token {
	return predicate.Token(func(s *sql.Selector) {
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
func Not(p predicate.Token) predicate.Token {
	return predicate.Token(func(s *sql.Selector) {
		p(s.Not())
	})
}
