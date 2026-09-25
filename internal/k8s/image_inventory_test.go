package k8s

import (
	"context"
	"slices"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestNodeImageInventoryMatchesWarmupIdentityAndRetainsDigests(t *testing.T) {
	client := useFakePods(t)
	digest := "ghcr.io/0rays/cbctf-worker@sha256:" + strings.Repeat("a", 64)
	node := &corev1.Node{
		Name: "worker-a",
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionTrue}},
			Images:     []corev1.ContainerImage{{Names: []string{"nginx", "docker.io/library/nginx:latest", digest, " "}}},
		},
	}
	if err := client.Tracker().Add(node); err != nil {
		t.Fatal(err)
	}
	images, ret := ListNodeImages(context.Background())
	if !ret.OK {
		t.Fatal(ret)
	}
	want := []string{"docker.io/library/nginx:latest", digest}
	if !slices.Equal(images[node.Name], want) {
		t.Fatalf("inventory = %v, want %v", images[node.Name], want)
	}
}
