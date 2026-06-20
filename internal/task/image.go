package task

import (
	"CBCTF/internal/log"
	"CBCTF/internal/utils"
	"context"
	"fmt"
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
	return enqueueTask(resizeImageTaskType, task, asynq.MaxRetry(3), asynq.Timeout(time.Minute))
}

func HandleResizeImageTask(_ context.Context, task *asynq.Task) error {
	var payload ResizeImagePayload
	if err := msgpack.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	log.Logger.Debugf("Resizing image: path=%s width=%d height=%d", payload.Path, payload.Width, payload.Height)
	if err := utils.ResizePicture(payload.Path, payload.Width, payload.Height); err != nil {
		return fmt.Errorf("resize image %s failed: %w", payload.Path, err)
	}
	log.Logger.Debugf("Image resized: path=%s width=%d height=%d", payload.Path, payload.Width, payload.Height)
	return nil
}
