package cds

import (
	"context"
	"time"
)

type Driver interface {
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	SetMaxIdleConns(n int)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	SetMaxOpenConns(n int)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	SetConnMaxLifetime(d time.Duration)

	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle.
	SetConnMaxIdleTime(d time.Duration)

	// PingContext verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	PingContext(context.Context) error

	// Underlay returns the underlay client dialect and client.
	Underlay(context.Context) (string, any, error)
}
