package cache

import (
	"context"
	"time"
)

type DriverStats struct {
	MaxOpenConnections int64 // Maximum number of open connections.
	IdleConnections    int64 // The number of idle connections.

	NewOpenCount int64 // The total number of free connection was not found in the pool.
	TimeoutCount int64 // The total number of timeout getting a connection from the pool.
	ClosedCount  int64 // The total number of connections closed in the pool.
}

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

	// Stats returns driver statistics.
	Stats() DriverStats
}
