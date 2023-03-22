package operator

import (
	"context"
	"io"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
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

	// GetKeys returns keys from the given resource.
	GetKeys(context.Context, *model.ApplicationResource) (*Keys, error)

	// GetStatus gets status of the given resource,
	// must returns GeneralStatusError if raises error.
	GetStatus(context.Context, *model.ApplicationResource) (*status.Status, error)

	// GetEndpoints gets endpoints of the given resource.
	GetEndpoints(context.Context, *model.ApplicationResource) ([]types.ApplicationResourceEndpoint, error)

	// Log gets logs from the given key.
	Log(context.Context, string, LogOptions) error

	// Exec executes commands to the given key.
	Exec(context.Context, string, ExecOptions) error

	// Label apply labels to the resource.
	Label(context.Context, *model.ApplicationResource, map[string]string) error
}

// Keys holds key for next query,
// it is a response block to assist frontend building multiple-levels selector.
type Keys struct {
	// Labels stores label of layer,
	// its length means each key contains levels with the same value as level.
	Labels []string `json:"labels,omitempty"`
	// Keys stores key in tree.
	Keys []Key `json:"keys,omitempty"`
}

// Key holds hierarchy query keys.
type Key struct {
	// Keys indicates the subordinate keys,
	// usually, it should not be valued in leaves.
	Keys []Key `json:"keys,omitempty"`
	// Name indicates the name of the key.
	Name string `json:"name"`
	// Value indicates the value of the key,
	// usually, it should be valued in leaves.
	Value string `json:"value,omitempty"`
	// Loggable indicates whether to be able to get log.
	Loggable *bool `json:"loggable,omitempty"`
	// Executable indicates whether to be able to execute remote command.
	Executable *bool `json:"executable,omitempty"`
}

// LogOptions holds the options of Operator's Log action.
type LogOptions struct {
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
