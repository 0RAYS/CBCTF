package resp

import (
	"CBCTF/internal/model"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/vmihailenco/msgpack/v5"
)

func compactJSON(value any) any {
	if value == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var normalized any
	if err = json.Unmarshal(raw, &normalized); err != nil {
		return string(raw)
	}
	return normalized
}

func normalizePayload(value any) any {
	switch v := value.(type) {
	case map[any]any:
		data := make(map[string]any, len(v))
		for key, item := range v {
			data[func(value any) string {
				if str, ok := value.(string); ok {
					return str
				}
				raw, err := json.Marshal(value)
				if err != nil {
					return "unknown"
				}
				return string(raw)
			}(key)] = normalizePayload(item)
		}
		return data
	case map[string]any:
		data := make(map[string]any, len(v))
		for key, item := range v {
			data[key] = normalizePayload(item)
		}
		return data
	case []any:
		data := make([]any, 0, len(v))
		for _, item := range v {
			data = append(data, normalizePayload(item))
		}
		return data
	default:
		return value
	}
}

func GetTaskResp(task model.Task) gin.H {
	return gin.H{
		"id":           task.ID,
		"task_id":      task.TaskID,
		"type":         task.Type,
		"queue":        task.Queue,
		"status":       task.Status,
		"payload":      compactJSON(task.Payload.V),
		"result":       compactJSON(task.Result.V),
		"error":        task.Error,
		"retry_count":  task.RetryCount,
		"max_retry":    task.MaxRetry,
		"processed_at": task.ProcessedAt,
		"created_at":   task.CreatedAt,
	}
}

func GetLiveTaskResp(task *asynq.TaskInfo) gin.H {
	return gin.H{
		"task_id": task.ID,
		"type":    task.Type,
		"queue":   task.Queue,
		"status":  task.State.String(),
		"payload": func(raw []byte) any {
			if len(raw) == 0 {
				return nil
			}
			var normalized any
			if err := json.Unmarshal(raw, &normalized); err == nil {
				return normalized
			}
			if err := msgpack.Unmarshal(raw, &normalized); err == nil {
				return normalizePayload(normalized)
			}
			return string(raw)
		}(task.Payload),
		"error":           task.LastErr,
		"retry_count":     task.Retried,
		"max_retry":       task.MaxRetry,
		"last_failed_at":  zeroTimeNil(task.LastFailedAt),
		"next_process_at": zeroTimeNil(task.NextProcessAt),
		"completed_at":    zeroTimeNil(task.CompletedAt),
		"timeout":         int64(task.Timeout.Seconds()),
		"retention":       int64(task.Retention.Seconds()),
		"is_orphaned":     task.IsOrphaned,
	}
}

func zeroTimeNil(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}
