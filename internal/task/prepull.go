package task

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
	corev1 "k8s.io/api/core/v1"

	"CBCTF/internal/k8s"
)

const prepullTaskType = "tasks:prepull"

type PrepullPayload struct {
	Images     []string
	Nodes      []string
	PullPolicy string
}

func EnqueuePrepullTask(images []string) error {
	return EnqueuePrepullTargets(images, nil, "IfNotPresent")
}

func EnqueuePrepullTargets(images, nodes []string, pullPolicy string) error {
	if pullPolicy == string(corev1.PullNever) {
		return nil
	}
	images = slices.Clone(images)
	for i, image := range images {
		images[i] = k8s.NormalizeImage(image)
	}
	slices.Sort(images)
	images = slices.DeleteFunc(slices.Compact(images), func(s string) bool { return s == "" })
	if len(images) == 0 {
		return nil
	}
	nodes = slices.Clone(nodes)
	slices.Sort(nodes)
	payload, err := msgpack.Marshal(PrepullPayload{Images: images, Nodes: slices.Compact(nodes), PullPolicy: pullPolicy})
	if err != nil {
		return err
	}
	options := []asynq.Option{asynq.MaxRetry(2), asynq.Timeout(11 * time.Minute)}
	if pullPolicy != string(corev1.PullAlways) {
		options = append(options, asynq.Unique(15*time.Minute))
	}
	_, err = enqueueTask(prepullTaskType, asynq.NewTask(prepullTaskType, payload), options...)
	if errors.Is(err, asynq.ErrDuplicateTask) {
		return nil
	}
	return err
}

func HandlePrepullTask(ctx context.Context, t *asynq.Task) error {
	var payload PrepullPayload
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	err := k8s.PrepullImages(ctx, payload.Images, payload.Nodes, payload.PullPolicy)
	if _, ok := errors.AsType[*k8s.PartialPrepullError](err); ok {
		return errors.Join(err, asynq.SkipRetry)
	}
	return err
}
