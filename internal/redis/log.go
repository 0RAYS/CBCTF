package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const logKey = "logs"

var MaxLogScanLimit int64 = 5000

var logLevelWeight = map[string]logrus.Level{
	"trace":   logrus.TraceLevel,
	"debug":   logrus.DebugLevel,
	"info":    logrus.InfoLevel,
	"warning": logrus.WarnLevel,
	"warn":    logrus.WarnLevel,
	"error":   logrus.ErrorLevel,
	"fatal":   logrus.FatalLevel,
	"panic":   logrus.PanicLevel,
}

type LogEntry struct {
	Level string `json:"level"`
	Type  string `json:"type"`
	Time  string `json:"time"`
	Line  string `json:"line"`
}

type LogHook struct {
	Key       string
	Max       int64
	Formatter logrus.Formatter
}

func NewLogHook(max int64, formatter logrus.Formatter) *LogHook {
	if max <= 0 {
		max = 1000
	}
	MaxLogScanLimit = max
	return &LogHook{Key: logKey, Max: MaxLogScanLimit, Formatter: formatter}
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
		formatter = &logrus.TextFormatter{ForceColors: true}
	}
	line, err := formatter.Format(entry)
	if err != nil {
		return nil
	}
	entryType, _ := entry.Data["Type"].(string)
	data, err := json.Marshal(LogEntry{
		Level: entry.Level.String(),
		Type:  entryType,
		Time:  entry.Time.Format(time.RFC3339Nano),
		Line:  string(line),
	})
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = RDB.LPush(ctx, h.Key, data).Err()
	_ = RDB.LTrim(ctx, h.Key, 0, h.Max-1).Err()
	return nil
}

func GetLogs(start, end int64, minLevel string) ([]string, model.RetVal) {
	if end < start {
		return []string{}, model.SuccessRetVal()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	logs, err := RDB.LRange(ctx, logKey, 0, MaxLogScanLimit).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get logs: %s", err)
		return logs, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": logKey, "Error": err.Error()}}
	}
	minWeight, ok := logLevelWeight[strings.ToLower(minLevel)]
	if minLevel == "" || !ok {
		minWeight = logLevelWeight["trace"]
	}
	filtered := make([]string, 0, len(logs))
	for _, raw := range logs {
		entry, ok := parseLogEntry(raw)
		if !ok {
			continue
		}
		if logLevelAtLeast(entry.Level, minWeight) {
			filtered = append(filtered, entry.Line)
		}
	}
	if start >= int64(len(filtered)) {
		return []string{}, model.SuccessRetVal()
	}
	if end >= int64(len(filtered)) {
		end = int64(len(filtered) - 1)
	}
	return filtered[start : end+1], model.SuccessRetVal()
}

func parseLogEntry(raw string) (LogEntry, bool) {
	var entry LogEntry
	if err := json.Unmarshal([]byte(raw), &entry); err != nil {
		return LogEntry{}, false
	}
	if entry.Level == "" || entry.Line == "" {
		return LogEntry{}, false
	}
	return entry, true
}

func logLevelAtLeast(level string, minWeight logrus.Level) bool {
	weight, ok := logLevelWeight[strings.ToLower(level)]
	return ok && weight <= minWeight
}
