package model

import (
	"CBCTF/internal/config"
	"testing"
)

func TestAttachmentGenerationsAreIsolated(t *testing.T) {
	previous := config.Env
	config.Env = &config.Config{Path: t.TempDir()}
	t.Cleanup(func() { config.Env = previous })
	challenge := Challenge{BaseModel: BaseModel{ID: 12}, GeneratorImage: "example/generator:v1"}
	base := challenge.AttachmentCachePathForRevision(7, []string{"flag{a}"}, "source-a")
	if base != challenge.AttachmentCachePathForRevision(7, []string{"flag{a}"}, "source-a") {
		t.Fatal("unstable cache identity")
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
