package token

import (
	"context"
	"fmt"
	"time"

	tokenbus "github.com/seal-io/walrus/pkg/bus/token"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/token"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/log"
)

type DeploymentExpiredCleanTask struct {
	logger      log.Logger
	modelClient model.ClientSet
}

func NewDeploymentExpiredCleanTask(mc model.ClientSet) (in *DeploymentExpiredCleanTask, err error) {
	in = &DeploymentExpiredCleanTask{
		logger:      log.WithName("task").WithName(in.Name()),
		modelClient: mc,
	}

	return
}

func (in *DeploymentExpiredCleanTask) Name() string {
	return "token-deployment-expired-clean"
}

func (in *DeploymentExpiredCleanTask) Process(ctx context.Context, args ...any) error {
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
