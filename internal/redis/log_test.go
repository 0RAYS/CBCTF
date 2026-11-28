package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"CBCTF/internal/log"
)

type logWriteHook struct {
	commands []goredis.Cmder
}

func (h *logWriteHook) DialHook(next goredis.DialHook) goredis.DialHook { return next }
func (h *logWriteHook) ProcessPipelineHook(_ goredis.ProcessPipelineHook) goredis.ProcessPipelineHook {
	return func(_ context.Context, cmds []goredis.Cmder) error {
		h.commands = append(h.commands, cmds...)
		return nil
	}
}
func (h *logWriteHook) ProcessHook(_ goredis.ProcessHook) goredis.ProcessHook {
	return func(_ context.Context, _ goredis.Cmder) error { return nil }
}

func TestLogHookWritesToLevelBucket(t *testing.T) {
	oldClient, oldLogger := RDB, log.Logger
	client := goredis.NewClient(&goredis.Options{})
	writer := &logWriteHook{}
	client.AddHook(writer)
	RDB = client
	log.Init()
	t.Cleanup(func() { RDB, log.Logger = oldClient, oldLogger; _ = client.Close() })

	entry := logrus.NewEntry(log.Logger).WithFields(logrus.Fields{"Type": log.GinLogType})
	entry.Level = logrus.ErrorLevel
	entry.Message = "boom"
	entry.Time = time.Now()
	if err := NewLogHook(100, log.NewFormatter(true)).Fire(entry); err != nil {
		t.Fatal(err)
	}

	if len(writer.commands) != 2 {
		t.Fatalf("commands = %d, want 2", len(writer.commands))
	}
	add := commandArgs(writer.commands[0])
	if len(add) < 4 || add[0] != "zadd" || add[1] != "logs:error" {
		t.Fatalf("zadd args = %v", add)
	}
	trim := commandArgs(writer.commands[1])
	if len(trim) != 4 || trim[0] != "zremrangebyrank" || trim[1] != "logs:error" || trim[2] != "0" || trim[3] != "-101" {
		t.Fatalf("trim args = %v", trim)
	}
}

func commandArgs(cmd goredis.Cmder) []string {
	args := cmd.Args()
	out := make([]string, 0, len(args))
	for _, arg := range args {
		out = append(out, fmt.Sprint(arg))
	}
	return out
}

func TestLogBucketsCoverRequestedLevelAndMoreSevere(t *testing.T) {
	for _, tc := range []struct {
		level string
		want  []string
	}{
		{level: "PANIC", want: []string{"logs:panic"}},
		{level: "FATAL", want: []string{"logs:panic", "logs:fatal"}},
		{level: "ERROR", want: []string{"logs:panic", "logs:fatal", "logs:error"}},
		{level: "WARNING", want: []string{"logs:panic", "logs:fatal", "logs:error", "logs:warning"}},
		{level: "INFO", want: []string{"logs:panic", "logs:fatal", "logs:error", "logs:warning", "logs:info"}},
		{level: "DEBUG", want: []string{"logs:panic", "logs:fatal", "logs:error", "logs:warning", "logs:info", "logs:debug"}},
		{level: "", want: []string{"logs:panic", "logs:fatal", "logs:error", "logs:warning", "logs:info", "logs:debug", "logs:trace"}},
	} {
		got := logBuckets(tc.level)
		if len(got) != len(tc.want) {
			t.Fatalf("logBuckets(%q) = %v, want %v", tc.level, got, tc.want)
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Fatalf("logBuckets(%q) = %v, want %v", tc.level, got, tc.want)
			}
		}
	}
}

func TestMergeLogEntriesOrdersByTimestamp(t *testing.T) {
	newest := goredis.NewZSliceCmd(context.Background())
	newest.SetVal([]goredis.Z{{Score: 30, Member: `{"line":"newest"}`}})
	older := goredis.NewZSliceCmd(context.Background())
	older.SetVal([]goredis.Z{{Score: 10, Member: `{"line":"older"}`}, {Score: 20, Member: `{"line":"middle"}`}})

	entries := mergeLogEntries([]*goredis.ZSliceCmd{newest, older})
	if len(entries) != 3 {
		t.Fatalf("entries = %d, want 3", len(entries))
	}
	for i := range entries {
		if _, ok := parseLogEntry(entries[i].member); !ok {
			t.Fatalf("entry %d is not parsable: %s", i, entries[i].member)
		}
		if entries[i].score != float64((3-i)*10) {
			t.Fatalf("entry %d score = %v", i, entries[i].score)
		}
	}
}

func TestParseLogLevelAcceptsWarningAlias(t *testing.T) {
	level, ok := parseLogLevel("WARNING")
	if !ok || level != logrus.WarnLevel {
		t.Fatalf("parseLogLevel(WARNING) = %v, %v", level, ok)
	}
	if _, ok := parseLogLevel("nonsense"); ok {
		t.Fatal("parseLogLevel accepted an unknown level")
	}
}
