package task

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/hibiken/asynq"

	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
)

const imageDetectTaskType = "tasks:detect"

// imageDetectUniqueTTL keeps a pending detection from being queued twice. It
// outlasts the detection timeout, and the scheduled runs are far enough apart
// to queue the same images again.
const imageDetectUniqueTTL = 15 * time.Minute

// ImageDetectPayload carries the images of one detection run.
type ImageDetectPayload struct {
	Images []string `msgpack:"images"`
}

// EnqueueImageDetect starts a detection run for the given images.
func EnqueueImageDetect(images []string) error {
	images = slices.Clone(images)
	for i := range images {
		images[i] = k8s.NormalizeImage(images[i])
	}
	slices.Sort(images)
	images = slices.DeleteFunc(slices.Compact(images), func(image string) bool { return image == "" })
	if len(images) == 0 {
		return nil
	}
	options := []asynq.Option{asynq.MaxRetry(1), asynq.Timeout(10 * time.Minute), asynq.Unique(imageDetectUniqueTTL)}
	_, err := enqueueMsgpackTask(imageDetectTaskType, ImageDetectPayload{Images: images}, options...)
	if errors.Is(err, asynq.ErrDuplicateTask) {
		return nil
	}
	return err
}

// HandleImageDetectTask probes the nodes and merges the result into the
// common detection cache.
func HandleImageDetectTask(ctx context.Context, t *asynq.Task) error {
	payload, err := decodeMsgpackTask[ImageDetectPayload](t)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 9*time.Minute)
	defer cancel()
	nodes, ret := k8s.ListSchedulableNodes(ctx)
	if !ret.OK {
		return fmt.Errorf("%w", ret.AsError())
	}
	detections, err := k8s.DetectNodeImages(ctx, payload.Images, nodes)
	if err != nil {
		return err
	}
	return redis.SaveImageDetection(ctx, model.ImageDetection{
		DetectedAt: time.Now(),
		Nodes:      detections,
	})
}
