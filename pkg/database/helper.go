package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/lib/pq" // Db = postgres.
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/seal/utils/log"
)

const (
	DialectPostgres = "postgres"
)

func GetDriverAndName(dataSourceAddress string) (dsd, dsn string, err error) {
	if dataSourceAddress == "" {
		err = errors.New("blank data source address")
		return
	}

	if strings.HasPrefix(dataSourceAddress, "postgres://") {
		dsd = DialectPostgres
		dsn = dataSourceAddress
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

	err = wait.PollUntilContextCancel(ctx, 2*time.Second, true,
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
