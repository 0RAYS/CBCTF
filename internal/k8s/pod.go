package k8s

import (
	"CBCTF/internal/log"
	"context"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// GetPods 获取所有 Pod
func GetPods() (*corev1.PodList, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	pods, err := Client.CoreV1().Pods(NamespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to get pods: %v", err)
		return &corev1.PodList{}, false, "GetPodError"
	}
	return pods, true, "Success"
}

// GetPod 依据 name 获取 Pod
func GetPod(name string) (*corev1.Pod, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	pod, err := Client.CoreV1().Pods(NamespaceName).Get(ctx, name, metav1.GetOptions{})
	if apierror.IsNotFound(err) {
		return &corev1.Pod{}, false, "PodNotFound"
	}
	if err != nil {
		log.Logger.Warningf("Failed to get pod: %v", err)
		return &corev1.Pod{}, false, "GetPodError"
	}
	return pod, true, "Success"
}

// DeletePod 依据 name 删除 Pod
func DeletePod(name string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	err := Client.CoreV1().Pods(NamespaceName).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete pod: %v", err)
		return false, "DeletePodError"
	}
	return true, "Success"
}
