package task

import (
	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

// enqueueMsgpackTask encodes the payload and enqueues it on the queue of the
// task type.
func enqueueMsgpackTask[T any](taskType string, payload T, options ...asynq.Option) (*asynq.TaskInfo, error) {
	data, err := msgpack.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return enqueueTask(taskType, asynq.NewTask(taskType, data), options...)
}

// decodeMsgpackTask decodes the payload of a task.
func decodeMsgpackTask[T any](t *asynq.Task) (T, error) {
	var payload T
	if err := msgpack.Unmarshal(t.Payload(), &payload); err != nil {
		return payload, err
	}
	return payload, nil
}
