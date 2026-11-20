package task

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"CBCTF/internal/log"
	"CBCTF/internal/utils"
)

const resizeImageTaskType = "tasks:image"

type ResizeImagePayload struct {
	Path   string
	Width  int
	Height int
}

func EnqueueResizeImageTask(path string, width, height int) (*asynq.TaskInfo, error) {
	return enqueueMsgpackTask(resizeImageTaskType, ResizeImagePayload{path, width, height}, asynq.MaxRetry(3), asynq.Timeout(time.Minute))
}

func HandleResizeImageTask(_ context.Context, task *asynq.Task) error {
	payload, err := decodeMsgpackTask[ResizeImagePayload](task)
	if err != nil {
		return err
	}
	log.Logger.Debugf("Resizing image: path=%s width=%d height=%d", payload.Path, payload.Width, payload.Height)
	if err := utils.ResizePicture(payload.Path, payload.Width, payload.Height); err != nil {
		return fmt.Errorf("resize image %s failed: %w", payload.Path, err)
	}
	log.Logger.Debugf("Image resized: path=%s width=%d height=%d", payload.Path, payload.Width, payload.Height)
	return nil
}
