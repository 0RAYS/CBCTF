package cron

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestStartupSchedule(t *testing.T) {
	now := time.Date(2026, time.September, 28, 12, 0, 0, 0, time.UTC)
	for _, immediate := range []bool{false, true} {
		t.Run(map[bool]string{false: "periodic", true: "startup"}[immediate], func(t *testing.T) {
			schedule := &startupSchedule{Schedule: cron.Every(time.Minute), pending: immediate}
			want := now.Add(time.Minute)
			if immediate {
				want = now
			}
			if got := schedule.Next(now); !got.Equal(want) {
				t.Fatalf("first execution = %s, want %s", got, want)
			}
			for range 3 {
				next := schedule.Next(want)
				if !next.Equal(want.Add(time.Minute)) {
					t.Fatalf("periodic execution = %s, want %s", next, want.Add(time.Minute))
				}
				want = next
			}
		})
	}
}

type skipLogger struct {
	skipped chan struct{}
}

func (l skipLogger) Info(msg string, _ ...any) {
	if msg == "skip" {
		select {
		case l.skipped <- struct{}{}:
		default:
		}
	}
}

func (skipLogger) Error(error, string, ...any) {}

func TestStartupExecutionUsesSchedulerLifecycle(t *testing.T) {
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	logger := skipLogger{skipped: make(chan struct{}, 1)}
	var runs atomic.Int32
	var releaseOnce sync.Once
	scheduler := cron.New(cron.WithChain(cron.Recover(logger), cron.SkipIfStillRunning(logger)))
	scheduler.Schedule(&startupSchedule{Schedule: cron.Every(time.Second), pending: true}, cron.FuncJob(func() {
		runs.Add(1)
		started <- struct{}{}
		<-release
	}))
	t.Cleanup(func() {
		releaseOnce.Do(func() { close(release) })
		<-scheduler.Stop().Done()
	})
	scheduler.Start()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("startup execution did not start")
	}
	select {
	case <-logger.skipped:
	case <-time.After(3 * time.Second):
		t.Fatal("periodic execution did not skip the running startup job")
	}
	ctx := scheduler.Stop()
	select {
	case <-ctx.Done():
		t.Fatal("scheduler stopped before startup execution completed")
	default:
	}
	releaseOnce.Do(func() { close(release) })
	select {
	case <-ctx.Done():
	case <-time.After(3 * time.Second):
		t.Fatal("scheduler did not finish after startup execution completed")
	}
	if got := runs.Load(); got != 1 {
		t.Fatalf("executions = %d, want 1", got)
	}
}
