//go:build !linux

package generator

import "os/exec"

func runCommand(command *exec.Cmd) error { return command.Run() }
