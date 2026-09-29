package model

import (
	"context"
	"testing"

	"CBCTF/internal/i18n"
)

func TestBatchResultReportsFailuresSkipsAndCancellation(t *testing.T) {
	b := NewBatch(4)
	b.Success("1", "queued")
	b.Skip("2", "already_active")
	b.Fail("3", "lookup", RetVal{Msg: i18n.Model.NotFound})
	ret := b.Result(context.Background())
	if ret.OK || b.Status != "partial" || b.NotAttempted != 1 || b.Failed != 1 {
		t.Fatalf("partial batch reported success: %+v", b)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if b.Result(ctx).OK || b.Status != "cancelled" {
		t.Fatal("cancelled batch reported success")
	}
	if BatchDependencyFailed(RetVal{Msg: i18n.Model.NotFound}) || !BatchDependencyFailed(RetVal{Msg: i18n.Model.GetError}) {
		t.Fatal("domain/dependency classification broken")
	}
	allFailed := NewBatch(1)
	allFailed.Fail("1", "write", RetVal{Msg: i18n.Model.CreateError})
	if allFailed.Result(context.Background()).OK || allFailed.Status != "failed" {
		t.Fatal("all-failed batch reported success")
	}
}
