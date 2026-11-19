package k8s

import (
	"context"
	"fmt"
	"slices"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"

	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

const (
	// detectLabel marks the probe Jobs of one detection run.
	detectLabel = "cbctf.io/detect"
	// detectionChunkSize keeps one probe Job within the cluster container limit.
	detectionChunkSize = 5
	// nodeImageReportLimit is the number of images the kubelet reports in
	// node.status.images. A node below the limit reports every image it holds,
	// a node at the limit only reports its largest images.
	nodeImageReportLimit = 50
)

// DetectNodeImages determines the images of every node from three sources:
// the image list the node reports, the images the node already runs and a
// probe with imagePullPolicy: Never.
//
// A node that holds few images reports all of them, so its list decides the
// result without a probe. The kubelet truncates the list of a busy node to its
// largest images, so there the list and the running Pods only prove presence,
// while a probe proves the rest: a probe container only starts when the node
// already holds the image, so a started container proves presence and
// ErrImageNeverPull proves absence. The probe never contacts a registry, so it
// never changes the node.
func DetectNodeImages(ctx context.Context, images []string, nodes []*corev1.Node) ([]model.NodeImageDetection, error) {
	probes := make([]string, 0, len(images))
	for _, image := range images {
		if image = NormalizeImage(image); image != "" {
			probes = append(probes, image)
		}
	}
	slices.Sort(probes)
	probes = slices.Compact(probes)
	detections := make([]model.NodeImageDetection, 0, len(nodes))
	for _, node := range nodes {
		detections = append(detections, model.NodeImageDetection{
			Node:       node.Name,
			Present:    []string{},
			Absent:     []string{},
			Unverified: []string{},
		})
	}
	if len(probes) == 0 || len(nodes) == 0 {
		return detections, nil
	}

	pods, err := cachedObjects(ctx, "pods")
	if err != nil {
		return nil, err
	}
	running := podImagesByNode(pods)
	present := make(map[string]map[string]struct{}, len(nodes))
	absent := make(map[string]map[string]struct{}, len(nodes))
	pending := make(map[pullTarget]bool)
	probesByNode := make(map[string][]string)
	for _, node := range nodes {
		held := nodeImageSet(node)
		for image := range running[node.Name] {
			held[image] = struct{}{}
		}
		present[node.Name] = held
		unresolved := make([]string, 0, len(probes))
		for _, image := range probes {
			if _, ok := held[image]; !ok {
				unresolved = append(unresolved, image)
			}
		}
		if len(node.Status.Images) > 0 && len(node.Status.Images) < nodeImageReportLimit {
			// The node holds few images, so its list reports every one of
			// them and decides the result without a probe. An empty list
			// proves nothing, so such a node still gets a probe.
			missing := make(map[string]struct{}, len(unresolved))
			for _, image := range unresolved {
				missing[image] = struct{}{}
			}
			absent[node.Name] = missing
			continue
		}
		probesByNode[node.Name] = unresolved
		for _, image := range unresolved {
			pending[pullTarget{node: node.Name, image: image}] = true
		}
	}
	if len(pending) == 0 {
		return buildDetections(detections, probes, present, absent, nil), nil
	}

	batch := utils.RandHexStr(16)
	observed := &pullResults{
		batch:   batch,
		label:   detectLabel,
		results: make(map[pullTarget]pullResult),
		changed: make(chan struct{}, 1),
	}
	// Subscribe before creating the probe Jobs so outcomes survive Pods leaving the cache.
	handler, err := pods.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    observed.observe,
		UpdateFunc: func(_, current any) { observed.observe(current) },
		DeleteFunc: observed.observe,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = pods.informer.RemoveEventHandler(handler) }()
	if !cache.WaitForCacheSync(ctx.Done(), handler.HasSynced) {
		return nil, fmt.Errorf("sync image detection observer: %w", ctx.Err())
	}

	submitted := 0
	for _, node := range nodes {
		unresolved := probesByNode[node.Name]
		for i := 0; i < len(unresolved); i += detectionChunkSize {
			chunk := unresolved[i:min(i+detectionChunkSize, len(unresolved))]
			_, ret := CreateJob(ctx, CreateJobOptions{
				Name:               "detect-" + utils.RandHexStr(20),
				Labels:             map[string]string{detectLabel: batch},
				Images:             chunk,
				PullPolicy:         string(corev1.PullNever),
				SelectedNode:       node.Name,
				ServiceAccountName: runtimeServiceAccount(),
			})
			if !ret.OK {
				// A failed probe leaves its images unverified instead of
				// answering for the whole node.
				log.Logger.Warningf("Failed to create the image detection probe: node=%s images=%v error=%v", node.Name, chunk, ret)
				for _, image := range chunk {
					delete(pending, pullTarget{node: node.Name, image: image})
				}
				continue
			}
			submitted++
		}
	}
	if submitted == 0 {
		return nil, fmt.Errorf("failed to create any image detection probe")
	}

	results := observed.snapshot
	for len(pending) > 0 {
		for target := range results() {
			if !pending[target] {
				continue
			}
			delete(pending, target)
		}
		if len(pending) == 0 {
			break
		}
		select {
		case <-ctx.Done():
			// The probes that never answered stay unverified.
			log.Logger.Warningf("Image detection ended with %d unanswered probes: %s", len(pending), ctx.Err())
			pending = nil
		case <-observed.changed:
		}
	}
	return buildDetections(detections, probes, present, absent, results()), nil
}

// buildDetections merges the three sources into one result per node.
func buildDetections(detections []model.NodeImageDetection, probes []string, present map[string]map[string]struct{}, absent map[string]map[string]struct{}, answers map[pullTarget]pullResult) []model.NodeImageDetection {
	for i := range detections {
		node := detections[i].Node
		for _, image := range probes {
			if _, ok := present[node][image]; ok {
				detections[i].Present = append(detections[i].Present, image)
				continue
			}
			if _, ok := absent[node][image]; ok {
				detections[i].Absent = append(detections[i].Absent, image)
				continue
			}
			result, ok := answers[pullTarget{node: node, image: image}]
			switch {
			case ok && result.pulled:
				detections[i].Present = append(detections[i].Present, image)
			case ok && result.absent:
				detections[i].Absent = append(detections[i].Absent, image)
			default:
				detections[i].Unverified = append(detections[i].Unverified, image)
			}
		}
	}
	return detections
}

// podImagesByNode collects the images the kubelet already resolved on every
// node. A container status carrying an image ID proves the node holds it,
// while a pending Pod proves nothing.
func podImagesByNode(pods *objectCache) map[string]map[string]struct{} {
	images := make(map[string]map[string]struct{})
	for _, obj := range pods.informer.GetStore().List() {
		pod, ok := obj.(*corev1.Pod)
		if !ok || pod.Spec.NodeName == "" {
			continue
		}
		node, ok := images[pod.Spec.NodeName]
		if !ok {
			node = make(map[string]struct{})
			images[pod.Spec.NodeName] = node
		}
		for _, status := range pod.Status.ContainerStatuses {
			if status.ImageID == "" {
				continue
			}
			if image := containerImage(pod.Spec.Containers, status.Name); image != "" {
				node[NormalizeImage(image)] = struct{}{}
			}
		}
		for _, status := range pod.Status.InitContainerStatuses {
			if status.ImageID == "" {
				continue
			}
			if image := containerImage(pod.Spec.InitContainers, status.Name); image != "" {
				node[NormalizeImage(image)] = struct{}{}
			}
		}
	}
	return images
}

func containerImage(containers []corev1.Container, name string) string {
	for _, container := range containers {
		if container.Name == name {
			return container.Image
		}
	}
	return ""
}
