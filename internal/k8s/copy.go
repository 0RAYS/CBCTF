package k8s

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// CopyToPod copies a file to a Pod
func CopyToPod(podName, containerName, src, dst string) error {
	var buf bytes.Buffer
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	_, err = io.Copy(&buf, file)
	if err != nil {
		return err
	}
	command := []string{"tee", dst}
	req := kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(globalNamespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     true,
			TTY:       false,
		}, scheme.ParameterCodec)
	executor, err := remotecommand.NewSPDYExecutor(kubeConfig, "POST", req.URL())
	if err != nil {
		return err
	}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin: &buf,
		Tty:   false,
	})
	if err != nil {
		return err
	}
	return nil
}

// CopyFromPod copies a file from a Pod
func CopyFromPod(podName, containerName, src, dst string) error {
	command := fmt.Sprintf("cat %s", src)
	stdout, _, err := Exec(podName, containerName, command, nil)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	_, err = io.Copy(file, &stdout)
	if err != nil {
		return err
	}
	return nil
}
