//go:build !windows

package sys

import (
	"os"
	"os/signal"
	"syscall"
)

func RegisterStopSignals(ch chan<- os.Signal) {
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
}

func RegisterRestartSignals(ch chan<- os.Signal) {
	signal.Notify(ch, syscall.SIGUSR1)
}
