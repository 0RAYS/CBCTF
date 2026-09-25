package k8s

import (
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"context"
	"crypto/sha256"
	"fmt"
	"slices"
	"time"

	corev1 "k8s.io/api/core/v1"
)

func imageFailureKey(image string) string {
	return fmt.Sprintf("prepull:%s:%x", globalNamespace, sha256.Sum256([]byte(image)))
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
	if len(excluded) == 0 {
		return nil, nil
	}
	slices.Sort(excluded)
	return &corev1.Affinity{NodeAffinity: &corev1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
			NodeSelectorTerms: []corev1.NodeSelectorTerm{{MatchFields: []corev1.NodeSelectorRequirement{{
				Key: "metadata.name", Operator: corev1.NodeSelectorOpNotIn, Values: excluded,
			}}}},
		},
	}}, nil
}

// PrepullImages submits all nodes before observing completion, so a slow node
// does not delay warming other nodes. ImageID, not command exit status, proves
// a successful pull (VM/distroless images need not contain echo or a shell).
func PrepullImages(ctx context.Context, images []string) error {
	nodes, ret := ListSchedulableNodes(ctx)
	if !ret.OK {
		return resourceError(ret)
	}
	if len(nodes) == 0 {
		return fmt.Errorf("no schedulable nodes for image warmup")
	}
	batch := utils.RandHexStr(16)
	selector := map[string]string{"cbctf.io/prepull": batch}
	pending := make(map[string]bool)
	for _, node := range nodes {
		for i := 0; i < len(images); i += 5 {
			chunk := images[i:min(i+5, len(images))]
			_, ret = CreateJob(ctx, CreateJobOptions{
				Name: "prepull-" + utils.RandHexStr(20), Labels: selector,
				Images: chunk, PullPolicy: string(corev1.PullIfNotPresent), SelectedNode: node.Name,
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
				} else if waiting := status.State.Waiting; waiting != nil && (waiting.Reason == "ErrImagePull" || waiting.Reason == "ImagePullBackOff" || waiting.Reason == "InvalidImageName") {
					if err := redis.RDB.HSet(ctx, imageFailureKey(image), pod.Spec.NodeName, waiting.Reason).Err(); err != nil {
						return err
					}
				}
			}
		}
		if len(pending) == 0 {
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
