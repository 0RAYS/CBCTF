package oa

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"CBCTF/internal/model"
)

type cancellationTransport struct{ called bool }

func (t *cancellationTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.called = true
	return nil, r.Context().Err()
}

func TestGithubEnrichmentUsesRequestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	transport := &cancellationTransport{}
	err := SetGithubEmail(ctx, model.Oauth{}, &http.Client{Transport: transport}, map[string]any{})
	if !transport.called || !errors.Is(err, context.Canceled) {
		t.Fatalf("lost cancellation: %v", err)
	}
}
