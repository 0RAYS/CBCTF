package task

import (
	"CBCTF/internal/k8s"
	"context"
	"errors"
	"slices"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
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
	images = slices.Clone(images)
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
	_, err = enqueueTask(prepullTaskType, asynq.NewTask(prepullTaskType, payload), asynq.Unique(15*time.Minute), asynq.MaxRetry(2), asynq.Timeout(11*time.Minute))
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
	return k8s.PrepullImages(ctx, payload.Images, payload.Nodes, payload.PullPolicy)
}
