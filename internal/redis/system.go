package redis

import (
	"CBCTF/internal/log"
	"context"
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

const systemMetricsKey = "system_metrics"

type SystemMetrics struct {
	Timestamp string  `json:"timestamp"`
	CPU       float64 `json:"cpu"`
	Mem       float64 `json:"mem"`
	Disk      float64 `json:"disk"`
}

func CollectMetrics() (*SystemMetrics, error) {
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

func SaveMetrics(metrics *SystemMetrics) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	err = RDB.RPush(ctx, systemMetricsKey, data).Err()
	if err != nil {
		return err
	}
	err = RDB.LTrim(ctx, systemMetricsKey, -900, -1).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetMetrics() []SystemMetrics {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	metrics := make([]SystemMetrics, 0)
	data, err := RDB.LRange(ctx, systemMetricsKey, 0, -1).Result()
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
