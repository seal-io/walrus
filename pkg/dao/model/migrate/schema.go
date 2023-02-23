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
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "environment_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "project_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ApplicationsTable holds the schema information for the "applications" table.
	ApplicationsTable = &schema.Table{
		Name:       "applications",
		Columns:    ApplicationsColumns,
		PrimaryKey: []*schema.Column{ApplicationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "applications_environments_applications",
				Columns:    []*schema.Column{ApplicationsColumns[6]},
				RefColumns: []*schema.Column{EnvironmentsColumns[0]},
				OnDelete:   schema.Restrict,
			},
			{
				Symbol:     "applications_projects_applications",
				Columns:    []*schema.Column{ApplicationsColumns[7]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "application_update_time",
				Unique:  false,
				Columns: []*schema.Column{ApplicationsColumns[5]},
			},
			{
				Name:    "application_project_id_name",
				Unique:  true,
				Columns: []*schema.Column{ApplicationsColumns[7], ApplicationsColumns[1]},
			},
		},
	}
	// ApplicationModuleRelationshipsColumns holds the columns for the "application_module_relationships" table.
	ApplicationModuleRelationshipsColumns = []*schema.Column{
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "variables", Type: field.TypeJSON, Nullable: true},
		{Name: "application_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "module_id", Type: field.TypeString},
	}
	// ApplicationModuleRelationshipsTable holds the schema information for the "application_module_relationships" table.
	ApplicationModuleRelationshipsTable = &schema.Table{
		Name:       "application_module_relationships",
		Columns:    ApplicationModuleRelationshipsColumns,
		PrimaryKey: []*schema.Column{ApplicationModuleRelationshipsColumns[4], ApplicationModuleRelationshipsColumns[5]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "application_module_relationships_applications_application",
				Columns:    []*schema.Column{ApplicationModuleRelationshipsColumns[4]},
				RefColumns: []*schema.Column{ApplicationsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "application_module_relationships_modules_module",
				Columns:    []*schema.Column{ApplicationModuleRelationshipsColumns[5]},
				RefColumns: []*schema.Column{ModulesColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "applicationmodulerelationship_update_time",
				Unique:  false,
				Columns: []*schema.Column{ApplicationModuleRelationshipsColumns[1]},
			},
		},
	}
	// ApplicationResourcesColumns holds the columns for the "application_resources" table.
	ApplicationResourcesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "module", Type: field.TypeString},
		{Name: "mode", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "application_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "connector_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ApplicationResourcesTable holds the schema information for the "application_resources" table.
	ApplicationResourcesTable = &schema.Table{
		Name:       "application_resources",
		Columns:    ApplicationResourcesColumns,
		PrimaryKey: []*schema.Column{ApplicationResourcesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "application_resources_applications_resources",
				Columns:    []*schema.Column{ApplicationResourcesColumns[9]},
				RefColumns: []*schema.Column{ApplicationsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "application_resources_connectors_resources",
				Columns:    []*schema.Column{ApplicationResourcesColumns[10]},
				RefColumns: []*schema.Column{ConnectorsColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "applicationresource_update_time",
				Unique:  false,
				Columns: []*schema.Column{ApplicationResourcesColumns[4]},
			},
			{
				Name:    "applicationresource_application_id_connector_id_module_mode_type_name",
				Unique:  true,
				Columns: []*schema.Column{ApplicationResourcesColumns[9], ApplicationResourcesColumns[10], ApplicationResourcesColumns[5], ApplicationResourcesColumns[6], ApplicationResourcesColumns[7], ApplicationResourcesColumns[8]},
			},
		},
	}
	// ApplicationRevisionsColumns holds the columns for the "application_revisions" table.
	ApplicationRevisionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "modules", Type: field.TypeJSON},
		{Name: "input_variables", Type: field.TypeJSON},
		{Name: "input_plan", Type: field.TypeString},
		{Name: "output", Type: field.TypeString},
		{Name: "application_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "environment_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ApplicationRevisionsTable holds the schema information for the "application_revisions" table.
	ApplicationRevisionsTable = &schema.Table{
		Name:       "application_revisions",
		Columns:    ApplicationRevisionsColumns,
		PrimaryKey: []*schema.Column{ApplicationRevisionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "application_revisions_applications_revisions",
				Columns:    []*schema.Column{ApplicationRevisionsColumns[9]},
				RefColumns: []*schema.Column{ApplicationsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "application_revisions_environments_revisions",
				Columns:    []*schema.Column{ApplicationRevisionsColumns[10]},
				RefColumns: []*schema.Column{EnvironmentsColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
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
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON, Nullable: true},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "type", Type: field.TypeString},
		{Name: "config_version", Type: field.TypeString},
		{Name: "config_data", Type: field.TypeJSON, Nullable: true},
		{Name: "enable_fin_ops", Type: field.TypeBool},
		{Name: "fin_ops_status", Type: field.TypeString, Nullable: true},
		{Name: "fin_ops_status_message", Type: field.TypeString, Nullable: true},
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
				Columns: []*schema.Column{ConnectorsColumns[7]},
			},
			{
				Name:    "connector_name",
				Unique:  true,
				Columns: []*schema.Column{ConnectorsColumns[1]},
			},
		},
	}
	// EnvironmentsColumns holds the columns for the "environments" table.
	EnvironmentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
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
				Columns: []*schema.Column{EnvironmentsColumns[5]},
			},
			{
				Name:    "environment_name",
				Unique:  true,
				Columns: []*schema.Column{EnvironmentsColumns[1]},
			},
		},
	}
	// EnvironmentConnectorRelationshipsColumns holds the columns for the "environment_connector_relationships" table.
	EnvironmentConnectorRelationshipsColumns = []*schema.Column{
		{Name: "create_time", Type: field.TypeTime},
		{Name: "environment_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "connector_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// EnvironmentConnectorRelationshipsTable holds the schema information for the "environment_connector_relationships" table.
	EnvironmentConnectorRelationshipsTable = &schema.Table{
		Name:       "environment_connector_relationships",
		Columns:    EnvironmentConnectorRelationshipsColumns,
		PrimaryKey: []*schema.Column{EnvironmentConnectorRelationshipsColumns[1], EnvironmentConnectorRelationshipsColumns[2]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "environment_connector_relationships_environments_environment",
				Columns:    []*schema.Column{EnvironmentConnectorRelationshipsColumns[1]},
				RefColumns: []*schema.Column{EnvironmentsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "environment_connector_relationships_connectors_connector",
				Columns:    []*schema.Column{EnvironmentConnectorRelationshipsColumns[2]},
				RefColumns: []*schema.Column{ConnectorsColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
	}
	// ModulesColumns holds the columns for the "modules" table.
	ModulesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON, Nullable: true},
		{Name: "source", Type: field.TypeString},
		{Name: "version", Type: field.TypeString, Nullable: true},
		{Name: "schema", Type: field.TypeJSON, Nullable: true},
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
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON, Nullable: true},
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
				Columns: []*schema.Column{ProjectsColumns[5]},
			},
			{
				Name:    "project_name",
				Unique:  true,
				Columns: []*schema.Column{ProjectsColumns[1]},
			},
		},
	}
	// RolesColumns holds the columns for the "roles" table.
	RolesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
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
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
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
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
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
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
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
		ApplicationModuleRelationshipsTable,
		ApplicationResourcesTable,
		ApplicationRevisionsTable,
		ConnectorsTable,
		EnvironmentsTable,
		EnvironmentConnectorRelationshipsTable,
		ModulesTable,
		ProjectsTable,
		RolesTable,
		SettingsTable,
		SubjectsTable,
		TokensTable,
	}
)

func init() {
	ApplicationsTable.ForeignKeys[0].RefTable = EnvironmentsTable
	ApplicationsTable.ForeignKeys[1].RefTable = ProjectsTable
	ApplicationModuleRelationshipsTable.ForeignKeys[0].RefTable = ApplicationsTable
	ApplicationModuleRelationshipsTable.ForeignKeys[1].RefTable = ModulesTable
	ApplicationResourcesTable.ForeignKeys[0].RefTable = ApplicationsTable
	ApplicationResourcesTable.ForeignKeys[1].RefTable = ConnectorsTable
	ApplicationRevisionsTable.ForeignKeys[0].RefTable = ApplicationsTable
	ApplicationRevisionsTable.ForeignKeys[1].RefTable = EnvironmentsTable
	EnvironmentConnectorRelationshipsTable.ForeignKeys[0].RefTable = EnvironmentsTable
	EnvironmentConnectorRelationshipsTable.ForeignKeys[1].RefTable = ConnectorsTable
}
