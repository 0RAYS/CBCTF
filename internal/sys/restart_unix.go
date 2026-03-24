//go:build !windows

package sys

import (
	"os"
	"syscall"
)

func Restart(proc *os.Process) error {
	return proc.Signal(syscall.SIGUSR1)
}
