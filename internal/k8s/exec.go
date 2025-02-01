package k8s

import (
	"CBCTF/internal/log"
	"bytes"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func ExecInPod(pod, container, command string, stdin io.Reader, stdout *bytes.Buffer, stderr *bytes.Buffer) error {
	cmd := []string{"sh", "-c", command}
	return ExecInPodWithStream(pod, container, cmd, stdin, stdout, stderr)
}

func ExecInPodWithStream(pod, container string, command []string, stdin io.Reader, stdout, stderr *bytes.Buffer) error {
	req := Client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(NamespaceName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdout:    stdout != nil,
			Stderr:    true,
			Stdin:     stdin != nil,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(Config, "POST", req.URL())
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
}
