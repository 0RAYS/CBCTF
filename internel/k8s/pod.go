package k8s

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func CreatePod(ctx context.Context, podName string, containers []corev1.Container, volumes []corev1.Volume, podIP string, dns map[string]string) (*corev1.Pod, bool, string) {
	var (
		pod *corev1.Pod
		err error
		ok  bool
	)
	if _, ok, _ := GetPod(ctx, podName); ok {
		DeletePod(ctx, podName)
	}
	pod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: NamespaceName,
			Labels: map[string]string{
				"victim": podName,
			},
			Annotations: map[string]string{
				"cni.projectcalico.org/ipAddrs": fmt.Sprintf("[\"%s\"]", podIP),
				"cni.projectcalico.org/podIP":   podIP,
				"cni.projectcalico.org/podIPs":  podIP,
			},
		},
		Spec: corev1.PodSpec{
			Containers:                    containers,
			Volumes:                       volumes,
			TerminationGracePeriodSeconds: ptr.To[int64](3),
			RestartPolicy:                 corev1.RestartPolicyNever,
			HostAliases: func() []corev1.HostAlias {
				aliases := make([]corev1.HostAlias, 0)
				tmp := make(map[string][]string)
				for k, v := range dns {
					if v == podIP {
						v = "127.0.0.1"
					}
					tmp[v] = append(tmp[v], k)
				}
				for ip, hostname := range tmp {
					aliases = append(aliases, corev1.HostAlias{
						IP:        ip,
						Hostnames: hostname,
					})
				}
				return aliases
			}(),
		},
	}
	pod, err = client.CoreV1().Pods(NamespaceName).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create Pod: %v", err)
		return nil, false, i18n.CreatePodError
	}
	for {
		pod, ok, _ = GetPod(ctx, podName)
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
	}
	return pod, true, i18n.Success
}

// GetPods 获取所有 Pod
func GetPods(ctx context.Context) (*corev1.PodList, bool, string) {
	pods, err := client.CoreV1().Pods(NamespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to get Pods: %v", err)
		return &corev1.PodList{}, false, i18n.GetPodError
	}
	return pods, true, i18n.Success
}

// GetPod 依据 name 获取 Pod
func GetPod(ctx context.Context, name string) (*corev1.Pod, bool, string) {
	pod, err := client.CoreV1().Pods(NamespaceName).Get(ctx, name, metav1.GetOptions{})
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
	err := client.CoreV1().Pods(NamespaceName).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierror.IsNotFound(err) {
		log.Logger.Warningf("Failed to delete Pod %s: %v", name, err)
		return false, i18n.DeletePodError
	}
	return true, i18n.Success
}
