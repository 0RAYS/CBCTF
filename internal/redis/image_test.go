package redis

import (
	"reflect"
	"testing"
	"time"

	"CBCTF/internal/model"
)

func TestMergeImageDetectionKeepsTheImagesTheRunDidNotProbe(t *testing.T) {
	cached := &model.ImageDetection{
		DetectedAt: time.Now(),
		Nodes: []model.NodeImageDetection{
			{Node: "worker-a", Present: []string{"app:v1"}, Absent: []string{"db:v1"}},
		},
	}
	run := model.ImageDetection{
		DetectedAt: time.Now().Add(time.Minute),
		Nodes: []model.NodeImageDetection{
			{Node: "worker-a", Present: []string{"web:v1"}, Absent: []string{"db:v1"}},
		},
	}
	merged := mergeImageDetection(cached, run)
	want := []model.NodeImageDetection{
		{Node: "worker-a", Present: []string{"app:v1", "web:v1"}, Absent: []string{"db:v1"}},
	}
	if !reflect.DeepEqual(merged.Nodes, want) {
		t.Fatalf("merged nodes = %+v, want %+v", merged.Nodes, want)
	}
	if !merged.DetectedAt.Equal(run.DetectedAt) {
		t.Fatalf("detected at = %v, want %v", merged.DetectedAt, run.DetectedAt)
	}
}

// A probe that never answered proves nothing, so the known state stays.
func TestMergeImageDetectionKeepsAKnownStateOverAnUnverifiedOne(t *testing.T) {
	cached := &model.ImageDetection{
		Nodes: []model.NodeImageDetection{{Node: "worker-a", Present: []string{"app:v1"}}},
	}
	run := model.ImageDetection{
		Nodes: []model.NodeImageDetection{{Node: "worker-a", Unverified: []string{"app:v1"}}},
	}
	merged := mergeImageDetection(cached, run)
	want := []model.NodeImageDetection{{Node: "worker-a", Present: []string{"app:v1"}}}
	if !reflect.DeepEqual(merged.Nodes, want) {
		t.Fatalf("merged nodes = %+v, want %+v", merged.Nodes, want)
	}
}

// The run reports the nodes it probed, so a node it did not see is gone.
func TestMergeImageDetectionDropsTheNodesTheRunDidNotSee(t *testing.T) {
	cached := &model.ImageDetection{
		Nodes: []model.NodeImageDetection{
			{Node: "worker-a", Present: []string{"app:v1"}},
			{Node: "worker-b", Present: []string{"app:v1"}},
		},
	}
	run := model.ImageDetection{
		Nodes: []model.NodeImageDetection{{Node: "worker-a", Present: []string{"app:v1"}}},
	}
	merged := mergeImageDetection(cached, run)
	want := []model.NodeImageDetection{{Node: "worker-a", Present: []string{"app:v1"}}}
	if !reflect.DeepEqual(merged.Nodes, want) {
		t.Fatalf("merged nodes = %+v, want %+v", merged.Nodes, want)
	}
}
