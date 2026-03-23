//go:build windows

package router

import "os"

func restartSignal(proc *os.Process) error {
	return proc.Kill()
}
