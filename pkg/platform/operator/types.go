package operator

import (
	"context"
	"io"

	"github.com/seal-io/seal/pkg/dao/model"
)

// Type indicates the type of Operator,
// e.g. Kubernetes, AWS, etc.
type Type = string

// Operator holds the actions that an operator must satisfy.
type Operator interface {
	// Type returns Type.
	Type() Type

	// IsConnected validates whether is connected.
	IsConnected(context.Context) (bool, error)

	// Log gets logs from the given resource.
	Log(context.Context, model.ApplicationResource, LogOptions) error

	// Exec executes commands to the given resource.
	Exec(context.Context, model.ApplicationResource, ExecOptions) error
}

// LogOptions holds the options of Operator's Log action.
type LogOptions struct {
	// Key indicates the key for accessing target,
	// parses by the Operator.
	Key string
	// Out receives the output.
	Out io.Writer
	// Previous indicates to get the previous log of the accessing target.
	Previous bool
	// Tail indicates to get the tail log of the accessing target.
	Tail bool
	// SinceSeconds returns logs newer than a relative duration.
	SinceSeconds *int64
	// Timestamps returns logs with RFC3339 or RFC3339Nano timestamp.
	Timestamps bool
}

// ExecOptions holds the options of Operator's Exec action.
type ExecOptions struct {
	// Key indicates the kye for accessing target,
	// parses by the Operator.
	Key string
	// Out receives the output.
	Out io.Writer
	// In passes the input.
	In io.Reader
	// Shell indicates to launch what kind of shell.
	Shell string
	// Resizer indicates to resize the size(width, height) of the terminal.
	Resizer TerminalResizer
}

// TerminalResizer holds the options to resize the terminal.
type TerminalResizer interface {
	// Next returns the new terminal size after the terminal has been resized.
	// It returns false when monitoring has been stopped.
	Next() (width, height uint16, ok bool)
}
