package k8s

import (
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestWarmupOnlyExcludesFailedNodes(t *testing.T) {
	if excludedNodeAffinity(nil) != nil {
		t.Fatal("healthy workloads must not carry affinity")
	}
	affinity := excludedNodeAffinity([]string{"worker-b", "worker-a"})
	node := affinity.NodeAffinity
	if len(node.PreferredDuringSchedulingIgnoredDuringExecution) != 0 {
		t.Fatal("unexpected positive preference")
	}
	terms := node.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	if len(terms) != 1 || len(terms[0].MatchExpressions) != 0 || len(terms[0].MatchFields) != 1 {
		t.Fatalf("unexpected exclusion: %+v", terms)
	}
	rule := terms[0].MatchFields[0]
	if rule.Key != "metadata.name" || rule.Operator != corev1.NodeSelectorOpNotIn || len(rule.Values) != 2 {
		t.Fatalf("not an exclusion: %+v", rule)
	}
	if imageFailureKey("nginx") != imageFailureKey("docker.io/library/nginx:latest") {
		t.Fatal("image aliases do not share failure state")
	}
}
