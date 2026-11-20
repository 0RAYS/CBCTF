package redis

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"

	"CBCTF/internal/model"
)

// imageDetectionTTL keeps the detection results of the common cache for a day.
const imageDetectionTTL = 24 * time.Hour

// imageDetectionKey addresses the common detection cache. Every detection run
// merges into it, so the system-wide view and the contest views read the same
// results.
const imageDetectionKey = "images:detect:common"

const (
	imageStatePresent    = "present"
	imageStateAbsent     = "absent"
	imageStateUnverified = "unverified"
)

// SaveImageDetection merges one detection run into the common cache. A run
// only covers the images it probed, so the state of the other images stays.
func SaveImageDetection(ctx context.Context, detection model.ImageDetection) error {
	existing, err := LoadImageDetection(ctx)
	if err != nil {
		return err
	}
	data, err := json.Marshal(mergeImageDetection(existing, detection))
	if err != nil {
		return err
	}
	return RDB.Set(ctx, imageDetectionKey, data, imageDetectionTTL).Err()
}

// LoadImageDetection returns the common detection result, if any.
func LoadImageDetection(ctx context.Context) (*model.ImageDetection, error) {
	data, err := RDB.Get(ctx, imageDetectionKey).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var detection model.ImageDetection
	if err := json.Unmarshal(data, &detection); err != nil {
		return nil, err
	}
	return &detection, nil
}

// mergeImageDetection combines the cached results with one detection run. The
// run answers for the nodes it probed, so the nodes it did not see are gone.
func mergeImageDetection(cached *model.ImageDetection, detection model.ImageDetection) model.ImageDetection {
	cachedNodes := make(map[string]model.NodeImageDetection)
	if cached != nil {
		for _, node := range cached.Nodes {
			cachedNodes[node.Node] = node
		}
	}
	merged := model.ImageDetection{
		DetectedAt: detection.DetectedAt,
		Nodes:      make([]model.NodeImageDetection, 0, len(detection.Nodes)),
	}
	for _, node := range detection.Nodes {
		merged.Nodes = append(merged.Nodes, mergeNodeDetection(cachedNodes[node.Node], node))
	}
	return merged
}

// mergeNodeDetection keeps the state of the images the run did not probe. A
// run that could not answer for an image leaves the known state alone.
func mergeNodeDetection(cached model.NodeImageDetection, node model.NodeImageDetection) model.NodeImageDetection {
	states := nodeImageStates(cached)
	for image, state := range nodeImageStates(node) {
		// A probe that never answered proves nothing, so a known state stays.
		if state == imageStateUnverified {
			if known, ok := states[image]; ok && known != imageStateUnverified {
				continue
			}
		}
		states[image] = state
	}
	merged := model.NodeImageDetection{Node: node.Node}
	for image, state := range states {
		switch state {
		case imageStatePresent:
			merged.Present = append(merged.Present, image)
		case imageStateAbsent:
			merged.Absent = append(merged.Absent, image)
		default:
			merged.Unverified = append(merged.Unverified, image)
		}
	}
	slices.Sort(merged.Present)
	slices.Sort(merged.Absent)
	slices.Sort(merged.Unverified)
	return merged
}

// nodeImageStates maps every image of one node to its detected state.
func nodeImageStates(node model.NodeImageDetection) map[string]string {
	states := make(map[string]string, len(node.Present)+len(node.Absent)+len(node.Unverified))
	for _, image := range node.Present {
		states[image] = imageStatePresent
	}
	for _, image := range node.Absent {
		states[image] = imageStateAbsent
	}
	for _, image := range node.Unverified {
		states[image] = imageStateUnverified
	}
	return states
}
