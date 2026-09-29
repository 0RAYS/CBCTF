package utils

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestGroupFailureCancelsWorkersWithoutCancellingParent(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	group := NewGroup(parent)
	started := make(chan struct{})
	want := errors.New("provisioning failed")
	group.Go(func() error { close(started); <-group.Context().Done(); return group.Context().Err() })
	group.Go(func() error { <-started; return want })
	if err := group.Wait(); !errors.Is(err, want) {
		t.Fatalf("error: %v", err)
	}
	if parent.Err() != nil {
		t.Fatal("group canceled the caller's context")
	}
}
