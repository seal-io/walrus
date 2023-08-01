// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package distributelock

import (
	"entgo.io/ent/dialect/sql"
	"golang.org/x/exp/slices"
)

const (
	// Label holds the string label denoting the distributelock type in the database.
	Label = "distribute_lock"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldExpireAt holds the string denoting the expireat field in the database.
	FieldExpireAt = "expire_at"
	// Table holds the table name of the distributelock in the database.
	Table = "distribute_locks"
)

// Columns holds all SQL columns for distributelock fields.
var Columns = []string{
	FieldID,
	FieldExpireAt,
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

var (
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(string) error
)

// OrderOption defines the ordering options for the DistributeLock queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByExpireAt orders the results by the expireAt field.
func ByExpireAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExpireAt, opts...).ToFunc()
}

// WithoutFields returns the fields ignored the given list.
func WithoutFields(ignores ...string) []string {
	if len(ignores) == 0 {
		return slices.Clone(Columns)
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
