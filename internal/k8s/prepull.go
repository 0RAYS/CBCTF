package k8s

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/distribution/reference"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"

	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
)

func imageFailureKey(image string) string {
	return fmt.Sprintf("prepull:%s:%x", globalNamespace, sha256.Sum256([]byte(NormalizeImage(image))))
}

// PartialPrepullError Some Jobs were already submitted. Retrying the original batch (especially
// Always) would repeat successful side effects; callers should target failures.
type PartialPrepullError struct{ Err error }

func (e *PartialPrepullError) Error() string { return e.Err.Error() }
func (e *PartialPrepullError) Unwrap() error { return e.Err }

// NormalizeImage uses one identity for warmup, failure tracking and image inventories.
func NormalizeImage(image string) string {
	image = strings.TrimSpace(image)
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
	if pullPolicy == string(corev1.PullNever) || len(images) == 0 {
		return nil
	}
	if pullPolicy == "" {
		pullPolicy = string(corev1.PullIfNotPresent)
	}
	if pullPolicy != string(corev1.PullIfNotPresent) && pullPolicy != string(corev1.PullAlways) {
		return fmt.Errorf("unsupported image pull policy: %s", pullPolicy)
	}
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
	selector := map[string]string{prepullLabel: batch}
	pods, err := cachedObjects(ctx, "pods")
	if err != nil {
		return err
	}
	observed := &pullResults{
		batch:   batch,
		results: make(map[pullTarget]pullResult),
		changed: make(chan struct{}, 1),
	}
	// Subscribe before creating Jobs to retain outcomes when Pods leave the
	// cache. This handler shares the namespace's existing watch.
	handler, err := pods.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    observed.observe,
		UpdateFunc: func(_, current any) { observed.observe(current) },
		DeleteFunc: observed.observe,
	})
	if err != nil {
		return err
	}
	defer func() { _ = pods.informer.RemoveEventHandler(handler) }()
	if !cache.WaitForCacheSync(ctx.Done(), handler.HasSynced) {
		return fmt.Errorf("sync image warmup observer: %w", ctx.Err())
	}
	pending := make(map[pullTarget]bool)
	var failures []error
	submitted := 0
	finish := func() error {
		err := errors.Join(failures...)
		if err != nil {
			err = fmt.Errorf("image warmup failed: %w", err)
		}
		if err != nil && submitted > 0 {
			return &PartialPrepullError{Err: err}
		}
		return err
	}
	for _, node := range nodes {
		if ctx.Err() != nil {
			failures = append(failures, ctx.Err())
			return finish()
		}
		if len(selectedNodes) > 0 && !slices.Contains(selectedNodes, node.Name) {
			continue
		}
		var missing []string
		inventory := nodeImageSet(node)
		for _, image := range images {
			_, present := inventory[NormalizeImage(image)]
			if pullPolicy == string(corev1.PullIfNotPresent) && present {
				if err := redis.RDB.HDel(ctx, imageFailureKey(image), node.Name).Err(); err != nil {
					failures = append(failures, fmt.Errorf("clear image status %s/%s: %w", node.Name, image, err))
				}
			} else {
				missing = append(missing, image)
			}
		}
		for i := 0; i < len(missing); i += 5 {
			if ctx.Err() != nil {
				failures = append(failures, ctx.Err())
				return finish()
			}
			chunk := missing[i:min(i+5, len(missing))]
			_, ret = CreateJob(ctx, CreateJobOptions{
				Name:         "prepull-" + utils.RandHexStr(20),
				Labels:       selector,
				Images:       chunk,
				PullPolicy:   pullPolicy,
				SelectedNode: node.Name,
			})
			if !ret.OK {
				failures = append(failures, fmt.Errorf("create warmup job node=%s images=%v: %w", node.Name, chunk, resourceError(ret)))
				continue
			}
			submitted++
			for _, image := range chunk {
				pending[pullTarget{node: node.Name, image: image}] = true
			}
		}
	}
	for len(pending) > 0 {
		for target, result := range observed.snapshot() {
			if !pending[target] {
				continue
			}
			if result.pulled {
				if err := redis.RDB.HDel(ctx, imageFailureKey(target.image), target.node).Err(); err != nil {
					failures = append(failures, fmt.Errorf("clear image status %s/%s: %w", target.node, target.image, err))
				}
			} else {
				if err := redis.RDB.HSet(ctx, imageFailureKey(target.image), target.node, result.reason).Err(); err != nil {
					failures = append(failures, fmt.Errorf("save image status %s/%s: %w", target.node, target.image, err))
				}
				failures = append(failures, fmt.Errorf("pull image %s/%s: %s", target.node, target.image, result.reason))
			}
			delete(pending, target)
		}
		if len(pending) == 0 {
			return finish()
		}
		select {
		case <-ctx.Done():
			failures = append(failures, fmt.Errorf("image warmup incomplete (%d node/image pairs): %w", len(pending), ctx.Err()))
			return finish()
		case <-observed.changed:
		}
	}
	return finish()
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
	for i, image := range images {
		images[i] = NormalizeImage(image)
	}
	slices.Sort(images)
	return slices.DeleteFunc(slices.Compact(images), func(image string) bool { return image == "" })
}
