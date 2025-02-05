package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"encoding/json"
	"errors"
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
	if !config.Env.Redis.On {
		return
	}
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
	if !config.Env.Redis.On {
		return nil, errors.New("redis off")
	}
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
	if !config.Env.Redis.On {
		return errors.New("redis off")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
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

func GetMetrics() ([]SystemMetrics, error) {
	if !config.Env.Redis.On {
		return nil, errors.New("redis off")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	data, err := RDB.LRange(ctx, "system_metrics", 0, -1).Result()
	if err != nil {
		cancel()
		return nil, err
	}
	cancel()
	var metrics []SystemMetrics
	for _, d := range data {
		var m SystemMetrics
		err = json.Unmarshal([]byte(d), &m)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}
