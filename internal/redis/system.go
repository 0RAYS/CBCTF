package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

const systemMetricsKey = "system:metrics"

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
	diskPath := config.Env.Path
	if diskPath == "" {
		diskPath = "."
	}
	d, err := disk.Usage(diskPath)
	if err != nil {
		return nil, err
	}
	return &SystemMetrics{
		Timestamp: time.Now().Format(time.RFC3339Nano),
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

func GetMetrics(ctx context.Context) ([]SystemMetrics, int, model.RetVal) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	metrics := make([]SystemMetrics, 0)
	data, err := RDB.LRange(ctx, systemMetricsKey, 0, -1).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get system metrics: %s", err)
		return metrics, 0, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": systemMetricsKey, "Error": err.Error()}}
	}
	skipped := 0
	for _, d := range data {
		var m SystemMetrics
		err = json.Unmarshal([]byte(d), &m)
		if err != nil {
			skipped++
			continue
		}
		metrics = append(metrics, m)
	}
	if skipped > 0 {
		log.Logger.Warningf("Skipped invalid metric samples: count=%d", skipped)
	}
	return metrics, skipped, model.SuccessRetVal()
}
