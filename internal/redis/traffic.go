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
func UpdateTraffics(ctx context.Context, victim model.Victim) model.RetVal {
	result, err := traffic.ReadPcapDirWithContext(ctx, victim.TrafficBasePath())
	if err != nil {
		return model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}
	data, err := msgpack.Marshal(result.Connections)
	if err == nil {
		ttl := 30 * time.Minute
		if victim.Status != model.StoppedVictimStatus {
			ttl = 15 * time.Second
		}
		err = RDB.Set(ctx, fmt.Sprintf(trafficsKeyTmpl, victim.ID), data, ttl).Err()
	}
	if err != nil {
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": trafficsKeyTmpl, "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func GetTraffic(ctx context.Context, victim model.Victim) ([]traffic.Connection, model.RetVal) {
	data, err := RDB.Get(ctx, fmt.Sprintf(trafficsKeyTmpl, victim.ID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, model.SuccessRetVal()
	}
	if err != nil {
		return nil, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": trafficsKeyTmpl, "Error": err.Error()}}
	}
	connections := make([]traffic.Connection, 0)
	if err = msgpack.Unmarshal(data, &connections); err != nil {
		return nil, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": trafficsKeyTmpl, "Error": err.Error()}}
	}
	return connections, model.SuccessRetVal()
}
