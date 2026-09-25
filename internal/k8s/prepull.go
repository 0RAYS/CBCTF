package k8s

import (
	"context"
	"crypto/sha256"
	"fmt"
	"slices"
	"time"

	"github.com/distribution/reference"
	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
)

func imageFailureKey(image string) string {
	return fmt.Sprintf("prepull:%s:%x", globalNamespace, sha256.Sum256([]byte(normalizedImage(image))))
}

func normalizedImage(image string) string {
	name, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return image
	}
	return reference.TagNameOnly(name).String()
}

// Only confirmed pull failures exclude a node. An absent/truncated Node image
// inventory is never evidence of failure. Success on a later warmup clears it.
func imageFailureAffinity(ctx context.Context, images []string) (*corev1.Affinity, error) {
	var excluded []string
	for _, image := range images {
		nodes, err := redis.RDB.HKeys(ctx, imageFailureKey(image)).Result()
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			if !slices.Contains(excluded, node) {
				excluded = append(excluded, node)
			}
		}
	}
	return excludedNodeAffinity(excluded), nil
}

func excludedNodeAffinity(excluded []string) *corev1.Affinity {
	if len(excluded) == 0 {
		return nil
	}
	slices.Sort(excluded)
	// MatchFields supports one value per requirement; requirements within a
	// term are ANDed, excluding every failed node without selecting a good one.
	fields := make([]corev1.NodeSelectorRequirement, 0, len(excluded))
	for _, node := range slices.Compact(excluded) {
		fields = append(fields, corev1.NodeSelectorRequirement{
			Key:      "metadata.name",
			Operator: corev1.NodeSelectorOpNotIn,
			Values:   []string{node},
		})
	}
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{{MatchFields: fields}},
			},
		},
	}
}

// PrepullImages submits all nodes before observing completion, so a slow node
// does not delay warming other nodes. ImageID, not command exit status, proves
// a successful pull (VM/distroless images need not contain echo or a shell).
func PrepullImages(ctx context.Context, images, selectedNodes []string, pullPolicy string) error {
	nodes, ret := ListSchedulableNodes(ctx)
	if !ret.OK {
		return resourceError(ret)
	}
	if len(nodes) == 0 {
		return fmt.Errorf("no schedulable nodes for image warmup")
	}
	for _, name := range selectedNodes {
		if !slices.ContainsFunc(nodes, func(n *corev1.Node) bool { return n.Name == name }) {
			return fmt.Errorf("node %s is no longer schedulable", name)
		}
	}
	batch := utils.RandHexStr(16)
	selector := map[string]string{"cbctf.io/prepull": batch}
	pending := make(map[string]bool)
	for _, node := range nodes {
		if len(selectedNodes) > 0 && !slices.Contains(selectedNodes, node.Name) {
			continue
		}
		var missing []string
		for _, image := range images {
			present := false
			if pullPolicy != string(corev1.PullAlways) {
				for _, entry := range node.Status.Images {
					for _, name := range entry.Names {
						if normalizedImage(name) == normalizedImage(image) {
							present = true
						}
					}
				}
			}
			if present {
				if err := redis.RDB.HDel(ctx, imageFailureKey(image), node.Name).Err(); err != nil {
					return err
				}
			} else {
				missing = append(missing, image)
			}
		}
		for i := 0; i < len(missing); i += 5 {
			chunk := missing[i:min(i+5, len(missing))]
			_, ret = CreateJob(ctx, CreateJobOptions{
				Name:         "prepull-" + utils.RandHexStr(20),
				Labels:       selector,
				Images:       chunk,
				PullPolicy:   pullPolicy,
				SelectedNode: node.Name,
			})
			if !ret.OK {
				return resourceError(ret)
			}
			for _, image := range chunk {
				pending[node.Name+"\x00"+image] = true
			}
		}
	}
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	failures := 0
	for len(pending) > 0 {
		pods, ret := ListPods(ctx, selector)
		if !ret.OK {
			return resourceError(ret)
		}
		for _, pod := range pods.Items {
			for _, status := range pod.Status.ContainerStatuses {
				image := ""
				for _, container := range pod.Spec.Containers {
					if container.Name == status.Name {
						image = container.Image
						break
					}
				}
				key := pod.Spec.NodeName + "\x00" + image
				if !pending[key] {
					continue
				}
				if status.ImageID != "" {
					if err := redis.RDB.HDel(ctx, imageFailureKey(image), pod.Spec.NodeName).Err(); err != nil {
						return err
					}
					delete(pending, key)
				} else if waiting := status.State.Waiting; waiting != nil &&
					(waiting.Reason == "ErrImagePull" ||
						waiting.Reason == "ImagePullBackOff" ||
						waiting.Reason == "InvalidImageName") {
					if err := redis.RDB.HSet(ctx, imageFailureKey(image), pod.Spec.NodeName, waiting.Reason).Err(); err != nil {
						return err
					}
					delete(pending, key)
					failures++
				}
			}
		}
		if len(pending) == 0 {
			if failures > 0 {
				return fmt.Errorf("image warmup failed for %d node/image pairs", failures)
			}
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("image warmup incomplete (%d node/image pairs): %w", len(pending), ctx.Err())
		case <-ticker.C:
		}
	}
	return nil
}

func containerImages(containers []corev1.Container) []string {
	images := make([]string, 0, len(containers))
	for _, container := range containers {
		if !slices.Contains(images, container.Image) {
			images = append(images, container.Image)
		}
	}
	return images
}

// ChallengeImages includes platform sidecars in the same warmup generation.
func ChallengeImages(challenge model.Challenge, sidecars ...string) []string {
	images := append([]string(nil), sidecars...)
	images = append(images, challenge.GeneratorImage)
	for _, pod := range challenge.Template.Pods {
		for _, container := range pod.Containers {
			images = append(images, container.Image)
		}
	}
	slices.Sort(images)
	return slices.DeleteFunc(slices.Compact(images), func(image string) bool { return image == "" })
}
