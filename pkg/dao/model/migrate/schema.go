// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ApplicationsColumns holds the columns for the "applications" table.
	ApplicationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "project_id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "environment_id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "modules", Type: field.TypeJSON},
	}
	// ApplicationsTable holds the schema information for the "applications" table.
	ApplicationsTable = &schema.Table{
		Name:       "applications",
		Columns:    ApplicationsColumns,
		PrimaryKey: []*schema.Column{ApplicationsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "application_update_time",
				Unique:  false,
				Columns: []*schema.Column{ApplicationsColumns[2]},
			},
		},
	}
	// ApplicationResourcesColumns holds the columns for the "application_resources" table.
	ApplicationResourcesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "application_id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "module", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
	}
	// ApplicationResourcesTable holds the schema information for the "application_resources" table.
	ApplicationResourcesTable = &schema.Table{
		Name:       "application_resources",
		Columns:    ApplicationResourcesColumns,
		PrimaryKey: []*schema.Column{ApplicationResourcesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "applicationresource_update_time",
				Unique:  false,
				Columns: []*schema.Column{ApplicationResourcesColumns[4]},
			},
		},
	}
	// ApplicationRevisionsColumns holds the columns for the "application_revisions" table.
	ApplicationRevisionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "application_id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "environment_id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "modules", Type: field.TypeJSON},
		{Name: "input_variables", Type: field.TypeJSON},
		{Name: "input_plan", Type: field.TypeString},
		{Name: "output", Type: field.TypeString},
	}
	// ApplicationRevisionsTable holds the schema information for the "application_revisions" table.
	ApplicationRevisionsTable = &schema.Table{
		Name:       "application_revisions",
		Columns:    ApplicationRevisionsColumns,
		PrimaryKey: []*schema.Column{ApplicationRevisionsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "applicationrevision_update_time",
				Unique:  false,
				Columns: []*schema.Column{ApplicationRevisionsColumns[4]},
			},
		},
	}
	// ConnectorsColumns holds the columns for the "connectors" table.
	ConnectorsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "driver", Type: field.TypeString},
		{Name: "config_version", Type: field.TypeString},
		{Name: "config_data", Type: field.TypeJSON, Nullable: true},
	}
	// ConnectorsTable holds the schema information for the "connectors" table.
	ConnectorsTable = &schema.Table{
		Name:       "connectors",
		Columns:    ConnectorsColumns,
		PrimaryKey: []*schema.Column{ConnectorsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "connector_update_time",
				Unique:  false,
				Columns: []*schema.Column{ConnectorsColumns[4]},
			},
		},
	}
	// EnvironmentsColumns holds the columns for the "environments" table.
	EnvironmentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "connector_ids", Type: field.TypeJSON, Nullable: true},
		{Name: "variables", Type: field.TypeJSON, Nullable: true},
	}
	// EnvironmentsTable holds the schema information for the "environments" table.
	EnvironmentsTable = &schema.Table{
		Name:       "environments",
		Columns:    EnvironmentsColumns,
		PrimaryKey: []*schema.Column{EnvironmentsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "environment_update_time",
				Unique:  false,
				Columns: []*schema.Column{EnvironmentsColumns[2]},
			},
		},
	}
	// ModulesColumns holds the columns for the "modules" table.
	ModulesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "source", Type: field.TypeString},
		{Name: "version", Type: field.TypeString},
		{Name: "input_schema", Type: field.TypeJSON, Nullable: true},
		{Name: "output_schema", Type: field.TypeJSON, Nullable: true},
	}
	// ModulesTable holds the schema information for the "modules" table.
	ModulesTable = &schema.Table{
		Name:       "modules",
		Columns:    ModulesColumns,
		PrimaryKey: []*schema.Column{ModulesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "module_update_time",
				Unique:  false,
				Columns: []*schema.Column{ModulesColumns[4]},
			},
		},
	}
	// ProjectsColumns holds the columns for the "projects" table.
	ProjectsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
	}
	// ProjectsTable holds the schema information for the "projects" table.
	ProjectsTable = &schema.Table{
		Name:       "projects",
		Columns:    ProjectsColumns,
		PrimaryKey: []*schema.Column{ProjectsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "project_update_time",
				Unique:  false,
				Columns: []*schema.Column{ProjectsColumns[2]},
			},
		},
	}
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
		ApplicationsTable,
		ApplicationResourcesTable,
		ApplicationRevisionsTable,
		ConnectorsTable,
		EnvironmentsTable,
		ModulesTable,
		ProjectsTable,
		RolesTable,
		SettingsTable,
		SubjectsTable,
		TokensTable,
	}
)

func init() {
}
