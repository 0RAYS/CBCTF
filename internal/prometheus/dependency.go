package prometheus

import (
	"context"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/version"
)

var processStart = time.Now()

// BuildCollector exposes the identity and uptime of the running binary.
type BuildCollector struct {
	buildInfoDesc *prometheus.Desc
	uptimeDesc    *prometheus.Desc
}

func NewBuildCollector() *BuildCollector {
	return &BuildCollector{
		buildInfoDesc: prometheus.NewDesc(
			"cbctf_build_info",
			"Build information of the running CBCTF binary",
			[]string{"version", "commit", "go_version"}, nil,
		),
		uptimeDesc: prometheus.NewDesc(
			"cbctf_uptime_seconds",
			"Seconds since the CBCTF process started",
			nil, nil,
		),
	}
}

func (c *BuildCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.buildInfoDesc
	ch <- c.uptimeDesc
}

func (c *BuildCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.buildInfoDesc, prometheus.GaugeValue, 1, version.Version, version.Commit, runtime.Version())
	ch <- prometheus.MustNewConstMetric(c.uptimeDesc, prometheus.GaugeValue, time.Since(processStart).Seconds())
}

// DependencyCollector probes the dependencies the platform cannot work without.
type DependencyCollector struct {
	readyDesc *prometheus.Desc
}

func NewDependencyCollector() *DependencyCollector {
	return &DependencyCollector{
		readyDesc: prometheus.NewDesc(
			"cbctf_dependency_ready",
			"Whether a dependency the platform needs is reachable.",
			[]string{"dependency"}, nil,
		),
	}
}

func (c *DependencyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.readyDesc
}

func (c *DependencyCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch <- prometheus.MustNewConstMetric(c.readyDesc, prometheus.GaugeValue, boolValue(c.postgresReady(ctx)), "postgres")
	ch <- prometheus.MustNewConstMetric(c.readyDesc, prometheus.GaugeValue, boolValue(c.valkeyReady(ctx)), "valkey")
	ch <- prometheus.MustNewConstMetric(c.readyDesc, prometheus.GaugeValue, boolValue(k8s.Ping(ctx) == nil), "kubernetes")
}

func (c *DependencyCollector) postgresReady(ctx context.Context) bool {
	if db.DB == nil {
		return false
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		return false
	}
	return sqlDB.PingContext(ctx) == nil
}

func (c *DependencyCollector) valkeyReady(ctx context.Context) bool {
	client := redisClientRef.Load()
	if client == nil {
		return false
	}
	return client.Ping(ctx).Err() == nil
}

func boolValue(ready bool) float64 {
	if ready {
		return 1
	}
	return 0
}
