// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// RolesColumns holds the columns for the "roles" table.
	RolesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "domain", Type: field.TypeString, Default: "system"},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Default: ""},
		{Name: "policies", Type: field.TypeJSON},
		{Name: "builtin", Type: field.TypeBool, Default: false},
		{Name: "session", Type: field.TypeBool, Default: false},
	}
	// RolesTable holds the schema information for the "roles" table.
	RolesTable = &schema.Table{
		Name:       "roles",
		Columns:    RolesColumns,
		PrimaryKey: []*schema.Column{RolesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "role_update_time",
				Unique:  false,
				Columns: []*schema.Column{RolesColumns[2]},
			},
			{
				Name:    "role_domain_name",
				Unique:  true,
				Columns: []*schema.Column{RolesColumns[3], RolesColumns[4]},
			},
		},
	}
	// SettingsColumns holds the columns for the "settings" table.
	SettingsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "value", Type: field.TypeString, Default: ""},
		{Name: "hidden", Type: field.TypeBool, Default: false},
		{Name: "editable", Type: field.TypeBool, Default: false},
		{Name: "private", Type: field.TypeBool, Default: false},
	}
	// SettingsTable holds the schema information for the "settings" table.
	SettingsTable = &schema.Table{
		Name:       "settings",
		Columns:    SettingsColumns,
		PrimaryKey: []*schema.Column{SettingsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "setting_update_time",
				Unique:  false,
				Columns: []*schema.Column{SettingsColumns[2]},
			},
			{
				Name:    "setting_name",
				Unique:  true,
				Columns: []*schema.Column{SettingsColumns[3]},
			},
		},
	}
	// SubjectsColumns holds the columns for the "subjects" table.
	SubjectsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "kind", Type: field.TypeString, Default: "user"},
		{Name: "group", Type: field.TypeString, Default: "default"},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Default: ""},
		{Name: "mount_to", Type: field.TypeBool, Default: false},
		{Name: "login_to", Type: field.TypeBool, Default: true},
		{Name: "roles", Type: field.TypeJSON},
		{Name: "paths", Type: field.TypeJSON},
		{Name: "builtin", Type: field.TypeBool, Default: false},
	}
	// SubjectsTable holds the schema information for the "subjects" table.
	SubjectsTable = &schema.Table{
		Name:       "subjects",
		Columns:    SubjectsColumns,
		PrimaryKey: []*schema.Column{SubjectsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "subject_update_time",
				Unique:  false,
				Columns: []*schema.Column{SubjectsColumns[2]},
			},
			{
				Name:    "subject_kind_group_name",
				Unique:  true,
				Columns: []*schema.Column{SubjectsColumns[3], SubjectsColumns[4], SubjectsColumns[5]},
			},
		},
	}
	// TokensColumns holds the columns for the "tokens" table.
	TokensColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "casdoor_token_name", Type: field.TypeString},
		{Name: "casdoor_token_owner", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "expiration", Type: field.TypeInt, Nullable: true},
	}
	// TokensTable holds the schema information for the "tokens" table.
	TokensTable = &schema.Table{
		Name:       "tokens",
		Columns:    TokensColumns,
		PrimaryKey: []*schema.Column{TokensColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "token_update_time",
				Unique:  false,
				Columns: []*schema.Column{TokensColumns[2]},
			},
			{
				Name:    "token_casdoor_token_name",
				Unique:  true,
				Columns: []*schema.Column{TokensColumns[3]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		RolesTable,
		SettingsTable,
		SubjectsTable,
		TokensTable,
	}
)

func init() {
}
