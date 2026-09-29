package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/traffic"
)

const trafficsKeyTmpl = "traffic:snapshot:%d"

// One atomic value preserves nanosecond ordering and represents an empty
// capture distinctly from a cache miss. No packet keys can expire independently.
func StoreTraffic(ctx context.Context, victim model.Victim, result *traffic.PcapDirResult) model.RetVal {
	data, err := msgpack.Marshal(result)
	if err == nil {
		ttl := 30 * time.Minute
		if victim.Status != model.StoppedVictimStatus || len(result.SourceIssues) > 0 {
			ttl = 15 * time.Second
		}
		err = RDB.Set(ctx, fmt.Sprintf(trafficsKeyTmpl, victim.ID), data, ttl).Err()
	}
	if err != nil {
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": trafficsKeyTmpl, "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func GetTraffic(ctx context.Context, victim model.Victim) (*traffic.PcapDirResult, model.RetVal) {
	data, err := RDB.Get(ctx, fmt.Sprintf(trafficsKeyTmpl, victim.ID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, model.SuccessRetVal()
	}
	if err != nil {
		return nil, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": trafficsKeyTmpl, "Error": err.Error()}}
	}
	var snapshot traffic.PcapDirResult
	if err = msgpack.Unmarshal(data, &snapshot); err != nil {
		return nil, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": trafficsKeyTmpl, "Error": err.Error()}}
	}
	return &snapshot, model.SuccessRetVal()
}
