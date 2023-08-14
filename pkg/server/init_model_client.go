package server

import (
	"context"
	"errors"

	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/log"
)

// configureModelClient enables the mutation hooks or interceptors for the model.Client.
func (r *Server) configureModelClient(ctx context.Context, opts initOptions) error {
	opts.ModelClient.Use(
		dispatchModelChange,
	)

	return nil
}

// dispatchModelChange intercepts almost all DAO writing operations,
// gains the change ID list from intercepting operation,
// then publishes to corresponding topic.
func dispatchModelChange(n model.Mutator) model.Mutator {
	type txer interface {
		Tx() (*model.Tx, error)
	}

	logger := log.WithName("dispatch").WithName("model")

	return model.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		v, err := n.Mutate(ctx, m)
		if err != nil {
			return nil, err
		}

		if !datamessage.IsAllowed(m.Type()) {
			// Return directly if not allowed.
			return v, nil
		}

		// Get ids notifier.
		notify, err := getIdsNotifier(ctx, m)
		if err != nil {
			// NB(thxCode): in order to keep consistency of operating,
			// e.g. delete error is still a write error, not a read error.
			// We only warn out the error to prevent change watching breaking the default behavior.
			logger.Errorf("error getting ids notifier: %v", err)
			return v, nil //nolint: nilerr
		}

		if notify == nil {
			// Return directly if not found.
			return v, nil
		}

		// Wrap the notifier to warn out if error raising.
		//nolint:unparam
		notifyOnly := func() error {
			// NB(thxCode): in order to keep final state of operating,
			// e.g. a deletion is main process, after main process is completed without error,
			// any other branch processes cannot change the main process.
			// We only warn out the error to prevent change watching breaking the final state.
			if err = notify(); err != nil {
				logger.Errorf("error notifying id list: %v", err)
			}

			return nil
		}

		// Notify after committed if processing in transactional session,
		// otherwise, execute immediately.
		t, ok := m.(txer)
		if ok {
			tx, _ := t.Tx()
			if tx != nil {
				tx.OnCommit(func(n model.Committer) model.Committer {
					return model.CommitFunc(func(ctx context.Context, tx *model.Tx) error {
						if err := n.Commit(ctx, tx); err != nil {
							return err
						}

						return notifyOnly()
					})
				})

				return v, nil
			}
		}

		return v, notifyOnly()
	})
}

// getIdsNotifier is a facade function to try the known types one-by-one via getIds function until matching,
// raises an `unknown id type` error after iterated all types.
func getIdsNotifier(ctx context.Context, m model.Mutation) (notify func() error, err error) {
	typ, op := m.Type(), m.Op()

	// Models used object.ID as ID type.
	oids, ok, err := getIds[object.ID](ctx, m)
	if err != nil {
		return
	}

	if ok {
		if len(oids) != 0 {
			notify = func() error { return datamessage.Publish(ctx, typ, op, oids) }
		}

		return
	}

	// Models used string as ID type.
	sids, ok, err := getIds[string](ctx, m)
	if err != nil {
		return
	}

	if ok {
		if len(sids) != 0 {
			notify = func() error { return datamessage.Publish(ctx, typ, op, sids) }
		}

		return
	}

	err = errors.New("unknown id type")

	return
}

type ider[T any] interface {
	ID() (T, bool)
	IDs(context.Context) ([]T, error)
}

func getIds[T any](ctx context.Context, m model.Mutation) (r []T, ok bool, err error) {
	t, ok := m.(ider[T])
	if !ok {
		return
	}

	if !m.Op().Is(ent.OpCreate) {
		// Delete/update ops.
		r, err = t.IDs(ctx)
		return
	}

	if v, exist := t.ID(); exist {
		// Create ops.
		r = []T{v}
	}

	return
}
