package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

const logKey = "logs"

// LogHook is a logrus Hook that mirrors log lines to a Redis list.
type LogHook struct {
	// Key is the Redis list key to push logs to.
	Key string
	// Max is the maximum number of log entries to retain (trimmed from the left).
	Max int
	// Formatter formats the entry to bytes before pushing to Redis.
	Formatter logrus.Formatter
}

// NewLogHook creates a new LogHook with sane defaults.
func NewLogHook(max int, formatter logrus.Formatter) *LogHook {
	if max <= 0 {
		max = 1000
	}
	return &LogHook{Key: logKey, Max: max, Formatter: formatter}
}

// Levels returns all log levels to mirror everything.
func (h *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire formats and pushes the log entry to Redis.
func (h *LogHook) Fire(entry *logrus.Entry) error {
	// If Redis is not initialized, skip silently.
	if RDB == nil {
		return nil
	}
	formatter := h.Formatter
	if formatter == nil {
		formatter = &logrus.TextFormatter{DisableTimestamp: true, DisableLevelTruncation: true}
	}
	// Use a shallow copy of the entry to avoid data races when formatting.
	cpy := *entry
	line, err := formatter.Format(&cpy)
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Best-effort push and trim; ignore errors to not block primary logging path.
	_ = RDB.RPush(ctx, h.Key, line).Err()
	_ = RDB.LTrim(ctx, h.Key, int64(-h.Max), -1).Err()
	return nil
}

func GetLogs(start, end int64) ([]string, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	logs, err := RDB.LRange(ctx, logKey, start, end).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get logs: %s", err)
		return logs, false, i18n.RedisError
	}
	return logs, true, i18n.Success
}
