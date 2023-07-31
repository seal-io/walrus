package cron

import (
	"context"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/distributelock"
	"github.com/seal-io/seal/utils/cron"
	"github.com/seal-io/seal/utils/log"
)

const (
	defaultExpiry = 15 * time.Minute
)

// NewLocker create locker with model client and options.
func NewLocker(client model.ClientSet, options ...Option) *Locker {
	l := &Locker{
		client: client,
		expiry: defaultExpiry,
		logger: log.GetLogger().WithName("locker"),
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

// Option is a function that configures a locker.
type Option func(*Locker)

// Locker implement the cronjob go-co-op/gocron Locker interface.
type Locker struct {
	client model.ClientSet
	expiry time.Duration
	logger log.Logger
}

// Lock try to lock the key provided, return error while failed to lock key.
func (l *Locker) Lock(ctx context.Context, key string) (lock cron.Lock, err error) {
	l.logger.V(6).Infof("try to lock key %s", key)

	defer func() {
		if err != nil && !sqlgraph.IsUniqueConstraintError(err) && !strings.Contains(err.Error(), "is locked") {
			l.logger.Warnf("error lock key %s: %v", key, err)
		}
	}()

	err = l.createLock(ctx, key)
	if err != nil {
		return nil, err
	}

	lock = &Lock{
		key:    key,
		client: l.client,
		logger: l.logger,
	}

	l.logger.V(6).Infof("success lock key %s", key)

	return
}

// createLock create lock with key and store it in the db.
func (l *Locker) createLock(ctx context.Context, key string) error {
	now := time.Now().Unix()
	expiry := now + int64(l.expiry.Seconds())

	return l.client.WithTx(ctx, func(tx *model.Tx) (err error) {
		// Check whether key is existed.
		// Key existed when exception occurred without proper unlocking, like instance been shutdown.
		current, err := tx.DistributeLock.Query().
			Where(distributelock.ID(key)).
			ForUpdate().
			Only(ctx)
		if err != nil && !model.IsNotFound(err) {
			return err
		}

		if current != nil {
			if float64(now-current.ExpireAt) >= l.expiry.Seconds() {
				// Key expired.
				err = tx.DistributeLock.DeleteOneID(key).
					Exec(ctx)
				if err != nil {
					return err
				}
			} else {
				// Key locked, will auto release after expired.
				unlockTime := time.Unix(current.ExpireAt, 0)
				return fmt.Errorf("key %s is locked, release at %s", key, unlockTime.String())
			}
		}

		// Create key.
		_, err = tx.DistributeLock.Create().
			SetID(key).
			SetExpireAt(expiry).
			Save(ctx)
		if err != nil {
			return err
		}

		return
	})
}

// Lock implement the cronjob go-co-op/gocron Lock interface.
type Lock struct {
	key    string
	client model.ClientSet
	logger log.Logger
}

// Unlock release the locked key.
func (l *Lock) Unlock(ctx context.Context) (err error) {
	l.logger.V(6).Infof("try to unlock key %s", l.key)

	defer func() {
		if err != nil {
			l.logger.Warnf("error unlock key %s: %v", l.key, err)
		}
	}()

	err = l.deleteLock(ctx, l.key)
	if err != nil {
		return err
	}

	l.logger.V(6).Infof("success unlock key %s", l.key)

	return nil
}

// deleteLock delete the lock with key.
func (l *Lock) deleteLock(ctx context.Context, key string) error {
	return l.client.WithTx(ctx, func(tx *model.Tx) (err error) {
		// Lock key for delete.
		current, err := tx.DistributeLock.Query().
			Where(distributelock.ID(key)).
			ForUpdate().
			Only(ctx)
		if err != nil && !model.IsNotFound(err) {
			return err
		}

		// Delete key.
		err = tx.DistributeLock.DeleteOne(current).
			Exec(ctx)
		if err != nil {
			return err
		}

		return
	})
}
