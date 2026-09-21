package redis

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func TestGeneratorKeyContract(t *testing.T) {
	for _, id := range []uint{1, 42, 1000000} {
		if got := fmt.Sprintf(GeneratorKeyTmpl, id); got != fmt.Sprintf("%s%d", generatorKeyPrefix, id) {
			t.Fatal(got)
		}
		if got := fmt.Sprintf(GeneratorAttachmentLockKeyTmpl, id); got != fmt.Sprintf("%s%d", generatorLockKeyPrefix, id) {
			t.Fatal(got)
		}
	}
	for _, fragment := range []string{"ARGV[3] .. id", "ARGV[4] .. id"} {
		if !strings.Contains(lockAvailableGeneratorScript, fragment) {
			t.Fatalf("Lua must use shared key prefixes: %s", fragment)
		}
	}
}

// Run against an isolated Redis: CBCTF_TEST_REDIS_ADDR=127.0.0.1:6379 go test ./internal/redis.
// Uses unique keys and never flushes the database.
func TestGeneratorLeaseIntegration(t *testing.T) {
	addr := os.Getenv("CBCTF_TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("CBCTF_TEST_REDIS_ADDR is not set")
	}
	log.Init()
	previous := RDB
	RDB = goredis.NewClient(&goredis.Options{Addr: addr})
	t.Cleanup(func() { _ = RDB.Close(); RDB = previous })
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id := uint(time.Now().UnixNano())
	generator := model.Generator{BaseModel: model.BaseModel{ID: id}, ChallengeID: id, Status: model.RunningGeneratorStatus}
	t.Cleanup(func() {
		RDB.Del(context.Background(), fmt.Sprintf(GeneratorKeyTmpl, id), fmt.Sprintf(GeneratorAttachmentLockKeyTmpl, id), generatorSetKey(0, false, id))
	})
	if err := RegisterGenerator(ctx, generator); err != nil {
		t.Fatal(err)
	}
	got, token, err := LockAvailableGenerator(ctx, 0, id)
	if err != nil || got.ID != id {
		t.Fatalf("first lease: %v %+v", err, got)
	}
	if _, _, err = LockAvailableGenerator(ctx, 0, id); !errors.Is(err, ErrNoAvailableGenerator) {
		t.Fatalf("duplicate lease: %v", err)
	}
	if valid, err := RefreshGeneratorAttachmentLock(ctx, id, token); err != nil || !valid {
		t.Fatalf("refresh: %t %v", valid, err)
	}
	if err = UnlockGeneratorAttachment(ctx, id, "wrong-owner"); err != nil {
		t.Fatal(err)
	}
	if err = UnregisterGenerator(ctx, generator); err != nil {
		t.Fatal(err)
	}
	if err = RegisterGenerator(ctx, generator); err != nil {
		t.Fatal(err)
	}
	if _, _, err = LockAvailableGenerator(ctx, 0, id); !errors.Is(err, ErrNoAvailableGenerator) {
		t.Fatalf("unregister released active lease: %v", err)
	}
	if err = UnlockGeneratorAttachment(ctx, id, token); err != nil {
		t.Fatal(err)
	}
	if _, _, err = LockAvailableGenerator(ctx, 0, id); err != nil {
		t.Fatalf("reacquire: %v", err)
	}
}
