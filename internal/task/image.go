package task

import (
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const resizeImageTaskType = "tasks:image"

type ResizeImagePayload struct {
	Path   string
	Width  int
	Height int
}

func EnqueueResizeImageTask(path string, width, height int) (*asynq.TaskInfo, error) {
	payload, err := msgpack.Marshal(ResizeImagePayload{path, width, height})
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(resizeImageTaskType, payload)
	info, err := client.Enqueue(task, asynq.Queue(resizeImageTaskType), asynq.MaxRetry(3), asynq.Timeout(time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(resizeImageTaskType)
	}
	return info, err
}

func HandleResizeImageTask(_ context.Context, task *asynq.Task) error {
	var payload ResizeImagePayload
	if err := msgpack.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return utils.ResizeImage(payload.Path, payload.Width, payload.Height)
}
