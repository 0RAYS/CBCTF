package router

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNodeImagesOnlyReportRequiredImagesAndRetainEmptyNodes(t *testing.T) {
	nodes := map[string][]string{
		"worker-b": {"unrelated:v1"},
		"worker-a": {"nginx:latest", "unrelated:v1", "app:v1"},
	}
	targets := []string{"app:v1", "nginx:latest", "absent-everywhere:v1"}
	want := []gin.H{
		{"node": "worker-a", "images": []string{"app:v1", "nginx:latest"}},
		{"node": "worker-b", "images": []string{}},
	}
	if got := formatNodeImages(nodes, targets); !reflect.DeepEqual(got, want) {
		t.Fatalf("node inventory = %v, want %v", got, want)
	}
	for _, node := range formatNodeImages(nodes, nil) {
		if len(node["images"].([]string)) != 0 {
			t.Fatalf("empty target list included unrelated images: %v", node)
		}
	}
}
