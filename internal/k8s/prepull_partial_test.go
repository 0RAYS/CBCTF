package k8s

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktesting "k8s.io/client-go/testing"
)

func TestWarmupJobFailureDoesNotSuppressOtherNodes(t *testing.T) {
	client := useFakePods(t)
	t.Cleanup(Stop)
	for _, name := range []string{"worker-a", "worker-b"} {
		if err := client.Tracker().Add(&corev1.Node{Name: name, Status: corev1.NodeStatus{Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}}}}); err != nil {
			t.Fatal(err)
		}
	}
	created := 0
	client.PrependReactor("create", "jobs", func(ktesting.Action) (bool, runtime.Object, error) {
		created++
		return true, nil, errors.New("job rejected")
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := PrepullImages(ctx, []string{"alpine:latest"}, nil, string(corev1.PullAlways))
	if err == nil || created != 2 || !strings.Contains(err.Error(), "worker-a") || !strings.Contains(err.Error(), "worker-b") {
		t.Fatalf("independent warmup targets lost: creates=%d error=%v", created, err)
	}
}
