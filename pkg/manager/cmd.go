package manager

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCommand returns a new cobra command for the manager.
func NewCommand(name, brief string) *cobra.Command {
	o := NewOptions()

	c := &cobra.Command{
		Use:   name,
		Short: brief,
		PreRunE: func(c *cobra.Command, args []string) error {
			return o.Validate(c.Context())
		},
		RunE: func(c *cobra.Command, args []string) error {
			cfg, err := o.Complete(c.Context())
			if err != nil {
				return fmt.Errorf("complete config: %w", err)
			}
			mgr, err := cfg.Apply(c.Context())
			if err != nil {
				return fmt.Errorf("apply config: %w", err)
			}
			err = mgr.Prepare(c.Context())
			if err != nil {
				return fmt.Errorf("prepare manager: %w", err)
			}
			return mgr.Start(c.Context())
		},
	}

	o.AddFlags(c.Flags())

	return c
}
