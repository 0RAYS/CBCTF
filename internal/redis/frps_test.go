package redis

import (
	"context"
	"errors"
	"testing"

	"CBCTF/internal/log"

	goredis "github.com/redis/go-redis/v9"
)

type frpsSnapshotHook struct {
	failStage     bool
	publishResult int64
	published     bool
	destructive   bool
}

func (*frpsSnapshotHook) DialHook(next goredis.DialHook) goredis.DialHook { return next }
func (h *frpsSnapshotHook) ProcessPipelineHook(_ goredis.ProcessPipelineHook) goredis.ProcessPipelineHook {
	return func(_ context.Context, cmds []goredis.Cmder) error {
		for _, cmd := range cmds {
			if cmd.Name() == "del" {
				h.destructive = true
			}
		}
		if h.failStage {
			return errors.New("staging failed")
		}
		return nil
	}
}
func (h *frpsSnapshotHook) ProcessHook(_ goredis.ProcessHook) goredis.ProcessHook {
	return func(_ context.Context, cmd goredis.Cmder) error {
		switch cmd.Name() {
		case "scan":
			cmd.(*goredis.ScanCmd).SetVal([]string{"frps:host:tcp", "frps:stale:tcp"}, 0)
		case "eval":
			h.published = true
			cmd.(*goredis.Cmd).SetVal(h.publishResult)
		case "del":
			h.destructive = true
		}
		return nil
	}
}

func TestFrpsStagingFailureKeepsExistingReservations(t *testing.T) {
	old := RDB
	client := goredis.NewClient(&goredis.Options{})
	RDB = client
	t.Cleanup(func() { RDB = old; _ = client.Close() })
	hook := &frpsSnapshotHook{failStage: true}
	client.AddHook(hook)
	_, _, ret := ReconcileFrpsPorts(map[string]map[string][]int32{"host": {"tcp": {10000}}}, "0")
	if ret.OK || hook.published || hook.destructive {
		t.Fatalf("failed staging changed live reservations: %+v %+v", ret, hook)
	}
}

func TestFrpsConcurrentAllocationDefersSnapshot(t *testing.T) {
	oldLogger := log.Logger
	log.Init()
	t.Cleanup(func() { log.Logger = oldLogger })
	old := RDB
	client := goredis.NewClient(&goredis.Options{})
	RDB = client
	t.Cleanup(func() { RDB = old; _ = client.Close() })
	hook := &frpsSnapshotHook{publishResult: 0}
	client.AddHook(hook)
	removed, kept, ret := ReconcileFrpsPorts(map[string]map[string][]int32{"host": {"tcp": {10000}}}, "old-revision")
	if !ret.OK || !hook.published || hook.destructive || removed != 0 || kept != 0 {
		t.Fatalf("concurrent allocation not protected: %+v %+v", ret, hook)
	}
}
