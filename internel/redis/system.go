package redis

import (
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

func StartCollect() {
	log.Logger.Info("Start collecting Redis metrics")
	for {
		metrics, err := collectMetrics()
		if err != nil {
			log.Logger.Warningf("Failed to collect Redis metrics: %s", err)
			continue
		}
		if err = saveMetrics(metrics); err != nil {
			log.Logger.Warningf("Failed to save Redis metrics: %s", err)
		}
		time.Sleep(1 * time.Second)
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
	d, err := disk.Usage(".")
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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
