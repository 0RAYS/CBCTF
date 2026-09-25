package k8s

import (
	"CBCTF/internal/model"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestMergedFilesPreserveSameBasenameAndContainerPaths(t *testing.T) {
	containers := []model.VictimContainerSpec{
		{VolumeMounts: model.XVolumes{{Path: "/root/flag", Content: "first"}, {Path: "/app/flag", Content: "second"}}},
		{VolumeMounts: model.XVolumes{{Path: "/root/flag", Content: "third"}}},
	}
	data, mounts := challengeFileMounts(containers)
	if len(data) != 3 {
		t.Fatal("same-basename files were merged")
	}
	for i, container := range containers {
		for j, file := range container.VolumeMounts {
			mount := mounts[i][j]
			if mount.MountPath != file.Path || data[mount.SubPath] != file.Content || mount.Name != "challenge-files" {
				t.Fatalf("mount changed: %+v", mount)
			}
		}
	}
}

func TestRequestsNeverExceedExplicitLimits(t *testing.T) {
	limits := corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("5m"), corev1.ResourceMemory: resource.MustParse("16Mi")}
	requests := workloadRequests(limits)
	for name, limit := range limits {
		request := requests[name]
		if request.Cmp(limit) > 0 {
			t.Fatalf("%s request exceeds limit", name)
		}
	}
	defaults := workloadRequests(nil)
	if len(defaults) != 2 {
		t.Fatal("unspecified workloads have no scheduling requests")
	}
}
