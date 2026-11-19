package cron

import (
	"context"
	"database/sql"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

const reindexPostgresLockKey = "cbctf:reindex_postgres"

// reindexPostgresTask rebuilds indexes concurrently to reduce table lock impact.
func reindexPostgresTask() model.RetVal {
	sqlDB, err := db.CronDB.DB()
	if err != nil {
		return model.FailRetVal(i18n.Cron.PostgresPoolError, "PostgreSQL", "", err)
	}

	ctx := context.Background()
	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		return model.FailRetVal(i18n.Cron.PostgresConnError, "PostgreSQL", "", err)
	}
	defer func(conn *sql.Conn) {
		_ = conn.Close()
	}(conn)

	var locked bool
	if err = conn.QueryRowContext(ctx, `SELECT pg_try_advisory_lock(hashtext($1))`, reindexPostgresLockKey).Scan(&locked); err != nil {
		return model.FailRetVal(i18n.Cron.PostgresLockError, "PostgreSQL", reindexPostgresLockKey, err)
	}
	if !locked {
		log.Logger.Info("PostgreSQL index rebuild skipped: another rebuild is already running")
		return model.SuccessRetVal()
	}
	defer func() {
		if _, unlockErr := conn.ExecContext(ctx, `SELECT pg_advisory_unlock(hashtext($1))`, reindexPostgresLockKey); unlockErr != nil {
			log.Logger.Warningf("Failed to release PostgreSQL index rebuild lock: %s", unlockErr)
		}
	}()

	rows, err := conn.QueryContext(ctx, `
		SELECT format('%I.%I', n.nspname, c.relname)
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind = 'i'
			AND n.nspname NOT IN ('pg_catalog', 'information_schema')
			AND n.nspname NOT LIKE 'pg_toast%'
			AND n.nspname NOT LIKE 'pg_temp_%'
		ORDER BY pg_relation_size(c.oid) DESC, n.nspname, c.relname
	`)
	if err != nil {
		return model.FailRetVal(i18n.Cron.PostgresListError, "PostgreSQL", "", err)
	}

	var indexes []string
	for rows.Next() {
		var indexName string
		if err = rows.Scan(&indexName); err != nil {
			_ = rows.Close()
			return model.FailRetVal(i18n.Cron.PostgresIndexError, "PostgreSQL", "", err)
		}
		indexes = append(indexes, indexName)
	}
	if err = rows.Err(); err != nil {
		_ = rows.Close()
		return model.FailRetVal(i18n.Cron.PostgresIndexError, "PostgreSQL", "", err)
	}
	if err = rows.Close(); err != nil {
		return model.FailRetVal(i18n.Cron.PostgresCloseError, "PostgreSQL", "", err)
	}
	if len(indexes) == 0 {
		log.Logger.Info("PostgreSQL index rebuild skipped: no user indexes found")
		return model.SuccessRetVal()
	}

	var rebuilt int
	failedIndexes := make(map[string]string)
	for _, indexName := range indexes {
		if _, err = conn.ExecContext(ctx, `REINDEX INDEX CONCURRENTLY `+indexName); err != nil {
			failedIndexes[indexName] = err.Error()
			log.Logger.Warningf("Failed to rebuild PostgreSQL index %s: %s", indexName, err)
			continue
		}
		rebuilt++
	}

	log.Logger.Infof("PostgreSQL indexes rebuilt: %d/%d", rebuilt, len(indexes))
	if len(failedIndexes) > 0 {
		return model.FailRetVal(i18n.Cron.PostgresRebuildFail, "PostgreSQL", "", nil).WithAttrs(map[string]any{
			"Rebuilt": rebuilt,
			"Failed":  len(failedIndexes),
			"Indexes": failedIndexes,
		})
	}

	return model.SuccessRetVal()
}
