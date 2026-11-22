package k8s

import (
	"context"
)

// Ping reports whether the Kubernetes API server is reachable.
func Ping(ctx context.Context) error {
	if kubeClient == nil {
		return errClientUnavailable
	}
	_, err := kubeClient.Discovery().ServerVersion()
	return err
}
