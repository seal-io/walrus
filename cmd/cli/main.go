package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/version"
)

var (
	cliVersion = version.Get()
	cliName    = "seal"
)

func main() {
	if err := Init(); err != nil {
		fmt.Fprintf(os.Stderr, "error init: %v\n", err)
		os.Exit(1)
	}

	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func Run() (returnErr error) {
	defer func() {
		if err := recover(); err != nil {
			log.Debugf("Error Stack: %s", string(debug.Stack()))

			if e, ok := err.(error); ok {
				returnErr = e
			} else {
				returnErr = fmt.Errorf("%v", err)
			}
		}
	}()

	if serverConfig.Server != "" {
		// Set log level to ignore debug log for generate sub command.
		log.SetLevel(log.WarnLevel)

		err := serverConfig.ValidateAndSetup()
		if err == nil {
			err = load(serverConfig, root, false)
			if err != nil {
				return err
			}
		}
	}

	return root.Execute()
}
