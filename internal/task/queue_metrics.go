package task

import (
	prom "github.com/prometheus/client_golang/prometheus"

	"CBCTF/internal/prometheus"
)

// QueueCollector 暴露 Asynq 队列的实时深度、等待时长与内存占用。
// 采集器放在本包是因为 Inspector 与任务队列同源，prometheus 包不能反向依赖这里。
type QueueCollector struct {
	collectionSuccessDesc *prom.Desc
}

func NewQueueCollector() *QueueCollector {
	return &QueueCollector{
		collectionSuccessDesc: prom.NewDesc(
			"cbctf_task_collection_success",
			"Whether the async task queue metrics were collected",
			nil, nil,
		),
	}
}

func (c *QueueCollector) Describe(ch chan<- *prom.Desc) {
	ch <- c.collectionSuccessDesc
}

func (c *QueueCollector) Collect(ch chan<- prom.Metric) {
	if inspector == nil {
		ch <- prom.MustNewConstMetric(c.collectionSuccessDesc, prom.GaugeValue, 0)
		return
	}
	queues, err := inspector.Queues()
	if err != nil {
		ch <- prom.MustNewConstMetric(c.collectionSuccessDesc, prom.GaugeValue, 0)
		return
	}
	success := float64(1)
	for _, queue := range queues {
		info, err := inspector.GetQueueInfo(queue)
		if err != nil {
			success = 0
			continue
		}
		for state, size := range map[string]int{
			"pending":     info.Pending,
			"active":      info.Active,
			"scheduled":   info.Scheduled,
			"retry":       info.Retry,
			"archived":    info.Archived,
			"completed":   info.Completed,
			"aggregating": info.Aggregating,
		} {
			prometheus.TaskQueueDepth.WithLabelValues(queue, state).Set(float64(size))
		}
		prometheus.TaskQueueLatency.WithLabelValues(queue).Set(info.Latency.Seconds())
		prometheus.TaskQueueMemory.WithLabelValues(queue).Set(float64(info.MemoryUsage))
	}
	ch <- prom.MustNewConstMetric(c.collectionSuccessDesc, prom.GaugeValue, success)
}
