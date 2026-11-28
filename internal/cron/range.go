package cron

import (
	"context"
	"fmt"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
)

// ensureRangePoolsTask 补齐靶场实例池。
// 队伍领取实例后池子变小,租约结束实例销毁,这里补回常驻数量。
func ensureRangePoolsTask() model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	root := db.CronDB.WithContext(ctx)
	challenges, _, ret := db.InitChallengeRepo(root).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.PodsChallengeType},
	})
	if !ret.OK {
		return ret
	}
	batch := model.NewBatch(len(challenges))
	for _, challenge := range challenges {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		key := fmt.Sprintf("challenge:%d", challenge.ID)
		if challenge.RangeConfig == nil {
			batch.Skip(key, "not_range")
			continue
		}
		if err := service.EnsureRangePool(ctx, challenge); err != nil {
			batch.Fail(key, "pool", model.FailRetVal(i18n.Common.UnknownError, "RangePool", key, err))
			continue
		}
		batch.Success(key, "ensured")
	}
	return batch.Result(ctx)
}
