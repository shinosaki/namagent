//go:build !windows

package utils

import (
	"os/exec"
	"syscall"
)

func SetSID(proc *exec.Cmd) {
	proc.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
}
