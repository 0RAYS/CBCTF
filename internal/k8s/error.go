package k8s

import (
	"errors"
)

// errMissingPodIdentity marks a deletion that was refused because the Pod
// could not be identified safely.
var errMissingPodIdentity = errors.New("refusing deletion without a Pod name and observed UID")

// errMissingNamespace marks a deletion that was refused because the target
// namespace was not given.
var errMissingNamespace = errors.New("refusing deletion without a namespace")
