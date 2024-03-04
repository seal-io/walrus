//go:build !windows

package signalx

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"k8s.io/klog/v2"
)

var registered = make(chan struct{})

// Handler registers for signals and returns a context.
func Handler() context.Context {
	close(registered) // Panics when called twice.

	logger := klog.Background().WithName("signal").WithName("handler")

	sigs := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	isRootProc := syscall.Getpid() == 1

	if isRootProc { // Reap child processes.
		sigs = append(sigs, syscall.SIGCHLD)
	}

	sigChan := make(chan os.Signal, len(sigs))
	reapChan := make(chan os.Signal)
	ctx, cancel := context.WithCancel(context.Background())

	// Register for signals.
	signal.Notify(sigChan, sigs...)

	// Process signals.
	go func() {
		var exited bool

		for sig := range sigChan {
			logger.V(5).
				Info("received", "signal", sig)

			if !exited && sig == syscall.SIGCHLD {
				select {
				case reapChan <- sig:
				default:
					// Don't block if the channel is full.
				}

				continue
			}

			if exited {
				os.Exit(1)
			}

			logger.Info("exiting")
			cancel()
			exited = true

			close(reapChan)
		}
	}()

	if !isRootProc {
		return ctx
	}

	// Reap child processes.
	go func() {
		for range reapChan {
			const zombieTimeWait = 500 * time.Millisecond

			for {
				runtime.Gosched()

				// Use waitid to get an exited child process,
				// not like waitpid, waitid can leave the child process in waitable state.
				var si Siginfo

				err := syscallWaitid(P_ALL, 0, &si, syscall.WEXITED|syscall.WNOWAIT, nil)
				for errors.Is(err, syscall.EINTR) { // Try again if got interrupted.
					err = syscallWaitid(P_ALL, 0, &si, syscall.WEXITED|syscall.WNOWAIT, nil)
				}

				pid := int(si.Pid)

				if pid <= 0 {
					if err != nil {
						logger.V(5).
							Error(err, "reap child process")
					}

					break // Cannot catch any child process, go back to poll.
				}

				if si.Exited() {
					// Since EXIT_ZOMBIE is before than EXIT_DEAD,
					// wait a period to avoid accidentally recycling the transitional process.
					time.Sleep(zombieTimeWait)

					if !isZombieProcess(pid) {
						continue // Go ahead to confirm next child process.
					}

					// Wait again if the process is still a zombie process,
					// so that the process can be recycled by its parent.
					time.Sleep(zombieTimeWait)
				}

				// Use wait4(waitpid) to reap the abnormal exited child process.
				wpid, err := syscall.Wait4(pid, nil, syscall.WNOHANG, nil)
				for errors.Is(err, syscall.EINTR) { // Try again if got interrupted.
					wpid, err = syscall.Wait4(pid, nil, syscall.WNOHANG, nil)
				}

				if wpid <= 0 {
					continue // Cannot recycle the candidate, go ahead to confirm next.
				}

				logger.V(5).
					Info("reaped child process", "pid", wpid)
			}
		}
	}()

	return ctx
}

// isZombieProcess returns true if the process is a zombie process,
// which reads the /proc/<pid>/stat file to check if the process is a zombie process.
func isZombieProcess(pid int) bool {
	bs, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return false
	}

	return bs[strings.IndexRune(string(bs), ')')+2] == 'Z'
}
