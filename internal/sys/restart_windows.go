//go:build windows

package sys

import "os"

func Restart(proc *os.Process) error {
	return proc.Kill()
}
