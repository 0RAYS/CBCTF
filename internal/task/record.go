package task

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

func recordTaskExecution(ctx context.Context, t *asynq.Task, status string, result any, err error) {
	if db.DB == nil {
		return
	}
	taskID, _ := asynq.GetTaskID(ctx)
	queue, _ := asynq.GetQueueName(ctx)
	retryCount, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if taskID == "" && t.ResultWriter() != nil {
		taskID = t.ResultWriter().TaskID()
	}
	if queue == "" {
		queue = t.Type()
	}
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	if _, ret := db.InitTaskRepo(db.DB).Create(db.CreateTaskOptions{
		TaskID:      taskID,
		Type:        t.Type(),
		Queue:       queue,
		Status:      status,
		Payload:     decodeTaskPayload(t.Payload()),
		Result:      result,
		Error:       errorMsg,
		RetryCount:  retryCount,
		MaxRetry:    maxRetry,
		ProcessedAt: time.Now(),
	}); !ret.OK {
		log.Logger.Warningf("Failed to record task %s history: %s", t.Type(), ret.Msg)
	}
}

func shouldRecordFailure(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, asynq.SkipRetry) || errors.Is(err, asynq.RevokeTask) {
		return true
	}
	retried, okRetry := asynq.GetRetryCount(ctx)
	maxRetry, okMax := asynq.GetMaxRetry(ctx)
	if !okRetry || !okMax {
		return true
	}
	return retried >= maxRetry
}

func decodeTaskPayload(payload []byte) any {
	if len(payload) == 0 {
		return nil
	}
	var normalized any
	if err := json.Unmarshal(payload, &normalized); err == nil {
		return normalized
	}
	if err := msgpack.Unmarshal(payload, &normalized); err == nil {
		return normalizeMsgpack(normalized)
	}
	return map[string]any{
		"encoding": "base64",
		"data":     base64.StdEncoding.EncodeToString(payload),
	}
}

func normalizeMsgpack(value any) any {
	switch v := value.(type) {
	case map[any]any:
		data := make(map[string]any, len(v))
		for key, item := range v {
			data[fmt.Sprint(key)] = normalizeMsgpack(item)
		}
		return data
	case map[string]any:
		data := make(map[string]any, len(v))
		for key, item := range v {
			data[key] = normalizeMsgpack(item)
		}
		return data
	case []any:
		data := make([]any, 0, len(v))
		for _, item := range v {
			data = append(data, normalizeMsgpack(item))
		}
		return data
	default:
		return value
	}
}
