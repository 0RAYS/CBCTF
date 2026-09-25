//go:build !linux

package worker

import "os/exec"

func runCommand(command *exec.Cmd) error { return command.Run() }
