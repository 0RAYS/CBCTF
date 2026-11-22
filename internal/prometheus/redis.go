package prometheus

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

// redisClientRef holds the shared client so the collector can read pool
// statistics without importing the redis package.
var redisClientRef atomic.Pointer[redis.Client]

// RedisCommandDuration 记录单条 Valkey/Redis 命令的耗时，命令名是低基数标签。
var RedisCommandDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "cbctf_redis_command_duration_seconds",
		Help:    "Histogram of Redis command execution time.",
		Buckets: prometheus.ExponentialBuckets(0.0001, 4, 8),
	},
	[]string{"command"},
)

// RedisCommandErrorsTotal 记录命令返回错误的次数。
var RedisCommandErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "cbctf_redis_command_errors_total",
		Help: "Total number of Redis commands that returned an error",
	},
	[]string{"command"},
)

type RedisCollector struct {
	poolTotalConnsDesc    *prometheus.Desc
	poolIdleConnsDesc     *prometheus.Desc
	poolStaleConnsDesc    *prometheus.Desc
	poolPendingDesc       *prometheus.Desc
	poolHitsDesc          *prometheus.Desc
	poolMissesDesc        *prometheus.Desc
	poolTimeoutsDesc      *prometheus.Desc
	poolWaitCountDesc     *prometheus.Desc
	poolWaitDurationDesc  *prometheus.Desc
	dependencyReadyDesc   *prometheus.Desc
	collectionSuccessDesc *prometheus.Desc
}

func NewRedisCollector() *RedisCollector {
	return &RedisCollector{
		poolTotalConnsDesc: prometheus.NewDesc(
			"cbctf_redis_pool_total_connections",
			"Number of established connections to Redis.",
			nil, nil,
		),
		poolIdleConnsDesc: prometheus.NewDesc(
			"cbctf_redis_pool_idle_connections",
			"Number of idle connections in the Redis pool.",
			nil, nil,
		),
		poolStaleConnsDesc: prometheus.NewDesc(
			"cbctf_redis_pool_stale_connections",
			"Number of stale connections removed from the Redis pool.",
			nil, nil,
		),
		poolPendingDesc: prometheus.NewDesc(
			"cbctf_redis_pool_pending_requests",
			"Number of requests waiting for a Redis connection.",
			nil, nil,
		),
		poolHitsDesc: prometheus.NewDesc(
			"cbctf_redis_pool_hits_total",
			"Total number of times a free connection was found in the pool.",
			nil, nil,
		),
		poolMissesDesc: prometheus.NewDesc(
			"cbctf_redis_pool_misses_total",
			"Total number of times a free connection was not found in the pool.",
			nil, nil,
		),
		poolTimeoutsDesc: prometheus.NewDesc(
			"cbctf_redis_pool_timeouts_total",
			"Total number of times a wait timeout occurred in the pool.",
			nil, nil,
		),
		poolWaitCountDesc: prometheus.NewDesc(
			"cbctf_redis_pool_wait_count_total",
			"Total number of times a connection was waited for.",
			nil, nil,
		),
		poolWaitDurationDesc: prometheus.NewDesc(
			"cbctf_redis_pool_wait_duration_seconds_total",
			"Total time in seconds spent waiting for a Redis connection.",
			nil, nil,
		),
		dependencyReadyDesc: prometheus.NewDesc(
			"cbctf_dependency_ready",
			"Whether a dependency the platform needs is reachable.",
			[]string{"dependency"}, nil,
		),
		collectionSuccessDesc: prometheus.NewDesc(
			"cbctf_redis_collection_success",
			"Whether the Redis collector reached the server.",
			nil, nil,
		),
	}
}

func (c *RedisCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.poolTotalConnsDesc
	ch <- c.poolIdleConnsDesc
	ch <- c.poolStaleConnsDesc
	ch <- c.poolPendingDesc
	ch <- c.poolHitsDesc
	ch <- c.poolMissesDesc
	ch <- c.poolTimeoutsDesc
	ch <- c.poolWaitCountDesc
	ch <- c.poolWaitDurationDesc
	ch <- c.dependencyReadyDesc
	ch <- c.collectionSuccessDesc
}

func (c *RedisCollector) Collect(ch chan<- prometheus.Metric) {
	client := redisClientRef.Load()
	if client == nil {
		ch <- prometheus.MustNewConstMetric(c.dependencyReadyDesc, prometheus.GaugeValue, 0, "valkey")
		ch <- prometheus.MustNewConstMetric(c.collectionSuccessDesc, prometheus.GaugeValue, 0)
		return
	}

	stats := client.PoolStats()
	ch <- prometheus.MustNewConstMetric(c.poolTotalConnsDesc, prometheus.GaugeValue, float64(stats.TotalConns))
	ch <- prometheus.MustNewConstMetric(c.poolIdleConnsDesc, prometheus.GaugeValue, float64(stats.IdleConns))
	ch <- prometheus.MustNewConstMetric(c.poolStaleConnsDesc, prometheus.GaugeValue, float64(stats.StaleConns))
	ch <- prometheus.MustNewConstMetric(c.poolPendingDesc, prometheus.GaugeValue, float64(stats.PendingRequests))
	ch <- prometheus.MustNewConstMetric(c.poolHitsDesc, prometheus.CounterValue, float64(stats.Hits))
	ch <- prometheus.MustNewConstMetric(c.poolMissesDesc, prometheus.CounterValue, float64(stats.Misses))
	ch <- prometheus.MustNewConstMetric(c.poolTimeoutsDesc, prometheus.CounterValue, float64(stats.Timeouts))
	ch <- prometheus.MustNewConstMetric(c.poolWaitCountDesc, prometheus.CounterValue, float64(stats.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.poolWaitDurationDesc, prometheus.CounterValue, float64(stats.WaitDurationNs)/float64(time.Second))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		ch <- prometheus.MustNewConstMetric(c.dependencyReadyDesc, prometheus.GaugeValue, 0, "valkey")
		ch <- prometheus.MustNewConstMetric(c.collectionSuccessDesc, prometheus.GaugeValue, 0)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.dependencyReadyDesc, prometheus.GaugeValue, 1, "valkey")
	ch <- prometheus.MustNewConstMetric(c.collectionSuccessDesc, prometheus.GaugeValue, 1)
}

// AttachRedisMetrics registers the command-level hook on the shared client.
// It must run again whenever the client is recreated.
func AttachRedisMetrics(client *redis.Client) {
	if client == nil {
		return
	}
	redisClientRef.Store(client)
	client.AddHook(redisCommandHook{})
}

type redisCommandHook struct{}

func (redisCommandHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (redisCommandHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		name := cmd.Name()
		RedisCommandDuration.WithLabelValues(name).Observe(time.Since(start).Seconds())
		if err != nil {
			RedisCommandErrorsTotal.WithLabelValues(name).Inc()
		}
		return err
	}
}

func (redisCommandHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmds)
		RedisCommandDuration.WithLabelValues("pipeline").Observe(time.Since(start).Seconds())
		if err != nil {
			RedisCommandErrorsTotal.WithLabelValues("pipeline").Inc()
		}
		return err
	}
}
