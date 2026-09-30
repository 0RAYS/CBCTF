package redis

import (
	"context"
	"testing"

	"CBCTF/internal/log"

	goredis "github.com/redis/go-redis/v9"
)

type metricsReadHook struct{}

func (metricsReadHook) DialHook(next goredis.DialHook) goredis.DialHook { return next }
func (metricsReadHook) ProcessPipelineHook(next goredis.ProcessPipelineHook) goredis.ProcessPipelineHook {
	return next
}
func (metricsReadHook) ProcessHook(_ goredis.ProcessHook) goredis.ProcessHook {
	return func(_ context.Context, cmd goredis.Cmder) error {
		cmd.(*goredis.StringSliceCmd).SetVal([]string{`{"timestamp":"first","cpu":1}`, `broken`, `{"timestamp":"last","cpu":2}`})
		return nil
	}
}
func TestMetricDecodeFailureDoesNotTruncateLaterSamples(t *testing.T) {
	old, oldLogger := RDB, log.Logger
	client := goredis.NewClient(&goredis.Options{})
	client.AddHook(metricsReadHook{})
	RDB = client
	log.Init()
	t.Cleanup(func() { RDB, log.Logger = old, oldLogger; _ = client.Close() })
	metrics, skipped, ret := GetMetrics(context.Background())
	if !ret.OK || skipped != 1 || len(metrics) != 2 || metrics[1].Timestamp != "last" {
		t.Fatalf("lost later samples: %+v skipped=%d %+v", metrics, skipped, ret)
	}
}
