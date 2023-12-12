package common

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// BindFlags bind the environment with flags.
func BindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "")

		// Use viper value when the flag is not set and environment variable has a value.
		if !f.Changed && viper.IsSet(configName) {
			val := viper.Get(configName)

			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				panic(err)
			}
		}
	})
}

// InputByFlags check whether user set flags.
func InputByFlags(cmd *cobra.Command) bool {
	var change bool

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			change = true
		}
	})

	return change
}
