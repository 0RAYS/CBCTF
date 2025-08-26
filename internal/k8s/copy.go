package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyToPod copies a file to a Pod
func CopyToPod(ctx context.Context, podName, containerName, src, dst string) error {
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
	_, _, err = Exec(ctx, podName, containerName, fmt.Sprintf("tee %s", dst), &buf)
	return err
}

// CopyFromPod copies a file from a Pod
func CopyFromPod(ctx context.Context, podName, containerName, src, dst string) error {
	command := fmt.Sprintf("cat %s", src)
	stdout, _, err := Exec(ctx, podName, containerName, command, nil)
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
	_, err = io.Copy(file, stdout)
	if err != nil {
		return err
	}
	return nil
}
