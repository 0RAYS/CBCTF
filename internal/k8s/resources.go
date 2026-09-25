package k8s

import (
	"CBCTF/internal/model"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func workloadRequests(limits corev1.ResourceList) corev1.ResourceList {
	requests := limits.DeepCopy()
	if requests == nil {
		requests = make(corev1.ResourceList)
	}
	if _, ok := requests[corev1.ResourceCPU]; !ok {
		requests[corev1.ResourceCPU] = resource.MustParse("100m")
	}
	if _, ok := requests[corev1.ResourceMemory]; !ok {
		requests[corev1.ResourceMemory] = resource.MustParse("64Mi")
	}
	return requests
}

func challengeFileMounts(containers []model.VictimContainerSpec) (map[string]string, [][]corev1.VolumeMount) {
	data := make(map[string]string)
	mounts := make([][]corev1.VolumeMount, len(containers))
	for i, container := range containers {
		for j, file := range container.VolumeMounts {
			key := fmt.Sprintf("c%d-f%d", i, j)
			data[key] = file.Content
			mounts[i] = append(mounts[i], corev1.VolumeMount{Name: "challenge-files", MountPath: file.Path, SubPath: key})
		}
	}
	return data, mounts
}

func sidecarResources() corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("10m"), corev1.ResourceMemory: resource.MustParse("32Mi")},
		Limits:   corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("500m"), corev1.ResourceMemory: resource.MustParse("256Mi")},
	}
}
