//go:build !windows

package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

func registerStopSignals(ch chan<- os.Signal) {
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
}

func registerRestartSignals(ch chan<- os.Signal) {
	signal.Notify(ch, syscall.SIGUSR1)
}
