//go:build !windows

package signals

import (
	"syscall"
	"unsafe"
)

// https://man7.org/linux/man-pages/man2/sigaction.2.html

const (
	P_ALL  = 0
	P_PID  = 1
	P_PGID = 2
)

const (
	cExited    = 1
	cKilled    = 2
	cDumped    = 3
	cTrapped   = 4
	cStopped   = 5
	cContinued = 6
)

type Siginfo struct {
	Signo  int32 /* Signal number. */
	Errno  int32 /* An errno value. */
	Code   int32 /* Signal code. */
	Trapno int32
	Pid    int32
	Uid    int32
	Status int32 /* Exit value or signal. */
	_      [100]byte
}

func (s Siginfo) Exited() bool {
	return s.Code == cExited
}

func (s Siginfo) ExitStatus() int {
	if !s.Exited() {
		return -1
	}

	return int(s.Status)
}

func (s Siginfo) Killed() bool {
	return s.Code == cKilled
}

func (s Siginfo) KillSignal() syscall.Signal {
	if !s.Killed() {
		return -1
	}

	return syscall.Signal(s.Errno)
}

func (s Siginfo) Stopped() bool {
	return s.Code == cStopped
}

func (s Siginfo) StopSignal() syscall.Signal {
	if !s.Stopped() {
		return -1
	}

	return syscall.Signal(s.Errno)
}

func (s Siginfo) Trapped() bool {
	return s.Code == cTrapped
}

func (s Siginfo) TrapCause() int {
	if !s.Trapped() {
		return -1
	}

	return int(s.Errno)
}

func (s Siginfo) Continued() bool {
	return s.Code == cContinued
}

func (s Siginfo) Dumped() bool {
	return s.Code == cDumped
}

func syscallWaitid(idType, id int, info *Siginfo, options int, rusage *syscall.Rusage) error {
	_, _, e1 := syscall.Syscall6(
		syscall.SYS_WAITID,
		uintptr(idType),
		uintptr(id),
		uintptr(unsafe.Pointer(info)),
		uintptr(options),
		uintptr(unsafe.Pointer(rusage)),
		0,
	)
	if e1 == 0 {
		return nil
	}

	return e1
}
