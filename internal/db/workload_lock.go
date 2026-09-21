package db

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// WithWorkloadLock serializes start/stop workers across processes without holding
// a DB transaction open during Kubernetes calls. PostgreSQL releases session
// locks if a worker dies. Requires direct/session-pooled PostgreSQL connections,
// not transaction-mode PgBouncer (the chart uses direct connections).
func WithWorkloadLock(ctx context.Context, root *gorm.DB, kind string, id uint, fn func() error) error {
	if root == nil || id == 0 || (kind != "victim" && kind != "generator") {
		return fmt.Errorf("invalid workload lock")
	}
	pool, err := root.DB()
	if err != nil {
		return err
	}
	conn, err := pool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	key := fmt.Sprintf("cbctf:%s:%d", kind, id)
	// On cancellation the server may have acquired the lock before the client
	// observed the result; discard the session rather than returning it locked.
	if _, err = conn.ExecContext(ctx, "SELECT pg_advisory_lock(hashtextextended($1, 0))", key); err != nil {
		_ = conn.Raw(func(any) error { return driver.ErrBadConn })
		return fmt.Errorf("lock %s: %w", key, err)
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var unlocked bool
		err := conn.QueryRowContext(unlockCtx, "SELECT pg_advisory_unlock(hashtextextended($1, 0))", key).Scan(&unlocked)
		if err != nil || !unlocked {
			_ = conn.Raw(func(any) error { return driver.ErrBadConn })
		}
	}()
	if err := ctx.Err(); err != nil {
		return err
	}
	return fn()
}
