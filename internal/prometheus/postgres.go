package prometheus

import (
	"context"
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
)

const postgresMetricPrefix = "postgres_"

type PostgresCollector struct {
	collectionSuccessDesc *prometheus.Desc
	replicationLagDesc    *prometheus.Desc
	postmasterStartedDesc *prometheus.Desc
	serverInfoDesc        *prometheus.Desc
	databaseSizeDesc      *prometheus.Desc

	xactCommitDesc       *prometheus.Desc
	xactRollbackDesc     *prometheus.Desc
	deadlocksDesc        *prometheus.Desc
	blksReadDesc         *prometheus.Desc
	blksHitDesc          *prometheus.Desc
	blkReadTimeDesc      *prometheus.Desc
	tempFilesDesc        *prometheus.Desc
	tempBytesDesc        *prometheus.Desc
	conflictsDesc        *prometheus.Desc
	activityBackendsDesc *prometheus.Desc

	checkpointsTimedDesc    *prometheus.Desc
	checkpointsReqDesc      *prometheus.Desc
	checkpointWriteTimeDesc *prometheus.Desc
	checkpointSyncTimeDesc  *prometheus.Desc
	buffersCheckpointDesc   *prometheus.Desc
	buffersBackendDesc      *prometheus.Desc
	buffersBackendFsyncDesc *prometheus.Desc
	buffersCleanDesc        *prometheus.Desc

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

	tableCache cacheMemo[postgresTableStats]
	viewsCache cacheMemo[postgresViews]
}

func NewPostgresCollector() *PostgresCollector {
	labels := prometheus.Labels{"db_name": config.Env.Gorm.Postgres.DB}
	tableLabels := []string{"datname", "schemaname", "relname"}
	databaseLabels := []string{"datname"}
	activityLabels := []string{"datname", "state"}

	return &PostgresCollector{
		collectionSuccessDesc: prometheus.NewDesc(postgresMetricPrefix+"collection_success", "Whether a PostgreSQL metric section was collected", []string{"section"}, labels),
		replicationLagDesc: prometheus.NewDesc(
			postgresMetricPrefix+"replication_lag_seconds",
			"Replication lag behind master in seconds",
			nil, labels,
		),
		postmasterStartedDesc: prometheus.NewDesc(
			postgresMetricPrefix+"start_time_seconds",
			"Time unix timestamp at which postmaster started",
			nil, labels,
		),
		serverInfoDesc: prometheus.NewDesc(
			postgresMetricPrefix+"server_info",
			"PostgreSQL server version information",
			[]string{"version"}, labels,
		),
		databaseSizeDesc: prometheus.NewDesc(
			postgresMetricPrefix+"size_bytes",
			"Size of database in bytes",
			databaseLabels, labels,
		),

		xactCommitDesc:       newPostgresDesc("stat_database_xact_commit_total", "Number of transactions committed", databaseLabels, labels),
		xactRollbackDesc:     newPostgresDesc("stat_database_xact_rollback_total", "Number of transactions rolled back", databaseLabels, labels),
		deadlocksDesc:        newPostgresDesc("stat_database_deadlocks_total", "Number of deadlocks detected", databaseLabels, labels),
		blksReadDesc:         newPostgresDesc("stat_database_blks_read_total", "Number of disk blocks read", databaseLabels, labels),
		blksHitDesc:          newPostgresDesc("stat_database_blks_hit_total", "Number of buffer hits", databaseLabels, labels),
		blkReadTimeDesc:      newPostgresDesc("stat_database_blk_read_time_seconds_total", "Time spent reading data file blocks", databaseLabels, labels),
		tempFilesDesc:        newPostgresDesc("stat_database_temp_files_total", "Number of temporary files created", databaseLabels, labels),
		tempBytesDesc:        newPostgresDesc("stat_database_temp_bytes_total", "Total bytes written to temporary files", databaseLabels, labels),
		conflictsDesc:        newPostgresDesc("stat_database_conflicts_total", "Number of queries cancelled by lock conflicts", databaseLabels, labels),
		activityBackendsDesc: newPostgresDesc("stat_activity_backends", "Number of backends by wait state", activityLabels, labels),

		checkpointsTimedDesc:    newPostgresDesc("bgwriter_checkpoints_timed_total", "Number of scheduled checkpoints", nil, labels),
		checkpointsReqDesc:      newPostgresDesc("bgwriter_checkpoints_req_total", "Number of requested checkpoints", nil, labels),
		checkpointWriteTimeDesc: newPostgresDesc("bgwriter_checkpoint_write_time_seconds_total", "Total time spent writing checkpoint files", nil, labels),
		checkpointSyncTimeDesc:  newPostgresDesc("bgwriter_checkpoint_sync_time_seconds_total", "Total time spent syncing checkpoint files", nil, labels),
		buffersCheckpointDesc:   newPostgresDesc("bgwriter_buffers_checkpoint_total", "Number of buffers written during checkpoints", nil, labels),
		buffersBackendDesc:      newPostgresDesc("bgwriter_buffers_backend_total", "Number of buffers written directly by a backend", nil, labels),
		buffersBackendFsyncDesc: newPostgresDesc("bgwriter_buffers_backend_fsync_total", "Number of buffers fsynced directly by a backend", nil, labels),
		buffersCleanDesc:        newPostgresDesc("bgwriter_buffers_clean_total", "Number of buffers written by the background writer", nil, labels),

		seqScanDesc:              newPostgresDesc("seq_scan_total", "Number of sequential scans initiated on this table", tableLabels, labels),
		seqTupReadDesc:           newPostgresDesc("seq_tup_read_total", "Number of live rows fetched by sequential scans", tableLabels, labels),
		idxScanDesc:              newPostgresDesc("idx_scan_total", "Number of index scans initiated on this table", tableLabels, labels),
		idxTupFetchDesc:          newPostgresDesc("idx_tup_fetch_total", "Number of live rows fetched by index scans", tableLabels, labels),
		nTupInsDesc:              newPostgresDesc("n_tup_ins_total", "Number of rows inserted", tableLabels, labels),
		nTupUpdDesc:              newPostgresDesc("n_tup_upd_total", "Number of rows updated", tableLabels, labels),
		nTupDelDesc:              newPostgresDesc("n_tup_del_total", "Number of rows deleted", tableLabels, labels),
		nTupHotUpdDesc:           newPostgresDesc("n_tup_hot_upd_total", "Number of rows HOT updated (i.e., with no separate index update required)", tableLabels, labels),
		nLiveTupDesc:             newPostgresDesc("n_live_tup", "Estimated number of live rows", tableLabels, labels),
		nDeadTupDesc:             newPostgresDesc("n_dead_tup", "Estimated number of dead rows", tableLabels, labels),
		nModSinceLastAnalyzeDesc: newPostgresDesc("n_mod_since_last_analyze", "Estimated number of rows changed since last analyze", tableLabels, labels),
		lastVacuumDesc:           newPostgresDesc("last_vacuum_timestamp_seconds", "Last time at which this table was manually vacuumed (not counting VACUUM FULL)", tableLabels, labels),
		lastAutovacuumDesc:       newPostgresDesc("last_autovacuum_timestamp_seconds", "Last time at which this table was vacuumed by the autovacuum daemon", tableLabels, labels),
		lastAnalyzeDesc:          newPostgresDesc("last_analyze_timestamp_seconds", "Last time at which this table was manually analyzed", tableLabels, labels),
		lastAutoanalyzeDesc:      newPostgresDesc("last_autoanalyze_timestamp_seconds", "Last time at which this table was analyzed by the autovacuum daemon", tableLabels, labels),
		vacuumCountDesc:          newPostgresDesc("vacuum_count_total", "Number of times this table has been manually vacuumed (not counting VACUUM FULL)", tableLabels, labels),
		autovacuumCountDesc:      newPostgresDesc("autovacuum_count_total", "Number of times this table has been vacuumed by the autovacuum daemon", tableLabels, labels),
		analyzeCountDesc:         newPostgresDesc("analyze_count_total", "Number of times this table has been manually analyzed", tableLabels, labels),
		autoanalyzeCountDesc:     newPostgresDesc("autoanalyze_count_total", "Number of times this table has been analyzed by the autovacuum daemon", tableLabels, labels),

		heapBlksReadDesc:     newPostgresDesc("heap_blks_read_total", "Number of disk blocks read from this table", tableLabels, labels),
		heapBlksHitDesc:      newPostgresDesc("heap_blks_hit_total", "Number of buffer hits in this table", tableLabels, labels),
		idxBlksReadDesc:      newPostgresDesc("idx_blks_read_total", "Number of disk blocks read from all indexes on this table", tableLabels, labels),
		idxBlksHitDesc:       newPostgresDesc("idx_blks_hit_total", "Number of buffer hits in all indexes on this table", tableLabels, labels),
		toastBlksReadDesc:    newPostgresDesc("toast_blks_read_total", "Number of disk blocks read from this table's TOAST table (if any)", tableLabels, labels),
		toastBlksHitDesc:     newPostgresDesc("toast_blks_hit_total", "Number of buffer hits in this table's TOAST table (if any)", tableLabels, labels),
		toastIdxBlksReadDesc: newPostgresDesc("toast_idx_blks_read_total", "Number of disk blocks read from this table's TOAST table indexes (if any)", tableLabels, labels),
		toastIdxBlksHitDesc:  newPostgresDesc("toast_idx_blks_hit_total", "Number of buffer hits in this table's TOAST table indexes (if any)", tableLabels, labels),
	}
}

func newPostgresDesc(metric, help string, variableLabels []string, constLabels prometheus.Labels) *prometheus.Desc {
	return prometheus.NewDesc(postgresMetricPrefix+metric, help, variableLabels, constLabels)
}

func (c *PostgresCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.collectionSuccessDesc
	descs := []*prometheus.Desc{
		c.replicationLagDesc, c.postmasterStartedDesc, c.serverInfoDesc, c.databaseSizeDesc,
		c.xactCommitDesc, c.xactRollbackDesc, c.deadlocksDesc, c.blksReadDesc, c.blksHitDesc,
		c.blkReadTimeDesc, c.tempFilesDesc, c.tempBytesDesc, c.conflictsDesc, c.activityBackendsDesc,
		c.checkpointsTimedDesc, c.checkpointsReqDesc, c.checkpointWriteTimeDesc, c.checkpointSyncTimeDesc,
		c.buffersCheckpointDesc, c.buffersBackendDesc, c.buffersBackendFsyncDesc, c.buffersCleanDesc,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, section := range []struct {
		name    string
		collect func(*gorm.DB, chan<- prometheus.Metric) bool
	}{
		{"replication_lag", c.collectReplicationLag}, {"postmaster_start", c.collectPostmasterStart},
		{"database_size", c.collectDatabaseSize}, {"database_stats", c.collectDatabaseStats},
		{"activity", c.collectActivity}, {"checkpointer", c.collectCheckpointer},
		{"table_stats", c.collectTableStats}, {"table_io", c.collectTableIOStats},
	} {
		success := float64(0)
		if db.DB != nil && section.collect(db.DB.WithContext(ctx), ch) {
			success = 1
		}
		ch <- prometheus.MustNewConstMetric(c.collectionSuccessDesc, prometheus.GaugeValue, success, section.name)
	}
}

func (c *PostgresCollector) collectReplicationLag(root *gorm.DB, ch chan<- prometheus.Metric) bool {
	var lag float64
	if err := root.Raw(
		"SELECT CASE WHEN NOT pg_is_in_recovery() THEN 0 ELSE GREATEST(0, EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp()))) END AS lag",
	).Scan(&lag).Error; err != nil {
		return false
	}
	ch <- prometheus.MustNewConstMetric(c.replicationLagDesc, prometheus.GaugeValue, lag)
	return true
}

func (c *PostgresCollector) collectPostmasterStart(root *gorm.DB, ch chan<- prometheus.Metric) bool {
	var started time.Time
	if err := root.Raw("SELECT pg_postmaster_start_time() AS start_time_seconds").Scan(&started).Error; err != nil {
		return false
	}
	ch <- prometheus.MustNewConstMetric(c.postmasterStartedDesc, prometheus.GaugeValue, float64(started.Unix()))
	var version struct {
		Version string `gorm:"column:version"`
	}
	if err := root.Raw("SELECT current_setting('server_version') AS version").Scan(&version).Error; err == nil && version.Version != "" {
		ch <- prometheus.MustNewConstMetric(c.serverInfoDesc, prometheus.GaugeValue, 1, version.Version)
	}
	return true
}

func (c *PostgresCollector) collectDatabaseSize(root *gorm.DB, ch chan<- prometheus.Metric) bool {
	type row struct {
		DatName   string `gorm:"column:datname"`
		SizeBytes int64  `gorm:"column:size_bytes"`
	}
	var rows []row
	if err := root.Raw("SELECT pg_database.datname, pg_database_size(pg_database.datname) AS size_bytes FROM pg_database WHERE datname NOT IN ('template0', 'template1')").Scan(&rows).Error; err != nil {
		return false
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(c.databaseSizeDesc, prometheus.GaugeValue, float64(row.SizeBytes), row.DatName)
	}
	return true
}

func (c *PostgresCollector) collectDatabaseStats(root *gorm.DB, ch chan<- prometheus.Metric) bool {
	type row struct {
		DatName      string  `gorm:"column:datname"`
		XactCommit   int64   `gorm:"column:xact_commit"`
		XactRollback int64   `gorm:"column:xact_rollback"`
		Deadlocks    int64   `gorm:"column:deadlocks"`
		BlksRead     int64   `gorm:"column:blks_read"`
		BlksHit      int64   `gorm:"column:blks_hit"`
		BlkReadTime  float64 `gorm:"column:blk_read_time"`
		TempFiles    int64   `gorm:"column:temp_files"`
		TempBytes    int64   `gorm:"column:temp_bytes"`
		Conflicts    int64   `gorm:"column:conflicts"`
	}
	var rows []row
	if err := root.Raw(`
		SELECT
			datname,
			xact_commit,
			xact_rollback,
			deadlocks,
			blks_read,
			blks_hit,
			blk_read_time,
			temp_files,
			temp_bytes,
			conflicts
		FROM pg_stat_database
		WHERE datname NOT IN ('template0', 'template1')`).Scan(&rows).Error; err != nil {
		return false
	}
	for _, row := range rows {
		labels := []string{row.DatName}
		ch <- prometheus.MustNewConstMetric(c.xactCommitDesc, prometheus.CounterValue, float64(row.XactCommit), labels...)
		ch <- prometheus.MustNewConstMetric(c.xactRollbackDesc, prometheus.CounterValue, float64(row.XactRollback), labels...)
		ch <- prometheus.MustNewConstMetric(c.deadlocksDesc, prometheus.CounterValue, float64(row.Deadlocks), labels...)
		ch <- prometheus.MustNewConstMetric(c.blksReadDesc, prometheus.CounterValue, float64(row.BlksRead), labels...)
		ch <- prometheus.MustNewConstMetric(c.blksHitDesc, prometheus.CounterValue, float64(row.BlksHit), labels...)
		ch <- prometheus.MustNewConstMetric(c.blkReadTimeDesc, prometheus.CounterValue, row.BlkReadTime/1000, labels...)
		ch <- prometheus.MustNewConstMetric(c.tempFilesDesc, prometheus.CounterValue, float64(row.TempFiles), labels...)
		ch <- prometheus.MustNewConstMetric(c.tempBytesDesc, prometheus.CounterValue, float64(row.TempBytes), labels...)
		ch <- prometheus.MustNewConstMetric(c.conflictsDesc, prometheus.CounterValue, float64(row.Conflicts), labels...)
	}
	return true
}

func (c *PostgresCollector) collectActivity(root *gorm.DB, ch chan<- prometheus.Metric) bool {
	type row struct {
		DatName string `gorm:"column:datname"`
		State   string `gorm:"column:state"`
		Count   int64  `gorm:"column:count"`
	}
	var rows []row
	if err := root.Raw(`
		SELECT current_database() AS datname, COALESCE(state, 'unknown') AS state, count(*) AS count
		FROM pg_stat_activity
		GROUP BY 1, 2`).Scan(&rows).Error; err != nil {
		return false
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(c.activityBackendsDesc, prometheus.GaugeValue, float64(row.Count), row.DatName, row.State)
	}
	return true
}

// postgresViews 记录统计视图的可用性。PostgreSQL 17 把 checkpoint 计数与
// 后端写入计数拆到了 pg_stat_checkpointer 和 pg_stat_io，老版本的
// pg_stat_bgwriter 已不再提供这些列。
type postgresViews struct {
	Checkpointer bool `gorm:"column:checkpointer"`
	StatIO       bool `gorm:"column:stat_io"`
}

func (c *PostgresCollector) pgViews(root *gorm.DB) (postgresViews, bool) {
	return c.viewsCache.get(func() (postgresViews, bool) {
		var views postgresViews
		if err := root.Raw(`
			SELECT to_regclass('pg_catalog.pg_stat_checkpointer') IS NOT NULL AS checkpointer,
			       to_regclass('pg_catalog.pg_stat_io') IS NOT NULL AS stat_io
		`).Scan(&views).Error; err != nil {
			return views, false
		}
		return views, true
	})
}

func (c *PostgresCollector) collectCheckpointer(root *gorm.DB, ch chan<- prometheus.Metric) bool {
	type row struct {
		CheckpointsTimed    int64   `gorm:"column:checkpoints_timed"`
		CheckpointsReq      int64   `gorm:"column:checkpoints_req"`
		CheckpointWriteTime float64 `gorm:"column:checkpoint_write_time"`
		CheckpointSyncTime  float64 `gorm:"column:checkpoint_sync_time"`
		BuffersCheckpoint   int64   `gorm:"column:buffers_checkpoint"`
		BuffersBackend      int64   `gorm:"column:buffers_backend"`
		BuffersBackendFsync int64   `gorm:"column:buffers_backend_fsync"`
	}
	views, ok := c.pgViews(root)
	if !ok {
		return false
	}
	var checkpoints []row
	if views.Checkpointer {
		if err := root.Raw(`
			SELECT num_timed AS checkpoints_timed, num_requested AS checkpoints_req,
			       write_time AS checkpoint_write_time, sync_time AS checkpoint_sync_time,
			       buffers_written AS buffers_checkpoint
			FROM pg_stat_checkpointer`).Scan(&checkpoints).Error; err != nil {
			return false
		}
	} else if err := root.Raw(`
		SELECT checkpoints_timed, checkpoints_req, checkpoint_write_time, checkpoint_sync_time,
		       buffers_checkpoint
		FROM pg_stat_bgwriter`).Scan(&checkpoints).Error; err != nil {
		return false
	}
	var backend []row
	if views.StatIO {
		if err := root.Raw(`
			SELECT COALESCE(SUM(writebacks), 0) AS buffers_backend,
			       COALESCE(SUM(fsyncs), 0) AS buffers_backend_fsync
			FROM pg_stat_io
			WHERE backend_type = 'client backend' AND object = 'relation' AND context = 'normal'`).Scan(&backend).Error; err != nil {
			return false
		}
	} else if err := root.Raw(`
		SELECT buffers_backend, buffers_backend_fsync
		FROM pg_stat_bgwriter`).Scan(&backend).Error; err != nil {
		return false
	}
	var buffersClean int64
	if err := root.Raw("SELECT buffers_clean FROM pg_stat_bgwriter").Scan(&buffersClean).Error; err != nil {
		return false
	}
	for _, item := range checkpoints {
		ch <- prometheus.MustNewConstMetric(c.checkpointsTimedDesc, prometheus.CounterValue, float64(item.CheckpointsTimed))
		ch <- prometheus.MustNewConstMetric(c.checkpointsReqDesc, prometheus.CounterValue, float64(item.CheckpointsReq))
		ch <- prometheus.MustNewConstMetric(c.checkpointWriteTimeDesc, prometheus.CounterValue, item.CheckpointWriteTime/1000)
		ch <- prometheus.MustNewConstMetric(c.checkpointSyncTimeDesc, prometheus.CounterValue, item.CheckpointSyncTime/1000)
		ch <- prometheus.MustNewConstMetric(c.buffersCheckpointDesc, prometheus.CounterValue, float64(item.BuffersCheckpoint))
	}
	for _, item := range backend {
		ch <- prometheus.MustNewConstMetric(c.buffersBackendDesc, prometheus.CounterValue, float64(item.BuffersBackend))
		ch <- prometheus.MustNewConstMetric(c.buffersBackendFsyncDesc, prometheus.CounterValue, float64(item.BuffersBackendFsync))
	}
	ch <- prometheus.MustNewConstMetric(c.buffersCleanDesc, prometheus.CounterValue, float64(buffersClean))
	return true
}

type postgresTableStats struct {
	stats []postgresTableStatRow
	io    []postgresTableIORow
}

type postgresTableStatRow struct {
	LastVacuum           sql.NullTime `gorm:"column:last_vacuum"`
	LastAutovacuum       sql.NullTime `gorm:"column:last_autovacuum"`
	LastAnalyze          sql.NullTime `gorm:"column:last_analyze"`
	LastAutoanalyze      sql.NullTime `gorm:"column:last_autoanalyze"`
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
	VacuumCount          int64        `gorm:"column:vacuum_count"`
	AutovacuumCount      int64        `gorm:"column:autovacuum_count"`
	AnalyzeCount         int64        `gorm:"column:analyze_count"`
	AutoanalyzeCount     int64        `gorm:"column:autoanalyze_count"`
}

type postgresTableIORow struct {
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

// tableStats collects the per-table statistics. They change slowly and the
// query scans every user table, so the result is cached between scrapes.
func (c *PostgresCollector) tableStats(root *gorm.DB) (postgresTableStats, bool) {
	return c.tableCache.get(func() (postgresTableStats, bool) {
		result := postgresTableStats{}
		if err := root.Raw(`
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
			FROM pg_stat_user_tables`).Scan(&result.stats).Error; err != nil {
			return result, false
		}
		if err := root.Raw(`
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
			FROM pg_statio_user_tables`).Scan(&result.io).Error; err != nil {
			return result, false
		}
		return result, true
	})
}

func (c *PostgresCollector) collectTableStats(root *gorm.DB, ch chan<- prometheus.Metric) bool {
	stats, ok := c.tableStats(root)
	if !ok {
		return false
	}
	for _, row := range stats.stats {
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
	return true
}

func (c *PostgresCollector) collectTableIOStats(root *gorm.DB, ch chan<- prometheus.Metric) bool {
	stats, ok := c.tableStats(root)
	if !ok {
		return false
	}
	for _, row := range stats.io {
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
	return true
}

func unixTime(value sql.NullTime) float64 {
	if !value.Valid {
		return 0
	}
	return float64(value.Time.Unix())
}
