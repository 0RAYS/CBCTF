package router

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/redis"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type systemLogsRedisHook struct {
	lines []string
	err   error
}

func (h systemLogsRedisHook) DialHook(next goredis.DialHook) goredis.DialHook {
	return next
}

func (h systemLogsRedisHook) ProcessPipelineHook(next goredis.ProcessPipelineHook) goredis.ProcessPipelineHook {
	return next
}

func (h systemLogsRedisHook) ProcessHook(_ goredis.ProcessHook) goredis.ProcessHook {
	return func(_ context.Context, cmd goredis.Cmder) error {
		cmd.(*goredis.StringSliceCmd).SetVal(h.lines)
		return h.err
	}
}

func TestGetLogsResponse(t *testing.T) {
	oldClient, oldLogger, oldBundle := redis.RDB, log.Logger, i18n.Bundle
	t.Cleanup(func() {
		redis.RDB, log.Logger, i18n.Bundle = oldClient, oldLogger, oldBundle
	})
	log.Logger = logrus.New()
	log.Logger.SetOutput(io.Discard)
	i18n.Init()

	for _, tc := range []struct {
		name   string
		lines  []string
		err    error
		offset string
		want   string
	}{
		{name: "Redis failure", err: errors.New("logs-test-unavailable"), offset: "0", want: "null"},
		{name: "empty Redis list", offset: "0", want: "[]"},
		{name: "empty page", lines: []string{`{"level":"info","line":"ready"}`}, offset: "100", want: "[]"},
		{name: "populated page", lines: []string{`{"level":"info","line":"ready"}`}, offset: "0", want: `["ready"]`},
	} {
		for _, lang := range []string{"en", "zh-CN"} {
			t.Run(tc.name+"/"+lang, func(t *testing.T) {
				client := goredis.NewClient(&goredis.Options{})
				client.AddHook(systemLogsRedisHook{lines: tc.lines, err: tc.err})
				t.Cleanup(func() { _ = client.Close() })
				redis.RDB = client

				recorder := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(recorder)
				ctx.Request = httptest.NewRequest("GET", "/admin/logs?limit=100&offset="+tc.offset+"&level=INFO", nil)
				ctx.Request.Header.Set("Accept-Language", lang)
				GetLogs(ctx)

				var body struct {
					Code int             `json:"code"`
					Msg  string          `json:"msg"`
					Data json.RawMessage `json:"data"`
				}
				if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if tc.err != nil {
					if body.Code == 200 || !strings.Contains(body.Msg, tc.err.Error()) {
						t.Fatalf("Redis failure was not preserved: %s", recorder.Body.String())
					}
				} else if body.Code != 200 {
					t.Fatalf("successful read failed: %s", recorder.Body.String())
				}
				if string(body.Data) != tc.want {
					t.Fatalf("data = %s, want %s", body.Data, tc.want)
				}
			})
		}
	}
}
