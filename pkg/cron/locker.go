package cron

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/distributelock"
	"github.com/seal-io/walrus/utils/cron"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
	"github.com/seal-io/walrus/utils/version"
)

const (
	defaultExpiry        = 15 * time.Second
	defaultRenewInterval = 5 * time.Second
	maxExecuteDuration   = 15 * time.Minute
)

// NewLocker create locker with model client and options.
func NewLocker(logger log.Logger, client *model.Client, options ...Option) *Locker {
	l := &Locker{
		logger:             logger.WithName("locker"),
		holder:             version.GetInstanceUUID(),
		client:             client,
		expiry:             defaultExpiry,
		renewInterval:      defaultRenewInterval,
		maxExecuteDuration: maxExecuteDuration,
	}
	for _, opt := range options {
		opt(l)
	}

	return l
}

// WithExpiryInterval can be used to set the expiry of a locker to set to every key.
func WithExpiryInterval(expiry time.Duration) Option {
	return func(s *Locker) {
		s.expiry = expiry
	}
}

// WithRenewConfig can be used to set renew interval and max execute duration of a locker.
func WithRenewConfig(renewInterval, maxExecuteDuration time.Duration) Option {
	return func(s *Locker) {
		s.renewInterval = renewInterval
		s.maxExecuteDuration = maxExecuteDuration
	}
}

// Option is a function that configures a locker.
type Option func(*Locker)

// Locker implement the cronjob go-co-op/gocron Locker interface.
type Locker struct {
	logger             log.Logger
	holder             string
	client             *model.Client
	expiry             time.Duration
	renewInterval      time.Duration
	maxExecuteDuration time.Duration
}

// Lock try to lock the key provided, return error while failed to lock key.
func (l *Locker) Lock(ctx context.Context, key string) (cron.Lock, error) {
	lk, err := l.acquire(ctx, key)
	if err != nil {
		var errLocked *lockedError

		switch {
		case !sqlgraph.IsUniqueConstraintError(err) && !errors.As(err, &errLocked):
			l.logger.Warnf("error acquiring lock for %q: %v", key, err)
		case errLocked != nil && errLocked.Holder == l.holder:
			// NB(thxCode): time-cost operation, need to be optimized.
			l.logger.WarnS("previous locking is not finished", "createdBy", key)
		}

		return nil, err
	}

	lk.logger.V(6).Info("succeeded lock")

	return lk, nil
}

type lockedError struct {
	Key, Holder string
}

func (e *lockedError) Error() string {
	return fmt.Sprintf("key %q is locked by %q", e.Key, e.Holder)
}

// acquire try to get the lock of key provided,
// return error while failed to acquire lock.
func (l *Locker) acquire(ctx context.Context, key string) (*Lock, error) {
	now := time.Now().Unix()
	expiry := now + int64(l.expiry.Seconds())

	err := l.client.WithTx(ctx, func(tx *model.Tx) error {
		// Check whether key is existed.
		// Key may fail to unlock when exception occurred, like instance been shutdown.
		current, err := tx.DistributeLock.Query().
			Where(distributelock.ID(key)).
			ForUpdate().
			Only(ctx)
		if err != nil && !model.IsNotFound(err) {
			return err
		}

		if current != nil {
			// Key locked, will auto release after expired.
			if float64(now-current.ExpireAt) < l.expiry.Seconds() {
				return &lockedError{
					Key:    key,
					Holder: current.Holder,
				}
			}

			// Key expired, release it.
			err = tx.DistributeLock.DeleteOneID(key).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		// Create key.
		return tx.DistributeLock.Create().
			SetID(key).
			SetExpireAt(expiry).
			SetHolder(l.holder).
			Exec(ctx)
	})
	if err != nil {
		return nil, err
	}

	lk := &Lock{
		logger: l.logger.WithValues("key", key, "id", strs.Hex(6)),
		holder: l.holder,
		client: l.client,
		expiry: l.expiry,
		key:    key,
	}

	// Start renew.
	renewCtx, renewCancel := context.WithCancel(ctx)
	lk.renewCancel = renewCancel
	lk.renew(renewCtx, l.renewInterval, l.maxExecuteDuration)

	return lk, nil
}

// Lock implement the cronjob go-co-op/gocron Lock interface.
type Lock struct {
	logger log.Logger
	holder string
	client model.ClientSet
	expiry time.Duration
	key    string

	renewCancel context.CancelFunc
}

// renew keeps alive the lock.
func (lk *Lock) renew(ctx context.Context, interval, timeout time.Duration) {
	gopool.Go(func() {
		_ = wait.PollUntilContextTimeout(ctx, interval, timeout, true, lk.doRenew)
		lk.logger.V(6).Info("stopped renew")
	})
}

// doRenew returns true if not found key,
// otherwise, keeps running until context canceled.
func (lk *Lock) doRenew(ctx context.Context) (done bool, err error) {
	err = lk.client.WithTx(ctx, func(tx *model.Tx) error {
		// Get key.
		key, err := tx.DistributeLock.Query().
			Where(
				distributelock.ID(lk.key),
				distributelock.Holder(lk.holder)).
			Only(ctx)
		if err != nil {
			return err
		}

		if key.ExpireAt < time.Now().Unix() {
			// Should not renew.
			return nil
		}

		// Renew key.
		expiry := time.Now().Unix() + int64(lk.expiry.Seconds())

		err = tx.DistributeLock.UpdateOne(key).
			SetExpireAt(expiry).
			Exec(ctx)
		if err != nil {
			return err
		}

		lk.logger.V(6).Info("succeeded renew")

		return nil
	})
	if err != nil && ctx.Err() == nil {
		if model.IsNotFound(err) {
			// Not found, should not renew again.
			return true, nil
		}

		lk.logger.Warnf("error renewing, try at next loop: %v", err)
	}

	return false, ctx.Err()
}

// Unlock release the lock.
func (lk *Lock) Unlock(ctx context.Context) error {
	// Stop renew.
	defer lk.renewCancel()

	// Delete key.
	if ctx.Err() != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
	}

	_, err := lk.client.DistributeLocks().Delete().
		Where(
			distributelock.ID(lk.key),
			distributelock.Holder(lk.holder)).
		Exec(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		lk.logger.Warnf("error unlocking: %v", err)
		return err
	}

	lk.logger.V(6).Info("succeeded unlock")

	return nil
}
