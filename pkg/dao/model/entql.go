// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/entql"
	"entgo.io/ent/schema/field"
)

// schemaGraph holds a representation of ent/schema at runtime.
var schemaGraph = func() *sqlgraph.Schema {
	graph := &sqlgraph.Schema{Nodes: make([]*sqlgraph.Node, 11)}
	graph.Nodes[0] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   application.Table,
			Columns: application.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: application.FieldID,
			},
		},
		Type: "Application",
		Fields: map[string]*sqlgraph.FieldSpec{
			application.FieldCreateTime:    {Type: field.TypeTime, Column: application.FieldCreateTime},
			application.FieldUpdateTime:    {Type: field.TypeTime, Column: application.FieldUpdateTime},
			application.FieldProjectID:     {Type: field.TypeOther, Column: application.FieldProjectID},
			application.FieldEnvironmentID: {Type: field.TypeOther, Column: application.FieldEnvironmentID},
			application.FieldModules:       {Type: field.TypeJSON, Column: application.FieldModules},
		},
	}
	graph.Nodes[1] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   applicationresource.Table,
			Columns: applicationresource.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: applicationresource.FieldID,
			},
		},
		Type: "ApplicationResource",
		Fields: map[string]*sqlgraph.FieldSpec{
			applicationresource.FieldStatus:        {Type: field.TypeString, Column: applicationresource.FieldStatus},
			applicationresource.FieldStatusMessage: {Type: field.TypeString, Column: applicationresource.FieldStatusMessage},
			applicationresource.FieldCreateTime:    {Type: field.TypeTime, Column: applicationresource.FieldCreateTime},
			applicationresource.FieldUpdateTime:    {Type: field.TypeTime, Column: applicationresource.FieldUpdateTime},
			applicationresource.FieldApplicationID: {Type: field.TypeOther, Column: applicationresource.FieldApplicationID},
			applicationresource.FieldModule:        {Type: field.TypeString, Column: applicationresource.FieldModule},
			applicationresource.FieldType:          {Type: field.TypeString, Column: applicationresource.FieldType},
		},
	}
	graph.Nodes[2] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   applicationrevision.Table,
			Columns: applicationrevision.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: applicationrevision.FieldID,
			},
		},
		Type: "ApplicationRevision",
		Fields: map[string]*sqlgraph.FieldSpec{
			applicationrevision.FieldStatus:         {Type: field.TypeString, Column: applicationrevision.FieldStatus},
			applicationrevision.FieldStatusMessage:  {Type: field.TypeString, Column: applicationrevision.FieldStatusMessage},
			applicationrevision.FieldCreateTime:     {Type: field.TypeTime, Column: applicationrevision.FieldCreateTime},
			applicationrevision.FieldUpdateTime:     {Type: field.TypeTime, Column: applicationrevision.FieldUpdateTime},
			applicationrevision.FieldApplicationID:  {Type: field.TypeOther, Column: applicationrevision.FieldApplicationID},
			applicationrevision.FieldEnvironmentID:  {Type: field.TypeOther, Column: applicationrevision.FieldEnvironmentID},
			applicationrevision.FieldModules:        {Type: field.TypeJSON, Column: applicationrevision.FieldModules},
			applicationrevision.FieldInputVariables: {Type: field.TypeJSON, Column: applicationrevision.FieldInputVariables},
			applicationrevision.FieldInputPlan:      {Type: field.TypeString, Column: applicationrevision.FieldInputPlan},
			applicationrevision.FieldOutput:         {Type: field.TypeString, Column: applicationrevision.FieldOutput},
		},
	}
	graph.Nodes[3] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   connector.Table,
			Columns: connector.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: connector.FieldID,
			},
		},
		Type: "Connector",
		Fields: map[string]*sqlgraph.FieldSpec{
			connector.FieldStatus:        {Type: field.TypeString, Column: connector.FieldStatus},
			connector.FieldStatusMessage: {Type: field.TypeString, Column: connector.FieldStatusMessage},
			connector.FieldCreateTime:    {Type: field.TypeTime, Column: connector.FieldCreateTime},
			connector.FieldUpdateTime:    {Type: field.TypeTime, Column: connector.FieldUpdateTime},
			connector.FieldDriver:        {Type: field.TypeString, Column: connector.FieldDriver},
			connector.FieldConfigVersion: {Type: field.TypeString, Column: connector.FieldConfigVersion},
			connector.FieldConfigData:    {Type: field.TypeJSON, Column: connector.FieldConfigData},
		},
	}
	graph.Nodes[4] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   environment.Table,
			Columns: environment.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: environment.FieldID,
			},
		},
		Type: "Environment",
		Fields: map[string]*sqlgraph.FieldSpec{
			environment.FieldCreateTime:   {Type: field.TypeTime, Column: environment.FieldCreateTime},
			environment.FieldUpdateTime:   {Type: field.TypeTime, Column: environment.FieldUpdateTime},
			environment.FieldConnectorIDs: {Type: field.TypeJSON, Column: environment.FieldConnectorIDs},
			environment.FieldVariables:    {Type: field.TypeJSON, Column: environment.FieldVariables},
		},
	}
	graph.Nodes[5] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   module.Table,
			Columns: module.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: module.FieldID,
			},
		},
		Type: "Module",
		Fields: map[string]*sqlgraph.FieldSpec{
			module.FieldStatus:        {Type: field.TypeString, Column: module.FieldStatus},
			module.FieldStatusMessage: {Type: field.TypeString, Column: module.FieldStatusMessage},
			module.FieldCreateTime:    {Type: field.TypeTime, Column: module.FieldCreateTime},
			module.FieldUpdateTime:    {Type: field.TypeTime, Column: module.FieldUpdateTime},
			module.FieldSource:        {Type: field.TypeString, Column: module.FieldSource},
			module.FieldVersion:       {Type: field.TypeString, Column: module.FieldVersion},
			module.FieldInputSchema:   {Type: field.TypeJSON, Column: module.FieldInputSchema},
			module.FieldOutputSchema:  {Type: field.TypeJSON, Column: module.FieldOutputSchema},
		},
	}
	graph.Nodes[6] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   project.Table,
			Columns: project.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: project.FieldID,
			},
		},
		Type: "Project",
		Fields: map[string]*sqlgraph.FieldSpec{
			project.FieldCreateTime: {Type: field.TypeTime, Column: project.FieldCreateTime},
			project.FieldUpdateTime: {Type: field.TypeTime, Column: project.FieldUpdateTime},
		},
	}
	graph.Nodes[7] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   role.Table,
			Columns: role.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: role.FieldID,
			},
		},
		Type: "Role",
		Fields: map[string]*sqlgraph.FieldSpec{
			role.FieldCreateTime:  {Type: field.TypeTime, Column: role.FieldCreateTime},
			role.FieldUpdateTime:  {Type: field.TypeTime, Column: role.FieldUpdateTime},
			role.FieldDomain:      {Type: field.TypeString, Column: role.FieldDomain},
			role.FieldName:        {Type: field.TypeString, Column: role.FieldName},
			role.FieldDescription: {Type: field.TypeString, Column: role.FieldDescription},
			role.FieldPolicies:    {Type: field.TypeJSON, Column: role.FieldPolicies},
			role.FieldBuiltin:     {Type: field.TypeBool, Column: role.FieldBuiltin},
			role.FieldSession:     {Type: field.TypeBool, Column: role.FieldSession},
		},
	}
	graph.Nodes[8] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   setting.Table,
			Columns: setting.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: setting.FieldID,
			},
		},
		Type: "Setting",
		Fields: map[string]*sqlgraph.FieldSpec{
			setting.FieldCreateTime: {Type: field.TypeTime, Column: setting.FieldCreateTime},
			setting.FieldUpdateTime: {Type: field.TypeTime, Column: setting.FieldUpdateTime},
			setting.FieldName:       {Type: field.TypeString, Column: setting.FieldName},
			setting.FieldValue:      {Type: field.TypeString, Column: setting.FieldValue},
			setting.FieldHidden:     {Type: field.TypeBool, Column: setting.FieldHidden},
			setting.FieldEditable:   {Type: field.TypeBool, Column: setting.FieldEditable},
			setting.FieldPrivate:    {Type: field.TypeBool, Column: setting.FieldPrivate},
		},
	}
	graph.Nodes[9] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   subject.Table,
			Columns: subject.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: subject.FieldID,
			},
		},
		Type: "Subject",
		Fields: map[string]*sqlgraph.FieldSpec{
			subject.FieldCreateTime:  {Type: field.TypeTime, Column: subject.FieldCreateTime},
			subject.FieldUpdateTime:  {Type: field.TypeTime, Column: subject.FieldUpdateTime},
			subject.FieldKind:        {Type: field.TypeString, Column: subject.FieldKind},
			subject.FieldGroup:       {Type: field.TypeString, Column: subject.FieldGroup},
			subject.FieldName:        {Type: field.TypeString, Column: subject.FieldName},
			subject.FieldDescription: {Type: field.TypeString, Column: subject.FieldDescription},
			subject.FieldMountTo:     {Type: field.TypeBool, Column: subject.FieldMountTo},
			subject.FieldLoginTo:     {Type: field.TypeBool, Column: subject.FieldLoginTo},
			subject.FieldRoles:       {Type: field.TypeJSON, Column: subject.FieldRoles},
			subject.FieldPaths:       {Type: field.TypeJSON, Column: subject.FieldPaths},
			subject.FieldBuiltin:     {Type: field.TypeBool, Column: subject.FieldBuiltin},
		},
	}
	graph.Nodes[10] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   token.Table,
			Columns: token.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: token.FieldID,
			},
		},
		Type: "Token",
		Fields: map[string]*sqlgraph.FieldSpec{
			token.FieldCreateTime:        {Type: field.TypeTime, Column: token.FieldCreateTime},
			token.FieldUpdateTime:        {Type: field.TypeTime, Column: token.FieldUpdateTime},
			token.FieldCasdoorTokenName:  {Type: field.TypeString, Column: token.FieldCasdoorTokenName},
			token.FieldCasdoorTokenOwner: {Type: field.TypeString, Column: token.FieldCasdoorTokenOwner},
			token.FieldName:              {Type: field.TypeString, Column: token.FieldName},
			token.FieldExpiration:        {Type: field.TypeInt, Column: token.FieldExpiration},
		},
	}
	return graph
}()

// predicateAdder wraps the addPredicate method.
// All update, update-one and query builders implement this interface.
type predicateAdder interface {
	addPredicate(func(s *sql.Selector))
}

// addPredicate implements the predicateAdder interface.
func (aq *ApplicationQuery) addPredicate(pred func(s *sql.Selector)) {
	aq.predicates = append(aq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the ApplicationQuery builder.
func (aq *ApplicationQuery) Filter() *ApplicationFilter {
	return &ApplicationFilter{config: aq.config, predicateAdder: aq}
}

// addPredicate implements the predicateAdder interface.
func (m *ApplicationMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the ApplicationMutation builder.
func (m *ApplicationMutation) Filter() *ApplicationFilter {
	return &ApplicationFilter{config: m.config, predicateAdder: m}
}

// ApplicationFilter provides a generic filtering capability at runtime for ApplicationQuery.
type ApplicationFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *ApplicationFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[0].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *ApplicationFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(application.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *ApplicationFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(application.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *ApplicationFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(application.FieldUpdateTime))
}

// WhereProjectID applies the entql other predicate on the projectID field.
func (f *ApplicationFilter) WhereProjectID(p entql.OtherP) {
	f.Where(p.Field(application.FieldProjectID))
}

// WhereEnvironmentID applies the entql other predicate on the environmentID field.
func (f *ApplicationFilter) WhereEnvironmentID(p entql.OtherP) {
	f.Where(p.Field(application.FieldEnvironmentID))
}

// WhereModules applies the entql json.RawMessage predicate on the modules field.
func (f *ApplicationFilter) WhereModules(p entql.BytesP) {
	f.Where(p.Field(application.FieldModules))
}

// addPredicate implements the predicateAdder interface.
func (arq *ApplicationResourceQuery) addPredicate(pred func(s *sql.Selector)) {
	arq.predicates = append(arq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the ApplicationResourceQuery builder.
func (arq *ApplicationResourceQuery) Filter() *ApplicationResourceFilter {
	return &ApplicationResourceFilter{config: arq.config, predicateAdder: arq}
}

// addPredicate implements the predicateAdder interface.
func (m *ApplicationResourceMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the ApplicationResourceMutation builder.
func (m *ApplicationResourceMutation) Filter() *ApplicationResourceFilter {
	return &ApplicationResourceFilter{config: m.config, predicateAdder: m}
}

// ApplicationResourceFilter provides a generic filtering capability at runtime for ApplicationResourceQuery.
type ApplicationResourceFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *ApplicationResourceFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[1].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *ApplicationResourceFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(applicationresource.FieldID))
}

// WhereStatus applies the entql string predicate on the status field.
func (f *ApplicationResourceFilter) WhereStatus(p entql.StringP) {
	f.Where(p.Field(applicationresource.FieldStatus))
}

// WhereStatusMessage applies the entql string predicate on the statusMessage field.
func (f *ApplicationResourceFilter) WhereStatusMessage(p entql.StringP) {
	f.Where(p.Field(applicationresource.FieldStatusMessage))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *ApplicationResourceFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(applicationresource.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *ApplicationResourceFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(applicationresource.FieldUpdateTime))
}

// WhereApplicationID applies the entql other predicate on the applicationID field.
func (f *ApplicationResourceFilter) WhereApplicationID(p entql.OtherP) {
	f.Where(p.Field(applicationresource.FieldApplicationID))
}

// WhereModule applies the entql string predicate on the module field.
func (f *ApplicationResourceFilter) WhereModule(p entql.StringP) {
	f.Where(p.Field(applicationresource.FieldModule))
}

// WhereType applies the entql string predicate on the type field.
func (f *ApplicationResourceFilter) WhereType(p entql.StringP) {
	f.Where(p.Field(applicationresource.FieldType))
}

// addPredicate implements the predicateAdder interface.
func (arq *ApplicationRevisionQuery) addPredicate(pred func(s *sql.Selector)) {
	arq.predicates = append(arq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the ApplicationRevisionQuery builder.
func (arq *ApplicationRevisionQuery) Filter() *ApplicationRevisionFilter {
	return &ApplicationRevisionFilter{config: arq.config, predicateAdder: arq}
}

// addPredicate implements the predicateAdder interface.
func (m *ApplicationRevisionMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the ApplicationRevisionMutation builder.
func (m *ApplicationRevisionMutation) Filter() *ApplicationRevisionFilter {
	return &ApplicationRevisionFilter{config: m.config, predicateAdder: m}
}

// ApplicationRevisionFilter provides a generic filtering capability at runtime for ApplicationRevisionQuery.
type ApplicationRevisionFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *ApplicationRevisionFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[2].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *ApplicationRevisionFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(applicationrevision.FieldID))
}

// WhereStatus applies the entql string predicate on the status field.
func (f *ApplicationRevisionFilter) WhereStatus(p entql.StringP) {
	f.Where(p.Field(applicationrevision.FieldStatus))
}

// WhereStatusMessage applies the entql string predicate on the statusMessage field.
func (f *ApplicationRevisionFilter) WhereStatusMessage(p entql.StringP) {
	f.Where(p.Field(applicationrevision.FieldStatusMessage))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *ApplicationRevisionFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(applicationrevision.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *ApplicationRevisionFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(applicationrevision.FieldUpdateTime))
}

// WhereApplicationID applies the entql other predicate on the applicationID field.
func (f *ApplicationRevisionFilter) WhereApplicationID(p entql.OtherP) {
	f.Where(p.Field(applicationrevision.FieldApplicationID))
}

// WhereEnvironmentID applies the entql other predicate on the environmentID field.
func (f *ApplicationRevisionFilter) WhereEnvironmentID(p entql.OtherP) {
	f.Where(p.Field(applicationrevision.FieldEnvironmentID))
}

// WhereModules applies the entql json.RawMessage predicate on the modules field.
func (f *ApplicationRevisionFilter) WhereModules(p entql.BytesP) {
	f.Where(p.Field(applicationrevision.FieldModules))
}

// WhereInputVariables applies the entql json.RawMessage predicate on the inputVariables field.
func (f *ApplicationRevisionFilter) WhereInputVariables(p entql.BytesP) {
	f.Where(p.Field(applicationrevision.FieldInputVariables))
}

// WhereInputPlan applies the entql string predicate on the inputPlan field.
func (f *ApplicationRevisionFilter) WhereInputPlan(p entql.StringP) {
	f.Where(p.Field(applicationrevision.FieldInputPlan))
}

// WhereOutput applies the entql string predicate on the output field.
func (f *ApplicationRevisionFilter) WhereOutput(p entql.StringP) {
	f.Where(p.Field(applicationrevision.FieldOutput))
}

// addPredicate implements the predicateAdder interface.
func (cq *ConnectorQuery) addPredicate(pred func(s *sql.Selector)) {
	cq.predicates = append(cq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the ConnectorQuery builder.
func (cq *ConnectorQuery) Filter() *ConnectorFilter {
	return &ConnectorFilter{config: cq.config, predicateAdder: cq}
}

// addPredicate implements the predicateAdder interface.
func (m *ConnectorMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the ConnectorMutation builder.
func (m *ConnectorMutation) Filter() *ConnectorFilter {
	return &ConnectorFilter{config: m.config, predicateAdder: m}
}

// ConnectorFilter provides a generic filtering capability at runtime for ConnectorQuery.
type ConnectorFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *ConnectorFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[3].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *ConnectorFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(connector.FieldID))
}

// WhereStatus applies the entql string predicate on the status field.
func (f *ConnectorFilter) WhereStatus(p entql.StringP) {
	f.Where(p.Field(connector.FieldStatus))
}

// WhereStatusMessage applies the entql string predicate on the statusMessage field.
func (f *ConnectorFilter) WhereStatusMessage(p entql.StringP) {
	f.Where(p.Field(connector.FieldStatusMessage))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *ConnectorFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(connector.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *ConnectorFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(connector.FieldUpdateTime))
}

// WhereDriver applies the entql string predicate on the driver field.
func (f *ConnectorFilter) WhereDriver(p entql.StringP) {
	f.Where(p.Field(connector.FieldDriver))
}

// WhereConfigVersion applies the entql string predicate on the configVersion field.
func (f *ConnectorFilter) WhereConfigVersion(p entql.StringP) {
	f.Where(p.Field(connector.FieldConfigVersion))
}

// WhereConfigData applies the entql json.RawMessage predicate on the configData field.
func (f *ConnectorFilter) WhereConfigData(p entql.BytesP) {
	f.Where(p.Field(connector.FieldConfigData))
}

// addPredicate implements the predicateAdder interface.
func (eq *EnvironmentQuery) addPredicate(pred func(s *sql.Selector)) {
	eq.predicates = append(eq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the EnvironmentQuery builder.
func (eq *EnvironmentQuery) Filter() *EnvironmentFilter {
	return &EnvironmentFilter{config: eq.config, predicateAdder: eq}
}

// addPredicate implements the predicateAdder interface.
func (m *EnvironmentMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the EnvironmentMutation builder.
func (m *EnvironmentMutation) Filter() *EnvironmentFilter {
	return &EnvironmentFilter{config: m.config, predicateAdder: m}
}

// EnvironmentFilter provides a generic filtering capability at runtime for EnvironmentQuery.
type EnvironmentFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *EnvironmentFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[4].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *EnvironmentFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(environment.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *EnvironmentFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(environment.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *EnvironmentFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(environment.FieldUpdateTime))
}

// WhereConnectorIDs applies the entql json.RawMessage predicate on the connectorIDs field.
func (f *EnvironmentFilter) WhereConnectorIDs(p entql.BytesP) {
	f.Where(p.Field(environment.FieldConnectorIDs))
}

// WhereVariables applies the entql json.RawMessage predicate on the variables field.
func (f *EnvironmentFilter) WhereVariables(p entql.BytesP) {
	f.Where(p.Field(environment.FieldVariables))
}

// addPredicate implements the predicateAdder interface.
func (mq *ModuleQuery) addPredicate(pred func(s *sql.Selector)) {
	mq.predicates = append(mq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the ModuleQuery builder.
func (mq *ModuleQuery) Filter() *ModuleFilter {
	return &ModuleFilter{config: mq.config, predicateAdder: mq}
}

// addPredicate implements the predicateAdder interface.
func (m *ModuleMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the ModuleMutation builder.
func (m *ModuleMutation) Filter() *ModuleFilter {
	return &ModuleFilter{config: m.config, predicateAdder: m}
}

// ModuleFilter provides a generic filtering capability at runtime for ModuleQuery.
type ModuleFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *ModuleFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[5].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *ModuleFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(module.FieldID))
}

// WhereStatus applies the entql string predicate on the status field.
func (f *ModuleFilter) WhereStatus(p entql.StringP) {
	f.Where(p.Field(module.FieldStatus))
}

// WhereStatusMessage applies the entql string predicate on the statusMessage field.
func (f *ModuleFilter) WhereStatusMessage(p entql.StringP) {
	f.Where(p.Field(module.FieldStatusMessage))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *ModuleFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(module.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *ModuleFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(module.FieldUpdateTime))
}

// WhereSource applies the entql string predicate on the source field.
func (f *ModuleFilter) WhereSource(p entql.StringP) {
	f.Where(p.Field(module.FieldSource))
}

// WhereVersion applies the entql string predicate on the version field.
func (f *ModuleFilter) WhereVersion(p entql.StringP) {
	f.Where(p.Field(module.FieldVersion))
}

// WhereInputSchema applies the entql json.RawMessage predicate on the inputSchema field.
func (f *ModuleFilter) WhereInputSchema(p entql.BytesP) {
	f.Where(p.Field(module.FieldInputSchema))
}

// WhereOutputSchema applies the entql json.RawMessage predicate on the outputSchema field.
func (f *ModuleFilter) WhereOutputSchema(p entql.BytesP) {
	f.Where(p.Field(module.FieldOutputSchema))
}

// addPredicate implements the predicateAdder interface.
func (pq *ProjectQuery) addPredicate(pred func(s *sql.Selector)) {
	pq.predicates = append(pq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the ProjectQuery builder.
func (pq *ProjectQuery) Filter() *ProjectFilter {
	return &ProjectFilter{config: pq.config, predicateAdder: pq}
}

// addPredicate implements the predicateAdder interface.
func (m *ProjectMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the ProjectMutation builder.
func (m *ProjectMutation) Filter() *ProjectFilter {
	return &ProjectFilter{config: m.config, predicateAdder: m}
}

// ProjectFilter provides a generic filtering capability at runtime for ProjectQuery.
type ProjectFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *ProjectFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[6].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *ProjectFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(project.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *ProjectFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(project.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *ProjectFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(project.FieldUpdateTime))
}

// addPredicate implements the predicateAdder interface.
func (rq *RoleQuery) addPredicate(pred func(s *sql.Selector)) {
	rq.predicates = append(rq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the RoleQuery builder.
func (rq *RoleQuery) Filter() *RoleFilter {
	return &RoleFilter{config: rq.config, predicateAdder: rq}
}

// addPredicate implements the predicateAdder interface.
func (m *RoleMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the RoleMutation builder.
func (m *RoleMutation) Filter() *RoleFilter {
	return &RoleFilter{config: m.config, predicateAdder: m}
}

// RoleFilter provides a generic filtering capability at runtime for RoleQuery.
type RoleFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *RoleFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[7].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *RoleFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(role.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *RoleFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(role.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *RoleFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(role.FieldUpdateTime))
}

// WhereDomain applies the entql string predicate on the domain field.
func (f *RoleFilter) WhereDomain(p entql.StringP) {
	f.Where(p.Field(role.FieldDomain))
}

// WhereName applies the entql string predicate on the name field.
func (f *RoleFilter) WhereName(p entql.StringP) {
	f.Where(p.Field(role.FieldName))
}

// WhereDescription applies the entql string predicate on the description field.
func (f *RoleFilter) WhereDescription(p entql.StringP) {
	f.Where(p.Field(role.FieldDescription))
}

// WherePolicies applies the entql json.RawMessage predicate on the policies field.
func (f *RoleFilter) WherePolicies(p entql.BytesP) {
	f.Where(p.Field(role.FieldPolicies))
}

// WhereBuiltin applies the entql bool predicate on the builtin field.
func (f *RoleFilter) WhereBuiltin(p entql.BoolP) {
	f.Where(p.Field(role.FieldBuiltin))
}

// WhereSession applies the entql bool predicate on the session field.
func (f *RoleFilter) WhereSession(p entql.BoolP) {
	f.Where(p.Field(role.FieldSession))
}

// addPredicate implements the predicateAdder interface.
func (sq *SettingQuery) addPredicate(pred func(s *sql.Selector)) {
	sq.predicates = append(sq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the SettingQuery builder.
func (sq *SettingQuery) Filter() *SettingFilter {
	return &SettingFilter{config: sq.config, predicateAdder: sq}
}

// addPredicate implements the predicateAdder interface.
func (m *SettingMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the SettingMutation builder.
func (m *SettingMutation) Filter() *SettingFilter {
	return &SettingFilter{config: m.config, predicateAdder: m}
}

// SettingFilter provides a generic filtering capability at runtime for SettingQuery.
type SettingFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *SettingFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[8].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *SettingFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(setting.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *SettingFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(setting.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *SettingFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(setting.FieldUpdateTime))
}

// WhereName applies the entql string predicate on the name field.
func (f *SettingFilter) WhereName(p entql.StringP) {
	f.Where(p.Field(setting.FieldName))
}

// WhereValue applies the entql string predicate on the value field.
func (f *SettingFilter) WhereValue(p entql.StringP) {
	f.Where(p.Field(setting.FieldValue))
}

// WhereHidden applies the entql bool predicate on the hidden field.
func (f *SettingFilter) WhereHidden(p entql.BoolP) {
	f.Where(p.Field(setting.FieldHidden))
}

// WhereEditable applies the entql bool predicate on the editable field.
func (f *SettingFilter) WhereEditable(p entql.BoolP) {
	f.Where(p.Field(setting.FieldEditable))
}

// WherePrivate applies the entql bool predicate on the private field.
func (f *SettingFilter) WherePrivate(p entql.BoolP) {
	f.Where(p.Field(setting.FieldPrivate))
}

// addPredicate implements the predicateAdder interface.
func (sq *SubjectQuery) addPredicate(pred func(s *sql.Selector)) {
	sq.predicates = append(sq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the SubjectQuery builder.
func (sq *SubjectQuery) Filter() *SubjectFilter {
	return &SubjectFilter{config: sq.config, predicateAdder: sq}
}

// addPredicate implements the predicateAdder interface.
func (m *SubjectMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the SubjectMutation builder.
func (m *SubjectMutation) Filter() *SubjectFilter {
	return &SubjectFilter{config: m.config, predicateAdder: m}
}

// SubjectFilter provides a generic filtering capability at runtime for SubjectQuery.
type SubjectFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *SubjectFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[9].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *SubjectFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(subject.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *SubjectFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(subject.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *SubjectFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(subject.FieldUpdateTime))
}

// WhereKind applies the entql string predicate on the kind field.
func (f *SubjectFilter) WhereKind(p entql.StringP) {
	f.Where(p.Field(subject.FieldKind))
}

// WhereGroup applies the entql string predicate on the group field.
func (f *SubjectFilter) WhereGroup(p entql.StringP) {
	f.Where(p.Field(subject.FieldGroup))
}

// WhereName applies the entql string predicate on the name field.
func (f *SubjectFilter) WhereName(p entql.StringP) {
	f.Where(p.Field(subject.FieldName))
}

// WhereDescription applies the entql string predicate on the description field.
func (f *SubjectFilter) WhereDescription(p entql.StringP) {
	f.Where(p.Field(subject.FieldDescription))
}

// WhereMountTo applies the entql bool predicate on the mountTo field.
func (f *SubjectFilter) WhereMountTo(p entql.BoolP) {
	f.Where(p.Field(subject.FieldMountTo))
}

// WhereLoginTo applies the entql bool predicate on the loginTo field.
func (f *SubjectFilter) WhereLoginTo(p entql.BoolP) {
	f.Where(p.Field(subject.FieldLoginTo))
}

// WhereRoles applies the entql json.RawMessage predicate on the roles field.
func (f *SubjectFilter) WhereRoles(p entql.BytesP) {
	f.Where(p.Field(subject.FieldRoles))
}

// WherePaths applies the entql json.RawMessage predicate on the paths field.
func (f *SubjectFilter) WherePaths(p entql.BytesP) {
	f.Where(p.Field(subject.FieldPaths))
}

// WhereBuiltin applies the entql bool predicate on the builtin field.
func (f *SubjectFilter) WhereBuiltin(p entql.BoolP) {
	f.Where(p.Field(subject.FieldBuiltin))
}

// addPredicate implements the predicateAdder interface.
func (tq *TokenQuery) addPredicate(pred func(s *sql.Selector)) {
	tq.predicates = append(tq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the TokenQuery builder.
func (tq *TokenQuery) Filter() *TokenFilter {
	return &TokenFilter{config: tq.config, predicateAdder: tq}
}

// addPredicate implements the predicateAdder interface.
func (m *TokenMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the TokenMutation builder.
func (m *TokenMutation) Filter() *TokenFilter {
	return &TokenFilter{config: m.config, predicateAdder: m}
}

// TokenFilter provides a generic filtering capability at runtime for TokenQuery.
type TokenFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *TokenFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[10].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql other predicate on the id field.
func (f *TokenFilter) WhereID(p entql.OtherP) {
	f.Where(p.Field(token.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the createTime field.
func (f *TokenFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(token.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the updateTime field.
func (f *TokenFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(token.FieldUpdateTime))
}

// WhereCasdoorTokenName applies the entql string predicate on the casdoorTokenName field.
func (f *TokenFilter) WhereCasdoorTokenName(p entql.StringP) {
	f.Where(p.Field(token.FieldCasdoorTokenName))
}

// WhereCasdoorTokenOwner applies the entql string predicate on the casdoorTokenOwner field.
func (f *TokenFilter) WhereCasdoorTokenOwner(p entql.StringP) {
	f.Where(p.Field(token.FieldCasdoorTokenOwner))
}

// WhereName applies the entql string predicate on the name field.
func (f *TokenFilter) WhereName(p entql.StringP) {
	f.Where(p.Field(token.FieldName))
}

// WhereExpiration applies the entql int predicate on the expiration field.
func (f *TokenFilter) WhereExpiration(p entql.IntP) {
	f.Where(p.Field(token.FieldExpiration))
}
