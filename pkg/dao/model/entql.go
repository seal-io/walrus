// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
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
	graph := &sqlgraph.Schema{Nodes: make([]*sqlgraph.Node, 4)}
	graph.Nodes[0] = &sqlgraph.Node{
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
	graph.Nodes[1] = &sqlgraph.Node{
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
	graph.Nodes[2] = &sqlgraph.Node{
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
	graph.Nodes[3] = &sqlgraph.Node{
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
		if err := schemaGraph.EvalP(schemaGraph.Nodes[0].Type, p, s); err != nil {
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
		if err := schemaGraph.EvalP(schemaGraph.Nodes[1].Type, p, s); err != nil {
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
		if err := schemaGraph.EvalP(schemaGraph.Nodes[2].Type, p, s); err != nil {
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
		if err := schemaGraph.EvalP(schemaGraph.Nodes[3].Type, p, s); err != nil {
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
