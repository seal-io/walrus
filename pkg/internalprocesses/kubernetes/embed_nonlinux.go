//go:build !linux

package kubernetes

import "syscall"

func getSysProcAttr() *syscall.SysProcAttr {
	return nil
}
