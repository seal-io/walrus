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
	entity   *model.Application
	delegate *model.ApplicationCreate
}

func (ac *WrappedApplicationCreate) Save(ctx context.Context) (created *model.Application, err error) {
	var mc = ac.delegate.Mutation().Client()

	// save entity.
	created, err = ac.delegate.Save(ctx)
	if err != nil {
		return
	}

	// construct relationships.
	var newRss = ac.entity.Edges.Modules
	var createRss = make([]*model.ApplicationModuleRelationshipCreate, len(newRss))
	for i, rs := range newRss {
		if rs == nil {
			return nil, errors.New("invalid input: nil relationship")
		}

		// required.
		var c = mc.ApplicationModuleRelationships().Create().
			SetApplicationID(created.ID).
			SetModuleID(rs.ModuleID).
			SetName(rs.Name)

		// optional.
		if rs.Variables != nil {
			c.SetVariables(rs.Variables)
		}
		createRss[i] = c
	}

	// save relationships.
	newRss, err = mc.ApplicationModuleRelationships().CreateBulk(createRss...).
		Save(ctx)
	if err != nil {
		return
	}
	created.Edges.Modules = newRss
	return
}

func (ac *WrappedApplicationCreate) Exec(ctx context.Context) error {
	var _, err = ac.Save(ctx)
	return err
}

func ApplicationCreates(mc model.ClientSet, input ...*model.Application) ([]*WrappedApplicationCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*WrappedApplicationCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// required.
		var c = mc.Applications().Create().
			SetName(r.Name).
			SetProjectID(r.ProjectID).
			SetEnvironmentID(r.EnvironmentID)

		// optional.
		c.SetDescription(r.Description)
		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}
		rrs[i] = &WrappedApplicationCreate{
			entity:   input[i],
			delegate: c,
		}
	}
	return rrs, nil
}

// WrappedApplicationUpdate is a wrapper for model.ApplicationUpdate
// to process the relationship with model.Module.
// TODO(thxCode): generate this with entc.
type WrappedApplicationUpdate struct {
	entity           *model.Application
	entityPredicates []predicate.Application
	delegate         *model.ApplicationUpdate
}

func (au *WrappedApplicationUpdate) Save(ctx context.Context) (updated int, err error) {
	var mc = au.delegate.Mutation().Client()

	if len(au.delegate.Mutation().Fields()) != 0 {
		// update entity.
		updated, err = au.delegate.Save(ctx)
		if err != nil {
			return
		}
	}
	if au.entity.Edges.Modules == nil {
		return
	}

	// get old relationships.
	oldEntity, err := mc.Applications().Query().
		Where(au.entityPredicates...).
		Select(application.FieldID).
		WithModules(func(rq *model.ApplicationModuleRelationshipQuery) {
			rq.Select(
				applicationmodulerelationship.FieldApplicationID,
				applicationmodulerelationship.FieldModuleID,
				applicationmodulerelationship.FieldName,
			)
		}).
		Only(ctx)
	if err != nil {
		return
	}

	// create new relationship or update relationship.
	var applicationID = oldEntity.ID
	var newRsKeys = sets.New[string]()
	var newRss = au.entity.Edges.Modules
	for _, rs := range newRss {
		newRsKeys.Insert(strs.Join("/", string(applicationID), rs.ModuleID, rs.Name))

		// required.
		var c = mc.ApplicationModuleRelationships().Create().
			SetApplicationID(applicationID).
			SetModuleID(rs.ModuleID).
			SetName(rs.Name)

		// optional.
		if rs.Variables != nil {
			c.SetVariables(rs.Variables)
		}

		err = c.OnConflict(
			sql.ConflictColumns(
				applicationmodulerelationship.FieldApplicationID,
				applicationmodulerelationship.FieldModuleID,
				applicationmodulerelationship.FieldName,
			)).
			Update(func(upsert *model.ApplicationModuleRelationshipUpsert) {
				upsert.UpdateVariables()
				upsert.UpdateUpdateTime()
			}).
			Exec(ctx)
		if err != nil {
			return
		}
	}

	// delete stale relationship.
	var oldRss = oldEntity.Edges.Modules
	for _, rs := range oldRss {
		if newRsKeys.Has(strs.Join("/", string(rs.ApplicationID), rs.ModuleID, rs.Name)) {
			continue
		}

		_, err = mc.ApplicationModuleRelationships().Delete().
			Where(
				applicationmodulerelationship.ApplicationID(rs.ApplicationID),
				applicationmodulerelationship.ModuleID(rs.ModuleID),
				applicationmodulerelationship.Name(rs.Name),
			).
			Exec(ctx)
		if err != nil {
			return
		}
	}

	return
}

func (au *WrappedApplicationUpdate) Exec(ctx context.Context) error {
	var _, err = au.Save(ctx)
	return err
}

func ApplicationUpdates(mc model.ClientSet, input ...*model.Application) ([]*WrappedApplicationUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*WrappedApplicationUpdate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// predicated.
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

		// conditional.
		var c = mc.Applications().Update().
			Where(ps...).
			SetDescription(r.Description)
		if r.Name != "" {
			c.SetName(r.Name)
		}
		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}
		rrs[i] = &WrappedApplicationUpdate{
			entity:           input[i],
			entityPredicates: ps,
			delegate:         c,
		}
	}
	return rrs, nil
}
