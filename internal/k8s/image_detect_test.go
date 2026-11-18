package k8s

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	ktesting "k8s.io/client-go/testing"
)

// truncatedNode returns a node whose image list is at the kubelet limit, so
// the list alone cannot decide the result.
func truncatedNode(name string) *corev1.Node {
	node := &corev1.Node{
		Name: name,
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
		},
	}
	for i := 0; i < nodeImageReportLimit; i++ {
		node.Status.Images = append(node.Status.Images, corev1.ContainerImage{
			Names: []string{fmt.Sprintf("docker.io/library/filler-%d:latest", i)},
		})
	}
	return node
}

// The detection must probe with imagePullPolicy: Never, so a node that holds
// the image starts the container and a node that lacks it fails to pull.
func TestDetectNodeImagesProbesWithPullPolicyNever(t *testing.T) {
	present := "docker.0rays.club/jbnrz/linux-sjk:latest"
	absent := "docker.0rays.club/jbnrz/missing:latest"
	client := useFakePods(t)
	t.Cleanup(Stop)
	if err := client.Tracker().Add(truncatedNode("worker-a")); err != nil {
		t.Fatal(err)
	}
	// The fake tracker has no resource-version replay, unlike the API server.
	// Establish its watch before simulating the probe Pod.
	watching := make(chan struct{})
	var once sync.Once
	client.PrependWatchReactor("pods", func(action ktesting.Action) (bool, watch.Interface, error) {
		watcher, err := client.Tracker().Watch(action.GetResource(), action.GetNamespace())
		once.Do(func() { close(watching) })
		return true, watcher, err
	})
	client.PrependReactor("create", "jobs", func(action ktesting.Action) (bool, runtime.Object, error) {
		select {
		case <-watching:
		case <-time.After(time.Second):
			return true, nil, fmt.Errorf("pod watch was not established")
		}
		job := action.(ktesting.CreateAction).GetObject().(*batchv1.Job)
		for _, container := range job.Spec.Template.Spec.Containers {
			if container.ImagePullPolicy != corev1.PullNever {
				t.Errorf("probe pull policy = %s, want %s", container.ImagePullPolicy, corev1.PullNever)
			}
		}
		pod := &corev1.Pod{
			Name:      job.Name + "-pod",
			Namespace: job.Namespace,
			UID:       "pod-uid",
			Labels:    job.Spec.Template.Labels,
			Spec:      job.Spec.Template.Spec,
		}
		pod.Spec.NodeName = "worker-a"
		statuses := make([]corev1.ContainerStatus, 0, len(job.Spec.Template.Spec.Containers))
		for _, container := range job.Spec.Template.Spec.Containers {
			status := corev1.ContainerStatus{Name: container.Name}
			if container.Image == present {
				status.ImageID = "sha256:ready"
			} else {
				status.State.Waiting = &corev1.ContainerStateWaiting{Reason: "ErrImageNeverPull"}
			}
			statuses = append(statuses, status)
		}
		pod.Status.ContainerStatuses = statuses
		pods := corev1.SchemeGroupVersion.WithResource("pods")
		if err := client.Tracker().Add(pod); err != nil {
			return true, nil, err
		}
		if err := client.Tracker().Update(pods, pod, pod.Namespace); err != nil {
			return true, nil, err
		}
		return true, job, nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	detections, err := DetectNodeImages(ctx, []string{present, absent}, []*corev1.Node{truncatedNode("worker-a")})
	if err != nil {
		t.Fatal(err)
	}
	if len(detections) != 1 {
		t.Fatalf("detections = %v, want one node", detections)
	}
	if !slices.Equal(detections[0].Present, []string{present}) {
		t.Fatalf("present = %v, want %v", detections[0].Present, []string{present})
	}
	if !slices.Equal(detections[0].Absent, []string{absent}) {
		t.Fatalf("absent = %v, want %v", detections[0].Absent, []string{absent})
	}
	if len(detections[0].Unverified) != 0 {
		t.Fatalf("unverified = %v, want none", detections[0].Unverified)
	}
}

// A probe that never answers must not be reported as absent.
func TestDetectNodeImagesKeepsUnansweredProbesUnverified(t *testing.T) {
	image := "docker.0rays.club/jbnrz/linux-sjk:latest"
	client := useFakePods(t)
	t.Cleanup(Stop)
	client.PrependReactor("create", "jobs", func(action ktesting.Action) (bool, runtime.Object, error) {
		job := action.(ktesting.CreateAction).GetObject().(*batchv1.Job)
		return true, job, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	detections, err := DetectNodeImages(ctx, []string{image}, []*corev1.Node{truncatedNode("worker-a")})
	if err != nil {
		t.Fatal(err)
	}
	if len(detections) != 1 || !slices.Equal(detections[0].Unverified, []string{image}) {
		t.Fatalf("detections = %v, want an unverified probe", detections)
	}
	if len(detections[0].Absent) != 0 || len(detections[0].Present) != 0 {
		t.Fatalf("detections = %v, want no answer", detections)
	}
}

// A node that holds few images reports all of them, so its list decides the
// result without a probe.
func TestDetectNodeImagesUsesTheCompleteNodeImageList(t *testing.T) {
	held := "docker.0rays.club/jbnrz/linux-sjk:latest"
	missing := "docker.0rays.club/jbnrz/missing:latest"
	client := useFakePods(t)
	t.Cleanup(Stop)
	node := &corev1.Node{
		Name: "worker-a",
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
			Images:     []corev1.ContainerImage{{Names: []string{held}}},
		},
	}
	if err := client.Tracker().Add(node); err != nil {
		t.Fatal(err)
	}
	probes := 0
	client.PrependReactor("create", "jobs", func(action ktesting.Action) (bool, runtime.Object, error) {
		probes++
		job := action.(ktesting.CreateAction).GetObject().(*batchv1.Job)
		return true, job, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	detections, err := DetectNodeImages(ctx, []string{held, missing}, []*corev1.Node{node})
	if err != nil {
		t.Fatal(err)
	}
	if probes != 0 {
		t.Fatalf("probes = %d, want none", probes)
	}
	if len(detections) != 1 {
		t.Fatalf("detections = %v, want one node", detections)
	}
	if !slices.Equal(detections[0].Present, []string{held}) {
		t.Fatalf("present = %v, want %v", detections[0].Present, []string{held})
	}
	if !slices.Equal(detections[0].Absent, []string{missing}) {
		t.Fatalf("absent = %v, want %v", detections[0].Absent, []string{missing})
	}
}

// A truncated list must not hide an image the node runs, so the running Pods
// prove its presence without a probe.
func TestDetectNodeImagesUsesTheRunningPodImages(t *testing.T) {
	image := "docker.0rays.club/jbnrz/linux-sjk:latest"
	client := useFakePods(t)
	t.Cleanup(Stop)
	if err := client.Tracker().Add(truncatedNode("worker-a")); err != nil {
		t.Fatal(err)
	}
	pod := &corev1.Pod{
		Name:      "challenge-pod",
		Namespace: globalNamespace,
		Spec: corev1.PodSpec{
			NodeName:   "worker-a",
			Containers: []corev1.Container{{Name: "challenge", Image: image}},
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{{Name: "challenge", ImageID: "sha256:ready"}},
		},
	}
	if err := client.Tracker().Add(pod); err != nil {
		t.Fatal(err)
	}
	probes := 0
	client.PrependReactor("create", "jobs", func(action ktesting.Action) (bool, runtime.Object, error) {
		probes++
		job := action.(ktesting.CreateAction).GetObject().(*batchv1.Job)
		return true, job, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	detections, err := DetectNodeImages(ctx, []string{image}, []*corev1.Node{truncatedNode("worker-a")})
	if err != nil {
		t.Fatal(err)
	}
	if probes != 0 {
		t.Fatalf("probes = %d, want none", probes)
	}
	if len(detections) != 1 || !slices.Equal(detections[0].Present, []string{image}) {
		t.Fatalf("detections = %v, want the running image to be present", detections)
	}
	if len(detections[0].Absent) != 0 || len(detections[0].Unverified) != 0 {
		t.Fatalf("detections = %v, want a present image", detections)
	}
}

// An empty image list proves nothing, so such a node still gets a probe.
func TestDetectNodeImagesProbesANodeWithoutImages(t *testing.T) {
	image := "docker.0rays.club/jbnrz/linux-sjk:latest"
	client := useFakePods(t)
	t.Cleanup(Stop)
	node := &corev1.Node{
		Name: "worker-a",
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
		},
	}
	if err := client.Tracker().Add(node); err != nil {
		t.Fatal(err)
	}
	probes := 0
	client.PrependReactor("create", "jobs", func(action ktesting.Action) (bool, runtime.Object, error) {
		probes++
		job := action.(ktesting.CreateAction).GetObject().(*batchv1.Job)
		return true, job, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	detections, err := DetectNodeImages(ctx, []string{image}, []*corev1.Node{node})
	if err != nil {
		t.Fatal(err)
	}
	if probes != 1 {
		t.Fatalf("probes = %d, want one probe", probes)
	}
	if len(detections) != 1 || !slices.Equal(detections[0].Unverified, []string{image}) {
		t.Fatalf("detections = %v, want the probe to decide", detections)
	}
	if len(detections[0].Absent) != 0 {
		t.Fatalf("detections = %v, want no answer", detections)
	}
}

// The detection must probe every node, so one node never answers for another.
func TestDetectNodeImagesProbesEveryNode(t *testing.T) {
	image := "docker.0rays.club/jbnrz/linux-sjk:latest"
	client := useFakePods(t)
	t.Cleanup(Stop)
	for _, name := range []string{"worker-a", "worker-b"} {
		if err := client.Tracker().Add(truncatedNode(name)); err != nil {
			t.Fatal(err)
		}
	}
	targets := make([]string, 0)
	client.PrependReactor("create", "jobs", func(action ktesting.Action) (bool, runtime.Object, error) {
		job := action.(ktesting.CreateAction).GetObject().(*batchv1.Job)
		values := job.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchFields[0].Values
		targets = append(targets, values...)
		return true, job, nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	nodes := []*corev1.Node{truncatedNode("worker-a"), truncatedNode("worker-b")}
	if _, err := DetectNodeImages(ctx, []string{image}, nodes); err != nil {
		t.Fatal(err)
	}
	slices.Sort(targets)
	if !slices.Equal(targets, []string{"worker-a", "worker-b"}) {
		t.Fatalf("probe targets = %v, want every node", targets)
	}
}

// A detection that cannot probe at all must fail instead of reporting every
// image as absent.
func TestDetectNodeImagesFailsWithoutAnyProbe(t *testing.T) {
	client := useFakePods(t)
	t.Cleanup(Stop)
	client.PrependReactor("create", "jobs", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, fmt.Errorf("probe rejected")
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := DetectNodeImages(ctx, []string{"docker.0rays.club/jbnrz/linux-sjk:latest"}, []*corev1.Node{truncatedNode("worker-a")}); err == nil {
		t.Fatal("a detection without any probe must fail")
	}
}
