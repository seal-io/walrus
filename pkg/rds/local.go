package rds

import (
	"context"
	"database/sql"
	"errors"
	"os/exec"
	"time"

	"github.com/seal-io/seal/utils/log"
)

const defaultDataSourceAddress = "postgres://root@127.0.0.1:5432/seal?sslmode=disable"

type Embedded struct{}

func (Embedded) Run(ctx context.Context) error {
	const (
		runUser = "postgres"
		cmdName = "postgres"
	)
	var cmdArgs = []string{
		runUser,
		cmdName,
		"-D",
		"/var/lib/postgresql/main",
		"-c",
		"config_file=/etc/postgresql/main/postgresql.conf",
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

func (Embedded) GetDriver(ctx context.Context) (string, string, *sql.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	var drvDialect, drv, err = LoadDriver(defaultDataSourceAddress)
	if err != nil {
		return "", "", nil, err
	}

	err = Wait(ctx, drv)
	if err != nil {
		return "", "", nil, err
	}

	return defaultDataSourceAddress, drvDialect, drv, err
}
