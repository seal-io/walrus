package token

import (
	"context"
	"fmt"
	"sync"
	"time"

	tokenbus "github.com/seal-io/seal/pkg/bus/token"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/utils/log"
)

type DeploymentExpiredCleanTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewDeploymentExpiredCleanTask(mc model.ClientSet) (*DeploymentExpiredCleanTask, error) {
	in := &DeploymentExpiredCleanTask{}
	in.modelClient = mc
	in.logger = log.WithName("task").WithName(in.Name())

	return in, nil
}

func (in *DeploymentExpiredCleanTask) Name() string {
	return "token-deployment-expired-clean"
}

func (in *DeploymentExpiredCleanTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}
	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	entities, err := in.modelClient.Tokens().Query().
		Where(
			token.Kind(types.TokenKindDeployment),
			token.ExpirationNotNil(),
			token.ExpirationLTE(time.Now())).
		Select(
			token.FieldID,
			token.FieldValue).
		All(ctx)
	if err != nil {
		return fmt.Errorf("error getting deployment expired token: %w", err)
	}

	if len(entities) == 0 {
		return nil
	}

	ids := make([]object.ID, len(entities))
	for i := range entities {
		ids[i] = entities[i].ID
	}

	_, err = in.modelClient.Tokens().Delete().
		Where(token.IDIn(ids...)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error cleaning deployment expired token: %w", err)
	}

	return tokenbus.Notify(ctx, entities)
}
