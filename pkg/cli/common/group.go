package common

import "github.com/spf13/cobra"

var (
	GroupManagement = CreateGroup("management", "Management Commands:")
	GroupAdvanced   = CreateGroup("advanced", "Advanced Commands:")
	GroupOther      = CreateGroup("other", "Other Commands:")
)

func CreateGroup(name, title string) *cobra.Group {
	return &cobra.Group{
		ID:    name,
		Title: title,
	}
}
