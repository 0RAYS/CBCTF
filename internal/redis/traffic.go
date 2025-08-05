package redis

import (
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"os"
	"strings"
	"time"
)

func GetTraffics(victim model.Victim) ([]utils.Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	key := fmt.Sprintf("traffic:%d", victim.ID)
	traffics := make([]utils.Connection, 0)
	results, err := RDB.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return traffics, err
	}

	pipe := RDB.Pipeline()
	for _, res := range results {
		pipe.Get(ctx, res.Member.(string))
	}
	cmds, _ := pipe.Exec(ctx)

	for _, cmd := range cmds {
		str, _ := cmd.(*redis.StringCmd).Bytes()
		var conn utils.Connection
		err = msgpack.Unmarshal(str, &conn)
		if err != nil {
			return traffics, err
		}
		traffics = append(traffics, conn)
	}
	return traffics, nil
}

func LoadTraffics(victim model.Victim) ([]utils.Connection, error) {
	dir, err := os.ReadDir(victim.TrafficBasePath())
	if err != nil {
		return nil, err
	}
	connections := make([]utils.Connection, 0)
	for _, file := range dir {
		if file.IsDir() || (!strings.HasSuffix(file.Name(), ".pcap") && !strings.HasSuffix(file.Name(), ".pcapng")) {
			continue
		}
		packet, err := utils.ReadPcap(fmt.Sprintf("%s/%s", victim.TrafficBasePath(), file.Name()))
		if err != nil {
			return nil, err
		}
		connections = append(connections, packet...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	key := fmt.Sprintf("traffic:%d", victim.ID)
	pipe := RDB.Pipeline()
	pipe.Del(ctx, key)
	for i, conn := range connections {
		timestamp := conn.Time.UnixMilli()
		score := float64(1e13 - timestamp)
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  score,
			Member: i,
		})
		data, _ := msgpack.Marshal(&conn)
		pipe.Set(ctx, fmt.Sprintf("%d", i), data, 10*time.Minute)
	}
	pipe.Expire(ctx, key, 10*time.Minute)
	_, err = pipe.Exec(ctx)
	return connections, err
}
