package prometheus

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const postgresMetricPrefix = "gorm_status_"

type PostgresCollector struct {
	labels prometheus.Labels

	replicationLagDesc    *prometheus.Desc
	postmasterStartedDesc *prometheus.Desc
	databaseSizeDesc      *prometheus.Desc
	recordCountDesc       *prometheus.Desc

	seqScanDesc              *prometheus.Desc
	seqTupReadDesc           *prometheus.Desc
	idxScanDesc              *prometheus.Desc
	idxTupFetchDesc          *prometheus.Desc
	nTupInsDesc              *prometheus.Desc
	nTupUpdDesc              *prometheus.Desc
	nTupDelDesc              *prometheus.Desc
	nTupHotUpdDesc           *prometheus.Desc
	nLiveTupDesc             *prometheus.Desc
	nDeadTupDesc             *prometheus.Desc
	nModSinceLastAnalyzeDesc *prometheus.Desc
	lastVacuumDesc           *prometheus.Desc
	lastAutovacuumDesc       *prometheus.Desc
	lastAnalyzeDesc          *prometheus.Desc
	lastAutoanalyzeDesc      *prometheus.Desc
	vacuumCountDesc          *prometheus.Desc
	autovacuumCountDesc      *prometheus.Desc
	analyzeCountDesc         *prometheus.Desc
	autoanalyzeCountDesc     *prometheus.Desc

	heapBlksReadDesc     *prometheus.Desc
	heapBlksHitDesc      *prometheus.Desc
	idxBlksReadDesc      *prometheus.Desc
	idxBlksHitDesc       *prometheus.Desc
	toastBlksReadDesc    *prometheus.Desc
	toastBlksHitDesc     *prometheus.Desc
	toastIdxBlksReadDesc *prometheus.Desc
	toastIdxBlksHitDesc  *prometheus.Desc
}

func NewPostgresCollector() *PostgresCollector {
	labels := prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB}
	tableLabels := []string{"datname", "schemaname", "relname"}

	return &PostgresCollector{
		labels: labels,

		replicationLagDesc: prometheus.NewDesc(
			postgresMetricPrefix+"lag",
			"Replication lag behind master in seconds",
			nil, labels,
		),
		postmasterStartedDesc: prometheus.NewDesc(
			postgresMetricPrefix+"start_time_seconds",
			"Time unix timestamp at which postmaster started",
			nil, labels,
		),
		databaseSizeDesc: prometheus.NewDesc(
			postgresMetricPrefix+"size_bytes",
			"Size of database in bytes",
			[]string{"datname"}, labels,
		),
		recordCountDesc: prometheus.NewDesc(
			postgresMetricPrefix+"rows_count",
			"Name of this table",
			[]string{"table_schema", "table_name"}, labels,
		),

		seqScanDesc:              newPostgresTableDesc("seq_scan", "Number of sequential scans initiated on this table", tableLabels, labels),
		seqTupReadDesc:           newPostgresTableDesc("seq_tup_read", "Number of live rows fetched by sequential scans", tableLabels, labels),
		idxScanDesc:              newPostgresTableDesc("idx_scan", "Number of index scans initiated on this table", tableLabels, labels),
		idxTupFetchDesc:          newPostgresTableDesc("idx_tup_fetch", "Number of live rows fetched by index scans", tableLabels, labels),
		nTupInsDesc:              newPostgresTableDesc("n_tup_ins", "Number of rows inserted", tableLabels, labels),
		nTupUpdDesc:              newPostgresTableDesc("n_tup_upd", "Number of rows updated", tableLabels, labels),
		nTupDelDesc:              newPostgresTableDesc("n_tup_del", "Number of rows deleted", tableLabels, labels),
		nTupHotUpdDesc:           newPostgresTableDesc("n_tup_hot_upd", "Number of rows HOT updated (i.e., with no separate index update required)", tableLabels, labels),
		nLiveTupDesc:             newPostgresTableDesc("n_live_tup", "Estimated number of live rows", tableLabels, labels),
		nDeadTupDesc:             newPostgresTableDesc("n_dead_tup", "Estimated number of dead rows", tableLabels, labels),
		nModSinceLastAnalyzeDesc: newPostgresTableDesc("n_mod_since_last_analyze", "Estimated number of rows changed since last analyze", tableLabels, labels),
		lastVacuumDesc:           newPostgresTableDesc("last_vacuum", "Last time at which this table was manually vacuumed (not counting VACUUM FULL)", tableLabels, labels),
		lastAutovacuumDesc:       newPostgresTableDesc("last_autovacuum", "Last time at which this table was vacuumed by the autovacuum daemon", tableLabels, labels),
		lastAnalyzeDesc:          newPostgresTableDesc("last_analyze", "Last time at which this table was manually analyzed", tableLabels, labels),
		lastAutoanalyzeDesc:      newPostgresTableDesc("last_autoanalyze", "Last time at which this table was analyzed by the autovacuum daemon", tableLabels, labels),
		vacuumCountDesc:          newPostgresTableDesc("vacuum_count", "Number of times this table has been manually vacuumed (not counting VACUUM FULL)", tableLabels, labels),
		autovacuumCountDesc:      newPostgresTableDesc("autovacuum_count", "Number of times this table has been vacuumed by the autovacuum daemon", tableLabels, labels),
		analyzeCountDesc:         newPostgresTableDesc("analyze_count", "Number of times this table has been manually analyzed", tableLabels, labels),
		autoanalyzeCountDesc:     newPostgresTableDesc("autoanalyze_count", "Number of times this table has been analyzed by the autovacuum daemon", tableLabels, labels),

		heapBlksReadDesc:     newPostgresTableDesc("heap_blks_read", "Number of disk blocks read from this table", tableLabels, labels),
		heapBlksHitDesc:      newPostgresTableDesc("heap_blks_hit", "Number of buffer hits in this table", tableLabels, labels),
		idxBlksReadDesc:      newPostgresTableDesc("idx_blks_read", "Number of disk blocks read from all indexes on this table", tableLabels, labels),
		idxBlksHitDesc:       newPostgresTableDesc("idx_blks_hit", "Number of buffer hits in all indexes on this table", tableLabels, labels),
		toastBlksReadDesc:    newPostgresTableDesc("toast_blks_read", "Number of disk blocks read from this table's TOAST table (if any)", tableLabels, labels),
		toastBlksHitDesc:     newPostgresTableDesc("toast_blks_hit", "Number of buffer hits in this table's TOAST table (if any)", tableLabels, labels),
		toastIdxBlksReadDesc: newPostgresTableDesc("toast_idx_blks_read", "Number of disk blocks read from this table's TOAST table indexes (if any)", tableLabels, labels),
		toastIdxBlksHitDesc:  newPostgresTableDesc("toast_idx_blks_hit", "Number of buffer hits in this table's TOAST table indexes (if any)", tableLabels, labels),
	}
}

func newPostgresTableDesc(metric, help string, variableLabels []string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(postgresMetricPrefix+metric, help, variableLabels, constLabels)
}

func (c *PostgresCollector) Describe(ch chan<- *prometheus.Desc) {
	descs := []*prometheus.Desc{
		c.replicationLagDesc, c.postmasterStartedDesc, c.databaseSizeDesc, c.recordCountDesc,
		c.seqScanDesc, c.seqTupReadDesc, c.idxScanDesc, c.idxTupFetchDesc, c.nTupInsDesc, c.nTupUpdDesc,
		c.nTupDelDesc, c.nTupHotUpdDesc, c.nLiveTupDesc, c.nDeadTupDesc, c.nModSinceLastAnalyzeDesc,
		c.lastVacuumDesc, c.lastAutovacuumDesc, c.lastAnalyzeDesc, c.lastAutoanalyzeDesc, c.vacuumCountDesc,
		c.autovacuumCountDesc, c.analyzeCountDesc, c.autoanalyzeCountDesc, c.heapBlksReadDesc, c.heapBlksHitDesc,
		c.idxBlksReadDesc, c.idxBlksHitDesc, c.toastBlksReadDesc, c.toastBlksHitDesc, c.toastIdxBlksReadDesc,
		c.toastIdxBlksHitDesc,
	}
	for _, desc := range descs {
		ch <- desc
	}
}

func (c *PostgresCollector) Collect(ch chan<- prometheus.Metric) {
	if db.DB == nil {
		return
	}
	c.collectReplicationLag(ch)
	c.collectPostmasterStart(ch)
	c.collectDatabaseSize(ch)
	c.collectTableStats(ch)
	c.collectTableIOStats(ch)
	c.collectRecordCount(ch)
}

func (c *PostgresCollector) collectReplicationLag(ch chan<- prometheus.Metric) {
	var lag float64
	if err := db.DB.Raw("SELECT CASE WHEN NOT pg_is_in_recovery() THEN 0 ELSE GREATEST(0, EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp()))) END AS lag").Scan(&lag).Error; err != nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(c.replicationLagDesc, prometheus.GaugeValue, lag)
}

func (c *PostgresCollector) collectPostmasterStart(ch chan<- prometheus.Metric) {
	var started time.Time
	if err := db.DB.Raw("SELECT pg_postmaster_start_time() AS start_time_seconds").Scan(&started).Error; err != nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(c.postmasterStartedDesc, prometheus.GaugeValue, float64(started.Unix()))
}

func (c *PostgresCollector) collectDatabaseSize(ch chan<- prometheus.Metric) {
	type row struct {
		DatName   string `gorm:"column:datname"`
		SizeBytes int64  `gorm:"column:size_bytes"`
	}
	var rows []row
	if err := db.DB.Raw("SELECT pg_database.datname, pg_database_size(pg_database.datname) AS size_bytes FROM pg_database").Scan(&rows).Error; err != nil {
		return
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(c.databaseSizeDesc, prometheus.GaugeValue, float64(row.SizeBytes), row.DatName)
	}
}

func (c *PostgresCollector) collectTableStats(ch chan<- prometheus.Metric) {
	type row struct {
		DatName              string       `gorm:"column:datname"`
		SchemaName           string       `gorm:"column:schemaname"`
		Relname              string       `gorm:"column:relname"`
		SeqScan              int64        `gorm:"column:seq_scan"`
		SeqTupRead           int64        `gorm:"column:seq_tup_read"`
		IdxScan              int64        `gorm:"column:idx_scan"`
		IdxTupFetch          int64        `gorm:"column:idx_tup_fetch"`
		NTupIns              int64        `gorm:"column:n_tup_ins"`
		NTupUpd              int64        `gorm:"column:n_tup_upd"`
		NTupDel              int64        `gorm:"column:n_tup_del"`
		NTupHotUpd           int64        `gorm:"column:n_tup_hot_upd"`
		NLiveTup             int64        `gorm:"column:n_live_tup"`
		NDeadTup             int64        `gorm:"column:n_dead_tup"`
		NModSinceLastAnalyze int64        `gorm:"column:n_mod_since_last_analyze"`
		LastVacuum           sql.NullTime `gorm:"column:last_vacuum"`
		LastAutovacuum       sql.NullTime `gorm:"column:last_autovacuum"`
		LastAnalyze          sql.NullTime `gorm:"column:last_analyze"`
		LastAutoanalyze      sql.NullTime `gorm:"column:last_autoanalyze"`
		VacuumCount          int64        `gorm:"column:vacuum_count"`
		AutovacuumCount      int64        `gorm:"column:autovacuum_count"`
		AnalyzeCount         int64        `gorm:"column:analyze_count"`
		AutoanalyzeCount     int64        `gorm:"column:autoanalyze_count"`
	}
	var rows []row
	if err := db.DB.Raw(`
		SELECT
			current_database() AS datname,
			schemaname,
			relname,
			seq_scan,
			seq_tup_read,
			idx_scan,
			idx_tup_fetch,
			n_tup_ins,
			n_tup_upd,
			n_tup_del,
			n_tup_hot_upd,
			n_live_tup,
			n_dead_tup,
			n_mod_since_analyze AS n_mod_since_last_analyze,
			last_vacuum,
			last_autovacuum,
			last_analyze,
			last_autoanalyze,
			vacuum_count,
			autovacuum_count,
			analyze_count,
			autoanalyze_count
		FROM pg_stat_user_tables`).Scan(&rows).Error; err != nil {
		return
	}
	for _, row := range rows {
		labels := []string{row.DatName, row.SchemaName, row.Relname}
		ch <- prometheus.MustNewConstMetric(c.seqScanDesc, prometheus.CounterValue, float64(row.SeqScan), labels...)
		ch <- prometheus.MustNewConstMetric(c.seqTupReadDesc, prometheus.CounterValue, float64(row.SeqTupRead), labels...)
		ch <- prometheus.MustNewConstMetric(c.idxScanDesc, prometheus.CounterValue, float64(row.IdxScan), labels...)
		ch <- prometheus.MustNewConstMetric(c.idxTupFetchDesc, prometheus.CounterValue, float64(row.IdxTupFetch), labels...)
		ch <- prometheus.MustNewConstMetric(c.nTupInsDesc, prometheus.CounterValue, float64(row.NTupIns), labels...)
		ch <- prometheus.MustNewConstMetric(c.nTupUpdDesc, prometheus.CounterValue, float64(row.NTupUpd), labels...)
		ch <- prometheus.MustNewConstMetric(c.nTupDelDesc, prometheus.CounterValue, float64(row.NTupDel), labels...)
		ch <- prometheus.MustNewConstMetric(c.nTupHotUpdDesc, prometheus.CounterValue, float64(row.NTupHotUpd), labels...)
		ch <- prometheus.MustNewConstMetric(c.nLiveTupDesc, prometheus.GaugeValue, float64(row.NLiveTup), labels...)
		ch <- prometheus.MustNewConstMetric(c.nDeadTupDesc, prometheus.GaugeValue, float64(row.NDeadTup), labels...)
		ch <- prometheus.MustNewConstMetric(c.nModSinceLastAnalyzeDesc, prometheus.GaugeValue, float64(row.NModSinceLastAnalyze), labels...)
		ch <- prometheus.MustNewConstMetric(c.lastVacuumDesc, prometheus.GaugeValue, unixTime(row.LastVacuum), labels...)
		ch <- prometheus.MustNewConstMetric(c.lastAutovacuumDesc, prometheus.GaugeValue, unixTime(row.LastAutovacuum), labels...)
		ch <- prometheus.MustNewConstMetric(c.lastAnalyzeDesc, prometheus.GaugeValue, unixTime(row.LastAnalyze), labels...)
		ch <- prometheus.MustNewConstMetric(c.lastAutoanalyzeDesc, prometheus.GaugeValue, unixTime(row.LastAutoanalyze), labels...)
		ch <- prometheus.MustNewConstMetric(c.vacuumCountDesc, prometheus.CounterValue, float64(row.VacuumCount), labels...)
		ch <- prometheus.MustNewConstMetric(c.autovacuumCountDesc, prometheus.CounterValue, float64(row.AutovacuumCount), labels...)
		ch <- prometheus.MustNewConstMetric(c.analyzeCountDesc, prometheus.CounterValue, float64(row.AnalyzeCount), labels...)
		ch <- prometheus.MustNewConstMetric(c.autoanalyzeCountDesc, prometheus.CounterValue, float64(row.AutoanalyzeCount), labels...)
	}
}

func (c *PostgresCollector) collectTableIOStats(ch chan<- prometheus.Metric) {
	type row struct {
		DatName       string `gorm:"column:datname"`
		SchemaName    string `gorm:"column:schemaname"`
		Relname       string `gorm:"column:relname"`
		HeapBlksRead  int64  `gorm:"column:heap_blks_read"`
		HeapBlksHit   int64  `gorm:"column:heap_blks_hit"`
		IdxBlksRead   int64  `gorm:"column:idx_blks_read"`
		IdxBlksHit    int64  `gorm:"column:idx_blks_hit"`
		ToastBlksRead int64  `gorm:"column:toast_blks_read"`
		ToastBlksHit  int64  `gorm:"column:toast_blks_hit"`
		TidxBlksRead  int64  `gorm:"column:tidx_blks_read"`
		TidxBlksHit   int64  `gorm:"column:tidx_blks_hit"`
	}
	var rows []row
	if err := db.DB.Raw(`
		SELECT
			current_database() AS datname,
			schemaname,
			relname,
			heap_blks_read,
			heap_blks_hit,
			idx_blks_read,
			idx_blks_hit,
			toast_blks_read,
			toast_blks_hit,
			tidx_blks_read,
			tidx_blks_hit
		FROM pg_statio_user_tables`).Scan(&rows).Error; err != nil {
		return
	}
	for _, row := range rows {
		labels := []string{row.DatName, row.SchemaName, row.Relname}
		ch <- prometheus.MustNewConstMetric(c.heapBlksReadDesc, prometheus.CounterValue, float64(row.HeapBlksRead), labels...)
		ch <- prometheus.MustNewConstMetric(c.heapBlksHitDesc, prometheus.CounterValue, float64(row.HeapBlksHit), labels...)
		ch <- prometheus.MustNewConstMetric(c.idxBlksReadDesc, prometheus.CounterValue, float64(row.IdxBlksRead), labels...)
		ch <- prometheus.MustNewConstMetric(c.idxBlksHitDesc, prometheus.CounterValue, float64(row.IdxBlksHit), labels...)
		ch <- prometheus.MustNewConstMetric(c.toastBlksReadDesc, prometheus.CounterValue, float64(row.ToastBlksRead), labels...)
		ch <- prometheus.MustNewConstMetric(c.toastBlksHitDesc, prometheus.CounterValue, float64(row.ToastBlksHit), labels...)
		ch <- prometheus.MustNewConstMetric(c.toastIdxBlksReadDesc, prometheus.CounterValue, float64(row.TidxBlksRead), labels...)
		ch <- prometheus.MustNewConstMetric(c.toastIdxBlksHitDesc, prometheus.CounterValue, float64(row.TidxBlksHit), labels...)
	}
}

func (c *PostgresCollector) collectRecordCount(ch chan<- prometheus.Metric) {
	type row struct {
		TableSchema string `gorm:"column:table_schema"`
		TableName   string `gorm:"column:table_name"`
		RowsCount   int64  `gorm:"column:rows_count"`
	}
	var rows []row
	if err := db.DB.Raw(`
		WITH tbl AS (
			SELECT table_schema, table_name
			FROM information_schema.tables
			WHERE table_name NOT LIKE 'pg_%' AND table_schema IN ('public')
		)
		SELECT
			table_schema,
			table_name,
			(xpath('/row/c/text()', query_to_xml(format('select count(*) as c from %I.%I', table_schema, table_name), false, true, '')))[1]::text::int AS rows_count
		FROM tbl
		ORDER BY 3 DESC`).Scan(&rows).Error; err != nil {
		return
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(c.recordCountDesc, prometheus.GaugeValue, float64(row.RowsCount), row.TableSchema, row.TableName)
	}
}

func unixTime(value sql.NullTime) float64 {
	if !value.Valid {
		return 0
	}
	return float64(value.Time.Unix())
}
