package router

import (
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func TestNodeImagesReportTheLastDetectionPerNode(t *testing.T) {
	detection := &model.ImageDetection{
		DetectedAt: time.Now(),
		Nodes: []model.NodeImageDetection{
			{Node: "worker-b", Present: []string{"unrelated:v1"}},
			{Node: "worker-a", Present: []string{"nginx:latest", "unrelated:v1", "app:v1"}, Absent: []string{"db:v1"}},
		},
	}
	targets := []string{"app:v1", "nginx:latest", "db:v1"}
	want := []gin.H{
		{"node": "worker-a", "images": []string{"app:v1", "nginx:latest"}, "missing": []string{"db:v1"}},
		{"node": "worker-b", "images": []string{}, "missing": []string{}},
	}
	if got := formatNodeImages([]string{"worker-a", "worker-b"}, detection, targets); !reflect.DeepEqual(got, want) {
		t.Fatalf("node inventory = %v, want %v", got, want)
	}
}

// Without a detection the platform must not claim that an image is missing.
func TestNodeImagesWithoutDetectionReportNothing(t *testing.T) {
	targets := []string{"app:v1"}
	want := []gin.H{{"node": "worker-a", "images": []string{}, "missing": []string{}}}
	if got := formatNodeImages([]string{"worker-a"}, nil, targets); !reflect.DeepEqual(got, want) {
		t.Fatalf("node inventory = %v, want %v", got, want)
	}
}

// A node that joined after the detection stays listed, without an answer.
func TestNodeImagesKeepNodesWithoutADetection(t *testing.T) {
	detection := &model.ImageDetection{
		DetectedAt: time.Now(),
		Nodes:      []model.NodeImageDetection{{Node: "worker-a", Present: []string{"app:v1"}}},
	}
	want := []gin.H{
		{"node": "worker-a", "images": []string{"app:v1"}, "missing": []string{}},
		{"node": "worker-b", "images": []string{}, "missing": []string{}},
	}
	if got := formatNodeImages([]string{"worker-a", "worker-b"}, detection, []string{"app:v1"}); !reflect.DeepEqual(got, want) {
		t.Fatalf("node inventory = %v, want %v", got, want)
	}
}
