package rds

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/seal-io/seal/pkg/consts"
	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

const defaultDataSourceAddress = "postgres://root@127.0.0.1:5432/seal?sslmode=disable"

type Embedded struct{}

func (Embedded) Run(ctx context.Context) error {
	// Create run data dir if not found.
	runDataPath := filepath.Join(consts.DataDir, "postgresql")
	if !files.Exists(runDataPath) {
		err := files.Copy(
			"/var/lib/postgresql/main",
			runDataPath,
			files.CopyWithTimes(),
			files.CopyWithOwner())
		if err != nil {
			return fmt.Errorf("error copy initial data: %w", err)
		}
	}

	const (
		runUser = "postgres"
		cmdName = "postgres"
	)
	logger := log.WithName(cmdName)
	cmdArgs := []string{
		runUser,
		cmdName,
		"-D", runDataPath,
		"-c", "config_file=/etc/postgresql/main/postgresql.conf",
	}
	logger.Infof("run: gosu %s", strs.Join(" ", cmdArgs...))
	cmd := exec.CommandContext(ctx, "gosu", cmdArgs...)
	cmd.Stdout = logger.V(5)
	cmd.Stderr = logger.V(5)
	err := cmd.Run()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (Embedded) GetDriver(ctx context.Context) (string, string, *sql.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	drvDialect, drv, err := LoadDriver(defaultDataSourceAddress)
	if err != nil {
		return "", "", nil, err
	}

	err = Wait(ctx, drv)
	if err != nil {
		return "", "", nil, err
	}

	return defaultDataSourceAddress, drvDialect, drv, err
}
