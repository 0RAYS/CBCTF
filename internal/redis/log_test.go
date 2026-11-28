package redis

import (
	"context"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

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
