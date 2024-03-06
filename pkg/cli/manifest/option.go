package manifest

import (
	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/config"
)

// ApplyOption is a type that represents the options for the manifest apply.
type ApplyOption struct {
	CommonOption

	// Apply in preview mode.
	Preview bool `json:"preview,omitempty"`

	// Comment.
	Comment string `json:"comment,omitempty"`
}

func (f *ApplyOption) AddFlags(cmd *cobra.Command) {
	f.CommonOption.AddFlags(cmd)

	cmd.Flags().BoolVarP(&f.Preview, "preview", "", false, "Applying changes will generate a preview instead of actual deployment")
	cmd.Flags().StringVarP(&f.Comment, "comment", "", "", "Add comment to the operation")
}

// CommonOption is a type that represents the options for the manifest operation.
type CommonOption struct {
	// Context flags.
	config.ScopeContext

	// File path or folder path.
	Filenames []string `json:"filenames,omitempty"`

	// Recursive apply.
	Recursive bool `json:"recursive,omitempty"`

	// Wait for the operation to complete.
	Wait bool `json:"wait,omitempty"`

	// Timeout in seconds for the wait operation.
	Timeout int `json:"timeout,omitempty"`

	// ValidateParameterAllSet is a flag that indicates whether to validate all parameters set.
	ValidateParametersSet bool `json:"validateParametersSet,omitempty"`
}

func (f *CommonOption) AddFlags(cmd *cobra.Command) {
	f.ScopeContext.AddFlags(cmd)

	cmd.Flags().StringSliceVarP(&f.Filenames, "filenames", "f", nil, "File path or folder path")
	cmd.Flags().BoolVarP(&f.Recursive, "recursive", "r", false, "Recursive apply")
	cmd.Flags().BoolVarP(&f.Wait, "wait", "", false, "Wait for the operation to complete")
	cmd.Flags().IntVarP(&f.Timeout, "timeout", "", 300, "Timeout in seconds for the wait operation")
	cmd.Flags().BoolVarP(&f.ValidateParametersSet, "validate-parameters-set", "", false, "Validate all parameters set")
}
