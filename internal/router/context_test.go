package router

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/redis"
)

type cancelledRequestHook struct {
	systemLogsRedisHook
	observed bool
}

func (h *cancelledRequestHook) ProcessHook(_ goredis.ProcessHook) goredis.ProcessHook {
	return func(ctx context.Context, _ goredis.Cmder) error {
		h.observed = errors.Is(ctx.Err(), context.Canceled)
		return ctx.Err()
	}
}

func TestRequestCancellationReachesRedisAndDatabaseWithoutGinFallback(t *testing.T) {
	oldDB, oldRedis, oldLogger, oldBundle := db.DB, redis.RDB, log.Logger, i18n.Bundle
	t.Cleanup(func() { db.DB, redis.RDB, log.Logger, i18n.Bundle = oldDB, oldRedis, oldLogger, oldBundle })
	log.Logger = logrus.New()
	log.Logger.SetOutput(io.Discard)
	i18n.Init()
	requestCtx, cancel := context.WithCancel(context.Background())
	cancel()
	recorder := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(recorder)
	engine.ContextWithFallback = false
	ctx.Request = httptest.NewRequest("GET", "/admin/logs?limit=100&offset=0&level=INFO", nil).WithContext(requestCtx)
	hook := &cancelledRequestHook{}
	redis.RDB = goredis.NewClient(&goredis.Options{})
	defer func(RDB *goredis.Client) {
		_ = RDB.Close()
	}(redis.RDB)
	redis.RDB.AddHook(hook)
	GetLogs(ctx)
	if !hook.observed {
		t.Fatal("Redis did not receive canceled request context")
	}
	var err error
	db.DB, err = gorm.Open(postgres.Open("host=127.0.0.1 dbname=test"), &gorm.Config{DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	pool, _ := db.DB.DB()
	defer func(pool *sql.DB) {
		_ = pool.Close()
	}(pool)
	observed := false
	err = db.DB.Callback().Query().Before("gorm:query").Register("test:context", func(query *gorm.DB) {
		observed = errors.Is(query.Statement.Context.Err(), context.Canceled)
		_ = query.AddError(context.Canceled)
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/admin/victims/1", nil).WithContext(requestCtx)
	ctx.Params = gin.Params{{Key: "victimID", Value: "1"}}
	middleware.SetVictim(ctx)
	if !observed {
		t.Fatal("authorization query lost request cancellation")
	}
}
