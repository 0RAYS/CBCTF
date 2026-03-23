//go:build windows

package cmd

import (
	"os"
	"os/signal"
)

func registerStopSignals(ch chan<- os.Signal) {
	signal.Notify(ch, os.Interrupt)
}

func registerRestartSignals(ch chan<- os.Signal) {
}
