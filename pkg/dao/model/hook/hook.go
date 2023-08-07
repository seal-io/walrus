// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/dao/model"
)

// The ConnectorFunc type is an adapter to allow the use of ordinary
// function as Connector mutator.
type ConnectorFunc func(context.Context, *model.ConnectorMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f ConnectorFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.ConnectorMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.ConnectorMutation", m)
}

// The CostReportFunc type is an adapter to allow the use of ordinary
// function as CostReport mutator.
type CostReportFunc func(context.Context, *model.CostReportMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f CostReportFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.CostReportMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.CostReportMutation", m)
}

// The DistributeLockFunc type is an adapter to allow the use of ordinary
// function as DistributeLock mutator.
type DistributeLockFunc func(context.Context, *model.DistributeLockMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f DistributeLockFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.DistributeLockMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.DistributeLockMutation", m)
}

// The EnvironmentFunc type is an adapter to allow the use of ordinary
// function as Environment mutator.
type EnvironmentFunc func(context.Context, *model.EnvironmentMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f EnvironmentFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.EnvironmentMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.EnvironmentMutation", m)
}

// The EnvironmentConnectorRelationshipFunc type is an adapter to allow the use of ordinary
// function as EnvironmentConnectorRelationship mutator.
type EnvironmentConnectorRelationshipFunc func(context.Context, *model.EnvironmentConnectorRelationshipMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f EnvironmentConnectorRelationshipFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.EnvironmentConnectorRelationshipMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.EnvironmentConnectorRelationshipMutation", m)
}

// The PerspectiveFunc type is an adapter to allow the use of ordinary
// function as Perspective mutator.
type PerspectiveFunc func(context.Context, *model.PerspectiveMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f PerspectiveFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.PerspectiveMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.PerspectiveMutation", m)
}

// The ProjectFunc type is an adapter to allow the use of ordinary
// function as Project mutator.
type ProjectFunc func(context.Context, *model.ProjectMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f ProjectFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.ProjectMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.ProjectMutation", m)
}

// The RoleFunc type is an adapter to allow the use of ordinary
// function as Role mutator.
type RoleFunc func(context.Context, *model.RoleMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f RoleFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.RoleMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.RoleMutation", m)
}

// The ServiceFunc type is an adapter to allow the use of ordinary
// function as Service mutator.
type ServiceFunc func(context.Context, *model.ServiceMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f ServiceFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.ServiceMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.ServiceMutation", m)
}

// The ServiceRelationshipFunc type is an adapter to allow the use of ordinary
// function as ServiceRelationship mutator.
type ServiceRelationshipFunc func(context.Context, *model.ServiceRelationshipMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f ServiceRelationshipFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.ServiceRelationshipMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.ServiceRelationshipMutation", m)
}

// The ServiceResourceFunc type is an adapter to allow the use of ordinary
// function as ServiceResource mutator.
type ServiceResourceFunc func(context.Context, *model.ServiceResourceMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f ServiceResourceFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.ServiceResourceMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.ServiceResourceMutation", m)
}

// The ServiceResourceRelationshipFunc type is an adapter to allow the use of ordinary
// function as ServiceResourceRelationship mutator.
type ServiceResourceRelationshipFunc func(context.Context, *model.ServiceResourceRelationshipMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f ServiceResourceRelationshipFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.ServiceResourceRelationshipMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.ServiceResourceRelationshipMutation", m)
}

// The ServiceRevisionFunc type is an adapter to allow the use of ordinary
// function as ServiceRevision mutator.
type ServiceRevisionFunc func(context.Context, *model.ServiceRevisionMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f ServiceRevisionFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.ServiceRevisionMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.ServiceRevisionMutation", m)
}

// The SettingFunc type is an adapter to allow the use of ordinary
// function as Setting mutator.
type SettingFunc func(context.Context, *model.SettingMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f SettingFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.SettingMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.SettingMutation", m)
}

// The SubjectFunc type is an adapter to allow the use of ordinary
// function as Subject mutator.
type SubjectFunc func(context.Context, *model.SubjectMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f SubjectFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.SubjectMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.SubjectMutation", m)
}

// The SubjectRoleRelationshipFunc type is an adapter to allow the use of ordinary
// function as SubjectRoleRelationship mutator.
type SubjectRoleRelationshipFunc func(context.Context, *model.SubjectRoleRelationshipMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f SubjectRoleRelationshipFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.SubjectRoleRelationshipMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.SubjectRoleRelationshipMutation", m)
}

// The TemplateFunc type is an adapter to allow the use of ordinary
// function as Template mutator.
type TemplateFunc func(context.Context, *model.TemplateMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f TemplateFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.TemplateMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.TemplateMutation", m)
}

// The TemplateVersionFunc type is an adapter to allow the use of ordinary
// function as TemplateVersion mutator.
type TemplateVersionFunc func(context.Context, *model.TemplateVersionMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f TemplateVersionFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.TemplateVersionMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.TemplateVersionMutation", m)
}

// The TokenFunc type is an adapter to allow the use of ordinary
// function as Token mutator.
type TokenFunc func(context.Context, *model.TokenMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f TokenFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.TokenMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.TokenMutation", m)
}

// The VariableFunc type is an adapter to allow the use of ordinary
// function as Variable mutator.
type VariableFunc func(context.Context, *model.VariableMutation) (model.Value, error)

// Mutate calls f(ctx, m).
func (f VariableFunc) Mutate(ctx context.Context, m model.Mutation) (model.Value, error) {
	if mv, ok := m.(*model.VariableMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *model.VariableMutation", m)
}

// Condition is a hook condition function.
type Condition func(context.Context, model.Mutation) bool

// And groups conditions with the AND operator.
func And(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m model.Mutation) bool {
		if !first(ctx, m) || !second(ctx, m) {
			return false
		}
		for _, cond := range rest {
			if !cond(ctx, m) {
				return false
			}
		}
		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m model.Mutation) bool {
		if first(ctx, m) || second(ctx, m) {
			return true
		}
		for _, cond := range rest {
			if cond(ctx, m) {
				return true
			}
		}
		return false
	}
}

// Not negates a given condition.
func Not(cond Condition) Condition {
	return func(ctx context.Context, m model.Mutation) bool {
		return !cond(ctx, m)
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op model.Op) Condition {
	return func(_ context.Context, m model.Mutation) bool {
		return m.Op().Is(op)
	}
}

// HasAddedFields is a condition validating `.AddedField` on fields.
func HasAddedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m model.Mutation) bool {
		if _, exists := m.AddedField(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.AddedField(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasClearedFields is a condition validating `.FieldCleared` on fields.
func HasClearedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m model.Mutation) bool {
		if exists := m.FieldCleared(field); !exists {
			return false
		}
		for _, field := range fields {
			if exists := m.FieldCleared(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasFields is a condition validating `.Field` on fields.
func HasFields(field string, fields ...string) Condition {
	return func(_ context.Context, m model.Mutation) bool {
		if _, exists := m.Field(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.Field(field); !exists {
				return false
			}
		}
		return true
	}
}

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
func If(hk model.Hook, cond Condition) model.Hook {
	return func(next model.Mutator) model.Mutator {
		return model.MutateFunc(func(ctx context.Context, m model.Mutation) (model.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// On executes the given hook only for the given operation.
//
//	hook.On(Log, model.Delete|model.Create)
func On(hk model.Hook, op model.Op) model.Hook {
	return If(hk, HasOp(op))
}

// Unless skips the given hook only for the given operation.
//
//	hook.Unless(Log, model.Update|model.UpdateOne)
func Unless(hk model.Hook, op model.Op) model.Hook {
	return If(hk, Not(HasOp(op)))
}

// FixedError is a hook returning a fixed error.
func FixedError(err error) model.Hook {
	return func(model.Mutator) model.Mutator {
		return model.MutateFunc(func(context.Context, model.Mutation) (model.Value, error) {
			return nil, err
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []model.Hook {
//		return []model.Hook{
//			Reject(model.Delete|model.Update),
//		}
//	}
func Reject(op model.Op) model.Hook {
	hk := FixedError(fmt.Errorf("%s operation is not allowed", op))
	return On(hk, op)
}

// Chain acts as a list of hooks and is effectively immutable.
// Once created, it will always hold the same set of hooks in the same order.
type Chain struct {
	hooks []model.Hook
}

// NewChain creates a new chain of hooks.
func NewChain(hooks ...model.Hook) Chain {
	return Chain{append([]model.Hook(nil), hooks...)}
}

// Hook chains the list of hooks and returns the final hook.
func (c Chain) Hook() model.Hook {
	return func(mutator model.Mutator) model.Mutator {
		for i := len(c.hooks) - 1; i >= 0; i-- {
			mutator = c.hooks[i](mutator)
		}
		return mutator
	}
}

// Append extends a chain, adding the specified hook
// as the last ones in the mutation flow.
func (c Chain) Append(hooks ...model.Hook) Chain {
	newHooks := make([]model.Hook, 0, len(c.hooks)+len(hooks))
	newHooks = append(newHooks, c.hooks...)
	newHooks = append(newHooks, hooks...)
	return Chain{newHooks}
}

// Extend extends a chain, adding the specified chain
// as the last ones in the mutation flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.hooks...)
}
