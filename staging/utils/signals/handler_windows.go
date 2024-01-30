package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/seal-io/walrus/utils/log"
)

var registered = make(chan struct{})

// Handler registers for signals and returns a context.
func Handler() context.Context {
	close(registered) // Panics when called twice.

	logger := log.WithName("signal").WithName("handler")

	sigs := []os.Signal{syscall.SIGINT}

	sigChan := make(chan os.Signal, len(sigs))
	ctx, cancel := context.WithCancel(context.Background())

	// Register for signals.
	signal.Notify(sigChan, sigs...)

	// Process signals.
	go func() {
		var exited bool

		for sig := range sigChan {
			logger.V(5).Infof("received signal %q", sig)

			if exited {
				os.Exit(1)
			}

			logger.Info("exiting")
			cancel()
			exited = true
		}
	}()

	return ctx
}
