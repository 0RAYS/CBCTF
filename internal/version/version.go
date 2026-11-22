// Package version carries the build identity of the running binary.
package version

// These values are injected at link time with -ldflags "-X CBCTF/internal/version.Version=...".
var (
	Version = "dev"
	Commit  = "unknown"
)
