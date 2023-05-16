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
	// ApplicationsColumns holds the columns for the "applications" table.
	ApplicationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "variables", Type: field.TypeOther, Nullable: true, SchemaType: map[string]string{"mysql": "json", "postgres": "jsonb", "sqlite3": "text"}},
		{Name: "project_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ApplicationsTable holds the schema information for the "applications" table.
	ApplicationsTable = &schema.Table{
		Name:       "applications",
		Columns:    ApplicationsColumns,
		PrimaryKey: []*schema.Column{ApplicationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
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
	// ApplicationInstancesColumns holds the columns for the "application_instances" table.
	ApplicationInstancesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "variables", Type: field.TypeOther, Nullable: true, SchemaType: map[string]string{"mysql": "json", "postgres": "jsonb", "sqlite3": "text"}},
		{Name: "status", Type: field.TypeJSON, Nullable: true},
		{Name: "application_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "environment_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ApplicationInstancesTable holds the schema information for the "application_instances" table.
	ApplicationInstancesTable = &schema.Table{
		Name:       "application_instances",
		Columns:    ApplicationInstancesColumns,
		PrimaryKey: []*schema.Column{ApplicationInstancesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "application_instances_applications_instances",
				Columns:    []*schema.Column{ApplicationInstancesColumns[6]},
				RefColumns: []*schema.Column{ApplicationsColumns[0]},
				OnDelete:   schema.Restrict,
			},
			{
				Symbol:     "application_instances_environments_instances",
				Columns:    []*schema.Column{ApplicationInstancesColumns[7]},
				RefColumns: []*schema.Column{EnvironmentsColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "applicationinstance_update_time",
				Unique:  false,
				Columns: []*schema.Column{ApplicationInstancesColumns[2]},
			},
			{
				Name:    "applicationinstance_application_id_environment_id_name",
				Unique:  true,
				Columns: []*schema.Column{ApplicationInstancesColumns[6], ApplicationInstancesColumns[7], ApplicationInstancesColumns[3]},
			},
		},
	}
	// ApplicationModuleRelationshipsColumns holds the columns for the "application_module_relationships" table.
	ApplicationModuleRelationshipsColumns = []*schema.Column{
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "version", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "attributes", Type: field.TypeOther, Nullable: true, SchemaType: map[string]string{"mysql": "json", "postgres": "jsonb", "sqlite3": "text"}},
		{Name: "application_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "module_id", Type: field.TypeString},
	}
	// ApplicationModuleRelationshipsTable holds the schema information for the "application_module_relationships" table.
	ApplicationModuleRelationshipsTable = &schema.Table{
		Name:       "application_module_relationships",
		Columns:    ApplicationModuleRelationshipsColumns,
		PrimaryKey: []*schema.Column{ApplicationModuleRelationshipsColumns[5], ApplicationModuleRelationshipsColumns[6], ApplicationModuleRelationshipsColumns[3]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "application_module_relationships_applications_application",
				Columns:    []*schema.Column{ApplicationModuleRelationshipsColumns[5]},
				RefColumns: []*schema.Column{ApplicationsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "application_module_relationships_modules_module",
				Columns:    []*schema.Column{ApplicationModuleRelationshipsColumns[6]},
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
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "module", Type: field.TypeString},
		{Name: "mode", Type: field.TypeString},
		{Name: "type", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "deployer_type", Type: field.TypeString},
		{Name: "status", Type: field.TypeJSON, Nullable: true},
		{Name: "instance_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "composition_id", Type: field.TypeString, Nullable: true, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "connector_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ApplicationResourcesTable holds the schema information for the "application_resources" table.
	ApplicationResourcesTable = &schema.Table{
		Name:       "application_resources",
		Columns:    ApplicationResourcesColumns,
		PrimaryKey: []*schema.Column{ApplicationResourcesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "application_resources_application_instances_resources",
				Columns:    []*schema.Column{ApplicationResourcesColumns[9]},
				RefColumns: []*schema.Column{ApplicationInstancesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "application_resources_application_resources_components",
				Columns:    []*schema.Column{ApplicationResourcesColumns[10]},
				RefColumns: []*schema.Column{ApplicationResourcesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "application_resources_connectors_resources",
				Columns:    []*schema.Column{ApplicationResourcesColumns[11]},
				RefColumns: []*schema.Column{ConnectorsColumns[0]},
				OnDelete:   schema.Restrict,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "applicationresource_update_time",
				Unique:  false,
				Columns: []*schema.Column{ApplicationResourcesColumns[2]},
			},
		},
	}
	// ApplicationRevisionsColumns holds the columns for the "application_revisions" table.
	ApplicationRevisionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "status", Type: field.TypeString, Nullable: true},
		{Name: "status_message", Type: field.TypeString, Nullable: true},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "modules", Type: field.TypeJSON},
		{Name: "secrets", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "blob", "postgres": "bytea", "sqlite3": "blob"}},
		{Name: "variables", Type: field.TypeOther, Nullable: true, SchemaType: map[string]string{"mysql": "json", "postgres": "jsonb", "sqlite3": "text"}},
		{Name: "input_variables", Type: field.TypeOther, SchemaType: map[string]string{"mysql": "json", "postgres": "jsonb", "sqlite3": "text"}},
		{Name: "input_plan", Type: field.TypeString},
		{Name: "output", Type: field.TypeString},
		{Name: "deployer_type", Type: field.TypeString, Default: "Terraform"},
		{Name: "duration", Type: field.TypeInt, Default: 0},
		{Name: "previous_required_providers", Type: field.TypeJSON},
		{Name: "tags", Type: field.TypeJSON},
		{Name: "instance_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "environment_id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
	}
	// ApplicationRevisionsTable holds the schema information for the "application_revisions" table.
	ApplicationRevisionsTable = &schema.Table{
		Name:       "application_revisions",
		Columns:    ApplicationRevisionsColumns,
		PrimaryKey: []*schema.Column{ApplicationRevisionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "application_revisions_application_instances_revisions",
				Columns:    []*schema.Column{ApplicationRevisionsColumns[14]},
				RefColumns: []*schema.Column{ApplicationInstancesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "application_revisions_environments_revisions",
				Columns:    []*schema.Column{ApplicationRevisionsColumns[15]},
				RefColumns: []*schema.Column{EnvironmentsColumns[0]},
				OnDelete:   schema.Restrict,
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
		{Name: "icon", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON},
		{Name: "source", Type: field.TypeString},
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
	// ModuleVersionsColumns holds the columns for the "module_versions" table.
	ModuleVersionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "version", Type: field.TypeString},
		{Name: "source", Type: field.TypeString},
		{Name: "schema", Type: field.TypeJSON},
		{Name: "module_id", Type: field.TypeString},
	}
	// ModuleVersionsTable holds the schema information for the "module_versions" table.
	ModuleVersionsTable = &schema.Table{
		Name:       "module_versions",
		Columns:    ModuleVersionsColumns,
		PrimaryKey: []*schema.Column{ModuleVersionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "module_versions_modules_versions",
				Columns:    []*schema.Column{ModuleVersionsColumns[6]},
				RefColumns: []*schema.Column{ModulesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "moduleversion_update_time",
				Unique:  false,
				Columns: []*schema.Column{ModuleVersionsColumns[2]},
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
		{Name: "id", Type: field.TypeString, SchemaType: map[string]string{"mysql": "bigint", "postgres": "bigint", "sqlite3": "integer"}},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "domain", Type: field.TypeString, Default: "system"},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
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
		{Name: "group", Type: field.TypeString, Default: "default"},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
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
		AllocationCostsTable,
		ApplicationsTable,
		ApplicationInstancesTable,
		ApplicationModuleRelationshipsTable,
		ApplicationResourcesTable,
		ApplicationRevisionsTable,
		ClusterCostsTable,
		ConnectorsTable,
		EnvironmentsTable,
		EnvironmentConnectorRelationshipsTable,
		ModulesTable,
		ModuleVersionsTable,
		PerspectivesTable,
		ProjectsTable,
		RolesTable,
		SecretsTable,
		SettingsTable,
		SubjectsTable,
		TokensTable,
	}
)

func init() {
	AllocationCostsTable.ForeignKeys[0].RefTable = ConnectorsTable
	ApplicationsTable.ForeignKeys[0].RefTable = ProjectsTable
	ApplicationInstancesTable.ForeignKeys[0].RefTable = ApplicationsTable
	ApplicationInstancesTable.ForeignKeys[1].RefTable = EnvironmentsTable
	ApplicationModuleRelationshipsTable.ForeignKeys[0].RefTable = ApplicationsTable
	ApplicationModuleRelationshipsTable.ForeignKeys[1].RefTable = ModulesTable
	ApplicationResourcesTable.ForeignKeys[0].RefTable = ApplicationInstancesTable
	ApplicationResourcesTable.ForeignKeys[1].RefTable = ApplicationResourcesTable
	ApplicationResourcesTable.ForeignKeys[2].RefTable = ConnectorsTable
	ApplicationRevisionsTable.ForeignKeys[0].RefTable = ApplicationInstancesTable
	ApplicationRevisionsTable.ForeignKeys[1].RefTable = EnvironmentsTable
	ClusterCostsTable.ForeignKeys[0].RefTable = ConnectorsTable
	EnvironmentConnectorRelationshipsTable.ForeignKeys[0].RefTable = EnvironmentsTable
	EnvironmentConnectorRelationshipsTable.ForeignKeys[1].RefTable = ConnectorsTable
	ModuleVersionsTable.ForeignKeys[0].RefTable = ModulesTable
	SecretsTable.ForeignKeys[0].RefTable = ProjectsTable
}
