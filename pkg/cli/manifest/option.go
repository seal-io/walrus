package manifest

import (
	"github.com/spf13/cobra"

	"github.com/seal-io/walrus/pkg/cli/config"
)

// ApplyOption is a type that represents the options for the manifest apply.
type ApplyOption struct {
	CommonOption

	// ChangeComment.
	ChangeComment string `json:"changeComment,omitempty"`
}

func (f *ApplyOption) AddFlags(cmd *cobra.Command) {
	f.CommonOption.AddFlags(cmd)

	cmd.Flags().StringVarP(&f.ChangeComment, "change-comment", "", "", "Add comment to the operation")
}

// PreviewOption is a type that represents the options for the manifest preview.
type PreviewOption struct {
	CommonOption

	// Apply.
	Apply bool `json:"apply,omitempty"`
	// ChangeComment.
	ChangeComment string `json:"changeComment,omitempty"`
	// RunLabels.
	RunLabels map[string]string `json:"runLabels,omitempty"`
}

func (f *PreviewOption) AddFlags(cmd *cobra.Command) {
	f.CommonOption.AddFlags(cmd)

	cmd.Flags().BoolVarP(&f.Apply, "apply", "", false, "Apply previews with the provided labels")
	cmd.Flags().StringVarP(&f.ChangeComment, "change-comment", "", "", "Add comment to previews")
	cmd.Flags().StringToStringVar(&f.RunLabels, "run-labels", nil, "Labels for resource runs")
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
