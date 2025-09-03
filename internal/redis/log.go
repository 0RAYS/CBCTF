package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

const logKey = "logs"

type LogHook struct {
	Key       string
	Max       int
	Formatter logrus.Formatter
}

func NewLogHook(max int, formatter logrus.Formatter) *LogHook {
	if max <= 0 {
		max = 1000
	}
	return &LogHook{Key: logKey, Max: max, Formatter: formatter}
}

func (h *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LogHook) Fire(entry *logrus.Entry) error {
	if RDB == nil {
		return nil
	}
	formatter := h.Formatter
	if formatter == nil {
		formatter = &logrus.TextFormatter{DisableTimestamp: true, DisableLevelTruncation: true}
	}
	cpy := *entry
	line, err := formatter.Format(&cpy)
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = RDB.LPush(ctx, h.Key, line).Err()
	_ = RDB.LTrim(ctx, h.Key, 0, int64(h.Max-1)).Err()
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
