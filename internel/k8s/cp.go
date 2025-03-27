package k8s

import (
	"bytes"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"path/filepath"
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
	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(NamespaceName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     true,
			TTY:       false,
		}, scheme.ParameterCodec)
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
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
	var buf bytes.Buffer
	err := ExecInPod(podName, containerName, command, nil, &buf, nil)
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
	_, err = io.Copy(file, &buf)
	if err != nil {
		return err
	}
	return nil
}
