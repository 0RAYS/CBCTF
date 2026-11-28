package k8s

// Ping reports whether the Kubernetes API server is reachable.
func Ping() error {
	if kubeClient == nil {
		return errClientUnavailable
	}
	_, err := kubeClient.Discovery().ServerVersion()
	return err
}
