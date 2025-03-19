package k8s

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"time"
)

func CreatePod(ctx context.Context, docker model.Docker, usage model.Usage, containers []corev1.Container) (*corev1.Pod, bool, string) {
	var (
		pod *corev1.Pod
		err error
		ok  bool
	)
	pod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      docker.PodName,
			Namespace: NamespaceName,
			Labels: map[string]string{
				"app": docker.PodName,
			},
		},
		Spec: corev1.PodSpec{
			Containers:                    containers,
			TerminationGracePeriodSeconds: ptr.To[int64](3),
			RestartPolicy:                 corev1.RestartPolicyNever,
		},
	}
	pod, err = Client.CoreV1().Pods(NamespaceName).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		log.Logger.Warningf("Failed to create pod: %v", err)
		return nil, false, "CreatePodError"
	}
	for {
		pod, ok, _ = GetPod(pod.Name)
		if !ok {
			log.Logger.Warningf("Failed to get pod: %v", err)
			return nil, false, "GetPodError"
		}
		if pod.Status.Phase == corev1.PodRunning {
			break
		}
		if pod.Status.Phase != corev1.PodPending {
			log.Logger.Warningf("Pod %s:%s failed to run", usage.Name, pod.Name)
			return nil, false, "CreatePodError"
		}
	}
	return pod, true, "Success"
}

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
