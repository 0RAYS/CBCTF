package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

// Logs live in one sorted set per level, so a level filter reads only the buckets it
// needs instead of scanning the whole history.
const (
	logKeyPrefix = "logs:"
	legacyLogKey = "logs"
	// logBucketLimit caps the entries kept per level.
	logBucketLimit = 5000
)

// logLevels orders the buckets from the most severe to the least severe.
var logLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
	logrus.TraceLevel,
}

var logSequence atomic.Uint64
var logWriteFailures atomic.Int64

// LogEntry carries the rendered line; the sequence and the timestamp keep the sorted
// set members unique, a repeated member would overwrite an earlier log line.
type LogEntry struct {
	Seq  uint64 `json:"seq"`
	Time string `json:"time"`
	Line string `json:"line"`
}

type LogHook struct {
	Max       int64
	Formatter logrus.Formatter
}

func NewLogHook(max int64, formatter logrus.Formatter) *LogHook {
	if max <= 0 {
		max = 1000
	}
	return &LogHook{Max: max, Formatter: formatter}
}

func (*LogHook) Levels() []logrus.Level {
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
	payload, err := json.Marshal(LogEntry{
		Seq:  logSequence.Add(1),
		Time: entry.Time.Format(time.RFC3339Nano),
		Line: string(line),
	})
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	key := logKeyPrefix + entry.Level.String()
	pipe := RDB.Pipeline()
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(entry.Time.UnixMicro()), Member: payload})
	pipe.ZRemRangeByRank(ctx, key, 0, -(h.Max + 1))
	if _, err := pipe.Exec(ctx); err != nil {
		count := logWriteFailures.Add(1)
		if count == 1 || count%100 == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "failed to mirror log to redis (%d failures): %v\n", count, err)
		}
	}
	return nil
}

func GetLogs(ctx context.Context, start, end int64, minLevel string) ([]string, model.RetVal) {
	if end < start {
		return []string{}, model.SuccessRetVal()
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	buckets := logBuckets(minLevel)
	pipe := RDB.Pipeline()
	cmds := make([]*redis.ZSliceCmd, 0, len(buckets))
	for _, bucket := range buckets {
		cmds = append(cmds, pipe.ZRevRangeWithScores(ctx, bucket, 0, end))
	}
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		log.Logger.Warningf("Failed to get logs: %s", err)
		return []string{}, model.FailRetVal(i18n.Redis.GetError, "Cache", logKeyPrefix, err)
	}
	entries := mergeLogEntries(cmds)
	if start >= int64(len(entries)) {
		return []string{}, model.SuccessRetVal()
	}
	if end >= int64(len(entries)) {
		end = int64(len(entries) - 1)
	}
	lines := make([]string, 0, end-start+1)
	for _, entry := range entries[start : end+1] {
		parsed, ok := parseLogEntry(entry.member)
		if !ok {
			continue
		}
		lines = append(lines, parsed.Line)
	}
	return lines, model.SuccessRetVal()
}

type logEntryRef struct {
	score  float64
	member string
}

// mergeLogEntries orders the buckets by timestamp, the newest entry first.
func mergeLogEntries(cmds []*redis.ZSliceCmd) []logEntryRef {
	entries := make([]logEntryRef, 0, 64)
	for _, cmd := range cmds {
		values, err := cmd.Result()
		if err != nil {
			continue
		}
		for _, value := range values {
			member, ok := value.Member.(string)
			if !ok {
				continue
			}
			entries = append(entries, logEntryRef{score: value.Score, member: member})
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].score != entries[j].score {
			return entries[i].score > entries[j].score
		}
		return entries[i].member > entries[j].member
	})
	return entries
}

// logBuckets returns the buckets holding the requested level and the more severe ones.
func logBuckets(level string) []string {
	weight, ok := parseLogLevel(level)
	if !ok {
		weight = logrus.TraceLevel
	}
	buckets := make([]string, 0, len(logLevels))
	for _, l := range logLevels {
		if l <= weight {
			buckets = append(buckets, logKeyPrefix+l.String())
		}
	}
	return buckets
}

func parseLogLevel(level string) (logrus.Level, bool) {
	normalized := strings.ToLower(strings.TrimSpace(level))
	if normalized == "warning" {
		return logrus.WarnLevel, true
	}
	parsed, err := logrus.ParseLevel(normalized)
	if err != nil {
		return logrus.TraceLevel, false
	}
	return parsed, true
}

func parseLogEntry(raw string) (LogEntry, bool) {
	var entry LogEntry
	if err := json.Unmarshal([]byte(raw), &entry); err != nil {
		return LogEntry{}, false
	}
	if entry.Line == "" {
		return LogEntry{}, false
	}
	return entry, true
}
