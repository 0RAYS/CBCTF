package k8s

import (
	"bytes"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	"os"

	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

func test() {
	var err error
	namespace := "default"
	podName := "my-pod"
	containerName := "my-container"
	srcPath := "/path/to/file/in/pod"
	destPath := "/path/to/destination/on/host"

	// 从 Pod 复制文件到主机
	err = copyFileFromPod(namespace, podName, containerName, srcPath, destPath)
	if err != nil {
		fmt.Println(err)
	}

	// 从主机复制文件到 Pod
	err = copyFileToPod(namespace, podName, containerName, "/path/to/file/on/host", "/path/to/destination/in/pod")
	if err != nil {
		fmt.Println(err)
	}
}

func copyFileFromPod(namespace, podName, containerName, srcPath, destPath string) error {
	req := Client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"cat", srcPath},
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(Config, "POST", req.URL())
	if err != nil {
		return err
	}

	var stdout, stderr bytes.Buffer
	err = executor.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return fmt.Errorf("failed to stream: %v, stderr: %s", err, stderr.String())
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, &stdout)
	if err != nil {
		return err
	}

	fmt.Println("File copied successfully from pod")
	return nil
}

func copyFileToPod(namespace, podName, containerName, srcPath, destPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	req := Client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"sh", "-c", fmt.Sprintf("cat > %s", destPath)},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(Config, "POST", req.URL())
	if err != nil {
		return err
	}

	var stdout, stderr bytes.Buffer
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:  file,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return fmt.Errorf("failed to stream: %v, stderr: %s", err, stderr.String())
	}

	fmt.Println("File copied successfully to pod")
	return nil
}
