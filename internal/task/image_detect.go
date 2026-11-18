package task

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"

	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
)

const imageDetectTaskType = "tasks:detect"

// ImageDetectPayload carries the scope and the images of one detection run.
type ImageDetectPayload struct {
	Scope  string   `msgpack:"scope"`
	Images []string `msgpack:"images"`
}

// EnqueueImageDetect starts a detection run for one scope.
func EnqueueImageDetect(scope string, images []string) error {
	images = slices.Clone(images)
	for i := range images {
		images[i] = k8s.NormalizeImage(images[i])
	}
	slices.Sort(images)
	images = slices.DeleteFunc(slices.Compact(images), func(image string) bool { return image == "" })
	if len(images) == 0 {
		return nil
	}
	payload, err := msgpack.Marshal(ImageDetectPayload{Scope: scope, Images: images})
	if err != nil {
		return err
	}
	options := []asynq.Option{asynq.MaxRetry(1), asynq.Timeout(10 * time.Minute)}
	_, err = enqueueTask(imageDetectTaskType, asynq.NewTask(imageDetectTaskType, payload), options...)
	if errors.Is(err, asynq.ErrDuplicateTask) {
		return nil
	}
	return err
}

// HandleImageDetectTask probes the nodes and keeps the last result only.
func HandleImageDetectTask(ctx context.Context, t *asynq.Task) error {
	var payload ImageDetectPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 9*time.Minute)
	defer cancel()
	nodes, ret := k8s.ListSchedulableNodes(ctx)
	if !ret.OK {
		return fmt.Errorf("%s: %v", ret.Msg, ret.Attr)
	}
	detections, err := k8s.DetectNodeImages(ctx, payload.Images, nodes)
	if err != nil {
		return err
	}
	return redis.SaveImageDetection(ctx, payload.Scope, model.ImageDetection{
		DetectedAt: time.Now(),
		Nodes:      detections,
	})
}
