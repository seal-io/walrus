//go:build !linux

package k8s

import "syscall"

func getSysProcAttr() *syscall.SysProcAttr {
	return nil
}
