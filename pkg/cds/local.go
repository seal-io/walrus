package cds

import (
	"context"
	"errors"
	"os/exec"
	"time"

	"github.com/seal-io/seal/utils/log"
)

const defaultCacheSourceAddress = "redis://default:Default123@127.0.0.1:6379"

type Embedded struct{}

func (Embedded) Run(ctx context.Context) error {
	const (
		runUser = "redis"
		cmdName = "redis-server"
	)
	var cmdArgs = []string{
		runUser,
		cmdName,
		"--save",
		"\"\"",
		"--appendonly",
		"no",
		"--databases",
		"1",
		"--requirepass",
		"Default123",
	}
	var cmd = exec.CommandContext(ctx, "gosu", cmdArgs...) // switch by gosu.
	var logger = log.WithName(cmdName).V(5)
	cmd.Stdout = logger
	cmd.Stderr = logger
	var err = cmd.Run()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (Embedded) GetDriver(ctx context.Context) (string, string, Driver, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	var drvDialect, drv, err = LoadDriver(defaultCacheSourceAddress)
	if err != nil {
		return "", "", nil, err
	}

	err = Wait(ctx, drv)
	if err != nil {
		return "", "", nil, err
	}

	return defaultCacheSourceAddress, drvDialect, drv, err
}
