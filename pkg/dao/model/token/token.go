// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package token

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the token type in the database.
	Label = "token"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldCasdoorTokenName holds the string denoting the casdoortokenname field in the database.
	FieldCasdoorTokenName = "casdoor_token_name"
	// FieldCasdoorTokenOwner holds the string denoting the casdoortokenowner field in the database.
	FieldCasdoorTokenOwner = "casdoor_token_owner"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldExpiration holds the string denoting the expiration field in the database.
	FieldExpiration = "expiration"
	// Table holds the table name of the token in the database.
	Table = "tokens"
)

// Columns holds all SQL columns for token fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldCasdoorTokenName,
	FieldCasdoorTokenOwner,
	FieldName,
	FieldExpiration,
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
	// CasdoorTokenNameValidator is a validator for the "casdoorTokenName" field. It is called by the builders before save.
	CasdoorTokenNameValidator func(string) error
	// CasdoorTokenOwnerValidator is a validator for the "casdoorTokenOwner" field. It is called by the builders before save.
	CasdoorTokenOwnerValidator func(string) error
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
)

// OrderOption defines the ordering options for the Token queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreateTime orders the results by the createTime field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the updateTime field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByCasdoorTokenName orders the results by the casdoorTokenName field.
func ByCasdoorTokenName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCasdoorTokenName, opts...).ToFunc()
}

// ByCasdoorTokenOwner orders the results by the casdoorTokenOwner field.
func ByCasdoorTokenOwner(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCasdoorTokenOwner, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByExpiration orders the results by the expiration field.
func ByExpiration(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExpiration, opts...).ToFunc()
}

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
