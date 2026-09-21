package k8s

import (
	"CBCTF/internal/log"
	"context"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	clientfeatures "k8s.io/client-go/features"
	clientfeaturestesting "k8s.io/client-go/features/testing"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestPodStartupComplete(t *testing.T) {
	for _, tc := range []struct {
		name          string
		status        corev1.PodStatus
		ready, failed bool
	}{
		{name: "pending", status: corev1.PodStatus{Phase: corev1.PodPending}},
		{name: "running is not ready", status: corev1.PodStatus{Phase: corev1.PodRunning}},
		{name: "ready", status: corev1.PodStatus{Phase: corev1.PodRunning, Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}}}, ready: true},
		{name: "completed workload", status: corev1.PodStatus{Phase: corev1.PodSucceeded}, failed: true},
		{name: "failed", status: corev1.PodStatus{Phase: corev1.PodFailed}, failed: true},
		{name: "sidecar hides workload exit", status: corev1.PodStatus{Phase: corev1.PodRunning, ContainerStatuses: []corev1.ContainerStatus{{Name: "web", State: corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{ExitCode: 1}}}}}, failed: true},
		{name: "image pull is retryable", status: corev1.PodStatus{Phase: corev1.PodPending, ContainerStatuses: []corev1.ContainerStatus{{State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "ImagePullBackOff"}}}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ready, err := podStartupComplete(&corev1.Pod{Status: tc.status})
			if ready != tc.ready || (err != nil) != tc.failed {
				t.Fatalf("ready=%t error=%v", ready, err)
			}
		})
	}
}

func useFakePods(t *testing.T, pods ...*corev1.Pod) *fake.Clientset {
	t.Helper()
	clientfeaturestesting.SetFeatureDuringTest(t, clientfeatures.WatchListClient, false)
	previous, namespace := kubeClient, globalNamespace
	client := fake.NewClientset()
	kubeClient, globalNamespace = client, "test"
	log.Init()
	for _, pod := range pods {
		if err := client.Tracker().Add(pod); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() { kubeClient, globalNamespace = previous, namespace })
	return client
}

func TestWaitPodReadyTracksUIDAndLatestState(t *testing.T) {
	for _, tc := range []struct {
		name             string
		exists, replaced bool
		want             string
	}{
		{name: "ready from relist", exists: true},
		{name: "deleted before watch", want: "deleted"},
		{name: "replaced namesake", exists: true, replaced: true, want: "replaced"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			created := &corev1.Pod{Name: "workload", Namespace: "test", UID: "original"}
			client := useFakePods(t)
			if tc.exists {
				current := created.DeepCopy()
				current.Status = corev1.PodStatus{Phase: corev1.PodRunning, Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}}}
				if tc.replaced {
					current.UID = "replacement"
				}
				if err := client.Tracker().Add(current); err != nil {
					t.Fatal(err)
				}
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			pod, err := waitPodReady(ctx, created)
			if tc.want == "" {
				if err != nil || !podReady(pod) {
					t.Fatalf("%v %v", pod, err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %s: %v", tc.want, err)
			}
		})
	}
}

func TestWaitPodReadyTimeoutRetainsDiagnostics(t *testing.T) {
	pod := &corev1.Pod{Name: "slow", Namespace: "test", UID: "uid", Status: corev1.PodStatus{Phase: corev1.PodPending, ContainerStatuses: []corev1.ContainerStatus{{Name: "generator", State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "ImagePullBackOff"}}}}}}
	useFakePods(t, pod)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := waitPodReady(ctx, pod)
	if err == nil || !strings.Contains(err.Error(), "ImagePullBackOff") {
		t.Fatalf("missing diagnostics: %v", err)
	}
}

func TestCreatePodNeverDeletesExistingPod(t *testing.T) {
	pod := &corev1.Pod{Name: "existing", Namespace: "test"}
	client := useFakePods(t, pod)
	_, ret := CreatePod(context.Background(), CreatePodOptions{Name: pod.Name})
	if ret.OK {
		t.Fatal("expected AlreadyExists")
	}
	for _, action := range client.Actions() {
		if action.GetVerb() != "create" {
			t.Fatalf("unexpected API call: %s", action.GetVerb())
		}
	}
}

func TestDeletePodAndWaitIsIdempotent(t *testing.T) {
	useFakePods(t, &corev1.Pod{Name: "done", Namespace: "test", UID: "done-uid"})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for range 2 {
		if ret := DeletePodAndWait(ctx, "done", "done-uid"); !ret.OK {
			t.Fatalf("delete: %+v", ret)
		}
	}
}

func TestWaitPodReadyRelistsAfterExpiredWatch(t *testing.T) {
	pod := &corev1.Pod{Name: "recover", Namespace: "test", UID: "uid", Status: corev1.PodStatus{Phase: corev1.PodPending}}
	client := useFakePods(t, pod)
	first := true
	client.PrependWatchReactor("pods", func(action ktesting.Action) (bool, watch.Interface, error) {
		if !first {
			return false, nil, nil
		}
		first = false
		ready := pod.DeepCopy()
		ready.Status = corev1.PodStatus{Phase: corev1.PodRunning, Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}}}
		if err := client.Tracker().Update(corev1.SchemeGroupVersion.WithResource("pods"), ready, "test"); err != nil {
			return true, nil, err
		}
		watcher := watch.NewRaceFreeFake()
		watcher.Error(&metav1.Status{Code: 410, Reason: metav1.StatusReasonExpired})
		return true, watcher, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ready, err := waitPodReady(ctx, pod)
	if err != nil || !podReady(ready) {
		t.Fatalf("expired watch did not recover: %v %v", ready, err)
	}
}
