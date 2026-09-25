package model

import (
	"CBCTF/internal/config"
	"testing"
	"time"
)

func TestAttachmentGenerationsAreIsolated(t *testing.T) {
	previous := config.Env
	config.Env = &config.Config{Path: t.TempDir()}
	t.Cleanup(func() { config.Env = previous })
	challenge := Challenge{BaseModel: BaseModel{ID: 12}, GeneratorImage: "example/generator:v1"}
	challenge.UpdatedAt = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	base := challenge.AttachmentCachePathForRevision(7, []string{"flag{a}"}, "source-a")
	if base != challenge.AttachmentCachePathForRevision(7, []string{"flag{a}"}, "source-a") {
		t.Fatal("unstable cache identity")
	}
	challenge.UpdatedAt = challenge.UpdatedAt.In(time.FixedZone("UTC+8", 8*3600))
	if base != challenge.AttachmentCachePathForRevision(7, []string{"flag{a}"}, "source-a") {
		t.Fatal("replica timezone changed cache identity")
	}
	for _, path := range []string{
		challenge.AttachmentCachePathForRevision(8, []string{"flag{a}"}, "source-a"),
		challenge.AttachmentCachePathForRevision(7, []string{"flag{b}"}, "source-a"),
		challenge.AttachmentCachePathForRevision(7, []string{"flag{a}"}, "source-b"),
	} {
		if path == base {
			t.Fatal("team, flags or source revision reused a cache generation")
		}
	}
}
