package redis

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"context"
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type SystemMetrics struct {
	Timestamp string  `json:"timestamp"`
	CPU       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	Disk      float64 `json:"disk"`
}

func StartCollect(ctx context.Context) {
	log.Logger.Info("Start collecting Redis metrics")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			metrics, err := collectMetrics()
			if err != nil {
				log.Logger.Error("Failed to collect Redis metrics: ", err)
				continue
			}
			if err := saveMetrics(metrics); err != nil {
				log.Logger.Error("Failed to save Redis metrics: ", err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func collectMetrics() (*SystemMetrics, error) {
	c, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	m, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	d, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	return &SystemMetrics{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		CPU:       c[0],
		Mem:       m.UsedPercent,
		Disk:      d.UsedPercent,
	}, nil
}

func saveMetrics(metrics *SystemMetrics) error {
	ctx := context.Background()
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	err = RDB.RPush(ctx, "system_metrics", data).Err()
	if err != nil {
		return err
	}
	err = RDB.LTrim(ctx, "system_metrics", -900, -1).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetMetrics() []SystemMetrics {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	metrics := make([]SystemMetrics, 0)
	data, err := RDB.LRange(ctx, "system_metrics", 0, -1).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get system metrics: %s", err)
		return metrics
	}
	for _, d := range data {
		var m SystemMetrics
		err = json.Unmarshal([]byte(d), &m)
		if err != nil {
			log.Logger.Warningf("Failed to parse system metrics: %s", err)
			return metrics
		}
		metrics = append(metrics, m)
	}
	return metrics
}
