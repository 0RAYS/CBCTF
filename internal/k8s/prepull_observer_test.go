package k8s

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"

	"CBCTF/internal/redis"
)

func TestPullResultsRetainDeletedPodsAndIsolateBatches(t *testing.T) {
	observed := &pullResults{
		batch:   "batch-one",
		results: make(map[pullTarget]pullResult),
		changed: make(chan struct{}, 1),
	}
	pod := &corev1.Pod{
		Labels: map[string]string{prepullLabel: "batch-one"},
		Spec: corev1.PodSpec{
			NodeName:   "worker-a",
			Containers: []corev1.Container{{Name: "pull", Image: "example/image:v1"}},
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{{Name: "pull", ImageID: "sha256:ready"}},
		},
	}
	other := pod.DeepCopy()
	other.Labels[prepullLabel] = "batch-two"
	observed.observe(other)
	observed.observe((*corev1.Pod)(nil))
	if len(observed.snapshot()) != 0 {
		t.Fatal("unrelated Pod was observed")
	}
	observed.observe(cache.DeletedFinalStateUnknown{Obj: pod})
	target := pullTarget{node: "worker-a", image: "example/image:v1"}
	if !observed.snapshot()[target].pulled {
		t.Fatal("deleted Pod outcome was lost")
	}
	late := pod.DeepCopy()
	late.Status.ContainerStatuses[0].ImageID = ""
	late.Status.ContainerStatuses[0].State.Waiting = &corev1.ContainerStateWaiting{Reason: "ErrImagePull"}
	observed.observe(late)
	if !observed.snapshot()[target].pulled {
		t.Fatal("late status overwrote a successful pull")
	}
	snapshot := observed.snapshot()
	delete(snapshot, target)
	if len(observed.snapshot()) != 1 {
		t.Fatal("snapshot mutated retained outcomes")
	}
}

type pullRedisHook struct {
	commands []string
}

func (*pullRedisHook) DialHook(next goredis.DialHook) goredis.DialHook {
	return next
}

func (*pullRedisHook) ProcessPipelineHook(next goredis.ProcessPipelineHook) goredis.ProcessPipelineHook {
	return next
}

func (h *pullRedisHook) ProcessHook(_ goredis.ProcessHook) goredis.ProcessHook {
	return func(_ context.Context, cmd goredis.Cmder) error {
		if cmd.Name() != "hset" && cmd.Name() != "hdel" {
			return fmt.Errorf("unexpected Redis command: %s", cmd.Name())
		}
		h.commands = append(h.commands, cmd.Name())
		cmd.(*goredis.IntCmd).SetVal(1)
		return nil
	}
}

func TestPrepullObservesResultsBeforeImmediateTTLCollection(t *testing.T) {
	for _, tc := range []struct {
		name     string
		imageID  string
		reason   string
		exitCode int32
	}{
		{name: "successful Job", imageID: "sha256:ready"},
		{name: "failed command after successful pull", imageID: "sha256:ready", exitCode: 127},
		{name: "failed pull", reason: "ErrImagePull"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := useFakePods(t)
			t.Cleanup(Stop)
			previous := redis.RDB
			rdb := goredis.NewClient(&goredis.Options{})
			hook := &pullRedisHook{}
			rdb.AddHook(hook)
			redis.RDB = rdb
			t.Cleanup(func() {
				redis.RDB = previous
				_ = rdb.Close()
			})
			node := &corev1.Node{
				Name: "worker-a",
				Status: corev1.NodeStatus{
					Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
				},
			}
			if err := client.Tracker().Add(node); err != nil {
				t.Fatal(err)
			}
			// The fake tracker has no resource-version replay, unlike the API server.
			// Establish its watch before simulating the very short-lived Job Pod.
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
					return true, nil, fmt.Errorf("Pod watch was not established")
				}
				job := action.(ktesting.CreateAction).GetObject().(*batchv1.Job)
				if job.Spec.TTLSecondsAfterFinished == nil || *job.Spec.TTLSecondsAfterFinished != 0 {
					return true, nil, fmt.Errorf("Job is missing immediate TTL collection")
				}
				job.UID = "job-uid"
				pod := &corev1.Pod{
					Name:      job.Name + "-pod",
					Namespace: job.Namespace,
					UID:       "pod-uid",
					Labels:    job.Spec.Template.Labels,
					OwnerReferences: []metav1.OwnerReference{
						{APIVersion: "batch/v1", Kind: "Job", Name: job.Name, UID: job.UID},
					},
					Spec: job.Spec.Template.Spec,
				}
				pod.Spec.NodeName = node.Name
				if err := client.Tracker().Add(job); err != nil {
					return true, nil, err
				}
				if err := client.Tracker().Add(pod); err != nil {
					return true, nil, err
				}
				status := corev1.ContainerStatus{Name: pod.Spec.Containers[0].Name, ImageID: tc.imageID}
				if tc.reason != "" {
					status.State.Waiting = &corev1.ContainerStateWaiting{Reason: tc.reason}
				} else {
					status.State.Terminated = &corev1.ContainerStateTerminated{ExitCode: tc.exitCode}
				}
				pod.Status.ContainerStatuses = []corev1.ContainerStatus{status}
				pods := corev1.SchemeGroupVersion.WithResource("pods")
				if err := client.Tracker().Update(pods, pod, pod.Namespace); err != nil {
					return true, nil, err
				}
				// Simulate TTL-controller/GC removal before Create returns. No cleanup
				// requests are issued by the application; the fake client has no GC.
				if err := client.Tracker().Delete(pods, pod.Namespace, pod.Name); err != nil {
					return true, nil, err
				}
				if err := client.Tracker().Delete(action.GetResource(), job.Namespace, job.Name); err != nil {
					return true, nil, err
				}
				return true, job, nil
			})

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := PrepullImages(ctx, []string{"example/image:v1"}, nil, string(corev1.PullAlways))
			wantCommand := "hdel"
			if tc.reason != "" {
				wantCommand = "hset"
				if err == nil || !strings.Contains(err.Error(), "image warmup failed") {
					t.Fatalf("missing pull failure: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(hook.commands, []string{wantCommand}) {
				t.Fatalf("pull outcome not persisted: %v", hook.commands)
			}
			for _, action := range client.Actions() {
				if slices.Contains([]string{"delete", "deletecollection", "patch", "update"}, action.GetVerb()) {
					t.Fatalf("application manually changed cleanup state: %s %s", action.GetVerb(), action.GetResource())
				}
			}
		})
	}
}
