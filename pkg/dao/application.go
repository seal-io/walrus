package dao

import (
	"context"
	"errors"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/utils/strs"
)

// WrappedApplicationCreate is a wrapper for model.ApplicationCreate
// to process the relationship with model.Module.
// TODO(thxCode): generate this with entc.
type WrappedApplicationCreate struct {
	*model.ApplicationCreate

	entity *model.Application
}

func (ac *WrappedApplicationCreate) Save(ctx context.Context) (created *model.Application, err error) {
	mc := ac.ApplicationCreate.Mutation().Client()

	// Save entity.
	created, err = ac.ApplicationCreate.Save(ctx)
	if err != nil {
		return
	}

	// Construct relationships.
	newRss := ac.entity.Edges.Modules
	createRss := make([]*model.ApplicationModuleRelationshipCreate, len(newRss))

	for i, rs := range newRss {
		if rs == nil {
			return nil, errors.New("invalid input: nil relationship")
		}

		// Required.
		c := mc.ApplicationModuleRelationships().Create().
			SetApplicationID(created.ID).
			SetModuleID(rs.ModuleID).
			SetVersion(rs.Version).
			SetName(rs.Name)

		// Optional.
		if rs.Attributes != nil {
			c.SetAttributes(rs.Attributes)
		}
		createRss[i] = c
	}

	// Save relationships.
	newRss, err = mc.ApplicationModuleRelationships().CreateBulk(createRss...).
		Save(ctx)
	if err != nil {
		return
	}
	created.Edges.Modules = newRss

	return created, nil
}

func (ac *WrappedApplicationCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

func ApplicationCreates(mc model.ClientSet, input ...*model.Application) ([]*WrappedApplicationCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*WrappedApplicationCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Applications().Create().
			SetName(r.Name).
			SetProjectID(r.ProjectID)

		// Optional.
		c.SetDescription(r.Description)

		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}

		if r.Variables != nil {
			c.SetVariables(r.Variables)
		}
		rrs[i] = &WrappedApplicationCreate{
			ApplicationCreate: c,
			entity:            input[i],
		}
	}

	return rrs, nil
}

// WrappedApplicationUpdate is a wrapper for model.ApplicationUpdate
// to process the relationship with model.Module.
// TODO(thxCode): generate this with entc.
type WrappedApplicationUpdate struct {
	*model.ApplicationUpdate

	entity           *model.Application
	entityPredicates []predicate.Application
}

func (au *WrappedApplicationUpdate) Save(ctx context.Context) (updated int, err error) {
	mc := au.ApplicationUpdate.Mutation().Client()

	if len(au.ApplicationUpdate.Mutation().Fields()) != 0 {
		// Update entity.
		updated, err = au.ApplicationUpdate.Save(ctx)
		if err != nil {
			return
		}
	}

	// Get old relationships.
	oldEntity, err := mc.Applications().Query().
		Where(au.entityPredicates...).
		Select(application.FieldID).
		WithModules(func(rq *model.ApplicationModuleRelationshipQuery) {
			rq.Select(
				applicationmodulerelationship.FieldApplicationID,
				applicationmodulerelationship.FieldModuleID,
				applicationmodulerelationship.FieldVersion,
				applicationmodulerelationship.FieldName,
			)
		}).
		Only(ctx)
	if err != nil {
		return
	}

	// Create new relationship or update relationship.
	applicationID := oldEntity.ID
	newRsKeys := sets.New[string]()
	newRss := au.entity.Edges.Modules

	for _, rs := range newRss {
		newRsKeys.Insert(strs.Join("/", string(applicationID), rs.ModuleID, rs.Name))

		// Required.
		c := mc.ApplicationModuleRelationships().Create().
			SetApplicationID(applicationID).
			SetModuleID(rs.ModuleID).
			SetVersion(rs.Version).
			SetName(rs.Name)

		// Optional.
		if rs.Attributes != nil {
			c.SetAttributes(rs.Attributes)
		}

		err = c.OnConflict(
			sql.ConflictColumns(
				applicationmodulerelationship.FieldApplicationID,
				applicationmodulerelationship.FieldModuleID,
				applicationmodulerelationship.FieldName,
			)).
			Update(func(upsert *model.ApplicationModuleRelationshipUpsert) {
				if rs.Attributes != nil {
					upsert.UpdateAttributes()
				}
				upsert.UpdateVersion()
				upsert.UpdateUpdateTime()
			}).
			Exec(ctx)
		if err != nil {
			return
		}
	}

	// Delete stale relationship.
	oldRss := oldEntity.Edges.Modules
	for _, rs := range oldRss {
		if newRsKeys.Has(strs.Join("/", string(rs.ApplicationID), rs.ModuleID, rs.Name)) {
			continue
		}

		_, err = mc.ApplicationModuleRelationships().Delete().
			Where(
				applicationmodulerelationship.ApplicationID(rs.ApplicationID),
				applicationmodulerelationship.Name(rs.Name),
			).
			Exec(ctx)
		if err != nil {
			return
		}
	}

	return updated, nil
}

func (au *WrappedApplicationUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

func ApplicationUpdates(mc model.ClientSet, input ...*model.Application) ([]*WrappedApplicationUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*WrappedApplicationUpdate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Predicated.
		var ps []predicate.Application

		switch {
		case r.ID.IsNaive():
			ps = append(ps, application.ID(r.ID))
		case r.ProjectID != "" && r.Name != "":
			ps = append(ps, application.And(
				application.ProjectID(r.ProjectID),
				application.Name(r.Name),
			))
		}

		if len(ps) == 0 {
			return nil, errors.New("invalid input: illegal predicates")
		}

		// Conditional.
		c := mc.Applications().Update().
			Where(ps...).
			SetDescription(r.Description)
		if r.Name != "" {
			c.SetName(r.Name)
		}

		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}

		if r.Variables != nil {
			c.SetVariables(r.Variables)
		}
		rrs[i] = &WrappedApplicationUpdate{
			ApplicationUpdate: c,
			entity:            input[i],
			entityPredicates:  ps,
		}
	}

	return rrs, nil
}
