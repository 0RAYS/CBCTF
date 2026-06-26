package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"database/sql"
)

const reindexPostgresLockKey = "cbctf:reindex_postgres"

// reindexPostgresTask rebuilds indexes concurrently to reduce table lock impact.
func reindexPostgresTask() model.RetVal {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return model.RetVal{Msg: "Failed to get PostgreSQL connection pool", Attr: map[string]any{"Error": err.Error()}}
	}

	ctx := context.Background()
	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		return model.RetVal{Msg: "Failed to get PostgreSQL connection", Attr: map[string]any{"Error": err.Error()}}
	}
	defer func(conn *sql.Conn) {
		_ = conn.Close()
	}(conn)

	var locked bool
	if err = conn.QueryRowContext(ctx, `SELECT pg_try_advisory_lock(hashtext($1))`, reindexPostgresLockKey).Scan(&locked); err != nil {
		return model.RetVal{Msg: "Failed to acquire PostgreSQL index rebuild lock", Attr: map[string]any{"Error": err.Error()}}
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
		return model.RetVal{Msg: "Failed to list PostgreSQL indexes", Attr: map[string]any{"Error": err.Error()}}
	}

	var indexes []string
	for rows.Next() {
		var indexName string
		if err = rows.Scan(&indexName); err != nil {
			_ = rows.Close()
			return model.RetVal{Msg: "Failed to read PostgreSQL index", Attr: map[string]any{"Error": err.Error()}}
		}
		indexes = append(indexes, indexName)
	}
	if err = rows.Err(); err != nil {
		_ = rows.Close()
		return model.RetVal{Msg: "Failed to iterate PostgreSQL indexes", Attr: map[string]any{"Error": err.Error()}}
	}
	if err = rows.Close(); err != nil {
		return model.RetVal{Msg: "Failed to close PostgreSQL index list", Attr: map[string]any{"Error": err.Error()}}
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
		return model.RetVal{Msg: "Failed to rebuild some PostgreSQL indexes", Attr: map[string]any{"Rebuilt": rebuilt, "Failed": len(failedIndexes), "Indexes": failedIndexes}}
	}

	return model.SuccessRetVal()
}
