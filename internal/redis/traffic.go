package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	trafficsKey = "traffics:%d"
	trafficKey  = "traffic:%d:%d"
)

func UpdateTraffics(victim model.Victim) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	connections, err := utils.ReadPcapDir(victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read pcap: %s", err)
		return model.RetVal{Msg: i18n.Model.File.ReadPcapError, Attr: map[string]any{"Error": err.Error()}}
	}

	key := fmt.Sprintf(trafficsKey, victim.ID)
	pipe := RDB.Pipeline()
	pipe.Del(ctx, key)

	for i, conn := range connections {
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(conn.Time.UnixNano()),
			Member: fmt.Sprintf(trafficKey, victim.ID, i),
		})
		data, _ := msgpack.Marshal(&conn)
		pipe.Set(ctx, fmt.Sprintf(trafficKey, victim.ID, i), data, 30*time.Minute)
	}
	pipe.Expire(ctx, key, 30*time.Minute)
	if _, err = pipe.Exec(ctx); err != nil {
		log.Logger.Warningf("Failed to update traffics: %s", err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": trafficsKey, "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func GetTraffic(victim model.Victim) ([]utils.Connection, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	connections := make([]utils.Connection, 0)
	results, err := RDB.ZRangeWithScores(ctx, fmt.Sprintf(trafficsKey, victim.ID), 0, -1).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get traffic: %s", err)
		return nil, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": trafficsKey, "Error": err.Error()}}
	}
	pipe := RDB.Pipeline()
	for _, res := range results {
		memberKey, _ := res.Member.(string)
		pipe.Get(ctx, memberKey)
	}
	cmds, _ := pipe.Exec(ctx)

	for _, cmd := range cmds {
		str, _ := cmd.(*redis.StringCmd).Bytes()
		var conn utils.Connection
		if err = msgpack.Unmarshal(str, &conn); err != nil {
			log.Logger.Warningf("Failed to unmarshal: %s", err)
			return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
		}
		connections = append(connections, conn)
	}
	return connections, model.SuccessRetVal()
}
