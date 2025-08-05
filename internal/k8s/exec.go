package k8s

import (
	"CBCTF/internal/log"
	"bytes"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// Exec executes a command in a Pod
func Exec(pod, container, command string, stdin io.Reader) (bytes.Buffer, bytes.Buffer, error) {
	cmd := []string{"sh", "-c", command}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := func(pod, container string, command []string, stdin io.Reader, stdout, stderr *bytes.Buffer) error {
		req := kubeClient.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(pod).
			Namespace(globalNamespace).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Container: container,
				Command:   command,
				Stdout:    stdout != nil,
				Stderr:    true,
				Stdin:     stdin != nil,
				TTY:       false,
			}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(kubeConfig, "POST", req.URL())
		if err != nil {
			log.Logger.Warningf("Failed to create SPDY executor: %v", err)
			return err
		}
		if stderr == nil {
			stderr = new(bytes.Buffer)
		}
		if stdout == nil {
			stdout = new(bytes.Buffer)
		}

		return exec.Stream(remotecommand.StreamOptions{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		})
	}(pod, container, cmd, stdin, &stdout, &stderr)
	return stdout, stderr, err
}
