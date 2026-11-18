//go:build !windows

package k8s

import (
	"os"
	"syscall"
)

// captureDirOwner reports the owner of a file on unix like systems.
func captureDirOwner(info os.FileInfo) (int, int, bool) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, false
	}
	return int(stat.Uid), int(stat.Gid), true
}
