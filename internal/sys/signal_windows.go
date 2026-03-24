//go:build windows

package sys

import (
	"os"
	"os/signal"
)

func RegisterStopSignals(ch chan<- os.Signal) {
	signal.Notify(ch, os.Interrupt)
}

func RegisterRestartSignals(_ chan<- os.Signal) {
}
