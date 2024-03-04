package args

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/gengo/args"
)

// CustomArgs is used by the gengo framework to pass args specific to this generator.
type CustomArgs struct{}

// NewDefaults returns default arguments for the generator.
func NewDefaults() (*args.GeneratorArgs, *CustomArgs) {
	genericArgs := args.Default().WithoutDefaultFlagParsing()

	customArgs := &CustomArgs{}
	genericArgs.CustomArgs = customArgs

	genericArgs.OutputFileBaseName = "apireg_generated"

	return genericArgs, customArgs
}

// AddFlags add the generator flags to the flag set.
func (ca *CustomArgs) AddFlags(fs *pflag.FlagSet) {
}

// Validate checks the given arguments.
func Validate(genericArgs *args.GeneratorArgs) error {
	if _, ok := genericArgs.CustomArgs.(*CustomArgs); !ok {
		return fmt.Errorf("expected CustomArgs to be of type %T", &CustomArgs{})
	}

	return nil
}
