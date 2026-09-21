package db

import (
	"context"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestWorkloadLockRejectsInvalidInput(t *testing.T) {
	if err := WithWorkloadLock(context.Background(), nil, "victim", 1, func() error { t.Fatal("invalid lock ran callback"); return nil }); err == nil {
		t.Fatal("expected validation error")
	}
}

// No schema or records are created; requires a direct PostgreSQL test connection.
func TestWorkloadLockIntegration(t *testing.T) {
	dsn := os.Getenv("CBCTF_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("CBCTF_TEST_POSTGRES_DSN is not set")
	}
	root, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	pool, err := root.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := uint(time.Now().UnixNano())
	err = WithWorkloadLock(ctx, root, "victim", id, func() error {
		blocked, stop := context.WithTimeout(ctx, 150*time.Millisecond)
		defer stop()
		err := WithWorkloadLock(blocked, root, "victim", id, func() error { t.Error("same workload lock overlapped"); return nil })
		if err == nil {
			t.Error("expected canceled lock wait")
		}
		return WithWorkloadLock(ctx, root, "generator", id, func() error { return nil })
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := WithWorkloadLock(ctx, root, "victim", id, func() error { return nil }); err != nil {
		t.Fatal("lock was not released:", err)
	}
}
