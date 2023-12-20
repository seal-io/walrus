package schema

import "github.com/spf13/cobra"

const (
	flagSchemaTemplateDirName = "dir"
)

type GenerateOption struct {
	Dir string
}

func (o *GenerateOption) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&o.Dir,
		flagSchemaTemplateDirName,
		"",
		"Template dir to generate",
	)
}
