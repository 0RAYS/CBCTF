package k8s

import (
	"CBCTF/internal/log"
	"bytes"
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// Exec executes a command in a Pod
func Exec(ctx context.Context, pod, container, command string, stdin io.Reader) (*bytes.Buffer, *bytes.Buffer, error) {
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
			Stdin:     stdin != nil,
			TTY:       false,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(kubeConfig, "POST", req.URL())
	if err != nil {
		log.Logger.Warningf("Failed to create SPDY executor: %s", err)
		return nil, nil, err
	}
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	return stdout, stderr, err
}
