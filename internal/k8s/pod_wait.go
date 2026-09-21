package k8s

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	watchtools "k8s.io/client-go/tools/watch"
)

func podReady(pod *corev1.Pod) bool {
	if pod.DeletionTimestamp != nil || pod.Status.Phase != corev1.PodRunning {
		return false
	}
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

// Only status fields: never include a Pod spec, environment or command in errors.
func podSummary(pod *corev1.Pod) string {
	details := []string{string(pod.Status.Phase)}
	if pod.Status.Reason != "" {
		details = append(details, pod.Status.Reason)
	}
	for _, condition := range pod.Status.Conditions {
		if condition.Status == corev1.ConditionFalse && condition.Reason != "" {
			details = append(details, string(condition.Type)+"="+condition.Reason)
		}
	}
	for _, statuses := range [][]corev1.ContainerStatus{pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses} {
		for _, status := range statuses {
			if state := status.State.Waiting; state != nil {
				details = append(details, status.Name+"="+state.Reason)
			}
			if state := status.State.Terminated; state != nil {
				details = append(details, fmt.Sprintf("%s=%s(exit %d)", status.Name, state.Reason, state.ExitCode))
			}
		}
	}
	return strings.Join(details, "; ")
}

func podStartupComplete(pod *corev1.Pod) (bool, error) {
	if pod.DeletionTimestamp != nil {
		return false, fmt.Errorf("pod %s is terminating", pod.Name)
	}
	if pod.Status.Phase == corev1.PodFailed || pod.Status.Phase == corev1.PodSucceeded {
		return false, fmt.Errorf("pod %s exited before readiness: %s", pod.Name, podSummary(pod))
	}
	// These are long-lived workloads using RestartPolicyNever, not batch Jobs.
	for _, status := range pod.Status.ContainerStatuses {
		if status.State.Terminated != nil {
			return false, fmt.Errorf("pod %s container exited: %s", pod.Name, podSummary(pod))
		}
	}
	for _, status := range pod.Status.InitContainerStatuses {
		if state := status.State.Terminated; state != nil && state.ExitCode != 0 {
			return false, fmt.Errorf("pod %s init failed: %s", pod.Name, podSummary(pod))
		}
	}
	return podReady(pod), nil
}

// Relist after expired resource versions and retry disconnected watches. Pin the
// UID so a deleted/recreated namesake cannot satisfy an earlier start operation.
func waitPodReady(ctx context.Context, created *corev1.Pod) (*corev1.Pod, error) {
	last := created
	if done, err := podStartupComplete(last); done || err != nil {
		return last, err
	}
	selector := fields.OneTermEqualSelector("metadata.name", created.Name).String()
	pods := kubeClient.CoreV1().Pods(globalNamespace)
	lw := &cache.ListWatch{
		ListWithContextFunc: func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = selector
			return pods.List(ctx, options)
		},
		WatchFuncWithContext: func(ctx context.Context, options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = selector
			return pods.Watch(ctx, options)
		},
	}
	check := func(pod *corev1.Pod) (bool, error) {
		if pod.UID != created.UID {
			return false, fmt.Errorf("pod %s was replaced", created.Name)
		}
		last = pod
		return podStartupComplete(pod)
	}
	_, err := watchtools.UntilWithSync(ctx, lw, &corev1.Pod{}, func(store cache.Store) (bool, error) {
		obj, exists, err := store.GetByKey(globalNamespace + "/" + created.Name)
		if err != nil {
			return false, err
		}
		if !exists {
			return false, fmt.Errorf("pod %s was deleted", created.Name)
		}
		return check(obj.(*corev1.Pod))
	}, func(event watch.Event) (bool, error) {
		if event.Type == watch.Deleted {
			return false, fmt.Errorf("pod %s was deleted", created.Name)
		}
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			return false, nil
		}
		return check(pod)
	})
	if err != nil {
		return nil, fmt.Errorf("waiting for pod %s (%s): %w", created.Name, podSummary(last), err)
	}
	return last, nil
}
