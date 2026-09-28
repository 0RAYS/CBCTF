package k8s

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktesting "k8s.io/client-go/testing"

	"CBCTF/internal/redis"
)

func TestPrepullPoliciesCompareEachTargetNodeInventory(t *testing.T) {
	digest := "example/app@sha256:" + strings.Repeat("a", 64)
	for _, tc := range []struct {
		name       string
		policy     string
		images     []string
		present    []string
		wantPull   []string
		wantClears int
	}{
		{name: "present aliases and digests", policy: "IfNotPresent", images: []string{"nginx", digest}, present: []string{" docker.io/library/nginx:latest ", digest}, wantClears: 2},
		{name: "only missing images", policy: "IfNotPresent", images: []string{"nginx", "app:v2"}, present: []string{"docker.io/library/nginx:latest", "app:v1"}, wantPull: []string{"app:v2"}, wantClears: 1},
		{name: "another node cannot satisfy target", policy: "IfNotPresent", images: []string{"nginx"}, wantPull: []string{"nginx"}},
		{name: "always includes present images", policy: "Always", images: []string{"nginx", "app:v2"}, present: []string{"nginx"}, wantPull: []string{"nginx", "app:v2"}},
		{name: "never skips missing images", policy: "Never", images: []string{"nginx"}},
		{name: "never skips present images", policy: "Never", images: []string{"nginx"}, present: []string{"nginx"}},
		{name: "default skips present images", images: []string{"nginx"}, present: []string{"nginx"}, wantClears: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := useFakePods(t)
			t.Cleanup(Stop)
			previous := redis.RDB
			rdb := goredis.NewClient(&goredis.Options{})
			hook := &pullRedisHook{}
			rdb.AddHook(hook)
			redis.RDB = rdb
			t.Cleanup(func() { redis.RDB = previous; _ = rdb.Close() })
			for name, images := range map[string][]string{"target": tc.present, "other": {"nginx", "app:v2"}} {
				if err := client.Tracker().Add(&corev1.Node{
					Name: name,
					Status: corev1.NodeStatus{
						Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
						Images:     []corev1.ContainerImage{{Names: images}},
					},
				}); err != nil {
					t.Fatal(err)
				}
			}
			var pulled []string
			client.PrependReactor("create", "jobs", func(action ktesting.Action) (bool, runtime.Object, error) {
				job := action.(ktesting.CreateAction).GetObject().(*batchv1.Job)
				for _, container := range job.Spec.Template.Spec.Containers {
					pulled = append(pulled, container.Image)
					if string(container.ImagePullPolicy) != tc.policy {
						t.Errorf("Job pull policy = %s, want %s", container.ImagePullPolicy, tc.policy)
					}
				}
				terms := job.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
				if !slices.Equal(terms[0].MatchFields[0].Values, []string{"target"}) {
					t.Errorf("Job targeted the wrong node: %v", terms)
				}
				// Stop after inspecting the submitted Job; completion is covered by
				// the observer tests without requiring a live Kubernetes controller.
				return true, nil, fmt.Errorf("test job submission observed")
			})
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err := PrepullImages(ctx, tc.images, []string{"target"}, tc.policy)
			if len(tc.wantPull) == 0 && err != nil {
				t.Fatal(err)
			}
			if len(tc.wantPull) > 0 && (err == nil || !strings.Contains(err.Error(), "test job submission observed")) {
				t.Fatalf("expected Job submission, got %v", err)
			}
			if !slices.Equal(pulled, tc.wantPull) || len(hook.commands) != tc.wantClears {
				t.Fatalf("pulled=%v, cleared=%v; want pulls=%v, clears=%d", pulled, hook.commands, tc.wantPull, tc.wantClears)
			}
			if tc.policy == "Never" && len(client.Actions()) != 0 {
				t.Fatalf("Never contacted Kubernetes: %v", client.Actions())
			}
		})
	}
}
