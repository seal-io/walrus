package server

import (
	"context"

	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/topic/datamessage"
)

type Mutation interface {
	ID() (types.ID, bool)
	IDs(context.Context) ([]types.ID, error)
	Tx() (*model.Tx, error)
}

type PublishOptions struct {
	MutationType string
	IDs          []types.ID
	Op           model.Op
	Client       model.ClientSet
}

func (r *Server) initDispatches(ctx context.Context, opts initOptions) error {
	opts.ModelClient.Use(
		func(next model.Mutator) model.Mutator {
			return model.MutateFunc(func(ctx context.Context, m model.Mutation) (model.Value, error) {
				var (
					err error
					ids []types.ID
				)

				hm, ok := m.(Mutation)
				if !ok {
					return next.Mutate(ctx, m)
				}

				if !m.Op().Is(ent.OpCreate) {
					ids, err = hm.IDs(ctx)
					if err != nil {
						return nil, err
					}
				}

				// action before mutate.

				// do mutate.
				value, err := next.Mutate(ctx, m)
				if err != nil {
					return nil, err
				}
				if m.Op() == ent.OpCreate {
					id, ok := hm.ID()
					if ok {
						ids = []types.ID{id}
					}
				}
				publishOpts := PublishOptions{
					MutationType: m.Type(),
					IDs:          ids,
					Op:           m.Op(),
					Client:       opts.ModelClient,
				}

				// action after mutate.
				tx, _ := hm.Tx()
				if tx != nil {
					tx.OnCommit(func(next model.Committer) model.Committer {
						return model.CommitFunc(func(ctx context.Context, tx *model.Tx) error {
							if err = next.Commit(ctx, tx); err != nil {
								return err
							}
							return publish(ctx, publishOpts)
						})
					})
				} else {
					if err = publish(ctx, publishOpts); err != nil {
						return nil, err
					}
				}

				return value, nil
			})
		},
	)

	return nil
}

func publish(ctx context.Context, opts PublishOptions) error {
	// publish application change event when application instance changed.
	if opts.MutationType == string(datamessage.ApplicationInstance) {
		applicationIDs, err := getInstancesApplicationIDs(ctx, opts)
		if err != nil {
			return err
		}
		err = datamessage.Publish(ctx, string(datamessage.Application), model.OpUpdate, applicationIDs)
		if err != nil {
			return err
		}
	}

	return datamessage.Publish(ctx, opts.MutationType, opts.Op, opts.IDs)
}

func getInstancesApplicationIDs(ctx context.Context, opts PublishOptions) ([]types.ID, error) {
	instances, err := opts.Client.ApplicationInstances().Query().
		Select(applicationinstance.FieldApplicationID).
		Where(applicationinstance.IDIn(opts.IDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	applicationIDs := make([]types.ID, 0, len(instances))
	for _, instance := range instances {
		applicationIDs = append(applicationIDs, instance.ApplicationID)
	}

	return applicationIDs, nil
}
