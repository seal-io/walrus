package server

import (
	"context"

	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/topic/datamessage"
)

type Mutation interface {
	IDs(ctx context.Context) ([]types.ID, error)
	Tx() (*model.Tx, error)
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

				// action after mutate.
				tx, _ := hm.Tx()
				if tx != nil {
					tx.OnCommit(func(next model.Committer) model.Committer {
						return model.CommitFunc(func(ctx context.Context, tx *model.Tx) error {
							if err = next.Commit(ctx, tx); err != nil {
								return err
							}
							return datamessage.Publish(ctx, m.Type(), m.Op(), ids)
						})
					})
				} else {
					if err = datamessage.Publish(ctx, m.Type(), m.Op(), ids); err != nil {
						return nil, err
					}
				}

				return value, nil
			})
		},
	)

	return nil
}
