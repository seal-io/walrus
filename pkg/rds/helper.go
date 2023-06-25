package rds

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // Db = mysql.
	_ "github.com/lib/pq"              // Db = postgres.
	_ "github.com/mattn/go-sqlite3"    // Db = sqlite3.
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/seal/utils/log"
)

func GetDriverAndName(dataSourceAddress string) (dsd, dsn string, err error) {
	if dataSourceAddress == "" {
		err = errors.New("blank data source address")
		return
	}

	switch {
	case strings.HasPrefix(dataSourceAddress, "postgres://"):
		dsd = "postgres"
		dsn = dataSourceAddress
	case strings.HasPrefix(dataSourceAddress, "file:"):
		dsd = "sqlite3"
		dsn = dataSourceAddress
	case strings.HasPrefix(dataSourceAddress, "mysql://"):
		dsd = "mysql"
		dsn = strings.TrimPrefix(dataSourceAddress, "mysql://")
	}

	if dsd == "" {
		err = errors.New("cannot recognize driver from data source address")
	}

	return
}

func LoadDriver(dataSourceAddress string) (drvDialect string, drv *sql.DB, err error) {
	drvDialect, drvSource, err := GetDriverAndName(dataSourceAddress)
	if err != nil {
		return
	}
	drv, err = sql.Open(drvDialect, drvSource)

	return
}

func Wait(ctx context.Context, drv *sql.DB) (err error) {
	var lastErr error

	err = wait.PollImmediateUntilWithContext(ctx, 2*time.Second,
		func(ctx context.Context) (bool, error) {
			lastErr = IsConnected(ctx, drv)
			if lastErr != nil {
				log.Warnf("waiting for database to be ready: %v", lastErr)
			}

			return lastErr == nil, ctx.Err()
		},
	)
	if err != nil && lastErr != nil {
		err = lastErr // Use last error to overwrite context error while existed.
	}

	return
}

func IsConnected(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}
