package casdoor

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/req"
)

func Wait(ctx context.Context, serverUrl string) (err error) {
	endpoint.Set(serverUrl)

	var lastErr error

	err = wait.PollImmediateUntilWithContext(ctx, 2*time.Second,
		func(ctx context.Context) (bool, error) {
			lastErr = IsConnected(ctx)
			if lastErr != nil {
				log.Warnf("waiting for casdoor to be ready: %v", lastErr)
			}

			return lastErr == nil, ctx.Err()
		},
	)
	if err != nil && lastErr != nil {
		err = lastErr // Use last error to overwrite context error while existed.
	}

	return
}

func IsConnected(ctx context.Context) error {
	logoutURL := fmt.Sprintf("%s/api/health", endpoint.Get())

	return req.HTTPRequest().
		GetWithContext(ctx, logoutURL).
		Error()
}
