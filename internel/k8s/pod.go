package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/utils"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type CreatePodOptions struct {
	Name        string
	PodIP       string
	Labels      map[string]string
	Containers  []corev1.Container
	Volumes     []corev1.Volume
	HostAliases []corev1.HostAlias
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
			Name:      options.Name,
			Namespace: namespaceName,
			Labels:    options.Labels,
		},
		Spec: corev1.PodSpec{
			EnableServiceLinks:            utils.Ptr(false),
			Containers:                    options.Containers,
			Volumes:                       options.Volumes,
			TerminationGracePeriodSeconds: utils.Ptr(int64(3)),
			RestartPolicy:                 corev1.RestartPolicyNever,
			HostAliases:                   options.HostAliases,
		},
	}
	if options.PodIP != "" {
		pod.ObjectMeta.Annotations = map[string]string{
			"cni.projectcalico.org/ipAddrs": fmt.Sprintf("[\"%s\"]", options.PodIP),
			"cni.projectcalico.org/podIP":   options.PodIP,
			"cni.projectcalico.org/podIPs":  options.PodIP,
		}
	}
	pod, err = kubeClient.CoreV1().Pods(namespaceName).Create(ctx, pod, metav1.CreateOptions{})
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

// GetPods 获取所有 Pod
func GetPods(ctx context.Context) (*corev1.PodList, bool, string) {
	pods, err := kubeClient.CoreV1().Pods(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to get Pods: %v", err)
		return &corev1.PodList{}, false, i18n.GetPodError
	}
	return pods, true, i18n.Success
}

// GetPod 依据 name 获取 Pod
func GetPod(ctx context.Context, name string) (*corev1.Pod, bool, string) {
	pod, err := kubeClient.CoreV1().Pods(namespaceName).Get(ctx, name, metav1.GetOptions{})
	if apierror.IsNotFound(err) {
		return &corev1.Pod{}, false, i18n.PodNotFound
	}
	if err != nil {
		log.Logger.Warningf("Failed to get Pod %s: %v", name, err)
		return &corev1.Pod{}, false, i18n.GetPodError
	}
	return pod, true, i18n.Success
}

// DeletePod 依据 name 删除 Pod
func DeletePod(ctx context.Context, name string) (bool, string) {
	err := kubeClient.CoreV1().Pods(namespaceName).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Pod %s: %v", name, err)
		return false, i18n.DeletePodError
	}
	return true, i18n.Success
}
