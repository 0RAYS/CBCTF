package k8s

import (
	"CBCTF/internal/log"
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// Exec executes a command in a Pod, draining output without buffering it.
func Exec(ctx context.Context, pod, container, command string) error {
	cmd := []string{"sh", "-c", command}
	req := kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(globalNamespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   cmd,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(kubeConfig, "POST", req.URL())
	if err != nil {
		log.Logger.Warningf("Failed to create SPDY executor: %s", err)
		return err
	}
	return exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
}
