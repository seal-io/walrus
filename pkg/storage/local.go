package storage

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/walrus/pkg/consts"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/req"
	"github.com/seal-io/walrus/utils/vars"
)

const DefaultS3SourceAddress = "s3://minio:minio123@localhost:9000/walrus?sslmode=disable"

var endpoint = &vars.SetOnce[string]{}

const (
	defaultEmbeddedMinioEndpointAddress = "localhost:9000"
	defaultEmbeddedMinioUser            = "minio"
	defaultEmbeddedMinioPassword        = "minio123"
)

type Embedded struct{}

func (Embedded) Run(ctx context.Context) error {
	// Create run data dir if not found.
	runDataPath := filepath.Join(consts.DataDir, "minio")

	const cmdName = "minio"
	logger := log.WithName(cmdName)
	cmdAgs := []string{
		"server",
		runDataPath,
	}

	envs := []string{
		"MINIO_ROOT_USER=" + defaultEmbeddedMinioUser,
		"MINIO_ROOT_PASSWORD=" + defaultEmbeddedMinioPassword,
	}

	logger.Infof("run: %s %s", cmdName, cmdAgs)
	cmd := exec.CommandContext(ctx, cmdName, cmdAgs...)

	cmd.Env = append(os.Environ(), envs...)

	cmd.Stdout = logger.V(5)
	cmd.Stderr = logger.V(5)

	err := cmd.Run()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (Embedded) GetAddress(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	err := Wait(ctx, defaultEmbeddedMinioEndpointAddress)
	if err != nil {
		return "", err
	}

	return defaultEmbeddedMinioEndpointAddress, nil
}

func Wait(ctx context.Context, serverUrl string) (err error) {
	endpoint.Set(serverUrl)

	var lastErr error

	err = wait.PollUntilContextCancel(ctx, 2*time.Second, true,
		func(ctx context.Context) (bool, error) {
			lastErr = IsConnected(ctx)
			if lastErr != nil {
				log.Warnf("waiting for minio to be ready: %v", lastErr)
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
	healthURL := fmt.Sprintf("http://%s/minio/health/ready", endpoint.Get())

	return req.HTTPRequest().
		GetWithContext(ctx, healthURL).
		Error()
}
