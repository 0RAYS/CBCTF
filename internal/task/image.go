package task

import (
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/utils"
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

const ResizeImageTaskType = "tasks:image"

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
	task := asynq.NewTask(ResizeImageTaskType, payload)
	info, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(time.Minute))
	if err == nil {
		prometheus.RecordTaskEnqueued(ResizeImageTaskType)
	}
	return info, err
}

func HandleResizeImageTask(_ context.Context, task *asynq.Task) error {
	var payload ResizeImagePayload
	if err := msgpack.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	if err := utils.ResizeImage(payload.Path, payload.Width, payload.Height); err != nil {
		log.Logger.Warningf("Failed to resize image: %v", err)
		return err
	}
	return nil
}
