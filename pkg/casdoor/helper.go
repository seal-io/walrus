package casdoor

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/req"
)

func Wait(ctx context.Context, serverUrl string) error {
	endpoint.Set(serverUrl)
	logoutURL := fmt.Sprintf("%s/api/logout", serverUrl)
	return wait.PollImmediateUntilWithContext(ctx, 2*time.Second,
		func(ctx context.Context) (bool, error) {
			err := req.HTTPRequest().
				Post(logoutURL).
				Error()
			if err != nil {
				log.Warnf("waiting for casdoor to be ready: %v", err)
			}
			return err == nil, ctx.Err()
		},
	)
}
