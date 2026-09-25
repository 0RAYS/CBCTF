//go:build !linux

package generatorworker

import "os/exec"

func runCommand(command *exec.Cmd) error { return command.Run() }
