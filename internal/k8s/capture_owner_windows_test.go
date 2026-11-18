//go:build windows

package k8s

import "os"

// captureDirOwner cannot report the owner of a file on windows.
func captureDirOwner(_ os.FileInfo) (int, int, bool) {
	return 0, 0, false
}
