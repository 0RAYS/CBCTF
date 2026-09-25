package k8s

import (
	"maps"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

const prepullLabel = "cbctf.io/prepull"

type pullTarget struct {
	node  string
	image string
}

type pullResult struct {
	pulled bool
	reason string
}

// Keep observed outcomes until the task consumes them. TTL-zero Jobs and their
// Pods can disappear from the informer store before the task reads that store.
type pullResults struct {
	batch   string
	mu      sync.Mutex
	results map[pullTarget]pullResult
	changed chan struct{}
}

func (p *pullResults) observe(obj any) {
	if deleted, ok := obj.(cache.DeletedFinalStateUnknown); ok {
		obj = deleted.Obj
	}
	pod, ok := obj.(*corev1.Pod)
	if !ok || pod == nil || pod.Labels[prepullLabel] != p.batch || pod.Spec.NodeName == "" {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	for _, status := range pod.Status.ContainerStatuses {
		result := pullResult{pulled: status.ImageID != ""}
		if !result.pulled {
			waiting := status.State.Waiting
			if waiting == nil {
				continue
			}
			switch waiting.Reason {
			case "ErrImagePull", "ImagePullBackOff", "InvalidImageName":
				result.reason = waiting.Reason
			default:
				continue
			}
		}
		for _, container := range pod.Spec.Containers {
			if container.Name != status.Name {
				continue
			}
			target := pullTarget{node: pod.Spec.NodeName, image: container.Image}
			// A later transient status must not undo proof of a successful pull.
			if !p.results[target].pulled {
				p.results[target] = result
			}
			select {
			case p.changed <- struct{}{}:
			default:
			}
			break
		}
	}
}

func (p *pullResults) snapshot() map[pullTarget]pullResult {
	p.mu.Lock()
	defer p.mu.Unlock()
	return maps.Clone(p.results)
}
