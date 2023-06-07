package types

import "context"

type ExecutableResource interface {
	Exec(ctx context.Context, key string, opts ExecOptions) error
	Supported(ctx context.Context, key string) (bool, error)
}

type LoggableResource interface {
	Log(ctx context.Context, key string, opts LogOptions) error
}
