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

func UpdateTraffics(victim model.Victim) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	connections, err := utils.ReadPcapDir(victim.TrafficBasePath())
	if err != nil {
		log.Logger.Warningf("Failed to read pcap: %s", err)
		return false, i18n.ReadPcapError
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
		return false, i18n.RedisError
	}
	return true, i18n.Success
}

func GetTraffic(victim model.Victim) ([]utils.Connection, bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	connections := make([]utils.Connection, 0)
	results, err := RDB.ZRangeWithScores(ctx, fmt.Sprintf(trafficsKey, victim.ID), 0, -1).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get traffic: %s", err)
		return nil, false, i18n.RedisError
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
			return nil, false, i18n.UnknownError
		}
		connections = append(connections, conn)
	}
	return connections, true, i18n.Success
}
