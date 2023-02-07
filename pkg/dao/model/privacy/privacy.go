// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package privacy

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/dao/model"

	"entgo.io/ent/entql"
	"entgo.io/ent/privacy"
)

var (
	// Allow may be returned by rules to indicate that the policy
	// evaluation should terminate with allow decision.
	Allow = privacy.Allow

	// Deny may be returned by rules to indicate that the policy
	// evaluation should terminate with deny decision.
	Deny = privacy.Deny

	// Skip may be returned by rules to indicate that the policy
	// evaluation should continue to the next rule.
	Skip = privacy.Skip
)

// Allowf returns an formatted wrapped Allow decision.
func Allowf(format string, a ...any) error {
	return fmt.Errorf(format+": %w", append(a, Allow)...)
}

// Denyf returns an formatted wrapped Deny decision.
func Denyf(format string, a ...any) error {
	return fmt.Errorf(format+": %w", append(a, Deny)...)
}

// Skipf returns an formatted wrapped Skip decision.
func Skipf(format string, a ...any) error {
	return fmt.Errorf(format+": %w", append(a, Skip)...)
}

// DecisionContext creates a new context from the given parent context with
// a policy decision attach to it.
func DecisionContext(parent context.Context, decision error) context.Context {
	return privacy.DecisionContext(parent, decision)
}

// DecisionFromContext retrieves the policy decision from the context.
func DecisionFromContext(ctx context.Context) (error, bool) {
	return privacy.DecisionFromContext(ctx)
}

type (
	// Policy groups query and mutation policies.
	Policy = privacy.Policy

	// QueryRule defines the interface deciding whether a
	// query is allowed and optionally modify it.
	QueryRule = privacy.QueryRule
	// QueryPolicy combines multiple query rules into a single policy.
	QueryPolicy = privacy.QueryPolicy

	// MutationRule defines the interface which decides whether a
	// mutation is allowed and optionally modifies it.
	MutationRule = privacy.MutationRule
	// MutationPolicy combines multiple mutation rules into a single policy.
	MutationPolicy = privacy.MutationPolicy
)

// QueryRuleFunc type is an adapter to allow the use of
// ordinary functions as query rules.
type QueryRuleFunc func(context.Context, model.Query) error

// Eval returns f(ctx, q).
func (f QueryRuleFunc) EvalQuery(ctx context.Context, q model.Query) error {
	return f(ctx, q)
}

// MutationRuleFunc type is an adapter which allows the use of
// ordinary functions as mutation rules.
type MutationRuleFunc func(context.Context, model.Mutation) error

// EvalMutation returns f(ctx, m).
func (f MutationRuleFunc) EvalMutation(ctx context.Context, m model.Mutation) error {
	return f(ctx, m)
}

// QueryMutationRule is an interface which groups query and mutation rules.
type QueryMutationRule interface {
	QueryRule
	MutationRule
}

// AlwaysAllowRule returns a rule that returns an allow decision.
func AlwaysAllowRule() QueryMutationRule {
	return fixedDecision{Allow}
}

// AlwaysDenyRule returns a rule that returns a deny decision.
func AlwaysDenyRule() QueryMutationRule {
	return fixedDecision{Deny}
}

type fixedDecision struct {
	decision error
}

func (f fixedDecision) EvalQuery(context.Context, model.Query) error {
	return f.decision
}

func (f fixedDecision) EvalMutation(context.Context, model.Mutation) error {
	return f.decision
}

type contextDecision struct {
	eval func(context.Context) error
}

// ContextQueryMutationRule creates a query/mutation rule from a context eval func.
func ContextQueryMutationRule(eval func(context.Context) error) QueryMutationRule {
	return contextDecision{eval}
}

func (c contextDecision) EvalQuery(ctx context.Context, _ model.Query) error {
	return c.eval(ctx)
}

func (c contextDecision) EvalMutation(ctx context.Context, _ model.Mutation) error {
	return c.eval(ctx)
}

// OnMutationOperation evaluates the given rule only on a given mutation operation.
func OnMutationOperation(rule MutationRule, op model.Op) MutationRule {
	return MutationRuleFunc(func(ctx context.Context, m model.Mutation) error {
		if m.Op().Is(op) {
			return rule.EvalMutation(ctx, m)
		}
		return Skip
	})
}

// DenyMutationOperationRule returns a rule denying specified mutation operation.
func DenyMutationOperationRule(op model.Op) MutationRule {
	rule := MutationRuleFunc(func(_ context.Context, m model.Mutation) error {
		return Denyf("model/privacy: operation %s is not allowed", m.Op())
	})
	return OnMutationOperation(rule, op)
}

// The RoleQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type RoleQueryRuleFunc func(context.Context, *model.RoleQuery) error

// EvalQuery return f(ctx, q).
func (f RoleQueryRuleFunc) EvalQuery(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.RoleQuery); ok {
		return f(ctx, q)
	}
	return Denyf("model/privacy: unexpected query type %T, expect *model.RoleQuery", q)
}

// The RoleMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type RoleMutationRuleFunc func(context.Context, *model.RoleMutation) error

// EvalMutation calls f(ctx, m).
func (f RoleMutationRuleFunc) EvalMutation(ctx context.Context, m model.Mutation) error {
	if m, ok := m.(*model.RoleMutation); ok {
		return f(ctx, m)
	}
	return Denyf("model/privacy: unexpected mutation type %T, expect *model.RoleMutation", m)
}

// The SettingQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SettingQueryRuleFunc func(context.Context, *model.SettingQuery) error

// EvalQuery return f(ctx, q).
func (f SettingQueryRuleFunc) EvalQuery(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.SettingQuery); ok {
		return f(ctx, q)
	}
	return Denyf("model/privacy: unexpected query type %T, expect *model.SettingQuery", q)
}

// The SettingMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SettingMutationRuleFunc func(context.Context, *model.SettingMutation) error

// EvalMutation calls f(ctx, m).
func (f SettingMutationRuleFunc) EvalMutation(ctx context.Context, m model.Mutation) error {
	if m, ok := m.(*model.SettingMutation); ok {
		return f(ctx, m)
	}
	return Denyf("model/privacy: unexpected mutation type %T, expect *model.SettingMutation", m)
}

// The SubjectQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SubjectQueryRuleFunc func(context.Context, *model.SubjectQuery) error

// EvalQuery return f(ctx, q).
func (f SubjectQueryRuleFunc) EvalQuery(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.SubjectQuery); ok {
		return f(ctx, q)
	}
	return Denyf("model/privacy: unexpected query type %T, expect *model.SubjectQuery", q)
}

// The SubjectMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SubjectMutationRuleFunc func(context.Context, *model.SubjectMutation) error

// EvalMutation calls f(ctx, m).
func (f SubjectMutationRuleFunc) EvalMutation(ctx context.Context, m model.Mutation) error {
	if m, ok := m.(*model.SubjectMutation); ok {
		return f(ctx, m)
	}
	return Denyf("model/privacy: unexpected mutation type %T, expect *model.SubjectMutation", m)
}

// The TokenQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type TokenQueryRuleFunc func(context.Context, *model.TokenQuery) error

// EvalQuery return f(ctx, q).
func (f TokenQueryRuleFunc) EvalQuery(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.TokenQuery); ok {
		return f(ctx, q)
	}
	return Denyf("model/privacy: unexpected query type %T, expect *model.TokenQuery", q)
}

// The TokenMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type TokenMutationRuleFunc func(context.Context, *model.TokenMutation) error

// EvalMutation calls f(ctx, m).
func (f TokenMutationRuleFunc) EvalMutation(ctx context.Context, m model.Mutation) error {
	if m, ok := m.(*model.TokenMutation); ok {
		return f(ctx, m)
	}
	return Denyf("model/privacy: unexpected mutation type %T, expect *model.TokenMutation", m)
}

type (
	// Filter is the interface that wraps the Where function
	// for filtering nodes in queries and mutations.
	Filter interface {
		// Where applies a filter on the executed query/mutation.
		Where(entql.P)
	}

	// The FilterFunc type is an adapter that allows the use of ordinary
	// functions as filters for query and mutation types.
	FilterFunc func(context.Context, Filter) error
)

// EvalQuery calls f(ctx, q) if the query implements the Filter interface, otherwise it is denied.
func (f FilterFunc) EvalQuery(ctx context.Context, q model.Query) error {
	fr, err := queryFilter(q)
	if err != nil {
		return err
	}
	return f(ctx, fr)
}

// EvalMutation calls f(ctx, q) if the mutation implements the Filter interface, otherwise it is denied.
func (f FilterFunc) EvalMutation(ctx context.Context, m model.Mutation) error {
	fr, err := mutationFilter(m)
	if err != nil {
		return err
	}
	return f(ctx, fr)
}

var _ QueryMutationRule = FilterFunc(nil)

func queryFilter(q model.Query) (Filter, error) {
	switch q := q.(type) {
	case *model.RoleQuery:
		return q.Filter(), nil
	case *model.SettingQuery:
		return q.Filter(), nil
	case *model.SubjectQuery:
		return q.Filter(), nil
	case *model.TokenQuery:
		return q.Filter(), nil
	default:
		return nil, Denyf("model/privacy: unexpected query type %T for query filter", q)
	}
}

func mutationFilter(m model.Mutation) (Filter, error) {
	switch m := m.(type) {
	case *model.RoleMutation:
		return m.Filter(), nil
	case *model.SettingMutation:
		return m.Filter(), nil
	case *model.SubjectMutation:
		return m.Filter(), nil
	case *model.TokenMutation:
		return m.Filter(), nil
	default:
		return nil, Denyf("model/privacy: unexpected mutation type %T for mutation filter", m)
	}
}
