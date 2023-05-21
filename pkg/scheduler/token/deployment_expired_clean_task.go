package token

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/types"
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

	_, err := in.modelClient.Tokens().Delete().
		Where(
			token.Kind(types.TokenKindDeployment),
			token.ExpirationNotNil(),
			token.ExpirationLTE(time.Now())).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error clean deployment expired token: %w", err)
	}

	return nil
}
