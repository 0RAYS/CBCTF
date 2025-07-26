package k8s

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
)

type CreatePodOptions struct {
	Name        string
	Labels      map[string]string
	Annotations map[string]string
	Containers  []corev1.Container
	Volumes     []corev1.Volume
}

func CreatePod(ctx context.Context, options CreatePodOptions) (*corev1.Pod, bool, string) {
	var (
		pod *corev1.Pod
		ok  bool
		err error
	)
	if _, ok, _ = GetPod(ctx, options.Name); ok {
		DeletePod(ctx, options.Name)
	}
	pod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        options.Name,
			Namespace:   GlobalNamespace,
			Labels:      options.Labels,
			Annotations: options.Annotations,
		},
		Spec: corev1.PodSpec{
			EnableServiceLinks:            utils.Ptr(false),
			AutomountServiceAccountToken:  utils.Ptr(false),
			Containers:                    options.Containers,
			Volumes:                       options.Volumes,
			TerminationGracePeriodSeconds: utils.Ptr(int64(3)),
			RestartPolicy:                 corev1.RestartPolicyNever,
		},
	}
	pod, err = kubeClient.CoreV1().Pods(GlobalNamespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Pod: %v", err)
		return nil, false, i18n.CreatePodError
	}
	for {
		pod, ok, _ = GetPod(ctx, options.Name)
		if !ok {
			return nil, false, i18n.GetPodError
		}
		if pod.Status.Phase == corev1.PodRunning {
			break
		}
		if pod.Status.Phase != corev1.PodPending {
			log.Logger.Warningf("Pod %s failed to run", pod.Name)
			return nil, false, i18n.CreatePodError
		}
		time.Sleep(500 * time.Millisecond)
	}
	return pod, true, i18n.Success
}

// GetPod 依据 name 获取 Pod
func GetPod(ctx context.Context, name string) (*corev1.Pod, bool, string) {
	pod, err := kubeClient.CoreV1().Pods(GlobalNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierror.IsNotFound(err) {
			return nil, false, i18n.PodNotFound
		}
		log.Logger.Warningf("Failed to get Pod %s: %v", name, err)
		return nil, false, i18n.GetPodError
	}
	return pod, true, i18n.Success
}

func GetPodList(ctx context.Context, labels ...map[string]string) (*corev1.PodList, bool, string) {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	podList, err := kubeClient.CoreV1().Pods(GlobalNamespace).List(ctx, options)
	if err != nil {
		log.Logger.Warningf("Failed to list Pods: %v", err)
		return nil, false, i18n.GetPodError
	}
	return podList, true, i18n.Success
}

// DeletePod 依据 name 删除 Pod
func DeletePod(ctx context.Context, name string) (bool, string) {
	err := kubeClient.CoreV1().Pods(GlobalNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Pod: %v", err)
		return false, i18n.DeletePodError
	}
	return true, i18n.Success
}

func DeletePodList(ctx context.Context, labels ...map[string]string) (bool, string) {
	var options metav1.ListOptions
	if len(labels) > 0 {
		var selector string
		for k, v := range labels[0] {
			selector += fmt.Sprintf("%s=%s,", k, v)
		}
		options = metav1.ListOptions{
			LabelSelector: strings.TrimSuffix(selector, ","),
		}
	}
	err := kubeClient.CoreV1().Pods(GlobalNamespace).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Pod: %v", err)
		return false, i18n.DeletePodError
	}
	return true, i18n.Success
}
