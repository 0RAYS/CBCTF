//go:build !windows

package router

import (
	"os"
	"syscall"
)

func restartSignal(proc *os.Process) error {
	return proc.Signal(syscall.SIGUSR1)
}
