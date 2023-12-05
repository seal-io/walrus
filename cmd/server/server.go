package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/seal-io/walrus/pkg/server"
	"github.com/seal-io/walrus/utils/clis"
	"github.com/seal-io/walrus/utils/log"
)

func main() {
	cmd := server.Command()

	app := clis.AsApp(cmd)
	if err := app.RunContext(withSignalHandler(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func withSignalHandler() context.Context {
	logger := log.WithName("signal").WithName("handler")

	sigs := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	if syscall.Getpid() == 1 {
		// Reap child processes if we are PID 1.
		sigs = append(sigs, syscall.SIGCHLD)
	}

	// Register for signals.
	sigChan := make(chan os.Signal, len(sigs))
	signal.Notify(sigChan, sigs...)

	// Process signals.
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		var exited bool

		for sig := range sigChan {
			logger.V(5).Infof("received signal %q", sig)

			if !exited && sig == syscall.SIGCHLD {
				logger.V(5).Info("reaping child process")

				for {
					var ws syscall.WaitStatus

					// Reap all child processes.
					pid, err := syscall.Wait4(-1, &ws, syscall.WNOHANG, nil)
					for errors.Is(err, syscall.EINTR) {
						pid, err = syscall.Wait4(pid, &ws, syscall.WNOHANG, nil)
					}

					if pid == 0 || errors.Is(err, syscall.ECHILD) {
						break
					}

					logger.Infof("reaped child process %d with exit code %d", pid, ws.ExitStatus())
				}

				continue
			}

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
