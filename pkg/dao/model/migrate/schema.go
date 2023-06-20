// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AllocationCostsColumns holds the columns for the "allocation_costs" table.
	AllocationCostsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "start_time", Type: field.TypeTime},
		{Name: "end_time", Type: field.TypeTime},
		{Name: "minutes", Type: field.TypeFloat64},
		{Name: "name", Type: field.TypeString},
		{Name: "fingerprint", Type: field.TypeString},
		{Name: "cluster_name", Type: field.TypeString, Nullable: true},
		{Name: "namespace", Type: field.TypeString, Nullable: true},
		{Name: "node", Type: field.TypeString, Nullable: true},
		{Name: "controller", Type: field.TypeString, Nullable: true},
		{Name: "controller_kind", Type: field.TypeString, Nullable: true},
		{Name: "pod", Type: field.TypeString, Nullable: true},
		{Name: "container", Type: field.TypeString, Nullable: true},
		{Name: "pvs", Type: field.TypeJSON},
		{Name: "labels", Type: field.TypeJSON},
		{Name: "total_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "currency", Type: field.TypeInt, Nullable: true},
		{Name: "cpu_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "cpu_core_request", Type: field.TypeFloat64, Default: 0},
		{Name: "gpu_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "gpu_count", Type: field.TypeFloat64, Default: 0},
		{Name: "ram_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "ram_byte_request", Type: field.TypeFloat64, Default: 0},
		{Name: "pv_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "pv_bytes", Type: field.TypeFloat64, Default: 0},
		{Name: "load_balancer_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "cpu_core_usage_average", Type: field.TypeFloat64, Default: 0},
		{Name: "cpu_core_usage_max", Type: field.TypeFloat64, Default: 0},
		{Name: "ram_byte_usage_average", Type: field.TypeFloat64, Default: 0},
		{Name: "ram_byte_usage_max", Type: field.TypeFloat64, Default: 0},
		{Name: "connector_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// AllocationCostsTable holds the schema information for the "allocation_costs" table.
	AllocationCostsTable = &schema.Table{
		Name:       "allocation_costs",
		Columns:    AllocationCostsColumns,
		PrimaryKey: []*schema.Column{AllocationCostsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "allocation_costs_connectors_allocationCosts",
				Columns:    []*schema.Column{AllocationCostsColumns[30]},
				RefColumns: []*schema.Column{ConnectorsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "allocationcost_start_time_end_time_connector_id",
				Unique:  false,
				Columns: []*schema.Column{AllocationCostsColumns[1], AllocationCostsColumns[2], AllocationCostsColumns[30]},
			},
			{
				Name:    "allocationcost_start_time_end_time_connector_id_fingerprint",
				Unique:  true,
				Columns: []*schema.Column{AllocationCostsColumns[1], AllocationCostsColumns[2], AllocationCostsColumns[30], AllocationCostsColumns[5]},
			},
		},
	}
	// ClusterCostsColumns holds the columns for the "cluster_costs" table.
	ClusterCostsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "start_time", Type: field.TypeTime},
		{Name: "end_time", Type: field.TypeTime},
		{Name: "minutes", Type: field.TypeFloat64},
		{Name: "cluster_name", Type: field.TypeString},
		{Name: "total_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "currency", Type: field.TypeInt, Nullable: true},
		{Name: "allocation_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "idle_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "management_cost", Type: field.TypeFloat64, Default: 0},
		{Name: "connector_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ClusterCostsTable holds the schema information for the "cluster_costs" table.
	ClusterCostsTable = &schema.Table{
		Name:       "cluster_costs",
		Columns:    ClusterCostsColumns,
		PrimaryKey: []*schema.Column{ClusterCostsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "cluster_costs_connectors_clusterCosts",
				Columns:    []*schema.Column{ClusterCostsColumns[10]},
				RefColumns: []*schema.Column{ConnectorsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "clustercost_start_time_end_time_connector_id",
				Unique:  true,
				Columns: []*schema.Column{ClusterCostsColumns[1], ClusterCostsColumns[2], ClusterCostsColumns[10]},
			},
		},
	}
	// ConnectorsColumns holds the columns for the "connectors" table.
	ConnectorsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "status", Type: field.TypeJSON, Nullable: true},
		{Name: "type", Type: field.TypeString},
		{Name: "config_version", Type: field.TypeString},
		{Name: "config_data", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "blob", "postgres": "bytea", "sqlite3": "blob"}},
		{Name: "enable_fin_ops", Type: field.TypeBool},
		{Name: "fin_ops_custom_pricing", Type: field.TypeJSON, Nullable: true},
		{Name: "category", Type: field.TypeString},
		{Name: "project_id", Type: field.TypeString, Nullable: true, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ConnectorsTable holds the schema information for the "connectors" table.
	ConnectorsTable = &schema.Table{
		Name:       "connectors",
		Columns:    ConnectorsColumns,
		PrimaryKey: []*schema.Column{ConnectorsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "connectors_projects_connectors",
				Columns:    []*schema.Column{ConnectorsColumns[13]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "connector_update_time",
				Unique:  false,
				Columns: []*schema.Column{ConnectorsColumns[5]},
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
		{Name: "labels", Type: field.TypeJSON},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "project_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// EnvironmentsTable holds the schema information for the "environments" table.
	EnvironmentsTable = &schema.Table{
		Name:       "environments",
		Columns:    EnvironmentsColumns,
		PrimaryKey: []*schema.Column{EnvironmentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "environments_projects_environments",
				Columns:    []*schema.Column{EnvironmentsColumns[6]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "environment_update_time",
				Unique:  false,
				Columns: []*schema.Column{EnvironmentsColumns[5]},
			},
			{
				Name:    "environment_project_id_name",
				Unique:  true,
				Columns: []*schema.Column{EnvironmentsColumns[6], EnvironmentsColumns[1]},
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
	// PerspectivesColumns holds the columns for the "perspectives" table.
	PerspectivesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "start_time", Type: field.TypeString},
		{Name: "end_time", Type: field.TypeString},
		{Name: "builtin", Type: field.TypeBool, Default: false},
		{Name: "allocation_queries", Type: field.TypeJSON},
	}
	// PerspectivesTable holds the schema information for the "perspectives" table.
	PerspectivesTable = &schema.Table{
		Name:       "perspectives",
		Columns:    PerspectivesColumns,
		PrimaryKey: []*schema.Column{PerspectivesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "perspective_update_time",
				Unique:  false,
				Columns: []*schema.Column{PerspectivesColumns[2]},
			},
			{
				Name:    "perspective_name",
				Unique:  true,
				Columns: []*schema.Column{PerspectivesColumns[3]},
			},
		},
	}
	// ProjectsColumns holds the columns for the "projects" table.
	ProjectsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON},
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
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "kind", Type: field.TypeString, Default: "system"},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "policies", Type: field.TypeJSON},
		{Name: "session", Type: field.TypeBool, Default: false},
		{Name: "builtin", Type: field.TypeBool, Default: false},
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
		},
	}
	// SecretsColumns holds the columns for the "secrets" table.
	SecretsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "value", Type: field.TypeString, SchemaType: map[string]string{"mysql": "blob", "postgres": "bytea", "sqlite3": "blob"}},
		{Name: "project_id", Type: field.TypeString, Nullable: true, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// SecretsTable holds the schema information for the "secrets" table.
	SecretsTable = &schema.Table{
		Name:       "secrets",
		Columns:    SecretsColumns,
		PrimaryKey: []*schema.Column{SecretsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "secrets_projects_secrets",
				Columns:    []*schema.Column{SecretsColumns[5]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "secret_update_time",
				Unique:  false,
				Columns: []*schema.Column{SecretsColumns[2]},
			},
			{
				Name:    "secret_project_id_name",
				Unique:  true,
				Columns: []*schema.Column{SecretsColumns[5], SecretsColumns[3]},
				Annotation: &entsql.IndexAnnotation{
					Where: "project_id IS NOT NULL",
				},
			},
			{
				Name:    "secret_name",
				Unique:  true,
				Columns: []*schema.Column{SecretsColumns[3]},
				Annotation: &entsql.IndexAnnotation{
					Where: "project_id IS NULL",
				},
			},
		},
	}
	// ServicesColumns holds the columns for the "services" table.
	ServicesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "template", Type: field.TypeJSON},
		{Name: "attributes", Type: field.TypeOther, Nullable: true, SchemaType: map[string]string{"mysql": "json", "postgres": "jsonb", "sqlite3": "text"}},
		{Name: "status", Type: field.TypeJSON, Nullable: true},
		{Name: "environment_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "project_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ServicesTable holds the schema information for the "services" table.
	ServicesTable = &schema.Table{
		Name:       "services",
		Columns:    ServicesColumns,
		PrimaryKey: []*schema.Column{ServicesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "services_environments_services",
				Columns:    []*schema.Column{ServicesColumns[9]},
				RefColumns: []*schema.Column{EnvironmentsColumns[0]},
				OnDelete:   schema.Restrict,
			},
			{
				Symbol:     "services_projects_services",
				Columns:    []*schema.Column{ServicesColumns[10]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "service_update_time",
				Unique:  false,
				Columns: []*schema.Column{ServicesColumns[5]},
			},
			{
				Name:    "service_environment_id_name",
				Unique:  true,
				Columns: []*schema.Column{ServicesColumns[9], ServicesColumns[1]},
			},
		},
	}
	// ServiceResourcesColumns holds the columns for the "service_resources" table.
	ServiceResourcesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "project_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "mode", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "deployer_type", Type: field.TypeString},
		{Name: "status", Type: field.TypeJSON, Nullable: true},
		{Name: "connector_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "service_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "composition_id", Type: field.TypeString, Nullable: true, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ServiceResourcesTable holds the schema information for the "service_resources" table.
	ServiceResourcesTable = &schema.Table{
		Name:       "service_resources",
		Columns:    ServiceResourcesColumns,
		PrimaryKey: []*schema.Column{ServiceResourcesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "service_resources_connectors_resources",
				Columns:    []*schema.Column{ServiceResourcesColumns[9]},
				RefColumns: []*schema.Column{ConnectorsColumns[0]},
				OnDelete:   schema.Restrict,
			},
			{
				Symbol:     "service_resources_services_resources",
				Columns:    []*schema.Column{ServiceResourcesColumns[10]},
				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "service_resources_service_resources_components",
				Columns:    []*schema.Column{ServiceResourcesColumns[11]},
				RefColumns: []*schema.Column{ServiceResourcesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "serviceresource_update_time",
				Unique:  false,
				Columns: []*schema.Column{ServiceResourcesColumns[3]},
			},
		},
	}
	// ServiceRevisionsColumns holds the columns for the "service_revisions" table.
	ServiceRevisionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "template_id", Type: field.TypeString},
		{Name: "template_version", Type: field.TypeString},
		{Name: "attributes", Type: field.TypeOther, Nullable: true, SchemaType: map[string]string{"mysql": "json", "postgres": "jsonb", "sqlite3": "text"}},
		{Name: "secrets", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "blob", "postgres": "bytea", "sqlite3": "blob"}},
		{Name: "input_plan", Type: field.TypeString},
		{Name: "output", Type: field.TypeString},
		{Name: "deployer_type", Type: field.TypeString, Default: "Terraform"},
		{Name: "duration", Type: field.TypeInt, Default: 0},
		{Name: "previous_required_providers", Type: field.TypeJSON},
		{Name: "tags", Type: field.TypeJSON},
		{Name: "environment_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "project_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "service_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ServiceRevisionsTable holds the schema information for the "service_revisions" table.
	ServiceRevisionsTable = &schema.Table{
		Name:       "service_revisions",
		Columns:    ServiceRevisionsColumns,
		PrimaryKey: []*schema.Column{ServiceRevisionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "service_revisions_environments_serviceRevisions",
				Columns:    []*schema.Column{ServiceRevisionsColumns[14]},
				RefColumns: []*schema.Column{EnvironmentsColumns[0]},
				OnDelete:   schema.Restrict,
			},
			{
				Symbol:     "service_revisions_projects_serviceRevisions",
				Columns:    []*schema.Column{ServiceRevisionsColumns[15]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.Restrict,
			},
			{
				Symbol:     "service_revisions_services_revisions",
				Columns:    []*schema.Column{ServiceRevisionsColumns[16]},
				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// SettingsColumns holds the columns for the "settings" table.
	SettingsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "value", Type: field.TypeString},
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
		{Name: "domain", Type: field.TypeString, Default: "builtin"},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
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
				Name:    "subject_kind_domain_name",
				Unique:  true,
				Columns: []*schema.Column{SubjectsColumns[3], SubjectsColumns[4], SubjectsColumns[5]},
			},
		},
	}
	// SubjectRoleRelationshipsColumns holds the columns for the "subject_role_relationships" table.
	SubjectRoleRelationshipsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "project_id", Type: field.TypeString, Nullable: true, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "subject_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "role_id", Type: field.TypeString},
	}
	// SubjectRoleRelationshipsTable holds the schema information for the "subject_role_relationships" table.
	SubjectRoleRelationshipsTable = &schema.Table{
		Name:       "subject_role_relationships",
		Columns:    SubjectRoleRelationshipsColumns,
		PrimaryKey: []*schema.Column{SubjectRoleRelationshipsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "subject_role_relationships_projects_subjectRoles",
				Columns:    []*schema.Column{SubjectRoleRelationshipsColumns[2]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "subject_role_relationships_subjects_subject",
				Columns:    []*schema.Column{SubjectRoleRelationshipsColumns[3]},
				RefColumns: []*schema.Column{SubjectsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "subject_role_relationships_roles_role",
				Columns:    []*schema.Column{SubjectRoleRelationshipsColumns[4]},
				RefColumns: []*schema.Column{RolesColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "subjectrolerelationship_project_id_subject_id_role_id",
				Unique:  true,
				Columns: []*schema.Column{SubjectRoleRelationshipsColumns[2], SubjectRoleRelationshipsColumns[3], SubjectRoleRelationshipsColumns[4]},
				Annotation: &entsql.IndexAnnotation{
					Where: "project_id IS NOT NULL",
				},
			},
			{
				Name:    "subjectrolerelationship_subject_id_role_id",
				Unique:  true,
				Columns: []*schema.Column{SubjectRoleRelationshipsColumns[3], SubjectRoleRelationshipsColumns[4]},
				Annotation: &entsql.IndexAnnotation{
					Where: "project_id IS NULL",
				},
			},
		},
	}
	// TemplatesColumns holds the columns for the "templates" table.
	TemplatesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "icon", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON},
		{Name: "source", Type: field.TypeString},
	}
	// TemplatesTable holds the schema information for the "templates" table.
	TemplatesTable = &schema.Table{
		Name:       "templates",
		Columns:    TemplatesColumns,
		PrimaryKey: []*schema.Column{TemplatesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "template_update_time",
				Unique:  false,
				Columns: []*schema.Column{TemplatesColumns[4]},
			},
		},
	}
	// TemplateVersionsColumns holds the columns for the "template_versions" table.
	TemplateVersionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "version", Type: field.TypeString},
		{Name: "source", Type: field.TypeString},
		{Name: "schema", Type: field.TypeJSON},
		{Name: "template_id", Type: field.TypeString},
	}
	// TemplateVersionsTable holds the schema information for the "template_versions" table.
	TemplateVersionsTable = &schema.Table{
		Name:       "template_versions",
		Columns:    TemplateVersionsColumns,
		PrimaryKey: []*schema.Column{TemplateVersionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "template_versions_templates_versions",
				Columns:    []*schema.Column{TemplateVersionsColumns[6]},
				RefColumns: []*schema.Column{TemplatesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "templateversion_update_time",
				Unique:  false,
				Columns: []*schema.Column{TemplateVersionsColumns[2]},
			},
		},
	}
	// TokensColumns holds the columns for the "tokens" table.
	TokensColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "kind", Type: field.TypeString, Default: "api"},
		{Name: "name", Type: field.TypeString},
		{Name: "expiration", Type: field.TypeTime, Nullable: true},
		{Name: "value", Type: field.TypeString, SchemaType: map[string]string{"mysql": "blob", "postgres": "bytea", "sqlite3": "blob"}},
		{Name: "subject_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// TokensTable holds the schema information for the "tokens" table.
	TokensTable = &schema.Table{
		Name:       "tokens",
		Columns:    TokensColumns,
		PrimaryKey: []*schema.Column{TokensColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "tokens_subjects_tokens",
				Columns:    []*schema.Column{TokensColumns[6]},
				RefColumns: []*schema.Column{SubjectsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AllocationCostsTable,
		ClusterCostsTable,
		ConnectorsTable,
		EnvironmentsTable,
		EnvironmentConnectorRelationshipsTable,
		PerspectivesTable,
		ProjectsTable,
		RolesTable,
		SecretsTable,
		ServicesTable,
		ServiceResourcesTable,
		ServiceRevisionsTable,
		SettingsTable,
		SubjectsTable,
		SubjectRoleRelationshipsTable,
		TemplatesTable,
		TemplateVersionsTable,
		TokensTable,
	}
)

func init() {
	AllocationCostsTable.ForeignKeys[0].RefTable = ConnectorsTable
	ClusterCostsTable.ForeignKeys[0].RefTable = ConnectorsTable
	ConnectorsTable.ForeignKeys[0].RefTable = ProjectsTable
	EnvironmentsTable.ForeignKeys[0].RefTable = ProjectsTable
	EnvironmentConnectorRelationshipsTable.ForeignKeys[0].RefTable = EnvironmentsTable
	EnvironmentConnectorRelationshipsTable.ForeignKeys[1].RefTable = ConnectorsTable
	SecretsTable.ForeignKeys[0].RefTable = ProjectsTable
	ServicesTable.ForeignKeys[0].RefTable = EnvironmentsTable
	ServicesTable.ForeignKeys[1].RefTable = ProjectsTable
	ServiceResourcesTable.ForeignKeys[0].RefTable = ConnectorsTable
	ServiceResourcesTable.ForeignKeys[1].RefTable = ServicesTable
	ServiceResourcesTable.ForeignKeys[2].RefTable = ServiceResourcesTable
	ServiceRevisionsTable.ForeignKeys[0].RefTable = EnvironmentsTable
	ServiceRevisionsTable.ForeignKeys[1].RefTable = ProjectsTable
	ServiceRevisionsTable.ForeignKeys[2].RefTable = ServicesTable
	SubjectRoleRelationshipsTable.ForeignKeys[0].RefTable = ProjectsTable
	SubjectRoleRelationshipsTable.ForeignKeys[1].RefTable = SubjectsTable
	SubjectRoleRelationshipsTable.ForeignKeys[2].RefTable = RolesTable
	TemplateVersionsTable.ForeignKeys[0].RefTable = TemplatesTable
	TokensTable.ForeignKeys[0].RefTable = SubjectsTable
}
