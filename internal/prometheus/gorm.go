package prometheus

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"

	"github.com/prometheus/client_golang/prometheus"
)

// GormCollector exports database/sql connection pool statistics for GORM.
type GormCollector struct {
	maxOpenConnectionsDesc *prometheus.Desc
	openConnectionsDesc    *prometheus.Desc
	inUseConnectionsDesc   *prometheus.Desc
	idleConnectionsDesc    *prometheus.Desc
	waitCountDesc          *prometheus.Desc
	waitDurationDesc       *prometheus.Desc
	maxIdleClosedDesc      *prometheus.Desc
	maxIdleTimeClosedDesc  *prometheus.Desc
	maxLifetimeClosedDesc  *prometheus.Desc
}

func NewGormCollector() *GormCollector {
	return &GormCollector{
		maxOpenConnectionsDesc: prometheus.NewDesc(
			"gorm_dbstats_max_open_connections",
			"Maximum number of open connections to the database.",
			nil, prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB},
		),
		openConnectionsDesc: prometheus.NewDesc(
			"gorm_dbstats_open_connections",
			"The number of established connections both in use and idle.",
			nil, prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB},
		),
		inUseConnectionsDesc: prometheus.NewDesc(
			"gorm_dbstats_in_use",
			"The number of connections currently in use.",
			nil, prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB},
		),
		idleConnectionsDesc: prometheus.NewDesc(
			"gorm_dbstats_idle",
			"The number of idle connections.",
			nil, prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB},
		),
		waitCountDesc: prometheus.NewDesc(
			"gorm_dbstats_wait_count",
			"The total number of connections waited for.",
			nil, prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB},
		),
		waitDurationDesc: prometheus.NewDesc(
			"gorm_dbstats_wait_duration",
			"The total time blocked waiting for a new connection.",
			nil, prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB},
		),
		maxIdleClosedDesc: prometheus.NewDesc(
			"gorm_dbstats_max_idle_closed",
			"The total number of connections closed due to SetMaxIdleConns.",
			nil, prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB},
		),
		maxIdleTimeClosedDesc: prometheus.NewDesc(
			"gorm_dbstats_max_idletime_closed",
			"The total number of connections closed due to SetConnMaxIdleTime.",
			nil, prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB},
		),
		maxLifetimeClosedDesc: prometheus.NewDesc(
			"gorm_dbstats_max_lifetime_closed",
			"The total number of connections closed due to SetConnMaxLifetime.",
			nil, prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB},
		),
	}
}

func (c *GormCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenConnectionsDesc
	ch <- c.openConnectionsDesc
	ch <- c.inUseConnectionsDesc
	ch <- c.idleConnectionsDesc
	ch <- c.waitCountDesc
	ch <- c.waitDurationDesc
	ch <- c.maxIdleClosedDesc
	ch <- c.maxIdleTimeClosedDesc
	ch <- c.maxLifetimeClosedDesc
}

func (c *GormCollector) Collect(ch chan<- prometheus.Metric) {
	if db.DB == nil {
		return
	}
	sqlDB, err := db.DB.DB()
	if err != nil {
		return
	}
	stats := sqlDB.Stats()

	ch <- prometheus.MustNewConstMetric(c.maxOpenConnectionsDesc, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(c.openConnectionsDesc, prometheus.GaugeValue, float64(stats.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.inUseConnectionsDesc, prometheus.GaugeValue, float64(stats.InUse))
	ch <- prometheus.MustNewConstMetric(c.idleConnectionsDesc, prometheus.GaugeValue, float64(stats.Idle))
	ch <- prometheus.MustNewConstMetric(c.waitCountDesc, prometheus.GaugeValue, float64(stats.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.waitDurationDesc, prometheus.GaugeValue, float64(stats.WaitDuration))
	ch <- prometheus.MustNewConstMetric(c.maxIdleClosedDesc, prometheus.GaugeValue, float64(stats.MaxIdleClosed))
	ch <- prometheus.MustNewConstMetric(c.maxIdleTimeClosedDesc, prometheus.GaugeValue, float64(stats.MaxIdleTimeClosed))
	ch <- prometheus.MustNewConstMetric(c.maxLifetimeClosedDesc, prometheus.GaugeValue, float64(stats.MaxLifetimeClosed))
}
